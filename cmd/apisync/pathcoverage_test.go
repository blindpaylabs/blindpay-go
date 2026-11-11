package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const coverageFixtureClient = `package widgets

import (
	"context"
	"net/http"

	"example.com/fixture/internal/request"
)

type Client struct{ instanceID string }

func (c *Client) Get(ctx context.Context, id string) error {
	path := request.Pathf("/instances/%s/widgets/%s", c.instanceID, id)
	_, err := request.Do[struct{}](nil, ctx, "GET", path, nil)
	return err
}

func (c *Client) List(ctx context.Context) error {
	_, err := request.DoList[struct{}](nil, ctx, c.listPath())
	return err
}

func (c *Client) listPath() string { return request.Pathf("/instances/%s/widgets", c.instanceID) }

func (c *Client) Upload(ctx context.Context) error {
	_, err := request.DoMultipart[struct{}](nil, ctx, "/upload?instance_id=x", nil, "")
	return err
}

func (c *Client) Yield(ctx context.Context) error {
	_, err := request.Do[struct{}](nil, ctx, http.MethodPost, request.Pathf("/instances/%s/yield", c.instanceID), nil)
	return err
}
`

func TestComputeCoverageGaps_MethodAware(t *testing.T) {
	repoRoot := writeFixture(t, map[string]string{
		"go.mod":            "module example.com/fixture\n\ngo 1.21\n",
		"widgets/client.go": coverageFixtureClient,
	})

	ops, err := sdkOperations(repoRoot)
	require.NoError(t, err)
	require.Equal(t, map[string]bool{
		"GET /instances/%s/widgets/%s": true,
		"GET /instances/%s/widgets":    true, // path resolved through the listPath helper
		"POST /upload":                 true,
		"POST /instances/%s/yield":     true,
	}, ops)

	sdkOp := map[string]any{"x-sdk": true}
	spec := specWithPaths(t, map[string]string{}, map[string]any{
		"/v1/instances/{instance_id}/widgets/{id}": map[string]any{
			"get":   sdkOp,
			"patch": sdkOp, // same path, different method: must still be a gap
		},
		"/v1/instances/{instance_id}/widgets": map[string]any{"get": sdkOp},
		"/v1/upload":                          map[string]any{"post": sdkOp},
		"/v1/instances/{instance_id}/yield":   map[string]any{"post": map[string]any{"x-sdk": false}},
		"/v1/presign":                         map[string]any{"post": sdkOp},
	})

	require.Equal(t, []string{
		"SDK implements x-sdk: false operation: POST /v1/instances/{instance_id}/yield",
		"spec operation with no SDK method: PATCH /v1/instances/{instance_id}/widgets/{id}",
		"spec operation with no SDK method: POST /v1/presign",
	}, computeCoverageGaps(repoRoot, spec))
}
