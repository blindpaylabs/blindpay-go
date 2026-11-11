package blindpaytest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
)

// FixturesDir returns the absolute path of the shared spec-shaped fixtures
// (testdata/fixtures, copied verbatim from the cross-SDK reference set; see
// testdata/fixtures/README.md). The directory is part of the module (Go
// module zips include testdata/), so a missing directory fails the test.
func FixturesDir(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("blindpaytest: cannot locate fixtures directory")
	}
	dir := filepath.Join(filepath.Dir(file), "..", "..", "testdata", "fixtures")
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("blindpaytest: fixtures not available: %v", err)
	}
	return dir
}

// ReadFixture returns the raw bytes of a fixture file, relative to
// FixturesDir (e.g. "errors/rate_limited_429.json").
func ReadFixture(t *testing.T, rel string) []byte {
	t.Helper()

	b, err := os.ReadFile(filepath.Join(FixturesDir(t), filepath.FromSlash(rel)))
	if err != nil {
		t.Fatalf("blindpaytest: reading fixture %s: %v", rel, err)
	}
	return b
}

// Shape tells which branch of a response a fixture variant exercises.
type Shape int

const (
	// ShapeSingle is a non-list response (variant names without a suffix).
	ShapeSingle Shape = iota
	// ShapeEnvelope is the {data, pagination} list branch (".b0" variants).
	ShapeEnvelope
	// ShapeBareArray is the bare-array list branch (".b1" variants).
	ShapeBareArray
)

// Variant is one named response body from a response fixture.
type Variant struct {
	Name  string
	Body  json.RawMessage
	Shape Shape
}

func shapeOf(name string) Shape {
	switch {
	case strings.HasSuffix(name, ".b0"):
		return ShapeEnvelope
	case strings.HasSuffix(name, ".b1"):
		return ShapeBareArray
	default:
		return ShapeSingle
	}
}

// ResponseVariants returns every variant of the success response fixture for
// the operation key used in index.json (e.g.
// "GET /v1/instances/{instance_id}/payouts"), sorted by name. For list
// operations, names ending in ".b0" are the {data, pagination} envelope and
// names ending in ".b1" are the bare array.
func ResponseVariants(t *testing.T, operation string) []Variant {
	t.Helper()

	var index struct {
		Operations map[string]string `json:"operations"`
	}
	if err := json.Unmarshal(ReadFixture(t, "index.json"), &index); err != nil {
		t.Fatalf("blindpaytest: decoding index.json: %v", err)
	}
	rel, ok := index.Operations[operation]
	if !ok {
		t.Fatalf("blindpaytest: operation %q not in fixtures index", operation)
	}

	var fixture struct {
		Responses map[string]struct {
			Variants map[string]json.RawMessage `json:"variants"`
		} `json:"responses"`
	}
	if err := json.Unmarshal(ReadFixture(t, rel), &fixture); err != nil {
		t.Fatalf("blindpaytest: decoding %s: %v", rel, err)
	}

	var variants map[string]json.RawMessage
	for _, status := range []string{"200", "201"} {
		if r, ok := fixture.Responses[status]; ok {
			variants = r.Variants
			break
		}
	}
	if len(variants) == 0 {
		t.Fatalf("blindpaytest: %s has no 200/201 variants", rel)
	}

	out := make([]Variant, 0, len(variants))
	for name, body := range variants {
		out = append(out, Variant{Name: name, Body: body, Shape: shapeOf(name)})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// ErrorFixture is one file from testdata/fixtures/errors.
type ErrorFixture struct {
	Status  int
	Headers map[string]string
	// Body is the JSON error body; RawBody is set instead for non-JSON or
	// empty bodies.
	Body    json.RawMessage
	RawBody *string
	// SpecGeneratedExample is an extra spec-derived body (only on some files).
	SpecGeneratedExample json.RawMessage
}

// ResponseBody returns the bytes the server sends for this fixture.
func (f ErrorFixture) ResponseBody() []byte {
	if f.RawBody != nil {
		return []byte(*f.RawBody)
	}
	return f.Body
}

// LoadErrorFixture decodes testdata/fixtures/errors/<name>.
func LoadErrorFixture(t *testing.T, name string) ErrorFixture {
	t.Helper()

	// The wire struct is local so cmd/contractcheck (which scans named
	// types) does not mistake fixture metadata keys for API fields.
	var raw struct {
		Status               int               `json:"status"`
		Headers              map[string]string `json:"headers"`
		Body                 json.RawMessage   `json:"body"`
		RawBody              *string           `json:"raw_body"`
		SpecGeneratedExample json.RawMessage   `json:"spec_generated_example"`
	}
	if err := json.Unmarshal(ReadFixture(t, "errors/"+name), &raw); err != nil {
		t.Fatalf("blindpaytest: decoding errors/%s: %v", name, err)
	}
	return ErrorFixture(raw)
}

// RowCount returns the number of rows in a list variant, for either shape.
func (v Variant) RowCount(t *testing.T) int {
	t.Helper()

	var rows []json.RawMessage
	// Detect from the body: always-envelope lists (transfers) have no
	// .b0/.b1 suffix.
	if trimmed := strings.TrimSpace(string(v.Body)); strings.HasPrefix(trimmed, "{") {
		var envelope struct {
			Data []json.RawMessage `json:"data"`
		}
		if err := json.Unmarshal(v.Body, &envelope); err != nil {
			t.Fatalf("blindpaytest: variant %s: %v", v.Name, err)
		}
		rows = envelope.Data
	} else if err := json.Unmarshal(v.Body, &rows); err != nil {
		t.Fatalf("blindpaytest: variant %s: %v", v.Name, err)
	}
	return len(rows)
}

// FindVariant looks up one named variant (e.g. "full.b0") of an operation's
// success response fixture.
func FindVariant(t *testing.T, operation, name string) Variant {
	t.Helper()
	for _, v := range ResponseVariants(t, operation) {
		if v.Name == name {
			return v
		}
	}
	t.Fatalf("blindpaytest: variant %q not found for %s", name, operation)
	return Variant{}
}

// NextPage returns the envelope's pagination.next_page ("" when absent).
func (v Variant) NextPage(t *testing.T) string {
	t.Helper()
	var envelope struct {
		Pagination struct {
			NextPage *string `json:"next_page"`
		} `json:"pagination"`
	}
	if err := json.Unmarshal(v.Body, &envelope); err != nil {
		t.Fatalf("blindpaytest: variant %s: %v", v.Name, err)
	}
	if envelope.Pagination.NextPage == nil {
		return ""
	}
	return *envelope.Pagination.NextPage
}

// AsLastPage returns the envelope body with pagination rewritten to a last
// page (has_more false, next_page null), keeping the fixture's rows.
func (v Variant) AsLastPage(t *testing.T) json.RawMessage {
	t.Helper()
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(v.Body, &envelope); err != nil {
		t.Fatalf("blindpaytest: variant %s: %v", v.Name, err)
	}
	envelope["pagination"] = json.RawMessage(`{"has_more":false,"next_page":null,"prev_page":null}`)
	b, err := json.Marshal(envelope)
	if err != nil {
		t.Fatalf("blindpaytest: variant %s: %v", v.Name, err)
	}
	return b
}
