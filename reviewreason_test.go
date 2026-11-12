package blindpay

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/stretchr/testify/require"
)

const fixtureReviewReason = "The bank statement does not cover the last 3 months. Send statements from August to October showing the income that supports the requested limit."

// TestReviewReason_GetLimitIncrease decodes review_reason on the GET
// limit-increase response: populated on the full variant, nil on
// nullables_null and required_only.
func TestReviewReason_GetLimitIncrease(t *testing.T) {
	const op = "GET /v1/instances/{instance_id}/customers/{customer_id}/limit-increase"
	cases := map[string]*string{
		"full":           ptr(fixtureReviewReason),
		"nullables_null": nil,
		"required_only":  nil,
	}
	for name, want := range cases {
		t.Run(name, func(t *testing.T) {
			v := blindpaytest.FindVariant(t, op, name)
			c, err := New("test_key", "in_000000000000",
				WithBaseURL("https://api.blindpay.com/v1"),
				WithHTTPClient(&http.Client{Transport: &blindpaytest.RoundTripper{
					T:      t,
					Out:    v.Body,
					Method: http.MethodGet,
					Path:   "/v1/instances/in_000000000000/customers/cu_000000000000/limit-increase",
				}}))
			require.NoError(t, err)

			got, err := c.Customers.GetLimitIncreaseRequests(context.Background(), "cu_000000000000")
			require.NoError(t, err)
			require.NotEmpty(t, got)
			require.Equal(t, want, got[0].ReviewReason)
		})
	}
}

// TestReviewReason_LimitIncreaseWebhooks decodes review_reason on both
// limitIncrease webhooks through ParseWebhookEvent.
func TestReviewReason_LimitIncreaseWebhooks(t *testing.T) {
	for _, event := range []string{"limitIncrease.new", "limitIncrease.update"} {
		var fixture struct {
			Variants map[string]json.RawMessage `json:"variants"`
		}
		require.NoError(t, json.Unmarshal(blindpaytest.ReadFixture(t, "webhooks/"+event+".json"), &fixture))

		for name, want := range map[string]*string{
			"full":           ptr(fixtureReviewReason),
			"nullables_null": nil,
			"required_only":  nil,
		} {
			t.Run(event+"/"+name, func(t *testing.T) {
				body, ok := fixture.Variants[name]
				require.True(t, ok, "variant %s missing", name)

				payload, err := ParseWebhookEvent(body)
				require.NoError(t, err)
				wh, ok := payload.(*LimitIncreaseWebhook)
				require.True(t, ok, "got %T", payload)
				require.Equal(t, WebhookEvent(event), wh.EventType())
				require.Equal(t, want, wh.ReviewReason)
			})
		}
	}
}
