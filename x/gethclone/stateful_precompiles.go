package main

import (
	"fmt"
	"go/ast"
	"math/big"

	"github.com/ava-labs/subnet-evm/x/gethclone/astpatch"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"golang.org/x/tools/go/ast/astutil"
)

// statefulPrecompiles is a `patchSet` that modifies the way the EVM call
// methods dispatch to precompiled contracts, allowing for integration with
// Avalanche stateful precompiles.
type statefulPrecompiles struct {
	patchedMethods map[string]bool
}

func (*statefulPrecompiles) name() string {
	return "stateful-precompiles"
}

func (*statefulPrecompiles) evmCallMethods() []string {
	return []string{"Call", "CallCode", "DelegateCall", "StaticCall"}
}

func (p *statefulPrecompiles) register(reg astpatch.PatchRegistry) {
	if p.patchedMethods == nil {
		p.patchedMethods = make(map[string]bool)
	}

	for _, method := range p.evmCallMethods() {
		reg.Add(
			geth("core/vm"),
			astpatch.Method("EVM", method, p.patchRunPrecompiledCalls),
		)
	}
}

// validate returns nil iff all `evmCallMethods()` were patched.
func (p *statefulPrecompiles) validate() error {
	for _, m := range p.evmCallMethods() {
		if !p.patchedMethods[m] {
			return fmt.Errorf("%T.%s() not patched", (*vm.EVM)(nil), m)
		}
	}
	return nil
}

// patchRunPrecompiledCalls finds all `RunPrecompiledContract()` calls inside
// `fn` and changes them to (a) call a different function; and (b) also
// propagate fn's first argument (the caller) and `evm.interpreter.readOnly`.
//
//	RunPrecompiledContract(p, input gas)
//	// becomes
//	RunStatefulPrecompiledContract(p, input, gas, caller, evm.interpreter.readOnly)
//
// The definition of `RunStatefulPrecompiledContract()` SHOULD be implemented as
// regular Go code. The determination of whether `p` is stateful or not can be
// achieved with a type switch.
func (p *statefulPrecompiles) patchRunPrecompiledCalls(_ *astutil.Cursor, fn *ast.FuncDecl) error {
	{
		// This block only locks in the assumptions we're making in the patch
		// that follows. By doing so, we (a) communicate intent should a merge
		// conflict arise in the future, and (b) ensure that assumptions can be
		// programatically verified (by the compiler).

		type (
			// We need to propagate the caller and do so by assuming that it's
			// the first parameter.
			callWithValue    func(caller vm.ContractRef, _ common.Address, input []byte, gas uint64, value *big.Int) ([]byte, uint64, error)
			callWithoutValue func(vm.ContractRef, common.Address, []byte, uint64) ([]byte, uint64, error)
		)

		var (
			evm = (*vm.EVM)(nil)

			_, _ callWithValue    = evm.Call, evm.CallCode
			_, _ callWithoutValue = evm.DelegateCall, evm.StaticCall

			// We simply extend the parameter list of the regular calls.
			_ func(_ vm.PrecompiledContract, input []byte, gas uint64) ([]byte, uint64, error) = vm.RunPrecompiledContract
		)
	}

	_, err := astpatch.Apply(fn,
		astpatch.UnqualifiedCall("RunPrecompiledContract", func(_ *astutil.Cursor, call *ast.CallExpr) error {
			call.Fun = ast.NewIdent("RunStatefulPrecompiledContract")

			call.Args = append(
				call.Args,
				ast.NewIdent(fn.Type.Params.List[0].Names[0].Name),
				&ast.SelectorExpr{
					X:   ast.NewIdent(fn.Recv.List[0].Names[0].Name),
					Sel: ast.NewIdent("interpreter.readOnly"),
				},
			)

			p.patchedMethods[fn.Name.Name] = true
			return nil
		}),
		nil,
	)
	return err
}
