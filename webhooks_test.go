package blindpay

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	svix "github.com/svix/svix-webhooks/go"

	"github.com/stretchr/testify/require"
)

// generateTestSignature creates a valid Svix signature for testing
func generateTestSignature(t *testing.T, secret, id, payload string) (timestamp, signature string) {
	wh, err := svix.NewWebhook(secret)
	require.NoError(t, err)

	ts := time.Now()
	sig, err := wh.Sign(id, ts, []byte(payload))
	require.NoError(t, err)

	return fmt.Sprintf("%d", ts.Unix()), sig
}

func TestVerifyWebhookSignature(t *testing.T) {
	// Using a test secret in the whsec_ format
	secret := "whsec_MfKQ9r8GKYqrTwjUPD8ILPZIo2LaLaSw"

	// Test payload
	payload := `{"test": 2432232314}`
	id := "msg_p5jXN8AQM9LWM0D4loKWxJek"

	// Generate valid signature with current timestamp
	timestamp, signature := generateTestSignature(t, secret, id, payload)

	tests := []struct {
		name          string
		secret        string
		id            string
		timestamp     string
		payload       string
		signature     string
		expectedValid bool
	}{
		{
			name:          "valid signature",
			secret:        secret,
			id:            id,
			timestamp:     timestamp,
			payload:       payload,
			signature:     signature,
			expectedValid: true,
		},
		{
			name:          "invalid signature",
			secret:        secret,
			id:            id,
			timestamp:     timestamp,
			payload:       payload,
			signature:     "v1,invalid_signature",
			expectedValid: false,
		},
		{
			name:          "tampered payload",
			secret:        secret,
			id:            id,
			timestamp:     timestamp,
			payload:       `{"test": 9999999999}`,
			signature:     signature,
			expectedValid: false,
		},
		{
			name:          "tampered timestamp",
			secret:        secret,
			id:            id,
			timestamp:     "1704067999",
			payload:       payload,
			signature:     signature,
			expectedValid: false,
		},
		{
			name:          "empty secret",
			secret:        "",
			id:            id,
			timestamp:     timestamp,
			payload:       payload,
			signature:     signature,
			expectedValid: false,
		},
		{
			name:          "empty id",
			secret:        secret,
			id:            "",
			timestamp:     timestamp,
			payload:       payload,
			signature:     signature,
			expectedValid: false,
		},
		{
			name:          "empty timestamp",
			secret:        secret,
			id:            id,
			timestamp:     "",
			payload:       payload,
			signature:     signature,
			expectedValid: false,
		},
		{
			name:          "empty signature",
			secret:        secret,
			id:            id,
			timestamp:     timestamp,
			payload:       payload,
			signature:     "",
			expectedValid: false,
		},
		{
			name:          "invalid secret format",
			secret:        "invalid_secret",
			id:            id,
			timestamp:     timestamp,
			payload:       payload,
			signature:     signature,
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := VerifyWebhookSignature(tt.secret, tt.id, tt.timestamp, tt.payload, tt.signature)
			require.Equal(t, tt.expectedValid, valid)
		})
	}
}

func TestVerifyWebhookSignature_RealWorldExample(t *testing.T) {
	secret := "whsec_MfKQ9r8GKYqrTwjUPD8ILPZIo2LaLaSw"
	id := "msg_2ZNTsTTSoODpWdYiMvNJn6yJOmL"
	payload := `{
  "webhook_event": "blockchainWallet.new",
  "id": "bw_000000000000",
  "name": "Wallet Display Name",
  "network": "polygon",
  "address": "0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C",
  "signature_tx_hash": "0x3c499c542cef5e3811e1192ce70d8cc03d5c3359",
  "is_account_abstraction": false,
  "receiver_id": "re_000000000000"
}`

	// Generate signature with current timestamp
	timestamp, signature := generateTestSignature(t, secret, id, payload)

	valid := VerifyWebhookSignature(secret, id, timestamp, payload, signature)
	require.True(t, valid, "Real-world webhook signature should be valid")

	// Test with modified payload (should fail)
	tamperedPayload := `{
  "webhook_event": "blockchainWallet.new",
  "id": "bw_tampered_id",
  "name": "Wallet Display Name",
  "network": "polygon",
  "address": "0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C",
  "signature_tx_hash": "0x3c499c542cef5e3811e1192ce70d8cc03d5c3359",
  "is_account_abstraction": false,
  "receiver_id": "re_000000000000"
}`

	valid = VerifyWebhookSignature(secret, id, timestamp, tamperedPayload, signature)
	require.False(t, valid, "Tampered webhook should be invalid")
}

func TestVerifyWebhookSignature_TimingAttackResistance(t *testing.T) {
	secret := "whsec_MfKQ9r8GKYqrTwjUPD8ILPZIo2LaLaSw"
	id := "msg_timing"
	payload := `{"test":"data"}`

	// Generate correct signature with current timestamp
	timestamp, correctSignature := generateTestSignature(t, secret, id, payload)

	// Test correct signature first
	valid := VerifyWebhookSignature(secret, id, timestamp, payload, correctSignature)
	require.True(t, valid, "Correct signature should be accepted")

	// Test almost-correct signatures (modify the signature slightly)
	almostCorrectOne := correctSignature[:len(correctSignature)-2] + "X="
	almostCorrectTwo := "v1,X" + correctSignature[4:]
	almostCorrectThree := correctSignature[:20] + "XXXXX" + correctSignature[25:]

	for _, sig := range []string{almostCorrectOne, almostCorrectTwo, almostCorrectThree} {
		valid := VerifyWebhookSignature(secret, id, timestamp, payload, sig)
		require.False(t, valid, "Almost-correct signature should be rejected")
	}
}

func TestVerifyWebhookSignature_WithHTTPRequest(t *testing.T) {
	// This test simulates how the function would be used with an actual HTTP request
	secret := "whsec_MfKQ9r8GKYqrTwjUPD8ILPZIo2LaLaSw"
	id := "msg_http_test"
	payload := `{"webhook_event":"payment.completed","amount":1000}`

	// Generate signature with current timestamp
	timestamp, signature := generateTestSignature(t, secret, id, payload)

	// Simulate extracting headers from HTTP request
	req, _ := http.NewRequest("POST", "/webhook", nil)
	req.Header.Set("webhook-id", id)
	req.Header.Set("webhook-timestamp", timestamp)
	req.Header.Set("webhook-signature", signature)

	// Verify using extracted header values
	valid := VerifyWebhookSignature(
		secret,
		req.Header.Get("webhook-id"),
		req.Header.Get("webhook-timestamp"),
		payload,
		req.Header.Get("webhook-signature"),
	)
	require.True(t, valid, "Webhook signature from HTTP request should be valid")
}
