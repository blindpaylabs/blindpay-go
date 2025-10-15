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
	receiverID := "re_000000000000"
	instanceID := "in_000000000000"
	id := "va_000000000000"
	walletID := "bw_000000000000"

	outJson := `{
		"id": "va_000000000000",
		"us": {
			"ach": {
				"routing_number": "123456789",
				"account_number": "123456789"
			},
			"wire": {
				"routing_number": "123456789",
				"account_number": "123456789"
			},
			"rtp": {
				"routing_number": "123456789",
				"account_number": "123456789"
			},
			"swift_bic_code": "CHASUS33",
			"account_type": "Business checking",
			"beneficiary": {
				"name": "Receiver Name",
				"address_line_1": "8 The Green, #19364",
				"address_line_2": "Dover, DE 19901"
			},
			"receiving_bank": {
				"name": "JPMorgan Chase",
				"address_line_1": "270 Park Ave",
				"address_line_2": "New York, NY, 10017-2070"
			}
		},
		"token": "USDC",
		"blockchain_wallet_id": "bw_000000000000"
	}`

	inJson := `{
		"blockchain_wallet_id": "bw_000000000000",
		"token": "USDC"
	}`

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				In:     json.RawMessage(inJson),
				Out:    json.RawMessage(outJson),
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
	receiverID := "re_000000000000"
	instanceID := "in_000000000000"
	id := "va_000000000000"
	walletID := "bw_000000000000"

	outJson := `{
		"id": "va_000000000000",
		"us": {
			"ach": {
				"routing_number": "123456789",
				"account_number": "123456789"
			},
			"wire": {
				"routing_number": "123456789",
				"account_number": "123456789"
			},
			"rtp": {
				"routing_number": "123456789",
				"account_number": "123456789"
			},
			"swift_bic_code": "CHASUS33",
			"account_type": "Business checking",
			"beneficiary": {
				"name": "Receiver Name",
				"address_line_1": "8 The Green, #19364",
				"address_line_2": "Dover, DE 19901"
			},
			"receiving_bank": {
				"name": "JPMorgan Chase",
				"address_line_1": "270 Park Ave",
				"address_line_2": "New York, NY, 10017-2070"
			}
		},
		"token": "USDC",
		"blockchain_wallet_id": "bw_000000000000"
	}`

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(outJson),
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
	require.Equal(t, walletID, account.BlockchainWalletID)
}

func TestVirtualAccounts_Update(t *testing.T) {
	receiverID := "re_000000000000"
	instanceID := "in_000000000000"
	newWalletID := "bw_000000000000"

	inJson := `{
		"blockchain_wallet_id": "bw_000000000000",
		"token": "USDC"
	}`
	outJson := `{"data": null}`

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				In:     json.RawMessage(inJson),
				Out:    json.RawMessage(outJson),
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
		Token:              types.StablecoinTokenUSDC,
	})
	require.NoError(t, err)
}
