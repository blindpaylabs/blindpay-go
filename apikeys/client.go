package apikeys

import (
	"context"
	"fmt"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
)

// APIKey represents an API key.
type APIKey struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Permission  string     `json:"permission"`
	Token       string     `json:"token"`
	IPWhitelist []string   `json:"ip_whitelist,omitempty"`
	UnkeyID     string     `json:"unkey_id"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	InstanceID  string     `json:"instance_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateParams represents parameters for creating an API key.
type CreateParams struct {
	Name        string   `json:"name"`
	Permission  string   `json:"permission"`
	IPWhitelist []string `json:"ip_whitelist,omitempty"`
}

// CreateResponse represents the response when creating an API key.
type CreateResponse struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// Client handles API key-related operations.
type Client struct {
	cfg        *request.Config
	instanceID string
}

// NewClient creates a new API keys client.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:        cfg.ToRequestConfig(),
		instanceID: cfg.InstanceID,
	}
}

// List retrieves all API keys for the instance.
func (c *Client) List(ctx context.Context) ([]APIKey, error) {
	path := fmt.Sprintf("/instances/%s/api-keys", c.instanceID)
	return request.Do[[]APIKey](c.cfg, ctx, "GET", path, nil)
}

// Create creates a new API key.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*CreateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := fmt.Sprintf("/instances/%s/api-keys", c.instanceID)
	return request.Do[*CreateResponse](c.cfg, ctx, "POST", path, params)
}

// Get retrieves a specific API key by ID.
func (c *Client) Get(ctx context.Context, id string) (*APIKey, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/api-keys/%s", c.instanceID, id)
	return request.Do[*APIKey](c.cfg, ctx, "GET", path, nil)
}

// Delete deletes an API key.
func (c *Client) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/api-keys/%s", c.instanceID, id)
	_, err := request.Do[struct{}](c.cfg, ctx, "DELETE", path, nil)
	return err
}
