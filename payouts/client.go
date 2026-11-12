// Package payouts handles payouts (stablecoin in, fiat out) on EVM chains,
// Solana and Stellar, from a customer wallet (SenderWalletAddress) or a
// BlindPay-managed wallet (WalletID): creation from a quote, listing with
// auto-paging, tracking and supporting documents. Amounts are integer cents
// (1000 = 10.00). Use it through blindpay.Client.Payouts.
package payouts

import (
	"context"
	"fmt"
	"time"

	"github.com/blindpaylabs/blindpay-go/bankaccounts"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/blindpaylabs/blindpay-go/option"
	"github.com/blindpaylabs/blindpay-go/pagination"
)

// Payout represents a payout transaction.
type Payout struct {
	CustomerID          string                     `json:"customer_id"`
	ID                  string                     `json:"id"`
	Status              types.TransactionStatus    `json:"status"`
	SenderWalletAddress string                     `json:"sender_wallet_address"`
	SignedTransaction   string                     `json:"signed_transaction"`
	QuoteID             string                     `json:"quote_id"`
	InstanceID          string                     `json:"instance_id"`
	TrackingTransaction *types.TrackingTransaction `json:"tracking_transaction"`
	TrackingPayment     *types.TrackingPayment     `json:"tracking_payment"`
	TrackingLiquidity   *types.TrackingLiquidity   `json:"tracking_liquidity"`
	TrackingComplete    *types.TrackingComplete    `json:"tracking_complete"`
	TrackingPartnerFee  *types.TrackingPartnerFee  `json:"tracking_partner_fee"`
	CreatedAt           time.Time                  `json:"created_at"`
	UpdatedAt           time.Time                  `json:"updated_at"`
	ImageURL            string                     `json:"image_url"`
	FirstName           string                     `json:"first_name"`
	LastName            string                     `json:"last_name"`
	LegalName           string                     `json:"legal_name"`
	Network             types.Network              `json:"network"`
	Token               types.StablecoinToken      `json:"token"`
	Description         string                     `json:"description"`
	// SenderAmount is in integer cents (1000 = 10.00).
	SenderAmount float64 `json:"sender_amount"`
	// ReceiverAmount is in integer cents (1000 = 10.00).
	ReceiverAmount float64 `json:"receiver_amount"`
	// PartnerFeeAmount is in integer cents (1000 = 10.00).
	PartnerFeeAmount    float64 `json:"partner_fee_amount"`
	CommercialQuotation float64 `json:"commercial_quotation"`
	BlindpayQuotation   float64 `json:"blindpay_quotation"`
	TotalFeeAmount      float64 `json:"total_fee_amount"`
	// CustomerLocalAmount is in integer cents (1000 = 10.00).
	CustomerLocalAmount        float64                         `json:"customer_local_amount"`
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
	PartnerFeeID               *string                         `json:"partner_fee_id,omitempty"`
	TedBankCode                *string                         `json:"ted_bank_code,omitempty"`
	TedBranchCode              *string                         `json:"ted_branch_code,omitempty"`
	TedCpfCnpj                 *string                         `json:"ted_cpf_cnpj,omitempty"`
	// BillingFeeAmount is in integer cents (1000 = 10.00).
	BillingFeeAmount                      *float64                 `json:"billing_fee_amount,omitempty"`
	CpnPaymentID                          *string                  `json:"cpn_payment_id,omitempty"`
	SenderLegalName                       *string                  `json:"sender_legal_name,omitempty"`
	TrackingDocuments                     *types.TrackingDocuments `json:"tracking_documents,omitempty"`
	BankAccountID                         *string                  `json:"bank_account_id,omitempty"`
	FasterPaymentsAccountNumber           *string                  `json:"faster_payments_account_number,omitempty"`
	FasterPaymentsBeneficiaryAddressLine1 *string                  `json:"faster_payments_beneficiary_address_line_1,omitempty"`
	FasterPaymentsBeneficiaryCity         *string                  `json:"faster_payments_beneficiary_city,omitempty"`
	FasterPaymentsBeneficiaryCountry      *string                  `json:"faster_payments_beneficiary_country,omitempty"`
	FasterPaymentsBeneficiaryLegalName    *string                  `json:"faster_payments_beneficiary_legal_name,omitempty"`
	FasterPaymentsBeneficiaryPostalCode   *string                  `json:"faster_payments_beneficiary_postal_code,omitempty"`
	FasterPaymentsSortCode                *string                  `json:"faster_payments_sort_code,omitempty"`
	PartnerFee                            *int                     `json:"partner_fee,omitempty"`
	PayableID                             *string                  `json:"payable_id,omitempty"`
	PixSafeBankCode                       *string                  `json:"pix_safe_bank_code,omitempty"`
	PixSafeBranchCode                     *string                  `json:"pix_safe_branch_code,omitempty"`
	PixSafeCpfCnpj                        *string                  `json:"pix_safe_cpf_cnpj,omitempty"`
	SepaBeneficiaryAddressLine1           *string                  `json:"sepa_beneficiary_address_line_1,omitempty"`
	SepaBeneficiaryCity                   *string                  `json:"sepa_beneficiary_city,omitempty"`
	SepaBeneficiaryCountry                *string                  `json:"sepa_beneficiary_country,omitempty"`
	SepaBeneficiaryLegalName              *string                  `json:"sepa_beneficiary_legal_name,omitempty"`
	SepaBeneficiaryPostalCode             *string                  `json:"sepa_beneficiary_postal_code,omitempty"`
	SepaIban                              *string                  `json:"sepa_iban,omitempty"`
	// TransactionFeeAmount is in integer cents (1000 = 10.00).
	TransactionFeeAmount *float64 `json:"transaction_fee_amount,omitempty"`
}

// ListParams represents parameters for listing payouts.
type ListParams struct {
	// Limit is the page size (10, 50, 100, 200, 500 or 1000). Setting Limit,
	// StartingAfter or EndingBefore makes the API return one page; with none
	// of them it returns every payout.
	Limit int `json:"limit,omitempty"`
	// Deprecated: ignored by the API; use StartingAfter.
	Offset int `json:"offset,omitempty"`
	// StartingAfter returns the page after this payout ID (the previous
	// page's Pagination.NextPage).
	StartingAfter string `json:"-"`
	// EndingBefore returns the page before this payout ID.
	EndingBefore  string                  `json:"-"`
	CustomerID    string                  `json:"customer_id,omitempty"`
	Status        types.TransactionStatus `json:"status,omitempty"`
	CustomerName  string                  `json:"-"`
	BankAccountID string                  `json:"-"`
	Country       string                  `json:"-"`
	PaymentMethod types.Rail              `json:"-"`
	Network       types.Network           `json:"-"`
	Token         types.StablecoinToken   `json:"-"`
}

// ExportParams represents parameters for exporting payouts.
//
// Deprecated: only used by Export, which always fails. Use ListParams.
type ExportParams struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// AuthorizeStellarTokenParams represents parameters for authorizing a Stellar token.
type AuthorizeStellarTokenParams struct {
	QuoteID string `json:"quote_id"`
	// SenderWalletAddress and WalletID: set exactly one. WalletID pays from
	// a BlindPay managed wallet; SenderWalletAddress from the customer's own
	// wallet. The API rejects both or neither (400).
	SenderWalletAddress string `json:"sender_wallet_address,omitempty"`
	WalletID            string `json:"wallet_id,omitempty"`
}

// AuthorizeStellarTokenResponse represents the response when authorizing a Stellar token.
type AuthorizeStellarTokenResponse struct {
	TransactionHash string `json:"transaction_hash"`
}

// CreateStellarParams represents parameters for creating a Stellar payout.
type CreateStellarParams struct {
	QuoteID string `json:"quote_id"`
	// SenderWalletAddress and WalletID: set exactly one. WalletID pays from
	// a BlindPay managed wallet; SenderWalletAddress from the customer's own
	// wallet. The API rejects both or neither (400).
	SenderWalletAddress string  `json:"sender_wallet_address,omitempty"`
	WalletID            string  `json:"wallet_id,omitempty"`
	SignedTransaction   *string `json:"signed_transaction,omitempty"`
}

// CreateSolanaParams represents parameters for creating a Solana payout.
type CreateSolanaParams struct {
	QuoteID string `json:"quote_id"`
	// SenderWalletAddress and WalletID: set exactly one. WalletID pays from
	// a BlindPay managed wallet; SenderWalletAddress from the customer's own
	// wallet. The API rejects both or neither (400).
	SenderWalletAddress string  `json:"sender_wallet_address,omitempty"`
	WalletID            string  `json:"wallet_id,omitempty"`
	SignedTransaction   *string `json:"signed_transaction"`
}

// CreateEvmParams represents parameters for creating an EVM payout.
type CreateEvmParams struct {
	QuoteID string `json:"quote_id"`
	// SenderWalletAddress and WalletID: set exactly one. WalletID pays from
	// a BlindPay managed wallet; SenderWalletAddress from the customer's own
	// wallet. The API rejects both or neither (400).
	SenderWalletAddress string `json:"sender_wallet_address,omitempty"`
	WalletID            string `json:"wallet_id,omitempty"`
}

// SubmitDocumentsParams represents parameters for submitting payout documents.
type SubmitDocumentsParams struct {
	PayoutID                string                        `json:"-"`
	TransactionDocumentType types.TransactionDocumentType `json:"transaction_document_type"`
	TransactionDocumentID   string                        `json:"transaction_document_id"`
	TransactionDocumentFile string                        `json:"transaction_document_file"`
	Description             string                        `json:"description,omitempty"`
}

// SubmitDocumentsResponse represents the response when submitting payout documents.
type SubmitDocumentsResponse struct {
	// Deprecated: the API returns {"success": true}; ID is always empty.
	// Use Success.
	ID      string `json:"id"`
	Success bool   `json:"success"`
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
	TrackingDocuments   *types.TrackingDocuments   `json:"tracking_documents,omitempty"`
	SenderWalletAddress string                     `json:"sender_wallet_address"`
	BankAccountID       *string                    `json:"bank_account_id,omitempty"`
	// BillingFeeAmount is in integer cents (1000 = 10.00).
	BillingFeeAmount *float64 `json:"billing_fee_amount,omitempty"`
	CustomerID       *string  `json:"customer_id,omitempty"`
	OfframpWalletID  *string  `json:"offramp_wallet_id,omitempty"`
	PartnerFee       *int     `json:"partner_fee,omitempty"`
	PayableID        *string  `json:"payable_id,omitempty"`
	// TransactionFeeAmount is in integer cents (1000 = 10.00).
	TransactionFeeAmount *float64 `json:"transaction_fee_amount,omitempty"`
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

// List retrieves payouts with optional filters.
//
// The API returns every payout as a bare array when no pagination param is
// sent, and one page wrapped in {data, pagination} when Limit,
// StartingAfter or EndingBefore is set; both shapes are decoded and the
// rows are returned. Use ListPage for the pagination metadata, or
// ListAutoPaging to iterate over every payout.
func (c *Client) List(ctx context.Context, params *ListParams, opts ...option.RequestOption) ([]Payout, error) {
	return request.DoList[Payout](c.cfg, ctx, c.listPath(params), opts...)
}

// ListPage retrieves one page of payouts with its pagination metadata.
// Set Limit (or a cursor) to get a page; without them the API returns every
// payout and Pagination is the zero value.
func (c *Client) ListPage(ctx context.Context, params *ListParams, opts ...option.RequestOption) (*pagination.Page[Payout], error) {
	return request.DoPage[Payout](c.cfg, ctx, c.listPath(params), opts...)
}

// ListAutoPaging iterates over every payout matching params, one page at a
// time (Limit defaults to 100), starting after params.StartingAfter. It only
// pages forward: EndingBefore makes the pager fail with
// pagination.ErrEndingBeforeUnsupported. Each page request is retried per
// the client's retry policy.
func (c *Client) ListAutoPaging(ctx context.Context, params *ListParams, opts ...option.RequestOption) *pagination.Pager[Payout] {
	p := ListParams{}
	if params != nil {
		p = *params
	}
	if p.EndingBefore != "" {
		return pagination.NewErrPager[Payout](pagination.ErrEndingBeforeUnsupported)
	}
	if p.Limit == 0 {
		p.Limit = pagination.DefaultLimit
	}
	return pagination.NewPager(ctx, p.StartingAfter, func(ctx context.Context, cursor string) (*pagination.Page[Payout], error) {
		page := p
		page.StartingAfter = cursor
		return c.ListPage(ctx, &page, opts...)
	})
}

func (c *Client) listPath(params *ListParams) string {
	path := request.Pathf("/instances/%s/payouts", c.instanceID)
	if params == nil {
		return path
	}
	var q request.Query
	q.Int("limit", params.Limit)
	q.Int("offset", params.Offset)
	q.String("starting_after", params.StartingAfter)
	q.String("ending_before", params.EndingBefore)
	q.String("customer_id", params.CustomerID)
	q.String("status", string(params.Status))
	q.String("customer_name", params.CustomerName)
	q.String("bank_account_id", params.BankAccountID)
	q.String("country", params.Country)
	q.String("payment_method", string(params.PaymentMethod))
	q.String("network", string(params.Network))
	q.String("token", string(params.Token))
	return q.Apply(path)
}

// Export retrieves all payouts for export with optional pagination.
//
// Deprecated: the API has no GET /instances/{id}/export/payouts route, so this
// method always fails (404). It is kept only for source compatibility.
// Use List instead.
func (c *Client) Export(ctx context.Context, params *ExportParams, opts ...option.RequestOption) ([]Payout, error) {
	path := request.Pathf("/instances/%s/export/payouts", c.instanceID)

	if params != nil {
		var q request.Query
		q.Int("limit", params.Limit)
		q.Int("offset", params.Offset)
		path = q.Apply(path)
	}

	return request.Do[[]Payout](c.cfg, ctx, "GET", path, nil, opts...)
}

// Get retrieves a specific payout by ID.
func (c *Client) Get(ctx context.Context, payoutID string, opts ...option.RequestOption) (*Payout, error) {
	if payoutID == "" {
		return nil, fmt.Errorf("payout ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/payouts/%s", c.instanceID, payoutID)
	return request.Do[*Payout](c.cfg, ctx, "GET", path, nil, opts...)
}

// GetTrack retrieves tracking information for a payout.
//
// This is BlindPay's public tracking endpoint for its web and iOS apps:
// server-side calls get 401 in production. Use Client.Get instead.
func (c *Client) GetTrack(ctx context.Context, payoutID string, opts ...option.RequestOption) (*Payout, error) {
	if payoutID == "" {
		return nil, fmt.Errorf("payout ID cannot be empty")
	}

	path := request.Pathf("/e/payouts/%s", payoutID)
	return request.Do[*Payout](c.cfg, ctx, "GET", path, nil, opts...)
}

// CreateEvm creates an EVM payout.
func (c *Client) CreateEvm(ctx context.Context, params *CreateEvmParams, opts ...option.RequestOption) (*CreateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := request.Pathf("/instances/%s/payouts/evm", c.instanceID)
	return request.Do[*CreateResponse](c.cfg, ctx, "POST", path, params, opts...)
}

// CreateStellar creates a Stellar payout.
func (c *Client) CreateStellar(ctx context.Context, params *CreateStellarParams, opts ...option.RequestOption) (*CreateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := request.Pathf("/instances/%s/payouts/stellar", c.instanceID)
	return request.Do[*CreateResponse](c.cfg, ctx, "POST", path, params, opts...)
}

// CreateSolana creates a Solana payout.
func (c *Client) CreateSolana(ctx context.Context, params *CreateSolanaParams, opts ...option.RequestOption) (*CreateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := request.Pathf("/instances/%s/payouts/solana", c.instanceID)
	return request.Do[*CreateResponse](c.cfg, ctx, "POST", path, params, opts...)
}

// SubmitDocuments submits documents for a payout.
func (c *Client) SubmitDocuments(ctx context.Context, params *SubmitDocumentsParams, opts ...option.RequestOption) (*SubmitDocumentsResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.PayoutID == "" {
		return nil, fmt.Errorf("payout ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/payouts/%s/documents", c.instanceID, params.PayoutID)

	return request.Do[*SubmitDocumentsResponse](c.cfg, ctx, "POST", path, params, opts...)
}

// AuthorizeStellarToken authorizes a Stellar token for payout.
func (c *Client) AuthorizeStellarToken(ctx context.Context, params *AuthorizeStellarTokenParams, opts ...option.RequestOption) (*AuthorizeStellarTokenResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := request.Pathf("/instances/%s/payouts/stellar/authorize", c.instanceID)
	return request.Do[*AuthorizeStellarTokenResponse](c.cfg, ctx, "POST", path, params, opts...)
}
