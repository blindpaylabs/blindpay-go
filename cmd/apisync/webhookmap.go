package main

import (
	"fmt"
	"sort"
)

// WebhookMapping ties a spec webhook payload schema (*WebhookOut) to the Go
// structs whose json fields together make up the typed payload in
// webhook_events.go: the wrapper (which embeds WebhookHeader, i.e.
// webhook_event) and the resource struct it embeds. apisync's struct parser
// does not follow embedded fields, so every struct is listed explicitly.
type WebhookMapping struct {
	Spec string    `json:"spec"`
	SDK  []SDKSite `json:"sdk"`
	Note string    `json:"note,omitempty"`
}

// checkWebhookCoverage makes webhook-only drift visible: every top-level
// property of a mapped webhook schema (except the webhook_event header)
// must be a json field on one of its SDK structs, or be excused in
// unmodeled.json properties[] under the webhook schema's name. Types, enum
// members and nested shapes are reconciled through the resource schema's
// types[] mapping; this check is about properties that exist only on the
// webhook payload (e.g. a field the API adds to a webhook before the REST
// response). A missing property is NEEDS_HUMAN: a human decides whether it
// belongs on the shared resource struct or on the wrapper. A field that a
// pending types[] action adds to one of the sites (the usual case: the
// property is new on the REST schema and the webhook at once) counts as
// modeled, so -apply still works unattended.
func checkWebhookCoverage(repoRoot string, sm *SpecMap, um *Unmodeled, spec *specDoc, pending []action) []string {
	planned := map[string]bool{}
	for _, a := range pending {
		if f := a.addsField; f != nil {
			planned[f.File+"\x00"+f.Symbol+"\x00"+f.JSONName] = true
		}
	}
	var issues []string
	for _, w := range sm.Webhooks {
		schema, _, ok := spec.schema(w.Spec)
		if !ok {
			issues = append(issues, fmt.Sprintf("NEEDS_HUMAN: webhook %s: spec schema not found in the target spec (removed or renamed?)", w.Spec))
			continue
		}
		modeled := map[string]bool{"webhook_event": true}
		for _, site := range w.SDK {
			shape, err := findStructShape(repoRoot, site.File, site.Symbol)
			if err != nil {
				continue // reported by checkMapValidity
			}
			for _, f := range shape.Fields {
				modeled[f.JSONName] = true
			}
		}
		for name := range properties(schema) {
			if modeled[name] || um.excusesProperty(w.Spec, name) || plannedOnAnySite(planned, w.SDK, name) {
				continue
			}
			issues = append(issues, fmt.Sprintf("NEEDS_HUMAN: webhook %s.%s: property is not a json field on any of %s; add it (usually to the embedded resource struct) or excuse it in unmodeled.json", w.Spec, name, siteSymbols(w.SDK)))
		}
	}
	sort.Strings(issues)
	return issues
}

func siteSymbols(sites []SDKSite) string {
	out := ""
	for i, s := range sites {
		if i > 0 {
			out += ", "
		}
		out += s.Symbol
	}
	return out
}

func plannedOnAnySite(planned map[string]bool, sites []SDKSite, jsonName string) bool {
	for _, s := range sites {
		if planned[s.File+"\x00"+s.Symbol+"\x00"+jsonName] {
			return true
		}
	}
	return false
}
