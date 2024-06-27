package astpatch

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/ast/astutil"
)

// Method returns a `TypePatcher` that only applies to the specific method on
// the specific type.
//
// The `patch` argument functions like a regular `Patch` except that its
// parameters are extended to also accept the methods's AST declaration as its
// concrete type (i.e. `astutil.Cursor.Node().(*ast.FuncDecl)`).
//
//	// Method declaration
//	func (x *Thing) Do() { ... }
//	// Patched with
//	astpatch.Method("Thing", "Do", ...)
func Method(receiverType, methodName string, patch func(*astutil.Cursor, *ast.FuncDecl) error) TypePatcher {
	return method(nil, receiverType, methodName, patch)
}

// PointerMethod is identical to `Method()` except that it only matches methods
// with pointer receivers.
func PointerMethod(receiverType, methodName string, patch func(*astutil.Cursor, *ast.FuncDecl) error) TypePatcher {
	ptr := true
	return method(&ptr, receiverType, methodName, patch)
}

// ValueMethod is identical to `Method()` except that it only matches methods
// with value receivers.
func ValueMethod(receiverType, methodName string, patch func(*astutil.Cursor, *ast.FuncDecl) error) TypePatcher {
	ptr := false
	return method(&ptr, receiverType, methodName, patch)
}

func method(pointerReceiver *bool, receiverType, methodName string, patch func(*astutil.Cursor, *ast.FuncDecl) error) TypePatcher {
	return typePatcher{
		typ: (*ast.FuncDecl)(nil),
		patch: func(c *astutil.Cursor) error {
			fn, ok := c.Node().(*ast.FuncDecl)
			if !ok || fn.Recv == nil /*not a method*/ || fn.Name.Name != methodName {
				return nil
			}
			if n := len(fn.Recv.List); n != 1 {
				return fmt.Errorf("func receiver list length = %d (%v)", n, fn.Name)
			}

			var rcvTypeName *ast.Ident

			switch rcvType := fn.Recv.List[0].Type.(type) {
			case *ast.Ident:
				if pointerReceiver != nil && *pointerReceiver {
					return nil
				}
				rcvTypeName = rcvType

			case *ast.StarExpr:
				if pointerReceiver != nil && !*pointerReceiver {
					return nil
				}
				id, ok := rcvType.X.(*ast.Ident)
				if !ok {
					return fmt.Errorf("func receiver %T.X is not %T", rcvType, rcvTypeName)
				}
				rcvTypeName = id

			default:
				return fmt.Errorf("unsupported %T.Recv.List.Type type %T", fn, rcvType)
			}

			if rcvTypeName.Name != receiverType {
				return nil
			}
			return patch(c, fn)
		},
	}
}

// UnqualifiedCall returns a patch that only applies to a call to the specific,
// unqualified function. A qualified function is one that has additional
// qualifiers before the selector (e.g. `foo.Bar()` or `pkg.Bar()`); an
// unqualified function lacks any such qualifiers and applies to builtin and
// package-internal functions.
//
// The `patch` argument functions like a regular `Patch` except that its
// parameters are extended to also accept the call's AST declaration as its
// concrete type (i.e. `astutil.Cursor.Node().(*ast.CallExpr)`).
func UnqualifiedCall(name string, patch func(*astutil.Cursor, *ast.CallExpr) error) Patch {
	return func(c *astutil.Cursor) error {
		call, ok := c.Node().(*ast.CallExpr)
		if !ok {
			return nil
		}
		fn, ok := call.Fun.(*ast.Ident)
		if !ok || fn.Name != name {
			return nil
		}
		return patch(c, call)
	}
}

// Function returns a `TypePatcher` that only applies to the specific function
// declaration.
//
// The `patch` argument functions like a regular `Patch` except that its
// parameters are extended to also accept the methods's AST declaration as its
// concrete type (i.e. `astutil.Cursor.Node().(*ast.FuncDecl)`).
func Function(name string, patch func(*astutil.Cursor, *ast.FuncDecl) error) TypePatcher {
	return typePatcher{
		typ: (*ast.FuncDecl)(nil),
		patch: func(c *astutil.Cursor) error {
			fn, ok := c.Node().(*ast.FuncDecl)
			if !ok || fn.Name.Name != name {
				return nil
			}
			return patch(c, fn)
		},
	}
}

// RenameFunction does what it says on the tin.
func RenameFunction(from, to string) TypePatcher {
	return Function(from, func(c *astutil.Cursor, fn *ast.FuncDecl) error {
		fn.Name.Name = to
		return nil
	})
}
