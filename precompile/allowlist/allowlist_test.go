// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package allowlist

import (
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ precompileconfig.Config = &dummyConfig{}
	_ contract.Configurator   = &dummyConfigurator{}

	dummyAddr = common.Address{1}
)

type dummyConfig struct {
	*AllowListConfig
}

func (d *dummyConfig) Key() string         { return "dummy" }
func (d *dummyConfig) Timestamp() *big.Int { return common.Big0 }
func (d *dummyConfig) IsDisabled() bool    { return false }
func (d *dummyConfig) Equal(other precompileconfig.Config) bool {
	return d.AllowListConfig.Equal(other.(*dummyConfig).AllowListConfig)
}

type dummyConfigurator struct{}

func (d *dummyConfigurator) NewConfig() precompileconfig.Config {
	return &dummyConfig{}
}

func (d *dummyConfigurator) Configure(
	chainConfig contract.ChainConfig,
	precompileConfig precompileconfig.Config,
	state contract.StateDB,
	blockContext contract.BlockContext,
) error {
	cfg := precompileConfig.(*dummyConfig)
	return cfg.Configure(state, dummyAddr)
}

func TestAllowListRun(t *testing.T) {
	dummyModule := modules.Module{
		Address:      dummyAddr,
		Contract:     CreateAllowListPrecompile(dummyAddr),
		Configurator: &dummyConfigurator{},
	}
	RunTestsWithAllowListSetup(t, dummyModule, state.NewTestStateDB, AllowListTests(dummyModule))
}
