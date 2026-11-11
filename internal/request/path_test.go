package request

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPathf(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		segments []string
		want     string
	}{
		{"plain ids", "/instances/%s/payouts/%s", []string{"in_000000000000", "po_000000000000"}, "/instances/in_000000000000/payouts/po_000000000000"},
		{"slash", "/instances/%s/payouts/%s", []string{"in_1", "a/b"}, "/instances/in_1/payouts/a%2Fb"},
		{"dot segments", "/instances/%s/payouts/%s", []string{"in_1", "../customers"}, "/instances/in_1/payouts/..%2Fcustomers"},
		{"query chars", "/x/%s", []string{"a?b#c"}, "/x/a%3Fb%23c"},
		{"space", "/x/%s", []string{"a b"}, "/x/a%20b"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, Pathf(tc.format, tc.segments...))
		})
	}
}
