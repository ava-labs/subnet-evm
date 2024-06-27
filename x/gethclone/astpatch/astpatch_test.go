package astpatch

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/ast/astutil"
)

type patchSpy struct {
	gotFuncs, gotStructs []string
}

const errorIfFuncName = "ErrorFuncName"

var errFuncName = fmt.Errorf("encountered sentinel function %q", errorIfFuncName)

func (s *patchSpy) funcRecorder(c *astutil.Cursor) error {
	name := c.Node().(*ast.FuncDecl).Name.String()
	if name == errorIfFuncName {
		return errFuncName
	}
	s.gotFuncs = append(s.gotFuncs, name)
	return nil
}

func (s *patchSpy) structRecorder(c *astutil.Cursor) error {
	switch p := c.Parent().(type) {
	case *ast.TypeSpec: // it's a `type x struct` not, for example, a `map[T]struct{}`
		s.gotStructs = append(s.gotStructs, p.Name.String())
	}
	return nil
}

func TestPatchRegistry(t *testing.T) {
	tests := []struct {
		name                   string
		src                    string
		wantErr                error
		wantFuncs, wantStructs []string
	}{
		{
			name: "happy path",
			src: `package thepackage

func FnA(){}

func FnB(){}

type StructA struct{}

type StructB struct{}
`,
			wantFuncs:   []string{"FnA", "FnB"},
			wantStructs: []string{"StructA", "StructB"},
		},
		{
			name: "error propagation",
			src: `package thepackage
			
func HappyFn() {}

func ` + errorIfFuncName + `() {}
`,
			wantErr:   errFuncName,
			wantFuncs: []string{"HappyFn"},
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
			assert.Empty(t, cmp.Diff(tt.wantFuncs, spy.gotFuncs), "encountered function declarations (-want +got)")
			assert.Empty(t, cmp.Diff(tt.wantStructs, spy.gotStructs), "encountered struct-type declarations (-want +got)")
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
