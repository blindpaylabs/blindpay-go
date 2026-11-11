package wallets

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/stretchr/testify/require"
)

func updateTestConfig(rt http.RoundTripper) *config.Config {
	return &config.Config{
		BaseURL:    "https://api.blindpay.com/v1",
		APIKey:     "test_key",
		InstanceID: "in_000000000000",
		HTTPClient: &http.Client{Transport: rt},
		UserAgent:  "blindpay-go/test",
	}
}

func TestWallets_SetSelfCustody(t *testing.T) {
	v := blindpaytest.FindVariant(t, "PATCH /v1/instances/{instance_id}/customers/{customer_id}/blockchain-wallets/{id}", "full")
	client := NewClient(updateTestConfig(&blindpaytest.RoundTripper{
		T:      t,
		Method: http.MethodPatch,
		Path:   "/v1/instances/in_000000000000/customers/cu_000000000000/blockchain-wallets/bw_000000000000",
		In:     json.RawMessage(`{"is_self_custody":false}`),
		Out:    v.Body,
	}))

	wallet, err := client.SetSelfCustody(context.Background(), "cu_000000000000", "bw_000000000000", false)
	require.NoError(t, err)
	require.NotEmpty(t, wallet.ID)
}

func TestOfframpWallets_UpdateTriState(t *testing.T) {
	fee := "fe_000000000000"
	for name, tc := range map[string]struct {
		params UpdateOfframpWalletParams
		body   string
	}{
		"set":   {UpdateOfframpWalletParams{PartnerFeeID: &fee}, `{"partner_fee_id":"fe_000000000000"}`},
		"clear": {UpdateOfframpWalletParams{ClearPartnerFee: true}, `{"partner_fee_id":null}`},
		"keep":  {UpdateOfframpWalletParams{}, `{}`},
	} {
		tc := tc
		t.Run(name, func(t *testing.T) {
			client := NewOfframpClient(updateTestConfig(&blindpaytest.RoundTripper{
				T:      t,
				Method: http.MethodPut,
				Path:   "/v1/instances/in_000000000000/customers/cu_1/bank-accounts/ba_1/offramp-wallets/ow_1",
				In:     json.RawMessage(tc.body),
				Out:    json.RawMessage(`{"success":true}`),
			}))
			p := tc.params
			p.CustomerID, p.BankAccountID, p.ID = "cu_1", "ba_1", "ow_1"
			require.NoError(t, client.Update(context.Background(), &p))
		})
	}

	client := NewOfframpClient(updateTestConfig(nil))
	err := client.Update(context.Background(), &UpdateOfframpWalletParams{CustomerID: "c", BankAccountID: "b", ID: "i", PartnerFeeID: &fee, ClearPartnerFee: true})
	require.Error(t, err)
}
