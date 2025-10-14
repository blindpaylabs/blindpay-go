package payouts

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

func TestPayouts_List(t *testing.T) {
	instanceID := "inst_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
					"data":[
						{
							"id":"payout_123",
							"receiver_id":"recv_123",
							"status":"pending",
							"sender_wallet_address":"0x123",
							"quote_id":"quote_123",
							"instance_id":"inst_123",
							"network":"ethereum_mainnet",
							"token":"USDC",
							"description":"Test payout",
							"sender_amount":100.0,
							"receiver_amount":95.0,
							"partner_fee_amount":0,
							"commercial_quotation":1.05,
							"blindpay_quotation":1.04,
							"total_fee_amount":5.0,
							"receiver_local_amount":500.0,
							"currency":"BRL",
							"name":"Test Bank Account",
							"type":"pix",
							"has_virtual_account":false,
							"created_at":"2024-01-01T00:00:00Z",
							"updated_at":"2024-01-01T00:00:00Z"
						}
					],
					"pagination":{"has_more":false,"next_page":0,"prev_page":0}
				}`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/payouts", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.List(context.Background(), &ListParams{})
	require.NoError(t, err)
	require.Len(t, response.Data, 1)
	require.Equal(t, "payout_123", response.Data[0].ID)
	require.Equal(t, types.TransactionStatusPending, response.Data[0].Status)
}

func TestPayouts_Get(t *testing.T) {
	instanceID := "inst_123"
	id := "payout_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(fmt.Sprintf(`{
					"id":"%s",
					"receiver_id":"recv_123",
					"status":"completed",
					"sender_wallet_address":"0x123",
					"quote_id":"quote_123",
					"instance_id":"inst_123",
					"network":"ethereum_mainnet",
					"token":"USDC",
					"description":"Test payout",
					"sender_amount":100.0,
					"receiver_amount":95.0,
					"partner_fee_amount":0,
					"commercial_quotation":1.05,
					"blindpay_quotation":1.04,
					"total_fee_amount":5.0,
					"receiver_local_amount":500.0,
					"currency":"BRL",
					"name":"Test Bank Account",
					"type":"pix",
					"has_virtual_account":false,
					"created_at":"2024-01-01T00:00:00Z",
					"updated_at":"2024-01-01T00:00:00Z"
				}`, id)),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/payouts/%s", instanceID, id),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	payout, err := client.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, payout.ID)
	require.Equal(t, types.TransactionStatusCompleted, payout.Status)
}

func TestPayouts_ExecuteEvm(t *testing.T) {
	instanceID := "inst_123"
	quoteID := "quote_123"
	payoutID := "payout_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"quote_id":"quote_123",
					"sender_wallet_address":"0x123"
				}`),
				Out: json.RawMessage(fmt.Sprintf(`{
					"id":"%s",
					"receiver_id":"recv_123",
					"status":"pending",
					"sender_wallet_address":"0x123",
					"signed_transaction":"0xabc",
					"quote_id":"quote_123",
					"instance_id":"inst_123",
					"network":"ethereum_mainnet",
					"token":"USDC",
					"has_virtual_account":false,
					"created_at":"2024-01-01T00:00:00Z",
					"updated_at":"2024-01-01T00:00:00Z"
				}`, payoutID)),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/payouts/evm", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	payout, err := client.CreateEvm(context.Background(), &CreateEvmParams{
		QuoteID:             quoteID,
		SenderWalletAddress: "0x123",
	})
	require.NoError(t, err)
	require.Equal(t, payoutID, payout.ID)
}

func TestPayouts_ExecuteStellar(t *testing.T) {
	instanceID := "inst_123"
	quoteID := "quote_123"
	payoutID := "payout_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"quote_id":"quote_123",
					"sender_wallet_address":"GABC123",
					"signed_transaction":"AAAA"
				}`),
				Out: json.RawMessage(fmt.Sprintf(`{
					"id":"%s",
					"receiver_id":"recv_123",
					"status":"pending",
					"sender_wallet_address":"GABC123",
					"signed_transaction":"AAAA",
					"quote_id":"quote_123",
					"instance_id":"inst_123",
					"network":"stellar_mainnet",
					"token":"USDC",
					"has_virtual_account":false,
					"created_at":"2024-01-01T00:00:00Z",
					"updated_at":"2024-01-01T00:00:00Z"
				}`, payoutID)),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/payouts/stellar", instanceID),
			},
		},
		UserAgent: "test",
	}

	signedTx := "AAAA"
	client := NewClient(cfg)
	payout, err := client.CreateStellar(context.Background(), &CreateStellarParams{
		QuoteID:             quoteID,
		SenderWalletAddress: "GABC123",
		SignedTransaction:   &signedTx,
	})
	require.NoError(t, err)
	require.Equal(t, payoutID, payout.ID)
}

func TestPayouts_AuthorizeStellarToken(t *testing.T) {
	instanceID := "inst_123"
	quoteID := "quote_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"quote_id":"quote_123",
					"sender_wallet_address":"GABC123"
				}`),
				Out: json.RawMessage(`{
					"transaction_hash":"abc123def456"
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/payouts/stellar/authorize-token", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.AuthorizeStellarToken(context.Background(), &AuthorizeStellarTokenParams{
		QuoteID:             quoteID,
		SenderWalletAddress: "GABC123",
	})
	require.NoError(t, err)
	require.Equal(t, "abc123def456", response.TransactionHash)
}
