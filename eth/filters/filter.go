// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2014 The go-ethereum Authors
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

package filters

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/types"
	"github.com/ava-labs/subnet-evm/core/bloombits"
	"github.com/ava-labs/subnet-evm/plugin/evm/customlogs"
	"github.com/ava-labs/subnet-evm/rpc"
)

// Filter can be used to retrieve and filter logs.
type Filter struct {
	sys *FilterSystem

	addresses []common.Address
	topics    [][]common.Hash

	block      *common.Hash // Block hash if filtering a single block
	begin, end int64        // Range interval if filtering multiple blocks

	matcher *bloombits.Matcher
}

// NewRangeFilter creates a new filter which uses a bloom filter on blocks to
// figure out whether a particular block is interesting or not.
func (sys *FilterSystem) NewRangeFilter(begin, end int64, addresses []common.Address, topics [][]common.Hash) *Filter {
	// Flatten the address and topic filter clauses into a single bloombits filter
	// system. Since the bloombits are not positional, nil topics are permitted,
	// which get flattened into a nil byte slice.
	var filters [][][]byte
	if len(addresses) > 0 {
		filter := make([][]byte, len(addresses))
		for i, address := range addresses {
			filter[i] = address.Bytes()
		}
		filters = append(filters, filter)
	}
	for _, topicList := range topics {
		filter := make([][]byte, len(topicList))
		for i, topic := range topicList {
			filter[i] = topic.Bytes()
		}
		filters = append(filters, filter)
	}
	size, _ := sys.backend.BloomStatus()

	// Create a generic filter and convert it into a range filter
	filter := newFilter(sys, addresses, topics)

	filter.matcher = bloombits.NewMatcher(size, filters)
	filter.begin = begin
	filter.end = end

	return filter
}

// NewBlockFilter creates a new filter which directly inspects the contents of
// a block to figure out whether it is interesting or not.
func (sys *FilterSystem) NewBlockFilter(block common.Hash, addresses []common.Address, topics [][]common.Hash) *Filter {
	// Create a generic filter and convert it into a block filter
	filter := newFilter(sys, addresses, topics)
	filter.block = &block
	return filter
}

// newFilter creates a generic filter that can either filter based on a block hash,
// or based on range queries. The search criteria needs to be explicitly set.
func newFilter(sys *FilterSystem, addresses []common.Address, topics [][]common.Hash) *Filter {
	return &Filter{
		sys:       sys,
		addresses: addresses,
		topics:    topics,
	}
}

// Logs searches the blockchain for matching log entries, returning all from the
// first block that contains matches, updating the start of the filter accordingly.
func (f *Filter) Logs(ctx context.Context) ([]*types.Log, error) {
	// If we're doing singleton block filtering, execute and return
	if f.block != nil {
		header, err := f.sys.backend.HeaderByHash(ctx, *f.block)
		if err != nil {
			return nil, err
		}
		if header == nil {
			return nil, errors.New("unknown block")
		}
		return f.blockLogs(ctx, header)
	}

	// Disallow blocks past the last accepted block if the backend does not
	// allow unfinalized queries.
	allowUnfinalizedQueries := f.sys.backend.IsAllowUnfinalizedQueries()
	acceptedBlock := f.sys.backend.LastAcceptedBlock()
	if !allowUnfinalizedQueries && acceptedBlock != nil {
		lastAccepted := acceptedBlock.Number().Int64()
		if f.begin >= 0 && f.begin > lastAccepted {
			return nil, fmt.Errorf("requested from block %d after last accepted block %d", f.begin, lastAccepted)
		}
		if f.end >= 0 && f.end > lastAccepted {
			return nil, fmt.Errorf("requested to block %d after last accepted block %d", f.end, lastAccepted)
		}
	}

	var (
		beginPending = f.begin == rpc.PendingBlockNumber.Int64()
		endPending   = f.end == rpc.PendingBlockNumber.Int64()
		endSet       = f.end >= 0
	)

	// special case for pending logs
	if beginPending && !endPending {
		return nil, errInvalidBlockRange
	}

	// Short-cut if all we care about is pending logs
	if beginPending && endPending {
		return nil, nil
	}

	resolveSpecial := func(number int64) (int64, error) {
		var hdr *types.Header
		switch number {
		case rpc.LatestBlockNumber.Int64(), rpc.PendingBlockNumber.Int64():
			// we should return head here since we've already captured
			// that we need to get the pending logs in the pending boolean above
			hdr, _ = f.sys.backend.HeaderByNumber(ctx, rpc.LatestBlockNumber)
			if hdr == nil {
				return 0, errors.New("latest header not found")
			}
		case rpc.FinalizedBlockNumber.Int64():
			hdr, _ = f.sys.backend.HeaderByNumber(ctx, rpc.FinalizedBlockNumber)
			if hdr == nil {
				return 0, errors.New("finalized header not found")
			}
		case rpc.SafeBlockNumber.Int64():
			hdr, _ = f.sys.backend.HeaderByNumber(ctx, rpc.SafeBlockNumber)
			if hdr == nil {
				return 0, errors.New("safe header not found")
			}
		default:
			return number, nil
		}
		return hdr.Number.Int64(), nil
	}

	var err error
	// range query need to resolve the special begin/end block number
	if f.begin, err = resolveSpecial(f.begin); err != nil {
		return nil, err
	}
	if f.end, err = resolveSpecial(f.end); err != nil {
		return nil, err
	}

	// When querying unfinalized data without a populated end block, it is
	// possible that the begin will be greater than the end.
	//
	// We error in this case to prevent a bad UX where the caller thinks there
	// are no logs from the specified beginning to end (when in reality there may
	// be some).
	if endSet && f.end < f.begin {
		return nil, fmt.Errorf("begin block %d is greater than end block %d", f.begin, f.end)
	}

	// If the requested range of blocks exceeds the maximum number of blocks allowed by the backend
	// return an error instead of searching for the logs.
	if maxBlocks := f.sys.backend.GetMaxBlocksPerRequest(); f.end-f.begin >= maxBlocks && maxBlocks > 0 {
		return nil, fmt.Errorf("requested too many blocks from %d to %d, maximum is set to %d", f.begin, f.end, maxBlocks)
	}
	// Gather all indexed logs, and finish with non indexed ones
	logChan, errChan := f.rangeLogsAsync(ctx)
	var logs []*types.Log
	for {
		select {
		case log := <-logChan:
			logs = append(logs, log)
		case err := <-errChan:
			if err != nil {
				// if an error occurs during extraction, we do return the extracted data
				return logs, err
			}
			return logs, nil
		}
	}
}

// rangeLogsAsync retrieves block-range logs that match the filter criteria asynchronously,
// it creates and returns two channels: one for delivering log data, and one for reporting errors.
func (f *Filter) rangeLogsAsync(ctx context.Context) (chan *types.Log, chan error) {
	var (
		logChan = make(chan *types.Log)
		errChan = make(chan error)
	)

	go func() {
		defer func() {
			close(errChan)
			close(logChan)
		}()

		// Gather all indexed logs, and finish with non indexed ones
		var (
			end            = uint64(f.end)
			size, sections = f.sys.backend.BloomStatus()
			err            error
		)
		if indexed := sections * size; indexed > uint64(f.begin) {
			if indexed > end {
				indexed = end + 1
			}
			if err = f.indexedLogs(ctx, indexed-1, logChan); err != nil {
				errChan <- err
				return
			}
		}

		if err := f.unindexedLogs(ctx, end, logChan); err != nil {
			errChan <- err
			return
		}

		errChan <- nil
	}()

	return logChan, errChan
}

// indexedLogs returns the logs matching the filter criteria based on the bloom
// bits indexed available locally or via the network.
func (f *Filter) indexedLogs(ctx context.Context, end uint64, logChan chan *types.Log) error {
	// Create a matcher session and request servicing from the backend
	matches := make(chan uint64, 64)

	session, err := f.matcher.Start(ctx, uint64(f.begin), end, matches)
	if err != nil {
		return err
	}
	defer session.Close()

	f.sys.backend.ServiceFilter(ctx, session)

	for {
		select {
		case number, ok := <-matches:
			// Abort if all matches have been fulfilled
			if !ok {
				err := session.Error()
				if err == nil {
					f.begin = int64(end) + 1
				}
				return err
			}
			f.begin = int64(number) + 1

			// Retrieve the suggested block and pull any truly matching logs
			header, err := f.sys.backend.HeaderByNumber(ctx, rpc.BlockNumber(number))
			if header == nil || err != nil {
				return err
			}
			found, err := f.checkMatches(ctx, header)
			if err != nil {
				return err
			}
			for _, log := range found {
				logChan <- log
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// unindexedLogs returns the logs matching the filter criteria based on raw block
// iteration and bloom matching.
func (f *Filter) unindexedLogs(ctx context.Context, end uint64, logChan chan *types.Log) error {
	for ; f.begin <= int64(end); f.begin++ {
		header, err := f.sys.backend.HeaderByNumber(ctx, rpc.BlockNumber(f.begin))
		if header == nil || err != nil {
			return err
		}
		found, err := f.blockLogs(ctx, header)
		if err != nil {
			return err
		}
		for _, log := range found {
			select {
			case logChan <- log:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	return nil
}

// blockLogs returns the logs matching the filter criteria within a single block.
func (f *Filter) blockLogs(ctx context.Context, header *types.Header) ([]*types.Log, error) {
	if bloomFilter(header.Bloom, f.addresses, f.topics) {
		return f.checkMatches(ctx, header)
	}
	return nil, nil
}

// checkMatches checks if the receipts belonging to the given header contain any log events that
// match the filter criteria. This function is called when the bloom filter signals a potential match.
func (f *Filter) checkMatches(ctx context.Context, header *types.Header) ([]*types.Log, error) {
	logsList, err := f.sys.getLogs(ctx, header.Hash(), header.Number.Uint64())
	if err != nil {
		return nil, err
	}

	unfiltered := customlogs.FlattenLogs(logsList)
	logs := filterLogs(unfiltered, nil, nil, f.addresses, f.topics)
	if len(logs) == 0 {
		return nil, nil
	}
	// Most backends will deliver un-derived logs, but check nevertheless.
	if len(logs) > 0 && logs[0].TxHash != (common.Hash{}) {
		return logs, nil
	}
	// We have matching logs, check if we need to resolve full logs via the light client
	receipts, err := f.sys.backend.GetReceipts(ctx, header.Hash())
	if err != nil {
		return nil, err
	}
	unfiltered = unfiltered[:0]
	for _, receipt := range receipts {
		unfiltered = append(unfiltered, receipt.Logs...)
	}
	logs = filterLogs(unfiltered, nil, nil, f.addresses, f.topics)

	return logs, nil
}

// includes returns true if the element is present in the list.
func includes[T comparable](things []T, element T) bool {
	for _, thing := range things {
		if thing == element {
			return true
		}
	}
	return false
}

// filterLogs creates a slice of logs matching the given criteria.
func filterLogs(logs []*types.Log, fromBlock, toBlock *big.Int, addresses []common.Address, topics [][]common.Hash) []*types.Log {
	check := func(log *types.Log) bool {
		if fromBlock != nil && fromBlock.Int64() >= 0 && fromBlock.Uint64() > log.BlockNumber {
			return false
		}
		if toBlock != nil && toBlock.Int64() >= 0 && toBlock.Uint64() < log.BlockNumber {
			return false
		}
		if len(addresses) > 0 && !includes(addresses, log.Address) {
			return false
		}
		// If the to filtered topics is greater than the amount of topics in logs, skip.
		if len(topics) > len(log.Topics) {
			return false
		}
		for i, sub := range topics {
			if len(sub) == 0 {
				continue // empty rule set == wildcard
			}
			if !includes(sub, log.Topics[i]) {
				return false
			}
		}
		return true
	}
	var ret []*types.Log
	for _, log := range logs {
		if check(log) {
			ret = append(ret, log)
		}
	}
	return ret
}

func bloomFilter(bloom types.Bloom, addresses []common.Address, topics [][]common.Hash) bool {
	if len(addresses) > 0 {
		var included bool
		for _, addr := range addresses {
			if types.BloomLookup(bloom, addr) {
				included = true
				break
			}
		}
		if !included {
			return false
		}
	}

	for _, sub := range topics {
		included := len(sub) == 0 // empty rule set == wildcard
		for _, topic := range sub {
			if types.BloomLookup(bloom, topic) {
				included = true
				break
			}
		}
		if !included {
			return false
		}
	}
	return true
}
