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
	instanceID := "in_000000000000"

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
							"receiver_id":"re_000000000000",
							"id":"pa_000000000000",
							"status":"processing",
							"sender_wallet_address":"0x123...890",
							"signed_transaction":"AAA...Zey8y0A",
							"quote_id":"qu_000000000000",
							"instance_id":"in_000000000000",
							"tracking_transaction":{
								"step":"processing",
								"status":"failed",
								"transaction_hash":"0x123...890",
								"completed_at":"2011-10-05T14:48:00.000Z"
							},
							"tracking_payment":{
								"step":"on_hold",
								"provider_name":"blockchain",
								"provider_transaction_id":"0x123...890",
								"provider_status":"canceled",
								"estimated_time_of_arrival":"5_min",
								"completed_at":"2011-10-05T14:48:00.000Z"
							},
							"tracking_liquidity":{
								"step":"processing",
								"provider_transaction_id":"0x123...890",
								"provider_status":"deposited",
								"estimated_time_of_arrival":"1_business_day",
								"completed_at":"2011-10-05T14:48:00.000Z"
							},
							"tracking_complete":{
								"step":"on_hold",
								"status":"tokens_refunded",
								"transaction_hash":"0x123...890",
								"completed_at":"2011-10-05T14:48:00.000Z"
							},
							"tracking_partner_fee":{
								"step":"on_hold",
								"transaction_hash":"0x123...890",
								"completed_at":"2011-10-05T14:48:00.000Z"
							},
							"created_at":"2021-01-01T00:00:00Z",
							"updated_at":"2021-01-01T00:00:00Z",
							"image_url":"https://example.com/image.png",
							"first_name":"John",
							"last_name":"Doe",
							"legal_name":"Company Name Inc.",
							"network":"sepolia",
							"token":"USDC",
							"description":"Memo code or description, only works with USD and BRL",
							"sender_amount":1010,
							"receiver_amount":5240,
							"partner_fee_amount":150,
							"commercial_quotation":495,
							"blindpay_quotation":485,
							"total_fee_amount":1.5,
							"receiver_local_amount":1000,
							"currency":"BRL",
							"transaction_document_file":"https://example.com/image.png",
							"transaction_document_type":"invoice",
							"transaction_document_id":"1234567890",
							"name":"Bank Account Name",
							"type":"wire",
							"pix_key":"14947677768",
							"account_number":"1001001234",
							"routing_number":"012345678",
							"country":"US",
							"account_class":"individual",
							"address_line_1":"Address line 1",
							"address_line_2":"Address line 2",
							"city":"City",
							"state_province_region":"State/Province/Region",
							"postal_code":"Postal code",
							"account_type":"checking",
							"ach_cop_beneficiary_first_name":"Fernando",
							"ach_cop_bank_account":"12345678",
							"ach_cop_bank_code":"051",
							"ach_cop_beneficiary_last_name":"Guzman Alarcón",
							"ach_cop_document_id":"1661105408",
							"ach_cop_document_type":"CC",
							"ach_cop_email":"fernando.guzman@gmail.com",
							"beneficiary_name":"Individual full name or business name",
							"spei_clabe":"5482347403740546",
							"spei_protocol":"clabe",
							"spei_institution_code":"40002",
							"swift_beneficiary_country":"MX",
							"swift_code_bic":"123456789",
							"swift_account_holder_name":"John Doe",
							"swift_account_number_iban":"123456789",
							"transfers_account":"BM123123123123",
							"transfers_type":"CVU",
							"has_virtual_account":true
						}
					],
					"pagination":{
						"has_more":true,
						"next_page":3,
						"prev_page":1
					}
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
	require.Equal(t, "pa_000000000000", response.Data[0].ID)
	require.Equal(t, types.TransactionStatusProcessing, response.Data[0].Status)
}

func TestPayouts_Export(t *testing.T) {
	instanceID := "in_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`[
					{
						"receiver_id":"re_000000000000",
						"id":"pa_000000000000",
						"status":"processing",
						"sender_wallet_address":"0x123...890",
						"signed_transaction":"AAA...Zey8y0A",
						"quote_id":"qu_000000000000",
						"instance_id":"in_000000000000",
						"tracking_transaction":{
							"step":"processing",
							"status":"failed",
							"transaction_hash":"0x123...890",
							"completed_at":"2011-10-05T14:48:00.000Z"
						},
						"tracking_payment":{
							"step":"on_hold",
							"provider_name":"blockchain",
							"provider_transaction_id":"0x123...890",
							"provider_status":"canceled",
							"estimated_time_of_arrival":"5_min",
							"completed_at":"2011-10-05T14:48:00.000Z"
						},
						"tracking_liquidity":{
							"step":"processing",
							"provider_transaction_id":"0x123...890",
							"provider_status":"deposited",
							"estimated_time_of_arrival":"1_business_day",
							"completed_at":"2011-10-05T14:48:00.000Z"
						},
						"tracking_complete":{
							"step":"on_hold",
							"status":"tokens_refunded",
							"transaction_hash":"0x123...890",
							"completed_at":"2011-10-05T14:48:00.000Z"
						},
						"tracking_partner_fee":{
							"step":"on_hold",
							"transaction_hash":"0x123...890",
							"completed_at":"2011-10-05T14:48:00.000Z"
						},
						"created_at":"2021-01-01T00:00:00Z",
						"updated_at":"2021-01-01T00:00:00Z",
						"image_url":"https://example.com/image.png",
						"first_name":"John",
						"last_name":"Doe",
						"legal_name":"Company Name Inc.",
						"network":"sepolia",
						"token":"USDC",
						"description":"Memo code or description, only works with USD and BRL",
						"sender_amount":1010,
						"receiver_amount":5240,
						"partner_fee_amount":150,
						"commercial_quotation":495,
						"blindpay_quotation":485,
						"total_fee_amount":1.5,
						"receiver_local_amount":1000,
						"currency":"BRL",
						"transaction_document_file":"https://example.com/image.png",
						"transaction_document_type":"invoice",
						"transaction_document_id":"1234567890",
						"name":"Bank Account Name",
						"type":"wire",
						"pix_key":"14947677768",
						"account_number":"1001001234",
						"routing_number":"012345678",
						"country":"US",
						"account_class":"individual",
						"address_line_1":"Address line 1",
						"address_line_2":"Address line 2",
						"city":"City",
						"state_province_region":"State/Province/Region",
						"postal_code":"Postal code",
						"account_type":"checking",
						"ach_cop_beneficiary_first_name":"Fernando",
						"ach_cop_bank_account":"12345678",
						"ach_cop_bank_code":"051",
						"ach_cop_beneficiary_last_name":"Guzman Alarcón",
						"ach_cop_document_id":"1661105408",
						"ach_cop_document_type":"CC",
						"ach_cop_email":"fernando.guzman@gmail.com",
						"beneficiary_name":"Individual full name or business name",
						"spei_clabe":"5482347403740546",
						"spei_protocol":"clabe",
						"spei_institution_code":"40002",
						"swift_beneficiary_country":"MX",
						"swift_code_bic":"123456789",
						"swift_account_holder_name":"John Doe",
						"swift_account_number_iban":"123456789",
						"transfers_account":"BM123123123123",
						"transfers_type":"CVU",
						"has_virtual_account":true
					},
					{
						"receiver_id":"re_111111111111",
						"id":"pa_111111111111",
						"status":"completed",
						"sender_wallet_address":"0x456...abc",
						"signed_transaction":"BBB...Xyz9z1B",
						"quote_id":"qu_111111111111",
						"instance_id":"in_000000000000",
						"tracking_transaction":{
							"step":"completed",
							"status":"success",
							"transaction_hash":"0x456...abc",
							"completed_at":"2011-10-05T15:48:00.000Z"
						},
						"tracking_payment":{
							"step":"completed",
							"provider_name":"blockchain",
							"provider_transaction_id":"0x456...abc",
							"provider_status":"completed",
							"estimated_time_of_arrival":"5_min",
							"completed_at":"2011-10-05T15:48:00.000Z"
						},
						"tracking_liquidity":{
							"step":"completed",
							"provider_transaction_id":"0x456...abc",
							"provider_status":"completed",
							"estimated_time_of_arrival":"1_business_day",
							"completed_at":"2011-10-05T15:48:00.000Z"
						},
						"tracking_complete":{
							"step":"completed",
							"status":"completed",
							"transaction_hash":"0x456...abc",
							"completed_at":"2011-10-05T15:48:00.000Z"
						},
						"tracking_partner_fee":{
							"step":"completed",
							"transaction_hash":"0x456...abc",
							"completed_at":"2011-10-05T15:48:00.000Z"
						},
						"created_at":"2021-01-02T00:00:00Z",
						"updated_at":"2021-01-02T00:00:00Z",
						"image_url":"https://example.com/image2.png",
						"first_name":"Jane",
						"last_name":"Smith",
						"legal_name":"Another Company LLC.",
						"network":"sepolia",
						"token":"USDC",
						"description":"Another payment description",
						"sender_amount":2020,
						"receiver_amount":10480,
						"partner_fee_amount":300,
						"commercial_quotation":495,
						"blindpay_quotation":485,
						"total_fee_amount":3.0,
						"receiver_local_amount":2000,
						"currency":"BRL",
						"transaction_document_file":"https://example.com/image2.png",
						"transaction_document_type":"receipt",
						"transaction_document_id":"9876543210",
						"name":"Another Bank Account",
						"type":"pix",
						"pix_key":"jane@example.com",
						"has_virtual_account":false
					}
				]`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/export/payouts", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	payouts, err := client.Export(context.Background(), &ExportParams{
		Limit:  100,
		Offset: 0,
	})
	require.NoError(t, err)
	require.Len(t, payouts, 2)
	require.Equal(t, "pa_000000000000", payouts[0].ID)
	require.Equal(t, types.TransactionStatusProcessing, payouts[0].Status)
	require.Equal(t, "pa_111111111111", payouts[1].ID)
	require.Equal(t, types.TransactionStatusCompleted, payouts[1].Status)
}

func TestPayouts_Get(t *testing.T) {
	instanceID := "in_000000000000"
	id := "pa_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
					"receiver_id":"re_000000000000",
					"id":"pa_000000000000",
					"status":"processing",
					"sender_wallet_address":"0x123...890",
					"signed_transaction":"AAA...Zey8y0A",
					"quote_id":"qu_000000000000",
					"instance_id":"in_000000000000",
					"tracking_transaction":{
						"step":"processing",
						"status":"failed",
						"transaction_hash":"0x123...890",
						"completed_at":"2011-10-05T14:48:00.000Z"
					},
					"tracking_payment":{
						"step":"on_hold",
						"provider_name":"blockchain",
						"provider_transaction_id":"0x123...890",
						"provider_status":"canceled",
						"estimated_time_of_arrival":"5_min",
						"completed_at":"2011-10-05T14:48:00.000Z"
					},
					"tracking_liquidity":{
						"step":"processing",
						"provider_transaction_id":"0x123...890",
						"provider_status":"deposited",
						"estimated_time_of_arrival":"1_business_day",
						"completed_at":"2011-10-05T14:48:00.000Z"
					},
					"tracking_complete":{
						"step":"on_hold",
						"status":"tokens_refunded",
						"transaction_hash":"0x123...890",
						"completed_at":"2011-10-05T14:48:00.000Z"
					},
					"tracking_partner_fee":{
						"step":"on_hold",
						"transaction_hash":"0x123...890",
						"completed_at":"2011-10-05T14:48:00.000Z"
					},
					"created_at":"2021-01-01T00:00:00Z",
					"updated_at":"2021-01-01T00:00:00Z",
					"image_url":"https://example.com/image.png",
					"first_name":"John",
					"last_name":"Doe",
					"legal_name":"Company Name Inc.",
					"network":"sepolia",
					"token":"USDC",
					"description":"Memo code or description, only works with USD and BRL",
					"sender_amount":1010,
					"receiver_amount":5240,
					"partner_fee_amount":150,
					"commercial_quotation":495,
					"blindpay_quotation":485,
					"total_fee_amount":1.5,
					"receiver_local_amount":1000,
					"currency":"BRL",
					"transaction_document_file":"https://example.com/image.png",
					"transaction_document_type":"invoice",
					"transaction_document_id":"1234567890",
					"name":"Bank Account Name",
					"type":"wire",
					"pix_key":"14947677768",
					"account_number":"1001001234",
					"routing_number":"012345678",
					"country":"US",
					"account_class":"individual",
					"address_line_1":"Address line 1",
					"address_line_2":"Address line 2",
					"city":"City",
					"state_province_region":"State/Province/Region",
					"postal_code":"Postal code",
					"account_type":"checking",
					"ach_cop_beneficiary_first_name":"Fernando",
					"ach_cop_bank_account":"12345678",
					"ach_cop_bank_code":"051",
					"ach_cop_beneficiary_last_name":"Guzman Alarcón",
					"ach_cop_document_id":"1661105408",
					"ach_cop_document_type":"CC",
					"ach_cop_email":"fernando.guzman@gmail.com",
					"beneficiary_name":"Individual full name or business name",
					"spei_clabe":"5482347403740546",
					"spei_protocol":"clabe",
					"spei_institution_code":"40002",
					"swift_beneficiary_country":"MX",
					"swift_code_bic":"123456789",
					"swift_account_holder_name":"John Doe",
					"swift_account_number_iban":"123456789",
					"transfers_account":"BM123123123123",
					"transfers_type":"CVU",
					"has_virtual_account":true
				}`),
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
	require.Equal(t, types.TransactionStatusProcessing, payout.Status)
}

func TestPayouts_ExecuteEvm(t *testing.T) {
	instanceID := "in_000000000000"
	quoteID := "qu_000000000000"
	payoutID := "pa_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"quote_id":"qu_000000000000",
					"sender_wallet_address":"0x123...890"
				}`),
				Out: json.RawMessage(fmt.Sprintf(`{
					"id":"%s",
					"status":"processing",
					"sender_wallet_address":"0x123...890",
					"tracking_complete":{
						"step":"on_hold",
						"status":"tokens_refunded",
						"transaction_hash":"0x123...890",
						"completed_at":"2011-10-05T14:48:00.000Z"
					},
					"tracking_payment":{
						"step":"on_hold",
						"provider_name":"blockchain",
						"provider_transaction_id":"0x123...890",
						"provider_status":"canceled",
						"estimated_time_of_arrival":"5_min",
						"completed_at":"2011-10-05T14:48:00.000Z"
					},
					"tracking_transaction":{
						"step":"processing",
						"status":"failed",
						"transaction_hash":"0x123...890",
						"completed_at":"2011-10-05T14:48:00.000Z"
					},
					"tracking_partner_fee":{
						"step":"on_hold",
						"transaction_hash":"0x123...890",
						"completed_at":"2011-10-05T14:48:00.000Z"
					},
					"tracking_liquidity":{
						"step":"processing",
						"provider_transaction_id":"0x123...890",
						"provider_status":"deposited",
						"estimated_time_of_arrival":"1_business_day",
						"completed_at":"2011-10-05T14:48:00.000Z"
					},
					"receiver_id":"re_000000000000"
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
		SenderWalletAddress: "0x123...890",
	})
	require.NoError(t, err)
	require.Equal(t, payoutID, payout.ID)
}

func TestPayouts_ExecuteStellar(t *testing.T) {
	instanceID := "in_000000000000"
	quoteID := "qu_000000000000"
	payoutID := "pa_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"quote_id":"qu_000000000000",
					"sender_wallet_address":"0x123...890"
				}`),
				Out: json.RawMessage(fmt.Sprintf(`{
					"id":"%s",
					"status":"processing",
					"sender_wallet_address":"0x123...890",
					"tracking_complete":{
						"step":"on_hold",
						"status":"tokens_refunded",
						"transaction_hash":"0x123...890",
						"completed_at":"2011-10-05T14:48:00.000Z"
					},
					"tracking_payment":{
						"step":"on_hold",
						"provider_name":"blockchain",
						"provider_transaction_id":"0x123...890",
						"provider_status":"canceled",
						"estimated_time_of_arrival":"5_min",
						"completed_at":"2011-10-05T14:48:00.000Z"
					},
					"tracking_transaction":{
						"step":"processing",
						"status":"failed",
						"transaction_hash":"0x123...890",
						"completed_at":"2011-10-05T14:48:00.000Z"
					},
					"tracking_partner_fee":{
						"step":"on_hold",
						"transaction_hash":"0x123...890",
						"completed_at":"2011-10-05T14:48:00.000Z"
					},
					"tracking_liquidity":{
						"step":"processing",
						"provider_transaction_id":"0x123...890",
						"provider_status":"deposited",
						"estimated_time_of_arrival":"1_business_day",
						"completed_at":"2011-10-05T14:48:00.000Z"
					},
					"receiver_id":"re_000000000000"
				}`, payoutID)),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/payouts/stellar", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	payout, err := client.CreateStellar(context.Background(), &CreateStellarParams{
		QuoteID:             quoteID,
		SenderWalletAddress: "0x123...890",
	})
	require.NoError(t, err)
	require.Equal(t, payoutID, payout.ID)
}

func TestPayouts_AuthorizeStellarToken(t *testing.T) {
	instanceID := "in_000000000000"
	quoteID := "qu_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"quote_id":"qu_000000000000",
					"sender_wallet_address":"0x123...890"
				}`),
				Out: json.RawMessage(`{
					"transaction_hash":"string"
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
		SenderWalletAddress: "0x123...890",
	})
	require.NoError(t, err)
	require.Equal(t, "string", response.TransactionHash)
}
