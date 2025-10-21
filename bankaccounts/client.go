package bankaccounts

import (
	"context"
	"fmt"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
)

// ArgentinaTransfers represents Argentina transfer types.
type ArgentinaTransfers string

const (
	ArgentinaTransfersCVU   ArgentinaTransfers = "CVU"
	ArgentinaTransfersCBU   ArgentinaTransfers = "CBU"
	ArgentinaTransfersALIAS ArgentinaTransfers = "ALIAS"
)

// AchCopDocument represents Colombia ACH document types.
type AchCopDocument string

const (
	AchCopDocumentCC   AchCopDocument = "CC"
	AchCopDocumentCE   AchCopDocument = "CE"
	AchCopDocumentNIT  AchCopDocument = "NIT"
	AchCopDocumentPASS AchCopDocument = "PASS"
	AchCopDocumentPEP  AchCopDocument = "PEP"
)

// SpeiProtocol represents SPEI protocol types.
type SpeiProtocol string

const (
	SpeiProtocolClabe     SpeiProtocol = "clabe"
	SpeiProtocolDebitcard SpeiProtocol = "debitcard"
	SpeiProtocolPhonenum  SpeiProtocol = "phonenum"
)

// BankAccount represents a bank account with all possible fields.
type BankAccount struct {
	ID                                     string                `json:"id"`
	Type                                   types.Rail            `json:"type"`
	Name                                   string                `json:"name"`
	PixKey                                 string                `json:"pix_key,omitempty"`
	BeneficiaryName                        string                `json:"beneficiary_name,omitempty"`
	RoutingNumber                          string                `json:"routing_number,omitempty"`
	AccountNumber                          string                `json:"account_number,omitempty"`
	AccountType                            types.BankAccountType `json:"account_type,omitempty"`
	AccountClass                           types.AccountClass    `json:"account_class,omitempty"`
	AddressLine1                           string                `json:"address_line_1,omitempty"`
	AddressLine2                           string                `json:"address_line_2,omitempty"`
	City                                   string                `json:"city,omitempty"`
	StateProvinceRegion                    string                `json:"state_province_region,omitempty"`
	Country                                types.Country         `json:"country,omitempty"`
	PostalCode                             string                `json:"postal_code,omitempty"`
	SpeiProtocol                           string                `json:"spei_protocol,omitempty"`
	SpeiInstitutionCode                    string                `json:"spei_institution_code,omitempty"`
	SpeiClabe                              string                `json:"spei_clabe,omitempty"`
	TransfersType                          ArgentinaTransfers    `json:"transfers_type,omitempty"`
	TransfersAccount                       string                `json:"transfers_account,omitempty"`
	AchCopBeneficiaryFirstName             string                `json:"ach_cop_beneficiary_first_name,omitempty"`
	AchCopBeneficiaryLastName              string                `json:"ach_cop_beneficiary_last_name,omitempty"`
	AchCopDocumentID                       string                `json:"ach_cop_document_id,omitempty"`
	AchCopDocumentType                     AchCopDocument        `json:"ach_cop_document_type,omitempty"`
	AchCopEmail                            string                `json:"ach_cop_email,omitempty"`
	AchCopBankCode                         string                `json:"ach_cop_bank_code,omitempty"`
	AchCopBankAccount                      string                `json:"ach_cop_bank_account,omitempty"`
	SwiftCodeBic                           string                `json:"swift_code_bic,omitempty"`
	SwiftAccountHolderName                 string                `json:"swift_account_holder_name,omitempty"`
	SwiftAccountNumberIban                 string                `json:"swift_account_number_iban,omitempty"`
	SwiftBeneficiaryAddressLine1           string                `json:"swift_beneficiary_address_line_1,omitempty"`
	SwiftBeneficiaryAddressLine2           string                `json:"swift_beneficiary_address_line_2,omitempty"`
	SwiftBeneficiaryCountry                types.Country         `json:"swift_beneficiary_country,omitempty"`
	SwiftBeneficiaryCity                   string                `json:"swift_beneficiary_city,omitempty"`
	SwiftBeneficiaryStateProvinceRegion    string                `json:"swift_beneficiary_state_province_region,omitempty"`
	SwiftBeneficiaryPostalCode             string                `json:"swift_beneficiary_postal_code,omitempty"`
	SwiftBankName                          string                `json:"swift_bank_name,omitempty"`
	SwiftBankAddressLine1                  string                `json:"swift_bank_address_line_1,omitempty"`
	SwiftBankAddressLine2                  string                `json:"swift_bank_address_line_2,omitempty"`
	SwiftBankCountry                       types.Country         `json:"swift_bank_country,omitempty"`
	SwiftBankCity                          string                `json:"swift_bank_city,omitempty"`
	SwiftBankStateProvinceRegion           string                `json:"swift_bank_state_province_region,omitempty"`
	SwiftBankPostalCode                    string                `json:"swift_bank_postal_code,omitempty"`
	SwiftIntermediaryBankSwiftCodeBic      string                `json:"swift_intermediary_bank_swift_code_bic,omitempty"`
	SwiftIntermediaryBankAccountNumberIban string                `json:"swift_intermediary_bank_account_number_iban,omitempty"`
	SwiftIntermediaryBankName              string                `json:"swift_intermediary_bank_name,omitempty"`
	SwiftIntermediaryBankCountry           types.Country         `json:"swift_intermediary_bank_country,omitempty"`
	TronWalletHash                         string                `json:"tron_wallet_hash,omitempty"`
	OfframpWallets                         []OfframpWalletInfo   `json:"offramp_wallets,omitempty"`
	CreatedAt                              time.Time             `json:"created_at"`
}

// OfframpWalletInfo represents offramp wallet information.
type OfframpWalletInfo struct {
	Address    string `json:"address"`
	ID         string `json:"id"`
	Network    string `json:"network"`
	ExternalID string `json:"external_id"`
}

// GetResponse represents a detailed bank account response.
type GetResponse struct {
	ID                string                `json:"id"`
	ReceiverID        string                `json:"receiver_id"`
	AccountHolderName string                `json:"account_holder_name"`
	AccountNumber     string                `json:"account_number"`
	RoutingNumber     string                `json:"routing_number"`
	AccountType       types.BankAccountType `json:"account_type"`
	BankName          string                `json:"bank_name"`
	SwiftCode         *string               `json:"swift_code"`
	IBAN              *string               `json:"iban"`
	IsPrimary         bool                  `json:"is_primary"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
}

// CreatePixParams represents parameters for creating a PIX bank account.
type CreatePixParams struct {
	ReceiverID string `json:"-"`
	Name       string `json:"name"`
	PixKey     string `json:"pix_key"`
}

// CreatePixResponse represents the response when creating a PIX bank account.
type CreatePixResponse struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	PixKey    string    `json:"pix_key"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateAchParams represents parameters for creating an ACH bank account.
type CreateAchParams struct {
	ReceiverID      string                `json:"-"`
	Name            string                `json:"name"`
	AccountClass    types.AccountClass    `json:"account_class"`
	AccountNumber   string                `json:"account_number"`
	AccountType     types.BankAccountType `json:"account_type"`
	BeneficiaryName string                `json:"beneficiary_name"`
	RoutingNumber   string                `json:"routing_number"`
}

// CreateWireParams represents parameters for creating a Wire bank account.
type CreateWireParams struct {
	ReceiverID          string        `json:"-"`
	Name                string        `json:"name"`
	AccountNumber       string        `json:"account_number"`
	BeneficiaryName     string        `json:"beneficiary_name"`
	RoutingNumber       string        `json:"routing_number"`
	AddressLine1        string        `json:"address_line_1"`
	AddressLine2        string        `json:"address_line_2,omitempty"`
	City                string        `json:"city"`
	StateProvinceRegion string        `json:"state_province_region"`
	Country             types.Country `json:"country"`
	PostalCode          string        `json:"postal_code"`
}

// CreateArgentinaTransfersParams represents parameters for creating an Argentina transfers bank account.
type CreateArgentinaTransfersParams struct {
	ReceiverID       string             `json:"-"`
	Name             string             `json:"name"`
	BeneficiaryName  string             `json:"beneficiary_name"`
	TransfersAccount string             `json:"transfers_account"`
	TransfersType    ArgentinaTransfers `json:"transfers_type"`
}

// CreateSpeiParams represents parameters for creating a SPEI bank account.
type CreateSpeiParams struct {
	ReceiverID          string       `json:"-"`
	BeneficiaryName     string       `json:"beneficiary_name"`
	Name                string       `json:"name"`
	SpeiClabe           string       `json:"spei_clabe"`
	SpeiInstitutionCode string       `json:"spei_institution_code"`
	SpeiProtocol        SpeiProtocol `json:"spei_protocol"`
}

// CreateColombiaAchParams represents parameters for creating a Colombia ACH bank account.
type CreateColombiaAchParams struct {
	ReceiverID                 string                `json:"-"`
	Name                       string                `json:"name"`
	AccountType                types.BankAccountType `json:"account_type"`
	AchCopBeneficiaryFirstName string                `json:"ach_cop_beneficiary_first_name"`
	AchCopBeneficiaryLastName  string                `json:"ach_cop_beneficiary_last_name"`
	AchCopDocumentID           string                `json:"ach_cop_document_id"`
	AchCopDocumentType         AchCopDocument        `json:"ach_cop_document_type"`
	AchCopEmail                string                `json:"ach_cop_email"`
	AchCopBankCode             string                `json:"ach_cop_bank_code"`
	AchCopBankAccount          string                `json:"ach_cop_bank_account"`
}

// CreateInternationalSwiftParams represents parameters for creating an international SWIFT bank account.
type CreateInternationalSwiftParams struct {
	ReceiverID                             string         `json:"-"`
	Name                                   string         `json:"name"`
	SwiftAccountHolderName                 string         `json:"swift_account_holder_name"`
	SwiftAccountNumberIban                 string         `json:"swift_account_number_iban"`
	SwiftBankAddressLine1                  string         `json:"swift_bank_address_line_1"`
	SwiftBankAddressLine2                  string         `json:"swift_bank_address_line_2,omitempty"`
	SwiftBankCity                          string         `json:"swift_bank_city"`
	SwiftBankCountry                       types.Country  `json:"swift_bank_country"`
	SwiftBankName                          string         `json:"swift_bank_name"`
	SwiftBankPostalCode                    string         `json:"swift_bank_postal_code"`
	SwiftBankStateProvinceRegion           string         `json:"swift_bank_state_province_region"`
	SwiftBeneficiaryAddressLine1           string         `json:"swift_beneficiary_address_line_1"`
	SwiftBeneficiaryAddressLine2           string         `json:"swift_beneficiary_address_line_2,omitempty"`
	SwiftBeneficiaryCity                   string         `json:"swift_beneficiary_city"`
	SwiftBeneficiaryCountry                types.Country  `json:"swift_beneficiary_country"`
	SwiftBeneficiaryPostalCode             string         `json:"swift_beneficiary_postal_code"`
	SwiftBeneficiaryStateProvinceRegion    string         `json:"swift_beneficiary_state_province_region"`
	SwiftCodeBic                           string         `json:"swift_code_bic"`
	SwiftIntermediaryBankAccountNumberIban *string        `json:"swift_intermediary_bank_account_number_iban"`
	SwiftIntermediaryBankCountry           *types.Country `json:"swift_intermediary_bank_country"`
	SwiftIntermediaryBankName              *string        `json:"swift_intermediary_bank_name"`
	SwiftIntermediaryBankSwiftCodeBic      *string        `json:"swift_intermediary_bank_swift_code_bic"`
}

// CreateRtpParams represents parameters for creating an RTP (Real-Time Payments) bank account.
type CreateRtpParams struct {
	ReceiverID          string        `json:"-"`
	Name                string        `json:"name"`
	BeneficiaryName     string        `json:"beneficiary_name"`
	RoutingNumber       string        `json:"routing_number"`
	AccountNumber       string        `json:"account_number"`
	AddressLine1        string        `json:"address_line_1"`
	AddressLine2        string        `json:"address_line_2,omitempty"`
	City                string        `json:"city"`
	StateProvinceRegion string        `json:"state_province_region"`
	Country             types.Country `json:"country"`
	PostalCode          string        `json:"postal_code"`
}

// Client handles bank account-related operations.
type Client struct {
	cfg        *request.Config
	instanceID string
}

// NewClient creates a new bank accounts client.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:        cfg.ToRequestConfig(),
		instanceID: cfg.InstanceID,
	}
}

// List retrieves all bank accounts for a receiver.
func (c *Client) List(ctx context.Context, receiverID string) ([]BankAccount, error) {
	if receiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", c.instanceID, receiverID)
	return request.Do[[]BankAccount](c.cfg, ctx, "GET", path, nil)
}

// Get retrieves a specific bank account.
func (c *Client) Get(ctx context.Context, receiverID, id string) (*GetResponse, error) {
	if receiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts/%s", c.instanceID, receiverID, id)
	return request.Do[*GetResponse](c.cfg, ctx, "GET", path, nil)
}

// Delete deletes a bank account.
func (c *Client) Delete(ctx context.Context, receiverID, id string) error {
	if receiverID == "" {
		return fmt.Errorf("receiver ID cannot be empty")
	}
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts/%s", c.instanceID, receiverID, id)
	_, err := request.Do[struct{}](c.cfg, ctx, "DELETE", path, nil)
	return err
}

// CreatePix creates a PIX bank account.
func (c *Client) CreatePix(ctx context.Context, params *CreatePixParams) (*CreatePixResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.ReceiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", c.instanceID, params.ReceiverID)

	body := map[string]any{
		"type":    "pix",
		"name":    params.Name,
		"pix_key": params.PixKey,
	}

	return request.Do[*CreatePixResponse](c.cfg, ctx, "POST", path, body)
}

// CreateAch creates an ACH bank account.
func (c *Client) CreateAch(ctx context.Context, params *CreateAchParams) (*BankAccount, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.ReceiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", c.instanceID, params.ReceiverID)

	body := map[string]any{
		"type":             "ach",
		"name":             params.Name,
		"account_class":    params.AccountClass,
		"account_number":   params.AccountNumber,
		"account_type":     params.AccountType,
		"beneficiary_name": params.BeneficiaryName,
		"routing_number":   params.RoutingNumber,
	}

	return request.Do[*BankAccount](c.cfg, ctx, "POST", path, body)
}

// CreateWire creates a Wire bank account.
func (c *Client) CreateWire(ctx context.Context, params *CreateWireParams) (*BankAccount, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.ReceiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", c.instanceID, params.ReceiverID)

	body := map[string]any{
		"type":                  "wire",
		"name":                  params.Name,
		"account_number":        params.AccountNumber,
		"beneficiary_name":      params.BeneficiaryName,
		"routing_number":        params.RoutingNumber,
		"address_line_1":        params.AddressLine1,
		"city":                  params.City,
		"state_province_region": params.StateProvinceRegion,
		"country":               params.Country,
		"postal_code":           params.PostalCode,
	}

	if params.AddressLine2 != "" {
		body["address_line_2"] = params.AddressLine2
	}

	return request.Do[*BankAccount](c.cfg, ctx, "POST", path, body)
}

// CreateArgentinaTransfers creates an Argentina transfers bank account.
func (c *Client) CreateArgentinaTransfers(ctx context.Context, params *CreateArgentinaTransfersParams) (*BankAccount, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.ReceiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", c.instanceID, params.ReceiverID)

	body := map[string]any{
		"type":              "transfers_bitso",
		"name":              params.Name,
		"beneficiary_name":  params.BeneficiaryName,
		"transfers_account": params.TransfersAccount,
		"transfers_type":    params.TransfersType,
	}

	return request.Do[*BankAccount](c.cfg, ctx, "POST", path, body)
}

// CreateSpei creates a SPEI bank account.
func (c *Client) CreateSpei(ctx context.Context, params *CreateSpeiParams) (*BankAccount, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.ReceiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", c.instanceID, params.ReceiverID)

	body := map[string]any{
		"type":                  "spei_bitso",
		"beneficiary_name":      params.BeneficiaryName,
		"name":                  params.Name,
		"spei_clabe":            params.SpeiClabe,
		"spei_institution_code": params.SpeiInstitutionCode,
		"spei_protocol":         params.SpeiProtocol,
	}

	return request.Do[*BankAccount](c.cfg, ctx, "POST", path, body)
}

// CreateColombiaAch creates a Colombia ACH bank account.
func (c *Client) CreateColombiaAch(ctx context.Context, params *CreateColombiaAchParams) (*BankAccount, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.ReceiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", c.instanceID, params.ReceiverID)

	body := map[string]any{
		"type":                           "ach_cop_bitso",
		"name":                           params.Name,
		"account_type":                   params.AccountType,
		"ach_cop_beneficiary_first_name": params.AchCopBeneficiaryFirstName,
		"ach_cop_beneficiary_last_name":  params.AchCopBeneficiaryLastName,
		"ach_cop_document_id":            params.AchCopDocumentID,
		"ach_cop_document_type":          params.AchCopDocumentType,
		"ach_cop_email":                  params.AchCopEmail,
		"ach_cop_bank_code":              params.AchCopBankCode,
		"ach_cop_bank_account":           params.AchCopBankAccount,
	}

	return request.Do[*BankAccount](c.cfg, ctx, "POST", path, body)
}

// CreateInternationalSwift creates an international SWIFT bank account.
func (c *Client) CreateInternationalSwift(ctx context.Context, params *CreateInternationalSwiftParams) (*BankAccount, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.ReceiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", c.instanceID, params.ReceiverID)

	body := map[string]any{
		"type":                                        "international_swift",
		"name":                                        params.Name,
		"swift_account_holder_name":                   params.SwiftAccountHolderName,
		"swift_account_number_iban":                   params.SwiftAccountNumberIban,
		"swift_bank_address_line_1":                   params.SwiftBankAddressLine1,
		"swift_bank_city":                             params.SwiftBankCity,
		"swift_bank_country":                          params.SwiftBankCountry,
		"swift_bank_name":                             params.SwiftBankName,
		"swift_bank_postal_code":                      params.SwiftBankPostalCode,
		"swift_bank_state_province_region":            params.SwiftBankStateProvinceRegion,
		"swift_beneficiary_address_line_1":            params.SwiftBeneficiaryAddressLine1,
		"swift_beneficiary_city":                      params.SwiftBeneficiaryCity,
		"swift_beneficiary_country":                   params.SwiftBeneficiaryCountry,
		"swift_beneficiary_postal_code":               params.SwiftBeneficiaryPostalCode,
		"swift_beneficiary_state_province_region":     params.SwiftBeneficiaryStateProvinceRegion,
		"swift_code_bic":                              params.SwiftCodeBic,
		"swift_intermediary_bank_account_number_iban": params.SwiftIntermediaryBankAccountNumberIban,
		"swift_intermediary_bank_country":             params.SwiftIntermediaryBankCountry,
		"swift_intermediary_bank_name":                params.SwiftIntermediaryBankName,
		"swift_intermediary_bank_swift_code_bic":      params.SwiftIntermediaryBankSwiftCodeBic,
	}

	// Add optional address_line_2 fields only if not empty
	if params.SwiftBankAddressLine2 != "" {
		body["swift_bank_address_line_2"] = params.SwiftBankAddressLine2
	}
	if params.SwiftBeneficiaryAddressLine2 != "" {
		body["swift_beneficiary_address_line_2"] = params.SwiftBeneficiaryAddressLine2
	}

	return request.Do[*BankAccount](c.cfg, ctx, "POST", path, body)
}

// CreateRtp creates an RTP (Real-Time Payments) bank account.
func (c *Client) CreateRtp(ctx context.Context, params *CreateRtpParams) (*BankAccount, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.ReceiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", c.instanceID, params.ReceiverID)

	body := map[string]any{
		"type":                  "rtp",
		"name":                  params.Name,
		"beneficiary_name":      params.BeneficiaryName,
		"routing_number":        params.RoutingNumber,
		"account_number":        params.AccountNumber,
		"address_line_1":        params.AddressLine1,
		"city":                  params.City,
		"state_province_region": params.StateProvinceRegion,
		"country":               params.Country,
		"postal_code":           params.PostalCode,
	}

	if params.AddressLine2 != "" {
		body["address_line_2"] = params.AddressLine2
	}

	return request.Do[*BankAccount](c.cfg, ctx, "POST", path, body)
}
