package blindpay_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/blindpaylabs/blindpay-go"
	"github.com/stretchr/testify/require"
)

// TestWebhookVectors_PublicAPI runs all 14 shared conformance vectors through
// the exported API only (package blindpay_test): VerifyWebhook and
// ConstructWebhookEvent with WithWebhookClock set to the vector's now.
func TestWebhookVectors_PublicAPI(t *testing.T) {
	data, err := os.ReadFile("testdata/webhook-vectors.json")
	require.NoError(t, err)
	var file struct {
		ToleranceSeconds int64 `json:"tolerance_seconds"`
		Vectors          []struct {
			Name    string            `json:"name"`
			Secret  string            `json:"secret"`
			Headers map[string]string `json:"headers"`
			Payload string            `json:"payload"`
			Now     int64             `json:"now"`
			Expect  string            `json:"expect"`
			Event   string            `json:"event"`
		} `json:"vectors"`
	}
	require.NoError(t, json.Unmarshal(data, &file))
	require.Len(t, file.Vectors, 14)
	require.Equal(t, int64(blindpay.DefaultWebhookTolerance/time.Second), file.ToleranceSeconds)

	passed := 0
	for _, v := range file.Vectors {
		v := v
		t.Run(v.Name, func(t *testing.T) {
			header := http.Header{}
			for name, value := range v.Headers {
				header.Set(name, value)
			}
			at := time.Unix(v.Now, 0)
			clock := blindpay.WithWebhookClock(func() time.Time { return at })

			valid := v.Expect == "valid"
			require.Equal(t, valid, blindpay.VerifyWebhook([]byte(v.Payload), header, v.Secret, clock))

			event, err := blindpay.ConstructWebhookEvent([]byte(v.Payload), header, v.Secret, clock)
			if !valid {
				require.True(t, errors.Is(err, blindpay.ErrInvalidWebhookSignature), "got %v", err)
				return
			}
			require.NoError(t, err)
			if v.Event != "" {
				require.Equal(t, blindpay.WebhookEvent(v.Event), event.EventType())
			}
			passed++
		})
	}
	t.Logf("webhook vectors through the public API: %d valid accepted, %d total", passed, len(file.Vectors))
}

// TestWithWebhookClock_BothDirections: the clock bounds the timestamp in the
// past and in the future, and a nil clock keeps time.Now.
func TestWithWebhookClock_BothDirections(t *testing.T) {
	data, err := os.ReadFile("testdata/webhook-vectors.json")
	require.NoError(t, err)
	var file struct {
		Vectors []struct {
			Name    string            `json:"name"`
			Secret  string            `json:"secret"`
			Headers map[string]string `json:"headers"`
			Payload string            `json:"payload"`
			Now     int64             `json:"now"`
		} `json:"vectors"`
	}
	require.NoError(t, json.Unmarshal(data, &file))
	v := file.Vectors[0]
	require.Equal(t, "valid", v.Name)
	header := http.Header{}
	for name, value := range v.Headers {
		header.Set(name, value)
	}
	at := func(d time.Duration) blindpay.WebhookOption {
		return blindpay.WithWebhookClock(func() time.Time { return time.Unix(v.Now, 0).Add(d) })
	}
	verify := func(opts ...blindpay.WebhookOption) bool {
		return blindpay.VerifyWebhook([]byte(v.Payload), header, v.Secret, opts...)
	}
	require.True(t, verify(at(0)))
	require.True(t, verify(at(299*time.Second)))
	require.False(t, verify(at(301*time.Second)), "timestamp too old")
	require.False(t, verify(at(-301*time.Second)), "timestamp too new")
	require.False(t, verify(blindpay.WithWebhookClock(nil)), "nil keeps time.Now: the vector is stale")
}
