package blindpay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// DefaultWebhookTolerance is the maximum allowed distance between a
// webhook's timestamp and the current time.
const DefaultWebhookTolerance = 5 * time.Minute

// ErrInvalidWebhookSignature is wrapped by ConstructWebhookEvent when a
// webhook fails verification.
var ErrInvalidWebhookSignature = errors.New("blindpay: invalid webhook signature")

// WebhookOption configures webhook verification.
type WebhookOption func(*webhookConfig)

type webhookConfig struct {
	tolerance time.Duration
	now       func() time.Time
}

// WithWebhookTolerance overrides the allowed timestamp distance
// (default 5 minutes).
func WithWebhookTolerance(d time.Duration) WebhookOption {
	return func(c *webhookConfig) {
		c.tolerance = d
	}
}

// withWebhookClock is the test seam for the current time.
func withWebhookClock(now func() time.Time) WebhookOption {
	return func(c *webhookConfig) {
		c.now = now
	}
}

// VerifyWebhook reports whether payload is an authentic BlindPay webhook.
// It needs no client: pass the raw request body (never re-serialised JSON),
// the request headers and the endpoint's signing secret ("whsec_...").
//
// BlindPay signs webhooks with Svix: the svix-id, svix-timestamp and
// svix-signature headers (webhook-* aliases are accepted, names are
// case-insensitive). The signature is HMAC-SHA256 over
// "{id}.{timestamp}.{body}", compared in constant time against every "v1,"
// entry, and the timestamp must be within the tolerance. VerifyWebhook never
// panics; any malformed input returns false.
//
//	body, _ := io.ReadAll(r.Body)
//	if !blindpay.VerifyWebhook(body, r.Header, secret) {
//		http.Error(w, "invalid signature", http.StatusBadRequest)
//		return
//	}
//
// Deduplicate deliveries on the svix-id header, and never log the secret.
func VerifyWebhook(payload []byte, header http.Header, secret string, opts ...WebhookOption) bool {
	return checkWebhook(payload, header, secret, opts) == nil
}

// VerifyWebhookSignature verifies a webhook from its individual header
// values.
//
// Deprecated: use VerifyWebhook, which takes the request headers.
//
// Parameters:
//   - secret: the webhook secret from the BlindPay dashboard
//   - id: the svix-id (or webhook-id) header
//   - timestamp: the svix-timestamp (or webhook-timestamp) header
//   - payload: the raw request body
//   - signature: the svix-signature (or webhook-signature) header
func VerifyWebhookSignature(secret, id, timestamp, payload, signature string) bool {
	h := http.Header{}
	h.Set("Svix-Id", id)
	h.Set("Svix-Timestamp", timestamp)
	h.Set("Svix-Signature", signature)
	return VerifyWebhook([]byte(payload), h, secret)
}

func checkWebhook(payload []byte, header http.Header, secret string, opts []WebhookOption) error {
	cfg := webhookConfig{tolerance: DefaultWebhookTolerance, now: time.Now}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	id := headerValue(header, "svix-id", "webhook-id")
	ts := headerValue(header, "svix-timestamp", "webhook-timestamp")
	sigs := headerValue(header, "svix-signature", "webhook-signature")
	if id == "" || ts == "" || sigs == "" {
		return errors.New("missing signature headers")
	}

	for _, r := range ts {
		if r < '0' || r > '9' {
			return errors.New("invalid timestamp")
		}
	}
	unix, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return errors.New("invalid timestamp")
	}
	if diff := math.Abs(float64(cfg.now().Unix() - unix)); diff > cfg.tolerance.Seconds() {
		return errors.New("timestamp outside tolerance")
	}

	key, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(secret, "whsec_"))
	if err != nil || len(key) == 0 {
		return errors.New("invalid secret")
	}

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(id + "." + ts + "."))
	mac.Write(payload)
	expected := []byte(base64.StdEncoding.EncodeToString(mac.Sum(nil)))

	for _, part := range strings.Fields(sigs) {
		version, value, ok := strings.Cut(part, ",")
		if !ok || version != "v1" {
			continue
		}
		// hmac.Equal is constant time and safe on length mismatch.
		if hmac.Equal([]byte(value), expected) {
			return nil
		}
	}
	return errors.New("no matching signature")
}

// headerValue returns the first non-empty header among names, matching
// names case-insensitively (also for non-canonical map keys).
func headerValue(header http.Header, names ...string) string {
	for _, name := range names {
		if v := header.Get(name); v != "" {
			return v
		}
		for k, values := range header {
			if strings.EqualFold(k, name) && len(values) > 0 && values[0] != "" {
				return values[0]
			}
		}
	}
	return ""
}

// ConstructWebhookEvent verifies the webhook and decodes it into a typed
// event (see WebhookPayload). An unknown event name returns *UnknownWebhook
// with the raw body instead of failing.
//
//	event, err := blindpay.ConstructWebhookEvent(body, r.Header, secret)
//	if err != nil { ... }
//	switch e := event.(type) {
//	case *blindpay.PayoutWebhook:
//		if e.Status == blindpay.TransactionStatusFailed { ... } // failures arrive on payout.update and payout.complete
//	case *blindpay.UnknownWebhook:
//		// newer event; e.Raw holds the payload
//	}
func ConstructWebhookEvent(payload []byte, header http.Header, secret string, opts ...WebhookOption) (WebhookPayload, error) {
	if err := checkWebhook(payload, header, secret, opts); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidWebhookSignature, err)
	}
	return ParseWebhookEvent(payload)
}
