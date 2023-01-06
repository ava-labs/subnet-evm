// (c) 2019-2020, Ava Labs, Inc.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package bind generates Ethereum contract Go bindings.
//
// Detailed usage document and tutorial available on the go-ethereum Wiki page:
// https://github.com/ethereum/go-ethereum/wiki/Native-DApps:-Go-bindings-to-Ethereum-contracts
package bind

import (
	"errors"
	"fmt"
)

// PrecompileBind generates a Go binding for a precompiled contract. It returns config binding and contract binding.
func PrecompileBind(types []string, abis []string, bytecodes []string, fsigs []map[string]string, pkg string, lang Lang, libs map[string]string, aliases map[string]string, abifilename string) (string, string, error) {
	// create hooks
	configHook := createPrecompileHook(abifilename, tmplSourcePrecompileConfigGo)
	contractHook := createPrecompileHook(abifilename, tmplSourcePrecompileContractGo)

	configBind, err := bindHelper(types, abis, bytecodes, fsigs, pkg, lang, libs, aliases, configHook)
	if err != nil {
		return "", "", err
	}
	contractBind, err := bindHelper(types, abis, bytecodes, fsigs, pkg, lang, libs, aliases, contractHook)

	return configBind, contractBind, err
}

func createPrecompileHook(abifilename string, template string) BindHook {
	var bindHook BindHook = func(lang Lang, pkg string, types []string, contracts map[string]*tmplContract, structs map[string]*tmplStruct) (interface{}, string, error) {
		// verify first
		if lang != LangGo {
			return nil, "", errors.New("only GoLang binding for precompiled contracts is supported yet")
		}

		if len(types) != 1 {
			return nil, "", errors.New("cannot generate more than 1 contract")
		}
		funcs := make(map[string]*tmplMethod)

		contract := contracts[types[0]]

		for k, v := range contract.Transacts {
			if err := checkOutputName(*v); err != nil {
				return nil, "", err
			}
			funcs[k] = v
		}

		for k, v := range contract.Calls {
			if err := checkOutputName(*v); err != nil {
				return nil, "", err
			}
			funcs[k] = v
		}
		isAllowList := allowListEnabled(funcs)
		if isAllowList {
			// remove these functions as we will directly inherit AllowList
			delete(funcs, readAllowListFuncKey)
			delete(funcs, setAdminFuncKey)
			delete(funcs, setEnabledFuncKey)
			delete(funcs, setNoneFuncKey)
		}

		precompileContract := &tmplPrecompileContract{
			tmplContract: contract,
			AllowList:    isAllowList,
			Funcs:        funcs,
			ABIFilename:  abifilename,
		}

		data := &tmplPrecompileData{
			Contract: precompileContract,
			Structs:  structs,
			Package:  pkg,
		}
		return data, template, nil
	}
	return bindHook
}

func allowListEnabled(funcs map[string]*tmplMethod) bool {
	keys := []string{readAllowListFuncKey, setAdminFuncKey, setEnabledFuncKey, setNoneFuncKey}
	for _, key := range keys {
		if _, ok := funcs[key]; !ok {
			return false
		}
	}
	return true
}

func checkOutputName(method tmplMethod) error {
	for _, output := range method.Original.Outputs {
		if output.Name == "" {
			return fmt.Errorf("ABI outputs for %s require a name to generate the precompile binding, re-generate the ABI from a Solidity source file with all named outputs", method.Original.Name)
		}
	}
	return nil
}
