package webhookendpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestWebhookEndpoints_Create(t *testing.T) {
	instanceID := "in_000000000000"
	id := "we_000000000000"
	url := "https://example.com/webhook"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"url":"https://example.com/webhook",
					"events":["receiver.new"]
				}`),
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s"}`, id)),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/webhook-endpoints", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.Create(context.Background(), &CreateParams{
		URL: url,
		Events: []types.WebhookEvent{
			types.WebhookEventReceiverNew,
		},
	})
	require.NoError(t, err)
	require.Equal(t, id, response.ID)
}

func TestWebhookEndpoints_List(t *testing.T) {
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
						"id":"we_000000000000",
						"url":"https://example.com/webhook",
						"events":["receiver.new"],
						"last_event_at":"2024-01-01T00:00:00.000Z",
						"instance_id":"in_000000000000",
						"created_at":"2021-01-01T00:00:00Z",
						"updated_at":"2021-01-01T00:00:00Z"
					}
				]`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/webhook-endpoints", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	endpoints, err := client.List(context.Background())
	require.NoError(t, err)
	require.Len(t, endpoints, 1)
	require.Equal(t, "we_000000000000", endpoints[0].ID)
}

func TestWebhookEndpoints_Delete(t *testing.T) {
	instanceID := "in_000000000000"
	id := "we_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{"data":null}`),
				Method: http.MethodDelete,
				Path:   fmt.Sprintf("/instances/%s/webhook-endpoints/%s", instanceID, id),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	err := client.Delete(context.Background(), id)
	require.NoError(t, err)
}

func TestWebhookEndpoints_GetSecret(t *testing.T) {
	instanceID := "in_000000000000"
	id := "we_000000000000"
	secret := "whsec_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(fmt.Sprintf(`{"key":"%s"}`, secret)),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/webhook-endpoints/%s/secret", instanceID, id),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.GetSecret(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, secret, response.Key)
}

func TestWebhookEndpoints_GetPortalAccessURL(t *testing.T) {
	instanceID := "in_000000000000"
	portalURL := "https://example.com/webhook"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(fmt.Sprintf(`{"url":"%s"}`, portalURL)),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/webhook-endpoints/portal-access", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.GetPortalAccessURL(context.Background())
	require.NoError(t, err)
	require.Equal(t, portalURL, response.URL)
}
