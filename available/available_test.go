package available

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestAvailable_GetBankDetails(t *testing.T) {
	rail := "pix"
	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: "inst_123",
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`[{"label":"PIX Key","regex":"^.+$","key":"pix_key","required":true}]`),
				Method: http.MethodGet,
				Path:   "/available/bank-details",
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	details, err := client.GetBankDetails(context.Background(), types.Rail(rail))
	require.NoError(t, err)
	require.Len(t, details, 1)
	require.Equal(t, "pix_key", details[0].Key)
	require.Equal(t, "PIX Key", details[0].Label)
	require.True(t, details[0].Required)
}

func TestAvailable_GetRails(t *testing.T) {
	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: "inst_123",
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`[
					{"label":"PIX","value":"pix","country":"BR"},
					{"label":"ACH","value":"ach","country":"US"}
				]`),
				Method: http.MethodGet,
				Path:   "/available/rails",
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	rails, err := client.GetRails(context.Background())
	require.NoError(t, err)
	require.Len(t, rails, 2)
	require.Equal(t, "pix", rails[0].Value)
	require.Equal(t, "PIX", rails[0].Label)
	require.Equal(t, "BR", rails[0].Country)
}
