// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"fmt"

	"github.com/ava-labs/avalanchego/cache"
	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/database/prefixdb"
	"github.com/ava-labs/avalanchego/database/versiondb"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/ava-labs/avalanchego/vms/platformvm/teleporter"
)

var (
	_ WarpBackend = &warpMessagesDB{}

	dbPrefix = []byte("warp_messages")
)

// WarpBackend keeps track of messages that are accepted by the warp precompiles and add them into a database.
// The backend is also used to query for warp message signatures by the signature request handler.
type WarpBackend interface {
	// AddMessage is called in the precompile OnAccept, to add warp messages into the database.
	AddMessage(ctx context.Context, unsignedMessage *teleporter.UnsignedMessage) error

	// GetSignature returns the signature of the requested message hash.
	GetSignature(ctx context.Context, messageHash ids.ID) ([]byte, error)
}

// warpMessagesDB implements WarpBackend, keeping track of warp messages, and generating message signatures.
type warpMessagesDB struct {
	database.Database
	snowCtx        *snow.Context
	signatureCache *cache.LRU
}

// NewWarpMessagesDB creates a new warpMessagesDB, and initializes the signature cache and message tracking database.
func NewWarpMessagesDB(snowCtx *snow.Context, vmDB *versiondb.Database, signatureCacheSize int) WarpBackend {
	return &warpMessagesDB{
		Database:       prefixdb.New(dbPrefix, vmDB),
		snowCtx:        snowCtx,
		signatureCache: &cache.LRU{Size: signatureCacheSize},
	}
}

func (w *warpMessagesDB) AddMessage(ctx context.Context, unsignedMessage *teleporter.UnsignedMessage) error {
	messageHashBytes := hashing.ComputeHash256(unsignedMessage.Bytes())
	messageHash, err := ids.ToID(messageHashBytes)
	if err != nil {
		return fmt.Errorf("failed to generate message hash for warp message db: %w", err)
	}

	signature, err := w.snowCtx.TeleporterSigner.Sign(unsignedMessage)
	if err != nil {
		return fmt.Errorf("failed to sign warp message %s: %w", messageHash.String(), err)
	}

	return w.Put(messageHash[:], signature)
}

func (w *warpMessagesDB) GetSignature(ctx context.Context, messageHash ids.ID) ([]byte, error) {
	if sig, ok := w.signatureCache.Get(messageHash[:]); ok {
		return sig.([]byte), nil
	}

	signature, err := w.Get(messageHash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to get warp signature for message %s from db: %w", messageHash.String(), err)
	}

	w.signatureCache.Put(messageHash[:], signature)
	return signature, nil
}
