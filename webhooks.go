package blindpay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

// VerifyWebhookSignature verifies the BlindPay webhook signature.
//
// The signature header format is: "v1,sig1 v1,sig2 ..." where each signature
// is prefixed with the version (v1) and we need to match at least one.
//
// Parameters:
//   - secret: The webhook secret from BlindPay dashboard (format: "whsec_<base64>")
//   - id: The value of the `svix-id` header
//   - timestamp: The value of the `svix-timestamp` header
//   - payload: The raw request body
//   - svixSignature: The value of the `svix-signature` header
//
// Returns true if the signature is valid, false otherwise.
func VerifyWebhookSignature(secret, id, timestamp, payload, svixSignature string) (bool, error) {
	if secret == "" || id == "" || timestamp == "" || payload == "" || svixSignature == "" {
		return false, fmt.Errorf("all parameters must be provided")
	}

	parts := strings.Split(secret, "_")
	if len(parts) != 2 || parts[0] != "whsec" {
		return false, fmt.Errorf("invalid secret format, expected 'whsec_<base64>'")
	}

	secretBytes, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, fmt.Errorf("failed to decode secret: %w", err)
	}

	signedContent := id + "." + timestamp + "." + payload

	h := hmac.New(sha256.New, secretBytes)
	h.Write([]byte(signedContent))
	expectedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	signatures := strings.Fields(svixSignature)
	for _, sig := range signatures {
		sigParts := strings.SplitN(sig, ",", 2)
		if len(sigParts) == 2 && sigParts[0] == "v1" {
			if hmac.Equal([]byte(sigParts[1]), []byte(expectedSignature)) {
				return true, nil
			}
		}
	}

	if hmac.Equal([]byte(svixSignature), []byte(expectedSignature)) {
		return true, nil
	}

	return false, nil
}
