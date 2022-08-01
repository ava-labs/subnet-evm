// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package bind

// tmplData is the data structure required to fill the binding template.
type tmplPrecompileData struct {
	Contract *tmplPrecompileContract // The contract to generate into this file
}

// tmplContract contains the data needed to generate an individual contract binding.
type tmplPrecompileContract struct {
	*tmplContract
	AllowList bool // Indicator whether the contract uses AllowList precompile/
}

// tmplSourceGo is the Go precompiled source template.
const tmplSourcePrecompileGo = `
// Code generated
// This file is a generated precompile with stubbed abstract functions.

package precompile

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
)

const (
	{{range .Contract.Calls}}
	// {{.Normalized.Name}}GasCost int = SET A GAS COST HERE
	{{end}}
	{{range .Contract.Transacts}}
	// {{.Normalized.Name}}GasCost int = SET A GAS COST HERE
	{{end}}
)

// Singleton StatefulPrecompiledContract and signatures.
var (
	_ StatefulPrecompileConfig = &{{.Contract.Type}}Config{}
	{{.Contract.Type}}Precompile StatefulPrecompiledContract = createNativeMinterPrecompile({{.Contract.Type}}Address)

	{{range .Contract.Calls}}
	{{.Normalized.Name}}Signature = CalculateFunctionSelector("{{.Original.Sig}}")
	{{end}}
	{{range .Contract.Transacts}}
	{{.Normalized.Name}}Signature = CalculateFunctionSelector("{{.Original.Sig}}")

	ErrCannot{{.Normalized.Name}} = errors.New("non-enabled cannot {{.Original.Name}}")
	{{end}}
)

// {{.Contract.Type}}Config {{if .Contract.AllowList}}wraps [AllowListConfig] and uses it to implement {{else}}implements{{end}} the StatefulPrecompileConfig
// interface while adding in the {{.Contract.Type}} specific precompile address.
type {{.Contract.Type}}Config struct {
	{{if .Contract.AllowList}}
	AllowListConfig
	{{end}}
}

// Address returns the address of the {{.Contract.Type}}. Addresses reside under the precompile/params.go
func (c *{{.Contract.Type}}Config) Address() common.Address {
	return {{.Contract.Type}}Address
}

// Configure configures [state] with the initial configuration.
func (c *{{.Contract.Type}}Config) Configure(_ ChainConfig, state StateDB, _ BlockContext) {
	{{if .Contract.AllowList}}c.AllowListConfig.Configure(state, {{.Contract.Type}}Address){{end}}
	// YOUR CODE STARTS HERE
}

// Contract returns the singleton stateful precompiled contract to be used for {{.Contract.Type}}.
func (c *{{.Contract.Type}}Config) Contract() StatefulPrecompiledContract {
	return {{.Contract.Type}}Precompile
}

{{if .Contract.AllowList}}
// Get{{.Contract.Type}}Status returns the role of [address] for the {{.Contract.Type}} list.
func Get{{.Contract.Type}}Status(stateDB StateDB, address common.Address) AllowListRole {
	return getAllowListStatus(stateDB, {{.Contract.Type}}Address, address)
}

// Set{{.Contract.Type}}Status sets the permissions of [address] to [role] for the
// {{.Contract.Type}} list. Assumes [role] has already been verified as valid.
func Set{{.Contract.Type}}Status(stateDB StateDB, address common.Address, role AllowListRole) {
	setAllowListRole(stateDB, {{.Contract.Type}}Address, address, role)
}
{{end}}

{{$contract := .Contract}}
{{range .Contract.Transacts}}
func {{.Normalized.Name}}(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = deductGas(suppliedGas, {{.Normalized.Name}}GasCost); err != nil {
		return nil, 0, err
	}

	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}

	{{if $contract.AllowList}}
	stateDB := accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to modify it
	callerStatus := getAllowListStatus(stateDB, {{$contract.Type}}Address, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("%w: %s", ErrCannot{{$contract.Type}}, caller)
	}
  {{end}}

	// YOUR CODE STARTS HERE

	// Return an empty output and the remaining gas
	return []byte{}, remainingGas, nil
}

{{end}}

{{range .Contract.Calls}}
func {{.Normalized.Name}}(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = deductGas(suppliedGas, {{.Normalized.Name}}GasCost); err != nil {
		return nil, 0, err
	}

	// YOUR CODE STARTS HERE
	// output := ...

	// Return the output and the remaining gas
	return output, remainingGas, err
}
{{end}}

// create{{.Contract.Type}}Precompile returns a StatefulPrecompiledContract
// with getters and setters for the precompile.
{{if .Contract.AllowList}} //Access to the getters/setters is controlled by an allow list for [precompileAddr].{{end}}
func create{{.Contract.Type}}Precompile(precompileAddr common.Address) StatefulPrecompiledContract {
	{{.Contract.Type}}AllowListFunctions := createAllowListFunctions(precompileAddr)

	{{range .Contract.Calls}}
	{{.Normalized.Name}}Func := newStatefulPrecompileFunction({{.Normalized.Name}}Signature, {{.Normalized.Name}})
	{{end}}
	{{range .Contract.Transacts}}
	{{.Normalized.Name}}Func := newStatefulPrecompileFunction({{.Normalized.Name}}Signature, {{.Normalized.Name}})
	{{end}}

	{{.Contract.Type}}Functions := append({{.Contract.Type}}AllowListFunctions, setFeeConfigFunc, getFeeConfigFunc, getFeeConfigLastChangedAtFunc)
	// Construct the contract with no fallback function.
	contract := newStatefulPrecompileWithFunctionSelectors(nil, feeConfigManagerFunctions)
	return contract
}
`
