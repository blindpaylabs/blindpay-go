package blindpay

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestParseWebhookEvent_Fixtures decodes every variant of every shared
// webhook fixture (testdata/fixtures/webhooks): each known event must decode
// into its typed payload with the right EventType, and an unknown event into
// *UnknownWebhook carrying the raw body.
func TestParseWebhookEvent_Fixtures(t *testing.T) {
	files, err := filepath.Glob("testdata/fixtures/webhooks/*.json")
	require.NoError(t, err)

	seen := map[WebhookEvent]bool{}
	for _, file := range files {
		file := file
		b, err := os.ReadFile(file)
		require.NoError(t, err)
		var fixture struct {
			Event    string                     `json:"event"`
			Variants map[string]json.RawMessage `json:"variants"`
		}
		require.NoError(t, json.Unmarshal(b, &fixture))
		require.NotEmpty(t, fixture.Variants, file)

		names := make([]string, 0, len(fixture.Variants))
		for name := range fixture.Variants {
			names = append(names, name)
		}
		sort.Strings(names)

		event := WebhookEvent(fixture.Event)
		_, known := webhookTypes[event]
		if known {
			seen[event] = true
		}
		for _, name := range names {
			body := fixture.Variants[name]
			t.Run(strings.TrimSuffix(filepath.Base(file), ".json")+"/"+name, func(t *testing.T) {
				payload, err := ParseWebhookEvent(body)
				require.NoError(t, err)
				require.Equal(t, event, payload.EventType())

				if !known {
					unknown, ok := payload.(*UnknownWebhook)
					require.True(t, ok, "got %T", payload)
					require.JSONEq(t, string(body), string(unknown.Raw))
					return
				}
				require.IsType(t, webhookTypes[event](), payload)
			})
		}
	}

	// Every event in the spec has a fixture and a typed payload.
	require.Len(t, webhookTypes, 28)
	for event := range webhookTypes {
		require.True(t, seen[event], "no fixture for %s", event)
	}
}

func TestParseWebhookEvent_TypedFields(t *testing.T) {
	payload, err := ParseWebhookEvent([]byte(`{"webhook_event":"payout.update","id":"po_000000000000","status":"failed","offramp_wallet_id":"ow_000000000000"}`))
	require.NoError(t, err)
	payout, ok := payload.(*PayoutWebhook)
	require.True(t, ok)
	require.Equal(t, "po_000000000000", payout.ID)
	require.Equal(t, TransactionStatusFailed, payout.Status)
	require.Equal(t, "ow_000000000000", *payout.OfframpWalletID)

	payload, err = ParseWebhookEvent([]byte(`{"webhook_event":"wallet.inbound","id":"wa_000000000000","address":"0x0","network":"base","token":{"id":"tk_000000000000","address":"0x1","symbol":"USDC","amount":10.5},"deposit_id":"de_000000000000","transaction_hash":"0x2"}`))
	require.NoError(t, err)
	inbound, ok := payload.(*WalletInboundWebhook)
	require.True(t, ok)
	require.Equal(t, StablecoinToken("USDC"), inbound.Token.Symbol)
	require.Equal(t, 10.5, inbound.Token.Amount)
}

func TestParseWebhookEvent_InvalidJSON(t *testing.T) {
	_, err := ParseWebhookEvent([]byte(`not json`))
	require.Error(t, err)
}
