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

// The RFI answer is the whole body (a flat map), not wrapped in
// {"response": ...}, and GET returns null when no RFI is open.
func TestRfi_SubmitSendsFlatBody(t *testing.T) {
	body := json.RawMessage(`{"source_of_funds":"salary","documents":["https://x/1.pdf"]}`)
	answer := rfi.Response{"source_of_funds": "salary", "documents": []string{"https://x/1.pdf"}}

	c := rtClient(t, &blindpaytest.RoundTripper{T: t, Method: http.MethodPost, Path: "/v1/instances/in_000000000000/customers/cu_000000000000/rfi", In: body, Out: json.RawMessage(`{"success":true}`)})
	res, err := c.Customers.SubmitRfi(context.Background(), tCustomer, answer)
	require.NoError(t, err)
	require.True(t, res.Success)

	c = rtClient(t, &blindpaytest.RoundTripper{T: t, Method: http.MethodPost, Path: "/v1/instances/in_000000000000/rfi", In: body, Out: json.RawMessage(`{"success":false,"reason":"no_pending_rfi"}`)})
	res, err = c.Instances.SubmitRfi(context.Background(), answer)
	require.NoError(t, err)
	require.False(t, res.Success)
	require.Equal(t, "no_pending_rfi", *res.Reason)
}

func TestRfi_GetNullMeansNoOpenRfi(t *testing.T) {
	c := rtClient(t, &blindpaytest.RoundTripper{T: t, Method: http.MethodGet, Path: "/v1/instances/in_000000000000/customers/cu_000000000000/rfi", Out: json.RawMessage(`null`)})
	got, err := c.Customers.GetRfi(context.Background(), tCustomer)
	require.NoError(t, err)
	require.Nil(t, got)

	c = rtClient(t, &blindpaytest.RoundTripper{T: t, Method: http.MethodGet, Path: "/v1/instances/in_000000000000/rfi", Out: json.RawMessage(`null`)})
	inst, err := c.Instances.GetRfi(context.Background())
	require.NoError(t, err)
	require.Nil(t, inst)
}

func TestRfi_GetDecodesSections(t *testing.T) {
	v := blindpaytest.FindVariant(t, "GET /v1/instances/{instance_id}/customers/{customer_id}/rfi", "full")
	c := rtClient(t, &blindpaytest.RoundTripper{T: t, Out: v.Body})
	got, err := c.Customers.GetRfi(context.Background(), tCustomer)
	require.NoError(t, err)
	require.NotEmpty(t, got.ID)
	require.NotEmpty(t, got.Request)
	require.NotEmpty(t, got.Request[0].Fields[0].Key)
}
