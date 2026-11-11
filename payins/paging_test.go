package payins

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

const pagingOperation = "GET /v1/instances/{instance_id}/payins"

func pagingClient(rt http.RoundTripper) *Client {
	return NewClient(&config.Config{
		BaseURL:    "https://api.blindpay.com/v1",
		APIKey:     "test_key",
		InstanceID: "in_000000000000",
		HTTPClient: &http.Client{Transport: rt},
		UserAgent:  "blindpay-go/test",
	})
}

func TestPayins_ListFilterQuery(t *testing.T) {
	v := blindpaytest.FindVariant(t, pagingOperation, "full.b0")
	client := pagingClient(&blindpaytest.RoundTripper{
		T:      t,
		Method: http.MethodGet,
		Path:   "/v1/instances/in_000000000000/payins",
		Query:  url.Values{"limit": {"10"}, "starting_after": {"pi_1"}, "customer_id": {"re_1"}, "status": {"completed"}, "customer_name": {"Ana"}, "bank_account_id": {"ba_1"}, "country": {"BR"}, "payment_method": {"pix"}, "network": {"base"}, "token": {"USDC"}},
		Out:    v.Body,
	})

	page, err := client.ListPage(context.Background(), &ListParams{Limit: 10, StartingAfter: "pi_1", CustomerID: "re_1", Status: types.TransactionStatusCompleted, CustomerName: "Ana", BankAccountID: "ba_1", Country: "BR", PaymentMethod: types.PayinPaymentMethodPix, Network: types.NetworkBase, Token: types.StablecoinToken("USDC")})
	require.NoError(t, err)
	require.Len(t, page.Data, v.RowCount(t))
	require.True(t, page.HasMore())
	require.Equal(t, v.NextPage(t), page.NextCursor())
}

func TestPayins_ListPageBareArray(t *testing.T) {
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

func TestPayins_ListAutoPaging(t *testing.T) {
	first := blindpaytest.FindVariant(t, pagingOperation, "full.b0")
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

func TestPayins_ListAutoPagingStopsOnError(t *testing.T) {
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

func TestPayins_ListAutoPagingRejectsEndingBefore(t *testing.T) {
	pager := pagingClient(&blindpaytest.RoundTripper{T: t}).ListAutoPaging(context.Background(), &ListParams{EndingBefore: "x"})
	require.False(t, pager.Next())
	require.True(t, errors.Is(pager.Err(), pagination.ErrEndingBeforeUnsupported))
}

func blindpaytestHeader() option.RequestOption {
	return option.WithHeader("X-Test", "1")
}
