package payins

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/blindpaylabs/blindpay-go/option"
	"github.com/blindpaylabs/blindpay-go/pagination"
)

// Payin represents a payin transaction.
type Payin struct {
	CustomerID            string                       `json:"customer_id"`
	ID                    string                       `json:"id"`
	PixCode               string                       `json:"pix_code,omitempty"`
	MemoCode              string                       `json:"memo_code,omitempty"`
	Clabe                 string                       `json:"clabe,omitempty"`
	Status                types.TransactionStatus      `json:"status"`
	PayinQuoteID          string                       `json:"payin_quote_id"`
	InstanceID            string                       `json:"instance_id"`
	TrackingTransaction   *types.TrackingTransaction   `json:"tracking_transaction,omitempty"`
	TrackingPayment       *types.TrackingPayment       `json:"tracking_payment,omitempty"`
	TrackingComplete      *types.TrackingComplete      `json:"tracking_complete,omitempty"`
	TrackingPartnerFee    *types.TrackingPartnerFee    `json:"tracking_partner_fee,omitempty"`
	CreatedAt             time.Time                    `json:"created_at"`
	UpdatedAt             time.Time                    `json:"updated_at"`
	ImageURL              string                       `json:"image_url,omitempty"`
	FirstName             string                       `json:"first_name,omitempty"`
	LastName              string                       `json:"last_name,omitempty"`
	LegalName             string                       `json:"legal_name,omitempty"`
	Type                  string                       `json:"type"`
	PaymentMethod         string                       `json:"payment_method"`
	SenderAmount          float64                      `json:"sender_amount"`
	ReceiverAmount        float64                      `json:"receiver_amount"`
	Token                 types.StablecoinToken        `json:"token"`
	PartnerFeeAmount      float64                      `json:"partner_fee_amount"`
	TotalFeeAmount        float64                      `json:"total_fee_amount"`
	CommercialQuotation   float64                      `json:"commercial_quotation"`
	BlindpayQuotation     float64                      `json:"blindpay_quotation"`
	Currency              string                       `json:"currency"`
	BillingFee            float64                      `json:"billing_fee"`
	Name                  string                       `json:"name"`
	Address               string                       `json:"address"`
	Network               types.Network                `json:"network"`
	BlindpayBankDetails   BankDetails                  `json:"blindpay_bank_details"`
	IsOtc                 *bool                        `json:"is_otc,omitempty"`
	BillingFeeAmount      *float64                     `json:"billing_fee_amount,omitempty"`
	PseDocumentType       *string                      `json:"pse_document_type,omitempty"`
	PseFullName           *string                      `json:"pse_full_name,omitempty"`
	PsePaymentLink        *string                      `json:"pse_payment_link,omitempty"`
	PseTaxId              *string                      `json:"pse_tax_id,omitempty"`
	PartnerFeeID          *string                      `json:"partner_fee_id,omitempty"`
	ManualExecutionStatus *types.ManualExecutionStatus `json:"manual_execution_status,omitempty"`
	PayerRules            *PayerRules                  `json:"payer_rules,omitempty"`
	// ManualConcludedAt (an ISO 8601 timestamp) and ManualConcludedBy
	// record when and by whom a manual execution was concluded.
	ManualConcludedAt    *types.LooseString `json:"manual_concluded_at,omitempty"`
	ManualConcludedBy    *string            `json:"manual_concluded_by,omitempty"`
	PartnerFee           *int               `json:"partner_fee,omitempty"`
	TransactionFeeAmount *float64           `json:"transaction_fee_amount,omitempty"`
}

// BankDetails represents bank details for a payin.
type BankDetails struct {
	RoutingNumber string `json:"routing_number"`
	AccountNumber string `json:"account_number"`
	AccountType   string `json:"account_type"`
	SwiftBicCode  string `json:"swift_bic_code"`
	ACH           struct {
		RoutingNumber string `json:"routing_number"`
		AccountNumber string `json:"account_number"`
	} `json:"ach"`
	Wire struct {
		RoutingNumber string `json:"routing_number"`
		AccountNumber string `json:"account_number"`
	} `json:"wire"`
	RTP struct {
		RoutingNumber string `json:"routing_number"`
		AccountNumber string `json:"account_number"`
	} `json:"rtp"`
	Beneficiary struct {
		Name         string `json:"name"`
		AddressLine1 string `json:"address_line_1"`
		AddressLine2 string `json:"address_line_2"`
	} `json:"beneficiary"`
	ReceivingBank struct {
		Name         string `json:"name"`
		AddressLine1 string `json:"address_line_1"`
		AddressLine2 string `json:"address_line_2"`
	} `json:"receiving_bank"`
}

// ListParams represents parameters for listing payins.
type ListParams struct {
	// Limit is the page size (10, 50, 100, 200, 500 or 1000). Setting Limit,
	// StartingAfter or EndingBefore makes the API return one page; with none
	// of them it returns every payin.
	Limit int `json:"limit,omitempty"`
	// Deprecated: ignored by the API; use StartingAfter.
	Offset int `json:"offset,omitempty"`
	// StartingAfter returns the page after this payin ID (the previous
	// page's Pagination.NextPage).
	StartingAfter string `json:"-"`
	// EndingBefore returns the page before this payin ID.
	EndingBefore  string                   `json:"-"`
	CustomerID    string                   `json:"customer_id,omitempty"`
	Status        types.TransactionStatus  `json:"status,omitempty"`
	CustomerName  string                   `json:"-"`
	BankAccountID string                   `json:"-"`
	Country       string                   `json:"-"`
	PaymentMethod types.PayinPaymentMethod `json:"-"`
	Network       types.Network            `json:"-"`
	Token         types.StablecoinToken    `json:"-"`
}

// CreateEvmResponse represents the response when creating an EVM payin.
type CreateEvmResponse struct {
	ID                   string                     `json:"id"`
	Status               types.TransactionStatus    `json:"status"`
	PixCode              string                     `json:"pix_code,omitempty"`
	MemoCode             string                     `json:"memo_code,omitempty"`
	Clabe                string                     `json:"clabe,omitempty"`
	TrackingComplete     *types.TrackingComplete    `json:"tracking_complete,omitempty"`
	TrackingPayment      *types.TrackingPayment     `json:"tracking_payment,omitempty"`
	TrackingTransaction  *types.TrackingTransaction `json:"tracking_transaction,omitempty"`
	TrackingPartnerFee   *types.TrackingPartnerFee  `json:"tracking_partner_fee,omitempty"`
	BlindpayBankDetails  BankDetails                `json:"blindpay_bank_details"`
	CustomerID           string                     `json:"customer_id"`
	ReceiverAmount       float64                    `json:"receiver_amount"`
	BillingFeeAmount     *float64                   `json:"billing_fee_amount,omitempty"`
	PartnerFee           *int                       `json:"partner_fee,omitempty"`
	PaymentMethod        *string                    `json:"payment_method,omitempty"`
	SenderAmount         *float64                   `json:"sender_amount,omitempty"`
	TransactionFeeAmount *float64                   `json:"transaction_fee_amount,omitempty"`
}

// Client handles payin-related operations.
type Client struct {
	cfg        *request.Config
	instanceID string
	Quotes     *QuotesClient
}

// NewClient creates a new payins client.
func NewClient(cfg *config.Config) *Client {
	reqCfg := cfg.ToRequestConfig()

	return &Client{
		cfg:        reqCfg,
		instanceID: cfg.InstanceID,
		Quotes:     &QuotesClient{cfg: reqCfg, instanceID: cfg.InstanceID},
	}
}

// List retrieves payins with optional filters.
//
// The API returns every payin as a bare array when no pagination param is
// sent, and one page wrapped in {data, pagination} when Limit,
// StartingAfter or EndingBefore is set; both shapes are decoded and the
// rows are returned. Use ListPage for the pagination metadata, or
// ListAutoPaging to iterate over every payin.
func (c *Client) List(ctx context.Context, params *ListParams, opts ...option.RequestOption) ([]Payin, error) {
	return request.DoList[Payin](c.cfg, ctx, c.listPath(params), opts...)
}

// ListPage retrieves one page of payins with its pagination metadata.
// Set Limit (or a cursor) to get a page; without them the API returns every
// payin and Pagination is the zero value.
func (c *Client) ListPage(ctx context.Context, params *ListParams, opts ...option.RequestOption) (*pagination.Page[Payin], error) {
	list, err := request.Do[request.List[Payin]](c.cfg, ctx, "GET", c.listPath(params), nil, opts...)
	if err != nil {
		return nil, err
	}
	page := &pagination.Page[Payin]{Data: list.Data}
	if list.Pagination != nil {
		page.Pagination = *list.Pagination
	}
	return page, nil
}

// ListAutoPaging iterates over every payin matching params, one page at a
// time (Limit defaults to 100), starting after params.StartingAfter. It only
// pages forward: EndingBefore makes the pager fail with
// pagination.ErrEndingBeforeUnsupported. Each page request is retried per
// the client's retry policy.
func (c *Client) ListAutoPaging(ctx context.Context, params *ListParams, opts ...option.RequestOption) *pagination.Pager[Payin] {
	p := ListParams{}
	if params != nil {
		p = *params
	}
	if p.EndingBefore != "" {
		return pagination.NewErrPager[Payin](pagination.ErrEndingBeforeUnsupported)
	}
	if p.Limit == 0 {
		p.Limit = pagination.DefaultLimit
	}
	return pagination.NewPager(ctx, p.StartingAfter, func(ctx context.Context, cursor string) (*pagination.Page[Payin], error) {
		page := p
		page.StartingAfter = cursor
		return c.ListPage(ctx, &page, opts...)
	})
}

func (c *Client) listPath(params *ListParams) string {
	path := request.Pathf("/instances/%s/payins", c.instanceID)
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

// Get retrieves a specific payin by ID.
func (c *Client) Get(ctx context.Context, payinID string, opts ...option.RequestOption) (*Payin, error) {
	if payinID == "" {
		return nil, fmt.Errorf("payin ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/payins/%s", c.instanceID, payinID)
	return request.Do[*Payin](c.cfg, ctx, "GET", path, nil, opts...)
}

// GetTrack retrieves tracking information for a payin.
func (c *Client) GetTrack(ctx context.Context, payinID string, opts ...option.RequestOption) (*Payin, error) {
	if payinID == "" {
		return nil, fmt.Errorf("payin ID cannot be empty")
	}

	path := request.Pathf("/e/payins/%s", payinID)
	return request.Do[*Payin](c.cfg, ctx, "GET", path, nil, opts...)
}

// Export exports payins with filters.
//
// Deprecated: the API has no GET /instances/{id}/export/payins route, so this
// method always fails (404). It will be removed in the next major version.
// Use List instead.
func (c *Client) Export(ctx context.Context, status types.TransactionStatus, limit, offset int, opts ...option.RequestOption) ([]Payin, error) {
	path := request.Pathf("/instances/%s/export/payins", c.instanceID)

	q := url.Values{}
	if status != "" {
		q.Set("status", string(status))
	}
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	if offset > 0 {
		q.Set("offset", fmt.Sprintf("%d", offset))
	}
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	return request.Do[[]Payin](c.cfg, ctx, "GET", path, nil, opts...)
}

// CreateEvm creates an EVM payin.
func (c *Client) CreateEvm(ctx context.Context, payinQuoteID string, opts ...option.RequestOption) (*CreateEvmResponse, error) {
	if payinQuoteID == "" {
		return nil, fmt.Errorf("payin quote ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/payins/evm", c.instanceID)
	body := map[string]string{
		"payin_quote_id": payinQuoteID,
	}

	return request.Do[*CreateEvmResponse](c.cfg, ctx, "POST", path, body, opts...)
}
