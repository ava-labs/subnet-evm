// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/database/memdb"

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

func GetRandomValues(n int) ([][]byte, error) {
	values := [][]byte{}
	for i := 0; i < n; i++ {
		msg := make([]byte, rand.Intn(100))
		_, err := rand.Read(msg)
		if err != nil {
			return nil, fmt.Errorf("failed to generate random values: %w", err)
		}

		values = append(values, msg)
	}
	return values, nil
}

//test that duplicate messages are added again correctly
func TestPruneDuplicate(t * testing.T) {
	db := memdb.New()
	snowCtx := snow.DefaultContextTest()
	sk, err := bls.NewSecretKey()
	require.NoError(t, err)
	snowCtx.WarpSigner = avalancheWarp.NewSigner(sk, sourceChainID)

	maxDbSize := 5
	backend := NewWarpBackend(snowCtx, db, 0).(*warpBackend)
	backenddb := backend.warpdb.(*warpDb)
	backenddb.size = uint64(maxDbSize)

	msg := []byte("test")
	unsignedMsg, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, msg)
	backend.AddMessage(unsignedMsg)
	backend.AddMessage(unsignedMsg)

	_, err = backenddb.countdb.Get([]byte(database.PackUInt64(1)))
	require.NoError(t, err)

	_, err = backenddb.countdb.Get([]byte(database.PackUInt64(0)))
	require.Error(t, database.ErrNotFound)
	require.EqualValues(t, backenddb.count, 1)

}

func FuzzTestDb(f *testing.F) {
	f.Add()
}

//potential update to testaddandgetvalidmessage
func TestEntryAdditionNoPruning(t *testing.T) {
	db := memdb.New()
	snowCtx := snow.DefaultContextTest()
	sk, err := bls.NewSecretKey()
	require.NoError(t, err)
	snowCtx.WarpSigner = avalancheWarp.NewSigner(sk, sourceChainID)

	maxDbSize := 10
	backend := NewWarpBackend(snowCtx, db, 0).(*warpBackend)
	backenddb := backend.warpdb.(*warpDb)
	backenddb.size = uint64(maxDbSize)

	values, err := GetRandomValues(maxDbSize)
	require.NoError(t, err)

	// Create a new unsigned message and add it to the warp backend.

	for i := 0; i < maxDbSize; i++ {
		unsignedMsg, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, values[i])
		require.NoError(t, err)
		
		err = backend.AddMessage(unsignedMsg)
		require.NoError(t, err)
		require.EqualValues(t, backenddb.count, i+1)

		//Go back through all messages that were added, ensure nothing was deleted
		for j := 0; j <= i; j++ {
			countBytes := database.PackUInt64(uint64(j))
			messageIDBytes, err := backenddb.countdb.Get(countBytes)
			require.NoError(t, err)

			messageID, err := ids.ToID(messageIDBytes)

			prevUnsignedMsg, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, values[j])
			require.NoError(t, err)

			expectedMessageIDBytes := hashing.ComputeHash256(prevUnsignedMsg.Bytes())
			expectedMessageID, err := ids.ToID(expectedMessageIDBytes)
			require.NoError(t, err)
			require.Equal(t, messageID, expectedMessageID)

			signature, err := backend.GetSignature(expectedMessageID)
			require.NoError(t, err)

			expectedSig, err := snowCtx.WarpSigner.Sign(prevUnsignedMsg)
			require.NoError(t, err)
			require.Equal(t, signature[:], expectedSig)
		}
	}

	//ensure that there are exactly maxDbSize values
	countIter := backenddb.countdb.NewIterator()
	entries := 0
	for countIter.Next() {
		entries++
	}
	require.EqualValues(t, entries, maxDbSize)

	msgIter := backenddb.msgdb.NewIterator()
	entries = 0
	for msgIter.Next() {
		entries++
	}
	require.EqualValues(t, entries, maxDbSize)
}

func TestEntryAdditionPruning(t *testing.T) {
	db := memdb.New()

	snowCtx := snow.DefaultContextTest()
	sk, err := bls.NewSecretKey()
	require.NoError(t, err)
	snowCtx.WarpSigner = avalancheWarp.NewSigner(sk, sourceChainID)

	maxDbSize := 20

	backend := NewWarpBackend(snowCtx, db, 0).(*warpBackend)
	backenddb := backend.warpdb.(*warpDb)
	backenddb.size = uint64(maxDbSize)

	values, err := GetRandomValues(maxDbSize*2)
	require.NoError(t, err)

	// Add twice the max db to the db, ensuring that some should get pruned
	for i := 0; i < maxDbSize*2; i++ {
		unsignedMsg, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, values[i])
		require.NoError(t, err)
		
		err = backend.AddMessage(unsignedMsg)
		require.NoError(t, err)
	}

	//Go back through all messages that should stay in the db and ensure their presence
	for i := maxDbSize+1; i < maxDbSize*2; i++ {
		countBytes := database.PackUInt64(uint64(i))
		prevUnsignedMsg, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, values[i])
		require.NoError(t, err)

		messageIDBytes, err := backenddb.countdb.Get(countBytes)
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
	for i := 0; i < maxDbSize; i++ {
		countBytes := database.PackUInt64(uint64(i))
		prevUnsignedMsg, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, values[i])
		require.NoError(t, err)

		_, err = backenddb.countdb.Get(countBytes)
		require.ErrorIs(t, err, database.ErrNotFound)

		messageIDBytes := hashing.ComputeHash256(prevUnsignedMsg.Bytes())
		messageID, err := ids.ToID(messageIDBytes)
		require.NoError(t, err)

		_, err = backend.GetSignature(messageID)
		require.ErrorIs(t, err, database.ErrNotFound)
	}

	//ensure that there are exactly maxDbSize values
	countIter := backenddb.countdb.NewIterator()
	entries := 0
	for countIter.Next() {
		entries++
	}
	require.EqualValues(t, entries, maxDbSize)

	msgIter := backenddb.msgdb.NewIterator()
	entries = 0
	for msgIter.Next() {
		entries++
	}
	require.EqualValues(t, entries, maxDbSize)
	require.EqualValues(t, backenddb.count, entries)
}