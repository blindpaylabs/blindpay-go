package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckWebhookCoverage(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"go.mod":            "module example.com/fixture\n",
		"webhook_events.go": "package fixture\n\ntype WidgetWebhook struct {\n\tHeader\n\tWidget\n\tExtra *string `json:\"extra,omitempty\"`\n}\n\ntype Header struct{}\n",
		"widget.go":         "package fixture\n\ntype Widget struct {\n\tID string `json:\"id\"`\n}\n",
	})
	sm := &SpecMap{Webhooks: []WebhookMapping{{Spec: "WidgetWebhookOut", SDK: []SDKSite{
		{File: "webhook_events.go", Symbol: "WidgetWebhook"},
		{File: "widget.go", Symbol: "Widget"},
	}}}}
	spec := func(props ...string) *specDoc {
		p := map[string]any{"webhook_event": map[string]any{"type": "string"}}
		for _, name := range props {
			p[name] = map[string]any{"type": "string"}
		}
		return loadSpecDoc(map[string]any{"components": map[string]any{"schemas": map[string]any{
			"WidgetWebhookOut": map[string]any{"type": "object", "properties": p},
		}}})
	}
	um := &Unmodeled{}

	require.Empty(t, checkWebhookCoverage(root, sm, um, spec("id", "extra"), nil), "fields on the wrapper and the embedded struct both count")

	issues := checkWebhookCoverage(root, sm, um, spec("id", "review_reason"), nil)
	require.Len(t, issues, 1, "webhook-only drift is reported")
	require.Contains(t, issues[0], "WidgetWebhookOut.review_reason")

	pending := []action{{addsField: &plannedField{File: "widget.go", Symbol: "Widget", JSONName: "review_reason"}}}
	require.Empty(t, checkWebhookCoverage(root, sm, um, spec("id", "review_reason"), pending), "a field a pending action adds counts as modeled")

	um.Properties = []PropertyExclusion{{Schema: "WidgetWebhookOut", Property: "review_reason", Reason: "r", Owner: "o"}}
	require.Empty(t, checkWebhookCoverage(root, sm, um, spec("id", "review_reason"), nil), "an unmodeled.json exclusion excuses it")

	require.Len(t, checkWebhookCoverage(root, sm, &Unmodeled{}, loadSpecDoc(map[string]any{}), nil), 1, "a missing schema is reported")
}
