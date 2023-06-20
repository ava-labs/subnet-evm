// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"fmt"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/database/prefixdb"
	"github.com/ava-labs/avalanchego/database/versiondb"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/hashing"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
)

var _ WarpDb = &warpDb{}

var (
	DefaultMaxDbSize = 10^6
)

const incrementSize	= 8

var (
	countPrefix			= []byte("c")
	messageIDPrefix		= []byte("m")
)

// WarpDb is a db that contains n messages, and automatically deletes
// the oldest entries to make space for new ones
type WarpDb interface {
	AddMessage(unsignedMessage *avalancheWarp.UnsignedMessage) error
	GetUnsignedMessage(messageID ids.ID) ([]byte, []byte, error)
}

type warpDb struct {
	basedb				*versiondb.Database
	countdb				database.Database 	//maps messageID -> unsignedMessage
	msgdb				database.Database	//maps count -> messageID, to keep track of old messages

	incrementer			uint64	//increments every time a new entry is published
	count				uint64	//tracks total number of entries
	size				uint64	//maximum limit on number of entries
}

func NewWarpDb(
	db		database.Database,
	size	uint64,
) WarpDb {
	w := warpDb{
		incrementer:	0,
		count: 			0,
		size:			size,
	}
	w.basedb		= versiondb.New(db)

	w.msgdb			= prefixdb.New(messageIDPrefix, w.basedb)

	w.countdb 		= prefixdb.New(countPrefix, w.basedb)
	return &w
}

// Record current status
func TakeSnapshot(w *warpDb) (uint64, uint64) {
	return w.incrementer, w.count
}

// If an operation fails, abort all changes and revert to previous status
func RevertToSnapshot(w *warpDb, incrementer uint64, count uint64) {
	w.basedb.Abort()
	w.incrementer = incrementer
	w.count = count
}

// Msg entries have a structure of [incrementbytes][unsignedmessagebytes], 
// this combines them
func prependIncrement(countbytes []byte, unsignedMessage []byte) []byte {
	return append(countbytes, unsignedMessage...)
}

// Msg entries have a structure of [incrementbytes][unsignedmessagebytes],
// this breaks them back into two pieces
func splitIncrement(msgEntry []byte) ([]byte, []byte){
	return msgEntry[:incrementSize], msgEntry[incrementSize:]
}

// Prune entries until there are at least maxDbsize entries.
// There is a guarantee that the oldest entries will be pruned, because
// Their incrememnt id will always be the lowest, and iterators will
// Organize keys lexographically.
func PruneEntries(w *warpDb) error {
	iter := w.countdb.NewIterator()
	for w.count > w.size {
		iter.Next()
		
		inc, messageID := iter.Key(), iter.Value()
		if err := w.msgdb.Delete(messageID); err != nil {
			return fmt.Errorf("failed to delete messageID from msgdb: %w", err) 
		}
		if err := w.countdb.Delete(inc); err != nil {
			return fmt.Errorf("failed to delte inc from countdb: %w", err)
		}
		w.count--
	}

	return nil
}

// Return the increment identifier and unsigned message associated with the messageID
func (w* warpDb) GetUnsignedMessage(messageID ids.ID) ([]byte, []byte, error) {
	msgEntry, err := w.msgdb.Get(messageID[:])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get unsigned message %s from warp database: %w", messageID.String(), err)
	}

	increment, unsignedMessage := splitIncrement(msgEntry)
	return increment, unsignedMessage, nil
}

// In the case when a node restarts, and possibly changes its bls key, the cache gets emptied but the database does not.
// So to avoid having incorrect signatures saved in the database after a bls key change, we save the full message in the database.
// Whereas for the cache, after the node restart, the cache would be emptied so we can directly save the signatures.
func (w* warpDb) AddMessage(unsignedMessage *avalancheWarp.UnsignedMessage) error {
	var (
		messageID = hashing.ComputeHash256Array(unsignedMessage.Bytes())
		countbytes = database.PackUInt64(w.incrementer)
	)

	sInc, sCount := TakeSnapshot(w) //revert to these values in case of a failure

	defer w.basedb.Abort()

	has, err := w.msgdb.Has(messageID[:]) 
	if err != nil {
		RevertToSnapshot(w, sInc, sCount)
		return fmt.Errorf("failed to check if message in db: %w", err)
	}

	switch has {
	case true:	//if this message has already appeared, simply put it in the new countdb with a higher increment
		increment, _, err := w.GetUnsignedMessage(messageID)
		if err != nil {
			RevertToSnapshot(w, sInc, sCount)
			return fmt.Errorf("failed to get unsigned message from db: %w", err)
		}
		if err := w.countdb.Delete(increment); err != nil {
			RevertToSnapshot(w, sInc, sCount)
			return fmt.Errorf("failed to delete item from countdb: %w", err)
		}
		if err := w.countdb.Put(countbytes, messageID[:]); err != nil {
			RevertToSnapshot(w, sInc, sCount)
			return fmt.Errorf("failed to put item into countdb: %w", err)
		}
		
	case false:
		msgEntry := prependIncrement(countbytes, unsignedMessage.Bytes())
		if err := w.msgdb.Put(messageID[:], msgEntry); err != nil {
			RevertToSnapshot(w, sInc, sCount)
			return fmt.Errorf("failed to put warp signature in db: %w", err)
		}
	
		if err := w.countdb.Put(countbytes, messageID[:]); err != nil {
			RevertToSnapshot(w, sInc, sCount)
			return fmt.Errorf("failed to put timestamp signature in db: %w", err)
		}

		w.count++
		if w.count > w.size {
			if err := PruneEntries(w); err != nil {
				RevertToSnapshot(w, sInc, sCount)
				return fmt.Errorf("failed to prune old entries in db: %w", err)
			}
		}
	}

	if err := w.basedb.Commit(); err != nil {
		RevertToSnapshot(w, sInc, sCount)
		return fmt.Errorf("failed to commit changes to db")
	}

	w.incrementer++
	return nil
}
