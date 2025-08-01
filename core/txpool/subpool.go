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
// Copyright 2023 The go-ethereum Authors
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

package txpool

import (
	"math/big"
	"time"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/types"
	"github.com/ava-labs/libevm/event"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/holiman/uint256"
)

// LazyTransaction contains a small subset of the transaction properties that is
// enough for the miner and other APIs to handle large batches of transactions;
// and supports pulling up the entire transaction when really needed.
type LazyTransaction struct {
	Pool LazyResolver       // Transaction resolver to pull the real transaction up
	Hash common.Hash        // Transaction hash to pull up if needed
	Tx   *types.Transaction // Transaction if already resolved

	Time      time.Time    // Time when the transaction was first seen
	GasFeeCap *uint256.Int // Maximum fee per gas the transaction may consume
	GasTipCap *uint256.Int // Maximum miner tip per gas the transaction can pay

	Gas     uint64 // Amount of gas required by the transaction
	BlobGas uint64 // Amount of blob gas required by the transaction
}

// Resolve retrieves the full transaction belonging to a lazy handle if it is still
// maintained by the transaction pool.
//
// Note, the method will *not* cache the retrieved transaction if the original
// pool has not cached it. The idea being, that if the tx was too big to insert
// originally, silently saving it will cause more trouble down the line (and
// indeed seems to have caused a memory bloat in the original implementation
// which did just that).
func (ltx *LazyTransaction) Resolve() *types.Transaction {
	if ltx.Tx != nil {
		return ltx.Tx
	}
	return ltx.Pool.Get(ltx.Hash)
}

// LazyResolver is a minimal interface needed for a transaction pool to satisfy
// resolving lazy transactions. It's mostly a helper to avoid the entire sub-
// pool being injected into the lazy transaction.
type LazyResolver interface {
	// Get returns a transaction if it is contained in the pool, or nil otherwise.
	Get(hash common.Hash) *types.Transaction
}

// AddressReserver is passed by the main transaction pool to subpools, so they
// may request (and relinquish) exclusive access to certain addresses.
type AddressReserver func(addr common.Address, reserve bool) error

// PendingFilter is a collection of filter rules to allow retrieving a subset
// of transactions for announcement or mining.
//
// Note, the entries here are not arbitrary useful filters, rather each one has
// a very specific call site in mind and each one can be evaluated very cheaply
// by the pool implementations. Only add new ones that satisfy those constraints.
type PendingFilter struct {
	MinTip  *uint256.Int // Minimum miner tip required to include a transaction
	BaseFee *uint256.Int // Minimum 1559 basefee needed to include a transaction
	BlobFee *uint256.Int // Minimum 4844 blobfee needed to include a blob transaction

	OnlyPlainTxs bool // Return only plain EVM transactions (peer-join announces, block space filling)
	OnlyBlobTxs  bool // Return only blob transactions (block blob-space filling)
}

// SubPool represents a specialized transaction pool that lives on its own (e.g.
// blob pool). Since independent of how many specialized pools we have, they do
// need to be updated in lockstep and assemble into one coherent view for block
// production, this interface defines the common methods that allow the primary
// transaction pool to manage the subpools.
type SubPool interface {
	// Filter is a selector used to decide whether a transaction would be added
	// to this particular subpool.
	Filter(tx *types.Transaction) bool

	// Init sets the base parameters of the subpool, allowing it to load any saved
	// transactions from disk and also permitting internal maintenance routines to
	// start up.
	//
	// These should not be passed as a constructor argument - nor should the pools
	// start by themselves - in order to keep multiple subpools in lockstep with
	// one another.
	Init(gasTip uint64, head *types.Header, reserve AddressReserver) error

	// Close terminates any background processing threads and releases any held
	// resources.
	Close() error

	// Reset retrieves the current state of the blockchain and ensures the content
	// of the transaction pool is valid with regard to the chain state.
	Reset(oldHead, newHead *types.Header)

	// SetGasTip updates the minimum price required by the subpool for a new
	// transaction, and drops all transactions below this threshold.
	SetGasTip(tip *big.Int)
	SetMinFee(fee *big.Int)

	// Has returns an indicator whether subpool has a transaction cached with the
	// given hash.
	Has(hash common.Hash) bool
	HasLocal(hash common.Hash) bool

	// Get returns a transaction if it is contained in the pool, or nil otherwise.
	Get(hash common.Hash) *types.Transaction

	// Add enqueues a batch of transactions into the pool if they are valid. Due
	// to the large transaction churn, add may postpone fully integrating the tx
	// to a later point to batch multiple ones together.
	Add(txs []*types.Transaction, local bool, sync bool) []error

	// Pending retrieves all currently processable transactions, grouped by origin
	// account and sorted by nonce.
	//
	// The transactions can also be pre-filtered by the dynamic fee components to
	// reduce allocations and load on downstream subsystems.
	Pending(filter PendingFilter) map[common.Address][]*LazyTransaction
	IteratePending(f func(tx *types.Transaction) bool) bool // Returns false if iteration was interrupted.

	// SubscribeTransactions subscribes to new transaction events. The subscriber
	// can decide whether to receive notifications only for newly seen transactions
	// or also for reorged out ones.
	SubscribeTransactions(ch chan<- core.NewTxsEvent, reorgs bool) event.Subscription

	// Nonce returns the next nonce of an account, with all transactions executable
	// by the pool already applied on top.
	Nonce(addr common.Address) uint64

	// Stats retrieves the current pool stats, namely the number of pending and the
	// number of queued (non-executable) transactions.
	Stats() (int, int)

	// Content retrieves the data content of the transaction pool, returning all the
	// pending as well as queued transactions, grouped by account and sorted by nonce.
	Content() (map[common.Address][]*types.Transaction, map[common.Address][]*types.Transaction)

	// ContentFrom retrieves the data content of the transaction pool, returning the
	// pending as well as queued transactions of this address, grouped by nonce.
	ContentFrom(addr common.Address) ([]*types.Transaction, []*types.Transaction)

	// Locals retrieves the accounts currently considered local by the pool.
	Locals() []common.Address

	// Status returns the known status (unknown/pending/queued) of a transaction
	// identified by their hashes.
	Status(hash common.Hash) TxStatus
}
