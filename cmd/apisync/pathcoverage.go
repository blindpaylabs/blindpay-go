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

// operationSet returns every "METHOD /path" key in the spec's paths object.
func operationSet(spec *specDoc) map[string]bool {
	out := map[string]bool{}
	paths, _ := spec.raw["paths"].(map[string]any)
	for p, methodsRaw := range paths {
		methods, ok := methodsRaw.(map[string]any)
		if !ok {
			continue
		}
		for m := range methods {
			switch m {
			case "get", "post", "put", "patch", "delete":
				out[strings.ToUpper(m)+" "+p] = true
			}
		}
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
}

// sdkPathLiterals scans every non-test .go file for string literals that
// look like path templates (leading "/", %s placeholders already in the
// SDK's own interpolation style).
func sdkPathLiterals(repoRoot string) (map[string]bool, error) {
	out := map[string]bool{}
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
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, 0)
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
				val = val[:idx] // strip a literal query string, e.g. "...?rail=%s"
			}
			if pathLiteralRE.MatchString(val) {
				out[val] = true
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

// normalizeSpecPath strips the "/v1" version prefix and turns "{param}"
// segments into "%s", matching this SDK's fmt.Sprintf path-template style.
var pathParamRE = regexp.MustCompile(`\{[^}]+\}`)

func normalizeSpecPath(p string) string {
	p = strings.TrimPrefix(p, "/v1")
	return pathParamRE.ReplaceAllString(p, "%s")
}

// computeCoverageGaps returns a human-readable, sorted, non-blocking list
// of (a) spec operations with no matching SDK path literal and (b) SDK path
// literals with no matching spec operation (e.g. this repo's known
// pre-existing drift: the two spec-absent /export/* paths, which have no
// corresponding operation in the current public spec). Scope note: this is a path-shape match
// only (method is not correlated back to the literal, since the SDK builds
// path variables separately from the request.Do method argument); good
// enough for a non-blocking survey, not a substitute for the map-validity
// or removal checks above. Known false positive: upload.Client.Upload
// builds its URL as "%s/upload?..." with the base URL prepended (not a
// path-only literal starting with "/"), so POST /v1/upload is reported as
// an apparent gap even though it is implemented; harmless since this
// report never fails the build.
func computeCoverageGaps(repoRoot string, spec *specDoc) []string {
	sdkPaths, err := sdkPathLiterals(repoRoot)
	if err != nil {
		return []string{fmt.Sprintf("coverage report unavailable: %v", err)}
	}

	specNormalized := map[string]string{} // normalized -> original spec path
	paths, _ := spec.raw["paths"].(map[string]any)
	for p := range paths {
		specNormalized[normalizeSpecPath(p)] = p
	}

	var gaps []string
	for norm, orig := range specNormalized {
		if !sdkPaths[norm] {
			gaps = append(gaps, fmt.Sprintf("spec operation with no SDK path: %s (normalized %s)", orig, norm))
		}
	}
	for sdkPath := range sdkPaths {
		if _, ok := specNormalized[sdkPath]; !ok {
			gaps = append(gaps, fmt.Sprintf("SDK path with no matching spec operation: %s", sdkPath))
		}
	}
	sort.Strings(gaps)
	return gaps
}

func runCoverage(repoRoot string) error {
	_, newRaw, err := readSpecFile(repoRoot, ".api-sync/spec-snapshot.json")
	if err != nil {
		return err
	}
	gaps := computeCoverageGaps(repoRoot, loadSpecDoc(newRaw))
	fmt.Println("=== apisync coverage report (non-blocking) ===")
	if len(gaps) == 0 {
		fmt.Println("no gaps found")
		return nil
	}
	for _, g := range gaps {
		fmt.Println("  -", g)
	}
	return nil
}
