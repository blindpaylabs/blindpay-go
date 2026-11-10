package request

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	mathrand "math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// MaxIdempotencyKeyLength is the longest Idempotency-Key the API accepts.
const MaxIdempotencyKeyLength = 255

// ErrInvalidIdempotencyKey is returned, before any request is sent, when a
// caller-supplied idempotency key is empty or longer than 255 characters.
// The API silently ignores an empty key, so sending one would look
// protected without being protected.
var ErrInvalidIdempotencyKey = errors.New("blindpay: idempotency key must be 1-255 characters")

// NetworkError is returned when no HTTP response was received (connection
// failure, DNS, TLS, per-attempt timeout, or caller cancellation). Use
// errors.Is on it for context.Canceled / context.DeadlineExceeded.
type NetworkError struct {
	Err error
	// IdempotencyKey is the key the request carried, if any. Resending with
	// it is safe: the API replays the stored response or reports that the
	// first attempt is still in flight.
	IdempotencyKey string
	// Retryable is false only when the caller cancelled the context or the
	// request cannot be retried safely (a mutation without a key, multipart).
	Retryable bool
}

func (e *NetworkError) Error() string {
	return "blindpay: request failed: " + e.Err.Error()
}

func (e *NetworkError) Unwrap() error { return e.Err }

// terminalCodes are 5xx codes the spec's error table classifies terminal.
var terminalCodes = map[string]bool{
	"PAYINS_CARD_PAYMENT_PENDING": true,
}

type outgoing struct {
	method      string
	path        string
	body        []byte
	hasBody     bool
	contentType string
	multipart   bool
}

// Do sends a JSON request and decodes a 2xx body into T. A 204 or empty body
// returns the zero value. Non-2xx responses return *APIError and transport
// failures *NetworkError. Retries follow the A2 decision table (see send).
func Do[T any](cfg *Config, ctx context.Context, method, path string, body any, opts ...RequestOption) (T, error) {
	var zero T

	req := outgoing{method: method, path: path}
	if !isNil(body) {
		b, err := json.Marshal(body)
		if err != nil {
			return zero, fmt.Errorf("failed to marshal request body: %w", err)
		}
		req.body, req.hasBody, req.contentType = b, true, "application/json"
	}

	return sendAndDecode[T](cfg, ctx, req, opts)
}

// DoMultipart sends a pre-encoded multipart/form-data body. It never sends
// Idempotency-Key (the API skips idempotency for multipart) and never
// retries, because Go cannot tell whether any byte reached the server.
func DoMultipart[T any](cfg *Config, ctx context.Context, path string, body []byte, contentType string, opts ...RequestOption) (T, error) {
	req := outgoing{
		method:      http.MethodPost,
		path:        path,
		body:        body,
		hasBody:     true,
		contentType: contentType,
		multipart:   true,
	}
	return sendAndDecode[T](cfg, ctx, req, opts)
}

func sendAndDecode[T any](cfg *Config, ctx context.Context, req outgoing, opts []RequestOption) (T, error) {
	var zero T

	status, respBody, err := send(cfg, ctx, req, applyOptions(opts))
	if err != nil {
		return zero, err
	}
	if status == http.StatusNoContent || len(bytes.TrimSpace(respBody)) == 0 {
		return zero, nil
	}

	var result T
	if err := json.Unmarshal(respBody, &result); err != nil {
		return zero, fmt.Errorf("failed to decode response: %w", err)
	}
	return result, nil
}

func isMutation(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	}
	return false
}

// send runs the attempt loop. Decision table (first match wins):
//
//  1. caller cancelled                       -> stop
//  2. multipart                              -> never retried
//  3. mutation without Idempotency-Key       -> never retried
//  4. no response (network error, timeout)   -> retry, backoff
//  5. 409 idempotency_key_in_flight          -> retry, Retry-After <= 60s else backoff
//  6. 429                                    -> retry, Retry-After <= 60s else backoff
//  7. mutation + 5xx                         -> stop, APIError.Retryable = true
//  8. GET + 5xx with a terminal code         -> stop
//  9. GET + 408/500/502/503/504              -> retry, Retry-After <= 60s else backoff
//  10. anything else                         -> stop
func send(cfg *Config, ctx context.Context, req outgoing, o RequestOptions) (int, []byte, error) {
	mutation := isMutation(req.method)

	key := ""
	if mutation && !req.multipart {
		switch {
		case o.IdempotencyKey != nil:
			if *o.IdempotencyKey == "" || len(*o.IdempotencyKey) > MaxIdempotencyKeyLength {
				return 0, nil, ErrInvalidIdempotencyKey
			}
			key = *o.IdempotencyKey
		case cfg.MaxRetries > 0:
			// Generated before attempt 0: the first attempt may already
			// execute, so a key added only on retry would not protect it.
			k, err := newKey(cfg)
			if err != nil {
				return 0, nil, fmt.Errorf("blindpay: generating idempotency key: %w", err)
			}
			key = k
		}
	}

	// Rows 2 and 3: requests that may never be retried.
	retryable := !req.multipart && (!mutation || key != "")

	timeout := cfg.Timeout
	if o.Timeout > 0 {
		timeout = o.Timeout
	}

	for attempt := 0; ; attempt++ {
		status, header, body, err := attemptOnce(cfg, ctx, req, o, key, timeout)
		canRetry := retryable && attempt < cfg.MaxRetries

		if err != nil {
			callerDone := ctx.Err() != nil
			netErr := &NetworkError{Err: err, IdempotencyKey: key, Retryable: retryable && !callerDone}
			if callerDone || !canRetry { // row 1
				return 0, nil, netErr
			}
			if serr := sleep(cfg, ctx, backoff(cfg, attempt)); serr != nil {
				return 0, nil, netErr
			}
			continue // row 4
		}

		if status >= 200 && status < 300 {
			return status, body, nil
		}

		apiErr := ParseAPIError(status, body)
		apiErr.Header = header
		apiErr.RequestID = header.Get("X-Request-Id")
		apiErr.IdempotencyKey = key
		apiErr.Retryable = isRetryableStatus(status, apiErr)

		if !canRetry || !shouldRetry(req.method, status, apiErr) {
			return 0, nil, apiErr
		}
		if serr := sleep(cfg, ctx, retryDelay(cfg, header, attempt)); serr != nil {
			return 0, nil, &NetworkError{Err: serr, IdempotencyKey: key}
		}
	}
}

func attemptOnce(cfg *Config, ctx context.Context, req outgoing, o RequestOptions, key string, timeout time.Duration) (int, http.Header, []byte, error) {
	attemptCtx := ctx
	if timeout > 0 {
		var cancel context.CancelFunc
		attemptCtx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	var reader io.Reader
	if req.hasBody {
		reader = bytes.NewReader(req.body)
	}
	httpReq, err := http.NewRequestWithContext(attemptCtx, req.method, cfg.BaseURL+req.path, reader)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	for name, values := range o.Header {
		for _, v := range values {
			httpReq.Header.Add(name, v)
		}
	}
	httpReq.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	httpReq.Header.Set("User-Agent", cfg.UserAgent)
	httpReq.Header.Set("Accept", "application/json")
	if req.hasBody {
		httpReq.Header.Set("Content-Type", req.contentType)
	}
	if key != "" {
		httpReq.Header.Set("Idempotency-Key", key)
	}

	resp, err := httpClient(cfg).Do(httpReq)
	if err != nil {
		return 0, nil, nil, err
	}
	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return 0, nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return resp.StatusCode, resp.Header, body, nil
}

func shouldRetry(method string, status int, apiErr *APIError) bool {
	switch {
	case status == http.StatusConflict && apiErr.Message == "idempotency_key_in_flight": // row 5
		return true
	case status == http.StatusTooManyRequests: // row 6
		return true
	case isMutation(method): // rows 7 and 10 (mutation 408)
		return false
	case status >= 500 && terminalCodes[apiErr.Code]: // row 8
		return false
	}
	switch status { // row 9
	case http.StatusRequestTimeout, http.StatusInternalServerError, http.StatusBadGateway,
		http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	}
	return false
}

func isRetryableStatus(status int, apiErr *APIError) bool {
	switch {
	case status == http.StatusConflict:
		return apiErr.Message == "idempotency_key_in_flight"
	case status == http.StatusRequestTimeout, status == http.StatusTooManyRequests:
		return true
	case status >= 500:
		return !terminalCodes[apiErr.Code]
	}
	return false
}

// retryDelay honours an integer Retry-After of at most 60 seconds; anything
// else (absent, HTTP-date, larger) falls back to backoff.
func retryDelay(cfg *Config, header http.Header, attempt int) time.Duration {
	if v := strings.TrimSpace(header.Get("Retry-After")); v != "" {
		if secs, err := strconv.Atoi(v); err == nil && secs >= 0 && secs <= 60 {
			return time.Duration(secs) * time.Second
		}
	}
	return backoff(cfg, attempt)
}

// backoff returns min(8s, 0.5s * 2^attempt) * jitter, jitter in [0.5, 1.0).
func backoff(cfg *Config, attempt int) time.Duration {
	base := 8 * time.Second
	if attempt < 4 {
		base = 500 * time.Millisecond << uint(attempt)
	}
	jitter := mathrand.Float64
	if cfg.Jitter != nil {
		jitter = cfg.Jitter
	}
	factor := 0.5 + jitter()/2
	return time.Duration(float64(base) * factor)
}

func sleep(cfg *Config, ctx context.Context, d time.Duration) error {
	if cfg.Sleep != nil {
		return cfg.Sleep(ctx, d)
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func newKey(cfg *Config) (string, error) {
	if cfg.NewKey != nil {
		return cfg.NewKey()
	}
	return NewUUIDv4()
}

// NewUUIDv4 returns a random RFC 4122 version 4 UUID.
func NewUUIDv4() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

func httpClient(cfg *Config) *http.Client {
	if cfg.HTTPClient != nil {
		return cfg.HTTPClient
	}
	return http.DefaultClient
}

func isNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Interface:
		return rv.IsNil()
	}
	return false
}
