package available

import (
	"context"
	"fmt"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
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
		cfg: &request.Config{
			BaseURL:    cfg.BaseURL,
			APIKey:     cfg.APIKey,
			HTTPClient: cfg.HTTPClient,
			UserAgent:  cfg.UserAgent,
		},
	}
}

// GetBankDetails retrieves the bank details configuration for a specific rail.
func (c *Client) GetBankDetails(ctx context.Context, rail types.Rail) ([]types.BankDetail, error) {
	if rail == "" {
		return nil, fmt.Errorf("rail cannot be empty")
	}

	path := fmt.Sprintf("/available/bank-details?rail=%s", rail)
	return request.Do[[]types.BankDetail](c.cfg, ctx, "GET", path, nil)
}

// GetRails retrieves all available payment rails.
func (c *Client) GetRails(ctx context.Context) ([]types.RailEntry, error) {
	return request.Do[[]types.RailEntry](c.cfg, ctx, "GET", "/available/rails", nil)
}

// GetSwiftCodeBankDetails retrieves the bank details of a specific swift code.
func (c *Client) GetSwiftCodeBankDetails(ctx context.Context, swift string) ([]GetSwiftCodeBankDetailsResponse, error) {
	if swift == "" {
		return nil, fmt.Errorf("swift code cannot be empty")
	}

	path := fmt.Sprintf("/available/swift/%s", swift)

	return request.Do[[]GetSwiftCodeBankDetailsResponse](c.cfg, ctx, "GET", path, nil)
}
