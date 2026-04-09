package transfers

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

func TestTransfers_CreateQuote(t *testing.T) {
	instanceID := "in_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
					"id":"tq_000000000000",
					"expires_at":1700000000,
					"commercial_quotation":1.0,
					"blindpay_quotation":0.99,
					"receiver_amount":100,
					"sender_amount":101,
					"partner_fee_amount":1,
					"flat_fee":0.5
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/transfer-quotes", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.Quotes.Create(context.Background(), &CreateQuoteParams{
		WalletID:              "wl_000000000000",
		SenderToken:           types.StablecoinTokenUSDC,
		ReceiverWalletAddress: "0x123...890",
		ReceiverToken:         types.StablecoinTokenUSDC,
		ReceiverNetwork:       types.NetworkBase,
		RequestAmount:         100,
		CoverFees:             true,
		AmountReference:       "sender",
	})
	require.NoError(t, err)
	require.Equal(t, "tq_000000000000", response.ID)
}

func TestTransfers_Create(t *testing.T) {
	instanceID := "in_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
					"id":"tr_000000000000",
					"status":"processing",
					"tracking_bridge_swap":{"step":"processing"},
					"tracking_complete":{"step":"processing"},
					"tracking_paymaster":{"step":"processing"},
					"tracking_transaction_monitoring":{"step":"processing"},
					"tracking_partner_fee":{"step":"processing"}
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/transfers", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.Create(context.Background(), &CreateParams{
		TransferQuoteID: "tq_000000000000",
	})
	require.NoError(t, err)
	require.Equal(t, "tr_000000000000", response.ID)
	require.Equal(t, types.TransactionStatusProcessing, response.Status)
}

func TestTransfers_List(t *testing.T) {
	instanceID := "in_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
					"data":[{
						"id":"tr_000000000000",
						"status":"processing",
						"transfer_quote_id":"tq_000000000000",
						"instance_id":"in_000000000000",
						"tracking_transaction_monitoring":{"step":"processing"},
						"tracking_paymaster":{"step":"processing"},
						"tracking_bridge_swap":{"step":"processing"},
						"tracking_complete":{"step":"processing"},
						"tracking_partner_fee":{"step":"processing"},
						"created_at":"2021-01-01T00:00:00Z",
						"updated_at":"2021-01-01T00:00:00Z",
						"wallet_id":"wl_000000000000",
						"sender_token":"USDC",
						"sender_amount":100,
						"receiver_amount":99,
						"receiver_network":"base",
						"receiver_token":"USDC",
						"receiver_wallet_address":"0x123...890"
					}],
					"pagination":{"has_more":false,"next_page":0,"prev_page":0}
				}`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/transfers", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.List(context.Background(), nil)
	require.NoError(t, err)
	require.Len(t, response.Data, 1)
	require.Equal(t, "tr_000000000000", response.Data[0].ID)
}

func TestTransfers_Get(t *testing.T) {
	instanceID := "in_000000000000"
	transferID := "tr_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
					"id":"tr_000000000000",
					"status":"completed",
					"transfer_quote_id":"tq_000000000000",
					"instance_id":"in_000000000000",
					"tracking_transaction_monitoring":{"step":"completed"},
					"tracking_paymaster":{"step":"completed"},
					"tracking_bridge_swap":{"step":"completed"},
					"tracking_complete":{"step":"completed"},
					"tracking_partner_fee":{"step":"completed"},
					"created_at":"2021-01-01T00:00:00Z",
					"updated_at":"2021-01-01T00:00:00Z",
					"wallet_id":"wl_000000000000",
					"sender_token":"USDC",
					"sender_amount":100,
					"receiver_amount":99,
					"receiver_network":"base",
					"receiver_token":"USDC",
					"receiver_wallet_address":"0x123...890"
				}`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/transfers/%s", instanceID, transferID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	transfer, err := client.Get(context.Background(), transferID)
	require.NoError(t, err)
	require.Equal(t, transferID, transfer.ID)
	require.Equal(t, types.TransactionStatusCompleted, transfer.Status)
}
