// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/choices"
	"github.com/ava-labs/avalanchego/snow/engine/snowman/block"
	"github.com/ava-labs/avalanchego/snow/validators"
	"github.com/ava-labs/avalanchego/utils"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/avalanchego/vms/components/chain"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/internal/ethapi"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/rpc"
<<<<<<< HEAD
=======
	subnetEVMUtils "github.com/ava-labs/subnet-evm/utils"
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
	predicateutils "github.com/ava-labs/subnet-evm/utils/predicate"
	warpPayload "github.com/ava-labs/subnet-evm/warp/payload"
	"github.com/ava-labs/subnet-evm/x/warp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

func TestSendWarpMessage(t *testing.T) {
<<<<<<< HEAD
	genesis := &core.Genesis{}
	if err := genesis.UnmarshalJSON([]byte(genesisJSONSubnetEVM)); err != nil {
		t.Fatal(err)
	}
	genesis.Config.GenesisPrecompiles = params.Precompiles{
		warp.ConfigKey: warp.NewDefaultConfig(big.NewInt(0)),
	}
	genesisJSON, err := genesis.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	issuer, vm, _, _ := GenesisVM(t, true, string(genesisJSON), "", "")

	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
	}()

	newTxPoolHeadChan := make(chan core.NewTxPoolReorgEvent, 1)
	reorgSub := vm.txPool.SubscribeNewReorgEvent(newTxPoolHeadChan)
	defer reorgSub.Unsubscribe()

=======
	require := require.New(t)
	genesis := &core.Genesis{}
	require.NoError(genesis.UnmarshalJSON([]byte(genesisJSONSubnetEVM)))
	genesis.Config.GenesisPrecompiles = params.Precompiles{
		warp.ConfigKey: warp.NewDefaultConfig(subnetEVMUtils.NewUint64(0)),
	}
	genesisJSON, err := genesis.MarshalJSON()
	require.NoError(err)
	issuer, vm, _, _ := GenesisVM(t, true, string(genesisJSON), "", "")

	defer func() {
		require.NoError(vm.Shutdown(context.Background()))
	}()

>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
	acceptedLogsChan := make(chan []*types.Log, 10)
	logsSub := vm.eth.APIBackend.SubscribeAcceptedLogsEvent(acceptedLogsChan)
	defer logsSub.Unsubscribe()

	payload := utils.RandomBytes(100)

	warpSendMessageInput, err := warp.PackSendWarpMessage(warp.SendWarpMessageInput{
<<<<<<< HEAD
		DestinationChainID: vm.ctx.CChainID,
		DestinationAddress: testEthAddrs[1].Hash(),
		Payload:            payload,
	})
	require.NoError(t, err)
=======
		DestinationChainID: common.Hash(vm.ctx.CChainID),
		DestinationAddress: testEthAddrs[1],
		Payload:            payload,
	})
	require.NoError(err)
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d

	// Submit a transaction to trigger sending a warp message
	tx0 := types.NewTransaction(uint64(0), warp.ContractAddress, big.NewInt(1), 100_000, big.NewInt(testMinGasPrice), warpSendMessageInput)
	signedTx0, err := types.SignTx(tx0, types.LatestSignerForChainID(vm.chainConfig.ChainID), testKeys[0])
<<<<<<< HEAD
	require.NoError(t, err)

	errs := vm.txPool.AddRemotesSync([]*types.Transaction{signedTx0})
	if err := errs[0]; err != nil {
		t.Fatalf("Failed to add tx at index: %s", err)
	}

	<-issuer
	blk, err := vm.BuildBlock(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if err := blk.Verify(context.Background()); err != nil {
		t.Fatal(err)
	}

	if status := blk.Status(); status != choices.Processing {
		t.Fatalf("Expected status of built block to be %s, but found %s", choices.Processing, status)
	}

	// Verify that the constructed block contains the expected log with an unsigned warp message in the log data
	ethBlock1 := blk.(*chain.BlockWrapper).Block.(*Block).ethBlock
	require.Len(t, ethBlock1.Transactions(), 1)
	receipts := rawdb.ReadReceipts(vm.chaindb, ethBlock1.Hash(), ethBlock1.NumberU64(), vm.chainConfig)
	require.NotNil(t, receipts)
	require.Len(t, receipts, 1)

	logData := receipts[0].Logs[0].Data
	unsignedMessage, err := avalancheWarp.ParseUnsignedMessage(logData)
	require.NoError(t, err)
=======
	require.NoError(err)

	errs := vm.txPool.AddRemotesSync([]*types.Transaction{signedTx0})
	require.NoError(errs[0])

	<-issuer
	blk, err := vm.BuildBlock(context.Background())
	require.NoError(err)

	require.NoError(blk.Verify(context.Background()))

	require.Equal(choices.Processing, blk.Status())

	// Verify that the constructed block contains the expected log with an unsigned warp message in the log data
	ethBlock1 := blk.(*chain.BlockWrapper).Block.(*Block).ethBlock
	require.Len(ethBlock1.Transactions(), 1)
	receipts := rawdb.ReadReceipts(vm.chaindb, ethBlock1.Hash(), ethBlock1.NumberU64(), vm.chainConfig)
	require.Len(receipts, 1)

	require.Len(receipts[0].Logs, 1)
	expectedTopics := []common.Hash{
		warp.WarpABI.Events["SendWarpMessage"].ID,
		common.Hash(vm.ctx.CChainID),
		testEthAddrs[1].Hash(),
		testEthAddrs[0].Hash(),
	}
	require.Equal(expectedTopics, receipts[0].Logs[0].Topics)
	logData := receipts[0].Logs[0].Data
	unsignedMessage, err := avalancheWarp.ParseUnsignedMessage(logData)
	require.NoError(err)
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
	unsignedMessageID := unsignedMessage.ID()

	// Verify the signature cannot be fetched before the block is accepted
	_, err = vm.warpBackend.GetSignature(unsignedMessageID)
<<<<<<< HEAD
	require.Error(t, err)

	if err := vm.SetPreference(context.Background(), blk.ID()); err != nil {
		t.Fatal(err)
	}
	if err := blk.Accept(context.Background()); err != nil {
		t.Fatal(err)
	}
	rawSignatureBytes, err := vm.warpBackend.GetSignature(unsignedMessageID)
	require.NoError(t, err)
	blsSignature, err := bls.SignatureFromBytes(rawSignatureBytes[:])
	require.NoError(t, err)

	select {
	case acceptedLogs := <-acceptedLogsChan:
		require.Len(t, acceptedLogs, 1, "unexpected length of accepted logs")
		require.Equal(t, acceptedLogs[0], receipts[0].Logs[0])
	case <-time.After(time.Second):
		t.Fatal("Failed to read accepted logs from subscription")
	}
	logsSub.Unsubscribe()

	// Verify the produced signature is valid
	require.True(t, bls.Verify(vm.ctx.PublicKey, blsSignature, unsignedMessage.Bytes()))
}

func TestReceiveWarpMessage(t *testing.T) {
	genesis := &core.Genesis{}
	if err := genesis.UnmarshalJSON([]byte(genesisJSONSubnetEVM)); err != nil {
		t.Fatal(err)
	}
	genesis.Config.GenesisPrecompiles = params.Precompiles{
		warp.ConfigKey: warp.NewDefaultConfig(big.NewInt(0)),
	}
	genesisJSON, err := genesis.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	issuer, vm, _, _ := GenesisVM(t, true, string(genesisJSON), "", "")

	defer func() {
		if err := vm.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
	}()

	newTxPoolHeadChan := make(chan core.NewTxPoolReorgEvent, 1)
	reorgSub := vm.txPool.SubscribeNewReorgEvent(newTxPoolHeadChan)
	defer reorgSub.Unsubscribe()

=======
	require.Error(err)

	require.NoError(vm.SetPreference(context.Background(), blk.ID()))
	require.NoError(blk.Accept(context.Background()))
	vm.blockChain.DrainAcceptorQueue()
	rawSignatureBytes, err := vm.warpBackend.GetSignature(unsignedMessageID)
	require.NoError(err)
	blsSignature, err := bls.SignatureFromBytes(rawSignatureBytes[:])
	require.NoError(err)

	select {
	case acceptedLogs := <-acceptedLogsChan:
		require.Len(acceptedLogs, 1, "unexpected length of accepted logs")
		require.Equal(acceptedLogs[0], receipts[0].Logs[0])
	case <-time.After(time.Second):
		require.Fail("Failed to read accepted logs from subscription")
	}

	// Verify the produced signature is valid
	require.True(bls.Verify(vm.ctx.PublicKey, blsSignature, unsignedMessage.Bytes()))
}

func TestReceiveWarpMessage(t *testing.T) {
	require := require.New(t)
	genesis := &core.Genesis{}
	require.NoError(genesis.UnmarshalJSON([]byte(genesisJSONSubnetEVM)))
	genesis.Config.GenesisPrecompiles = params.Precompiles{
		warp.ConfigKey: warp.NewDefaultConfig(subnetEVMUtils.NewUint64(0)),
	}
	genesisJSON, err := genesis.MarshalJSON()
	require.NoError(err)
	issuer, vm, _, _ := GenesisVM(t, true, string(genesisJSON), "", "")

	defer func() {
		require.NoError(vm.Shutdown(context.Background()))
	}()

>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
	acceptedLogsChan := make(chan []*types.Log, 10)
	logsSub := vm.eth.APIBackend.SubscribeAcceptedLogsEvent(acceptedLogsChan)
	defer logsSub.Unsubscribe()

	payload := utils.RandomBytes(100)

	addressedPayload, err := warpPayload.NewAddressedPayload(
<<<<<<< HEAD
		ids.ID(testEthAddrs[0].Hash()),
		ids.ID(testEthAddrs[1].Hash()),
		payload,
	)
	require.NoError(t, err)
	unsignedMessage, err := avalancheWarp.NewUnsignedMessage(
		vm.ctx.ChainID,
		vm.ctx.CChainID,
		addressedPayload.Bytes(),
	)
	require.NoError(t, err)

	nodeID1 := ids.GenerateTestNodeID()
	blsSecretKey1, err := bls.NewSecretKey()
	require.NoError(t, err)
=======
		testEthAddrs[0],
		common.Hash(vm.ctx.CChainID),
		testEthAddrs[1],
		payload,
	)
	require.NoError(err)
	unsignedMessage, err := avalancheWarp.NewUnsignedMessage(
		vm.ctx.NetworkID,
		vm.ctx.ChainID,
		addressedPayload.Bytes(),
	)
	require.NoError(err)

	nodeID1 := ids.GenerateTestNodeID()
	blsSecretKey1, err := bls.NewSecretKey()
	require.NoError(err)
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
	blsPublicKey1 := bls.PublicFromSecretKey(blsSecretKey1)
	blsSignature1 := bls.Sign(blsSecretKey1, unsignedMessage.Bytes())

	nodeID2 := ids.GenerateTestNodeID()
	blsSecretKey2, err := bls.NewSecretKey()
<<<<<<< HEAD
	require.NoError(t, err)
=======
	require.NoError(err)
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
	blsPublicKey2 := bls.PublicFromSecretKey(blsSecretKey2)
	blsSignature2 := bls.Sign(blsSecretKey2, unsignedMessage.Bytes())

	blsAggregatedSignature, err := bls.AggregateSignatures([]*bls.Signature{blsSignature1, blsSignature2})
<<<<<<< HEAD
	require.NoError(t, err)
=======
	require.NoError(err)
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d

	vm.ctx.ValidatorState = &validators.TestState{
		GetSubnetIDF: func(ctx context.Context, chainID ids.ID) (ids.ID, error) {
			return ids.Empty, nil
		},
		GetValidatorSetF: func(ctx context.Context, height uint64, subnetID ids.ID) (map[ids.NodeID]*validators.GetValidatorOutput, error) {
			return map[ids.NodeID]*validators.GetValidatorOutput{
				nodeID1: {
					NodeID:    nodeID1,
					PublicKey: blsPublicKey1,
					Weight:    50,
				},
				nodeID2: {
					NodeID:    nodeID2,
					PublicKey: blsPublicKey2,
					Weight:    50,
				},
			}, nil
		},
	}

	signersBitSet := set.NewBits()
	signersBitSet.Add(0)
	signersBitSet.Add(1)

	warpSignature := &avalancheWarp.BitSetSignature{
		Signers: signersBitSet.Bytes(),
	}

	blsAggregatedSignatureBytes := bls.SignatureToBytes(blsAggregatedSignature)
	copy(warpSignature.Signature[:], blsAggregatedSignatureBytes)

	signedMessage, err := avalancheWarp.NewMessage(
		unsignedMessage,
		warpSignature,
	)
<<<<<<< HEAD
	require.NoError(t, err)

	getWarpMsgInput, err := warp.PackGetVerifiedWarpMessage()
	require.NoError(t, err)
	signedTx1, err := types.SignTx(
=======
	require.NoError(err)

	getWarpMsgInput, err := warp.PackGetVerifiedWarpMessage()
	require.NoError(err)
	getVerifiedWarpMessageTx, err := types.SignTx(
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
		predicateutils.NewPredicateTx(
			vm.chainConfig.ChainID,
			0,
			&warp.Module.Address,
			1_000_000,
			big.NewInt(225*params.GWei),
			big.NewInt(params.GWei),
			common.Big0,
			getWarpMsgInput,
			types.AccessList{},
			warp.ContractAddress,
			signedMessage.Bytes(),
		),
		types.LatestSignerForChainID(vm.chainConfig.ChainID),
		testKeys[0],
	)
<<<<<<< HEAD
	require.NoError(t, err)
	errs := vm.txPool.AddRemotesSync([]*types.Transaction{signedTx1})
	for i, err := range errs {
		if err != nil {
			t.Fatalf("Failed to add tx at index %d: %s", i, err)
		}
	}
	vm.clock.Set(vm.clock.Time().Add(2 * time.Second))
	<-issuer

	expectedOutput, err := warp.PackGetVerifiedWarpMessageOutput(warp.GetVerifiedWarpMessageOutput{
		Message: warp.WarpMessage{
			OriginChainID:       vm.ctx.ChainID,
			OriginSenderAddress: testEthAddrs[0].Hash(),
			DestinationChainID:  vm.ctx.CChainID,
			DestinationAddress:  testEthAddrs[1].Hash(),
=======
	require.NoError(err)
	errs := vm.txPool.AddRemotesSync([]*types.Transaction{getVerifiedWarpMessageTx})
	for i, err := range errs {
		require.NoError(err, "failed to add tx at index %d", i)
	}

	expectedOutput, err := warp.PackGetVerifiedWarpMessageOutput(warp.GetVerifiedWarpMessageOutput{
		Message: warp.WarpMessage{
			OriginChainID:       common.Hash(vm.ctx.ChainID),
			OriginSenderAddress: testEthAddrs[0],
			DestinationChainID:  common.Hash(vm.ctx.CChainID),
			DestinationAddress:  testEthAddrs[1],
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
			Payload:             payload,
		},
		Exists: true,
	})
<<<<<<< HEAD
	require.NoError(t, err)

	// Assert that DoCall returns the expected output
	hexGetWarpMsgInput := hexutil.Bytes(signedTx1.Data())
	hexGasLimit := hexutil.Uint64(signedTx1.Gas())
	accessList := signedTx1.AccessList()
=======
	require.NoError(err)

	// Assert that DoCall returns the expected output
	hexGetWarpMsgInput := hexutil.Bytes(getVerifiedWarpMessageTx.Data())
	hexGasLimit := hexutil.Uint64(getVerifiedWarpMessageTx.Gas())
	accessList := getVerifiedWarpMessageTx.AccessList()
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
	blockNum := new(rpc.BlockNumber)
	*blockNum = rpc.LatestBlockNumber

	executionRes, err := ethapi.DoCall(
		context.Background(),
		vm.eth.APIBackend,
		ethapi.TransactionArgs{
<<<<<<< HEAD
			To:         signedTx1.To(),
=======
			To:         getVerifiedWarpMessageTx.To(),
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
			Input:      &hexGetWarpMsgInput,
			AccessList: &accessList,
			Gas:        &hexGasLimit,
		},
		rpc.BlockNumberOrHash{BlockNumber: blockNum},
		nil,
		time.Second,
		10_000_000,
	)
<<<<<<< HEAD
	require.NoError(t, err)
	require.NoError(t, executionRes.Err)
	require.Equal(t, expectedOutput, executionRes.ReturnData)
=======
	require.NoError(err)
	require.NoError(executionRes.Err)
	require.Equal(expectedOutput, executionRes.ReturnData)
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d

	// Build, verify, and accept block with valid proposer context.
	validProposerCtx := &block.Context{
		PChainHeight: 10,
	}
<<<<<<< HEAD
	block2, err := vm.BuildBlockWithContext(context.Background(), validProposerCtx)
	require.NoError(t, err)

	block2VerifyWithCtx, ok := block2.(block.WithVerifyContext)
	require.True(t, ok)
	shouldVerifyWithCtx, err := block2VerifyWithCtx.ShouldVerifyWithContext(context.Background())
	require.NoError(t, err)
	require.True(t, shouldVerifyWithCtx)
	require.NoError(t, block2VerifyWithCtx.VerifyWithContext(context.Background(), validProposerCtx))
	require.Equal(t, choices.Processing, block2.Status())
	require.NoError(t, vm.SetPreference(context.Background(), block2.ID()))

	// Verify the block with another valid context
	require.NoError(t, block2VerifyWithCtx.VerifyWithContext(context.Background(), &block.Context{
		PChainHeight: 11,
	}))
	require.Equal(t, choices.Processing, block2.Status())
=======
	vm.clock.Set(vm.clock.Time().Add(2 * time.Second))
	<-issuer

	block2, err := vm.BuildBlockWithContext(context.Background(), validProposerCtx)
	require.NoError(err)

	block2VerifyWithCtx, ok := block2.(block.WithVerifyContext)
	require.True(ok)
	shouldVerifyWithCtx, err := block2VerifyWithCtx.ShouldVerifyWithContext(context.Background())
	require.NoError(err)
	require.True(shouldVerifyWithCtx)
	require.NoError(block2VerifyWithCtx.VerifyWithContext(context.Background(), validProposerCtx))
	require.Equal(choices.Processing, block2.Status())
	require.NoError(vm.SetPreference(context.Background(), block2.ID()))

	// Verify the block with another valid context
	require.NoError(block2VerifyWithCtx.VerifyWithContext(context.Background(), &block.Context{
		PChainHeight: 11,
	}))
	require.Equal(choices.Processing, block2.Status())
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d

	// Verify the block with a different context and modified ValidatorState so that it should fail verification
	testErr := errors.New("test error")
	vm.ctx.ValidatorState.(*validators.TestState).GetValidatorSetF = func(ctx context.Context, height uint64, subnetID ids.ID) (map[ids.NodeID]*validators.GetValidatorOutput, error) {
		return nil, testErr
	}
<<<<<<< HEAD
	require.ErrorIs(t, block2VerifyWithCtx.VerifyWithContext(context.Background(), &block.Context{
		PChainHeight: 9,
	}), testErr)
	require.Equal(t, choices.Processing, block2.Status())

	// Accept the block after performing multiple VerifyWithContext operations
	require.NoError(t, block2.Accept(context.Background()))

	ethBlock := block2.(*chain.BlockWrapper).Block.(*Block).ethBlock
	verifiedMessageReceipts := vm.blockChain.GetReceiptsByHash(ethBlock.Hash())
	require.Len(t, verifiedMessageReceipts, 1)
	verifiedMessageTxReceipt := verifiedMessageReceipts[0]
	require.Equal(t, types.ReceiptStatusSuccessful, verifiedMessageTxReceipt.Status)
=======
	require.ErrorIs(block2VerifyWithCtx.VerifyWithContext(context.Background(), &block.Context{
		PChainHeight: 9,
	}), testErr)
	require.Equal(choices.Processing, block2.Status())

	// Accept the block after performing multiple VerifyWithContext operations
	require.NoError(block2.Accept(context.Background()))
	vm.blockChain.DrainAcceptorQueue()

	ethBlock := block2.(*chain.BlockWrapper).Block.(*Block).ethBlock
	verifiedMessageReceipts := vm.blockChain.GetReceiptsByHash(ethBlock.Hash())
	require.Len(verifiedMessageReceipts, 1)
	verifiedMessageTxReceipt := verifiedMessageReceipts[0]
	require.Equal(types.ReceiptStatusSuccessful, verifiedMessageTxReceipt.Status)
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
}
