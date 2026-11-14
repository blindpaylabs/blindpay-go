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

// publicVector is one case of testdata/webhook-vectors.json. This file is
// an external test package (only the exported API), so it cannot share
// webhookVector from webhooks_vectors_test.go.
type publicVector struct {
	Name    string            `json:"name"`
	Secret  string            `json:"secret"`
	Headers map[string]string `json:"headers"`
	Payload string            `json:"payload"`
	Now     int64             `json:"now"`
	Expect  string            `json:"expect"`
	Event   string            `json:"event"`
}

func loadPublicVectors(t *testing.T) (vectors []publicVector, toleranceSeconds int64) {
	t.Helper()
	data, err := os.ReadFile("testdata/webhook-vectors.json")
	require.NoError(t, err)
	var file struct {
		ToleranceSeconds int64          `json:"tolerance_seconds"`
		Vectors          []publicVector `json:"vectors"`
	}
	require.NoError(t, json.Unmarshal(data, &file))
	require.Len(t, file.Vectors, 14)
	return file.Vectors, file.ToleranceSeconds
}

func (v publicVector) header() http.Header {
	h := http.Header{}
	for name, value := range v.Headers {
		h.Set(name, value)
	}
	return h
}

// clock returns WithWebhookClock fixed at the vector's now shifted by d.
func (v publicVector) clock(d time.Duration) blindpay.WebhookOption {
	at := time.Unix(v.Now, 0).Add(d)
	return blindpay.WithWebhookClock(func() time.Time { return at })
}

// TestWebhookVectors_PublicAPI runs all 14 shared conformance vectors through
// the exported API only: VerifyWebhook and ConstructWebhookEvent with
// WithWebhookClock set to the vector's now.
func TestWebhookVectors_PublicAPI(t *testing.T) {
	vectors, tolerance := loadPublicVectors(t)
	require.Equal(t, int64(blindpay.DefaultWebhookTolerance/time.Second), tolerance)

	for _, v := range vectors {
		v := v
		t.Run(v.Name, func(t *testing.T) {
			valid := v.Expect == "valid"
			require.Equal(t, valid, blindpay.VerifyWebhook([]byte(v.Payload), v.header(), v.Secret, v.clock(0)))

			event, err := blindpay.ConstructWebhookEvent([]byte(v.Payload), v.header(), v.Secret, v.clock(0))
			if !valid {
				require.True(t, errors.Is(err, blindpay.ErrInvalidWebhookSignature), "got %v", err)
				return
			}
			require.NoError(t, err)
			if v.Event != "" {
				require.Equal(t, blindpay.WebhookEvent(v.Event), event.EventType())
			}
		})
	}
}

// TestWithWebhookClock_BothDirections: the clock bounds the timestamp in the
// past and in the future, and a nil clock keeps time.Now.
func TestWithWebhookClock_BothDirections(t *testing.T) {
	vectors, _ := loadPublicVectors(t)
	v := vectors[0]
	require.Equal(t, "valid", v.Name)
	verify := func(opt blindpay.WebhookOption) bool {
		return blindpay.VerifyWebhook([]byte(v.Payload), v.header(), v.Secret, opt)
	}
	require.True(t, verify(v.clock(0)))
	require.True(t, verify(v.clock(299*time.Second)))
	require.False(t, verify(v.clock(301*time.Second)), "timestamp too old")
	require.False(t, verify(v.clock(-301*time.Second)), "timestamp too new")
	require.False(t, verify(blindpay.WithWebhookClock(nil)), "nil keeps time.Now: the vector is stale")
}
