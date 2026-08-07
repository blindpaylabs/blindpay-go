// Command contractcheck is a mechanical, dependency-free wire-contract check.
//
// It parses .api-sync/spec-snapshot.json (a committed copy of the OpenAPI spec
// this SDK targets) and every non-test .go file in the repo, then:
//
//   - Direction A (hard failure): flags any json struct tag the SDK declares
//     that does not exist as an object property name anywhere in the spec
//     snapshot, unless the exact Package.Type.field is in the allow-list.
//   - Direction B (enum, hard failure): flags any WebhookEvent enum member
//     present in the spec's webhook events enum but missing from the SDK's
//     WebhookEvent Go consts.
//   - Direction B (field, warning only): reports spec properties the SDK
//     does not model at all. Informational, never fails the build.
//
// Run it with: go run ./cmd/contractcheck
package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	specPath      = ".api-sync/spec-snapshot.json"
	allowlistPath = ".api-sync/contract-allowlist.json"
)

// allowEntry is one row of the committed allow-list.
type allowEntry struct {
	Schema string `json:"schema"` // "package.TypeName"
	Field  string `json:"field"`  // wire (json) field name
	Reason string `json:"reason"`
	Owner  string `json:"owner"`
}

// wireField is one json-tagged struct field the SDK declares.
type wireField struct {
	schema string // "package.TypeName"
	field  string // json field name
	file   string
	line   int
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "contractcheck:", err)
		os.Exit(1)
	}
}

func run() error {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return err
	}

	specProps, webhookEnum, err := loadSpec(filepath.Join(repoRoot, specPath))
	if err != nil {
		return fmt.Errorf("loading spec snapshot: %w", err)
	}

	allowlist, err := loadAllowlist(filepath.Join(repoRoot, allowlistPath))
	if err != nil {
		return fmt.Errorf("loading allow-list: %w", err)
	}

	fields, sdkWebhookEvents, err := scanGoSources(repoRoot)
	if err != nil {
		return fmt.Errorf("scanning go sources: %w", err)
	}

	var hardFailures []string
	usedAllow := map[string]bool{}

	// Direction A.
	for _, f := range fields {
		if specProps[f.field] {
			continue
		}
		key := f.schema + "." + f.field
		if _, ok := allowlist[key]; ok {
			usedAllow[key] = true
			continue
		}
		hardFailures = append(hardFailures, fmt.Sprintf(
			"%s:%d: %s declares wire field %q which does not exist anywhere in %s",
			f.file, f.line, f.schema, f.field, specPath))
	}

	// Direction B, enum (hard failure).
	for _, ev := range webhookEnum {
		if !sdkWebhookEvents[ev] {
			hardFailures = append(hardFailures, fmt.Sprintf(
				"internal/types/types.go: WebhookEvent enum is missing spec member %q", ev))
		}
	}

	// Report unused allow-list entries (hygiene, non-fatal).
	var staleAllow []string
	for key := range allowlist {
		if !usedAllow[key] {
			staleAllow = append(staleAllow, key)
		}
	}
	sort.Strings(staleAllow)
	for _, key := range staleAllow {
		fmt.Printf("WARNING: allow-list entry %q did not match any current wire field (stale?)\n", key)
	}

	// Direction B, field (warning only).
	sdkFieldNames := map[string]bool{}
	for _, f := range fields {
		sdkFieldNames[f.field] = true
	}
	var unmodeled []string
	for name := range specProps {
		if !sdkFieldNames[name] {
			unmodeled = append(unmodeled, name)
		}
	}
	sort.Strings(unmodeled)
	fmt.Printf("INFO: %d spec properties are not modeled by any SDK struct (warning only, not a failure)\n", len(unmodeled))

	if len(hardFailures) > 0 {
		sort.Strings(hardFailures)
		fmt.Println("\nFAIL: contract-check found wire-contract mismatches:")
		for _, f := range hardFailures {
			fmt.Println("  -", f)
		}
		return fmt.Errorf("%d hard failure(s)", len(hardFailures))
	}

	fmt.Println("OK: contract-check passed")
	return nil
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found above %s", dir)
		}
		dir = parent
	}
}

// loadSpec returns the set of every object property name that appears
// anywhere in the spec document, plus the webhook endpoint events enum.
func loadSpec(path string) (map[string]bool, []string, error) {
	data, err := os.ReadFile(path) //#nosec G304 -- path is developer-supplied (CLI arg or repo-relative), this is a local codegen tool
	if err != nil {
		return nil, nil, err
	}

	var doc any
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, nil, err
	}

	props := map[string]bool{}
	collectProps(doc, props)
	collectParamNames(doc, props)

	webhookEnum := extractWebhookEventEnum(doc)
	if len(webhookEnum) == 0 {
		return nil, nil, fmt.Errorf("could not find webhook events enum in spec (WebhookEndpoint/WebhookEndpointIn.properties.events.items.enum)")
	}

	return props, webhookEnum, nil
}

func collectProps(node any, acc map[string]bool) {
	switch v := node.(type) {
	case map[string]any:
		if props, ok := v["properties"].(map[string]any); ok {
			for k := range props {
				acc[k] = true
			}
		}
		for _, child := range v {
			collectProps(child, acc)
		}
	case []any:
		for _, child := range v {
			collectProps(child, acc)
		}
	}
}

// collectParamNames adds OpenAPI "parameter" object names (query/path
// params, e.g. limit/offset/customer_id) to acc. Parameters are a distinct
// wire concept from schema properties but are just as real a contract.
func collectParamNames(node any, acc map[string]bool) {
	switch v := node.(type) {
	case map[string]any:
		if _, hasIn := v["in"]; hasIn {
			if name, ok := v["name"].(string); ok {
				acc[name] = true
			}
		}
		for _, child := range v {
			collectParamNames(child, acc)
		}
	case []any:
		for _, child := range v {
			collectParamNames(child, acc)
		}
	}
}

// extractWebhookEventEnum walks components.schemas.WebhookEndpoint (falling
// back to WebhookEndpointIn) for properties.events.items.enum.
func extractWebhookEventEnum(doc any) []string {
	root, ok := doc.(map[string]any)
	if !ok {
		return nil
	}
	components, _ := root["components"].(map[string]any)
	if components == nil {
		return nil
	}
	schemas, _ := components["schemas"].(map[string]any)
	if schemas == nil {
		return nil
	}

	for _, name := range []string{"WebhookEndpoint", "WebhookEndpointIn", "WebhookEndpointOut"} {
		schema, ok := schemas[name].(map[string]any)
		if !ok {
			continue
		}
		props, ok := schema["properties"].(map[string]any)
		if !ok {
			continue
		}
		events, ok := props["events"].(map[string]any)
		if !ok {
			continue
		}
		items, ok := events["items"].(map[string]any)
		if !ok {
			continue
		}
		enumRaw, ok := items["enum"].([]any)
		if !ok {
			continue
		}
		var out []string
		for _, e := range enumRaw {
			if s, ok := e.(string); ok {
				out = append(out, s)
			}
		}
		if len(out) > 0 {
			return out
		}
	}
	return nil
}

func loadAllowlist(path string) (map[string]allowEntry, error) {
	data, err := os.ReadFile(path) //#nosec G304 -- path is developer-supplied (CLI arg or repo-relative), this is a local codegen tool
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]allowEntry{}, nil
		}
		return nil, err
	}

	var entries []allowEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("%s: invalid JSON: %w", path, err)
	}

	out := map[string]allowEntry{}
	for i, e := range entries {
		if e.Schema == "" || e.Field == "" {
			return nil, fmt.Errorf("%s: entry %d missing schema or field", path, i)
		}
		if e.Reason == "" || e.Owner == "" {
			return nil, fmt.Errorf("%s: entry %d (%s.%s) missing reason or owner", path, i, e.Schema, e.Field)
		}
		key := e.Schema + "." + e.Field
		if _, dup := out[key]; dup {
			return nil, fmt.Errorf("%s: duplicate entry %s", path, key)
		}
		out[key] = e
	}
	return out, nil
}

// skipDirs are not scanned for wire-declaring structs: generated/vendored
// code, examples (consume the public API, don't declare wire shapes), and
// this checker's own package.
var skipDirs = map[string]bool{
	".git":      true,
	"examples":  true,
	"cmd":       true,
	".api-sync": true,
}

func scanGoSources(repoRoot string) ([]wireField, map[string]bool, error) {
	var fields []wireField
	webhookEvents := map[string]bool{}

	err := filepath.WalkDir(repoRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if skipDirs[d.Name()] && path != repoRoot {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		rel, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return err
		}

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", rel, err)
		}

		pkgName := file.Name.Name

		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				collectWebhookConsts(decl, webhookEvents)
				continue
			}
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}
				schemaName := pkgName + "." + typeSpec.Name.Name
				fields = append(fields, extractFields(structType, schemaName, rel, fset)...)
			}
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return fields, webhookEvents, nil
}

// extractFields recursively walks a struct (including nested anonymous
// structs) collecting json-tagged fields.
func extractFields(st *ast.StructType, schemaName, file string, fset *token.FileSet) []wireField {
	var out []wireField
	if st.Fields == nil {
		return out
	}
	for _, f := range st.Fields.List {
		if len(f.Names) == 0 {
			// Embedded field; if it's itself an inline struct, recurse.
			if nested, ok := f.Type.(*ast.StructType); ok {
				out = append(out, extractFields(nested, schemaName, file, fset)...)
			}
			continue
		}
		if nested, ok := f.Type.(*ast.StructType); ok {
			out = append(out, extractFields(nested, schemaName, file, fset)...)
		}
		if f.Tag == nil {
			continue
		}
		tagVal, err := strconv.Unquote(f.Tag.Value)
		if err != nil {
			continue
		}
		jsonName := jsonTagName(tagVal)
		if jsonName == "" || jsonName == "-" {
			continue
		}
		out = append(out, wireField{
			schema: schemaName,
			field:  jsonName,
			file:   file,
			line:   fset.Position(f.Pos()).Line,
		})
	}
	return out
}

// jsonTagName extracts the name portion of a `json:"name,omitempty"` tag.
func jsonTagName(tag string) string {
	const key = "json:\""
	idx := strings.Index(tag, key)
	if idx == -1 {
		return ""
	}
	rest := tag[idx+len(key):]
	end := strings.Index(rest, "\"")
	if end == -1 {
		return ""
	}
	value := rest[:end]
	if comma := strings.Index(value, ","); comma != -1 {
		value = value[:comma]
	}
	return value
}

// collectWebhookConsts picks up `const ( X WebhookEvent = "foo.bar" )` blocks.
func collectWebhookConsts(decl ast.Decl, acc map[string]bool) {
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok || genDecl.Tok != token.CONST {
		return
	}
	for _, spec := range genDecl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		ident, ok := valueSpec.Type.(*ast.Ident)
		if !ok || ident.Name != "WebhookEvent" {
			continue
		}
		for _, v := range valueSpec.Values {
			lit, ok := v.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				continue
			}
			val, err := strconv.Unquote(lit.Value)
			if err != nil {
				continue
			}
			acc[val] = true
		}
	}
}
