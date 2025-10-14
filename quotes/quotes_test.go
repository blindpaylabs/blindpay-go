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
	instanceID := "inst_123"
	id := "quote_123"
	bankAccountID := "ba_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"bank_account_id":"ba_123",
					"currency_type":"sender",
					"request_amount":100.50,
					"transaction_document_type":"invoice"
				}`),
				Out: json.RawMessage(fmt.Sprintf(`{
					"id":"%s",
					"expires_at":1704067200,
					"commercial_quotation":1.05,
					"blindpay_quotation":1.04,
					"receiver_amount":96.15,
					"sender_amount":100.50,
					"partner_fee_amount":0,
					"flat_fee":2.50,
					"contract":{"abi":[],"address":"0x123","functionName":"transfer","blindpayContractAddress":"0x456","amount":"100500000","network":{"name":"ethereum","chainId":1}},
					"receiver_local_amount":96.15,
					"description":""
				}`, id)),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/quotes", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	quote, err := client.Create(context.Background(), &CreateParams{
		BankAccountID:           bankAccountID,
		CurrencyType:            types.CurrencyTypeSender,
		RequestAmount:           100.50,
		TransactionDocumentType: types.TransactionDocumentTypeInvoice,
	})
	require.NoError(t, err)
	require.Equal(t, id, quote.ID)
	require.Equal(t, 100.50, quote.SenderAmount)
	require.Equal(t, 96.15, quote.ReceiverAmount)
}

func TestQuotes_GetFxRate(t *testing.T) {
	instanceID := "inst_123"

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
					"request_amount":100.0
				}`),
				Out: json.RawMessage(`{
					"commercial_quotation":5.25,
					"blindpay_quotation":5.20,
					"result_amount":520.0,
					"instance_flat_fee":2.5,
					"instance_percentage_fee":0.02
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
		RequestAmount: 100.0,
	})
	require.NoError(t, err)
	require.Equal(t, 5.25, response.CommercialQuotation)
	require.Equal(t, 520.0, response.ResultAmount)
}
