// The gethclone binary clones ethereum/go-ethereum Go packages, applying
// semantic patches.
package main

import (
	"context"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/subnet-evm/x/gethclone/astpatch"
	"github.com/spf13/pflag"
	"go.uber.org/multierr"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"

	_ "embed"

	// TODO(arr4n): change to using a git sub-module
	_ "github.com/ethereum/go-ethereum/common"
)

func main() {
	c := config{
		astPatches: make(astpatch.PatchRegistry),
	}

	pflag.StringSliceVar(&c.packages, "packages", []string{"core/vm"}, `Geth packages to clone, with or without "github.com/ethereum/go-ethereum" prefix.`)
	pflag.StringVar(&c.outputGoMod, "output_go_mod", "", "go.mod file of the destination to which geth will be cloned.")
	pflag.Parse()

	log.SetOutput(os.Stderr)
	log.Print("START")
	if err := c.run(context.Background()); err != nil {
		log.Fatal(err)
	}
	log.Print("DONE")
}

type config struct {
	// Externally configurable (e.g. flags)
	packages    []string
	outputGoMod string

	// Internal
	outputModule *modfile.Module
	astPatches   astpatch.PatchRegistry

	processed map[string]bool
}

const geth = "github.com/ethereum/go-ethereum"

func (c *config) run(ctx context.Context) error {
	for i, p := range c.packages {
		if !strings.HasPrefix(p, geth) {
			c.packages[i] = path.Join(geth, p)
		}
	}

	mod, err := parseGoMod(c.outputGoMod)
	fmt.Println(err)
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

	pkgConfig := &packages.Config{
		Context: ctx,
		Mode:    packages.NeedName | packages.NeedCompiledGoFiles,
	}
	pkgs, err := packages.Load(pkgConfig, patterns...)
	if err != nil {
		return fmt.Errorf("packages.Load(..., %q): %v", c.packages, err)
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
func (c *config) parse(ctx context.Context, pkg *packages.Package, fset *token.FileSet) error {
	if len(pkg.Errors) != 0 {
		var err error
		for _, e := range pkg.Errors {
			multierr.AppendInto(&err, e)
		}
		return err
	}

	if c.processed[pkg.PkgPath] {
		return nil
	}
	c.processed[pkg.PkgPath] = true

	outDir := filepath.Join(
		filepath.Dir(c.outputGoMod),
		strings.TrimPrefix(pkg.PkgPath, geth),
	)
	log.Printf("Cloning %q into %q", pkg.PkgPath, outDir)

	allGethImports := set.NewSet[string](0)
	for _, fName := range pkg.CompiledGoFiles {
		file, err := parser.ParseFile(fset, fName, nil, parser.ParseComments|parser.SkipObjectResolution)
		if err != nil {
			return fmt.Errorf("parser.ParseFile(... %q ...): %v", fName, err)
		}

		gethImports, err := c.transformGethImports(fset, file)
		if err != nil {
			return nil
		}
		allGethImports.Union(gethImports)

		file.Comments = append([]*ast.CommentGroup{{
			List: []*ast.Comment{{Text: copyrightHeader}},
		}}, file.Comments...)

		if err := c.astPatches.Apply(pkg.PkgPath, file); err != nil {
			return fmt.Errorf("apply AST patches to %q: %v", pkg.PkgPath, err)
		}

		if false {
			// TODO(arr4n): enable writing when a suitable strategy exists; in
			// the meantime, this code is demonstrative of intent. Do we want to
			// simply overwrite and then check `git diff`? Do we want to check
			// the diffs here?
			outFile := filepath.Join(outDir, filepath.Base(fName))
			f, err := os.Create(outFile)
			if err != nil {
				return err
			}
			if err := format.Node(f, fset, file); err != nil {
				return err
			}
		}
	}

	return c.loadAndParse(ctx, fset, allGethImports.List()...)
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
