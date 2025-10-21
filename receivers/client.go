package receivers

import (
	"context"
	"fmt"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
)

// ProofOfAddressDocType represents proof of address document types.
type ProofOfAddressDocType string

const (
	ProofOfAddressDocTypeUtilityBill              ProofOfAddressDocType = "UTILITY_BILL"
	ProofOfAddressDocTypeBankStatement            ProofOfAddressDocType = "BANK_STATEMENT"
	ProofOfAddressDocTypeRentalAgreement          ProofOfAddressDocType = "RENTAL_AGREEMENT"
	ProofOfAddressDocTypeTaxDocument              ProofOfAddressDocType = "TAX_DOCUMENT"
	ProofOfAddressDocTypeGovernmentCorrespondence ProofOfAddressDocType = "GOVERNMENT_CORRESPONDENCE"
)

// PurposeOfTransactions represents the purpose of transactions.
type PurposeOfTransactions string

const (
	PurposeBusinessTransactions            PurposeOfTransactions = "business_transactions"
	PurposeCharitableDonations             PurposeOfTransactions = "charitable_donations"
	PurposeInvestmentPurposes              PurposeOfTransactions = "investment_purposes"
	PurposePaymentsToFriendsOrFamilyAbroad PurposeOfTransactions = "payments_to_friends_or_family_abroad"
	PurposePersonalOrLivingExpenses        PurposeOfTransactions = "personal_or_living_expenses"
	PurposeProtectWealth                   PurposeOfTransactions = "protect_wealth"
	PurposePurchaseGoodAndServices         PurposeOfTransactions = "purchase_good_and_services"
	PurposeReceivePaymentForFreelancing    PurposeOfTransactions = "receive_payment_for_freelancing"
	PurposeReceiveSalary                   PurposeOfTransactions = "receive_salary"
	PurposeOther                           PurposeOfTransactions = "other"
)

// SourceOfFundsDocType represents source of funds document types.
type SourceOfFundsDocType string

const (
	SourceOfFundsBusinessIncome         SourceOfFundsDocType = "business_income"
	SourceOfFundsGamblingProceeds       SourceOfFundsDocType = "gambling_proceeds"
	SourceOfFundsGifts                  SourceOfFundsDocType = "gifts"
	SourceOfFundsGovernmentBenefits     SourceOfFundsDocType = "government_benefits"
	SourceOfFundsInheritance            SourceOfFundsDocType = "inheritance"
	SourceOfFundsInvestmentLoans        SourceOfFundsDocType = "investment_loans"
	SourceOfFundsPensionRetirement      SourceOfFundsDocType = "pension_retirement"
	SourceOfFundsSalary                 SourceOfFundsDocType = "salary"
	SourceOfFundsSaleOfAssetsRealEstate SourceOfFundsDocType = "sale_of_assets_real_estate"
	SourceOfFundsSavings                SourceOfFundsDocType = "savings"
	SourceOfFundsESOPs                  SourceOfFundsDocType = "esops"
	SourceOfFundsInvestmentProceeds     SourceOfFundsDocType = "investment_proceeds"
	SourceOfFundsSomeoneElseFunds       SourceOfFundsDocType = "someone_else_funds"
)

// LimitIncreaseRequestStatus represents the status of a limit increase request.
type LimitIncreaseRequestStatus string

const (
	LimitIncreaseRequestStatusInReview LimitIncreaseRequestStatus = "in_review"
	LimitIncreaseRequestStatusApproved LimitIncreaseRequestStatus = "approved"
	LimitIncreaseRequestStatusRejected LimitIncreaseRequestStatus = "rejected"
)

// LimitIncreaseRequestSupportingDocumentType represents supporting document types for limit increase requests.
type LimitIncreaseRequestSupportingDocumentType string

const (
	LimitIncreaseRequestSupportingDocumentTypeIndividualBankStatement     LimitIncreaseRequestSupportingDocumentType = "individual_bank_statement"
	LimitIncreaseRequestSupportingDocumentTypeIndividualTaxReturn         LimitIncreaseRequestSupportingDocumentType = "individual_tax_return"
	LimitIncreaseRequestSupportingDocumentTypeIndividualProofOfIncome     LimitIncreaseRequestSupportingDocumentType = "individual_proof_of_income"
	LimitIncreaseRequestSupportingDocumentTypeBusinessBankStatement       LimitIncreaseRequestSupportingDocumentType = "business_bank_statement"
	LimitIncreaseRequestSupportingDocumentTypeBusinessFinancialStatements LimitIncreaseRequestSupportingDocumentType = "business_financial_statements"
	LimitIncreaseRequestSupportingDocumentTypeBusinessTaxReturn           LimitIncreaseRequestSupportingDocumentType = "business_tax_return"
)

// IdentificationDocument represents identification document types.
type IdentificationDocument string

const (
	IdentificationDocumentPassport IdentificationDocument = "PASSPORT"
	IdentificationDocumentIDCard   IdentificationDocument = "ID_CARD"
	IdentificationDocumentDrivers  IdentificationDocument = "DRIVERS"
)

// KycType represents KYC verification levels.
type KycType string

const (
	KycTypeLight    KycType = "light"
	KycTypeStandard KycType = "standard"
	KycTypeEnhanced KycType = "enhanced"
)

// OwnerRole represents the role of an owner in a business.
type OwnerRole string

const (
	OwnerRoleBeneficialControlling OwnerRole = "beneficial_controlling"
	OwnerRoleBeneficialOwner       OwnerRole = "beneficial_owner"
	OwnerRoleControllingPerson     OwnerRole = "controlling_person"
)

// Owner represents a business owner.
type Owner struct {
	ID                    string                 `json:"id,omitempty"`
	InstanceID            string                 `json:"instance_id,omitempty"`
	ReceiverID            string                 `json:"receiver_id,omitempty"`
	Role                  OwnerRole              `json:"role"`
	FirstName             string                 `json:"first_name"`
	LastName              string                 `json:"last_name"`
	DateOfBirth           string                 `json:"date_of_birth"`
	TaxID                 string                 `json:"tax_id"`
	AddressLine1          string                 `json:"address_line_1"`
	AddressLine2          *string                `json:"address_line_2"`
	City                  string                 `json:"city"`
	StateProvinceRegion   string                 `json:"state_province_region"`
	Country               types.Country          `json:"country"`
	PostalCode            string                 `json:"postal_code"`
	IDDocCountry          types.Country          `json:"id_doc_country"`
	IDDocType             IdentificationDocument `json:"id_doc_type"`
	IDDocFrontFile        string                 `json:"id_doc_front_file"`
	IDDocBackFile         *string                `json:"id_doc_back_file"`
	ProofOfAddressDocType ProofOfAddressDocType  `json:"proof_of_address_doc_type"`
	ProofOfAddressDocFile string                 `json:"proof_of_address_doc_file"`
}

// KycWarning represents a KYC warning.
type KycWarning struct {
	Code             *string `json:"code"`
	Message          *string `json:"message"`
	ResolutionStatus *string `json:"resolution_status"`
	WarningID        *string `json:"warning_id"`
}

// Limits represents transaction limits.
type Limits struct {
	PerTransaction float64 `json:"per_transaction"`
	Daily          float64 `json:"daily"`
	Monthly        float64 `json:"monthly"`
}

// Receiver represents a receiver with KYC information.
type Receiver struct {
	ID                    string                `json:"id"`
	Type                  types.AccountClass    `json:"type"`
	KycType               KycType               `json:"kyc_type"`
	KycStatus             string                `json:"kyc_status"`
	KycWarnings           []KycWarning          `json:"kyc_warnings"`
	Email                 string                `json:"email"`
	TaxID                 string                `json:"tax_id"`
	AddressLine1          string                `json:"address_line_1"`
	AddressLine2          *string               `json:"address_line_2"`
	City                  string                `json:"city"`
	StateProvinceRegion   string                `json:"state_province_region"`
	Country               types.Country         `json:"country"`
	PostalCode            string                `json:"postal_code"`
	IPAddress             *string               `json:"ip_address"`
	ImageURL              *string               `json:"image_url"`
	PhoneNumber           *string               `json:"phone_number"`
	ProofOfAddressDocType ProofOfAddressDocType `json:"proof_of_address_doc_type"`
	ProofOfAddressDocFile string                `json:"proof_of_address_doc_file"`
	AipriseValidationKey  string                `json:"aiprise_validation_key"`
	InstanceID            string                `json:"instance_id"`
	TosID                 *string               `json:"tos_id"`
	CreatedAt             time.Time             `json:"created_at"`
	UpdatedAt             time.Time             `json:"updated_at"`
	Limit                 Limits                `json:"limit"`
	// Individual fields
	FirstName      string                 `json:"first_name,omitempty"`
	LastName       string                 `json:"last_name,omitempty"`
	DateOfBirth    string                 `json:"date_of_birth,omitempty"`
	IDDocCountry   types.Country          `json:"id_doc_country,omitempty"`
	IDDocType      IdentificationDocument `json:"id_doc_type,omitempty"`
	IDDocFrontFile string                 `json:"id_doc_front_file,omitempty"`
	IDDocBackFile  *string                `json:"id_doc_back_file,omitempty"`
	// Enhanced KYC fields
	SourceOfFundsDocType             string                `json:"source_of_funds_doc_type,omitempty"`
	SourceOfFundsDocFile             string                `json:"source_of_funds_doc_file,omitempty"`
	IndividualHoldingDocFrontFile    string                `json:"individual_holding_doc_front_file,omitempty"`
	PurposeOfTransactions            PurposeOfTransactions `json:"purpose_of_transactions,omitempty"`
	PurposeOfTransactionsExplanation *string               `json:"purpose_of_transactions_explanation,omitempty"`
	// Business fields
	LegalName               string  `json:"legal_name,omitempty"`
	AlternateName           *string `json:"alternate_name,omitempty"`
	FormationDate           string  `json:"formation_date,omitempty"`
	Website                 *string `json:"website,omitempty"`
	Owners                  []Owner `json:"owners,omitempty"`
	IncorporationDocFile    string  `json:"incorporation_doc_file,omitempty"`
	ProofOfOwnershipDocFile string  `json:"proof_of_ownership_doc_file,omitempty"`
	ExternalID              *string `json:"external_id,omitempty"`
}

// CreateIndividualStandardParams represents parameters for creating an individual with standard KYC.
type CreateIndividualStandardParams struct {
	ExternalID            *string                `json:"external_id,omitempty"`
	AddressLine1          string                 `json:"address_line_1"`
	AddressLine2          *string                `json:"address_line_2,omitempty"`
	City                  string                 `json:"city"`
	Country               types.Country          `json:"country"`
	DateOfBirth           string                 `json:"date_of_birth"`
	Email                 string                 `json:"email"`
	FirstName             string                 `json:"first_name"`
	PhoneNumber           *string                `json:"phone_number"`
	IDDocCountry          types.Country          `json:"id_doc_country"`
	IDDocFrontFile        string                 `json:"id_doc_front_file"`
	IDDocType             IdentificationDocument `json:"id_doc_type"`
	IDDocBackFile         *string                `json:"id_doc_back_file,omitempty"`
	LastName              string                 `json:"last_name"`
	PostalCode            string                 `json:"postal_code"`
	ProofOfAddressDocFile string                 `json:"proof_of_address_doc_file"`
	ProofOfAddressDocType ProofOfAddressDocType  `json:"proof_of_address_doc_type"`
	StateProvinceRegion   string                 `json:"state_province_region"`
	TaxID                 string                 `json:"tax_id"`
	TosID                 string                 `json:"tos_id"`
}

// CreateIndividualEnhancedParams represents parameters for creating an individual with enhanced KYC.
type CreateIndividualEnhancedParams struct {
	ExternalID                       *string                `json:"external_id,omitempty"`
	AddressLine1                     string                 `json:"address_line_1"`
	AddressLine2                     *string                `json:"address_line_2,omitempty"`
	City                             string                 `json:"city"`
	Country                          types.Country          `json:"country"`
	DateOfBirth                      string                 `json:"date_of_birth"`
	Email                            string                 `json:"email"`
	FirstName                        string                 `json:"first_name"`
	IDDocCountry                     types.Country          `json:"id_doc_country"`
	IDDocFrontFile                   string                 `json:"id_doc_front_file"`
	IDDocType                        IdentificationDocument `json:"id_doc_type"`
	IDDocBackFile                    *string                `json:"id_doc_back_file,omitempty"`
	IndividualHoldingDocFrontFile    string                 `json:"individual_holding_doc_front_file"`
	LastName                         string                 `json:"last_name"`
	PostalCode                       string                 `json:"postal_code"`
	PhoneNumber                      *string                `json:"phone_number"`
	ProofOfAddressDocFile            string                 `json:"proof_of_address_doc_file"`
	ProofOfAddressDocType            ProofOfAddressDocType  `json:"proof_of_address_doc_type"`
	PurposeOfTransactions            PurposeOfTransactions  `json:"purpose_of_transactions"`
	SourceOfFundsDocFile             string                 `json:"source_of_funds_doc_file"`
	SourceOfFundsDocType             SourceOfFundsDocType   `json:"source_of_funds_doc_type"`
	PurposeOfTransactionsExplanation *string                `json:"purpose_of_transactions_explanation,omitempty"`
	StateProvinceRegion              string                 `json:"state_province_region"`
	TaxID                            string                 `json:"tax_id"`
	TosID                            string                 `json:"tos_id"`
}

// CreateBusinessStandardParams represents parameters for creating a business with standard KYB.
type CreateBusinessStandardParams struct {
	ExternalID              *string               `json:"external_id,omitempty"`
	AddressLine1            string                `json:"address_line_1"`
	AddressLine2            *string               `json:"address_line_2,omitempty"`
	AlternateName           string                `json:"alternate_name"`
	City                    string                `json:"city"`
	Country                 types.Country         `json:"country"`
	Email                   string                `json:"email"`
	FormationDate           string                `json:"formation_date"`
	IncorporationDocFile    string                `json:"incorporation_doc_file"`
	LegalName               string                `json:"legal_name"`
	Owners                  []Owner               `json:"owners"`
	PostalCode              string                `json:"postal_code"`
	ProofOfAddressDocFile   string                `json:"proof_of_address_doc_file"`
	ProofOfAddressDocType   ProofOfAddressDocType `json:"proof_of_address_doc_type"`
	ProofOfOwnershipDocFile string                `json:"proof_of_ownership_doc_file"`
	StateProvinceRegion     string                `json:"state_province_region"`
	TaxID                   string                `json:"tax_id"`
	TosID                   string                `json:"tos_id"`
	Website                 *string               `json:"website,omitempty"`
}

// CreateResponse represents the response when creating a receiver.
type CreateResponse struct {
	ID string `json:"id"`
}

// UpdateParams represents parameters for updating a receiver.
type UpdateParams struct {
	ReceiverID                       string                  `json:"-"`
	Email                            *string                 `json:"email,omitempty"`
	TaxID                            *string                 `json:"tax_id,omitempty"`
	AddressLine1                     *string                 `json:"address_line_1,omitempty"`
	AddressLine2                     *string                 `json:"address_line_2,omitempty"`
	City                             *string                 `json:"city,omitempty"`
	StateProvinceRegion              *string                 `json:"state_province_region,omitempty"`
	Country                          *types.Country          `json:"country,omitempty"`
	PostalCode                       *string                 `json:"postal_code,omitempty"`
	IPAddress                        *string                 `json:"ip_address,omitempty"`
	ImageURL                         *string                 `json:"image_url,omitempty"`
	PhoneNumber                      *string                 `json:"phone_number,omitempty"`
	ProofOfAddressDocType            *ProofOfAddressDocType  `json:"proof_of_address_doc_type,omitempty"`
	ProofOfAddressDocFile            *string                 `json:"proof_of_address_doc_file,omitempty"`
	FirstName                        *string                 `json:"first_name,omitempty"`
	LastName                         *string                 `json:"last_name,omitempty"`
	DateOfBirth                      *string                 `json:"date_of_birth,omitempty"`
	IDDocCountry                     *types.Country          `json:"id_doc_country,omitempty"`
	IDDocType                        *IdentificationDocument `json:"id_doc_type,omitempty"`
	IDDocFrontFile                   *string                 `json:"id_doc_front_file,omitempty"`
	IDDocBackFile                    *string                 `json:"id_doc_back_file,omitempty"`
	LegalName                        *string                 `json:"legal_name,omitempty"`
	AlternateName                    *string                 `json:"alternate_name,omitempty"`
	FormationDate                    *string                 `json:"formation_date,omitempty"`
	Website                          *string                 `json:"website,omitempty"`
	Owners                           []Owner                 `json:"owners,omitempty"`
	IncorporationDocFile             *string                 `json:"incorporation_doc_file,omitempty"`
	ProofOfOwnershipDocFile          *string                 `json:"proof_of_ownership_doc_file,omitempty"`
	SourceOfFundsDocType             *SourceOfFundsDocType   `json:"source_of_funds_doc_type,omitempty"`
	SourceOfFundsDocFile             *string                 `json:"source_of_funds_doc_file,omitempty"`
	IndividualHoldingDocFrontFile    *string                 `json:"individual_holding_doc_front_file,omitempty"`
	PurposeOfTransactions            *PurposeOfTransactions  `json:"purpose_of_transactions,omitempty"`
	PurposeOfTransactionsExplanation *string                 `json:"purpose_of_transactions_explanation,omitempty"`
	ExternalID                       *string                 `json:"external_id,omitempty"`
	TosID                            *string                 `json:"tos_id,omitempty"`
}

// LimitsResponse represents receiver limits.
type LimitsResponse struct {
	Limits struct {
		Payin struct {
			Daily   float64 `json:"daily"`
			Monthly float64 `json:"monthly"`
		} `json:"payin"`
		Payout struct {
			Daily   float64 `json:"daily"`
			Monthly float64 `json:"monthly"`
		} `json:"payout"`
	} `json:"limits"`
}

// LimitIncreaseRequest represents a limit increase request for a receiver.
type LimitIncreaseRequest struct {
	ID                     string                                     `json:"id"`
	ReceiverID             string                                     `json:"receiver_id"`
	Status                 LimitIncreaseRequestStatus                 `json:"status"`
	Daily                  float64                                    `json:"daily"`
	Monthly                float64                                    `json:"monthly"`
	PerTransaction         float64                                    `json:"per_transaction"`
	SupportingDocumentFile string                                     `json:"supporting_document_file"`
	SupportingDocumentType LimitIncreaseRequestSupportingDocumentType `json:"supporting_document_type"`
	CreatedAt              string                                     `json:"created_at"`
	UpdatedAt              string                                     `json:"updated_at"`
}

// RequestLimitIncreaseParams represents parameters for requesting a limit increase.
type RequestLimitIncreaseParams struct {
	ReceiverID             string                                     `json:"-"`
	Daily                  float64                                    `json:"daily"`
	Monthly                float64                                    `json:"monthly"`
	PerTransaction         float64                                    `json:"per_transaction"`
	SupportingDocumentFile string                                     `json:"supporting_document_file"`
	SupportingDocumentType LimitIncreaseRequestSupportingDocumentType `json:"supporting_document_type"`
}

// RequestLimitIncreaseResponse represents the response when requesting a limit increase.
type RequestLimitIncreaseResponse struct {
	ID string `json:"id"`
}

// Client handles receiver-related operations.
type Client struct {
	cfg        *request.Config
	instanceID string
}

// NewClient creates a new receivers client.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:        cfg.ToRequestConfig(),
		instanceID: cfg.InstanceID,
	}
}

// List retrieves all receivers for the instance.
func (c *Client) List(ctx context.Context) ([]Receiver, error) {
	path := fmt.Sprintf("/instances/%s/receivers", c.instanceID)
	return request.Do[[]Receiver](c.cfg, ctx, "GET", path, nil)
}

// CreateIndividualWithStandardKYC creates an individual receiver with standard KYC.
func (c *Client) CreateIndividualWithStandardKYC(ctx context.Context, params *CreateIndividualStandardParams) (*CreateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := fmt.Sprintf("/instances/%s/receivers", c.instanceID)

	body := map[string]any{
		"kyc_type":                  "standard",
		"type":                      "individual",
		"address_line_1":            params.AddressLine1,
		"city":                      params.City,
		"country":                   params.Country,
		"date_of_birth":             params.DateOfBirth,
		"email":                     params.Email,
		"first_name":                params.FirstName,
		"id_doc_country":            params.IDDocCountry,
		"id_doc_front_file":         params.IDDocFrontFile,
		"id_doc_type":               params.IDDocType,
		"last_name":                 params.LastName,
		"postal_code":               params.PostalCode,
		"proof_of_address_doc_file": params.ProofOfAddressDocFile,
		"proof_of_address_doc_type": params.ProofOfAddressDocType,
		"state_province_region":     params.StateProvinceRegion,
		"tax_id":                    params.TaxID,
		"tos_id":                    params.TosID,
	}

	if params.AddressLine2 != nil {
		body["address_line_2"] = params.AddressLine2
	}
	if params.PhoneNumber != nil {
		body["phone_number"] = params.PhoneNumber
	}
	if params.IDDocBackFile != nil {
		body["id_doc_back_file"] = params.IDDocBackFile
	}
	if params.ExternalID != nil {
		body["external_id"] = params.ExternalID
	}

	return request.Do[*CreateResponse](c.cfg, ctx, "POST", path, body)
}

// CreateIndividualWithEnhancedKYC creates an individual receiver with enhanced KYC.
func (c *Client) CreateIndividualWithEnhancedKYC(ctx context.Context, params *CreateIndividualEnhancedParams) (*CreateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := fmt.Sprintf("/instances/%s/receivers", c.instanceID)

	body := map[string]any{
		"kyc_type":                          "enhanced",
		"type":                              "individual",
		"address_line_1":                    params.AddressLine1,
		"city":                              params.City,
		"country":                           params.Country,
		"date_of_birth":                     params.DateOfBirth,
		"email":                             params.Email,
		"first_name":                        params.FirstName,
		"id_doc_country":                    params.IDDocCountry,
		"id_doc_front_file":                 params.IDDocFrontFile,
		"id_doc_type":                       params.IDDocType,
		"individual_holding_doc_front_file": params.IndividualHoldingDocFrontFile,
		"last_name":                         params.LastName,
		"postal_code":                       params.PostalCode,
		"proof_of_address_doc_file":         params.ProofOfAddressDocFile,
		"proof_of_address_doc_type":         params.ProofOfAddressDocType,
		"purpose_of_transactions":           params.PurposeOfTransactions,
		"source_of_funds_doc_file":          params.SourceOfFundsDocFile,
		"source_of_funds_doc_type":          params.SourceOfFundsDocType,
		"state_province_region":             params.StateProvinceRegion,
		"tax_id":                            params.TaxID,
		"tos_id":                            params.TosID,
	}

	if params.AddressLine2 != nil {
		body["address_line_2"] = params.AddressLine2
	}
	if params.PhoneNumber != nil {
		body["phone_number"] = params.PhoneNumber
	}
	if params.IDDocBackFile != nil {
		body["id_doc_back_file"] = params.IDDocBackFile
	}
	if params.PurposeOfTransactionsExplanation != nil {
		body["purpose_of_transactions_explanation"] = params.PurposeOfTransactionsExplanation
	}
	if params.ExternalID != nil {
		body["external_id"] = params.ExternalID
	}

	return request.Do[*CreateResponse](c.cfg, ctx, "POST", path, body)
}

// CreateBusinessWithStandardKYB creates a business receiver with standard KYB.
func (c *Client) CreateBusinessWithStandardKYB(ctx context.Context, params *CreateBusinessStandardParams) (*CreateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := fmt.Sprintf("/instances/%s/receivers", c.instanceID)

	body := map[string]any{
		"kyc_type":                    "standard",
		"type":                        "business",
		"address_line_1":              params.AddressLine1,
		"alternate_name":              params.AlternateName,
		"city":                        params.City,
		"country":                     params.Country,
		"email":                       params.Email,
		"formation_date":              params.FormationDate,
		"incorporation_doc_file":      params.IncorporationDocFile,
		"legal_name":                  params.LegalName,
		"owners":                      params.Owners,
		"postal_code":                 params.PostalCode,
		"proof_of_address_doc_file":   params.ProofOfAddressDocFile,
		"proof_of_address_doc_type":   params.ProofOfAddressDocType,
		"proof_of_ownership_doc_file": params.ProofOfOwnershipDocFile,
		"state_province_region":       params.StateProvinceRegion,
		"tax_id":                      params.TaxID,
		"tos_id":                      params.TosID,
	}

	if params.AddressLine2 != nil {
		body["address_line_2"] = params.AddressLine2
	}
	if params.Website != nil {
		body["website"] = params.Website
	}
	if params.ExternalID != nil {
		body["external_id"] = params.ExternalID
	}

	return request.Do[*CreateResponse](c.cfg, ctx, "POST", path, body)
}

// Get retrieves a specific receiver by ID.
func (c *Client) Get(ctx context.Context, receiverID string) (*Receiver, error) {
	if receiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s", c.instanceID, receiverID)
	return request.Do[*Receiver](c.cfg, ctx, "GET", path, nil)
}

// Update updates a receiver.
func (c *Client) Update(ctx context.Context, params *UpdateParams) error {
	if params == nil {
		return fmt.Errorf("params cannot be nil")
	}
	if params.ReceiverID == "" {
		return fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s", c.instanceID, params.ReceiverID)

	// Build the request body with only non-nil fields
	body := make(map[string]any)

	if params.Email != nil {
		body["email"] = params.Email
	}
	if params.TaxID != nil {
		body["tax_id"] = params.TaxID
	}
	if params.AddressLine1 != nil {
		body["address_line_1"] = params.AddressLine1
	}
	if params.AddressLine2 != nil {
		body["address_line_2"] = params.AddressLine2
	}
	if params.City != nil {
		body["city"] = params.City
	}
	if params.StateProvinceRegion != nil {
		body["state_province_region"] = params.StateProvinceRegion
	}
	if params.Country != nil {
		body["country"] = params.Country
	}
	if params.PostalCode != nil {
		body["postal_code"] = params.PostalCode
	}
	if params.IPAddress != nil {
		body["ip_address"] = params.IPAddress
	}
	if params.ImageURL != nil {
		body["image_url"] = params.ImageURL
	}
	if params.PhoneNumber != nil {
		body["phone_number"] = params.PhoneNumber
	}
	if params.ProofOfAddressDocType != nil {
		body["proof_of_address_doc_type"] = params.ProofOfAddressDocType
	}
	if params.ProofOfAddressDocFile != nil {
		body["proof_of_address_doc_file"] = params.ProofOfAddressDocFile
	}
	if params.FirstName != nil {
		body["first_name"] = params.FirstName
	}
	if params.LastName != nil {
		body["last_name"] = params.LastName
	}
	if params.DateOfBirth != nil {
		body["date_of_birth"] = params.DateOfBirth
	}
	if params.IDDocCountry != nil {
		body["id_doc_country"] = params.IDDocCountry
	}
	if params.IDDocType != nil {
		body["id_doc_type"] = params.IDDocType
	}
	if params.IDDocFrontFile != nil {
		body["id_doc_front_file"] = params.IDDocFrontFile
	}
	if params.IDDocBackFile != nil {
		body["id_doc_back_file"] = params.IDDocBackFile
	}
	if params.LegalName != nil {
		body["legal_name"] = params.LegalName
	}
	if params.AlternateName != nil {
		body["alternate_name"] = params.AlternateName
	}
	if params.FormationDate != nil {
		body["formation_date"] = params.FormationDate
	}
	if params.Website != nil {
		body["website"] = params.Website
	}
	if params.Owners != nil {
		body["owners"] = params.Owners
	}
	if params.IncorporationDocFile != nil {
		body["incorporation_doc_file"] = params.IncorporationDocFile
	}
	if params.ProofOfOwnershipDocFile != nil {
		body["proof_of_ownership_doc_file"] = params.ProofOfOwnershipDocFile
	}
	if params.SourceOfFundsDocType != nil {
		body["source_of_funds_doc_type"] = params.SourceOfFundsDocType
	}
	if params.SourceOfFundsDocFile != nil {
		body["source_of_funds_doc_file"] = params.SourceOfFundsDocFile
	}
	if params.IndividualHoldingDocFrontFile != nil {
		body["individual_holding_doc_front_file"] = params.IndividualHoldingDocFrontFile
	}
	if params.PurposeOfTransactions != nil {
		body["purpose_of_transactions"] = params.PurposeOfTransactions
	}
	if params.PurposeOfTransactionsExplanation != nil {
		body["purpose_of_transactions_explanation"] = params.PurposeOfTransactionsExplanation
	}
	if params.ExternalID != nil {
		body["external_id"] = params.ExternalID
	}
	if params.TosID != nil {
		body["tos_id"] = params.TosID
	}

	_, err := request.Do[struct{}](c.cfg, ctx, "PATCH", path, body)
	return err
}

// Delete deletes a receiver.
func (c *Client) Delete(ctx context.Context, receiverID string) error {
	if receiverID == "" {
		return fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s", c.instanceID, receiverID)
	_, err := request.Do[struct{}](c.cfg, ctx, "DELETE", path, nil)
	return err
}

// GetLimits retrieves transaction limits for a receiver.
func (c *Client) GetLimits(ctx context.Context, receiverID string) (*LimitsResponse, error) {
	if receiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/limits/receivers/%s", c.instanceID, receiverID)
	return request.Do[*LimitsResponse](c.cfg, ctx, "GET", path, nil)
}

// GetLimitIncreaseRequests retrieves all limit increase requests for a receiver.
func (c *Client) GetLimitIncreaseRequests(ctx context.Context, receiverID string) ([]LimitIncreaseRequest, error) {
	if receiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/limit-increase", c.instanceID, receiverID)
	return request.Do[[]LimitIncreaseRequest](c.cfg, ctx, "GET", path, nil)
}

// RequestLimitIncrease creates a new limit increase request for a receiver.
func (c *Client) RequestLimitIncrease(ctx context.Context, params *RequestLimitIncreaseParams) (*RequestLimitIncreaseResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.ReceiverID == "" {
		return nil, fmt.Errorf("receiver ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/receivers/%s/limit-increase", c.instanceID, params.ReceiverID)

	body := map[string]any{
		"daily":                    params.Daily,
		"monthly":                  params.Monthly,
		"per_transaction":          params.PerTransaction,
		"supporting_document_file": params.SupportingDocumentFile,
		"supporting_document_type": params.SupportingDocumentType,
	}

	return request.Do[*RequestLimitIncreaseResponse](c.cfg, ctx, "POST", path, body)
}
