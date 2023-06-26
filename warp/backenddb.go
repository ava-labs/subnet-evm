// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/database/prefixdb"
	"github.com/ava-labs/avalanchego/database/versiondb"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/ava-labs/avalanchego/utils/linkedhashmap"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
)

var _ WarpDb = &warpDb{}

var (
	tsPrefix			= []byte("t")
	messageIDPrefix		= []byte("m")
)

var (
	keyExistsError = errors.New("key exists")
)

const (
	tsSize 			= 8
	messageIDSize 	= 32

	Nanosecond 		= uint64(time.Nanosecond)
	Microsecond 	= uint64(time.Microsecond)
	Millisecond 	= uint64(time.Millisecond)
	Second 			= uint64(time.Second)
	Minute			= uint64(time.Minute)
	Hour			= uint64(time.Hour)
	Day				= Hour * 24
	Week			= Day * 7
)

type prunedMap linkedhashmap.LinkedHashmap[int, []byte]

func NewPrunedMap() prunedMap {
	return linkedhashmap.New[int, []byte]()
}

type warpDbConfig struct {
	autroprune		bool
	maxPruneAmt		uint
}

func GetDefaultWarpDbConfig() warpDbConfig {
	return warpDbConfig{
		true,
		10,
	}
}

// WarpDb is a db that contains n messages, and automatically deletes
// the oldest entries to make space for new ones
type WarpDb interface {
	AddMessage(unsignedMessage *avalancheWarp.UnsignedMessage) (int, prunedMap, error)
	GetUnsignedMessage(messageID ids.ID) ([]byte, []byte, error)
	PruneEntries(amt uint) (int, prunedMap, error)
}

type warpDb struct {
	basedb			*versiondb.Database
	msgdb			database.Database 	//maps messageID -> unsignedMessage
	tsdb			database.Database	//maps ts -> messageID, to keep track of old messages
	lock			sync.RWMutex
	config			warpDbConfig			

	count			uint64	//tracks total number of entries
	timeLimit		uint64	//maximum time limit for an entry to be relevant
}

func NewWarpDb(
	db			database.Database,
	timeLimit	uint64,
	config		warpDbConfig,
) WarpDb {
	w := warpDb{
		count: 			0,
		timeLimit:		timeLimit,
	}
	w.basedb		= versiondb.New(db)
	w.msgdb			= prefixdb.New(messageIDPrefix, w.basedb)
	w.tsdb 			= prefixdb.New(tsPrefix, w.basedb)
	w.config = config
	return &w
}

// Only put entry in database if the key does not exist in the database.
// If key does exist, throw an error.
func dBPutSafe(db database.Database, key []byte, value []byte) error {
	has, err :=  db.Has(key)
	if err != nil {
		return err
	}
	if has {
		return keyExistsError
	}
	err = db.Put(key, value)
	if err != nil {
		return err
	}

	return nil
}

// Msg entries have a structure of [timestampbytes][unsignedmessagebytes], 
// this combines them
func prependTimestamp(tsBytes []byte, unsignedMessage []byte) []byte {
	return append(tsBytes, unsignedMessage...)
}

// Msg entries have a structure of [timestampbytes][unsignedmessagebytes],
// this breaks them back into two pieces
func splitTimestamp(msgEntry []byte) ([]byte, []byte){
	return msgEntry[:tsSize], msgEntry[tsSize:]
}

func pruneEntries(w *warpDb, max uint) (int, prunedMap, error) {
	var (
		threshold = uint64(time.Now().UnixNano()) - w.timeLimit
		startCount = w.count
		prunedMsgs = NewPrunedMap()
	)
	iter := w.tsdb.NewIterator()
	defer iter.Release()

	for i := 0; iter.Next() && i < int(w.config.maxPruneAmt); i++ {
		
		ts, messageID := iter.Key(), iter.Value()
		
		msgTs, err := database.ParseUInt64(ts)
		if err != nil {
			return int(startCount - w.count), prunedMsgs, err
		}

		// ts iterator sorts items by increasing value, so when
		// messages that are within the threshold are reached
		// the rest of the messages are not old, and pruning can stop
		if msgTs >= threshold {
			break
		}

		if err := w.msgdb.Delete(messageID); err != nil {
			return int(startCount - w.count), prunedMsgs, fmt.Errorf("failed to delete messageID %s from msgdb: %w",messageID, err)
		}
		if err := w.tsdb.Delete(ts); err != nil {
			return int(startCount - w.count), prunedMsgs, fmt.Errorf("failed to delete inc from tsdb: %w", err)
		}
		
		w.count--
	}

	return int(startCount - w.count), prunedMsgs, nil
}

// Prune entries until there are at least maxDbsize entries.
// There is a guarantee that the oldest entries will be pruned, because
// Their incrememnt id will always be the lowest, and iterators will
// Organize keys lexographically.
func (w *warpDb) PruneEntries(amt uint) (int, prunedMap, error) {
	w.lock.RLock()
	defer w.lock.RUnlock()
	defer w.basedb.Abort()

	numPruned, prunedMsgs, err := pruneEntries(w, amt)
	if err != nil {
		return numPruned, prunedMsgs, fmt.Errorf("error pruning entries: %w", err)
	}

	if err := w.basedb.Commit(); err != nil {
		return numPruned, prunedMsgs, fmt.Errorf("failed to commit changes to db: %w", err)
	}
	return numPruned, prunedMsgs, nil
}

// Return the timestamp identifier and unsigned message associated with the messageID
func (w* warpDb) GetUnsignedMessage(messageID ids.ID) ([]byte, []byte, error) {
	w.lock.RLock()
	defer w.lock.RUnlock()

	msgEntry, err := w.msgdb.Get(messageID[:])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get unsigned message %s from warp database: %w", messageID.String(), err)
	}

	timestamp, unsignedMessage := splitTimestamp(msgEntry)
	return timestamp, unsignedMessage, nil
}

// In the case when a node restarts, and possibly changes its bls key, the cache gets emptied but the database does not.
// So to avoid having incorrect signatures saved in the database after a bls key change, we save the full message in the database.
// Whereas for the cache, after the node restart, the cache would be emptied so we can directly save the signatures.
func (w* warpDb) AddMessage(unsignedMessage *avalancheWarp.UnsignedMessage) (int, prunedMap, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	defer w.basedb.Abort()

	var (
		messageID = hashing.ComputeHash256Array(unsignedMessage.Bytes())
		numPruned = 0
		prunedMsgs = NewPrunedMap()
	)

	has, err := w.msgdb.Has(messageID[:])
	if err != nil {
		return numPruned, prunedMsgs, err
	}

	if has {
	// if this message has already appeared, delete the old ts
	// also add to pruned msgs for debugging purposes
		ts, _, err := w.GetUnsignedMessage(messageID)
		if err != nil {
			return numPruned, prunedMsgs, err
		}
		if err := w.tsdb.Delete(ts); err != nil {
			return numPruned, prunedMsgs, err
		}

		msgTs, err := database.ParseUInt64(ts)
		if err != nil {
			return numPruned, prunedMsgs, err
		}
		
		prunedMsgs.Put(int(msgTs), messageID[:])
		w.count--
	}

	if w.config.autroprune {
		numPruned, prunedMsgs, err = pruneEntries(w, w.config.maxPruneAmt)
		if err != nil {
			return numPruned, prunedMsgs, fmt.Errorf("failed to prune old entries in db: %w", err)
		}
	}

	// Sleep one second to ensure ts will not exist in db
	time.Sleep(time.Nanosecond)
	tsBytes := database.PackUInt64(uint64(time.Now().UnixNano()))
	msgEntry := prependTimestamp(tsBytes, unsignedMessage.Bytes())

	if err := w.msgdb.Put(messageID[:], msgEntry); err != nil {
		return numPruned, prunedMsgs, err
	}

	//put safe to ensure that timestamps are never overwritten
	if err := dBPutSafe(w.tsdb, tsBytes, messageID[:]); err != nil {
		return numPruned, prunedMsgs, err
	}

	if err := w.basedb.Commit(); err != nil {
		return numPruned, prunedMsgs, err
	}
	
	w.count++
	return numPruned, prunedMsgs, nil
}
