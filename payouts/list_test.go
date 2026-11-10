package payouts

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/stretchr/testify/require"
)

// TestPayouts_ListShapes decodes every shared fixture variant through List.
// The API returns the {data, pagination} envelope (.b0) when limit is sent
// and a bare array (.b1) otherwise (blindpay-v2 payouts/payins/customers list handlers); before the fix, List failed on
// every .b0 variant.
func TestPayouts_ListShapes(t *testing.T) {
	instanceID := "in_000000000000"

	for _, v := range blindpaytest.ResponseVariants(t, "GET /v1/instances/{instance_id}/payouts") {
		v := v
		t.Run(v.Name, func(t *testing.T) {
			params := &ListParams{CustomerID: "re_000000000000"}
			query := url.Values{"customer_id": {"re_000000000000"}}
			if v.IsEnvelope() {
				params.Limit = 10
				query.Set("limit", "10")
			}

			client := NewClient(&config.Config{
				BaseURL:    "https://api.blindpay.com/v1",
				APIKey:     "test_key",
				InstanceID: instanceID,
				HTTPClient: &http.Client{Transport: &blindpaytest.RoundTripper{
					T:      t,
					Method: http.MethodGet,
					Path:   fmt.Sprintf("/v1/instances/%s/payouts", instanceID),
					Query:  query,
					Out:    v.Body,
				}},
				UserAgent: "blindpay-go/test",
			})

			payouts, err := client.List(context.Background(), params)
			require.NoError(t, err)
			require.Len(t, payouts, v.RowCount(t))
			require.NotEmpty(t, payouts)
			require.NotEmpty(t, payouts[0].ID)
		})
	}
}

func TestPayouts_ListNoParamsSendsNoQuery(t *testing.T) {
	client := NewClient(&config.Config{
		BaseURL:    "https://api.blindpay.com/v1",
		APIKey:     "test_key",
		InstanceID: "in_000000000000",
		HTTPClient: &http.Client{Transport: &blindpaytest.RoundTripper{
			T:     t,
			Query: url.Values{},
			Out:   []byte(`[]`),
		}},
		UserAgent: "blindpay-go/test",
	})

	payouts, err := client.List(context.Background(), nil)
	require.NoError(t, err)
	require.Empty(t, payouts)
}
