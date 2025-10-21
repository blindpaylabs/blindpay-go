package payins

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
)

// Payin represents a payin transaction.
type Payin struct {
	ReceiverID          string                     `json:"receiver_id"`
	ID                  string                     `json:"id"`
	PixCode             string                     `json:"pix_code,omitempty"`
	MemoCode            string                     `json:"memo_code,omitempty"`
	Clabe               string                     `json:"clabe,omitempty"`
	Status              types.TransactionStatus    `json:"status"`
	PayinQuoteID        string                     `json:"payin_quote_id"`
	InstanceID          string                     `json:"instance_id"`
	TrackingTransaction *types.TrackingTransaction `json:"tracking_transaction,omitempty"`
	TrackingPayment     *types.TrackingPayment     `json:"tracking_payment,omitempty"`
	TrackingComplete    *types.TrackingComplete    `json:"tracking_complete,omitempty"`
	TrackingPartnerFee  *types.TrackingPartnerFee  `json:"tracking_partner_fee,omitempty"`
	CreatedAt           time.Time                  `json:"created_at"`
	UpdatedAt           time.Time                  `json:"updated_at"`
	ImageURL            string                     `json:"image_url,omitempty"`
	FirstName           string                     `json:"first_name,omitempty"`
	LastName            string                     `json:"last_name,omitempty"`
	LegalName           string                     `json:"legal_name,omitempty"`
	Type                string                     `json:"type"`
	PaymentMethod       string                     `json:"payment_method"`
	SenderAmount        float64                    `json:"sender_amount"`
	ReceiverAmount      float64                    `json:"receiver_amount"`
	Token               types.StablecoinToken      `json:"token"`
	PartnerFeeAmount    float64                    `json:"partner_fee_amount"`
	TotalFeeAmount      float64                    `json:"total_fee_amount"`
	CommercialQuotation float64                    `json:"commercial_quotation"`
	BlindpayQuotation   float64                    `json:"blindpay_quotation"`
	Currency            string                     `json:"currency"`
	BillingFee          float64                    `json:"billing_fee"`
	Name                string                     `json:"name"`
	Address             string                     `json:"address"`
	Network             types.Network              `json:"network"`
	BlindpayBankDetails BankDetails                `json:"blindpay_bank_details"`
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
	Status     types.TransactionStatus `json:"status,omitempty"`
	ReceiverID string                  `json:"receiver_id,omitempty"`
	Limit      int                     `json:"limit,omitempty"`
	Offset     int                     `json:"offset,omitempty"`
}

// ListResponse represents the response when listing payins.
type ListResponse struct {
	Data       []Payin                  `json:"data"`
	Pagination types.PaginationMetadata `json:"pagination"`
}

// CreateEvmResponse represents the response when creating an EVM payin.
type CreateEvmResponse struct {
	ID                  string                     `json:"id"`
	Status              types.TransactionStatus    `json:"status"`
	PixCode             string                     `json:"pix_code,omitempty"`
	MemoCode            string                     `json:"memo_code,omitempty"`
	Clabe               string                     `json:"clabe,omitempty"`
	TrackingComplete    *types.TrackingComplete    `json:"tracking_complete,omitempty"`
	TrackingPayment     *types.TrackingPayment     `json:"tracking_payment,omitempty"`
	TrackingTransaction *types.TrackingTransaction `json:"tracking_transaction,omitempty"`
	TrackingPartnerFee  *types.TrackingPartnerFee  `json:"tracking_partner_fee,omitempty"`
	BlindpayBankDetails BankDetails                `json:"blindpay_bank_details"`
	ReceiverID          string                     `json:"receiver_id"`
	ReceiverAmount      float64                    `json:"receiver_amount"`
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

// List retrieves all payins with optional filters.
func (c *Client) List(ctx context.Context, params *ListParams) (*ListResponse, error) {
	path := fmt.Sprintf("/instances/%s/payins", c.instanceID)

	if params != nil {
		q := url.Values{}
		if params.Status != "" {
			q.Set("status", string(params.Status))
		}
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

// Get retrieves a specific payin by ID.
func (c *Client) Get(ctx context.Context, payinID string) (*Payin, error) {
	if payinID == "" {
		return nil, fmt.Errorf("payin ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/payins/%s", c.instanceID, payinID)
	return request.Do[*Payin](c.cfg, ctx, "GET", path, nil)
}

// GetTrack retrieves tracking information for a payin.
func (c *Client) GetTrack(ctx context.Context, payinID string) (*Payin, error) {
	if payinID == "" {
		return nil, fmt.Errorf("payin ID cannot be empty")
	}

	path := fmt.Sprintf("/e/payins/%s", payinID)
	return request.Do[*Payin](c.cfg, ctx, "GET", path, nil)
}

// Export exports payins with filters.
func (c *Client) Export(ctx context.Context, status types.TransactionStatus, limit, offset int) ([]Payin, error) {
	path := fmt.Sprintf("/instances/%s/export/payins", c.instanceID)

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

	return request.Do[[]Payin](c.cfg, ctx, "GET", path, nil)
}

// CreateEvm creates an EVM payin.
func (c *Client) CreateEvm(ctx context.Context, payinQuoteID string) (*CreateEvmResponse, error) {
	if payinQuoteID == "" {
		return nil, fmt.Errorf("payin quote ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/payins/evm", c.instanceID)
	body := map[string]string{
		"payin_quote_id": payinQuoteID,
	}

	return request.Do[*CreateEvmResponse](c.cfg, ctx, "POST", path, body)
}
