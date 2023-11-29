// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package allowlist

import (
	"testing"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/contract"
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
	precompileconfig.Upgrade
	AllowListConfig
}

func (d *dummyConfig) Key() string      { return "dummy" }
func (d *dummyConfig) IsDisabled() bool { return false }
func (d *dummyConfig) Verify(chainConfig precompileconfig.ChainConfig) error {
	return d.AllowListConfig.Verify(chainConfig, d.Upgrade)
}

func (d *dummyConfig) Equal(cfg precompileconfig.Config) bool {
	other, ok := (cfg).(*dummyConfig)
	if !ok {
		return false
	}
	return d.AllowListConfig.Equal(&other.AllowListConfig)
}

type dummyConfigurator struct{}

func (d *dummyConfigurator) MakeConfig() precompileconfig.Config {
	return &dummyConfig{}
}

func (d *dummyConfigurator) Configure(
	chainConfig precompileconfig.ChainConfig,
	precompileConfig precompileconfig.Config,
	state contract.StateDB,
	blockContext contract.ConfigurationBlockContext,
) error {
	cfg := precompileConfig.(*dummyConfig)
	return cfg.AllowListConfig.Configure(chainConfig, dummyAddr, state, blockContext)
}

func TestAllowListRun(t *testing.T) {
	dummyModule := modules.Module{
		Address:      dummyAddr,
		Contract:     CreateAllowListPrecompile(dummyAddr),
		Configurator: &dummyConfigurator{},
		ConfigKey:    "dummy",
	}
	RunPrecompileWithAllowListTests(t, dummyModule, state.NewTestStateDB, nil)
}

func BenchmarkAllowList(b *testing.B) {
	dummyModule := modules.Module{
		Address:      dummyAddr,
		Contract:     CreateAllowListPrecompile(dummyAddr),
		Configurator: &dummyConfigurator{},
		ConfigKey:    "dummy",
	}
	BenchPrecompileWithAllowList(b, dummyModule, state.NewTestStateDB, nil)
}

func TestSignatures(t *testing.T) {
	setAdminSignature := contract.CalculateFunctionSelector("setAdmin(address)")
	setAdminABI := AllowListABI.Methods["setAdmin"]
	require.Equal(t, setAdminSignature, setAdminABI.ID)

	setManagerSignature := contract.CalculateFunctionSelector("setManager(address)")
	setManagerABI := AllowListABI.Methods["setManager"]
	require.Equal(t, setManagerSignature, setManagerABI.ID)

	setEnabledSignature := contract.CalculateFunctionSelector("setEnabled(address)")
	setEnabledABI := AllowListABI.Methods["setEnabled"]
	require.Equal(t, setEnabledSignature, setEnabledABI.ID)

	setNoneSignature := contract.CalculateFunctionSelector("setNone(address)")
	setNoneABI := AllowListABI.Methods["setNone"]
	require.Equal(t, setNoneSignature, setNoneABI.ID)

	readAllowlistSignature := contract.CalculateFunctionSelector("readAllowList(address)")
	readAllowlistABI := AllowListABI.Methods["readAllowList"]
	require.Equal(t, readAllowlistSignature, readAllowlistABI.ID)
}
