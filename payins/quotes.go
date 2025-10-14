package payins

import (
	"context"
	"fmt"

	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
)

// CreateQuoteParams represents parameters for creating a payin quote.
type CreateQuoteParams struct {
	BlockchainWalletID string                `json:"blockchain_wallet_id"`
	CurrencyType       types.CurrencyType    `json:"currency_type"`
	PaymentMethod      string                `json:"payment_method"`
	RequestAmount      float64               `json:"request_amount"`
	Token              types.StablecoinToken `json:"token"`
	CoverFees          bool                  `json:"cover_fees"`
	PartnerFeeID       *string               `json:"partner_fee_id,omitempty"`
	PayerRules         PayerRules            `json:"payer_rules"`
}

// PayerRules represents payer rules for payin quotes.
type PayerRules struct {
	PixAllowedTaxIDs []string `json:"pix_allowed_tax_ids"`
}

// CreateQuoteResponse represents the response when creating a payin quote.
type CreateQuoteResponse struct {
	ID                  string   `json:"id"`
	ExpiresAt           int64    `json:"expires_at"`
	CommercialQuotation float64  `json:"commercial_quotation"`
	BlindpayQuotation   float64  `json:"blindpay_quotation"`
	ReceiverAmount      float64  `json:"receiver_amount"`
	SenderAmount        float64  `json:"sender_amount"`
	PartnerFeeAmount    *float64 `json:"partner_fee_amount,omitempty"`
	FlatFee             *float64 `json:"flat_fee,omitempty"`
}

// GetFxRateParams represents parameters for getting payin FX rates.
type GetFxRateParams struct {
	CurrencyType  types.CurrencyType `json:"currency_type"`
	From          types.Currency     `json:"from"`
	To            types.Currency     `json:"to"`
	RequestAmount float64            `json:"request_amount"`
}

// GetFxRateResponse represents the payin FX rate response.
type GetFxRateResponse struct {
	CommercialQuotation   float64 `json:"commercial_quotation"`
	BlindpayQuotation     float64 `json:"blindpay_quotation"`
	ResultAmount          float64 `json:"result_amount"`
	InstanceFlatFee       float64 `json:"instance_flat_fee"`
	InstancePercentageFee float64 `json:"instance_percentage_fee"`
}

// QuotesClient handles payin quote-related operations.
type QuotesClient struct {
	cfg        *request.Config
	instanceID string
}

// Create creates a new payin quote.
func (c *QuotesClient) Create(ctx context.Context, params *CreateQuoteParams) (*CreateQuoteResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := fmt.Sprintf("/instances/%s/payin-quotes", c.instanceID)
	return request.Do[*CreateQuoteResponse](c.cfg, ctx, "POST", path, params)
}

// GetFxRate retrieves the current FX rate for a payin currency pair.
func (c *QuotesClient) GetFxRate(ctx context.Context, params *GetFxRateParams) (*GetFxRateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := fmt.Sprintf("/instances/%s/payin-quotes/fx", c.instanceID)
	return request.Do[*GetFxRateResponse](c.cfg, ctx, "POST", path, params)
}
