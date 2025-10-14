package wallets

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

func TestWallets_CreateWithAddress(t *testing.T) {
	receiverID := "recv_123"
	instanceID := "inst_123"
	id := "wallet_123"
	address := "0x1234567890abcdef"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"name":"My Wallet",
					"network":"ethereum_mainnet",
					"address":"0x1234567890abcdef",
					"is_account_abstraction":true
				}`),
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"My Wallet","network":"ethereum_mainnet","address":"%s","is_account_abstraction":false,"receiver_id":"%s"}`, id, address, receiverID)),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/blockchain-wallets", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	wallet, err := client.CreateWithAddress(context.Background(), &CreateWithAddressParams{
		ReceiverID: receiverID,
		Name:       "My Wallet",
		Network:    types.NetworkEthereumMainnet,
		Address:    address,
	})
	require.NoError(t, err)
	require.Equal(t, id, wallet.ID)
	require.Equal(t, address, wallet.Address)
}

func TestWallets_List(t *testing.T) {
	receiverID := "recv_123"
	instanceID := "inst_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`[
					{"id":"wallet_123","name":"Wallet 1","network":"ethereum_mainnet","address":"0x123","is_account_abstraction":false,"receiver_id":"recv_123"},
					{"id":"wallet_456","name":"Wallet 2","network":"polygon_mainnet","address":"0x456","is_account_abstraction":false,"receiver_id":"recv_123"}
				]`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/blockchain-wallets", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	wallets, err := client.List(context.Background(), receiverID)
	require.NoError(t, err)
	require.Len(t, wallets, 2)
	require.Equal(t, "wallet_123", wallets[0].ID)
}

func TestWallets_Get(t *testing.T) {
	receiverID := "recv_123"
	id := "wallet_123"
	instanceID := "inst_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"My Wallet","network":"ethereum_mainnet","address":"0x123","is_account_abstraction":false,"receiver_id":"%s"}`, id, receiverID)),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/blockchain-wallets/%s", instanceID, receiverID, id),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	wallet, err := client.Get(context.Background(), receiverID, id)
	require.NoError(t, err)
	require.Equal(t, id, wallet.ID)
}

func TestWallets_Delete(t *testing.T) {
	receiverID := "recv_123"
	id := "wallet_123"
	instanceID := "inst_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{}`),
				Method: http.MethodDelete,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/blockchain-wallets/%s", instanceID, receiverID, id),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	err := client.Delete(context.Background(), receiverID, id)
	require.NoError(t, err)
}

func TestWallets_GetWalletMessage(t *testing.T) {
	receiverID := "recv_123"
	instanceID := "inst_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{"message":"Sign this message to verify wallet ownership"}`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/blockchain-wallets/sign-message", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.GetWalletMessage(context.Background(), receiverID)
	require.NoError(t, err)
	require.Equal(t, "Sign this message to verify wallet ownership", response.Message)
}
