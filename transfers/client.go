package transfers

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
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
	ReceiverNetwork               types.Network                 `json:"receiver_network"`
	ReceiverToken                 types.StablecoinToken         `json:"receiver_token"`
	ReceiverWalletAddress         string                        `json:"receiver_wallet_address"`
	PartnerFeeAmount              *float64                      `json:"partner_fee_amount"`
	ExternalID                    *string                       `json:"external_id,omitempty"`
}

// CreateQuoteParams represents parameters for creating a transfer quote.
type CreateQuoteParams struct {
	WalletID              string                `json:"wallet_id"`
	SenderToken           types.StablecoinToken `json:"sender_token"`
	ReceiverWalletAddress string                `json:"receiver_wallet_address"`
	ReceiverToken         types.StablecoinToken `json:"receiver_token"`
	ReceiverNetwork       types.Network         `json:"receiver_network"`
	RequestAmount         float64               `json:"request_amount"`
	CoverFees             bool                  `json:"cover_fees"`
	AmountReference       string                `json:"amount_reference"`
	PartnerFeeID          *string               `json:"partner_fee_id,omitempty"`
}

// CreateQuoteResponse represents the response when creating a transfer quote.
type CreateQuoteResponse struct {
	ID                  string   `json:"id"`
	ExpiresAt           int64    `json:"expires_at"`
	CommercialQuotation float64  `json:"commercial_quotation"`
	BlindpayQuotation   float64  `json:"blindpay_quotation"`
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

// ListParams represents parameters for listing transfers.
type ListParams struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
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
func (c *QuotesClient) Create(ctx context.Context, params *CreateQuoteParams) (*CreateQuoteResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := fmt.Sprintf("/instances/%s/transfer-quotes", c.instanceID)
	return request.Do[*CreateQuoteResponse](c.cfg, ctx, "POST", path, params)
}

// Create creates a new transfer.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*CreateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := fmt.Sprintf("/instances/%s/transfers", c.instanceID)
	return request.Do[*CreateResponse](c.cfg, ctx, "POST", path, params)
}

// List retrieves all transfers with optional pagination.
func (c *Client) List(ctx context.Context, params *ListParams) (*ListResponse, error) {
	path := fmt.Sprintf("/instances/%s/transfers", c.instanceID)

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

	return request.Do[*ListResponse](c.cfg, ctx, "GET", path, nil)
}

// Get retrieves a specific transfer by ID.
func (c *Client) Get(ctx context.Context, transferID string) (*Transfer, error) {
	if transferID == "" {
		return nil, fmt.Errorf("transfer ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/transfers/%s", c.instanceID, transferID)
	return request.Do[*Transfer](c.cfg, ctx, "GET", path, nil)
}

// GetTrack retrieves tracking information for a transfer (public endpoint).
func (c *Client) GetTrack(ctx context.Context, transferID string) (*Transfer, error) {
	if transferID == "" {
		return nil, fmt.Errorf("transfer ID cannot be empty")
	}

	path := fmt.Sprintf("/e/transfers/%s", transferID)
	return request.Do[*Transfer](c.cfg, ctx, "GET", path, nil)
}
