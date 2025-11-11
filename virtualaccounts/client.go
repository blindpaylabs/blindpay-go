package virtualaccounts

import (
	"context"
	"fmt"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
)

// VirtualAccountUS represents US-specific virtual account details.
type VirtualAccountUS struct {
	ACH struct {
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
	SwiftBICCode string `json:"swift_bic_code"`
	AccountType  string `json:"account_type"`
	Beneficiary  struct {
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

// VirtualAccount represents a virtual account.
type VirtualAccount struct {
	ID                 string                `json:"id"`
	US                 VirtualAccountUS      `json:"us"`
	Token              types.StablecoinToken `json:"token"`
	BlockchainWalletID string                `json:"blockchain_wallet_id"`
}

// CreateParams represents parameters for creating a virtual account.
type CreateParams struct {
	ReceiverID         string                `json:"-"`
	BlockchainWalletID string                `json:"blockchain_wallet_id"`
	Token              types.StablecoinToken `json:"token"`
}

// UpdateParams represents parameters for updating a virtual account.
type UpdateParams struct {
	ReceiverID         string                `json:"-"`
	BlockchainWalletID string                `json:"blockchain_wallet_id"`
	Token              types.StablecoinToken `json:"token"`
}

// Client handles virtual account-related operations.
type Client struct {
	cfg        *request.Config
	instanceID string
}

// NewClient creates a new virtual accounts client.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:        cfg.ToRequestConfig(),
		instanceID: cfg.InstanceID,
	}
}

// Create creates a new virtual account.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*VirtualAccount, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.ReceiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/virtual-accounts", c.instanceID, params.ReceiverID)

	body := struct {
		BlockchainWalletID string                `json:"blockchain_wallet_id"`
		Token              types.StablecoinToken `json:"token"`
	}{
		BlockchainWalletID: params.BlockchainWalletID,
		Token:              params.Token,
	}

	return request.Do[*VirtualAccount](c.cfg, ctx, "POST", path, body)
}

// Get retrieves a virtual account for a receiver.
func (c *Client) Get(ctx context.Context, receiverID string) (*VirtualAccount, error) {
	if receiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/virtual-accounts", c.instanceID, receiverID)
	return request.Do[*VirtualAccount](c.cfg, ctx, "GET", path, nil)
}

// Update updates a virtual account.
func (c *Client) Update(ctx context.Context, params *UpdateParams) error {
	if params == nil {
		return fmt.Errorf("params cannot be nil")
	}
	if params.ReceiverID == "" {
		return fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/virtual-accounts", c.instanceID, params.ReceiverID)

	body := struct {
		BlockchainWalletID string                `json:"blockchain_wallet_id"`
		Token              types.StablecoinToken `json:"token"`
	}{
		BlockchainWalletID: params.BlockchainWalletID,
		Token:              params.Token,
	}

	_, err := request.Do[struct{}](c.cfg, ctx, "PUT", path, body)
	return err
}
