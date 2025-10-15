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
	instanceID := "in_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
					"data": [
						{
							"receiver_id": "re_000000000000",
							"id": "re_000000000000",
							"pix_code":"00020101021226790014br.gov.bcb.pix2557brcode.starkinfra.com/v2/bcf07f6c4110454e9fd6f120bab13e835204000053039865802BR5915Blind Pay, Inc.6010Vila Velha62070503***6304BCAB",
							"memo_code": "8K45GHBNT6BQ6462",
							"clabe": "014027000000000008",
							"status": "processing",
							"payin_quote_id": "pq_000000000000",
							"instance_id": "in_000000000000",
							"tracking_transaction": {
								"step": "processing",
								"status": "failed",
								"transaction_hash": "0x123...890",
								"completed_at": "2011-10-05T14:48:00.000Z"
							},
							"tracking_payment": {
								"step": "on_hold",
								"provider_name": "blockchain",
								"provider_transaction_id": "tx_123456789",
								"provider_status": "confirmed",
								"estimated_time_of_arrival": "2011-10-05T15:00:00.000Z",
								"completed_at": "2011-10-05T14:48:00.000Z"
							},
							"tracking_complete": {
								"step": "on_hold",
								"status": "completed",
								"transaction_hash": "0x123...890",
								"completed_at": "2011-10-05T14:48:00.000Z"
							},
							"tracking_partner_fee": {
								"step": "on_hold",
								"transaction_hash": "0x123...890",
								"completed_at": "2011-10-05T14:48:00.000Z"
							},
							"created_at": "2021-01-01T00:00:00Z",
							"updated_at": "2021-01-01T00:00:00Z",
							"image_url": "https://example.com/image.png",
							"first_name": "John",
							"last_name": "Doe",
							"legal_name": "Company Name Inc.",
							"type": "individual",
							"payment_method": "pix",
							"sender_amount": 5240,
							"receiver_amount": 1010,
							"token": "USDC",
							"partner_fee_amount": 150,
							"total_fee_amount": 1.53,
							"commercial_quotation": 495,
							"blindpay_quotation": 505,
							"currency": "BRL",
							"billing_fee": 100,
							"name": "Wallet Display Name",
							"address": "0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C",
							"network": "polygon",
							"blindpay_bank_details": {
								"routing_number": "121145349",
								"account_number": "621327727210181",
								"account_type": "Business checking",
								"swift_bic_code": "CHASUS33",
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
								"beneficiary": {
									"name": "BlindPay, Inc.",
									"address_line_1": "8 The Green, #19364",
									"address_line_2": "Dover, DE 19901"
								},
								"receiving_bank": {
									"name": "Column NA - Brex",
									"address_line_1": "1 Letterman Drive, Building A, Suite A4-700",
									"address_line_2": "San Francisco, CA 94129"
								}
							}
						}
					],
					"pagination": {
						"has_more": true,
						"next_page": 3,
						"prev_page": 1
					}
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
	require.Equal(t, "re_000000000000", response.Data[0].ID)
	require.Equal(t, types.TransactionStatusProcessing, response.Data[0].Status)
}

func TestPayins_Get(t *testing.T) {
	instanceID := "in_000000000000"
	id := "re_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
					"receiver_id": "re_000000000000",
					"id": "re_000000000000",
					"pix_code":"00020101021226790014br.gov.bcb.pix2557brcode.starkinfra.com/v2/bcf07f6c4110454e9fd6f120bab13e835204000053039865802BR5915Blind Pay, Inc.6010Vila Velha62070503***6304BCAB",
					"memo_code": "8K45GHBNT6BQ6462",
					"clabe": "014027000000000008",
					"status": "processing",
					"payin_quote_id": "pq_000000000000",
					"instance_id": "in_000000000000",
					"tracking_transaction": {
						"step": "processing",
						"status": "failed",
						"transaction_hash": "0x123...890",
						"completed_at": "2011-10-05T14:48:00.000Z"
					},
					"tracking_payment": {
						"step": "on_hold",
						"provider_name": "blockchain",
						"provider_transaction_id": "tx_123456789",
						"provider_status": "confirmed",
						"estimated_time_of_arrival": "2011-10-05T15:00:00.000Z",
						"completed_at": "2011-10-05T14:48:00.000Z"
					},
					"tracking_complete": {
						"step": "on_hold",
						"status": "completed",
						"transaction_hash": "0x123...890",
						"completed_at": "2011-10-05T14:48:00.000Z"
					},
					"tracking_partner_fee": {
						"step": "on_hold",
						"transaction_hash": "0x123...890",
						"completed_at": "2011-10-05T14:48:00.000Z"
					},
					"created_at": "2021-01-01T00:00:00Z",
					"updated_at": "2021-01-01T00:00:00Z",
					"image_url": "https://example.com/image.png",
					"first_name": "John",
					"last_name": "Doe",
					"legal_name": "Company Name Inc.",
					"type": "individual",
					"payment_method": "pix",
					"sender_amount": 5240,
					"receiver_amount": 1010,
					"token": "USDC",
					"partner_fee_amount": 150,
					"total_fee_amount": 1.53,
					"commercial_quotation": 495,
					"blindpay_quotation": 505,
					"currency": "BRL",
					"billing_fee": 100,
					"name": "Wallet Display Name",
					"address": "0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C",
					"network": "polygon",
					"blindpay_bank_details": {
						"routing_number": "121145349",
						"account_number": "621327727210181",
						"account_type": "Business checking",
						"swift_bic_code": "CHASUS33",
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
						"beneficiary": {
							"name": "BlindPay, Inc.",
							"address_line_1": "8 The Green, #19364",
							"address_line_2": "Dover, DE 19901"
						},
						"receiving_bank": {
							"name": "Column NA - Brex",
							"address_line_1": "1 Letterman Drive, Building A, Suite A4-700",
							"address_line_2": "San Francisco, CA 94129"
						}
					}
				}`),
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
	require.Equal(t, types.TransactionStatusProcessing, payin.Status)
}

func TestPayins_CreateQuote(t *testing.T) {
	instanceID := "in_000000000000"
	quoteID := "pq_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"blockchain_wallet_id":"wallet_123",
					"currency_type":"sender",
					"payment_method":"pix",
					"request_amount":5240,
					"token":"USDC",
					"cover_fees":false,
					"payer_rules":{"pix_allowed_tax_ids":null}
				}`),
				Out: json.RawMessage(fmt.Sprintf(`{
					"id":"%s",
					"commercial_quotation":495,
					"blindpay_quotation":505,
					"receiver_amount":1010,
					"sender_amount":5240,
					"partner_fee_amount":150,
					"flat_fee":100,
					"image_url":"https://example.com/image.png",
					"pix_code":"00020101021226790014br.gov.bcb.pix2557brcode.starkinfra.com/v2/bcf07f6c4110454e9fd6f120bab13e835204000053039865802BR5915Blind Pay, Inc.6010Vila Velha62070503***6304BCAB"
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
		RequestAmount:      5240,
		Token:              types.StablecoinTokenUSDC,
		CoverFees:          false,
		PayerRules:         PayerRules{},
	})
	require.NoError(t, err)
	require.Equal(t, quoteID, quote.ID)
	require.Equal(t, 5240.0, quote.SenderAmount)
}
