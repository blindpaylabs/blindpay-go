// Package rfi holds the request-for-information (RFI) types shared by
// Customers.GetRfi / SubmitRfi and Instances.GetRfi / SubmitRfi.
//
// When compliance needs more information it opens an RFI: a list of
// sections, each with fields to fill in. Read it with GetRfi (nil when no RFI
// is open) and answer with SubmitRfi, passing a Response keyed by field key.
package rfi

import "time"

// Status is the state of an RFI. Unknown future values decode as-is.
type Status string

const (
	StatusPending   Status = "pending"
	StatusSubmitted Status = "submitted"
	StatusExpired   Status = "expired"
	StatusCancelled Status = "cancelled"
)

// CustomerType is the kind of customer an RFI was opened for.
type CustomerType string

const (
	CustomerTypeIndividual CustomerType = "individual"
	CustomerTypeBusiness   CustomerType = "business"
)

// SubmitReason explains a SubmitResponse with Success false.
type SubmitReason string

const (
	SubmitReasonNoPendingRfi    SubmitReason = "no_pending_rfi"
	SubmitReasonAlreadyResolved SubmitReason = "already_resolved"
)

// FieldItem is one choice of a select-style field.
type FieldItem struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// Field is one answer the RFI asks for. Answer it in Response under Key.
type Field struct {
	Key      string `json:"key"`
	Label    string `json:"label"`
	Required bool   `json:"required"`
	// Regex, when set, is the pattern a text answer must match.
	Regex *string `json:"regex,omitempty"`
	// Items lists the allowed values of a select field; Multiple allows
	// several of them (send a []string).
	Items    []FieldItem `json:"items,omitempty"`
	Multiple *bool       `json:"multiple,omitempty"`
}

// Section groups related fields.
type Section struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	// SupportingDocument is the URL of a document that explains the
	// section, if any.
	SupportingDocument *string `json:"supporting_document,omitempty"`
	Fields             []Field `json:"fields"`
}

// Response is the answer to an RFI: a flat map from Field.Key to its value
// (a string, a []string for multi-select fields, or an uploaded file URL).
// It is sent as the whole JSON request body.
type Response map[string]any

// CustomerRfi is an RFI opened for a customer.
type CustomerRfi struct {
	ID                string         `json:"id"`
	CustomerID        string         `json:"customer_id"`
	InstanceID        string         `json:"instance_id"`
	Status            Status         `json:"status"`
	Request           []Section      `json:"request"`
	Response          map[string]any `json:"response"`
	ExpiresAt         time.Time      `json:"expires_at"`
	SubmittedAt       *time.Time     `json:"submitted_at,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	CustomerType      *CustomerType  `json:"customer_type,omitempty"`
	CustomerKycStatus *string        `json:"customer_kyc_status,omitempty"`
}

// InstanceRfi is an RFI opened for the instance itself.
type InstanceRfi struct {
	ID          string         `json:"id"`
	InstanceID  string         `json:"instance_id"`
	Status      Status         `json:"status"`
	Request     []Section      `json:"request"`
	Response    map[string]any `json:"response"`
	ExpiresAt   time.Time      `json:"expires_at"`
	SubmittedAt *time.Time     `json:"submitted_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
}

// SubmitResponse is the result of SubmitRfi. Success is false, with Reason
// "no_pending_rfi" or "already_resolved", when there was nothing to answer.
type SubmitResponse struct {
	Success bool          `json:"success"`
	Reason  *SubmitReason `json:"reason,omitempty"`
}
