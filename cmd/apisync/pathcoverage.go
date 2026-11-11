package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// operationSet returns every SDK "METHOD /path" key in the spec's paths
// object (x-sdk: false operations are skipped, see operationEntries).
func operationSet(spec *specDoc) map[string]bool {
	out := map[string]bool{}
	for key := range operationEntries(spec) {
		out[key] = true
	}
	return out
}

// operationSetChanged reports whether the spec's operation set differs
// between old and new, feeding "any operation change -> minor" bump
// classification. It does not by itself fail the build (see runCoverage
// for the non-blocking SDK-coverage report).
func operationSetChanged(oldSpec, newSpec *specDoc) bool {
	oldOps, newOps := operationSet(oldSpec), operationSet(newSpec)
	if len(oldOps) != len(newOps) {
		return true
	}
	for k := range oldOps {
		if !newOps[k] {
			return true
		}
	}
	return false
}

var pathLiteralRE = regexp.MustCompile(`^/[a-zA-Z0-9_%./-]*$`)

var skipDirsCoverage = map[string]bool{
	".git":      true,
	"examples":  true,
	"cmd":       true,
	".api-sync": true,
	"testdata":  true,
}

// funcFacts is what sdkOperations needs to know about one function.
type funcFacts struct {
	methods map[string]bool // HTTP methods passed to request.Do* calls
	paths   map[string]bool // path-template string literals
	callees []string        // same-package functions/methods it calls
}

// sdkOperations returns every "METHOD /path-template" the SDK sends, keyed
// as the spec would be after normalizeSpecPath ("GET /instances/%s/payouts").
//
// It pairs, per function, the HTTP method of each request.Do* call
// (request.Do's method argument; DoList and DoPage are GET, DoMultipart is POST) with
// the path literals of that function and of the same-package functions it
// calls (e.g. a listPath helper). Methods are what make the match exact:
// GET and PATCH on the same path are two operations.
func sdkOperations(repoRoot string) (map[string]bool, error) {
	byPkg := map[string]map[string]*funcFacts{}
	err := filepath.WalkDir(repoRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if skipDirsCoverage[d.Name()] && path != repoRoot {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			return err
		}
		dir := filepath.Dir(path)
		if byPkg[dir] == nil {
			byPkg[dir] = map[string]*funcFacts{}
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			facts := byPkg[dir][fn.Name.Name]
			if facts == nil {
				facts = &funcFacts{methods: map[string]bool{}, paths: map[string]bool{}}
				byPkg[dir][fn.Name.Name] = facts
			}
			collectFuncFacts(fn.Body, facts)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	out := map[string]bool{}
	for _, funcs := range byPkg {
		for _, facts := range funcs {
			if len(facts.methods) == 0 {
				continue
			}
			paths := map[string]bool{}
			seen := map[*funcFacts]bool{}
			var visit func(f *funcFacts)
			visit = func(f *funcFacts) {
				if f == nil || seen[f] {
					return
				}
				seen[f] = true
				for p := range f.paths {
					paths[p] = true
				}
				for _, callee := range f.callees {
					visit(funcs[callee])
				}
			}
			visit(facts)
			for m := range facts.methods {
				for p := range paths {
					out[m+" "+p] = true
				}
			}
		}
	}
	return out, nil
}

func collectFuncFacts(body ast.Node, facts *funcFacts) {
	ast.Inspect(body, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.BasicLit:
			if v.Kind != token.STRING {
				return true
			}
			val, err := strconv.Unquote(v.Value)
			if err != nil {
				return true
			}
			if idx := strings.Index(val, "?"); idx != -1 {
				val = val[:idx] // strip a literal query string, e.g. "/upload?"
			}
			if len(val) > 1 && pathLiteralRE.MatchString(val) {
				facts.paths[val] = true
			}
		case *ast.CallExpr:
			if m := transportMethod(v); m != "" {
				facts.methods[m] = true
				return true
			}
			switch fun := v.Fun.(type) {
			case *ast.Ident:
				facts.callees = append(facts.callees, fun.Name)
			case *ast.SelectorExpr:
				if x, ok := fun.X.(*ast.Ident); ok && x.Name != "request" {
					facts.callees = append(facts.callees, fun.Sel.Name)
				}
			}
		}
		return true
	})
}

// transportMethod returns the HTTP method of a request.Do* call, or "".
func transportMethod(call *ast.CallExpr) string {
	fun := call.Fun
	switch f := fun.(type) {
	case *ast.IndexExpr:
		fun = f.X
	case *ast.IndexListExpr:
		fun = f.X
	}
	sel, ok := fun.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	if pkg, ok := sel.X.(*ast.Ident); !ok || pkg.Name != "request" {
		return ""
	}
	switch sel.Sel.Name {
	case "DoList", "DoPage":
		return "GET"
	case "DoMultipart":
		return "POST"
	case "Do":
		if len(call.Args) < 3 {
			return ""
		}
		switch a := call.Args[2].(type) {
		case *ast.BasicLit:
			if m, err := strconv.Unquote(a.Value); err == nil {
				return strings.ToUpper(m)
			}
		case *ast.SelectorExpr: // http.MethodGet etc.
			return strings.ToUpper(strings.TrimPrefix(a.Sel.Name, "Method"))
		}
	}
	return ""
}

// normalizeSpecPath strips the "/v1" version prefix and turns "{param}"
// segments into "%s", matching this SDK's path-template style.
var pathParamRE = regexp.MustCompile(`\{[^}]+\}`)

func normalizeSpecPath(p string) string {
	p = strings.TrimPrefix(p, "/v1")
	return pathParamRE.ReplaceAllString(p, "%s")
}

// specOperationKeys returns normalized "METHOD /path" -> spec key for the
// spec's SDK operations (sdk) and for its x-sdk: false operations (excluded).
func specOperationKeys(spec *specDoc) (sdk, excluded map[string]string) {
	sdk, excluded = map[string]string{}, map[string]string{}
	paths, _ := spec.raw["paths"].(map[string]any)
	for p, methodsRaw := range paths {
		methods, ok := methodsRaw.(map[string]any)
		if !ok {
			continue
		}
		for m, opRaw := range methods {
			if !isHTTPMethod(m) {
				continue
			}
			op, _ := opRaw.(map[string]any)
			method := strings.ToUpper(m)
			norm := method + " " + normalizeSpecPath(p)
			if excludedFromSDK(op) {
				excluded[norm] = method + " " + p
			} else {
				sdk[norm] = method + " " + p
			}
		}
	}
	return sdk, excluded
}

// computeCoverageGaps returns a sorted, non-blocking list of (a) spec SDK
// operations no SDK method sends, (b) x-sdk: false operations the SDK does
// send (it must not), and (c) SDK method+path pairs with no spec operation
// (e.g. the deprecated Export methods' /export/* routes). Matching is on
// METHOD + path, so PUT and GET on one path are distinct operations.
func computeCoverageGaps(repoRoot string, spec *specDoc) []string {
	sdkOps, err := sdkOperations(repoRoot)
	if err != nil {
		return []string{fmt.Sprintf("coverage report unavailable: %v", err)}
	}
	specOps, excluded := specOperationKeys(spec)

	var gaps []string
	for norm, orig := range specOps {
		if !sdkOps[norm] {
			gaps = append(gaps, fmt.Sprintf("spec operation with no SDK method: %s", orig))
		}
	}
	for norm := range sdkOps {
		if orig, ok := excluded[norm]; ok {
			gaps = append(gaps, fmt.Sprintf("SDK implements x-sdk: false operation: %s", orig))
			continue
		}
		if _, ok := specOps[norm]; !ok {
			gaps = append(gaps, fmt.Sprintf("SDK operation with no matching spec operation: %s", norm))
		}
	}
	sort.Strings(gaps)
	return gaps
}

// runCoverage prints the coverage report against specPath (default: the
// committed snapshot).
func runCoverage(repoRoot, specPath string) error {
	if specPath == "" {
		specPath = ".api-sync/spec-snapshot.json"
	}
	_, newRaw, err := readSpecFile(repoRoot, specPath)
	if err != nil {
		return err
	}
	spec := loadSpecDoc(newRaw)
	gaps := computeCoverageGaps(repoRoot, spec)
	sdkOps, _ := specOperationKeys(spec)
	fmt.Println("=== apisync coverage report (non-blocking) ===")
	fmt.Printf("spec SDK operations: %d, covered: %d\n", len(sdkOps), len(sdkOps)-countPrefix(gaps, "spec operation with no SDK method"))
	if len(gaps) == 0 {
		fmt.Println("no gaps found")
		return nil
	}
	for _, g := range gaps {
		fmt.Println("  -", g)
	}
	return nil
}

func countPrefix(lines []string, prefix string) int {
	n := 0
	for _, l := range lines {
		if strings.HasPrefix(l, prefix) {
			n++
		}
	}
	return n
}
