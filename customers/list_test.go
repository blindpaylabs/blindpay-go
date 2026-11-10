package customers

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

// TestCustomers_ListWithParamsShapes decodes every shared fixture variant
// through ListWithParams. The API returns the {data, pagination} envelope
// (.b0) only when limit/starting_after/ending_before is sent; with filters
// alone it returns a bare array (.b1), which ListWithParams used to reject
// (blindpay-v2 payouts/payins/customers list handlers).
func TestCustomers_ListWithParamsShapes(t *testing.T) {
	instanceID := "in_000000000000"

	for _, v := range blindpaytest.ResponseVariants(t, "GET /v1/instances/{instance_id}/customers") {
		v := v
		t.Run(v.Name, func(t *testing.T) {
			params := &ListParams{Country: "BR"}
			query := url.Values{"country": {"BR"}}
			if v.IsEnvelope() {
				params.Limit = "100"
				query.Set("limit", "100")
			}

			client := NewClient(&config.Config{
				BaseURL:    "https://api.blindpay.com/v1",
				APIKey:     "test_key",
				InstanceID: instanceID,
				HTTPClient: &http.Client{Transport: &blindpaytest.RoundTripper{
					T:      t,
					Method: http.MethodGet,
					Path:   fmt.Sprintf("/v1/instances/%s/customers", instanceID),
					Query:  query,
					Out:    v.Body,
				}},
				UserAgent: "blindpay-go/test",
			})

			resp, err := client.ListWithParams(context.Background(), params)
			require.NoError(t, err)
			require.Len(t, resp.Data, v.RowCount(t))
			require.NotEmpty(t, resp.Data)
			require.NotEmpty(t, resp.Data[0].ID)
			if !v.IsEnvelope() {
				require.False(t, resp.Pagination.HasMore)
				require.Nil(t, resp.Pagination.NextPage)
			}
		})
	}
}

func TestCustomers_ListWithParamsEnvelopePagination(t *testing.T) {
	client := NewClient(&config.Config{
		BaseURL:    "https://api.blindpay.com/v1",
		APIKey:     "test_key",
		InstanceID: "in_000000000000",
		HTTPClient: &http.Client{Transport: &blindpaytest.RoundTripper{
			T:     t,
			Query: url.Values{"limit": {"10"}, "starting_after": {"cu_000000000001"}},
			Out:   []byte(`{"data":[{"id":"cu_000000000002"}],"pagination":{"has_more":true,"next_page":"cu_000000000002","prev_page":"cu_000000000001"}}`),
		}},
		UserAgent: "blindpay-go/test",
	})

	resp, err := client.ListWithParams(context.Background(), &ListParams{Limit: "10", StartingAfter: "cu_000000000001"})
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)
	require.True(t, resp.Pagination.HasMore)
	require.NotNil(t, resp.Pagination.NextPage)
	require.Equal(t, "cu_000000000002", *resp.Pagination.NextPage)
	require.NotNil(t, resp.Pagination.PrevPage)
	require.Equal(t, "cu_000000000001", *resp.Pagination.PrevPage)
}

func TestCustomers_ListShapes(t *testing.T) {
	for _, v := range blindpaytest.ResponseVariants(t, "GET /v1/instances/{instance_id}/customers") {
		v := v
		if v.IsEnvelope() {
			continue // List sends no pagination params, so the API only returns the bare array.
		}
		t.Run(v.Name, func(t *testing.T) {
			client := NewClient(&config.Config{
				BaseURL:    "https://api.blindpay.com/v1",
				APIKey:     "test_key",
				InstanceID: "in_000000000000",
				HTTPClient: &http.Client{Transport: &blindpaytest.RoundTripper{
					T:     t,
					Query: url.Values{},
					Out:   v.Body,
				}},
				UserAgent: "blindpay-go/test",
			})

			customers, err := client.List(context.Background())
			require.NoError(t, err)
			require.Len(t, customers, v.RowCount(t))
		})
	}
}
