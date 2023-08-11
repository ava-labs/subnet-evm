// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package load

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils"
	"github.com/ava-labs/avalanchego/utils/math"
	"github.com/ava-labs/avalanchego/vms/platformvm"
	pwarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/cmd/simulator/config"
	"github.com/ava-labs/subnet-evm/cmd/simulator/txs"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/exp/maps"
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

	threshold := uint64(4) // TODO: should not be hardcoded
	// TODO: should not be hardcoded like this
	expectedMessages := int(config.TxsPerWorker) * config.Workers
	warpRelay, err := NewWarpRelay(ctx, subnetA, threshold, expectedMessages)
	if err != nil {
		return nil, err
	}
	go func() {
		err := warpRelay.Run(ctx)
		if err != nil {
			log.Error("warp relay failed", "err", err)
		}
	}()

	// Each worker will listen for signed warp messages that are
	// ready to be issued
	txSequences := make([]txs.TxSequence[*AwmTx], config.Workers)
	for i := 0; i < config.Workers; i++ {
		txSequences[i] = NewWarpRelayTxSequence(ctx, warpRelay.signedMessages, chainID, pks[i], startingNonces[i])
	}
	return txSequences, nil
}

type validatorInfo map[ids.NodeID]int // nodeID -> index in bls validator set

func getValidatorIndexes(ctx context.Context, nodeURI string, subnetID ids.ID) (validatorInfo, error) {
	client := platformvm.NewClient(nodeURI)
	height, err := client.GetHeight(ctx)
	if err != nil {
		return nil, err
	}
	vdrSet, err := client.GetValidatorsAt(ctx, subnetID, height)
	if err != nil {
		return nil, err
	}
	// TODO: should factor this code out in avalanchego
	var (
		vdrs        = make(map[string]*pwarp.Validator, len(vdrSet))
		totalWeight uint64
	)
	for _, vdr := range vdrSet {
		totalWeight, err = math.Add64(totalWeight, vdr.Weight)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", pwarp.ErrWeightOverflow, err)
		}

		if vdr.PublicKey == nil {
			continue
		}

		pkBytes := vdr.PublicKey.Serialize()
		uniqueVdr, ok := vdrs[string(pkBytes)]
		if !ok {
			uniqueVdr = &pwarp.Validator{
				PublicKey:      vdr.PublicKey,
				PublicKeyBytes: pkBytes,
			}
			vdrs[string(pkBytes)] = uniqueVdr
		}

		uniqueVdr.Weight += vdr.Weight // Impossible to overflow here
		uniqueVdr.NodeIDs = append(uniqueVdr.NodeIDs, vdr.NodeID)
	}

	// Sort validators by public key
	vdrList := maps.Values(vdrs)
	utils.Sort(vdrList)
	log.Info("Got validator set", "numValidators", len(vdrs), "subnetID", subnetID)

	indexMap := make(map[ids.NodeID]int, len(vdrSet))
	for i, vdr := range vdrList {
		for _, nodeID := range vdr.NodeIDs {
			indexMap[nodeID] = i
			log.Info(
				"validator bls info",
				"nodeID", nodeID,
				"index", i,
				"weight", vdr.Weight,
				"pk", common.Bytes2Hex(vdr.PublicKeyBytes[0:5]),
			)
		}
	}
	return indexMap, nil
}
