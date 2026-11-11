package instances

import (
	"context"

	"github.com/blindpaylabs/blindpay-go/customers"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/blindpaylabs/blindpay-go/option"
)

// Type is the kind of instance. Development instances settle with the test
// token USDB and move no real money.
type Type string

const (
	TypeProduction  Type = "production"
	TypeDevelopment Type = "development"
)

// OnboardingStep is how far the instance's onboarding has progressed.
type OnboardingStep string

const (
	OnboardingStepInitial         OnboardingStep = "initial"
	OnboardingStepBusinessDetails OnboardingStep = "business_details"
	OnboardingStepBusinessProfile OnboardingStep = "business_profile"
	OnboardingStepDocuments       OnboardingStep = "documents"
	OnboardingStepApplicant       OnboardingStep = "applicant"
	OnboardingStepCompleted       OnboardingStep = "completed"
	OnboardingStepBlocked         OnboardingStep = "blocked"
)

// Instance is the instance the client is bound to.
type Instance struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Type           Type           `json:"type"`
	OnboardingStep OnboardingStep `json:"onboarding_step"`
}

// BusinessDetails is the instance's onboarding business details.
type BusinessDetails struct {
	LegalName           *string        `json:"legal_name,omitempty"`
	TaxID               *string        `json:"tax_id,omitempty"`
	FormationDate       *string        `json:"formation_date,omitempty"`
	TradingNames        []string       `json:"trading_names,omitempty"`
	AddressLine1        *string        `json:"address_line_1,omitempty"`
	AddressLine2        *string        `json:"address_line_2,omitempty"`
	City                *string        `json:"city,omitempty"`
	StateProvinceRegion *string        `json:"state_province_region,omitempty"`
	PostalCode          *string        `json:"postal_code,omitempty"`
	Country             *types.Country `json:"country,omitempty"`
	// EntityType is e.g. "company", "corporation", "sole_proprietor".
	EntityType           *string `json:"entity_type,omitempty"`
	StateOfIncorporation *string `json:"state_of_incorporation,omitempty"`
}

// BusinessProfile is the instance's onboarding business profile.
type BusinessProfile struct {
	// Scenario is "custodial", "custodial_mtl", "non_custodial" or "other".
	Scenario            *string `json:"scenario,omitempty"`
	ScenarioDescription *string `json:"scenario_description,omitempty"`
	HasKycScreening     *bool   `json:"has_kyc_screening,omitempty"`
	IsRegulated         *bool   `json:"is_regulated,omitempty"`
	// BusinessIndustry is the instance's segment ("wallet", "exchange", ...),
	// not a NAICS code.
	BusinessIndustry   *string         `json:"business_industry,omitempty"`
	Website            *string         `json:"website,omitempty"`
	TopVolumeByCountry []types.Country `json:"top_volume_by_country,omitempty"`
	// EstimatedVolumePerMonth is a bucket such as "up_to_1000".
	EstimatedVolumePerMonth *string `json:"estimated_volume_per_month,omitempty"`
}

// OwnershipDocuments is the instance's onboarding documents and owners.
type OwnershipDocuments struct {
	IncorporationDocFile    *string                          `json:"incorporation_doc_file,omitempty"`
	ProofOfOwnershipDocFile *string                          `json:"proof_of_ownership_doc_file,omitempty"`
	ProofOfAddressDocType   *customers.ProofOfAddressDocType `json:"proof_of_address_doc_type,omitempty"`
	ProofOfAddressDocFile   *string                          `json:"proof_of_address_doc_file,omitempty"`
	Owners                  []OnboardingOwner                `json:"owners,omitempty"`
}

// OnboardingOwner is a beneficial owner or controlling person of the
// instance's business. OwnershipPercentage is 0 to 100.
type OnboardingOwner struct {
	ID                    *string                           `json:"id,omitempty"`
	InstanceID            *string                           `json:"instance_id,omitempty"`
	Role                  *customers.OwnerRole              `json:"role,omitempty"`
	FirstName             *string                           `json:"first_name,omitempty"`
	LastName              *string                           `json:"last_name,omitempty"`
	OwnershipPercentage   *float64                          `json:"ownership_percentage,omitempty"`
	DateOfBirth           *string                           `json:"date_of_birth,omitempty"`
	TaxID                 *string                           `json:"tax_id,omitempty"`
	AddressLine1          *string                           `json:"address_line_1,omitempty"`
	AddressLine2          *string                           `json:"address_line_2,omitempty"`
	City                  *string                           `json:"city,omitempty"`
	StateProvinceRegion   *string                           `json:"state_province_region,omitempty"`
	Country               *types.Country                    `json:"country,omitempty"`
	PostalCode            *string                           `json:"postal_code,omitempty"`
	IDDocCountry          *types.Country                    `json:"id_doc_country,omitempty"`
	IDDocType             *customers.IdentificationDocument `json:"id_doc_type,omitempty"`
	IDDocFrontFile        *string                           `json:"id_doc_front_file,omitempty"`
	IDDocBackFile         *string                           `json:"id_doc_back_file,omitempty"`
	ProofOfAddressDocType *customers.ProofOfAddressDocType  `json:"proof_of_address_doc_type,omitempty"`
	ProofOfAddressDocFile *string                           `json:"proof_of_address_doc_file,omitempty"`
}

// Get retrieves the instance the client is bound to.
func (c *Client) Get(ctx context.Context, opts ...option.RequestOption) (*Instance, error) {
	path := request.Pathf("/instances/%s", c.instanceID)
	return request.Do[*Instance](c.cfg, ctx, "GET", path, nil, opts...)
}

// GetBusinessDetails retrieves the instance's onboarding business details.
func (c *Client) GetBusinessDetails(ctx context.Context, opts ...option.RequestOption) (*BusinessDetails, error) {
	path := request.Pathf("/instances/%s/onboarding/business_details", c.instanceID)
	return request.Do[*BusinessDetails](c.cfg, ctx, "GET", path, nil, opts...)
}

// GetBusinessProfile retrieves the instance's onboarding business profile.
func (c *Client) GetBusinessProfile(ctx context.Context, opts ...option.RequestOption) (*BusinessProfile, error) {
	path := request.Pathf("/instances/%s/onboarding/business_profile", c.instanceID)
	return request.Do[*BusinessProfile](c.cfg, ctx, "GET", path, nil, opts...)
}

// GetOwnershipDocuments retrieves the instance's onboarding documents and
// owners.
func (c *Client) GetOwnershipDocuments(ctx context.Context, opts ...option.RequestOption) (*OwnershipDocuments, error) {
	path := request.Pathf("/instances/%s/onboarding/ownership_documents", c.instanceID)
	return request.Do[*OwnershipDocuments](c.cfg, ctx, "GET", path, nil, opts...)
}
