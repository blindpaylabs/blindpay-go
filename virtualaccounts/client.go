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
	BankingPartner     *types.BankingPartner `json:"banking_partner,omitempty"`
	KycStatus          *string               `json:"kyc_status,omitempty"`
	BlockchainWallet   interface{}           `json:"blockchain_wallet,omitempty"`
}

// CreateParams represents parameters for creating a virtual account.
type CreateParams struct {
	ReceiverID            string                `json:"-"`
	BlockchainWalletID    string                `json:"blockchain_wallet_id"`
	Token                 types.StablecoinToken `json:"token"`
	BankingPartner        types.BankingPartner  `json:"banking_partner"`
	SoleProprietorDocType *string               `json:"sole_proprietor_doc_type,omitempty"`
	SoleProprietorDocFile *string               `json:"sole_proprietor_doc_file,omitempty"`
}

// UpdateParams represents parameters for updating a virtual account.
type UpdateParams struct {
	ReceiverID         string                `json:"-"`
	VirtualAccountID   string                `json:"-"`
	BlockchainWalletID string                `json:"blockchain_wallet_id"`
	Token              types.StablecoinToken `json:"token"`
}

// ListResponse represents the response when listing virtual accounts.
type ListResponse struct {
	Data       []VirtualAccount         `json:"data"`
	Pagination types.PaginationMetadata `json:"pagination"`
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
		BlockchainWalletID    string                `json:"blockchain_wallet_id"`
		Token                 types.StablecoinToken `json:"token"`
		BankingPartner        types.BankingPartner  `json:"banking_partner"`
		SoleProprietorDocType *string               `json:"sole_proprietor_doc_type,omitempty"`
		SoleProprietorDocFile *string               `json:"sole_proprietor_doc_file,omitempty"`
	}{
		BlockchainWalletID:    params.BlockchainWalletID,
		Token:                 params.Token,
		BankingPartner:        params.BankingPartner,
		SoleProprietorDocType: params.SoleProprietorDocType,
		SoleProprietorDocFile: params.SoleProprietorDocFile,
	}

	return request.Do[*VirtualAccount](c.cfg, ctx, "POST", path, body)
}

// Get retrieves a virtual account by ID.
func (c *Client) Get(ctx context.Context, receiverID, virtualAccountID string) (*VirtualAccount, error) {
	if receiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}
	if virtualAccountID == "" {
		return nil, fmt.Errorf("virtual account ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/virtual-accounts/%s", c.instanceID, receiverID, virtualAccountID)
	return request.Do[*VirtualAccount](c.cfg, ctx, "GET", path, nil)
}

// List retrieves all virtual accounts for a receiver.
func (c *Client) List(ctx context.Context, receiverID string) ([]VirtualAccount, error) {
	if receiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/virtual-accounts", c.instanceID, receiverID)
	return request.Do[[]VirtualAccount](c.cfg, ctx, "GET", path, nil)
}

// Update updates a virtual account.
func (c *Client) Update(ctx context.Context, params *UpdateParams) error {
	if params == nil {
		return fmt.Errorf("params cannot be nil")
	}
	if params.ReceiverID == "" {
		return fmt.Errorf("receiver ID cannot be empty")
	}
	if params.VirtualAccountID == "" {
		return fmt.Errorf("virtual account ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/virtual-accounts/%s", c.instanceID, params.ReceiverID, params.VirtualAccountID)

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
