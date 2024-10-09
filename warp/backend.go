// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"context"
	"errors"
	"fmt"

	"github.com/ava-labs/avalanchego/cache"
	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/consensus/snowman"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"
	"github.com/ava-labs/subnet-evm/warp/messages"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
)

var (
	_                         Backend = &backend{}
	errParsingOffChainMessage         = errors.New("failed to parse off-chain message")
)

const batchSize = ethdb.IdealBatchSize

type BlockClient interface {
	GetAcceptedBlock(ctx context.Context, blockID ids.ID) (snowman.Block, error)
}

// Backend tracks signature-eligible warp messages and provides an interface to fetch them.
// The backend is also used to query for warp message signatures by the signature request handler.
type Backend interface {
	// AddMessage signs [unsignedMessage] and adds it to the warp backend database
	AddMessage(unsignedMessage *avalancheWarp.UnsignedMessage) error

	// GetMessageSignature validates the message and returns the signature of the requested message.
	GetMessageSignature(message *avalancheWarp.UnsignedMessage) ([bls.SignatureLen]byte, error)

	// GetBlockSignature validates blockID and returns the signature of the requested message hash.
	GetBlockSignature(blockID ids.ID) ([bls.SignatureLen]byte, error)

	// GetMessage retrieves the [unsignedMessage] from the warp backend database if available
	// TODO: After E-Upgrade, the backend no longer needs to store the mapping from messageHash
	// to unsignedMessage (and this method can be removed).
	GetMessage(messageHash ids.ID) (*avalancheWarp.UnsignedMessage, error)

	// ValidateMessage validates the [unsignedMessage] and returns an error if the message is invalid.
	ValidateMessage(unsignedMessage *avalancheWarp.UnsignedMessage) error

	// ValidateBlockMessage validates the block message with the given [blockID] and returns an error if the message is invalid.
	ValidateBlockMessage(blockID ids.ID) error

	// SignMessage signs the [unsignedMessage] and returns the signature.
	SignMessage(unsignedMessage *avalancheWarp.UnsignedMessage) ([bls.SignatureLen]byte, error)

	// SignBlock signs the block message with the given [blockID] and returns the signature.
	SignBlock(blockID ids.ID) ([bls.SignatureLen]byte, error)

	// Clear clears the entire db
	Clear() error
}

// backend implements Backend, keeps track of warp messages, and generates message signatures.
type backend struct {
	networkID                 uint32
	sourceChainID             ids.ID
	db                        database.Database
	warpSigner                avalancheWarp.Signer
	blockClient               BlockClient
	messageSignatureCache     *cache.LRU[ids.ID, [bls.SignatureLen]byte]
	blockSignatureCache       *cache.LRU[ids.ID, [bls.SignatureLen]byte]
	messageCache              *cache.LRU[ids.ID, *avalancheWarp.UnsignedMessage]
	offchainAddressedCallMsgs map[ids.ID]*avalancheWarp.UnsignedMessage
}

// NewBackend creates a new Backend, and initializes the signature cache and message tracking database.
func NewBackend(
	networkID uint32,
	sourceChainID ids.ID,
	warpSigner avalancheWarp.Signer,
	blockClient BlockClient,
	db database.Database,
	cacheSize int,
	offchainMessages [][]byte,
) (Backend, error) {
	b := &backend{
		networkID:                 networkID,
		sourceChainID:             sourceChainID,
		db:                        db,
		warpSigner:                warpSigner,
		blockClient:               blockClient,
		messageSignatureCache:     &cache.LRU[ids.ID, [bls.SignatureLen]byte]{Size: cacheSize},
		blockSignatureCache:       &cache.LRU[ids.ID, [bls.SignatureLen]byte]{Size: cacheSize},
		messageCache:              &cache.LRU[ids.ID, *avalancheWarp.UnsignedMessage]{Size: cacheSize},
		offchainAddressedCallMsgs: make(map[ids.ID]*avalancheWarp.UnsignedMessage),
	}
	return b, b.initOffChainMessages(offchainMessages)
}

func (b *backend) initOffChainMessages(offchainMessages [][]byte) error {
	for i, offchainMsg := range offchainMessages {
		unsignedMsg, err := avalancheWarp.ParseUnsignedMessage(offchainMsg)
		if err != nil {
			return fmt.Errorf("%w at index %d: %w", errParsingOffChainMessage, i, err)
		}

		if unsignedMsg.NetworkID != b.networkID {
			return fmt.Errorf("%w at index %d", avalancheWarp.ErrWrongNetworkID, i)
		}

		if unsignedMsg.SourceChainID != b.sourceChainID {
			return fmt.Errorf("%w at index %d", avalancheWarp.ErrWrongSourceChainID, i)
		}

		_, err = payload.ParseAddressedCall(unsignedMsg.Payload)
		if err != nil {
			return fmt.Errorf("%w at index %d as AddressedCall: %w", errParsingOffChainMessage, i, err)
		}
		b.offchainAddressedCallMsgs[unsignedMsg.ID()] = unsignedMsg
	}

	return nil
}

func (b *backend) Clear() error {
	b.messageSignatureCache.Flush()
	b.blockSignatureCache.Flush()
	b.messageCache.Flush()
	return database.Clear(b.db, batchSize)
}

func (b *backend) AddMessage(unsignedMessage *avalancheWarp.UnsignedMessage) error {
	messageID := unsignedMessage.ID()

	// In the case when a node restarts, and possibly changes its bls key, the cache gets emptied but the database does not.
	// So to avoid having incorrect signatures saved in the database after a bls key change, we save the full message in the database.
	// Whereas for the cache, after the node restart, the cache would be emptied so we can directly save the signatures.
	if err := b.db.Put(messageID[:], unsignedMessage.Bytes()); err != nil {
		return fmt.Errorf("failed to put warp signature in db: %w", err)
	}

	var signature [bls.SignatureLen]byte
	sig, err := b.warpSigner.Sign(unsignedMessage)
	if err != nil {
		return fmt.Errorf("failed to sign warp message: %w", err)
	}

	copy(signature[:], sig)
	b.messageSignatureCache.Put(messageID, signature)
	log.Debug("Adding warp message to backend", "messageID", messageID)
	return nil
}

func (b *backend) GetMessageSignature(unsignedMessage *avalancheWarp.UnsignedMessage) ([bls.SignatureLen]byte, error) {
	messageID := unsignedMessage.ID()

	log.Debug("Getting warp message from backend", "messageID", messageID)
	if sig, ok := b.messageSignatureCache.Get(messageID); ok {
		return sig, nil
	}

	if err := b.ValidateMessage(unsignedMessage); err != nil {
		return [bls.SignatureLen]byte{}, fmt.Errorf("failed to validate warp message: %w", err)
	}
	return b.signMessage(unsignedMessage)
}

func (b *backend) SignMessage(unsignedMessage *avalancheWarp.UnsignedMessage) ([bls.SignatureLen]byte, error) {
	messageID := unsignedMessage.ID()

	if sig, ok := b.messageSignatureCache.Get(messageID); ok {
		return sig, nil
	}

	return b.signMessage(unsignedMessage)
}

func (b *backend) ValidateMessage(unsignedMessage *avalancheWarp.UnsignedMessage) error {
	messageID := unsignedMessage.ID()

	if _, ok := b.messageSignatureCache.Get(messageID); ok {
		return nil
	}
	// Known on-chain messages should be signed
	if _, err := b.GetMessage(messageID); err == nil {
		return nil
	}

	// Try to parse the payload as an AddressedCall
	addressedCall, err := payload.ParseAddressedCall(unsignedMessage.Payload)
	if err != nil {
		return fmt.Errorf("failed to parse unknown message as AddressedCall: %w", err)
	}

	// Further, parse the payload to see if it is a known type.
	parsed, err := messages.Parse(addressedCall.Payload)
	if err != nil {
		return fmt.Errorf("failed to parse unknown message: %w", err)
	}

	// Check if the message is a known type that can be signed on demand
	signable, ok := parsed.(messages.Signable)
	if !ok {
		return fmt.Errorf("parsed message is not Signable: %T", signable)
	}

	// Check if the message should be signed according to its type
	if err := signable.VerifyMesssage(addressedCall.SourceAddress); err != nil {
		return fmt.Errorf("failed to verify Signable message: %w", err)
	}
	return nil
}

func (b *backend) ValidateBlockMessage(blockID ids.ID) error {
	if _, ok := b.blockSignatureCache.Get(blockID); ok {
		return nil
	}

	_, err := b.blockClient.GetAcceptedBlock(context.TODO(), blockID)
	if err != nil {
		return fmt.Errorf("failed to get block %s: %w", blockID, err)
	}

	return nil
}

func (b *backend) GetBlockSignature(blockID ids.ID) ([bls.SignatureLen]byte, error) {
	log.Debug("Getting block from backend", "blockID", blockID)
	if sig, ok := b.blockSignatureCache.Get(blockID); ok {
		return sig, nil
	}

	if err := b.ValidateBlockMessage(blockID); err != nil {
		return [bls.SignatureLen]byte{}, fmt.Errorf("failed to validate block message: %w", err)
	}

	return b.signBlock(blockID)
}

func (b *backend) SignBlock(blockID ids.ID) ([bls.SignatureLen]byte, error) {
	if sig, ok := b.blockSignatureCache.Get(blockID); ok {
		return sig, nil
	}

	return b.signBlock(blockID)
}

func (b *backend) GetMessage(messageID ids.ID) (*avalancheWarp.UnsignedMessage, error) {
	if message, ok := b.messageCache.Get(messageID); ok {
		return message, nil
	}
	if message, ok := b.offchainAddressedCallMsgs[messageID]; ok {
		return message, nil
	}

	unsignedMessageBytes, err := b.db.Get(messageID[:])
	if err != nil {
		return nil, fmt.Errorf("failed to get warp message %s from db: %w", messageID.String(), err)
	}

	unsignedMessage, err := avalancheWarp.ParseUnsignedMessage(unsignedMessageBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse unsigned message %s: %w", messageID.String(), err)
	}
	b.messageCache.Put(messageID, unsignedMessage)

	return unsignedMessage, nil
}

func (b *backend) signMessage(unsignedMessage *avalancheWarp.UnsignedMessage) ([bls.SignatureLen]byte, error) {
	sig, err := b.warpSigner.Sign(unsignedMessage)
	if err != nil {
		return [bls.SignatureLen]byte{}, fmt.Errorf("failed to sign warp message: %w", err)
	}

	var signature [bls.SignatureLen]byte
	copy(signature[:], sig)
	b.messageSignatureCache.Put(unsignedMessage.ID(), signature)
	return signature, nil
}

func (b *backend) signBlock(blockID ids.ID) ([bls.SignatureLen]byte, error) {
	blockHashPayload, err := payload.NewHash(blockID)
	if err != nil {
		return [bls.SignatureLen]byte{}, fmt.Errorf("failed to create new block hash payload: %w", err)
	}
	unsignedMessage, err := avalancheWarp.NewUnsignedMessage(b.networkID, b.sourceChainID, blockHashPayload.Bytes())
	if err != nil {
		return [bls.SignatureLen]byte{}, fmt.Errorf("failed to create new unsigned warp message: %w", err)
	}
	sig, err := b.warpSigner.Sign(unsignedMessage)
	if err != nil {
		return [bls.SignatureLen]byte{}, fmt.Errorf("failed to sign warp message: %w", err)
	}

	var signature [bls.SignatureLen]byte
	copy(signature[:], sig)
	b.blockSignatureCache.Put(blockID, signature)
	return signature, nil
}
