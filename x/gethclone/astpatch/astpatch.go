// Package astpatch provides functionality for traversing and modifying Go
// syntax trees. It extends the astutil package with reusable "patches".
package astpatch

import (
	"go/ast"
	"reflect"

	"golang.org/x/tools/go/ast/astutil"
)

type (
	// A Patch (optionally) modifies an AST node; it is equivalent to an
	// `astutil.ApplyFunc` except that it returns an error instead of a boolean.
	// A non-nil error is equivalent to returning false and will also abort all
	// further calls to other patches.
	Patch func(*astutil.Cursor) error
	// A PatchRegistry maps [Go package path] -> [ast.Node concrete types] ->
	// [all `Patch` functions that must be applied to said node types in said
	// package].
	//
	// The special `pkgPath` value "*" will match all package paths.
	PatchRegistry map[string]map[reflect.Type][]Patch
)

// Add is a convenience wrapper for registering a new `Patch` in the registry.
// The `zeroNode` can be any type (including nil pointers) that implements
// `ast.Node`.
//
// The special `pkgPath` value "*" will match all package paths. While there is
// no specific requirement for `pkgPath` other than it matching the equivalent
// argument passed to `Apply()`, it is typically sourced from
// `golang.org/x/tools/go/packages.Package.PkgPath`.
func (r PatchRegistry) Add(pkgPath string, zeroNode ast.Node, fn Patch) {
	pkg, ok := r[pkgPath]
	if !ok {
		pkg = make(map[reflect.Type][]Patch)
		r[pkgPath] = pkg
	}

	t := nodeType(zeroNode)
	pkg[t] = append(pkg[t], fn)
}

// Apply calls `astutil.Apply()` on `node`, calling the appropriate `Patch`
// functions as the syntax tree is traversed. Patches are applied as the `pre`
// argument to `astutil.Apply()`.
//
// Global `pkgPath` matches (i.e. to those registered with "*") will be applied
// before package-specific matches.
//
// If any `Patch` returns an error then no further patches will be called, and
// the error will be returned by `Apply()`.
func (r PatchRegistry) Apply(pkgPath string, node ast.Node) error {
	var err error
	astutil.Apply(node, func(c *astutil.Cursor) bool {
		if err != nil {
			return false
		}
		if err = r.applyToCursor("*", c); err != nil {
			return false
		}
		err = r.applyToCursor(pkgPath, c)
		return err == nil
	}, nil)
	return err
}

func (r PatchRegistry) applyToCursor(pkgPath string, c *astutil.Cursor) error {
	if c.Node() == nil {
		return nil
	}

	pkg, ok := r[pkgPath]
	if !ok {
		return nil
	}
	for _, fn := range pkg[nodeType(c.Node())] {
		if err := fn(c); err != nil {
			return err
		}
	}
	return nil
}

// nodeType returns the `reflect.Type` of the _concrete_ type implementing
// `ast.Node`. Simpy calling `reflect.TypeOf(n)` would be incorrect as it would
// reflect the interface (and not match any nodes).
func nodeType(n ast.Node) reflect.Type {
	return reflect.ValueOf(n).Type()
}
