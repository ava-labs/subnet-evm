package astpatch

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/ast/astutil"
)

type patchSpy struct {
	visited
}

type visited struct {
	Funcs, Structs, Calls []string
	FuncParams            [][]string
}

// assertEqual asserts that `v == want`, reporting a diff otherwise.
func (v visited) assertEqual(tb testing.TB, want visited) {
	tb.Helper()
	assert.Empty(tb, cmp.Diff(want, v, cmpopts.EquateEmpty()), "visited nodes diff (-want +got)")
}

const errorIfFuncName = "ErrorFuncName"

var errFuncName = fmt.Errorf("encountered sentinel function %q", errorIfFuncName)

func (s *patchSpy) funcRecorder(c *astutil.Cursor) error {
	fn, ok := c.Node().(*ast.FuncDecl)
	if !ok {
		return fmt.Errorf("%T.funcRecorder() called with %T not %T", s, c.Node(), fn)
	}

	name := fn.Name.String()
	if name == errorIfFuncName {
		return errFuncName
	}
	s.Funcs = append(s.Funcs, name)

	var params []string
	for _, p := range fn.Type.Params.List {
		// Params of the same type but different name are grouped together in
		// AST nodes
		for _, n := range p.Names {
			params = append(params, n.Name)
		}
	}
	s.FuncParams = append(s.FuncParams, params)

	return nil
}

func (s *patchSpy) structRecorder(c *astutil.Cursor) error {
	switch p := c.Parent().(type) {
	case *ast.TypeSpec: // it's a `type x struct` not, for example, a `map[T]struct{}`
		s.Structs = append(s.Structs, p.Name.String())
	}
	return nil
}

func (s *patchSpy) funcDeclRecorder(c *astutil.Cursor, fn *ast.FuncDecl) error {
	if !reflect.DeepEqual(c.Node(), fn) {
		return fmt.Errorf("reflect.DeepEqual(%T.Node(), %T) = false; want true", c, fn)
	}
	return s.funcRecorder(c)
}

func (s *patchSpy) callRecorder(c *astutil.Cursor, call *ast.CallExpr) error {
	if !reflect.DeepEqual(c.Node(), call) {
		return fmt.Errorf("reflect.DeepEqual(%T.Node(), %T) = false; want true", c, call)
	}

	var name string
	switch fn := call.Fun.(type) {
	case *ast.Ident:
		name = fn.Name
	default:
		return fmt.Errorf("incomplete test double: %T.callRecorder() called with %T.Fun of unsupported type %T", s, call, call.Fun)
	}
	s.Calls = append(s.Calls, name)

	return nil
}

func TestPatchRegistry(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantErr error
		want    visited
	}{
		{
			name: "happy path",
			src: `package thepackage

func FnA(){}

func FnB(){}

type StructA struct{}

type StructB struct{}
`,
			want: visited{
				Funcs:      []string{"FnA", "FnB"},
				FuncParams: [][]string{{}, {}},
				Structs:    []string{"StructA", "StructB"},
			},
		},
		{
			name: "error propagation",
			src: `package thepackage
			
func HappyFn() {}

func ` + errorIfFuncName + `() {}
`,
			wantErr: errFuncName,
			want: visited{
				Funcs:      []string{"HappyFn"},
				FuncParams: [][]string{{}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var spy patchSpy
			reg := make(PatchRegistry)

			reg.AddForType("*", &ast.FuncDecl{}, spy.funcRecorder)
			const pkgPath = `github.com/the/repo/thepackage`
			reg.AddForType(pkgPath, &ast.StructType{}, spy.structRecorder)

			reg.AddForType("unknown/package/path", &ast.FuncDecl{}, func(c *astutil.Cursor) error {
				t.Errorf("unexpected call to %T with different package path", (Patch)(nil))
				return nil
			})

			file := parseGoFile(t, token.NewFileSet(), tt.src)
			bestEffortLogAST(t, file)

			// None of the `require.Equal*()` variants provide a check for exact
			// match (i.e. equivalent to ==) of the identical error being
			// propagated.
			if _, gotErr := reg.Apply(pkgPath, file); gotErr != tt.wantErr {
				t.Fatalf("%T.Apply(...) got err %v; want %v", reg, gotErr, tt.wantErr)
			}
			spy.visited.assertEqual(t, tt.want)
		})
	}
}

func parseGoFile(tb testing.TB, fset *token.FileSet, src string) *ast.File {
	tb.Helper()
	f, err := parser.ParseFile(fset, "", src, parser.SkipObjectResolution)
	require.NoError(tb, err, "Parsing Go source as file: parser.ParseFile([see logged source])")
	return f
}

func bestEffortLogAST(tb testing.TB, x any) {
	tb.Helper()

	var buf bytes.Buffer
	if err := ast.Fprint(&buf, nil, x, nil); err != nil {
		return
	}
	tb.Logf("AST of parsed source:\n\n%s", buf.String())
}

func TestTypePatchers(t *testing.T) {
	const src = `package box

type (
	TypeA struct {}
	TypeB int
)

func (TypeA) ValueMethod(a int) {}
func (TypeA) valueMethod(a int) {}
func (TypeB) ValueMethod(b int) {}
func (TypeB) valueMethod(b int) {}

func (*TypeA) PointerMethod(a int) {}
func (*TypeA) pointerMethod(a int) {}
func (*TypeB) PointerMethod(b int) {}
func (*TypeB) pointerMethod(b int) {}

func Fn() {}
func fn() {}

func calledFn() {}
func notCalledFn() {}
func init() {
	calledFn()
}
`

	tests := []struct {
		name    string
		patcher func(*patchSpy) TypePatcher
		want    visited
	}{
		{
			name: "method agnostic to pointer/value receiver",
			patcher: func(s *patchSpy) TypePatcher {
				return Method("TypeA", "valueMethod", s.funcDeclRecorder)
			},
			want: visited{
				Funcs:      []string{"valueMethod"},
				FuncParams: [][]string{{"a"}},
			},
		},
		{
			name: "exported method, otherwise same as earlier test",
			patcher: func(s *patchSpy) TypePatcher {
				return Method("TypeA", "ValueMethod", s.funcDeclRecorder)
			},
			want: visited{
				Funcs:      []string{"ValueMethod"},
				FuncParams: [][]string{{"a"}},
			},
		},
		{
			name: "method on different type, otherwise same as earlier test",
			patcher: func(s *patchSpy) TypePatcher {
				return Method("TypeB", "valueMethod", s.funcDeclRecorder)
			},
			want: visited{
				Funcs:      []string{"valueMethod"},
				FuncParams: [][]string{{"b"}},
			},
		},
		{
			name: "PointerMethod() with pointer receiver matches",
			patcher: func(s *patchSpy) TypePatcher {
				return PointerMethod("TypeA", "pointerMethod", s.funcDeclRecorder)
			},
			want: visited{
				Funcs:      []string{"pointerMethod"},
				FuncParams: [][]string{{"a"}},
			},
		},
		{
			name: "PointerMethod() with value receiver ignores",
			patcher: func(s *patchSpy) TypePatcher {
				return PointerMethod("TypeA", "valueMethod", s.funcDeclRecorder)
			},
			want: visited{},
		},
		{
			name: "ValueMethod() with value receiver matches",
			patcher: func(s *patchSpy) TypePatcher {
				return ValueMethod("TypeA", "valueMethod", s.funcDeclRecorder)
			},
			want: visited{
				Funcs:      []string{"valueMethod"},
				FuncParams: [][]string{{"a"}},
			},
		},
		{
			name: "ValueMethod() with pointer receiver ignores",
			patcher: func(s *patchSpy) TypePatcher {
				return ValueMethod("TypeA", "pointerMethod", s.funcDeclRecorder)
			},
			want: visited{},
		},
		{
			name: "function (not method) declaration",
			patcher: func(s *patchSpy) TypePatcher {
				return Function("fn", s.funcDeclRecorder)
			},
			want: visited{
				Funcs:      []string{"fn"},
				FuncParams: [][]string{{}},
			},
		},
		{
			name: "function of different name to earlier test",
			patcher: func(s *patchSpy) TypePatcher {
				return Function("Fn", s.funcDeclRecorder)
			},
			want: visited{
				Funcs:      []string{"Fn"},
				FuncParams: [][]string{{}},
			},
		},
		{
			name: "unqualified function call",
			patcher: func(s *patchSpy) TypePatcher {
				return typePatcher{
					typ:   new(ast.CallExpr),
					patch: UnqualifiedCall("calledFn", s.callRecorder),
				}
			},
			want: visited{Calls: []string{"calledFn"}},
		},
		{
			name: "UnqualifiedCall() for function not actually called",
			patcher: func(s *patchSpy) TypePatcher {
				return typePatcher{
					typ:   new(ast.CallExpr),
					patch: UnqualifiedCall("notCalledFn", s.callRecorder),
				}
			},
			want: visited{Calls: nil},
		},
	}

	file := parseGoFile(t, token.NewFileSet(), src)
	bestEffortLogAST(t, file)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var spy patchSpy
			reg := make(PatchRegistry)
			reg.Add("*", tt.patcher(&spy))

			_, err := reg.Apply("", file)
			require.NoErrorf(t, err, `%T.Apply(...)`, reg)
			spy.visited.assertEqual(t, tt.want)
		})
	}
}

func TestRefactoring(t *testing.T) {
	tests := []struct {
		name, src string
		patcher   TypePatcher
		want      string
	}{
		{
			name: `RenameFunction("foo", "phew")`,
			src: `
package tape

func foo() {}
func bar() {}
`,
			patcher: RenameFunction("foo", "phew"),
			want: `
package tape

func phew() {}
func bar() {}
`,
		},
		{
			name: `RenameFunction("bar", "pub")`,
			src: `
package tape

func foo() {}
func bar() {}
`,
			patcher: RenameFunction("bar", "pub"),
			want: `
package tape

func foo() {}
func pub() {}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := make(PatchRegistry)
			reg.Add("*", tt.patcher)

			fset := token.NewFileSet()

			gotNode, err := reg.Apply("", parseGoFile(t, fset, tt.src))
			require.NoErrorf(t, err, "%T.Apply(...)", reg)

			var got bytes.Buffer
			require.NoErrorf(t, format.Node(&got, fset, gotNode), "format.Node(..., [output of %T.Apply(...)])", reg)

			assert.Equal(t, got.String(), formatGo(t, tt.want), "output of format.Node() after patching")
		})
	}
}

// formatGo parses the file represented by `src` into AST and returns the output
// of `format.Node()`. This allows expected test values to be formatted
// incorrectly without affecting equality checks.
func formatGo(tb testing.TB, src string) string {
	tb.Helper()

	fset := token.NewFileSet()
	file := parseGoFile(tb, fset, src)

	var buf bytes.Buffer
	require.NoError(tb, format.Node(&buf, fset, file), "format.Node(parser.ParseFile(...)) round trip")
	return buf.String()
}
