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
	instanceID := "inst_test123"
	walletID := "bw_abc123"

	reqBody := &CreateQuoteParams{
		BlockchainWalletID: walletID,
		CurrencyType:       types.CurrencyTypeSender,
		PaymentMethod:      "pix",
		RequestAmount:      100.00,
		Token:              types.StablecoinTokenUSDC,
		CoverFees:          true,
		PayerRules: PayerRules{
			PixAllowedTaxIDs: []string{"12345678901", "98765432109"},
		},
	}

	respBody := &CreateQuoteResponse{
		ID:                  "pq_quote123",
		ExpiresAt:           1704153600,
		CommercialQuotation: 5.50,
		BlindpayQuotation:   5.45,
		ReceiverAmount:      545.00,
		SenderAmount:        100.00,
	}

	reqJSON, _ := json.Marshal(reqBody)
	respJSON, _ := json.Marshal(respBody)

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
		cfg: &request.Config{
			BaseURL:    cfg.BaseURL,
			APIKey:     cfg.APIKey,
			HTTPClient: cfg.HTTPClient,
			UserAgent:  cfg.UserAgent,
		},
		instanceID: instanceID,
	}

	quote, err := client.Create(context.Background(), reqBody)
	require.NoError(t, err)
	require.NotNil(t, quote)
	require.Equal(t, "pq_quote123", quote.ID)
	require.Equal(t, int64(1704153600), quote.ExpiresAt)
	require.Equal(t, 5.50, quote.CommercialQuotation)
	require.Equal(t, 5.45, quote.BlindpayQuotation)
	require.Equal(t, 545.00, quote.ReceiverAmount)
	require.Equal(t, 100.00, quote.SenderAmount)
}

func TestQuotesClient_Create_WithPartnerFee(t *testing.T) {
	instanceID := "inst_test123"
	walletID := "bw_abc123"
	partnerFeeID := "pf_fee123"
	partnerFeeAmount := 5.00
	flatFee := 2.50

	reqBody := &CreateQuoteParams{
		BlockchainWalletID: walletID,
		CurrencyType:       types.CurrencyTypeReceiver,
		PaymentMethod:      "pix",
		RequestAmount:      500.00,
		Token:              types.StablecoinTokenUSDT,
		CoverFees:          false,
		PartnerFeeID:       &partnerFeeID,
		PayerRules: PayerRules{
			PixAllowedTaxIDs: []string{"12345678901"},
		},
	}

	respBody := &CreateQuoteResponse{
		ID:                  "pq_quote456",
		ExpiresAt:           1704153700,
		CommercialQuotation: 5.50,
		BlindpayQuotation:   5.48,
		ReceiverAmount:      500.00,
		SenderAmount:        91.24,
		PartnerFeeAmount:    &partnerFeeAmount,
		FlatFee:             &flatFee,
	}

	reqJSON, _ := json.Marshal(reqBody)
	respJSON, _ := json.Marshal(respBody)

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
		cfg: &request.Config{
			BaseURL:    cfg.BaseURL,
			APIKey:     cfg.APIKey,
			HTTPClient: cfg.HTTPClient,
			UserAgent:  cfg.UserAgent,
		},
		instanceID: instanceID,
	}

	quote, err := client.Create(context.Background(), reqBody)
	require.NoError(t, err)
	require.NotNil(t, quote)
	require.Equal(t, "pq_quote456", quote.ID)
	require.NotNil(t, quote.PartnerFeeAmount)
	require.Equal(t, 5.00, *quote.PartnerFeeAmount)
	require.NotNil(t, quote.FlatFee)
	require.Equal(t, 2.50, *quote.FlatFee)
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
	instanceID := "inst_test123"

	reqBody := &GetFxRateParams{
		CurrencyType:  types.CurrencyTypeSender,
		From:          types.CurrencyUSD,
		To:            types.CurrencyBRL,
		RequestAmount: 100.00,
	}

	respBody := &GetFxRateResponse{
		CommercialQuotation:   5.50,
		BlindpayQuotation:     5.45,
		ResultAmount:          545.00,
		InstanceFlatFee:       2.00,
		InstancePercentageFee: 0.5,
	}

	reqJSON, _ := json.Marshal(reqBody)
	respJSON, _ := json.Marshal(respBody)

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
		cfg: &request.Config{
			BaseURL:    cfg.BaseURL,
			APIKey:     cfg.APIKey,
			HTTPClient: cfg.HTTPClient,
			UserAgent:  cfg.UserAgent,
		},
		instanceID: instanceID,
	}

	fxRate, err := client.GetFxRate(context.Background(), reqBody)
	require.NoError(t, err)
	require.NotNil(t, fxRate)
	require.Equal(t, 5.50, fxRate.CommercialQuotation)
	require.Equal(t, 5.45, fxRate.BlindpayQuotation)
	require.Equal(t, 545.00, fxRate.ResultAmount)
	require.Equal(t, 2.00, fxRate.InstanceFlatFee)
	require.Equal(t, 0.5, fxRate.InstancePercentageFee)
}

func TestQuotesClient_GetFxRate_ReceiverCurrency(t *testing.T) {
	instanceID := "inst_test123"

	reqBody := &GetFxRateParams{
		CurrencyType:  types.CurrencyTypeReceiver,
		From:          types.CurrencyBRL,
		To:            types.CurrencyUSD,
		RequestAmount: 500.00,
	}

	respBody := &GetFxRateResponse{
		CommercialQuotation:   5.50,
		BlindpayQuotation:     5.48,
		ResultAmount:          91.24,
		InstanceFlatFee:       1.50,
		InstancePercentageFee: 0.3,
	}

	reqJSON, _ := json.Marshal(reqBody)
	respJSON, _ := json.Marshal(respBody)

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
		cfg: &request.Config{
			BaseURL:    cfg.BaseURL,
			APIKey:     cfg.APIKey,
			HTTPClient: cfg.HTTPClient,
			UserAgent:  cfg.UserAgent,
		},
		instanceID: instanceID,
	}

	fxRate, err := client.GetFxRate(context.Background(), reqBody)
	require.NoError(t, err)
	require.NotNil(t, fxRate)
	require.Equal(t, 5.50, fxRate.CommercialQuotation)
	require.Equal(t, 5.48, fxRate.BlindpayQuotation)
	require.Equal(t, 91.24, fxRate.ResultAmount)
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
