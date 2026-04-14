package sandbox

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

const instanceID = "in_000000000000"

func newClient(t *testing.T, method, path string, out json.RawMessage, in json.RawMessage) *Client {
	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    out,
				In:     in,
				Method: method,
				Path:   path,
			},
		},
		UserAgent: "test",
	}
	return NewClient(cfg)
}

func TestSandbox_List(t *testing.T) {
	client := newClient(t,
		http.MethodGet,
		fmt.Sprintf("/instances/%s/sandbox", instanceID),
		json.RawMessage(`[{"id":"sb_000000000001","name":"Test Item","status":"active"}]`),
		nil,
	)

	items, err := client.List(context.Background())
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, "sb_000000000001", items[0].ID)
	require.Equal(t, "Test Item", items[0].Name)
	require.Equal(t, SandboxStatusActive, items[0].Status)
}

func TestSandbox_Get(t *testing.T) {
	sandboxID := "sb_000000000001"
	client := newClient(t,
		http.MethodGet,
		fmt.Sprintf("/instances/%s/sandbox/%s", instanceID, sandboxID),
		json.RawMessage(`{"id":"sb_000000000001","name":"Test Item","status":"active"}`),
		nil,
	)

	item, err := client.Get(context.Background(), sandboxID)
	require.NoError(t, err)
	require.Equal(t, sandboxID, item.ID)
	require.Equal(t, "Test Item", item.Name)
	require.Equal(t, SandboxStatusActive, item.Status)
}

func TestSandbox_Get_EmptyID(t *testing.T) {
	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{},
		UserAgent:  "test",
	}
	client := NewClient(cfg)

	_, err := client.Get(context.Background(), "")
	require.ErrorContains(t, err, "id cannot be empty")
}

func TestSandbox_Create(t *testing.T) {
	client := newClient(t,
		http.MethodPost,
		fmt.Sprintf("/instances/%s/sandbox", instanceID),
		json.RawMessage(`{"id":"sb_000000000001","name":"My Sandbox Item","status":"active"}`),
		json.RawMessage(`{"name":"My Sandbox Item"}`),
	)

	resp, err := client.Create(context.Background(), &CreateParams{
		Name: "My Sandbox Item",
	})
	require.NoError(t, err)
	require.Equal(t, "sb_000000000001", resp.ID)
	require.Equal(t, "My Sandbox Item", resp.Name)
	require.Equal(t, SandboxStatusActive, resp.Status)
}

func TestSandbox_Create_NilParams(t *testing.T) {
	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{},
		UserAgent:  "test",
	}
	client := NewClient(cfg)

	_, err := client.Create(context.Background(), nil)
	require.ErrorContains(t, err, "params cannot be nil")
}

func TestSandbox_Update(t *testing.T) {
	sandboxID := "sb_000000000001"
	name := "Updated Name"
	client := newClient(t,
		http.MethodPatch,
		fmt.Sprintf("/instances/%s/sandbox/%s", instanceID, sandboxID),
		json.RawMessage(`{"id":"sb_000000000001","name":"Updated Name","status":"active"}`),
		json.RawMessage(`{"name":"Updated Name"}`),
	)

	item, err := client.Update(context.Background(), &UpdateParams{
		ID:   sandboxID,
		Name: &name,
	})
	require.NoError(t, err)
	require.Equal(t, sandboxID, item.ID)
	require.Equal(t, "Updated Name", item.Name)
}

func TestSandbox_Delete(t *testing.T) {
	sandboxID := "sb_000000000001"
	client := newClient(t,
		http.MethodDelete,
		fmt.Sprintf("/instances/%s/sandbox/%s", instanceID, sandboxID),
		json.RawMessage(`{}`),
		nil,
	)

	err := client.Delete(context.Background(), sandboxID)
	require.NoError(t, err)
}

func TestSandbox_Delete_EmptyID(t *testing.T) {
	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{},
		UserAgent:  "test",
	}
	client := NewClient(cfg)

	err := client.Delete(context.Background(), "")
	require.ErrorContains(t, err, "id cannot be empty")
}
