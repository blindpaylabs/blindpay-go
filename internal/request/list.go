package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/blindpaylabs/blindpay-go/pagination"
)

// List is the decoded body of a list endpoint whose response shape depends
// on the query: payouts, payins, payables and customers return a bare JSON
// array unless limit, starting_after or ending_before is sent, and
// {"data": [...], "pagination": {...}} when any of them is (blindpay-v2
// modules/{payouts,payins,customers}/endpoints list handlers).
//
// The shape is detected from the body (first non-space byte), never inferred
// from the request params. Pagination is nil for a bare array.
type List[T any] struct {
	Data       []T
	Pagination *types.PaginationMetadata
}

// UnmarshalJSON implements json.Unmarshaler for both list shapes.
func (l *List[T]) UnmarshalJSON(data []byte) error {
	*l = List[T]{}

	trimmed := bytes.TrimLeft(data, " \t\r\n")
	if len(trimmed) == 0 {
		return fmt.Errorf("blindpay: empty list response")
	}

	switch trimmed[0] {
	case '[':
		return json.Unmarshal(trimmed, &l.Data)
	case '{':
		var envelope struct {
			Data       []T                       `json:"data"`
			Pagination *types.PaginationMetadata `json:"pagination"`
		}
		if err := json.Unmarshal(trimmed, &envelope); err != nil {
			return err
		}
		l.Data = envelope.Data
		l.Pagination = envelope.Pagination
		return nil
	case 'n':
		// JSON null: no rows.
		return json.Unmarshal(trimmed, &l.Data)
	default:
		return fmt.Errorf("blindpay: list response is neither an array nor an object")
	}
}

// DoList GETs a list endpoint that may answer with either shape (see List)
// and returns the rows.
func DoList[T any](cfg *Config, ctx context.Context, path string, opts ...RequestOption) ([]T, error) {
	list, err := Do[List[T]](cfg, ctx, "GET", path, nil, opts...)
	if err != nil {
		return nil, err
	}
	return list.Data, nil
}

// DoPage GETs one page of a list endpoint (either shape, see List). A bare
// array (every row) yields a Page with zero-value Pagination.
func DoPage[T any](cfg *Config, ctx context.Context, path string, opts ...RequestOption) (*pagination.Page[T], error) {
	list, err := Do[List[T]](cfg, ctx, "GET", path, nil, opts...)
	if err != nil {
		return nil, err
	}
	page := &pagination.Page[T]{Data: list.Data}
	if list.Pagination != nil {
		page.Pagination = *list.Pagination
	}
	return page, nil
}
