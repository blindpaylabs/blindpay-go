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
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// MaxIdempotencyKeyLength is the longest Idempotency-Key the API accepts.
const MaxIdempotencyKeyLength = 255

// ErrInvalidIdempotencyKey is returned, before any request is sent, when a
// caller-supplied idempotency key (WithIdempotencyKey or an Idempotency-Key
// header) is empty, has leading or trailing whitespace, is longer than 255
// characters, or when two different keys are given.
// The API silently ignores an empty key, so sending one would look
// protected without being protected.
var ErrInvalidIdempotencyKey = errors.New("blindpay: idempotency key must be 1-255 characters")

// NetworkError is returned when no usable HTTP response was received
// (connection failure, DNS, TLS, per-attempt timeout, caller cancellation, a
// 2xx whose body could not be read or was larger than the size limit) or
// when the request could not be built (bad base URL, invalid header). Use
// errors.Is on it for context.Canceled / context.DeadlineExceeded /
// ErrResponseTooLarge. Its message never includes the request's query
// string.
type NetworkError struct {
	Err error
	// IdempotencyKey is the key the request carried, if any. Resending with
	// it is safe: the API replays the stored response or reports that the
	// first attempt is still in flight.
	IdempotencyKey string
	// Retryable is false when the caller cancelled the context, the request
	// could not be built, or it cannot be resent safely (a mutation without
	// a key, multipart).
	Retryable bool
}

func (e *NetworkError) Error() string {
	return "blindpay: request failed: " + e.Err.Error()
}

func (e *NetworkError) Unwrap() error { return e.Err }

// ErrInvalidTimeout is returned, before any request is sent, when a
// per-request timeout is zero or negative.
var ErrInvalidTimeout = errors.New("blindpay: request timeout must be positive")

// DefaultMaxResponseBodyBytes is the largest response body the SDK reads
// when Config.MaxResponseBodyBytes is zero (32 MiB).
const DefaultMaxResponseBodyBytes int64 = 32 << 20

// ErrResponseTooLarge is wrapped by the error returned when a response body
// is larger than the configured limit. A 2xx response that is too large
// returns a *NetworkError (not retried: the same request would return the
// same body); an error response that is too large returns an *APIError
// without the body.
var ErrResponseTooLarge = errors.New("blindpay: response body exceeds the size limit")

// DecodeError is returned when a 2xx response body cannot be decoded into the
// method's result type. The request succeeded on the server: do not resend it
// without checking; IdempotencyKey replays the same response.
type DecodeError struct {
	StatusCode     int
	RequestID      string
	IdempotencyKey string
	RawBody        []byte
	Err            error
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("blindpay: failed to decode response (status %d): %v", e.StatusCode, e.Err)
}

func (e *DecodeError) Unwrap() error { return e.Err }

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
// failures *NetworkError. Retries follow the rules documented on send.
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

	res, err := send(cfg, ctx, req, applyOptions(opts))
	if err != nil {
		return zero, err
	}
	if res.status == http.StatusNoContent || len(bytes.TrimSpace(res.body)) == 0 {
		return zero, nil
	}

	var result T
	if err := json.Unmarshal(res.body, &result); err != nil {
		return zero, &DecodeError{StatusCode: res.status, RequestID: res.header.Get("X-Request-Id"),
			IdempotencyKey: res.key, RawBody: res.body, Err: err}
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

// send runs the attempt loop. The rules, first match wins:
//
//   - the caller cancelled the context: stop;
//   - the request could not be built (bad URL, invalid header): stop;
//   - multipart, or a POST/PUT/PATCH/DELETE without an Idempotency-Key:
//     never resent;
//   - no response (connection error, per-attempt timeout before the status
//     line): resend after a backoff;
//   - 409 idempotency_key_in_flight or 429: resend, waiting Retry-After when
//     it is at most 60s, else the backoff;
//   - POST/PUT/PATCH/DELETE with a 5xx: stop with APIError.Retryable set (the
//     API released the key, so a resend would run the operation again);
//   - GET with 408, 500, 502, 503 or 504 (unless the code is terminal):
//     resend;
//   - anything else: stop.
//
// Once a status line arrives the attempt is classified by that status, even
// if reading the body then fails.
// success is a 2xx response.
type success struct {
	status int
	header http.Header
	body   []byte
	key    string
}

func send(cfg *Config, ctx context.Context, req outgoing, o RequestOptions) (success, error) {
	mutation := isMutation(req.method)

	callerKey, err := callerIdempotencyKey(&o)
	if err != nil {
		return success{}, err
	}
	if o.TimeoutSet && o.Timeout <= 0 {
		return success{}, ErrInvalidTimeout
	}

	key := ""
	if mutation && !req.multipart { // multipart is never keyed (the API skips it)
		switch {
		case callerKey != nil:
			key = *callerKey
		case cfg.MaxRetries > 0:
			// Generated before attempt 0: the first attempt may already
			// execute, so a key added only on retry would not protect it.
			k, err := newKey(cfg)
			if err != nil {
				return success{}, fmt.Errorf("blindpay: generating idempotency key: %w", err)
			}
			key = k
		}
	}

	resendable := !req.multipart && (!mutation || key != "")

	timeout := cfg.Timeout
	if o.Timeout > 0 {
		timeout = o.Timeout
	}

	for attempt := 0; ; attempt++ {
		res := attemptOnce(cfg, ctx, req, o, key, timeout)
		canRetry := resendable && attempt < cfg.MaxRetries

		if res.buildErr != nil {
			return success{}, &NetworkError{Err: res.buildErr, IdempotencyKey: key}
		}

		if res.status == 0 { // no response
			callerDone := ctx.Err() != nil
			netErr := &NetworkError{Err: res.err, IdempotencyKey: key, Retryable: resendable && !callerDone}
			if callerDone || !canRetry {
				return success{}, netErr
			}
			if serr := sleep(cfg, ctx, backoff(cfg, attempt)); serr != nil {
				return success{}, netErr
			}
			continue
		}

		if res.status >= 200 && res.status < 300 {
			if res.err == nil {
				return success{status: res.status, header: res.header, body: res.body, key: key}, nil
			}
			if errors.Is(res.err, ErrResponseTooLarge) {
				// Deterministic: resending returns the same body.
				return success{}, &NetworkError{Err: res.err, IdempotencyKey: key}
			}
			// Success, but the body was lost. Resending is safe for a GET and
			// for a keyed mutation (the API replays the stored response).
			netErr := &NetworkError{Err: res.err, IdempotencyKey: key, Retryable: resendable && ctx.Err() == nil}
			if !canRetry || ctx.Err() != nil {
				return success{}, netErr
			}
			if serr := sleep(cfg, ctx, backoff(cfg, attempt)); serr != nil {
				return success{}, netErr
			}
			continue
		}

		apiErr := ParseAPIError(res.status, res.body)
		if res.err != nil {
			apiErr.Message += " (response body unreadable: " + res.err.Error() + ")"
		}
		apiErr.Header = res.header
		apiErr.RequestID = res.header.Get("X-Request-Id")
		apiErr.IdempotencyKey = key
		transient, autoRetry := classifyStatus(req.method, res.status, apiErr)
		// Multipart and unkeyed mutations are never resent, so they are
		// never Retryable either.
		apiErr.Retryable = transient && resendable

		if !canRetry || !autoRetry {
			return success{}, apiErr
		}
		if serr := sleep(cfg, ctx, retryDelay(cfg, res.header, attempt)); serr != nil {
			// The caller gave up while waiting: report what the API said,
			// wrapped with the context error (errors.As finds *APIError).
			return success{}, fmt.Errorf("%w: %w", serr, apiErr)
		}
	}
}

// callerIdempotencyKey returns the caller's key, taken from
// WithIdempotencyKey or from an Idempotency-Key header (any casing; the
// header is removed so it is sent once, through the same validation). A key
// must be 1-255 characters without surrounding whitespace.
func callerIdempotencyKey(o *RequestOptions) (*string, error) {
	key := o.IdempotencyKey
	for name, values := range o.Header {
		if !strings.EqualFold(name, "Idempotency-Key") {
			continue
		}
		delete(o.Header, name)
		if len(values) != 1 || (key != nil && *key != values[0]) {
			return nil, ErrInvalidIdempotencyKey
		}
		v := values[0]
		key = &v
	}
	if key == nil {
		return nil, nil
	}
	k := *key
	if k == "" || k != strings.TrimSpace(k) || len(k) > MaxIdempotencyKeyLength {
		return nil, ErrInvalidIdempotencyKey
	}
	return key, nil
}

// attemptResult is the outcome of one attempt. status is 0 when no response
// arrived; err may be set together with a status when reading the body
// failed.
type attemptResult struct {
	status   int
	header   http.Header
	body     []byte
	err      error
	buildErr error
}

func attemptOnce(cfg *Config, ctx context.Context, req outgoing, o RequestOptions, key string, timeout time.Duration) attemptResult {
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
		return attemptResult{buildErr: fmt.Errorf("failed to create request: %w", redactURLError(err))}
	}

	for name, values := range o.Header {
		for _, v := range values {
			if !validHeaderName(name) || !validHeaderValue(v) {
				return attemptResult{buildErr: fmt.Errorf("blindpay: invalid request header %q", name)}
			}
			httpReq.Header.Add(name, v)
		}
	}
	if cfg.APIKey != nil {
		httpReq.Header.Set("Authorization", "Bearer "+*cfg.APIKey)
	}
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
		return attemptResult{err: redactURLError(err)}
	}
	body, err := readBody(resp.Body, maxResponseBodyBytes(cfg))
	_ = resp.Body.Close()
	res := attemptResult{status: resp.StatusCode, header: resp.Header, body: body}
	if err != nil {
		res.err = fmt.Errorf("failed to read response body: %w", err)
	}
	return res
}

func maxResponseBodyBytes(cfg *Config) int64 {
	if cfg.MaxResponseBodyBytes > 0 {
		return cfg.MaxResponseBodyBytes
	}
	return DefaultMaxResponseBodyBytes
}

// readBody reads at most limit bytes. It reads one byte past the limit so a
// body of exactly limit bytes is accepted and a longer one is rejected
// (returning no body) with ErrResponseTooLarge.
func readBody(r io.Reader, limit int64) ([]byte, error) {
	body, err := io.ReadAll(io.LimitReader(r, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(body)) > limit {
		return nil, fmt.Errorf("%w (%d bytes)", ErrResponseTooLarge, limit)
	}
	return body, nil
}

// redactURLError drops the query string and fragment from the URL of the
// *url.Error that net/http returns, so error messages and logs never contain
// query parameters (filter values such as customer names are PII). The
// chain is kept: errors.Is/As still reach the context error and the
// *url.Error, whose Timeout method still works.
func redactURLError(err error) error {
	urlErr, ok := err.(*url.Error) //nolint:errorlint // net/http returns *url.Error unwrapped
	if !ok {
		return err
	}
	return &url.Error{Op: urlErr.Op, URL: stripQuery(urlErr.URL), Err: urlErr.Err}
}

func stripQuery(u string) string {
	if i := strings.IndexAny(u, "?#"); i >= 0 {
		return u[:i]
	}
	return u
}

// validHeaderName reports whether name is an RFC 7230 token, as net/http
// requires (it would otherwise fail the request without sending it).
func validHeaderName(name string) bool {
	if name == "" {
		return false
	}
	for i := 0; i < len(name); i++ {
		c := name[i]
		if c <= ' ' || c >= 0x7f || strings.IndexByte("\"(),/:;<=>?@[\\]{}", c) >= 0 {
			return false
		}
	}
	return true
}

// validHeaderValue rejects control characters other than tab (CR/LF would
// split the header).
func validHeaderValue(v string) bool {
	for i := 0; i < len(v); i++ {
		if c := v[i]; (c < ' ' && c != '\t') || c == 0x7f {
			return false
		}
	}
	return true
}

// classifyStatus decides an error response. transient reports a failure that
// may succeed when resent (APIError.Retryable, which send also clears for
// multipart and unkeyed mutations): 429, 409 idempotency_key_in_flight, any
// 5xx whose code is not terminal, and 408 on GET.
// autoRetry reports whether the SDK resends it itself: 429 and the in-flight
// 409 for every method; 408, 500, 502, 503 and 504 for GET only.
func classifyStatus(method string, status int, apiErr *APIError) (transient, autoRetry bool) {
	switch {
	case status == http.StatusTooManyRequests,
		status == http.StatusConflict && apiErr.Message == "idempotency_key_in_flight":
		return true, true
	case status == http.StatusRequestTimeout:
		mutation := isMutation(method)
		return !mutation, !mutation
	case status >= 500:
		if terminalCodes[apiErr.Code] {
			return false, false
		}
		switch status {
		case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
			return true, !isMutation(method)
		}
		return true, false
	}
	return false, false
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
