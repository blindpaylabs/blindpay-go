package quotes

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

func TestQuotes_Create(t *testing.T) {
	instanceID := "in_000000000000"
	id := "qu_000000000000"
	bankAccountID := "ba_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"bank_account_id":"ba_000000000000",
					"currency_type":"sender",
					"network":"sepolia",
					"request_amount":1000,
					"token":"USDC",
					"cover_fees":true,
					"description":"Memo code or description, only works with USD and BRL",
					"partner_fee_id":"pf_000000000000"
				}`),
				Out: json.RawMessage(`{
					"id":"qu_000000000000",
					"expires_at":1712958191,
					"commercial_quotation":495,
					"blindpay_quotation":485,
					"receiver_amount":5240,
					"sender_amount":1010,
					"partner_fee_amount":150,
					"flat_fee":50,
					"contract":{
						"abi":[{}],
						"address":"0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238",
						"functionName":"approve",
						"blindpayContractAddress":"0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238",
						"amount":"1000000000000000000",
						"network":{
							"name":"Ethereum",
							"chainId":1
						}
					},
					"receiver_local_amount":1000,
					"description":"Memo code or description, only works with USD and BRL"
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/quotes", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	coverFees := true
	description := "Memo code or description, only works with USD and BRL"
	partnerFeeID := "pf_000000000000"
	quote, err := client.Create(context.Background(), &CreateParams{
		BankAccountID: bankAccountID,
		CurrencyType:  types.CurrencyTypeSender,
		Network:       types.Network("sepolia"),
		RequestAmount: 1000,
		Token:         types.StablecoinToken("USDC"),
		CoverFees:     &coverFees,
		Description:   &description,
		PartnerFeeID:  &partnerFeeID,
	})
	require.NoError(t, err)
	require.Equal(t, id, quote.ID)
	require.Equal(t, 1010.0, *quote.SenderAmount)
	require.Equal(t, 5240.0, *quote.ReceiverAmount)
}

func TestQuotes_GetFxRate(t *testing.T) {
	instanceID := "in_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"currency_type":"sender",
					"from":"USD",
					"to":"BRL",
					"request_amount":1000
				}`),
				Out: json.RawMessage(`{
					"commercial_quotation":495,
					"blindpay_quotation":485,
					"result_amount":1,
					"instance_flat_fee":50,
					"instance_percentage_fee":0
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/quotes/fx", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.GetFxRate(context.Background(), &GetFxRateParams{
		CurrencyType:  types.CurrencyTypeSender,
		From:          types.CurrencyUSD,
		To:            types.CurrencyBRL,
		RequestAmount: 1000,
	})
	require.NoError(t, err)
	require.Equal(t, 495.0, *response.CommercialQuotation)
	require.Equal(t, 1.0, response.ResultAmount)
}
