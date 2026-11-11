package custodialwallets

import (
	"context"
	"fmt"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/blindpaylabs/blindpay-go/option"
	"github.com/blindpaylabs/blindpay-go/pagination"
)

// DepositStatus is the state of a deposit into a managed wallet.
type DepositStatus string

const (
	DepositStatusProcessing DepositStatus = "processing"
	DepositStatusCompleted  DepositStatus = "completed"
	DepositStatusOnHold     DepositStatus = "on_hold"
	DepositStatusRefunded   DepositStatus = "refunded"
)

// Deposit is an inbound transfer into a managed wallet.
type Deposit struct {
	ID       string                `json:"id"`
	WalletID string                `json:"wallet_id"`
	Network  types.Network         `json:"network"`
	Token    types.StablecoinToken `json:"token"`
	// Amount is in integer cents of the token (1000 = 10.00 USDC).
	Amount int64 `json:"amount"`
	// TokenAmount is the exact on-chain amount in whole tokens, as a
	// decimal string.
	TokenAmount     *string       `json:"token_amount,omitempty"`
	Status          DepositStatus `json:"status"`
	TransactionHash string        `json:"transaction_hash"`
	SenderAddress   *string       `json:"sender_address,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
}

// ListDepositsParams filters and paginates a wallet's deposits.
type ListDepositsParams struct {
	// Limit is the page size (5, 10, 50, 100 or 200; the API defaults to 10).
	Limit int
	// StartingAfter returns the page after this deposit ID (the previous
	// page's Pagination.NextPage).
	StartingAfter string
	// EndingBefore returns the page before this deposit ID.
	EndingBefore string
	Token        types.StablecoinToken
	Status       DepositStatus
}

// ListDeposits retrieves one page of a managed wallet's deposits (the API
// always paginates this list).
func (c *Client) ListDeposits(ctx context.Context, customerID, walletID string, params *ListDepositsParams, opts ...option.RequestOption) (*pagination.Page[Deposit], error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}
	if walletID == "" {
		return nil, fmt.Errorf("wallet ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/wallets/%s/deposits", c.instanceID, customerID, walletID)
	if params != nil {
		var q request.Query
		q.Int("limit", params.Limit)
		q.String("starting_after", params.StartingAfter)
		q.String("ending_before", params.EndingBefore)
		q.String("token", string(params.Token))
		q.String("status", string(params.Status))
		path = q.Apply(path)
	}

	list, err := request.Do[request.List[Deposit]](c.cfg, ctx, "GET", path, nil, opts...)
	if err != nil {
		return nil, err
	}
	page := &pagination.Page[Deposit]{Data: list.Data}
	if list.Pagination != nil {
		page.Pagination = *list.Pagination
	}
	return page, nil
}

// ListDepositsAutoPaging iterates over every deposit of a managed wallet,
// one page at a time (Limit defaults to 100), starting after
// params.StartingAfter. It only pages forward: EndingBefore makes the pager
// fail with pagination.ErrEndingBeforeUnsupported.
func (c *Client) ListDepositsAutoPaging(ctx context.Context, customerID, walletID string, params *ListDepositsParams, opts ...option.RequestOption) *pagination.Pager[Deposit] {
	p := ListDepositsParams{}
	if params != nil {
		p = *params
	}
	if p.EndingBefore != "" {
		return pagination.NewErrPager[Deposit](pagination.ErrEndingBeforeUnsupported)
	}
	if p.Limit == 0 {
		p.Limit = pagination.DefaultLimit
	}
	return pagination.NewPager(ctx, p.StartingAfter, func(ctx context.Context, cursor string) (*pagination.Page[Deposit], error) {
		page := p
		page.StartingAfter = cursor
		return c.ListDeposits(ctx, customerID, walletID, &page, opts...)
	})
}
