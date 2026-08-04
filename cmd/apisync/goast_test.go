package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindStructShape_PointerSiblingStyle(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"pkg/client.go": `package pkg

type CreateParams struct {
	Name   string  ` + "`json:\"name\"`" + `
	Token  *string ` + "`json:\"token,omitempty\"`" + `
	Amount *int    ` + "`json:\"amount,omitempty\"`" + `
}
`,
	})

	shape, err := findStructShape(root, "pkg/client.go", "CreateParams")
	require.NoError(t, err)
	require.Len(t, shape.Fields, 3)
	require.True(t, shape.preferPointerStyle(), "majority of optional siblings are pointer-typed")
	require.True(t, shape.hasJSONField("token"))
	require.False(t, shape.hasJSONField("missing"))
}

func TestFindStructShape_BareSiblingStyle(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"pkg/client.go": `package pkg

type BankAccount struct {
	ID       string ` + "`json:\"id\"`" + `
	PixKey   string ` + "`json:\"pix_key,omitempty\"`" + `
	SpeiCode string ` + "`json:\"spei_code,omitempty\"`" + `
}
`,
	})

	shape, err := findStructShape(root, "pkg/client.go", "BankAccount")
	require.NoError(t, err)
	require.False(t, shape.preferPointerStyle(), "majority of optional siblings are bare+omitempty")
}

func TestFindStructShape_NoOptionalSiblingsDefaultsToPointer(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"pkg/client.go": `package pkg

type Minimal struct {
	ID string ` + "`json:\"id\"`" + `
}
`,
	})

	shape, err := findStructShape(root, "pkg/client.go", "Minimal")
	require.NoError(t, err)
	require.True(t, shape.preferPointerStyle(), "no optional siblings to learn from: defaults to this repo's general pointer convention")
}

func TestFindStructShape_NotFound(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"pkg/client.go": `package pkg

type Foo struct {
	ID string ` + "`json:\"id\"`" + `
}
`,
	})

	_, err := findStructShape(root, "pkg/client.go", "Bar")
	require.Error(t, err)
}

func TestFindEnumShape(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"internal/types/enums.go": `package types

type Color string

const (
	ColorRed  Color = "red"
	ColorBlue Color = "blue"
)
`,
	})

	shape, err := findEnumShape(root, "internal/types/enums.go", "Color")
	require.NoError(t, err)
	require.True(t, shape.hasMember("red"))
	require.True(t, shape.hasMember("blue"))
	require.False(t, shape.hasMember("green"))
	require.Equal(t, 7, shape.LastMemberLine) // line of "ColorBlue Color = \"blue\""
}

func TestFindReexportShape(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"types.go": `package fixture

import "example.com/fixture/internal/types"

type Color = types.Color

const (
	ColorRed  = types.ColorRed
	ColorBlue = types.ColorBlue
)
`,
	})

	shape, err := findReexportShape(root, "types.go", "Color")
	require.NoError(t, err)
	require.Contains(t, shape.Members, "ColorRed")
	require.Contains(t, shape.Members, "ColorBlue")
}

func TestFindReexportShape_MissingBlockIsAnError(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"types.go": `package fixture

import "example.com/fixture/internal/types"

type BusinessIndustry = types.BusinessIndustry
`,
	})

	_, err := findReexportShape(root, "types.go", "BusinessIndustry")
	require.Error(t, err, "type alias with zero const re-exports must be reported, not silently treated as empty")
}

func TestPascalCase(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"sepa", "Sepa"},
		{"refund_wallet_address", "RefundWalletAddress"},
		{"business_industry", "BusinessIndustry"},
		{"a", "A"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			require.Equal(t, tt.want, pascalCase(tt.in))
		})
	}
}

func TestApplyInsertions_SameLineOrderingAndLineShiftSafety(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"f.go": "line1\nline2\nline3\n",
	})

	err := applyInsertions(root, []insertion{
		{File: "f.go", Line: 1, Text: "insertedAfter1"},
		{File: "f.go", Line: 3, Text: "insertedAfter3-a"},
		{File: "f.go", Line: 3, Text: "insertedAfter3-b"},
	})
	require.NoError(t, err)

	got := readFixture(t, root, "f.go")
	require.Equal(t, "line1\ninsertedAfter1\nline2\nline3\ninsertedAfter3-a\ninsertedAfter3-b\n", got)
}
