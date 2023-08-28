// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package load

import (
	"sync"
	"time"

	pwarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/utils/predicate"
	"github.com/ava-labs/subnet-evm/warp/payload"
	"github.com/ava-labs/subnet-evm/x/warp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

type txTracker struct {
	lock   sync.Mutex
	closed bool
	done   chan struct{}

	issued   map[common.Hash]time.Time
	observer func(float64)
}

func newTxTracker(observer func(float64)) *txTracker {
	return &txTracker{
		issued:   make(map[common.Hash]time.Time),
		observer: observer,
		done:     make(chan struct{}),
	}
}

// removeMethodID removes the first 4 bytes of data, which is the method id
func removeMethodID(data []byte) []byte {
	if len(data) < 4 {
		log.Error("invalid warp message data", "data", common.Bytes2Hex(data))
		panic("invalid warp message data")
	}
	return data[4:]
}

// getSendTxAwmID returns a unique identifier for the awm message contained in
// the transaction. This is calculated by hashing the payload of the warp
// message. If the transaction is not a well formed warp message, this panics.
func getSendTxAwmID(tx *types.Transaction) common.Hash {
	input := removeMethodID(tx.Data())
	parsedInput, err := warp.UnpackSendWarpMessageInput(input)
	if err != nil {
		log.Error("failed to parse warp message input", "err", err)
		panic(err)
	}
	return crypto.Keccak256Hash(parsedInput.Payload)
}

// getReceiveTxAwmID returns a unique identifier for the awm message contained
// in the transaction. This is calculated by hashing the payload of the warp
// message. If the transaction is not a well formed warp message, this panics.
func getReceiveTxAwmID(tx *types.Transaction) common.Hash {
	storageSlots := make([]byte, 0)
	for _, tuple := range tx.AccessList() {
		if tuple.Address != warp.ContractAddress {
			continue
		}
		storageSlots = append(
			storageSlots, predicate.HashSliceToBytes(tuple.StorageKeys)...)
	}
	unpackedPredicateBytes, err := predicate.UnpackPredicate(storageSlots)
	if err != nil {
		log.Error("failed to unpack predicate bytes", "err", err)
		panic(err)
	}
	msg, err := pwarp.ParseMessage(unpackedPredicateBytes)
	if err != nil {
		log.Error("failed to parse warp message", "err", err)
		panic(err)
	}
	payload, err := payload.ParseAddressedPayload(msg.Payload)
	if err != nil {
		log.Error("failed to parse addressed payload", "err", err)
		panic(err)
	}
	return crypto.Keccak256Hash(payload.Payload)
}

func (tt *txTracker) IssueTx(tx *types.Transaction) {
	tt.lock.Lock()
	defer tt.lock.Unlock()

	id := getSendTxAwmID(tx)
	tt.issued[id] = time.Now()
}

func (tt *txTracker) ConfirmTx(tx *types.Transaction) {
	tt.lock.Lock()
	defer tt.lock.Unlock()

	id := getReceiveTxAwmID(tx)
	start, ok := tt.issued[id]
	if !ok {
		panic("unexpected confirm " + id.Hex())
	}
	duration := time.Since(start)
	tt.observer(duration.Seconds())

	delete(tt.issued, id)
	tt.checkDone()
}

func (tt *txTracker) Close() {
	tt.lock.Lock()
	defer tt.lock.Unlock()

	tt.closed = true
	tt.checkDone()
}

// assumes lock is held
func (tt *txTracker) checkDone() {
	if !tt.closed || len(tt.issued) > 0 {
		return
	}
	close(tt.done)
}
