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

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/core/vm"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed TrieStressTest.bin
	stressBinStr string
	//go:embed TrieStressTest.abi
	stressABIStr string
)

func BenchmarkTrie(t *testing.B) {
	benchInsertChain(t, true, stressTestTrieDb(t, 100, 6, 50, 1202102))
}

func applyTransactionAndGetResult(msg *Message, config *params.ChainConfig, gp *GasPool, statedb *state.StateDB, blockNumber *big.Int, blockHash common.Hash, tx *types.Transaction, usedGas *uint64, evm *vm.EVM) (*types.Receipt, *ExecutionResult, error) {
	// Create a new context to be used in the EVM environment.
	txContext := NewEVMTxContext(msg)
	evm.Reset(txContext, statedb)

	// Apply the transaction to the current state (included in the env).
	result, err := ApplyMessage(evm, msg, gp)
	if err != nil {
		return nil, nil, err
	}

	// Update the state with pending changes.
	var root []byte
	if config.IsByzantium(blockNumber) {
		statedb.Finalise(true)
	} else {
		root = statedb.IntermediateRoot(config.IsEIP158(blockNumber)).Bytes()
	}
	*usedGas += result.UsedGas

	// Create a new receipt for the transaction, storing the intermediate root and gas used
	// by the tx.
	receipt := &types.Receipt{Type: tx.Type(), PostState: root, CumulativeGasUsed: *usedGas}
	if result.Failed() {
		receipt.Status = types.ReceiptStatusFailed
	} else {
		receipt.Status = types.ReceiptStatusSuccessful
	}
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = result.UsedGas

	// If the transaction created a contract, store the creation address in the receipt.
	if msg.To == nil {
		receipt.ContractAddress = crypto.CreateAddress(evm.TxContext.Origin, tx.Nonce())
	}

	// Set the receipt logs and create the bloom filter.
	receipt.Logs = statedb.GetLogs(tx.Hash(), blockNumber.Uint64(), blockHash)
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	receipt.BlockHash = blockHash
	receipt.BlockNumber = blockNumber
	receipt.TransactionIndex = uint(statedb.TxIndex())
	return receipt, result, err
}

func ApplyTransactionAndGetResult(config *params.ChainConfig, bc ChainContext, blockContext vm.BlockContext, gp *GasPool, statedb *state.StateDB, header *types.Header, tx *types.Transaction, usedGas *uint64, cfg vm.Config) (*types.Receipt, *ExecutionResult, error) {
	msg, err := TransactionToMessage(tx, types.MakeSigner(config, header.Number, header.Time), header.BaseFee)
	if err != nil {
		return nil, nil, err
	}
	// Create a new context to be used in the EVM environment
	vmenv := vm.NewEVM(blockContext, vm.TxContext{}, statedb, config, cfg)
	return applyTransactionAndGetResult(msg, config, gp, statedb, header.Number, header.Hash(), tx, usedGas, vmenv)
}

func (b *BlockGen) AddTxOrFail(tx *types.Transaction) (*types.Receipt, error) {
	if b.gasPool == nil {
		b.SetCoinbase(common.Address{})
	}
	b.statedb.SetTxContext(tx.Hash(), len(b.txs))
	blockContext := NewEVMBlockContext(b.header, nil, &b.header.Coinbase)
	receipt, result, err := ApplyTransactionAndGetResult(b.config, nil, blockContext, b.gasPool, b.statedb, b.header, tx, &b.header.GasUsed, vm.Config{})
	if err != nil {
		return nil, err
	}

	b.txs = append(b.txs, tx)
	b.receipts = append(b.receipts, receipt)

	if result.Err != nil {
		return receipt, result.Err
	}
	return receipt, nil
}

func stressTestTrieDb(t *testing.B, numContracts int, callsPerBlock int, elements int64, gasTxLimit uint64) func(int, *BlockGen) {
	require := require.New(t)
	contractAddr := make([]common.Address, numContracts)
	contractTxs := make([]*types.Transaction, numContracts)

	gasPrice := big.NewInt(225000000000)
	gasCreation := uint64(258000)
	deployedContracts := 0

	for i := 0; i < numContracts; i++ {
		nonce := uint64(i)
		tx, _ := types.SignTx(types.NewContractCreation(nonce, big.NewInt(0), gasCreation, gasPrice, common.FromHex(stressBinStr)), signer, testKey)
		sender, _ := types.Sender(signer, tx)
		contractTxs[i] = tx
		contractAddr[i] = crypto.CreateAddress(sender, nonce)
	}

	stressABI := contract.ParseABI(stressABIStr)
	txPayload, _ := stressABI.Pack(
		"writeValues",
		big.NewInt(elements),
	)

	return func(i int, gen *BlockGen) {
		if len(contractTxs) != deployedContracts {
			block := gen.PrevBlock(i - 1)
			gas := block.GasLimit()
			for ; deployedContracts < len(contractTxs) && gasCreation < gas; deployedContracts++ {
				_, err := gen.AddTxOrFail(contractTxs[deployedContracts])
				require.NoError(err)
				gas -= gasCreation
			}
			return
		}

		for e := 0; e < callsPerBlock; e++ {
			contractId := (i + e) % deployedContracts
			tx, err := types.SignTx(types.NewTransaction(gen.TxNonce(benchRootAddr), contractAddr[contractId], big.NewInt(0), gasTxLimit, gasPrice, txPayload), signer, testKey)
			require.NoError(err)
			_, err = gen.AddTxOrFail(tx)
			require.NoError(err)
		}
	}
}
