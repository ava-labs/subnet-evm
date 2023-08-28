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
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/exp/maps"
)

func GetWarpSendTxSequences(
	ctx context.Context, config config.Config, chainID *big.Int,
	pks []*ecdsa.PrivateKey, startingNonces []uint64,
) ([]txs.TxSequence[*types.Transaction], error) {
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
	signedMessages chan *pwarp.Message,
) ([]txs.TxSequence[*types.Transaction], error) {
	// Each worker will listen for signed warp messages that are
	// ready to be issued
	txSequences := make([]txs.TxSequence[*types.Transaction], config.Workers)
	for i := 0; i < config.Workers; i++ {
		txSequences[i] = NewWarpRelayTxSequence(ctx, signedMessages, chainID, pks[i], startingNonces[i])
	}
	return txSequences, nil
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
	vdrSet, err := client.GetValidatorsAt(ctx, subnetID, height)
	if err != nil {
		return nil, 0, err
	}
	// TODO: use the factored out code in avalanchego when the new version is
	// released.
	var (
		vdrs        = make(map[string]*pwarp.Validator, len(vdrSet))
		totalWeight uint64
	)
	for _, vdr := range vdrSet {
		totalWeight, err = math.Add64(totalWeight, vdr.Weight)
		if err != nil {
			return nil, 0, fmt.Errorf("%w: %v", pwarp.ErrWeightOverflow, err)
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
