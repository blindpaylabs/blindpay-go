package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestReconcileTypes_FieldTypeCompatibility covers the spec-vs-actual-Go-type
// check: does a property the SDK already models still correspond to a
// compatible Go type, given the CURRENT spec (a standing state assertion,
// not only a check for drift against the old snapshot). Every case here
// uses the SAME schema for old and new, so a failure can only come from
// the state check, never the separate old-vs-new type-change detector.
func TestReconcileTypes_FieldTypeCompatibility(t *testing.T) {
	tests := []struct {
		name           string
		goField        string // Go struct field declaration (type only, name fixed as "Value")
		specProperty   string // raw JSON for the "value" property definition
		wantNeedsHuman bool
	}{
		{
			name:           "string spec vs int Go field is incompatible (the exact swift-found bug shape)",
			goField:        `int`,
			specProperty:   `{"type":"string"}`,
			wantNeedsHuman: true,
		},
		{
			name:           "nullable to non-nullable of the same base type is compatible",
			goField:        `string`,
			specProperty:   `{"type":"string"}`,
			wantNeedsHuman: false,
		},
		{
			name:           "an enum-typed property degrading to a bare string is compatible (base type unchanged)",
			goField:        `CustomColor`,
			specProperty:   `{"type":"string"}`,
			wantNeedsHuman: false,
		},
		{
			// Deliberately compatible: JSON Schema "integer" and "number" are
			// both a JSON number token on the wire; encoding/json unmarshals
			// either into a float64 field without error, so this is not
			// flagged even though the spec's declared type text differs.
			name:           "spec integer against a Go float64 field is deliberately treated as compatible",
			goField:        `float64`,
			specProperty:   `{"type":"integer"}`,
			wantNeedsHuman: false,
		},
		{
			name:           "spec boolean against a Go string field is incompatible (CustomerOut.email-shaped bug)",
			goField:        `string`,
			specProperty:   `{"type":"boolean"}`,
			wantNeedsHuman: true,
		},
		{
			name:           "spec number against a Go *string field is incompatible (the real PayinOut.billing_fee_amount bug)",
			goField:        `*string`,
			specProperty:   `{"type":["number","null"]}`,
			wantNeedsHuman: true,
		},
		{
			name:           "spec array against a Go slice is compatible",
			goField:        `[]string`,
			specProperty:   `{"type":"array","items":{"type":"string"}}`,
			wantNeedsHuman: false,
		},
		{
			name:           "spec array against a Go bare string is incompatible",
			goField:        `string`,
			specProperty:   `{"type":"array","items":{"type":"string"}}`,
			wantNeedsHuman: true,
		},
		{
			name:           "a Go interface{} field (free-form, e.g. VirtualAccount.BlockchainWallet) accepts any spec type",
			goField:        `interface{}`,
			specProperty:   `{"type":"integer"}`,
			wantNeedsHuman: false,
		},
		{
			name:           "an unconstrained spec property (no type key) is always compatible",
			goField:        `int`,
			specProperty:   `{"example":42}`,
			wantNeedsHuman: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := writeFixture(t, map[string]string{
				widgetGoFile: "package widgets\n\ntype Widget struct {\n\tValue " + tt.goField + " `json:\"value\"`\n}\n",
			})
			sm := &SpecMap{Types: []TypeMapping{{
				Spec: "Widget",
				SDK:  []SDKSite{{File: widgetGoFile, Symbol: "Widget"}},
			}}}
			um := &Unmodeled{}
			spec := mustSpec(t, map[string]string{
				"Widget": `{"type":"object","properties":{"value":` + tt.specProperty + `},"required":[]}`,
			})

			plan := reconcileTypes(root, sm, um, spec, spec)
			if tt.wantNeedsHuman {
				require.True(t, containsString(plan.NeedsHuman, "value"), "expected a needs-human issue mentioning 'value', got: %v", plan.NeedsHuman)
				require.True(t, containsString(plan.NeedsHuman, "encoding/json would fail"))
			} else {
				require.Empty(t, plan.NeedsHuman, "expected no needs-human issue")
			}
			require.Empty(t, plan.Actions, "an already-modeled field is never itself an applicable action")
		})
	}
}

func TestReconcileTypes_FieldTypeMismatchExcusedByUnmodeledJSON(t *testing.T) {
	root := writeFixture(t, map[string]string{
		widgetGoFile: `package widgets

type Widget struct {
	Value int ` + "`json:\"value\"`" + `
}
`,
	})
	sm := &SpecMap{Types: []TypeMapping{{
		Spec: "Widget",
		SDK:  []SDKSite{{File: widgetGoFile, Symbol: "Widget"}},
	}}}
	um := &Unmodeled{Properties: []PropertyExclusion{
		{Schema: "Widget", Property: "value", Reason: "known pre-existing bug, tracked for its own fix", Owner: "eric"},
	}}
	spec := mustSpec(t, map[string]string{
		"Widget": `{"type":"object","properties":{"value":{"type":"string"}},"required":[]}`,
	})

	plan := reconcileTypes(root, sm, um, spec, spec)
	require.Empty(t, plan.NeedsHuman, "unmodeled.json excuses a known type mismatch the same way it excuses a missing property")
}
