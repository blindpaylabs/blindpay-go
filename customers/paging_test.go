package customers

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/option"
	"github.com/blindpaylabs/blindpay-go/pagination"
	"github.com/stretchr/testify/require"
)

const pagingOperation = "GET /v1/instances/{instance_id}/customers"

func pagingClient(rt http.RoundTripper) *Client {
	return NewClient(&config.Config{
		BaseURL:    "https://api.blindpay.com/v1",
		APIKey:     "test_key",
		InstanceID: "in_000000000000",
		HTTPClient: &http.Client{Transport: rt},
		UserAgent:  "blindpay-go/test",
	})
}

func TestCustomers_ListFilterQuery(t *testing.T) {
	v := blindpaytest.FindVariant(t, pagingOperation, "full.b0")
	client := pagingClient(&blindpaytest.RoundTripper{
		T:      t,
		Method: http.MethodGet,
		Path:   "/v1/instances/in_000000000000/customers",
		Query:  url.Values{"limit": {"10"}, "starting_after": {"cu_1"}, "full_name": {"Ana Lima"}, "customer_name": {"Ana"}, "status": {"approved"}, "customer_id": {"re_1"}, "bank_account_id": {"ba_1"}, "country": {"BR"}},
		Out:    v.Body,
	})

	page, err := client.ListPage(context.Background(), &ListParams{Limit: "10", StartingAfter: "cu_1", FullName: "Ana Lima", CustomerName: "Ana", Status: KycStatus("approved"), CustomerID: "re_1", BankAccountID: "ba_1", Country: "BR"})
	require.NoError(t, err)
	require.Len(t, page.Data, v.RowCount(t))
	require.True(t, page.HasMore())
	require.Equal(t, v.NextPage(t), page.NextCursor())
}

func TestCustomers_ListPageBareArray(t *testing.T) {
	v := blindpaytest.FindVariant(t, pagingOperation, "full.b1")
	if v.Body == nil {
		t.Skip("no bare-array variant")
	}
	client := pagingClient(&blindpaytest.RoundTripper{T: t, Query: url.Values{}, Out: v.Body})

	page, err := client.ListPage(context.Background(), nil)
	require.NoError(t, err)
	require.Len(t, page.Data, v.RowCount(t))
	require.False(t, page.HasMore())
	require.Empty(t, page.NextCursor())
}

func TestCustomers_ListAutoPaging(t *testing.T) {
	first := blindpaytest.FindVariant(t, pagingOperation, "full.b0")
	cursor := first.NextPage(t)
	require.NotEmpty(t, cursor)

	seq := &blindpaytest.Sequence{T: t, Steps: []*blindpaytest.RoundTripper{
		{T: t, Method: http.MethodGet, Query: url.Values{"limit": {"100"}, "country": {"BR"}}, Out: first.Body,
			RequestHeader: http.Header{"X-Test": {"1"}}},
		{T: t, Method: http.MethodGet, Query: url.Values{"limit": {"100"}, "country": {"BR"}, "starting_after": {cursor}}, Out: first.AsLastPage(t),
			RequestHeader: http.Header{"X-Test": {"1"}}},
	}}
	client := pagingClient(seq)

	pager := client.ListAutoPaging(context.Background(), &ListParams{Country: "BR"}, blindpaytestHeader())
	n := 0
	for pager.Next() {
		require.NotEmpty(t, pager.Current().ID)
		n++
	}
	require.NoError(t, pager.Err())
	require.Equal(t, 2*first.RowCount(t), n)
	require.Equal(t, 2, seq.Calls())
}

func TestCustomers_ListAutoPagingStopsOnError(t *testing.T) {
	first := blindpaytest.FindVariant(t, pagingOperation, "full.b0")
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

func TestCustomers_ListAutoPagingRejectsEndingBefore(t *testing.T) {
	pager := pagingClient(&blindpaytest.RoundTripper{T: t}).ListAutoPaging(context.Background(), &ListParams{EndingBefore: "x"})
	require.False(t, pager.Next())
	require.True(t, errors.Is(pager.Err(), pagination.ErrEndingBeforeUnsupported))
}

func blindpaytestHeader() option.RequestOption {
	return option.WithHeader("X-Test", "1")
}
