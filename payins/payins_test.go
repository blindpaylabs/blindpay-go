package payins

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

func TestPayins_List(t *testing.T) {
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
							"id":"payin_123",
							"receiver_id":"recv_123",
							"status":"pending",
							"payin_quote_id":"quote_123",
							"instance_id":"inst_123",
							"type":"pix",
							"payment_method":"bank_transfer",
							"sender_amount":100.0,
							"receiver_amount":95.0,
							"token":"USDC",
							"partner_fee_amount":0,
							"total_fee_amount":5.0,
							"commercial_quotation":1.05,
							"blindpay_quotation":1.04,
							"currency":"BRL",
							"billing_fee":2.5,
							"name":"Test Payin",
							"address":"0x123",
							"network":"ethereum_mainnet",
							"created_at":"2024-01-01T00:00:00Z",
							"updated_at":"2024-01-01T00:00:00Z"
						}
					],
					"pagination":{"has_more":false,"next_page":0,"prev_page":0}
				}`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/payins", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.List(context.Background(), &ListParams{})
	require.NoError(t, err)
	require.Len(t, response.Data, 1)
	require.Equal(t, "payin_123", response.Data[0].ID)
	require.Equal(t, types.TransactionStatusPending, response.Data[0].Status)
}

func TestPayins_Get(t *testing.T) {
	instanceID := "inst_123"
	id := "payin_123"

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
					"payin_quote_id":"quote_123",
					"instance_id":"inst_123",
					"type":"pix",
					"payment_method":"bank_transfer",
					"sender_amount":100.0,
					"receiver_amount":95.0,
					"token":"USDC",
					"partner_fee_amount":0,
					"total_fee_amount":5.0,
					"commercial_quotation":1.05,
					"blindpay_quotation":1.04,
					"currency":"BRL",
					"billing_fee":2.5,
					"name":"Test Payin",
					"address":"0x123",
					"network":"ethereum_mainnet",
					"created_at":"2024-01-01T00:00:00Z",
					"updated_at":"2024-01-01T00:00:00Z"
				}`, id)),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/payins/%s", instanceID, id),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	payin, err := client.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, payin.ID)
	require.Equal(t, types.TransactionStatusCompleted, payin.Status)
}

func TestPayins_CreateQuote(t *testing.T) {
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
					"blockchain_wallet_id":"wallet_123",
					"currency_type":"sender",
					"payment_method":"pix",
					"request_amount":100.0,
					"token":"USDC",
					"cover_fees":false,
					"payer_rules":{"pix_allowed_tax_ids":null}
				}`),
				Out: json.RawMessage(fmt.Sprintf(`{
					"id":"%s",
					"commercial_quotation":5.25,
					"blindpay_quotation":5.20,
					"receiver_amount":95.0,
					"sender_amount":100.0,
					"partner_fee_amount":0,
					"flat_fee":2.5,
					"image_url":"https://example.com/qr.png",
					"pix_code":"00020126330014BR.GOV.BCB.PIX"
				}`, quoteID)),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/payin-quotes", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	quote, err := client.Quotes.Create(context.Background(), &CreateQuoteParams{
		BlockchainWalletID: "wallet_123",
		CurrencyType:       types.CurrencyTypeSender,
		PaymentMethod:      "pix",
		RequestAmount:      100.0,
		Token:              types.StablecoinTokenUSDC,
		CoverFees:          false,
		PayerRules:         PayerRules{},
	})
	require.NoError(t, err)
	require.Equal(t, quoteID, quote.ID)
	require.Equal(t, 100.0, quote.SenderAmount)
}
