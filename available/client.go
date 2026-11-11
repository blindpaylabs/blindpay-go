package available

import (
	"context"
	"fmt"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/blindpaylabs/blindpay-go/option"
)

// Client handles available-related operations.
type Client struct {
	cfg *request.Config
}

// GetSwiftCodeBankDetailsResponse represents the response when retrieving the bank details of a specific swift code.
type GetSwiftCodeBankDetailsResponse struct {
	ID            string `json:"id"`
	Bank          string `json:"bank"`
	City          string `json:"city"`
	Branch        string `json:"branch"`
	SwiftCode     string `json:"swiftCode"`
	SwiftCodeLink string `json:"swiftCodeLink"`
	Country       string `json:"country"`
	CountrySlug   string `json:"countrySlug"`
}

// NewClient creates a new available client.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg: cfg.ToRequestConfig(),
	}
}

// GetBankDetails retrieves the bank details configuration for a specific rail.
func (c *Client) GetBankDetails(ctx context.Context, rail types.Rail, opts ...option.RequestOption) ([]types.BankDetail, error) {
	if rail == "" {
		return nil, fmt.Errorf("rail cannot be empty")
	}

	var q request.Query
	q.String("rail", string(rail))
	path := q.Apply("/available/bank-details")
	return request.Do[[]types.BankDetail](c.cfg, ctx, "GET", path, nil, opts...)
}

// GetRails retrieves all available payment rails.
func (c *Client) GetRails(ctx context.Context, opts ...option.RequestOption) ([]types.RailEntry, error) {
	return request.Do[[]types.RailEntry](c.cfg, ctx, "GET", "/available/rails", nil, opts...)
}

// NaicsCode represents a NAICS business industry code.
type NaicsCode struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// GetNaicsCodes retrieves the available NAICS business industry codes.
func (c *Client) GetNaicsCodes(ctx context.Context, opts ...option.RequestOption) ([]NaicsCode, error) {
	return request.Do[[]NaicsCode](c.cfg, ctx, "GET", "/available/naics", nil, opts...)
}

// GetSwiftCodeBankDetails retrieves the bank details of a specific swift code.
func (c *Client) GetSwiftCodeBankDetails(ctx context.Context, swift string, opts ...option.RequestOption) ([]GetSwiftCodeBankDetailsResponse, error) {
	if swift == "" {
		return nil, fmt.Errorf("swift code cannot be empty")
	}

	path := request.Pathf("/available/swift/%s", swift)

	return request.Do[[]GetSwiftCodeBankDetailsResponse](c.cfg, ctx, "GET", path, nil, opts...)
}
