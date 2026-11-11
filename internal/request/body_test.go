package request

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBodyWith(t *testing.T) {
	type params struct {
		ID     string   `json:"-"`
		Name   string   `json:"name"`
		Amount float64  `json:"amount"`
		Note   *string  `json:"note,omitempty"`
		Tags   []string `json:"tags,omitempty"`
	}
	body, err := BodyWith(&params{ID: "x", Name: "n", Amount: 0.1}, map[string]any{"type": "ach"})
	require.NoError(t, err)
	b, err := json.Marshal(body)
	require.NoError(t, err)
	require.JSONEq(t, `{"name":"n","amount":0.1,"type":"ach"}`, string(b))
}
