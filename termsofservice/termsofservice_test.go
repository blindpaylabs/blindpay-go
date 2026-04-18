package termsofservice

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

func TestTermsOfService_Initiate(t *testing.T) {
	instanceID := "in_000000000000"
	idempotencyKey := "123e4567-e89b-12d3-a456-426614174000"
	url := "https://app.blindpay.com/e/terms-of-service?session_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"idempotency_key":"123e4567-e89b-12d3-a456-426614174000",
					"receiver_id":null,
					"redirect_url":null
				}`),
				Out: json.RawMessage(`{
					"url":"https://app.blindpay.com/e/terms-of-service?session_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/e/instances/%s/tos", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.Initiate(context.Background(), &InitiateParams{
		IdempotencyKey: idempotencyKey,
		ReceiverID:     nil,
		RedirectURL:    nil,
	})
	require.NoError(t, err)
	require.Equal(t, url, response.URL)
}
