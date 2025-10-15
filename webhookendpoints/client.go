package webhookendpoints

import (
	"context"
	"fmt"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/internal/types"
)

// WebhookEndpoint represents a webhook endpoint.
type WebhookEndpoint struct {
	ID          string               `json:"id"`
	URL         string               `json:"url"`
	Events      []types.WebhookEvent `json:"events"`
	LastEventAt time.Time            `json:"last_event_at"`
	InstanceID  string               `json:"instance_id"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

// CreateParams represents parameters for creating a webhook endpoint.
type CreateParams struct {
	URL    string               `json:"url"`
	Events []types.WebhookEvent `json:"events"`
}

// CreateResponse represents the response when creating a webhook endpoint.
type CreateResponse struct {
	ID string `json:"id"`
}

// GetSecretResponse represents the webhook endpoint secret response.
type GetSecretResponse struct {
	Key string `json:"key"`
}

// GetPortalAccessURLResponse represents the portal access URL response.
type GetPortalAccessURLResponse struct {
	URL string `json:"url"`
}

// Client handles webhook endpoint-related operations.
type Client struct {
	cfg        *request.Config
	instanceID string
}

// NewClient creates a new webhook endpoints client.
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

// List retrieves all webhook endpoints for the instance.
func (c *Client) List(ctx context.Context) ([]WebhookEndpoint, error) {
	path := fmt.Sprintf("/instances/%s/webhook-endpoints", c.instanceID)
	return request.Do[[]WebhookEndpoint](c.cfg, ctx, "GET", path, nil)
}

// Create creates a new webhook endpoint.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*CreateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := fmt.Sprintf("/instances/%s/webhook-endpoints", c.instanceID)
	return request.Do[*CreateResponse](c.cfg, ctx, "POST", path, params)
}

// Delete deletes a webhook endpoint.
func (c *Client) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/webhook-endpoints/%s", c.instanceID, id)
	_, err := request.Do[struct{}](c.cfg, ctx, "DELETE", path, nil)
	return err
}

// GetSecret retrieves the secret for a webhook endpoint.
func (c *Client) GetSecret(ctx context.Context, id string) (*GetSecretResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/webhook-endpoints/%s/secret", c.instanceID, id)
	return request.Do[*GetSecretResponse](c.cfg, ctx, "GET", path, nil)
}

// GetPortalAccessURL retrieves the portal access URL.
func (c *Client) GetPortalAccessURL(ctx context.Context) (*GetPortalAccessURLResponse, error) {
	path := fmt.Sprintf("/instances/%s/webhook-endpoints/portal-access", c.instanceID)
	return request.Do[*GetPortalAccessURLResponse](c.cfg, ctx, "GET", path, nil)
}
