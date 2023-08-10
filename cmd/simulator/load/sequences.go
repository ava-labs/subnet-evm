// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package load

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/platformvm"
	"github.com/ava-labs/subnet-evm/cmd/simulator/config"
	"github.com/ava-labs/subnet-evm/cmd/simulator/txs"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/warp"
	"github.com/ethereum/go-ethereum/log"
)

type TxSequenceGetter func(
	ctx context.Context, config config.Config, chainID *big.Int,
	pks []*ecdsa.PrivateKey, client ethclient.Client,
) ([]txs.TxSequence[*types.Transaction], error)

// func GetEVMTxSequences(
// 	ctx context.Context, config config.Config, chainID *big.Int,
// 	pks []*ecdsa.PrivateKey, client ethclient.Client,
//	tracker *awmTimeTracker,
// ) ([]txs.TxSequence[TrackableTx], error) {
// 	bigGwei := big.NewInt(params.GWei)
// 	gasTipCap := new(big.Int).Mul(bigGwei, big.NewInt(config.MaxTipCap))
// 	gasFeeCap := new(big.Int).Mul(bigGwei, big.NewInt(config.MaxFeeCap))
//
// 	// Normal EVM txs
// 	signer := types.LatestSignerForChainID(chainID)
// 	txGenerator := func(key *ecdsa.PrivateKey, nonce uint64) (*types.Transaction, error) {
// 		addr := ethcrypto.PubkeyToAddress(key.PublicKey)
// 		tx, err := types.SignNewTx(key, signer, &types.DynamicFeeTx{
// 			ChainID:   chainID,
// 			Nonce:     nonce,
// 			GasTipCap: gasTipCap,
// 			GasFeeCap: gasFeeCap,
// 			Gas:       params.TxGas,
// 			To:        &addr,
// 			Data:      nil,
// 			Value:     common.Big0,
// 		})
// 		if err != nil {
// 			return nil, err
// 		}
// 		return tx, nil
// 	}
// 	return txs.GenerateTxSequences(ctx, txGenerator, client, pks, config.TxsPerWorker)
// }

func GetWarpSendTxSequences(
	ctx context.Context, config config.Config, chainID *big.Int,
	pks []*ecdsa.PrivateKey, startingNonces []uint64,
) ([]txs.TxSequence[*AwmTx], error) {
	bigGwei := big.NewInt(params.GWei)
	gasTipCap := new(big.Int).Mul(bigGwei, big.NewInt(config.MaxTipCap))
	gasFeeCap := new(big.Int).Mul(bigGwei, big.NewInt(config.MaxFeeCap))

	subnetBStr := config.Subnets[1].BlockchainID.String()
	subnetB, err := ids.FromString(subnetBStr)
	if err != nil {
		return nil, err
	}
	txGenerator := MkSendWarpTxGenerator(chainID, subnetB, gasFeeCap, gasTipCap)
	return txs.GenerateTxSequences(ctx, txGenerator, pks, startingNonces, config.TxsPerWorker)
}

func GetWarpReceiveTxSequences(
	ctx context.Context, config config.Config, chainID *big.Int,
	pks []*ecdsa.PrivateKey, startingNonces []uint64,
) ([]txs.TxSequence[*AwmTx], error) {
	subnetA := config.Subnets[0]
	// We need the validator set of subnet A to determine the index of
	// each validator in the bit set.
	validatorIndexes, err := getValidatorIndexes(ctx, subnetA.ValidatorURIs[0], subnetA.SubnetID)
	if err != nil {
		return nil, err
	}

	ch := make(chan warpSignature) // channel for incoming signatures
	// We will need to aggregate signatures for messages that are sent on
	// subnet A. So we will subscribe to the subnet A's accepted logs.
	endpoints := toWebsocketURIs(subnetA)
	for i, endpoint := range endpoints {
		client, err := ethclient.Dial(endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to dial client at %s: %w", endpoint, err)
		}
		log.Info("Connected to client", "client", endpoint, "idx", i)

		warpClient, err := warp.NewWarpClient(subnetA.ValidatorURIs[i], subnetA.BlockchainID.String())
		if err != nil {
			return nil, err
		}
		// TODO: properly shutdown warp clients
		bitsetIndex, ok := validatorIndexes[subnetA.NodeIDs[i]]
		if !ok {
			return nil, fmt.Errorf("validator %s not found in validator set", subnetA.NodeIDs[i])
		}
		_ = NewWarpRelayClient(ctx, client, warpClient, ch, bitsetIndex)
	}

	threshold := uint64(4) // TODO: should not be hardcoded
	// TODO: should not be hardcoded like this
	expectedMessages := int(config.TxsPerWorker) * config.Workers
	warpRelay := NewWarpRelay(ctx, threshold, ch, expectedMessages)
	// Each worker will listen for signed warp messages that are
	// ready to be issued
	txSequences := make([]txs.TxSequence[*AwmTx], config.Workers)
	for i := 0; i < config.Workers; i++ {
		txSequences[i] = NewWarpRelayTxSequence(ctx, warpRelay.signedMessages, chainID, pks[i], startingNonces[i])
	}
	return txSequences, nil
}

func getValidatorIndexes(ctx context.Context, nodeURI string, subnetID ids.ID) (map[ids.NodeID]int, error) {
	client := platformvm.NewClient(nodeURI)
	vdrs, err := client.GetCurrentValidators(ctx, subnetID, nil)
	if err != nil {
		return nil, err
	}
	log.Info("Got validator set", "numValidators", len(vdrs), "subnetID", subnetID)

	indexMap := make(map[ids.NodeID]int, len(vdrs))
	for i, vdr := range vdrs {
		indexMap[vdr.NodeID] = i
		log.Info("Validator", "nodeID", vdr.NodeID, "index", i)
	}

	return indexMap, nil
}
