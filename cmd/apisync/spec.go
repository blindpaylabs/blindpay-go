package main

import (
	"fmt"
	"strings"
)

// specDoc wraps a parsed OpenAPI document (generic map, no schema library).
type specDoc struct {
	raw map[string]any
}

func loadSpecDoc(raw map[string]any) *specDoc {
	return &specDoc{raw: raw}
}

func (d *specDoc) schemas() map[string]any {
	if d == nil || d.raw == nil {
		return nil
	}
	components, _ := d.raw["components"].(map[string]any)
	if components == nil {
		return nil
	}
	schemas, _ := components["schemas"].(map[string]any)
	return schemas
}

// schema returns the named component schema, trying fallbacks in order.
func (d *specDoc) schema(name string, fallbacks ...string) (map[string]any, string, bool) {
	schemas := d.schemas()
	if schemas == nil {
		return nil, "", false
	}
	for _, n := range append([]string{name}, fallbacks...) {
		if s, ok := schemas[n].(map[string]any); ok {
			return s, n, true
		}
	}
	return nil, "", false
}

func properties(schema map[string]any) map[string]any {
	if schema == nil {
		return nil
	}
	props, _ := schema["properties"].(map[string]any)
	return props
}

func requiredSet(schema map[string]any) map[string]bool {
	out := map[string]bool{}
	if schema == nil {
		return out
	}
	req, _ := schema["required"].([]any)
	for _, r := range req {
		if s, ok := r.(string); ok {
			out[s] = true
		}
	}
	return out
}

// resolveProperty walks a dotted property path (e.g. "tracking_payment.step")
// starting from a schema's top-level properties, descending into nested
// "properties" objects. Returns the leaf property definition.
func resolveProperty(schema map[string]any, dottedPath string) map[string]any {
	cur := schema
	parts := strings.Split(dottedPath, ".")
	var leaf map[string]any
	for i, part := range parts {
		props := properties(cur)
		if props == nil {
			return nil
		}
		def, ok := props[part].(map[string]any)
		if !ok {
			return nil
		}
		leaf = def
		if i < len(parts)-1 {
			cur = def
		}
	}
	return leaf
}

// enumMembers extracts the string enum members for a SpecPointer: from
// property.enum (Items false) or property.items.enum (Items true).
func enumMembers(schema map[string]any, ptr SpecPointer) ([]string, error) {
	def := schema
	if ptr.Property != "" {
		def = resolveProperty(schema, ptr.Property)
		if def == nil {
			return nil, fmt.Errorf("property %q not found on schema %q", ptr.Property, ptr.Schema)
		}
	}
	if ptr.Items {
		items, _ := def["items"].(map[string]any)
		def = items
	}
	if def == nil {
		return nil, fmt.Errorf("enum anchor resolved to nothing for schema %q property %q", ptr.Schema, ptr.Property)
	}
	raw, ok := def["enum"].([]any)
	if !ok {
		return nil, fmt.Errorf("no enum array at schema %q property %q (items=%v)", ptr.Schema, ptr.Property, ptr.Items)
	}
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out, nil
}

// baseTypes normalizes a JSON Schema "type" value (string, array of
// strings, or absent) into a sorted, deduplicated set of non-null base
// types. An absent "type" (untyped/free-form in this spec's dialect, e.g.
// early created_at/updated_at fields before format tightening) yields an
// empty set, meaning "unconstrained" -- never treated as a mismatch against
// a later, more specific type, only against another concrete type.
func baseTypes(def map[string]any) []string {
	if def == nil {
		return nil
	}
	var raw []any
	switch t := def["type"].(type) {
	case string:
		raw = []any{t}
	case []any:
		raw = t
	default:
		return nil
	}
	set := map[string]bool{}
	for _, v := range raw {
		s, ok := v.(string)
		if !ok || s == "null" {
			continue
		}
		set[s] = true
	}
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	return out
}

func sameStringSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	m := map[string]bool{}
	for _, s := range a {
		m[s] = true
	}
	for _, s := range b {
		if !m[s] {
			return false
		}
	}
	return true
}
