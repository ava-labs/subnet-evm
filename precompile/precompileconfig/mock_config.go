// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// TODO: replace with gomock
package precompileconfig

import (
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ Config      = &noopStatefulPrecompileConfig{}
	_ ChainConfig = &mockChainConfig{}
)

type noopStatefulPrecompileConfig struct{}

func NewNoopStatefulPrecompileConfig() *noopStatefulPrecompileConfig {
	return &noopStatefulPrecompileConfig{}
}

func (n *noopStatefulPrecompileConfig) Key() string {
	return ""
}

func (n *noopStatefulPrecompileConfig) Address() common.Address {
	return common.Address{}
}

func (n *noopStatefulPrecompileConfig) Timestamp() *uint64 {
	return nil
}

func (n *noopStatefulPrecompileConfig) IsDisabled() bool {
	return false
}

func (n *noopStatefulPrecompileConfig) Equal(Config) bool {
	return false
}

func (n *noopStatefulPrecompileConfig) Verify(ChainConfig) error {
	return nil
}

type mockChainConfig struct {
	feeConfig            commontype.FeeConfig
	allowedFeeRecipients bool
	dUpgradeTimestamp    *uint64
}

func (m *mockChainConfig) GetFeeConfig() commontype.FeeConfig { return m.feeConfig }
func (m *mockChainConfig) AllowedFeeRecipients() bool         { return m.allowedFeeRecipients }
func (m *mockChainConfig) IsDUpgrade(time uint64) bool {
	return utils.IsTimestampForked(m.dUpgradeTimestamp, time)
}

func NewMockChainConfig(feeConfig commontype.FeeConfig, allowedFeeRecipients bool, dUpgradeTimestamp *uint64) *mockChainConfig {
	return &mockChainConfig{
		feeConfig:            feeConfig,
		allowedFeeRecipients: allowedFeeRecipients,
		dUpgradeTimestamp:    dUpgradeTimestamp,
	}
}

