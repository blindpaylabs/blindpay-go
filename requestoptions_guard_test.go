package blindpay

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestEveryMethodForwardsRequestOptions is the options-forwarding guard: every
// exported resource method whose first parameter is ctx context.Context must
// end with "opts ...option.RequestOption" and pass "opts..." to every
// transport call (request.Do, DoList, DoMultipart) or delegated method call
// (including the per-page fetch inside ListAutoPaging).
func TestEveryMethodForwardsRequestOptions(t *testing.T) {
	skip := map[string]bool{"cmd": true, "internal": true, "examples": true, "testdata": true, "option": true, "pagination": true}

	entries, err := os.ReadDir(".")
	require.NoError(t, err)

	checked := 0
	for _, e := range entries {
		if !e.IsDir() || skip[e.Name()] || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		files, err := filepath.Glob(filepath.Join(e.Name(), "*.go"))
		require.NoError(t, err)
		for _, file := range files {
			if strings.HasSuffix(file, "_test.go") {
				continue
			}
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, file, nil, 0)
			require.NoError(t, err)

			for _, decl := range f.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Recv == nil || !fn.Name.IsExported() || !takesCtx(fn) {
					continue
				}
				checked++
				where := fset.Position(fn.Pos()).String() + " " + fn.Name.Name

				params := fn.Type.Params.List
				last := params[len(params)-1]
				ell, ok := last.Type.(*ast.Ellipsis)
				require.True(t, ok, "%s: last parameter must be opts ...option.RequestOption", where)
				require.Equal(t, "option.RequestOption", exprString(ell.Elt), "%s: variadic type", where)
				require.Equal(t, "opts", last.Names[0].Name, "%s: variadic name", where)

				calls := 0
				ast.Inspect(fn.Body, func(n ast.Node) bool {
					call, ok := n.(*ast.CallExpr)
					if !ok || !isTransportCall(call) {
						return true
					}
					calls++
					lastArg, _ := call.Args[len(call.Args)-1].(*ast.Ident)
					require.True(t, call.Ellipsis.IsValid() && lastArg != nil && lastArg.Name == "opts",
						"%s: %s must forward opts...", where, exprString(call.Fun))
					return true
				})
				require.Positive(t, calls, "%s: no transport call found", where)
			}
		}
	}
	require.Greater(t, checked, 80, "guard test scanned too few methods")
}

func takesCtx(fn *ast.FuncDecl) bool {
	params := fn.Type.Params.List
	if len(params) == 0 || len(params[0].Names) == 0 {
		return false
	}
	return params[0].Names[0].Name == "ctx" && exprString(params[0].Type) == "context.Context"
}

// isTransportCall matches request.Do* and delegation to another exported
// method on the same receiver (c.ListPage(...)), which must forward opts too.
func isTransportCall(call *ast.CallExpr) bool {
	name := exprString(call.Fun)
	if strings.HasPrefix(name, "request.Do") {
		return true
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || !sel.Sel.IsExported() {
		return false
	}
	recv, ok := sel.X.(*ast.Ident)
	return ok && recv.Name == "c"
}

func exprString(e ast.Expr) string {
	switch v := e.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return exprString(v.X) + "." + v.Sel.Name
	case *ast.IndexExpr:
		return exprString(v.X)
	case *ast.IndexListExpr:
		return exprString(v.X)
	case *ast.StarExpr:
		return "*" + exprString(v.X)
	}
	return ""
}
