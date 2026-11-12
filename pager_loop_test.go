package blindpay

import (
	"context"
	"encoding/json"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/blindpaylabs/blindpay-go/customers"
	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/payouts"
	"github.com/stretchr/testify/require"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// stuckServer answers every request with the same fixed body (one row and
// the same next_page), like a server stuck on a cursor. blindpaytest's
// RoundTripper returns the fixed body and asserts the request; it does not
// count requests, so it is wrapped to count them.
func stuckServer(t *testing.T, path, body string) (*Client, *int32) {
	rt := &blindpaytest.RoundTripper{T: t, Method: http.MethodGet, Path: path, Out: json.RawMessage(body)}
	var calls int32
	c, err := New("test_key", "in_000000000000", WithHTTPClient(&http.Client{Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		atomic.AddInt32(&calls, 1)
		return rt.RoundTrip(r)
	})}))
	require.NoError(t, err)
	return c, &calls
}

// TestAutoPaging_StopsWhenCursorRepeats: a server that returns the same
// next_page twice in a row ends the auto-pager with an error instead of an
// endless loop, for every resource that pages through the shared Pager.
func TestAutoPaging_StopsWhenCursorRepeats(t *testing.T) {
	c, calls := stuckServer(t, "/v1/instances/in_000000000000/payouts",
		`{"data":[{"id":"po_000000000001"}],"pagination":{"has_more":true,"next_page":"po_000000000001"}}`)
	pager := c.Payouts.ListAutoPaging(context.Background(), &payouts.ListParams{})
	var ids []string
	for pager.Next() {
		ids = append(ids, pager.Current().ID)
	}
	require.ErrorContains(t, pager.Err(), "cursor did not advance")
	require.Equal(t, []string{"po_000000000001", "po_000000000001"}, ids, "items already fetched are still yielded")
	require.Equal(t, int32(2), atomic.LoadInt32(calls), "first page, then the repeated cursor: no third request")

	c, calls = stuckServer(t, "/v1/instances/in_000000000000/customers",
		`{"data":[{"id":"cu_000000000001"}],"pagination":{"has_more":true,"next_page":"cu_000000000001"}}`)
	cp := c.Customers.ListAutoPaging(context.Background(), &customers.ListParams{StartingAfter: "cu_000000000001"})
	n := 0
	for cp.Next() {
		n++
	}
	require.ErrorContains(t, cp.Err(), "cursor did not advance", "a cursor equal to the caller's StartingAfter also stops")
	require.Equal(t, 1, n)
	require.Equal(t, int32(1), atomic.LoadInt32(calls))
}
