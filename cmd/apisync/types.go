package main

// SDKSite is a single {file, symbol} anchor into this repo's Go source.
type SDKSite struct {
	File   string `json:"file"`
	Symbol string `json:"symbol"`
}

// SpecPointer anchors a spec-current.json construct: a component schema,
// optionally a property on it (possibly dotted, e.g. "tracking_payment.step"
// to reach a nested inline object), optionally its array items (for enums
// declared as items.enum on an array property, e.g. WebhookEndpoint.events).
type SpecPointer struct {
	Schema   string `json:"schema"`
	Property string `json:"property,omitempty"`
	Items    bool   `json:"items,omitempty"`
}

// EnumMapping ties one SDK const-block enum to the spec construct it mirrors.
type EnumMapping struct {
	SDK                 SDKSite     `json:"sdk"`
	Reexport            *SDKSite    `json:"reexport"`
	Spec                SpecPointer `json:"spec"`
	FallbackSpecSchemas []string    `json:"fallback_spec_schemas,omitempty"`
	Note                string      `json:"note,omitempty"`
}

// NestedMapping is a types[].nested entry: an inline sub-object (never
// $ref'd) reconciled against its own SDK site(s) independent of the parent.
type NestedMapping struct {
	SDK []SDKSite `json:"sdk"`
}

// TypeMapping ties one spec schema to one or more SDK struct sites.
//
// Policy governs how a property is checked/applied across multiple SDK
// sites (irrelevant, and left "", when len(SDK) == 1):
//
//   - "uniform": every site must contain the property; a missing optional
//     property is added to every site independently (never just the first).
//   - "union": a property may exist on a named subset of sites only, listed
//     explicitly in PropertySites. This must never be inferred from which
//     sites currently happen to have the property -- a property absent from
//     PropertySites is NEEDS_HUMAN, not silently treated as "satisfied
//     anywhere" or "add to site zero".
type TypeMapping struct {
	Spec          string                   `json:"spec"`
	SDK           []SDKSite                `json:"sdk"`
	Policy        string                   `json:"policy,omitempty"`
	PropertySites map[string][]string      `json:"property_sites,omitempty"`
	Nested        map[string]NestedMapping `json:"nested,omitempty"`
	Discriminator []string                 `json:"discriminator,omitempty"`
	Note          string                   `json:"note,omitempty"`
}

// sitesFor returns the SDK sites a property applies to, per the mapping's
// fan-out policy. ok is false when a union-policy mapping has no explicit
// PropertySites entry for this property -- the caller must treat that as
// NEEDS_HUMAN, never guess.
func (t TypeMapping) sitesFor(property string) (symbols []string, ok bool) {
	if t.Policy != "union" {
		symbols = make([]string, 0, len(t.SDK))
		for _, s := range t.SDK {
			symbols = append(symbols, s.Symbol)
		}
		return symbols, true
	}
	symbols, ok = t.PropertySites[property]
	return symbols, ok
}

// IgnoreEntry documents a spec schema this SDK deliberately does not model.
type IgnoreEntry struct {
	Schema string `json:"schema"`
	Reason string `json:"reason"`
}

// SpecMap is the parsed .api-sync/spec-map.json.
type SpecMap struct {
	Enums  []EnumMapping `json:"enums"`
	Types  []TypeMapping `json:"types"`
	Ignore struct {
		Schemas []IgnoreEntry `json:"schemas"`
	} `json:"ignore"`
}

// PropertyExclusion is one recorded, reasoned, owned property gap that keeps
// -check green. Schema is normally a spec component schema name; for a
// property inside a types[].nested sub-object, schema is that nested
// mapping's SDK symbol name instead (see unmodeled.json's $schema_note).
type PropertyExclusion struct {
	Schema   string `json:"schema"`
	Property string `json:"property"`
	Reason   string `json:"reason"`
	Owner    string `json:"owner"`
}

// EnumExclusion is one recorded, reasoned, owned enum-member gap.
type EnumExclusion struct {
	Symbol string `json:"symbol"`
	Member string `json:"member"`
	Reason string `json:"reason"`
	Owner  string `json:"owner"`
}

// Unmodeled is the parsed .api-sync/unmodeled.json.
type Unmodeled struct {
	Properties []PropertyExclusion `json:"properties"`
	Enums      []EnumExclusion     `json:"enums"`
	// EnumCoverage records an enum-constrained property (on a mapped schema,
	// possibly a dotted nested path e.g. "tracking_payment.step") that is
	// deliberately not tied to any mapped enum symbol -- see
	// checkEnumCoverage. Schema is always the top-level mapped spec schema
	// name (t.Spec), never a nested SDK symbol, matching EnumMapping.Spec's
	// own dotted-property convention.
	EnumCoverage []PropertyExclusion `json:"enum_coverage,omitempty"`
	// NestedObjects records an inline object / array-item-object shape (on a
	// mapped, reachable schema) that is deliberately not given a nested[]
	// map entry -- see checkNestedObjectCoverage. Schema follows the same
	// convention as Properties: the top-level spec schema name, or (for a
	// shape found one level inside an already-mapped nested[] entry) that
	// nested mapping's SDK symbol name.
	NestedObjects []PropertyExclusion `json:"nested_objects,omitempty"`
}

func (u *Unmodeled) excusesProperty(schema, property string) bool {
	for _, p := range u.Properties {
		if p.Schema == schema && p.Property == property {
			return true
		}
	}
	return false
}

func (u *Unmodeled) excusesEnumMember(symbol, member string) bool {
	for _, e := range u.Enums {
		if e.Symbol == symbol && e.Member == member {
			return true
		}
	}
	return false
}

func (u *Unmodeled) excusesEnumCoverage(schema, property string) bool {
	for _, p := range u.EnumCoverage {
		if p.Schema == schema && p.Property == property {
			return true
		}
	}
	return false
}

func (u *Unmodeled) excusesNestedObject(schema, property string) bool {
	for _, p := range u.NestedObjects {
		if p.Schema == schema && p.Property == property {
			return true
		}
	}
	return false
}
