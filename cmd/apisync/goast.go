package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// FieldShape is one json-tagged field already declared on a struct.
type FieldShape struct {
	GoName    string
	JSONName  string
	Omitempty bool
	Pointer   bool
}

// StructShape locates a struct type declaration precisely enough to splice
// a new field into it.
type StructShape struct {
	File             string
	Symbol           string
	Fields           []FieldShape
	ClosingBraceLine int // 1-indexed line of the struct's closing "}"
	LastFieldLine    int // 1-indexed line of the last existing field; insert after this
	Indent           string
}

func findStructShape(repoRoot, file, symbol string) (*StructShape, error) {
	path := filepath.Join(repoRoot, file)
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", file, err)
	}

	for _, decl := range astFile.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != symbol {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok || structType.Fields == nil {
				return nil, fmt.Errorf("%s: %s is not a struct type", file, symbol)
			}
			shape := &StructShape{
				File:             file,
				Symbol:           symbol,
				ClosingBraceLine: fset.Position(structType.Fields.Closing).Line,
				Indent:           "\t",
			}
			for _, f := range structType.Fields.List {
				if line := fset.Position(f.End()).Line; line > shape.LastFieldLine {
					shape.LastFieldLine = line
				}
				if len(f.Names) == 0 || f.Tag == nil {
					continue
				}
				tagVal, err := strconv.Unquote(f.Tag.Value)
				if err != nil {
					continue
				}
				jsonName, omitempty := parseJSONTag(tagVal)
				if jsonName == "" || jsonName == "-" {
					continue
				}
				_, isPointer := f.Type.(*ast.StarExpr)
				shape.Fields = append(shape.Fields, FieldShape{
					GoName:    f.Names[0].Name,
					JSONName:  jsonName,
					Omitempty: omitempty,
					Pointer:   isPointer,
				})
			}
			if shape.LastFieldLine == 0 {
				shape.LastFieldLine = shape.ClosingBraceLine - 1
			}
			return shape, nil
		}
	}
	return nil, fmt.Errorf("%s: struct %s not found", file, symbol)
}

func (s *StructShape) hasJSONField(name string) bool {
	for _, f := range s.Fields {
		if f.JSONName == name {
			return true
		}
	}
	return false
}

// preferPointerStyle inspects the struct's existing optional (omitempty)
// siblings and returns true if the majority use a pointer type, matching
// "the sibling style within each struct" rather than a repo-wide default.
// Defaults to pointer style when there are no optional siblings to learn
// from, matching this repo's general convention (CLAUDE.md section 2).
func (s *StructShape) preferPointerStyle() bool {
	pointerCount, bareCount := 0, 0
	for _, f := range s.Fields {
		if !f.Omitempty {
			continue
		}
		if f.Pointer {
			pointerCount++
		} else {
			bareCount++
		}
	}
	if pointerCount == 0 && bareCount == 0 {
		return true
	}
	return pointerCount >= bareCount
}

// EnumShape locates a set of `X Symbol = "value"` (typed, internal enum) or
// `X = types.SymbolX` (untyped alias, root re-export) const declarations.
type EnumShape struct {
	File           string
	Symbol         string
	Members        map[string]string // spec value -> Go const name
	LastMemberLine int               // 1-indexed line to insert after
	Indent         string
}

func (e *EnumShape) hasMember(value string) bool {
	_, ok := e.Members[value]
	return ok
}

// findEnumShape locates the typed internal const block for symbol (each
// line explicitly repeats "Symbol" as the const type in this repo's style).
func findEnumShape(repoRoot, file, symbol string) (*EnumShape, error) {
	path := filepath.Join(repoRoot, file)
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", file, err)
	}

	shape := &EnumShape{File: file, Symbol: symbol, Members: map[string]string{}, Indent: "\t"}
	found := false

	for _, decl := range astFile.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}
		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok || len(valueSpec.Names) == 0 {
				continue
			}
			ident, ok := valueSpec.Type.(*ast.Ident)
			if !ok || ident.Name != symbol {
				continue
			}
			if len(valueSpec.Values) == 0 {
				continue
			}
			lit, ok := valueSpec.Values[0].(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				continue
			}
			val, err := strconv.Unquote(lit.Value)
			if err != nil {
				continue
			}
			shape.Members[val] = valueSpec.Names[0].Name
			found = true
			line := fset.Position(valueSpec.End()).Line
			if line > shape.LastMemberLine {
				shape.LastMemberLine = line
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("%s: no const block typed %s found", file, symbol)
	}
	return shape, nil
}

// findReexportShape locates root types.go's untyped alias block for an
// internal/types enum: lines of the form `X = types.X` where X shares the
// given symbol as its prefix.
func findReexportShape(repoRoot, file, symbol string) (*EnumShape, error) {
	path := filepath.Join(repoRoot, file)
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", file, err)
	}

	shape := &EnumShape{File: file, Symbol: symbol, Members: map[string]string{}, Indent: "\t"}
	found := false

	for _, decl := range astFile.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}
		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok || len(valueSpec.Names) == 0 || len(valueSpec.Values) == 0 {
				continue
			}
			sel, ok := valueSpec.Values[0].(*ast.SelectorExpr)
			if !ok {
				continue
			}
			pkgIdent, ok := sel.X.(*ast.Ident)
			if !ok || pkgIdent.Name != "types" {
				continue
			}
			if !strings.HasPrefix(sel.Sel.Name, symbol) {
				continue
			}
			found = true
			// Members here are keyed by the exported Go name suffix, not a
			// spec value (the re-export block has no wire values of its own).
			shape.Members[sel.Sel.Name] = valueSpec.Names[0].Name
			line := fset.Position(valueSpec.End()).Line
			if line > shape.LastMemberLine {
				shape.LastMemberLine = line
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("%s: no re-export const block for %s found", file, symbol)
	}
	return shape, nil
}

// parseJSONTag extracts (name, omitempty) from a `json:"name,omitempty"` tag.
func parseJSONTag(tag string) (string, bool) {
	const key = `json:"`
	idx := strings.Index(tag, key)
	if idx == -1 {
		return "", false
	}
	rest := tag[idx+len(key):]
	end := strings.Index(rest, `"`)
	if end == -1 {
		return "", false
	}
	value := rest[:end]
	omitempty := strings.Contains(value, ",omitempty")
	if comma := strings.Index(value, ","); comma != -1 {
		value = value[:comma]
	}
	return value, omitempty
}

// pascalCase converts a spec name into a PascalCase Go identifier, splitting
// on any run of non-alphanumeric characters -- snake_case property names
// ("refund_wallet_address") and dotted enum values alike ("customer.update",
// matching this repo's existing WebhookEventCustomerUpdate = "customer.update"
// naming). Purely mechanical (no acronym table): a mechanically produced
// name is always valid Go and additive; it just may not match a
// hand-tuned identifier's casing for acronyms (e.g. produces "Uetr" rather
// than "UETR"). Acceptable for an auto-applied, reviewable PR.
func pascalCase(name string) string {
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9')
	})
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]))
		if len(p) > 1 {
			b.WriteString(p[1:])
		}
	}
	return b.String()
}

// insertion is one planned text splice: insert Text immediately after
// 1-indexed line Line of File (paths relative to the repo root).
type insertion struct {
	File string
	Line int
	Text string
}

// applyInsertions splices every planned insertion into its file. Multiple
// insertions targeting the same (file, line) are joined in the order given
// (callers must pre-sort for determinism); insertions are applied from the
// bottom of each file upward so an earlier splice never shifts the line
// number of one still to be applied.
func applyInsertions(repoRoot string, insertions []insertion) error {
	byFile := map[string][]insertion{}
	for _, ins := range insertions {
		byFile[ins.File] = append(byFile[ins.File], ins)
	}

	for file, list := range byFile {
		path := filepath.Join(repoRoot, file)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		lines := strings.Split(string(data), "\n")

		byLine := map[int][]string{}
		var lineNums []int
		for _, ins := range list {
			if _, seen := byLine[ins.Line]; !seen {
				lineNums = append(lineNums, ins.Line)
			}
			byLine[ins.Line] = append(byLine[ins.Line], ins.Text)
		}
		sortIntsDesc(lineNums)

		for _, line := range lineNums {
			if line < 1 || line > len(lines) {
				return fmt.Errorf("%s: line %d out of range (file has %d lines)", path, line, len(lines))
			}
			var out []string
			out = append(out, lines[:line]...)
			out = append(out, byLine[line]...)
			out = append(out, lines[line:]...)
			lines = out
		}

		if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func sortIntsDesc(v []int) {
	for i := 1; i < len(v); i++ {
		for j := i; j > 0 && v[j-1] < v[j]; j-- {
			v[j-1], v[j] = v[j], v[j-1]
		}
	}
}
