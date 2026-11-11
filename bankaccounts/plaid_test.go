package bankaccounts

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/stretchr/testify/require"
)

func TestBankAccounts_CreatePlaidLink(t *testing.T) {
	v := blindpaytest.FindVariant(t, "POST /v1/instances/{instance_id}/customers/{customer_id}/bank-accounts/plaid", "full")
	redirect := "https://example.com/done"
	client := NewClient(&config.Config{
		BaseURL:    "https://api.blindpay.com/v1",
		APIKey:     "test_key",
		InstanceID: "in_000000000000",
		HTTPClient: &http.Client{Transport: &blindpaytest.RoundTripper{
			T:      t,
			Method: http.MethodPost,
			Path:   "/v1/instances/in_000000000000/customers/cu_000000000000/bank-accounts/plaid",
			In:     json.RawMessage(`{"redirect_url":"https://example.com/done"}`),
			Out:    v.Body,
		}},
		UserAgent: "blindpay-go/test",
	})

	link, err := client.CreatePlaidLink(context.Background(), &CreatePlaidLinkParams{CustomerID: "cu_000000000000", RedirectURL: &redirect})
	require.NoError(t, err)
	require.NotEmpty(t, link.URL)
}
