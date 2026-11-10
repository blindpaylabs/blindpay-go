package wallets

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/stretchr/testify/require"
)

const createWalletOperation = "POST /v1/instances/{instance_id}/customers/{customer_id}/blockchain-wallets"

func boolPtr(b bool) *bool { return &b }

// TestWallets_CreateIsSelfCustody checks that is_self_custody (required for
// BR customers since 2026-10-10) reaches the wire from both create
// methods, and that nil omits the key instead of sending null (the spec types
// it as a non-nullable boolean).
func TestWallets_CreateIsSelfCustody(t *testing.T) {
	instanceID := "in_000000000000"
	customerID := "re_000000000000"
	address := "0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C"
	hash := "0x3c499c542cef5e3811e1192ce70d8cc03d5c3359"

	tests := []struct {
		name          string
		isSelfCustody *bool
		wantField     string
	}{
		{name: "nil omits the field", isSelfCustody: nil, wantField: ""},
		{name: "true", isSelfCustody: boolPtr(true), wantField: `,"is_self_custody":true`},
		{name: "false is sent, not omitted", isSelfCustody: boolPtr(false), wantField: `,"is_self_custody":false`},
	}

	methods := []struct {
		name     string
		wantBody string
		call     func(c *Client, isSelfCustody *bool) (*BlockchainWallet, error)
	}{
		{
			name:     "CreateWithAddress",
			wantBody: `{"name":"Wallet","network":"polygon","address":"` + address + `","is_account_abstraction":true%s}`,
			call: func(c *Client, isSelfCustody *bool) (*BlockchainWallet, error) {
				return c.CreateWithAddress(context.Background(), &CreateWithAddressParams{
					CustomerID:    customerID,
					Name:          "Wallet",
					Network:       types.NetworkPolygon,
					Address:       address,
					IsSelfCustody: isSelfCustody,
				})
			},
		},
		{
			name:     "CreateWithHash",
			wantBody: `{"name":"Wallet","network":"polygon","signature_tx_hash":"` + hash + `","is_account_abstraction":false%s}`,
			call: func(c *Client, isSelfCustody *bool) (*BlockchainWallet, error) {
				return c.CreateWithHash(context.Background(), &CreateWithHashParams{
					CustomerID:      customerID,
					Name:            "Wallet",
					Network:         types.NetworkPolygon,
					SignatureTxHash: hash,
					IsSelfCustody:   isSelfCustody,
				})
			},
		},
	}

	var full json.RawMessage
	for _, v := range blindpaytest.ResponseVariants(t, createWalletOperation) {
		if v.Name == "full" {
			full = v.Body
		}
	}
	require.NotNil(t, full)

	for _, m := range methods {
		m := m
		for _, tc := range tests {
			tc := tc
			t.Run(m.name+"/"+tc.name, func(t *testing.T) {
				client := NewClient(&config.Config{
					BaseURL:    "https://api.blindpay.com/v1",
					APIKey:     "test_key",
					InstanceID: instanceID,
					HTTPClient: &http.Client{Transport: &blindpaytest.RoundTripper{
						T:      t,
						Method: http.MethodPost,
						Path:   fmt.Sprintf("/v1/instances/%s/customers/%s/blockchain-wallets", instanceID, customerID),
						Query:  url.Values{},
						In:     json.RawMessage(fmt.Sprintf(m.wantBody, tc.wantField)),
						Out:    full,
					}},
					UserAgent: "blindpay-go/test",
				})

				wallet, err := m.call(client, tc.isSelfCustody)
				require.NoError(t, err)
				require.NotEmpty(t, wallet.ID)
			})
		}
	}
}

// TestWallets_CreateResponseFixtures decodes every shared fixture variant of
// the create response, including is_self_custody (true or null).
func TestWallets_CreateResponseFixtures(t *testing.T) {
	for _, v := range blindpaytest.ResponseVariants(t, createWalletOperation) {
		v := v
		t.Run(v.Name, func(t *testing.T) {
			var raw map[string]any
			require.NoError(t, json.Unmarshal(v.Body, &raw))

			var wallet BlockchainWallet
			require.NoError(t, json.Unmarshal(v.Body, &wallet))
			require.NotEmpty(t, wallet.ID)

			switch got := raw["is_self_custody"].(type) {
			case bool:
				require.NotNil(t, wallet.IsSelfCustody)
				require.Equal(t, got, *wallet.IsSelfCustody)
			default:
				require.Nil(t, wallet.IsSelfCustody)
			}
		})
	}
}
