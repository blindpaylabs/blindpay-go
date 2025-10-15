package apikeys

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/stretchr/testify/require"
)

func TestAPIKeys_Create(t *testing.T) {
	instanceID := "in_000000000000"
	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"name": "test",
					"permission": "full_access"
				}`),
				Out: json.RawMessage(`{
					"id": "ap_000000000000",
					"token": "token"
				}`),
				Method: http.MethodPost,
				Path:   "/instances/in_000000000000/api-keys",
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	apiKey, err := client.Create(context.Background(), &CreateParams{
		Name:       "test",
		Permission: "full_access",
	})
	require.NoError(t, err)
	require.Equal(t, "ap_000000000000", apiKey.ID)
	require.Equal(t, "token", apiKey.Token)
}

func TestAPIKeys_List(t *testing.T) {
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
						"id": "ap_000000000000",
						"token": "token",
						"name": "test",
						"permission": "full_access",
						"ip_whitelist": ["127.0.0.1"],
						"unkey_id": "key_123456789",
						"last_used_at": "2024-01-01T00:00:00.000Z",
						"instance_id": "in_000000000000",
						"created_at": "2021-01-01T00:00:00Z",
						"updated_at": "2021-01-01T00:00:00Z"
					}
				]`),
				Method: http.MethodGet,
				Path:   "/instances/in_000000000000/api-keys",
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	keys, err := client.List(context.Background())
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.Equal(t, "ap_000000000000", keys[0].ID)
	require.Equal(t, "test", keys[0].Name)
	require.Equal(t, "token", keys[0].Token)
	require.Equal(t, "full_access", keys[0].Permission)
	require.Equal(t, []string{"127.0.0.1"}, keys[0].IPWhitelist)
	require.Equal(t, "key_123456789", keys[0].UnkeyID)
	require.NotNil(t, keys[0].LastUsedAt)
	require.Equal(t, "2024-01-01T00:00:00.000Z", keys[0].LastUsedAt.Format("2006-01-02T15:04:05.000Z"))
	require.Equal(t, "in_000000000000", keys[0].InstanceID)
	require.Equal(t, "2021-01-01T00:00:00Z", keys[0].CreatedAt.Format("2006-01-02T15:04:05Z"))
	require.Equal(t, "2021-01-01T00:00:00Z", keys[0].UpdatedAt.Format("2006-01-02T15:04:05Z"))
}

func TestAPIKeys_Get(t *testing.T) {
	id := "ap_000000000000"
	instanceID := "in_000000000000"
	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
					"id": "ap_000000000000",
					"token": "token",
					"name": "test",
					"permission": "full_access",
					"ip_whitelist": ["127.0.0.1"],
					"unkey_id": "key_123456789",
					"last_used_at": "2024-01-01T00:00:00.000Z",
					"instance_id": "in_000000000000",
					"created_at": "2021-01-01T00:00:00Z",
					"updated_at": "2021-01-01T00:00:00Z"
				}`),
				Method: http.MethodGet,
				Path:   "/instances/in_000000000000/api-keys/ap_000000000000",
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	apiKey, err := client.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, "ap_000000000000", apiKey.ID)
	require.Equal(t, "token", apiKey.Token)
	require.Equal(t, "test", apiKey.Name)
	require.Equal(t, "full_access", apiKey.Permission)
	require.Equal(t, []string{"127.0.0.1"}, apiKey.IPWhitelist)
	require.Equal(t, "key_123456789", apiKey.UnkeyID)
	require.NotNil(t, apiKey.LastUsedAt)
	require.Equal(t, "2024-01-01T00:00:00.000Z", apiKey.LastUsedAt.Format("2006-01-02T15:04:05.000Z"))
	require.Equal(t, "in_000000000000", apiKey.InstanceID)
	require.Equal(t, "2021-01-01T00:00:00Z", apiKey.CreatedAt.Format("2006-01-02T15:04:05Z"))
	require.Equal(t, "2021-01-01T00:00:00Z", apiKey.UpdatedAt.Format("2006-01-02T15:04:05Z"))
}

func TestAPIKeys_Delete(t *testing.T) {
	id := "ap_000000000000"
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
				Path:   "/instances/in_000000000000/api-keys/ap_000000000000",
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	err := client.Delete(context.Background(), id)
	require.NoError(t, err)
}

func TestAPIKeys_Create_Error(t *testing.T) {
	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: "in_000000000000",
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Status: http.StatusBadRequest,
				Out: json.RawMessage(`{
					"error": {
						"message": "Invalid parameters",
						"type": "invalid_request_error"
					}
				}`),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	_, err := client.Create(context.Background(), &CreateParams{})
	require.Error(t, err)
}
