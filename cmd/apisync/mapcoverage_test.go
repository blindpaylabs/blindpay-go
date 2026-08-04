package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// reachableSpecWithPaths is specWithPaths plus a synthetic path referencing
// every given schema, so checkEnumCoverage/checkNestedObjectCoverage's
// reachability gate does not silently exempt the fixture's schemas.
func reachableSpecWithPaths(t *testing.T, schemas map[string]string) *specDoc {
	t.Helper()
	paths := map[string]any{}
	i := 0
	for name := range schemas {
		i++
		paths[fmt.Sprintf("/v1/fixture%d", i)] = map[string]any{"get": getOp(name)}
	}
	return specWithPaths(t, schemas, paths)
}

func TestIsEnumConstrained(t *testing.T) {
	items, ok := isEnumConstrained(map[string]any{"type": "string", "enum": []any{"a", "b"}})
	require.True(t, ok)
	require.False(t, items)

	items, ok = isEnumConstrained(map[string]any{"type": "array", "items": map[string]any{"type": "string", "enum": []any{"a"}}})
	require.True(t, ok)
	require.True(t, items)

	items, ok = isEnumConstrained(map[string]any{"anyOf": []any{
		map[string]any{"type": "string", "enum": []any{"a"}},
		map[string]any{"type": "string", "enum": []any{"b"}},
	}})
	require.True(t, ok)
	require.False(t, items)

	_, ok = isEnumConstrained(map[string]any{"type": "string"})
	require.False(t, ok)

	_, ok = isEnumConstrained(map[string]any{"anyOf": []any{
		map[string]any{"type": "string", "enum": []any{"a"}},
		map[string]any{"type": "integer"},
	}})
	require.False(t, ok, "a mixed anyOf (not every branch enum-only) is not enum-constrained")
}

func TestCheckEnumCoverage_CoveredPropertyIsClean(t *testing.T) {
	sm := &SpecMap{
		Types: []TypeMapping{{Spec: "Widget", SDK: []SDKSite{{File: "x", Symbol: "Y"}}}},
		Enums: []EnumMapping{{SDK: SDKSite{File: "a", Symbol: "Color"}, Spec: SpecPointer{Schema: "Widget", Property: "color"}}},
	}
	um := &Unmodeled{}
	spec := reachableSpecWithPaths(t, map[string]string{
		"Widget": `{"type":"object","properties":{"color":{"type":"string","enum":["red","blue"]}}}`,
	})

	issues := checkEnumCoverage(sm, um, spec)
	require.Empty(t, issues)
}

// TestCheckEnumCoverage_MutationNewUnmappedEnumPropertyFailsCI proves the
// gate actually blocks: a brand new enum-constrained property with no
// enums[] mapping and no unmodeled.json enum_coverage exclusion must fail.
func TestCheckEnumCoverage_MutationNewUnmappedEnumPropertyFailsCI(t *testing.T) {
	sm := &SpecMap{Types: []TypeMapping{{Spec: "Widget", SDK: []SDKSite{{File: "x", Symbol: "Y"}}}}}
	um := &Unmodeled{}
	spec := reachableSpecWithPaths(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"status":{"type":"string","enum":["on","off"]}
		}}`,
	})

	issues := checkEnumCoverage(sm, um, spec)
	require.Len(t, issues, 1)
	require.Contains(t, issues[0], "Widget.status")
	require.Contains(t, issues[0], "NEEDS_HUMAN")
}

func TestCheckEnumCoverage_RecordedExclusionIsClean(t *testing.T) {
	sm := &SpecMap{Types: []TypeMapping{{Spec: "Widget", SDK: []SDKSite{{File: "x", Symbol: "Y"}}}}}
	um := &Unmodeled{EnumCoverage: []PropertyExclusion{
		{Schema: "Widget", Property: "status", Reason: "baseline 2026-08-04: not yet modeled, needs review", Owner: "eric"},
	}}
	spec := reachableSpecWithPaths(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"status":{"type":"string","enum":["on","off"]}
		}}`,
	})

	issues := checkEnumCoverage(sm, um, spec)
	require.Empty(t, issues)
}

func TestCheckEnumCoverage_DiscriminatorPropertyIsExempt(t *testing.T) {
	sm := &SpecMap{Types: []TypeMapping{{
		Spec:          "Widget",
		SDK:           []SDKSite{{File: "x", Symbol: "Y"}},
		Discriminator: []string{"type"},
	}}}
	um := &Unmodeled{}
	spec := reachableSpecWithPaths(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"type":{"type":"string","enum":["a","b"]}
		}}`,
	})

	issues := checkEnumCoverage(sm, um, spec)
	require.Empty(t, issues, "a discriminator-injected property is a different, already-reviewed design decision, not a coverage gap")
}

func TestCheckEnumCoverage_UnreachableSchemaIsExempt(t *testing.T) {
	sm := &SpecMap{Types: []TypeMapping{{Spec: "Widget", SDK: []SDKSite{{File: "x", Symbol: "Y"}}}}}
	um := &Unmodeled{}
	// mustSpec (unlike reachableSpecWithPaths) has no paths at all, so
	// Widget is an orphan: no work either way.
	spec := mustSpec(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"status":{"type":"string","enum":["on","off"]}
		}}`,
	})

	issues := checkEnumCoverage(sm, um, spec)
	require.Empty(t, issues)
}

func TestCheckEnumCoverage_NestednEnumPropertyUsesDottedPath(t *testing.T) {
	sm := &SpecMap{Types: []TypeMapping{{
		Spec: "Widget",
		SDK:  []SDKSite{{File: "x", Symbol: "Y"}},
		Nested: map[string]NestedMapping{
			"tracking": {SDK: []SDKSite{{File: "x", Symbol: "Tracking"}}},
		},
	}}}
	um := &Unmodeled{}
	spec := reachableSpecWithPaths(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"tracking":{"type":"object","properties":{
				"step":{"type":"string","enum":["queued","done"]}
			}}
		}}`,
	})

	issues := checkEnumCoverage(sm, um, spec)
	require.Len(t, issues, 1)
	require.Contains(t, issues[0], `property "tracking.step"`)

	sm.Enums = []EnumMapping{{SDK: SDKSite{File: "a", Symbol: "Step"}, Spec: SpecPointer{Schema: "Widget", Property: "tracking.step"}}}
	issues = checkEnumCoverage(sm, um, spec)
	require.Empty(t, issues, "a dotted-path enums[] mapping covers the nested property")
}

func TestCheckNestedObjectCoverage_MappedShapeIsClean(t *testing.T) {
	sm := &SpecMap{Types: []TypeMapping{{
		Spec: "Widget",
		SDK:  []SDKSite{{File: "x", Symbol: "Y"}},
		Nested: map[string]NestedMapping{
			"tracking": {SDK: []SDKSite{{File: "x", Symbol: "Tracking"}}},
		},
	}}}
	um := &Unmodeled{}
	spec := reachableSpecWithPaths(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"tracking":{"type":"object","properties":{"step":{"type":"string"}}}
		}}`,
	})

	issues := checkNestedObjectCoverage(sm, um, spec)
	require.Empty(t, issues)
}

// TestCheckNestedObjectCoverage_MutationNewUnmappedInlineObjectFailsCI proves
// the gate blocks a brand new inline object shape with no nested[] entry and
// no unmodeled.json nested_objects exclusion.
func TestCheckNestedObjectCoverage_MutationNewUnmappedInlineObjectFailsCI(t *testing.T) {
	sm := &SpecMap{Types: []TypeMapping{{Spec: "Widget", SDK: []SDKSite{{File: "x", Symbol: "Y"}}}}}
	um := &Unmodeled{}
	spec := reachableSpecWithPaths(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"metadata":{"type":"object","properties":{"key":{"type":"string"}}}
		}}`,
	})

	issues := checkNestedObjectCoverage(sm, um, spec)
	require.Len(t, issues, 1)
	require.Contains(t, issues[0], "Widget.metadata")
	require.Contains(t, issues[0], "NEEDS_HUMAN")
}

func TestCheckNestedObjectCoverage_ArrayItemObjectShapeIsDetected(t *testing.T) {
	sm := &SpecMap{Types: []TypeMapping{{Spec: "Widget", SDK: []SDKSite{{File: "x", Symbol: "Y"}}}}}
	um := &Unmodeled{}
	spec := reachableSpecWithPaths(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"owners":{"type":"array","items":{"type":"object","properties":{"name":{"type":"string"}}}}
		}}`,
	})

	issues := checkNestedObjectCoverage(sm, um, spec)
	require.Len(t, issues, 1)
	require.Contains(t, issues[0], "Widget.owners")
}

func TestCheckNestedObjectCoverage_RecordedExclusionIsClean(t *testing.T) {
	sm := &SpecMap{Types: []TypeMapping{{Spec: "Widget", SDK: []SDKSite{{File: "x", Symbol: "Y"}}}}}
	um := &Unmodeled{NestedObjects: []PropertyExclusion{
		{Schema: "Widget", Property: "metadata", Reason: "baseline 2026-08-04: not yet modeled, needs review", Owner: "eric"},
	}}
	spec := reachableSpecWithPaths(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"metadata":{"type":"object","properties":{"key":{"type":"string"}}}
		}}`,
	})

	issues := checkNestedObjectCoverage(sm, um, spec)
	require.Empty(t, issues)
}

func TestCheckNestedObjectCoverage_RefdObjectIsNotInline(t *testing.T) {
	sm := &SpecMap{Types: []TypeMapping{{Spec: "Widget", SDK: []SDKSite{{File: "x", Symbol: "Y"}}}}}
	um := &Unmodeled{}
	spec := reachableSpecWithPaths(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"owner":{"$ref":"#/components/schemas/Owner"}
		}}`,
		"Owner": `{"type":"object","properties":{"name":{"type":"string"}}}`,
	})

	issues := checkNestedObjectCoverage(sm, um, spec)
	require.Empty(t, issues, "a $ref'd schema is a separate top-level mapping concern, not an inline shape")
}

func TestCheckNestedObjectCoverage_DepthTwoShapeIsLedgerOnly(t *testing.T) {
	sm := &SpecMap{Types: []TypeMapping{{
		Spec: "Widget",
		SDK:  []SDKSite{{File: "x", Symbol: "Y"}},
		Nested: map[string]NestedMapping{
			"tracking": {SDK: []SDKSite{{File: "x", Symbol: "Tracking"}}},
		},
	}}}
	um := &Unmodeled{}
	spec := reachableSpecWithPaths(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"tracking":{"type":"object","properties":{
				"detail":{"type":"object","properties":{"note":{"type":"string"}}}
			}}
		}}`,
	})

	issues := checkNestedObjectCoverage(sm, um, spec)
	require.Len(t, issues, 1)
	require.Contains(t, issues[0], "Tracking", "depth-2 shapes are keyed by the covering nested mapping's SDK symbol")

	um.NestedObjects = []PropertyExclusion{
		{Schema: "Tracking", Property: "detail", Reason: "baseline 2026-08-04: not yet modeled, needs review", Owner: "eric"},
	}
	issues = checkNestedObjectCoverage(sm, um, spec)
	require.Empty(t, issues)
}

// TestIntegration_MutationNewUnmappedEnumPropertyFailsCICheck is the
// end-to-end proof for checkEnumCoverage: -check against a target spec that
// adds a brand new enum-constrained property, with no enums[] mapping and
// no unmodeled.json enum_coverage exclusion, must fail (not silently pass).
func TestIntegration_MutationNewUnmappedEnumPropertyFailsCICheck(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"go.mod":      "module example.com/mapcoverage\n\ngo 1.21\n",
		"blindpay.go": "package mapcoverage\n\nconst Version = \"1.0.0\"\n",
		"widgets/client.go": `package widgets

type Widget struct {
	ID string ` + "`json:\"id\"`" + `
}
`,
		".api-sync/spec-map.json": `{
			"enums": [],
			"types": [{"spec": "Widget", "sdk": [{"file": "widgets/client.go", "symbol": "Widget"}]}],
			"ignore": {"schemas": []}
		}`,
		".api-sync/unmodeled.json": `{"properties": [], "enums": []}`,
		".api-sync/spec-snapshot.json": `{
			"paths": {"/v1/widgets": {"get": {"responses": {"200": {"content": {"application/json": {"schema": {"$ref": "#/components/schemas/Widget"}}}}}}}},
			"components": {"schemas": {"Widget": {"type": "object", "properties": {"id": {"type": "string"}}}}}
		}`,
	})

	quiet, err := reconcileAndMaybeApply(root, syncOptions{Check: true})
	require.NoError(t, err)
	require.True(t, quiet, "baseline with no enum-constrained properties must be clean")

	// Mutate: the target spec adds a new enum-constrained property with no
	// map entry and no ledger exclusion.
	require.NoError(t, os.WriteFile(filepath.Join(root, ".api-sync", "spec-current.json"), []byte(`{
		"paths": {"/v1/widgets": {"get": {"responses": {"200": {"content": {"application/json": {"schema": {"$ref": "#/components/schemas/Widget"}}}}}}}},
		"components": {"schemas": {"Widget": {"type": "object", "properties": {
			"id": {"type": "string"},
			"status": {"type": "string", "enum": ["on", "off"]}
		}}}}
	}`), 0o644))

	_, err = reconcileAndMaybeApply(root, syncOptions{Check: true, SpecPath: ".api-sync/spec-current.json"})
	require.Error(t, err, "checking the mutated spec directly must fail")

	// This is how it actually gates in production: the api-sync workflow
	// runs -apply against the fetched spec-current.json (see
	// .github/workflows/api-sync.yml); NEEDS_HUMAN must block it, atomically,
	// before anything is written or the snapshot is refreshed.
	before := readFixture(t, root, ".api-sync/spec-snapshot.json")
	_, err = reconcileAndMaybeApply(root, syncOptions{Apply: true})
	require.Error(t, err, "an unmapped, unrecorded enum-constrained property must block -apply")
	require.Equal(t, before, readFixture(t, root, ".api-sync/spec-snapshot.json"), "a blocked apply must never refresh the snapshot")
}

// TestIntegration_MutationNewUnmappedNestedObjectFailsCICheck is the same
// end-to-end proof for checkNestedObjectCoverage.
func TestIntegration_MutationNewUnmappedNestedObjectFailsCICheck(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"go.mod":      "module example.com/mapcoverage2\n\ngo 1.21\n",
		"blindpay.go": "package mapcoverage2\n\nconst Version = \"1.0.0\"\n",
		"widgets/client.go": `package widgets

type Widget struct {
	ID string ` + "`json:\"id\"`" + `
}
`,
		".api-sync/spec-map.json": `{
			"enums": [],
			"types": [{"spec": "Widget", "sdk": [{"file": "widgets/client.go", "symbol": "Widget"}]}],
			"ignore": {"schemas": []}
		}`,
		".api-sync/unmodeled.json": `{"properties": [], "enums": []}`,
		".api-sync/spec-snapshot.json": `{
			"paths": {"/v1/widgets": {"get": {"responses": {"200": {"content": {"application/json": {"schema": {"$ref": "#/components/schemas/Widget"}}}}}}}},
			"components": {"schemas": {"Widget": {"type": "object", "properties": {"id": {"type": "string"}}}}}
		}`,
	})

	quiet, err := reconcileAndMaybeApply(root, syncOptions{Check: true})
	require.NoError(t, err)
	require.True(t, quiet)

	// Mutate: the target spec adds a new inline object shape with no
	// nested[] map entry and no ledger exclusion.
	require.NoError(t, os.WriteFile(filepath.Join(root, ".api-sync", "spec-current.json"), []byte(`{
		"paths": {"/v1/widgets": {"get": {"responses": {"200": {"content": {"application/json": {"schema": {"$ref": "#/components/schemas/Widget"}}}}}}}},
		"components": {"schemas": {"Widget": {"type": "object", "properties": {
			"id": {"type": "string"},
			"metadata": {"type": "object", "properties": {"key": {"type": "string"}}}
		}}}}
	}`), 0o644))

	_, err = reconcileAndMaybeApply(root, syncOptions{Check: true, SpecPath: ".api-sync/spec-current.json"})
	require.Error(t, err, "checking the mutated spec directly must fail")

	before := readFixture(t, root, ".api-sync/spec-snapshot.json")
	_, err = reconcileAndMaybeApply(root, syncOptions{Apply: true})
	require.Error(t, err, "an unmapped, unrecorded inline object shape must block -apply")
	require.Equal(t, before, readFixture(t, root, ".api-sync/spec-snapshot.json"), "a blocked apply must never refresh the snapshot")
}
