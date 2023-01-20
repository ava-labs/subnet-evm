// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"fmt"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/database/prefixdb"
	"github.com/ava-labs/avalanchego/database/versiondb"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/vms/platformvm/teleporter"
	lru "github.com/hashicorp/golang-lru"
)

var (
	_ WarpBackend = &WarpMessagesDB{}

	dbPrefix = []byte("warp_messages")
)

const (
	signatureCacheSize = 500
)

// WarpBackend keeps track of messages that are accepted by the warp precompiles and add them into a database.
// The backend is also used to query for warp message signatures by the signature request handler.
type WarpBackend interface {
	// AddMessage is called in the precompile OnAccept, to add warp messages into the database.
	AddMessage(ctx context.Context, unsignedMessage *teleporter.UnsignedMessage) error

	// GetSignature returns the signature of the requested message hash.
	GetSignature(ctx context.Context, messageHash ids.ID) ([]byte, error)
}

// WarpMessagesDB implements WarpBackend, keeping track of warp messages, and generating message signatures.
type WarpMessagesDB struct {
	database.Database
	snowCtx        *snow.Context
	signatureCache *lru.Cache
}

// NewWarpMessagesDB creates a new WarpMessagesDB, and initializes the signature cache and message tracking database.
func NewWarpMessagesDB(snowCtx *snow.Context, vmDB *versiondb.Database) (*WarpMessagesDB, error) {
	signatureCache, err := lru.New(signatureCacheSize)
	if err != nil {
		return nil, err
	}

	db := &WarpMessagesDB{
		Database:       prefixdb.New(dbPrefix, vmDB),
		snowCtx:        snowCtx,
		signatureCache: signatureCache,
	}

	return db, nil
}

func (w *WarpMessagesDB) AddMessage(ctx context.Context, unsignedMessage *teleporter.UnsignedMessage) error {
	messageHash, err := ids.ToID(unsignedMessage.Bytes())
	if err != nil {
		return fmt.Errorf("failed to add message with key %s to warp database: %w", messageHash.String(), err)
	}

	return w.Put(messageHash[:], unsignedMessage.Bytes())
}

func (w *WarpMessagesDB) GetSignature(ctx context.Context, messageHash ids.ID) ([]byte, error) {
	if sig, ok := w.signatureCache.Get(messageHash[:]); ok {
		return sig.([]byte), nil
	}

	messageBytes, err := w.Get(messageHash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to get warp message %s from db: %w", messageHash.String(), err)
	}

	unsignedMessage, err := teleporter.ParseUnsignedMessage(messageBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse warp message %s: %w", messageHash.String(), err)
	}

	signature, err := w.snowCtx.TeleporterSigner.Sign(unsignedMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to sign warp message %s: %w", messageHash.String(), err)
	}

	w.signatureCache.Add(messageHash[:], signature)
	return signature, nil
}
