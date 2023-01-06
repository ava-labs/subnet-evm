// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"math/big"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ethereum/go-ethereum/common"
)

// TODO: replace with gomock library

var (
	_ BlockContext              = &mockBlockContext{}
	_ PrecompileAccessibleState = &mockAccessibleState{}
	_ ChainConfig               = &mockChainConfig{}
	_ StatefulPrecompileConfig  = &noopStatefulPrecompileConfig{}
)

type mockBlockContext struct {
	blockNumber *big.Int
	timestamp   uint64
}

func NewMockBlockContext(blockNumber *big.Int, timestamp uint64) *mockBlockContext {
	return &mockBlockContext{
		blockNumber: blockNumber,
		timestamp:   timestamp,
	}
}

func (mb *mockBlockContext) Number() *big.Int    { return mb.blockNumber }
func (mb *mockBlockContext) Timestamp() *big.Int { return new(big.Int).SetUint64(mb.timestamp) }

type mockAccessibleState struct {
	state        StateDB
	blockContext *mockBlockContext
	snowContext  *snow.Context
}

func NewMockAccessibleState(state StateDB, blockContext *mockBlockContext, snowContext *snow.Context) *mockAccessibleState {
	return &mockAccessibleState{
		state:        state,
		blockContext: blockContext,
		snowContext:  snowContext,
	}
}

func (m *mockAccessibleState) GetStateDB() StateDB { return m.state }

func (m *mockAccessibleState) GetBlockContext() BlockContext { return m.blockContext }

func (m *mockAccessibleState) GetSnowContext() *snow.Context { return m.snowContext }

func (m *mockAccessibleState) CallFromPrecompile(caller common.Address, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	return nil, 0, nil
}

type mockChainConfig struct {
	feeConfig            commontype.FeeConfig
	allowedFeeRecipients bool
}

func NewMockChainConfig(feeConfig commontype.FeeConfig, allowedFeeRecipients bool) *mockChainConfig {
	return &mockChainConfig{
		feeConfig:            feeConfig,
		allowedFeeRecipients: allowedFeeRecipients,
	}
}

func (m *mockChainConfig) GetFeeConfig() commontype.FeeConfig { return m.feeConfig }

func (m *mockChainConfig) AllowedFeeRecipients() bool { return m.allowedFeeRecipients }

type noopStatefulPrecompileConfig struct {
}

func NewNoopStatefulPrecompileConfig() *noopStatefulPrecompileConfig {
	return &noopStatefulPrecompileConfig{}
}

func (n *noopStatefulPrecompileConfig) Address() common.Address {
	return common.Address{}
}

func (n *noopStatefulPrecompileConfig) Timestamp() *big.Int {
	return new(big.Int)
}

func (n *noopStatefulPrecompileConfig) IsDisabled() bool {
	return false
}

func (n *noopStatefulPrecompileConfig) Equal(StatefulPrecompileConfig) bool {
	return false
}

func (n *noopStatefulPrecompileConfig) Verify() error {
	return nil
}

func (n *noopStatefulPrecompileConfig) Configure(ChainConfig, StateDB, BlockContext) error {
	return nil
}

func (n *noopStatefulPrecompileConfig) Contract() StatefulPrecompiledContract {
	return nil
}

func (n *noopStatefulPrecompileConfig) String() string {
	return ""
}

func (n *noopStatefulPrecompileConfig) Key() string {
	return ""
}

func (noopStatefulPrecompileConfig) New() StatefulPrecompileConfig {
	return new(noopStatefulPrecompileConfig)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *noopStatefulPrecompileConfig) UnmarshalJSON(b []byte) error {
	return nil
}
