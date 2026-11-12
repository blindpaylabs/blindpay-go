package partnerfees

import (
	"context"
	"fmt"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/option"
)

// PartnerFee represents a partner fee configuration.
type PartnerFee struct {
	ID         string `json:"id"`
	InstanceID string `json:"instance_id"`
	Name       string `json:"name"`
	// PayoutPercentageFee is in basis points (100 = 1%).
	PayoutPercentageFee float64 `json:"payout_percentage_fee"`
	// PayoutFlatFee is in integer cents (1000 = 10.00).
	PayoutFlatFee float64 `json:"payout_flat_fee"`
	// PayinPercentageFee is in basis points (100 = 1%).
	PayinPercentageFee float64 `json:"payin_percentage_fee"`
	// PayinFlatFee is in integer cents (1000 = 10.00).
	PayinFlatFee         float64 `json:"payin_flat_fee"`
	EVMWalletAddress     string  `json:"evm_wallet_address,omitempty"`
	StellarWalletAddress string  `json:"stellar_wallet_address,omitempty"`
}

// CreatePartnerFeeParams represents parameters for creating a partner fee.
type CreatePartnerFeeParams struct {
	VirtualAccountSet *bool  `json:"virtual_account_set,omitempty"`
	EVMWalletAddress  string `json:"evm_wallet_address"`
	Name              string `json:"name"`
	// PayinFlatFee is in integer cents (1000 = 10.00).
	PayinFlatFee float64 `json:"payin_flat_fee"`
	// PayinPercentageFee is in basis points (100 = 1%).
	PayinPercentageFee float64 `json:"payin_percentage_fee"`
	// PayoutFlatFee is in integer cents (1000 = 10.00).
	PayoutFlatFee float64 `json:"payout_flat_fee"`
	// PayoutPercentageFee is in basis points (100 = 1%).
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
		cfg:        cfg.ToRequestConfig(),
		instanceID: cfg.InstanceID,
	}
}

// List retrieves all partner fees for the instance.
func (c *Client) List(ctx context.Context, opts ...option.RequestOption) ([]PartnerFee, error) {
	path := request.Pathf("/instances/%s/partner-fees", c.instanceID)
	return request.Do[[]PartnerFee](c.cfg, ctx, "GET", path, nil, opts...)
}

// Create creates a new partner fee configuration.
func (c *Client) Create(ctx context.Context, params *CreatePartnerFeeParams, opts ...option.RequestOption) (*PartnerFee, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := request.Pathf("/instances/%s/partner-fees", c.instanceID)
	return request.Do[*PartnerFee](c.cfg, ctx, "POST", path, params, opts...)
}

// Get retrieves a specific partner fee by ID.
func (c *Client) Get(ctx context.Context, id string, opts ...option.RequestOption) (*PartnerFee, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	path := request.Pathf("/instances/%s/partner-fees/%s", c.instanceID, id)
	return request.Do[*PartnerFee](c.cfg, ctx, "GET", path, nil, opts...)
}

// Delete deletes a partner fee configuration.
func (c *Client) Delete(ctx context.Context, id string, opts ...option.RequestOption) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	path := request.Pathf("/instances/%s/partner-fees/%s", c.instanceID, id)
	_, err := request.Do[struct{}](c.cfg, ctx, "DELETE", path, nil, opts...)
	return err
}
