// Package payables manages payables: bills (boletos, PIX codes, invoices)
// that a customer pays out of BlindPay. The Payable type is also the payload
// of the payable.new, payable.update and payable.complete webhooks.
package payables

import (
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/types"
)

// Status is a payable status. Unknown future values decode as-is.
type Status string

const (
	StatusDraft      Status = "draft"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusCanceled   Status = "canceled"
)

// LastAttemptStatus is the outcome of the latest failed payment attempt.
type LastAttemptStatus string

const (
	LastAttemptStatusFailed   LastAttemptStatus = "failed"
	LastAttemptStatusRefunded LastAttemptStatus = "refunded"
)

// BankAccountType is the rail of a payable's beneficiary bank account.
type BankAccountType string

const (
	BankAccountTypeAch  BankAccountType = "ach"
	BankAccountTypeWire BankAccountType = "wire"
)

// Party is the payer (From) or the beneficiary (To) of a payable.
type Party struct {
	LegalName           *string `json:"legal_name,omitempty"`
	AddressLine1        *string `json:"address_line_1,omitempty"`
	AddressLine2        *string `json:"address_line_2,omitempty"`
	City                *string `json:"city,omitempty"`
	StateProvinceRegion *string `json:"state_province_region,omitempty"`
	PostalCode          *string `json:"postal_code,omitempty"`
	Country             *string `json:"country,omitempty"`
	TaxID               *string `json:"tax_id,omitempty"`
}

// LineItem is one line of a payable. Price is in integer cents.
type LineItem struct {
	Name     string `json:"name"`
	Quantity int64  `json:"quantity"`
	Price    *int64 `json:"price"`
}

// BankAccount is the beneficiary bank account the payable pays into.
type BankAccount struct {
	Type            BankAccountType `json:"type"`
	BeneficiaryName *string         `json:"beneficiary_name,omitempty"`
	RoutingNumber   *string         `json:"routing_number,omitempty"`
	AccountNumber   *string         `json:"account_number,omitempty"`
	AccountType     *string         `json:"account_type,omitempty"`
}

// Payable is a bill to pay. Amount, Discount, Taxes and line item prices are
// integer cents.
type Payable struct {
	ID                string             `json:"id"`
	InstanceID        string             `json:"instance_id"`
	CustomerID        string             `json:"customer_id"`
	From              *Party             `json:"from,omitempty"`
	To                *Party             `json:"to,omitempty"`
	Currency          types.Currency     `json:"currency"`
	LineItems         []LineItem         `json:"line_items,omitempty"`
	Note              *string            `json:"note,omitempty"`
	Discount          *int64             `json:"discount,omitempty"`
	Taxes             *int64             `json:"taxes,omitempty"`
	BoletoBarcode     *string            `json:"boleto_barcode,omitempty"`
	PixQrcode         *string            `json:"pix_qrcode,omitempty"`
	DocumentFile      *string            `json:"document_file,omitempty"`
	DueDate           *string            `json:"due_date,omitempty"`
	ScheduledAt       *string            `json:"scheduled_at,omitempty"`
	Status            Status             `json:"status"`
	CreatedAt         *time.Time         `json:"created_at,omitempty"`
	UpdatedAt         *time.Time         `json:"updated_at,omitempty"`
	Amount            *int64             `json:"amount"`
	PayoutID          *string            `json:"payout_id,omitempty"`
	BankAccount       *BankAccount       `json:"bank_account,omitempty"`
	LastAttemptStatus *LastAttemptStatus `json:"last_attempt_status,omitempty"`
	FirstName         *string            `json:"first_name,omitempty"`
	LastName          *string            `json:"last_name,omitempty"`
	LegalName         *string            `json:"legal_name,omitempty"`
}
