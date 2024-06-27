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

	_ "github.com/ethereum/go-ethereum/common"
)

type config struct {
	// Externally configurable (e.g. flags)
	packages    []string
	outputGoMod string
	runGo       goRunner

	// Internal
	log          *zap.SugaredLogger
	outputModule *modfile.Module
	astPatches   astpatch.PatchRegistry
	patchSets    []patchSet

	processed set.Set[string]
}

// A patchSet registers one or more patches on a `patch.PatchRegistry` and later
// validates that they were correctly applied. Validation is necessary because
// an error-free application of the registry doesn't guarantee that all expected
// nodes were actually visited.
type patchSet interface {
	name() string
	register(astpatch.PatchRegistry)
	validate() error
}

const gethMod = "github.com/ethereum/go-ethereum"

// geth returns `gethMod`+`pkg` unless `pkg` already has `gethMod` as a prefix,
// in which case `pkg` is returned unchanged.
func geth(pkg string) string {
	if strings.HasPrefix(pkg, gethMod) {
		return pkg
	}
	return path.Join(gethMod, strings.TrimLeft(pkg, `/`))
}

func (c *config) run(ctx context.Context, logOpts ...zap.Option) (retErr error) {
	l, err := zap.NewDevelopment(logOpts...)
	if err != nil {
		return err
	}
	c.log = l.Sugar()
	c.runGo.log = c.log
	defer c.log.Sync()

	for i, p := range c.packages {
		c.packages[i] = geth(p)
	}
	for _, ps := range c.patchSets {
		ps.register(c.astPatches)
	}

	mod, err := parseGoMod(c.outputGoMod)
	if err != nil {
		return nil
	}
	c.outputModule = mod.Module

	c.processed = make(set.Set[string])
	if err := c.loadAndParse(ctx, token.NewFileSet(), c.packages...); err != nil {
		return err
	}

	for _, ps := range c.patchSets {
		if err := ps.validate(); err != nil {
			return fmt.Errorf("patch-set %q validation: %v", ps.name(), err)
		}
	}
	return nil
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

	// TODO(arr4n): most of the time is spent here, listing patterns. Although
	// the `processed` set gets rid of most of the duplication, occasionally a
	// package is still `list`ed (but not parse()d) twice. If the overhead
	// becomes problematic, this is where to look first.
	ps := set.Of(patterns...)
	ps.Difference(c.processed)
	pkgs, err := c.runGo.list(ctx, ps.List()...)
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
	if c.processed.Contains(pkg.ImportPath) {
		c.log.Debugf("Already processed %q", pkg.ImportPath)
		return nil
	}
	c.processed.Add(pkg.ImportPath)

	shortPkgPath := strings.TrimPrefix(pkg.ImportPath, gethMod)

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

		if _, err := c.astPatches.Apply(pkg.ImportPath, file); err != nil {
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
		if !strings.HasPrefix(p, gethMod) {
			continue
		}

		imports.Add(p)
		if !astutil.RewriteImport(fset, file, p, strings.Replace(p, gethMod, c.outputModule.Mod.String(), 1)) {
			return nil, fmt.Errorf("failed to rewrite import %q", p)
		}
	}
	return imports, nil
}
