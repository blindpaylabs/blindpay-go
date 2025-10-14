package apikeys

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

func TestAPIKeys_Create(t *testing.T) {
	name := "Test API Key"
	permission := "read"
	id := "key_123"
	token := "bp_test_token_123"
	instanceID := "inst_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				In:     json.RawMessage(fmt.Sprintf(`{"name":"%s","permission":"%s"}`, name, permission)),
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","token":"%s"}`, id, token)),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/api-keys", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	apiKey, err := client.Create(context.Background(), &CreateParams{
		Name:       name,
		Permission: permission,
	})
	require.NoError(t, err)
	require.Equal(t, id, apiKey.ID)
	require.Equal(t, token, apiKey.Token)
}

func TestAPIKeys_List(t *testing.T) {
	instanceID := "inst_123"
	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`[
					{"id":"key_123","name":"Test Key 1","permission":"read","token":"","unkey_id":"unkey_1","instance_id":"inst_123","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z"},
					{"id":"key_456","name":"Test Key 2","permission":"write","token":"","unkey_id":"unkey_2","instance_id":"inst_123","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z"}
				]`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/api-keys", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	keys, err := client.List(context.Background())
	require.NoError(t, err)
	require.Len(t, keys, 2)
	require.Equal(t, "key_123", keys[0].ID)
	require.Equal(t, "Test Key 1", keys[0].Name)
}

func TestAPIKeys_Get(t *testing.T) {
	id := "key_123"
	instanceID := "inst_123"
	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","name":"Test Key","permission":"read","token":"","unkey_id":"unkey_1","instance_id":"inst_123","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z"}`, id)),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/api-keys/%s", instanceID, id),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	apiKey, err := client.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, apiKey.ID)
	require.Equal(t, "Test Key", apiKey.Name)
}

func TestAPIKeys_Delete(t *testing.T) {
	id := "key_123"
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
				Path:   fmt.Sprintf("/instances/%s/api-keys/%s", instanceID, id),
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
		APIKey:     "test_key",
		InstanceID: "inst_123",
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
