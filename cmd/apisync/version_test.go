package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBumpVersion(t *testing.T) {
	tests := []struct {
		name  string
		start string
		class string
		want  string
	}{
		{"minor resets patch", "1.16.0", "minor", "1.17.0"},
		{"minor resets a nonzero patch", "1.16.3", "minor", "1.17.0"},
		{"patch increments", "1.16.0", "patch", "1.16.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := writeFixture(t, map[string]string{
				"blindpay.go": "package fixture\n\nconst Version = \"" + tt.start + "\"\n",
			})
			got, err := bumpVersion(root, tt.class)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
			require.Contains(t, readFixture(t, root, "blindpay.go"), `const Version = "`+tt.want+`"`)
		})
	}
}

func TestBumpVersion_RejectsUnknownClass(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"blindpay.go": "package fixture\n\nconst Version = \"1.0.0\"\n",
	})
	_, err := bumpVersion(root, "major")
	require.Error(t, err, "major is never automatic")
}
