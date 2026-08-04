package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// specWithPaths builds a minimal spec document with the given components
// schemas and paths, for operation-change tests.
func specWithPaths(t *testing.T, schemas map[string]string, paths map[string]any) *specDoc {
	t.Helper()
	sm := mustSpec(t, schemas)
	sm.raw["paths"] = paths
	return sm
}

func getOp(responseSchemaRef string) map[string]any {
	op := map[string]any{"tags": []any{"Payouts"}}
	if responseSchemaRef != "" {
		op["responses"] = map[string]any{
			"200": map[string]any{"content": map[string]any{"application/json": map[string]any{
				"schema": map[string]any{"$ref": "#/components/schemas/" + responseSchemaRef},
			}}},
		}
	} else {
		op["responses"] = map[string]any{"200": map[string]any{}}
	}
	return op
}

func TestCheckOperationChanges_NewOperationHardFailsByDefault(t *testing.T) {
	// No fixture SDK packages exist at all, so the deterministic generator
	// can never resolve a target package for this route and must fall back
	// to NEEDS_HUMAN, exactly like the pre-generator behavior.
	repoRoot := writeFixture(t, map[string]string{"go.mod": "module example.com/fixture\n\ngo 1.21\n"})

	sm := &SpecMap{Types: []TypeMapping{{Spec: "PayoutOut", SDK: []SDKSite{{File: "x", Symbol: "Y"}}}}}
	oldSpec := specWithPaths(t, map[string]string{"PayoutOut": `{"type":"object","properties":{}}`}, map[string]any{})
	newSpec := specWithPaths(t, map[string]string{"PayoutOut": `{"type":"object","properties":{}}`}, map[string]any{
		"/v1/instances/{instance_id}/payouts/new-thing": map[string]any{
			"get": getOp("PayoutOut"), // tagged Payouts, referencing an already-MAPPED (not ignored) schema
		},
	})

	issues := checkOperationChanges(repoRoot, sm, oldSpec, newSpec).NeedsHuman
	require.True(t, containsString(issues, "GET /v1/instances/{instance_id}/payouts/new-thing"))
	require.True(t, containsString(issues, "NEEDS_HUMAN"))
	require.True(t, containsString(issues, "hand-written client wiring"))
}

func TestCheckOperationChanges_NewOperationInAnIgnoredFamilyIsNotBlocking(t *testing.T) {
	repoRoot := writeFixture(t, map[string]string{"go.mod": "module example.com/fixture\n\ngo 1.21\n"})
	sm := &SpecMap{}
	sm.Ignore.Schemas = []IgnoreEntry{{Schema: "Rfi", Reason: "not modeled"}}

	tests := []struct {
		name string
		op   map[string]any
	}{
		{
			name: "referenced response schema is in ignore.schemas",
			op:   getOp("Rfi"),
		},
		{
			name: "no useful schema ref, but the tag exactly matches an ignored schema name",
			op:   map[string]any{"tags": []any{"Rfi"}, "responses": map[string]any{"200": map[string]any{}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldSpec := specWithPaths(t, map[string]string{"Rfi": `{"type":"object","properties":{}}`}, map[string]any{})
			newSpec := specWithPaths(t, map[string]string{"Rfi": `{"type":"object","properties":{}}`}, map[string]any{
				"/v1/instances/{instance_id}/rfi/new-endpoint": map[string]any{"get": tt.op},
			})

			plan := checkOperationChanges(repoRoot, sm, oldSpec, newSpec)
			require.Empty(t, plan.NeedsHuman, "a new operation reachable only through an already-ignored schema/tag must not block the pipeline")
		})
	}
}

func TestCheckOperationChanges_TagPrefixAloneIsNotEnoughToExempt(t *testing.T) {
	// UploadIn/UploadOut are real, mapped schemas alongside the ignored
	// UploadAnalyzeIn/UploadAnalyzeOut; a tag "Upload" must not exempt an
	// operation by loosely prefix-matching the ignored name, since that
	// would also wrongly exempt a genuinely new Upload-family operation.
	repoRoot := writeFixture(t, map[string]string{"go.mod": "module example.com/fixture\n\ngo 1.21\n"})
	sm := &SpecMap{}
	sm.Ignore.Schemas = []IgnoreEntry{{Schema: "UploadAnalyzeOut", Reason: "not modeled"}}

	oldSpec := specWithPaths(t, map[string]string{"UploadOut": `{"type":"object","properties":{}}`}, map[string]any{})
	newSpec := specWithPaths(t, map[string]string{"UploadOut": `{"type":"object","properties":{}}`}, map[string]any{
		"/v1/upload/new-variant": map[string]any{
			"post": map[string]any{"tags": []any{"Upload"}, "responses": map[string]any{
				"200": map[string]any{"content": map[string]any{"application/json": map[string]any{
					"schema": map[string]any{"$ref": "#/components/schemas/UploadOut"},
				}}},
			}},
		},
	})

	issues := checkOperationChanges(repoRoot, sm, oldSpec, newSpec).NeedsHuman
	require.True(t, containsString(issues, "POST /v1/upload/new-variant"), "referencing a MAPPED schema (UploadOut) must still hard-fail even though the tag loosely resembles an ignored family name")
}

func TestCheckOperationChanges_RemovedOperationIsAlwaysNeedsHuman(t *testing.T) {
	repoRoot := writeFixture(t, map[string]string{"go.mod": "module example.com/fixture\n\ngo 1.21\n"})
	sm := &SpecMap{}
	oldSpec := specWithPaths(t, map[string]string{}, map[string]any{
		"/v1/instances/{instance_id}/widgets": map[string]any{"get": map[string]any{"responses": map[string]any{"200": map[string]any{}}}},
	})
	newSpec := specWithPaths(t, map[string]string{}, map[string]any{})

	issues := checkOperationChanges(repoRoot, sm, oldSpec, newSpec).NeedsHuman
	require.True(t, containsString(issues, "GET /v1/instances/{instance_id}/widgets"))
	require.True(t, containsString(issues, "removal is always a hard fail"))
}

func TestCheckOperationChanges_NoChangeIsClean(t *testing.T) {
	repoRoot := writeFixture(t, map[string]string{"go.mod": "module example.com/fixture\n\ngo 1.21\n"})
	sm := &SpecMap{}
	spec := specWithPaths(t, map[string]string{}, map[string]any{
		"/v1/widgets": map[string]any{"get": map[string]any{"responses": map[string]any{"200": map[string]any{}}}},
	})
	plan := checkOperationChanges(repoRoot, sm, spec, spec)
	require.Empty(t, plan.NeedsHuman)
	require.Empty(t, plan.Actions)
}

func TestCheckUnclassifiedSchemas_NewSchemaAbsentFromCommittedSnapshotIsNeedsHuman(t *testing.T) {
	// The broader existing check (every schema in the target spec must be
	// mapped or ignored) already covers "a new reachable schema" as a
	// special case: it does not matter whether the schema is brand new or
	// was simply never classified, both are NEEDS_HUMAN either way. This
	// test pins the old-vs-new framing specifically.
	sm := &SpecMap{Types: []TypeMapping{{Spec: "Widget"}}}
	sm.Ignore.Schemas = []IgnoreEntry{{Schema: "Ignored", Reason: "not modeled"}}

	newSpec := specWithPaths(t, map[string]string{
		"Widget":        `{"type":"object","properties":{}}`,
		"Ignored":       `{"type":"object","properties":{}}`,
		"BrandNewThing": `{"type":"object","properties":{}}`, // present only in the new spec
	}, map[string]any{
		"/v1/widgets": map[string]any{"get": getOp("Widget")},
		"/v1/ignored": map[string]any{"get": getOp("Ignored")},
		"/v1/new":     map[string]any{"get": getOp("BrandNewThing")},
	})

	issues := checkUnclassifiedSchemas(sm, newSpec, nil)
	require.True(t, containsString(issues, "BrandNewThing"))
	require.True(t, containsString(issues, "NEEDS_HUMAN"))
}
