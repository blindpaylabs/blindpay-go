package instances

import (
	"context"
	"fmt"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
)

// InstanceMemberRole represents the role of an instance member.
type InstanceMemberRole string

const (
	InstanceMemberRoleOwner      InstanceMemberRole = "owner"
	InstanceMemberRoleAdmin      InstanceMemberRole = "admin"
	InstanceMemberRoleFinance    InstanceMemberRole = "finance"
	InstanceMemberRoleChecker    InstanceMemberRole = "checker"
	InstanceMemberRoleOperations InstanceMemberRole = "operations"
	InstanceMemberRoleDeveloper  InstanceMemberRole = "developer"
	InstanceMemberRoleViewer     InstanceMemberRole = "viewer"
)

// InstanceMember represents a member of an instance.
type InstanceMember struct {
	ID         string             `json:"id"`
	Email      string             `json:"email"`
	FirstName  string             `json:"first_name"`
	MiddleName string             `json:"middle_name"`
	LastName   string             `json:"last_name"`
	ImageURL   string             `json:"image_url"`
	CreatedAt  time.Time          `json:"created_at"`
	Role       InstanceMemberRole `json:"role"`
}

// UpdateParams represents parameters for updating an instance.
type UpdateParams struct {
	Name                      string  `json:"name"`
	ReceiverInviteRedirectURL *string `json:"receiver_invite_redirect_url,omitempty"`
}

// UpdateMemberRoleParams represents parameters for updating a member's role.
type UpdateMemberRoleParams struct {
	MemberID string             `json:"-"`
	Role     InstanceMemberRole `json:"role"`
}

// Client handles instance-related operations.
type Client struct {
	cfg        *request.Config
	instanceID string
}

// NewClient creates a new instances client.
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

// GetMembers retrieves all members of the instance.
func (c *Client) GetMembers(ctx context.Context) ([]InstanceMember, error) {
	path := fmt.Sprintf("/instances/%s/members", c.instanceID)
	return request.Do[[]InstanceMember](c.cfg, ctx, "GET", path, nil)
}

// Update updates the instance settings.
func (c *Client) Update(ctx context.Context, params *UpdateParams) error {
	if params == nil {
		return fmt.Errorf("params cannot be nil")
	}

	path := fmt.Sprintf("/instances/%s", c.instanceID)

	body := map[string]interface{}{
		"name": params.Name,
	}

	if params.ReceiverInviteRedirectURL != nil {
		body["receiver_invite_redirect_url"] = params.ReceiverInviteRedirectURL
	}

	_, err := request.Do[struct{}](c.cfg, ctx, "PUT", path, body)
	return err
}

// Delete deletes the instance.
func (c *Client) Delete(ctx context.Context) error {
	path := fmt.Sprintf("/instances/%s", c.instanceID)
	_, err := request.Do[struct{}](c.cfg, ctx, "DELETE", path, nil)
	return err
}

// DeleteMember removes a member from the instance.
func (c *Client) DeleteMember(ctx context.Context, memberID string) error {
	if memberID == "" {
		return fmt.Errorf("member ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/members/%s", c.instanceID, memberID)
	_, err := request.Do[struct{}](c.cfg, ctx, "DELETE", path, nil)
	return err
}

// UpdateMemberRole updates a member's role in the instance.
func (c *Client) UpdateMemberRole(ctx context.Context, params *UpdateMemberRoleParams) error {
	if params == nil {
		return fmt.Errorf("params cannot be nil")
	}
	if params.MemberID == "" {
		return fmt.Errorf("member ID cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/members/%s", c.instanceID, params.MemberID)

	body := map[string]interface{}{
		"role": params.Role,
	}

	_, err := request.Do[struct{}](c.cfg, ctx, "PUT", path, body)
	return err
}
