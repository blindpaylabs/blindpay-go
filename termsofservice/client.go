package termsofservice

import (
	"context"
	"fmt"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
)

// InitiateParams represents parameters for initiating terms of service.
type InitiateParams struct {
	IdempotencyKey string  `json:"idempotency_key"`
	ReceiverID     *string `json:"receiver_id"`
	RedirectURL    *string `json:"redirect_url"`
}

// InitiateResponse represents the response from initiating terms of service.
type InitiateResponse struct {
	URL string `json:"url"`
}

// Client handles terms of service-related operations.
type Client struct {
	cfg        *request.Config
	instanceID string
}

// NewClient creates a new terms of service client.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg: &request.Config{
			BaseURL:    cfg.BaseURL,
			APIKey:     cfg.APIKey,
			HTTPClient: cfg.HTTPClient,
			UserAgent:  cfg.UserAgent,
		},
		instanceID: cfg.InstanceID,
	}
}

// Initiate initiates the terms of service flow and returns a URL.
func (c *Client) Initiate(ctx context.Context, params *InitiateParams) (*InitiateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.IdempotencyKey == "" {
		return nil, fmt.Errorf("idempotency key cannot be empty")
	}

	path := fmt.Sprintf("/e/instances/%s/tos", c.instanceID)
	return request.Do[*InitiateResponse](c.cfg, ctx, "POST", path, params)
}
