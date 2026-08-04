package main

import (
	"sort"
	"strings"
)

// operationEntry is one method+path operation from a spec's paths object,
// carrying just enough to decide whether it belongs to an already-ignored
// subsystem.
type operationEntry struct {
	Method string
	Path   string
	Tags   []string
	Op     map[string]any
}

func (e operationEntry) key() string { return e.Method + " " + e.Path }

// operationEntries returns every operation in a spec, keyed by "METHOD /path".
func operationEntries(spec *specDoc) map[string]operationEntry {
	out := map[string]operationEntry{}
	paths, _ := spec.raw["paths"].(map[string]any)
	for p, methodsRaw := range paths {
		methods, ok := methodsRaw.(map[string]any)
		if !ok {
			continue
		}
		for m, opRaw := range methods {
			switch m {
			case "get", "post", "put", "patch", "delete":
			default:
				continue
			}
			op, _ := opRaw.(map[string]any)
			var tags []string
			if rawTags, ok := op["tags"].([]any); ok {
				for _, t := range rawTags {
					if s, ok := t.(string); ok {
						tags = append(tags, s)
					}
				}
			}
			method := strings.ToUpper(m)
			entry := operationEntry{Method: method, Path: p, Tags: tags, Op: op}
			out[entry.key()] = entry
		}
	}
	return out
}

// collectRefSchemaNames recursively walks a JSON subtree collecting every
// "#/components/schemas/X" $ref target it finds, however deeply nested
// (anyOf/oneOf branches, array items, etc.).
func collectRefSchemaNames(node any) []string {
	var out []string
	const prefix = "#/components/schemas/"
	var walk func(n any)
	walk = func(n any) {
		switch v := n.(type) {
		case map[string]any:
			if ref, ok := v["$ref"].(string); ok && strings.HasPrefix(ref, prefix) {
				out = append(out, strings.TrimPrefix(ref, prefix))
			}
			for _, child := range v {
				walk(child)
			}
		case []any:
			for _, child := range v {
				walk(child)
			}
		}
	}
	walk(node)
	return out
}

// referencedSchemas returns every component schema name reachable from an
// operation's requestBody and 2xx responses.
func referencedSchemas(op map[string]any) []string {
	var names []string
	if rb, ok := op["requestBody"].(map[string]any); ok {
		names = append(names, collectRefSchemaNames(rb)...)
	}
	if responses, ok := op["responses"].(map[string]any); ok {
		for code, r := range responses {
			if !strings.HasPrefix(code, "2") {
				continue
			}
			names = append(names, collectRefSchemaNames(r)...)
		}
	}
	return names
}

// operationIsIgnored decides whether a new operation belongs to a subsystem
// this SDK already deliberately does not model, so a spec addition there
// must not wedge the whole pipeline. Two independent signals, either
// sufficient on its own:
//
//  1. Schema reachability (primary, mechanical): every component schema the
//     operation's requestBody/responses reference is already in
//     ignore.schemas, and there is at least one such schema (an operation
//     with zero named schema refs is never exempted this way -- an
//     untyped/generic body is not evidence of anything).
//  2. Tag match (secondary, for operations whose bodies are too generic to
//     carry a useful schema ref, e.g. `{"type":"object","additionalProperties":{}}`):
//     an exact, case-insensitive match between one of the operation's
//     OpenAPI tags and an ignored schema name. Exact match only, never a
//     prefix: e.g. UploadIn/UploadOut are real, mapped schemas alongside
//     the ignored UploadAnalyzeIn/UploadAnalyzeOut, so a tag "Upload"
//     prefix-matching "UploadAnalyzeIn" would have wrongly exempted an
//     operation actually belonging to the mapped Upload family.
func operationIsIgnored(op operationEntry, sm *SpecMap) bool {
	ignored := map[string]bool{}
	for _, e := range sm.Ignore.Schemas {
		ignored[e.Schema] = true
	}

	if refs := referencedSchemas(op.Op); len(refs) > 0 {
		allIgnored := true
		for _, r := range refs {
			if !ignored[r] {
				allIgnored = false
				break
			}
		}
		if allIgnored {
			return true
		}
	}

	for _, tag := range op.Tags {
		for name := range ignored {
			if strings.EqualFold(tag, name) {
				return true
			}
		}
	}

	return false
}

// reachableSchemas computes every component schema transitively reachable
// from the spec's actual surface: every path operation, every webhook, and
// every non-schema components section (parameters, requestBodies,
// responses, headers, etc. -- anything that can itself hold a $ref). A
// schema that is not reachable is an orphan: nothing in the spec's surface
// can ever produce or consume it, so a patcher has no work to do for it
// either way (map it or ignore it), and must not report it as an
// unclassified schema needing a human decision.
func reachableSchemas(spec *specDoc) map[string]bool {
	reachable := map[string]bool{}
	var queue []string

	if paths, ok := spec.raw["paths"].(map[string]any); ok {
		queue = append(queue, collectRefSchemaNames(paths)...)
	}
	if webhooks, ok := spec.raw["webhooks"].(map[string]any); ok {
		queue = append(queue, collectRefSchemaNames(webhooks)...)
	}
	if components, ok := spec.raw["components"].(map[string]any); ok {
		for section, val := range components {
			if section == "schemas" {
				continue
			}
			queue = append(queue, collectRefSchemaNames(val)...)
		}
	}

	schemas := spec.schemas()
	for len(queue) > 0 {
		name := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		if reachable[name] {
			continue
		}
		reachable[name] = true
		if def, ok := schemas[name]; ok {
			queue = append(queue, collectRefSchemaNames(def)...)
		}
	}
	return reachable
}

// checkOperationChanges reconciles every operation added or removed
// relative to the committed snapshot. Removal is always a hard fail.
// Addition, unless it belongs to an already-ignored subsystem (see
// operationIsIgnored), is handed to classifyAndPlanOperation: STANDARD
// (fully generated method + structs + any spec-map.json registrations,
// via actions) or NON-STANDARD (a precise NEEDS_HUMAN reason -- the
// generator refuses to guess naming, grouping, or a shape it cannot
// express, same posture as every other NEEDS_HUMAN in this patcher).
func checkOperationChanges(repoRoot string, sm *SpecMap, oldSpec, newSpec *specDoc) *planResult {
	oldOps, newOps := operationEntries(oldSpec), operationEntries(newSpec)
	p := &planResult{PendingSchemas: map[string]bool{}}

	newKeys := make([]string, 0, len(newOps))
	for key := range newOps {
		newKeys = append(newKeys, key)
	}
	sort.Strings(newKeys)

	for _, key := range newKeys {
		entry := newOps[key]
		if _, existed := oldOps[key]; existed {
			continue
		}
		if operationIsIgnored(entry, sm) {
			continue
		}

		acts, pending, reason := classifyAndPlanOperation(repoRoot, sm, newSpec, entry)
		if reason != "" {
			p.addNeedsHuman(
				"NEEDS_HUMAN: new operation %s is not present in the committed snapshot and the deterministic generator could not classify it as STANDARD (%s); adding SDK support needs hand-written client wiring (naming, grouping, request/response shape) that a patcher must not invent (tags=%v)",
				key, reason, entry.Tags)
			continue
		}
		p.Actions = append(p.Actions, acts...)
		for _, name := range pending {
			p.PendingSchemas[name] = true
		}
	}

	oldKeys := make([]string, 0, len(oldOps))
	for key := range oldOps {
		oldKeys = append(oldKeys, key)
	}
	sort.Strings(oldKeys)
	for _, key := range oldKeys {
		if _, stillPresent := newOps[key]; !stillPresent {
			p.addNeedsHuman(
				"NEEDS_HUMAN: operation %s present in the committed snapshot is absent from the target spec (removal is always a hard fail)", key)
		}
	}

	sort.Strings(p.NeedsHuman)
	return p
}
