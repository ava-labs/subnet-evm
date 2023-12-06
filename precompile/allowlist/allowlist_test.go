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

	// AllowList function signatures
	setAdminSignature      = contract.CalculateFunctionSelector("setAdmin(address)")
	setManagerSignature    = contract.CalculateFunctionSelector("setManager(address)")
	setEnabledSignature    = contract.CalculateFunctionSelector("setEnabled(address)")
	setNoneSignature       = contract.CalculateFunctionSelector("setNone(address)")
	readAllowListSignature = contract.CalculateFunctionSelector("readAllowList(address)")
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

func TestFunctionSignatures(t *testing.T) {
	require := require.New(t)
	setAdminABI := AllowListABI.Methods["setAdmin"]
	require.Equal(setAdminSignature, setAdminABI.ID)

	setManagerABI := AllowListABI.Methods["setManager"]
	require.Equal(setManagerSignature, setManagerABI.ID)

	setEnabledABI := AllowListABI.Methods["setEnabled"]
	require.Equal(setEnabledSignature, setEnabledABI.ID)

	setNoneABI := AllowListABI.Methods["setNone"]
	require.Equal(setNoneSignature, setNoneABI.ID)

	readAllowlistABI := AllowListABI.Methods["readAllowList"]
	require.Equal(readAllowListSignature, readAllowlistABI.ID)
}

func FuzzPackReadAllowlistTest(f *testing.F) {
	f.Add(common.Address{}.Bytes())
	key, err := crypto.GenerateKey()
	require.NoError(f, err)
	addr := crypto.PubkeyToAddress(key.PublicKey)
	f.Add(addr.Bytes())
	f.Fuzz(func(t *testing.T, b []byte) {
		testPackReadAllowlistTest(t, common.BytesToAddress(b))
	})
}

func FuzzPackReadAllowlistTestSkipCheck(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		require := require.New(t)
		res, err := UnpackReadAllowListInput(b, false)
		oldRes, oldErr := OldUnpackReadAllowList(b)
		if oldErr != nil {
			require.ErrorContains(err, oldErr.Error())
		} else {
			require.NoError(err)
		}
		require.Equal(oldRes, res)
	})
}

func TestPackReadAllowlistTest(f *testing.T) {
	testPackReadAllowlistTest(f, common.Address{})
}

func testPackReadAllowlistTest(t *testing.T, address common.Address) {
	t.Helper()
	require := require.New(t)
	t.Run(fmt.Sprintf("TestPackReadAllowlistTest, address %v", address), func(t *testing.T) {
		// use new Pack/Unpack methods
		input, err := PackReadAllowList(address)
		require.NoError(err)
		// exclude 4 bytes for function selector
		input = input[4:]
		unpacked, err := UnpackReadAllowListInput(input, false)
		require.NoError(err)
		require.Equal(address, unpacked)

		// use old Pack/Unpack methods
		input = OldPackReadAllowList(address)
		// exclude 4 bytes for function selector
		input = input[4:]
		require.NoError(err)
		unpacked, err = OldUnpackReadAllowList(input)
		require.NoError(err)
		require.Equal(address, unpacked)

		// now mix and match old and new methods
		input, err = PackReadAllowList(address)
		require.NoError(err)
		// exclude 4 bytes for function selector
		input = input[4:]
		input2 := OldPackReadAllowList(address)
		// exclude 4 bytes for function selector
		input2 = input2[4:]
		require.Equal(input, input2)
		unpacked, err = UnpackReadAllowListInput(input2, false)
		require.NoError(err)
		unpacked2, err := OldUnpackReadAllowList(input)
		require.NoError(err)
		require.Equal(unpacked, unpacked2)
	})
}

func OldPackReadAllowList(address common.Address) []byte {
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

func FuzzPackModifyAllowListTest(f *testing.F) {
	f.Add(common.Address{}.Bytes(), uint(0))
	key, err := crypto.GenerateKey()
	require.NoError(f, err)
	addr := crypto.PubkeyToAddress(key.PublicKey)
	f.Add(addr.Bytes(), uint(0))
	f.Fuzz(func(t *testing.T, b []byte, roleIndex uint) {
		testPackModifyAllowListTest(t, common.BytesToAddress(b), getRole(roleIndex))
	})
}

func FuzzPackModifyAllowlistTestSkipCheck(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		require := require.New(t)
		res, err := UnpackModifyAllowListInput(b, AdminRole, false)
		oldRes, oldErr := OldUnpackModifyAllowList(b, AdminRole)
		if oldErr != nil {
			require.ErrorContains(err, oldErr.Error())
		} else {
			require.NoError(err)
		}
		require.Equal(oldRes, res)
	})
}

func testPackModifyAllowListTest(t *testing.T, address common.Address, role Role) {
	t.Helper()
	require := require.New(t)
	t.Run(fmt.Sprintf("TestPackModifyAllowlistTest, address %v, role %s", address, role.String()), func(t *testing.T) {
		// use new Pack/Unpack methods
		input, err := PackModifyAllowList(address, role)
		require.NoError(err)
		// exclude 4 bytes for function selector
		input = input[4:]
		unpacked, err := UnpackModifyAllowListInput(input, role, false)
		require.NoError(err)
		require.Equal(address, unpacked)

		// use old Pack/Unpack methods
		input, err = OldPackModifyAllowList(address, role)
		require.NoError(err)
		// exclude 4 bytes for function selector
		input = input[4:]
		require.NoError(err)

		unpacked, err = OldUnpackModifyAllowList(input, role)
		require.NoError(err)

		require.Equal(address, unpacked)

		// now mix and match new and old methods
		input, err = PackModifyAllowList(address, role)
		require.NoError(err)
		// exclude 4 bytes for function selector
		input = input[4:]
		input2, err := OldPackModifyAllowList(address, role)
		require.NoError(err)
		// exclude 4 bytes for function selector
		input2 = input2[4:]
		require.Equal(input, input2)
		unpacked, err = UnpackModifyAllowListInput(input2, role, false)
		require.NoError(err)
		unpacked2, err := OldUnpackModifyAllowList(input, role)
		require.NoError(err)
		require.Equal(unpacked, unpacked2)
	})
}

func OldPackModifyAllowList(address common.Address, role Role) ([]byte, error) {
	// function selector (4 bytes) + hash for address
	input := make([]byte, 0, contract.SelectorLen+common.HashLength)

	switch role {
	case AdminRole:
		input = append(input, setAdminSignature...)
	case ManagerRole:
		input = append(input, setManagerSignature...)
	case EnabledRole:
		input = append(input, setEnabledSignature...)
	case NoRole:
		input = append(input, setNoneSignature...)
	default:
		return nil, fmt.Errorf("cannot pack modify list input with invalid role: %s", role)
	}

	input = append(input, address.Hash().Bytes()...)
	return input, nil
}

func OldUnpackModifyAllowList(input []byte, role Role) (common.Address, error) {
	if len(input) != allowListInputLen {
		return common.Address{}, fmt.Errorf("invalid input length for modifying allow list: %d", len(input))
	}
	return common.BytesToAddress(input), nil
}

func FuzzPackReadAllowListOutputTest(f *testing.F) {
	f.Fuzz(func(t *testing.T, roleIndex uint) {
		role := getRole(roleIndex)
		packedOutput, err := PackReadAllowListOutput(role.Big())
		require.NoError(t, err)
		require.Equal(t, packedOutput, role.Bytes())
	})
}

func getRole(roleIndex uint) Role {
	index := roleIndex % 4
	switch index {
	case 0:
		return NoRole
	case 1:
		return EnabledRole
	case 2:
		return AdminRole
	case 3:
		return ManagerRole
	default:
		panic("unknown role")
	}
}
