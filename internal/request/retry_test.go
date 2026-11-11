package request

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/stretchr/testify/require"
)

// step is one scripted transport outcome: an error (no response) or a
// response with status, headers and body.
type step struct {
	err    error
	status int
	header map[string]string
	body   string
	// bodyErr, when set, fails the body read after body has been read
	// (a connection reset mid-body).
	bodyErr error
}

type failingReader struct{ err error }

func (f failingReader) Read([]byte) (int, error) { return 0, f.err }

type recorded struct {
	method string
	header http.Header
	body   []byte
}

// scriptedTransport replays steps in order and records every request.
type scriptedTransport struct {
	mu       sync.Mutex
	steps    []step
	requests []recorded
}

func (s *scriptedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
	}
	s.requests = append(s.requests, recorded{method: req.Method, header: req.Header.Clone(), body: body})

	if len(s.requests) > len(s.steps) {
		return nil, errors.New("scriptedTransport: unexpected extra request")
	}
	st := s.steps[len(s.requests)-1]
	if st.err != nil {
		return nil, st.err
	}
	h := http.Header{}
	for k, v := range st.header {
		h.Set(k, v)
	}
	var respBody io.Reader = strings.NewReader(st.body)
	if st.bodyErr != nil {
		respBody = io.MultiReader(respBody, failingReader{st.bodyErr})
	}
	return &http.Response{
		StatusCode: st.status,
		Header:     h,
		Body:       io.NopCloser(respBody),
		Request:    req,
	}, nil
}

type fakeSleeper struct {
	mu     sync.Mutex
	delays []time.Duration
}

func (f *fakeSleeper) Sleep(ctx context.Context, d time.Duration) error {
	f.mu.Lock()
	f.delays = append(f.delays, d)
	f.mu.Unlock()
	return ctx.Err()
}

func retryConfig(tr http.RoundTripper, sl *fakeSleeper, maxRetries int) *Config {
	n := 0
	return &Config{
		BaseURL:    "https://api.blindpay.com/v1",
		APIKey:     "test_key",
		HTTPClient: &http.Client{Transport: tr},
		UserAgent:  "blindpay-go/test",
		MaxRetries: maxRetries,
		Sleep:      sl.Sleep,
		Jitter:     func() float64 { return 1 }, // factor 1.0 -> deterministic backoff
		NewKey: func() (string, error) {
			n++
			return "auto-key-" + string(rune('0'+n)), nil
		},
	}
}

var errConnReset = errors.New("connection reset by peer")

const okBody = `{"id":"po_000000000000"}`

func TestDo_RetryDecisionTable(t *testing.T) {
	rateLimited := blindpaytest.LoadErrorFixture(t, "rate_limited_429.json")
	inFlight := blindpaytest.LoadErrorFixture(t, "idempotency_in_flight_409.json")
	mismatch := blindpaytest.LoadErrorFixture(t, "idempotency_payload_mismatch_422.json")
	internal := blindpaytest.LoadErrorFixture(t, "internal_error_no_errors_500.json")
	nonJSON := blindpaytest.LoadErrorFixture(t, "non_json_502.json")

	tests := []struct {
		name       string
		method     string
		maxRetries int
		opts       []RequestOption
		steps      []step
		wantCalls  int
		wantDelays []time.Duration
		wantOK     bool
		check      func(t *testing.T, err error)
	}{
		{
			name: "GET retried on 503 then success", method: http.MethodGet, maxRetries: 2,
			steps:     []step{{status: 503, body: `{"message":"unavailable"}`}, {status: 200, body: okBody}},
			wantCalls: 2, wantDelays: []time.Duration{500 * time.Millisecond}, wantOK: true,
		},
		{
			name: "GET retried on non-JSON 502 and 504, then exhausted", method: http.MethodGet, maxRetries: 2,
			steps: []step{
				{status: 502, header: nonJSON.Headers, body: string(nonJSON.ResponseBody())},
				{status: 504, body: ``},
				{status: 502, header: nonJSON.Headers, body: string(nonJSON.ResponseBody())},
			},
			wantCalls: 3, wantDelays: []time.Duration{500 * time.Millisecond, time.Second},
			check: func(t *testing.T, err error) {
				apiErr := asAPIError(t, err)
				require.Equal(t, 502, apiErr.StatusCode)
				require.True(t, apiErr.Retryable)
				require.Empty(t, apiErr.IdempotencyKey, "GET carries no key")
			},
		},
		{
			name: "GET retried on 408", method: http.MethodGet, maxRetries: 1,
			steps:     []step{{status: 408, body: `{"message":"timeout"}`}, {status: 200, body: okBody}},
			wantCalls: 2, wantOK: true, wantDelays: []time.Duration{500 * time.Millisecond},
		},
		{
			name: "GET 5xx with terminal code is not retried", method: http.MethodGet, maxRetries: 2,
			steps:     []step{{status: 502, body: `{"success":false,"code":"PAYINS_CARD_PAYMENT_PENDING","message":"card_payment_pending"}`}},
			wantCalls: 1,
			check: func(t *testing.T, err error) {
				require.False(t, asAPIError(t, err).Retryable)
			},
		},
		{
			name: "GET retried on network error", method: http.MethodGet, maxRetries: 2,
			steps:     []step{{err: errConnReset}, {status: 200, body: okBody}},
			wantCalls: 2, wantOK: true, wantDelays: []time.Duration{500 * time.Millisecond},
		},
		{
			name: "mutation NOT retried on 500", method: http.MethodPost, maxRetries: 2,
			steps:     []step{{status: 500, body: string(internal.ResponseBody())}},
			wantCalls: 1,
			check: func(t *testing.T, err error) {
				apiErr := asAPIError(t, err)
				require.True(t, apiErr.Retryable)
				require.Equal(t, "auto-key-1", apiErr.IdempotencyKey)
				require.Equal(t, "INTERNAL_ERROR", apiErr.Code)
			},
		},
		{
			name: "mutation NOT retried on 502", method: http.MethodPost, maxRetries: 2,
			steps:     []step{{status: 502, header: nonJSON.Headers, body: string(nonJSON.ResponseBody())}},
			wantCalls: 1,
			check: func(t *testing.T, err error) {
				apiErr := asAPIError(t, err)
				require.True(t, apiErr.Retryable)
				require.Equal(t, "auto-key-1", apiErr.IdempotencyKey)
			},
		},
		{
			name: "mutation NOT retried on 503", method: http.MethodDelete, maxRetries: 2,
			steps:     []step{{status: 503, body: `{"message":"unavailable"}`}},
			wantCalls: 1,
			check: func(t *testing.T, err error) {
				require.True(t, asAPIError(t, err).Retryable)
			},
		},
		{
			name: "mutation NOT retried on 408", method: http.MethodPut, maxRetries: 2,
			steps:     []step{{status: 408, body: `{"message":"timeout"}`}},
			wantCalls: 1,
		},
		{
			name: "mutation retried on network error with the same key", method: http.MethodPost, maxRetries: 2,
			steps:     []step{{err: errConnReset}, {status: 200, body: okBody}},
			wantCalls: 2, wantOK: true, wantDelays: []time.Duration{500 * time.Millisecond},
		},
		{
			name: "mutation retried on 429 honouring Retry-After", method: http.MethodPost, maxRetries: 2,
			steps:     []step{{status: 429, header: rateLimited.Headers, body: string(rateLimited.ResponseBody())}, {status: 200, body: okBody}},
			wantCalls: 2, wantOK: true, wantDelays: []time.Duration{2 * time.Second},
		},
		{
			name: "409 in flight retried with Retry-After 1", method: http.MethodPost, maxRetries: 2,
			steps:     []step{{status: 409, header: inFlight.Headers, body: string(inFlight.ResponseBody())}, {status: 200, body: okBody}},
			wantCalls: 2, wantOK: true, wantDelays: []time.Duration{time.Second},
		},
		{
			name: "409 in flight with Retry-After 120 falls back to backoff", method: http.MethodPost, maxRetries: 2,
			steps: []step{
				{status: 409, header: map[string]string{"Retry-After": "120"}, body: string(inFlight.ResponseBody())},
				{status: 409, header: map[string]string{"Retry-After": "120"}, body: string(inFlight.ResponseBody())},
				{status: 409, header: map[string]string{"Retry-After": "120"}, body: string(inFlight.ResponseBody())},
			},
			wantCalls: 3, wantDelays: []time.Duration{500 * time.Millisecond, time.Second},
			check: func(t *testing.T, err error) {
				apiErr := asAPIError(t, err)
				require.True(t, apiErr.Retryable)
				require.Equal(t, "auto-key-1", apiErr.IdempotencyKey)
			},
		},
		{
			name: "Retry-After as HTTP date falls back to backoff", method: http.MethodGet, maxRetries: 1,
			steps:     []step{{status: 429, header: map[string]string{"Retry-After": "Wed, 21 Oct 2015 07:28:00 GMT"}, body: `{"message":"rate_limited"}`}, {status: 200, body: okBody}},
			wantCalls: 2, wantOK: true, wantDelays: []time.Duration{500 * time.Millisecond},
		},
		{
			name: "other 409 not retried", method: http.MethodPost, maxRetries: 2,
			steps:     []step{{status: 409, body: `{"success":false,"code":"WALLETS_SELF_CUSTODY_ALREADY_SET","message":"self_custody_already_set"}`}},
			wantCalls: 1,
			check: func(t *testing.T, err error) {
				require.False(t, asAPIError(t, err).Retryable)
			},
		},
		{
			name: "422 not retried", method: http.MethodPost, maxRetries: 2,
			steps:     []step{{status: 422, body: string(mismatch.ResponseBody())}},
			wantCalls: 1,
		},
		{
			name: "400 not retried", method: http.MethodGet, maxRetries: 2,
			steps:     []step{{status: 400, body: `{"message":"bad_request"}`}},
			wantCalls: 1,
		},
		{
			name: "max_retries=0 disables retries and auto keys", method: http.MethodPost, maxRetries: 0,
			steps:     []step{{err: errConnReset}},
			wantCalls: 1,
			check: func(t *testing.T, err error) {
				var netErr *NetworkError
				require.True(t, errors.As(err, &netErr))
				require.Empty(t, netErr.IdempotencyKey)
				require.False(t, netErr.Retryable, "a mutation without a key is not safely retryable")
				require.ErrorIs(t, err, errConnReset)
			},
		},
		{
			name: "max_retries=0 GET not retried", method: http.MethodGet, maxRetries: 0,
			steps:     []step{{status: 503, body: `{}`}},
			wantCalls: 1,
		},
		{
			name: "caller key is used and kept on every attempt", method: http.MethodPost, maxRetries: 2,
			opts:      []RequestOption{WithIdempotencyKey("caller-key")},
			steps:     []step{{err: errConnReset}, {status: 429, body: `{}`}, {status: 200, body: okBody}},
			wantCalls: 3, wantOK: true, wantDelays: []time.Duration{500 * time.Millisecond, time.Second},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tr := &scriptedTransport{steps: tc.steps}
			sl := &fakeSleeper{}
			cfg := retryConfig(tr, sl, tc.maxRetries)

			var body any
			if tc.method == http.MethodPost || tc.method == http.MethodPut {
				body = map[string]any{"quote_id": "qu_000000000000"}
			}
			_, err := Do[map[string]any](cfg, context.Background(), tc.method, "/instances/in_000000000000/payouts", body, tc.opts...)

			require.Len(t, tr.requests, tc.wantCalls)
			require.Equal(t, tc.wantDelays, sl.delays)
			if tc.wantOK {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
			if tc.check != nil {
				tc.check(t, err)
			}

			// Idempotency-Key: present iff the request is a mutation with a
			// key, identical on every attempt; body bytes identical too.
			first := tr.requests[0]
			for _, r := range tr.requests[1:] {
				require.Equal(t, first.header.Get("Idempotency-Key"), r.header.Get("Idempotency-Key"))
				require.Equal(t, first.body, r.body)
			}
			if tc.method == http.MethodGet {
				require.Empty(t, first.header.Values("Idempotency-Key"))
			}
		})
	}
}

func asAPIError(t *testing.T, err error) *APIError {
	t.Helper()
	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr), "expected *APIError, got %T: %v", err, err)
	return apiErr
}

func TestDo_AutoKeyGeneratedBeforeFirstAttempt(t *testing.T) {
	tr := &scriptedTransport{steps: []step{{status: 200, body: okBody}}}
	cfg := retryConfig(tr, &fakeSleeper{}, 2)

	_, err := Do[map[string]any](cfg, context.Background(), http.MethodDelete, "/instances/in_000000000000/payables/pa_000000000000", nil)
	require.NoError(t, err)
	require.Len(t, tr.requests, 1)
	require.Equal(t, "auto-key-1", tr.requests[0].header.Get("Idempotency-Key"))
	require.Empty(t, tr.requests[0].header.Values("Content-Type"), "no Content-Type without a body")
}

func TestDo_InvalidCallerKeyRejectedWithoutRequest(t *testing.T) {
	for _, key := range []string{"", strings.Repeat("k", MaxIdempotencyKeyLength+1)} {
		tr := &scriptedTransport{}
		cfg := retryConfig(tr, &fakeSleeper{}, 2)

		_, err := Do[map[string]any](cfg, context.Background(), http.MethodPost, "/instances/in_000000000000/payouts/evm", map[string]any{}, WithIdempotencyKey(key))
		require.ErrorIs(t, err, ErrInvalidIdempotencyKey)
		require.Empty(t, tr.requests)
	}

	// 255 characters is accepted.
	tr := &scriptedTransport{steps: []step{{status: 200, body: okBody}}}
	cfg := retryConfig(tr, &fakeSleeper{}, 2)
	key := strings.Repeat("k", MaxIdempotencyKeyLength)
	_, err := Do[map[string]any](cfg, context.Background(), http.MethodPost, "/x", map[string]any{}, WithIdempotencyKey(key))
	require.NoError(t, err)
	require.Equal(t, key, tr.requests[0].header.Get("Idempotency-Key"))
}

// blockingTransport blocks until the request context is done, standing in
// for a server that never answers.
type blockingTransport struct {
	mu    sync.Mutex
	calls int
	keys  []string
	then  *step
}

func (b *blockingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	b.mu.Lock()
	b.calls++
	n := b.calls
	b.keys = append(b.keys, req.Header.Get("Idempotency-Key"))
	b.mu.Unlock()

	if n > 1 && b.then != nil {
		return &http.Response{StatusCode: b.then.status, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(b.then.body)), Request: req}, nil
	}
	<-req.Context().Done()
	return nil, req.Context().Err()
}

func TestDo_PerAttemptTimeoutThenSuccessReusesKey(t *testing.T) {
	tr := &blockingTransport{then: &step{status: 200, body: okBody}}
	cfg := retryConfig(tr, &fakeSleeper{}, 2)
	cfg.Timeout = 20 * time.Millisecond

	_, err := Do[map[string]any](cfg, context.Background(), http.MethodPost, "/instances/in_000000000000/payouts/evm", map[string]any{"quote_id": "qu"})
	require.NoError(t, err)
	require.Equal(t, 2, tr.calls)
	require.Equal(t, []string{"auto-key-1", "auto-key-1"}, tr.keys)
}

func TestDo_RequestTimeoutOptionOverridesClient(t *testing.T) {
	tr := &blockingTransport{}
	cfg := retryConfig(tr, &fakeSleeper{}, 0)
	cfg.Timeout = time.Hour

	start := time.Now()
	_, err := Do[map[string]any](cfg, context.Background(), http.MethodGet, "/x", nil, WithTimeout(20*time.Millisecond))
	require.Error(t, err)
	require.Less(t, time.Since(start), 10*time.Second)
	require.ErrorIs(t, err, context.DeadlineExceeded)

	var netErr *NetworkError
	require.True(t, errors.As(err, &netErr))
	require.True(t, netErr.Retryable)
}

func TestDo_CallerCancellationStopsRetries(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	tr := &scriptedTransport{steps: []step{{status: 503, body: `{}`}, {status: 503, body: `{}`}, {status: 503, body: `{}`}}}
	sl := &fakeSleeper{}
	cfg := retryConfig(tr, sl, 2)
	cfg.Sleep = func(ctx context.Context, d time.Duration) error {
		cancel() // the caller gives up while the SDK waits
		return sl.Sleep(ctx, d)
	}

	_, err := Do[map[string]any](cfg, ctx, http.MethodGet, "/x", nil)
	require.ErrorIs(t, err, context.Canceled)
	require.Len(t, tr.requests, 1, "no attempt after the caller cancelled")
}

func TestDo_DefaultSleepStopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cfg := &Config{}
	done := make(chan error, 1)
	go func() { done <- sleep(cfg, ctx, time.Hour) }()
	cancel()
	select {
	case err := <-done:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(5 * time.Second):
		t.Fatal("sleep did not stop on cancellation")
	}
}

func TestDo_HeadersAndRequestIDOnError(t *testing.T) {
	tr := &scriptedTransport{steps: []step{{status: 404, header: map[string]string{"X-Request-Id": "req_123"}, body: `{"success":false,"code":"CUSTOMERS_NOT_FOUND","message":"customer_not_found"}`}}}
	cfg := retryConfig(tr, &fakeSleeper{}, 2)

	_, err := Do[map[string]any](cfg, context.Background(), http.MethodGet, "/x", nil, WithHeader("X-Trace", "abc"))
	apiErr := asAPIError(t, err)
	require.Equal(t, "req_123", apiErr.RequestID)
	require.Equal(t, "req_123", apiErr.Header.Get("X-Request-Id"))
	require.False(t, apiErr.Retryable)

	h := tr.requests[0].header
	require.Equal(t, "abc", h.Get("X-Trace"))
	require.Equal(t, "Bearer test_key", h.Get("Authorization"))
	require.Equal(t, "blindpay-go/test", h.Get("User-Agent"))
	require.Empty(t, h.Values("Content-Type"))
	require.NotContains(t, apiErr.Error(), "test_key")
}

func TestDo_JSONBodySerialisedOnceWithContentType(t *testing.T) {
	tr := &scriptedTransport{steps: []step{{status: 200, body: okBody}}}
	cfg := retryConfig(tr, &fakeSleeper{}, 0)

	var nilBody *struct{ A string }
	_, err := Do[map[string]any](cfg, context.Background(), http.MethodPost, "/x", nilBody)
	require.NoError(t, err)
	require.Empty(t, tr.requests[0].body, "typed nil body sends no body")

	tr = &scriptedTransport{steps: []step{{status: 200, body: okBody}}}
	cfg = retryConfig(tr, &fakeSleeper{}, 0)
	_, err = Do[map[string]any](cfg, context.Background(), http.MethodPost, "/x", map[string]string{"a": "b"})
	require.NoError(t, err)
	require.Equal(t, "application/json", tr.requests[0].header.Get("Content-Type"))
	require.True(t, bytes.Equal([]byte(`{"a":"b"}`), tr.requests[0].body))
	require.Empty(t, tr.requests[0].header.Values("Idempotency-Key"), "no key when retries are off and none supplied")
}

func TestDo_MultipartNeverRetriedNorKeyed(t *testing.T) {
	tr := &scriptedTransport{steps: []step{{err: errConnReset}}}
	cfg := retryConfig(tr, &fakeSleeper{}, 2)

	_, err := DoMultipart[map[string]any](cfg, context.Background(), "/upload", []byte("--b--"), "multipart/form-data; boundary=b", WithIdempotencyKey("ignored"))
	require.Error(t, err)
	require.Len(t, tr.requests, 1)
	require.Empty(t, tr.requests[0].header.Values("Idempotency-Key"))
	require.Equal(t, "multipart/form-data; boundary=b", tr.requests[0].header.Get("Content-Type"))
}

func TestNewUUIDv4(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		k, err := NewUUIDv4()
		require.NoError(t, err)
		require.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`, k)
		require.False(t, seen[k])
		seen[k] = true
	}
}

func TestBackoff(t *testing.T) {
	cfg := &Config{Jitter: func() float64 { return 1 }}
	require.Equal(t, 500*time.Millisecond, backoff(cfg, 0))
	require.Equal(t, time.Second, backoff(cfg, 1))
	require.Equal(t, 4*time.Second, backoff(cfg, 3))
	require.Equal(t, 8*time.Second, backoff(cfg, 4))
	require.Equal(t, 8*time.Second, backoff(cfg, 30))

	cfg.Jitter = func() float64 { return 0 }
	require.Equal(t, 250*time.Millisecond, backoff(cfg, 0))
}
