package blindpay

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type countingTransport struct {
	mu     sync.Mutex
	calls  int
	status []int
	keys   []string
}

func (c *countingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	status := c.status[c.calls]
	c.calls++
	c.keys = append(c.keys, req.Header.Get("Idempotency-Key"))
	h := http.Header{}
	// Retry-After: 0 keeps the retry test from sleeping.
	h.Set("Retry-After", "0")
	return &http.Response{StatusCode: status, Header: h, Body: io.NopCloser(strings.NewReader(`{"id":"po_000000000000"}`)), Request: req}, nil
}

func TestNew_Defaults(t *testing.T) {
	c, err := New("test_key", "in_000000000000")
	require.NoError(t, err)
	require.Equal(t, DefaultTimeout, c.timeout)
	require.Equal(t, DefaultMaxRetries, c.maxRetries)
	require.Zero(t, c.httpClient.Timeout, "the default client must not cap the per-attempt timeout")
	require.Equal(t, "https://api.blindpay.com/v1", c.baseURL)
}

func TestNew_Options(t *testing.T) {
	hc := &http.Client{}
	c, err := New("test_key", "in_000000000000",
		WithHTTPClient(hc), WithBaseURL("http://localhost:1"), WithTimeout(5e9), WithMaxRetries(0))
	require.NoError(t, err)
	require.Same(t, hc, c.httpClient)
	require.Zero(t, hc.Timeout, "a caller-supplied client is never modified")
	require.EqualValues(t, 5e9, c.timeout)
	require.Zero(t, c.maxRetries)

	c, err = New("test_key", "in_000000000000", WithMaxRetries(-1), WithTimeout(-1))
	require.NoError(t, err)
	require.Zero(t, c.maxRetries)
	require.Zero(t, c.timeout)
}

func TestNew_RetriesReachResourceClients(t *testing.T) {
	tr := &countingTransport{status: []int{503, 200}}
	c, err := New("test_key", "in_000000000000", WithHTTPClient(&http.Client{Transport: tr}))
	require.NoError(t, err)

	p, err := c.Payouts.Get(context.Background(), "po_000000000000")
	require.NoError(t, err)
	require.Equal(t, "po_000000000000", p.ID)
	require.Equal(t, 2, tr.calls)

	tr = &countingTransport{status: []int{503}}
	c, err = New("test_key", "in_000000000000", WithHTTPClient(&http.Client{Transport: tr}), WithMaxRetries(0))
	require.NoError(t, err)
	_, err = c.Payouts.Get(context.Background(), "po_000000000000")
	require.Error(t, err)
	require.Equal(t, 1, tr.calls)
}

func TestClient_StringRedactsAPIKey(t *testing.T) {
	c, err := New("sk_secret_value", "in_000000000000")
	require.NoError(t, err)
	for _, s := range []string{fmt.Sprint(c), fmt.Sprintf("%v", c), fmt.Sprintf("%+v", c), fmt.Sprintf("%#v", c), c.String()} {
		require.NotContains(t, s, "sk_secret_value")
		require.Contains(t, s, "[REDACTED]")
	}
}
