package virtualaccounts

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

func TestVirtualAccounts_Create(t *testing.T) {
	receiverID := "recv_123"
	instanceID := "inst_123"
	id := "va_123"
	walletID := "wallet_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"blockchain_wallet_id":"wallet_123",
					"token":"USDC"
				}`),
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","blockchain_wallet_id":"%s","token":"USDC","us":{"ach":{"routing_number":"123456789","account_number":"987654321"}}}`, id, walletID)),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/virtual-accounts", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	account, err := client.Create(context.Background(), &CreateParams{
		ReceiverID:         receiverID,
		BlockchainWalletID: walletID,
		Token:              types.StablecoinTokenUSDC,
	})
	require.NoError(t, err)
	require.Equal(t, id, account.ID)
	require.Equal(t, walletID, account.BlockchainWalletID)
}

func TestVirtualAccounts_Get(t *testing.T) {
	receiverID := "recv_123"
	instanceID := "inst_123"
	id := "va_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","blockchain_wallet_id":"wallet_123","token":"USDC","us":{"ach":{"routing_number":"123456789","account_number":"987654321"}}}`, id)),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/virtual-accounts", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	account, err := client.Get(context.Background(), receiverID)
	require.NoError(t, err)
	require.Equal(t, id, account.ID)
}

func TestVirtualAccounts_Update(t *testing.T) {
	receiverID := "recv_123"
	instanceID := "inst_123"
	newWalletID := "wallet_456"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"blockchain_wallet_id":"wallet_456",
					"token":"USDT"
				}`),
				Out:    json.RawMessage(`{}`),
				Method: http.MethodPut,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/virtual-accounts", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	err := client.Update(context.Background(), &UpdateParams{
		ReceiverID:         receiverID,
		BlockchainWalletID: newWalletID,
		Token:              types.StablecoinTokenUSDT,
	})
	require.NoError(t, err)
}

// Note: Delete method doesn't exist in the virtual accounts client
