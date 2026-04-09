package custodialwallets

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestCustodialWallets_List(t *testing.T) {
	instanceID := "in_000000000000"
	receiverID := "re_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`[{
					"id":"cw_000000000000",
					"name":"My Wallet",
					"external_id":null,
					"address":"0x123...890",
					"network":"base",
					"created_at":"2021-01-01T00:00:00Z"
				}]`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/wallets", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	wallets, err := client.List(context.Background(), receiverID)
	require.NoError(t, err)
	require.Len(t, wallets, 1)
	require.Equal(t, "cw_000000000000", wallets[0].ID)
}

func TestCustodialWallets_Create(t *testing.T) {
	instanceID := "in_000000000000"
	receiverID := "re_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
					"id":"cw_000000000000",
					"name":"My Wallet",
					"external_id":null,
					"address":"0x123...890",
					"network":"base",
					"created_at":"2021-01-01T00:00:00Z"
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/wallets", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	wallet, err := client.Create(context.Background(), &CreateParams{
		ReceiverID: receiverID,
		Name:       "My Wallet",
		Network:    types.NetworkBase,
	})
	require.NoError(t, err)
	require.Equal(t, "cw_000000000000", wallet.ID)
	require.Equal(t, types.NetworkBase, wallet.Network)
}

func TestCustodialWallets_GetBalance(t *testing.T) {
	instanceID := "in_000000000000"
	receiverID := "re_000000000000"
	walletID := "cw_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
					"USDC":{"address":"0x123","id":"tok_1","symbol":"USDC","amount":100.5},
					"USDT":{"address":"0x456","id":"tok_2","symbol":"USDT","amount":50.25},
					"USDB":{"address":"0x789","id":"tok_3","symbol":"USDB","amount":0}
				}`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/wallets/%s/balance", instanceID, receiverID, walletID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	balance, err := client.GetBalance(context.Background(), receiverID, walletID)
	require.NoError(t, err)
	require.Equal(t, 100.5, balance.USDC.Amount)
	require.Equal(t, 50.25, balance.USDT.Amount)
}
