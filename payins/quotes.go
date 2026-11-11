package payins

import (
	"context"
	"fmt"

	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/blindpaylabs/blindpay-go/option"
)

// CreateQuoteParams represents parameters for creating a payin quote.
type CreateQuoteParams struct {
	// BlockchainWalletID or WalletID (a managed wallet) receives the tokens.
	BlockchainWalletID string                `json:"blockchain_wallet_id,omitempty"`
	CurrencyType       types.CurrencyType    `json:"currency_type"`
	PaymentMethod      string                `json:"payment_method"`
	RequestAmount      float64               `json:"request_amount"`
	Token              types.StablecoinToken `json:"token"`
	CoverFees          bool                  `json:"cover_fees,omitempty"`
	PartnerFeeID       *string               `json:"partner_fee_id,omitempty"`
	PayerRules         PayerRules            `json:"payer_rules"`
	IsOtc              *bool                 `json:"is_otc,omitempty"`
	WalletID           *string               `json:"wallet_id,omitempty"`
	// FundingBankAccountID is the customer's linked bank account to pull
	// from (ACH pull payins).
	FundingBankAccountID *string `json:"funding_bank_account_id,omitempty"`
}

// PayerRules represents payer rules for payin quotes.
type PayerRules struct {
	PixAllowedTaxIDs      []string `json:"pix_allowed_tax_ids"`
	TransfersAllowedTaxID *string  `json:"transfers_allowed_tax_id,omitempty"`
	PseAllowedTaxIDs      []string `json:"pse_allowed_tax_ids,omitempty"`
	PseFullName           *string  `json:"pse_full_name,omitempty"`
	// PseDocumentType is "CC" or "NIT".
	PseDocumentType   *string `json:"pse_document_type,omitempty"`
	PseDocumentNumber *string `json:"pse_document_number,omitempty"`
	PseEmail          *string `json:"pse_email,omitempty"`
	PsePhone          *string `json:"pse_phone,omitempty"`
	PseBankCode       *string `json:"pse_bank_code,omitempty"`
}

// CardInstruction tells a browser checkout how to collect a card payment.
type CardInstruction struct {
	// Provider names the card gateway, so the checkout loads its
	// tokenization library.
	Provider          string  `json:"provider"`
	TokenizationKey   string  `json:"tokenization_key"`
	CheckoutPublicKey *string `json:"checkout_public_key,omitempty"`
	CheckoutToken     *string `json:"checkout_token,omitempty"`
	CheckoutURL       *string `json:"checkout_url,omitempty"`
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
	// CardInstruction is set for card payins: what a browser checkout needs
	// to tokenize the card.
	CardInstruction *CardInstruction `json:"card_instruction,omitempty"`
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
func (c *QuotesClient) Create(ctx context.Context, params *CreateQuoteParams, opts ...option.RequestOption) (*CreateQuoteResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := request.Pathf("/instances/%s/payin-quotes", c.instanceID)
	return request.Do[*CreateQuoteResponse](c.cfg, ctx, "POST", path, params, opts...)
}

// GetFxRate retrieves the current FX rate for a payin currency pair.
func (c *QuotesClient) GetFxRate(ctx context.Context, params *GetFxRateParams, opts ...option.RequestOption) (*GetFxRateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := request.Pathf("/instances/%s/payin-quotes/fx", c.instanceID)
	return request.Do[*GetFxRateResponse](c.cfg, ctx, "POST", path, params, opts...)
}
