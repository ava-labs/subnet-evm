package vm

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
)

type (
	// RunFunc is the signature of a precompiled contract run function
	// Consider passing caller as ContractRef instead of common.Address
	RunFunc func(caller common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error)
)

var defaultEVMFactory EvmFactory = &evmFactory{}

type evmFactory struct{}

type ChainConfig interface {
	IsEIP158(blockNum *big.Int) bool
	IsSubnetEVM(timestamp uint64) bool
	IsCancun(blockNum *big.Int, timestamp uint64) bool
	IsPrecompileEnabled(addr common.Address, timestamp uint64) bool
	Rules(blockNum *big.Int, timestamp uint64) params.Rules
}

type EvmFactory interface {
	NewEVM(blockCtx BlockContext, txCtx TxContext, statedb StateDB, chainConfig ChainConfig, config Config) *EVM
}

func DefaultEVMFactory() EvmFactory {
	return defaultEVMFactory
}

func SetDefaultEVMFactory(factory EvmFactory) {
	defaultEVMFactory = factory
}

func NewEVM(blockCtx BlockContext, txCtx TxContext, statedb StateDB, chainConfig ChainConfig, config Config) *EVM {
	return DefaultEVMFactory().NewEVM(blockCtx, txCtx, statedb, chainConfig, config)
}
