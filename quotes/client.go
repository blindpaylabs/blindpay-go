package quotes

import (
	"context"
	"fmt"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
)

// CreateParams represents parameters for creating a quote.
type CreateParams struct {
	BankAccountID           string                        `json:"bank_account_id"`
	CurrencyType            types.CurrencyType            `json:"currency_type"`
	Network                 *types.Network                `json:"network,omitempty"`
	RequestAmount           float64                       `json:"request_amount"`
	Token                   *types.StablecoinToken        `json:"token,omitempty"`
	CoverFees               *bool                         `json:"cover_fees,omitempty"`
	Description             *string                       `json:"description,omitempty"`
	PartnerFeeID            *string                       `json:"partner_fee_id,omitempty"`
	TransactionDocumentFile *string                       `json:"transaction_document_file"`
	TransactionDocumentID   *string                       `json:"transaction_document_id"`
	TransactionDocumentType types.TransactionDocumentType `json:"transaction_document_type"`
}

// QuoteContract represents contract information in a quote.
type QuoteContract struct {
	ABI                     []map[string]any `json:"abi"`
	Address                 string           `json:"address"`
	FunctionName            string           `json:"functionName"`
	BlindpayContractAddress string           `json:"blindpayContractAddress"`
	Amount                  string           `json:"amount"`
	Network                 struct {
		Name    string `json:"name"`
		ChainID int    `json:"chainId"`
	} `json:"network"`
}

// CreateResponse represents the response when creating a quote.
type CreateResponse struct {
	ID                  string        `json:"id"`
	ExpiresAt           int64         `json:"expires_at"`
	CommercialQuotation float64       `json:"commercial_quotation"`
	BlindpayQuotation   float64       `json:"blindpay_quotation"`
	ReceiverAmount      float64       `json:"receiver_amount"`
	SenderAmount        float64       `json:"sender_amount"`
	PartnerFeeAmount    float64       `json:"partner_fee_amount"`
	FlatFee             float64       `json:"flat_fee"`
	Contract            QuoteContract `json:"contract"`
	ReceiverLocalAmount float64       `json:"receiver_local_amount"`
	Description         string        `json:"description"`
}

// GetFxRateParams represents parameters for getting FX rates.
type GetFxRateParams struct {
	CurrencyType  types.CurrencyType `json:"currency_type"`
	From          types.Currency     `json:"from"`
	To            types.Currency     `json:"to"`
	RequestAmount float64            `json:"request_amount"`
}

// GetFxRateResponse represents the FX rate response.
type GetFxRateResponse struct {
	CommercialQuotation   float64 `json:"commercial_quotation"`
	BlindpayQuotation     float64 `json:"blindpay_quotation"`
	ResultAmount          float64 `json:"result_amount"`
	InstanceFlatFee       float64 `json:"instance_flat_fee"`
	InstancePercentageFee float64 `json:"instance_percentage_fee"`
}

// Client handles quote-related operations.
type Client struct {
	cfg        *request.Config
	instanceID string
}

// NewClient creates a new quotes client.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:        cfg.ToRequestConfig(),
		instanceID: cfg.InstanceID,
	}
}

// Create creates a new quote for a payout.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*CreateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := fmt.Sprintf("/instances/%s/quotes", c.instanceID)
	return request.Do[*CreateResponse](c.cfg, ctx, "POST", path, params)
}

// GetFxRate retrieves the current FX rate for a currency pair.
func (c *Client) GetFxRate(ctx context.Context, params *GetFxRateParams) (*GetFxRateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := fmt.Sprintf("/instances/%s/quotes/fx", c.instanceID)
	return request.Do[*GetFxRateResponse](c.cfg, ctx, "POST", path, params)
}
