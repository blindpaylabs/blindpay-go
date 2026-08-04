package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// fanoutFixture builds a minimal repo with one schema ("Variant") mapped to
// N SDK struct sites (variantNames), mirroring the shape of this SDK's real
// multi-site mappings: BankAccountOut fans out to 1 base struct + 11
// per-rail CreateXxxResponse structs (12 total), CreateCustomerIn to 3 KYC/
// KYB variants. Every site starts with only "id"; oldSchema/newSchema (JSON
// Schema property maps, as raw JSON) drive what -apply should add.
func fanoutFixture(t *testing.T, variantNames []string, policy string, propertySites map[string][]string, oldSchema, newSchema string) string {
	t.Helper()

	var structsGo strings.Builder
	structsGo.WriteString("package variants\n\n")
	for _, name := range variantNames {
		structsGo.WriteString(fmt.Sprintf("type %s struct {\n\tID string `json:\"id\"`\n}\n\n", name))
	}

	sdkSites := make([]map[string]any, 0, len(variantNames))
	for _, name := range variantNames {
		sdkSites = append(sdkSites, map[string]any{"file": "variants/client.go", "symbol": name})
	}
	typeMapping := map[string]any{
		"spec": "Variant",
		"sdk":  sdkSites,
	}
	if policy != "" {
		typeMapping["policy"] = policy
	}
	if propertySites != nil {
		typeMapping["property_sites"] = propertySites
	}
	specMap := map[string]any{
		"enums":  []any{},
		"types":  []any{typeMapping},
		"ignore": map[string]any{"schemas": []any{}},
	}
	specMapJSON, err := json.MarshalIndent(specMap, "", "  ")
	require.NoError(t, err)

	var oldSchemaVal any
	require.NoError(t, json.Unmarshal([]byte(oldSchema), &oldSchemaVal))
	oldSpecDoc := map[string]any{
		"components": map[string]any{"schemas": map[string]any{"Variant": oldSchemaVal}},
		"paths":      map[string]any{},
	}
	oldSpecJSON, err := json.MarshalIndent(oldSpecDoc, "", "  ")
	require.NoError(t, err)

	root := writeFixture(t, map[string]string{
		"go.mod":                       "module example.com/fanoutfixture\n\ngo 1.21\n",
		"blindpay.go":                  "package fanoutfixture\n\nconst Version = \"1.0.0\"\n",
		"variants/client.go":           structsGo.String(),
		".api-sync/spec-map.json":      string(specMapJSON),
		".api-sync/unmodeled.json":     `{"properties": [], "enums": []}`,
		".api-sync/spec-snapshot.json": string(oldSpecJSON),
	})

	var newSchemaVal any
	require.NoError(t, json.Unmarshal([]byte(newSchema), &newSchemaVal))
	newSpecDoc := map[string]any{
		"components": map[string]any{"schemas": map[string]any{"Variant": newSchemaVal}},
		"paths":      map[string]any{},
	}
	newSpecJSON, err := json.MarshalIndent(newSpecDoc, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(root, ".api-sync", "spec-current.json"), newSpecJSON, 0o644))

	return root
}

func elevenBankAccountVariants() []string {
	return []string{
		"BankAccount", "CreatePixResponse", "CreateAchResponse", "CreateWireResponse",
		"CreateArgentinaTransfersResponse", "CreateSpeiResponse", "CreateColombiaAchResponse",
		"CreateInternationalSwiftResponse", "CreatePixSafeResponse", "CreateRtpResponse",
		"CreateTedResponse", "CreateSepaResponse",
	}
}

func threeCustomerVariants() []string {
	return []string{
		"CreateIndividualStandardParams", "CreateIndividualEnhancedParams", "CreateBusinessStandardParams",
	}
}

// TestFanout_Uniform_BankAccountShapedSet_AddedFieldLandsOnEverySite is the
// regression test for the reported bug: the patcher used to treat a
// property as satisfied if ANY mapped struct had it, and applied a new
// field only to the FIRST mapped struct. With a 12-site mapping shaped like
// BankAccountOut (1 base + 11 per-rail variants) and default ("uniform")
// policy, a newly added optional property must be inserted into every one
// of the 12 sites, not just the first.
func TestFanout_Uniform_BankAccountShapedSet_AddedFieldLandsOnEverySite(t *testing.T) {
	variants := elevenBankAccountVariants()
	oldSchema := `{"type":"object","properties":{"id":{"type":"string"}},"required":["id"]}`
	newSchema := `{"type":"object","properties":{"id":{"type":"string"},"business_industry":{"type":"string"}},"required":["id"]}`

	root := fanoutFixture(t, variants, "", nil, oldSchema, newSchema)

	quiet, err := reconcileAndMaybeApply(root, syncOptions{Apply: true})
	require.NoError(t, err)
	require.False(t, quiet)

	src := readFixture(t, root, "variants/client.go")
	for _, name := range variants {
		start := strings.Index(src, "type "+name+" struct {")
		require.NotEqual(t, -1, start, "struct %s not found", name)
		end := strings.Index(src[start:], "}")
		block := src[start : start+end]
		require.Contains(t, block, `json:"business_industry,omitempty"`, "site %s must receive the new uniform-policy field", name)
	}

	// Second apply against the same target spec must be a byte-identical
	// no-op: the patcher must never re-derive or duplicate an insertion it
	// already made.
	before := readFixture(t, root, "variants/client.go")
	beforeVersion := readFixture(t, root, "blindpay.go")

	quiet, err = reconcileAndMaybeApply(root, syncOptions{Apply: true})
	require.NoError(t, err)
	require.False(t, quiet)
	require.Equal(t, before, readFixture(t, root, "variants/client.go"), "second apply must be byte-identical (no duplicate insertions)")
	require.Equal(t, beforeVersion, readFixture(t, root, "blindpay.go"), "second apply must not bump the version again")

	quiet, err = reconcileAndMaybeApply(root, syncOptions{Check: true})
	require.NoError(t, err)
	require.True(t, quiet, "-check must be clean once every uniform site has caught up")
}

// TestFanout_Uniform_CustomerKYCShapedSet_AddedFieldLandsOnEverySite is the
// same regression test at the 3-site scale of CreateCustomerIn's
// individual-standard/individual-enhanced/business-standard KYC/KYB
// variants.
func TestFanout_Uniform_CustomerKYCShapedSet_AddedFieldLandsOnEverySite(t *testing.T) {
	variants := threeCustomerVariants()
	oldSchema := `{"type":"object","properties":{"id":{"type":"string"}},"required":["id"]}`
	newSchema := `{"type":"object","properties":{"id":{"type":"string"},"tax_id":{"type":"string"}},"required":["id"]}`

	root := fanoutFixture(t, variants, "", nil, oldSchema, newSchema)

	quiet, err := reconcileAndMaybeApply(root, syncOptions{Apply: true})
	require.NoError(t, err)
	require.False(t, quiet)

	src := readFixture(t, root, "variants/client.go")
	for _, name := range variants {
		start := strings.Index(src, "type "+name+" struct {")
		require.NotEqual(t, -1, start, "struct %s not found", name)
		end := strings.Index(src[start:], "}")
		block := src[start : start+end]
		require.Contains(t, block, `json:"tax_id,omitempty"`, "site %s must receive the new uniform-policy field", name)
	}

	before := readFixture(t, root, "variants/client.go")
	quiet, err = reconcileAndMaybeApply(root, syncOptions{Apply: true})
	require.NoError(t, err)
	require.False(t, quiet)
	require.Equal(t, before, readFixture(t, root, "variants/client.go"), "second apply must be byte-identical (no-op)")
}

// TestFanout_Union_AddedFieldLandsOnlyOnNamedSites verifies the other half
// of the policy: under "union", a new optional property with an explicit
// property_sites entry is added only to the named sites, never to every
// mapped struct (the old shapes[0]-only bug's mirror-image failure mode
// would be "adds to every struct regardless").
func TestFanout_Union_AddedFieldLandsOnlyOnNamedSites(t *testing.T) {
	variants := elevenBankAccountVariants()
	oldSchema := `{"type":"object","properties":{"id":{"type":"string"}},"required":["id"]}`
	newSchema := `{"type":"object","properties":{"id":{"type":"string"},"pix_key":{"type":"string"}},"required":["id"]}`

	root := fanoutFixture(t, variants, "union", map[string][]string{
		"id":      variants,
		"pix_key": {"BankAccount", "CreatePixResponse"},
	}, oldSchema, newSchema)

	_, err := reconcileAndMaybeApply(root, syncOptions{Apply: true})
	require.NoError(t, err)

	src := readFixture(t, root, "variants/client.go")
	for _, name := range variants {
		start := strings.Index(src, "type "+name+" struct {")
		require.NotEqual(t, -1, start)
		end := strings.Index(src[start:], "}")
		block := src[start : start+end]
		want := name == "BankAccount" || name == "CreatePixResponse"
		require.Equal(t, want, strings.Contains(block, `json:"pix_key,omitempty"`), "site %s union-policy field placement", name)
	}
}

// TestFanout_Union_NewPropertyWithoutPropertySitesEntryIsNeedsHuman asserts
// the "never infer" rule: a union-policy mapping must never guess which
// sites a brand new property belongs on just because it is missing
// everywhere (that would be indistinguishable from "belongs nowhere yet,
// decide later"); it must block until a human records an explicit
// property_sites entry (or an unmodeled.json exclusion).
func TestFanout_Union_NewPropertyWithoutPropertySitesEntryIsNeedsHuman(t *testing.T) {
	variants := elevenBankAccountVariants()
	oldSchema := `{"type":"object","properties":{"id":{"type":"string"}},"required":["id"]}`
	newSchema := `{"type":"object","properties":{"id":{"type":"string"},"swift_ifsc_branch_code":{"type":"string"}},"required":["id"]}`

	root := fanoutFixture(t, variants, "union", map[string][]string{
		"id": variants,
	}, oldSchema, newSchema)

	before := readFixture(t, root, "variants/client.go")
	_, err := reconcileAndMaybeApply(root, syncOptions{Apply: true})
	require.Error(t, err, "a union-policy property with no property_sites entry must block, not be silently skipped or added everywhere")
	require.Equal(t, before, readFixture(t, root, "variants/client.go"), "a needs-human issue must block every applicable change, atomically")
}

// TestFanout_Union_PropertySitesNamingAnUnknownSiteIsNeedsHuman guards
// against a typo'd or stale property_sites entry silently doing nothing.
func TestFanout_Union_PropertySitesNamingAnUnknownSiteIsNeedsHuman(t *testing.T) {
	variants := threeCustomerVariants()
	oldSchema := `{"type":"object","properties":{"id":{"type":"string"}},"required":["id"]}`

	root := fanoutFixture(t, variants, "union", map[string][]string{
		"id":     variants,
		"tax_id": {"CreateIndividualStandardParams", "CreateIndividualTypoParams"},
	}, oldSchema, oldSchema)

	require.NoError(t, os.WriteFile(filepath.Join(root, ".api-sync", "spec-current.json"),
		[]byte(readFixture(t, root, ".api-sync/spec-snapshot.json")), 0o644))

	_, err := reconcileAndMaybeApply(root, syncOptions{Check: true})
	require.Error(t, err)
	require.Contains(t, err.Error(), "1 needs-human")
}
