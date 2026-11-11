package transfers

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/blindpaylabs/blindpay-go/option"
	"github.com/blindpaylabs/blindpay-go/pagination"
	"github.com/stretchr/testify/require"
)

const pagingOperation = "GET /v1/instances/{instance_id}/transfers"

func pagingClient(rt http.RoundTripper) *Client {
	return NewClient(&config.Config{
		BaseURL:    "https://api.blindpay.com/v1",
		APIKey:     "test_key",
		InstanceID: "in_000000000000",
		HTTPClient: &http.Client{Transport: rt},
		UserAgent:  "blindpay-go/test",
	})
}

func TestTransfers_ListFilterQuery(t *testing.T) {
	v := blindpaytest.FindVariant(t, pagingOperation, "full")
	client := pagingClient(&blindpaytest.RoundTripper{
		T:      t,
		Method: http.MethodGet,
		Path:   "/v1/instances/in_000000000000/transfers",
		Query:  url.Values{"limit": {"10"}, "starting_after": {"tr_1"}, "customer_id": {"re_1"}, "wallet_id": {"wa_1"}, "status": {"completed"}, "network": {"base"}, "token": {"USDC"}},
		Out:    v.Body,
	})

	page, err := client.ListPage(context.Background(), &ListParams{Limit: 10, StartingAfter: "tr_1", CustomerID: "re_1", WalletID: "wa_1", Status: types.TransactionStatusCompleted, Network: types.NetworkBase, Token: types.StablecoinToken("USDC")})
	require.NoError(t, err)
	require.Len(t, page.Data, v.RowCount(t))
	require.True(t, page.HasMore())
	require.Equal(t, v.NextPage(t), page.NextCursor())
}

func TestTransfers_ListAutoPaging(t *testing.T) {
	first := blindpaytest.FindVariant(t, pagingOperation, "full")
	cursor := first.NextPage(t)
	require.NotEmpty(t, cursor)

	seq := &blindpaytest.Sequence{T: t, Steps: []*blindpaytest.RoundTripper{
		{T: t, Method: http.MethodGet, Query: url.Values{"limit": {"100"}, "status": {"completed"}}, Out: first.Body,
			RequestHeader: http.Header{"X-Test": {"1"}}},
		{T: t, Method: http.MethodGet, Query: url.Values{"limit": {"100"}, "status": {"completed"}, "starting_after": {cursor}}, Out: first.AsLastPage(t),
			RequestHeader: http.Header{"X-Test": {"1"}}},
	}}
	client := pagingClient(seq)

	pager := client.ListAutoPaging(context.Background(), &ListParams{Status: types.TransactionStatusCompleted}, blindpaytestHeader())
	n := 0
	for pager.Next() {
		require.NotEmpty(t, pager.Current().ID)
		n++
	}
	require.NoError(t, pager.Err())
	require.Equal(t, 2*first.RowCount(t), n)
	require.Equal(t, 2, seq.Calls())
}

func TestTransfers_ListAutoPagingStopsOnError(t *testing.T) {
	first := blindpaytest.FindVariant(t, pagingOperation, "full")
	seq := &blindpaytest.Sequence{T: t, Steps: []*blindpaytest.RoundTripper{
		{T: t, Out: first.Body},
		{T: t, Status: http.StatusBadRequest, Out: []byte(`{"success":false,"code":"VALIDATION_ERROR","message":"bad_request"}`)},
	}}
	pager := pagingClient(seq).ListAutoPaging(context.Background(), nil)
	n := 0
	for pager.Next() {
		n++
	}
	require.Equal(t, first.RowCount(t), n)
	require.Error(t, pager.Err())
	require.Equal(t, 2, seq.Calls())
}

func TestTransfers_ListAutoPagingRejectsEndingBefore(t *testing.T) {
	pager := pagingClient(&blindpaytest.RoundTripper{T: t}).ListAutoPaging(context.Background(), &ListParams{EndingBefore: "x"})
	require.False(t, pager.Next())
	require.True(t, errors.Is(pager.Err(), pagination.ErrEndingBeforeUnsupported))
}

func blindpaytestHeader() option.RequestOption {
	return option.WithHeader("X-Test", "1")
}
