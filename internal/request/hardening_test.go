package request

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadBody_Limit(t *testing.T) {
	body, err := readBody(strings.NewReader("12345"), 5)
	require.NoError(t, err, "a body of exactly the limit is accepted")
	require.Equal(t, "12345", string(body))

	body, err = readBody(strings.NewReader("123456"), 5)
	require.ErrorIs(t, err, ErrResponseTooLarge)
	require.Nil(t, body)
}

func TestDo_ResponseTooLarge(t *testing.T) {
	t.Run("2xx is a NetworkError and is not retried", func(t *testing.T) {
		tr := &scriptedTransport{steps: []step{
			{status: 200, body: `{"id":"po_000000000000"}`},
			{status: 200, body: okBody},
		}}
		cfg := retryConfig(tr, &fakeSleeper{}, 2)
		cfg.MaxResponseBodyBytes = 10

		_, err := Do[map[string]any](cfg, context.Background(), http.MethodGet, "/payouts", nil)
		require.ErrorIs(t, err, ErrResponseTooLarge)
		var netErr *NetworkError
		require.ErrorAs(t, err, &netErr)
		require.False(t, netErr.Retryable)
		require.Len(t, tr.requests, 1, "an oversized body is deterministic: one attempt only")
	})

	t.Run("exact limit decodes", func(t *testing.T) {
		tr := &scriptedTransport{steps: []step{{status: 200, body: okBody}}}
		cfg := retryConfig(tr, &fakeSleeper{}, 0)
		cfg.MaxResponseBodyBytes = int64(len(okBody))

		got, err := Do[map[string]any](cfg, context.Background(), http.MethodGet, "/payouts", nil)
		require.NoError(t, err)
		require.Equal(t, "po_000000000000", got["id"])
	})

	t.Run("error response keeps the status", func(t *testing.T) {
		tr := &scriptedTransport{steps: []step{{status: 400, body: `{"code":"QUOTES_EXPIRED","message":"quote_expired"}`}}}
		cfg := retryConfig(tr, &fakeSleeper{}, 0)
		cfg.MaxResponseBodyBytes = 10

		_, err := Do[map[string]any](cfg, context.Background(), http.MethodPost, "/payouts", map[string]any{})
		apiErr := asAPIError(t, err)
		require.Equal(t, 400, apiErr.StatusCode)
		require.Contains(t, apiErr.Message, ErrResponseTooLarge.Error())
	})

	t.Run("zero means the default", func(t *testing.T) {
		require.Equal(t, DefaultMaxResponseBodyBytes, maxResponseBodyBytes(&Config{}))
		require.Equal(t, int64(32<<20), DefaultMaxResponseBodyBytes)
	})
}

// TestNetworkError_RedactsQuery: transport errors must not carry query
// strings (customer names, ids) into messages and logs (C13).
func TestNetworkError_RedactsQuery(t *testing.T) {
	const pii = "Jane+Doe+PII"
	path := "/instances/in_000000000000/payouts?customer_id=cu_000000000000&customer_name=" + pii

	tr := &scriptedTransport{steps: []step{{err: &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connection refused")}}}}
	cfg := retryConfig(tr, &fakeSleeper{}, 0)

	_, err := Do[map[string]any](cfg, context.Background(), http.MethodGet, path, nil)
	var netErr *NetworkError
	require.ErrorAs(t, err, &netErr)
	require.NotContains(t, err.Error(), pii)
	require.NotContains(t, err.Error(), "customer_id")
	require.Contains(t, err.Error(), "/instances/in_000000000000/payouts")

	var urlErr *url.Error
	require.ErrorAs(t, err, &urlErr, "the *url.Error stays in the chain")
	require.NotContains(t, urlErr.URL, "?")
	var opErr *net.OpError
	require.ErrorAs(t, err, &opErr, "the cause stays in the chain")

	t.Run("caller cancellation still matches", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		cfg := retryConfig(&scriptedTransport{steps: []step{{status: 200, body: okBody}}}, &fakeSleeper{}, 0)
		cfg.HTTPClient = &http.Client{}
		_, err := Do[map[string]any](cfg, ctx, http.MethodGet, path, nil)
		require.ErrorIs(t, err, context.Canceled)
		require.NotContains(t, err.Error(), pii)
	})

	t.Run("request build errors", func(t *testing.T) {
		cfg := retryConfig(&scriptedTransport{}, &fakeSleeper{}, 0)
		cfg.BaseURL = "https://api.blindpay.com/v1\x7f"
		_, err := Do[map[string]any](cfg, context.Background(), http.MethodGet, path, nil)
		require.Error(t, err)
		require.NotContains(t, err.Error(), pii)
	})
}

func ptrTo[T any](v T) *T { return &v }
