package payouts

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/blindpaylabs/blindpay-go/bankaccounts"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
)

// Payout represents a payout transaction.
type Payout struct {
	ReceiverID                 string                          `json:"receiver_id"`
	ID                         string                          `json:"id"`
	Status                     types.TransactionStatus         `json:"status"`
	SenderWalletAddress        string                          `json:"sender_wallet_address"`
	SignedTransaction          string                          `json:"signed_transaction"`
	QuoteID                    string                          `json:"quote_id"`
	InstanceID                 string                          `json:"instance_id"`
	TrackingTransaction        *types.TrackingTransaction      `json:"tracking_transaction"`
	TrackingPayment            *types.TrackingPayment          `json:"tracking_payment"`
	TrackingLiquidity          *types.TrackingLiquidity        `json:"tracking_liquidity"`
	TrackingComplete           *types.TrackingComplete         `json:"tracking_complete"`
	TrackingPartnerFee         *types.TrackingPartnerFee       `json:"tracking_partner_fee"`
	CreatedAt                  time.Time                       `json:"created_at"`
	UpdatedAt                  time.Time                       `json:"updated_at"`
	ImageURL                   string                          `json:"image_url"`
	FirstName                  string                          `json:"first_name"`
	LastName                   string                          `json:"last_name"`
	LegalName                  string                          `json:"legal_name"`
	Network                    types.Network                   `json:"network"`
	Token                      types.StablecoinToken           `json:"token"`
	Description                string                          `json:"description"`
	SenderAmount               float64                         `json:"sender_amount"`
	ReceiverAmount             float64                         `json:"receiver_amount"`
	PartnerFeeAmount           float64                         `json:"partner_fee_amount"`
	CommercialQuotation        float64                         `json:"commercial_quotation"`
	BlindpayQuotation          float64                         `json:"blindpay_quotation"`
	TotalFeeAmount             float64                         `json:"total_fee_amount"`
	ReceiverLocalAmount        float64                         `json:"receiver_local_amount"`
	Currency                   types.Currency                  `json:"currency"`
	TransactionDocumentFile    string                          `json:"transaction_document_file"`
	TransactionDocumentType    types.TransactionDocumentType   `json:"transaction_document_type"`
	TransactionDocumentID      string                          `json:"transaction_document_id"`
	Name                       string                          `json:"name"`
	Type                       types.Rail                      `json:"type"`
	PixKey                     string                          `json:"pix_key,omitempty"`
	AccountNumber              string                          `json:"account_number,omitempty"`
	RoutingNumber              string                          `json:"routing_number,omitempty"`
	Country                    types.Country                   `json:"country,omitempty"`
	AccountClass               types.AccountClass              `json:"account_class,omitempty"`
	AddressLine1               string                          `json:"address_line_1,omitempty"`
	AddressLine2               string                          `json:"address_line_2,omitempty"`
	City                       string                          `json:"city,omitempty"`
	StateProvinceRegion        string                          `json:"state_province_region,omitempty"`
	PostalCode                 string                          `json:"postal_code,omitempty"`
	AccountType                types.BankAccountType           `json:"account_type,omitempty"`
	AchCopBeneficiaryFirstName string                          `json:"ach_cop_beneficiary_first_name,omitempty"`
	AchCopBankAccount          string                          `json:"ach_cop_bank_account,omitempty"`
	AchCopBankCode             string                          `json:"ach_cop_bank_code,omitempty"`
	AchCopBeneficiaryLastName  string                          `json:"ach_cop_beneficiary_last_name,omitempty"`
	AchCopDocumentID           string                          `json:"ach_cop_document_id,omitempty"`
	AchCopDocumentType         string                          `json:"ach_cop_document_type,omitempty"`
	AchCopEmail                string                          `json:"ach_cop_email,omitempty"`
	BeneficiaryName            string                          `json:"beneficiary_name,omitempty"`
	SpeiClabe                  string                          `json:"spei_clabe,omitempty"`
	SpeiProtocol               bankaccounts.SpeiProtocol       `json:"spei_protocol,omitempty"`
	SpeiInstitutionCode        string                          `json:"spei_institution_code,omitempty"`
	SwiftBeneficiaryCountry    types.Country                   `json:"swift_beneficiary_country,omitempty"`
	SwiftCodeBic               string                          `json:"swift_code_bic,omitempty"`
	SwiftAccountHolderName     string                          `json:"swift_account_holder_name,omitempty"`
	SwiftAccountNumberIban     string                          `json:"swift_account_number_iban,omitempty"`
	TransfersAccount           string                          `json:"transfers_account,omitempty"`
	TransfersType              bankaccounts.ArgentinaTransfers `json:"transfers_type,omitempty"`
	HasVirtualAccount          bool                            `json:"has_virtual_account"`
}

// ListParams represents parameters for listing payouts.
type ListParams struct {
	ReceiverID string `json:"receiver_id,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Offset     int    `json:"offset,omitempty"`
}

// ListResponse represents the response when listing payouts.
type ListResponse struct {
	Data       []Payout                 `json:"data"`
	Pagination types.PaginationMetadata `json:"pagination"`
}

// ExportParams represents parameters for exporting payouts.
type ExportParams struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// AuthorizeStellarTokenParams represents parameters for authorizing a Stellar token.
type AuthorizeStellarTokenParams struct {
	QuoteID             string `json:"quote_id"`
	SenderWalletAddress string `json:"sender_wallet_address"`
}

// AuthorizeStellarTokenResponse represents the response when authorizing a Stellar token.
type AuthorizeStellarTokenResponse struct {
	TransactionHash string `json:"transaction_hash"`
}

// CreateStellarParams represents parameters for creating a Stellar payout.
type CreateStellarParams struct {
	QuoteID             string  `json:"quote_id"`
	SenderWalletAddress string  `json:"sender_wallet_address"`
	SignedTransaction   *string `json:"signed_transaction,omitempty"`
}

// CreateEvmParams represents parameters for creating an EVM payout.
type CreateEvmParams struct {
	QuoteID             string `json:"quote_id"`
	SenderWalletAddress string `json:"sender_wallet_address"`
}

// CreateResponse represents the response when creating a payout.
type CreateResponse struct {
	ID                  string                     `json:"id"`
	Status              types.TransactionStatus    `json:"status"`
	TrackingTransaction *types.TrackingTransaction `json:"tracking_transaction,omitempty"`
	TrackingPayment     *types.TrackingPayment     `json:"tracking_payment,omitempty"`
	TrackingLiquidity   *types.TrackingLiquidity   `json:"tracking_liquidity,omitempty"`
	TrackingComplete    *types.TrackingComplete    `json:"tracking_complete,omitempty"`
	TrackingPartnerFee  *types.TrackingPartnerFee  `json:"tracking_partner_fee,omitempty"`
}

// Client handles payout-related operations.
type Client struct {
	cfg        *request.Config
	instanceID string
}

// NewClient creates a new payouts client.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:        cfg.ToRequestConfig(),
		instanceID: cfg.InstanceID,
	}
}

// List retrieves all payouts with optional filters.
func (c *Client) List(ctx context.Context, params *ListParams) (*ListResponse, error) {
	path := fmt.Sprintf("/instances/%s/payouts", c.instanceID)

	if params != nil {
		q := url.Values{}
		if params.ReceiverID != "" {
			q.Set("receiver_id", params.ReceiverID)
		}
		if params.Limit > 0 {
			q.Set("limit", fmt.Sprintf("%d", params.Limit))
		}
		if params.Offset > 0 {
			q.Set("offset", fmt.Sprintf("%d", params.Offset))
		}
		if len(q) > 0 {
			path += "?" + q.Encode()
		}
	}

	return request.Do[*ListResponse](c.cfg, ctx, "GET", path, nil)
}

// Export retrieves all payouts for export with optional pagination.
func (c *Client) Export(ctx context.Context, params *ExportParams) ([]Payout, error) {
	path := fmt.Sprintf("/instances/%s/export/payouts", c.instanceID)

	if params != nil {
		q := url.Values{}
		if params.Limit > 0 {
			q.Set("limit", fmt.Sprintf("%d", params.Limit))
		}
		if params.Offset > 0 {
			q.Set("offset", fmt.Sprintf("%d", params.Offset))
		}
		if len(q) > 0 {
			path += "?" + q.Encode()
		}
	}

	return request.Do[[]Payout](c.cfg, ctx, "GET", path, nil)
}

// Get retrieves a specific payout by ID.
func (c *Client) Get(ctx context.Context, payoutID string) (*Payout, error) {
	if payoutID == "" {
		return nil, fmt.Errorf("payout ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/payouts/%s", c.instanceID, payoutID)
	return request.Do[*Payout](c.cfg, ctx, "GET", path, nil)
}

// GetTrack retrieves tracking information for a payout.
func (c *Client) GetTrack(ctx context.Context, payoutID string) (*Payout, error) {
	if payoutID == "" {
		return nil, fmt.Errorf("payout ID cannot be empty")
	}

	path := fmt.Sprintf("/e/payouts/%s", payoutID)
	return request.Do[*Payout](c.cfg, ctx, "GET", path, nil)
}

// CreateEvm creates an EVM payout.
func (c *Client) CreateEvm(ctx context.Context, params *CreateEvmParams) (*CreateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := fmt.Sprintf("/instances/%s/payouts/evm", c.instanceID)
	return request.Do[*CreateResponse](c.cfg, ctx, "POST", path, params)
}

// CreateStellar creates a Stellar payout.
func (c *Client) CreateStellar(ctx context.Context, params *CreateStellarParams) (*CreateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := fmt.Sprintf("/instances/%s/payouts/stellar", c.instanceID)
	return request.Do[*CreateResponse](c.cfg, ctx, "POST", path, params)
}

// AuthorizeStellarToken authorizes a Stellar token for payout.
func (c *Client) AuthorizeStellarToken(ctx context.Context, params *AuthorizeStellarTokenParams) (*AuthorizeStellarTokenResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := fmt.Sprintf("/instances/%s/payouts/stellar/authorize-token", c.instanceID)
	return request.Do[*AuthorizeStellarTokenResponse](c.cfg, ctx, "POST", path, params)
}
