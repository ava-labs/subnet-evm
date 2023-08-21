// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package load

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/validators"
	"github.com/ava-labs/avalanchego/vms/platformvm"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/cmd/simulator/config"
	"github.com/ava-labs/subnet-evm/cmd/simulator/txs"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ethereum/go-ethereum/log"
)

var _ warp.ValidatorState = (*validatorState)(nil)

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
	signedMessages chan *warp.Message,
) ([]txs.TxSequence[*AwmTx], error) {
	// Each worker will listen for signed warp messages that are
	// ready to be issued
	txSequences := make([]txs.TxSequence[*AwmTx], config.Workers)
	for i := 0; i < config.Workers; i++ {
		txSequences[i] = NewWarpRelayTxSequence(ctx, signedMessages, chainID, pks[i], startingNonces[i])
	}
	return txSequences, nil
}

type validatorState struct {
	client platformvm.Client
}

func (v *validatorState) GetValidatorSet(
	ctx context.Context, height uint64, subnetID ids.ID,
) (map[ids.NodeID]*validators.GetValidatorOutput, error) {
	return v.client.GetValidatorsAt(ctx, subnetID, height)
}

type validatorInfo map[ids.NodeID]validator

type validator struct {
	index  int
	weight uint64
}

// getValidatorInfo returns a map of nodeID to validator index and the total
// weight of the validator set
func getValidatorInfo(
	ctx context.Context, nodeURI string, subnetID ids.ID,
) (validatorInfo, uint64, error) {
	client := platformvm.NewClient(nodeURI)
	height, err := client.GetHeight(ctx)
	if err != nil {
		return nil, 0, err
	}
	vdrList, totalWeight, err := warp.GetCanonicalValidatorSet(
		ctx, &validatorState{client: client}, height, subnetID)
	if err != nil {
		return nil, 0, err
	}
	log.Info(
		"Got canonical validator set",
		"numValidators", len(vdrList),
		"subnetID", subnetID,
	)

	indexMap := make(validatorInfo, len(vdrList))
	for i, vdr := range vdrList {
		for _, nodeID := range vdr.NodeIDs {
			indexMap[nodeID] = validator{
				index:  i,
				weight: vdr.Weight,
			}
			log.Info(
				"validator bls info",
				"nodeID", nodeID,
				"index", i,
				"weight", vdr.Weight,
			)
			// In case of duplicate BLS keys, the caller will use the first
			// nodeID on the list with the weight of all the nodes on the list.
			break
		}
	}
	return indexMap, totalWeight, nil
}
