// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"fmt"

	"github.com/ava-labs/avalanchego/cache"
	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ethereum/go-ethereum/log"
)

var _ Backend = &backend{}

// Backend tracks signature-eligible warp messages and provides an interface to fetch them.
// The backend is also used to query for warp message signatures by the signature request handler.
type Backend interface {
	// AddMessage signs [unsignedMessage] and adds it to the warp backend database
	AddMessage(unsignedMessage *avalancheWarp.UnsignedMessage) error

	// GetSignature returns the signature of the requested message hash.
	GetSignature(messageHash ids.ID) ([bls.SignatureLen]byte, error)

	// GetMessage retrieves the [unsignedMessage] from the warp backend database if available
	GetMessage(messageHash ids.ID) (*avalancheWarp.UnsignedMessage, error)

	// Clear clears the entire db
	Clear() error
}

// backend implements Backend, keeps track of warp messages, and generates message signatures.
type backend struct {
	db             database.Database
	snowCtx        *snow.Context
	signatureCache *cache.LRU[ids.ID, [bls.SignatureLen]byte]
	messageCache   *cache.LRU[ids.ID, *avalancheWarp.UnsignedMessage]
}

// NewBackend creates a new Backend, and initializes the signature cache and message tracking database.
func NewBackend(snowCtx *snow.Context, db database.Database, cacheSize int) Backend {
	return &backend{
		db:             db,
		snowCtx:        snowCtx,
		signatureCache: &cache.LRU[ids.ID, [bls.SignatureLen]byte]{Size: cacheSize},
		messageCache:   &cache.LRU[ids.ID, *avalancheWarp.UnsignedMessage]{Size: cacheSize},
	}
}

func (w *backend) Clear() error {
	w.signatureCache.Flush()
	return database.Clear(w.db, w.db)
}

func (w *backend) AddMessage(unsignedMessage *avalancheWarp.UnsignedMessage) error {
	messageID := unsignedMessage.ID()

	// In the case when a node restarts, and possibly changes its bls key, the cache gets emptied but the database does not.
	// So to avoid having incorrect signatures saved in the database after a bls key change, we save the full message in the database.
	// Whereas for the cache, after the node restart, the cache would be emptied so we can directly save the signatures.
	if err := w.db.Put(messageID[:], unsignedMessage.Bytes()); err != nil {
		return fmt.Errorf("failed to put warp signature in db: %w", err)
	}

	var signature [bls.SignatureLen]byte
	sig, err := w.snowCtx.WarpSigner.Sign(unsignedMessage)
	if err != nil {
		return fmt.Errorf("failed to sign warp message: %w", err)
	}

	copy(signature[:], sig)
	w.signatureCache.Put(messageID, signature)
	log.Debug("Adding warp message to backend", "messageID", messageID)
	return nil
}

func (w *backend) GetSignature(messageID ids.ID) ([bls.SignatureLen]byte, error) {
	log.Debug("Getting warp message from backend", "messageID", messageID)
	if sig, ok := w.signatureCache.Get(messageID); ok {
		return sig, nil
	}

	unsignedMessage, err := w.GetMessage(messageID)
	if err != nil {
		return [bls.SignatureLen]byte{}, fmt.Errorf("failed to get warp message %s from db: %w", messageID.String(), err)
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

func (w *backend) GetMessage(messageID ids.ID) (*avalancheWarp.UnsignedMessage, error) {
	if message, ok := w.messageCache.Get(messageID); ok {
		return message, nil
	}

	unsignedMessageBytes, err := w.db.Get(messageID[:])
	if err != nil {
		return nil, fmt.Errorf("failed to get warp message %s from db: %w", messageID.String(), err)
	}

	unsignedMessage, err := avalancheWarp.ParseUnsignedMessage(unsignedMessageBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse unsigned message %s: %w", messageID.String(), err)
	}
	w.messageCache.Put(messageID, unsignedMessage)

	return unsignedMessage, nil
}
