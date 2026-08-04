package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// rawRef is a component schema definition that $refs another schema, for
// building transitive-reachability fixtures.
func rawRef(target string) map[string]any {
	return map[string]any{"$ref": "#/components/schemas/" + target}
}

func specDocWithRaw(schemas map[string]any, extra map[string]any) *specDoc {
	raw := map[string]any{"components": map[string]any{"schemas": schemas}}
	for k, v := range extra {
		raw[k] = v
	}
	return loadSpecDoc(raw)
}

func TestReachableSchemas_PathOnly(t *testing.T) {
	spec := specDocWithRaw(map[string]any{
		"Reached": map[string]any{"type": "object", "properties": map[string]any{}},
		"Orphan":  map[string]any{"type": "object", "properties": map[string]any{}},
	}, map[string]any{
		"paths": map[string]any{
			"/v1/widgets": map[string]any{"get": map[string]any{
				"responses": map[string]any{"200": map[string]any{"content": map[string]any{"application/json": map[string]any{
					"schema": rawRef("Reached"),
				}}}},
			}},
		},
	})

	reachable := reachableSchemas(spec)
	require.True(t, reachable["Reached"])
	require.False(t, reachable["Orphan"])
}

func TestReachableSchemas_WebhookOnly(t *testing.T) {
	spec := specDocWithRaw(map[string]any{
		"Reached": map[string]any{"type": "object", "properties": map[string]any{}},
		"Orphan":  map[string]any{"type": "object", "properties": map[string]any{}},
	}, map[string]any{
		"webhooks": map[string]any{
			"customer.updated": map[string]any{"post": map[string]any{
				"requestBody": map[string]any{"content": map[string]any{"application/json": map[string]any{
					"schema": rawRef("Reached"),
				}}},
			}},
		},
	})

	reachable := reachableSchemas(spec)
	require.True(t, reachable["Reached"])
	require.False(t, reachable["Orphan"])
}

func TestReachableSchemas_ParameterOnly(t *testing.T) {
	// A shared components.parameters entry (not itself a components.schemas
	// section) can carry a $ref to a schema without any path ever directly
	// referencing that schema in its own requestBody/responses.
	spec := specDocWithRaw(map[string]any{
		"Reached": map[string]any{"type": "object", "properties": map[string]any{}},
		"Orphan":  map[string]any{"type": "object", "properties": map[string]any{}},
	}, map[string]any{
		"components": map[string]any{
			"parameters": map[string]any{
				"FilterParam": map[string]any{"schema": rawRef("Reached")},
			},
		},
	})

	reachable := reachableSchemas(spec)
	require.True(t, reachable["Reached"])
	require.False(t, reachable["Orphan"])
}

func TestReachableSchemas_Transitive(t *testing.T) {
	// Root is reachable from a path; Root refs Middle; Middle refs Leaf.
	// Both Middle and Leaf must end up reachable even though neither is
	// directly referenced by any path/webhook/parameter.
	spec := specDocWithRaw(map[string]any{
		"Root":   map[string]any{"type": "object", "properties": map[string]any{"middle": rawRef("Middle")}},
		"Middle": map[string]any{"type": "object", "properties": map[string]any{"leaf": rawRef("Leaf")}},
		"Leaf":   map[string]any{"type": "object", "properties": map[string]any{}},
		"Orphan": map[string]any{"type": "object", "properties": map[string]any{}},
	}, map[string]any{
		"paths": map[string]any{
			"/v1/widgets": map[string]any{"get": map[string]any{
				"responses": map[string]any{"200": map[string]any{"content": map[string]any{"application/json": map[string]any{
					"schema": rawRef("Root"),
				}}}},
			}},
		},
	})

	reachable := reachableSchemas(spec)
	require.True(t, reachable["Root"])
	require.True(t, reachable["Middle"])
	require.True(t, reachable["Leaf"])
	require.False(t, reachable["Orphan"])
}

func TestReachableSchemas_Orphan(t *testing.T) {
	// No paths, no webhooks, no non-schema components at all: every schema
	// is an orphan.
	spec := specDocWithRaw(map[string]any{
		"Orphan": map[string]any{"type": "object", "properties": map[string]any{}},
	}, map[string]any{
		"paths": map[string]any{},
	})

	reachable := reachableSchemas(spec)
	require.False(t, reachable["Orphan"])
}

// TestCheckUnclassifiedSchemas_OrphanSchemaProducesNoWork is the integration
// point: an orphan schema that is neither mapped nor ignored must not be
// reported by checkUnclassifiedSchemas at all (previously the patcher
// stopped on any unclassified schema, reachable or not).
func TestCheckUnclassifiedSchemas_OrphanSchemaProducesNoWork(t *testing.T) {
	sm := &SpecMap{Types: []TypeMapping{{Spec: "Widget"}}}
	spec := specDocWithRaw(map[string]any{
		"Widget":       map[string]any{"type": "object", "properties": map[string]any{}},
		"OrphanUnused": map[string]any{"type": "object", "properties": map[string]any{}},
	}, map[string]any{
		"paths": map[string]any{
			"/v1/widgets": map[string]any{"get": map[string]any{
				"responses": map[string]any{"200": map[string]any{"content": map[string]any{"application/json": map[string]any{
					"schema": rawRef("Widget"),
				}}}},
			}},
		},
	})

	issues := checkUnclassifiedSchemas(sm, spec)
	require.Empty(t, issues, "an orphan schema unreachable from any path/webhook/parameter must produce no work")
}
