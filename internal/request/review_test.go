package request

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// A status line was received, so a body-read failure is classified by the
// status: a POST that got a 500 is never resent.
func TestDo_MutationWith5xxAndBrokenBodyIsNotResent(t *testing.T) {
	tr := &scriptedTransport{steps: []step{{status: 500, body: `{"success":false,"mess`, bodyErr: errConnReset}}}
	cfg := retryConfig(tr, &fakeSleeper{}, 2)

	_, err := Do[map[string]any](cfg, context.Background(), http.MethodPost, "/instances/in_000000000000/payouts/evm", map[string]any{"quote_id": "qu"})
	apiErr := asAPIError(t, err)
	require.Len(t, tr.requests, 1, "a mutation that got a 5xx is never resent")
	require.Equal(t, 500, apiErr.StatusCode)
	require.True(t, apiErr.Retryable)
	require.Equal(t, "auto-key-1", apiErr.IdempotencyKey)
	require.Contains(t, apiErr.Message, "response body unreadable")
}

func TestDo_GetWith503AndBrokenBodyIsResent(t *testing.T) {
	tr := &scriptedTransport{steps: []step{{status: 503, bodyErr: errConnReset}, {status: 200, body: okBody}}}
	cfg := retryConfig(tr, &fakeSleeper{}, 2)

	_, err := Do[map[string]any](cfg, context.Background(), http.MethodGet, "/x", nil)
	require.NoError(t, err)
	require.Len(t, tr.requests, 2)
}

func TestDo_SuccessWithBrokenBody(t *testing.T) {
	// GET: resent. Keyed POST: resent with the same key (the API replays).
	tr := &scriptedTransport{steps: []step{{status: 200, bodyErr: errConnReset}, {status: 200, body: okBody}}}
	cfg := retryConfig(tr, &fakeSleeper{}, 2)
	_, err := Do[map[string]any](cfg, context.Background(), http.MethodPost, "/x", map[string]any{})
	require.NoError(t, err)
	require.Len(t, tr.requests, 2)
	require.Equal(t, tr.requests[0].header.Get("Idempotency-Key"), tr.requests[1].header.Get("Idempotency-Key"))

	// Unkeyed POST (retries off): a NetworkError, not resent.
	tr = &scriptedTransport{steps: []step{{status: 200, bodyErr: errConnReset}}}
	cfg = retryConfig(tr, &fakeSleeper{}, 0)
	_, err = Do[map[string]any](cfg, context.Background(), http.MethodPost, "/x", map[string]any{})
	var netErr *NetworkError
	require.True(t, errors.As(err, &netErr))
	require.False(t, netErr.Retryable)
	require.Len(t, tr.requests, 1)
}

func TestDo_RequestBuildErrorsAreNotRetried(t *testing.T) {
	tr := &scriptedTransport{}
	cfg := retryConfig(tr, &fakeSleeper{}, 2)
	cfg.BaseURL = "://not a url"
	_, err := Do[map[string]any](cfg, context.Background(), http.MethodGet, "/x", nil)
	var netErr *NetworkError
	require.True(t, errors.As(err, &netErr))
	require.False(t, netErr.Retryable)
	require.Empty(t, tr.requests)

	cfg = retryConfig(tr, &fakeSleeper{}, 2)
	for _, opt := range []RequestOption{WithHeader("X-Bad", "a\r\nInjected: 1"), WithHeader("Bad Name", "v")} {
		_, err = Do[map[string]any](cfg, context.Background(), http.MethodGet, "/x", nil, opt)
		require.True(t, errors.As(err, &netErr))
		require.False(t, netErr.Retryable)
		require.Empty(t, tr.requests)
	}
}

func TestDo_IdempotencyKeyValidation(t *testing.T) {
	for _, key := range []string{"", "   ", " k", "k\t", strings.Repeat("k", MaxIdempotencyKeyLength+1)} {
		tr := &scriptedTransport{}
		cfg := retryConfig(tr, &fakeSleeper{}, 2)
		_, err := Do[map[string]any](cfg, context.Background(), http.MethodPost, "/x", map[string]any{}, WithIdempotencyKey(key))
		require.ErrorIs(t, err, ErrInvalidIdempotencyKey, "%q", key)
		_, err = Do[map[string]any](cfg, context.Background(), http.MethodPost, "/x", map[string]any{}, WithHeader("Idempotency-Key", key))
		require.ErrorIs(t, err, ErrInvalidIdempotencyKey, "header %q", key)
		require.Empty(t, tr.requests)
	}
}

// A key passed as a header (any casing) is the caller's key: sent once,
// reported on errors, never replaced by the auto key.
func TestDo_IdempotencyKeyHeaderIsTheCallerKey(t *testing.T) {
	tr := &scriptedTransport{steps: []step{{status: 500, body: `{"message":"boom"}`}}}
	cfg := retryConfig(tr, &fakeSleeper{}, 2)

	_, err := Do[map[string]any](cfg, context.Background(), http.MethodPost, "/x", map[string]any{}, WithHeader("idempotency-key", "order-1"))
	apiErr := asAPIError(t, err)
	require.Equal(t, "order-1", apiErr.IdempotencyKey)
	require.Equal(t, []string{"order-1"}, tr.requests[0].header.Values("Idempotency-Key"))

	// Conflicting header and option values are rejected.
	tr = &scriptedTransport{}
	cfg = retryConfig(tr, &fakeSleeper{}, 2)
	_, err = Do[map[string]any](cfg, context.Background(), http.MethodPost, "/x", map[string]any{}, WithIdempotencyKey("a"), WithHeader("Idempotency-Key", "b"))
	require.ErrorIs(t, err, ErrInvalidIdempotencyKey)
	require.Empty(t, tr.requests)
}

func TestDo_MultipartNeverSendsAHeaderKey(t *testing.T) {
	tr := &scriptedTransport{steps: []step{{status: 200, body: okBody}}}
	cfg := retryConfig(tr, &fakeSleeper{}, 2)
	_, err := DoMultipart[map[string]any](cfg, context.Background(), "/upload", []byte("--b--"), "multipart/form-data; boundary=b", WithHeader("Idempotency-Key", "k"))
	require.NoError(t, err)
	require.Empty(t, tr.requests[0].header.Values("Idempotency-Key"))
}

// Cancelling while waiting out a 429 returns the API's error.
func TestDo_CancelDuringRateLimitWaitReturnsAPIError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	tr := &scriptedTransport{steps: []step{{status: 429, header: map[string]string{"Retry-After": "30"}, body: `{"code":"AUTH_RATE_LIMITED","message":"too_many_requests"}`}}}
	sl := &fakeSleeper{}
	cfg := retryConfig(tr, sl, 2)
	cfg.Sleep = func(ctx context.Context, d time.Duration) error {
		cancel()
		return sl.Sleep(ctx, d)
	}

	_, err := Do[map[string]any](cfg, ctx, http.MethodPost, "/x", map[string]any{})
	require.ErrorIs(t, err, context.Canceled)
	apiErr := asAPIError(t, err)
	require.Equal(t, 429, apiErr.StatusCode)
	require.True(t, apiErr.Retryable)
	require.Equal(t, "auto-key-1", apiErr.IdempotencyKey)
	require.Len(t, tr.requests, 1)
}
