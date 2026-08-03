package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// integrationFixture builds a minimal but complete repo: go.mod, an enum
// (with root re-export), a struct, blindpay.go's Version const, and a
// matching spec-map.json/unmodeled.json/spec-snapshot.json. specSchemas
// seeds spec-snapshot.json (the "old"/committed baseline); a test then
// writes spec-current.json separately to represent "new".
func integrationFixture(t *testing.T, schemas map[string]string) string {
	t.Helper()
	spec := map[string]any{}
	for name, raw := range schemas {
		var v any
		require.NoError(t, json.Unmarshal([]byte(raw), &v))
		spec[name] = v
	}
	specDoc := map[string]any{
		"components": map[string]any{"schemas": spec},
		"paths":      map[string]any{},
	}
	specJSON, err := json.MarshalIndent(specDoc, "", "  ")
	require.NoError(t, err)

	specMap := `{
		"enums": [
			{
				"sdk": {"file": "internal/types/enums.go", "symbol": "Color"},
				"reexport": {"file": "types.go", "symbol": "Color"},
				"spec": {"schema": "Widget", "property": "color"}
			}
		],
		"types": [
			{
				"spec": "Widget",
				"sdk": [{"file": "widgets/client.go", "symbol": "Widget"}]
			}
		],
		"ignore": {"schemas": []}
	}`
	unmodeled := `{"properties": [], "enums": []}`

	root := writeFixture(t, map[string]string{
		"go.mod": "module example.com/fixture\n\ngo 1.21\n",
		"blindpay.go": `package fixture

const Version = "1.0.0"
`,
		"internal/types/enums.go": `package types

type Color string

const (
	ColorRed  Color = "red"
	ColorBlue Color = "blue"
)
`,
		"types.go": `package fixture

import "example.com/fixture/internal/types"

type Color = types.Color

const (
	ColorRed  = types.ColorRed
	ColorBlue = types.ColorBlue
)
`,
		"widgets/client.go": `package widgets

type Widget struct {
	ID    string  ` + "`json:\"id\"`" + `
	Name  *string ` + "`json:\"name,omitempty\"`" + `
	Color *string ` + "`json:\"color,omitempty\"`" + `
}
`,
		".api-sync/spec-map.json":      specMap,
		".api-sync/unmodeled.json":     unmodeled,
		".api-sync/spec-snapshot.json": string(specJSON),
	})
	return root
}

func writeSpecCurrent(t *testing.T, root string, schemas map[string]string) {
	t.Helper()
	spec := map[string]any{}
	for name, raw := range schemas {
		var v any
		require.NoError(t, json.Unmarshal([]byte(raw), &v))
		spec[name] = v
	}
	specDoc := map[string]any{
		"components": map[string]any{"schemas": spec},
		"paths":      map[string]any{},
	}
	data, err := json.MarshalIndent(specDoc, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(root, ".api-sync", "spec-current.json"), data, 0o644))
}

const baselineWidgetSchema = `{"type":"object","properties":{
	"id":{"type":"string"},
	"name":{"type":["string","null"]},
	"color":{"type":"string","enum":["red","blue"]}
},"required":["id"]}`

func TestIntegration_CheckIsCleanAgainstItsOwnBaseline(t *testing.T) {
	root := integrationFixture(t, map[string]string{"Widget": baselineWidgetSchema})

	quiet, err := reconcileAndMaybeApply(root, syncOptions{Check: true})
	require.NoError(t, err)
	require.True(t, quiet)
}

func TestIntegration_ApplyThenIdempotent(t *testing.T) {
	root := integrationFixture(t, map[string]string{"Widget": baselineWidgetSchema})
	writeSpecCurrent(t, root, map[string]string{
		"Widget": `{"type":"object","properties":{
			"id":{"type":"string"},
			"name":{"type":["string","null"]},
			"color":{"type":"string","enum":["red","blue","green"]},
			"score":{"type":"integer"}
		},"required":["id"]}`,
	})

	quiet, err := reconcileAndMaybeApply(root, syncOptions{Apply: true})
	require.NoError(t, err)
	require.False(t, quiet)

	enumsGo := readFixture(t, root, "internal/types/enums.go")
	require.Contains(t, enumsGo, `ColorGreen Color = "green"`)
	typesGo := readFixture(t, root, "types.go")
	require.Contains(t, typesGo, "ColorGreen = types.ColorGreen")
	widgetGo := readFixture(t, root, "widgets/client.go")
	require.Contains(t, widgetGo, `json:"score,omitempty"`)
	blindpayGo := readFixture(t, root, "blindpay.go")
	require.Contains(t, blindpayGo, `const Version = "1.1.0"`, "enum member change bumps minor")

	// The snapshot must now match spec-current.json's content, so a second
	// -check against the (refreshed) baseline is clean...
	quiet, err = reconcileAndMaybeApply(root, syncOptions{Check: true})
	require.NoError(t, err)
	require.True(t, quiet)

	// ...and a second -apply against the same spec-current.json is a no-op:
	// running the patcher twice must never produce a second change.
	beforeEnums := readFixture(t, root, "internal/types/enums.go")
	beforeTypes := readFixture(t, root, "types.go")
	beforeWidget := readFixture(t, root, "widgets/client.go")
	beforeVersion := readFixture(t, root, "blindpay.go")

	quiet, err = reconcileAndMaybeApply(root, syncOptions{Apply: true})
	require.NoError(t, err)
	require.False(t, quiet) // "nothing to do" is a nil error, not the -check quiet path

	require.Equal(t, beforeEnums, readFixture(t, root, "internal/types/enums.go"))
	require.Equal(t, beforeTypes, readFixture(t, root, "types.go"))
	require.Equal(t, beforeWidget, readFixture(t, root, "widgets/client.go"))
	require.Equal(t, beforeVersion, readFixture(t, root, "blindpay.go"), "idempotent apply must not bump the version again")
}

func TestIntegration_RefreshedSnapshotIsAByteForByteCopyNotAReserialization(t *testing.T) {
	// Deliberately odd formatting (2-space indent, unsorted keys, an escaped
	// unicode char) that json.Marshal would never reproduce on its own, to
	// catch any regression to re-marshaling the parsed document instead of
	// copying the source file's bytes verbatim -- a whole-file diff on
	// every future no-op sync, and a snapshot that stops byte-matching
	// what blindpay-v2 actually ships as spec-current.json.
	oddSpecJSON := "{\n  \"paths\": {},\n  \"components\": {\n    \"schemas\": {\n      \"Widget\": {\n        \"required\": [\"id\"],\n        \"type\": \"object\",\n        \"properties\": {\n          \"name\": {\"type\": [\"string\", \"null\"]},\n          \"id\": {\"type\": \"string\"},\n          \"color\": {\"type\": \"string\", \"enum\": [\"red\", \"blue\"], \"example\": \"caf\\u00e9\"},\n          \"score\": {\"type\": \"integer\"}\n        }\n      }\n    }\n  }\n}"

	root := integrationFixture(t, map[string]string{"Widget": baselineWidgetSchema})
	specCurrentPath := filepath.Join(root, ".api-sync", "spec-current.json")
	require.NoError(t, os.WriteFile(specCurrentPath, []byte(oddSpecJSON), 0o644))

	_, err := reconcileAndMaybeApply(root, syncOptions{Apply: true})
	require.NoError(t, err)

	snapshot, err := os.ReadFile(filepath.Join(root, ".api-sync", "spec-snapshot.json"))
	require.NoError(t, err)
	require.Equal(t, oddSpecJSON, string(snapshot), "refreshSnapshot must copy spec-current.json's bytes verbatim, never re-marshal them")
}

func TestIntegration_FieldOnlyChangeBumpsPatch(t *testing.T) {
	root := integrationFixture(t, map[string]string{"Widget": baselineWidgetSchema})
	writeSpecCurrent(t, root, map[string]string{
		"Widget": `{"type":"object","properties":{
			"id":{"type":"string"},
			"name":{"type":["string","null"]},
			"color":{"type":"string","enum":["red","blue"]},
			"score":{"type":"integer"}
		},"required":["id"]}`,
	})

	_, err := reconcileAndMaybeApply(root, syncOptions{Apply: true})
	require.NoError(t, err)
	require.Contains(t, readFixture(t, root, "blindpay.go"), `const Version = "1.0.1"`)
}

func TestIntegration_OperationSetChangeAloneBumpsMinorEvenForAFieldOnlyAction(t *testing.T) {
	root := integrationFixture(t, map[string]string{"Widget": baselineWidgetSchema})
	spec := map[string]any{
		"components": map[string]any{"schemas": map[string]any{
			"Widget": mustUnmarshal(t, `{"type":"object","properties":{
				"id":{"type":"string"},
				"name":{"type":["string","null"]},
				"color":{"type":"string","enum":["red","blue"]},
				"score":{"type":"integer"}
			},"required":["id"]}`),
		}},
		"paths": map[string]any{
			"/widgets": map[string]any{"get": map[string]any{}},
		},
	}
	data, err := json.MarshalIndent(spec, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(root, ".api-sync", "spec-current.json"), data, 0o644))

	_, err = reconcileAndMaybeApply(root, syncOptions{Apply: true})
	require.NoError(t, err)
	require.Contains(t, readFixture(t, root, "blindpay.go"), `const Version = "1.1.0"`, "an operation-set change alone escalates a field-only change to minor")
}

func TestIntegration_ApplyRefusesToWriteAnythingWhenAnyNeedsHumanIssueExists(t *testing.T) {
	root := integrationFixture(t, map[string]string{"Widget": baselineWidgetSchema})
	writeSpecCurrent(t, root, map[string]string{
		"Widget": `{"type":"object","properties":{
			"id":{"type":"string"},
			"name":{"type":["string","null"]},
			"color":{"type":"string","enum":["red","blue","green"]},
			"score":{"type":"integer"},
			"meta":{"type":"object","properties":{}}
		},"required":["id"]}`,
	})

	beforeEnums := readFixture(t, root, "internal/types/enums.go")
	beforeWidget := readFixture(t, root, "widgets/client.go")
	beforeVersion := readFixture(t, root, "blindpay.go")

	_, err := reconcileAndMaybeApply(root, syncOptions{Apply: true})
	require.Error(t, err)

	require.Equal(t, beforeEnums, readFixture(t, root, "internal/types/enums.go"), "a needs-human issue anywhere must block every applicable change too, atomically")
	require.Equal(t, beforeWidget, readFixture(t, root, "widgets/client.go"))
	require.Equal(t, beforeVersion, readFixture(t, root, "blindpay.go"))
}

func mustUnmarshal(t *testing.T, raw string) any {
	t.Helper()
	var v any
	require.NoError(t, json.Unmarshal([]byte(raw), &v))
	return v
}
