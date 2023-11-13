// (c) 2020-2021, Ava Labs, Inc.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2014 The go-ethereum Authors
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

package core

import (
	_ "embed"
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	//go:embed TrieStressTest.bin
	stressBinStr string
	//go:embed TrieStressTest.abi
	stressABIStr string
)

func BenchmarkTrie(t *testing.B) {
	benchInsertChain(t, true, generateTx(5000))
}

func generateTx(elements int64) func(int, *BlockGen) {
	return func(i int, gen *BlockGen) {
		gasPrice := big.NewInt(225000000000)
		nonce := gen.TxNonce(benchRootAddr)
		tx := types.NewContractCreation(nonce, big.NewInt(0), 3000000, gasPrice, common.FromHex(stressBinStr))
		tx, _ = types.SignTx(tx, signer, testKey)
		sender, _ := types.Sender(signer, tx)
		gen.AddTx(tx)

		contractAddr := crypto.CreateAddress(sender, nonce)

		stressABI := contract.ParseABI(stressABIStr)
		txPayload, _ := stressABI.Pack(
			"writeValues",
			big.NewInt(elements),
		)
		tx = types.NewTransaction(gen.TxNonce(benchRootAddr), contractAddr, big.NewInt(0), 3000000, gasPrice, txPayload)
		tx, _ = types.SignTx(tx, signer, testKey)
		gen.AddTx(tx)
	}
}
