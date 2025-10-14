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
	instanceID := "in_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`[
					{
						"id": "us_000000000000",
						"email": "email@example.com",
						"first_name": "Harry",
						"middle_name": "James",
						"last_name": "Potter",
						"image_url": "https://example.com/image.png",
						"created_at": "2021-01-01T00:00:00Z",
						"role": "admin"
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
	require.Len(t, members, 1)
	require.Equal(t, "us_000000000000", members[0].ID)
	require.Equal(t, "email@example.com", members[0].Email)
	require.Equal(t, InstanceMemberRoleAdmin, members[0].Role)
}

func TestInstances_Update(t *testing.T) {
	instanceID := "in_000000000000"
	redirectURL := "https://example.com/invite"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{"data":null}`),
				Method: http.MethodPut,
				Path:   fmt.Sprintf("/instances/%s", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	err := client.Update(context.Background(), &UpdateParams{
		Name:                      "New Instance Name",
		ReceiverInviteRedirectURL: &redirectURL,
	})
	require.NoError(t, err)
}

func TestInstances_Delete(t *testing.T) {
	instanceID := "in_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{"data":null}`),
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
	instanceID := "in_000000000000"
	memberID := "us_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{"data":null}`),
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
	instanceID := "in_000000000000"
	memberID := "us_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{"data":null}`),
				Method: http.MethodPut,
				Path:   fmt.Sprintf("/instances/%s/members/%s", instanceID, memberID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	err := client.UpdateMemberRole(context.Background(), &UpdateMemberRoleParams{
		MemberID: memberID,
		Role:     InstanceMemberRoleChecker,
	})
	require.NoError(t, err)
}
