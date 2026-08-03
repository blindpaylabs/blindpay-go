package main

import "strings"

// typeCategory is a coarse wire-shape bucket used to compare a spec
// property's declared JSON Schema type against the Go type of the field
// already modeling it, independent of pointer-vs-bare or nullable-vs-not
// (this repo uses both optionality styles interchangeably, see
// StructShape.preferPointerStyle).
type typeCategory int

const (
	catUnconstrained typeCategory = iota // spec has no "type" key at all
	catString
	catNumeric // JSON Schema "integer" and "number" are both a JSON number
	// token on the wire; encoding/json unmarshals either into a Go
	// int-family or float-family field without error (only a genuinely
	// fractional value into an int-family field ever fails, which this
	// schema-level check can't observe), so they are one category, not two.
	catBool
	catArray
	catObject
	catWildcard    // Go interface{}/any/map[string]T: accepts anything
	catNamedString // a bare named Go identifier: could be a string-backed
	// enum (every enum in this repo is `type X string`) or a named struct;
	// ambiguous without full cross-package type resolution, so treated as
	// compatible with string or object, incompatible with numeric/bool/array.
)

// specTypeCategory classifies a spec property definition's declared type,
// ignoring "null" (nullability is handled separately, and never itself a
// mismatch -- see baseTypes' doc comment on the same rationale).
func specTypeCategory(def map[string]any) typeCategory {
	bases := baseTypes(def)
	if len(bases) == 0 {
		return catUnconstrained
	}
	// Prefer the first non-null base type; a spec property should not
	// legitimately declare more than one non-null base type.
	switch bases[0] {
	case "string":
		return catString
	case "integer", "number":
		return catNumeric
	case "boolean":
		return catBool
	case "array":
		return catArray
	case "object":
		return catObject
	default:
		return catUnconstrained
	}
}

// goTypeCategory classifies a Go field's type expression text (as produced
// by typeExprString), stripping a leading pointer.
func goTypeCategory(goType string) typeCategory {
	t := strings.TrimPrefix(goType, "*")
	switch {
	case t == "string":
		return catString
	case t == "time.Time":
		return catString // wire shape is a JSON string (RFC3339)
	case t == "int" || t == "int8" || t == "int16" || t == "int32" || t == "int64" ||
		t == "uint" || t == "uint8" || t == "uint16" || t == "uint32" || t == "uint64" ||
		t == "float32" || t == "float64":
		return catNumeric
	case t == "bool":
		return catBool
	case strings.HasPrefix(t, "[]"):
		return catArray
	case t == "interface{}" || t == "any" || strings.HasPrefix(t, "map["):
		return catWildcard
	case strings.HasPrefix(t, "struct"):
		return catObject
	default:
		// A bare (possibly package-qualified) identifier: every enum in
		// this repo is `type X string`, but a named struct is equally
		// possible (e.g. a nested resource type); genuinely ambiguous
		// without loading and type-checking every package.
		return catNamedString
	}
}

// typesCompatible reports whether a spec-declared type and a Go field's
// type can plausibly represent the same wire value. False means a real,
// dangerous mismatch (e.g. spec integer against Go string: encoding/json
// will fail to unmarshal that at runtime for every consumer, silently,
// unless this is caught here at sync time).
func typesCompatible(spec, goCat typeCategory) bool {
	if spec == catUnconstrained || goCat == catWildcard {
		return true
	}
	if goCat == catNamedString {
		return spec == catString || spec == catObject
	}
	return spec == goCat
}
