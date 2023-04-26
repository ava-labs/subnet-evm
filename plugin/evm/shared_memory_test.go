// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"math/big"
	"testing"

	"github.com/ava-labs/avalanchego/chains/atomic"
	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/consensus/snowman"
	engCommon "github.com/ava-labs/avalanchego/snow/engine/common"
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/ava-labs/avalanchego/vms/components/avax"
	"github.com/ava-labs/avalanchego/vms/components/chain"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/subnet-evm/consensus/dummy"
	"github.com/ava-labs/subnet-evm/constants"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethdb/memorydb"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contracts/sharedmemory"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ava-labs/subnet-evm/utils/codec"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func setupVMWithSharedMemory(t *testing.T) (*snow.Context, *VM, chan engCommon.Message, *atomic.Memory) {
	t.Helper()
	// create genesis with the shared memory precompile enabled
	genesis := &core.Genesis{}
	require.NoError(t, genesis.UnmarshalJSON([]byte(genesisJSONSubnetEVM)))
	genesis.Config.GenesisPrecompiles = params.Precompiles{
		sharedmemory.ConfigKey: sharedmemory.NewConfig(common.Big0),
	}
	genesisJSON, err := genesis.MarshalJSON()
	require.NoError(t, err)

	ctx, dbManager, genesisBytes, issuer, atomicMemory := setupGenesis(t, string(genesisJSON))

	// initialize and configure the VM
	vm := &VM{}
	err = vm.Initialize(
		context.Background(),
		ctx,
		dbManager,
		genesisBytes,
		nil,
		nil,
		issuer,
		[]*engCommon.Fx{},
		nil,
	)
	require.NoError(t, err)

	return ctx, vm, issuer, atomicMemory
}

type exportTest struct {
	assetID    ids.ID
	utxoAmount uint64
	avaxSent   uint64
}

func (et exportTest) run(t *testing.T) {
	require := require.New(t)

	// Initialize the VM
	ctx, vm, issuer, atomicMemory := setupVMWithSharedMemory(t)
	defer func() { require.NoError(vm.Shutdown(context.Background())) }()
	require.NoError(vm.SetState(context.Background(), snow.Bootstrapping))
	require.NoError(vm.SetState(context.Background(), snow.NormalOp))
	newTxPoolHeadChan := make(chan core.NewTxPoolReorgEvent, 1)
	vm.txPool.SubscribeNewReorgEvent(newTxPoolHeadChan)

	// Get the starting balance so we can check that the correct final balance.
	state, err := vm.blockChain.State()
	require.NoError(err)
	startingBalance := state.GetBalance(testEthAddrs[0])

	// Note: the output address specified here must correspond to the account that should control
	// the exported funds on the recipient chain. If this were to be imported to the X-Chain, we
	// would need to use an address derived from an X-Chain private key.
	outAddr := common.Address{0xff}

	// Prepare the data for the transaction
	var txData []byte
	if et.assetID == vm.ctx.AVAXAssetID {
		txData, err = sharedmemory.PackExportAVAX(sharedmemory.ExportAVAXInput{
			DestinationChainID: testXChainID,
			Locktime:           uint64(0),
			Threshold:          uint64(1),
			Addrs:              []common.Address{outAddr},
		})
	} else {
		txData, err = sharedmemory.PackExportUTXO(sharedmemory.ExportUTXOInput{
			DestinationChainID: testXChainID,
			Locktime:           uint64(0),
			Threshold:          uint64(1),
			Addrs:              []common.Address{outAddr},
			Amount:             et.utxoAmount,
		})
	}
	require.NoError(err)

	// Sign and submit the transaction
	tx := types.NewTransaction(
		uint64(0),
		sharedmemory.ContractAddress,
		new(big.Int).SetUint64(et.avaxSent),
		200_000,
		big.NewInt(testMinGasPrice),
		txData,
	)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(vm.chainConfig.ChainID), testKeys[0])
	require.NoError(err)
	txErrors := vm.txPool.AddRemotesSync([]*types.Transaction{signedTx})
	for _, err := range txErrors {
		require.NoError(err)
	}

	// Subscribe to logs so we can verify the expected log
	logsCh := make(chan []*types.Log, 1)
	vm.blockChain.SubscribeLogsEvent(logsCh)

	// Build and accept the block
	blk := issueAndAccept(t, issuer, vm)
	newHead := <-newTxPoolHeadChan
	require.Equal(newHead.Head.Hash(), common.Hash(blk.ID()))

	// Verify the expected log
	logs := <-logsCh
	require.Len(logs, 1)
	ethBlock := blk.(*chain.BlockWrapper).Block.(*Block).ethBlock
	require.Equal(ethBlock.Transactions()[0].Hash(), logs[0].TxHash)
	// TODO: maybe verify this through the Accepter interface?

	// Find the X-Chain's shared memory and verify it contains the expected UTXO.
	xChainSharedMemory := atomicMemory.NewSharedMemory(ctx.XChainID)
	values, _, _, err := xChainSharedMemory.Indexed(ctx.ChainID, [][]byte{outAddr[:]}, nil, nil, 100)
	require.NoError(err)
	require.Len(values, 1)
	utxo := &avax.UTXO{}
	version, err := codec.Codec.Unmarshal(values[0], utxo)
	require.NoError(err)
	require.Equal(codec.CodecVersion, version)

	expectedUTXO := &avax.UTXO{
		// Derive unique UTXOID from txHash and log index
		UTXOID: avax.UTXOID{
			TxID:        ids.ID(signedTx.Hash()),
			OutputIndex: uint32(0),
		},
		Asset: avax.Asset{ID: et.assetID},
		Out: &secp256k1fx.TransferOutput{
			Amt: et.utxoAmount,
			OutputOwners: secp256k1fx.OutputOwners{
				Locktime:  0,
				Threshold: uint32(1),
				Addrs:     []ids.ShortID{ids.ShortID(outAddr)},
			},
		},
	}

	expectedUTXOBytes, err := codec.Codec.Marshal(codec.CodecVersion, expectedUTXO)
	require.NoError(err)
	require.Equal(expectedUTXOBytes, values[0])

	// Check the balance is has decreased by expected amount of exported
	// AVAX and the fees paid.
	gasPrice := new(big.Int).Add(
		ethBlock.BaseFee(),
		tx.EffectiveGasTipValue(ethBlock.BaseFee()),
	)
	feesPaid := new(big.Int).Mul(
		gasPrice,
		new(big.Int).SetUint64(ethBlock.GasUsed()), // Note there is only 1 tx in the block
	)
	state, err = vm.blockChain.State()
	require.NoError(err)
	balance := state.GetBalance(testEthAddrs[0])
	expected := new(big.Int).Sub(startingBalance, new(big.Int).SetUint64(et.avaxSent))
	expected.Sub(expected, feesPaid)
	require.Equal(expected, balance)
}

func TestExportAssets(t *testing.T) {
	tests := map[string]exportTest{
		"export AVAX": {
			assetID:    testAvaxAssetID,
			utxoAmount: uint64(1_000_000_000),
			avaxSent:   params.Ether,
		},
		"export non-AVAX": {
			assetID:    ids.ID(sharedmemory.CalculateANTAssetID(common.Hash(testCChainID), testEthAddrs[0])),
			utxoAmount: 1,
			avaxSent:   0,
		},
	}
	for name, et := range tests {
		t.Run(name, et.run)
	}
}

type importTest struct {
	assetID            ids.ID
	amount             uint64
	expectedAVAXImport uint64
}

// Testing imports is a bit more complicated, since the VM
// needs to account to prevent double spends in transactions
// that occur in the same block prior to the current tx, and
// also to prevent double spends in transactions that occur
// previously verified but yet unaccessed ancestor blocks.
func (it importTest) run(t *testing.T) {
	require := require.New(t)

	// Initialize the VM
	ctx, vm, _, atomicMemory := setupVMWithSharedMemory(t)
	defer func() { require.NoError(vm.Shutdown(context.Background())) }()
	require.NoError(vm.SetState(context.Background(), snow.Bootstrapping))
	require.NoError(vm.SetState(context.Background(), snow.NormalOp))
	newTxPoolHeadChan := make(chan core.NewTxPoolReorgEvent, 1)
	vm.txPool.SubscribeNewReorgEvent(newTxPoolHeadChan)

	// Get the starting balance so we can check that the correct amount was imported
	// minus fees.
	state, err := vm.blockChain.State()
	require.NoError(err)
	startingBalance := state.GetBalance(testEthAddrs[0])

	// Derive a txid from the test key's address
	txID, err := ids.ToID(hashing.ComputeHash256(testEthAddrs[0][:]))
	require.NoError(err)

	// Add specified asset to the X Chain's shared memory
	utxo, err := addUTXO(atomicMemory, ctx, txID, 0, it.assetID, it.amount, ids.ShortID(testEthAddrs[0]))
	require.NoError(err)

	// Prepare the data for the transaction
	var txData []byte
	if it.assetID == ctx.AVAXAssetID {
		txData, err = sharedmemory.PackImportAVAX(sharedmemory.ImportAVAXInput{
			SourceChain: ctx.XChainID,
			UtxoID:      utxo.ID,
		})
	} else {
		txData, err = sharedmemory.PackImportUTXO(sharedmemory.ImportUTXOInput{
			SourceChain: ctx.XChainID,
			UtxoID:      utxo.ID,
		})
	}
	require.NoError(err)

	// We need to name the utxo we are importing in the access list.
	atomicPredicate := &sharedmemory.AtomicPredicate{
		SourceChain:   ctx.XChainID,
		ImportedUTXOs: []*avax.UTXO{utxo}, // TODO: I think this should be a slice of UTXOIDs
	}
	atomicPredicateBytes, err := codec.Codec.Marshal(codec.CodecVersion, atomicPredicate)
	require.NoError(err)
	accessList := types.AccessList{
		types.AccessTuple{
			Address:     sharedmemory.ContractAddress,
			StorageKeys: utils.BytesToHashSlice(utils.PackPredicate(atomicPredicateBytes)),
		},
	}

	// This function creates and signs the type of transaction that we will be testing:
	mkTx := func(nonce uint64, value *big.Int) *types.Transaction {
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:    vm.chainConfig.ChainID,
			Nonce:      nonce,
			Gas:        200_000,
			GasTipCap:  common.Big0,
			GasFeeCap:  big.NewInt(testMinGasPrice),
			To:         &sharedmemory.ContractAddress,
			Value:      value,
			Data:       txData,
			AccessList: accessList,
		})
		signedTx, err := types.SignTx(tx, types.LatestSigner(vm.chainConfig), testKeys[0])
		require.NoError(err)
		return signedTx
	}

	// Let's try to spend the same UTXO multiple times:
	// - In the same block
	// - In the next block, too
	//
	// We use GenerateChain to make these blocks since the VM will not
	// build the second block through TX issuance.
	var (
		numBlocks  = 2
		txPerBlock = 2
		gap        = uint64(2)
		tempDB     = memorydb.New()
		genesis    = vm.ethConfig.Genesis.ToBlock(tempDB)
	)
	blocks, allReceipts, err := core.GenerateChain(vm.chainConfig, genesis, dummy.NewETHFaker(), tempDB, numBlocks, gap, func(n int, b *core.BlockGen) {
		// Block must have proper coinbase address to pass syntactic validation
		b.SetCoinbase(constants.BlackholeAddr)
		for i := 0; i < txPerBlock; i++ {
			// Each tx will attempt to spend the same UTXO
			b.AddTx(mkTx(uint64(n*txPerBlock+i), common.Big0))
		}
	})
	require.NoError(err)
	require.Len(blocks, numBlocks)

	for i, receipts := range allReceipts {
		for j, receipt := range receipts {
			if i == 0 && j == 0 {
				// The first tx in the first block should succeed
				require.Equal(types.ReceiptStatusSuccessful, receipt.Status)
				continue
			}
			// All other txs should fail
			require.Equal(types.ReceiptStatusFailed, receipt.Status)
		}
	}

	// Subscribe to logs so we can verify the expected log
	logsCh := make(chan []*types.Log, 1)
	vm.blockChain.SubscribeLogsEvent(logsCh)

	// Now we can verify these blocks
	vmBlks := make([]snowman.Block, len(blocks))
	for i, blk := range blocks {
		ctx := context.Background()
		vmBlks[i] = vm.newBlock(blk)
		require.NoError(vmBlks[i].Verify(ctx))
		require.NoError(vm.SetPreference(ctx, vmBlks[i].ID()))

		// The block should process as head in the tx pool as well.
		newHead := <-newTxPoolHeadChan
		require.Equal(common.Hash(vmBlks[i].ID()), newHead.Head.Hash())
	}

	// Verify the expected log
	logs := <-logsCh
	require.Len(logs, 1)
	require.Equal(blocks[0].Transactions()[0].Hash(), logs[0].TxHash)
	// TODO: maybe verify this through the Accepter interface?

	// The mempool will not accept a tx that spends this UTXO,
	// even though the blocks are not accepted yet.
	nonce, err := vm.GetCurrentNonce(testEthAddrs[0])
	require.NoError(err)
	tx := mkTx(nonce, common.Big0)
	errs := vm.txPool.AddRemotesSync([]*types.Transaction{tx})
	require.Len(errs, 1)
	require.ErrorIs(errs[0], vmerrs.ErrNamedUTXOSpent)

	// Verify the UTXO is still in the shared memory before block acceptance
	inputID := utxo.InputID()
	_, err = vm.ctx.SharedMemory.Get(vm.ctx.XChainID, [][]byte{inputID[:]})
	require.NoError(err)

	// Now we can accept the blocks
	for _, vmBlk := range vmBlks {
		require.NoError(vmBlk.Accept(context.Background()))
	}

	// Verify the UTXO was removed after block acceptance
	_, err = vm.ctx.SharedMemory.Get(vm.ctx.XChainID, [][]byte{inputID[:]})
	require.ErrorIs(err, database.ErrNotFound)

	// The mempool will not accept a tx that spends this UTXO,
	// even though the blocks are not accepted yet.
	tx = mkTx(nonce, common.Big0)
	errs = vm.txPool.AddRemotesSync([]*types.Transaction{tx})
	require.Len(errs, 1)
	require.ErrorIs(errs[0], vmerrs.ErrNamedUTXOSpent)

	// Check the balance is has increased by expected amount of imported
	// AVAX minus the fees paid.
	feesPaid := new(big.Int)
	for i, block := range blocks {
		feesPaid.Add(feesPaid, totalFees(block, allReceipts[i]))
	}
	state, err = vm.blockChain.State()
	require.NoError(err)
	balance := state.GetBalance(testEthAddrs[0])
	expected := new(big.Int).Add(startingBalance, new(big.Int).SetUint64(it.expectedAVAXImport))
	expected.Sub(expected, feesPaid)
	require.Equal(expected, balance)
}

func TestImportAssets(t *testing.T) {
	tests := map[string]importTest{
		"import AVAX": {
			assetID:            testAvaxAssetID,
			amount:             uint64(1_000_000_000),
			expectedAVAXImport: params.Ether,
		},
		"import non-AVAX": {
			assetID:            ids.ID(sharedmemory.CalculateANTAssetID(common.Hash(testCChainID), testEthAddrs[0])),
			amount:             1,
			expectedAVAXImport: 0,
		},
	}
	for name, it := range tests {
		t.Run(name, it.run)
	}
}

func addUTXO(sharedMemory *atomic.Memory, ctx *snow.Context, txID ids.ID, index uint32, assetID ids.ID, amount uint64, addr ids.ShortID) (*avax.UTXO, error) {
	utxo := &avax.UTXO{
		UTXOID: avax.UTXOID{
			TxID:        txID,
			OutputIndex: index,
		},
		Asset: avax.Asset{ID: assetID},
		Out: &secp256k1fx.TransferOutput{
			Amt: amount,
			OutputOwners: secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs:     []ids.ShortID{addr},
			},
		},
	}
	utxoBytes, err := codec.Codec.Marshal(codec.CodecVersion, utxo)
	if err != nil {
		return nil, err
	}

	xChainSharedMemory := sharedMemory.NewSharedMemory(ctx.XChainID)
	inputID := utxo.InputID()
	if err := xChainSharedMemory.Apply(map[ids.ID]*atomic.Requests{ctx.ChainID: {PutRequests: []*atomic.Element{{
		Key:   inputID[:],
		Value: utxoBytes,
		Traits: [][]byte{
			addr.Bytes(),
		},
	}}}}); err != nil {
		return nil, err
	}

	return utxo, nil
}

func totalFees(block *types.Block, receipts []*types.Receipt) *big.Int {
	feesWei := new(big.Int)
	for i, tx := range block.Transactions() {
		var minerFee *big.Int
		if baseFee := block.BaseFee(); baseFee != nil {
			// Note in coreth the coinbase payment is (baseFee + effectiveGasTip) * gasUsed
			minerFee = new(big.Int).Add(baseFee, tx.EffectiveGasTipValue(baseFee))
		} else {
			// Prior to activation of EIP-1559, the coinbase payment was gasPrice * gasUsed
			minerFee = tx.GasPrice()
		}
		feesWei.Add(feesWei, new(big.Int).Mul(new(big.Int).SetUint64(receipts[i].GasUsed), minerFee))
	}
	return feesWei
}
