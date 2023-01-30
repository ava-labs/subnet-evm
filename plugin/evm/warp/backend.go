// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"context"
	"fmt"

	"github.com/ava-labs/avalanchego/cache"
	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/ava-labs/avalanchego/vms/platformvm/teleporter"
)

var _ WarpBackend = &warpBackend{}

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
	db             database.Database
	snowCtx        *snow.Context
	signatureCache *cache.LRU
}

// NewWarpBackend creates a new WarpBackend, and initializes the signature cache and message tracking database.
func NewWarpBackend(snowCtx *snow.Context, db database.Database, signatureCacheSize int) WarpBackend {
	return &warpBackend{
		db:             db,
		snowCtx:        snowCtx,
		signatureCache: &cache.LRU{Size: signatureCacheSize},
	}
}

func (w *warpBackend) AddMessage(ctx context.Context, unsignedMessage *teleporter.UnsignedMessage) error {
	messageID := hashing.ComputeHash256Array(unsignedMessage.Bytes())

	// We save the message instead of signature for db, in case for bls key changes.
	if err := w.db.Put(messageID[:], unsignedMessage.Bytes()); err != nil {
		return fmt.Errorf("failed to put warp signature in db: %w", err)
	}

	signature, err := w.snowCtx.TeleporterSigner.Sign(unsignedMessage)
	if err != nil {
		return fmt.Errorf("failed to sign warp message: %w", err)
	}

	w.signatureCache.Put(ids.ID(messageID), signature)
	return nil
}

func (w *warpBackend) GetSignature(ctx context.Context, messageID ids.ID) ([]byte, error) {
	if sig, ok := w.signatureCache.Get(messageID); ok {
		return sig.([]byte), nil
	}

	unsignedMessageBytes, err := w.db.Get(messageID[:])
	if err != nil {
		return nil, fmt.Errorf("failed to get warp message %s from db: %w", messageID.String(), err)
	}

	unsignedMessage, err := teleporter.ParseUnsignedMessage(unsignedMessageBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse unsigned message %s: %w", messageID.String(), err)
	}

	signature, err := w.snowCtx.TeleporterSigner.Sign(unsignedMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to sign warp message: %w", err)
	}

	w.signatureCache.Put(messageID[:], signature)
	return signature, nil
}
