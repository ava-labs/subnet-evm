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
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	_ "embed"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/accounts/abi/bind/precompilebind"
	"github.com/ava-labs/subnet-evm/internal/flags"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
)

var (
	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit = ""
	gitDate   = ""

	app *cli.App

	//go:embed template-readme.md
	readme string
)

var (
	// Flags needed by abigen
	abiFlag = &cli.StringFlag{
		Name:  "abi",
		Usage: "Path to the contract ABI json to generate, - for STDIN",
	}
	typeFlag = &cli.StringFlag{
		Name:  "type",
		Usage: "Struct name for the precompile (default = {abi file name})",
	}
	pkgFlag = &cli.StringFlag{
		Name:  "pkg",
		Usage: "Go package name to generate the precompile into (default = {type})",
	}
	outFlag = &cli.StringFlag{
		Name:  "out",
		Usage: "Output folder for the generated precompile files, - for STDOUT (default = ./precompile/contracts/{pkg}). Test files won't be generated if STDOUT is used",
	}
)

func init() {
	app = flags.NewApp(gitCommit, gitDate, "subnet-evm precompile generator tool")
	app.Name = "precompilegen"
	app.Flags = []cli.Flag{
		abiFlag,
		outFlag,
		pkgFlag,
		typeFlag,
	}
	app.Action = precompilegen
}

func precompilegen(c *cli.Context) error {
	outFlagStr := c.String(outFlag.Name)
	isOutStdout := outFlagStr == "-"

	if isOutStdout && !c.IsSet(typeFlag.Name) {
		utils.Fatalf("type (--type) should be set explicitly for STDOUT ")
	}
	lang := bind.LangGo
	// If the entire solidity code was specified, build and bind based on that
	var (
		abis    []string
		bins    []string
		types   []string
		sigs    []map[string]string
		libs    = make(map[string]string)
		aliases = make(map[string]string)
	)
	if c.String(abiFlag.Name) == "" {
		utils.Fatalf("no abi path is specified (--abi)")
	}
	// Load up the ABI
	var (
		abi []byte
		err error
	)

	input := c.String(abiFlag.Name)
	if input == "-" {
		abi, err = io.ReadAll(os.Stdin)
	} else {
		abi, err = os.ReadFile(input)
	}
	if err != nil {
		utils.Fatalf("Failed to read input ABI: %v", err)
	}
	abis = append(abis, string(abi))

	bins = append(bins, "")

	kind := c.String(typeFlag.Name)
	if kind == "" {
		fn := filepath.Base(input)
		kind = strings.TrimSuffix(fn, filepath.Ext(fn))
		kind = strings.TrimSpace(kind)
	}
	types = append(types, kind)

	pkg := c.String(pkgFlag.Name)
	if pkg == "" {
		pkg = strings.ToLower(kind)
	}

	if outFlagStr == "" {
		outFlagStr = filepath.Join("./precompile/contracts", pkg)
	}

	abifilename := ""
	abipath := ""
	// we should not generate the abi file if output is set to stdout
	if !isOutStdout {
		// get file name from the output path
		abifilename = "contract.abi"
		abipath = filepath.Join(outFlagStr, abifilename)
	}
	// if output is set to stdout, we should not generate the test codes
	generateTests := !isOutStdout

	// Generate the contract precompile
	bindedFiles, err := precompilebind.PrecompileBind(types, abis, bins, sigs, pkg, lang, libs, aliases, abifilename, generateTests)
	if err != nil {
		utils.Fatalf("Failed to generate precompile: %v", err)
	}

	configCode := bindedFiles.Config
	contractCode := bindedFiles.Contract
	moduleCode := bindedFiles.Module

	// Either flush it out to a file or display on the standard output
	// Skip displaying test codes here.
	if isOutStdout {
		fmt.Print("-----Config Code-----\n")
		fmt.Printf("%s\n", configCode)
		fmt.Print("-----Contract Code-----\n")
		fmt.Printf("%s\n", contractCode)
		fmt.Print("-----Module Code-----\n")
		fmt.Printf("%s\n", moduleCode)
		return nil
	}

	if _, err := os.Stat(outFlagStr); os.IsNotExist(err) {
		os.MkdirAll(outFlagStr, 0o700) // Create your file
	}

	// Write the generated config code to the output folder
	configCodeOut := filepath.Join(outFlagStr, "config.go")
	if err := os.WriteFile(configCodeOut, []byte(configCode), 0o600); err != nil {
		utils.Fatalf("Failed to write generated config code: %v", err)
	}

	// Write the generated contract code to the output folder
	contractCodeOut := filepath.Join(outFlagStr, "contract.go")
	if err := os.WriteFile(contractCodeOut, []byte(contractCode), 0o600); err != nil {
		utils.Fatalf("Failed to write generated contract code: %v", err)
	}

	// Write the generated module code to the output folder
	moduleCodeOut := filepath.Join(outFlagStr, "module.go")
	if err := os.WriteFile(moduleCodeOut, []byte(moduleCode), 0o600); err != nil {
		utils.Fatalf("Failed to write generated module code: %v", err)
	}

	// Write the ABI to the output folder
	if err := os.WriteFile(abipath, []byte(abis[0]), 0o600); err != nil {
		utils.Fatalf("Failed to write ABI: %v", err)
	}

	// Write the README to the output folder
	readmeOut := filepath.Join(outFlagStr, "README.md")
	if err := os.WriteFile(readmeOut, []byte(readme), 0o600); err != nil {
		utils.Fatalf("Failed to write README: %v", err)
	}

	// Write the test code to the output folder
	if generateTests {
		configTestCode := bindedFiles.ConfigTest
		configTestCodeOut := filepath.Join(outFlagStr, "config_test.go")
		if err := os.WriteFile(configTestCodeOut, []byte(configTestCode), 0o600); err != nil {
			utils.Fatalf("Failed to write generated test code: %v", err)
		}

		contractTestCode := bindedFiles.ContractTest
		contractTestCodeOut := filepath.Join(outFlagStr, "contract_test.go")
		if err := os.WriteFile(contractTestCodeOut, []byte(contractTestCode), 0o600); err != nil {
			utils.Fatalf("Failed to write generated test code: %v", err)
		}
	}

	fmt.Println("Precompile files generated successfully at: ", outFlagStr)
	return nil
}

func main() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
