// (c) 2019-2021, Ava Labs, Inc.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2015 The go-ethereum Authors
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
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/consensus"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/core/vm"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ava-labs/subnet-evm/stateupgrade"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

// StateProcessor is a basic Processor, which takes care of transitioning
// state from one point to another.
//
// StateProcessor implements Processor.
type StateProcessor struct {
	config *params.ChainConfig // Chain configuration options
	bc     *BlockChain         // Canonical block chain
	engine consensus.Engine    // Consensus engine used for block rewards
}

// NewStateProcessor initialises a new StateProcessor.
func NewStateProcessor(config *params.ChainConfig, bc *BlockChain, engine consensus.Engine) *StateProcessor {
	return &StateProcessor{
		config: config,
		bc:     bc,
		engine: engine,
	}
}

// Process processes the state changes according to the Ethereum rules by running
// the transaction messages using the statedb and applying any rewards to both
// the processor (coinbase) and any included uncles.
//
// Process returns the receipts and logs accumulated during the process and
// returns the amount of gas that was used in the process. If any of the
// transactions failed to execute due to insufficient gas it will return an error.
func (p *StateProcessor) Process(block *types.Block, parent *types.Header, statedb *state.StateDB, cfg vm.Config) (types.Receipts, []*types.Log, uint64, error) {
	var (
		receipts    types.Receipts
		usedGas     = new(uint64)
		header      = block.Header()
		blockHash   = block.Hash()
		blockNumber = block.Number()
		allLogs     []*types.Log
		gp          = new(GasPool).AddGas(block.GasLimit())
		timestamp   = new(big.Int).SetUint64(header.Time)
	)

	// Configure any upgrades that should go into effect during this block.
	err := ApplyUpgrades(p.config, new(big.Int).SetUint64(parent.Time), block, statedb)
	if err != nil {
		log.Error("failed to configure precompiles processing block", "hash", block.Hash(), "number", block.NumberU64(), "timestamp", block.Time(), "err", err)
		return nil, nil, 0, err
	}

	blockContext := NewEVMBlockContext(header, p.bc, nil)
	vmenv := vm.NewEVM(blockContext, vm.TxContext{}, statedb, p.config, cfg)
	// Iterate over and process the individual transactions
	for i, tx := range block.Transactions() {
		msg, err := tx.AsMessage(types.MakeSigner(p.config, header.Number, timestamp), header.BaseFee)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("could not apply tx %d [%v]: %w", i, tx.Hash().Hex(), err)
		}
		statedb.Prepare(tx.Hash(), i)
		receipt, err := applyTransaction(msg, p.config, nil, gp, statedb, blockNumber, blockHash, tx, usedGas, vmenv)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("could not apply tx %d [%v]: %w", i, tx.Hash().Hex(), err)
		}
		receipts = append(receipts, receipt)
		allLogs = append(allLogs, receipt.Logs...)
	}
	// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)
	if err := p.engine.Finalize(p.bc, block, parent, statedb, receipts); err != nil {
		return nil, nil, 0, fmt.Errorf("engine finalization check failed: %w", err)
	}

	return receipts, allLogs, *usedGas, nil
}

func applyTransaction(msg types.Message, config *params.ChainConfig, author *common.Address, gp *GasPool, statedb *state.StateDB, blockNumber *big.Int, blockHash common.Hash, tx *types.Transaction, usedGas *uint64, evm *vm.EVM) (*types.Receipt, error) {
	// Create a new context to be used in the EVM environment.
	txContext := NewEVMTxContext(msg)
	evm.Reset(txContext, statedb)

	// Apply the transaction to the current state (included in the env).
	result, err := ApplyMessage(evm, msg, gp)
	if err != nil {
		return nil, err
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
	if msg.To() == nil {
		receipt.ContractAddress = crypto.CreateAddress(evm.TxContext.Origin, tx.Nonce())
	}

	// Set the receipt logs and create the bloom filter.
	receipt.Logs = statedb.GetLogs(tx.Hash(), blockHash)
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	receipt.BlockHash = blockHash
	receipt.BlockNumber = blockNumber
	receipt.TransactionIndex = uint(statedb.TxIndex())
	return receipt, err
}

// ApplyTransaction attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. It returns the receipt
// for the transaction, gas used and an error if the transaction failed,
// indicating the block was invalid.
func ApplyTransaction(config *params.ChainConfig, bc ChainContext, author *common.Address, gp *GasPool, statedb *state.StateDB, header *types.Header, tx *types.Transaction, usedGas *uint64, cfg vm.Config) (*types.Receipt, error) {
	msg, err := tx.AsMessage(types.MakeSigner(config, header.Number, new(big.Int).SetUint64(header.Time)), header.BaseFee)
	if err != nil {
		return nil, err
	}
	// Create a new context to be used in the EVM environment
	blockContext := NewEVMBlockContext(header, bc, author)
	vmenv := vm.NewEVM(blockContext, vm.TxContext{}, statedb, config, cfg)
	return applyTransaction(msg, config, author, gp, statedb, header.Number, header.Hash(), tx, usedGas, vmenv)
}

// ApplyPrecompileActivations checks if any of the precompiles specified by the chain config are enabled or disabled by the block
// transition from [parentTimestamp] to the timestamp set in [blockContext]. If this is the case, it calls [Configure]
// to apply the necessary state transitions for the upgrade.
// This function is called within genesis setup to configure the starting state for precompiles enabled at genesis.
// In block processing and building, ApplyUpgrades is called instead which also applies state upgrades.
func ApplyPrecompileActivations(c *params.ChainConfig, parentTimestamp *big.Int, blockContext contract.BlockContext, statedb *state.StateDB) error {
	blockTimestamp := blockContext.Timestamp()
	// Note: RegisteredModules returns precompiles sorted by module addresses.
	// This ensures that the order we call Configure for each precompile is consistent.
	// This ensures even if precompiles read/write state other than their own they will observe
	// an identical global state in a deterministic order when they are configured.
	for _, module := range modules.RegisteredModules() {
		key := module.ConfigKey
		for _, activatingConfig := range c.GetActivatingPrecompileConfigs(module.Address, parentTimestamp, blockTimestamp, c.PrecompileUpgrades) {
			// If this transition activates the upgrade, configure the stateful precompile.
			// (or deconfigure it if it is being disabled.)
			if activatingConfig.IsDisabled() {
				log.Info("Disabling precompile", "name", key)
				statedb.Suicide(module.Address)
				// Calling Finalise here effectively commits Suicide call and wipes the contract state.
				// This enables re-configuration of the same contract state in the same block.
				// Without an immediate Finalise call after the Suicide, a reconfigured precompiled state can be wiped out
				// since Suicide will be committed after the reconfiguration.
				statedb.Finalise(true)
			} else {
				module, ok := modules.GetPrecompileModule(key)
				if !ok {
					return fmt.Errorf("could not find module for activating precompile, name: %s", key)
				}
				log.Info("Activating new precompile", "name", key, "config", activatingConfig)
				// Set the nonce of the precompile's address (as is done when a contract is created) to ensure
				// that it is marked as non-empty and will not be cleaned up when the statedb is finalized.
				statedb.SetNonce(module.Address, 1)
				// Set the code of the precompile's address to a non-zero length byte slice to ensure that the precompile
				// can be called from within Solidity contracts. Solidity adds a check before invoking a contract to ensure
				// that it does not attempt to invoke a non-existent contract.
				statedb.SetCode(module.Address, []byte{0x1})
				if err := module.Configure(c, activatingConfig, statedb, blockContext); err != nil {
					return fmt.Errorf("could not configure precompile, name: %s, reason: %w", key, err)
				}
			}
		}
	}
	return nil
}

// applyStateUpgrades checks if any of the state upgrades specified by the chain config are activated by the block
// transition from [parentTimestamp] to the timestamp set in [header]. If this is the case, it calls [Configure]
// to apply the necessary state transitions for the upgrade.
func applyStateUpgrades(c *params.ChainConfig, parentTimestamp *big.Int, blockContext contract.BlockContext, statedb *state.StateDB) error {
	// Apply state upgrades
	for _, upgrade := range c.GetActivatingStateUpgrades(parentTimestamp, blockContext.Timestamp(), c.StateUpgrades) {
		log.Info("Applying state upgrade", "blockNumber", blockContext.Number(), "upgrade", upgrade)
		if err := stateupgrade.Configure(&upgrade, c, statedb, blockContext); err != nil {
			return fmt.Errorf("could not configure state upgrade: %w", err)
		}
	}
	return nil
}

// ApplyUpgrades checks if any of the precompile or state upgrades specified by the chain config are activated by the block
// transition from [parentTimestamp] to the timestamp set in [header]. If this is the case, it calls [Configure]
// to apply the necessary state transitions for the upgrade.
// This function is called:
// - in block processing to update the state when processing a block.
// - in the miner to apply the state upgrades when producing a block.
func ApplyUpgrades(c *params.ChainConfig, parentTimestamp *big.Int, blockContext contract.BlockContext, statedb *state.StateDB) error {
	if err := ApplyPrecompileActivations(c, parentTimestamp, blockContext, statedb); err != nil {
		return err
	}
	return applyStateUpgrades(c, parentTimestamp, blockContext, statedb)
}
