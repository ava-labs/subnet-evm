// Copyright 2023 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package ethapi

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/ava-labs/subnet-evm/accounts"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/consensus"
	"github.com/ava-labs/subnet-evm/consensus/dummy"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/bloombits"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/core/vm"
	"github.com/ava-labs/subnet-evm/internal/blocktest"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/rpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

func TestTransaction_RoundTripRpcJSON(t *testing.T) {
	var (
		config = params.TestChainConfig
		signer = types.LatestSigner(config)
		key, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		tests  = allTransactionTypes(common.Address{0xde, 0xad}, config)
	)
	t.Parallel()
	for i, tt := range tests {
		var tx2 types.Transaction
		tx, err := types.SignNewTx(key, signer, tt.Tx)
		if err != nil {
			t.Fatalf("test %d: signing failed: %v", i, err)
		}
		// Regular transaction
		if data, err := json.Marshal(tx); err != nil {
			t.Fatalf("test %d: marshalling failed; %v", i, err)
		} else if err = tx2.UnmarshalJSON(data); err != nil {
			t.Fatalf("test %d: sunmarshal failed: %v", i, err)
		} else if want, have := tx.Hash(), tx2.Hash(); want != have {
			t.Fatalf("test %d: stx changed, want %x have %x", i, want, have)
		}

		// rpcTransaction
		rpcTx := newRPCTransaction(tx, common.Hash{}, 0, 0, 0, nil, config)
		if data, err := json.Marshal(rpcTx); err != nil {
			t.Fatalf("test %d: marshalling failed; %v", i, err)
		} else if err = tx2.UnmarshalJSON(data); err != nil {
			t.Fatalf("test %d: unmarshal failed: %v", i, err)
		} else if want, have := tx.Hash(), tx2.Hash(); want != have {
			t.Fatalf("test %d: tx changed, want %x have %x", i, want, have)
		} else {
			want, have := tt.Want, string(data)
			require.JSONEqf(t, want, have, "test %d: rpc json not match, want %s have %s", i, want, have)
		}
	}
}

type txData struct {
	Tx   types.TxData
	Want string
}

func allTransactionTypes(addr common.Address, config *params.ChainConfig) []txData {
	return []txData{
		{
			Tx: &types.LegacyTx{
				Nonce:    5,
				GasPrice: big.NewInt(6),
				Gas:      7,
				To:       &addr,
				Value:    big.NewInt(8),
				Data:     []byte{0, 1, 2, 3, 4},
				V:        big.NewInt(9),
				R:        big.NewInt(10),
				S:        big.NewInt(11),
			},
			Want: `{
				"blockHash": null,
				"blockNumber": null,
				"from": "0x71562b71999873db5b286df957af199ec94617f7",
				"gas": "0x7",
				"gasPrice": "0x6",
				"hash": "0x3fa586d2448ae279279fa7036da74eb932763661543428c1a0aba21b95b37bdb",
				"input": "0x0001020304",
				"nonce": "0x5",
				"to": "0xdead000000000000000000000000000000000000",
				"transactionIndex": null,
				"value": "0x8",
				"type": "0x0",
				"chainId": "0x1",
				"v": "0x25",
				"r": "0xac639f4319e9268898e29444b97101f1225e2a0837151626da23e73dda2443fc",
				"s": "0x4fcc3f4c3a75f70ee45bb42d4b0aad432cc8c0140efb3e2611d6a6dda8460907"
			}`,
		}, {
			Tx: &types.LegacyTx{
				Nonce:    5,
				GasPrice: big.NewInt(6),
				Gas:      7,
				To:       nil,
				Value:    big.NewInt(8),
				Data:     []byte{0, 1, 2, 3, 4},
				V:        big.NewInt(32),
				R:        big.NewInt(10),
				S:        big.NewInt(11),
			},
			Want: `{
				"blockHash": null,
				"blockNumber": null,
				"from": "0x71562b71999873db5b286df957af199ec94617f7",
				"gas": "0x7",
				"gasPrice": "0x6",
				"hash": "0x617a316c6ff7ed2aa6ead1b4bb28a1322c2156c1c72f376a976d2d2adb1748ee",
				"input": "0x0001020304",
				"nonce": "0x5",
				"to": null,
				"transactionIndex": null,
				"value": "0x8",
				"type": "0x0",
				"chainId": "0x1",
				"v": "0x25",
				"r": "0xee8e89b513778d4815ae5969f3d55e0f7590f31b08f2a2290d5bc4ae37fce299",
				"s": "0x663db5c74c10e2b6525e7026e7cfd569b819ec91a042322655ff2b35060784b1"
			}`,
		},
		{
			Tx: &types.AccessListTx{
				ChainID:  config.ChainID,
				Nonce:    5,
				GasPrice: big.NewInt(6),
				Gas:      7,
				To:       &addr,
				Value:    big.NewInt(8),
				Data:     []byte{0, 1, 2, 3, 4},
				AccessList: types.AccessList{
					types.AccessTuple{
						Address:     common.Address{0x2},
						StorageKeys: []common.Hash{types.EmptyRootHash},
					},
				},
				V: big.NewInt(32),
				R: big.NewInt(10),
				S: big.NewInt(11),
			},
			Want: `{
				"blockHash": null,
				"blockNumber": null,
				"from": "0x71562b71999873db5b286df957af199ec94617f7",
				"gas": "0x7",
				"gasPrice": "0x6",
				"hash": "0x6becb7b9c171aa0d6d0a90dcd97bc3529c4d521f9cc9b7e31616aa9afc178c10",
				"input": "0x0001020304",
				"nonce": "0x5",
				"to": "0xdead000000000000000000000000000000000000",
				"transactionIndex": null,
				"value": "0x8",
				"type": "0x1",
				"accessList": [
					{
						"address": "0x0200000000000000000000000000000000000000",
						"storageKeys": [
							"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"
						]
					}
				],
				"chainId": "0x1",
				"v": "0x1",
				"r": "0xea2289ca0243beecbe69d337bbc53c618e1fb6bd2ec69fd4121df47149a4d4be",
				"s": "0x2dc134b6bc43abbdfebef0b2d62c175459fc1e8ddff60c8e740c461d7ea1522f",
				"yParity": "0x1"
			}`,
		}, {
			Tx: &types.AccessListTx{
				ChainID:  config.ChainID,
				Nonce:    5,
				GasPrice: big.NewInt(6),
				Gas:      7,
				To:       nil,
				Value:    big.NewInt(8),
				Data:     []byte{0, 1, 2, 3, 4},
				AccessList: types.AccessList{
					types.AccessTuple{
						Address:     common.Address{0x2},
						StorageKeys: []common.Hash{types.EmptyRootHash},
					},
				},
				V: big.NewInt(32),
				R: big.NewInt(10),
				S: big.NewInt(11),
			},
			Want: `{
				"blockHash": null,
				"blockNumber": null,
				"from": "0x71562b71999873db5b286df957af199ec94617f7",
				"gas": "0x7",
				"gasPrice": "0x6",
				"hash": "0x22fbf81bae4640511c706e2c72d2f2ef1abc1e7861f2b82c4cae5b102a40709c",
				"input": "0x0001020304",
				"nonce": "0x5",
				"to": null,
				"transactionIndex": null,
				"value": "0x8",
				"type": "0x1",
				"accessList": [
					{
						"address": "0x0200000000000000000000000000000000000000",
						"storageKeys": [
							"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"
						]
					}
				],
				"chainId": "0x1",
				"v": "0x1",
				"r": "0xc50e18edd861639735ec69ca12d82fcbb2c1921d2e2a8fd3a75f408d2d4b8118",
				"s": "0x32a908d1bc2db0229354f4dd392ffc37934e24ae8b18a620c6588c72660b6238",
				"yParity": "0x1"
			}`,
		}, {
			Tx: &types.DynamicFeeTx{
				ChainID:   config.ChainID,
				Nonce:     5,
				GasTipCap: big.NewInt(6),
				GasFeeCap: big.NewInt(9),
				Gas:       7,
				To:        &addr,
				Value:     big.NewInt(8),
				Data:      []byte{0, 1, 2, 3, 4},
				AccessList: types.AccessList{
					types.AccessTuple{
						Address:     common.Address{0x2},
						StorageKeys: []common.Hash{types.EmptyRootHash},
					},
				},
				V: big.NewInt(32),
				R: big.NewInt(10),
				S: big.NewInt(11),
			},
			Want: `{
				"blockHash": null,
				"blockNumber": null,
				"from": "0x71562b71999873db5b286df957af199ec94617f7",
				"gas": "0x7",
				"gasPrice": "0x9",
				"maxFeePerGas": "0x9",
				"maxPriorityFeePerGas": "0x6",
				"hash": "0xc5763d2ce6af3f694dcda8a9a50d4f75005a711edd382e993dd0406e0c54cfde",
				"input": "0x0001020304",
				"nonce": "0x5",
				"to": "0xdead000000000000000000000000000000000000",
				"transactionIndex": null,
				"value": "0x8",
				"type": "0x2",
				"accessList": [
					{
						"address": "0x0200000000000000000000000000000000000000",
						"storageKeys": [
							"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"
						]
					}
				],
				"chainId": "0x1",
				"v": "0x0",
				"r": "0x740eb1e3bc206760182993845b7815fd8cf7a42f1a1ed26027f736e9eccfd20f",
				"s": "0x31da567e2b3a83e58e42f7902c3706926c926625f6978c24fdaa21b9d143bbf7",
				"yParity": "0x0"
			}`,
		}, {
			Tx: &types.DynamicFeeTx{
				ChainID:    config.ChainID,
				Nonce:      5,
				GasTipCap:  big.NewInt(6),
				GasFeeCap:  big.NewInt(9),
				Gas:        7,
				To:         nil,
				Value:      big.NewInt(8),
				Data:       []byte{0, 1, 2, 3, 4},
				AccessList: types.AccessList{},
				V:          big.NewInt(32),
				R:          big.NewInt(10),
				S:          big.NewInt(11),
			},
			Want: `{
				"blockHash": null,
				"blockNumber": null,
				"from": "0x71562b71999873db5b286df957af199ec94617f7",
				"gas": "0x7",
				"gasPrice": "0x9",
				"maxFeePerGas": "0x9",
				"maxPriorityFeePerGas": "0x6",
				"hash": "0x85545f69b2410640fbbb7157b9a79adece45bac4b2803733d250d049e9501a28",
				"input": "0x0001020304",
				"nonce": "0x5",
				"to": null,
				"transactionIndex": null,
				"value": "0x8",
				"type": "0x2",
				"accessList": [],
				"chainId": "0x1",
				"v": "0x1",
				"r": "0x5004538adbe499313737033b22eb2b50a9450f02fab3971a591e6d57761b2cdf",
				"s": "0x5f7b1f5d11bd467d84f32beb2e89629351b96c5204c4f72d5d2040bee369a73a",
				"yParity": "0x1"
			  }`,
		},
	}
}

type testBackend struct {
	db    ethdb.Database
	chain *core.BlockChain
}

func newTestBackend(t *testing.T, n int, gspec *core.Genesis, generator func(i int, b *core.BlockGen)) *testBackend {
	var (
		engine      = dummy.NewCoinbaseFaker()
		cacheConfig = &core.CacheConfig{
			TrieCleanLimit: 256,
			TrieDirtyLimit: 256,
			SnapshotLimit:  0,
			Pruning:        false, // Archive mode
		}
	)
	// Generate blocks for testing
	db, blocks, _, _ := core.GenerateChainWithGenesis(gspec, engine, n, 10, generator)
	chain, err := core.NewBlockChain(db, cacheConfig, gspec, engine, vm.Config{}, gspec.ToBlock().Hash(), false)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	if n, err := chain.InsertChain(blocks); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}
	for _, block := range blocks {
		if err := chain.Accept(block); err != nil {
			t.Fatalf("block %d: failed to accept into chain: %v", block.NumberU64(), err)
		}
	}

	backend := &testBackend{db: db, chain: chain}
	return backend
}

func (b testBackend) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return big.NewInt(0), nil
}
func (b testBackend) FeeHistory(ctx context.Context, blockCount uint64, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*big.Int, [][]*big.Int, []*big.Int, []float64, error) {
	return nil, nil, nil, nil, nil
}
func (b testBackend) ChainDb() ethdb.Database                    { return b.db }
func (b testBackend) AccountManager() *accounts.Manager          { return nil }
func (b testBackend) ExtRPCEnabled() bool                        { return false }
func (b testBackend) RPCGasCap() uint64                          { return 10000000 }
func (b testBackend) RPCEVMTimeout() time.Duration               { return time.Second }
func (b testBackend) RPCTxFeeCap() float64                       { return 0 }
func (b testBackend) UnprotectedAllowed(*types.Transaction) bool { return false }
func (b testBackend) SetHead(number uint64)                      {}
func (b testBackend) HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error) {
	if number == rpc.LatestBlockNumber {
		return b.chain.CurrentBlock(), nil
	}
	return b.chain.GetHeaderByNumber(uint64(number)), nil
}
func (b testBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	return b.chain.GetHeaderByHash(hash), nil
}
func (b testBackend) HeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Header, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return b.HeaderByNumber(ctx, blockNr)
	}
	if blockHash, ok := blockNrOrHash.Hash(); ok {
		return b.HeaderByHash(ctx, blockHash)
	}
	panic("unknown type rpc.BlockNumberOrHash")
}
func (b testBackend) CurrentHeader() *types.Header { return b.chain.CurrentBlock() }
func (b testBackend) CurrentBlock() *types.Header  { return b.chain.CurrentBlock() }
func (b testBackend) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error) {
	if number == rpc.LatestBlockNumber {
		head := b.chain.CurrentBlock()
		return b.chain.GetBlock(head.Hash(), head.Number.Uint64()), nil
	}
	return b.chain.GetBlockByNumber(uint64(number)), nil
}
func (b testBackend) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return b.chain.GetBlockByHash(hash), nil
}
func (b testBackend) BlockByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Block, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return b.BlockByNumber(ctx, blockNr)
	}
	if blockHash, ok := blockNrOrHash.Hash(); ok {
		return b.BlockByHash(ctx, blockHash)
	}
	panic("unknown type rpc.BlockNumberOrHash")
}
func (b testBackend) GetBody(ctx context.Context, hash common.Hash, number rpc.BlockNumber) (*types.Body, error) {
	return b.chain.GetBlock(hash, uint64(number.Int64())).Body(), nil
}
func (b testBackend) StateAndHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	if number == rpc.PendingBlockNumber {
		panic("pending state not implemented")
	}
	header, err := b.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, nil, err
	}
	if header == nil {
		return nil, nil, errors.New("header not found")
	}
	stateDb, err := b.chain.StateAt(header.Root)
	return stateDb, header, err
}
func (b testBackend) StateAndHeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*state.StateDB, *types.Header, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return b.StateAndHeaderByNumber(ctx, blockNr)
	}
	panic("only implemented for number")
}
func (b testBackend) PendingBlockAndReceipts() (*types.Block, types.Receipts) { panic("implement me") }
func (b testBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	header, err := b.HeaderByHash(ctx, hash)
	if header == nil || err != nil {
		return nil, err
	}
	receipts := rawdb.ReadReceipts(b.db, hash, header.Number.Uint64(), header.Time, b.chain.Config())
	return receipts, nil
}
func (b testBackend) GetEVM(ctx context.Context, msg *core.Message, state *state.StateDB, header *types.Header, vmConfig *vm.Config, blockContext *vm.BlockContext) (*vm.EVM, func() error) {
	vmError := func() error { return nil }
	if vmConfig == nil {
		vmConfig = b.chain.GetVMConfig()
	}
	txContext := core.NewEVMTxContext(msg)
	context := core.NewEVMBlockContext(header, b.chain, nil)
	if blockContext != nil {
		context = *blockContext
	}
	return vm.NewEVM(context, txContext, state, b.chain.Config(), *vmConfig), vmError
}
func (b testBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	panic("implement me")
}
func (b testBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	panic("implement me")
}
func (b testBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	panic("implement me")
}
func (b testBackend) GetFeeConfigAt(parent *types.Header) (commontype.FeeConfig, *big.Int, error) {
	panic("implement me")
}
func (b testBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	panic("implement me")
}
func (b testBackend) GetTransaction(ctx context.Context, txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, error) {
	tx, blockHash, blockNumber, index := rawdb.ReadTransaction(b.db, txHash)
	return tx, blockHash, blockNumber, index, nil
}
func (b testBackend) GetPoolTransactions() (types.Transactions, error)         { panic("implement me") }
func (b testBackend) GetPoolTransaction(txHash common.Hash) *types.Transaction { panic("implement me") }
func (b testBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	panic("implement me")
}
func (b testBackend) Stats() (pending int, queued int) { panic("implement me") }
func (b testBackend) TxPoolContent() (map[common.Address][]*types.Transaction, map[common.Address][]*types.Transaction) {
	panic("implement me")
}
func (b testBackend) TxPoolContentFrom(addr common.Address) ([]*types.Transaction, []*types.Transaction) {
	panic("implement me")
}
func (b testBackend) SubscribeNewTxsEvent(events chan<- core.NewTxsEvent) event.Subscription {
	panic("implement me")
}
func (b testBackend) ChainConfig() *params.ChainConfig { return b.chain.Config() }
func (b testBackend) Engine() consensus.Engine         { return b.chain.Engine() }
func (b testBackend) GetLogs(ctx context.Context, blockHash common.Hash, number uint64) ([][]*types.Log, error) {
	panic("implement me")
}
func (b testBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	panic("implement me")
}
func (b testBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	panic("implement me")
}
func (b testBackend) SubscribePendingLogsEvent(ch chan<- []*types.Log) event.Subscription {
	panic("implement me")
}
func (b testBackend) BloomStatus() (uint64, uint64) { panic("implement me") }
func (b testBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	panic("implement me")
}
func (b testBackend) BadBlocks() ([]*types.Block, []*core.BadBlockReason) { return nil, nil }
func (b testBackend) EstimateBaseFee(ctx context.Context) (*big.Int, error) {
	panic("implement me")
}
func (b testBackend) LastAcceptedBlock() *types.Block { panic("implement me") }
func (b testBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	panic("implement me")
}

func TestEstimateGas(t *testing.T) {
	t.Parallel()
	// Initialize test accounts
	var (
		accounts = newAccounts(2)
		genesis  = &core.Genesis{
			Config: params.TestChainConfig,
			Alloc: core.GenesisAlloc{
				accounts[0].addr: {Balance: big.NewInt(params.Ether)},
				accounts[1].addr: {Balance: big.NewInt(params.Ether)},
			},
		}
		genBlocks      = 10
		signer         = types.HomesteadSigner{}
		randomAccounts = newAccounts(2)
	)
	api := NewBlockChainAPI(newTestBackend(t, genBlocks, genesis, func(i int, b *core.BlockGen) {
		// Transfer from account[0] to account[1]
		//    value: 1000 wei
		//    fee:   0 wei
		tx, _ := types.SignTx(types.NewTx(&types.LegacyTx{Nonce: uint64(i), To: &accounts[1].addr, Value: big.NewInt(1000), Gas: params.TxGas, GasPrice: b.BaseFee(), Data: nil}), signer, accounts[0].key)
		b.AddTx(tx)
	}))
	var testSuite = []struct {
		blockNumber rpc.BlockNumber
		call        TransactionArgs
		overrides   StateOverride
		expectErr   error
		want        uint64
	}{
		// simple transfer on latest block
		{
			blockNumber: rpc.LatestBlockNumber,
			call: TransactionArgs{
				From:  &accounts[0].addr,
				To:    &accounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			expectErr: nil,
			want:      21000,
		},
		// simple transfer with insufficient funds on latest block
		{
			blockNumber: rpc.LatestBlockNumber,
			call: TransactionArgs{
				From:  &randomAccounts[0].addr,
				To:    &accounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			expectErr: core.ErrInsufficientFunds,
			want:      21000,
		},
		// empty create
		{
			blockNumber: rpc.LatestBlockNumber,
			call:        TransactionArgs{},
			expectErr:   nil,
			want:        53000,
		},
		{
			blockNumber: rpc.LatestBlockNumber,
			call:        TransactionArgs{},
			overrides: StateOverride{
				randomAccounts[0].addr: OverrideAccount{Balance: newRPCBalance(new(big.Int).Mul(big.NewInt(1), big.NewInt(params.Ether)))},
			},
			expectErr: nil,
			want:      53000,
		},
		{
			blockNumber: rpc.LatestBlockNumber,
			call: TransactionArgs{
				From:  &randomAccounts[0].addr,
				To:    &randomAccounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			overrides: StateOverride{
				randomAccounts[0].addr: OverrideAccount{Balance: newRPCBalance(big.NewInt(0))},
			},
			expectErr: core.ErrInsufficientFunds,
		},
	}
	for i, tc := range testSuite {
		result, err := api.EstimateGas(context.Background(), tc.call, &rpc.BlockNumberOrHash{BlockNumber: &tc.blockNumber}, &tc.overrides)
		if tc.expectErr != nil {
			if err == nil {
				t.Errorf("test %d: want error %v, have nothing", i, tc.expectErr)
				continue
			}
			if !errors.Is(err, tc.expectErr) {
				t.Errorf("test %d: error mismatch, want %v, have %v", i, tc.expectErr, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("test %d: want no error, have %v", i, err)
			continue
		}
		if uint64(result) != tc.want {
			t.Errorf("test %d, result mismatch, have\n%v\n, want\n%v\n", i, uint64(result), tc.want)
		}
	}
}

func TestCall(t *testing.T) {
	t.Parallel()
	// Initialize test accounts
	var (
		accounts = newAccounts(3)
		genesis  = &core.Genesis{
			Config: params.TestChainConfig,
			Alloc: core.GenesisAlloc{
				accounts[0].addr: {Balance: big.NewInt(params.Ether)},
				accounts[1].addr: {Balance: big.NewInt(params.Ether)},
				accounts[2].addr: {Balance: big.NewInt(params.Ether)},
			},
		}
		genBlocks = 10
		signer    = types.HomesteadSigner{}
	)
	api := NewBlockChainAPI(newTestBackend(t, genBlocks, genesis, func(i int, b *core.BlockGen) {
		// Transfer from account[0] to account[1]
		//    value: 1000 wei
		//    fee:   0 wei
		tx, _ := types.SignTx(types.NewTx(&types.LegacyTx{Nonce: uint64(i), To: &accounts[1].addr, Value: big.NewInt(1000), Gas: params.TxGas, GasPrice: b.BaseFee(), Data: nil}), signer, accounts[0].key)
		b.AddTx(tx)
	}))
	randomAccounts := newAccounts(3)
	var testSuite = []struct {
		blockNumber    rpc.BlockNumber
		overrides      StateOverride
		call           TransactionArgs
		blockOverrides BlockOverrides
		expectErr      error
		want           string
	}{
		// transfer on genesis
		{
			blockNumber: rpc.BlockNumber(0),
			call: TransactionArgs{
				From:  &accounts[0].addr,
				To:    &accounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			expectErr: nil,
			want:      "0x",
		},
		// transfer on the head
		{
			blockNumber: rpc.BlockNumber(genBlocks),
			call: TransactionArgs{
				From:  &accounts[0].addr,
				To:    &accounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			expectErr: nil,
			want:      "0x",
		},
		// transfer on a non-existent block, error expects
		{
			blockNumber: rpc.BlockNumber(genBlocks + 1),
			call: TransactionArgs{
				From:  &accounts[0].addr,
				To:    &accounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			expectErr: errors.New("header not found"),
		},
		// transfer on the latest block
		{
			blockNumber: rpc.LatestBlockNumber,
			call: TransactionArgs{
				From:  &accounts[0].addr,
				To:    &accounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			expectErr: nil,
			want:      "0x",
		},
		// Call which can only succeed if state is state overridden
		{
			blockNumber: rpc.LatestBlockNumber,
			call: TransactionArgs{
				From:  &randomAccounts[0].addr,
				To:    &randomAccounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			overrides: StateOverride{
				randomAccounts[0].addr: OverrideAccount{Balance: newRPCBalance(new(big.Int).Mul(big.NewInt(1), big.NewInt(params.Ether)))},
			},
			want: "0x",
		},
		// Invalid call without state overriding
		{
			blockNumber: rpc.LatestBlockNumber,
			call: TransactionArgs{
				From:  &randomAccounts[0].addr,
				To:    &randomAccounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			expectErr: core.ErrInsufficientFunds,
		},
		// Successful simple contract call
		//
		// // SPDX-License-Identifier: GPL-3.0
		//
		//  pragma solidity >=0.7.0 <0.8.0;
		//
		//  /**
		//   * @title Storage
		//   * @dev Store & retrieve value in a variable
		//   */
		//  contract Storage {
		//      uint256 public number;
		//      constructor() {
		//          number = block.number;
		//      }
		//  }
		{
			blockNumber: rpc.LatestBlockNumber,
			call: TransactionArgs{
				From: &randomAccounts[0].addr,
				To:   &randomAccounts[2].addr,
				Data: hex2Bytes("8381f58a"), // call number()
			},
			overrides: StateOverride{
				randomAccounts[2].addr: OverrideAccount{
					Code:      hex2Bytes("6080604052348015600f57600080fd5b506004361060285760003560e01c80638381f58a14602d575b600080fd5b60336049565b6040518082815260200191505060405180910390f35b6000548156fea2646970667358221220eab35ffa6ab2adfe380772a48b8ba78e82a1b820a18fcb6f59aa4efb20a5f60064736f6c63430007040033"),
					StateDiff: &map[common.Hash]common.Hash{{}: common.BigToHash(big.NewInt(123))},
				},
			},
			want: "0x000000000000000000000000000000000000000000000000000000000000007b",
		},
		// Block overrides should work
		{
			blockNumber: rpc.LatestBlockNumber,
			call: TransactionArgs{
				From: &accounts[1].addr,
				Input: &hexutil.Bytes{
					0x43,             // NUMBER
					0x60, 0x00, 0x52, // MSTORE offset 0
					0x60, 0x20, 0x60, 0x00, 0xf3,
				},
			},
			blockOverrides: BlockOverrides{Number: (*hexutil.Big)(big.NewInt(11))},
			want:           "0x000000000000000000000000000000000000000000000000000000000000000b",
		},
	}
	for i, tc := range testSuite {
		result, err := api.Call(context.Background(), tc.call, rpc.BlockNumberOrHash{BlockNumber: &tc.blockNumber}, &tc.overrides, &tc.blockOverrides)
		if tc.expectErr != nil {
			if err == nil {
				t.Errorf("test %d: want error %v, have nothing", i, tc.expectErr)
				continue
			}
			if !errors.Is(err, tc.expectErr) {
				// Second try
				if !reflect.DeepEqual(err, tc.expectErr) {
					t.Errorf("test %d: error mismatch, want %v, have %v", i, tc.expectErr, err)
				}
			}
			continue
		}
		if err != nil {
			t.Errorf("test %d: want no error, have %v", i, err)
			continue
		}
		if !reflect.DeepEqual(result.String(), tc.want) {
			t.Errorf("test %d, result mismatch, have\n%v\n, want\n%v\n", i, result.String(), tc.want)
		}
	}
}

type Account struct {
	key  *ecdsa.PrivateKey
	addr common.Address
}

func newAccounts(n int) (accounts []Account) {
	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)
		accounts = append(accounts, Account{key: key, addr: addr})
	}
	slices.SortFunc(accounts, func(a, b Account) int { return a.addr.Cmp(b.addr) })
	return accounts
}

func newRPCBalance(balance *big.Int) **hexutil.Big {
	rpcBalance := (*hexutil.Big)(balance)
	return &rpcBalance
}

func hex2Bytes(str string) *hexutil.Bytes {
	rpcBytes := hexutil.Bytes(common.Hex2Bytes(str))
	return &rpcBytes
}

func TestRPCMarshalBlock(t *testing.T) {
	t.Parallel()
	var (
		txs []*types.Transaction
		to  = common.BytesToAddress([]byte{0x11})
	)
	for i := uint64(1); i <= 4; i++ {
		var tx *types.Transaction
		if i%2 == 0 {
			tx = types.NewTx(&types.LegacyTx{
				Nonce:    i,
				GasPrice: big.NewInt(11111),
				Gas:      1111,
				To:       &to,
				Value:    big.NewInt(111),
				Data:     []byte{0x11, 0x11, 0x11},
			})
		} else {
			tx = types.NewTx(&types.AccessListTx{
				ChainID:  big.NewInt(1337),
				Nonce:    i,
				GasPrice: big.NewInt(11111),
				Gas:      1111,
				To:       &to,
				Value:    big.NewInt(111),
				Data:     []byte{0x11, 0x11, 0x11},
			})
		}
		txs = append(txs, tx)
	}
	block := types.NewBlock(&types.Header{Number: big.NewInt(100)}, txs, nil, nil, blocktest.NewHasher())

	var testSuite = []struct {
		inclTx bool
		fullTx bool
		want   string
	}{
		// without txs
		{
			inclTx: false,
			fullTx: false,
			want: `{
				"difficulty": "0x0",
				"extraData": "0x",
				"gasLimit": "0x0",
				"gasUsed": "0x0",
				"hash": "0x9b73c83b25d0faf7eab854e3684c7e394336d6e135625aafa5c183f27baa8fee",
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"miner": "0x0000000000000000000000000000000000000000",
				"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce": "0x0000000000000000",
				"number": "0x64",
				"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
				"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"size": "0x296",
				"stateRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"timestamp": "0x0",
				"transactionsRoot": "0x661a9febcfa8f1890af549b874faf9fa274aede26ef489d9db0b25daa569450e",
				"uncles": []
			}`,
		},
		// only tx hashes
		{
			inclTx: true,
			fullTx: false,
			want: `{
				"difficulty": "0x0",
				"extraData": "0x",
				"gasLimit": "0x0",
				"gasUsed": "0x0",
				"hash": "0x9b73c83b25d0faf7eab854e3684c7e394336d6e135625aafa5c183f27baa8fee",
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"miner": "0x0000000000000000000000000000000000000000",
				"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce": "0x0000000000000000",
				"number": "0x64",
				"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
				"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"size": "0x296",
				"stateRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"timestamp": "0x0",
				"transactions": [
					"0x7d39df979e34172322c64983a9ad48302c2b889e55bda35324afecf043a77605",
					"0x9bba4c34e57c875ff57ac8d172805a26ae912006985395dc1bdf8f44140a7bf4",
					"0x98909ea1ff040da6be56bc4231d484de1414b3c1dac372d69293a4beb9032cb5",
					"0x12e1f81207b40c3bdcc13c0ee18f5f86af6d31754d57a0ea1b0d4cfef21abef1"
				],
				"transactionsRoot": "0x661a9febcfa8f1890af549b874faf9fa274aede26ef489d9db0b25daa569450e",
				"uncles": []
			}`,
		},
		// full tx details
		{
			inclTx: true,
			fullTx: true,
			want: `{
				"difficulty": "0x0",
				"extraData": "0x",
				"gasLimit": "0x0",
				"gasUsed": "0x0",
				"hash": "0x9b73c83b25d0faf7eab854e3684c7e394336d6e135625aafa5c183f27baa8fee",
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"miner": "0x0000000000000000000000000000000000000000",
				"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce": "0x0000000000000000",
				"number": "0x64",
				"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
				"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"size": "0x296",
				"stateRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"timestamp": "0x0",
				"transactions": [
					{
						"blockHash": "0x9b73c83b25d0faf7eab854e3684c7e394336d6e135625aafa5c183f27baa8fee",
						"blockNumber": "0x64",
						"from": "0x0000000000000000000000000000000000000000",
						"gas": "0x457",
						"gasPrice": "0x2b67",
						"hash": "0x7d39df979e34172322c64983a9ad48302c2b889e55bda35324afecf043a77605",
						"input": "0x111111",
						"nonce": "0x1",
						"to": "0x0000000000000000000000000000000000000011",
						"transactionIndex": "0x0",
						"value": "0x6f",
						"type": "0x1",
						"accessList": [],
						"chainId": "0x539",
						"v": "0x0",
						"r": "0x0",
						"s": "0x0",
						"yParity": "0x0"
					},
					{
						"blockHash": "0x9b73c83b25d0faf7eab854e3684c7e394336d6e135625aafa5c183f27baa8fee",
						"blockNumber": "0x64",
						"from": "0x0000000000000000000000000000000000000000",
						"gas": "0x457",
						"gasPrice": "0x2b67",
						"hash": "0x9bba4c34e57c875ff57ac8d172805a26ae912006985395dc1bdf8f44140a7bf4",
						"input": "0x111111",
						"nonce": "0x2",
						"to": "0x0000000000000000000000000000000000000011",
						"transactionIndex": "0x1",
						"value": "0x6f",
						"type": "0x0",
						"chainId": "0x7fffffffffffffee",
						"v": "0x0",
						"r": "0x0",
						"s": "0x0"
					},
					{
						"blockHash": "0x9b73c83b25d0faf7eab854e3684c7e394336d6e135625aafa5c183f27baa8fee",
						"blockNumber": "0x64",
						"from": "0x0000000000000000000000000000000000000000",
						"gas": "0x457",
						"gasPrice": "0x2b67",
						"hash": "0x98909ea1ff040da6be56bc4231d484de1414b3c1dac372d69293a4beb9032cb5",
						"input": "0x111111",
						"nonce": "0x3",
						"to": "0x0000000000000000000000000000000000000011",
						"transactionIndex": "0x2",
						"value": "0x6f",
						"type": "0x1",
						"accessList": [],
						"chainId": "0x539",
						"v": "0x0",
						"r": "0x0",
						"s": "0x0",
						"yParity": "0x0"
					},
					{
						"blockHash": "0x9b73c83b25d0faf7eab854e3684c7e394336d6e135625aafa5c183f27baa8fee",
						"blockNumber": "0x64",
						"from": "0x0000000000000000000000000000000000000000",
						"gas": "0x457",
						"gasPrice": "0x2b67",
						"hash": "0x12e1f81207b40c3bdcc13c0ee18f5f86af6d31754d57a0ea1b0d4cfef21abef1",
						"input": "0x111111",
						"nonce": "0x4",
						"to": "0x0000000000000000000000000000000000000011",
						"transactionIndex": "0x3",
						"value": "0x6f",
						"type": "0x0",
						"chainId": "0x7fffffffffffffee",
						"v": "0x0",
						"r": "0x0",
						"s": "0x0"
					}
				],
				"transactionsRoot": "0x661a9febcfa8f1890af549b874faf9fa274aede26ef489d9db0b25daa569450e",
				"uncles": []
			}`,
		},
	}

	for i, tc := range testSuite {
		resp := RPCMarshalBlock(block, tc.inclTx, tc.fullTx, params.TestSubnetEVMConfig)
		out, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("test %d: json marshal error: %v", i, err)
			continue
		}
		assert.JSONEqf(t, tc.want, string(out), "test %d", i)
	}
}

func TestRPCGetBlockOrHeader(t *testing.T) {
	t.Parallel()

	// Initialize test accounts
	var (
		acc1Key, _ = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		acc2Key, _ = crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
		acc1Addr   = crypto.PubkeyToAddress(acc1Key.PublicKey)
		acc2Addr   = crypto.PubkeyToAddress(acc2Key.PublicKey)
		genesis    = &core.Genesis{
			Config: params.TestSubnetEVMConfig,
			Alloc: core.GenesisAlloc{
				acc1Addr: {Balance: big.NewInt(params.Ether)},
				acc2Addr: {Balance: big.NewInt(params.Ether)},
			},
		}
		genBlocks = 10
		signer    = types.HomesteadSigner{}
		tx        = types.NewTx(&types.LegacyTx{
			Nonce:    11,
			GasPrice: big.NewInt(11111),
			Gas:      1111,
			To:       &acc2Addr,
			Value:    big.NewInt(111),
			Data:     []byte{0x11, 0x11, 0x11},
		})
		pending = types.NewBlock(&types.Header{Number: big.NewInt(11), Time: 42}, []*types.Transaction{tx}, nil, nil, blocktest.NewHasher())
	)
	backend := newTestBackend(t, genBlocks, genesis, func(i int, b *core.BlockGen) {
		// Transfer from account[0] to account[1]
		//    value: 1000 wei
		//    fee:   0 wei
		tx, _ := types.SignTx(types.NewTx(&types.LegacyTx{Nonce: uint64(i), To: &acc2Addr, Value: big.NewInt(1000), Gas: params.TxGas, GasPrice: b.BaseFee(), Data: nil}), signer, acc1Key)
		b.AddTx(tx)
	})
	api := NewBlockChainAPI(backend)
	blockHashes := make([]common.Hash, genBlocks+1)
	ctx := context.Background()
	for i := 0; i <= genBlocks; i++ {
		header, err := backend.HeaderByNumber(ctx, rpc.BlockNumber(i))
		if err != nil {
			t.Errorf("failed to get block: %d err: %v", i, err)
		}
		blockHashes[i] = header.Hash()
	}
	pendingHash := pending.Hash()

	var testSuite = []struct {
		blockNumber rpc.BlockNumber
		blockHash   *common.Hash
		fullTx      bool
		reqHeader   bool
		want        string
		expectErr   error
	}{
		// 0. latest header
		{
			blockNumber: rpc.LatestBlockNumber,
			reqHeader:   true,
			want: `{
				"baseFeePerGas": "0x5d21dba00",
				"blockGasCost": "0x0",
				"difficulty": "0x1",
				"extraData": "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"gasLimit": "0x7a1200",
				"gasUsed": "0x5208",
				"hash": "0x1ec39e7ec46f8df1fb31cfca53fbf71a01869af8bd8f9a1bccbffc16ffa1461d",
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"miner": "0x0000000000000000000000000000000000000000",
				"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce": "0x0000000000000000",
				"number": "0xa",
				"parentHash": "0x0583a9d630632001771b4ecc7d62574aec3825aff47e2a680b0ea4ddb79e7365",
				"receiptsRoot": "0x056b23fbba480696b65fe5a59b8f2148a1299103c4f57df839233af2cf4ca2d2",
				"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"stateRoot": "0x7e06187d15d50badf60930290fb292ebe43e79553ad8b7d8f1b614316631def7",
				"timestamp": "0x64",
				"totalDifficulty": "0xa",
				"transactionsRoot": "0x69ff8003291e1cd08f75d174f070618f7291e4540b2e33f60b3375743e3fda01"
			  }`,
		},
		// 1. genesis header
		{
			blockNumber: rpc.BlockNumber(0),
			reqHeader:   true,
			want: `{
				"baseFeePerGas": "0x5d21dba00",
				"difficulty": "0x20000",
				"extraData": "0x",
				"gasLimit": "0x47e7c4",
				"gasUsed": "0x0",
				"hash": "0x3ead7923676a44500c46ad2192a0fc084aa42063b1703e6866f138a47fb1a9ca",
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"miner": "0x0000000000000000000000000000000000000000",
				"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce": "0x0000000000000000",
				"number": "0x0",
				"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
				"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"stateRoot": "0xfe168c5e9584a85927212e5bea5304bb7d0d8a893453b4b2c52176a72f585ae2",
				"timestamp": "0x0",
				"totalDifficulty": "0x0",
				"transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"
			}`,
		},
		// 2. #1 header
		{
			blockNumber: rpc.BlockNumber(1),
			reqHeader:   true,
			want: `{
				"baseFeePerGas": "0x5d21dba00",
				"blockGasCost": "0x0",
				"difficulty": "0x1",
				"extraData": "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"gasLimit": "0x7a1200",
				"gasUsed": "0x5208",
				"hash": "0x0f67ad1fc8052afad4c24551748600c164091cf37e068adef76315025d3c78e7",
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"miner": "0x0000000000000000000000000000000000000000",
				"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce": "0x0000000000000000",
				"number": "0x1",
				"parentHash": "0x3ead7923676a44500c46ad2192a0fc084aa42063b1703e6866f138a47fb1a9ca",
				"receiptsRoot": "0x056b23fbba480696b65fe5a59b8f2148a1299103c4f57df839233af2cf4ca2d2",
				"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"stateRoot": "0x6b830601767ac4968163193facbe20123435180e325910b2c50efa21f778c697",
				"timestamp": "0xa",
				"totalDifficulty": "0x1",
				"transactionsRoot": "0x87c65a3f1a98dafe282ace11eaf88b8f31bf41fe6794d401d2f986c1af84bcd5"
			}`,
		},
		// 3. latest-1 header
		{
			blockNumber: rpc.BlockNumber(9),
			reqHeader:   true,
			want: `{
				"baseFeePerGas": "0x5d21dba00",
				"blockGasCost": "0x0",
				"difficulty": "0x1",
				"extraData": "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"gasLimit": "0x7a1200",
				"gasUsed": "0x5208",
				"hash": "0x0583a9d630632001771b4ecc7d62574aec3825aff47e2a680b0ea4ddb79e7365",
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"miner": "0x0000000000000000000000000000000000000000",
				"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce": "0x0000000000000000",
				"number": "0x9",
				"parentHash": "0x2fab5c6892c66668842683ced6b384c2ee83bfd6096a58f451290cabaf57a63e",
				"receiptsRoot": "0x056b23fbba480696b65fe5a59b8f2148a1299103c4f57df839233af2cf4ca2d2",
				"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"stateRoot": "0x3703d70c6443e809ce035c2a8212dbf9813f6b7d1b0f597766e9023867a852f5",
				"timestamp": "0x5a",
				"totalDifficulty": "0x9",
				"transactionsRoot": "0xe16929d9c7efab0f962c1ed8c1295ddff42d3026779ed1318ea079ca580ee4cb"
			}`,
		},
		// 4. latest+1 header
		{
			blockNumber: rpc.BlockNumber(11),
			reqHeader:   true,
			want:        "null",
		},
		// 5. pending header
		{
			blockNumber: rpc.PendingBlockNumber,
			reqHeader:   true,
			want:        "null",
		},
		// 6. latest block
		{
			blockNumber: rpc.LatestBlockNumber,
			want: `{
				"baseFeePerGas": "0x5d21dba00",
				"blockGasCost": "0x0",
				"difficulty": "0x1",
				"extraData": "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"gasLimit": "0x7a1200",
				"gasUsed": "0x5208",
				"hash": "0x1ec39e7ec46f8df1fb31cfca53fbf71a01869af8bd8f9a1bccbffc16ffa1461d",
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"miner": "0x0000000000000000000000000000000000000000",
				"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce": "0x0000000000000000",
				"number": "0xa",
				"parentHash": "0x0583a9d630632001771b4ecc7d62574aec3825aff47e2a680b0ea4ddb79e7365",
				"receiptsRoot": "0x056b23fbba480696b65fe5a59b8f2148a1299103c4f57df839233af2cf4ca2d2",
				"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"size": "0x2bb",
				"stateRoot": "0x7e06187d15d50badf60930290fb292ebe43e79553ad8b7d8f1b614316631def7",
				"timestamp": "0x64",
				"totalDifficulty": "0xa",
				"transactions": [
				  "0x71be223424ab6e3457513a760b196d43b094414c32a70ff929b2b720a16b832d"
				],
				"transactionsRoot": "0x69ff8003291e1cd08f75d174f070618f7291e4540b2e33f60b3375743e3fda01",
				"uncles": []
			}`,
		},
		// 7. genesis block
		{
			blockNumber: rpc.BlockNumber(0),
			want: `{
				"baseFeePerGas": "0x5d21dba00",
				"difficulty": "0x20000",
				"extraData": "0x",
				"gasLimit": "0x47e7c4",
				"gasUsed": "0x0",
				"hash": "0x3ead7923676a44500c46ad2192a0fc084aa42063b1703e6866f138a47fb1a9ca",
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"miner": "0x0000000000000000000000000000000000000000",
				"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce": "0x0000000000000000",
				"number": "0x0",
				"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
				"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"size": "0x201",
				"stateRoot": "0xfe168c5e9584a85927212e5bea5304bb7d0d8a893453b4b2c52176a72f585ae2",
				"timestamp": "0x0",
				"totalDifficulty": "0x0",
				"transactions": [],
				"transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
				"uncles": []
			}`,
		},
		// 8. #1 block
		{
			blockNumber: rpc.BlockNumber(1),
			want: `{
				"baseFeePerGas": "0x5d21dba00",
				"blockGasCost": "0x0",
				"difficulty": "0x1",
				"extraData": "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"gasLimit": "0x7a1200",
				"gasUsed": "0x5208",
				"hash": "0x0f67ad1fc8052afad4c24551748600c164091cf37e068adef76315025d3c78e7",
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"miner": "0x0000000000000000000000000000000000000000",
				"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce": "0x0000000000000000",
				"number": "0x1",
				"parentHash": "0x3ead7923676a44500c46ad2192a0fc084aa42063b1703e6866f138a47fb1a9ca",
				"receiptsRoot": "0x056b23fbba480696b65fe5a59b8f2148a1299103c4f57df839233af2cf4ca2d2",
				"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"size": "0x2bb",
				"stateRoot": "0x6b830601767ac4968163193facbe20123435180e325910b2c50efa21f778c697",
				"timestamp": "0xa",
				"totalDifficulty": "0x1",
				"transactions": [
				  "0xdf92bc7c4c0341ecbdcd2a3ca7011fe9e21df4b8553bf0c8caabe6cb4a1aee26"
				],
				"transactionsRoot": "0x87c65a3f1a98dafe282ace11eaf88b8f31bf41fe6794d401d2f986c1af84bcd5",
				"uncles": []
			}`,
		},
		// 9. latest-1 block
		{
			blockNumber: rpc.BlockNumber(9),
			fullTx:      true,
			want: `{
				"baseFeePerGas": "0x5d21dba00",
				"blockGasCost": "0x0",
				"difficulty": "0x1",
				"extraData": "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"gasLimit": "0x7a1200",
				"gasUsed": "0x5208",
				"hash": "0x0583a9d630632001771b4ecc7d62574aec3825aff47e2a680b0ea4ddb79e7365",
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"miner": "0x0000000000000000000000000000000000000000",
				"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce": "0x0000000000000000",
				"number": "0x9",
				"parentHash": "0x2fab5c6892c66668842683ced6b384c2ee83bfd6096a58f451290cabaf57a63e",
				"receiptsRoot": "0x056b23fbba480696b65fe5a59b8f2148a1299103c4f57df839233af2cf4ca2d2",
				"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"size": "0x2bb",
				"stateRoot": "0x3703d70c6443e809ce035c2a8212dbf9813f6b7d1b0f597766e9023867a852f5",
				"timestamp": "0x5a",
				"totalDifficulty": "0x9",
				"transactions": [
					{
						"blockHash": "0x0583a9d630632001771b4ecc7d62574aec3825aff47e2a680b0ea4ddb79e7365",
						"blockNumber": "0x9",
						"from": "0x703c4b2bd70c169f5717101caee543299fc946c7",
						"gas": "0x5208",
						"gasPrice": "0x5d21dba00",
						"hash": "0x237f95840187a93f8aaf8d6f1515f8a8ac9d9359fcb0c220cdb3d642d6b9a19a",
						"input": "0x",
						"nonce": "0x8",
						"to": "0x0d3ab14bbad3d99f4203bd7a11acb94882050e7e",
						"transactionIndex": "0x0",
						"value": "0x3e8",
						"type": "0x0",
						"v": "0x1c",
						"r": "0xd7cdc527490b7ba29c515aae3bbe80c67729cda7f736e6515652cfc40e9da68f",
						"s": "0x4d0a4a59bef165b16f910bdadd41efaaad1b73549bacc35eaf6d073eb1fb92b7"
					}
				],
				"transactionsRoot": "0xe16929d9c7efab0f962c1ed8c1295ddff42d3026779ed1318ea079ca580ee4cb",
				"uncles": []
			}`,
		},
		// 10. latest+1 block
		{
			blockNumber: rpc.BlockNumber(11),
			fullTx:      true,
			want:        "null",
		},
		// 11. pending block
		{
			blockNumber: rpc.PendingBlockNumber,
			want:        "null",
		},
		// 12. pending block + fullTx
		{
			blockNumber: rpc.PendingBlockNumber,
			fullTx:      true,
			want:        "null",
		},
		// 13. latest header by hash
		{
			blockHash: &blockHashes[len(blockHashes)-1],
			reqHeader: true,
			want: `{
				"baseFeePerGas": "0x5d21dba00",
				"blockGasCost": "0x0",
				"difficulty": "0x1",
				"extraData": "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"gasLimit": "0x7a1200",
				"gasUsed": "0x5208",
				"hash": "0x1ec39e7ec46f8df1fb31cfca53fbf71a01869af8bd8f9a1bccbffc16ffa1461d",
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"miner": "0x0000000000000000000000000000000000000000",
				"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce": "0x0000000000000000",
				"number": "0xa",
				"parentHash": "0x0583a9d630632001771b4ecc7d62574aec3825aff47e2a680b0ea4ddb79e7365",
				"receiptsRoot": "0x056b23fbba480696b65fe5a59b8f2148a1299103c4f57df839233af2cf4ca2d2",
				"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"stateRoot": "0x7e06187d15d50badf60930290fb292ebe43e79553ad8b7d8f1b614316631def7",
				"timestamp": "0x64",
				"totalDifficulty": "0xa",
				"transactionsRoot": "0x69ff8003291e1cd08f75d174f070618f7291e4540b2e33f60b3375743e3fda01"
			}`,
		},
		// 14. genesis header by hash
		{
			blockHash: &blockHashes[0],
			reqHeader: true,
			want: `{
				"baseFeePerGas": "0x5d21dba00",
				"difficulty": "0x20000",
				"extraData": "0x",
				"gasLimit": "0x47e7c4",
				"gasUsed": "0x0",
				"hash": "0x3ead7923676a44500c46ad2192a0fc084aa42063b1703e6866f138a47fb1a9ca",
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"miner": "0x0000000000000000000000000000000000000000",
				"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce": "0x0000000000000000",
				"number": "0x0",
				"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
				"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"stateRoot": "0xfe168c5e9584a85927212e5bea5304bb7d0d8a893453b4b2c52176a72f585ae2",
				"timestamp": "0x0",
				"totalDifficulty": "0x0",
				"transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"
			}`,
		},
		// 15. #1 header
		{
			blockHash: &blockHashes[1],
			reqHeader: true,
			want: `{
				"baseFeePerGas": "0x5d21dba00",
				"blockGasCost": "0x0",
				"difficulty": "0x1",
				"extraData": "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"gasLimit": "0x7a1200",
				"gasUsed": "0x5208",
				"hash": "0x0f67ad1fc8052afad4c24551748600c164091cf37e068adef76315025d3c78e7",
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"miner": "0x0000000000000000000000000000000000000000",
				"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce": "0x0000000000000000",
				"number": "0x1",
				"parentHash": "0x3ead7923676a44500c46ad2192a0fc084aa42063b1703e6866f138a47fb1a9ca",
				"receiptsRoot": "0x056b23fbba480696b65fe5a59b8f2148a1299103c4f57df839233af2cf4ca2d2",
				"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"stateRoot": "0x6b830601767ac4968163193facbe20123435180e325910b2c50efa21f778c697",
				"timestamp": "0xa",
				"totalDifficulty": "0x1",
				"transactionsRoot": "0x87c65a3f1a98dafe282ace11eaf88b8f31bf41fe6794d401d2f986c1af84bcd5"
			}`,
		},
		// 16. latest-1 header
		{
			blockHash: &blockHashes[len(blockHashes)-2],
			reqHeader: true,
			want: `{
				"baseFeePerGas": "0x5d21dba00",
				"blockGasCost": "0x0",
				"difficulty": "0x1",
				"extraData": "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"gasLimit": "0x7a1200",
				"gasUsed": "0x5208",
				"hash": "0x0583a9d630632001771b4ecc7d62574aec3825aff47e2a680b0ea4ddb79e7365",
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"miner": "0x0000000000000000000000000000000000000000",
				"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce": "0x0000000000000000",
				"number": "0x9",
				"parentHash": "0x2fab5c6892c66668842683ced6b384c2ee83bfd6096a58f451290cabaf57a63e",
				"receiptsRoot": "0x056b23fbba480696b65fe5a59b8f2148a1299103c4f57df839233af2cf4ca2d2",
				"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"stateRoot": "0x3703d70c6443e809ce035c2a8212dbf9813f6b7d1b0f597766e9023867a852f5",
				"timestamp": "0x5a",
				"totalDifficulty": "0x9",
				"transactionsRoot": "0xe16929d9c7efab0f962c1ed8c1295ddff42d3026779ed1318ea079ca580ee4cb"
			}`,
		},
		// 17. empty hash
		{
			blockHash: &common.Hash{},
			reqHeader: true,
			want:      "null",
		},
		// 18. pending hash
		{
			blockHash: &pendingHash,
			reqHeader: true,
			want:      `null`,
		},
		// 19. latest block
		{
			blockHash: &blockHashes[len(blockHashes)-1],
			want: `{
				"baseFeePerGas": "0x5d21dba00",
				"blockGasCost": "0x0",
				"difficulty": "0x1",
				"extraData": "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"gasLimit": "0x7a1200",
				"gasUsed": "0x5208",
				"hash": "0x1ec39e7ec46f8df1fb31cfca53fbf71a01869af8bd8f9a1bccbffc16ffa1461d",
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"miner": "0x0000000000000000000000000000000000000000",
				"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce": "0x0000000000000000",
				"number": "0xa",
				"parentHash": "0x0583a9d630632001771b4ecc7d62574aec3825aff47e2a680b0ea4ddb79e7365",
				"receiptsRoot": "0x056b23fbba480696b65fe5a59b8f2148a1299103c4f57df839233af2cf4ca2d2",
				"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"size": "0x2bb",
				"stateRoot": "0x7e06187d15d50badf60930290fb292ebe43e79553ad8b7d8f1b614316631def7",
				"timestamp": "0x64",
				"totalDifficulty": "0xa",
				"transactions": [
					"0x71be223424ab6e3457513a760b196d43b094414c32a70ff929b2b720a16b832d"
				],
				"transactionsRoot": "0x69ff8003291e1cd08f75d174f070618f7291e4540b2e33f60b3375743e3fda01",
				"uncles": []
			}`,
		},
		// 20. genesis block
		{
			blockHash: &blockHashes[0],
			want: `{
				"baseFeePerGas": "0x5d21dba00",
				"difficulty": "0x20000",
				"extraData": "0x",
				"gasLimit": "0x47e7c4",
				"gasUsed": "0x0",
				"hash": "0x3ead7923676a44500c46ad2192a0fc084aa42063b1703e6866f138a47fb1a9ca",
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"miner": "0x0000000000000000000000000000000000000000",
				"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce": "0x0000000000000000",
				"number": "0x0",
				"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
				"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"size": "0x201",
				"stateRoot": "0xfe168c5e9584a85927212e5bea5304bb7d0d8a893453b4b2c52176a72f585ae2",
				"timestamp": "0x0",
				"totalDifficulty": "0x0",
				"transactions": [],
				"transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
				"uncles": []
			}`,
		},
		// 21. #1 block
		{
			blockHash: &blockHashes[1],
			want: `{
				"baseFeePerGas": "0x5d21dba00",
				"blockGasCost": "0x0",
				"difficulty": "0x1",
				"extraData": "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"gasLimit": "0x7a1200",
				"gasUsed": "0x5208",
				"hash": "0x0f67ad1fc8052afad4c24551748600c164091cf37e068adef76315025d3c78e7",
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"miner": "0x0000000000000000000000000000000000000000",
				"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce": "0x0000000000000000",
				"number": "0x1",
				"parentHash": "0x3ead7923676a44500c46ad2192a0fc084aa42063b1703e6866f138a47fb1a9ca",
				"receiptsRoot": "0x056b23fbba480696b65fe5a59b8f2148a1299103c4f57df839233af2cf4ca2d2",
				"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"size": "0x2bb",
				"stateRoot": "0x6b830601767ac4968163193facbe20123435180e325910b2c50efa21f778c697",
				"timestamp": "0xa",
				"totalDifficulty": "0x1",
				"transactions": [
					"0xdf92bc7c4c0341ecbdcd2a3ca7011fe9e21df4b8553bf0c8caabe6cb4a1aee26"
				],
				"transactionsRoot": "0x87c65a3f1a98dafe282ace11eaf88b8f31bf41fe6794d401d2f986c1af84bcd5",
				"uncles": []
			}`,
		},
		// 22. latest-1 block
		{
			blockHash: &blockHashes[len(blockHashes)-2],
			fullTx:    true,
			want: `{
				"baseFeePerGas": "0x5d21dba00",
				"blockGasCost": "0x0",
				"difficulty": "0x1",
				"extraData": "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"gasLimit": "0x7a1200",
				"gasUsed": "0x5208",
				"hash": "0x0583a9d630632001771b4ecc7d62574aec3825aff47e2a680b0ea4ddb79e7365",
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"miner": "0x0000000000000000000000000000000000000000",
				"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce": "0x0000000000000000",
				"number": "0x9",
				"parentHash": "0x2fab5c6892c66668842683ced6b384c2ee83bfd6096a58f451290cabaf57a63e",
				"receiptsRoot": "0x056b23fbba480696b65fe5a59b8f2148a1299103c4f57df839233af2cf4ca2d2",
				"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"size": "0x2bb",
				"stateRoot": "0x3703d70c6443e809ce035c2a8212dbf9813f6b7d1b0f597766e9023867a852f5",
				"timestamp": "0x5a",
				"totalDifficulty": "0x9",
				"transactions": [
					{
						"blockHash": "0x0583a9d630632001771b4ecc7d62574aec3825aff47e2a680b0ea4ddb79e7365",
						"blockNumber": "0x9",
						"from": "0x703c4b2bd70c169f5717101caee543299fc946c7",
						"gas": "0x5208",
						"gasPrice": "0x5d21dba00",
						"hash": "0x237f95840187a93f8aaf8d6f1515f8a8ac9d9359fcb0c220cdb3d642d6b9a19a",
						"input": "0x",
						"nonce": "0x8",
						"to": "0x0d3ab14bbad3d99f4203bd7a11acb94882050e7e",
						"transactionIndex": "0x0",
						"value": "0x3e8",
						"type": "0x0",
						"v": "0x1c",
						"r": "0xd7cdc527490b7ba29c515aae3bbe80c67729cda7f736e6515652cfc40e9da68f",
						"s": "0x4d0a4a59bef165b16f910bdadd41efaaad1b73549bacc35eaf6d073eb1fb92b7"
					}
				],
				"transactionsRoot": "0xe16929d9c7efab0f962c1ed8c1295ddff42d3026779ed1318ea079ca580ee4cb",
				"uncles": []
			}`,
		},
		// 23. empty hash + body
		{
			blockHash: &common.Hash{},
			fullTx:    true,
			want:      "null",
		},
		// 24. pending block
		{
			blockHash: &pendingHash,
			want:      `null`,
		},
		// 25. pending block + fullTx
		{
			blockHash: &pendingHash,
			fullTx:    true,
			want:      `null`,
		},
	}

	for i, tt := range testSuite {
		var (
			result map[string]interface{}
			err    error
		)
		if tt.blockHash != nil {
			if tt.reqHeader {
				result = api.GetHeaderByHash(context.Background(), *tt.blockHash)
			} else {
				result, err = api.GetBlockByHash(context.Background(), *tt.blockHash, tt.fullTx)
			}
		} else {
			if tt.reqHeader {
				result, err = api.GetHeaderByNumber(context.Background(), tt.blockNumber)
			} else {
				result, err = api.GetBlockByNumber(context.Background(), tt.blockNumber, tt.fullTx)
			}
		}
		if tt.expectErr != nil {
			if err == nil {
				t.Errorf("test %d: want error %v, have nothing", i, tt.expectErr)
				continue
			}
			if !errors.Is(err, tt.expectErr) {
				t.Errorf("test %d: error mismatch, want %v, have %v", i, tt.expectErr, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("test %d: want no error, have %v", i, err)
			continue
		}
		data, err := json.Marshal(result)
		if err != nil {
			t.Errorf("test %d: json marshal error", i)
			continue
		}
		want, have := tt.want, string(data)
		require.JSONEqf(t, want, have, "test %d: json not match, want: %s, have: %s", i, want, have)
	}
}

func TestRPCGetTransactionReceipt(t *testing.T) {
	t.Parallel()

	// Initialize test accounts
	var (
		acc1Key, _ = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		acc2Key, _ = crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
		acc1Addr   = crypto.PubkeyToAddress(acc1Key.PublicKey)
		acc2Addr   = crypto.PubkeyToAddress(acc2Key.PublicKey)
		contract   = common.HexToAddress("0000000000000000000000000000000000031ec7")
		genesis    = &core.Genesis{
			Config: params.TestSubnetEVMConfig,
			Alloc: core.GenesisAlloc{
				acc1Addr: {Balance: big.NewInt(params.Ether)},
				acc2Addr: {Balance: big.NewInt(params.Ether)},
				// // SPDX-License-Identifier: GPL-3.0
				// pragma solidity >=0.7.0 <0.9.0;
				//
				// contract Token {
				//     event Transfer(address indexed from, address indexed to, uint256 value);
				//     function transfer(address to, uint256 value) public returns (bool) {
				//         emit Transfer(msg.sender, to, value);
				//         return true;
				//     }
				// }
				contract: {Balance: big.NewInt(params.Ether), Code: common.FromHex("0x608060405234801561001057600080fd5b506004361061002b5760003560e01c8063a9059cbb14610030575b600080fd5b61004a6004803603810190610045919061016a565b610060565b60405161005791906101c5565b60405180910390f35b60008273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516100bf91906101ef565b60405180910390a36001905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610101826100d6565b9050919050565b610111816100f6565b811461011c57600080fd5b50565b60008135905061012e81610108565b92915050565b6000819050919050565b61014781610134565b811461015257600080fd5b50565b6000813590506101648161013e565b92915050565b60008060408385031215610181576101806100d1565b5b600061018f8582860161011f565b92505060206101a085828601610155565b9150509250929050565b60008115159050919050565b6101bf816101aa565b82525050565b60006020820190506101da60008301846101b6565b92915050565b6101e981610134565b82525050565b600060208201905061020460008301846101e0565b9291505056fea2646970667358221220b469033f4b77b9565ee84e0a2f04d496b18160d26034d54f9487e57788fd36d564736f6c63430008120033")},
			},
		}
		genBlocks = 5
		signer    = types.LatestSignerForChainID(params.TestChainConfig.ChainID)
		txHashes  = make([]common.Hash, genBlocks)
	)
	backend := newTestBackend(t, genBlocks, genesis, func(i int, b *core.BlockGen) {
		var (
			tx  *types.Transaction
			err error
		)
		switch i {
		case 0:
			// transfer 1000wei
			tx, err = types.SignTx(types.NewTx(&types.LegacyTx{Nonce: uint64(i), To: &acc2Addr, Value: big.NewInt(1000), Gas: params.TxGas, GasPrice: b.BaseFee(), Data: nil}), types.HomesteadSigner{}, acc1Key)
		case 1:
			// create contract
			tx, err = types.SignTx(types.NewTx(&types.LegacyTx{Nonce: uint64(i), To: nil, Gas: 53100, GasPrice: b.BaseFee(), Data: common.FromHex("0x60806040")}), signer, acc1Key)
		case 2:
			// with logs
			// transfer(address to, uint256 value)
			data := fmt.Sprintf("0xa9059cbb%s%s", common.HexToHash(common.BigToAddress(big.NewInt(int64(i + 1))).Hex()).String()[2:], common.BytesToHash([]byte{byte(i + 11)}).String()[2:])
			tx, err = types.SignTx(types.NewTx(&types.LegacyTx{Nonce: uint64(i), To: &contract, Gas: 60000, GasPrice: b.BaseFee(), Data: common.FromHex(data)}), signer, acc1Key)
		case 3:
			// dynamic fee with logs
			// transfer(address to, uint256 value)
			data := fmt.Sprintf("0xa9059cbb%s%s", common.HexToHash(common.BigToAddress(big.NewInt(int64(i + 1))).Hex()).String()[2:], common.BytesToHash([]byte{byte(i + 11)}).String()[2:])
			fee := big.NewInt(500)
			fee.Add(fee, b.BaseFee())
			tx, err = types.SignTx(types.NewTx(&types.DynamicFeeTx{Nonce: uint64(i), To: &contract, Gas: 60000, Value: big.NewInt(1), GasTipCap: big.NewInt(500), GasFeeCap: fee, Data: common.FromHex(data)}), signer, acc1Key)
		case 4:
			// access list with contract create
			accessList := types.AccessList{{
				Address:     contract,
				StorageKeys: []common.Hash{{0}},
			}}
			tx, err = types.SignTx(types.NewTx(&types.AccessListTx{Nonce: uint64(i), To: nil, Gas: 58100, GasPrice: b.BaseFee(), Data: common.FromHex("0x60806040"), AccessList: accessList}), signer, acc1Key)
		}
		if err != nil {
			t.Errorf("failed to sign tx: %v", err)
		}
		if tx != nil {
			b.AddTx(tx)
			txHashes[i] = tx.Hash()
		}
	})
	api := NewTransactionAPI(backend, new(AddrLocker))
	blockHashes := make([]common.Hash, genBlocks+1)
	ctx := context.Background()
	for i := 0; i <= genBlocks; i++ {
		header, err := backend.HeaderByNumber(ctx, rpc.BlockNumber(i))
		if err != nil {
			t.Errorf("failed to get block: %d err: %v", i, err)
		}
		blockHashes[i] = header.Hash()
	}

	var testSuite = []struct {
		txHash common.Hash
		want   string
	}{
		// 0. normal success
		{
			txHash: txHashes[0],
			want: `{
				"blockHash": "0xcc27e155b6eadfa892992a8cd8adaf3c929a6ec6d98c4dfbc60258883c73568e",
				"blockNumber": "0x1",
				"contractAddress": null,
				"cumulativeGasUsed": "0x5208",
				"effectiveGasPrice": "0x5d21dba00",
				"from": "0x703c4b2bd70c169f5717101caee543299fc946c7",
				"gasUsed": "0x5208",
				"logs": [],
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"status": "0x1",
				"to": "0x0d3ab14bbad3d99f4203bd7a11acb94882050e7e",
				"transactionHash": "0xdf92bc7c4c0341ecbdcd2a3ca7011fe9e21df4b8553bf0c8caabe6cb4a1aee26",
				"transactionIndex": "0x0",
				"type": "0x0"
			  }`,
		},
		// 1. create contract
		{
			txHash: txHashes[1],
			want: `{
				"blockHash": "0xbea66b509ec6e5639279ca696def697d47d0c40ecfa00bbcdb5e31a492491c83",
				"blockNumber": "0x2",
				"contractAddress": "0xae9bea628c4ce503dcfd7e305cab4e29e7476592",
				"cumulativeGasUsed": "0xcf4e",
				"effectiveGasPrice": "0x5d21dba00",
				"from": "0x703c4b2bd70c169f5717101caee543299fc946c7",
				"gasUsed": "0xcf4e",
				"logs": [],
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"status": "0x1",
				"to": null,
				"transactionHash": "0x22aa617165f83a9f8c191c2b7724ae43eeb1249bee06c98c03c7624c21d27dc8",
				"transactionIndex": "0x0",
				"type": "0x0"
			}`,
		},
		// 2. with logs success
		{
			txHash: txHashes[2],
			want: `{
				"blockHash": "0x139eee6b02792c6bee20be4d0aa72b3876f22cf8fe8e2bf45e1a0cee94aa3cf1",
				"blockNumber": "0x3",
				"contractAddress": null,
				"cumulativeGasUsed": "0x5e28",
				"effectiveGasPrice": "0x5d21dba00",
				"from": "0x703c4b2bd70c169f5717101caee543299fc946c7",
				"gasUsed": "0x5e28",
				"logs": [
					{
						"address": "0x0000000000000000000000000000000000031ec7",
						"topics": [
							"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
							"0x000000000000000000000000703c4b2bd70c169f5717101caee543299fc946c7",
							"0x0000000000000000000000000000000000000000000000000000000000000003"
						],
						"data": "0x000000000000000000000000000000000000000000000000000000000000000d",
						"blockNumber": "0x3",
						"transactionHash": "0x7366a7738f47e32f5b6d292ca064b6b66f295d3931533a3745975be1191fccdf",
						"transactionIndex": "0x0",
						"blockHash": "0x139eee6b02792c6bee20be4d0aa72b3876f22cf8fe8e2bf45e1a0cee94aa3cf1",
						"logIndex": "0x0",
						"removed": false
					}
				],
				"logsBloom": "0x00000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000800000000000000008000000000000000000000000000000000020000000080000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000000400000000002000000000000800000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000",
				"status": "0x1",
				"to": "0x0000000000000000000000000000000000031ec7",
				"transactionHash": "0x7366a7738f47e32f5b6d292ca064b6b66f295d3931533a3745975be1191fccdf",
				"transactionIndex": "0x0",
				"type": "0x0"
			}`,
		},
		// 3. dynamic tx with logs success
		{
			txHash: txHashes[3],
			want: `{
				"blockHash": "0xba48c351b0aa848ab2ec889f5794f0db779e0840af80472d1c29df54b22288c8",
				"blockNumber": "0x4",
				"contractAddress": null,
				"cumulativeGasUsed": "0x538d",
				"effectiveGasPrice": "0x5d21dbbf4",
				"from": "0x703c4b2bd70c169f5717101caee543299fc946c7",
				"gasUsed": "0x538d",
				"logs": [],
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"status": "0x0",
				"to": "0x0000000000000000000000000000000000031ec7",
				"transactionHash": "0x4e1e9194ca6f9d4e1736e9e441f66104f273548ed6d91b236a5f9c2ea10fa06d",
				"transactionIndex": "0x0",
				"type": "0x2"
			}`,
		},
		// 4. access list tx with create contract
		{
			txHash: txHashes[4],
			want: `{
				"blockHash": "0x83f2712dfaeab6ab6239bf060bccfc49652e4afdc3b80b22a8373816a2047bd3",
				"blockNumber": "0x5",
				"contractAddress": "0xfdaa97661a584d977b4d3abb5370766ff5b86a18",
				"cumulativeGasUsed": "0xe01a",
				"effectiveGasPrice": "0x5d21dba00",
				"from": "0x703c4b2bd70c169f5717101caee543299fc946c7",
				"gasUsed": "0xe01a",
				"logs": [],
				"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				"status": "0x1",
				"to": null,
				"transactionHash": "0x8afe030574f663fe5096371d6f58a6287bfb3e0c73a5050220f5775a08e7abc9",
				"transactionIndex": "0x0",
				"type": "0x1"
			}`,
		},
		// 5. txhash empty
		{
			txHash: common.Hash{},
			want:   `null`,
		},
		// 6. txhash not found
		{
			txHash: common.HexToHash("deadbeef"),
			want:   `null`,
		},
	}

	for i, tt := range testSuite {
		var (
			result interface{}
			err    error
		)
		result, err = api.GetTransactionReceipt(context.Background(), tt.txHash)
		if err != nil {
			t.Errorf("test %d: want no error, have %v", i, err)
			continue
		}
		data, err := json.Marshal(result)
		if err != nil {
			t.Errorf("test %d: json marshal error", i)
			continue
		}
		want, have := tt.want, string(data)
		fmt.Printf("%d: have: %s\n", i, have)
		_ = want
		// require.JSONEqf(t, want, have, "test %d: json not match, want: %s, have: %s", i, want, have)
	}
}

func setupReceiptBackend(t *testing.T, genBlocks int) (*testBackend, []common.Hash) {
	// Initialize test accounts
	var (
		acc1Key, _ = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		acc2Key, _ = crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
		acc1Addr   = crypto.PubkeyToAddress(acc1Key.PublicKey)
		acc2Addr   = crypto.PubkeyToAddress(acc2Key.PublicKey)
		contract   = common.HexToAddress("0000000000000000000000000000000000031ec7")
		genesis    = &core.Genesis{
			Config: params.TestChainConfig,
			Alloc: core.GenesisAlloc{
				acc1Addr: {Balance: big.NewInt(params.Ether)},
				acc2Addr: {Balance: big.NewInt(params.Ether)},
				// // SPDX-License-Identifier: GPL-3.0
				// pragma solidity >=0.7.0 <0.9.0;
				//
				// contract Token {
				//     event Transfer(address indexed from, address indexed to, uint256 value);
				//     function transfer(address to, uint256 value) public returns (bool) {
				//         emit Transfer(msg.sender, to, value);
				//         return true;
				//     }
				// }
				contract: {Balance: big.NewInt(params.Ether), Code: common.FromHex("0x608060405234801561001057600080fd5b506004361061002b5760003560e01c8063a9059cbb14610030575b600080fd5b61004a6004803603810190610045919061016a565b610060565b60405161005791906101c5565b60405180910390f35b60008273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516100bf91906101ef565b60405180910390a36001905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610101826100d6565b9050919050565b610111816100f6565b811461011c57600080fd5b50565b60008135905061012e81610108565b92915050565b6000819050919050565b61014781610134565b811461015257600080fd5b50565b6000813590506101648161013e565b92915050565b60008060408385031215610181576101806100d1565b5b600061018f8582860161011f565b92505060206101a085828601610155565b9150509250929050565b60008115159050919050565b6101bf816101aa565b82525050565b60006020820190506101da60008301846101b6565b92915050565b6101e981610134565b82525050565b600060208201905061020460008301846101e0565b9291505056fea2646970667358221220b469033f4b77b9565ee84e0a2f04d496b18160d26034d54f9487e57788fd36d564736f6c63430008120033")},
			},
		}
		signer   = types.LatestSignerForChainID(params.TestChainConfig.ChainID)
		txHashes = make([]common.Hash, genBlocks)
	)
	backend := newTestBackend(t, genBlocks, genesis, func(i int, b *core.BlockGen) {
		var (
			tx  *types.Transaction
			err error
		)
		switch i {
		case 0:
			// transfer 1000wei
			tx, err = types.SignTx(types.NewTx(&types.LegacyTx{Nonce: uint64(i), To: &acc2Addr, Value: big.NewInt(1000), Gas: params.TxGas, GasPrice: b.BaseFee(), Data: nil}), types.HomesteadSigner{}, acc1Key)
		case 1:
			// create contract
			tx, err = types.SignTx(types.NewTx(&types.LegacyTx{Nonce: uint64(i), To: nil, Gas: 53100, GasPrice: b.BaseFee(), Data: common.FromHex("0x60806040")}), signer, acc1Key)
		case 2:
			// with logs
			// transfer(address to, uint256 value)
			data := fmt.Sprintf("0xa9059cbb%s%s", common.HexToHash(common.BigToAddress(big.NewInt(int64(i + 1))).Hex()).String()[2:], common.BytesToHash([]byte{byte(i + 11)}).String()[2:])
			tx, err = types.SignTx(types.NewTx(&types.LegacyTx{Nonce: uint64(i), To: &contract, Gas: 60000, GasPrice: b.BaseFee(), Data: common.FromHex(data)}), signer, acc1Key)
		case 3:
			// dynamic fee with logs
			// transfer(address to, uint256 value)
			data := fmt.Sprintf("0xa9059cbb%s%s", common.HexToHash(common.BigToAddress(big.NewInt(int64(i + 1))).Hex()).String()[2:], common.BytesToHash([]byte{byte(i + 11)}).String()[2:])
			fee := big.NewInt(500)
			fee.Add(fee, b.BaseFee())
			tx, err = types.SignTx(types.NewTx(&types.DynamicFeeTx{Nonce: uint64(i), To: &contract, Gas: 60000, Value: big.NewInt(1), GasTipCap: big.NewInt(500), GasFeeCap: fee, Data: common.FromHex(data)}), signer, acc1Key)
		case 4:
			// access list with contract create
			accessList := types.AccessList{{
				Address:     contract,
				StorageKeys: []common.Hash{{0}},
			}}
			tx, err = types.SignTx(types.NewTx(&types.AccessListTx{Nonce: uint64(i), To: nil, Gas: 58100, GasPrice: b.BaseFee(), Data: common.FromHex("0x60806040"), AccessList: accessList}), signer, acc1Key)
		}
		if err != nil {
			t.Errorf("failed to sign tx: %v", err)
		}
		if tx != nil {
			b.AddTx(tx)
			txHashes[i] = tx.Hash()
		}
	})
	return backend, txHashes
}

func TestRPCGetBlockReceipts(t *testing.T) {
	t.Parallel()

	var (
		genBlocks  = 5
		backend, _ = setupReceiptBackend(t, genBlocks)
		api        = NewBlockChainAPI(backend)
	)
	blockHashes := make([]common.Hash, genBlocks+1)
	ctx := context.Background()
	for i := 0; i <= genBlocks; i++ {
		header, err := backend.HeaderByNumber(ctx, rpc.BlockNumber(i))
		if err != nil {
			t.Errorf("failed to get block: %d err: %v", i, err)
		}
		blockHashes[i] = header.Hash()
	}

	var testSuite = []struct {
		test rpc.BlockNumberOrHash
		want string
	}{
		// 0. block without any txs(hash)
		{
			test: rpc.BlockNumberOrHashWithHash(blockHashes[0], false),
			want: `[]`,
		},
		// 1. block without any txs(number)
		{
			test: rpc.BlockNumberOrHashWithNumber(0),
			want: `[]`,
		},
		// 2. earliest tag
		{
			test: rpc.BlockNumberOrHashWithNumber(rpc.EarliestBlockNumber),
			want: `[]`,
		},
		// 3. latest tag
		{
			test: rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber),
			want: `[{"blockHash":"0x581862f55d194bd8945820d0d5f28c3a53a72fffa17c0e0376e43478c2e2a75e","blockNumber":"0x5","contractAddress":"0xfdaa97661a584d977b4d3abb5370766ff5b86a18","cumulativeGasUsed":"0xe01c","effectiveGasPrice":"0x5d21dba00","from":"0x703c4b2bd70c169f5717101caee543299fc946c7","gasUsed":"0xe01c","logs":[],"logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","status":"0x1","to":null,"transactionHash":"0x8afe030574f663fe5096371d6f58a6287bfb3e0c73a5050220f5775a08e7abc9","transactionIndex":"0x0","type":"0x1"}]`,
		},
		// 4. block with legacy transfer tx(hash)
		{
			test: rpc.BlockNumberOrHashWithHash(blockHashes[1], false),
			want: `[{"blockHash":"0xcc27e155b6eadfa892992a8cd8adaf3c929a6ec6d98c4dfbc60258883c73568e","blockNumber":"0x1","contractAddress":null,"cumulativeGasUsed":"0x5208","effectiveGasPrice":"0x5d21dba00","from":"0x703c4b2bd70c169f5717101caee543299fc946c7","gasUsed":"0x5208","logs":[],"logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","status":"0x1","to":"0x0d3ab14bbad3d99f4203bd7a11acb94882050e7e","transactionHash":"0xdf92bc7c4c0341ecbdcd2a3ca7011fe9e21df4b8553bf0c8caabe6cb4a1aee26","transactionIndex":"0x0","type":"0x0"}]`,
		},
		// 5. block with contract create tx(number)
		{
			test: rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(2)),
			want: `[{"blockHash":"0x5dcac7b6d680226ef9429ad15f726aa051d422365adaae65247cec38b42788cd","blockNumber":"0x2","contractAddress":"0xae9bea628c4ce503dcfd7e305cab4e29e7476592","cumulativeGasUsed":"0xcf50","effectiveGasPrice":"0x5d21dba00","from":"0x703c4b2bd70c169f5717101caee543299fc946c7","gasUsed":"0xcf50","logs":[],"logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","status":"0x1","to":null,"transactionHash":"0x22aa617165f83a9f8c191c2b7724ae43eeb1249bee06c98c03c7624c21d27dc8","transactionIndex":"0x0","type":"0x0"}]`,
		},
		// 6. block with legacy contract call tx(hash)
		{
			test: rpc.BlockNumberOrHashWithHash(blockHashes[3], false),
			want: `[{"blockHash":"0x2517532d72847621eb37804ae02257e8b3b9bf9668d3e6f7d20dbd24d7dfe863","blockNumber":"0x3","contractAddress":null,"cumulativeGasUsed":"0x5e28","effectiveGasPrice":"0x5d21dba00","from":"0x703c4b2bd70c169f5717101caee543299fc946c7","gasUsed":"0x5e28","logs":[{"address":"0x0000000000000000000000000000000000031ec7","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","0x000000000000000000000000703c4b2bd70c169f5717101caee543299fc946c7","0x0000000000000000000000000000000000000000000000000000000000000003"],"data":"0x000000000000000000000000000000000000000000000000000000000000000d","blockNumber":"0x3","transactionHash":"0x7366a7738f47e32f5b6d292ca064b6b66f295d3931533a3745975be1191fccdf","transactionIndex":"0x0","blockHash":"0x2517532d72847621eb37804ae02257e8b3b9bf9668d3e6f7d20dbd24d7dfe863","logIndex":"0x0","removed":false}],"logsBloom":"0x00000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000800000000000000008000000000000000000000000000000000020000000080000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000000400000000002000000000000800000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000","status":"0x1","to":"0x0000000000000000000000000000000000031ec7","transactionHash":"0x7366a7738f47e32f5b6d292ca064b6b66f295d3931533a3745975be1191fccdf","transactionIndex":"0x0","type":"0x0"}]`,
		},
		// 7. block with dynamic fee tx(number)
		{
			test: rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(4)),
			want: `[{"blockHash":"0xf3ec30848127f600b396b9cf7ead62b3baacbdf0fbf32a419cc4837b362d4e33","blockNumber":"0x4","contractAddress":null,"cumulativeGasUsed":"0x538d","effectiveGasPrice":"0x5d21dbbf4","from":"0x703c4b2bd70c169f5717101caee543299fc946c7","gasUsed":"0x538d","logs":[],"logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","status":"0x0","to":"0x0000000000000000000000000000000000031ec7","transactionHash":"0x4e1e9194ca6f9d4e1736e9e441f66104f273548ed6d91b236a5f9c2ea10fa06d","transactionIndex":"0x0","type":"0x2"}]`,
		},
		// 8. block is empty
		{
			test: rpc.BlockNumberOrHashWithHash(common.Hash{}, false),
			want: `null`,
		},
		// 9. block is not found
		{
			test: rpc.BlockNumberOrHashWithHash(common.HexToHash("deadbeef"), false),
			want: `null`,
		},
		// 10. block is not found
		{
			test: rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(genBlocks + 1)),
			want: `null`,
		},
	}

	for i, tt := range testSuite {
		var (
			result interface{}
			err    error
		)
		result, err = api.GetBlockReceipts(context.Background(), tt.test)
		if err != nil {
			t.Errorf("test %d: want no error, have %v", i, err)
			continue
		}
		data, err := json.Marshal(result)
		if err != nil {
			t.Errorf("test %d: json marshal error", i)
			continue
		}
		want, have := tt.want, string(data)
		require.JSONEqf(t, want, have, "test %d: json not match, want: %s, have: %s", i, want, have)
	}
}
