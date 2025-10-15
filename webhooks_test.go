package blindpay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVerifyWebhookSignature(t *testing.T) {
	// Create a test secret
	secretKey := []byte("test_secret_key_12345")
	secretB64 := base64.StdEncoding.EncodeToString(secretKey)
	secret := "whsec_" + secretB64

	id := "msg_test123"
	timestamp := "1704067200"
	payload := `{"webhook_event":"blockchainWallet.new","id":"bw_000000000000","name":"Wallet Display Name","network":"polygon","address":"0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C","signature_tx_hash":"0x3c499c542cef5e3811e1192ce70d8cc03d5c3359","is_account_abstraction":false,"receiver_id":"re_000000000000"}`

	// Generate valid signature
	signedContent := id + "." + timestamp + "." + payload
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(signedContent))
	expectedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	tests := []struct {
		name          string
		secret        string
		id            string
		timestamp     string
		payload       string
		signature     string
		expectedValid bool
		expectedError bool
	}{
		{
			name:          "valid signature",
			secret:        secret,
			id:            id,
			timestamp:     timestamp,
			payload:       payload,
			signature:     expectedSignature,
			expectedValid: true,
			expectedError: false,
		},
		{
			name:          "invalid signature",
			secret:        secret,
			id:            id,
			timestamp:     timestamp,
			payload:       payload,
			signature:     "invalid_signature",
			expectedValid: false,
			expectedError: false,
		},
		{
			name:          "tampered payload",
			secret:        secret,
			id:            id,
			timestamp:     timestamp,
			payload:       `{"webhook_event":"blockchainWallet.new","id":"bw_999999999999","name":"Wallet Display Name","network":"polygon","address":"0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C","signature_tx_hash":"0x3c499c542cef5e3811e1192ce70d8cc03d5c3359","is_account_abstraction":false,"receiver_id":"re_000000000000"}`,
			signature:     expectedSignature,
			expectedValid: false,
			expectedError: false,
		},
		{
			name:          "tampered timestamp",
			secret:        secret,
			id:            id,
			timestamp:     "1704067999",
			payload:       payload,
			signature:     expectedSignature,
			expectedValid: false,
			expectedError: false,
		},
		{
			name:          "empty secret",
			secret:        "",
			id:            id,
			timestamp:     timestamp,
			payload:       payload,
			signature:     expectedSignature,
			expectedValid: false,
			expectedError: true,
		},
		{
			name:          "empty id",
			secret:        secret,
			id:            "",
			timestamp:     timestamp,
			payload:       payload,
			signature:     expectedSignature,
			expectedValid: false,
			expectedError: true,
		},
		{
			name:          "empty timestamp",
			secret:        secret,
			id:            id,
			timestamp:     "",
			payload:       payload,
			signature:     expectedSignature,
			expectedValid: false,
			expectedError: true,
		},
		{
			name:          "empty payload",
			secret:        secret,
			id:            id,
			timestamp:     timestamp,
			payload:       "",
			signature:     expectedSignature,
			expectedValid: false,
			expectedError: true,
		},
		{
			name:          "empty signature",
			secret:        secret,
			id:            id,
			timestamp:     timestamp,
			payload:       payload,
			signature:     "",
			expectedValid: false,
			expectedError: true,
		},
		{
			name:          "invalid secret format - no prefix",
			secret:        secretB64,
			id:            id,
			timestamp:     timestamp,
			payload:       payload,
			signature:     expectedSignature,
			expectedValid: false,
			expectedError: true,
		},
		{
			name:          "invalid secret format - wrong prefix",
			secret:        "wrong_" + secretB64,
			id:            id,
			timestamp:     timestamp,
			payload:       payload,
			signature:     expectedSignature,
			expectedValid: false,
			expectedError: true,
		},
		{
			name:          "invalid secret format - invalid base64",
			secret:        "whsec_invalid!!!base64",
			id:            id,
			timestamp:     timestamp,
			payload:       payload,
			signature:     expectedSignature,
			expectedValid: false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := VerifyWebhookSignature(tt.secret, tt.id, tt.timestamp, tt.payload, tt.signature)

			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectedValid, valid)
		})
	}
}

func TestVerifyWebhookSignature_RealWorldExample(t *testing.T) {
	secretKey := []byte("my_webhook_secret_key_very_secure_123456")
	secretB64 := base64.StdEncoding.EncodeToString(secretKey)
	secret := "whsec_" + secretB64
	id := "msg_2ZNTsTTSoODpWdYiMvNJn6yJOmL"
	timestamp := "1704153600"
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

	signedContent := id + "." + timestamp + "." + payload
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(signedContent))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	valid, err := VerifyWebhookSignature(secret, id, timestamp, payload, signature)
	require.NoError(t, err)
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

	valid, err = VerifyWebhookSignature(secret, id, timestamp, tamperedPayload, signature)
	require.NoError(t, err)
	require.False(t, valid, "Tampered webhook should be invalid")
}

func TestVerifyWebhookSignature_TimingAttackResistance(t *testing.T) {
	secretKey := []byte("timing_test_secret")
	secretB64 := base64.StdEncoding.EncodeToString(secretKey)
	secret := "whsec_" + secretB64

	id := "msg_timing"
	timestamp := "1704153600"
	payload := `{"test":"data"}`

	signedContent := id + "." + timestamp + "." + payload
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(signedContent))
	correctSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	almostCorrectOne := correctSignature[:len(correctSignature)-1] + "X"
	almostCorrectTwo := "X" + correctSignature[1:]
	almostCorrectThree := correctSignature[:len(correctSignature)/2] + "XXXXX"

	for _, sig := range []string{almostCorrectOne, almostCorrectTwo, almostCorrectThree} {
		valid, err := VerifyWebhookSignature(secret, id, timestamp, payload, sig)
		require.NoError(t, err)
		require.False(t, valid, "Almost-correct signature should be rejected")
	}
}
