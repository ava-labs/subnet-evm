package testutils

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMustNotImport fails if a package in plugin/evm imports the testutils package.
func TestMustNotImport(t *testing.T) {
	var evmPackagePaths []string
	filepath.WalkDir("../", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			return nil
		} else if path == "../testutils" {
			return nil
		}
		evmPackagePaths = append(evmPackagePaths, path)
		return nil
	})

	for _, packagePath := range evmPackagePaths {
		packageName := "github.com/ava-labs/subnet-evm/plugin/evm/" + strings.TrimPrefix(packagePath, "../")
		imports, err := getPackageImports(packagePath)
		require.NoError(t, err)
		_, ok := imports["github.com/ava-labs/subnet-evm/plugin/evm/testutils"]
		assert.Falsef(t, ok, "package %s imports testutils: testutils should only be used outside plugin/evm.", packageName)
	}
}

func getPackageImports(packagePath string) (imports map[string]struct{}, err error) {
	imports = make(map[string]struct{})

	err = filepath.Walk(packagePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if path != packagePath && info.IsDir() {
			return fs.SkipDir
		} else if !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}
		err = parseImportsFromFile(path, imports)
		if err != nil {
			return fmt.Errorf("failed to parse imports: %s", err)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk through package files: %v", err)
	}

	return imports, nil
}

func parseImportsFromFile(filePath string, imports map[string]struct{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	node, err := parser.ParseFile(token.NewFileSet(), filePath, file, parser.ImportsOnly)
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %v", filePath, err)
	}

	for _, nodeImport := range node.Imports {
		imports[strings.Trim(nodeImport.Path.Value, `"`)] = struct{}{}
	}

	return nil
}
