package blindpay

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/blindpaylabs/blindpay-go/customers"
	"github.com/blindpaylabs/blindpay-go/payouts"
	"github.com/stretchr/testify/require"
)

// staticCursorTransport always answers with one row and the same
// next_page, like a server stuck on a cursor.
type staticCursorTransport struct {
	calls atomic.Int32
	body  string
}

func (s *staticCursorTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	s.calls.Add(1)
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}},
		Body: io.NopCloser(strings.NewReader(s.body)), Request: r}, nil
}

// TestAutoPaging_StopsWhenCursorRepeats: a server that returns the same
// next_page twice in a row ends the auto-pager with an error instead of an
// endless loop, for every resource that pages through the shared Pager.
func TestAutoPaging_StopsWhenCursorRepeats(t *testing.T) {
	newClient := func(body string) (*Client, *staticCursorTransport) {
		tr := &staticCursorTransport{body: body}
		c, err := New("test_key", "in_000000000000", WithHTTPClient(&http.Client{Transport: tr}))
		require.NoError(t, err)
		return c, tr
	}

	c, tr := newClient(`{"data":[{"id":"po_000000000001"}],"pagination":{"has_more":true,"next_page":"po_000000000001"}}`)
	pager := c.Payouts.ListAutoPaging(context.Background(), &payouts.ListParams{})
	var ids []string
	for pager.Next() {
		ids = append(ids, pager.Current().ID)
	}
	require.ErrorContains(t, pager.Err(), "cursor did not advance")
	require.Equal(t, []string{"po_000000000001", "po_000000000001"}, ids, "items already fetched are still yielded")
	require.Equal(t, int32(2), tr.calls.Load(), "first page, then the repeated cursor: no third request")

	c, tr = newClient(`{"data":[{"id":"cu_000000000001"}],"pagination":{"has_more":true,"next_page":"cu_000000000001"}}`)
	cp := c.Customers.ListAutoPaging(context.Background(), &customers.ListParams{StartingAfter: "cu_000000000001"})
	n := 0
	for cp.Next() {
		n++
	}
	require.ErrorContains(t, cp.Err(), "cursor did not advance", "a cursor equal to the caller's StartingAfter also stops")
	require.Equal(t, 1, n)
	require.Equal(t, int32(1), tr.calls.Load())
}
