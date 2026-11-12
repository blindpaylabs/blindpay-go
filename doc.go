// Package blindpay is the Go SDK for the BlindPay API, which moves money
// between stablecoins and local bank rails (payins, payouts, transfers,
// virtual accounts, managed wallets).
//
// # Quick start
//
// Create one Client per API key and instance; it is immutable and safe for
// concurrent use. Each resource is a field of the client (Payouts, Payins,
// Quotes, Customers, ...), and every method takes a context first and
// optional per-request options last.
//
//	client, err := blindpay.New(apiKey, instanceID)
//	if err != nil {
//		log.Fatal(err)
//	}
//	rails, err := client.Available.GetRails(ctx)
//
// Amounts on payouts, payins, quotes and limits are integer cents
// (1000 = 10.00) even where the Go field is a float64; partner percentage
// fees are basis points (100 = 1%) and wallet balances are whole tokens.
//
// # Errors
//
// A non-2xx response is an *APIError. Branch on its Code, the stable error
// code, with the generated ErrorCode constants; RetryClass gives the API's
// retry guidance and LookupErrorCode the table's metadata. Message stays
// exposed for the bodies that have no code. A transport failure is a
// *NetworkError, a 2xx that cannot be decoded a *DecodeError.
//
//	var apiErr *blindpay.APIError
//	if errors.As(err, &apiErr) {
//		switch blindpay.ErrorCode(apiErr.Code) {
//		case blindpay.ErrorCodeTransfersSwapAmountBelowMinimum:
//			// swaps between different tokens need at least 100 USDC
//		}
//	}
//
// # Retries and idempotency
//
// Requests have a 60-second per-attempt timeout (WithTimeout) and are
// retried up to twice (WithMaxRetries) with exponential backoff, honouring
// Retry-After up to 60 seconds. GET is retried on network errors, 408, 429,
// 500, 502, 503 and 504. POST, PUT, PATCH and DELETE carry an Idempotency-Key
// (generated, or yours via WithIdempotencyKey) and are retried only on
// network errors, 429 and an in-flight 409, never on a 5xx: look the
// resource up, then resend with APIError.IdempotencyKey. Indicative FX quote
// endpoints are rate limited per instance; a limited call returns 429
// AUTH_RATE_LIMITED with a Retry-After header, which the SDK honours when it
// retries.
//
// # Webhooks
//
// Verify a delivery on the raw request body with ConstructWebhookEvent,
// which checks the Svix signature and timestamp and returns a typed payload
// (an *UnknownWebhook for an event newer than this SDK). Deduplicate on the
// svix-id header.
//
//	body, _ := io.ReadAll(r.Body)
//	event, err := blindpay.ConstructWebhookEvent(body, r.Header, secret)
//	if err != nil {
//		http.Error(w, "invalid signature", http.StatusBadRequest)
//		return
//	}
//	switch e := event.(type) {
//	case *blindpay.PayoutWebhook:
//		log.Println(e.ID, e.Status)
//	}
//
// The API key never appears when a client is printed, and transport error
// messages never include query strings.
package blindpay
