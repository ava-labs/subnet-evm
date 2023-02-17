// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package allowlist

import (
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/contracts/test_utils"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
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

	tests := map[string]test_utils.PrecompileTest{
		"initial config sets admins": {
			Config: &dummyConfig{
				&AllowListConfig{
					AdminAddresses: []common.Address{NoRoleAddr, EnabledAddr},
				},
			},
			SuppliedGas: 0,
			ReadOnly:    false,
			AfterHook: func(t *testing.T, state contract.StateDB) {
				require.Equal(t, AdminRole, GetAllowListStatus(state, dummyAddr, NoRoleAddr))
				require.Equal(t, AdminRole, GetAllowListStatus(state, dummyAddr, EnabledAddr))
			},
		},
		"initial config sets enabled": {
			Config: &dummyConfig{
				&AllowListConfig{
					EnabledAddresses: []common.Address{NoRoleAddr, AdminAddr},
				},
			},
			SuppliedGas: 0,
			ReadOnly:    false,
			AfterHook: func(t *testing.T, state contract.StateDB) {
				require.Equal(t, EnabledRole, GetAllowListStatus(state, dummyAddr, AdminAddr))
				require.Equal(t, EnabledRole, GetAllowListStatus(state, dummyAddr, NoRoleAddr))
			},
		},
	}

	RunTestsWithAllowListSetup(t, dummyModule, state.NewTestStateDB, AllowListTests(dummyModule))
	RunTestsWithAllowListSetup(t, dummyModule, state.NewTestStateDB, tests)
}
