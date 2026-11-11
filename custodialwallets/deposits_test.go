package custodialwallets

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/pagination"
	"github.com/stretchr/testify/require"
)

const depositsOperation = "GET /v1/instances/{instance_id}/customers/{customer_id}/wallets/{id}/deposits"

func depositsClient(rt http.RoundTripper) *Client {
	return NewClient(&config.Config{
		BaseURL:    "https://api.blindpay.com/v1",
		APIKey:     "test_key",
		InstanceID: "in_000000000000",
		HTTPClient: &http.Client{Transport: rt},
		UserAgent:  "blindpay-go/test",
	})
}

func TestCustodialWallets_ListDeposits(t *testing.T) {
	v := blindpaytest.FindVariant(t, depositsOperation, "full")
	client := depositsClient(&blindpaytest.RoundTripper{
		T:      t,
		Method: http.MethodGet,
		Path:   "/v1/instances/in_000000000000/customers/cu_1/wallets/wa_1/deposits",
		Query:  url.Values{"limit": {"50"}, "starting_after": {"de_1"}, "token": {"USDC"}, "status": {"completed"}},
		Out:    v.Body,
	})

	page, err := client.ListDeposits(context.Background(), "cu_1", "wa_1", &ListDepositsParams{Limit: 50, StartingAfter: "de_1", Token: "USDC", Status: DepositStatusCompleted})
	require.NoError(t, err)
	require.Len(t, page.Data, v.RowCount(t))
	require.NotEmpty(t, page.Data[0].ID)
}

func TestCustodialWallets_ListDepositsAutoPaging(t *testing.T) {
	first := blindpaytest.FindVariant(t, depositsOperation, "full")
	cursor := first.NextPage(t)
	require.NotEmpty(t, cursor)

	seq := &blindpaytest.Sequence{T: t, Steps: []*blindpaytest.RoundTripper{
		{T: t, Query: url.Values{"limit": {"100"}}, Out: first.Body},
		{T: t, Query: url.Values{"limit": {"100"}, "starting_after": {cursor}}, Out: first.AsLastPage(t)},
	}}
	pager := depositsClient(seq).ListDepositsAutoPaging(context.Background(), "cu_1", "wa_1", nil)
	n := 0
	for pager.Next() {
		n++
	}
	require.NoError(t, pager.Err())
	require.Equal(t, 2*first.RowCount(t), n)

	pager = depositsClient(nil).ListDepositsAutoPaging(context.Background(), "cu_1", "wa_1", &ListDepositsParams{EndingBefore: "x"})
	require.False(t, pager.Next())
	require.True(t, errors.Is(pager.Err(), pagination.ErrEndingBeforeUnsupported))
}
