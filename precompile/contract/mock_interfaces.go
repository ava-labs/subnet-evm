// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package contract

import (
	"math/big"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/subnet-evm/commontype"
)

// TODO: replace with gomock library

var (
	_ BlockContext    = &mockBlockContext{}
	_ AccessibleState = &mockAccessibleState{}
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
	chainConfig  ChainConfig
}

func NewMockAccessibleState(state StateDB, blockContext *mockBlockContext, snowContext *snow.Context, chainConfig ChainConfig) *mockAccessibleState {
	return &mockAccessibleState{
		state:        state,
		blockContext: blockContext,
		snowContext:  snowContext,
		chainConfig:  chainConfig,
	}
}

func (m *mockAccessibleState) GetStateDB() StateDB { return m.state }

func (m *mockAccessibleState) GetBlockContext() BlockContext { return m.blockContext }

func (m *mockAccessibleState) GetSnowContext() *snow.Context { return m.snowContext }

func (m *mockAccessibleState) GetChainConfig() ChainConfig { return m.chainConfig }

type mockChainConfig struct {
	feeConfig            commontype.FeeConfig
	allowedFeeRecipients bool
	dUpgradeTimestamp    *big.Int
}

func (m *mockChainConfig) GetFeeConfig() commontype.FeeConfig { return m.feeConfig }
func (m *mockChainConfig) AllowedFeeRecipients() bool         { return m.allowedFeeRecipients }
func (m *mockChainConfig) IsDUpgrade(timestamp *big.Int) bool {
	if m.dUpgradeTimestamp == nil {
		return false
	}
	return m.dUpgradeTimestamp.Cmp(timestamp) <= 0
}

func NewMockChainConfig(feeConfig commontype.FeeConfig, allowedFeeRecipients bool, dUpgradeTimestamp *big.Int) *mockChainConfig {
	return &mockChainConfig{
		feeConfig:            feeConfig,
		allowedFeeRecipients: allowedFeeRecipients,
		dUpgradeTimestamp:    dUpgradeTimestamp,
	}
}
