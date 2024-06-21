package main

import (
	"context"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/subnet-evm/x/gethclone/astpatch"
	"go.uber.org/zap"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/ast/astutil"

	_ "embed"

	// TODO(arr4n): change to using a git sub-module
	_ "github.com/ethereum/go-ethereum/common"
)

type config struct {
	// Externally configurable (e.g. flags)
	packages    []string
	outputGoMod string
	goBinary    string

	// Internal
	log          *zap.SugaredLogger
	outputModule *modfile.Module
	astPatches   astpatch.PatchRegistry

	processed map[string]bool
}

const geth = "github.com/ethereum/go-ethereum"

func (c *config) run(ctx context.Context, logOpts ...zap.Option) (retErr error) {
	l, err := zap.NewDevelopment(logOpts...)
	if err != nil {
		return err
	}
	c.log = l.Sugar()
	defer c.log.Sync()

	for i, p := range c.packages {
		if !strings.HasPrefix(p, geth) {
			c.packages[i] = path.Join(geth, p)
		}
	}

	mod, err := parseGoMod(c.outputGoMod)
	if err != nil {
		return nil
	}
	c.outputModule = mod.Module

	c.processed = make(map[string]bool)
	return c.loadAndParse(ctx, token.NewFileSet(), c.packages...)
}

func parseGoMod(filePath string) (*modfile.File, error) {
	buf, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("parsing output go.mod: os.ReadFile(%q): %v", filePath, err)
	}
	return modfile.ParseLax(filePath, buf, nil)
}

// loadAndParse loads all packages that match the `patterns` and individually
// passes them to `c.parse()`.
func (c *config) loadAndParse(ctx context.Context, fset *token.FileSet, patterns ...string) error {
	if len(patterns) == 0 {
		return nil
	}

	pkgs, err := c.goList(ctx, patterns...)
	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		if err := c.parse(ctx, pkg, fset); err != nil {
			return err
		}
	}
	return nil
}

//go:embed copyright.go.txt
var copyrightHeader string

// parse parses all `pkg.Files` into `fset`, transforms each according to
// semantic patches, and passes all geth imports back to `c.loadAndParse()` for
// recursive handling.
func (c *config) parse(ctx context.Context, pkg *PackagePublic, fset *token.FileSet) error {
	if c.processed[pkg.ImportPath] {
		c.log.Debugf("Already processed %q", pkg.ImportPath)
		return nil
	}
	c.processed[pkg.ImportPath] = true

	shortPkgPath := strings.TrimPrefix(pkg.ImportPath, geth)

	outDir := filepath.Join(filepath.Dir(c.outputGoMod), shortPkgPath)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("create directory for %q: %v", shortPkgPath, err)
	}
	c.log.Infof("Cloning %q into %q", pkg.ImportPath, outDir)

	allGethImports := set.NewSet[string](0)
	for _, fName := range concat(pkg.GoFiles, pkg.TestGoFiles) {
		fPath := filepath.Join(pkg.Dir, fName)
		file, err := parser.ParseFile(fset, fPath, nil, parser.ParseComments|parser.SkipObjectResolution)
		if err != nil {
			return fmt.Errorf("parser.ParseFile(... %q ...): %v", fPath, err)
		}

		gethImports, err := c.transformGethImports(fset, file)
		if err != nil {
			return nil
		}
		allGethImports.Union(gethImports)

		file.Comments = append([]*ast.CommentGroup{{
			List: []*ast.Comment{{Text: copyrightHeader}},
		}}, file.Comments...)

		if err := c.astPatches.Apply(pkg.ImportPath, file); err != nil {
			return fmt.Errorf("apply AST patches to %q: %v", pkg.ImportPath, err)
		}

		outFile := fmt.Sprintf("%s.gethclone", filepath.Join(outDir, filepath.Base(fName)))
		f, err := os.Create(outFile)
		if err != nil {
			return err
		}
		if err := format.Node(f, fset, file); err != nil {
			return fmt.Errorf("format.Node(..., %T): %v", file, err)
		}
		c.log.Infof("Cloned %q", filepath.Join(shortPkgPath, fName))
	}

	return c.loadAndParse(ctx, fset, allGethImports.List()...)
}

func concat(strs ...[]string) []string {
	var out []string
	for _, s := range strs {
		out = append(out, s...)
	}
	return out
}

// transformGethImports finds all `ethereum/go-ethereum` imports in the file,
// converts their path to `c.outputModule`, and returns the set of transformed
// import paths.
func (c *config) transformGethImports(fset *token.FileSet, file *ast.File) (set.Set[string], error) {
	imports := set.NewSet[string](len(file.Imports))
	for _, im := range file.Imports {
		p := strings.Trim(im.Path.Value, `"`)
		if !strings.HasPrefix(p, geth) {
			continue
		}

		imports.Add(p)
		if !astutil.RewriteImport(fset, file, p, strings.Replace(p, geth, c.outputModule.Mod.String(), 1)) {
			return nil, fmt.Errorf("failed to rewrite import %q", p)
		}
	}
	return imports, nil
}
