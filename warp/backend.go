// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"fmt"

	"github.com/ava-labs/avalanchego/cache"
	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/database/versiondb"
	"github.com/ava-labs/avalanchego/database/prefixdb"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/hashing"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
)

var _ WarpBackend = &warpBackend{}

var (
	DefaultMaxDbSize = 10^6
)

var (
	countPrefix			= []byte("c")
	messageIDPrefix		= []byte("m")
)

type warpBackendConfig struct {
	MaxDbSize	uint64
}

func NewWarpBackendConfig(MaxDbSize uint64) warpBackendConfig {
	return warpBackendConfig{
		MaxDbSize: MaxDbSize,
	}
}

func DefaultWarpBackendConfig() warpBackendConfig {
	return NewWarpBackendConfig(uint64(DefaultMaxDbSize))
}

// WarpBackend tracks signature eligible warp messages and provides an interface to fetch them.
// The backend is also used to query for warp message signatures by the signature request handler.
type WarpBackend interface {
	// AddMessage signs [unsignedMessage] and adds it to the warp backend database
	AddMessage(unsignedMessage *avalancheWarp.UnsignedMessage) error

	// GetSignature returns the signature of the requested message hash.
	GetSignature(messageHash ids.ID) ([bls.SignatureLen]byte, error)
}

// warpBackend implements WarpBackend, keeps track of warp messages, and generates message signatures.
type warpBackend struct {
	db				*versiondb.Database
	countdb			database.Database
	msgdb			database.Database
	snowCtx        *snow.Context
	signatureCache *cache.LRU[ids.ID, [bls.SignatureLen]byte]
	msgCount		uint64
	config			warpBackendConfig
}

// NewWarpBackend creates a new WarpBackend, and initializes the signature cache and message tracking database.
func NewWarpBackend(
	snowCtx 			*snow.Context,
	db 					database.Database,
	signatureCacheSize 	int,
) WarpBackend {
	w := &warpBackend {
		snowCtx: 			snowCtx,
		signatureCache: 	&cache.LRU[ids.ID, [bls.SignatureLen]byte]{Size: signatureCacheSize},
		msgCount: 			0,
	}
	//versiondb to ensure that msgdb & countdb are updated atomically
	w.db		= versiondb.New(db)
	//maps messageID -> unsignedMessage
	w.msgdb		= prefixdb.New(messageIDPrefix, w.db)
	//maps count -> messageID, to keep track of old messages
	w.countdb 	= prefixdb.New(countPrefix, w.db)
	w.config	= DefaultWarpBackendConfig()
	return w
}

func (w *warpBackend) AddMessage(unsignedMessage *avalancheWarp.UnsignedMessage) error {
	var (
		messageID 	= hashing.ComputeHash256Array(unsignedMessage.Bytes())
		countbytes	= database.PackUInt64(w.msgCount)
	)
	defer w.db.Abort()
	// In the case when a node restarts, and possibly changes its bls key, the cache gets emptied but the database does not.
	// So to avoid having incorrect signatures saved in the database after a bls key change, we save the full message in the database.
	// Whereas for the cache, after the node restart, the cache would be emptied so we can directly save the signatures.
	if err := w.msgdb.Put(messageID[:], unsignedMessage.Bytes()); err != nil {
		defer w.db.Abort()
		return fmt.Errorf("failed to put warp signature in db: %w", err)
	}

	// Add the message count -> messageID mapping
	if err := w.countdb.Put(countbytes, messageID[:]); err != nil {
		defer w.db.Abort()
		return fmt.Errorf("failed to put timestamp signature in db: %w", err)
	}

	if w.config.MaxDbSize <= w.msgCount {
		//offset by 1, because msg count should only be updated after committing
		if err := PruneEntry(w, w.msgCount-w.config.MaxDbSize+1); err != nil {
			defer w.db.Abort()
			return fmt.Errorf("failed to prune db")
		}
	}
	if err :=w.db.Commit(); err != nil {
		return fmt.Errorf("failed to commit changes to database")
	}
	w.msgCount++

	var signature [bls.SignatureLen]byte
	sig, err := w.snowCtx.WarpSigner.Sign(unsignedMessage)
	if err != nil {
		return fmt.Errorf("failed to sign warp message: %w", err)
	}

	copy(signature[:], sig)
	w.signatureCache.Put(messageID, signature)
	return nil
}

func (w *warpBackend) GetSignature(messageID ids.ID) ([bls.SignatureLen]byte, error) {
	if sig, ok := w.signatureCache.Get(messageID); ok {
		return sig, nil
	}

	unsignedMessageBytes, err := w.msgdb.Get(messageID[:])
	if err != nil {
		return [bls.SignatureLen]byte{}, fmt.Errorf("failed to get warp message %s from db: %w", messageID.String(), err)
	}

	unsignedMessage, err := avalancheWarp.ParseUnsignedMessage(unsignedMessageBytes)
	if err != nil {
		return [bls.SignatureLen]byte{}, fmt.Errorf("failed to parse unsigned message %s: %w", messageID.String(), err)
	}

	var signature [bls.SignatureLen]byte
	sig, err := w.snowCtx.WarpSigner.Sign(unsignedMessage)
	if err != nil {
		return [bls.SignatureLen]byte{}, fmt.Errorf("failed to sign warp message: %w", err)
	}

	copy(signature[:], sig)
	w.signatureCache.Put(messageID, signature)
	return signature, nil
}

func PruneEntry(w *warpBackend, msgCount uint64 ) error { 
	oldCount := database.PackUInt64(msgCount)
	has, err := w.countdb.Has(oldCount);
	if err != nil {
		return fmt.Errorf("Error reading countdb: %w", err)
	}

	if has {
		messageID, err := w.countdb.Get(oldCount)

		if err != nil {
			return fmt.Errorf("Error fetching from countdb: %w", err)
		}
		if err := w.msgdb.Delete(messageID); err != nil {
			return fmt.Errorf("Error deleting from messagedb: %w", err)
		}
		if err := w.countdb.Delete(oldCount); err != nil {
			return fmt.Errorf("Error deleting from countdb: %w", err)
		}
	}

	return nil
}

//in the event that the maxDBSize changes, prune all old entries
func PruneAllOldEntries(w *warpBackend) error {
	defer w.db.Abort()
	iter := w.countdb.NewIterator()

	for iter.Next() {
		count, err := database.ParseUInt64(iter.Key())
		if err != nil {
			return fmt.Errorf("Error parsing count: %w", err)
		}
		if w.msgCount - count >= w.config.MaxDbSize {
			if err := PruneEntry(w, w.msgCount - count); err != nil {
				return fmt.Errorf("Error pruning msg: %w", err)
			}
		}
	}
	if err := w.db.Commit(); err != nil {
		return fmt.Errorf("Error committing prunes: %w", err)
	}
	
	return nil
}