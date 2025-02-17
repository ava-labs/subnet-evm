package testutils

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMustNotImport fails if:
// - a package in plugin/evm imports the testutils package
// - a non-test file imports the testutils package
func TestMustNotImport(t *testing.T) {
	prodOnly := false
	checkNotImporting(t, "../", prodOnly, "testutils should only be used outside plugin/evm")
	prodOnly = true
	checkNotImporting(t, getProjectRoot(t), prodOnly, "testutils should only be used in test files")
}

// checkNotImporting checks any package in plugin/evm does not import this package.
func checkNotImporting(t *testing.T, searchPath string, prodOnly bool, message string) {
	var packagePaths []string
	err := filepath.WalkDir(searchPath, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			return nil
		} else if d.Name() == ".git" {
			return fs.SkipDir
		}
		packagePaths = append(packagePaths, path)
		return nil
	})
	require.NoError(t, err)

	projectRootPath := getProjectRoot(t)
	const projectName = "github.com/ava-labs/subnet-evm"
	packageNamePrefix := getPackageNamePrefix(t, projectName, projectRootPath, searchPath)
	testUtilsPackage := getTestUtilsPackage(t, projectName, projectRootPath)

	for _, packagePath := range packagePaths {
		packageName := packageNamePrefix + strings.TrimPrefix(packagePath, searchPath)
		imports, err := getPackageImports(packagePath, prodOnly)
		require.NoError(t, err)
		_, ok := imports[testUtilsPackage]
		assert.Falsef(t, ok, "package %s imports %s: %s", packageName, testUtilsPackage, message)
	}
}

func getPackageImports(packagePath string, prodOnly bool) (imports map[string]struct{}, err error) {
	imports = make(map[string]struct{})

	err = filepath.Walk(packagePath, func(path string, info os.FileInfo, err error) error {
		switch {
		case err != nil:
			return err
		case info.IsDir() && path != packagePath:
			return fs.SkipDir
		case !strings.HasSuffix(info.Name(), ".go"):
			return nil
		case prodOnly && strings.HasSuffix(path, "_test.go"):
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

func getTestUtilsPackage(t *testing.T, projectName, projectRootPath string) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok, "runtimer.Caller failed")
	relativePath := strings.TrimPrefix(filepath.Dir(file), projectRootPath)
	return filepath.Join(projectName, relativePath)
}

func getPackageNamePrefix(t *testing.T, projectName, projectRootPath, searchPath string) string {
	t.Helper()
	searchPath, err := filepath.Abs(searchPath)
	require.NoError(t, err)
	packageNameRoot := strings.TrimPrefix(searchPath, projectRootPath)
	packageNameRoot = strings.TrimPrefix(packageNameRoot, "/")
	return filepath.Join(projectName, packageNameRoot)
}

func getProjectRoot(t *testing.T) (path string) {
	t.Helper()
	path, err := os.Getwd()
	require.NoError(t, err)
	for {
		goModPath := filepath.Join(path, "go.mod")
		_, err := os.Stat(goModPath)
		if err == nil {
			return path
		} else if !os.IsNotExist(err) {
			require.NoError(t, err)
		}
		parentDir := filepath.Dir(path)
		if parentDir == path {
			t.Fatal("go.mod file not found, unable to determine project root")
		}
		path = parentDir
	}
}
