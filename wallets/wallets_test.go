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
	receiverID := "re_000000000000"
	instanceID := "in_000000000000"
	id := "bw_000000000000"
	address := "0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
			"name":"Wallet Display Name",
			"network":"polygon",
			"address":"0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C",
			"is_account_abstraction":true
		}`),
				Out: json.RawMessage(`{
					"id":"bw_000000000000",
					"name":"Wallet Display Name",
					"network":"polygon",
					"address":"0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C",
					"signature_tx_hash":null,
					"is_account_abstraction":true,
					"receiver_id":"re_000000000000"
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/blockchain-wallets", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	wallet, err := client.CreateWithAddress(context.Background(), &CreateWithAddressParams{
		ReceiverID: receiverID,
		Name:       "Wallet Display Name",
		Network:    types.NetworkPolygon,
		Address:    address,
	})
	require.NoError(t, err)
	require.Equal(t, id, wallet.ID)
	require.Equal(t, address, wallet.Address)
	require.True(t, wallet.IsAccountAbstraction)
}

func TestWallets_List(t *testing.T) {
	receiverID := "re_000000000000"
	instanceID := "in_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`[
					{
						"id":"bw_000000000000",
						"name":"Wallet Display Name",
						"network":"polygon",
						"address":"0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C",
						"signature_tx_hash":"0x3c499c542cef5e3811e1192ce70d8cc03d5c3359",
						"is_account_abstraction":false,
						"receiver_id":"re_000000000000"
					}
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
	require.Len(t, wallets, 1)
	require.Equal(t, "bw_000000000000", wallets[0].ID)
	require.Equal(t, "Wallet Display Name", wallets[0].Name)
	require.Equal(t, types.NetworkPolygon, wallets[0].Network)
	require.Equal(t, "0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C", wallets[0].Address)
	require.Equal(t, "0x3c499c542cef5e3811e1192ce70d8cc03d5c3359", wallets[0].SignatureTxHash)
	require.False(t, wallets[0].IsAccountAbstraction)
	require.Equal(t, "re_000000000000", wallets[0].ReceiverID)
}

func TestWallets_Get(t *testing.T) {
	receiverID := "re_000000000000"
	id := "bw_000000000000"
	instanceID := "in_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
					"id":"bw_000000000000",
					"name":"Wallet Display Name",
					"network":"polygon",
					"address":"0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C",
					"signature_tx_hash":"0x3c499c542cef5e3811e1192ce70d8cc03d5c3359",
					"is_account_abstraction":false,
					"receiver_id":"re_000000000000"
				}`),
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
	require.Equal(t, "Wallet Display Name", wallet.Name)
	require.Equal(t, types.NetworkPolygon, wallet.Network)
	require.Equal(t, "0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C", wallet.Address)
	require.Equal(t, "0x3c499c542cef5e3811e1192ce70d8cc03d5c3359", wallet.SignatureTxHash)
	require.False(t, wallet.IsAccountAbstraction)
	require.Equal(t, "re_000000000000", wallet.ReceiverID)
}

func TestWallets_Delete(t *testing.T) {
	receiverID := "re_000000000000"
	id := "bw_000000000000"
	instanceID := "in_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{"data":null}`),
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
	receiverID := "re_000000000000"
	instanceID := "in_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{"message":"random"}`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/blockchain-wallets/sign-message", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.GetWalletMessage(context.Background(), receiverID)
	require.NoError(t, err)
	require.Equal(t, "random", response.Message)
}
