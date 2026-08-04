package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const widgetGoFile = "widgets/client.go"

func widgetFixture(t *testing.T) string {
	return writeFixture(t, map[string]string{
		widgetGoFile: `package widgets

type Widget struct {
	ID   string  ` + "`json:\"id\"`" + `
	Name *string ` + "`json:\"name,omitempty\"`" + `
}
`,
	})
}

func widgetTypeMapping() TypeMapping {
	return TypeMapping{
		Spec: "Widget",
		SDK:  []SDKSite{{File: widgetGoFile, Symbol: "Widget"}},
	}
}

func TestReconcileTypes_AppliesNewOptionalScalarField(t *testing.T) {
	root := widgetFixture(t)
	sm := &SpecMap{Types: []TypeMapping{widgetTypeMapping()}}
	um := &Unmodeled{}

	spec := mustSpec(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"id":{"type":"string"},
			"name":{"type":["string","null"]},
			"score":{"type":"integer"}
		},"required":["id"]}`,
	})

	plan := reconcileTypes(root, sm, um, spec, spec)
	require.Empty(t, plan.NeedsHuman)
	require.Len(t, plan.Actions, 1)
	require.Contains(t, plan.Actions[0].description, "score")
	require.Equal(t, "patch", plan.Actions[0].bump)
	require.Contains(t, plan.Actions[0].insert.Text, "Score")
	require.Contains(t, plan.Actions[0].insert.Text, `json:"score,omitempty"`)
}

func TestReconcileTypes_RequiredMissingPropertyIsNeedsHuman(t *testing.T) {
	root := widgetFixture(t)
	sm := &SpecMap{Types: []TypeMapping{widgetTypeMapping()}}
	um := &Unmodeled{}

	spec := mustSpec(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"id":{"type":"string"},
			"name":{"type":["string","null"]},
			"status":{"type":"string"}
		},"required":["id","status"]}`,
	})

	plan := reconcileTypes(root, sm, um, spec, spec)
	require.Empty(t, plan.Actions)
	require.True(t, containsString(plan.NeedsHuman, "status"))
	require.True(t, containsString(plan.NeedsHuman, "required"))
}

func TestReconcileTypes_NonScalarPropertyIsNeedsHuman(t *testing.T) {
	root := widgetFixture(t)
	sm := &SpecMap{Types: []TypeMapping{widgetTypeMapping()}}
	um := &Unmodeled{}

	spec := mustSpec(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"id":{"type":"string"},
			"name":{"type":["string","null"]},
			"meta":{"type":"object","properties":{"k":{"type":"string"}}}
		},"required":["id"]}`,
	})

	plan := reconcileTypes(root, sm, um, spec, spec)
	require.Empty(t, plan.Actions)
	require.True(t, containsString(plan.NeedsHuman, "meta"))
	require.True(t, containsString(plan.NeedsHuman, "non-scalar"))
}

func TestReconcileTypes_UnmodeledJSONExcusesTheGap(t *testing.T) {
	root := widgetFixture(t)
	sm := &SpecMap{Types: []TypeMapping{widgetTypeMapping()}}
	um := &Unmodeled{Properties: []PropertyExclusion{
		{Schema: "Widget", Property: "status", Reason: "needs design review", Owner: "eric"},
	}}

	spec := mustSpec(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"id":{"type":"string"},
			"name":{"type":["string","null"]},
			"status":{"type":"string"}
		},"required":["id","status"]}`,
	})

	plan := reconcileTypes(root, sm, um, spec, spec)
	require.Empty(t, plan.Actions)
	require.Empty(t, plan.NeedsHuman, "an unmodeled.json entry excuses even a required, otherwise-hard-failing gap")
}

func TestReconcileTypes_PropertyRemovalIsAlwaysNeedsHuman(t *testing.T) {
	root := widgetFixture(t)
	sm := &SpecMap{Types: []TypeMapping{widgetTypeMapping()}}
	um := &Unmodeled{}

	oldSpec := mustSpec(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"id":{"type":"string"},
			"name":{"type":["string","null"]},
			"legacy_field":{"type":"string"}
		},"required":["id"]}`,
	})
	newSpec := mustSpec(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"id":{"type":"string"},
			"name":{"type":["string","null"]}
		},"required":["id"]}`,
	})

	plan := reconcileTypes(root, sm, um, oldSpec, newSpec)
	require.True(t, containsString(plan.NeedsHuman, "legacy_field"))
	require.True(t, containsString(plan.NeedsHuman, "removal"))
}

func TestReconcileTypes_RequiredNessChangeIsNeedsHuman(t *testing.T) {
	root := widgetFixture(t)
	sm := &SpecMap{Types: []TypeMapping{widgetTypeMapping()}}
	um := &Unmodeled{}

	oldSpec := mustSpec(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"id":{"type":"string"},
			"name":{"type":["string","null"]}
		},"required":["id"]}`,
	})
	newSpec := mustSpec(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"id":{"type":"string"},
			"name":{"type":["string","null"]}
		},"required":["id","name"]}`,
	})

	plan := reconcileTypes(root, sm, um, oldSpec, newSpec)
	require.True(t, containsString(plan.NeedsHuman, "required-ness changed"))
}

func TestReconcileTypes_TypeChangeIsNeedsHuman(t *testing.T) {
	root := widgetFixture(t)
	sm := &SpecMap{Types: []TypeMapping{widgetTypeMapping()}}
	um := &Unmodeled{}

	oldSpec := mustSpec(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"id":{"type":"string"},
			"name":{"type":["string","null"]}
		},"required":["id"]}`,
	})
	newSpec := mustSpec(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"id":{"type":"string"},
			"name":{"type":["integer","null"]}
		},"required":["id"]}`,
	})

	plan := reconcileTypes(root, sm, um, oldSpec, newSpec)
	require.True(t, containsString(plan.NeedsHuman, "type changed"))
}

func TestReconcileTypes_NullableWideningIsNotATypeChange(t *testing.T) {
	// An untyped ("free-form") old property gaining an explicit nullable
	// type in the new spec is a documentation tightening, not a break; see
	// baseTypes' treatment of an absent "type" as unconstrained. This
	// mirrors the real created_at/updated_at drift observed between
	// spec-snapshot.json and spec-current.json in this repo.
	root := widgetFixture(t)
	sm := &SpecMap{Types: []TypeMapping{widgetTypeMapping()}}
	um := &Unmodeled{}

	oldSpec := mustSpec(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"id":{"type":"string"},
			"name":{"example":"foo"}
		},"required":["id"]}`,
	})
	newSpec := mustSpec(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"id":{"type":"string"},
			"name":{"type":["string","null"],"format":"date-time","example":"foo"}
		},"required":["id"]}`,
	})

	plan := reconcileTypes(root, sm, um, oldSpec, newSpec)
	require.Empty(t, plan.NeedsHuman)
}

func TestReconcileTypes_DiscriminatorAndNestedPropertiesAreExcluded(t *testing.T) {
	root := writeFixture(t, map[string]string{
		widgetGoFile: `package widgets

type CreateParams struct {
	Name string ` + "`json:\"name\"`" + `
}

type Tracking struct {
	Step string ` + "`json:\"step\"`" + `
}
`,
	})
	sm := &SpecMap{Types: []TypeMapping{{
		Spec:          "CreateWidgetIn",
		Discriminator: []string{"type"},
		SDK:           []SDKSite{{File: widgetGoFile, Symbol: "CreateParams"}},
		Nested: map[string]NestedMapping{
			"tracking": {SDK: []SDKSite{{File: widgetGoFile, Symbol: "Tracking"}}},
		},
	}}}
	um := &Unmodeled{}

	spec := mustSpec(t, map[string]string{
		"CreateWidgetIn": `{"type":"object","properties":{
			"type":{"type":"string","enum":["a","b"]},
			"name":{"type":"string"},
			"tracking":{"type":"object","properties":{"step":{"type":"string"}}}
		},"required":["type","name"]}`,
	})

	plan := reconcileTypes(root, sm, um, spec, spec)
	require.Empty(t, plan.NeedsHuman, "discriminator 'type' and nested 'tracking' must not be treated as missing top-level fields")
	require.Empty(t, plan.Actions)
}

func TestReconcileTypes_NestedSubObjectReconciledIndependently(t *testing.T) {
	root := writeFixture(t, map[string]string{
		widgetGoFile: `package widgets

type Widget struct {
	ID       string    ` + "`json:\"id\"`" + `
	Tracking *Tracking ` + "`json:\"tracking,omitempty\"`" + `
}

type Tracking struct {
	Step string ` + "`json:\"step\"`" + `
}
`,
	})
	sm := &SpecMap{Types: []TypeMapping{{
		Spec: "Widget",
		SDK:  []SDKSite{{File: widgetGoFile, Symbol: "Widget"}},
		Nested: map[string]NestedMapping{
			"tracking": {SDK: []SDKSite{{File: widgetGoFile, Symbol: "Tracking"}}},
		},
	}}}
	um := &Unmodeled{}

	spec := mustSpec(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"id":{"type":"string"},
			"tracking":{"type":"object","properties":{
				"step":{"type":"string"},
				"provider_reference":{"type":["string","null"]}
			}}
		},"required":["id"]}`,
	})

	plan := reconcileTypes(root, sm, um, spec, spec)
	require.Empty(t, plan.NeedsHuman)
	require.Len(t, plan.Actions, 1)
	require.Contains(t, plan.Actions[0].description, "provider_reference")
	require.Equal(t, widgetGoFile, plan.Actions[0].insert.File, "the nested field is added to the nested struct's own file, not the parent's")
}

func TestReconcileEnums_AppliesMissingMemberWithReexport(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"internal/types/enums.go": `package types

type Color string

const (
	ColorRed  Color = "red"
	ColorBlue Color = "blue"
)
`,
		"types.go": `package fixture

import "example.com/fixture/internal/types"

type Color = types.Color

const (
	ColorRed  = types.ColorRed
	ColorBlue = types.ColorBlue
)
`,
	})
	sm := &SpecMap{Enums: []EnumMapping{{
		SDK:      SDKSite{File: "internal/types/enums.go", Symbol: "Color"},
		Reexport: &SDKSite{File: "types.go", Symbol: "Color"},
		Spec:     SpecPointer{Schema: "Widget", Property: "color"},
	}}}
	um := &Unmodeled{}

	spec := mustSpec(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"color":{"type":"string","enum":["red","blue","green"]}
		}}`,
	})

	plan := reconcileEnums(root, sm, um, spec, spec)
	require.Empty(t, plan.NeedsHuman)
	require.Len(t, plan.Actions, 2, "one const insertion plus one root re-export insertion")

	var sawConst, sawReexport bool
	for _, a := range plan.Actions {
		require.Equal(t, "minor", a.bump)
		if a.insert.File == "internal/types/enums.go" {
			sawConst = true
			require.Contains(t, a.insert.Text, "ColorGreen")
			require.Contains(t, a.insert.Text, `"green"`)
		}
		if a.insert.File == "types.go" {
			sawReexport = true
			require.Contains(t, a.insert.Text, "ColorGreen = types.ColorGreen")
		}
	}
	require.True(t, sawConst)
	require.True(t, sawReexport)
}

func TestReconcileEnums_UnmodeledJSONExcusesMissingMember(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"internal/types/enums.go": `package types

type Color string

const (
	ColorRed Color = "red"
)
`,
	})
	sm := &SpecMap{Enums: []EnumMapping{{
		SDK:  SDKSite{File: "internal/types/enums.go", Symbol: "Color"},
		Spec: SpecPointer{Schema: "Widget", Property: "color"},
	}}}
	um := &Unmodeled{Enums: []EnumExclusion{
		{Symbol: "Color", Member: "green", Reason: "needs its own PR", Owner: "eric"},
	}}

	spec := mustSpec(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"color":{"type":"string","enum":["red","green"]}
		}}`,
	})

	plan := reconcileEnums(root, sm, um, spec, spec)
	require.Empty(t, plan.Actions)
	require.Empty(t, plan.NeedsHuman)
}

func TestReconcileEnums_MemberRemovalIsAlwaysNeedsHuman(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"internal/types/enums.go": `package types

type Color string

const (
	ColorRed   Color = "red"
	ColorGreen Color = "green"
)
`,
	})
	sm := &SpecMap{Enums: []EnumMapping{{
		SDK:  SDKSite{File: "internal/types/enums.go", Symbol: "Color"},
		Spec: SpecPointer{Schema: "Widget", Property: "color"},
	}}}
	um := &Unmodeled{}

	oldSpec := mustSpec(t, map[string]string{
		"Widget": `{"type":"object","properties":{"color":{"type":"string","enum":["red","green"]}}}`,
	})
	newSpec := mustSpec(t, map[string]string{
		"Widget": `{"type":"object","properties":{"color":{"type":"string","enum":["red"]}}}`,
	})

	plan := reconcileEnums(root, sm, um, oldSpec, newSpec)
	require.True(t, containsString(plan.NeedsHuman, "green"))
	require.True(t, containsString(plan.NeedsHuman, "removal"))
}

func TestReconcileEnums_InlineSchemaPointerIsSkippedNotFailed(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"internal/types/enums.go": `package types

type Currency string

const (
	CurrencyUSD Currency = "USD"
)
`,
	})
	sm := &SpecMap{Enums: []EnumMapping{{
		SDK:  SDKSite{File: "internal/types/enums.go", Symbol: "Currency"},
		Spec: SpecPointer{Schema: "inline:POST /fx#requestBody", Property: "to"},
	}}}
	um := &Unmodeled{}
	spec := mustSpec(t, map[string]string{})

	plan := reconcileEnums(root, sm, um, spec, spec)
	require.Empty(t, plan.Actions)
	require.Empty(t, plan.NeedsHuman, "an inline: schema pointer is informational-only, never a hard failure")
}

func TestReconcileEnums_WebhookEventItemsStyle(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"internal/types/enums.go": `package types

type WebhookEvent string

const (
	WebhookEventCustomerNew WebhookEvent = "customer.new"
)
`,
	})
	sm := &SpecMap{Enums: []EnumMapping{{
		SDK:  SDKSite{File: "internal/types/enums.go", Symbol: "WebhookEvent"},
		Spec: SpecPointer{Schema: "WebhookEndpoint", Property: "events", Items: true},
	}}}
	um := &Unmodeled{}

	spec := mustSpec(t, map[string]string{
		"WebhookEndpoint": `{"type":"object","properties":{
			"events":{"type":"array","items":{"type":"string","enum":["customer.new","customer.update"]}}
		}}`,
	})

	plan := reconcileEnums(root, sm, um, spec, spec)
	require.Empty(t, plan.NeedsHuman)
	require.Len(t, plan.Actions, 1)
	require.Contains(t, plan.Actions[0].insert.Text, "WebhookEventCustomerUpdate")
}

func TestCheckMapValidity_ReportsUnresolvedSymbol(t *testing.T) {
	root := widgetFixture(t)
	sm := &SpecMap{Types: []TypeMapping{{
		Spec: "Widget",
		SDK:  []SDKSite{{File: widgetGoFile, Symbol: "DoesNotExist"}},
	}}}

	issues := checkMapValidity(root, sm)
	require.True(t, containsString(issues, "DoesNotExist"))
}

func TestCheckUnclassifiedSchemas_FlagsAnUnmappedSchema(t *testing.T) {
	sm := &SpecMap{Types: []TypeMapping{{Spec: "Widget"}}}
	sm.Ignore.Schemas = []IgnoreEntry{{Schema: "Ignored", Reason: "not modeled"}}

	spec := specWithPaths(t, map[string]string{
		"Widget":  `{"type":"object","properties":{}}`,
		"Ignored": `{"type":"object","properties":{}}`,
		"Mystery": `{"type":"object","properties":{}}`,
	}, map[string]any{
		"/v1/widgets": map[string]any{"get": getOp("Widget")},
		"/v1/ignored": map[string]any{"get": getOp("Ignored")},
		"/v1/mystery": map[string]any{"get": getOp("Mystery")},
	})

	issues := checkUnclassifiedSchemas(sm, spec, nil)
	require.True(t, containsString(issues, "Mystery"))
	require.False(t, containsString(issues, "Widget"))
	require.False(t, containsString(issues, `"Ignored"`))
}
