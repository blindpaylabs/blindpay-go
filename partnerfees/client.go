package partnerfees

import (
	"context"
	"fmt"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
)

// PartnerFee represents a partner fee configuration.
type PartnerFee struct {
	ID                   string  `json:"id"`
	InstanceID           string  `json:"instance_id"`
	Name                 string  `json:"name"`
	PayoutPercentageFee  float64 `json:"payout_percentage_fee"`
	PayoutFlatFee        float64 `json:"payout_flat_fee"`
	PayinPercentageFee   float64 `json:"payin_percentage_fee"`
	PayinFlatFee         float64 `json:"payin_flat_fee"`
	EVMWalletAddress     string  `json:"evm_wallet_address,omitempty"`
	StellarWalletAddress string  `json:"stellar_wallet_address,omitempty"`
}

// CreatePartnerFeeParams represents parameters for creating a partner fee.
type CreatePartnerFeeParams struct {
	VirtualAccountSet    *bool   `json:"virtual_account_set,omitempty"`
	EVMWalletAddress     string  `json:"evm_wallet_address"`
	Name                 string  `json:"name"`
	PayinFlatFee         float64 `json:"payin_flat_fee"`
	PayinPercentageFee   float64 `json:"payin_percentage_fee"`
	PayoutFlatFee        float64 `json:"payout_flat_fee"`
	PayoutPercentageFee  float64 `json:"payout_percentage_fee"`
	StellarWalletAddress *string `json:"stellar_wallet_address,omitempty"`
}

// Client handles partner fee-related operations.
type Client struct {
	cfg        *request.Config
	instanceID string
}

// NewClient creates a new partner fees client.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg: &request.Config{
			BaseURL:    cfg.BaseURL,
			APIKey:     cfg.APIKey,
			HTTPClient: cfg.HTTPClient,
			UserAgent:  cfg.UserAgent,
		},
		instanceID: cfg.InstanceID,
	}
}

// List retrieves all partner fees for the instance.
func (c *Client) List(ctx context.Context) ([]PartnerFee, error) {
	path := fmt.Sprintf("/instances/%s/partner-fees", c.instanceID)
	return request.Do[[]PartnerFee](c.cfg, ctx, "GET", path, nil)
}

// Create creates a new partner fee configuration.
func (c *Client) Create(ctx context.Context, params *CreatePartnerFeeParams) (*PartnerFee, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := fmt.Sprintf("/instances/%s/partner-fees", c.instanceID)
	return request.Do[*PartnerFee](c.cfg, ctx, "POST", path, params)
}

// Get retrieves a specific partner fee by ID.
func (c *Client) Get(ctx context.Context, id string) (*PartnerFee, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/partner-fees/%s", c.instanceID, id)
	return request.Do[*PartnerFee](c.cfg, ctx, "GET", path, nil)
}

// Delete deletes a partner fee configuration.
func (c *Client) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/partner-fees/%s", c.instanceID, id)
	_, err := request.Do[struct{}](c.cfg, ctx, "DELETE", path, nil)
	return err
}
