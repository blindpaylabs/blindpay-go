package payables

import (
	"context"
	"fmt"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/blindpaylabs/blindpay-go/option"
	"github.com/blindpaylabs/blindpay-go/pagination"
)

// Type is the kind of bill a payable represents (list filter).
type Type string

const (
	TypeBoleto  Type = "boleto"
	TypePixCode Type = "pix_code"
	TypeInvoice Type = "invoice"
)

// CreateParams creates a payable. Amounts (Discount, Taxes, line item
// prices) are integer cents. Set the field that matches the bill:
// BoletoBarcode, PixQrcode, or DocumentFile plus BankAccount for an invoice.
type CreateParams struct {
	CustomerID    string         `json:"customer_id"`
	Currency      types.Currency `json:"currency"`
	From          *Party         `json:"from,omitempty"`
	To            *Party         `json:"to,omitempty"`
	LineItems     []LineItem     `json:"line_items,omitempty"`
	Note          *string        `json:"note,omitempty"`
	Discount      *int64         `json:"discount,omitempty"`
	Taxes         *int64         `json:"taxes,omitempty"`
	BoletoBarcode *string        `json:"boleto_barcode,omitempty"`
	PixQrcode     *string        `json:"pix_qrcode,omitempty"`
	DocumentFile  *string        `json:"document_file,omitempty"`
	// DueDate and ScheduledAt are dates as the API accepts them
	// (e.g. "2026-12-31").
	DueDate     *string `json:"due_date,omitempty"`
	ScheduledAt *string `json:"scheduled_at,omitempty"`
	// BankAccount is the beneficiary account, in the same shape as the body
	// of a bank account create (e.g. {"type": "ach", "name": ...,
	// "routing_number": ..., "account_number": ...}).
	BankAccount map[string]any `json:"bank_account,omitempty"`
}

// ListParams filters and paginates payables.
type ListParams struct {
	// Limit is the page size (10, 50, 100, 200, 500 or 1000). Setting Limit,
	// StartingAfter or EndingBefore makes the API return one page; with none
	// of them it returns every payable.
	Limit int
	// StartingAfter returns the page after this payable ID (the previous
	// page's Pagination.NextPage).
	StartingAfter string
	// EndingBefore returns the page before this payable ID.
	EndingBefore    string
	Status          Status
	CustomerID      string
	CustomerName    string
	Type            Type
	BeneficiaryName string
}

// Client handles payable operations.
type Client struct {
	cfg        *request.Config
	instanceID string
}

// NewClient creates a new payables client.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:        cfg.ToRequestConfig(),
		instanceID: cfg.InstanceID,
	}
}

// Create creates a payable.
func (c *Client) Create(ctx context.Context, params *CreateParams, opts ...option.RequestOption) (*Payable, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := request.Pathf("/instances/%s/payables", c.instanceID)
	return request.Do[*Payable](c.cfg, ctx, "POST", path, params, opts...)
}

// List retrieves payables matching params. Without Limit or a cursor the
// API returns every payable; use ListAutoPaging for large lists.
func (c *Client) List(ctx context.Context, params *ListParams, opts ...option.RequestOption) ([]Payable, error) {
	return request.DoList[Payable](c.cfg, ctx, c.listPath(params), opts...)
}

// ListPage retrieves one page of payables with its pagination metadata.
// Set Limit (or a cursor) to get a page; without them the API returns every
// payable and Pagination is the zero value.
func (c *Client) ListPage(ctx context.Context, params *ListParams, opts ...option.RequestOption) (*pagination.Page[Payable], error) {
	list, err := request.Do[request.List[Payable]](c.cfg, ctx, "GET", c.listPath(params), nil, opts...)
	if err != nil {
		return nil, err
	}
	page := &pagination.Page[Payable]{Data: list.Data}
	if list.Pagination != nil {
		page.Pagination = *list.Pagination
	}
	return page, nil
}

// ListAutoPaging iterates over every payable matching params, one page at a
// time (Limit defaults to 100), starting after params.StartingAfter. It only
// pages forward: EndingBefore makes the pager fail with
// pagination.ErrEndingBeforeUnsupported.
func (c *Client) ListAutoPaging(ctx context.Context, params *ListParams, opts ...option.RequestOption) *pagination.Pager[Payable] {
	p := ListParams{}
	if params != nil {
		p = *params
	}
	if p.EndingBefore != "" {
		return pagination.NewErrPager[Payable](pagination.ErrEndingBeforeUnsupported)
	}
	if p.Limit == 0 {
		p.Limit = pagination.DefaultLimit
	}
	return pagination.NewPager(ctx, p.StartingAfter, func(ctx context.Context, cursor string) (*pagination.Page[Payable], error) {
		page := p
		page.StartingAfter = cursor
		return c.ListPage(ctx, &page, opts...)
	})
}

func (c *Client) listPath(params *ListParams) string {
	path := request.Pathf("/instances/%s/payables", c.instanceID)
	if params == nil {
		return path
	}
	var q request.Query
	q.Int("limit", params.Limit)
	q.String("starting_after", params.StartingAfter)
	q.String("ending_before", params.EndingBefore)
	q.String("status", string(params.Status))
	q.String("customer_id", params.CustomerID)
	q.String("customer_name", params.CustomerName)
	q.String("type", string(params.Type))
	q.String("beneficiary_name", params.BeneficiaryName)
	return q.Apply(path)
}

// Get retrieves a payable by ID.
func (c *Client) Get(ctx context.Context, id string, opts ...option.RequestOption) (*Payable, error) {
	if id == "" {
		return nil, fmt.Errorf("payable ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/payables/%s", c.instanceID, id)
	return request.Do[*Payable](c.cfg, ctx, "GET", path, nil, opts...)
}

// Delete deletes a payable and returns it as it was before deletion.
func (c *Client) Delete(ctx context.Context, id string, opts ...option.RequestOption) (*Payable, error) {
	if id == "" {
		return nil, fmt.Errorf("payable ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/payables/%s", c.instanceID, id)
	return request.Do[*Payable](c.cfg, ctx, "DELETE", path, nil, opts...)
}
