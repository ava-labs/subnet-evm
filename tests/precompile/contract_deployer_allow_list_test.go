// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"context"

	"github.com/ava-labs/libevm/accounts/abi/bind"
	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/types"

	"github.com/ava-labs/subnet-evm/contracts/bindings"
	"github.com/ava-labs/subnet-evm/tests/precompile/contracttest"
	"github.com/ava-labs/subnet-evm/tests/precompile/solidity"
	"github.com/ava-labs/subnet-evm/tests/utils"

	ginkgo "github.com/onsi/ginkgo/v2"
	requirePkg "github.com/stretchr/testify/require"
)

// Register Go-based contract deployer allow list tests
// Tests within this suite run serially to avoid nonce conflicts on the same subnet
func init() {
	ginkgo.Describe("[Go] Contract Deployer Allow List", ginkgo.Serial, ginkgo.Label("Precompile", "ContractDeployerAllowList", "Go"), func() {
		var (
			require        *requirePkg.Assertions
			backend        *contracttest.TmpnetBackend
			allowList      *bindings.IAllowList
			waitForReceipt func(tx *types.Transaction) *types.Receipt
		)

		ginkgo.BeforeEach(func() {
			require = requirePkg.New(ginkgo.GinkgoT())

			blockchainID := solidity.SubnetsSuite.GetBlockchainID("contract_deployer_allow_list")
			rpcURL := utils.GetDefaultChainURI(blockchainID)

			backend = contracttest.NewTmpnetBackend(ginkgo.GinkgoTB(), rpcURL)
			require.NotNil(backend, "failed to create tmpnet backend")

			var err error
			allowList, err = bindings.NewIAllowList(contracttest.ContractDeployerAllowListAddress, backend.Client)
			require.NoError(err, "failed to bind contract deployer allow list precompile")

			waitForReceipt = func(tx *types.Transaction) *types.Receipt {
				receipt, err := bind.WaitMined(context.Background(), backend.Client, tx)
				require.NoError(err)
				return receipt
			}
		})

		ginkgo.AfterEach(func() {
			if backend != nil {
				backend.Close()
				backend = nil
			}
			allowList = nil
			waitForReceipt = nil
			require = nil
		})

		ginkgo.It("should emit RoleSet event when setting admin role", func() {
			// TODO: revert to using the TypeScript address for this after migration (ensures OldRole is None)
			// testAddress := common.HexToAddress("0x0111000000000000000000000000000000000001")
			testAddress := common.HexToAddress("0x1111111111111111111111111111111111111111")

			tx, err := allowList.SetAdmin(backend.Admin.Auth, testAddress)
			require.NoError(err)

			receipt := waitForReceipt(tx)
			require.Equal(types.ReceiptStatusSuccessful, receipt.Status, "transaction failed")
			require.Len(receipt.Logs, 1, "should have exactly one log")

			event, err := allowList.ParseRoleSet(*receipt.Logs[0])
			require.NoError(err, "failed to parse RoleSet event")

			require.Equal(uint64(contracttest.RoleAdmin), event.Role.Uint64(), "role should be Admin")
			require.Equal(testAddress, event.Account, "account should match test address")
			require.Equal(backend.Admin.Address, event.Sender, "sender should be admin")
			require.Equal(uint64(contracttest.RoleNone), event.OldRole.Uint64(), "old role should be None")
		})

		ginkgo.It("should emit RoleSet event when setting manager role", func() {
			// TODO: revert to using the TypeScript address for this after migration (ensures OldRole is None)
			// testAddress := common.HexToAddress("0x0222000000000000000000000000000000000002")
			testAddress := common.HexToAddress("0x2222222222222222222222222222222222222222")

			tx, err := allowList.SetManager(backend.Admin.Auth, testAddress)
			require.NoError(err)

			receipt := waitForReceipt(tx)
			require.Equal(types.ReceiptStatusSuccessful, receipt.Status)
			require.Len(receipt.Logs, 1)

			event, err := allowList.ParseRoleSet(*receipt.Logs[0])
			require.NoError(err)

			require.Equal(uint64(contracttest.RoleManager), event.Role.Uint64())
			require.Equal(testAddress, event.Account)
			require.Equal(backend.Admin.Address, event.Sender)
			require.Equal(uint64(contracttest.RoleNone), event.OldRole.Uint64())
		})

		ginkgo.It("should emit RoleSet event when setting enabled role", func() {
			testAddress := common.HexToAddress("0x0333000000000000000000000000000000000003")

			tx, err := allowList.SetEnabled(backend.Admin.Auth, testAddress)
			require.NoError(err)

			receipt := waitForReceipt(tx)
			require.Equal(types.ReceiptStatusSuccessful, receipt.Status)
			require.Len(receipt.Logs, 1)

			event, err := allowList.ParseRoleSet(*receipt.Logs[0])
			require.NoError(err)

			require.Equal(uint64(contracttest.RoleEnabled), event.Role.Uint64())
			require.Equal(testAddress, event.Account)
			require.Equal(backend.Admin.Address, event.Sender)
			require.Equal(uint64(contracttest.RoleNone), event.OldRole.Uint64())
		})

		ginkgo.It("should emit RoleSet event when setting none role", func() {
			testAddress := common.HexToAddress("0x0333000000000000000000000000000000000003")

			tx, err := allowList.SetEnabled(backend.Admin.Auth, testAddress)
			require.NoError(err)
			receipt := waitForReceipt(tx)
			require.Equal(types.ReceiptStatusSuccessful, receipt.Status)

			tx, err = allowList.SetNone(backend.Admin.Auth, testAddress)
			require.NoError(err)

			receipt = waitForReceipt(tx)
			require.Equal(types.ReceiptStatusSuccessful, receipt.Status)
			require.Len(receipt.Logs, 1)

			event, err := allowList.ParseRoleSet(*receipt.Logs[0])
			require.NoError(err)

			require.Equal(uint64(contracttest.RoleNone), event.Role.Uint64())
			require.Equal(testAddress, event.Account)
			require.Equal(backend.Admin.Address, event.Sender)
			require.Equal(uint64(contracttest.RoleEnabled), event.OldRole.Uint64(), "old role should be Enabled")
		})

		ginkgo.It("should verify deployer list shows admin has admin role via precompile", func() {
			deployerListAddr, tx, deployerList, err := bindings.DeployExampleDeployerList(backend.Admin.Auth, backend.Client)
			require.NoError(err)

			receipt := waitForReceipt(tx)
			require.Equal(types.ReceiptStatusSuccessful, receipt.Status)

			tx, err = allowList.SetAdmin(backend.Admin.Auth, deployerListAddr)
			require.NoError(err)
			receipt = waitForReceipt(tx)
			require.Equal(types.ReceiptStatusSuccessful, receipt.Status)

			role, err := allowList.ReadAllowList(nil, deployerListAddr)
			require.NoError(err)
			require.Equal(uint64(contracttest.RoleAdmin), role.Uint64())

			isAdmin, err := deployerList.IsAdmin(nil, deployerListAddr)
			require.NoError(err)
			require.True(isAdmin, "contract should report itself as admin")
		})

		ginkgo.It("should verify new address has no role", func() {
			deployerListAddr, tx, _, err := bindings.DeployExampleDeployerList(backend.Admin.Auth, backend.Client)
			require.NoError(err)

			receipt := waitForReceipt(tx)
			require.Equal(types.ReceiptStatusSuccessful, receipt.Status)

			role, err := allowList.ReadAllowList(nil, deployerListAddr)
			require.NoError(err)
			require.Equal(uint64(contracttest.RoleNone), role.Uint64(), "new contract should have no role")
		})

		ginkgo.It("should verify contract correctly reports admin status", func() {
			deployerListAddr, tx, deployerList, err := bindings.DeployExampleDeployerList(backend.Admin.Auth, backend.Client)
			require.NoError(err)

			receipt := waitForReceipt(tx)
			require.Equal(types.ReceiptStatusSuccessful, receipt.Status)

			owner, err := deployerList.Owner(nil)
			require.NoError(err)
			require.Equal(backend.Admin.Address, owner)

			tx, err = allowList.SetAdmin(backend.Admin.Auth, backend.Admin.Address)
			require.NoError(err)
			_ = waitForReceipt(tx)

			isAdmin, err := deployerList.IsAdmin(nil, backend.Admin.Address)
			require.NoError(err)
			require.True(isAdmin, "owner should be admin")

			isAdmin, err = deployerList.IsAdmin(nil, deployerListAddr)
			require.NoError(err)
			require.False(isAdmin, "contract with no role should not be admin")
		})

		ginkgo.It("should not let address with no role deploy contracts", func() {
			_, tx, deployerList, err := bindings.DeployExampleDeployerList(backend.Admin.Auth, backend.Client)
			require.NoError(err)

			receipt := waitForReceipt(tx)
			require.Equal(types.ReceiptStatusSuccessful, receipt.Status)

			require.NotNil(backend.Unprivileged, "missing other test account")
			otherAuth := *backend.Unprivileged.Auth
			// Override the gas limit so bind skips preflight estimation; otherwise the call
			// reverts during estimation and no transaction is sent, which would hide the
			// failure status we expect from the mined receipt.
			otherAuth.GasLimit = 500000

			tx, err = deployerList.DeployContract(&otherAuth)
			require.NoError(err, "transaction submission should succeed")

			receipt = waitForReceipt(tx)
			require.Equal(types.ReceiptStatusFailed, receipt.Status, "deploy should fail for address with no role")
		})

		ginkgo.It("should allow admin to add contract as admin via precompile", func() {
			deployerListAddr, tx, deployerList, err := bindings.DeployExampleDeployerList(backend.Admin.Auth, backend.Client)
			require.NoError(err)

			receipt := waitForReceipt(tx)
			require.Equal(types.ReceiptStatusSuccessful, receipt.Status)

			role, err := allowList.ReadAllowList(nil, deployerListAddr)
			require.NoError(err)
			require.Equal(uint64(contracttest.RoleNone), role.Uint64())

			tx, err = allowList.SetAdmin(backend.Admin.Auth, deployerListAddr)
			require.NoError(err)
			_ = waitForReceipt(tx)

			role, err = allowList.ReadAllowList(nil, deployerListAddr)
			require.NoError(err)
			require.Equal(uint64(contracttest.RoleAdmin), role.Uint64())

			isAdmin, err := deployerList.IsAdmin(nil, deployerListAddr)
			require.NoError(err)
			require.True(isAdmin)
		})

		ginkgo.It("should allow admin to add deployer via contract", func() {
			exampleAddr, tx, example, err := bindings.DeployExampleDeployerList(backend.Admin.Auth, backend.Client)
			require.NoError(err)
			_ = waitForReceipt(tx)

			otherAddr, tx, other, err := bindings.DeployExampleDeployerList(backend.Admin.Auth, backend.Client)
			require.NoError(err)
			_ = waitForReceipt(tx)

			tx, err = allowList.SetAdmin(backend.Admin.Auth, exampleAddr)
			require.NoError(err)
			_ = waitForReceipt(tx)

			tx, err = example.SetEnabled(backend.Admin.Auth, otherAddr)
			require.NoError(err)
			_ = waitForReceipt(tx)

			isEnabled, err := example.IsEnabled(nil, otherAddr)
			require.NoError(err)
			require.True(isEnabled, "other contract should be enabled")

			isEnabledSelf, err := other.IsEnabled(nil, otherAddr)
			require.NoError(err)
			require.True(isEnabledSelf, "contract should report itself as enabled")
		})

		ginkgo.It("should allow enabled address to deploy contracts", func() {
			exampleAddr, tx, example, err := bindings.DeployExampleDeployerList(backend.Admin.Auth, backend.Client)
			require.NoError(err)
			_ = waitForReceipt(tx)

			deployerAddr, tx, deployer, err := bindings.DeployExampleDeployerList(backend.Admin.Auth, backend.Client)
			require.NoError(err)
			_ = waitForReceipt(tx)

			tx, err = allowList.SetAdmin(backend.Admin.Auth, exampleAddr)
			require.NoError(err)
			_ = waitForReceipt(tx)

			tx, err = example.SetEnabled(backend.Admin.Auth, deployerAddr)
			require.NoError(err)
			_ = waitForReceipt(tx)

			isEnabled, err := example.IsEnabled(nil, deployerAddr)
			require.NoError(err)
			require.True(isEnabled)

			tx, err = deployer.DeployContract(backend.Admin.Auth)
			require.NoError(err)

			receipt := waitForReceipt(tx)
			require.Equal(types.ReceiptStatusSuccessful, receipt.Status, "enabled address should be able to deploy")
		})

		ginkgo.It("should allow admin to revoke deployer role", func() {
			exampleAddr, tx, example, err := bindings.DeployExampleDeployerList(backend.Admin.Auth, backend.Client)
			require.NoError(err)
			_ = waitForReceipt(tx)

			deployerAddr, tx, _, err := bindings.DeployExampleDeployerList(backend.Admin.Auth, backend.Client)
			require.NoError(err)
			_ = waitForReceipt(tx)

			tx, err = allowList.SetAdmin(backend.Admin.Auth, exampleAddr)
			require.NoError(err)
			_ = waitForReceipt(tx)

			tx, err = example.SetEnabled(backend.Admin.Auth, deployerAddr)
			require.NoError(err)
			_ = waitForReceipt(tx)

			isEnabled, err := example.IsEnabled(nil, deployerAddr)
			require.NoError(err)
			require.True(isEnabled)

			tx, err = example.Revoke(backend.Admin.Auth, deployerAddr)
			require.NoError(err)
			_ = waitForReceipt(tx)

			role, err := allowList.ReadAllowList(nil, deployerAddr)
			require.NoError(err)
			require.Equal(uint64(contracttest.RoleNone), role.Uint64(), "deployer should have no role after revoke")
		})
	})
}
