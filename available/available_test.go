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
		APIKey:     "test-key",
		InstanceID: "in_000000000000",
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`[
					{
						"label": "PIX Key",
						"regex": "",
						"key": "pix_key",
						"required": true
					}
				]`),
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
		APIKey:     "test-key",
		InstanceID: "in_000000000000",
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`[
					{
						"label": "Domestic Wire",
						"value": "wire",
						"country": "US"
					},
					{
						"label": "ACH",
						"value": "ach",
						"country": "US"
					},
					{
						"label": "PIX",
						"value": "pix",
						"country": "BR"
					},
					{
						"label": "SPEI",
						"value": "spei_bitso",
						"country": "MX"
					},
					{
						"label": "Transfers 3.0",
						"value": "transfers_bitso",
						"country": "AR"
					},
					{
						"label": "ACH Colombia",
						"value": "ach_cop_bitso",
						"country": "CO"
					},
					{
						"label": "International Swift",
						"value": "international_swift",
						"country": "US"
					},
					{
						"label": "RTP",
						"value": "rtp",
						"country": "US"
					}
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
	require.Len(t, rails, 8)
	require.Equal(t, "wire", rails[0].Value)
	require.Equal(t, "Domestic Wire", rails[0].Label)
	require.Equal(t, "US", rails[0].Country)

	require.Equal(t, "ach", rails[1].Value)
	require.Equal(t, "ACH", rails[1].Label)
	require.Equal(t, "US", rails[1].Country)

	require.Equal(t, "pix", rails[2].Value)
	require.Equal(t, "PIX", rails[2].Label)
	require.Equal(t, "BR", rails[2].Country)

	require.Equal(t, "spei_bitso", rails[3].Value)
	require.Equal(t, "SPEI", rails[3].Label)
	require.Equal(t, "MX", rails[3].Country)

	require.Equal(t, "transfers_bitso", rails[4].Value)
	require.Equal(t, "Transfers 3.0", rails[4].Label)
	require.Equal(t, "AR", rails[4].Country)

	require.Equal(t, "ach_cop_bitso", rails[5].Value)
	require.Equal(t, "ACH Colombia", rails[5].Label)
	require.Equal(t, "CO", rails[5].Country)

	require.Equal(t, "international_swift", rails[6].Value)
	require.Equal(t, "International Swift", rails[6].Label)
	require.Equal(t, "US", rails[6].Country)

	require.Equal(t, "rtp", rails[7].Value)
	require.Equal(t, "RTP", rails[7].Label)
	require.Equal(t, "US", rails[7].Country)
}
