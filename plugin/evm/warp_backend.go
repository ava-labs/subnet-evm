// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"fmt"

	"github.com/ava-labs/avalanchego/snow"

	"github.com/ava-labs/avalanchego/database/prefixdb"
	"github.com/ava-labs/avalanchego/database/versiondb"

	"github.com/ava-labs/avalanchego/database"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/platformvm/teleporter"
)

var (
	_ WarpBackend = &WarpMessagesDB{}

	dbPrefix = []byte("warp_messages")
)

type WarpBackend interface {
	AddMessage(ctx context.Context, unsignedMessage *teleporter.UnsignedMessage) error
	GetSignature(ctx context.Context, messageHash ids.ID) ([]byte, error)
}

type WarpMessagesDB struct {
	database.Database
	snowCtx *snow.Context
}

func NewWarpMessagesDB(snowCtx *snow.Context, vmDB *versiondb.Database) *WarpMessagesDB {
	return &WarpMessagesDB{
		Database: prefixdb.New(dbPrefix, vmDB),
		snowCtx:  snowCtx,
	}
}

func (w *WarpMessagesDB) AddMessage(ctx context.Context, unsignedMessage *teleporter.UnsignedMessage) error {
	messageHash, err := ids.ToID(unsignedMessage.Bytes())
	if err != nil {
		return fmt.Errorf("failed to add message with key %s to warp database, err: %e", messageHash.String(), err)
	}

	w.Put(messageHash[:], unsignedMessage.Bytes())
	return nil
}

func (w *WarpMessagesDB) GetSignature(ctx context.Context, messageHash ids.ID) ([]byte, error) {
	messageBytes, err := w.Get(messageHash[:])
	if err != nil {
		return nil, err
	}

	unsignedMessage, err := teleporter.ParseUnsignedMessage(messageBytes)
	if err != nil {
		return nil, err
	}

	return w.snowCtx.TeleporterSigner.Sign(unsignedMessage)
}
