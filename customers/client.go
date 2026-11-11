package customers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/blindpaylabs/blindpay-go/option"
	"github.com/blindpaylabs/blindpay-go/pagination"
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

// KycStatus represents KYC verification status.
type KycStatus string

const (
	KycStatusVerifying         KycStatus = "verifying"
	KycStatusApproved          KycStatus = "approved"
	KycStatusRejected          KycStatus = "rejected"
	KycStatusDeprecated        KycStatus = "deprecated"
	KycStatusPendingReview     KycStatus = "pending_review"
	KycStatusAwaitingContract  KycStatus = "awaiting_contract"
	KycStatusComplianceRequest KycStatus = "compliance_request"
)

// AccountPurpose represents the purpose of the account.
type AccountPurpose string

const (
	AccountPurposeCharitableDonations                AccountPurpose = "charitable_donations"
	AccountPurposeEcommerceRetailPayments            AccountPurpose = "ecommerce_retail_payments"
	AccountPurposeInvestmentPurposes                 AccountPurpose = "investment_purposes"
	AccountPurposeBusinessExpenses                   AccountPurpose = "business_expenses"
	AccountPurposePaymentsToFriendsOrFamilyAbroad    AccountPurpose = "payments_to_friends_or_family_abroad"
	AccountPurposePersonalOrLivingExpenses           AccountPurpose = "personal_or_living_expenses"
	AccountPurposeProtectWealth                      AccountPurpose = "protect_wealth"
	AccountPurposePurchaseGoodsAndServices           AccountPurpose = "purchase_goods_and_services"
	AccountPurposeReceivePaymentsForGoodsAndServices AccountPurpose = "receive_payments_for_goods_and_services"
	AccountPurposeTaxOptimization                    AccountPurpose = "tax_optimization"
	AccountPurposeThirdPartyMoneyTransmission        AccountPurpose = "third_party_money_transmission"
	AccountPurposePayroll                            AccountPurpose = "payroll"
	AccountPurposeTreasuryManagement                 AccountPurpose = "treasury_management"
	AccountPurposeOther                              AccountPurpose = "other"
)

// BusinessType represents the type of business entity.
type BusinessType string

const (
	BusinessTypeCorporation        BusinessType = "corporation"
	BusinessTypeLLC                BusinessType = "llc"
	BusinessTypePartnership        BusinessType = "partnership"
	BusinessTypeSoleProprietorship BusinessType = "sole_proprietorship"
	BusinessTypeTrust              BusinessType = "trust"
	BusinessTypeNonProfit          BusinessType = "non_profit"
)

// EstimatedAnnualRevenue represents estimated annual revenue ranges.
type EstimatedAnnualRevenue string

const (
	EstimatedAnnualRevenue0To99999            EstimatedAnnualRevenue = "0_99999"
	EstimatedAnnualRevenue100000To999999      EstimatedAnnualRevenue = "100000_999999"
	EstimatedAnnualRevenue1000000To9999999    EstimatedAnnualRevenue = "1000000_9999999"
	EstimatedAnnualRevenue10000000To49999999  EstimatedAnnualRevenue = "10000000_49999999"
	EstimatedAnnualRevenue50000000To249999999 EstimatedAnnualRevenue = "50000000_249999999"
	EstimatedAnnualRevenue250000000Plus       EstimatedAnnualRevenue = "250000000_plus"
	// Deprecated: the name carried an extra digit and the value the API rejects.
	// Use EstimatedAnnualRevenue250000000Plus.
	EstimatedAnnualRevenue2500000000Plus EstimatedAnnualRevenue = "250000000_plus"
)

// SourceOfWealth represents the source of wealth for a customer.
type SourceOfWealth string

const (
	SourceOfWealthBusinessDividendsOrProfits  SourceOfWealth = "business_dividends_or_profits"
	SourceOfWealthInvestments                 SourceOfWealth = "investments"
	SourceOfWealthAssetSales                  SourceOfWealth = "asset_sales"
	SourceOfWealthClientInvestorContributions SourceOfWealth = "client_investor_contributions"
	SourceOfWealthGambling                    SourceOfWealth = "gambling"
	SourceOfWealthCharitableContributions     SourceOfWealth = "charitable_contributions"
	SourceOfWealthInheritance                 SourceOfWealth = "inheritance"
	SourceOfWealthAffiliateOrRoyaltyIncome    SourceOfWealth = "affiliate_or_royalty_income"
)

// AmlStatus represents AML screening status.
type AmlStatus string

const (
	AmlStatusClear AmlStatus = "clear"
	AmlStatusHit   AmlStatus = "hit"
	AmlStatusError AmlStatus = "error"
)

// TaxType represents tax identifier types.
type TaxType string

const (
	TaxTypeSSN  TaxType = "SSN"
	TaxTypeITIN TaxType = "ITIN"
)

// FraudWarning represents a fraud warning.
type FraudWarning struct {
	ID        *string  `json:"id"`
	Name      *string  `json:"name"`
	Operation *string  `json:"operation"`
	Score     *float64 `json:"score"`
}

// AmlHits represents AML screening match details.
type AmlHits struct {
	HasSanctionMatch     bool `json:"has_sanction_match"`
	HasPepMatch          bool `json:"has_pep_match"`
	HasWatchlistMatch    bool `json:"has_watchlist_match"`
	HasCrimelistMatch    bool `json:"has_crimelist_match"`
	HasAdversemediaMatch bool `json:"has_adversemedia_match"`
}

// Owner represents a business owner.
type Owner struct {
	ID                    string                 `json:"id,omitempty"`
	InstanceID            string                 `json:"instance_id,omitempty"`
	CustomerID            string                 `json:"customer_id,omitempty"`
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
	OwnershipPercentage   *int                   `json:"ownership_percentage,omitempty"`
	Title                 *string                `json:"title,omitempty"`
	TaxType               TaxType                `json:"tax_type,omitempty"`
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

// Customer represents a customer with KYC information.
type Customer struct {
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
	SelfieFile                       string                `json:"selfie_file,omitempty"`
	PurposeOfTransactions            PurposeOfTransactions `json:"purpose_of_transactions,omitempty"`
	PurposeOfTransactionsExplanation *string               `json:"purpose_of_transactions_explanation,omitempty"`
	// Business fields
	LegalName               string                       `json:"legal_name,omitempty"`
	AlternateName           *string                      `json:"alternate_name,omitempty"`
	FormationDate           string                       `json:"formation_date,omitempty"`
	Website                 *string                      `json:"website,omitempty"`
	Owners                  []Owner                      `json:"owners,omitempty"`
	IncorporationDocFile    string                       `json:"incorporation_doc_file,omitempty"`
	ProofOfOwnershipDocFile string                       `json:"proof_of_ownership_doc_file,omitempty"`
	ExternalID              *string                      `json:"external_id,omitempty"`
	FraudWarnings           []FraudWarning               `json:"fraud_warnings,omitempty"`
	IsFbo                   *bool                        `json:"is_fbo,omitempty"`
	AccountPurpose          AccountPurpose               `json:"account_purpose,omitempty"`
	AccountPurposeOther     *string                      `json:"account_purpose_other,omitempty"`
	BusinessType            BusinessType                 `json:"business_type,omitempty"`
	BusinessDescription     *string                      `json:"business_description,omitempty"`
	BusinessIndustry        types.BusinessIndustry       `json:"business_industry,omitempty"`
	EstimatedAnnualRevenue  EstimatedAnnualRevenue       `json:"estimated_annual_revenue,omitempty"`
	SourceOfWealth          SourceOfWealth               `json:"source_of_wealth,omitempty"`
	PubliclyTraded          *bool                        `json:"publicly_traded,omitempty"`
	Occupation              *string                      `json:"occupation,omitempty"`
	RecipientRelationship   *types.RecipientRelationship `json:"recipient_relationship,omitempty"`
	SoleProprietorDocType   *types.SoleProprietorDocType `json:"sole_proprietor_doc_type,omitempty"`
	AmlStatus               AmlStatus                    `json:"aml_status,omitempty"`
	AmlHitsData             *AmlHits                     `json:"aml_hits,omitempty"`
	IsTosAccepted           *bool                        `json:"is_tos_accepted,omitempty"`
}

// CreateIndividualStandardParams represents parameters for creating an individual with standard KYC.
type CreateIndividualStandardParams struct {
	ExternalID            *string                      `json:"external_id,omitempty"`
	AddressLine1          string                       `json:"address_line_1"`
	AddressLine2          *string                      `json:"address_line_2,omitempty"`
	City                  string                       `json:"city"`
	Country               types.Country                `json:"country"`
	DateOfBirth           string                       `json:"date_of_birth"`
	Email                 string                       `json:"email"`
	FirstName             string                       `json:"first_name"`
	PhoneNumber           *string                      `json:"phone_number,omitempty"`
	IDDocCountry          types.Country                `json:"id_doc_country"`
	IDDocFrontFile        string                       `json:"id_doc_front_file"`
	IDDocType             IdentificationDocument       `json:"id_doc_type"`
	IDDocBackFile         *string                      `json:"id_doc_back_file,omitempty"`
	LastName              string                       `json:"last_name"`
	PostalCode            string                       `json:"postal_code"`
	ProofOfAddressDocFile string                       `json:"proof_of_address_doc_file"`
	ProofOfAddressDocType ProofOfAddressDocType        `json:"proof_of_address_doc_type"`
	StateProvinceRegion   string                       `json:"state_province_region"`
	TaxID                 string                       `json:"tax_id"`
	TosID                 string                       `json:"tos_id"`
	IPAddress             *string                      `json:"ip_address,omitempty"`
	ImageURL              *string                      `json:"image_url,omitempty"`
	Occupation            *string                      `json:"occupation,omitempty"`
	RecipientRelationship *types.RecipientRelationship `json:"recipient_relationship,omitempty"`
	SoleProprietorDocType *types.SoleProprietorDocType `json:"sole_proprietor_doc_type,omitempty"`
	// Latitude and Longitude are the customer's location, if known.
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	// AdditionalInfo attaches extra documents (label + file URL).
	AdditionalInfo                   []AdditionalInfoItem   `json:"additional_info,omitempty"`
	SelfieFile                       *string                `json:"selfie_file,omitempty"`
	SourceOfFundsDocType             *SourceOfFundsDocType  `json:"source_of_funds_doc_type,omitempty"`
	SourceOfFundsDocFile             *string                `json:"source_of_funds_doc_file,omitempty"`
	PurposeOfTransactions            *PurposeOfTransactions `json:"purpose_of_transactions,omitempty"`
	PurposeOfTransactionsExplanation *string                `json:"purpose_of_transactions_explanation,omitempty"`
}

// CreateIndividualEnhancedParams represents parameters for creating an individual with enhanced KYC.
type CreateIndividualEnhancedParams struct {
	ExternalID                       *string                      `json:"external_id,omitempty"`
	AddressLine1                     string                       `json:"address_line_1"`
	AddressLine2                     *string                      `json:"address_line_2,omitempty"`
	City                             string                       `json:"city"`
	Country                          types.Country                `json:"country"`
	DateOfBirth                      string                       `json:"date_of_birth"`
	Email                            string                       `json:"email"`
	FirstName                        string                       `json:"first_name"`
	IDDocCountry                     types.Country                `json:"id_doc_country"`
	IDDocFrontFile                   string                       `json:"id_doc_front_file"`
	IDDocType                        IdentificationDocument       `json:"id_doc_type"`
	IDDocBackFile                    *string                      `json:"id_doc_back_file,omitempty"`
	SelfieFile                       string                       `json:"selfie_file"`
	LastName                         string                       `json:"last_name"`
	PostalCode                       string                       `json:"postal_code"`
	PhoneNumber                      *string                      `json:"phone_number,omitempty"`
	ProofOfAddressDocFile            string                       `json:"proof_of_address_doc_file"`
	ProofOfAddressDocType            ProofOfAddressDocType        `json:"proof_of_address_doc_type"`
	PurposeOfTransactions            PurposeOfTransactions        `json:"purpose_of_transactions"`
	SourceOfFundsDocFile             string                       `json:"source_of_funds_doc_file"`
	SourceOfFundsDocType             SourceOfFundsDocType         `json:"source_of_funds_doc_type"`
	PurposeOfTransactionsExplanation *string                      `json:"purpose_of_transactions_explanation,omitempty"`
	StateProvinceRegion              string                       `json:"state_province_region"`
	TaxID                            string                       `json:"tax_id"`
	TosID                            string                       `json:"tos_id"`
	IPAddress                        *string                      `json:"ip_address,omitempty"`
	ImageURL                         *string                      `json:"image_url,omitempty"`
	Occupation                       *string                      `json:"occupation,omitempty"`
	RecipientRelationship            *types.RecipientRelationship `json:"recipient_relationship,omitempty"`
	SoleProprietorDocType            *types.SoleProprietorDocType `json:"sole_proprietor_doc_type,omitempty"`
	// Latitude and Longitude are the customer's location, if known.
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	// AdditionalInfo attaches extra documents (label + file URL).
	AdditionalInfo []AdditionalInfoItem `json:"additional_info,omitempty"`
}

// CreateBusinessStandardParams represents parameters for creating a business with standard KYB.
type CreateBusinessStandardParams struct {
	ExternalID              *string                      `json:"external_id,omitempty"`
	AddressLine1            string                       `json:"address_line_1"`
	AddressLine2            *string                      `json:"address_line_2,omitempty"`
	AlternateName           string                       `json:"alternate_name"`
	City                    string                       `json:"city"`
	Country                 types.Country                `json:"country"`
	Email                   string                       `json:"email"`
	FormationDate           string                       `json:"formation_date"`
	IncorporationDocFile    string                       `json:"incorporation_doc_file"`
	LegalName               string                       `json:"legal_name"`
	Owners                  []Owner                      `json:"owners"`
	PostalCode              string                       `json:"postal_code"`
	ProofOfAddressDocFile   string                       `json:"proof_of_address_doc_file"`
	ProofOfAddressDocType   ProofOfAddressDocType        `json:"proof_of_address_doc_type"`
	ProofOfOwnershipDocFile string                       `json:"proof_of_ownership_doc_file"`
	StateProvinceRegion     string                       `json:"state_province_region"`
	TaxID                   string                       `json:"tax_id"`
	TosID                   string                       `json:"tos_id"`
	Website                 *string                      `json:"website,omitempty"`
	PhoneNumber             *string                      `json:"phone_number,omitempty"`
	IPAddress               *string                      `json:"ip_address,omitempty"`
	ImageURL                *string                      `json:"image_url,omitempty"`
	AccountPurpose          *AccountPurpose              `json:"account_purpose,omitempty"`
	AccountPurposeOther     *string                      `json:"account_purpose_other,omitempty"`
	BusinessTypeField       *BusinessType                `json:"business_type,omitempty"`
	BusinessDescription     *string                      `json:"business_description,omitempty"`
	BusinessIndustry        *types.BusinessIndustry      `json:"business_industry,omitempty"`
	EstimatedAnnualRevenue  *EstimatedAnnualRevenue      `json:"estimated_annual_revenue,omitempty"`
	SourceOfWealth          *SourceOfWealth              `json:"source_of_wealth,omitempty"`
	PubliclyTraded          *bool                        `json:"publicly_traded,omitempty"`
	RecipientRelationship   *types.RecipientRelationship `json:"recipient_relationship,omitempty"`
	SoleProprietorDocType   *types.SoleProprietorDocType `json:"sole_proprietor_doc_type,omitempty"`
	// Latitude and Longitude are the customer's location, if known.
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	// AdditionalInfo attaches extra documents (label + file URL).
	AdditionalInfo                   []AdditionalInfoItem   `json:"additional_info,omitempty"`
	SelfieFile                       *string                `json:"selfie_file,omitempty"`
	SourceOfFundsDocType             *SourceOfFundsDocType  `json:"source_of_funds_doc_type,omitempty"`
	SourceOfFundsDocFile             *string                `json:"source_of_funds_doc_file,omitempty"`
	PurposeOfTransactions            *PurposeOfTransactions `json:"purpose_of_transactions,omitempty"`
	PurposeOfTransactionsExplanation *string                `json:"purpose_of_transactions_explanation,omitempty"`
}

// CreateResponse represents the response when creating a customer.
type CreateResponse struct {
	ID string `json:"id"`
}

// UpdateParams represents parameters for updating a customer.
type UpdateParams struct {
	CustomerID                       string                       `json:"-"`
	Email                            *string                      `json:"email,omitempty"`
	TaxID                            *string                      `json:"tax_id,omitempty"`
	AddressLine1                     *string                      `json:"address_line_1,omitempty"`
	AddressLine2                     *string                      `json:"address_line_2,omitempty"`
	City                             *string                      `json:"city,omitempty"`
	StateProvinceRegion              *string                      `json:"state_province_region,omitempty"`
	Country                          *types.Country               `json:"country,omitempty"`
	PostalCode                       *string                      `json:"postal_code,omitempty"`
	IPAddress                        *string                      `json:"ip_address,omitempty"`
	ImageURL                         *string                      `json:"image_url,omitempty"`
	PhoneNumber                      *string                      `json:"phone_number,omitempty"`
	ProofOfAddressDocType            *ProofOfAddressDocType       `json:"proof_of_address_doc_type,omitempty"`
	ProofOfAddressDocFile            *string                      `json:"proof_of_address_doc_file,omitempty"`
	FirstName                        *string                      `json:"first_name,omitempty"`
	LastName                         *string                      `json:"last_name,omitempty"`
	DateOfBirth                      *string                      `json:"date_of_birth,omitempty"`
	IDDocCountry                     *types.Country               `json:"id_doc_country,omitempty"`
	IDDocType                        *IdentificationDocument      `json:"id_doc_type,omitempty"`
	IDDocFrontFile                   *string                      `json:"id_doc_front_file,omitempty"`
	IDDocBackFile                    *string                      `json:"id_doc_back_file,omitempty"`
	LegalName                        *string                      `json:"legal_name,omitempty"`
	AlternateName                    *string                      `json:"alternate_name,omitempty"`
	FormationDate                    *string                      `json:"formation_date,omitempty"`
	Website                          *string                      `json:"website,omitempty"`
	Owners                           []Owner                      `json:"owners,omitempty"`
	IncorporationDocFile             *string                      `json:"incorporation_doc_file,omitempty"`
	ProofOfOwnershipDocFile          *string                      `json:"proof_of_ownership_doc_file,omitempty"`
	SourceOfFundsDocType             *SourceOfFundsDocType        `json:"source_of_funds_doc_type,omitempty"`
	SourceOfFundsDocFile             *string                      `json:"source_of_funds_doc_file,omitempty"`
	SelfieFile                       *string                      `json:"selfie_file,omitempty"`
	PurposeOfTransactions            *PurposeOfTransactions       `json:"purpose_of_transactions,omitempty"`
	PurposeOfTransactionsExplanation *string                      `json:"purpose_of_transactions_explanation,omitempty"`
	ExternalID                       *string                      `json:"external_id,omitempty"`
	TosID                            *string                      `json:"tos_id,omitempty"`
	AccountPurpose                   *AccountPurpose              `json:"account_purpose,omitempty"`
	AccountPurposeOther              *string                      `json:"account_purpose_other,omitempty"`
	BusinessTypeField                *BusinessType                `json:"business_type,omitempty"`
	BusinessDescription              *string                      `json:"business_description,omitempty"`
	BusinessIndustry                 *types.BusinessIndustry      `json:"business_industry,omitempty"`
	EstimatedAnnualRevenue           *EstimatedAnnualRevenue      `json:"estimated_annual_revenue,omitempty"`
	SourceOfWealth                   *SourceOfWealth              `json:"source_of_wealth,omitempty"`
	PubliclyTraded                   *bool                        `json:"publicly_traded,omitempty"`
	Occupation                       *string                      `json:"occupation,omitempty"`
	RecipientRelationship            *types.RecipientRelationship `json:"recipient_relationship,omitempty"`
	SoleProprietorDocType            *types.SoleProprietorDocType `json:"sole_proprietor_doc_type,omitempty"`
	// Latitude and Longitude are the customer's location, if known.
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	// AdditionalInfo attaches extra documents (label + file URL).
	AdditionalInfo []AdditionalInfoItem `json:"additional_info,omitempty"`
}

// LimitsResponse represents customer limits.
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

// LimitIncreaseRequest represents a limit increase request for a customer.
type LimitIncreaseRequest struct {
	ID                     string                                     `json:"id"`
	CustomerID             string                                     `json:"customer_id"`
	Status                 LimitIncreaseRequestStatus                 `json:"status"`
	Daily                  float64                                    `json:"daily"`
	Monthly                float64                                    `json:"monthly"`
	PerTransaction         float64                                    `json:"per_transaction"`
	ApprovedPerTransaction *int                                       `json:"approved_per_transaction,omitempty"`
	ApprovedDaily          *int                                       `json:"approved_daily,omitempty"`
	ApprovedMonthly        *int                                       `json:"approved_monthly,omitempty"`
	SupportingDocumentFile string                                     `json:"supporting_document_file"`
	SupportingDocumentType LimitIncreaseRequestSupportingDocumentType `json:"supporting_document_type"`
	CreatedAt              string                                     `json:"created_at"`
	UpdatedAt              string                                     `json:"updated_at"`
}

// AdditionalInfoLabel names the kind of an additional document (e.g.
// "BANK_STATEMENT"). Unknown future values decode as-is.
type AdditionalInfoLabel string

// AdditionalInfoItem is an extra document attached to a customer.
type AdditionalInfoItem struct {
	Label AdditionalInfoLabel `json:"label"`
	// Value is the document's file URL (from Upload).
	Value string `json:"value"`
}

// LimitIncreaseSupportingDocument is one document backing a limit increase
// request.
type LimitIncreaseSupportingDocument struct {
	Type LimitIncreaseRequestSupportingDocumentType `json:"type"`
	// File is the document's URL (from Upload).
	File string `json:"file"`
}

// RequestLimitIncreaseParams represents parameters for requesting a limit increase.
type RequestLimitIncreaseParams struct {
	CustomerID     string  `json:"-"`
	Daily          float64 `json:"daily"`
	Monthly        float64 `json:"monthly"`
	PerTransaction float64 `json:"per_transaction"`
	// SupportingDocumentFile and SupportingDocumentType attach one
	// document; SupportingDocuments attaches several.
	SupportingDocumentFile string                                     `json:"supporting_document_file,omitempty"`
	SupportingDocumentType LimitIncreaseRequestSupportingDocumentType `json:"supporting_document_type,omitempty"`
	SupportingDocuments    []LimitIncreaseSupportingDocument          `json:"supporting_documents,omitempty"`
}

// RequestLimitIncreaseResponse represents the response when requesting a limit increase.
type RequestLimitIncreaseResponse struct {
	ID string `json:"id"`
}

// ListParams represents parameters for listing customers with pagination.
type ListParams struct {
	// Limit is the page size ("10", "50", "100", "200", "500" or "1000").
	Limit string `json:"-"`
	// Deprecated: ignored by the API; use StartingAfter.
	Offset        string    `json:"-"`
	StartingAfter string    `json:"-"`
	EndingBefore  string    `json:"-"`
	FullName      string    `json:"-"`
	CustomerName  string    `json:"-"`
	Status        KycStatus `json:"-"`
	CustomerID    string    `json:"-"`
	BankAccountID string    `json:"-"`
	Country       string    `json:"-"`
}

// ListResponse represents a paginated list of customers.
type ListResponse struct {
	Data       []Customer               `json:"data"`
	Pagination types.PaginationMetadata `json:"pagination"`
}

// Client handles customer-related operations.
type Client struct {
	cfg        *request.Config
	instanceID string
}

// NewClient creates a new customers client.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:        cfg.ToRequestConfig(),
		instanceID: cfg.InstanceID,
	}
}

// List retrieves all customers for the instance.
func (c *Client) List(ctx context.Context, opts ...option.RequestOption) ([]Customer, error) {
	path := request.Pathf("/instances/%s/customers", c.instanceID)
	return request.DoList[Customer](c.cfg, ctx, path, opts...)
}

// ListWithParams retrieves customers with pagination and filtering.
//
// The API wraps the result in {data, pagination} only when Limit,
// StartingAfter or EndingBefore is set; with filters alone (or nil params) it
// returns every matching customer as a bare array. Both shapes are decoded;
// for a bare array Pagination is the zero value (HasMore false).
func (c *Client) ListWithParams(ctx context.Context, params *ListParams, opts ...option.RequestOption) (*ListResponse, error) {
	page, err := c.ListPage(ctx, params, opts...)
	if err != nil {
		return nil, err
	}
	return &ListResponse{Data: page.Data, Pagination: page.Pagination}, nil
}

// ListPage retrieves one page of customers with its pagination metadata
// (same request as ListWithParams).
func (c *Client) ListPage(ctx context.Context, params *ListParams, opts ...option.RequestOption) (*pagination.Page[Customer], error) {
	list, err := request.Do[request.List[Customer]](c.cfg, ctx, "GET", c.listPath(params), nil, opts...)
	if err != nil {
		return nil, err
	}
	page := &pagination.Page[Customer]{Data: list.Data}
	if list.Pagination != nil {
		page.Pagination = *list.Pagination
	}
	return page, nil
}

// ListAutoPaging iterates over every customer matching params, one page at
// a time (Limit defaults to "100"), starting after params.StartingAfter. It
// only pages forward: EndingBefore makes the pager fail with
// pagination.ErrEndingBeforeUnsupported.
func (c *Client) ListAutoPaging(ctx context.Context, params *ListParams, opts ...option.RequestOption) *pagination.Pager[Customer] {
	p := ListParams{}
	if params != nil {
		p = *params
	}
	if p.EndingBefore != "" {
		return pagination.NewErrPager[Customer](pagination.ErrEndingBeforeUnsupported)
	}
	if p.Limit == "" {
		p.Limit = strconv.Itoa(pagination.DefaultLimit)
	}
	return pagination.NewPager(ctx, p.StartingAfter, func(ctx context.Context, cursor string) (*pagination.Page[Customer], error) {
		page := p
		page.StartingAfter = cursor
		return c.ListPage(ctx, &page, opts...)
	})
}

func (c *Client) listPath(params *ListParams) string {
	path := request.Pathf("/instances/%s/customers", c.instanceID)
	if params == nil {
		return path
	}
	var q request.Query
	q.String("limit", params.Limit)
	q.String("offset", params.Offset)
	q.String("starting_after", params.StartingAfter)
	q.String("ending_before", params.EndingBefore)
	q.String("full_name", params.FullName)
	q.String("customer_name", params.CustomerName)
	q.String("status", string(params.Status))
	q.String("customer_id", params.CustomerID)
	q.String("bank_account_id", params.BankAccountID)
	q.String("country", params.Country)
	return q.Apply(path)
}

// CreateIndividualWithStandardKYC creates an individual customer with standard KYC.
func (c *Client) CreateIndividualWithStandardKYC(ctx context.Context, params *CreateIndividualStandardParams, opts ...option.RequestOption) (*CreateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := request.Pathf("/instances/%s/customers", c.instanceID)

	body, err := request.BodyWith(params, map[string]any{"kyc_type": "standard", "type": "individual"})
	if err != nil {
		return nil, err
	}

	return request.Do[*CreateResponse](c.cfg, ctx, "POST", path, body, opts...)
}

// CreateIndividualWithEnhancedKYC creates an individual customer with enhanced KYC.
func (c *Client) CreateIndividualWithEnhancedKYC(ctx context.Context, params *CreateIndividualEnhancedParams, opts ...option.RequestOption) (*CreateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := request.Pathf("/instances/%s/customers", c.instanceID)

	body, err := request.BodyWith(params, map[string]any{"kyc_type": "enhanced", "type": "individual"})
	if err != nil {
		return nil, err
	}

	return request.Do[*CreateResponse](c.cfg, ctx, "POST", path, body, opts...)
}

// CreateBusinessWithStandardKYB creates a business customer with standard KYB.
func (c *Client) CreateBusinessWithStandardKYB(ctx context.Context, params *CreateBusinessStandardParams, opts ...option.RequestOption) (*CreateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := request.Pathf("/instances/%s/customers", c.instanceID)

	body, err := request.BodyWith(params, map[string]any{"kyc_type": "standard", "type": "business"})
	if err != nil {
		return nil, err
	}

	return request.Do[*CreateResponse](c.cfg, ctx, "POST", path, body, opts...)
}

// Get retrieves a specific customer by ID.
func (c *Client) Get(ctx context.Context, customerID string, opts ...option.RequestOption) (*Customer, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s", c.instanceID, customerID)
	return request.Do[*Customer](c.cfg, ctx, "GET", path, nil, opts...)
}

// Update updates a customer.
func (c *Client) Update(ctx context.Context, params *UpdateParams, opts ...option.RequestOption) error {
	if params == nil {
		return fmt.Errorf("params cannot be nil")
	}
	if params.CustomerID == "" {
		return fmt.Errorf("customer ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s", c.instanceID, params.CustomerID)

	_, err := request.Do[struct{}](c.cfg, ctx, "PUT", path, params, opts...)
	return err
}

// Delete deletes a customer.
func (c *Client) Delete(ctx context.Context, customerID string, opts ...option.RequestOption) error {
	if customerID == "" {
		return fmt.Errorf("customer ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s", c.instanceID, customerID)
	_, err := request.Do[struct{}](c.cfg, ctx, "DELETE", path, nil, opts...)
	return err
}

// GetLimits retrieves transaction limits for a customer.
func (c *Client) GetLimits(ctx context.Context, customerID string, opts ...option.RequestOption) (*LimitsResponse, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/limits/customers/%s", c.instanceID, customerID)
	return request.Do[*LimitsResponse](c.cfg, ctx, "GET", path, nil, opts...)
}

// GetLimitIncreaseRequests retrieves all limit increase requests for a customer.
func (c *Client) GetLimitIncreaseRequests(ctx context.Context, customerID string, opts ...option.RequestOption) ([]LimitIncreaseRequest, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/limit-increase", c.instanceID, customerID)
	return request.Do[[]LimitIncreaseRequest](c.cfg, ctx, "GET", path, nil, opts...)
}

// RequestLimitIncrease creates a new limit increase request for a customer.
func (c *Client) RequestLimitIncrease(ctx context.Context, params *RequestLimitIncreaseParams, opts ...option.RequestOption) (*RequestLimitIncreaseResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.CustomerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/limit-increase", c.instanceID, params.CustomerID)

	return request.Do[*RequestLimitIncreaseResponse](c.cfg, ctx, "POST", path, params, opts...)
}
