package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// gadgetsFixture builds a tiny synthetic SDK package ("gadgets") with one
// existing List method, for classifier unit tests that don't need the
// full real repo.
func gadgetsFixture(t *testing.T) string {
	t.Helper()
	return writeFixture(t, map[string]string{
		"go.mod":                  "module example.com/fixture\n\ngo 1.21\n",
		".api-sync/spec-map.json": "{\n  \"enums\": [],\n  \"types\": [\n    { \"spec\": \"Placeholder\", \"sdk\": [{ \"file\": \"gadgets/client.go\", \"symbol\": \"Gadget\" }] }\n  ]\n}\n",
		"gadgets/client.go": `package gadgets

import (
	"context"
	"fmt"

	"github.com/blindpaylabs/blindpay-go/internal/request"
)

type Client struct {
	cfg        *request.Config
	instanceID string
}

func (c *Client) List(ctx context.Context) ([]Gadget, error) {
	path := fmt.Sprintf("/instances/%s/gadgets", c.instanceID)
	return request.Do[[]Gadget](c.cfg, ctx, "GET", path, nil)
}

type Gadget struct {
	ID   string ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}
`,
	})
}

func gadgetOp(method string, requestBody, response map[string]any) operationEntry {
	op := map[string]any{"tags": []any{"Gadgets"}}
	if requestBody != nil {
		op["requestBody"] = requestBody
	}
	if response != nil {
		op["responses"] = map[string]any{"200": response}
	} else {
		op["responses"] = map[string]any{"204": map[string]any{}}
	}
	return operationEntry{Method: method, Path: "/v1/instances/{instance_id}/gadgets/{id}/activate", Tags: []string{"Gadgets"}, Op: op}
}

func jsonContent(schema map[string]any) map[string]any {
	return map[string]any{"content": map[string]any{"application/json": map[string]any{"schema": schema}}}
}

func TestClassifyAndPlanOperation_InlineFlatBodyIsStandard(t *testing.T) {
	root := gadgetsFixture(t)
	spec := mustSpec(t, map[string]string{})

	op := gadgetOp("POST",
		jsonContent(map[string]any{"type": "object", "properties": map[string]any{"reason": map[string]any{"type": "string"}}, "required": []any{"reason"}}),
		map[string]any{"content": map[string]any{"application/json": map[string]any{
			"schema": map[string]any{"type": "object", "properties": map[string]any{"ok": map[string]any{"type": "boolean"}}, "required": []any{"ok"}},
		}}},
	)

	acts, pending, reason := classifyAndPlanOperation(root, &SpecMap{}, spec, op)
	require.Empty(t, reason)
	require.Empty(t, pending) // inline shapes are never spec-map registered
	require.Len(t, acts, 1)
	require.Contains(t, acts[0].insert.Text, "func (c *Client) Activate(")
	require.Contains(t, acts[0].insert.Text, `fmt.Sprintf("/instances/%s/gadgets/%s/activate", c.instanceID, id)`)
	require.Contains(t, acts[0].insert.Text, `request.Do[*ActivateResponse](c.cfg, ctx, "POST", path, params)`)
	require.Contains(t, acts[0].insert.Text, "type ActivateParams struct")
	require.Contains(t, acts[0].insert.Text, "Reason string `json:\"reason\"`")
	require.Contains(t, acts[0].insert.Text, "type ActivateResponse struct")
}

func TestClassifyAndPlanOperation_MultipartRequestIsNonStandard(t *testing.T) {
	root := gadgetsFixture(t)
	spec := mustSpec(t, map[string]string{})

	op := gadgetOp("POST",
		map[string]any{"content": map[string]any{"multipart/form-data": map[string]any{
			"schema": map[string]any{"type": "object", "properties": map[string]any{"file": map[string]any{"type": "string", "format": "binary"}}},
		}}},
		nil,
	)

	acts, _, reason := classifyAndPlanOperation(root, &SpecMap{}, spec, op)
	require.Nil(t, acts)
	require.Contains(t, reason, "non-JSON content type")
	require.Contains(t, reason, "multipart/form-data")
}

func TestClassifyAndPlanOperation_NonFlatInlineBodyIsNonStandard(t *testing.T) {
	root := gadgetsFixture(t)
	spec := mustSpec(t, map[string]string{})

	op := gadgetOp("POST",
		jsonContent(map[string]any{"type": "object", "properties": map[string]any{
			"nested": map[string]any{"type": "object", "properties": map[string]any{"x": map[string]any{"type": "string"}}},
		}}),
		nil,
	)

	acts, _, reason := classifyAndPlanOperation(root, &SpecMap{}, spec, op)
	require.Nil(t, acts)
	require.Contains(t, reason, "non-scalar property")
}

func TestClassifyAndPlanOperation_NewNamedRefSchemaIsSynthesizedAndRegistered(t *testing.T) {
	root := gadgetsFixture(t)
	spec := mustSpec(t, map[string]string{
		"GadgetStatusOut": `{"type":"object","properties":{"active":{"type":"boolean"}},"required":["active"]}`,
	})

	op := gadgetOp("POST", nil, map[string]any{"content": map[string]any{"application/json": map[string]any{
		"schema": map[string]any{"$ref": "#/components/schemas/GadgetStatusOut"},
	}}})

	acts, pending, reason := classifyAndPlanOperation(root, &SpecMap{}, spec, op)
	require.Empty(t, reason)
	require.Equal(t, []string{"GadgetStatusOut"}, pending)
	require.Len(t, acts, 2) // method+struct, and the spec-map.json registration
	require.Contains(t, acts[1].insert.Text, `"spec": "GadgetStatusOut"`)
	require.Contains(t, acts[1].insert.Text, `"symbol": "GadgetStatusOut"`)
}

func TestClassifyAndPlanOperation_MappedRefInAnotherPackageIsNonStandard(t *testing.T) {
	root := gadgetsFixture(t)
	spec := mustSpec(t, map[string]string{
		"GadgetStatusOut": `{"type":"object","properties":{"active":{"type":"boolean"}},"required":["active"]}`,
	})
	sm := &SpecMap{Types: []TypeMapping{{Spec: "GadgetStatusOut", SDK: []SDKSite{{File: "widgets/client.go", Symbol: "Widget"}}}}}

	op := gadgetOp("POST", nil, map[string]any{"content": map[string]any{"application/json": map[string]any{
		"schema": map[string]any{"$ref": "#/components/schemas/GadgetStatusOut"},
	}}})

	acts, _, reason := classifyAndPlanOperation(root, sm, spec, op)
	require.Nil(t, acts)
	require.Contains(t, reason, "cross-package type reuse is not supported")
}

func TestClassifyAndPlanOperation_MethodNameCollisionIsNonStandard(t *testing.T) {
	root := gadgetsFixture(t)
	spec := mustSpec(t, map[string]string{})

	op := operationEntry{
		Method: "GET",
		Path:   "/v1/instances/{instance_id}/gadgets",
		Tags:   []string{"Gadgets"},
		Op:     map[string]any{"tags": []any{"Gadgets"}, "responses": map[string]any{"204": map[string]any{}}},
	}

	acts, _, reason := classifyAndPlanOperation(root, &SpecMap{}, spec, op)
	require.Nil(t, acts)
	require.Contains(t, reason, `"List" already exists`)
}

func TestDeriveMethodName(t *testing.T) {
	cases := []struct {
		method string
		tail   []string
		want   string
	}{
		{"GET", nil, "List"},
		{"GET", []string{"{id}"}, "Get"},
		{"GET", []string{"{id}", "balance"}, "GetBalance"},
		{"POST", nil, "Create"},
		{"POST", []string{"{id}", "activate"}, "Activate"},
		{"PUT", nil, "Update"},
		{"DELETE", []string{"{id}"}, "Delete"},
	}
	for _, tc := range cases {
		require.Equal(t, tc.want, deriveMethodName(tc.method, tc.tail), "%s %v", tc.method, tc.tail)
	}
}

func TestCamelParam(t *testing.T) {
	require.Equal(t, "id", camelParam("id"))
	require.Equal(t, "customerID", camelParam("customer_id"))
	require.Equal(t, "userID", camelParam("user_id"))
	require.Equal(t, "swift", camelParam("swift"))
}
