package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuditTypes_ReportsEveryMismatchRegardlessOfUnmodeledJSON(t *testing.T) {
	root := widgetFixture(t)
	sm := &SpecMap{Types: []TypeMapping{widgetTypeMapping()}}

	spec := mustSpec(t, map[string]string{
		"Widget": `{"type":"object","properties":{
			"id":{"type":"integer"},
			"name":{"type":["string","null"]}
		},"required":["id"]}`,
	})

	t.Run("unexcused finding is reported and marked as such", func(t *testing.T) {
		um := &Unmodeled{}
		findings := auditTypes(root, sm, um, spec)
		require.Len(t, findings, 1)
		require.Contains(t, findings[0], "Widget.ID")
		require.Contains(t, findings[0], "NOT in unmodeled.json")
	})

	t.Run("excused finding is still reported, just annotated", func(t *testing.T) {
		um := &Unmodeled{Properties: []PropertyExclusion{
			{Schema: "Widget", Property: "id", Reason: "known, tracked", Owner: "eric"},
		}}
		findings := auditTypes(root, sm, um, spec)
		require.Len(t, findings, 1, "audit-types is non-blocking and unconditional: it must not go silent just because unmodeled.json excuses check/apply")
		require.Contains(t, findings[0], "recorded in unmodeled.json")
	})
}

func TestAuditTypes_NoMismatchesIsAnEmptyReport(t *testing.T) {
	root := widgetFixture(t)
	sm := &SpecMap{Types: []TypeMapping{widgetTypeMapping()}}
	um := &Unmodeled{}
	spec := mustSpec(t, map[string]string{
		"Widget": baselineWidgetSchema,
	})
	findings := auditTypes(root, sm, um, spec)
	require.Empty(t, findings)
}

func TestRunAuditTypes_AlwaysExitsCleanEvenWithMismatches(t *testing.T) {
	root := integrationFixture(t, map[string]string{"Widget": baselineWidgetSchema})
	writeSpecCurrent(t, root, map[string]string{"Widget": baselineWidgetSchema})
	// Mutate the committed snapshot (what -audit-types reads) to introduce a
	// real mismatch, then confirm the survey still exits nil (non-blocking).
	mutatedSnapshot := `{"components":{"schemas":{"Widget":{"type":"object","properties":{` +
		`"id":{"type":"integer"},"name":{"type":["string","null"]},"color":{"type":"string","enum":["red","blue"]}` +
		`},"required":["id"]}}},"paths":{}}`
	require.NoError(t, os.WriteFile(filepath.Join(root, ".api-sync", "spec-snapshot.json"), []byte(mutatedSnapshot), 0o644))

	err := runAuditTypes(root)
	require.NoError(t, err, "-audit-types must always exit 0, even when it finds real mismatches")
}
