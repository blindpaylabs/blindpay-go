package transfers

import (
	"context"
	"fmt"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/blindpaylabs/blindpay-go/option"
	"github.com/blindpaylabs/blindpay-go/pagination"
)

// TrackingStep represents a transfer tracking step.
type TrackingStep struct {
	Step            string  `json:"step"`
	TransactionHash *string `json:"transaction_hash"`
	GasFee          *string `json:"gas_fee"`
	CompletedAt     *string `json:"completed_at"`
	ErrorMessage    *string `json:"error_message"`
}

// TrackingTransactionMonitoring represents transaction monitoring tracking.
type TrackingTransactionMonitoring struct {
	Step                string   `json:"step"`
	BlockchainScreening *float64 `json:"blockchain_screening"`
	RiskScore           *float64 `json:"risk_score"`
	CompletedAt         *string  `json:"completed_at"`
}

// Transfer represents a transfer transaction.
type Transfer struct {
	ID                            string                        `json:"id"`
	Status                        types.TransactionStatus       `json:"status"`
	TransferQuoteID               string                        `json:"transfer_quote_id"`
	InstanceID                    string                        `json:"instance_id"`
	TrackingTransactionMonitoring TrackingTransactionMonitoring `json:"tracking_transaction_monitoring"`
	TrackingPaymaster             TrackingStep                  `json:"tracking_paymaster"`
	TrackingBridgeSwap            TrackingStep                  `json:"tracking_bridge_swap"`
	TrackingComplete              TrackingStep                  `json:"tracking_complete"`
	TrackingPartnerFee            TrackingStep                  `json:"tracking_partner_fee"`
	CreatedAt                     time.Time                     `json:"created_at"`
	UpdatedAt                     time.Time                     `json:"updated_at"`
	ImageURL                      string                        `json:"image_url,omitempty"`
	FirstName                     string                        `json:"first_name,omitempty"`
	LastName                      string                        `json:"last_name,omitempty"`
	LegalName                     string                        `json:"legal_name,omitempty"`
	WalletID                      string                        `json:"wallet_id"`
	SenderToken                   types.StablecoinToken         `json:"sender_token"`
	SenderAmount                  float64                       `json:"sender_amount"`
	ReceiverAmount                float64                       `json:"receiver_amount"`
	CustomerNetwork               types.Network                 `json:"customer_network"`
	CustomerToken                 types.StablecoinToken         `json:"customer_token"`
	CustomerWalletAddress         string                        `json:"customer_wallet_address"`
	PartnerFeeAmount              *float64                      `json:"partner_fee_amount"`
	ExternalID                    *string                       `json:"external_id,omitempty"`
	CustomerID                    string                        `json:"customer_id"`
	Address                       string                        `json:"address"`
	Network                       types.Network                 `json:"network"`
}

// CreateQuoteParams represents parameters for creating a transfer quote.
type CreateQuoteParams struct {
	WalletID              string                `json:"wallet_id"`
	SenderToken           types.StablecoinToken `json:"sender_token"`
	CustomerWalletAddress string                `json:"customer_wallet_address"`
	CustomerToken         types.StablecoinToken `json:"customer_token"`
	CustomerNetwork       types.Network         `json:"customer_network"`
	RequestAmount         float64               `json:"request_amount"`
	CoverFees             bool                  `json:"cover_fees,omitempty"`
	AmountReference       string                `json:"amount_reference"`
	PartnerFeeID          *string               `json:"partner_fee_id,omitempty"`
	// PartnerPercentageFee is in basis points (100 = 1%) and PartnerFlatFee
	// in integer cents; they override the partner fee for this quote.
	PartnerPercentageFee *int64 `json:"partner_percentage_fee,omitempty"`
	PartnerFlatFee       *int64 `json:"partner_flat_fee,omitempty"`
}

// CreateQuoteResponse represents the response when creating a transfer quote.
type CreateQuoteResponse struct {
	ID                  string   `json:"id"`
	ExpiresAt           *float64 `json:"expires_at"`
	CommercialQuotation *float64 `json:"commercial_quotation"`
	BlindpayQuotation   *float64 `json:"blindpay_quotation"`
	ReceiverAmount      float64  `json:"receiver_amount"`
	SenderAmount        float64  `json:"sender_amount"`
	PartnerFeeAmount    *float64 `json:"partner_fee_amount"`
	FlatFee             float64  `json:"flat_fee"`
}

// CreateParams represents parameters for creating a transfer.
type CreateParams struct {
	TransferQuoteID string `json:"transfer_quote_id"`
}

// CreateResponse represents the response when creating a transfer.
type CreateResponse struct {
	ID                            string                        `json:"id"`
	Status                        types.TransactionStatus       `json:"status"`
	TrackingBridgeSwap            TrackingStep                  `json:"tracking_bridge_swap"`
	TrackingComplete              TrackingStep                  `json:"tracking_complete"`
	TrackingPaymaster             TrackingStep                  `json:"tracking_paymaster"`
	TrackingTransactionMonitoring TrackingTransactionMonitoring `json:"tracking_transaction_monitoring"`
	TrackingPartnerFee            TrackingStep                  `json:"tracking_partner_fee"`
}

// ListParams represents parameters for listing transfers. The API always
// returns one page (default 10 rows).
type ListParams struct {
	// Limit is the page size (10, 50, 100, 200, 500 or 1000).
	Limit int `json:"limit,omitempty"`
	// Deprecated: ignored by the API; use StartingAfter.
	Offset int `json:"offset,omitempty"`
	// StartingAfter returns the page after this transfer ID (the previous
	// page's Pagination.NextPage).
	StartingAfter string `json:"-"`
	// EndingBefore returns the page before this transfer ID.
	EndingBefore string                  `json:"-"`
	CustomerID   string                  `json:"-"`
	WalletID     string                  `json:"-"`
	Status       types.TransactionStatus `json:"-"`
	Network      types.Network           `json:"-"`
	Token        types.StablecoinToken   `json:"-"`
}

// ListResponse represents the response when listing transfers.
type ListResponse struct {
	Data       []Transfer               `json:"data"`
	Pagination types.PaginationMetadata `json:"pagination"`
}

// Client handles transfer-related operations.
type Client struct {
	cfg        *request.Config
	instanceID string
	Quotes     *QuotesClient
}

// QuotesClient handles transfer quote operations.
type QuotesClient struct {
	cfg        *request.Config
	instanceID string
}

// NewClient creates a new transfers client.
func NewClient(cfg *config.Config) *Client {
	reqCfg := cfg.ToRequestConfig()
	return &Client{
		cfg:        reqCfg,
		instanceID: cfg.InstanceID,
		Quotes: &QuotesClient{
			cfg:        reqCfg,
			instanceID: cfg.InstanceID,
		},
	}
}

// Create creates a new transfer quote.
func (c *QuotesClient) Create(ctx context.Context, params *CreateQuoteParams, opts ...option.RequestOption) (*CreateQuoteResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := request.Pathf("/instances/%s/transfer-quotes", c.instanceID)
	return request.Do[*CreateQuoteResponse](c.cfg, ctx, "POST", path, params, opts...)
}

// Create creates a new transfer.
func (c *Client) Create(ctx context.Context, params *CreateParams, opts ...option.RequestOption) (*CreateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := request.Pathf("/instances/%s/transfers", c.instanceID)
	return request.Do[*CreateResponse](c.cfg, ctx, "POST", path, params, opts...)
}

// List retrieves one page of transfers (the API always paginates them).
func (c *Client) List(ctx context.Context, params *ListParams, opts ...option.RequestOption) (*ListResponse, error) {
	page, err := c.ListPage(ctx, params, opts...)
	if err != nil {
		return nil, err
	}
	return &ListResponse{Data: page.Data, Pagination: page.Pagination}, nil
}

// ListPage retrieves one page of transfers with its pagination metadata.
func (c *Client) ListPage(ctx context.Context, params *ListParams, opts ...option.RequestOption) (*pagination.Page[Transfer], error) {
	list, err := request.Do[request.List[Transfer]](c.cfg, ctx, "GET", c.listPath(params), nil, opts...)
	if err != nil {
		return nil, err
	}
	page := &pagination.Page[Transfer]{Data: list.Data}
	if list.Pagination != nil {
		page.Pagination = *list.Pagination
	}
	return page, nil
}

// ListAutoPaging iterates over every transfer matching params, one page at
// a time (Limit defaults to 100), starting after params.StartingAfter. It
// only pages forward: EndingBefore makes the pager fail with
// pagination.ErrEndingBeforeUnsupported.
func (c *Client) ListAutoPaging(ctx context.Context, params *ListParams, opts ...option.RequestOption) *pagination.Pager[Transfer] {
	p := ListParams{}
	if params != nil {
		p = *params
	}
	if p.EndingBefore != "" {
		return pagination.NewErrPager[Transfer](pagination.ErrEndingBeforeUnsupported)
	}
	if p.Limit == 0 {
		p.Limit = pagination.DefaultLimit
	}
	return pagination.NewPager(ctx, p.StartingAfter, func(ctx context.Context, cursor string) (*pagination.Page[Transfer], error) {
		page := p
		page.StartingAfter = cursor
		return c.ListPage(ctx, &page, opts...)
	})
}

func (c *Client) listPath(params *ListParams) string {
	path := request.Pathf("/instances/%s/transfers", c.instanceID)
	if params == nil {
		return path
	}
	var q request.Query
	q.Int("limit", params.Limit)
	q.Int("offset", params.Offset)
	q.String("starting_after", params.StartingAfter)
	q.String("ending_before", params.EndingBefore)
	q.String("customer_id", params.CustomerID)
	q.String("wallet_id", params.WalletID)
	q.String("status", string(params.Status))
	q.String("network", string(params.Network))
	q.String("token", string(params.Token))
	return q.Apply(path)
}

// Get retrieves a specific transfer by ID.
func (c *Client) Get(ctx context.Context, transferID string, opts ...option.RequestOption) (*Transfer, error) {
	if transferID == "" {
		return nil, fmt.Errorf("transfer ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/transfers/%s", c.instanceID, transferID)
	return request.Do[*Transfer](c.cfg, ctx, "GET", path, nil, opts...)
}

// GetTrack retrieves tracking information for a transfer (public endpoint).
func (c *Client) GetTrack(ctx context.Context, transferID string, opts ...option.RequestOption) (*Transfer, error) {
	if transferID == "" {
		return nil, fmt.Errorf("transfer ID cannot be empty")
	}

	path := request.Pathf("/e/transfers/%s", transferID)
	return request.Do[*Transfer](c.cfg, ctx, "GET", path, nil, opts...)
}
