package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const errorsTagHeader = "### Error codes\n\n| Code | Status | Retry | Description |\n| --- | --- | --- | --- |\n"

func specWithErrors(table string) map[string]any {
	return map[string]any{
		"tags": []any{
			map[string]any{"name": "Payouts"},
			map[string]any{"name": "Errors", "description": "## Error responses\n\nintro\n\n" + errorsTagHeader + table},
		},
	}
}

func TestParseErrorCodeTable(t *testing.T) {
	rows, found, err := parseErrorCodeTable(specWithErrors(
		"| `QUOTES_EXPIRED` | 400 | actionable | This quote has expired. |\n" +
			"| `AUTH_IP_BLOCKED` | 403 | terminal | Requests from this IP | address are not allowed. |\n"))
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, []errorCodeRow{
		{Code: "AUTH_IP_BLOCKED", Status: "403", Retry: "terminal", Description: "Requests from this IP | address are not allowed."},
		{Code: "QUOTES_EXPIRED", Status: "400", Retry: "actionable", Description: "This quote has expired."},
	}, rows, "sorted by code; a | inside the description is kept")
	require.Equal(t, "ErrorCodeAuthIPBlocked", errorCodeIdent("AUTH_IP_BLOCKED"))

	_, found, err = parseErrorCodeTable(map[string]any{"tags": []any{map[string]any{"name": "Payouts"}}})
	require.NoError(t, err)
	require.False(t, found)

	for name, table := range map[string]string{
		"zero rows":     "",
		"unknown retry": "| `QUOTES_EXPIRED` | 400 | maybe | x |\n",
		"duplicate":     "| `QUOTES_EXPIRED` | 400 | actionable | x |\n| `QUOTES_EXPIRED` | 400 | actionable | y |\n",
	} {
		_, found, err := parseErrorCodeTable(specWithErrors(table))
		require.True(t, found, name)
		require.Error(t, err, name)
	}
}

// TestParseErrorCodeTable_Snapshot: the committed snapshot's table parses
// to the 175 codes of blindpay-v2@7e49addd9.
func TestParseErrorCodeTable_Snapshot(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", ".api-sync", "spec-snapshot.json"))
	require.NoError(t, err)
	var raw map[string]any
	require.NoError(t, json.Unmarshal(data, &raw))
	rows, found, err := parseErrorCodeTable(raw)
	require.NoError(t, err)
	require.True(t, found)
	require.Len(t, rows, 175)
}

func TestPlanErrorCodes(t *testing.T) {
	one := "| `QUOTES_EXPIRED` | 400 | actionable | This quote has expired. |\n"
	two := one + "| `TRANSFERS_SWAP_AMOUNT_BELOW_MINIMUM` | 400 | actionable | Swaps need 100 USDC. |\n"

	root := writeFixture(t, map[string]string{"go.mod": "module example.com/fixture\n"})
	require.NoError(t, os.MkdirAll(filepath.Join(root, "internal", "request"), 0o755))

	// No Errors tag and no generated file: nothing to do (test fixtures).
	actions, nh := planErrorCodes(root, loadSpecDoc(map[string]any{}))
	require.Empty(t, actions)
	require.Empty(t, nh)

	// First generation: pending, minor.
	actions, nh = planErrorCodes(root, loadSpecDoc(specWithErrors(one)))
	require.Empty(t, nh)
	require.Len(t, actions, 1)
	require.Equal(t, "minor", actions[0].bump)
	writeAll(t, root, actions[0].writes)

	// Same spec: clean (determinism, -check silent).
	actions, nh = planErrorCodes(root, loadSpecDoc(specWithErrors(one)))
	require.Empty(t, actions)
	require.Empty(t, nh)

	// A new code: pending, minor, named in the description.
	actions, nh = planErrorCodes(root, loadSpecDoc(specWithErrors(two)))
	require.Empty(t, nh)
	require.Len(t, actions, 1)
	require.Equal(t, "minor", actions[0].bump)
	require.Contains(t, actions[0].description, "TRANSFERS_SWAP_AMOUNT_BELOW_MINIMUM")
	writeAll(t, root, actions[0].writes)
	require.Contains(t, readFixture(t, root, errorCodesRootFile),
		`ErrorCodeTransfersSwapAmountBelowMinimum ErrorCode = "TRANSFERS_SWAP_AMOUNT_BELOW_MINIMUM"`)

	// A description-only change: patch.
	actions, _ = planErrorCodes(root, loadSpecDoc(specWithErrors(strings.Replace(two, "100 USDC", "100 USDC on the USDC side", 1))))
	require.Len(t, actions, 1)
	require.Equal(t, "patch", actions[0].bump)

	// A removed code is NEEDS_HUMAN: its exported constant is public API.
	actions, nh = planErrorCodes(root, loadSpecDoc(specWithErrors(one)))
	require.Empty(t, actions)
	require.Len(t, nh, 1)
	require.Contains(t, nh[0], "TRANSFERS_SWAP_AMOUNT_BELOW_MINIMUM")

	// The Errors tag disappearing while the generated file exists is NEEDS_HUMAN.
	_, nh = planErrorCodes(root, loadSpecDoc(map[string]any{}))
	require.Len(t, nh, 1)

	// A table with no rows is NEEDS_HUMAN, never an empty file.
	actions, nh = planErrorCodes(root, loadSpecDoc(specWithErrors("")))
	require.Empty(t, actions)
	require.Len(t, nh, 1)
}

func writeAll(t *testing.T, root string, files map[string][]byte) {
	t.Helper()
	for file, data := range files {
		require.NoError(t, os.WriteFile(filepath.Join(root, file), data, 0o644))
	}
}
