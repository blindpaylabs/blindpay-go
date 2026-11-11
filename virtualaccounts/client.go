package virtualaccounts

import (
	"context"
	"fmt"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/blindpaylabs/blindpay-go/option"
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

// VirtualAccountBeneficiary is the name shown to senders.
type VirtualAccountBeneficiary struct {
	Name string `json:"name"`
}

// VirtualAccountMX holds a Mexican virtual account's SPEI details.
type VirtualAccountMX struct {
	Spei struct {
		Clabe string `json:"clabe"`
	} `json:"spei"`
	Beneficiary *VirtualAccountBeneficiary `json:"beneficiary,omitempty"`
}

// VirtualAccountCO holds a Colombian virtual account's Bre-B details.
type VirtualAccountCO struct {
	BreB struct {
		Key string `json:"key"`
	} `json:"bre_b"`
	Beneficiary *VirtualAccountBeneficiary `json:"beneficiary,omitempty"`
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
	CustomerID         *string               `json:"customer_id,omitempty"`
	// MX is set for Mexican (SPEI) accounts and CO for Colombian (Bre-B)
	// accounts; US is then absent.
	MX           *VirtualAccountMX `json:"mx,omitempty"`
	CO           *VirtualAccountCO `json:"co,omitempty"`
	PartnerFeeID *string           `json:"partner_fee_id,omitempty"`
}

// CreateParams represents parameters for creating a virtual account.
type CreateParams struct {
	CustomerID            string                `json:"-"`
	BlockchainWalletID    string                `json:"blockchain_wallet_id"`
	Token                 types.StablecoinToken `json:"token"`
	BankingPartner        types.BankingPartner  `json:"banking_partner"`
	SoleProprietorDocType *string               `json:"sole_proprietor_doc_type,omitempty"`
	SoleProprietorDocFile *string               `json:"sole_proprietor_doc_file,omitempty"`
	// PartnerFeeID charges this partner fee on the account's payins.
	PartnerFeeID *string `json:"partner_fee_id,omitempty"`
	// BreB is required for Colombian (bre_b) accounts.
	BreB *BreBParams `json:"bre_b,omitempty"`
}

// BreBGender is the account holder's gender for a Bre-B account.
type BreBGender string

const (
	BreBGenderMale   BreBGender = "male"
	BreBGenderFemale BreBGender = "female"
)

// BreBParams holds the Bre-B (Colombia) details of a virtual account.
type BreBParams struct {
	Gender BreBGender `json:"gender"`
}

// UpdateParams represents parameters for updating a virtual account.
type UpdateParams struct {
	CustomerID         string                `json:"-"`
	VirtualAccountID   string                `json:"-"`
	BlockchainWalletID string                `json:"blockchain_wallet_id"`
	Token              types.StablecoinToken `json:"token"`
	// PartnerFeeID charges this partner fee on the account's payins.
	PartnerFeeID *string `json:"partner_fee_id,omitempty"`
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
func (c *Client) Create(ctx context.Context, params *CreateParams, opts ...option.RequestOption) (*VirtualAccount, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.CustomerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/virtual-accounts", c.instanceID, params.CustomerID)

	return request.Do[*VirtualAccount](c.cfg, ctx, "POST", path, params, opts...)
}

// Get retrieves a virtual account by ID.
func (c *Client) Get(ctx context.Context, customerID, virtualAccountID string, opts ...option.RequestOption) (*VirtualAccount, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}
	if virtualAccountID == "" {
		return nil, fmt.Errorf("virtual account ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/virtual-accounts/%s", c.instanceID, customerID, virtualAccountID)
	return request.Do[*VirtualAccount](c.cfg, ctx, "GET", path, nil, opts...)
}

// List retrieves all virtual accounts for a customer.
func (c *Client) List(ctx context.Context, customerID string, opts ...option.RequestOption) ([]VirtualAccount, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/virtual-accounts", c.instanceID, customerID)
	return request.Do[[]VirtualAccount](c.cfg, ctx, "GET", path, nil, opts...)
}

// Update updates a virtual account.
func (c *Client) Update(ctx context.Context, params *UpdateParams, opts ...option.RequestOption) error {
	if params == nil {
		return fmt.Errorf("params cannot be nil")
	}
	if params.CustomerID == "" {
		return fmt.Errorf("customer ID cannot be empty")
	}
	if params.VirtualAccountID == "" {
		return fmt.Errorf("virtual account ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/virtual-accounts/%s", c.instanceID, params.CustomerID, params.VirtualAccountID)

	_, err := request.Do[struct{}](c.cfg, ctx, "PUT", path, params, opts...)
	return err
}
