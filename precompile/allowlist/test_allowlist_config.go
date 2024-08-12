// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package allowlist

import (
	"encoding/json"
	"testing"

	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/mock/gomock"
)

// mkConfigWithAllowList creates a new config with the correct type for [module]
// by marshalling [cfg] to JSON and then unmarshalling it into the config.
func mkConfigWithAllowList(module modules.Module, cfg *AllowListConfig) precompileconfig.Config {
	jsonBytes, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}

	moduleCfg := module.MakeConfig()
	err = json.Unmarshal(jsonBytes, moduleCfg)
	if err != nil {
		panic(err)
	}

	return moduleCfg
}

func mkConfigWithUpgradeAndAllowList(module modules.Module, cfg *AllowListConfig, update precompileconfig.Upgrade) precompileconfig.Config {
	jsonUpgradeBytes, err := json.Marshal(update)
	if err != nil {
		panic(err)
	}

	moduleCfg := mkConfigWithAllowList(module, cfg)
	err = json.Unmarshal(jsonUpgradeBytes, moduleCfg)
	if err != nil {
		panic(err)
	}
	return moduleCfg
}

func AllowListConfigVerifyTests(t testing.TB, module modules.Module) map[string]testutils.ConfigVerifyTest {
	return map[string]testutils.ConfigVerifyTest{
		"invalid allow list config with duplicate admins in allowlist": {
			Config: mkConfigWithAllowList(module, &AllowListConfig{
				AdminAddresses:   []common.Address{TestAdminAddr, TestAdminAddr},
				ManagerAddresses: nil,
				EnabledAddresses: nil,
			}),
			ExpectedError: "duplicate address in admin list",
		},
		"invalid allow list config with duplicate enableds in allowlist": {
			Config: mkConfigWithAllowList(module, &AllowListConfig{
				AdminAddresses:   nil,
				ManagerAddresses: nil,
				EnabledAddresses: []common.Address{TestEnabledAddr, TestEnabledAddr},
			}),
			ExpectedError: "duplicate address in enabled list",
		},
		"invalid allow list config with duplicate managers in allowlist": {
			Config: mkConfigWithAllowList(module, &AllowListConfig{
				AdminAddresses:   nil,
				ManagerAddresses: []common.Address{TestManagerAddr, TestManagerAddr},
				EnabledAddresses: nil,
			}),
			ExpectedError: "duplicate address in manager list",
		},
		"invalid allow list config with same admin and enabled in allowlist": {
			Config: mkConfigWithAllowList(module, &AllowListConfig{
				AdminAddresses:   []common.Address{TestAdminAddr},
				ManagerAddresses: nil,
				EnabledAddresses: []common.Address{TestAdminAddr},
			}),
			ExpectedError: "cannot set address as both admin and enabled",
		},
		"invalid allow list config with same admin and manager in allowlist": {
			Config: mkConfigWithAllowList(module, &AllowListConfig{
				AdminAddresses:   []common.Address{TestAdminAddr},
				ManagerAddresses: []common.Address{TestAdminAddr},
				EnabledAddresses: nil,
			}),
			ExpectedError: "cannot set address as both admin and manager",
		},
		"invalid allow list config with same manager and enabled in allowlist": {
			Config: mkConfigWithAllowList(module, &AllowListConfig{
				AdminAddresses:   nil,
				ManagerAddresses: []common.Address{TestManagerAddr},
				EnabledAddresses: []common.Address{TestManagerAddr},
			}),
			ExpectedError: "cannot set address as both enabled and manager",
		},
		"invalid allow list config with manager role before activation": {
			Config: mkConfigWithUpgradeAndAllowList(module, &AllowListConfig{
				AdminAddresses:   nil,
				ManagerAddresses: []common.Address{TestManagerAddr},
				EnabledAddresses: nil,
			}, precompileconfig.Upgrade{
				BlockTimestamp: utils.NewUint64(1),
			}),
			ChainConfig: func() precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(gomock.NewController(t))
				config.EXPECT().IsDurango(gomock.Any()).Return(false)
				return config
			}(),
			ExpectedError: ErrCannotAddManagersBeforeDurango.Error(),
		},
		"nil member allow list config in allowlist": {
			Config: mkConfigWithAllowList(module, &AllowListConfig{
				AdminAddresses:   nil,
				ManagerAddresses: nil,
				EnabledAddresses: nil,
			}),
			ExpectedError: "",
		},
		"empty member allow list config in allowlist": {
			Config: mkConfigWithAllowList(module, &AllowListConfig{
				AdminAddresses:   []common.Address{},
				ManagerAddresses: []common.Address{},
				EnabledAddresses: []common.Address{},
			}),
			ExpectedError: "",
		},
		"valid allow list config in allowlist": {
			Config: mkConfigWithAllowList(module, &AllowListConfig{
				AdminAddresses:   []common.Address{TestAdminAddr},
				ManagerAddresses: []common.Address{TestManagerAddr},
				EnabledAddresses: []common.Address{TestEnabledAddr},
			}),
			ExpectedError: "",
		},
	}
}

func AllowListConfigEqualTests(t testing.TB, module modules.Module) map[string]testutils.ConfigEqualTest {
	return map[string]testutils.ConfigEqualTest{
		"allowlist non-nil config and nil other": {
			Config: mkConfigWithAllowList(module, &AllowListConfig{
				AdminAddresses:   []common.Address{TestAdminAddr},
				ManagerAddresses: []common.Address{TestManagerAddr},
				EnabledAddresses: []common.Address{TestEnabledAddr},
			}),
			Other:    nil,
			Expected: false,
		},
		"allowlist different admin": {
			Config: mkConfigWithAllowList(module, &AllowListConfig{
				AdminAddresses:   []common.Address{TestAdminAddr},
				ManagerAddresses: []common.Address{TestManagerAddr},
				EnabledAddresses: []common.Address{TestEnabledAddr},
			}),
			Other: mkConfigWithAllowList(module, &AllowListConfig{
				AdminAddresses:   []common.Address{{3}},
				ManagerAddresses: []common.Address{TestManagerAddr},
				EnabledAddresses: []common.Address{TestEnabledAddr},
			}),
			Expected: false,
		},
		"allowlist different manager": {
			Config: mkConfigWithAllowList(module, &AllowListConfig{
				AdminAddresses:   []common.Address{TestAdminAddr},
				ManagerAddresses: []common.Address{TestManagerAddr},
				EnabledAddresses: []common.Address{TestEnabledAddr},
			}),
			Other: mkConfigWithAllowList(module, &AllowListConfig{
				AdminAddresses:   []common.Address{TestAdminAddr},
				ManagerAddresses: []common.Address{{3}},
				EnabledAddresses: []common.Address{TestEnabledAddr},
			}),
			Expected: false,
		},
		"allowlist different enabled": {
			Config: mkConfigWithAllowList(module, &AllowListConfig{
				AdminAddresses:   []common.Address{TestAdminAddr},
				ManagerAddresses: []common.Address{TestManagerAddr},
				EnabledAddresses: []common.Address{TestEnabledAddr},
			}),
			Other: mkConfigWithAllowList(module, &AllowListConfig{
				AdminAddresses:   []common.Address{TestAdminAddr},
				ManagerAddresses: []common.Address{TestManagerAddr},
				EnabledAddresses: []common.Address{{3}},
			}),
			Expected: false,
		},
		"allowlist same config": {
			Config: mkConfigWithAllowList(module, &AllowListConfig{
				AdminAddresses:   []common.Address{TestAdminAddr},
				ManagerAddresses: []common.Address{TestManagerAddr},
				EnabledAddresses: []common.Address{TestEnabledAddr},
			}),
			Other: mkConfigWithAllowList(module, &AllowListConfig{
				AdminAddresses:   []common.Address{TestAdminAddr},
				ManagerAddresses: []common.Address{TestManagerAddr},
				EnabledAddresses: []common.Address{TestEnabledAddr},
			}),
			Expected: true,
		},
	}
}

func VerifyPrecompileWithAllowListTests(t *testing.T, module modules.Module, verifyTests map[string]testutils.ConfigVerifyTest) {
	t.Helper()
	tests := AllowListConfigVerifyTests(t, module)
	// Add the contract specific tests to the map of tests to run.
	for name, test := range verifyTests {
		if _, exists := tests[name]; exists {
			t.Fatalf("duplicate test name: %s", name)
		}
		tests[name] = test
	}

	testutils.RunVerifyTests(t, tests)
}

func EqualPrecompileWithAllowListTests(t *testing.T, module modules.Module, equalTests map[string]testutils.ConfigEqualTest) {
	t.Helper()
	tests := AllowListConfigEqualTests(t, module)
	// Add the contract specific tests to the map of tests to run.
	for name, test := range equalTests {
		if _, exists := tests[name]; exists {
			t.Fatalf("duplicate test name: %s", name)
		}
		tests[name] = test
	}

	testutils.RunEqualTests(t, tests)
}
