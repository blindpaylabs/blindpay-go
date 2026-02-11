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

// AccountPurpose represents the account purpose.
type AccountPurpose string

const (
	AccountPurposeCharitableDonations              AccountPurpose = "charitable_donations"
	AccountPurposeEcommerceRetailPayments           AccountPurpose = "ecommerce_retail_payments"
	AccountPurposeInvestmentPurposes                AccountPurpose = "investment_purposes"
	AccountPurposeBusinessExpenses                  AccountPurpose = "business_expenses"
	AccountPurposePaymentsToFriendsOrFamilyAbroad   AccountPurpose = "payments_to_friends_or_family_abroad"
	AccountPurposePersonalOrLivingExpenses           AccountPurpose = "personal_or_living_expenses"
	AccountPurposeProtectWealth                     AccountPurpose = "protect_wealth"
	AccountPurposePurchaseGoodsAndServices           AccountPurpose = "purchase_goods_and_services"
	AccountPurposeReceivePaymentsForGoodsAndServices AccountPurpose = "receive_payments_for_goods_and_services"
	AccountPurposeTaxOptimization                   AccountPurpose = "tax_optimization"
	AccountPurposeThirdPartyMoneyTransmission        AccountPurpose = "third_party_money_transmission"
	AccountPurposeOther                             AccountPurpose = "other"
	AccountPurposePayroll                           AccountPurpose = "payroll"
	AccountPurposeTreasuryManagement                AccountPurpose = "treasury_management"
)

// BusinessType represents the legal structure of a business.
type BusinessType string

const (
	BusinessTypeCorporation        BusinessType = "corporation"
	BusinessTypeLLC                BusinessType = "llc"
	BusinessTypePartnership        BusinessType = "partnership"
	BusinessTypeSoleProprietorship BusinessType = "sole_proprietorship"
	BusinessTypeTrust              BusinessType = "trust"
	BusinessTypeNonProfit          BusinessType = "non_profit"
)

// BusinessIndustry represents the NAICS industry code.
type BusinessIndustry string

const (
	BusinessIndustry541511 BusinessIndustry = "541511"
	BusinessIndustry541512 BusinessIndustry = "541512"
	BusinessIndustry541519 BusinessIndustry = "541519"
	BusinessIndustry518210 BusinessIndustry = "518210"
	BusinessIndustry511210 BusinessIndustry = "511210"
	BusinessIndustry541611 BusinessIndustry = "541611"
	BusinessIndustry541618 BusinessIndustry = "541618"
	BusinessIndustry541330 BusinessIndustry = "541330"
	BusinessIndustry541990 BusinessIndustry = "541990"
	BusinessIndustry522110 BusinessIndustry = "522110"
	BusinessIndustry523110 BusinessIndustry = "523110"
	BusinessIndustry523920 BusinessIndustry = "523920"
	BusinessIndustry423430 BusinessIndustry = "423430"
	BusinessIndustry423690 BusinessIndustry = "423690"
	BusinessIndustry423110 BusinessIndustry = "423110"
	BusinessIndustry423830 BusinessIndustry = "423830"
	BusinessIndustry423840 BusinessIndustry = "423840"
	BusinessIndustry423510 BusinessIndustry = "423510"
	BusinessIndustry424210 BusinessIndustry = "424210"
	BusinessIndustry424690 BusinessIndustry = "424690"
	BusinessIndustry424990 BusinessIndustry = "424990"
	BusinessIndustry454110 BusinessIndustry = "454110"
	BusinessIndustry334111 BusinessIndustry = "334111"
	BusinessIndustry334118 BusinessIndustry = "334118"
	BusinessIndustry325412 BusinessIndustry = "325412"
	BusinessIndustry339112 BusinessIndustry = "339112"
	BusinessIndustry336111 BusinessIndustry = "336111"
	BusinessIndustry336390 BusinessIndustry = "336390"
	BusinessIndustry551112 BusinessIndustry = "551112"
	BusinessIndustry561499 BusinessIndustry = "561499"
	BusinessIndustry488510 BusinessIndustry = "488510"
	BusinessIndustry484121 BusinessIndustry = "484121"
	BusinessIndustry493110 BusinessIndustry = "493110"
	BusinessIndustry424410 BusinessIndustry = "424410"
	BusinessIndustry424480 BusinessIndustry = "424480"
	BusinessIndustry315990 BusinessIndustry = "315990"
	BusinessIndustry313110 BusinessIndustry = "313110"
	BusinessIndustry213112 BusinessIndustry = "213112"
	BusinessIndustry517110 BusinessIndustry = "517110"
	BusinessIndustry541214 BusinessIndustry = "541214"
)

// EstimatedAnnualRevenue represents the estimated annual revenue range.
type EstimatedAnnualRevenue string

const (
	EstimatedAnnualRevenue0to99999            EstimatedAnnualRevenue = "0_99999"
	EstimatedAnnualRevenue100000to999999      EstimatedAnnualRevenue = "100000_999999"
	EstimatedAnnualRevenue1000000to9999999    EstimatedAnnualRevenue = "1000000_9999999"
	EstimatedAnnualRevenue10000000to49999999  EstimatedAnnualRevenue = "10000000_49999999"
	EstimatedAnnualRevenue50000000to249999999 EstimatedAnnualRevenue = "50000000_249999999"
	EstimatedAnnualRevenue2500000000Plus      EstimatedAnnualRevenue = "2500000000_plus"
)

// SourceOfWealth represents the source of wealth.
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
	ProofOfAddressDocFile *string                `json:"proof_of_address_doc_file"`
	OwnershipPercentage   *float64               `json:"ownership_percentage"`
	Title                 *string                `json:"title"`
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
	Country               types.Country           `json:"country"`
	Email                 string                  `json:"email"`
	AccountPurpose        *AccountPurpose         `json:"account_purpose,omitempty"`
	AddressLine1          *string                 `json:"address_line_1,omitempty"`
	AddressLine2          *string                 `json:"address_line_2,omitempty"`
	City                  *string                 `json:"city,omitempty"`
	DateOfBirth           *string                 `json:"date_of_birth,omitempty"`
	ExternalID            *string                 `json:"external_id,omitempty"`
	FirstName             *string                 `json:"first_name,omitempty"`
	IDDocBackFile         *string                 `json:"id_doc_back_file,omitempty"`
	IDDocCountry          *types.Country          `json:"id_doc_country,omitempty"`
	IDDocFrontFile        *string                 `json:"id_doc_front_file,omitempty"`
	IDDocType             *IdentificationDocument `json:"id_doc_type,omitempty"`
	ImageURL              *string                 `json:"image_url,omitempty"`
	IPAddress             *string                 `json:"ip_address,omitempty"`
	LastName              *string                 `json:"last_name,omitempty"`
	PhoneNumber           *string                 `json:"phone_number,omitempty"`
	PostalCode            *string                 `json:"postal_code,omitempty"`
	ProofOfAddressDocFile *string                 `json:"proof_of_address_doc_file,omitempty"`
	ProofOfAddressDocType *ProofOfAddressDocType  `json:"proof_of_address_doc_type,omitempty"`
	SelfieFile            *string                 `json:"selfie_file,omitempty"`
	SourceOfFundsDocFile  *string                 `json:"source_of_funds_doc_file,omitempty"`
	SourceOfFundsDocType  *SourceOfFundsDocType   `json:"source_of_funds_doc_type,omitempty"`
	SourceOfWealth        *SourceOfWealth         `json:"source_of_wealth,omitempty"`
	StateProvinceRegion   *string                 `json:"state_province_region,omitempty"`
	TaxID                 *string                 `json:"tax_id,omitempty"`
	TosID                 *string                 `json:"tos_id,omitempty"`
}

// CreateIndividualEnhancedParams represents parameters for creating an individual with enhanced KYC.
type CreateIndividualEnhancedParams struct {
	Country                          types.Country           `json:"country"`
	Email                            string                  `json:"email"`
	AccountPurpose                   *AccountPurpose         `json:"account_purpose,omitempty"`
	AddressLine1                     *string                 `json:"address_line_1,omitempty"`
	AddressLine2                     *string                 `json:"address_line_2,omitempty"`
	City                             *string                 `json:"city,omitempty"`
	DateOfBirth                      *string                 `json:"date_of_birth,omitempty"`
	ExternalID                       *string                 `json:"external_id,omitempty"`
	FirstName                        *string                 `json:"first_name,omitempty"`
	IDDocBackFile                    *string                 `json:"id_doc_back_file,omitempty"`
	IDDocCountry                     *types.Country          `json:"id_doc_country,omitempty"`
	IDDocFrontFile                   *string                 `json:"id_doc_front_file,omitempty"`
	IDDocType                        *IdentificationDocument `json:"id_doc_type,omitempty"`
	ImageURL                         *string                 `json:"image_url,omitempty"`
	IPAddress                        *string                 `json:"ip_address,omitempty"`
	LastName                         *string                 `json:"last_name,omitempty"`
	PhoneNumber                      *string                 `json:"phone_number,omitempty"`
	PostalCode                       *string                 `json:"postal_code,omitempty"`
	ProofOfAddressDocFile            *string                 `json:"proof_of_address_doc_file,omitempty"`
	ProofOfAddressDocType            *ProofOfAddressDocType  `json:"proof_of_address_doc_type,omitempty"`
	PurposeOfTransactions            *PurposeOfTransactions  `json:"purpose_of_transactions,omitempty"`
	PurposeOfTransactionsExplanation *string                 `json:"purpose_of_transactions_explanation,omitempty"`
	SelfieFile                       *string                 `json:"selfie_file,omitempty"`
	SourceOfFundsDocFile             *string                 `json:"source_of_funds_doc_file,omitempty"`
	SourceOfFundsDocType             *SourceOfFundsDocType   `json:"source_of_funds_doc_type,omitempty"`
	SourceOfWealth                   *SourceOfWealth         `json:"source_of_wealth,omitempty"`
	StateProvinceRegion              *string                 `json:"state_province_region,omitempty"`
	TaxID                            *string                 `json:"tax_id,omitempty"`
	TosID                            *string                 `json:"tos_id,omitempty"`
}

// CreateBusinessStandardParams represents parameters for creating a business with standard KYB.
type CreateBusinessStandardParams struct {
	Country                 types.Country           `json:"country"`
	Email                   string                  `json:"email"`
	AccountPurpose          *AccountPurpose         `json:"account_purpose,omitempty"`
	AddressLine1            *string                 `json:"address_line_1,omitempty"`
	AddressLine2            *string                 `json:"address_line_2,omitempty"`
	AlternateName           *string                 `json:"alternate_name,omitempty"`
	BusinessDescription     *string                 `json:"business_description,omitempty"`
	BusinessIndustry        *BusinessIndustry       `json:"business_industry,omitempty"`
	BusinessType            *BusinessType           `json:"business_type,omitempty"`
	City                    *string                 `json:"city,omitempty"`
	EstimatedAnnualRevenue  *EstimatedAnnualRevenue `json:"estimated_annual_revenue,omitempty"`
	ExternalID              *string                 `json:"external_id,omitempty"`
	FormationDate           *string                 `json:"formation_date,omitempty"`
	ImageURL                *string                 `json:"image_url,omitempty"`
	IncorporationDocFile    *string                 `json:"incorporation_doc_file,omitempty"`
	IPAddress               *string                 `json:"ip_address,omitempty"`
	LegalName               *string                 `json:"legal_name,omitempty"`
	Owners                  []Owner                 `json:"owners,omitempty"`
	PhoneNumber             *string                 `json:"phone_number,omitempty"`
	PostalCode              *string                 `json:"postal_code,omitempty"`
	ProofOfAddressDocFile   *string                 `json:"proof_of_address_doc_file,omitempty"`
	ProofOfAddressDocType   *ProofOfAddressDocType  `json:"proof_of_address_doc_type,omitempty"`
	ProofOfOwnershipDocFile *string                 `json:"proof_of_ownership_doc_file,omitempty"`
	PubliclyTraded          *bool                   `json:"publicly_traded,omitempty"`
	SourceOfFundsDocFile    *string                 `json:"source_of_funds_doc_file,omitempty"`
	SourceOfFundsDocType    *SourceOfFundsDocType   `json:"source_of_funds_doc_type,omitempty"`
	SourceOfWealth          *SourceOfWealth         `json:"source_of_wealth,omitempty"`
	StateProvinceRegion     *string                 `json:"state_province_region,omitempty"`
	TaxID                   *string                 `json:"tax_id,omitempty"`
	TosID                   *string                 `json:"tos_id,omitempty"`
	Website                 *string                 `json:"website,omitempty"`
}

// CreateResponse represents the response when creating a receiver.
type CreateResponse struct {
	ID string `json:"id"`
}

// UpdateParams represents parameters for updating a receiver.
type UpdateParams struct {
	ReceiverID                       string                  `json:"-"`
	AccountPurpose                   *AccountPurpose         `json:"account_purpose,omitempty"`
	AddressLine1                     *string                 `json:"address_line_1,omitempty"`
	AddressLine2                     *string                 `json:"address_line_2,omitempty"`
	AlternateName                    *string                 `json:"alternate_name,omitempty"`
	BusinessDescription              *string                 `json:"business_description,omitempty"`
	BusinessIndustry                 *BusinessIndustry       `json:"business_industry,omitempty"`
	BusinessType                     *BusinessType           `json:"business_type,omitempty"`
	City                             *string                 `json:"city,omitempty"`
	Country                          *types.Country          `json:"country,omitempty"`
	DateOfBirth                      *string                 `json:"date_of_birth,omitempty"`
	Email                            *string                 `json:"email,omitempty"`
	EstimatedAnnualRevenue           *EstimatedAnnualRevenue `json:"estimated_annual_revenue,omitempty"`
	ExternalID                       *string                 `json:"external_id,omitempty"`
	FirstName                        *string                 `json:"first_name,omitempty"`
	FormationDate                    *string                 `json:"formation_date,omitempty"`
	IDDocBackFile                    *string                 `json:"id_doc_back_file,omitempty"`
	IDDocCountry                     *types.Country          `json:"id_doc_country,omitempty"`
	IDDocFrontFile                   *string                 `json:"id_doc_front_file,omitempty"`
	IDDocType                        *IdentificationDocument `json:"id_doc_type,omitempty"`
	ImageURL                         *string                 `json:"image_url,omitempty"`
	IncorporationDocFile             *string                 `json:"incorporation_doc_file,omitempty"`
	IPAddress                        *string                 `json:"ip_address,omitempty"`
	LastName                         *string                 `json:"last_name,omitempty"`
	LegalName                        *string                 `json:"legal_name,omitempty"`
	Owners                           []Owner                 `json:"owners,omitempty"`
	PhoneNumber                      *string                 `json:"phone_number,omitempty"`
	PostalCode                       *string                 `json:"postal_code,omitempty"`
	ProofOfAddressDocFile            *string                 `json:"proof_of_address_doc_file,omitempty"`
	ProofOfAddressDocType            *ProofOfAddressDocType  `json:"proof_of_address_doc_type,omitempty"`
	ProofOfOwnershipDocFile          *string                 `json:"proof_of_ownership_doc_file,omitempty"`
	PubliclyTraded                   *bool                   `json:"publicly_traded,omitempty"`
	PurposeOfTransactions            *PurposeOfTransactions  `json:"purpose_of_transactions,omitempty"`
	PurposeOfTransactionsExplanation *string                 `json:"purpose_of_transactions_explanation,omitempty"`
	SelfieFile                       *string                 `json:"selfie_file,omitempty"`
	SourceOfFundsDocFile             *string                 `json:"source_of_funds_doc_file,omitempty"`
	SourceOfFundsDocType             *SourceOfFundsDocType   `json:"source_of_funds_doc_type,omitempty"`
	SourceOfWealth                   *SourceOfWealth         `json:"source_of_wealth,omitempty"`
	StateProvinceRegion              *string                 `json:"state_province_region,omitempty"`
	TaxID                            *string                 `json:"tax_id,omitempty"`
	TosID                            *string                 `json:"tos_id,omitempty"`
	Website                          *string                 `json:"website,omitempty"`
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
		cfg: &request.Config{
			BaseURL:    cfg.BaseURL,
			APIKey:     cfg.APIKey,
			HTTPClient: cfg.HTTPClient,
			UserAgent:  cfg.UserAgent,
		},
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
		"kyc_type": "standard",
		"type":     "individual",
		"country":  params.Country,
		"email":    params.Email,
	}

	if params.AccountPurpose != nil {
		body["account_purpose"] = *params.AccountPurpose
	}
	if params.AddressLine1 != nil {
		body["address_line_1"] = *params.AddressLine1
	}
	if params.AddressLine2 != nil {
		body["address_line_2"] = *params.AddressLine2
	}
	if params.City != nil {
		body["city"] = *params.City
	}
	if params.DateOfBirth != nil {
		body["date_of_birth"] = *params.DateOfBirth
	}
	if params.ExternalID != nil {
		body["external_id"] = *params.ExternalID
	}
	if params.FirstName != nil {
		body["first_name"] = *params.FirstName
	}
	if params.IDDocBackFile != nil {
		body["id_doc_back_file"] = *params.IDDocBackFile
	}
	if params.IDDocCountry != nil {
		body["id_doc_country"] = *params.IDDocCountry
	}
	if params.IDDocFrontFile != nil {
		body["id_doc_front_file"] = *params.IDDocFrontFile
	}
	if params.IDDocType != nil {
		body["id_doc_type"] = *params.IDDocType
	}
	if params.ImageURL != nil {
		body["image_url"] = *params.ImageURL
	}
	if params.IPAddress != nil {
		body["ip_address"] = *params.IPAddress
	}
	if params.LastName != nil {
		body["last_name"] = *params.LastName
	}
	if params.PhoneNumber != nil {
		body["phone_number"] = *params.PhoneNumber
	}
	if params.PostalCode != nil {
		body["postal_code"] = *params.PostalCode
	}
	if params.ProofOfAddressDocFile != nil {
		body["proof_of_address_doc_file"] = *params.ProofOfAddressDocFile
	}
	if params.ProofOfAddressDocType != nil {
		body["proof_of_address_doc_type"] = *params.ProofOfAddressDocType
	}
	if params.SelfieFile != nil {
		body["selfie_file"] = *params.SelfieFile
	}
	if params.SourceOfFundsDocFile != nil {
		body["source_of_funds_doc_file"] = *params.SourceOfFundsDocFile
	}
	if params.SourceOfFundsDocType != nil {
		body["source_of_funds_doc_type"] = *params.SourceOfFundsDocType
	}
	if params.SourceOfWealth != nil {
		body["source_of_wealth"] = *params.SourceOfWealth
	}
	if params.StateProvinceRegion != nil {
		body["state_province_region"] = *params.StateProvinceRegion
	}
	if params.TaxID != nil {
		body["tax_id"] = *params.TaxID
	}
	if params.TosID != nil {
		body["tos_id"] = *params.TosID
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
		"kyc_type": "enhanced",
		"type":     "individual",
		"country":  params.Country,
		"email":    params.Email,
	}

	if params.AccountPurpose != nil {
		body["account_purpose"] = *params.AccountPurpose
	}
	if params.AddressLine1 != nil {
		body["address_line_1"] = *params.AddressLine1
	}
	if params.AddressLine2 != nil {
		body["address_line_2"] = *params.AddressLine2
	}
	if params.City != nil {
		body["city"] = *params.City
	}
	if params.DateOfBirth != nil {
		body["date_of_birth"] = *params.DateOfBirth
	}
	if params.ExternalID != nil {
		body["external_id"] = *params.ExternalID
	}
	if params.FirstName != nil {
		body["first_name"] = *params.FirstName
	}
	if params.IDDocBackFile != nil {
		body["id_doc_back_file"] = *params.IDDocBackFile
	}
	if params.IDDocCountry != nil {
		body["id_doc_country"] = *params.IDDocCountry
	}
	if params.IDDocFrontFile != nil {
		body["id_doc_front_file"] = *params.IDDocFrontFile
	}
	if params.IDDocType != nil {
		body["id_doc_type"] = *params.IDDocType
	}
	if params.ImageURL != nil {
		body["image_url"] = *params.ImageURL
	}
	if params.IPAddress != nil {
		body["ip_address"] = *params.IPAddress
	}
	if params.LastName != nil {
		body["last_name"] = *params.LastName
	}
	if params.PhoneNumber != nil {
		body["phone_number"] = *params.PhoneNumber
	}
	if params.PostalCode != nil {
		body["postal_code"] = *params.PostalCode
	}
	if params.ProofOfAddressDocFile != nil {
		body["proof_of_address_doc_file"] = *params.ProofOfAddressDocFile
	}
	if params.ProofOfAddressDocType != nil {
		body["proof_of_address_doc_type"] = *params.ProofOfAddressDocType
	}
	if params.PurposeOfTransactions != nil {
		body["purpose_of_transactions"] = *params.PurposeOfTransactions
	}
	if params.PurposeOfTransactionsExplanation != nil {
		body["purpose_of_transactions_explanation"] = *params.PurposeOfTransactionsExplanation
	}
	if params.SelfieFile != nil {
		body["selfie_file"] = *params.SelfieFile
	}
	if params.SourceOfFundsDocFile != nil {
		body["source_of_funds_doc_file"] = *params.SourceOfFundsDocFile
	}
	if params.SourceOfFundsDocType != nil {
		body["source_of_funds_doc_type"] = *params.SourceOfFundsDocType
	}
	if params.SourceOfWealth != nil {
		body["source_of_wealth"] = *params.SourceOfWealth
	}
	if params.StateProvinceRegion != nil {
		body["state_province_region"] = *params.StateProvinceRegion
	}
	if params.TaxID != nil {
		body["tax_id"] = *params.TaxID
	}
	if params.TosID != nil {
		body["tos_id"] = *params.TosID
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
		"kyc_type": "standard",
		"type":     "business",
		"country":  params.Country,
		"email":    params.Email,
	}

	if params.AccountPurpose != nil {
		body["account_purpose"] = *params.AccountPurpose
	}
	if params.AddressLine1 != nil {
		body["address_line_1"] = *params.AddressLine1
	}
	if params.AddressLine2 != nil {
		body["address_line_2"] = *params.AddressLine2
	}
	if params.AlternateName != nil {
		body["alternate_name"] = *params.AlternateName
	}
	if params.BusinessDescription != nil {
		body["business_description"] = *params.BusinessDescription
	}
	if params.BusinessIndustry != nil {
		body["business_industry"] = *params.BusinessIndustry
	}
	if params.BusinessType != nil {
		body["business_type"] = *params.BusinessType
	}
	if params.City != nil {
		body["city"] = *params.City
	}
	if params.EstimatedAnnualRevenue != nil {
		body["estimated_annual_revenue"] = *params.EstimatedAnnualRevenue
	}
	if params.ExternalID != nil {
		body["external_id"] = *params.ExternalID
	}
	if params.FormationDate != nil {
		body["formation_date"] = *params.FormationDate
	}
	if params.ImageURL != nil {
		body["image_url"] = *params.ImageURL
	}
	if params.IncorporationDocFile != nil {
		body["incorporation_doc_file"] = *params.IncorporationDocFile
	}
	if params.IPAddress != nil {
		body["ip_address"] = *params.IPAddress
	}
	if params.LegalName != nil {
		body["legal_name"] = *params.LegalName
	}
	if params.Owners != nil {
		body["owners"] = params.Owners
	}
	if params.PhoneNumber != nil {
		body["phone_number"] = *params.PhoneNumber
	}
	if params.PostalCode != nil {
		body["postal_code"] = *params.PostalCode
	}
	if params.ProofOfAddressDocFile != nil {
		body["proof_of_address_doc_file"] = *params.ProofOfAddressDocFile
	}
	if params.ProofOfAddressDocType != nil {
		body["proof_of_address_doc_type"] = *params.ProofOfAddressDocType
	}
	if params.ProofOfOwnershipDocFile != nil {
		body["proof_of_ownership_doc_file"] = *params.ProofOfOwnershipDocFile
	}
	if params.PubliclyTraded != nil {
		body["publicly_traded"] = *params.PubliclyTraded
	}
	if params.SourceOfFundsDocFile != nil {
		body["source_of_funds_doc_file"] = *params.SourceOfFundsDocFile
	}
	if params.SourceOfFundsDocType != nil {
		body["source_of_funds_doc_type"] = *params.SourceOfFundsDocType
	}
	if params.SourceOfWealth != nil {
		body["source_of_wealth"] = *params.SourceOfWealth
	}
	if params.StateProvinceRegion != nil {
		body["state_province_region"] = *params.StateProvinceRegion
	}
	if params.TaxID != nil {
		body["tax_id"] = *params.TaxID
	}
	if params.TosID != nil {
		body["tos_id"] = *params.TosID
	}
	if params.Website != nil {
		body["website"] = *params.Website
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

	if params.AccountPurpose != nil {
		body["account_purpose"] = *params.AccountPurpose
	}
	if params.AddressLine1 != nil {
		body["address_line_1"] = *params.AddressLine1
	}
	if params.AddressLine2 != nil {
		body["address_line_2"] = *params.AddressLine2
	}
	if params.AlternateName != nil {
		body["alternate_name"] = *params.AlternateName
	}
	if params.BusinessDescription != nil {
		body["business_description"] = *params.BusinessDescription
	}
	if params.BusinessIndustry != nil {
		body["business_industry"] = *params.BusinessIndustry
	}
	if params.BusinessType != nil {
		body["business_type"] = *params.BusinessType
	}
	if params.City != nil {
		body["city"] = *params.City
	}
	if params.Country != nil {
		body["country"] = *params.Country
	}
	if params.DateOfBirth != nil {
		body["date_of_birth"] = *params.DateOfBirth
	}
	if params.Email != nil {
		body["email"] = *params.Email
	}
	if params.EstimatedAnnualRevenue != nil {
		body["estimated_annual_revenue"] = *params.EstimatedAnnualRevenue
	}
	if params.ExternalID != nil {
		body["external_id"] = *params.ExternalID
	}
	if params.FirstName != nil {
		body["first_name"] = *params.FirstName
	}
	if params.FormationDate != nil {
		body["formation_date"] = *params.FormationDate
	}
	if params.IDDocBackFile != nil {
		body["id_doc_back_file"] = *params.IDDocBackFile
	}
	if params.IDDocCountry != nil {
		body["id_doc_country"] = *params.IDDocCountry
	}
	if params.IDDocFrontFile != nil {
		body["id_doc_front_file"] = *params.IDDocFrontFile
	}
	if params.IDDocType != nil {
		body["id_doc_type"] = *params.IDDocType
	}
	if params.ImageURL != nil {
		body["image_url"] = *params.ImageURL
	}
	if params.IncorporationDocFile != nil {
		body["incorporation_doc_file"] = *params.IncorporationDocFile
	}
	if params.IPAddress != nil {
		body["ip_address"] = *params.IPAddress
	}
	if params.LastName != nil {
		body["last_name"] = *params.LastName
	}
	if params.LegalName != nil {
		body["legal_name"] = *params.LegalName
	}
	if params.Owners != nil {
		body["owners"] = params.Owners
	}
	if params.PhoneNumber != nil {
		body["phone_number"] = *params.PhoneNumber
	}
	if params.PostalCode != nil {
		body["postal_code"] = *params.PostalCode
	}
	if params.ProofOfAddressDocFile != nil {
		body["proof_of_address_doc_file"] = *params.ProofOfAddressDocFile
	}
	if params.ProofOfAddressDocType != nil {
		body["proof_of_address_doc_type"] = *params.ProofOfAddressDocType
	}
	if params.ProofOfOwnershipDocFile != nil {
		body["proof_of_ownership_doc_file"] = *params.ProofOfOwnershipDocFile
	}
	if params.PubliclyTraded != nil {
		body["publicly_traded"] = *params.PubliclyTraded
	}
	if params.PurposeOfTransactions != nil {
		body["purpose_of_transactions"] = *params.PurposeOfTransactions
	}
	if params.PurposeOfTransactionsExplanation != nil {
		body["purpose_of_transactions_explanation"] = *params.PurposeOfTransactionsExplanation
	}
	if params.SelfieFile != nil {
		body["selfie_file"] = *params.SelfieFile
	}
	if params.SourceOfFundsDocFile != nil {
		body["source_of_funds_doc_file"] = *params.SourceOfFundsDocFile
	}
	if params.SourceOfFundsDocType != nil {
		body["source_of_funds_doc_type"] = *params.SourceOfFundsDocType
	}
	if params.SourceOfWealth != nil {
		body["source_of_wealth"] = *params.SourceOfWealth
	}
	if params.StateProvinceRegion != nil {
		body["state_province_region"] = *params.StateProvinceRegion
	}
	if params.TaxID != nil {
		body["tax_id"] = *params.TaxID
	}
	if params.TosID != nil {
		body["tos_id"] = *params.TosID
	}
	if params.Website != nil {
		body["website"] = *params.Website
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
