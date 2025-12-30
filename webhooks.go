package blindpay

import (
	svix "github.com/svix/svix-webhooks/go"
)

// VerifyWebhookSignature verifies the BlindPay webhook signature using Svix.
//
// Parameters:
//   - secret: The webhook secret from BlindPay dashboard
//   - id: The value of the `webhook-id` header
//   - timestamp: The value of the `webhook-timestamp` header
//   - payload: The raw request body as a string
//   - signature: The value of the `webhook-signature` header
//
// Returns true if the signature is valid, false otherwise.
func VerifyWebhookSignature(secret, id, timestamp, payload, signature string) bool {
	wh, err := svix.NewWebhook(secret)
	if err != nil {
		return false
	}

	// Keys must use canonical HTTP header casing
	headers := map[string][]string{
		"Webhook-Id":        {id},
		"Webhook-Timestamp": {timestamp},
		"Webhook-Signature": {signature},
	}

	err = wh.Verify([]byte(payload), headers)
	return err == nil
}
