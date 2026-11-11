package payables

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/pagination"
	"github.com/stretchr/testify/require"
)

const listOperation = "GET /v1/instances/{instance_id}/payables"

func testClient(rt http.RoundTripper) *Client {
	return NewClient(&config.Config{
		BaseURL:    "https://api.blindpay.com/v1",
		APIKey:     "test_key",
		InstanceID: "in_000000000000",
		HTTPClient: &http.Client{Transport: rt},
		UserAgent:  "blindpay-go/test",
	})
}

func TestPayables_Create(t *testing.T) {
	v := blindpaytest.FindVariant(t, "POST /v1/instances/{instance_id}/payables", "full")
	barcode := "34191790010104351004791020150008291070026000"
	discount := int64(500)
	client := testClient(&blindpaytest.RoundTripper{
		T:      t,
		Method: http.MethodPost,
		Path:   "/v1/instances/in_000000000000/payables",
		In: json.RawMessage(`{"customer_id":"cu_000000000000","currency":"BRL","boleto_barcode":"` + barcode + `","discount":500,
			"line_items":[{"name":"Rent","quantity":1,"price":100000}]}`),
		Out: v.Body,
	})

	price := int64(100000)
	payable, err := client.Create(context.Background(), &CreateParams{
		CustomerID:    "cu_000000000000",
		Currency:      "BRL",
		BoletoBarcode: &barcode,
		Discount:      &discount,
		LineItems:     []LineItem{{Name: "Rent", Quantity: 1, Price: &price}},
	})
	require.NoError(t, err)
	require.NotEmpty(t, payable.ID)
}

func TestPayables_GetAndDelete(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodDelete} {
		method := method
		t.Run(method, func(t *testing.T) {
			v := blindpaytest.FindVariant(t, method+" /v1/instances/{instance_id}/payables/{id}", "full")
			client := testClient(&blindpaytest.RoundTripper{T: t, Method: method, Path: "/v1/instances/in_000000000000/payables/pa%2F1", Out: v.Body})

			var payable *Payable
			var err error
			if method == http.MethodGet {
				payable, err = client.Get(context.Background(), "pa/1")
			} else {
				payable, err = client.Delete(context.Background(), "pa/1")
			}
			require.NoError(t, err)
			require.NotEmpty(t, payable.ID)
		})
	}

	_, err := testClient(nil).Get(context.Background(), "")
	require.Error(t, err)
}

func TestPayables_ListFilterQuery(t *testing.T) {
	v := blindpaytest.FindVariant(t, listOperation, "full.b0")
	client := testClient(&blindpaytest.RoundTripper{
		T:      t,
		Method: http.MethodGet,
		Path:   "/v1/instances/in_000000000000/payables",
		Query: url.Values{"limit": {"10"}, "starting_after": {"pa_1"}, "status": {"draft"}, "customer_id": {"cu_1"},
			"customer_name": {"Ana"}, "type": {"boleto"}, "beneficiary_name": {"Acme"}},
		Out: v.Body,
	})

	page, err := client.ListPage(context.Background(), &ListParams{Limit: 10, StartingAfter: "pa_1", Status: StatusDraft,
		CustomerID: "cu_1", CustomerName: "Ana", Type: TypeBoleto, BeneficiaryName: "Acme"})
	require.NoError(t, err)
	require.Len(t, page.Data, v.RowCount(t))
	require.True(t, page.HasMore())
	require.Equal(t, v.NextPage(t), page.NextCursor())
}

func TestPayables_ListBothShapes(t *testing.T) {
	for _, name := range []string{"full.b0", "full.b1"} {
		v := blindpaytest.FindVariant(t, listOperation, name)
		client := testClient(&blindpaytest.RoundTripper{T: t, Query: url.Values{}, Out: v.Body})
		rows, err := client.List(context.Background(), nil)
		require.NoError(t, err, name)
		require.Len(t, rows, v.RowCount(t), name)
	}
}

func TestPayables_ListAutoPaging(t *testing.T) {
	first := blindpaytest.FindVariant(t, listOperation, "full.b0")
	cursor := first.NextPage(t)
	require.NotEmpty(t, cursor)

	seq := &blindpaytest.Sequence{T: t, Steps: []*blindpaytest.RoundTripper{
		{T: t, Method: http.MethodGet, Query: url.Values{"limit": {"100"}, "status": {"completed"}}, Out: first.Body},
		{T: t, Method: http.MethodGet, Query: url.Values{"limit": {"100"}, "status": {"completed"}, "starting_after": {cursor}}, Out: first.AsLastPage(t)},
	}}
	pager := testClient(seq).ListAutoPaging(context.Background(), &ListParams{Status: StatusCompleted})
	n := 0
	for pager.Next() {
		n++
	}
	require.NoError(t, pager.Err())
	require.Equal(t, 2*first.RowCount(t), n)
	require.Equal(t, 2, seq.Calls())

	pager = testClient(nil).ListAutoPaging(context.Background(), &ListParams{EndingBefore: "pa_1"})
	require.False(t, pager.Next())
	require.True(t, errors.Is(pager.Err(), pagination.ErrEndingBeforeUnsupported))
}
