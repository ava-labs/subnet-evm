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
	_ WarpBackend = &warpBackend{}

	dbPrefix = []byte("warp")
)

// WarpBackend tracks signature eligible warp messages and provides an interface to fetch them.
// The backend is also used to query for warp message signatures by the signature request handler.
type WarpBackend interface {
	// AddMessage signs [unsignedMessage] and adds it to the warp backend database
	AddMessage(ctx context.Context, unsignedMessage *teleporter.UnsignedMessage) error

	// GetSignature returns the signature of the requested message hash.
	GetSignature(ctx context.Context, messageHash ids.ID) ([]byte, error)
}

// warpBackend implements WarpBackend, keeps track of warp messages, and generates message signatures.
type warpBackend struct {
	database.Database
	snowCtx        *snow.Context
	signatureCache *cache.LRU
}

// NewWarpBackend creates a new warpBackend, and initializes the signature cache and message tracking database.
func NewWarpBackend(snowCtx *snow.Context, vmDB *versiondb.Database, signatureCacheSize int) WarpBackend {
	return &warpBackend{
		Database:       prefixdb.New(dbPrefix, vmDB),
		snowCtx:        snowCtx,
		signatureCache: &cache.LRU{Size: signatureCacheSize},
	}
}

func (w *warpBackend) AddMessage(ctx context.Context, unsignedMessage *teleporter.UnsignedMessage) error {
	messageHashBytes := hashing.ComputeHash256(unsignedMessage.Bytes())
	messageHash, err := ids.ToID(messageHashBytes)
	if err != nil {
		return fmt.Errorf("failed to generate message hash for warp message db: %w", err)
	}

	// We generate the signature here and only save the signature in the db.
	// It is left to smart contracts built on top of Warp to save messages if required.
	signature, err := w.snowCtx.TeleporterSigner.Sign(unsignedMessage)
	if err != nil {
		return fmt.Errorf("failed to sign warp message %s: %w", messageHash.String(), err)
	}

	return w.Put(messageHash[:], signature)
}

func (w *warpBackend) GetSignature(ctx context.Context, messageHash ids.ID) ([]byte, error) {
	// Attempt to get the signature from cache before calling the db.
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
