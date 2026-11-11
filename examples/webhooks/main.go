// Command webhooks registers a webhook endpoint and serves a handler that
// verifies every delivery before trusting it (docs: "Webhooks").
package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/blindpaylabs/blindpay-go"
	"github.com/blindpaylabs/blindpay-go/webhookendpoints"
)

func main() {
	client, err := blindpay.New(os.Getenv("BLINDPAY_API_KEY"), os.Getenv("BLINDPAY_INSTANCE_ID"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	endpoint, err := client.WebhookEndpoints.Create(ctx, &webhookendpoints.CreateParams{
		URL:    "https://example.com/webhooks/blindpay",
		Events: []blindpay.WebhookEvent{blindpay.WebhookEventPayoutUpdate, blindpay.WebhookEventPayoutComplete},
	})
	if err != nil {
		log.Fatal(err)
	}
	secret, err := client.WebhookEndpoints.GetSecret(ctx, endpoint.ID)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/webhooks/blindpay", handler(secret.Key))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handler verifies the raw body, ignores repeated deliveries (same svix-id)
// and branches on the typed event. Never log the secret.
func handler(secret string) http.Handler {
	var mu sync.Mutex
	seen := map[string]bool{} // use a durable store in production

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body) // raw bytes: never re-encode the JSON
		if err != nil {
			http.Error(w, "read error", http.StatusBadRequest)
			return
		}
		event, err := blindpay.ConstructWebhookEvent(body, r.Header, secret)
		if err != nil {
			http.Error(w, "invalid signature", http.StatusBadRequest)
			return
		}

		mu.Lock()
		dup := seen[r.Header.Get("svix-id")]
		seen[r.Header.Get("svix-id")] = true
		mu.Unlock()
		if dup {
			w.WriteHeader(http.StatusOK)
			return
		}

		switch e := event.(type) {
		case *blindpay.PayoutWebhook:
			// A failed payout arrives on payout.update or payout.complete.
			log.Printf("payout %s: %s (%s)", e.ID, e.Status, e.EventType())
		case *blindpay.UnknownWebhook:
			log.Printf("unhandled event %s", e.EventType())
		}
		w.WriteHeader(http.StatusOK)
	})
}
