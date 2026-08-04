package main

import (
	"fmt"
	"sort"
	"strings"
)

// action is one fully-resolved APPLICABLE change: exactly where to splice
// what text, plus the version-bump class it implies.
type action struct {
	description string
	insert      insertion
	bump        string // "minor" or "patch"
}

// planResult is everything a single reconciliation pass produced.
type planResult struct {
	Actions    []action
	NeedsHuman []string
	// PendingSchemas is set only by checkOperationChanges: the spec schema
	// names its operation-insert actions are about to register in
	// spec-map.json, so checkUnclassifiedSchemas can treat them as already
	// classified this same run (the registration hasn't been written to
	// disk yet -- see reconcileAndMaybeApply).
	PendingSchemas map[string]bool
}

func (p *planResult) addNeedsHuman(format string, args ...any) {
	p.NeedsHuman = append(p.NeedsHuman, fmt.Sprintf(format, args...))
}

// checkFieldTypeCompatibility verifies that a property already modeled by
// the SDK still corresponds to the CURRENT spec's declared type -- not a
// drift check against the old snapshot, a standing state assertion. This
// is the dangerous direction contractcheck-style presence checks miss:
// a spec property silently changing shape (e.g. string -> integer) behind
// a Go field of the wrong kind compiles and lints clean, then fails at
// runtime for every consumer the moment encoding/json tries to unmarshal a
// number into that string field. Nullability changes and an enum
// property degrading to a bare string of the same base type are
// deliberately NOT flagged (see typesCompatible / specTypeCategory).
func (p *planResult) checkFieldTypeCompatibility(shapes []*StructShape, label, name string, specDef map[string]any) {
	specCat := specTypeCategory(specDef)
	for _, shape := range shapes {
		for _, f := range shape.fieldsNamed(name) {
			goCat := goTypeCategory(f.GoType)
			if typesCompatible(specCat, goCat) {
				continue
			}
			p.addNeedsHuman(
				"NEEDS_HUMAN: %s.%s: spec type %s is not compatible with %s.%s's Go type %q (%s); encoding/json would fail to unmarshal this at runtime",
				label, name, specWireTypeName(specDef), shape.Symbol, f.GoName, f.GoType, shape.File)
		}
	}
}

// specWireTypeName renders a property definition's declared type for an
// error message (e.g. "integer", "[string, null]" widened to "string").
func specWireTypeName(def map[string]any) string {
	bases := baseTypes(def)
	if len(bases) == 0 {
		return "(unconstrained)"
	}
	return strings.Join(bases, "|")
}

// checkMapValidity resolves every SDK anchor in the map (struct sites, enum
// sites, re-export sites). A resolution failure is always NEEDS_HUMAN: a
// stale map must fail loudly, never silently skip a mapping.
func checkMapValidity(repoRoot string, sm *SpecMap) []string {
	var issues []string

	for _, e := range sm.Enums {
		if _, err := findEnumShape(repoRoot, e.SDK.File, e.SDK.Symbol); err != nil {
			issues = append(issues, fmt.Sprintf("spec-map.json: enum %s: %v", e.SDK.Symbol, err))
		}
		if e.Reexport != nil {
			if _, err := findReexportShape(repoRoot, e.Reexport.File, e.SDK.Symbol); err != nil {
				issues = append(issues, fmt.Sprintf("spec-map.json: enum %s reexport: %v", e.SDK.Symbol, err))
			}
		}
	}

	for _, t := range sm.Types {
		knownSymbols := map[string]bool{}
		for _, site := range t.SDK {
			knownSymbols[site.Symbol] = true
			if _, err := findStructShape(repoRoot, site.File, site.Symbol); err != nil {
				issues = append(issues, fmt.Sprintf("spec-map.json: type %s: %v", t.Spec, err))
			}
		}
		if t.Policy != "" && t.Policy != "uniform" && t.Policy != "union" {
			issues = append(issues, fmt.Sprintf("spec-map.json: type %s: unknown policy %q (must be \"uniform\" or \"union\")", t.Spec, t.Policy))
		}
		for propName, siteSymbols := range t.PropertySites {
			for _, sym := range siteSymbols {
				if !knownSymbols[sym] {
					issues = append(issues, fmt.Sprintf("spec-map.json: type %s property_sites[%s]: %q is not one of this mapping's sdk sites", t.Spec, propName, sym))
				}
			}
		}
		for propName, nested := range t.Nested {
			for _, site := range nested.SDK {
				if _, err := findStructShape(repoRoot, site.File, site.Symbol); err != nil {
					issues = append(issues, fmt.Sprintf("spec-map.json: type %s nested %s: %v", t.Spec, propName, err))
				}
			}
		}
	}

	sort.Strings(issues)
	return issues
}

// checkUnclassifiedSchemas flags any component schema in the target spec
// that is neither mapped nor explicitly ignored -- a genuinely new spec
// construct that needs a human triage decision (map it or ignore it),
// per the design's "property added on a schema with no mapping" /
// "schema added" NEEDS_HUMAN categories. An unreachable orphan schema (not
// reachable from any path, webhook, or non-schema component section, even
// transitively through other schemas' $refs) produces no work at all: see
// reachableSchemas.
func checkUnclassifiedSchemas(sm *SpecMap, spec *specDoc, pending map[string]bool) []string {
	classified := map[string]bool{}
	for _, t := range sm.Types {
		classified[t.Spec] = true
	}
	for _, i := range sm.Ignore.Schemas {
		classified[i.Schema] = true
	}

	reachable := reachableSchemas(spec)

	var issues []string
	for name := range spec.schemas() {
		if !reachable[name] {
			continue // orphan: nothing in the spec's surface can reach it, no work either way
		}
		if !classified[name] && !pending[name] {
			issues = append(issues, fmt.Sprintf(
				"NEEDS_HUMAN: schema %q is not in spec-map.json (neither types[] nor ignore.schemas[]); "+
					"a human must decide whether to model it or ignore it with a reason", name))
		}
	}
	sort.Strings(issues)
	return issues
}

// reconcileEnums checks every mapped enum's spec members against the SDK
// const block, honoring unmodeled.json exclusions, and against removals
// relative to the old (snapshot) spec.
func reconcileEnums(repoRoot string, sm *SpecMap, um *Unmodeled, oldSpec, newSpec *specDoc) *planResult {
	p := &planResult{}

	for _, e := range sm.Enums {
		if strings.HasPrefix(e.Spec.Schema, "inline:") {
			// Not a component schema -- an inline (non-$ref'd) operation
			// body, informational only. See the entry's "note" in
			// spec-map.json for why; not resolvable or enforceable here.
			continue
		}

		enumShape, err := findEnumShape(repoRoot, e.SDK.File, e.SDK.Symbol)
		if err != nil {
			// Already reported by checkMapValidity; skip further work here.
			continue
		}

		newSchema, _, ok := newSpec.schema(e.Spec.Schema, e.FallbackSpecSchemas...)
		if !ok {
			p.addNeedsHuman("NEEDS_HUMAN: enum %s: spec schema %q not found in the target spec (removed or renamed?)", e.SDK.Symbol, e.Spec.Schema)
			continue
		}
		newMembers, err := enumMembers(newSchema, e.Spec)
		if err != nil {
			p.addNeedsHuman("NEEDS_HUMAN: enum %s: %v", e.SDK.Symbol, err)
			continue
		}

		// Removal detection against the old snapshot, scoped to this mapping.
		if oldSchema, _, ok := oldSpec.schema(e.Spec.Schema, e.FallbackSpecSchemas...); ok {
			if oldMembers, err := enumMembers(oldSchema, e.Spec); err == nil {
				newSet := map[string]bool{}
				for _, m := range newMembers {
					newSet[m] = true
				}
				for _, m := range oldMembers {
					if !newSet[m] {
						p.addNeedsHuman("NEEDS_HUMAN: enum %s: spec member %q present in the committed snapshot is absent from the target spec (removal is always a hard fail)", e.SDK.Symbol, m)
					}
				}
			}
		}

		var reexportShape *EnumShape
		if e.Reexport != nil {
			reexportShape, err = findReexportShape(repoRoot, e.Reexport.File, e.SDK.Symbol)
			if err != nil {
				continue // already reported by checkMapValidity
			}
		}

		sort.Strings(newMembers)
		for _, member := range newMembers {
			if enumShape.hasMember(member) {
				continue
			}
			if um.excusesEnumMember(e.SDK.Symbol, member) {
				continue
			}
			goName := e.SDK.Symbol + pascalCase(member)
			line := fmt.Sprintf("\t%s %s = %q", goName, e.SDK.Symbol, member)
			act := action{
				description: fmt.Sprintf("enum %s: add member %s = %q (%s)", e.SDK.Symbol, goName, member, e.SDK.File),
				insert:      insertion{File: e.SDK.File, Line: enumShape.LastMemberLine, Text: line},
				bump:        "minor",
			}
			p.Actions = append(p.Actions, act)

			if reexportShape != nil {
				reexportLine := fmt.Sprintf("\t%s = types.%s", goName, goName)
				p.Actions = append(p.Actions, action{
					description: fmt.Sprintf("types.go: re-export %s = types.%s", goName, goName),
					insert:      insertion{File: e.Reexport.File, Line: reexportShape.LastMemberLine, Text: reexportLine},
					bump:        "minor",
				})
			}
		}
	}

	sort.Strings(p.NeedsHuman)
	return p
}

// scalarGoType maps a JSON Schema base type to a Go scalar type. ok is
// false for anything beyond a single simple scalar (object/array/multiple
// base types), which this patcher refuses to guess a shape for.
func scalarGoType(def map[string]any) (goType string, ok bool) {
	bases := baseTypes(def)
	if len(bases) != 1 {
		return "", false
	}
	switch bases[0] {
	case "string":
		return "string", true
	case "boolean":
		return "bool", true
	case "integer":
		return "int", true
	case "number":
		return "float64", true
	default:
		return "", false
	}
}

// reconcileTypes checks every mapped schema's (and nested sub-object's)
// properties against the union of its SDK struct sites, honoring
// unmodeled.json exclusions, discriminator-injected properties, removals,
// and required-ness/type changes relative to the old (snapshot) spec.
func reconcileTypes(repoRoot string, sm *SpecMap, um *Unmodeled, oldSpec, newSpec *specDoc) *planResult {
	p := &planResult{}

	for _, t := range sm.Types {
		newSchema, _, ok := newSpec.schema(t.Spec)
		if !ok {
			p.addNeedsHuman("NEEDS_HUMAN: type %s: spec schema not found in the target spec (removed or renamed?)", t.Spec)
			continue
		}
		oldSchema, _, oldFound := oldSpec.schema(t.Spec)

		exclude := map[string]bool{}
		for _, d := range t.Discriminator {
			exclude[d] = true
		}
		for propName := range t.Nested {
			exclude[propName] = true
		}

		p.reconcileOneLevel(repoRoot, t.Spec, "", newSchema, oldSchema, oldFound, t.SDK, t.Policy, t.PropertySites, exclude, um)

		for propName, nested := range t.Nested {
			newSub := resolveProperty(newSchema, propName)
			if newSub == nil {
				p.addNeedsHuman("NEEDS_HUMAN: type %s: nested property %q not found in the target spec", t.Spec, propName)
				continue
			}
			var oldSub map[string]any
			oldSubFound := false
			if oldFound {
				oldSub = resolveProperty(oldSchema, propName)
				oldSubFound = oldSub != nil
			}
			// Nested exclusions are keyed by the nested SDK symbol (there is
			// exactly one nested SDK site expected per current usage).
			nestedSchemaKey := propName
			if len(nested.SDK) > 0 {
				nestedSchemaKey = nested.SDK[0].Symbol
			}
			p.reconcileOneLevel(repoRoot, nestedSchemaKey, t.Spec+"."+propName, newSub, oldSub, oldSubFound, nested.SDK, "", nil, nil, um)
		}
	}

	sort.Strings(p.NeedsHuman)
	return p
}

// reconcileOneLevel handles one flat property set (top-level schema or a
// resolved nested sub-object) against its SDK site(s).
//
// Every mapped site is reconciled independently (never "does ANY site have
// this property" / "add it to site zero"): for a single-site mapping there is
// only one target anyway; for a multi-site mapping, policy decides the
// target sites per property -- "uniform" (default) means every site, "union"
// means exactly the sites named in propertySites, which must be explicit
// (never inferred from current site membership).
func (p *planResult) reconcileOneLevel(repoRoot, unmodeledSchemaKey, humanLabelPrefix string, newSchema, oldSchema map[string]any, oldFound bool, sites []SDKSite, policy string, propertySites map[string][]string, exclude map[string]bool, um *Unmodeled) {
	label := unmodeledSchemaKey
	if humanLabelPrefix != "" {
		label = humanLabelPrefix
	}

	shapeBySymbol := map[string]*StructShape{}
	var shapes []*StructShape
	for _, site := range sites {
		shape, err := findStructShape(repoRoot, site.File, site.Symbol)
		if err != nil {
			continue // already reported by checkMapValidity
		}
		shapeBySymbol[site.Symbol] = shape
		shapes = append(shapes, shape)
	}
	if len(shapes) == 0 {
		return
	}
	multiSite := len(sites) > 1
	tm := TypeMapping{SDK: sites, Policy: policy, PropertySites: propertySites}

	newProps := properties(newSchema)
	newReq := requiredSet(newSchema)
	var oldProps map[string]any
	var oldReq map[string]bool
	if oldFound {
		oldProps = properties(oldSchema)
		oldReq = requiredSet(oldSchema)
	}

	names := make([]string, 0, len(newProps))
	for name := range newProps {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		if exclude[name] {
			continue
		}
		newDef, _ := newProps[name].(map[string]any)

		if oldFound {
			if oldDef, existedBefore := oldProps[name].(map[string]any); existedBefore {
				if oldReq[name] != newReq[name] {
					p.addNeedsHuman("NEEDS_HUMAN: %s.%s: required-ness changed (was required=%v, now required=%v)", label, name, oldReq[name], newReq[name])
				}
				oldBases, newBases := baseTypes(oldDef), baseTypes(newDef)
				if len(oldBases) > 0 && len(newBases) > 0 {
					sort.Strings(oldBases)
					sort.Strings(newBases)
					if !sameStringSet(oldBases, newBases) {
						p.addNeedsHuman("NEEDS_HUMAN: %s.%s: type changed (was %v, now %v)", label, name, oldBases, newBases)
					}
				}
			}
		}

		var targetSymbols []string
		if multiSite {
			var ok bool
			targetSymbols, ok = tm.sitesFor(name)
			if !ok {
				if !um.excusesProperty(unmodeledSchemaKey, name) {
					p.addNeedsHuman("NEEDS_HUMAN: %s.%s: union-policy mapping has no property_sites entry for this property; a human must record which SDK sites it belongs on (or exclude it in unmodeled.json)", label, name)
				}
				continue
			}
		} else {
			targetSymbols = []string{shapes[0].Symbol}
		}

		excused := um.excusesProperty(unmodeledSchemaKey, name)
		goType, scalarOK := scalarGoType(newDef)
		missingRequiredReported := false
		nonScalarReported := false

		for _, sym := range targetSymbols {
			shape, ok := shapeBySymbol[sym]
			if !ok {
				continue // property_sites names a site not among this mapping's resolved shapes
			}
			if shape.hasJSONField(name) {
				if !excused {
					p.checkFieldTypeCompatibility([]*StructShape{shape}, label, name, newDef)
				}
				continue
			}
			if excused {
				continue
			}
			if newReq[name] {
				if !missingRequiredReported {
					p.addNeedsHuman("NEEDS_HUMAN: %s.%s: required spec property is not modeled by every mapped SDK struct (missing on %s)", label, name, shape.Symbol)
					missingRequiredReported = true
				}
				continue
			}
			if !scalarOK {
				if !nonScalarReported {
					p.addNeedsHuman("NEEDS_HUMAN: %s.%s: new property has a non-scalar shape (object/array/union); needs a human-designed Go type", label, name)
					nonScalarReported = true
				}
				continue
			}

			goName := pascalCase(name)
			var fieldLine string
			if shape.preferPointerStyle() {
				fieldLine = fmt.Sprintf("\t%s *%s `json:\"%s,omitempty\"`", goName, goType, name)
			} else {
				fieldLine = fmt.Sprintf("\t%s %s `json:\"%s,omitempty\"`", goName, goType, name)
			}

			p.Actions = append(p.Actions, action{
				description: fmt.Sprintf("field %s.%s: add %s (%s)", shape.Symbol, goName, name, shape.File),
				insert:      insertion{File: shape.File, Line: shape.LastFieldLine, Text: fieldLine},
				bump:        "patch",
			})
		}
	}

	// Property removal relative to the old snapshot, scoped to this mapping.
	if oldFound {
		for name := range oldProps {
			if exclude[name] {
				continue
			}
			if _, stillPresent := newProps[name]; !stillPresent {
				p.addNeedsHuman("NEEDS_HUMAN: %s.%s: property present in the committed snapshot is absent from the target spec (removal is always a hard fail)", label, name)
			}
		}
	}
}
