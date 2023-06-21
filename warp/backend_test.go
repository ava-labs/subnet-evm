// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/database/memdb"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/ava-labs/avalanchego/utils/linkedhashmap"

	//"github.com/ava-labs/avalanchego/utils/linkedhashmap"
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
		msg := make([]byte, rand.Intn(100)+10)
		_, err := rand.Read(msg)
		if err != nil {
			return nil, fmt.Errorf("failed to generate random values: %w", err)
		}

		values = append(values, msg)
	}
	return values, nil
}

// test that duplicate messages are added again correctly
func TestPruneDuplicate(t *testing.T) {
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

const (
	opGetLiving = iota
	opGetDead
	opPutNew
	opPutExisting
	opMax
)

type previousEntry struct {
	//background information about the db
	incr    []byte
	payload []byte
}

type livingEntry struct {
	data          previousEntry
	msgsUntilDead int
}

type livingEntries []livingEntry
type deadEntries []previousEntry

type sessionTracker struct {
	backend WarpBackend
	dbFull  bool
	living  livingEntries
	dead    deadEntries
	mu      sync.RWMutex
}

func NewSessionTracker(backend WarpBackend, living livingEntries, dead deadEntries) sessionTracker {
	linkedhashmap.New[int, []byte]()

	return sessionTracker{
		backend: backend,
		dbFull:  false,
		living:  living,
		dead:    dead,
		mu:      sync.RWMutex{},
	}
}

func (s *sessionTracker) VerifyLivingEntry(i uint, signer avalancheWarp.Signer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	warpdb := s.backend.(*warpBackend).warpdb.(*warpDb)

	//living entry values
	entry := s.living[i].data
	entryincr := entry.incr
	entryUnsignedMessage, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, entry.payload)
	if err != nil {
		return fmt.Errorf("error making unsigned message: %w from living entry record: %s", err, entry.payload)
	}
	entryMessageID, err := ids.ToID(hashing.ComputeHash256(entryUnsignedMessage.Bytes()))
	if err != nil {
		return fmt.Errorf("error converting entry unsigned message %v to message ID: %w", entryUnsignedMessage, err)
	}
	entrySig, err := signer.Sign(entryUnsignedMessage)
	if err != nil {
		return fmt.Errorf("error converting entry unsigned message %v to signature: %w", entryUnsignedMessage, err)
	}

	//Signature from warpdb
	warpSig, err := s.backend.GetSignature(entryMessageID)
	if err != nil {
		return fmt.Errorf("error getting warp signature from warpdb.  Dead: %d, MessageID: %v, Error: %w", s.living[i].msgsUntilDead, entryMessageID, err)
	}

	//get all countdb values
	countDbMessageIDBytes, err := warpdb.countdb.Get(entryincr)
	if err != nil {
		return fmt.Errorf("error fetching messageID from countDB.  incr: %v, Error: %w", entryincr, err)
	}
	countDbMessageID, err := ids.ToID(countDbMessageIDBytes)
	if err != nil {
		return fmt.Errorf("error converting countDbMessageIDBytes to messageID. Messagebytes: %v, Error: %w", countDbMessageIDBytes, err)
	}

	//get all msgdb values
	msgDbKey, err := warpdb.msgdb.Get(entryMessageID[:])
	if err != nil {
		return fmt.Errorf("error fetching msgDbKey from msgDB.  MessageID: %v, Error: %w", entryMessageID, err)
	}
	msgDbincr, msgDbUnsignedMessage := splitIncrement(msgDbKey)

	// compare stored values
	if !bytes.Equal(entryincr, msgDbincr) {
		return fmt.Errorf("entry incr: %v different from stored msgDbincr: %v", entryincr, msgDbincr)
	}
	if countDbMessageID != entryMessageID {
		return fmt.Errorf("entry messageID: %v different from stored countDbMessageID: %v", entryMessageID, countDbMessageID)
	}
	if !bytes.Equal(entryincr, msgDbincr) {
		return fmt.Errorf("entryincr: %v different from stored msgDbincr: %v", entryincr, countDbMessageID)
	}

	if !bytes.Equal(entrySig, warpSig[:]) {
		return fmt.Errorf("entrySig: %v different from stored warpSig: %v", entrySig, warpSig)
	}
	if !bytes.Equal(entryUnsignedMessage.Bytes(), msgDbUnsignedMessage) {
		return fmt.Errorf("entryUnsignedMessage %v different from stored msgDbUnsignedMessage: %v", entryUnsignedMessage, msgDbUnsignedMessage)
	}

	return nil
}

func (s *sessionTracker) VerifyDeadEntry(i uint, signer avalancheWarp.Signer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	warpdb := s.backend.(*warpBackend).warpdb.(*warpDb)
	entry := s.dead[i]
	entryUnsignedMessage, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, entry.payload)
	if err != nil {
		return fmt.Errorf("error making unsigned message: %w from living entry record: %s", err, entry.payload)
	}
	entryMessageID, err := ids.ToID(hashing.ComputeHash256(entryUnsignedMessage.Bytes()))
	_, err = s.backend.GetSignature(entryMessageID)
	if err == nil {
		return fmt.Errorf("old message should have been deleted: %v", entry.payload)
	}

	if !errors.Is(database.ErrNotFound, err) {
		return fmt.Errorf("error getting signature: %w", err)
	}

	has, err := warpdb.countdb.Has(entry.incr)
	if has {
		return fmt.Errorf("old message count not deleted.  incr: %v", entry.incr)
	}
	if err != nil {
		return fmt.Errorf("error reading countdb:  %w", err)
	}

	return nil
}

func (s *sessionTracker) addNewLiving(payload []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	warpdb := s.backend.(*warpBackend).warpdb.(*warpDb)
	incr := database.PackUInt64(warpdb.incrementer) //increment value at the current state

	unsignedMessage, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, payload)
	if err != nil {
		return fmt.Errorf("error making new unsigned message: %w", err)
	}
	s.backend.AddMessage(unsignedMessage)
	if err != nil {
		return fmt.Errorf("error adding unsignedMessage to db: %w", err)
	}
	s.updateDeadCounter(-1)
	livingEntry := livingEntry{
		data: previousEntry{
			incr:    incr,
			payload: payload,
		},
		msgsUntilDead: int(warpdb.size),
	}
	s.living = append(s.living, livingEntry)
	return nil
}

func (s *sessionTracker) addOldLiving(i uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	warpdb := s.backend.(*warpBackend).warpdb.(*warpDb)
	incr := database.PackUInt64(warpdb.incrementer)
	entry := s.living[i].data

	unsignedMessage, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, entry.payload)
	if err != nil {
		return fmt.Errorf("error making new unsigned message: %w", err)
	}

	s.backend.AddMessage(unsignedMessage)
	if err != nil {
		return fmt.Errorf("error adding unsignedMessage to db: %w", err)
	}

	s.living[i].msgsUntilDead = int(warpdb.size)
	s.living[i].data.incr = incr
	return nil

}

func (s *sessionTracker) updateDeadCounter(i int) {
	for index, entry := range s.living {
		if entry.msgsUntilDead <= -i {
			s.dead = append(s.dead, entry.data)
			s.living = append(s.living[:index], s.living[index+1:]...)
		} else {
			s.living[index].msgsUntilDead += i
		}
	}
}

func Test(t *testing.T) {

	t.Logf("%d", database.PackUInt64(uint64(time.Now().Unix())))
	t.Fail()

}

func (s *sessionTracker) changeDbSize(i uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	warpdb := s.backend.(*warpBackend).warpdb.(*warpDb)
	changeAmt := i - warpdb.size
	warpdb.size = i
	s.updateDeadCounter(int(changeAmt))
}

func FuzzTestDb(f *testing.F) {
	db := memdb.New()
	snowCtx := snow.DefaultContextTest()
	sk, err := bls.NewSecretKey()
	require.NoError(f, err)
	snowCtx.WarpSigner = avalancheWarp.NewSigner(sk, sourceChainID)

	maxDbSize := 10000
	backend := NewWarpBackend(snowCtx, db, 0)
	backend.(*warpBackend).warpdb.(*warpDb).size = uint64(maxDbSize)

	livingEntries := make(livingEntries, 0)
	deadEntries := make(deadEntries, 0)

	tracker := NewSessionTracker(backend, livingEntries, deadEntries)

	f.Add([]byte("test"), uint(2), uint(10))
	f.Fuzz(func(t *testing.T, input []byte, op uint, index uint) {

		op = op % opMax
		t.Logf("Op: %d", op)

		switch op {
		case opGetLiving:
			t.Log("getting living")

			if len(tracker.living) > 0 {
				t.Logf("length of tracker: %d", len(tracker.living))
				err := tracker.VerifyLivingEntry(uint(rand.Intn(len(tracker.living))), snowCtx.WarpSigner)
				require.NoError(t, err)
			}

		case opGetDead:
			t.Log("getting dead")
			if len(tracker.dead) > 0 {
				err := tracker.VerifyDeadEntry(index%uint(len(tracker.dead)), snowCtx.WarpSigner)
				require.NoError(t, err)
			}
		case opPutNew:
			t.Log("putting new")
			
			err := tracker.addNewLiving(input)
			require.NoError(t, err)
			err = tracker.VerifyLivingEntry(uint(len(tracker.living)-1), snowCtx.WarpSigner)
			require.NoError(t, err)
		case opPutExisting:

			t.Log("putting existing")
			if len(tracker.living) > 0 {
				err := tracker.addOldLiving(index % uint(len(tracker.living)))
				require.NoError(t, err)
				err = tracker.VerifyLivingEntry(uint(len(tracker.living)-1), snowCtx.WarpSigner)
			}
		}
	})
}

// potential update to testaddandgetvalidmessage
func TestEntryAdditionNoPruning(t *testing.T) {
	db := memdb.New()
	snowCtx := snow.DefaultContextTest()
	sk, err := bls.NewSecretKey()
	require.NoError(t, err)
	snowCtx.WarpSigner = avalancheWarp.NewSigner(sk, sourceChainID)

	maxDbSize := 5
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

	maxDbSize := 5

	backend := NewWarpBackend(snowCtx, db, 0).(*warpBackend)
	backenddb := backend.warpdb.(*warpDb)
	backenddb.size = uint64(maxDbSize)

	values, err := GetRandomValues(maxDbSize * 2)
	require.NoError(t, err)

	// Add twice the max db to the db, ensuring that some should get pruned
	for i := 0; i < maxDbSize*2; i++ {
		unsignedMsg, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, values[i])
		require.NoError(t, err)

		err = backend.AddMessage(unsignedMsg)
		require.NoError(t, err)
	}

	//Go back through all messages that should stay in the db and ensure their presence
	for i := maxDbSize; i < maxDbSize*2; i++ {
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
