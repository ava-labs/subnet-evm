// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package precompilebind

// tmplSourcePrecompileConfigGo is the Go precompiled config source template.
const tmplSourcePrecompileContractTestGo = `
// Code generated
// This file is a generated precompile contract test with the skeleton of test functions.
// The file is generated by a template. Please inspect every code and comment in this file before use.

package {{.Package}}

import (
	"testing"
	"math/big"

	"github.com/ava-labs/subnet-evm/core/state"
	{{- if .Contract.AllowList}}
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	{{- end}}
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var (
	_ = vmerrs.ErrOutOfGas
	_ = big.NewInt
	_ = common.Big0
	_ = require.New
)

// These tests are run against the precompile contract directly with
// the given input and expected output. They're just a guide to
// help you write your own tests. These tests are for general cases like
// allowlist, readOnly behaviour, and gas cost. You should write your own
// tests for specific cases.
var(
	tests = map[string]testutils.PrecompileTest{
		{{- $contract := .Contract}}
		{{- $structs := .Structs}}
		{{- range .Contract.Funcs}}
		{{- $func := .}}
		{{- if $contract.AllowList}}
		{{- $roles := mkList "NoRole" "Enabled" "Admin"}}
		{{- range $role := $roles}}
		{{- $fail := and (not $func.Original.IsConstant) (eq $role "NoRole")}}
		"calling {{decapitalise $func.Normalized.Name}} from {{$role}} should {{- if $fail}} fail {{- else}} succeed{{- end}}":  {
			Caller:     allowlist.Test{{$role}}Addr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				{{- if len $func.Normalized.Inputs | lt 1}}
				// CUSTOM CODE STARTS HERE
				// populate test input here
				testInput := {{capitalise $func.Normalized.Name}}Input{}
				input, err := Pack{{$func.Normalized.Name}}(testInput)
				{{- else if len $func.Normalized.Inputs | eq 1 }}
				{{- $input := index $func.Normalized.Inputs 0}}
				// CUSTOM CODE STARTS HERE
				// set test input to a value here
				var testInput {{bindtype $input.Type $structs}}
				testInput = {{bindtypenew $input.Type $structs}}
				input, err := Pack{{$func.Normalized.Name}}(testInput)
				{{- else}}
				input, err := Pack{{$func.Normalized.Name}}()
				{{- end}}
				require.NoError(t, err)
				return input
			},
			{{- if not $fail}}
			// This test is for a successful call. You can set the expected output here.
			// CUSTOM CODE STARTS HERE
			ExpectedRes: func() []byte{
				{{- if len $func.Normalized.Outputs | eq 0}}
				// this function does not return an output, leave this one as is
				packedOutput := []byte{}
				{{- else}}
				{{- if len $func.Normalized.Outputs | lt 1}}
				var output {{capitalise $func.Normalized.Name}}Output // CUSTOM CODE FOR AN OUTPUT
				{{- else }}
				{{$output := index $func.Normalized.Outputs 0}}
				var output {{bindtype $output.Type $structs}} // CUSTOM CODE FOR AN OUTPUT
				output = {{bindtypenew $output.Type $structs}} // CUSTOM CODE FOR AN OUTPUT
				{{- end}}
				packedOutput, err := Pack{{$func.Normalized.Name}}Output(output)
				if err != nil {
					panic(err)
				}
				{{- end}}
				return packedOutput
			}(),
			{{- end}}
			SuppliedGas: {{$func.Normalized.Name}}GasCost,
			ReadOnly:    false,
			ExpectedErr: {{if $fail}} ErrCannot{{$func.Normalized.Name}}.Error() {{- else}} "" {{- end}},
		},
		{{- end}}
		{{- end}}
		{{- if not $func.Original.IsConstant}}
		"readOnly {{decapitalise $func.Normalized.Name}} should fail": {
			Caller:	common.Address{1},
			InputFn: func(t testing.TB) []byte {
				{{- if len $func.Normalized.Inputs | lt 1}}
				// CUSTOM CODE STARTS HERE
				// populate test input here
				testInput := {{capitalise $func.Normalized.Name}}Input{}
				input, err := Pack{{$func.Normalized.Name}}(testInput)
				{{- else if len $func.Normalized.Inputs | eq 1 }}
				{{- $input := index $func.Normalized.Inputs 0}}
				// CUSTOM CODE STARTS HERE
				// set test input to a value here
				var testInput {{bindtype $input.Type $structs}}
				testInput = {{bindtypenew $input.Type $structs}}
				input, err := Pack{{$func.Normalized.Name}}(testInput)
				{{- else}}
				input, err := Pack{{$func.Normalized.Name}}()
				{{- end}}
				require.NoError(t, err)
				return input
			},
			SuppliedGas:  {{$func.Normalized.Name}}GasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		{{- end}}
		"insufficient gas for {{decapitalise $func.Normalized.Name}} should fail": {
			Caller:	common.Address{1},
			InputFn: func(t testing.TB) []byte {
				{{- if len $func.Normalized.Inputs | lt 1}}
				// CUSTOM CODE STARTS HERE
				// populate test input here
				testInput := {{capitalise $func.Normalized.Name}}Input{}
				input, err := Pack{{$func.Normalized.Name}}(testInput)
				{{- else if len $func.Normalized.Inputs | eq 1 }}
				{{- $input := index $func.Normalized.Inputs 0}}
				// CUSTOM CODE STARTS HERE
				// set test input to a value here
				var testInput {{bindtype $input.Type $structs}}
				testInput = {{bindtypenew $input.Type $structs}}
				input, err := Pack{{$func.Normalized.Name}}(testInput)
				{{- else}}
				input, err := Pack{{$func.Normalized.Name}}()
				{{- end}}
				require.NoError(t, err)
				return input
			},
			SuppliedGas: {{$func.Normalized.Name}}GasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		{{- end}}
		{{- if .Contract.Fallback}}
		"insufficient gas for fallback should fail": {
			Caller:	common.Address{1},
			Input: []byte{},
			SuppliedGas: {{.Contract.Type}}FallbackGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"readOnly fallback should fail": {
			Caller:	common.Address{1},
			Input: []byte{},
			SuppliedGas: {{.Contract.Type}}FallbackGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"fallback should succeed": {
			Caller:	common.Address{1},
			Input: []byte{},
			SuppliedGas: {{.Contract.Type}}FallbackGasCost,
			ReadOnly:    false,
			ExpectedErr: "",
			// CUSTOM CODE STARTS HERE
			// set expected output here
			ExpectedRes: []byte{},
		},
		{{- end}}
	}
)

// Test{{.Contract.Type}}Run tests the Run function of the precompile contract.
func Test{{.Contract.Type}}Run(t *testing.T) {
	{{- if .Contract.AllowList}}
	// Run tests with allowlist tests.
	// This adds allowlist tests to your custom tests
	// and runs them all together.
	// Even if you don't add any custom tests, keep this. This will still
	// run the default allowlist tests.
	allowlist.RunPrecompileWithAllowListTests(t, Module, state.NewTestStateDB, tests)
	{{- else}}
	// Run tests.
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.Run(t, Module, state.NewTestStateDB(t))
		})
	}
	{{- end}}
}

func Benchmark{{.Contract.Type}}(b *testing.B) {
	{{- if .Contract.AllowList}}
	// Benchmark tests with allowlist tests.
	// This adds allowlist tests to your custom tests
	// and benchmarks them all together.
	// Even if you don't add any custom tests, keep this. This will still
	// run the default allowlist tests.
	allowlist.BenchPrecompileWithAllowList(b, Module, state.NewTestStateDB, tests)
	{{- else}}
	// Benchmark tests.
	for name, test := range tests {
		b.Run(name, func(b *testing.B) {
			test.Bench(b, Module, state.NewTestStateDB(b))
		})
	}
	{{- end}}
}

`