package warp

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/database/memdb"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/ava-labs/avalanchego/utils/linkedhashmap"

	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"

	"github.com/stretchr/testify/require"
)

type op uint64

const (
	opAdd op = iota
	opGet
	opPrune
	opMax //not an actual op, used for utility purposes
)

const maxMessageSize = 1000

type mockWarpDb struct {
	mem linkedhashmap.LinkedHashmap[int, []byte]
}

func NewMockWarpDb() mockWarpDb {
	return mockWarpDb{
		linkedhashmap.New[int, []byte](),
	}
}

func removeFromMockWarp(mwd mockWarpDb, m prunedMap) error {
	iter := m.NewIterator()
	for iter.Next() {
		key := iter.Key()
		_, has := mwd.mem.Get(key)
		if !has {
			return fmt.Errorf("mock db doesn't contain: key: %d, value: %s", key, iter.Value())
		}

		mwd.mem.Delete(key)
	}

	return nil
}

func addToMockWarp(mwd mockWarpDb, keyBytes []byte, value []byte) error {
	key, err := database.ParseUInt64(keyBytes)
	if err != nil {
		return err
	}
	mwd.mem.Put(int(key), value)
	
	return nil
}

func MakeRandomMessage(maxMessageSize uint) ([]byte, error) {
	msg := make([]byte, rand.Intn(int(maxMessageSize)))
	_, err := rand.Read(msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func FuzzTestWarpDb(f *testing.F) {

	f.Add(uint(50), int64(1337), Minute, uint(10), true)
	f.Fuzz(func(t *testing.T, numOps uint, randomSeed int64, threshold uint64, maxPruneSize uint, autoprune bool) {
		rand.Seed(randomSeed)

		db := memdb.New()
		snowCtx := snow.DefaultContextTest()
		sk, err := bls.NewSecretKey()
		require.NoError(t, err)
		snowCtx.WarpSigner = avalancheWarp.NewSigner(sk, sourceChainID)

		config := warpDbConfig {
			autoprune,
			maxPruneSize,
		}
		warpDb := NewWarpDb(db, threshold, config)
		mockWarpDb := NewMockWarpDb()

		for i := uint(0) ; i < numOps; i++ {
			op := op(rand.Intn(int(opMax)))

			switch op {
			case opAdd:
				msg, err:= MakeRandomMessage(maxMessageSize)
				require.NoError(t, err)
				unsignedMsg, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, msg)
				messageIDBytes := hashing.ComputeHash256Array(unsignedMsg.Bytes())

				_, prunedMap, err := warpDb.AddMessage(unsignedMsg)
				require.NoError(t, err)

				err = removeFromMockWarp(mockWarpDb, prunedMap)
				require.NoError(t, err)

				ts, _, err := warpDb.GetUnsignedMessage(messageIDBytes)
				require.NoError(t, err)

				err = addToMockWarp(mockWarpDb, ts, messageIDBytes[:])
				require.NoError(t, err)
			case opGet:
				if mockWarpDb.mem.Len() > 0 {
					index := rand.Intn(mockWarpDb.mem.Len())
					iter := mockWarpDb.mem.NewIterator()
					for index >= 0 {
						iter.Next()
						index--
					}
					messageID, err := ids.ToID(iter.Value())
					require.NoError(t, err)

					ts, _, err := warpDb.GetUnsignedMessage(messageID)
					require.NoError(t, err)

					require.Equal(t, ts, database.PackUInt64(uint64(iter.Key())))
				}
			case opPrune:
				_, prunedMap, err := warpDb.PruneEntries(maxPruneSize)
				require.NoError(t, err)

				err = removeFromMockWarp(mockWarpDb, prunedMap)
				require.NoError(t, err)
			}
		}
	})
}

func TestDbPutSafe(t *testing.T) {
	db := memdb.New()

	key1 := []byte("key1")
	key2 := []byte("key2")
	value1 := []byte("value1")
	value2 := []byte("value2")

	err := dBPutSafe(db, key1, value1)
	require.NoError(t, err)
	err = dBPutSafe(db, key2, value2)
	require.NoError(t, err)
	err = dBPutSafe(db, key1, value1)
	require.Equal(t, err, keyExistsError)
}

func TestPrependTimestamp(t *testing.T) {
	tsBytes := database.PackUInt64(256)
	testUnsignedMessage := []byte("test")
	expectedMessage := []byte{0, 0, 0, 0, 0, 0, 1, 0, 116, 101, 115, 116}

	msgEntry := prependTimestamp(tsBytes, testUnsignedMessage)
	require.Equal(t, expectedMessage, msgEntry)
}

func TestSplitTimestamp(t *testing.T) {
	msgEntry := []byte{0, 0, 0, 0, 0, 0, 1, 0, 116, 101, 115, 116}
	expectedTs := database.PackUInt64(256)
	expectedUnsignedMessage := []byte("test")

	ts, unsignedMessage := splitTimestamp(msgEntry)

	require.Equal(t, expectedTs, ts)
	require.Equal(t, unsignedMessage, expectedUnsignedMessage)
}

func TestAutoPrune(t *testing.T) {
	db := memdb.New()

	timeout := Second

	config := GetDefaultWarpDbConfig()
	warpDb := NewWarpDb(db, timeout, config)

	msg, err := GetRandomValues(2)
	require.NoError(t, err)

	unsignedMsg0, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, msg[0])
	require.NoError(t, err)
	expectedMessageID0 := hashing.ComputeHash256Array(unsignedMsg0.Bytes())

	unsignedMsg1, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, msg[1])
	require.NoError(t, err)
	expectedMessageID1 := hashing.ComputeHash256Array(unsignedMsg1.Bytes())

	_, _, err = warpDb.AddMessage(unsignedMsg0)
	require.NoError(t, err)

	_, returnedUnsignedMsg0, err := warpDb.GetUnsignedMessage(expectedMessageID0)
	require.NoError(t, err)
	require.Equal(t, unsignedMsg0.Bytes(), returnedUnsignedMsg0)

	time.Sleep(time.Duration(timeout))

	_, _, err = warpDb.AddMessage(unsignedMsg1)
	require.NoError(t, err)

	_, returnUnsignedMsg1, err := warpDb.GetUnsignedMessage(expectedMessageID1)
	require.Equal(t, unsignedMsg1.Bytes(), returnUnsignedMsg1)

	_, v, err := warpDb.GetUnsignedMessage(expectedMessageID0)
	t.Logf("%v", v)
	require.ErrorIs(t, err, database.ErrNotFound)
}

