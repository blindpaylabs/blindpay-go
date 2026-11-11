// Package pagination holds the page and auto-pager types returned by the
// cursor-paginated list methods (ListPage and ListAutoPaging on payouts,
// payins, customers and transfers).
//
//	pager := client.Payouts.ListAutoPaging(ctx, &payouts.ListParams{Status: blindpay.TransactionStatusCompleted})
//	for pager.Next() {
//		p := pager.Current()
//		// ...
//	}
//	if err := pager.Err(); err != nil {
//		// ...
//	}
package pagination

import (
	"context"
	"errors"

	"github.com/blindpaylabs/blindpay-go/internal/types"
)

// DefaultLimit is the page size the auto-pagers request when the caller
// sets none (valid for every list endpoint).
const DefaultLimit = 100

// ErrEndingBeforeUnsupported is returned by auto-pagers given EndingBefore:
// they only page forward.
var ErrEndingBeforeUnsupported = errors.New("blindpay: auto-pagination only pages forward; EndingBefore is not supported")

// Page is one page of a list.
type Page[T any] struct {
	Data []T
	// Pagination is the zero value when the API returned a bare array
	// (no limit/starting_after/ending_before sent), i.e. every row.
	Pagination types.PaginationMetadata
}

// HasMore reports whether another page follows.
func (p *Page[T]) HasMore() bool {
	return p.Pagination.HasMore
}

// NextCursor returns the value to pass as StartingAfter for the next page,
// or "" when this is the last page.
func (p *Page[T]) NextCursor() string {
	if !p.Pagination.HasMore || p.Pagination.NextPage == nil {
		return ""
	}
	return *p.Pagination.NextPage
}

// FetchFunc loads the page that starts after the given cursor ("" for the
// first page).
type FetchFunc[T any] func(ctx context.Context, startingAfter string) (*Page[T], error)

// Pager iterates over every item of a list, fetching pages on demand
// (following next_page as starting_after while has_more is true). Each page
// request goes through the client's normal retry policy. A Pager is not
// safe for concurrent use.
type Pager[T any] struct {
	ctx    context.Context
	fetch  FetchFunc[T]
	cursor string
	items  []T
	idx    int
	cur    T
	err    error
	done   bool
	// stopErr is reported once the current page's items are consumed.
	stopErr error
}

// NewPager returns a Pager that starts after startingAfter ("" for the
// beginning).
func NewPager[T any](ctx context.Context, startingAfter string, fetch FetchFunc[T]) *Pager[T] {
	return &Pager[T]{ctx: ctx, cursor: startingAfter, fetch: fetch}
}

// NewErrPager returns a Pager whose Next is false and whose Err is err.
func NewErrPager[T any](err error) *Pager[T] {
	return &Pager[T]{err: err, done: true}
}

// Next advances to the next item, fetching the next page when needed. It
// returns false at the end of the list or on error (see Err).
func (p *Pager[T]) Next() bool {
	for {
		if p.err != nil {
			return false
		}
		if p.idx < len(p.items) {
			p.cur = p.items[p.idx]
			p.idx++
			return true
		}
		if p.done {
			p.err = p.stopErr
			return false
		}

		page, err := p.fetch(p.ctx, p.cursor)
		if err != nil {
			p.err = err
			return false
		}
		if page == nil {
			p.done = true
			return false
		}
		p.items, p.idx = page.Data, 0

		next := page.NextCursor()
		switch {
		case next == "":
			p.done = true
		case next == p.cursor:
			p.done = true
			p.stopErr = errors.New("blindpay: pagination cursor did not advance")
		default:
			p.cursor = next
		}
	}
}

// Current returns the item Next advanced to.
func (p *Pager[T]) Current() T {
	return p.cur
}

// Err returns the error that stopped the iteration, if any.
func (p *Pager[T]) Err() error {
	return p.err
}
