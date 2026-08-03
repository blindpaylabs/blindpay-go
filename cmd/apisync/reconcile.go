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
}

func (p *planResult) addNeedsHuman(format string, args ...any) {
	p.NeedsHuman = append(p.NeedsHuman, fmt.Sprintf(format, args...))
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
		for _, site := range t.SDK {
			if _, err := findStructShape(repoRoot, site.File, site.Symbol); err != nil {
				issues = append(issues, fmt.Sprintf("spec-map.json: type %s: %v", t.Spec, err))
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
// "schema added" NEEDS_HUMAN categories.
func checkUnclassifiedSchemas(sm *SpecMap, spec *specDoc) []string {
	classified := map[string]bool{}
	for _, t := range sm.Types {
		classified[t.Spec] = true
	}
	for _, i := range sm.Ignore.Schemas {
		classified[i.Schema] = true
	}

	var issues []string
	for name := range spec.schemas() {
		if !classified[name] {
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

		p.reconcileOneLevel(repoRoot, t.Spec, "", newSchema, oldSchema, oldFound, t.SDK, exclude, um)

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
			p.reconcileOneLevel(repoRoot, nestedSchemaKey, t.Spec+"."+propName, newSub, oldSub, oldSubFound, nested.SDK, nil, um)
		}
	}

	sort.Strings(p.NeedsHuman)
	return p
}

// reconcileOneLevel handles one flat property set (top-level schema or a
// resolved nested sub-object) against its SDK site(s).
func (p *planResult) reconcileOneLevel(repoRoot, unmodeledSchemaKey, humanLabelPrefix string, newSchema, oldSchema map[string]any, oldFound bool, sites []SDKSite, exclude map[string]bool, um *Unmodeled) {
	label := unmodeledSchemaKey
	if humanLabelPrefix != "" {
		label = humanLabelPrefix
	}

	shapes := make([]*StructShape, 0, len(sites))
	for _, site := range sites {
		shape, err := findStructShape(repoRoot, site.File, site.Symbol)
		if err != nil {
			continue // already reported by checkMapValidity
		}
		shapes = append(shapes, shape)
	}
	if len(shapes) == 0 {
		return
	}

	newProps := properties(newSchema)
	newReq := requiredSet(newSchema)
	var oldProps map[string]any
	var oldReq map[string]bool
	if oldFound {
		oldProps = properties(oldSchema)
		oldReq = requiredSet(oldSchema)
	}

	sdkHas := func(jsonName string) bool {
		for _, s := range shapes {
			if s.hasJSONField(jsonName) {
				return true
			}
		}
		return false
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

		if sdkHas(name) {
			continue
		}
		if um.excusesProperty(unmodeledSchemaKey, name) {
			continue
		}

		if newReq[name] {
			p.addNeedsHuman("NEEDS_HUMAN: %s.%s: required spec property is not modeled by any mapped SDK struct", label, name)
			continue
		}

		goType, ok := scalarGoType(newDef)
		if !ok {
			p.addNeedsHuman("NEEDS_HUMAN: %s.%s: new property has a non-scalar shape (object/array/union); needs a human-designed Go type", label, name)
			continue
		}

		target := shapes[0]
		goName := pascalCase(name)
		var fieldLine string
		if target.preferPointerStyle() {
			fieldLine = fmt.Sprintf("\t%s *%s `json:\"%s,omitempty\"`", goName, goType, name)
		} else {
			fieldLine = fmt.Sprintf("\t%s %s `json:\"%s,omitempty\"`", goName, goType, name)
		}

		p.Actions = append(p.Actions, action{
			description: fmt.Sprintf("field %s.%s: add %s (%s)", target.Symbol, goName, name, target.File),
			insert:      insertion{File: target.File, Line: target.LastFieldLine, Text: fieldLine},
			bump:        "patch",
		})
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
