package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// realRepoRoot locates the actual checked-out blindpay-go repo (this test
// runs from cmd/apisync, so it's two levels up) -- the golden test needs
// the real partnerfees package and the real committed spec-snapshot.json,
// not a synthetic fixture, to prove the generator against actual SDK
// conventions and an actual spec.
func realRepoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	require.NoError(t, err)
	root := filepath.Join(wd, "..", "..")
	_, err = os.Stat(filepath.Join(root, "go.mod"))
	require.NoError(t, err, "expected a go.mod two levels above cmd/apisync")
	return root
}

// copyTree recursively copies src into dst, skipping .git (large, and
// irrelevant to a build/test run).
func copyTree(t *testing.T, src, dst string) {
	t.Helper()
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if parts := strings.Split(rel, string(filepath.Separator)); len(parts) > 0 && parts[0] == ".git" {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode())
	})
	require.NoError(t, err)
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return string(data)
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
}

// TestGolden_OperationInsertRegeneratesDeletedMethods is the operation-
// insert golden self-test: delete two representative, already-shipped
// operations from the real SDK (partnerfees.Get, a bare GET, and
// partnerfees.Create, a POST with a body) along with the structs/tests
// that belong solely to them, doctor a throwaway copy of the committed
// snapshot to also not have those two operations (simulating "the SDK is
// missing operations the spec already has"), then -apply against the
// real, undoctored snapshot as the target. The generator must reintroduce
// both routes with the right verb/path/params, the result must build,
// vet, and test clean, and a second -apply must be a byte-for-byte no-op.
func TestGolden_OperationInsertRegeneratesDeletedMethods(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go toolchain not on PATH")
	}

	real := realRepoRoot(t)
	work := t.TempDir()
	copyTree(t, real, work)

	clientPath := filepath.Join(work, "partnerfees", "client.go")
	src := readFile(t, clientPath)

	removals := []string{
		`// CreatePartnerFeeParams represents parameters for creating a partner fee.
type CreatePartnerFeeParams struct {
	VirtualAccountSet    *bool   ` + "`json:\"virtual_account_set,omitempty\"`" + `
	EVMWalletAddress     string  ` + "`json:\"evm_wallet_address\"`" + `
	Name                 string  ` + "`json:\"name\"`" + `
	PayinFlatFee         float64 ` + "`json:\"payin_flat_fee\"`" + `
	PayinPercentageFee   float64 ` + "`json:\"payin_percentage_fee\"`" + `
	PayoutFlatFee        float64 ` + "`json:\"payout_flat_fee\"`" + `
	PayoutPercentageFee  float64 ` + "`json:\"payout_percentage_fee\"`" + `
	StellarWalletAddress *string ` + "`json:\"stellar_wallet_address,omitempty\"`" + `
}

`,
		`// Create creates a new partner fee configuration.
func (c *Client) Create(ctx context.Context, params *CreatePartnerFeeParams) (*PartnerFee, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	path := fmt.Sprintf("/instances/%s/partner-fees", c.instanceID)
	return request.Do[*PartnerFee](c.cfg, ctx, "POST", path, params)
}

`,
		`// Get retrieves a specific partner fee by ID.
func (c *Client) Get(ctx context.Context, id string) (*PartnerFee, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	path := fmt.Sprintf("/instances/%s/partner-fees/%s", c.instanceID, id)
	return request.Do[*PartnerFee](c.cfg, ctx, "GET", path, nil)
}

`,
	}
	for _, r := range removals {
		require.Contains(t, src, r, "golden test's removal text is stale against the real partnerfees/client.go source")
		src = strings.Replace(src, r, "", 1)
	}
	writeFile(t, clientPath, src)

	// The Get/Create tests belong solely to the two deleted methods; drop
	// them (List/Delete's tests are untouched).
	testPath := filepath.Join(work, "partnerfees", "partnerfees_test.go")
	writeFile(t, testPath, partnerFeesGoldenTestFile)

	// Doctor a throwaway copy of the committed snapshot: remove the same
	// two operations, so checkOperationChanges sees them as "new" against
	// the real, undoctored snapshot passed as -spec below.
	snapshotPath := filepath.Join(work, ".api-sync", "spec-snapshot.json")
	var spec map[string]any
	require.NoError(t, json.Unmarshal([]byte(readFile(t, snapshotPath)), &spec))
	paths := spec["paths"].(map[string]any)
	collection := paths["/v1/instances/{instance_id}/partner-fees"].(map[string]any)
	byID := paths["/v1/instances/{instance_id}/partner-fees/{id}"].(map[string]any)
	require.Contains(t, collection, "post")
	require.Contains(t, byID, "get")
	delete(collection, "post")
	delete(byID, "get")
	doctored, err := json.MarshalIndent(spec, "", "  ")
	require.NoError(t, err)

	targetSpecPath := filepath.Join(t.TempDir(), "target-spec.json")
	writeFile(t, targetSpecPath, readFile(t, snapshotPath)) // the real, undoctored spec is the apply target
	writeFile(t, snapshotPath, string(doctored))            // the doctored copy becomes the committed baseline

	// A pending, un-applied operation-insert must fail -check.
	_, checkErr := reconcileAndMaybeApply(work, syncOptions{Check: true, SpecPath: targetSpecPath})
	require.Error(t, checkErr, "an un-applied operation-insert must fail -check")

	quiet, applyErr := reconcileAndMaybeApply(work, syncOptions{Apply: true, SpecPath: targetSpecPath})
	require.NoError(t, applyErr)
	require.False(t, quiet)

	regenerated := readFile(t, clientPath)
	require.Contains(t, regenerated, `func (c *Client) Get(ctx context.Context, id string)`)
	require.Contains(t, regenerated, `fmt.Sprintf("/instances/%s/partner-fees/%s", c.instanceID, id)`)
	require.Contains(t, regenerated, `request.Do[*GetResponse](c.cfg, ctx, "GET", path, nil)`)
	require.Contains(t, regenerated, `func (c *Client) Create(ctx context.Context, params *CreateParams)`)
	require.Contains(t, regenerated, `fmt.Sprintf("/instances/%s/partner-fees", c.instanceID)`)
	require.Contains(t, regenerated, `request.Do[*CreateResponse](c.cfg, ctx, "POST", path, params)`)

	blindpayGo := readFile(t, filepath.Join(work, "blindpay.go"))
	require.Contains(t, blindpayGo, `const Version = "1.19.0"`, "operationSetChanged escalates the bump to minor")

	run := func(name string, args ...string) string {
		cmd := exec.Command(name, args...)
		cmd.Dir = work
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "%s %v failed:\n%s", name, args, string(out))
		return string(out)
	}
	run("go", "build", "./...")
	run("go", "vet", "./...")
	// go/vet/test the SDK packages this change actually touches, never
	// cmd/apisync itself: `work` is a full copy of this very repo, so
	// running its own apisync test suite in-place would recursively
	// re-run (and re-copy) this exact golden test against an already-
	// mutated tree.
	pkgList := run("go", "list", "./...")
	var pkgs []string
	for _, p := range strings.Split(strings.TrimSpace(pkgList), "\n") {
		if p == "" || strings.Contains(p, "/cmd/") {
			continue
		}
		pkgs = append(pkgs, p)
	}
	run("go", append([]string{"test"}, pkgs...)...)

	// Second apply against the same target must be a byte-for-byte no-op.
	before := readFile(t, clientPath)
	beforeMap := readFile(t, filepath.Join(work, ".api-sync", "spec-map.json"))
	beforeVersion := readFile(t, filepath.Join(work, "blindpay.go"))

	quiet, applyErr = reconcileAndMaybeApply(work, syncOptions{Apply: true, SpecPath: targetSpecPath})
	require.NoError(t, applyErr)
	require.False(t, quiet) // "nothing to do" is a nil error, not the -check quiet path

	require.Equal(t, before, readFile(t, clientPath), "second -apply must not touch partnerfees/client.go again")
	require.Equal(t, beforeMap, readFile(t, filepath.Join(work, ".api-sync", "spec-map.json")))
	require.Equal(t, beforeVersion, readFile(t, filepath.Join(work, "blindpay.go")), "idempotent apply must not bump the version again")

	quiet, checkErr = reconcileAndMaybeApply(work, syncOptions{Check: true, SpecPath: targetSpecPath})
	require.NoError(t, checkErr)
	require.True(t, quiet, "-check must now be clean against the (refreshed) snapshot")
}

// TestGolden_MultipartOperationIsNeedsHuman is the golden test's other
// half: a synthetic operation on the real partnerfees package's own path
// prefix, but with a multipart/form-data body -- non-JSON, so the
// generator must refuse it and leave it to a human, never guess a shape.
func TestGolden_MultipartOperationIsNeedsHuman(t *testing.T) {
	real := realRepoRoot(t)
	work := t.TempDir()
	copyTree(t, real, work)

	sm, err := loadSpecMap(work)
	require.NoError(t, err)

	snapshotPath := filepath.Join(work, ".api-sync", "spec-snapshot.json")
	var spec map[string]any
	require.NoError(t, json.Unmarshal([]byte(readFile(t, snapshotPath)), &spec))
	oldSpecDoc := loadSpecDoc(spec)

	newSpecRaw := map[string]any{}
	require.NoError(t, json.Unmarshal([]byte(readFile(t, snapshotPath)), &newSpecRaw))
	newPaths := newSpecRaw["paths"].(map[string]any)
	newPaths["/v1/instances/{instance_id}/partner-fees/{id}/logo"] = map[string]any{
		"post": map[string]any{
			"tags": []any{"Partner Fees"},
			"requestBody": map[string]any{"content": map[string]any{"multipart/form-data": map[string]any{
				"schema": map[string]any{"type": "object", "properties": map[string]any{
					"file": map[string]any{"type": "string", "format": "binary"},
				}},
			}}},
			"responses": map[string]any{"204": map[string]any{}},
		},
	}
	newSpecDoc := loadSpecDoc(newSpecRaw)

	plan := checkOperationChanges(work, sm, oldSpecDoc, newSpecDoc)
	require.Empty(t, plan.Actions)
	require.Len(t, plan.NeedsHuman, 1)
	require.Contains(t, plan.NeedsHuman[0], "NEEDS_HUMAN")
	require.Contains(t, plan.NeedsHuman[0], "POST /v1/instances/{instance_id}/partner-fees/{id}/logo")
	require.Contains(t, plan.NeedsHuman[0], "non-JSON content type")
	require.Contains(t, plan.NeedsHuman[0], "multipart/form-data")
}

const partnerFeesGoldenTestFile = `package partnerfees

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/stretchr/testify/require"
)

func TestPartnerFees_List(t *testing.T) {
	instanceID := "in_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(` + "`" + `[
					{
						"id":"fe_000000000000",
						"instance_id":"in_000000000000",
						"name":"Display Name",
						"payout_percentage_fee":0,
						"payout_flat_fee":0,
						"payin_percentage_fee":0,
						"payin_flat_fee":0,
						"evm_wallet_address":"0x1234567890123456789012345678901234567890",
						"stellar_wallet_address":"GAB22222222222222222222222222222222222222222222222222222222222222"
					}
				]` + "`" + `),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/partner-fees", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	fees, err := client.List(context.Background())
	require.NoError(t, err)
	require.Len(t, fees, 1)
	require.Equal(t, "fe_000000000000", fees[0].ID)
	require.Equal(t, "Display Name", fees[0].Name)
}

func TestPartnerFees_Delete(t *testing.T) {
	instanceID := "in_000000000000"
	id := "fe_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(` + "`" + `{"data":null}` + "`" + `),
				Method: http.MethodDelete,
				Path:   fmt.Sprintf("/instances/%s/partner-fees/%s", instanceID, id),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	err := client.Delete(context.Background(), id)
	require.NoError(t, err)
}
`
