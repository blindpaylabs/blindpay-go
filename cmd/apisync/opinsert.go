package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// pkgTarget is one resolved landing spot for a brand-new operation: the
// existing SDK package directory it belongs to (found mechanically, by
// path-literal prefix, never a hand-maintained tag/dir alias table) plus
// the path segments left over after that package's own resource root.
type pkgTarget struct {
	Dir      string   // package directory, e.g. "partnerfees"
	Package  string   // Go package name declared in that directory
	TailSegs []string // original ({param}-form) path segments after the matched root
}

// sdkPathLiteralsByDir is pathcoverage.go's sdkPathLiterals, bucketed by
// top-level package directory instead of flattened -- the signal
// resolveTarget needs to decide *which* existing package a brand-new
// operation's path belongs to.
func sdkPathLiteralsByDir(repoRoot string) (map[string]map[string]bool, error) {
	out := map[string]map[string]bool{}
	err := filepath.WalkDir(repoRoot, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(repoRoot, p)
		if relErr != nil {
			return relErr
		}
		top := strings.SplitN(rel, string(filepath.Separator), 2)[0]
		if d.IsDir() {
			if rel != "." && skipDirsCoverage[top] {
				return filepath.SkipDir
			}
			return nil
		}
		if top == rel {
			return nil // a file directly at repo root is never a package dir we resolve against
		}
		if !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
			return nil
		}
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, p, nil, 0)
		if err != nil {
			return err
		}
		ast.Inspect(file, func(n ast.Node) bool {
			lit, ok := n.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			val, err := strconv.Unquote(lit.Value)
			if err != nil {
				return true
			}
			if idx := strings.Index(val, "?"); idx != -1 {
				val = val[:idx]
			}
			if pathLiteralRE.MatchString(val) {
				if out[top] == nil {
					out[top] = map[string]bool{}
				}
				out[top][val] = true
			}
			return true
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// resolveTarget finds the one existing SDK package directory whose own
// path literals share the longest leading run of segments with specPath
// (comparing "{param}" spec segments against that package's "%s"
// placeholders positionally, e.g. "/instances/%s/partner-fees" against
// "/instances/{instance_id}/partner-fees/{id}"). Zero or ambiguous matches
// are both refused -- a patcher must never guess which package a new
// route belongs to.
func resolveTarget(repoRoot, specPath string) (*pkgTarget, error) {
	segs := strings.Split(strings.TrimPrefix(specPath, "/"), "/")
	if len(segs) > 0 && segs[0] == "v1" {
		segs = segs[1:]
	}
	norm := make([]string, len(segs))
	for i, s := range segs {
		if strings.HasPrefix(s, "{") {
			norm[i] = "%s"
		} else {
			norm[i] = s
		}
	}

	byDir, err := sdkPathLiteralsByDir(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("scanning SDK path literals: %w", err)
	}

	// "instances/%s/..." is a fixed, universal prefix shared by nearly
	// every route (the instance_id every operation is scoped by); it
	// carries no resource-identifying signal at all, so it is stripped
	// before matching and re-added (as a fixed +2 offset) afterward.
	prefixLen := 0
	if len(norm) >= 2 && norm[0] == "instances" && norm[1] == "%s" {
		prefixLen = 2
	}
	normRest := norm[prefixLen:]

	bestDir, bestLen := "", 0
	ambiguous := false
	for dir, lits := range byDir {
		for lit := range lits {
			litSegs := strings.Split(strings.Trim(lit, "/"), "/")
			litPrefixLen := 0
			if len(litSegs) >= 2 && litSegs[0] == "instances" && litSegs[1] == "%s" {
				litPrefixLen = 2
			}
			litRest := litSegs[litPrefixLen:]

			// Cap the match at the literal's own resource root -- its first
			// "%s" placeholder -- so a sibling id-form route (e.g. an
			// existing Delete's ".../partner-fees/%s") never gets credit
			// for matching all the way through a *different* id value too,
			// which would wrongly swallow the new route's own trailing
			// {id} into the "root" and leave no tail to name a method from.
			litRootLen := 0
			for _, s := range litRest {
				if s == "%s" {
					break
				}
				litRootLen++
			}
			n := 0
			for n < len(normRest) && n < litRootLen && normRest[n] == litRest[n] {
				n++
			}
			if n == 0 {
				continue
			}
			switch {
			case n > bestLen:
				bestLen, bestDir, ambiguous = n, dir, false
			case n == bestLen && dir != bestDir:
				ambiguous = true
			}
		}
	}
	if bestLen == 0 {
		return nil, fmt.Errorf("no existing SDK package has a path literal sharing a leading path segment with %q; a human must add the new package/client wiring", specPath)
	}
	if ambiguous {
		return nil, fmt.Errorf("path %q matches more than one existing SDK package's path prefix equally well; ambiguous, a human must decide", specPath)
	}

	pkgName, err := packageNameOfDir(repoRoot, bestDir)
	if err != nil {
		return nil, err
	}
	return &pkgTarget{Dir: bestDir, Package: pkgName, TailSegs: segs[prefixLen+bestLen:]}, nil
}

func packageNameOfDir(repoRoot, dir string) (string, error) {
	absDir := filepath.Join(repoRoot, dir)
	entries, err := os.ReadDir(absDir)
	if err != nil {
		return "", err
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, filepath.Join(absDir, name), nil, parser.PackageClauseOnly)
		if err != nil {
			continue
		}
		return f.Name.Name, nil
	}
	return "", fmt.Errorf("%s: no non-test .go file found to determine the package name", dir)
}

// existingMethodNames returns every method already declared on dir's
// Client, so a generated name never collides with one already there.
func existingMethodNames(repoRoot, dir string) (map[string]bool, error) {
	out := map[string]bool{}
	path := filepath.Join(repoRoot, dir, "client.go")
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("parsing %s/client.go: %w", dir, err)
	}
	for _, decl := range astFile.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || len(fn.Recv.List) != 1 {
			continue
		}
		out[fn.Name.Name] = true
	}
	return out, nil
}

// allStructShapesInDir is goast.go's findStructShape generalized to every
// struct declared anywhere in dir (any non-test .go file), used to look
// for an existing struct a newly synthesized one can be deduplicated
// against (see findDedupStruct) instead of emitting a redundant type.
func allStructShapesInDir(repoRoot, dir string) ([]*StructShape, error) {
	absDir := filepath.Join(repoRoot, dir)
	entries, err := os.ReadDir(absDir)
	if err != nil {
		return nil, err
	}
	var out []*StructShape
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		rel := dir + "/" + name
		fset := token.NewFileSet()
		astFile, err := parser.ParseFile(fset, filepath.Join(repoRoot, rel), nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		for _, decl := range astFile.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok || structType.Fields == nil {
					continue
				}
				shape := &StructShape{
					File:             rel,
					Symbol:           typeSpec.Name.Name,
					ClosingBraceLine: fset.Position(structType.Fields.Closing).Line,
					Indent:           "\t",
				}
				for _, f := range structType.Fields.List {
					if line := fset.Position(f.End()).Line; line > shape.LastFieldLine {
						shape.LastFieldLine = line
					}
					if len(f.Names) == 0 || f.Tag == nil {
						continue
					}
					tagVal, err := strconv.Unquote(f.Tag.Value)
					if err != nil {
						continue
					}
					jsonName, omitempty := parseJSONTag(tagVal)
					if jsonName == "" || jsonName == "-" {
						continue
					}
					_, isPointer := f.Type.(*ast.StarExpr)
					shape.Fields = append(shape.Fields, FieldShape{
						GoName:    f.Names[0].Name,
						JSONName:  jsonName,
						Omitempty: omitempty,
						Pointer:   isPointer,
						GoType:    typeExprString(f.Type),
					})
				}
				if shape.LastFieldLine == 0 {
					shape.LastFieldLine = shape.ClosingBraceLine - 1
				}
				out = append(out, shape)
			}
		}
	}
	return out, nil
}

func lastLineOf(repoRoot, file string) (int, error) {
	data, err := os.ReadFile(filepath.Join(repoRoot, file)) //#nosec G304 -- path is developer-supplied (CLI arg or repo-relative), this is a local codegen tool
	if err != nil {
		return 0, err
	}
	// Mirrors applyInsertions' own strings.Split(data, "\n"); using the same
	// split keeps "insert after the last line" exactly equivalent to append.
	return len(strings.Split(string(data), "\n")), nil
}

// jsonArrayOpenLine returns the 1-indexed line of `"key": [` in a
// .api-sync/*.json config file, so a new entry can be spliced in as the
// array's new first element via the same insertion mechanism goast.go
// already uses for Go source -- these config files are hand-formatted
// JSON, not machine re-marshaled, so a line splice (not a JSON re-encode)
// is the only way to add an entry without reordering/reformatting
// everything else in the file.
func jsonArrayOpenLine(repoRoot, file, key string) (int, error) {
	data, err := os.ReadFile(filepath.Join(repoRoot, file)) //#nosec G304 -- path is developer-supplied (CLI arg or repo-relative), this is a local codegen tool
	if err != nil {
		return 0, err
	}
	pattern := fmt.Sprintf(`"%s": [`, key)
	for i, line := range strings.Split(string(data), "\n") {
		if strings.Contains(line, pattern) {
			return i + 1, nil
		}
	}
	return 0, fmt.Errorf("%s: %q array not found", file, key)
}

func renderTypeMappingJSON(tm TypeMapping) string {
	return fmt.Sprintf("    {\n      \"spec\": %q,\n      \"sdk\": [{ \"file\": %q, \"symbol\": %q }]\n    },",
		tm.Spec, tm.SDK[0].File, tm.SDK[0].Symbol)
}

// synthField is one property of a struct this generator is about to emit.
type synthField struct {
	JSONName string
	GoName   string
	GoType   string // "string", "[]string", never a pointer -- callers decide pointer-vs-bare per field
	Required bool
}

// flatGoType maps a JSON Schema property definition to a Go scalar or
// array-of-scalar type. ok is false for anything else (enum, object,
// union, array of non-scalar) -- exactly the shapes this deterministic
// generator refuses to guess at, same posture as reconcile.go's
// scalarGoType for hand-added fields on already-mapped types.
func flatGoType(def map[string]any) (goType string, ok bool) {
	if def == nil {
		return "", false
	}
	if _, hasEnum := def["enum"]; hasEnum {
		return "", false
	}
	bases := baseTypes(def)
	if len(bases) == 1 && bases[0] == "array" {
		items, _ := def["items"].(map[string]any)
		if items == nil {
			return "", false
		}
		if _, hasEnum := items["enum"]; hasEnum {
			return "", false
		}
		itemType, ok := scalarGoType(items)
		if !ok {
			return "", false
		}
		return "[]" + itemType, true
	}
	return scalarGoType(def)
}

// flatSchemaFields returns every property of schema as a synthField, or
// ok=false the moment any single property isn't flatGoType-expressible --
// a schema is either entirely synthesizable or not synthesized at all,
// never partially modeled.
func flatSchemaFields(schema map[string]any) ([]synthField, bool) {
	props := properties(schema)
	req := requiredSet(schema)
	names := make([]string, 0, len(props))
	for n := range props {
		names = append(names, n)
	}
	sort.Strings(names)

	fields := make([]synthField, 0, len(names))
	for _, n := range names {
		def, _ := props[n].(map[string]any)
		gt, ok := flatGoType(def)
		if !ok {
			return nil, false
		}
		fields = append(fields, synthField{JSONName: n, GoName: fieldGoName(n), GoType: gt, Required: req[n]})
	}
	return fields, true
}

// fieldGoName is pascalCase specialized for struct field names: an "id"
// segment (whole word, any position, e.g. "id" or "instance_id") is
// capitalized "ID" throughout, matching this repo's existing hand-written
// fields (PartnerFee.ID, VirtualAccountOut's InstanceID, etc.) rather than
// pascalCase's generic "Id".
func fieldGoName(jsonName string) string {
	var b strings.Builder
	for _, part := range strings.FieldsFunc(jsonName, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	}) {
		if strings.EqualFold(part, "id") {
			b.WriteString("ID")
			continue
		}
		b.WriteString(pascalCase(part))
	}
	return b.String()
}

func renderStruct(name string, fields []synthField) string {
	var b strings.Builder
	fmt.Fprintf(&b, "\ntype %s struct {\n", name)
	for _, f := range fields {
		isArray := strings.HasPrefix(f.GoType, "[]")
		switch {
		case f.Required || isArray:
			fmt.Fprintf(&b, "\t%s %s `json:\"%s%s\"`\n", f.GoName, f.GoType, f.JSONName, omitemptySuffix(!f.Required))
		default:
			fmt.Fprintf(&b, "\t%s *%s `json:\"%s,omitempty\"`\n", f.GoName, f.GoType, f.JSONName)
		}
	}
	b.WriteString("}\n")
	return b.String()
}

func omitemptySuffix(optional bool) string {
	if optional {
		return ",omitempty"
	}
	return ""
}

// canonicalSig is the (json name, base Go type) signature used to detect
// that two structs describe the same wire shape, ignoring pointer style.
func canonicalSig(jsonName, goType string) string {
	return jsonName + ":" + strings.TrimPrefix(goType, "*")
}

// findDedupStruct looks for an existing struct in dir whose field set is
// wire-identical to the one about to be synthesized, so a repeated inline
// shape (e.g. the same flat object reused as several sibling operations'
// response) reuses one type instead of each operation minting its own.
func findDedupStruct(repoRoot, dir string, fields []synthField) (string, bool) {
	want := map[string]bool{}
	for _, f := range fields {
		want[canonicalSig(f.JSONName, f.GoType)] = true
	}
	shapes, err := allStructShapesInDir(repoRoot, dir)
	if err != nil {
		return "", false
	}
	for _, s := range shapes {
		if len(s.Fields) != len(fields) {
			continue
		}
		got := map[string]bool{}
		for _, f := range s.Fields {
			got[canonicalSig(f.JSONName, f.GoType)] = true
		}
		if len(got) != len(want) {
			continue
		}
		match := true
		for k := range want {
			if !got[k] {
				match = false
				break
			}
		}
		if match {
			return s.Symbol, true
		}
	}
	return "", false
}

// bodySpec is the resolved shape of one operation's request or response
// body: the Go type to use (bare, never pointer -- callers pointerize
// object types, leave slices bare), plus optionally a brand-new struct's
// source text and/or a spec-map.json registration when the body's named
// schema had no existing SDK mapping.
type bodySpec struct {
	Kind       string // "none" or "type"
	GoType     string
	NewStruct  string
	NewMapping *TypeMapping
}

// resolveBody classifies one requestBody/response content map. Only
// application/json is ever expressible; everything else (multipart,
// binary, or any other media type) is refused with a precise reason.
func resolveBody(repoRoot string, sm *SpecMap, spec *specDoc, dir, methodName, roleSuffix string, content map[string]any) (*bodySpec, string) {
	if len(content) == 0 {
		return &bodySpec{Kind: "none"}, ""
	}
	appJSON, ok := content["application/json"].(map[string]any)
	if !ok || len(content) != 1 {
		var kinds []string
		for k := range content {
			kinds = append(kinds, k)
		}
		sort.Strings(kinds)
		return nil, fmt.Sprintf("non-JSON content type(s) %v (multipart/binary/streaming bodies are not expressible by this generator)", kinds)
	}
	schema, _ := appJSON["schema"].(map[string]any)
	if schema == nil {
		return nil, "application/json content has no schema"
	}

	if ref, ok := schema["$ref"].(string); ok {
		name := strings.TrimPrefix(ref, "#/components/schemas/")
		return resolveNamedSchema(repoRoot, sm, spec, dir, name, false)
	}
	if t, _ := schema["type"].(string); t == "array" {
		items, _ := schema["items"].(map[string]any)
		if items == nil {
			return nil, "array body has no items schema"
		}
		if ref, ok := items["$ref"].(string); ok {
			name := strings.TrimPrefix(ref, "#/components/schemas/")
			return resolveNamedSchema(repoRoot, sm, spec, dir, name, true)
		}
		if gt, ok := flatGoType(schema); ok {
			return &bodySpec{Kind: "type", GoType: gt}, ""
		}
		return nil, "array body items are neither a $ref to a component schema nor a flat scalar shape"
	}

	fields, ok := flatSchemaFields(schema)
	if !ok {
		return nil, "inline body schema has a non-scalar property (nested object, union, or enum-constrained field); needs a human-designed type"
	}
	structName := methodName + roleSuffix
	if reuse, ok := findDedupStruct(repoRoot, dir, fields); ok {
		return &bodySpec{Kind: "type", GoType: reuse}, ""
	}
	return &bodySpec{Kind: "type", GoType: structName, NewStruct: renderStruct(structName, fields)}, ""
}

// resolveNamedSchema handles a $ref (directly, or as an array's items) to
// a component schema. An already-mapped schema is reused as-is, but only
// when every one of its SDK sites already lives in the target package --
// cross-package type reuse would need import management this generator
// deliberately does not attempt. An unmapped schema is synthesized (same
// flat-shape-only rule as inline bodies) and queued for spec-map.json
// registration, so future drift on it is patchable by ordinary reconcile,
// not by a second round of operation-insert.
func resolveNamedSchema(repoRoot string, sm *SpecMap, spec *specDoc, dir, name string, isArray bool) (*bodySpec, string) {
	for _, tm := range sm.Types {
		if tm.Spec != name {
			continue
		}
		if len(tm.SDK) == 0 {
			return nil, fmt.Sprintf("schema %q is mapped in spec-map.json with no SDK site", name)
		}
		for _, site := range tm.SDK {
			if siteDir(site.File) != dir {
				return nil, fmt.Sprintf("schema %q is already mapped to %s (%s), outside the target package %q; cross-package type reuse is not supported by this generator", name, site.Symbol, site.File, dir)
			}
		}
		symbol := tm.SDK[0].Symbol
		gt := symbol
		if isArray {
			gt = "[]" + symbol
		}
		return &bodySpec{Kind: "type", GoType: gt}, ""
	}

	schemaDef, _, ok := spec.schema(name)
	if !ok {
		return nil, fmt.Sprintf("schema %q referenced by $ref could not be resolved in the target spec", name)
	}
	fields, ok := flatSchemaFields(schemaDef)
	if !ok {
		return nil, fmt.Sprintf("schema %q has a non-flat shape (nested object, union, or enum-constrained property); this generator only synthesizes brand-new named schemas that are entirely flat scalar/array-of-scalar", name)
	}

	structName := pascalCase(name)
	var newStruct string
	typeName := structName
	if reuse, ok := findDedupStruct(repoRoot, dir, fields); ok {
		typeName = reuse
	} else {
		newStruct = renderStruct(structName, fields)
	}
	gt := typeName
	if isArray {
		gt = "[]" + typeName
	}
	mapping := &TypeMapping{Spec: name, SDK: []SDKSite{{File: dir + "/client.go", Symbol: typeName}}}
	return &bodySpec{Kind: "type", GoType: gt, NewStruct: newStruct, NewMapping: mapping}, ""
}

func siteDir(file string) string {
	if i := strings.LastIndex(file, "/"); i >= 0 {
		return file[:i]
	}
	return ""
}

func asPtr(t string) string {
	if strings.HasPrefix(t, "[]") {
		return t
	}
	return "*" + t
}

// successResponseContent picks the operation's primary 2xx response body
// (preferring 200, then the lowest remaining 2xx code), returning ok=false
// when there is no response body to model (e.g. 204, or no 2xx at all).
func successResponseContent(op map[string]any) (map[string]any, bool) {
	responses, _ := op["responses"].(map[string]any)
	if responses == nil {
		return nil, false
	}
	var codes []string
	for c := range responses {
		if strings.HasPrefix(c, "2") {
			codes = append(codes, c)
		}
	}
	sort.Strings(codes)
	for _, c := range codes {
		r, _ := responses[c].(map[string]any)
		if content, ok := r["content"].(map[string]any); ok && len(content) > 0 {
			return content, true
		}
	}
	return nil, false
}

// deriveMethodName mechanically names the generated method from the HTTP
// verb and the path segments left over after the matched resource root,
// mirroring this repo's existing hand-written naming (List/Get/Create/
// Update/Delete for a bare resource root, Get<Tail>/<Tail>/Update<Tail>/
// Delete<Tail> for a literal sub-route such as ".../balance" or
// ".../sign-message").
func deriveMethodName(httpMethod string, tailSegs []string) string {
	var suffix strings.Builder
	hasParam := false
	for _, seg := range tailSegs {
		if strings.HasPrefix(seg, "{") {
			hasParam = true
			continue
		}
		suffix.WriteString(pascalCase(seg))
	}
	tail := suffix.String()

	switch httpMethod {
	case "GET":
		if tail != "" {
			return "Get" + tail
		}
		if hasParam {
			return "Get"
		}
		return "List"
	case "POST":
		if tail != "" {
			return tail
		}
		return "Create"
	case "PUT", "PATCH":
		if tail != "" {
			return "Update" + tail
		}
		return "Update"
	case "DELETE":
		if tail != "" {
			return "Delete" + tail
		}
		return "Delete"
	}
	return ""
}

// camelParam converts a snake_case path parameter name into this repo's Go
// parameter-naming convention: lowerCamel, with a standalone or trailing
// "id" segment capitalized to "ID" (matching existing hand-written
// signatures like "customerID, id string").
func camelParam(name string) string {
	parts := strings.Split(name, "_")
	var b strings.Builder
	for i, p := range parts {
		if p == "" {
			continue
		}
		switch {
		case i == 0 && strings.EqualFold(p, "id"):
			b.WriteString("id")
		case strings.EqualFold(p, "id"):
			b.WriteString("ID")
		case i == 0:
			b.WriteString(strings.ToLower(p[:1]) + strings.ToLower(p[1:]))
		default:
			b.WriteString(strings.ToUpper(p[:1]) + strings.ToLower(p[1:]))
		}
	}
	return b.String()
}

// buildPathTemplate turns a spec path (with "{param}" segments) into this
// SDK's fmt.Sprintf style ("%s" placeholders) plus the ordered Sprintf
// arguments -- "{instance_id}" always resolves to "c.instanceID" (this
// SDK's existing convention, never a function parameter), every other
// param consumes the next entry of pathArgs in path order.
func buildPathTemplate(specPath string, pathArgs []string) (string, []string) {
	trimmed := strings.TrimPrefix(specPath, "/v1")
	var fmtBuf strings.Builder
	var args []string
	pi := 0
	for i := 0; i < len(trimmed); {
		if trimmed[i] == '{' {
			j := strings.IndexByte(trimmed[i:], '}')
			name := trimmed[i+1 : i+j]
			fmtBuf.WriteString("%s")
			if name == "instance_id" {
				args = append(args, "c.instanceID")
			} else {
				args = append(args, pathArgs[pi])
				pi++
			}
			i += j + 1
			continue
		}
		fmtBuf.WriteByte(trimmed[i])
		i++
	}
	return fmtBuf.String(), args
}

func renderMethod(methodName, httpMethod, specPath string, pathArgs []string, req, resp *bodySpec) string {
	pathFmt, sprintfArgs := buildPathTemplate(specPath, pathArgs)

	sigParams := []string{"ctx context.Context"}
	for _, a := range pathArgs {
		sigParams = append(sigParams, a+" string")
	}
	bodyArg := "nil"
	if req.Kind != "none" {
		sigParams = append(sigParams, "params "+asPtr(req.GoType))
		bodyArg = "params"
	}

	var retType, callLine string
	if resp.Kind == "none" {
		retType = "error"
		callLine = fmt.Sprintf("_, err := request.Do[struct{}](c.cfg, ctx, %q, path, %s)\n\treturn err", httpMethod, bodyArg)
	} else {
		rt := asPtr(resp.GoType)
		retType = fmt.Sprintf("(%s, error)", rt)
		callLine = fmt.Sprintf("return request.Do[%s](c.cfg, ctx, %q, path, %s)", rt, httpMethod, bodyArg)
	}

	return fmt.Sprintf(`
// %s corresponds to %s %s (generated by cmd/apisync's operation-insert).
func (c *Client) %s(%s) %s {
	path := fmt.Sprintf(%q, %s)
	%s
}
`, methodName, httpMethod, specPath, methodName, strings.Join(sigParams, ", "), retType, pathFmt, strings.Join(sprintfArgs, ", "), callLine)
}

// classifyAndPlanOperation is checkOperationChanges' STANDARD/NON-STANDARD
// decision for exactly one brand-new operation. On success it returns the
// full set of actions to apply (method + any brand-new structs + any
// spec-map.json registrations) and the spec schema names those
// registrations cover (so the caller can keep checkUnclassifiedSchemas
// from re-flagging a schema this same run is about to map). On failure it
// returns a precise NON-STANDARD reason and no actions.
func classifyAndPlanOperation(repoRoot string, sm *SpecMap, spec *specDoc, entry operationEntry) (acts []action, pendingSchemas []string, reason string) {
	target, err := resolveTarget(repoRoot, entry.Path)
	if err != nil {
		return nil, nil, err.Error()
	}
	methods, err := existingMethodNames(repoRoot, target.Dir)
	if err != nil {
		return nil, nil, fmt.Sprintf("could not read %s/client.go: %v", target.Dir, err)
	}

	methodName := deriveMethodName(entry.Method, target.TailSegs)
	if methodName == "" {
		return nil, nil, fmt.Sprintf("could not derive a deterministic method name for %s (unsupported verb)", entry.Method)
	}
	if methods[methodName] {
		return nil, nil, fmt.Sprintf("derived method name %q already exists on %s.Client; a human must pick a non-colliding name", methodName, target.Package)
	}

	req := &bodySpec{Kind: "none"}
	if rb, ok := entry.Op["requestBody"].(map[string]any); ok {
		content, _ := rb["content"].(map[string]any)
		b, why := resolveBody(repoRoot, sm, spec, target.Dir, methodName, "Params", content)
		if why != "" {
			return nil, nil, "request body: " + why
		}
		req = b
	}

	// DELETE always discards the response body in this SDK's existing
	// convention (see e.g. partnerfees.Delete), regardless of what the
	// spec's DELETE response schema is.
	resp := &bodySpec{Kind: "none"}
	if entry.Method != "DELETE" {
		if content, ok := successResponseContent(entry.Op); ok {
			b, why := resolveBody(repoRoot, sm, spec, target.Dir, methodName, "Response", content)
			if why != "" {
				return nil, nil, "response body: " + why
			}
			resp = b
		}
	}

	var pathArgs []string
	for _, seg := range strings.Split(strings.TrimPrefix(entry.Path, "/v1"), "/") {
		if !strings.HasPrefix(seg, "{") {
			continue
		}
		name := strings.TrimSuffix(strings.TrimPrefix(seg, "{"), "}")
		if name == "instance_id" {
			continue
		}
		pathArgs = append(pathArgs, camelParam(name))
	}

	methodText := renderMethod(methodName, entry.Method, entry.Path, pathArgs, req, resp)

	var pieces []string
	if req.NewStruct != "" {
		pieces = append(pieces, req.NewStruct)
	}
	if resp.NewStruct != "" && resp.NewStruct != req.NewStruct {
		pieces = append(pieces, resp.NewStruct)
	}
	pieces = append(pieces, methodText)

	clientFile := target.Dir + "/client.go"
	lastLine, err := lastLineOf(repoRoot, clientFile)
	if err != nil {
		return nil, nil, fmt.Sprintf("could not read %s: %v", clientFile, err)
	}

	acts = append(acts, action{
		description: fmt.Sprintf("operation-insert: %s %s -> %s.%s (%s)", entry.Method, entry.Path, target.Package, methodName, clientFile),
		insert:      insertion{File: clientFile, Line: lastLine, Text: strings.Join(pieces, "\n")},
		bump:        "minor",
	})

	for _, m := range []*TypeMapping{req.NewMapping, resp.NewMapping} {
		if m == nil {
			continue
		}
		line, err := jsonArrayOpenLine(repoRoot, ".api-sync/spec-map.json", "types")
		if err != nil {
			return nil, nil, fmt.Sprintf("could not locate spec-map.json types[] array: %v", err)
		}
		acts = append(acts, action{
			description: fmt.Sprintf("operation-insert: spec-map.json: register schema %s -> %s.%s", m.Spec, target.Package, m.SDK[0].Symbol),
			insert:      insertion{File: ".api-sync/spec-map.json", Line: line, Text: renderTypeMappingJSON(*m)},
			bump:        "minor",
		})
		pendingSchemas = append(pendingSchemas, m.Spec)
	}

	return acts, pendingSchemas, ""
}
