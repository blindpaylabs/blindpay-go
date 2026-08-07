package main

import (
	"encoding/json"
	"os"
	"sort"
)

// Report is the -report JSON output: applied (or, in -check mode, pending)
// changes, needs-human issues, the version-bump class, and non-blocking
// coverage gaps. Slices are always sorted for determinism.
type Report struct {
	Applied      []string `json:"applied"`
	NeedsHuman   []string `json:"needs_human"`
	Bump         string   `json:"bump"`
	CoverageGaps []string `json:"coverage_gaps"`
}

func buildReport(actions []action, needsHuman []string, gaps []string, applyMode bool) Report {
	applied := make([]string, 0, len(actions))
	bump := "none"
	for _, a := range actions {
		applied = append(applied, a.description)
		if a.bump == "minor" {
			bump = "minor"
		} else if a.bump == "patch" && bump != "minor" {
			bump = "patch"
		}
	}
	sort.Strings(applied)

	nh := append([]string{}, needsHuman...)
	sort.Strings(nh)

	cg := append([]string{}, gaps...)
	sort.Strings(cg)

	// "applied" lists what was (or, in -check mode, would be) changed; the
	// caller distinguishes check-vs-apply via how the process was invoked.
	return Report{
		Applied:      applied,
		NeedsHuman:   nh,
		Bump:         bump,
		CoverageGaps: cg,
	}
}

func writeReport(path string, r Report) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o600)
}
