package fees

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

func TestFees_Get(t *testing.T) {
	instanceID := "in_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
					"id":"fee_000000000000",
					"instance_id":"in_000000000000",
					"ach":{"payin_flat":0,"payin_percentage":0,"payout_flat":5,"payout_percentage":0},
					"domestic_wire":{"payin_flat":0,"payin_percentage":0,"payout_flat":25,"payout_percentage":0},
					"rtp":{"payin_flat":0,"payin_percentage":0,"payout_flat":1,"payout_percentage":0},
					"international_swift":{"payin_flat":0,"payin_percentage":0,"payout_flat":35,"payout_percentage":0},
					"pix":{"payin_flat":0,"payin_percentage":0,"payout_flat":2,"payout_percentage":0},
					"pix_safe":{"payin_flat":0,"payin_percentage":0,"payout_flat":3,"payout_percentage":0},
					"ach_colombia":{"payin_flat":0,"payin_percentage":0,"payout_flat":5,"payout_percentage":0},
					"transfers_3":{"payin_flat":0,"payin_percentage":0,"payout_flat":4,"payout_percentage":0},
					"spei":{"payin_flat":0,"payin_percentage":0,"payout_flat":3,"payout_percentage":0},
					"tron":{"payin_flat":0,"payin_percentage":0,"payout_flat":1,"payout_percentage":0},
					"ethereum":{"payin_flat":0,"payin_percentage":0,"payout_flat":10,"payout_percentage":0},
					"polygon":{"payin_flat":0,"payin_percentage":0,"payout_flat":1,"payout_percentage":0},
					"base":{"payin_flat":0,"payin_percentage":0,"payout_flat":1,"payout_percentage":0},
					"arbitrum":{"payin_flat":0,"payin_percentage":0,"payout_flat":1,"payout_percentage":0},
					"stellar":{"payin_flat":0,"payin_percentage":0,"payout_flat":1,"payout_percentage":0},
					"solana":{"payin_flat":0,"payin_percentage":0,"payout_flat":1,"payout_percentage":0},
					"created_at":"2021-01-01T00:00:00Z",
					"updated_at":"2021-01-01T00:00:00Z"
				}`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/billing/fees", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, "fee_000000000000", response.ID)
	require.Equal(t, float64(5), response.ACH.PayoutFlat)
	require.Equal(t, float64(25), response.DomesticWire.PayoutFlat)
}
