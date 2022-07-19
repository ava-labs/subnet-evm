// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/snow/engine/common"
	"github.com/ava-labs/avalanchego/vms/components/chain"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/stretchr/testify/assert"
)

func TestVMUpgradeBytesPrecompile(t *testing.T) {
	// Get a json specifying a TxAllowListConfig upgrade at genesis
	// to apply as upgradeBytes.
	enableAllowListTimestamp := time.Unix(0, 0) // enable at genesis
	upgradeBytesConfig := &params.UpgradeBytesConfig{
		PrecompileUpgrades: []precompile.Upgrade{
			{
				TxAllowListConfig: &precompile.TxAllowListConfig{
					UpgradeableConfig: precompile.UpgradeableConfig{
						BlockTimestamp: big.NewInt(enableAllowListTimestamp.Unix()),
					},
					AllowListConfig: precompile.AllowListConfig{
						AllowListAdmins: testEthAddrs[0:1],
					},
				},
			},
		},
	}
	upgradeBytesJSON, err := json.Marshal(upgradeBytesConfig)
	if err != nil {
		t.Fatalf("could not marshal upgradeBytesConfig to json: %s", err)
	}

	// initialize the VM with these upgrade bytes
	issuer, vm, dbManager, appSender := GenesisVM(t, true, genesisJSONSubnetEVM, "", string(upgradeBytesJSON))

	// Submit a successful transaction
	tx0 := types.NewTransaction(uint64(0), testEthAddrs[0], big.NewInt(1), 21000, big.NewInt(testMinGasPrice), nil)
	signedTx0, err := types.SignTx(tx0, types.NewEIP155Signer(vm.chainConfig.ChainID), testKeys[0])
	assert.NoError(t, err)

	err = vm.chain.GetTxPool().AddRemote(signedTx0)
	if err != nil {
		t.Fatalf("Failed to add tx at index: %s", err)
	}

	// Submit a rejected transaction, should throw an error
	tx1 := types.NewTransaction(uint64(0), testEthAddrs[1], big.NewInt(2), 21000, big.NewInt(testMinGasPrice), nil)
	signedTx1, err := types.SignTx(tx1, types.NewEIP155Signer(vm.chainConfig.ChainID), testKeys[1])
	if err != nil {
		t.Fatal(err)
	}
	err = vm.chain.GetTxPool().AddRemote(signedTx1)
	if !errors.Is(err, precompile.ErrSenderAddressNotAllowListed) {
		t.Fatalf("expected ErrSenderAddressNotAllowListed, got: %s", err)
	}

	// shutdown the vm
	if err := vm.Shutdown(); err != nil {
		t.Fatal(err)
	}

	// prepare the new upgrade bytes to disable the TxAllowList
	disableAllowListTimestamp := enableAllowListTimestamp.Add(10 * time.Hour) // arbitrary choice
	upgradeBytesConfig.PrecompileUpgrades = append(
		upgradeBytesConfig.PrecompileUpgrades,
		precompile.Upgrade{
			TxAllowListConfig: &precompile.TxAllowListConfig{
				UpgradeableConfig: precompile.UpgradeableConfig{
					BlockTimestamp: big.NewInt(disableAllowListTimestamp.Unix()),
					Disable:        true,
				},
			},
		},
	)
	upgradeBytesJSON, err = json.Marshal(upgradeBytesConfig)
	if err != nil {
		t.Fatalf("could not marshal upgradeBytesConfig to json: %s", err)
	}

	// restart the vm
	ctx := NewContext()
	if err := vm.Initialize(
		ctx, dbManager, []byte(genesisJSONSubnetEVM), upgradeBytesJSON, []byte{}, issuer, []*common.Fx{}, appSender,
	); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := vm.Shutdown(); err != nil {
			t.Fatal(err)
		}
	}()
	newTxPoolHeadChan := make(chan core.NewTxPoolReorgEvent, 1)
	vm.chain.GetTxPool().SubscribeNewReorgEvent(newTxPoolHeadChan)
	vm.clock.Set(disableAllowListTimestamp)

	// Make a block, previous rules still apply (TxAllowList is active)
	// Submit a successful transaction
	err = vm.chain.GetTxPool().AddRemote(signedTx0)
	if err != nil {
		t.Fatalf("Failed to add tx at index: %s", err)
	}

	// Submit a rejected transaction, should throw an error
	err = vm.chain.GetTxPool().AddRemote(signedTx1)
	if !errors.Is(err, precompile.ErrSenderAddressNotAllowListed) {
		t.Fatalf("expected ErrSenderAddressNotAllowListed, got: %s", err)
	}

	blk := issueAndAccept(t, issuer, vm)

	// Verify that the constructed block only has the whitelisted tx
	block := blk.(*chain.BlockWrapper).Block.(*Block).ethBlock
	txs := block.Transactions()
	if txs.Len() != 1 {
		t.Fatalf("Expected number of txs to be %d, but found %d", 1, txs.Len())
	}
	assert.Equal(t, signedTx0.Hash(), txs[0].Hash())

	// verify the issued block is after the network upgrade
	assert.True(t, block.Timestamp().Cmp(big.NewInt(disableAllowListTimestamp.Unix())) >= 0)

	<-newTxPoolHeadChan // wait for new head in tx pool

	// retry the rejected Tx, which should now succeed
	errs := vm.chain.GetTxPool().AddRemotesSync([]*types.Transaction{signedTx1})
	if err := errs[0]; err != nil {
		t.Fatalf("Failed to add tx at index: %s", err)
	}

	vm.clock.Set(vm.clock.Time().Add(2 * time.Second)) // add 2 seconds for gas fee to adjust
	blk = issueAndAccept(t, issuer, vm)

	// Verify that the constructed block only has the previously rejected tx
	block = blk.(*chain.BlockWrapper).Block.(*Block).ethBlock
	txs = block.Transactions()
	if txs.Len() != 1 {
		t.Fatalf("Expected number of txs to be %d, but found %d", 1, txs.Len())
	}
	assert.Equal(t, signedTx1.Hash(), txs[0].Hash())
}

func TestVMUpgradeBytesNetworkUpgrades(t *testing.T) {
	// Get a json specifying a Network upgrade at genesis
	// to apply as upgradeBytes.
	subnetEVMTimestamp := time.Unix(10, 0)
	upgradeBytesConfig := &params.UpgradeBytesConfig{
		NetworkUpgrades: &params.NetworkUpgrades{
			SubnetEVMTimestamp: big.NewInt(subnetEVMTimestamp.Unix()),
		},
	}
	upgradeBytesJSON, err := json.Marshal(upgradeBytesConfig)
	if err != nil {
		t.Fatalf("could not marshal upgradeBytesConfig to json: %s", err)
	}

	// initialize the VM with these upgrade bytes
	issuer, vm, dbManager, appSender := GenesisVM(t, true, genesisJSONPreSubnetEVM, "", string(upgradeBytesJSON))
	vm.clock.Set(subnetEVMTimestamp)

	// verify upgrade is applied
	if !vm.chainConfig.IsSubnetEVM(big.NewInt(subnetEVMTimestamp.Unix())) {
		t.Fatal("expected subnet-evm network upgrade to have been enabled")
	}

	// Submit a successful transaction
	tx0 := types.NewTransaction(uint64(0), testEthAddrs[0], big.NewInt(1), 21000, big.NewInt(testMinGasPrice), nil)
	signedTx0, err := types.SignTx(tx0, types.NewEIP155Signer(vm.chainConfig.ChainID), testKeys[0])
	assert.NoError(t, err)
	err = vm.chain.GetTxPool().AddRemote(signedTx0)
	if err != nil {
		t.Fatalf("Failed to add tx at index: %s", err)
	}

	issueAndAccept(t, issuer, vm) // make a block

	if err := vm.Shutdown(); err != nil {
		t.Fatal(err)
	}

	// VM should not start again without proper upgrade bytes.
	ctx := NewContext()
	err = vm.Initialize(ctx, dbManager, []byte(genesisJSONPreSubnetEVM), []byte{}, []byte{}, issuer, []*common.Fx{}, appSender)
	assert.ErrorContains(t, err, "mismatching SubnetEVM fork block timestamp in database")

	// VM should not start if fork is moved back
	upgradeBytesConfig.NetworkUpgrades.SubnetEVMTimestamp = big.NewInt(0)
	upgradeBytesJSON, err = json.Marshal(upgradeBytesConfig)
	if err != nil {
		t.Fatalf("could not marshal upgradeBytesConfig to json: %s", err)
	}
	err = vm.Initialize(ctx, dbManager, []byte(genesisJSONPreSubnetEVM), upgradeBytesJSON, []byte{}, issuer, []*common.Fx{}, appSender)
	fmt.Println(err)
	assert.ErrorContains(t, err, "mismatching SubnetEVM fork block timestamp in database")

	// VM should not start if fork is moved forward
	upgradeBytesConfig.NetworkUpgrades.SubnetEVMTimestamp = big.NewInt(30)
	upgradeBytesJSON, err = json.Marshal(upgradeBytesConfig)
	if err != nil {
		t.Fatalf("could not marshal upgradeBytesConfig to json: %s", err)
	}
	err = vm.Initialize(ctx, dbManager, []byte(genesisJSONPreSubnetEVM), upgradeBytesJSON, []byte{}, issuer, []*common.Fx{}, appSender)
	fmt.Println(err)
	assert.ErrorContains(t, err, "mismatching SubnetEVM fork block timestamp in database")
}

func TestVMUpgradeBytesNetworkUpgradesWithGenesis(t *testing.T) {
	// make genesis w/ fork at block 5
	var genesis core.Genesis
	if err := json.Unmarshal([]byte(genesisJSONPreSubnetEVM), &genesis); err != nil {
		t.Fatalf("could not unmarshal genesis bytes: %s", err)
	}
	genesisSubnetEVMTimestamp := big.NewInt(5)
	genesis.Config.SubnetEVMTimestamp = genesisSubnetEVMTimestamp
	genesisBytes, err := json.Marshal(&genesis)
	if err != nil {
		t.Fatalf("could not unmarshal genesis bytes: %s", err)
	}

	// Get a json specifying a Network upgrade at genesis
	// to apply as upgradeBytes.
	subnetEVMTimestamp := time.Unix(10, 0)
	upgradeBytesConfig := &params.UpgradeBytesConfig{
		NetworkUpgrades: &params.NetworkUpgrades{
			SubnetEVMTimestamp: big.NewInt(subnetEVMTimestamp.Unix()),
		},
	}
	upgradeBytesJSON, err := json.Marshal(upgradeBytesConfig)
	if err != nil {
		t.Fatalf("could not marshal upgradeBytesConfig to json: %s", err)
	}

	// initialize the VM with these upgrade bytes
	_, vm, _, _ := GenesisVM(t, true, string(genesisBytes), "", string(upgradeBytesJSON))

	// verify upgrade is rescheduled
	assert.False(t, vm.chainConfig.IsSubnetEVM(genesisSubnetEVMTimestamp))
	assert.True(t, vm.chainConfig.IsSubnetEVM(big.NewInt(subnetEVMTimestamp.Unix())))

	if err := vm.Shutdown(); err != nil {
		t.Fatal(err)
	}

	// abort a fork specified in genesis
	upgradeBytesConfig.NetworkUpgrades.SubnetEVMTimestamp = nil
	upgradeBytesJSON, err = json.Marshal(upgradeBytesConfig)
	if err != nil {
		t.Fatalf("could not marshal upgradeBytesConfig to json: %s", err)
	}

	// initialize the VM with these upgrade bytes
	_, vm, _, _ = GenesisVM(t, true, string(genesisBytes), "", string(upgradeBytesJSON))

	// verify upgrade is aborted
	assert.False(t, vm.chainConfig.IsSubnetEVM(genesisSubnetEVMTimestamp))
	assert.False(t, vm.chainConfig.IsSubnetEVM(big.NewInt(subnetEVMTimestamp.Unix())))

	if err := vm.Shutdown(); err != nil {
		t.Fatal(err)
	}
}
