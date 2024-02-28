// Copyright (C) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/subnet-evm/cmd/simulator/key"
	"github.com/ava-labs/subnet-evm/cmd/simulator/load"
	"github.com/ava-labs/subnet-evm/cmd/simulator/metrics"
	"github.com/ava-labs/subnet-evm/cmd/simulator/txs"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/interfaces"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contracts/warp"
	"github.com/ava-labs/subnet-evm/predicate"
	warpBackend "github.com/ava-labs/subnet-evm/warp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

func ExecuteWarpLoadTest(
	ctx context.Context,
	sendingSubnetURI string,
	sendingSubnetID ids.ID,
	sendingSubnetBlockchainID ids.ID,
	receivingSubnetID ids.ID,
	numWorkers int,
	txsPerWorker uint64,
	batchSize uint64,
	fundedKey *ecdsa.PrivateKey,
	sendingClients []ethclient.Client,
	receivingClients []ethclient.Client,
) error {
	keys := make([]*key.Key, 0, numWorkers)
	privateKeys := make([]*ecdsa.PrivateKey, 0, numWorkers)
	prefundedKey := key.CreateKey(fundedKey)
	keys = append(keys, prefundedKey)
	for i := 1; i < numWorkers; i++ {
		newKey, err := key.Generate()
		if err != nil {
			return err
		}
		keys = append(keys, newKey)
		privateKeys = append(privateKeys, newKey.PrivKey)
	}

	loadMetrics := metrics.NewDefaultMetrics()

	log.Info("Distributing funds on sending subnet", "numKeys", len(keys))
	keys, err := load.DistributeFunds(ctx, sendingClients[0], keys, len(keys), new(big.Int).Mul(big.NewInt(100), big.NewInt(params.Ether)), loadMetrics)
	if err != nil {
		return err
	}

	log.Info("Distributing funds on receiving subnet", "numKeys", len(keys))
	_, err = load.DistributeFunds(ctx, receivingClients[0], keys, len(keys), new(big.Int).Mul(big.NewInt(100), big.NewInt(params.Ether)), loadMetrics)
	if err != nil {
		return err
	}

	log.Info("Creating workers for each subnet...")
	chainAWorkers := make([]txs.Worker[*types.Transaction], 0, len(keys))
	for i := range keys {
		chainAWorkers = append(chainAWorkers, load.NewTxReceiptWorker(ctx, sendingClients[i]))
	}
	chainBWorkers := make([]txs.Worker[*types.Transaction], 0, len(keys))
	for i := range keys {
		chainBWorkers = append(chainBWorkers, load.NewTxReceiptWorker(ctx, receivingClients[i]))
	}

	log.Info("Subscribing to warp send events on sending subnet")
	logs := make(chan types.Log, numWorkers*int(txsPerWorker))
	sub, err := sendingClients[0].SubscribeFilterLogs(ctx, interfaces.FilterQuery{
		Addresses: []common.Address{warp.Module.Address},
	}, logs)
	if err != nil {
		return err
	}
	defer func() {
		sub.Unsubscribe()
		err := <-sub.Err()
		if err != nil {
			log.Error("failed to unscubscribe filter logs", "err", err)
		}
	}()

	sendingSubnetChainID, err := sendingClients[0].ChainID(ctx)
	if err != nil {
		return err
	}
	sendingSubnetSigner := types.LatestSignerForChainID(sendingSubnetChainID)

	log.Info("Generating tx sequence to send warp messages...")
	warpSendSequences, err := txs.GenerateTxSequences(ctx, func(key *ecdsa.PrivateKey, nonce uint64) (*types.Transaction, error) {
		data, err := warp.PackSendWarpMessage([]byte(fmt.Sprintf("Jets %d-%d Dolphins", key.X.Int64(), nonce)))
		if err != nil {
			return nil, err
		}
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   sendingSubnetChainID,
			Nonce:     nonce,
			To:        &warp.Module.Address,
			Gas:       200_000,
			GasFeeCap: big.NewInt(225 * params.GWei),
			GasTipCap: big.NewInt(params.GWei),
			Value:     common.Big0,
			Data:      data,
		})
		return types.SignTx(tx, sendingSubnetSigner, key)
	}, sendingClients[0], privateKeys, txsPerWorker, false)
	if err != nil {
		return err
	}
	log.Info("Executing warp send loader...")
	warpSendLoader := load.New(chainAWorkers, warpSendSequences, batchSize, loadMetrics)
	// TODO: execute send and receive loaders concurrently.
	if err := warpSendLoader.Execute(ctx); err != nil {
		return err
	}
	if err := warpSendLoader.ConfirmReachedTip(ctx); err != nil {
		return err
	}

	warpClient, err := warpBackend.NewClient(sendingSubnetURI, sendingSubnetBlockchainID.String())
	if err != nil {
		return err
	}
	subnetIDStr := ""
	if sendingSubnetID == constants.PrimaryNetworkID {
		subnetIDStr = receivingSubnetID.String()
	}

	receivingSubnetChainID, err := receivingClients[0].ChainID(ctx)
	if err != nil {
		return err
	}
	receivingSubnetSigner := types.LatestSignerForChainID(receivingSubnetChainID)

	log.Info("Executing warp delivery sequences...")
	warpDeliverSequences, err := txs.GenerateTxSequences(ctx, func(key *ecdsa.PrivateKey, nonce uint64) (*types.Transaction, error) {
		// Wait for the next warp send log
		warpLog := <-logs

		unsignedMessage, err := warp.UnpackSendWarpEventDataToMessage(warpLog.Data)
		if err != nil {
			return nil, err
		}
		log.Info("Fetching addressed call aggregate signature via p2p API")

		signedWarpMessageBytes, err := warpClient.GetMessageAggregateSignature(ctx, unsignedMessage.ID(), warp.WarpDefaultQuorumNumerator, subnetIDStr)
		if err != nil {
			return nil, err
		}

		packedInput, err := warp.PackGetVerifiedWarpMessage(0)
		if err != nil {
			return nil, err
		}
		tx := predicate.NewPredicateTx(
			receivingSubnetChainID,
			nonce,
			&warp.Module.Address,
			5_000_000,
			big.NewInt(225*params.GWei),
			big.NewInt(params.GWei),
			common.Big0,
			packedInput,
			types.AccessList{},
			warp.ContractAddress,
			signedWarpMessageBytes,
		)
		return types.SignTx(tx, receivingSubnetSigner, key)
	}, receivingClients[0], privateKeys, txsPerWorker, true)
	if err != nil {
		return err
	}
	log.Info("Executing warp delivery...")
	warpDeliverLoader := load.New(chainBWorkers, warpDeliverSequences, batchSize, loadMetrics)
	if err := warpDeliverLoader.Execute(ctx); err != nil {
		return err
	}
	if err := warpSendLoader.ConfirmReachedTip(ctx); err != nil {
		return err
	}
	log.Info("Completed warp delivery successfully.")
	return nil
}
