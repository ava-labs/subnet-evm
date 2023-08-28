// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package load

import (
	"fmt"
	"sync"
	"time"

	pwarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/utils/predicate"
	"github.com/ava-labs/subnet-evm/warp/payload"
	"github.com/ava-labs/subnet-evm/x/warp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
func removeMethodID(data []byte) ([]byte, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("data is too short: %d", len(data))
	}
	return data[4:], nil
}

// getSendTxAwmID returns a unique identifier for the awm message contained in
// the transaction. This is calculated by hashing the payload of the warp
// message.
func getSendTxAwmID(tx *types.Transaction) (common.Hash, error) {
	input, err := removeMethodID(tx.Data())
	if err != nil {
		return common.Hash{}, err
	}
	parsedInput, err := warp.UnpackSendWarpMessageInput(input)
	if err != nil {
		return common.Hash{}, err
	}
	return crypto.Keccak256Hash(parsedInput.Payload), nil
}

// getReceiveTxAwmID returns a unique identifier for the awm message contained
// in the transaction. This is calculated by hashing the payload of the warp
// message.
func getReceiveTxAwmID(tx *types.Transaction) (common.Hash, error) {
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
		return common.Hash{}, fmt.Errorf("failed to unpack predicate bytes: %w", err)
	}
	msg, err := pwarp.ParseMessage(unpackedPredicateBytes)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to parse message: %w", err)
	}
	payload, err := payload.ParseAddressedPayload(msg.Payload)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to parse payload: %w", err)
	}
	return crypto.Keccak256Hash(payload.Payload), nil
}

func (tt *txTracker) IssueTx(tx *types.Transaction) error {
	tt.lock.Lock()
	defer tt.lock.Unlock()

	id, err := getSendTxAwmID(tx)
	if err != nil {
		return err
	}
	tt.issued[id] = time.Now()
	return nil
}

func (tt *txTracker) ConfirmTx(tx *types.Transaction) error {
	tt.lock.Lock()
	defer tt.lock.Unlock()

	id, err := getReceiveTxAwmID(tx)
	if err != nil {
		return err
	}
	start, ok := tt.issued[id]
	if !ok {
		panic("unexpected confirm " + id.Hex())
	}
	duration := time.Since(start)
	tt.observer(duration.Seconds())

	delete(tt.issued, id)
	tt.checkDone()
	return nil
}

func (tt *txTracker) Close() error {
	tt.lock.Lock()
	defer tt.lock.Unlock()

	tt.closed = true
	tt.checkDone()
	return nil
}

// assumes lock is held
func (tt *txTracker) checkDone() {
	if !tt.closed || len(tt.issued) > 0 {
		return
	}
	close(tt.done)
}
