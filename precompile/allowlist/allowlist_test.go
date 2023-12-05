// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package allowlist

import (
	"fmt"
	"testing"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

func FuzzPackReadAllowlistTest(f *testing.F) {
	f.Add(common.Address{}.Bytes())
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	f.Add(addr.Bytes())
	f.Fuzz(func(t *testing.T, b []byte) {
		testPackReadAllowlistTest(t, common.BytesToAddress(b))
	})
}

func FuzzPackReadAllowlistTestSkipCheck(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		res, err := UnpackReadAllowListInput(b, false)
		oldRes, oldErr := OldUnpackReadAllowList(b)
		if oldErr != nil {
			require.ErrorContains(t, err, oldErr.Error())
		} else {
			require.NoError(t, err)
		}
		require.Equal(t, oldRes, res)
	})
}

func TestPackReadAllowlistTest(f *testing.T) {
	testPackReadAllowlistTest(f, common.Address{})
}

func testPackReadAllowlistTest(t *testing.T, address common.Address) {
	t.Helper()
	t.Run(fmt.Sprintf("TestPackReadAllowlistTest, address %v", address), func(t *testing.T) {
		// Test PackGetFeeConfigOutputV2, UnpackGetFeeConfigOutputV2
		input, err := PackReadAllowList(address)
		require.NoError(t, err)
		// exclude 4 bytes for function selector
		input = input[4:]

		unpacked, err := UnpackReadAllowListInput(input, false)
		require.NoError(t, err)

		require.Equal(t, address, unpacked)

		// Test PackGetFeeConfigOutput, UnpackGetFeeConfigOutput
		input = OldPackReadAllowList(address)
		// exclude 4 bytes for function selector
		input = input[4:]
		require.NoError(t, err)

		unpacked, err = OldUnpackReadAllowList(input)
		require.NoError(t, err)

		require.Equal(t, address, unpacked)

		// // now mix and match
		// Test PackGetFeeConfigOutput, PackGetFeeConfigOutputV2
		input, err = PackReadAllowList(address)
		// exclude 4 bytes for function selector
		input = input[4:]
		require.NoError(t, err)
		input2 := OldPackReadAllowList(address)
		// exclude 4 bytes for function selector
		input2 = input2[4:]
		require.Equal(t, input, input2)

		// // Test UnpackGetFeeConfigOutput, UnpackGetFeeConfigOutputV2
		unpacked, err = UnpackReadAllowListInput(input2, false)
		require.NoError(t, err)
		unpacked2, err := OldUnpackReadAllowList(input)
		require.NoError(t, err)
		require.Equal(t, unpacked, unpacked2)
	})
}

func OldPackReadAllowList(address common.Address) []byte {
	readAllowListSignature := contract.CalculateFunctionSelector("readAllowList(address)")
	input := make([]byte, 0, contract.SelectorLen+common.HashLength)
	input = append(input, readAllowListSignature...)
	input = append(input, address.Hash().Bytes()...)
	return input
}

func OldUnpackReadAllowList(input []byte) (common.Address, error) {
	if len(input) != allowListInputLen {
		return common.Address{}, fmt.Errorf("invalid input length for read allow list: %d", len(input))
	}
	return common.BytesToAddress(input), nil
}
