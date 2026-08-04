package main

import (
	"fmt"
	"sort"
)

// isEnumConstrained reports whether a property definition constrains its
// value to an enumerated set, in every shape this spec uses: a direct
// "enum", an array's "items.enum" (items is true in that case), or an
// anyOf/oneOf where every branch is itself enum-only (a discriminated
// union of literal string variants, still effectively one enum).
func isEnumConstrained(def map[string]any) (items bool, ok bool) {
	if def == nil {
		return false, false
	}
	if _, has := def["enum"]; has {
		return false, true
	}
	if itemsDef, hasItems := def["items"].(map[string]any); hasItems {
		if _, has := itemsDef["enum"]; has {
			return true, true
		}
	}
	for _, key := range [...]string{"anyOf", "oneOf"} {
		branches, isSlice := def[key].([]any)
		if !isSlice || len(branches) == 0 {
			continue
		}
		allEnum := true
		for _, b := range branches {
			bm, isMap := b.(map[string]any)
			if !isMap {
				allEnum = false
				break
			}
			if _, has := bm["enum"]; !has {
				allEnum = false
				break
			}
		}
		if allEnum {
			return false, true
		}
	}
	return false, false
}

// checkEnumCoverage is a blocking check: every enum-constrained property
// (including items.enum and anyOf/oneOf-of-enums) on every mapped, reachable
// schema must resolve to a mapped enum symbol (an enums[] entry whose
// spec.schema/spec.property match, using the same dotted-path convention as
// nested enums like "tracking_payment.step") or a recorded
// unmodeled.json enum_coverage exclusion. discriminator-injected properties
// are exempt: the SDK deliberately models those as a literal, not a stored
// enum-typed field, which is a different, already-reviewed design decision.
func checkEnumCoverage(sm *SpecMap, um *Unmodeled, spec *specDoc) []string {
	reachable := reachableSchemas(spec)

	enumCovered := map[string]bool{}
	for _, e := range sm.Enums {
		enumCovered[e.Spec.Schema+"\x00"+e.Spec.Property] = true
	}

	var issues []string
	var walk func(schemaKey, label, pathPrefix string, schema map[string]any, exclude map[string]bool, nested map[string]NestedMapping)
	walk = func(schemaKey, label, pathPrefix string, schema map[string]any, exclude map[string]bool, nested map[string]NestedMapping) {
		props := properties(schema)
		names := make([]string, 0, len(props))
		for n := range props {
			names = append(names, n)
		}
		sort.Strings(names)

		for _, name := range names {
			if exclude[name] {
				continue
			}
			def, _ := props[name].(map[string]any)
			dottedProp := name
			if pathPrefix != "" {
				dottedProp = pathPrefix + "." + name
			}

			if _, enumLike := isEnumConstrained(def); enumLike {
				key := schemaKey + "\x00" + dottedProp
				if !enumCovered[key] && !um.excusesEnumCoverage(schemaKey, dottedProp) {
					issues = append(issues, fmt.Sprintf(
						"NEEDS_HUMAN: %s.%s: enum-constrained property is not tied to any mapped enum symbol; "+
							"add an enums[] mapping (spec %q property %q) or record an unmodeled.json enum_coverage exclusion",
						label, name, schemaKey, dottedProp))
				}
			}

			// Recurse into a mapped nested[] sub-object so an enum-shaped
			// field nested inside it is not silently skipped.
			if n, isNested := nested[name]; isNested {
				if sub := resolveProperty(schema, name); sub != nil {
					walk(schemaKey, label+"."+name, dottedProp, sub, nil, nil)
				}
				_ = n
			}
		}
	}

	for _, t := range sm.Types {
		if !reachable[t.Spec] {
			continue
		}
		schema, _, ok := spec.schema(t.Spec)
		if !ok {
			continue // already reported by reconcileTypes
		}
		exclude := map[string]bool{}
		for _, d := range t.Discriminator {
			exclude[d] = true
		}
		walk(t.Spec, t.Spec, "", schema, exclude, t.Nested)
	}

	sort.Strings(issues)
	return issues
}

// inlineObjectShape returns the object schema a property's own definition
// describes inline (never through a $ref): either the property itself when
// it is an object with declared properties, or an array property's item
// schema when the item is itself an inline object. ok is false for a $ref'd
// object (that's a separate top-level schema with its own mapping entry,
// not an inline shape) or anything that isn't object-shaped.
func inlineObjectShape(def map[string]any) (shape map[string]any, ok bool) {
	if def == nil {
		return nil, false
	}
	if _, hasRef := def["$ref"]; hasRef {
		return nil, false
	}
	isBase := func(d map[string]any, base string) bool {
		for _, b := range baseTypes(d) {
			if b == base {
				return true
			}
		}
		return false
	}
	if isBase(def, "object") {
		if _, hasProps := def["properties"]; hasProps {
			return def, true
		}
		return nil, false
	}
	if isBase(def, "array") {
		items, _ := def["items"].(map[string]any)
		if items == nil {
			return nil, false
		}
		if _, hasRef := items["$ref"]; hasRef {
			return nil, false
		}
		if isBase(items, "object") {
			if _, hasProps := items["properties"]; hasProps {
				return items, true
			}
		}
	}
	return nil, false
}

// checkNestedObjectCoverage is a blocking check: every inline object /
// array-item-object shape found under a mapped, reachable schema must have
// either a nested[] map entry or a recorded unmodeled.json nested_objects
// exclusion -- never silently skipped. A shape covered by a nested[] entry
// is itself walked one further level (ledger-only: the map format has no
// deeper "nested-of-nested" construct, so anything inline found there can
// only ever be excused, never mapped).
func checkNestedObjectCoverage(sm *SpecMap, um *Unmodeled, spec *specDoc) []string {
	reachable := reachableSchemas(spec)

	var issues []string
	var walk func(schemaKey, label string, schema map[string]any, exclude map[string]bool, nested map[string]NestedMapping)
	walk = func(schemaKey, label string, schema map[string]any, exclude map[string]bool, nested map[string]NestedMapping) {
		props := properties(schema)
		names := make([]string, 0, len(props))
		for n := range props {
			names = append(names, n)
		}
		sort.Strings(names)

		for _, name := range names {
			if exclude[name] {
				continue
			}
			def, _ := props[name].(map[string]any)
			shape, isInline := inlineObjectShape(def)
			if !isInline {
				continue
			}

			if n, isNested := nested[name]; isNested {
				nestedKey := name
				if len(n.SDK) > 0 {
					nestedKey = n.SDK[0].Symbol
				}
				walk(nestedKey, label+"."+name, shape, nil, nil)
				continue
			}

			if um.excusesNestedObject(schemaKey, name) {
				continue
			}
			issues = append(issues, fmt.Sprintf(
				"NEEDS_HUMAN: %s.%s: inline object shape has no nested[] map entry and no recorded unmodeled.json nested_objects exclusion (ledger schema %q property %q)",
				label, name, schemaKey, name))
		}
	}

	for _, t := range sm.Types {
		if !reachable[t.Spec] {
			continue
		}
		schema, _, ok := spec.schema(t.Spec)
		if !ok {
			continue // already reported by reconcileTypes
		}
		exclude := map[string]bool{}
		for _, d := range t.Discriminator {
			exclude[d] = true
		}
		walk(t.Spec, t.Spec, schema, exclude, t.Nested)
	}

	sort.Strings(issues)
	return issues
}
