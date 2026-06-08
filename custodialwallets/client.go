package custodialwallets

import (
	"context"
	"fmt"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
)

// CustodialWallet represents a custodial wallet.
type CustodialWallet struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	ExternalID *string       `json:"external_id"`
	Address    *string       `json:"address"`
	Network    types.Network `json:"network"`
	CreatedAt  time.Time     `json:"created_at"`
}

// CreateParams represents parameters for creating a custodial wallet.
type CreateParams struct {
	CustomerID string        `json:"-"`
	Name       string        `json:"name"`
	Network    types.Network `json:"network"`
	ExternalID *string       `json:"external_id,omitempty"`
}

// WalletTokenBalance represents a token balance in a wallet.
type WalletTokenBalance struct {
	Address string                `json:"address"`
	ID      string                `json:"id"`
	Symbol  types.StablecoinToken `json:"symbol"`
	Amount  float64               `json:"amount"`
}

// GetBalanceResponse represents the balance of a custodial wallet.
type GetBalanceResponse struct {
	USDC WalletTokenBalance `json:"USDC"`
	USDT WalletTokenBalance `json:"USDT"`
	USDB WalletTokenBalance `json:"USDB"`
}

// Client handles custodial wallet-related operations.
type Client struct {
	cfg        *request.Config
	instanceID string
}

// NewClient creates a new custodial wallets client.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:        cfg.ToRequestConfig(),
		instanceID: cfg.InstanceID,
	}
}

// List retrieves all custodial wallets for a customer.
func (c *Client) List(ctx context.Context, customerID string) ([]CustodialWallet, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/customers/%s/wallets", c.instanceID, customerID)
	return request.Do[[]CustodialWallet](c.cfg, ctx, "GET", path, nil)
}

// Get retrieves a specific custodial wallet.
func (c *Client) Get(ctx context.Context, customerID, id string) (*CustodialWallet, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/customers/%s/wallets/%s", c.instanceID, customerID, id)
	return request.Do[*CustodialWallet](c.cfg, ctx, "GET", path, nil)
}

// Create creates a new custodial wallet.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*CustodialWallet, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.CustomerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/customers/%s/wallets", c.instanceID, params.CustomerID)

	body := map[string]any{
		"name":    params.Name,
		"network": params.Network,
	}

	if params.ExternalID != nil {
		body["external_id"] = *params.ExternalID
	}

	return request.Do[*CustodialWallet](c.cfg, ctx, "POST", path, body)
}

// GetBalance retrieves the balance of a custodial wallet.
func (c *Client) GetBalance(ctx context.Context, customerID, id string) (*GetBalanceResponse, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/customers/%s/wallets/%s/balance", c.instanceID, customerID, id)
	return request.Do[*GetBalanceResponse](c.cfg, ctx, "GET", path, nil)
}

// Delete deletes a custodial wallet.
func (c *Client) Delete(ctx context.Context, customerID, id string) error {
	if customerID == "" {
		return fmt.Errorf("customer ID cannot be empty")
	}
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/customers/%s/wallets/%s", c.instanceID, customerID, id)
	_, err := request.Do[struct{}](c.cfg, ctx, "DELETE", path, nil)
	return err
}
