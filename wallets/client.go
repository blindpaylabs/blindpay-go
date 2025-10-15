package wallets

import (
	"context"
	"fmt"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
)

// BlockchainWallet represents a blockchain wallet.
type BlockchainWallet struct {
	ID                   string        `json:"id"`
	Name                 string        `json:"name"`
	Network              types.Network `json:"network"`
	Address              string        `json:"address,omitempty"`
	SignatureTxHash      string        `json:"signature_tx_hash,omitempty"`
	IsAccountAbstraction bool          `json:"is_account_abstraction"`
	ReceiverID           string        `json:"receiver_id"`
}

// GetMessageResponse represents the wallet message response.
type GetMessageResponse struct {
	Message string `json:"message"`
}

// CreateWithAddressParams represents parameters for creating a wallet with address.
type CreateWithAddressParams struct {
	ReceiverID string        `json:"receiver_id"`
	Name       string        `json:"name"`
	Network    types.Network `json:"network"`
	Address    string        `json:"address"`
}

// CreateWithHashParams represents parameters for creating a wallet with hash.
type CreateWithHashParams struct {
	ReceiverID      string        `json:"receiver_id"`
	Name            string        `json:"name"`
	Network         types.Network `json:"network"`
	SignatureTxHash string        `json:"signature_tx_hash"`
}

// CreateAssetTrustlineResponse represents the asset trustline response.
type CreateAssetTrustlineResponse struct {
	XDR string `json:"xdr"`
}

// MintUsdbStellarParams represents parameters for minting USDB on Stellar.
type MintUsdbStellarParams struct {
	Address   string `json:"address"`
	Amount    string `json:"amount"`
	SignedXDR string `json:"signedXdr"`
}

// Client handles blockchain wallet-related operations.
type Client struct {
	cfg        *request.Config
	instanceID string
}

// NewClient creates a new wallets client.
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

// List retrieves all blockchain wallets for a receiver.
func (c *Client) List(ctx context.Context, receiverID string) ([]BlockchainWallet, error) {
	if receiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/blockchain-wallets", c.instanceID, receiverID)
	return request.Do[[]BlockchainWallet](c.cfg, ctx, "GET", path, nil)
}

// CreateWithAddress creates a new blockchain wallet with an address.
func (c *Client) CreateWithAddress(ctx context.Context, params *CreateWithAddressParams) (*BlockchainWallet, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.ReceiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/blockchain-wallets", c.instanceID, params.ReceiverID)

	body := struct {
		Name                 string        `json:"name"`
		Network              types.Network `json:"network"`
		Address              string        `json:"address"`
		IsAccountAbstraction bool          `json:"is_account_abstraction"`
	}{
		Name:                 params.Name,
		Network:              params.Network,
		Address:              params.Address,
		IsAccountAbstraction: true,
	}

	return request.Do[*BlockchainWallet](c.cfg, ctx, "POST", path, body)
}

// CreateWithHash creates a new blockchain wallet with a signature transaction hash.
func (c *Client) CreateWithHash(ctx context.Context, params *CreateWithHashParams) (*BlockchainWallet, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.ReceiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/blockchain-wallets", c.instanceID, params.ReceiverID)

	body := struct {
		Name                 string        `json:"name"`
		Network              types.Network `json:"network"`
		SignatureTxHash      string        `json:"signature_tx_hash"`
		IsAccountAbstraction bool          `json:"is_account_abstraction"`
	}{
		Name:                 params.Name,
		Network:              params.Network,
		SignatureTxHash:      params.SignatureTxHash,
		IsAccountAbstraction: false,
	}

	return request.Do[*BlockchainWallet](c.cfg, ctx, "POST", path, body)
}

// GetWalletMessage retrieves the wallet message for signing.
func (c *Client) GetWalletMessage(ctx context.Context, receiverID string) (*GetMessageResponse, error) {
	if receiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/blockchain-wallets/sign-message", c.instanceID, receiverID)
	return request.Do[*GetMessageResponse](c.cfg, ctx, "GET", path, nil)
}

// Get retrieves a specific blockchain wallet.
func (c *Client) Get(ctx context.Context, receiverID, id string) (*BlockchainWallet, error) {
	if receiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/blockchain-wallets/%s", c.instanceID, receiverID, id)
	return request.Do[*BlockchainWallet](c.cfg, ctx, "GET", path, nil)
}

// Delete deletes a blockchain wallet.
func (c *Client) Delete(ctx context.Context, receiverID, id string) error {
	if receiverID == "" {
		return fmt.Errorf("receiver ID cannot be empty")
	}
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/blockchain-wallets/%s", c.instanceID, receiverID, id)
	_, err := request.Do[struct{}](c.cfg, ctx, "DELETE", path, nil)
	return err
}

// CreateAssetTrustline creates an asset trustline on Stellar.
func (c *Client) CreateAssetTrustline(ctx context.Context, address string) (*CreateAssetTrustlineResponse, error) {
	if address == "" {
		return nil, fmt.Errorf("address cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/create-asset-trustline", c.instanceID)
	body := struct {
		Address string `json:"address"`
	}{
		Address: address,
	}

	return request.Do[*CreateAssetTrustlineResponse](c.cfg, ctx, "POST", path, body)
}

// MintUsdbStellar mints USDB on Stellar.
func (c *Client) MintUsdbStellar(ctx context.Context, params *MintUsdbStellarParams) error {
	if params == nil {
		return fmt.Errorf("params cannot be nil")
	}

	path := fmt.Sprintf("/instances/%s/mint-usdb-stellar", c.instanceID)
	_, err := request.Do[struct{}](c.cfg, ctx, "POST", path, params)
	return err
}

// OfframpWallet represents an offramp wallet.
type OfframpWallet struct {
	ID            string    `json:"id"`
	ExternalID    string    `json:"external_id"`
	InstanceID    string    `json:"instance_id"`
	ReceiverID    string    `json:"receiver_id"`
	BankAccountID string    `json:"bank_account_id"`
	Network       string    `json:"network"`
	Address       string    `json:"address"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateOfframpWalletParams represents parameters for creating an offramp wallet.
type CreateOfframpWalletParams struct {
	ReceiverID    string `json:"receiver_id"`
	BankAccountID string `json:"bank_account_id"`
	ExternalID    string `json:"external_id"`
	Network       string `json:"network"`
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
		cfg: &request.Config{
			BaseURL:    cfg.BaseURL,
			APIKey:     cfg.APIKey,
			HTTPClient: cfg.HTTPClient,
			UserAgent:  cfg.UserAgent,
		},
		instanceID: cfg.InstanceID,
	}
}

// List retrieves all offramp wallets for a bank account.
func (c *OfframpClient) List(ctx context.Context, receiverID, bankAccountID string) ([]OfframpWallet, error) {
	if receiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}
	if bankAccountID == "" {
		return nil, fmt.Errorf("bank account ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts/%s/offramp-wallets",
		c.instanceID, receiverID, bankAccountID)
	return request.Do[[]OfframpWallet](c.cfg, ctx, "GET", path, nil)
}

// Create creates a new offramp wallet.
func (c *OfframpClient) Create(ctx context.Context, params *CreateOfframpWalletParams) (*CreateOfframpWalletResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.ReceiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}
	if params.BankAccountID == "" {
		return nil, fmt.Errorf("bank account ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts/%s/offramp-wallets",
		c.instanceID, params.ReceiverID, params.BankAccountID)

	body := struct {
		ExternalID string `json:"external_id"`
		Network    string `json:"network"`
	}{
		ExternalID: params.ExternalID,
		Network:    params.Network,
	}

	return request.Do[*CreateOfframpWalletResponse](c.cfg, ctx, "POST", path, body)
}

// Get retrieves a specific offramp wallet.
func (c *OfframpClient) Get(ctx context.Context, receiverID, bankAccountID, id string) (*OfframpWallet, error) {
	if receiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}
	if bankAccountID == "" {
		return nil, fmt.Errorf("bank account ID cannot be empty")
	}
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts/%s/offramp-wallets/%s",
		c.instanceID, receiverID, bankAccountID, id)
	return request.Do[*OfframpWallet](c.cfg, ctx, "GET", path, nil)
}
