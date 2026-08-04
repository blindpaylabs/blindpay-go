package main

import (
	"fmt"
	"sort"
)

// auditTypes is a full, STATE-based survey of every mapped property's spec
// type against its declared Go field type -- independent of, and strictly
// broader than, -check's blocking type-compatibility gate (which only
// reports what is not already excused by unmodeled.json). This exists
// because the blocking gate answers "did anything regress", not "does the
// SDK's current type surface already disagree with the spec"; a pre-
// existing mismatch that was excused (or simply missed) at some point
// would otherwise never surface again. Every finding is printed here
// regardless of unmodeled.json, annotated when it is already recorded
// there, so this is always a complete picture -- never gated on whether
// someone remembered to triage it.
func auditTypes(repoRoot string, sm *SpecMap, um *Unmodeled, spec *specDoc) []string {
	var findings []string

	for _, t := range sm.Types {
		schema, _, ok := spec.schema(t.Spec)
		if !ok {
			continue
		}
		exclude := map[string]bool{}
		for _, d := range t.Discriminator {
			exclude[d] = true
		}
		for propName := range t.Nested {
			exclude[propName] = true
		}
		findings = append(findings, auditOneLevel(repoRoot, t.Spec, t.Spec, schema, t.SDK, exclude, um)...)

		for propName, nested := range t.Nested {
			sub := resolveProperty(schema, propName)
			if sub == nil {
				continue
			}
			unmodeledKey := propName
			if len(nested.SDK) > 0 {
				unmodeledKey = nested.SDK[0].Symbol
			}
			findings = append(findings, auditOneLevel(repoRoot, unmodeledKey, t.Spec+"."+propName, sub, nested.SDK, nil, um)...)
		}
	}

	sort.Strings(findings)
	return findings
}

func auditOneLevel(repoRoot, unmodeledSchemaKey, label string, schema map[string]any, sites []SDKSite, exclude map[string]bool, um *Unmodeled) []string {
	var shapes []*StructShape
	for _, site := range sites {
		shape, err := findStructShape(repoRoot, site.File, site.Symbol)
		if err != nil {
			continue // map-invalid; already reported by -validate-map
		}
		shapes = append(shapes, shape)
	}
	if len(shapes) == 0 {
		return nil
	}

	props := properties(schema)
	names := make([]string, 0, len(props))
	for name := range props {
		names = append(names, name)
	}
	sort.Strings(names)

	var out []string
	for _, name := range names {
		if exclude[name] {
			continue
		}
		def, _ := props[name].(map[string]any)
		specCat := specTypeCategory(def)

		for _, shape := range shapes {
			for _, f := range shape.fieldsNamed(name) {
				goCat := goTypeCategory(f.GoType)
				if typesCompatible(specCat, goCat) {
					continue
				}
				line := fmt.Sprintf("%s.%s: spec %s vs %s.%s Go type %q (%s)",
					label, name, specWireTypeName(def), shape.Symbol, f.GoName, f.GoType, shape.File)
				if um.excusesProperty(unmodeledSchemaKey, name) {
					line += "  [recorded in unmodeled.json]"
				} else {
					line += "  [NOT in unmodeled.json]"
				}
				out = append(out, line)
			}
		}
	}
	return out
}

// runAuditTypes prints the full type-compatibility survey and always
// exits 0: non-blocking by design, since triaging every pre-existing
// mismatch is its own (Phase C) piece of work, not a gate on every sync.
func runAuditTypes(repoRoot string) error {
	sm, err := loadSpecMap(repoRoot)
	if err != nil {
		return fmt.Errorf("loading spec-map.json: %w", err)
	}
	um, err := loadUnmodeled(repoRoot)
	if err != nil {
		return fmt.Errorf("loading unmodeled.json: %w", err)
	}
	_, raw, err := readSpecFile(repoRoot, ".api-sync/spec-snapshot.json")
	if err != nil {
		return err
	}

	findings := auditTypes(repoRoot, sm, um, loadSpecDoc(raw))
	fmt.Println("=== apisync type audit (non-blocking) ===")
	if len(findings) == 0 {
		fmt.Println("no spec-vs-Go-type mismatches found")
		return nil
	}
	for _, f := range findings {
		fmt.Println("  -", f)
	}
	return nil
}
