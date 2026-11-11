package bankaccounts

import (
	"context"
	"fmt"

	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/option"
)

// CreatePlaidLinkParams starts a Plaid Link session for a customer.
type CreatePlaidLinkParams struct {
	CustomerID string `json:"-"`
	// RedirectURL is where Plaid sends the customer back after linking.
	RedirectURL *string `json:"redirect_url,omitempty"`
}

// PlaidLink is a hosted Plaid Link session.
type PlaidLink struct {
	// URL is the hosted page where the customer connects a US bank account
	// through Plaid.
	URL string `json:"url"`
}

// CreatePlaidLink creates a hosted Plaid Link URL the customer opens to
// connect a US bank account (used for ACH pull payins).
func (c *Client) CreatePlaidLink(ctx context.Context, params *CreatePlaidLinkParams, opts ...option.RequestOption) (*PlaidLink, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.CustomerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/bank-accounts/plaid", c.instanceID, params.CustomerID)
	return request.Do[*PlaidLink](c.cfg, ctx, "POST", path, params, opts...)
}
