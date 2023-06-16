// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"fmt"
	"testing"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/database/memdb"
	"github.com/ava-labs/avalanchego/database/prefixdb"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/hashing"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/stretchr/testify/require"
)

var (
	sourceChainID      = ids.GenerateTestID()
	destinationChainID = ids.GenerateTestID()
	payload            = []byte("test")
)

var (
	testWarpBackendConfig = warpBackendConfig{
		MaxDbSize: 5,
	}
)

func TestAddAndGetValidMessage(t *testing.T) {
	db := memdb.New()

	snowCtx := snow.DefaultContextTest()
	sk, err := bls.NewSecretKey()
	require.NoError(t, err)
	snowCtx.WarpSigner = avalancheWarp.NewSigner(sk, sourceChainID)
	backend := NewWarpBackend(snowCtx, db, 500)

	// Create a new unsigned message and add it to the warp backend.
	unsignedMsg, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, payload)
	require.NoError(t, err)
	err = backend.AddMessage(unsignedMsg)
	require.NoError(t, err)

	// Verify that a signature is returned successfully, and compare to expected signature.
	messageID := hashing.ComputeHash256Array(unsignedMsg.Bytes())
	signature, err := backend.GetSignature(messageID)
	require.NoError(t, err)

	expectedSig, err := snowCtx.WarpSigner.Sign(unsignedMsg)
	require.NoError(t, err)
	require.Equal(t, expectedSig, signature[:])
}

func TestAddAndGetUnknownMessage(t *testing.T) {
	db := memdb.New()

	backend := NewWarpBackend(snow.DefaultContextTest(), db, 500)
	unsignedMsg, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, payload)
	require.NoError(t, err)

	// Try getting a signature for a message that was not added.
	messageID := hashing.ComputeHash256Array(unsignedMsg.Bytes())
	_, err = backend.GetSignature(messageID)
	require.Error(t, err)
}

func TestZeroSizedCache(t *testing.T) {
	db := memdb.New()

	snowCtx := snow.DefaultContextTest()
	sk, err := bls.NewSecretKey()
	require.NoError(t, err)
	snowCtx.WarpSigner = avalancheWarp.NewSigner(sk, sourceChainID)

	// Verify zero sized cache works normally, because the lru cache will be initialized to size 1 for any size parameter <= 0.
	backend := NewWarpBackend(snowCtx, db, 0)

	// Create a new unsigned message and add it to the warp backend.
	unsignedMsg, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, payload)
	require.NoError(t, err)
	err = backend.AddMessage(unsignedMsg)
	require.NoError(t, err)

	// Verify that a signature is returned successfully, and compare to expected signature.
	messageID := hashing.ComputeHash256Array(unsignedMsg.Bytes())
	signature, err := backend.GetSignature(messageID)
	require.NoError(t, err)

	expectedSig, err := snowCtx.WarpSigner.Sign(unsignedMsg)
	require.NoError(t, err)
	require.Equal(t, expectedSig, signature[:])
}

func TestDatabase(t *testing.T) {

	db := memdb.New()
	test := prefixdb.New([]byte("hello"), db)

	test.Put([]byte("test"), []byte("test"))
	x, err := test.Get([]byte("tes"))
	t.Log(x)
	t.Log(fmt.Printf("%T", err))
	t.Error()
	
}

//potential update to testaddandgetvalidmessage
func TestPruneEntry(t *testing.T) {
	db := memdb.New()
	snowCtx := snow.DefaultContextTest()
	sk, err := bls.NewSecretKey()
	require.NoError(t, err)
	snowCtx.WarpSigner = avalancheWarp.NewSigner(sk, sourceChainID)

	maxDbSize := 5
	backendConfig := warpBackendConfig{MaxDbSize: uint64(maxDbSize)}

	backend := NewWarpBackend(snowCtx, db, 500).(*warpBackend)
	backend.config = backendConfig

	// Create a new unsigned message and add it to the warp backend.

	for i := 0; i < maxDbSize; i++ {
		msg := append(payload, database.PackUInt64(uint64(i))...) //results in test0, test1 etc.
		unsignedMsg, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, msg)
		require.NoError(t, err)
		
		err = backend.AddMessage(unsignedMsg)
		require.NoError(t, err)
		require.EqualValues(t, backend.msgCount, i+1)

		//Go back through all messages that were added, ensure nothing was deleted
		for j := 0; j <= i; j++ {
			msgCountBytes := database.PackUInt64(uint64(j))
			prevMsg := append(payload, msgCountBytes...)
			prevUnsignedMsg, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, prevMsg)
			require.NoError(t, err)

			messageIDBytes, err := backend.countdb.Get(msgCountBytes)
			require.NoError(t, err)

			messageID, err := ids.ToID(messageIDBytes)
			require.NoError(t, err)

			expectedMessageIDBytes := hashing.ComputeHash256(prevUnsignedMsg.Bytes())
			expectedMessageID, err := ids.ToID(expectedMessageIDBytes)
			require.NoError(t, err)
			require.Equal(t, messageID, expectedMessageID)

			signature, err := backend.GetSignature(messageID)
			require.NoError(t, err)

			expectedSig, err := snowCtx.WarpSigner.Sign(prevUnsignedMsg)
			require.NoError(t, err)
			require.Equal(t, signature[:], expectedSig)
		}
	}

	//ensure that there are exactly maxDbSize values
	countIter := backend.countdb.NewIterator()
	entries := 0
	for countIter.Next() {
		entries++
	}
	require.EqualValues(t, entries, maxDbSize)

	msgIter := backend.msgdb.NewIterator()
	entries = 0
	for msgIter.Next() {
		entries++
	}
	require.EqualValues(t, entries, maxDbSize)
}

func TestPruneEntry2(t *testing.T) {
	db := memdb.New()

	snowCtx := snow.DefaultContextTest()
	sk, err := bls.NewSecretKey()
	require.NoError(t, err)
	snowCtx.WarpSigner = avalancheWarp.NewSigner(sk, sourceChainID)

	maxDbSize := 5
	backendConfig := warpBackendConfig{MaxDbSize: uint64(maxDbSize)}

	backend := NewWarpBackend(snowCtx, db, 500).(*warpBackend)
	backend.config = backendConfig

	// Add twice the max db to the db, ensuring that some should get pruned
	for i := 0; i < maxDbSize*2; i++ {
		msg := append(payload, database.PackUInt64(uint64(i))...) //results in test0, test1 etc.
		unsignedMsg, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, msg)
		require.NoError(t, err)
		
		err = backend.AddMessage(unsignedMsg)
		require.NoError(t, err)
		require.EqualValues(t, backend.msgCount, i+1)
	}

	//Go back through all messages that should stay in the db and ensure their presence
	for i := maxDbSize+1; i < maxDbSize*2; i++ {
		msgCountBytes := database.PackUInt64(uint64(i))
		prevMsg := append(payload, msgCountBytes...)
		prevUnsignedMsg, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, prevMsg)
		require.NoError(t, err)

		messageIDBytes, err := backend.countdb.Get(msgCountBytes)
		require.NoError(t, err)

		messageID, err := ids.ToID(messageIDBytes)
		require.NoError(t, err)

		expectedMessageIDBytes := hashing.ComputeHash256(prevUnsignedMsg.Bytes())
		expectedMessageID, err := ids.ToID(expectedMessageIDBytes)
		require.NoError(t, err)
		require.Equal(t, messageID, expectedMessageID)

		signature, err := backend.GetSignature(messageID)
		require.NoError(t, err)

		expectedSig, err := snowCtx.WarpSigner.Sign(prevUnsignedMsg)
		require.NoError(t, err)
		require.Equal(t, signature[:], expectedSig)
	}

	//Go back through messages that should have been deleted, ensure they are not present
	for i := maxDbSize+1; i <= maxDbSize; i++ {
		msgCountBytes := database.PackUInt64(uint64(i))
		prevMsg := append(payload, msgCountBytes...)
		prevUnsignedMsg, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, prevMsg)
		require.NoError(t, err)

		_, err = backend.countdb.Get(msgCountBytes)
		require.ErrorIs(t, err, database.ErrNotFound)

		messageIDBytes := hashing.ComputeHash256(prevUnsignedMsg.Bytes())
		messageID, err := ids.ToID(messageIDBytes)
		require.NoError(t, err)

		_, err = backend.GetSignature(messageID)
		require.ErrorIs(t, err, database.ErrNotFound)
	}

	//ensure that there are exactly maxDbSize values
	countIter := backend.countdb.NewIterator()
	entries := 0
	for countIter.Next() {
		entries++
	}
	require.EqualValues(t, entries, maxDbSize)

	msgIter := backend.msgdb.NewIterator()
	entries = 0
	for msgIter.Next() {
		entries++
	}
	require.EqualValues(t, entries, maxDbSize)
}