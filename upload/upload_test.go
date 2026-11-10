package upload

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/stretchr/testify/require"
)

func newTestClient(rt *blindpaytest.RoundTripper) *Client {
	return NewClient(&config.Config{
		BaseURL:    "https://api.blindpay.com/v1",
		APIKey:     "test_key",
		InstanceID: "in_000000000000",
		HTTPClient: &http.Client{Transport: rt},
		UserAgent:  "blindpay-go/test",
	})
}

func TestUpload_Upload(t *testing.T) {
	client := newTestClient(&blindpaytest.RoundTripper{
		T:      t,
		Method: http.MethodPost,
		Path:   "/v1/upload",
		Query:  url.Values{"instance_id": {"in_000000000000"}},
		RequestHeader: http.Header{
			"Authorization": {"Bearer test_key"},
		},
		Out: json.RawMessage(`{"file_url":"https://example.com/file.png"}`),
	})

	resp, err := client.Upload(context.Background(), &UploadParams{
		File:     strings.NewReader("data"),
		FileName: "file.png",
		Bucket:   BucketOnboarding,
	})
	require.NoError(t, err)
	require.Equal(t, "https://example.com/file.png", resp.FileURL)
}

func TestUpload_UploadAPIError(t *testing.T) {
	f := blindpaytest.LoadErrorFixture(t, "api_error_envelope.json")
	client := newTestClient(&blindpaytest.RoundTripper{
		T:      t,
		Method: http.MethodPost,
		Path:   "/v1/upload",
		Status: http.StatusBadRequest,
		Out:    f.SpecGeneratedExample,
	})

	_, err := client.Upload(context.Background(), &UploadParams{
		File:     strings.NewReader("data"),
		FileName: "file.png",
		Bucket:   BucketOnboarding,
	})
	require.Error(t, err)

	var apiErr *request.APIError
	require.True(t, errors.As(err, &apiErr), "expected *APIError, got %T", err)
	require.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	require.Equal(t, "PAYOUTS_INSUFFICIENT_BALANCE", apiErr.Code)
	require.Equal(t, "not_enough_balance", apiErr.Message)
}
