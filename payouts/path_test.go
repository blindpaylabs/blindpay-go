package payouts

import (
	"context"
	"net/http"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/stretchr/testify/require"
)

// TestPayouts_PathSegmentsAreEscaped: an ID can never change the route.
func TestPayouts_PathSegmentsAreEscaped(t *testing.T) {
	var full []byte
	for _, v := range blindpaytest.ResponseVariants(t, "GET /v1/instances/{instance_id}/payouts/{id}") {
		if v.Name == "full" {
			full = v.Body
		}
	}
	require.NotNil(t, full)

	tests := []struct {
		id   string
		path string
	}{
		{id: "po_000000000000", path: "/v1/instances/in_000000000000/payouts/po_000000000000"},
		{id: "a/b", path: "/v1/instances/in_000000000000/payouts/a%2Fb"},
		{id: "../customers", path: "/v1/instances/in_000000000000/payouts/..%2Fcustomers"},
		{id: "a?b", path: "/v1/instances/in_000000000000/payouts/a%3Fb"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.id, func(t *testing.T) {
			client := NewClient(&config.Config{
				BaseURL:    "https://api.blindpay.com/v1",
				APIKey:     "test_key",
				InstanceID: "in_000000000000",
				HTTPClient: &http.Client{Transport: &blindpaytest.RoundTripper{
					T:      t,
					Method: http.MethodGet,
					Path:   tc.path,
					Out:    full,
				}},
				UserAgent: "blindpay-go/test",
			})
			_, err := client.Get(context.Background(), tc.id)
			require.NoError(t, err)
		})
	}
}
