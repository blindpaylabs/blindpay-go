package wallets

import (
	"context"
	"fmt"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/blindpaylabs/blindpay-go/option"
)

// BlockchainWallet represents a blockchain wallet.
type BlockchainWallet struct {
	ID                   string        `json:"id"`
	Name                 string        `json:"name"`
	Network              types.Network `json:"network"`
	Address              string        `json:"address,omitempty"`
	SignatureTxHash      string        `json:"signature_tx_hash,omitempty"`
	IsAccountAbstraction bool          `json:"is_account_abstraction"`
	CustomerID           string        `json:"customer_id"`
	// IsSelfCustody reports whether the customer holds the wallet's private
	// keys. nil means it was never answered; once set it cannot be changed.
	IsSelfCustody *bool `json:"is_self_custody,omitempty"`
}

// GetMessageResponse represents the wallet message response.
type GetMessageResponse struct {
	Message string `json:"message"`
}

// CreateWithAddressParams represents parameters for creating a wallet with address.
type CreateWithAddressParams struct {
	CustomerID string        `json:"-"`
	Name       string        `json:"name"`
	Network    types.Network `json:"network"`
	Address    string        `json:"address"`
	// IsSelfCustody declares whether the customer holds the wallet's private
	// keys. Required when the customer's country is BR (the API returns
	// 400 self_custody_required without it since 2026-10-10); optional
	// otherwise. nil omits the field.
	IsSelfCustody *bool `json:"is_self_custody,omitempty"`
}

// CreateWithHashParams represents parameters for creating a wallet with hash.
type CreateWithHashParams struct {
	CustomerID      string        `json:"-"`
	Name            string        `json:"name"`
	Network         types.Network `json:"network"`
	SignatureTxHash string        `json:"signature_tx_hash"`
	// IsSelfCustody declares whether the customer holds the wallet's private
	// keys. Required when the customer's country is BR (the API returns
	// 400 self_custody_required without it since 2026-10-10); optional
	// otherwise. nil omits the field.
	IsSelfCustody *bool `json:"is_self_custody,omitempty"`
}

// CreateAssetTrustlineResponse represents the asset trustline response.
type CreateAssetTrustlineResponse struct {
	XDR string `json:"xdr"`
}

// MintUsdbStellarParams represents parameters for minting USDB on Stellar.
type MintUsdbStellarParams struct {
	Address   string `json:"address"`
	Amount    string `json:"amount"`
	SignedXDR string `json:"signedXdr,omitempty"`
}

// MintUsdbSolanaParams represents parameters for minting USDB on Solana.
type MintUsdbSolanaParams struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
}

// MintUsdbSolanaResponse represents the response from minting USDB on Solana.
type MintUsdbSolanaResponse struct {
	Success   bool   `json:"success"`
	Signature string `json:"signature"`
	Error     string `json:"error"`
}

// PrepareSolanaDelegationTransactionParams represents parameters for preparing a Solana delegation transaction.
type PrepareSolanaDelegationTransactionParams struct {
	// QuoteID (recommended) makes the API take the amount and token from
	// the quote; otherwise set TokenAddress and Amount.
	QuoteID      string `json:"quote_id,omitempty"`
	TokenAddress string `json:"token_address,omitempty"`
	Amount       string `json:"amount,omitempty"`
	OwnerAddress string `json:"owner_address"`
}

// PrepareSolanaDelegationTransactionResponse represents the response from preparing a Solana delegation transaction.
type PrepareSolanaDelegationTransactionResponse struct {
	Success     bool   `json:"success"`
	Transaction string `json:"transaction"`
}

// Client handles blockchain wallet-related operations.
type Client struct {
	cfg        *request.Config
	instanceID string
}

// NewClient creates a new wallets client.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:        cfg.ToRequestConfig(),
		instanceID: cfg.InstanceID,
	}
}

// List retrieves all blockchain wallets for a customer.
func (c *Client) List(ctx context.Context, customerID string, opts ...option.RequestOption) ([]BlockchainWallet, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/blockchain-wallets", c.instanceID, customerID)
	return request.Do[[]BlockchainWallet](c.cfg, ctx, "GET", path, nil, opts...)
}

// CreateWithAddress creates a new blockchain wallet with an address.
func (c *Client) CreateWithAddress(ctx context.Context, params *CreateWithAddressParams, opts ...option.RequestOption) (*BlockchainWallet, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.CustomerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/blockchain-wallets", c.instanceID, params.CustomerID)

	body, err := request.BodyWith(params, map[string]any{"is_account_abstraction": true})
	if err != nil {
		return nil, err
	}

	return request.Do[*BlockchainWallet](c.cfg, ctx, "POST", path, body, opts...)
}

// CreateWithHash creates a new blockchain wallet with a signature transaction hash.
func (c *Client) CreateWithHash(ctx context.Context, params *CreateWithHashParams, opts ...option.RequestOption) (*BlockchainWallet, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.CustomerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/blockchain-wallets", c.instanceID, params.CustomerID)

	body, err := request.BodyWith(params, map[string]any{"is_account_abstraction": false})
	if err != nil {
		return nil, err
	}

	return request.Do[*BlockchainWallet](c.cfg, ctx, "POST", path, body, opts...)
}

// GetWalletMessage retrieves the wallet message for signing.
func (c *Client) GetWalletMessage(ctx context.Context, customerID string, opts ...option.RequestOption) (*GetMessageResponse, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/blockchain-wallets/sign-message", c.instanceID, customerID)
	return request.Do[*GetMessageResponse](c.cfg, ctx, "GET", path, nil, opts...)
}

// Get retrieves a specific blockchain wallet.
func (c *Client) Get(ctx context.Context, customerID, id string, opts ...option.RequestOption) (*BlockchainWallet, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/blockchain-wallets/%s", c.instanceID, customerID, id)
	return request.Do[*BlockchainWallet](c.cfg, ctx, "GET", path, nil, opts...)
}

// Delete deletes a blockchain wallet.
func (c *Client) Delete(ctx context.Context, customerID, id string, opts ...option.RequestOption) error {
	if customerID == "" {
		return fmt.Errorf("customer ID cannot be empty")
	}
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/blockchain-wallets/%s", c.instanceID, customerID, id)
	_, err := request.Do[struct{}](c.cfg, ctx, "DELETE", path, nil, opts...)
	return err
}

// CreateAssetTrustline creates an asset trustline on Stellar.
func (c *Client) CreateAssetTrustline(ctx context.Context, address string, opts ...option.RequestOption) (*CreateAssetTrustlineResponse, error) {
	if address == "" {
		return nil, fmt.Errorf("address cannot be empty")
	}

	path := request.Pathf("/instances/%s/create-asset-trustline", c.instanceID)
	body := struct {
		Address string `json:"address"`
	}{
		Address: address,
	}

	return request.Do[*CreateAssetTrustlineResponse](c.cfg, ctx, "POST", path, body, opts...)
}

// MintUsdbStellar mints USDB on Stellar.
func (c *Client) MintUsdbStellar(ctx context.Context, params *MintUsdbStellarParams, opts ...option.RequestOption) error {
	if params == nil {
		return fmt.Errorf("params cannot be nil")
	}

	path := request.Pathf("/instances/%s/mint-usdb-stellar", c.instanceID)
	_, err := request.Do[struct{}](c.cfg, ctx, "POST", path, params, opts...)
	return err
}

// MintUsdbSolana mints USDB on Solana.
func (c *Client) MintUsdbSolana(ctx context.Context, params *MintUsdbSolanaParams, opts ...option.RequestOption) (*MintUsdbSolanaResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := request.Pathf("/instances/%s/mint-usdb-solana", c.instanceID)
	return request.Do[*MintUsdbSolanaResponse](c.cfg, ctx, "POST", path, params, opts...)
}

// PrepareSolanaDelegationTransaction prepares a Solana delegation transaction.
func (c *Client) PrepareSolanaDelegationTransaction(ctx context.Context, params *PrepareSolanaDelegationTransactionParams, opts ...option.RequestOption) (*PrepareSolanaDelegationTransactionResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := request.Pathf("/instances/%s/prepare-delegate-solana", c.instanceID)
	return request.Do[*PrepareSolanaDelegationTransactionResponse](c.cfg, ctx, "POST", path, params, opts...)
}

// OfframpWallet represents an offramp wallet.
type OfframpWallet struct {
	ID            string    `json:"id"`
	ExternalID    string    `json:"external_id"`
	InstanceID    string    `json:"instance_id"`
	CustomerID    string    `json:"customer_id"`
	BankAccountID string    `json:"bank_account_id"`
	Network       string    `json:"network"`
	Address       string    `json:"address"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateOfframpWalletParams represents parameters for creating an offramp wallet.
type CreateOfframpWalletParams struct {
	CustomerID    string `json:"-"`
	BankAccountID string `json:"-"`
	ExternalID    string `json:"external_id"`
	Network       string `json:"network"`
	// PartnerFeeID charges this partner fee on the wallet's offramps.
	PartnerFeeID *string `json:"partner_fee_id,omitempty"`
}

// CreateOfframpWalletResponse represents the response when creating an offramp wallet.
type CreateOfframpWalletResponse struct {
	ID         string `json:"id"`
	ExternalID string `json:"external_id"`
	Network    string `json:"network"`
	Address    string `json:"address"`
}

// OfframpClient handles offramp wallet-related operations.
type OfframpClient struct {
	cfg        *request.Config
	instanceID string
}

// NewOfframpClient creates a new offramp wallets client.
func NewOfframpClient(cfg *config.Config) *OfframpClient {
	return &OfframpClient{
		cfg:        cfg.ToRequestConfig(),
		instanceID: cfg.InstanceID,
	}
}

// List retrieves all offramp wallets for a bank account.
func (c *OfframpClient) List(ctx context.Context, customerID, bankAccountID string, opts ...option.RequestOption) ([]OfframpWallet, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}
	if bankAccountID == "" {
		return nil, fmt.Errorf("bank account ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/bank-accounts/%s/offramp-wallets",
		c.instanceID, customerID, bankAccountID)
	return request.Do[[]OfframpWallet](c.cfg, ctx, "GET", path, nil, opts...)
}

// Create creates a new offramp wallet.
func (c *OfframpClient) Create(ctx context.Context, params *CreateOfframpWalletParams, opts ...option.RequestOption) (*CreateOfframpWalletResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.CustomerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}
	if params.BankAccountID == "" {
		return nil, fmt.Errorf("bank account ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/bank-accounts/%s/offramp-wallets",
		c.instanceID, params.CustomerID, params.BankAccountID)

	return request.Do[*CreateOfframpWalletResponse](c.cfg, ctx, "POST", path, params, opts...)
}

// Get retrieves a specific offramp wallet.
func (c *OfframpClient) Get(ctx context.Context, customerID, bankAccountID, id string, opts ...option.RequestOption) (*OfframpWallet, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}
	if bankAccountID == "" {
		return nil, fmt.Errorf("bank account ID cannot be empty")
	}
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/bank-accounts/%s/offramp-wallets/%s",
		c.instanceID, customerID, bankAccountID, id)
	return request.Do[*OfframpWallet](c.cfg, ctx, "GET", path, nil, opts...)
}
