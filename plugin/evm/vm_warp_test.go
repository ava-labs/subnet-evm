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
	subnetEVMUtils "github.com/ava-labs/subnet-evm/utils"
	predicateutils "github.com/ava-labs/subnet-evm/utils/predicate"
	warpPayload "github.com/ava-labs/subnet-evm/warp/payload"
	"github.com/ava-labs/subnet-evm/x/warp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

func TestSendWarpMessage(t *testing.T) {
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

	acceptedLogsChan := make(chan []*types.Log, 10)
	logsSub := vm.eth.APIBackend.SubscribeAcceptedLogsEvent(acceptedLogsChan)
	defer logsSub.Unsubscribe()

	payload := utils.RandomBytes(100)

	warpSendMessageInput, err := warp.PackSendWarpMessage(warp.SendWarpMessageInput{
		DestinationChainID: common.Hash(vm.ctx.CChainID),
		DestinationAddress: testEthAddrs[1],
		Payload:            payload,
	})
	require.NoError(err)

	// Submit a transaction to trigger sending a warp message
	tx0 := types.NewTransaction(uint64(0), warp.ContractAddress, big.NewInt(1), 100_000, big.NewInt(testMinGasPrice), warpSendMessageInput)
	signedTx0, err := types.SignTx(tx0, types.LatestSignerForChainID(vm.chainConfig.ChainID), testKeys[0])
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
	unsignedMessageID := unsignedMessage.ID()

	// Verify the signature cannot be fetched before the block is accepted
	_, err = vm.warpBackend.GetSignature(unsignedMessageID)
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

	acceptedLogsChan := make(chan []*types.Log, 10)
	logsSub := vm.eth.APIBackend.SubscribeAcceptedLogsEvent(acceptedLogsChan)
	defer logsSub.Unsubscribe()

	payload := utils.RandomBytes(100)

	addressedPayload, err := warpPayload.NewAddressedPayload(
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
	blsPublicKey1 := bls.PublicFromSecretKey(blsSecretKey1)
	blsSignature1 := bls.Sign(blsSecretKey1, unsignedMessage.Bytes())

	nodeID2 := ids.GenerateTestNodeID()
	blsSecretKey2, err := bls.NewSecretKey()
	require.NoError(err)
	blsPublicKey2 := bls.PublicFromSecretKey(blsSecretKey2)
	blsSignature2 := bls.Sign(blsSecretKey2, unsignedMessage.Bytes())

	blsAggregatedSignature, err := bls.AggregateSignatures([]*bls.Signature{blsSignature1, blsSignature2})
	require.NoError(err)

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
	require.NoError(err)

	getWarpMsgInput, err := warp.PackGetVerifiedWarpMessage()
	require.NoError(err)
	getVerifiedWarpMessageTx, err := types.SignTx(
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
			Payload:             payload,
		},
		Exists: true,
	})
	require.NoError(err)

	// Assert that DoCall returns the expected output
	hexGetWarpMsgInput := hexutil.Bytes(getVerifiedWarpMessageTx.Data())
	hexGasLimit := hexutil.Uint64(getVerifiedWarpMessageTx.Gas())
	accessList := getVerifiedWarpMessageTx.AccessList()
	blockNum := new(rpc.BlockNumber)
	*blockNum = rpc.LatestBlockNumber

	executionRes, err := ethapi.DoCall(
		context.Background(),
		vm.eth.APIBackend,
		ethapi.TransactionArgs{
			To:         getVerifiedWarpMessageTx.To(),
			Input:      &hexGetWarpMsgInput,
			AccessList: &accessList,
			Gas:        &hexGasLimit,
		},
		rpc.BlockNumberOrHash{BlockNumber: blockNum},
		nil,
		time.Second,
		10_000_000,
	)
	require.NoError(err)
	require.NoError(executionRes.Err)
	require.Equal(expectedOutput, executionRes.ReturnData)

	// Build, verify, and accept block with valid proposer context.
	validProposerCtx := &block.Context{
		PChainHeight: 10,
	}
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

	// Verify the block with a different context and modified ValidatorState so that it should fail verification
	testErr := errors.New("test error")
	vm.ctx.ValidatorState.(*validators.TestState).GetValidatorSetF = func(ctx context.Context, height uint64, subnetID ids.ID) (map[ids.NodeID]*validators.GetValidatorOutput, error) {
		return nil, testErr
	}
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
}
