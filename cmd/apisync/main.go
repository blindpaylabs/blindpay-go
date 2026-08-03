// Command apisync is a mechanical, dependency-free patcher that reconciles
// this SDK against the BlindPay OpenAPI spec, mirroring cmd/contractcheck's
// style (go/ast + go/parser, stdlib only, no external deps).
//
// It computes a state-reconciliation plan against a target spec (every
// mapped enum/schema in .api-sync/spec-map.json must be fully present in
// the SDK, or excused by .api-sync/unmodeled.json), plus an old-vs-new diff
// against the committed .api-sync/spec-snapshot.json for two things state
// alone can't tell you: removal detection (always a hard fail) and version
// bump classification.
//
//   - go run ./cmd/apisync -check    state assertion; exit 1 with precise
//     messages on any pending drift or NEEDS_HUMAN issue; silent when clean.
//   - go run ./cmd/apisync -apply    same computation; if anything is
//     NEEDS_HUMAN, changes nothing and exits 1; otherwise applies every
//     APPLICABLE change, gofmts the touched files, bumps blindpay.go's
//     Version, and refreshes the snapshot.
//   - -spec <path>   spec file to reconcile against. Defaults to
//     .api-sync/spec-snapshot.json for -check (so plain `-check` works in
//     ordinary CI where spec-current.json does not exist) and to
//     .api-sync/spec-current.json for -apply.
//   - -report <path>   write a JSON report (applied, needs-human, bump,
//     coverage gaps) to this path.
//   - -coverage   print the non-blocking spec-operations-with-no-SDK-path
//     report and exit 0. Independent of -check/-apply; never fails the build.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "apisync:", err)
		os.Exit(1)
	}
}

func run() error {
	checkMode := flag.Bool("check", false, "assert state, exit 1 on any pending drift or needs-human issue")
	applyMode := flag.Bool("apply", false, "apply every applicable change; exit 1 without writing anything if any needs-human issue exists")
	specFlag := flag.String("spec", "", "spec file to reconcile against (default: spec-snapshot.json for -check, spec-current.json for -apply)")
	reportPath := flag.String("report", "", "write a JSON report to this path")
	coverage := flag.Bool("coverage", false, "print the non-blocking spec-operations-with-no-SDK-path coverage report and exit 0")
	validateMap := flag.Bool("validate-map", false, "resolve every spec-map.json anchor (files, symbols) and exit 1 on the first one that does not; a fast, standalone preflight")
	flag.Parse()

	repoRoot, err := findRepoRoot()
	if err != nil {
		return err
	}

	if *coverage {
		return runCoverage(repoRoot)
	}

	if *validateMap {
		return runValidateMap(repoRoot)
	}

	if *checkMode == *applyMode {
		return fmt.Errorf("exactly one of -check or -apply is required")
	}

	quiet, err := reconcileAndMaybeApply(repoRoot, syncOptions{
		Check:      *checkMode,
		Apply:      *applyMode,
		SpecPath:   *specFlag,
		ReportPath: *reportPath,
	})
	if quiet {
		return nil
	}
	return err
}

// syncOptions carries every -check/-apply behavioral input, kept separate
// from the flag package so the core logic is directly unit-testable.
type syncOptions struct {
	Check      bool
	Apply      bool
	SpecPath   string // empty means apply the mode-based default
	ReportPath string // empty means do not write a report
}

// reconcileAndMaybeApply is the whole tool's core, independent of the CLI
// flag package: load the map/exclusions, reconcile against the target
// spec, and either report (check) or write out every applicable change
// (apply). quiet is true exactly for -check's "silent when clean" path,
// letting callers distinguish "nothing to report" from "nothing failed".
func reconcileAndMaybeApply(repoRoot string, opts syncOptions) (quiet bool, err error) {
	sm, err := loadSpecMap(repoRoot)
	if err != nil {
		return false, fmt.Errorf("loading spec-map.json: %w", err)
	}
	um, err := loadUnmodeled(repoRoot)
	if err != nil {
		return false, fmt.Errorf("loading unmodeled.json: %w", err)
	}

	specPath := opts.SpecPath
	if specPath == "" {
		if opts.Check {
			specPath = ".api-sync/spec-snapshot.json"
		} else {
			specPath = ".api-sync/spec-current.json"
		}
	}

	newRaw, err := readSpecFile(repoRoot, specPath)
	if err != nil {
		return false, fmt.Errorf("reading %s: %w", specPath, err)
	}
	oldRaw, err := readSpecFile(repoRoot, ".api-sync/spec-snapshot.json")
	if err != nil {
		return false, fmt.Errorf("reading .api-sync/spec-snapshot.json: %w", err)
	}
	newSpec := loadSpecDoc(newRaw)
	oldSpec := loadSpecDoc(oldRaw)

	var needsHuman []string
	needsHuman = append(needsHuman, checkMapValidity(repoRoot, sm)...)
	needsHuman = append(needsHuman, checkUnclassifiedSchemas(sm, newSpec)...)

	enumPlan := reconcileEnums(repoRoot, sm, um, oldSpec, newSpec)
	typePlan := reconcileTypes(repoRoot, sm, um, oldSpec, newSpec)
	needsHuman = append(needsHuman, enumPlan.NeedsHuman...)
	needsHuman = append(needsHuman, typePlan.NeedsHuman...)
	sort.Strings(needsHuman)

	actions := append(append([]action{}, enumPlan.Actions...), typePlan.Actions...)
	sort.Slice(actions, func(i, j int) bool { return actions[i].description < actions[j].description })

	gaps := computeCoverageGaps(repoRoot, newSpec)

	report := buildReport(actions, needsHuman, gaps, opts.Apply)

	if opts.ReportPath != "" {
		if err := writeReport(opts.ReportPath, report); err != nil {
			return false, fmt.Errorf("writing report: %w", err)
		}
	}

	if opts.Check {
		if len(needsHuman) == 0 && len(actions) == 0 {
			return true, nil // silent when clean
		}
		for _, msg := range needsHuman {
			fmt.Println(msg)
		}
		for _, a := range actions {
			fmt.Printf("PENDING (applicable, run -apply): %s\n", a.description)
		}
		return false, fmt.Errorf("%d needs-human issue(s), %d pending applicable change(s)", len(needsHuman), len(actions))
	}

	// -apply
	if len(needsHuman) > 0 {
		for _, msg := range needsHuman {
			fmt.Println(msg)
		}
		return false, fmt.Errorf("%d needs-human issue(s); nothing was changed", len(needsHuman))
	}
	if len(actions) == 0 {
		return false, nil // nothing to do
	}

	insertions := make([]insertion, 0, len(actions))
	touchedFiles := map[string]bool{}
	bump := "patch"
	for _, a := range actions {
		insertions = append(insertions, a.insert)
		touchedFiles[a.insert.File] = true
		if a.bump == "minor" {
			bump = "minor"
		}
	}
	if operationSetChanged(oldSpec, newSpec) {
		bump = "minor"
	}

	if err := applyInsertions(repoRoot, insertions); err != nil {
		return false, fmt.Errorf("applying changes: %w", err)
	}
	for file := range touchedFiles {
		if err := gofmtFile(filepath.Join(repoRoot, file)); err != nil {
			return false, fmt.Errorf("gofmt %s: %w", file, err)
		}
	}

	newVersion, err := bumpVersion(repoRoot, bump)
	if err != nil {
		return false, fmt.Errorf("bumping version: %w", err)
	}
	fmt.Printf("applied %d change(s), version bumped to %s (%s)\n", len(actions), newVersion, bump)

	if err := refreshSnapshot(repoRoot, newRaw); err != nil {
		return false, fmt.Errorf("refreshing snapshot: %w", err)
	}

	return false, nil
}

// runValidateMap resolves every spec-map.json anchor with no spec or
// SDK-state reconciliation at all: a fast, standalone "is the map itself
// well-formed" preflight, distinct from -check's full state assertion
// (which also runs this internally, since a stale map must never silently
// skip a mapping either way).
func runValidateMap(repoRoot string) error {
	sm, err := loadSpecMap(repoRoot)
	if err != nil {
		return fmt.Errorf("loading spec-map.json: %w", err)
	}
	if _, err := loadUnmodeled(repoRoot); err != nil {
		return fmt.Errorf("loading unmodeled.json: %w", err)
	}
	issues := checkMapValidity(repoRoot, sm)
	if len(issues) == 0 {
		return nil
	}
	for _, msg := range issues {
		fmt.Println(msg)
	}
	return fmt.Errorf("%d spec-map.json anchor(s) did not resolve", len(issues))
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

func readSpecFile(repoRoot, relPath string) (map[string]any, error) {
	data, err := os.ReadFile(filepath.Join(repoRoot, relPath))
	if err != nil {
		return nil, err
	}
	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return doc, nil
}

func loadSpecMap(repoRoot string) (*SpecMap, error) {
	data, err := os.ReadFile(filepath.Join(repoRoot, ".api-sync", "spec-map.json"))
	if err != nil {
		return nil, err
	}
	var sm SpecMap
	if err := json.Unmarshal(data, &sm); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return &sm, nil
}

func loadUnmodeled(repoRoot string) (*Unmodeled, error) {
	data, err := os.ReadFile(filepath.Join(repoRoot, ".api-sync", "unmodeled.json"))
	if err != nil {
		return nil, err
	}
	var um Unmodeled
	if err := json.Unmarshal(data, &um); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	for i, p := range um.Properties {
		if p.Schema == "" || p.Property == "" || p.Reason == "" || p.Owner == "" {
			return nil, fmt.Errorf("properties[%d] missing schema/property/reason/owner", i)
		}
	}
	for i, e := range um.Enums {
		if e.Symbol == "" || e.Member == "" || e.Reason == "" || e.Owner == "" {
			return nil, fmt.Errorf("enums[%d] missing symbol/member/reason/owner", i)
		}
	}
	return &um, nil
}

func gofmtFile(path string) error {
	cmd := exec.Command("gofmt", "-w", path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w", string(out), err)
	}
	return nil
}
