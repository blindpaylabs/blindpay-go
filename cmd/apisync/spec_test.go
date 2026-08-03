package main

import (
	"encoding/json"
	"strings"
)

// mustSpec builds a minimal spec document ({"components":{"schemas":{...}}})
// from a schemaName -> raw-JSON-schema-body map, for reconciliation tests.
func mustSpec(t testingT, schemas map[string]string) *specDoc {
	t.Helper()
	parsed := map[string]any{}
	for name, raw := range schemas {
		var v any
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			t.Fatalf("schema %s: invalid JSON: %v", name, err)
		}
		parsed[name] = v
	}
	return loadSpecDoc(map[string]any{
		"components": map[string]any{"schemas": parsed},
		"paths":      map[string]any{},
	})
}

// testingT is the subset of *testing.T mustSpec needs.
type testingT interface {
	Helper()
	Fatalf(format string, args ...any)
}

// containsString reports whether any element of list contains substr.
func containsString(list []string, substr string) bool {
	for _, s := range list {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
