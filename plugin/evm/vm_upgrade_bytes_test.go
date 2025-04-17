// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/api/metrics"
	"github.com/ava-labs/avalanchego/snow"
	commonEng "github.com/ava-labs/avalanchego/snow/engine/common"
	"github.com/ava-labs/avalanchego/snow/engine/enginetest"
	"github.com/ava-labs/avalanchego/upgrade"
	"github.com/ava-labs/avalanchego/vms/components/chain"
	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/common/hexutil"
	"github.com/ava-labs/libevm/common/math"
	"github.com/ava-labs/libevm/core/types"
	"github.com/ava-labs/libevm/crypto"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/params/extras"
	"github.com/ava-labs/subnet-evm/plugin/evm/vmerrors"
	"github.com/ava-labs/subnet-evm/precompile/contracts/txallowlist"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var DefaultEtnaTime = uint64(upgrade.GetConfig(testNetworkID).EtnaTime.Unix())

func TestVMUpgradeBytesPrecompile(t *testing.T) {
	// Make a TxAllowListConfig upgrade at genesis and convert it to JSON to apply as upgradeBytes.
	enableAllowListTimestamp := upgrade.InitiallyActiveTime // enable at initial time
	upgradeConfig := &extras.UpgradeConfig{
		PrecompileUpgrades: []extras.PrecompileUpgrade{
			{
				Config: txallowlist.NewConfig(utils.TimeToNewUint64(enableAllowListTimestamp), testEthAddrs[0:1], nil, nil),
			},
		},
	}
	upgradeBytesJSON, err := json.Marshal(upgradeConfig)
	if err != nil {
		t.Fatalf("could not marshal upgradeConfig to json: %s", err)
	}

	// initialize the VM with these upgrade bytes
	issuer, vm, dbManager, appSender := GenesisVM(t, true, genesisJSONSubnetEVM, "", string(upgradeBytesJSON))
	vm.clock.Set(enableAllowListTimestamp)

	// Submit a successful transaction
	tx0 := types.NewTransaction(uint64(0), testEthAddrs[0], big.NewInt(1), 21000, big.NewInt(testMinGasPrice), nil)
	signedTx0, err := types.SignTx(tx0, types.NewEIP155Signer(vm.chainConfig.ChainID), testKeys[0])
	assert.NoError(t, err)

	errs := vm.txPool.AddRemotesSync([]*types.Transaction{signedTx0})
	if err := errs[0]; err != nil {
		t.Fatalf("Failed to add tx at index: %s", err)
	}

	// Submit a rejected transaction, should throw an error
	tx1 := types.NewTransaction(uint64(0), testEthAddrs[1], big.NewInt(2), 21000, big.NewInt(testMinGasPrice), nil)
	signedTx1, err := types.SignTx(tx1, types.NewEIP155Signer(vm.chainConfig.ChainID), testKeys[1])
	if err != nil {
		t.Fatal(err)
	}
	errs = vm.txPool.AddRemotesSync([]*types.Transaction{signedTx1})
	if err := errs[0]; !errors.Is(err, vmerrors.ErrSenderAddressNotAllowListed) {
		t.Fatalf("expected ErrSenderAddressNotAllowListed, got: %s", err)
	}

	// shutdown the vm
	if err := vm.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}

	// prepare the new upgrade bytes to disable the TxAllowList
	disableAllowListTimestamp := vm.clock.Time().Add(10 * time.Hour) // arbitrary choice
	upgradeConfig.PrecompileUpgrades = append(
		upgradeConfig.PrecompileUpgrades,
		extras.PrecompileUpgrade{
			Config: txallowlist.NewDisableConfig(utils.TimeToNewUint64(disableAllowListTimestamp)),
		},
	)
	upgradeBytesJSON, err = json.Marshal(upgradeConfig)
	if err != nil {
		t.Fatalf("could not marshal upgradeConfig to json: %s", err)
	}

	// restart the vm

	// Reset metrics to allow re-initialization
	vm.ctx.Metrics = metrics.NewPrefixGatherer()

	if err := vm.Initialize(
		context.Background(), vm.ctx, dbManager, []byte(genesisJSONSubnetEVM), upgradeBytesJSON, []byte{}, issuer, []*commonEng.Fx{}, appSender,
	); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
	}()
	// Set the VM's state to NormalOp to initialize the tx pool.
	if err := vm.SetState(context.Background(), snow.Bootstrapping); err != nil {
		t.Fatal(err)
	}
	if err := vm.SetState(context.Background(), snow.NormalOp); err != nil {
		t.Fatal(err)
	}
	newTxPoolHeadChan := make(chan core.NewTxPoolReorgEvent, 1)
	vm.txPool.SubscribeNewReorgEvent(newTxPoolHeadChan)
	vm.clock.Set(disableAllowListTimestamp)

	// Make a block, previous rules still apply (TxAllowList is active)
	// Submit a successful transaction
	errs = vm.txPool.AddRemotesSync([]*types.Transaction{signedTx0})
	if err := errs[0]; err != nil {
		t.Fatalf("Failed to add tx at index: %s", err)
	}

	// Submit a rejected transaction, should throw an error
	errs = vm.txPool.AddRemotesSync([]*types.Transaction{signedTx1})
	if err := errs[0]; !errors.Is(err, vmerrors.ErrSenderAddressNotAllowListed) {
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
	assert.GreaterOrEqual(t, int64(block.Time()), disableAllowListTimestamp.Unix())

	<-newTxPoolHeadChan // wait for new head in tx pool

	// retry the rejected Tx, which should now succeed
	errs = vm.txPool.AddRemotesSync([]*types.Transaction{signedTx1})
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

func TestNetworkUpgradesOverriden(t *testing.T) {
	var genesis core.Genesis
	if err := json.Unmarshal([]byte(genesisJSONPreSubnetEVM), &genesis); err != nil {
		t.Fatalf("could not unmarshal genesis bytes: %s", err)
	}
	genesisBytes, err := json.Marshal(&genesis)
	if err != nil {
		t.Fatalf("could not unmarshal genesis bytes: %s", err)
	}

	upgradeBytesJSON := `{
			"networkUpgradeOverrides": {
				"subnetEVMTimestamp": 2,
				"durangoTimestamp": 1607144402
			}
		}`

	vm := &VM{}
	ctx, dbManager, genesisBytes, issuer, _ := setupGenesis(t, string(genesisBytes))
	appSender := &enginetest.Sender{T: t}
	appSender.CantSendAppGossip = true
	appSender.SendAppGossipF = func(context.Context, commonEng.SendConfig, []byte) error { return nil }
	err = vm.Initialize(
		context.Background(),
		ctx,
		dbManager,
		genesisBytes,
		[]byte(upgradeBytesJSON),
		nil,
		issuer,
		[]*commonEng.Fx{},
		appSender,
	)
	require.NoError(t, err, "error initializing GenesisVM")

	require.NoError(t, vm.SetState(context.Background(), snow.Bootstrapping))
	require.NoError(t, vm.SetState(context.Background(), snow.NormalOp))

	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
	}()

	// verify upgrade overrides
	require.False(t, vm.chainConfigExtra().IsSubnetEVM(0))
	require.True(t, vm.chainConfigExtra().IsSubnetEVM(2))
	require.False(t, vm.chainConfigExtra().IsDurango(0))
	require.False(t, vm.chainConfigExtra().IsDurango(uint64(upgrade.InitiallyActiveTime.Unix())))
	require.True(t, vm.chainConfigExtra().IsDurango(1607144402))
}

func mustMarshal(t *testing.T, v interface{}) string {
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return string(b)
}

func TestVMStateUpgrade(t *testing.T) {
	// modify genesis to add a key to the state
	genesis := &core.Genesis{}
	err := json.Unmarshal([]byte(genesisJSONSubnetEVM), genesis)
	require.NoError(t, err)
	genesisAccount, ok := genesis.Alloc[testEthAddrs[0]]
	require.True(t, ok)
	storageKey := common.HexToHash("0x1234")
	genesisAccount.Storage = map[common.Hash]common.Hash{storageKey: common.HexToHash("0x5555")}
	genesisCode, err := hexutil.Decode("0xabcd")
	require.NoError(t, err)
	genesisAccount.Code = genesisCode
	genesisAccount.Nonce = 2                        // set to a non-zero value to test that it is preserved
	genesis.Alloc[testEthAddrs[0]] = genesisAccount // have to assign this back to the map for changes to take effect.
	genesisStr := mustMarshal(t, genesis)

	upgradedCodeStr := "0xdeadbeef" // this code will be applied during the upgrade
	upgradedCode, err := hexutil.Decode(upgradedCodeStr)
	// This modification will be applied to an existing account
	genesisAccountUpgrade := &extras.StateUpgradeAccount{
		BalanceChange: (*math.HexOrDecimal256)(big.NewInt(100)),
		Storage:       map[common.Hash]common.Hash{storageKey: {}},
		Code:          upgradedCode,
	}

	// This modification will be applied to a new account
	newAccount := common.Address{42}
	require.NoError(t, err)
	newAccountUpgrade := &extras.StateUpgradeAccount{
		BalanceChange: (*math.HexOrDecimal256)(big.NewInt(100)),
		Storage:       map[common.Hash]common.Hash{storageKey: common.HexToHash("0x6666")},
		Code:          upgradedCode,
	}

	// Create a new account just for testing nonce setting
	nonceTestAccount := common.Address{99}
	explicitNonce := uint64(10)
	nonceTestAccountUpgrade := &params.StateUpgradeAccount{
		BalanceChange: (*math.HexOrDecimal256)(big.NewInt(50)),
		Nonce:         &explicitNonce,
		// No code for this account - should still set the nonce directly
	}

	upgradeTimestamp := upgrade.InitiallyActiveTime.Add(10 * time.Hour)
	upgradeBytesJSON := fmt.Sprintf(
		`{
			"stateUpgrades": [
				{
					"blockTimestamp": %d,
					"accounts": {
						"%s": %s,
						"%s": %s,
						"%s": %s
					}
				}
			]
		}`,
		upgradeTimestamp.Unix(),
		testEthAddrs[0].Hex(),
		mustMarshal(t, genesisAccountUpgrade),
		newAccount.Hex(),
		mustMarshal(t, newAccountUpgrade),
		nonceTestAccount.Hex(),
		mustMarshal(t, nonceTestAccountUpgrade),
	)
	require.Contains(t, upgradeBytesJSON, upgradedCodeStr)

	// initialize the VM with these upgrade bytes
	issuer, vm, _, _ := GenesisVM(t, true, genesisStr, "", upgradeBytesJSON)
	defer func() { require.NoError(t, vm.Shutdown(context.Background())) }()

	// Verify the new account doesn't exist yet
	genesisState, err := vm.blockChain.State()
	require.NoError(t, err)
	require.Equal(t, common.U2560, genesisState.GetBalance(newAccount))

	// Advance the chain to the upgrade time
	vm.clock.Set(upgradeTimestamp)

	// Submit a successful (unrelated) transaction, so we can build a block
	// in this tx, testEthAddrs[1] sends 1 wei to itself.
	tx0 := types.NewTransaction(uint64(0), testEthAddrs[1], big.NewInt(1), 21000, big.NewInt(testMinGasPrice), nil)
	signedTx0, err := types.SignTx(tx0, types.NewEIP155Signer(vm.chainConfig.ChainID), testKeys[1])
	require.NoError(t, err)

	errs := vm.txPool.AddRemotesSync([]*types.Transaction{signedTx0})
	require.NoError(t, errs[0], "Failed to add tx")

	blk := issueAndAccept(t, issuer, vm)
	require.NotNil(t, blk)
	require.EqualValues(t, 1, blk.Height())

	// Verify the state upgrade was applied
	state, err := vm.blockChain.State()
	require.NoError(t, err)

	// Existing account
	expectedGenesisAccountBalance := uint256.MustFromBig(
		new(big.Int).Add(
			genesisAccount.Balance,
			(*big.Int)(genesisAccountUpgrade.BalanceChange),
		),
	)
	require.Equal(t, state.GetBalance(testEthAddrs[0]), expectedGenesisAccountBalance)
	require.Equal(t, state.GetState(testEthAddrs[0], storageKey), genesisAccountUpgrade.Storage[storageKey])
	require.Equal(t, state.GetCode(testEthAddrs[0]), upgradedCode)
	require.Equal(t, state.GetCodeHash(testEthAddrs[0]), crypto.Keccak256Hash(upgradedCode))
	require.Equal(t, state.GetNonce(testEthAddrs[0]), genesisAccount.Nonce) // Nonce should be preserved since it was non-zero

	// New account
	expectedNewAccountBalance := uint256.MustFromBig((*big.Int)(newAccountUpgrade.BalanceChange))
	require.Equal(t, state.GetBalance(newAccount), expectedNewAccountBalance)
	require.Equal(t, state.GetCode(newAccount), upgradedCode)
	require.Equal(t, state.GetCodeHash(newAccount), crypto.Keccak256Hash(upgradedCode))
	require.Equal(t, state.GetNonce(newAccount), uint64(1)) // Nonce should be set to 1 when code is set if nonce was 0
	require.Equal(t, state.GetState(newAccount, storageKey), newAccountUpgrade.Storage[storageKey])

	// Test the nonce-specific account
	expectedNonceTestAccountBalance := uint256.MustFromBig((*big.Int)(nonceTestAccountUpgrade.BalanceChange))
	require.Equal(t, state.GetNonce(nonceTestAccount), explicitNonce)
	require.Equal(t, state.GetBalance(nonceTestAccount), expectedNonceTestAccountBalance)
	require.Empty(t, state.GetCode(nonceTestAccount)) // No code should be set
}

func TestVMEupgradeActivatesCancun(t *testing.T) {
	tests := []struct {
		name        string
		genesisJSON string
		upgradeJSON string
		check       func(*testing.T, *VM) // function to check the VM state
	}{
		{
			name:        "Etna activates Cancun",
			genesisJSON: genesisJSONEtna,
			check: func(t *testing.T, vm *VM) {
				require.True(t, vm.chainConfig.IsCancun(common.Big0, DefaultEtnaTime))
			},
		},
		{
			name:        "Later Etna activates Cancun",
			genesisJSON: genesisJSONDurango,
			upgradeJSON: func() string {
				upgrade := &extras.UpgradeConfig{
					NetworkUpgradeOverrides: &extras.NetworkUpgrades{
						EtnaTimestamp: utils.NewUint64(DefaultEtnaTime + 2),
					},
				}
				b, err := json.Marshal(upgrade)
				require.NoError(t, err)
				return string(b)
			}(),
			check: func(t *testing.T, vm *VM) {
				require.False(t, vm.chainConfig.IsCancun(common.Big0, DefaultEtnaTime))
				require.True(t, vm.chainConfig.IsCancun(common.Big0, DefaultEtnaTime+2))
			},
		},
		{
			name:        "Changed Etna changes Cancun",
			genesisJSON: genesisJSONEtna,
			upgradeJSON: func() string {
				upgrade := &extras.UpgradeConfig{
					NetworkUpgradeOverrides: &extras.NetworkUpgrades{
						EtnaTimestamp: utils.NewUint64(DefaultEtnaTime + 2),
					},
				}
				b, err := json.Marshal(upgrade)
				require.NoError(t, err)
				return string(b)
			}(),
			check: func(t *testing.T, vm *VM) {
				require.False(t, vm.chainConfig.IsCancun(common.Big0, DefaultEtnaTime))
				require.True(t, vm.chainConfig.IsCancun(common.Big0, DefaultEtnaTime+2))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, vm, _, _ := GenesisVM(t, true, test.genesisJSON, "", test.upgradeJSON)
			defer func() { require.NoError(t, vm.Shutdown(context.Background())) }()
			test.check(t, vm)
		})
	}
}
