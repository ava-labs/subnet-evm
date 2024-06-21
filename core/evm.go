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

package core

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/subnet-evm/consensus"
	"github.com/ava-labs/subnet-evm/consensus/misc/eip4844"
	"github.com/ava-labs/subnet-evm/constants"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/contracts/deployerallowlist"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/predicate"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	gethparams "github.com/ethereum/go-ethereum/params"
	//"github.com/ethereum/go-ethereum/log"
)

// ChainContext supports retrieving headers and consensus parameters from the
// current blockchain to be used during transaction processing.
type ChainContext interface {
	// Engine retrieves the chain's consensus engine.
	Engine() consensus.Engine

	// GetHeader returns the header corresponding to the hash/number argument pair.
	GetHeader(common.Hash, uint64) *types.Header
}

// NewEVMBlockContext creates a new context for use in the EVM.
func NewEVMBlockContext(header *types.Header, chain ChainContext, author *common.Address) vm.BlockContext {
	predicateBytes, ok := predicate.GetPredicateResultBytes(header.Extra)
	if !ok {
		return newEVMBlockContext(header, chain, author, nil)
	}
	// Prior to Durango, the VM enforces the extra data is smaller than or
	// equal to this size. After Durango, the VM pre-verifies the extra
	// data past the dynamic fee rollup window is valid.
	predicateResults, err := predicate.ParseResults(predicateBytes)
	if err != nil {
		log.Error("failed to parse predicate results creating new block context", "err", err, "extra", header.Extra)
		// As mentioned above, we pre-verify the extra data to ensure this never happens.
		// If we hit an error, construct a new block context rather than use a potentially half initialized value
		// as defense in depth.
		return newEVMBlockContext(header, chain, author, nil)
	}
	return newEVMBlockContext(header, chain, author, predicateResults)
}

// NewEVMBlockContextWithPredicateResults creates a new context for use in the EVM with an override for the predicate results that is not present
// in header.Extra.
// This function is used to create a BlockContext when the header Extra data is not fully formed yet and it's more efficient to pass in predicateResults
// directly rather than re-encode the latest results when executing each individaul transaction.
func NewEVMBlockContextWithPredicateResults(header *types.Header, chain ChainContext, author *common.Address, predicateResults *predicate.Results) vm.BlockContext {
	return newEVMBlockContext(header, chain, author, predicateResults)
}

func newEVMBlockContext(header *types.Header, chain ChainContext, author *common.Address, predicateResults *predicate.Results) vm.BlockContext {
	var (
		beneficiary common.Address
		baseFee     *big.Int
		blobBaseFee *big.Int
	)

	// If we don't have an explicit author (i.e. not mining), extract from the header
	if author == nil {
		beneficiary, _ = chain.Engine().Author(header) // Ignore error, we're past header validation
	} else {
		beneficiary = *author
	}
	if header.BaseFee != nil {
		baseFee = new(big.Int).Set(header.BaseFee)
	}
	if header.ExcessBlobGas != nil {
		blobBaseFee = eip4844.CalcBlobFee(*header.ExcessBlobGas)
	}
	return vm.BlockContext{
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		GetHash:     GetHashFn(header, chain),
		Coinbase:    beneficiary,
		BlockNumber: new(big.Int).Set(header.Number),
		Time:        header.Time,
		Difficulty:  new(big.Int).Set(header.Difficulty),
		BaseFee:     baseFee,
		BlobBaseFee: blobBaseFee,
		GasLimit:    header.GasLimit,
		Extra:       predicateResults,
	}
}

// NewEVMTxContext creates a new transaction context for a single transaction.
func NewEVMTxContext(msg *Message) vm.TxContext {
	ctx := vm.TxContext{
		Origin:     msg.From,
		GasPrice:   new(big.Int).Set(msg.GasPrice),
		BlobHashes: msg.BlobHashes,
	}
	if msg.BlobGasFeeCap != nil {
		ctx.BlobFeeCap = new(big.Int).Set(msg.BlobGasFeeCap)
	}
	return ctx
}

// GetHashFn returns a GetHashFunc which retrieves header hashes by number
func GetHashFn(ref *types.Header, chain ChainContext) func(n uint64) common.Hash {
	// Cache will initially contain [refHash.parent],
	// Then fill up with [refHash.p, refHash.pp, refHash.ppp, ...]
	var cache []common.Hash

	return func(n uint64) common.Hash {
		if ref.Number.Uint64() <= n {
			// This situation can happen if we're doing tracing and using
			// block overrides.
			return common.Hash{}
		}
		// If there's no hash cache yet, make one
		if len(cache) == 0 {
			cache = append(cache, ref.ParentHash)
		}
		if idx := ref.Number.Uint64() - n - 1; idx < uint64(len(cache)) {
			return cache[idx]
		}
		// No luck in the cache, but we can start iterating from the last element we already know
		lastKnownHash := cache[len(cache)-1]
		lastKnownNumber := ref.Number.Uint64() - uint64(len(cache))

		for {
			header := chain.GetHeader(lastKnownHash, lastKnownNumber)
			if header == nil {
				break
			}
			cache = append(cache, header.ParentHash)
			lastKnownHash = header.ParentHash
			lastKnownNumber = header.Number.Uint64() - 1
			if n == lastKnownNumber {
				return lastKnownHash
			}
		}
		return common.Hash{}
	}
}

// CanTransfer checks whether there are enough funds in the address' account to make a transfer.
// This does not take the necessary gas in to account to make the transfer valid.
func CanTransfer(db vm.StateDB, addr common.Address, amount *big.Int) bool {
	return db.GetBalance(addr).Cmp(amount) >= 0
}

// Transfer subtracts amount from sender and adds amount to recipient using the given Db
func Transfer(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {
	db.SubBalance(sender, amount)
	db.AddBalance(recipient, amount)
}

type EVM struct {
	*vm.EVM

	chainConfig *params.ChainConfig
	stateDB     vmStateDB
}

func NewEVM(blockCtx vm.BlockContext, txCtx vm.TxContext, statedb vmStateDB, chainConfig *params.ChainConfig, config vm.Config) *EVM {
	evm := &EVM{
		chainConfig: chainConfig,
		stateDB:     statedb,
	}

	rules := chainConfig.Rules(blockCtx.BlockNumber, blockCtx.Time)
	switch {
	case rules.IsCancun:
		config.JumpTable = &vm.SubnetEVMCancunInstructionSet
	case rules.IsDurango:
		config.JumpTable = &vm.SubnetEVMDurangoInstructionSet
	case rules.IsSubnetEVM:
		config.JumpTable = &vm.SubnetEVMInstructionSet
	}
	config.ActivePrecompiles = ActivePrecompiles(rules)
	config.IsProhibited = func(addr common.Address) error {
		if IsProhibited(addr) {
			return vmerrs.ErrAddrProhibited
		}
		return nil
	}
	config.CanDeploy = func(origin common.Address) error {
		// If the allow list is enabled, check that [origin] has permission to deploy a contract.
		if rules.IsPrecompileEnabled(deployerallowlist.ContractAddress) {
			allowListRole := deployerallowlist.GetContractDeployerAllowListStatus(evm.stateDB, origin)
			if !allowListRole.IsEnabled() {
				return fmt.Errorf("tx.origin %s is not authorized to deploy a contract", origin)
			}
		}
		return nil
	}

	config.CustomPrecompiles = make(map[common.Address]vm.RunFunc)

	// stateful precompiles
	var precompiles map[common.Address]contract.StatefulPrecompiledContract
	for addr, precompile := range precompiles {
		addr, precompile := addr, precompile
		config.CustomPrecompiles[addr] = func(caller common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
			ret, remainingGas, err = precompile.Run(evm, caller, addr, input, suppliedGas, readOnly)
			return ret, remainingGas, fromVMErr(err)
		}
	}

	// module precompiles
	for addr := range rules.ActivePrecompiles {
		addr := addr
		module, ok := modules.GetPrecompileModuleByAddress(addr)
		if !ok {
			continue
		}
		config.CustomPrecompiles[addr] = func(caller common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
			ret, remainingGas, err = module.Contract.Run(evm, caller, addr, input, suppliedGas, readOnly)
			return ret, remainingGas, fromVMErr(err)
		}
	}

	evm.EVM = vm.NewEVM(blockCtx, txCtx, &stateDBWrapper{evm.stateDB}, &chainConfigWrapper{chainConfig}, config)
	return evm
}

func fromVMErr(err error) error {
	switch err {
	case vmerrs.ErrExecutionReverted:
		return vm.ErrExecutionReverted
	case vmerrs.ErrOutOfGas:
		return vm.ErrOutOfGas
	case vmerrs.ErrInsufficientBalance:
		return vm.ErrInsufficientBalance
	case vmerrs.ErrWriteProtection:
		return vm.ErrWriteProtection
	}
	return err
}

type blockContext struct {
	*vm.BlockContext
}

func (bc *blockContext) GetPredicateResults(txHash common.Hash, address common.Address) []byte {
	pr := bc.BlockContext.Extra.(*predicate.Results)
	if pr == nil {
		return nil
	}
	return pr.GetResults(txHash, address)
}

func (evm *EVM) GetBlockContext() contract.BlockContext {
	return &blockContext{&evm.EVM.Context}
}

func (evm *EVM) GetChainConfig() precompileconfig.ChainConfig {
	return evm.chainConfig
}

func (evm *EVM) GetSnowContext() *snow.Context {
	return evm.chainConfig.AvalancheContext.SnowCtx
}

func (evm *EVM) GetStateDB() contract.StateDB {
	return evm.stateDB
}

type stateDBWrapper struct {
	StateDB
}

func (s *stateDBWrapper) AddLog(log *gethtypes.Log) {
	s.StateDB.AddLog(log.Address, log.Topics, log.Data, log.BlockNumber)
}

type vmStateDB interface {
	StateDB
	contract.StateDB
}

type chainConfigWrapper struct {
	*params.ChainConfig
}

func (c *chainConfigWrapper) IsLondon(blockNum *big.Int) bool {
	panic("should not be called")
}

func (c *chainConfigWrapper) Rules(blockNum *big.Int, isMerge bool, timestamp uint64) gethparams.Rules {
	rules := c.ChainConfig.Rules(blockNum, timestamp)
	return rules.AsGeth()
}

func unwrapStateDB(db vm.StateDB) StateDB {
	return db.(*stateDBWrapper).StateDB
}

// IsProhibited returns true if [addr] is in the prohibited list of addresses which should
// not be allowed as an EOA or newly created contract address.
func IsProhibited(addr common.Address) bool {
	if addr == constants.BlackholeAddr {
		return true
	}

	return modules.ReservedAddress(addr)
}

var BuiltinAddr = common.Address{
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}
