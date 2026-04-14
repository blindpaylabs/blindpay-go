package sandbox

import (
	"context"
	"fmt"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
)

// SandboxStatus represents the status of a sandbox item.
type SandboxStatus string

const (
	SandboxStatusActive   SandboxStatus = "active"
	SandboxStatusInactive SandboxStatus = "inactive"
)

// Sandbox represents a sandbox item.
type Sandbox struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Status    SandboxStatus     `json:"status"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt *time.Time        `json:"created_at,omitempty"`
	UpdatedAt *time.Time        `json:"updated_at,omitempty"`
}

// CreateParams represents parameters for creating a sandbox item.
type CreateParams struct {
	Name     string            `json:"name"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// CreateResponse represents the response when creating a sandbox item.
type CreateResponse struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Status   SandboxStatus     `json:"status"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// UpdateParams represents parameters for updating a sandbox item.
type UpdateParams struct {
	ID       string            `json:"-"`
	Name     *string           `json:"name,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// Client handles sandbox-related operations.
type Client struct {
	cfg        *request.Config
	instanceID string
}

// NewClient creates a new sandbox client.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:        cfg.ToRequestConfig(),
		instanceID: cfg.InstanceID,
	}
}

// List retrieves all sandbox items for the instance.
func (c *Client) List(ctx context.Context) ([]Sandbox, error) {
	path := fmt.Sprintf("/instances/%s/sandbox", c.instanceID)
	return request.Do[[]Sandbox](c.cfg, ctx, "GET", path, nil)
}

// Get retrieves a specific sandbox item by ID.
func (c *Client) Get(ctx context.Context, id string) (*Sandbox, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/sandbox/%s", c.instanceID, id)
	return request.Do[*Sandbox](c.cfg, ctx, "GET", path, nil)
}

// Create creates a new sandbox item.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*CreateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/sandbox", c.instanceID)
	return request.Do[*CreateResponse](c.cfg, ctx, "POST", path, params)
}

// Update updates a sandbox item.
func (c *Client) Update(ctx context.Context, params *UpdateParams) (*Sandbox, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.ID == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/sandbox/%s", c.instanceID, params.ID)
	return request.Do[*Sandbox](c.cfg, ctx, "PATCH", path, params)
}

// Delete deletes a sandbox item.
func (c *Client) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/sandbox/%s", c.instanceID, id)
	_, err := request.Do[struct{}](c.cfg, ctx, "DELETE", path, nil)
	return err
}
