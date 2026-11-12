// Package custodialwallets manages BlindPay-managed (custodial) wallets for a
// customer: create, list, get, delete, balances and the deposit history
// (ListDeposits, ListDepositsAutoPaging). Wallet balances are whole tokens
// (10.5 = 10.5 USDC). Use it through blindpay.Client.CustodialWallets.
package custodialwallets

import (
	"context"
	"fmt"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/blindpaylabs/blindpay-go/option"
)

// YieldStatus is whether the wallet earns yield.
type YieldStatus string

const (
	YieldStatusDisabled  YieldStatus = "disabled"
	YieldStatusEnabled   YieldStatus = "enabled"
	YieldStatusDisabling YieldStatus = "disabling"
)

// CustodialWallet represents a custodial wallet.
type CustodialWallet struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	ExternalID  *string       `json:"external_id"`
	Address     *string       `json:"address"`
	Network     types.Network `json:"network"`
	CreatedAt   time.Time     `json:"created_at"`
	YieldStatus *YieldStatus  `json:"yield_status,omitempty"`
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
	// Amounts are whole tokens (10.5 = 10.5 USDC). TotalAmount is Amount
	// plus yield earned (EarningAmount).
	Address       string                `json:"address"`
	ID            string                `json:"id"`
	Symbol        types.StablecoinToken `json:"symbol"`
	Amount        float64               `json:"amount"`
	EarningAmount float64               `json:"earning_amount"`
	TotalAmount   float64               `json:"total_amount"`
}

// GetBalanceResponse represents the balance of a custodial wallet.
type GetBalanceResponse struct {
	USDC WalletTokenBalance `json:"USDC"`
	USDT WalletTokenBalance `json:"USDT"`
	USDB WalletTokenBalance `json:"USDB"`
	// EURC and EURB are present on networks that hold them.
	EURC *WalletTokenBalance `json:"EURC,omitempty"`
	EURB *WalletTokenBalance `json:"EURB,omitempty"`
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
func (c *Client) List(ctx context.Context, customerID string, opts ...option.RequestOption) ([]CustodialWallet, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/wallets", c.instanceID, customerID)
	return request.Do[[]CustodialWallet](c.cfg, ctx, "GET", path, nil, opts...)
}

// Get retrieves a specific custodial wallet.
func (c *Client) Get(ctx context.Context, customerID, id string, opts ...option.RequestOption) (*CustodialWallet, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/wallets/%s", c.instanceID, customerID, id)
	return request.Do[*CustodialWallet](c.cfg, ctx, "GET", path, nil, opts...)
}

// Create creates a new custodial wallet.
func (c *Client) Create(ctx context.Context, params *CreateParams, opts ...option.RequestOption) (*CustodialWallet, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.CustomerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/wallets", c.instanceID, params.CustomerID)

	return request.Do[*CustodialWallet](c.cfg, ctx, "POST", path, params, opts...)
}

// GetBalance retrieves the balance of a custodial wallet.
func (c *Client) GetBalance(ctx context.Context, customerID, id string, opts ...option.RequestOption) (*GetBalanceResponse, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/wallets/%s/balance", c.instanceID, customerID, id)
	return request.Do[*GetBalanceResponse](c.cfg, ctx, "GET", path, nil, opts...)
}

// Delete deletes a custodial wallet.
func (c *Client) Delete(ctx context.Context, customerID, id string, opts ...option.RequestOption) error {
	if customerID == "" {
		return fmt.Errorf("customer ID cannot be empty")
	}
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/wallets/%s", c.instanceID, customerID, id)
	_, err := request.Do[struct{}](c.cfg, ctx, "DELETE", path, nil, opts...)
	return err
}
