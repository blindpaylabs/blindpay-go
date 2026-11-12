package blindpay

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// webhookVector is one case of testdata/webhook-vectors.json, the cross-SDK
// conformance set generated and self-checked with the official svix package.
type webhookVector struct {
	Name    string            `json:"name"`
	Secret  string            `json:"secret"`
	Headers map[string]string `json:"headers"`
	Payload string            `json:"payload"`
	Now     int64             `json:"now"`
	Expect  string            `json:"expect"`
	Event   string            `json:"event"`
}

func loadWebhookVectors(t *testing.T) []webhookVector {
	t.Helper()
	b, err := os.ReadFile("testdata/webhook-vectors.json")
	require.NoError(t, err)
	var file struct {
		Vectors []webhookVector `json:"vectors"`
	}
	require.NoError(t, json.Unmarshal(b, &file))
	require.Len(t, file.Vectors, 14)
	return file.Vectors
}

func (v webhookVector) header() http.Header {
	h := http.Header{}
	for name, value := range v.Headers {
		h.Set(name, value)
	}
	return h
}

func (v webhookVector) clock() WebhookOption {
	now := time.Unix(v.Now, 0)
	return WithWebhookClock(func() time.Time { return now })
}

func TestVerifyWebhook_Vectors(t *testing.T) {
	for _, v := range loadWebhookVectors(t) {
		v := v
		t.Run(v.Name, func(t *testing.T) {
			require.Contains(t, []string{"valid", "invalid"}, v.Expect)
			got := VerifyWebhook([]byte(v.Payload), v.header(), v.Secret, v.clock())
			require.Equal(t, v.Expect == "valid", got)
		})
	}
}

func TestConstructWebhookEvent_Vectors(t *testing.T) {
	for _, v := range loadWebhookVectors(t) {
		v := v
		t.Run(v.Name, func(t *testing.T) {
			event, err := ConstructWebhookEvent([]byte(v.Payload), v.header(), v.Secret, v.clock())
			if v.Expect == "invalid" {
				require.Nil(t, event)
				require.True(t, errors.Is(err, ErrInvalidWebhookSignature), "got %v", err)
				return
			}
			require.NoError(t, err)
			var head struct {
				Event WebhookEvent `json:"webhook_event"`
			}
			require.NoError(t, json.Unmarshal([]byte(v.Payload), &head))
			require.Equal(t, head.Event, event.EventType())
			if _, known := webhookTypes[head.Event]; !known {
				unknown, ok := event.(*UnknownWebhook)
				require.True(t, ok, "unknown event must decode as *UnknownWebhook, got %T", event)
				require.JSONEq(t, v.Payload, string(unknown.Raw))
			}
		})
	}
}

func TestVerifyWebhook_HeaderAliasesAndCase(t *testing.T) {
	v := loadWebhookVectors(t)[0]
	require.Equal(t, "valid", v.Expect)

	aliases := http.Header{}
	aliases.Set("Webhook-Id", v.Headers["svix-id"])
	aliases.Set("Webhook-Timestamp", v.Headers["svix-timestamp"])
	aliases.Set("Webhook-Signature", v.Headers["svix-signature"])
	require.True(t, VerifyWebhook([]byte(v.Payload), aliases, v.Secret, v.clock()))

	// A map built by hand (not via Header.Set) keeps non-canonical keys.
	raw := http.Header{
		"svix-id":        {v.Headers["svix-id"]},
		"SVIX-TIMESTAMP": {v.Headers["svix-timestamp"]},
		"svix-signature": {v.Headers["svix-signature"]},
	}
	require.True(t, VerifyWebhook([]byte(v.Payload), raw, v.Secret, v.clock()))
}

func TestVerifyWebhook_Tolerance(t *testing.T) {
	v := loadWebhookVectors(t)[0]
	late := time.Unix(v.Now+600, 0)
	clock := WithWebhookClock(func() time.Time { return late })

	require.False(t, VerifyWebhook([]byte(v.Payload), v.header(), v.Secret, clock))
	require.True(t, VerifyWebhook([]byte(v.Payload), v.header(), v.Secret, clock, WithWebhookTolerance(10*time.Minute)))
}

func TestVerifyWebhook_MalformedInputNeverPanics(t *testing.T) {
	v := loadWebhookVectors(t)[0]
	for name, secret := range map[string]string{
		"empty":          "",
		"prefix only":    "whsec_",
		"not base64":     "whsec_!!!",
		"no prefix ok":   v.Secret[len("whsec_"):],
		"nil header map": v.Secret,
	} {
		name, secret := name, secret
		t.Run(name, func(t *testing.T) {
			h := v.header()
			if name == "nil header map" {
				h = nil
			}
			got := VerifyWebhook([]byte(v.Payload), h, secret, v.clock())
			require.Equal(t, name == "no prefix ok", got)
		})
	}
}
