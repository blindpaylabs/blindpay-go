package payins

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestQuotesClient_Create(t *testing.T) {
	instanceID := "in_000000000000"
	walletID := "bw_000000000000"
	partnerFeeID := "pf_000000000000"

	reqJSON := json.RawMessage(`{
		"blockchain_wallet_id": "bw_000000000000",
		"currency_type": "sender",
		"cover_fees": true,
		"request_amount": 1000,
		"payment_method": "pix",
		"token": "USDC",
		"partner_fee_id": "pf_000000000000",
		"payer_rules": {
			"pix_allowed_tax_ids": ["149.476.037-68"]
		}
	}`)

	respJSON := json.RawMessage(`{
		"id": "qu_000000000000",
		"expires_at": 1712958191,
		"commercial_quotation": 495,
		"blindpay_quotation": 505,
		"receiver_amount": 1010,
		"sender_amount": 5240,
		"partner_fee_amount": 150,
		"flat_fee": 50
	}`)

	rt := &blindpaytest.RoundTripper{
		T:      t,
		Method: http.MethodPost,
		Path:   "/v1/instances/" + instanceID + "/payin-quotes",
		Status: http.StatusOK,
		In:     reqJSON,
		Out:    respJSON,
	}

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com/v1",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{Transport: rt},
		UserAgent:  "test",
	}

	client := &QuotesClient{
		cfg:        cfg.ToRequestConfig(),
		instanceID: instanceID,
	}

	params := &CreateQuoteParams{
		BlockchainWalletID: walletID,
		CurrencyType:       types.CurrencyTypeSender,
		CoverFees:          true,
		RequestAmount:      1000,
		PaymentMethod:      "pix",
		Token:              types.StablecoinTokenUSDC,
		PartnerFeeID:       &partnerFeeID,
		PayerRules: PayerRules{
			PixAllowedTaxIDs: []string{"149.476.037-68"},
		},
	}

	quote, err := client.Create(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, quote)
	require.Equal(t, "qu_000000000000", quote.ID)
	require.Equal(t, int64(1712958191), quote.ExpiresAt)
	require.Equal(t, float64(495), quote.CommercialQuotation)
	require.Equal(t, float64(505), quote.BlindpayQuotation)
	require.Equal(t, float64(1010), quote.ReceiverAmount)
	require.Equal(t, float64(5240), quote.SenderAmount)
	require.NotNil(t, quote.PartnerFeeAmount)
	require.Equal(t, float64(150), *quote.PartnerFeeAmount)
	require.NotNil(t, quote.FlatFee)
	require.Equal(t, float64(50), *quote.FlatFee)
}

func TestQuotesClient_Create_WithPartnerFee(t *testing.T) {
	instanceID := "in_000000000000"
	walletID := "bw_000000000000"
	partnerFeeID := "pf_000000000000"

	reqJSON := json.RawMessage(`{
		"blockchain_wallet_id": "bw_000000000000",
		"currency_type": "sender",
		"cover_fees": true,
		"request_amount": 1000,
		"payment_method": "pix",
		"token": "USDC",
		"partner_fee_id": "pf_000000000000",
		"payer_rules": {
			"pix_allowed_tax_ids": ["149.476.037-68"]
		}
	}`)

	respJSON := json.RawMessage(`{
		"id": "qu_000000000000",
		"expires_at": 1712958191,
		"commercial_quotation": 495,
		"blindpay_quotation": 505,
		"receiver_amount": 1010,
		"sender_amount": 5240,
		"partner_fee_amount": 150,
		"flat_fee": 50
	}`)

	rt := &blindpaytest.RoundTripper{
		T:      t,
		Method: http.MethodPost,
		Path:   "/v1/instances/" + instanceID + "/payin-quotes",
		Status: http.StatusOK,
		In:     reqJSON,
		Out:    respJSON,
	}

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com/v1",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{Transport: rt},
		UserAgent:  "test",
	}

	client := &QuotesClient{
		cfg:        cfg.ToRequestConfig(),
		instanceID: instanceID,
	}

	params := &CreateQuoteParams{
		BlockchainWalletID: walletID,
		CurrencyType:       types.CurrencyTypeSender,
		CoverFees:          true,
		RequestAmount:      1000,
		PaymentMethod:      "pix",
		Token:              types.StablecoinTokenUSDC,
		PartnerFeeID:       &partnerFeeID,
		PayerRules: PayerRules{
			PixAllowedTaxIDs: []string{"149.476.037-68"},
		},
	}

	quote, err := client.Create(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, quote)
	require.Equal(t, "qu_000000000000", quote.ID)
	require.NotNil(t, quote.PartnerFeeAmount)
	require.Equal(t, float64(150), *quote.PartnerFeeAmount)
	require.NotNil(t, quote.FlatFee)
	require.Equal(t, float64(50), *quote.FlatFee)
}

func TestQuotesClient_Create_NilParams(t *testing.T) {
	client := &QuotesClient{
		cfg:        &request.Config{},
		instanceID: "inst_test",
	}

	quote, err := client.Create(context.Background(), nil)
	require.Error(t, err)
	require.Nil(t, quote)
	require.Contains(t, err.Error(), "params cannot be nil")
}

func TestQuotesClient_GetFxRate(t *testing.T) {
	instanceID := "in_000000000000"

	reqJSON := json.RawMessage(`{
		"currency_type": "sender",
		"from": "USD",
		"to": "BRL",
		"request_amount": 1000
	}`)

	respJSON := json.RawMessage(`{
		"commercial_quotation": 495,
		"blindpay_quotation": 505,
		"result_amount": 1,
		"instance_flat_fee": 50,
		"instance_percentage_fee": 0
	}`)

	rt := &blindpaytest.RoundTripper{
		T:      t,
		Method: http.MethodPost,
		Path:   "/v1/instances/" + instanceID + "/payin-quotes/fx",
		Status: http.StatusOK,
		In:     reqJSON,
		Out:    respJSON,
	}

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com/v1",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{Transport: rt},
		UserAgent:  "test",
	}

	client := &QuotesClient{
		cfg:        cfg.ToRequestConfig(),
		instanceID: instanceID,
	}

	params := &GetFxRateParams{
		CurrencyType:  types.CurrencyTypeSender,
		From:          types.CurrencyUSD,
		To:            types.CurrencyBRL,
		RequestAmount: 1000,
	}

	fxRate, err := client.GetFxRate(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, fxRate)
	require.Equal(t, float64(495), fxRate.CommercialQuotation)
	require.Equal(t, float64(505), fxRate.BlindpayQuotation)
	require.Equal(t, float64(1), fxRate.ResultAmount)
	require.Equal(t, float64(50), fxRate.InstanceFlatFee)
	require.Equal(t, float64(0), fxRate.InstancePercentageFee)
}

func TestQuotesClient_GetFxRate_ReceiverCurrency(t *testing.T) {

	instanceID := "in_000000000000"

	reqJSON := json.RawMessage(`{
		"currency_type": "sender",
		"from": "USD",
		"to": "BRL",
		"request_amount": 1000
	}`)

	respJSON := json.RawMessage(`{
		"commercial_quotation": 495,
		"blindpay_quotation": 505,
		"result_amount": 1,
		"instance_flat_fee": 50,
		"instance_percentage_fee": 0
	}`)

	rt := &blindpaytest.RoundTripper{
		T:      t,
		Method: http.MethodPost,
		Path:   "/v1/instances/" + instanceID + "/payin-quotes/fx",
		Status: http.StatusOK,
		In:     reqJSON,
		Out:    respJSON,
	}

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com/v1",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{Transport: rt},
		UserAgent:  "test",
	}

	client := &QuotesClient{
		cfg:        cfg.ToRequestConfig(),
		instanceID: instanceID,
	}

	params := &GetFxRateParams{
		CurrencyType:  types.CurrencyTypeSender,
		From:          types.CurrencyUSD,
		To:            types.CurrencyBRL,
		RequestAmount: 1000,
	}

	fxRate, err := client.GetFxRate(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, fxRate)
	require.Equal(t, float64(495), fxRate.CommercialQuotation)
	require.Equal(t, float64(505), fxRate.BlindpayQuotation)
	require.Equal(t, float64(1), fxRate.ResultAmount)
}

func TestQuotesClient_GetFxRate_NilParams(t *testing.T) {
	client := &QuotesClient{
		cfg:        &request.Config{},
		instanceID: "inst_test",
	}

	fxRate, err := client.GetFxRate(context.Background(), nil)
	require.Error(t, err)
	require.Nil(t, fxRate)
	require.Contains(t, err.Error(), "params cannot be nil")
}
