// (c) 2019-2020, Ava Labs, Inc.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2021 The go-ethereum Authors
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

package gasprice

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"slices"

	"github.com/ava-labs/subnet-evm/consensus/misc/eip4844"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/rpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

var (
	errInvalidPercentile     = errors.New("invalid reward percentile")
	errRequestBeyondHead     = errors.New("request beyond head block")
	errBeyondHistoricalLimit = errors.New("request beyond historical limit")
)

const (
	// maxQueryLimit is the max number of requested percentiles.
	maxQueryLimit = 100
)

// txGasAndReward is sorted in ascending order based on reward
type txGasAndReward struct {
	gasUsed uint64
	reward  *big.Int
}

type slimBlock struct {
	GasUsed  uint64
	GasLimit uint64
	BaseFee  *big.Int
	Txs      []txGasAndReward

	BlobGasUsed, ExcessBlobGas *uint64
}

// processBlock prepares a [slimBlock] from a retrieved block and list of
// receipts. This slimmed block can be cached and used for future calls.
func processBlock(block *types.Block, receipts types.Receipts) *slimBlock {
	var sb slimBlock
	if sb.BaseFee = block.BaseFee(); sb.BaseFee == nil {
		sb.BaseFee = new(big.Int)
	}
	sb.GasUsed = block.GasUsed()
	sb.GasLimit = block.GasLimit()
	sorter := make([]txGasAndReward, len(block.Transactions()))
	for i, tx := range block.Transactions() {
		reward, _ := tx.EffectiveGasTip(sb.BaseFee)
		sorter[i] = txGasAndReward{gasUsed: receipts[i].GasUsed, reward: reward}
	}
	slices.SortStableFunc(sorter, func(a, b txGasAndReward) int {
		return a.reward.Cmp(b.reward)
	})
	sb.Txs = sorter
	sb.BlobGasUsed, sb.ExcessBlobGas = block.BlobGasUsed(), block.ExcessBlobGas()
	return &sb
}

// processPercentiles returns baseFee, gasUsedRatio, and optionally reward percentiles (if any are
// requested)
func (sb *slimBlock) processPercentiles(percentiles []float64) ([]*big.Int, *big.Int, float64, *big.Int, float64) {
	// XXX: what happens to the next block's base fee here?
	// Fill in blob base fee and next blob base fee.
	blobBaseFee := new(big.Int)
	if excessBlobGas := sb.ExcessBlobGas; excessBlobGas != nil {
		blobBaseFee = eip4844.CalcBlobFee(*excessBlobGas)
	}

	// Compute gas used ratio for normal and blob gas.
	var (
		gasUsedRatio     = float64(sb.GasUsed) / float64(sb.GasLimit)
		blobGasUsedRatio float64
	)
	if blobGasUsed := sb.BlobGasUsed; blobGasUsed != nil {
		blobGasUsedRatio = float64(*blobGasUsed) / params.MaxBlobGasPerBlock
	}
	if len(percentiles) == 0 {
		// rewards were not requested
		return nil, sb.BaseFee, gasUsedRatio, blobBaseFee, blobGasUsedRatio
	}

	txLen := len(sb.Txs)
	reward := make([]*big.Int, len(percentiles))
	if txLen == 0 {
		// return an all zero row if there are no transactions to gather data from
		for i := range reward {
			reward[i] = new(big.Int)
		}
		return reward, sb.BaseFee, gasUsedRatio, blobBaseFee, blobGasUsedRatio
	}

	// sb transactions are already sorted by tip, so we don't need to re-sort
	var txIndex int
	sumGasUsed := sb.Txs[0].gasUsed
	for i, p := range percentiles {
		thresholdGasUsed := uint64(float64(sb.GasUsed) * p / 100)
		for sumGasUsed < thresholdGasUsed && txIndex < txLen-1 {
			txIndex++
			sumGasUsed += sb.Txs[txIndex].gasUsed
		}
		reward[i] = sb.Txs[txIndex].reward
	}
	return reward, sb.BaseFee, gasUsedRatio, blobBaseFee, blobGasUsedRatio
}

// resolveBlockRange resolves the specified block range to absolute block numbers while also
// enforcing backend specific limitations.
// Note: an error is only returned if retrieving the head header has failed. If there are no
// retrievable blocks in the specified range then zero block count is returned with no error.
func (oracle *Oracle) resolveBlockRange(ctx context.Context, lastBlock rpc.BlockNumber, blocks uint64) (uint64, uint64, error) {
	// Query either pending block or head header and set headBlock
	if lastBlock == rpc.PendingBlockNumber {
		// Pending block not supported by backend, process until latest block
		lastBlock = rpc.LatestBlockNumber
		blocks--
	}
	if blocks == 0 {
		return 0, 0, nil
	}

	lastAcceptedBlock := rpc.BlockNumber(oracle.backend.LastAcceptedBlock().NumberU64())
	maxQueryDepth := rpc.BlockNumber(oracle.maxBlockHistory) - 1
	if lastBlock.IsAccepted() {
		lastBlock = lastAcceptedBlock
	} else if lastAcceptedBlock > maxQueryDepth && lastAcceptedBlock-maxQueryDepth > lastBlock {
		// If the requested last block reaches further back than [oracle.maxBlockHistory] past the last accepted block return an error
		// Note: this allows some blocks past this point to be fetched since it will start fetching [blocks] from this point.
		return 0, 0, fmt.Errorf("%w: requested %d, head %d", errBeyondHistoricalLimit, lastBlock, lastAcceptedBlock)
	} else if lastBlock > lastAcceptedBlock {
		// If the requested block is above the accepted block return an error
		return 0, 0, fmt.Errorf("%w: requested %d, head %d", errRequestBeyondHead, lastBlock, lastAcceptedBlock)
	}
	// Ensure not trying to retrieve before genesis
	if rpc.BlockNumber(blocks) > lastBlock+1 {
		blocks = uint64(lastBlock + 1)
	}
	// Truncate blocks range if extending past [oracle.maxBlockHistory]
	oldestQueriedIndex := lastBlock - rpc.BlockNumber(blocks) + 1
	if queryDepth := lastAcceptedBlock - oldestQueriedIndex; queryDepth > maxQueryDepth {
		overage := uint64(queryDepth - maxQueryDepth)
		blocks -= overage
	}
	// It is not possible that [blocks] could be <= 0 after
	// truncation as the [lastBlock] requested will at least by fetchable.
	// Otherwise, we would've returned an error earlier.
	return uint64(lastBlock), blocks, nil
}

// FeeHistory returns data relevant for fee estimation based on the specified range of blocks.
// The range can be specified either with absolute block numbers or ending with the latest
// or pending block. Backends may or may not support gathering data from the pending block
// or blocks older than a certain age (specified in maxHistory). The first block of the
// actually processed range is returned to avoid ambiguity when parts of the requested range
// are not available or when the head has changed during processing this request.
// Five arrays are returned based on the processed blocks:
//   - reward: the requested percentiles of effective priority fees per gas of transactions in each
//     block, sorted in ascending order and weighted by gas used.
//   - baseFee: base fee per gas in the given block
//   - gasUsedRatio: gasUsed/gasLimit in the given block
//   - blobBaseFee: the blob base fee per gas in the given block
//   - blobGasUsedRatio: blobGasUsed/blobGasLimit in the given block
//
// Note: baseFee and blobBaseFee both include the next block after the newest of the returned range,
// because this value can be derived from the newest block.
func (oracle *Oracle) FeeHistory(ctx context.Context, blocks uint64, unresolvedLastBlock rpc.BlockNumber, rewardPercentiles []float64) (*big.Int, [][]*big.Int, []*big.Int, []float64, []*big.Int, []float64, error) {
	if blocks < 1 {
		return common.Big0, nil, nil, nil, nil, nil, nil // returning with no data and no error means there are no retrievable blocks
	}
	maxFeeHistory := oracle.maxCallBlockHistory // XXX: rename to maxHeaderHistory ?
	if len(rewardPercentiles) != 0 {
		maxFeeHistory = oracle.maxCallBlockHistory
	}
	if len(rewardPercentiles) > maxQueryLimit {
		return common.Big0, nil, nil, nil, nil, nil, fmt.Errorf("%w: over the query limit %d", errInvalidPercentile, maxQueryLimit)
	}
	if blocks > maxFeeHistory {
		log.Warn("Sanitizing fee history length", "requested", blocks, "truncated", maxFeeHistory)
		blocks = maxFeeHistory
	}
	for i, p := range rewardPercentiles {
		if p < 0 || p > 100 {
			return common.Big0, nil, nil, nil, nil, nil, fmt.Errorf("%w: %f", errInvalidPercentile, p)
		}
		if i > 0 && p <= rewardPercentiles[i-1] {
			return common.Big0, nil, nil, nil, nil, nil, fmt.Errorf("%w: #%d:%f >= #%d:%f", errInvalidPercentile, i-1, rewardPercentiles[i-1], i, p)
		}
	}
	lastBlock, blocks, err := oracle.resolveBlockRange(ctx, unresolvedLastBlock, blocks)
	if err != nil || blocks == 0 {
		return common.Big0, nil, nil, nil, nil, nil, err
	}
	oldestBlock := lastBlock + 1 - blocks

	var (
		reward           = make([][]*big.Int, blocks)
		baseFee          = make([]*big.Int, blocks)
		gasUsedRatio     = make([]float64, blocks)
		blobGasUsedRatio = make([]float64, blocks)
		blobBaseFee      = make([]*big.Int, blocks)
		firstMissing     = blocks
	)

	for blockNumber := oldestBlock; blockNumber < oldestBlock+blocks; blockNumber++ {
		// Check if the context has errored
		if err := ctx.Err(); err != nil {
			return common.Big0, nil, nil, nil, nil, nil, err
		}

		i := blockNumber - oldestBlock
		var sb *slimBlock
		if sbCache, ok := oracle.historyCache.Get(blockNumber); ok {
			sb = sbCache
		} else {
			block, err := oracle.backend.BlockByNumber(ctx, rpc.BlockNumber(blockNumber))
			if err != nil {
				return common.Big0, nil, nil, nil, nil, nil, err
			}
			// getting no block and no error means we are requesting into the future (might happen because of a reorg)
			if block == nil {
				if i == 0 {
					return common.Big0, nil, nil, nil, nil, nil, nil
				}
				firstMissing = i
				break
			}
			receipts, err := oracle.backend.GetReceipts(ctx, block.Hash())
			if err != nil {
				return common.Big0, nil, nil, nil, nil, nil, err
			}
			sb = processBlock(block, receipts)
			oracle.historyCache.Add(blockNumber, sb)
		}
		reward[i], baseFee[i], gasUsedRatio[i], blobBaseFee[i], blobGasUsedRatio[i] = sb.processPercentiles(rewardPercentiles)
	}
	if len(rewardPercentiles) != 0 {
		reward = reward[:firstMissing]
	} else {
		reward = nil
	}
	// XXX: why is upstream using [+1] indexing for firstMissing?
	baseFee, gasUsedRatio = baseFee[:firstMissing], gasUsedRatio[:firstMissing]
	blobBaseFee, blobGasUsedRatio = blobBaseFee[:firstMissing], blobGasUsedRatio[:firstMissing]
	return new(big.Int).SetUint64(oldestBlock), reward, baseFee, gasUsedRatio, blobBaseFee, blobGasUsedRatio, nil
}
