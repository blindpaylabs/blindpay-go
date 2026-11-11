package blindpay

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/rfi"
	"github.com/stretchr/testify/require"
)

const (
	opCustomerRfiGet    = "GET /v1/instances/{instance_id}/customers/{customer_id}/rfi"
	opCustomerRfiSubmit = "POST /v1/instances/{instance_id}/customers/{customer_id}/rfi"
	opInstanceRfiGet    = "GET /v1/instances/{instance_id}/rfi"
	opInstanceRfiSubmit = "POST /v1/instances/{instance_id}/rfi"
)

// The RFI answer is the whole body (a flat map), not wrapped in
// {"response": ...}. Both response branches (success, and success:false with
// a reason) decode.
func TestRfi_SubmitSendsFlatBody(t *testing.T) {
	body := json.RawMessage(`{"source_of_funds":"salary","documents":["https://x/1.pdf"]}`)
	answer := rfi.Response{"source_of_funds": "salary", "documents": []string{"https://x/1.pdf"}}

	for _, op := range []string{opCustomerRfiSubmit, opInstanceRfiSubmit} {
		for _, v := range blindpaytest.ResponseVariants(t, op) {
			op, v := op, v
			t.Run(op+"/"+v.Name, func(t *testing.T) {
				var want rfi.SubmitResponse
				require.NoError(t, json.Unmarshal(v.Body, &want))

				c := rtClient(t, &blindpaytest.RoundTripper{T: t, Method: http.MethodPost, In: body, Out: v.Body})
				var got *rfi.SubmitResponse
				var err error
				if op == opCustomerRfiSubmit {
					got, err = c.Customers.SubmitRfi(context.Background(), tCustomer, answer)
				} else {
					got, err = c.Instances.SubmitRfi(context.Background(), answer)
				}
				require.NoError(t, err)
				require.Equal(t, want, *got)
			})
		}
	}
}

func TestRfi_GetDecodesFixtures(t *testing.T) {
	for _, v := range blindpaytest.ResponseVariants(t, opCustomerRfiGet) {
		c := rtClient(t, &blindpaytest.RoundTripper{T: t, Method: http.MethodGet, Out: v.Body})
		got, err := c.Customers.GetRfi(context.Background(), tCustomer)
		require.NoError(t, err, v.Name)
		require.NotEmpty(t, got.ID, v.Name)
		require.NotEmpty(t, got.Request[0].Fields[0].Key, v.Name)
	}
	for _, v := range blindpaytest.ResponseVariants(t, opInstanceRfiGet) {
		c := rtClient(t, &blindpaytest.RoundTripper{T: t, Method: http.MethodGet, Out: v.Body})
		got, err := c.Instances.GetRfi(context.Background())
		require.NoError(t, err, v.Name)
		require.NotEmpty(t, got.ID, v.Name)
	}
}

// The API answers 200 with a JSON null when no RFI is open (the fixtures
// only cover the open case).
func TestRfi_GetNullMeansNoOpenRfi(t *testing.T) {
	c := rtClient(t, &blindpaytest.RoundTripper{T: t, Out: json.RawMessage(`null`)})
	got, err := c.Customers.GetRfi(context.Background(), tCustomer)
	require.NoError(t, err)
	require.Nil(t, got)

	inst, err := c.Instances.GetRfi(context.Background())
	require.NoError(t, err)
	require.Nil(t, inst)
}
