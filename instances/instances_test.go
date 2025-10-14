package instances

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/stretchr/testify/require"
)

func TestInstances_GetMembers(t *testing.T) {
	instanceID := "inst_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`[
					{
						"id":"mem_123",
						"email":"admin@example.com",
						"first_name":"John",
						"middle_name":"",
						"last_name":"Doe",
						"image_url":"https://example.com/avatar.jpg",
						"created_at":"2024-01-01T00:00:00Z",
						"role":"owner"
					},
					{
						"id":"mem_456",
						"email":"dev@example.com",
						"first_name":"Jane",
						"middle_name":"",
						"last_name":"Smith",
						"image_url":"https://example.com/avatar2.jpg",
						"created_at":"2024-01-02T00:00:00Z",
						"role":"developer"
					}
				]`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/members", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	members, err := client.GetMembers(context.Background())
	require.NoError(t, err)
	require.Len(t, members, 2)
	require.Equal(t, "mem_123", members[0].ID)
	require.Equal(t, "admin@example.com", members[0].Email)
	require.Equal(t, InstanceMemberRoleOwner, members[0].Role)
}

func TestInstances_Update(t *testing.T) {
	instanceID := "inst_123"
	redirectURL := "https://example.com/invite"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{}`),
				Method: http.MethodPut,
				Path:   fmt.Sprintf("/instances/%s", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	err := client.Update(context.Background(), &UpdateParams{
		Name:                      "My Updated Instance",
		ReceiverInviteRedirectURL: &redirectURL,
	})
	require.NoError(t, err)
}

func TestInstances_Delete(t *testing.T) {
	instanceID := "inst_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{}`),
				Method: http.MethodDelete,
				Path:   fmt.Sprintf("/instances/%s", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	err := client.Delete(context.Background())
	require.NoError(t, err)
}

func TestInstances_DeleteMember(t *testing.T) {
	instanceID := "inst_123"
	memberID := "mem_456"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{}`),
				Method: http.MethodDelete,
				Path:   fmt.Sprintf("/instances/%s/members/%s", instanceID, memberID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	err := client.DeleteMember(context.Background(), memberID)
	require.NoError(t, err)
}

func TestInstances_UpdateMemberRole(t *testing.T) {
	instanceID := "inst_123"
	memberID := "mem_456"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{}`),
				Method: http.MethodPut,
				Path:   fmt.Sprintf("/instances/%s/members/%s", instanceID, memberID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	err := client.UpdateMemberRole(context.Background(), &UpdateMemberRoleParams{
		MemberID: memberID,
		Role:     InstanceMemberRoleAdmin,
	})
	require.NoError(t, err)
}
