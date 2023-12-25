// (c) 2021-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package statesync

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"runtime/pprof"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/database/memdb"
	"github.com/ava-labs/avalanchego/x/merkledb"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethdb"
	"github.com/ava-labs/subnet-evm/ethdb/memorydb"
	"github.com/ava-labs/subnet-evm/plugin/evm/message"
	"github.com/ava-labs/subnet-evm/plugin/mdb"
	statesyncclient "github.com/ava-labs/subnet-evm/sync/client"
	"github.com/ava-labs/subnet-evm/sync/handlers"
	handlerstats "github.com/ava-labs/subnet-evm/sync/handlers/stats"
	"github.com/ava-labs/subnet-evm/sync/syncutils"
	"github.com/ava-labs/subnet-evm/trie"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSyncTimeout = 30 * time.Second

var errInterrupted = errors.New("interrupted sync")

type syncTest struct {
	ctx               context.Context
	prepareForTest    func(t *testing.T) (clientDB ethdb.Database, serverDB ethdb.Database, serverTrieDB *trie.Database, syncRoot common.Hash)
	expectedError     error
	GetLeafsIntercept func(message.LeafsRequest, message.LeafsResponse) (message.LeafsResponse, error)
	GetCodeIntercept  func([]common.Hash, [][]byte) ([][]byte, error)
	onFinish          func(ethdb.Database)
}

type bothDBs struct {
	ethdb.Database
	db2 ethdb.Database
}

func newClientDB() *bothDBs {
	db := memorydb.New()
	merkleDB, err := merkledb.New(context.Background(), memdb.New(), mdb.NewBasicConfig())
	if err != nil {
		panic(err)
	}
	return &bothDBs{
		Database: memorydb.New(),
		db2:      mdb.NewWithMerkleDB(db, merkleDB, nil),
	}
}

func testSync(t *testing.T, test syncTest) {
	t.Helper()
	ctx := context.Background()
	if test.ctx != nil {
		ctx = test.ctx
	}
	clientDB, serverDB, serverTrieDB, root := test.prepareForTest(t)
	if both, ok := clientDB.(*bothDBs); ok {
		// note this blocks in a different goroutine then this test
		// continues with the memorydb (previous behavior)
		t.Run(
			fmt.Sprintf("%s_MerkleDB", t.Name()),
			func(t *testing.T) {
				syncTestCopy := test
				syncTestCopy.prepareForTest = func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
					return both.db2, serverDB, serverTrieDB, root
				}
				testSync(t, syncTestCopy)
			},
		)
	}

	leafsRequestHandler := handlers.NewLeafsRequestHandler(serverTrieDB, nil, message.Codec, handlerstats.NewNoopHandlerStats())
	codeRequestHandler := handlers.NewCodeRequestHandler(serverDB, message.Codec, handlerstats.NewNoopHandlerStats())
	mockClient := statesyncclient.NewMockClient(message.Codec, leafsRequestHandler, codeRequestHandler, nil)
	// Set intercept functions for the mock client
	mockClient.GetLeafsIntercept = test.GetLeafsIntercept
	mockClient.GetCodeIntercept = test.GetCodeIntercept

	s, err := NewStateSyncer(&StateSyncerConfig{
		Client:                   mockClient,
		Root:                     root,
		DB:                       clientDB,
		BatchSize:                1000, // Use a lower batch size in order to get test coverage of batches being written early.
		NumCodeFetchingWorkers:   DefaultNumCodeFetchingWorkers,
		MaxOutstandingCodeHashes: DefaultMaxOutstandingCodeHashes,
		RequestSize:              1024,
	})
	if err != nil {
		t.Fatal(err)
	}
	// begin sync
	s.Start(ctx)
	waitFor(t, s.Done(), test.expectedError, testSyncTimeout)
	if test.onFinish != nil {
		defer test.onFinish(clientDB)
	}
	if test.expectedError != nil {
		return
	}

	assertDBConsistency(t, root, clientDB, serverTrieDB, trie.NewDatabase(clientDB))
}

// waitFor waits for a result on the [result] channel to match [expected], or a timeout.
func waitFor(t *testing.T, result <-chan error, expected error, timeout time.Duration) {
	t.Helper()
	select {
	case err := <-result:
		if expected != nil {
			if err == nil {
				t.Fatalf("Expected error %s, but got nil", expected)
			}
			assert.Contains(t, err.Error(), expected.Error())
		} else if err != nil {
			t.Fatal("unexpected error waiting for sync result", err)
		}
	case <-time.After(timeout):
		// print a stack trace to assist with debugging
		// if the test times out.
		var stackBuf bytes.Buffer
		pprof.Lookup("goroutine").WriteTo(&stackBuf, 2)
		t.Log(stackBuf.String())
		// fail the test
		t.Fatal("unexpected timeout waiting for sync result")
	}
}

func TestSimpleSyncCases(t *testing.T) {
	rand := rand.New(rand.NewSource(1))
	var (
		numAccounts      = 250
		numAccountsSmall = 10
		clientErr        = errors.New("dummy client error")
	)
	tests := map[string]syncTest{
		"accounts": {
			prepareForTest: func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
				serverDB := memorydb.New()
				serverTrieDB := trie.NewDatabase(serverDB)
				root, _ := trie.FillAccounts(t, rand, serverTrieDB, common.Hash{}, numAccounts, nil)
				return newClientDB(), serverDB, serverTrieDB, root
			},
		},
		"accounts with code": {
			prepareForTest: func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
				serverDB := memorydb.New()
				serverTrieDB := trie.NewDatabase(serverDB)
				root, _ := trie.FillAccounts(t, rand, serverTrieDB, common.Hash{}, numAccounts, func(t *testing.T, index int, account types.StateAccount) types.StateAccount {
					if index%3 == 0 {
						codeBytes := make([]byte, 256)
						_, err := rand.Read(codeBytes)
						if err != nil {
							t.Fatalf("error reading random code bytes: %v", err)
						}

						codeHash := crypto.Keccak256Hash(codeBytes)
						rawdb.WriteCode(serverDB, codeHash, codeBytes)
						account.CodeHash = codeHash[:]
					}
					return account
				})
				return newClientDB(), serverDB, serverTrieDB, root
			},
		},
		"accounts with code and storage": {
			prepareForTest: func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
				serverDB := memorydb.New()
				serverTrieDB := trie.NewDatabase(serverDB)
				root := fillAccountsWithStorage(t, rand, serverDB, serverTrieDB, common.Hash{}, numAccounts)
				return newClientDB(), serverDB, serverTrieDB, root
			},
		},
		"accounts with storage": {
			prepareForTest: func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
				serverDB := memorydb.New()
				serverTrieDB := trie.NewDatabase(serverDB)
				root, _ := trie.FillAccounts(t, rand, serverTrieDB, common.Hash{}, numAccounts, func(t *testing.T, i int, account types.StateAccount) types.StateAccount {
					if i%5 == 0 {
						account.Root, _, _ = trie.GenerateTrie(t, rand, serverTrieDB, 16, common.HashLength)
					}

					return account
				})
				return newClientDB(), serverDB, serverTrieDB, root
			},
		},
		"accounts with overlapping storage": {
			prepareForTest: func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
				serverDB := memorydb.New()
				serverTrieDB := trie.NewDatabase(serverDB)
				root, _ := FillAccountsWithOverlappingStorage(t, rand, serverTrieDB, common.Hash{}, numAccounts, 3)
				return newClientDB(), serverDB, serverTrieDB, root
			},
		},
		"failed to fetch leafs": {
			prepareForTest: func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
				serverDB := memorydb.New()
				serverTrieDB := trie.NewDatabase(serverDB)
				root, _ := trie.FillAccounts(t, rand, serverTrieDB, common.Hash{}, numAccountsSmall, nil)
				return newClientDB(), serverDB, serverTrieDB, root
			},
			GetLeafsIntercept: func(_ message.LeafsRequest, _ message.LeafsResponse) (message.LeafsResponse, error) {
				return message.LeafsResponse{}, clientErr
			},
			expectedError: clientErr,
		},
		"failed to fetch code": {
			prepareForTest: func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
				serverDB := memorydb.New()
				serverTrieDB := trie.NewDatabase(serverDB)
				root := fillAccountsWithStorage(t, rand, serverDB, serverTrieDB, common.Hash{}, numAccountsSmall)
				return newClientDB(), serverDB, serverTrieDB, root
			},
			GetCodeIntercept: func(_ []common.Hash, _ [][]byte) ([][]byte, error) {
				return nil, clientErr
			},
			expectedError: clientErr,
		},
	}
	for name, test := range tests {
		rand.Seed(1)
		t.Run(name, func(t *testing.T) {
			testSync(t, test)
		})
	}
}

func TestCancelSync(t *testing.T) {
	rand := rand.New(rand.NewSource(1))
	serverDB := memorydb.New()
	serverTrieDB := trie.NewDatabase(serverDB)
	// Create trie with 2000 accounts (more than one leaf request)
	root := fillAccountsWithStorage(t, rand, serverDB, serverTrieDB, common.Hash{}, 2000)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	testSync(t, syncTest{
		ctx: ctx,
		prepareForTest: func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
			return newClientDB(), serverDB, serverTrieDB, root
		},
		expectedError: context.Canceled,
		GetLeafsIntercept: func(_ message.LeafsRequest, lr message.LeafsResponse) (message.LeafsResponse, error) {
			cancel()
			return lr, nil
		},
	})
}

// interruptLeafsIntercept provides the parameters to the getLeafsIntercept
// function which returns [errInterrupted] after passing through [numRequests]
// leafs requests for [root].
type interruptLeafsIntercept struct {
	numRequests    uint32
	interruptAfter uint32
	root           common.Hash
}

// getLeafsIntercept can be passed to mockClient and returns an unmodified
// response for the first [numRequest] requests for leafs from [root].
// After that, all requests for leafs from [root] return [errInterrupted].
func (i *interruptLeafsIntercept) getLeafsIntercept(request message.LeafsRequest, response message.LeafsResponse) (message.LeafsResponse, error) {
	if request.Root == i.root {
		if numRequests := atomic.AddUint32(&i.numRequests, 1); numRequests > i.interruptAfter {
			return message.LeafsResponse{}, errInterrupted
		}
	}
	return response, nil
}

func TestResumeSyncAccountsTrieInterrupted(t *testing.T) {
	rand := rand.New(rand.NewSource(1))
	serverDB := memorydb.New()
	serverTrieDB := trie.NewDatabase(serverDB)
	root, _ := FillAccountsWithOverlappingStorage(t, rand, serverTrieDB, common.Hash{}, 2000, 3)
	clientDB := newClientDB()
	intercept := &interruptLeafsIntercept{
		root:           root,
		interruptAfter: 1,
	}
	testSync(t, syncTest{
		prepareForTest: func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
			return clientDB, serverDB, serverTrieDB, root
		},
		expectedError:     errInterrupted,
		GetLeafsIntercept: intercept.getLeafsIntercept,
		onFinish: func(clientDB ethdb.Database) {
			assert.EqualValues(t, 2, intercept.numRequests)
			intercept.numRequests = 0
		},
	})

	testSync(t, syncTest{
		prepareForTest: func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
			return clientDB, serverDB, serverTrieDB, root
		},
	})
}

func TestResumeSyncLargeStorageTrieInterrupted(t *testing.T) {
	rand := rand.New(rand.NewSource(1))
	serverDB := memorydb.New()
	serverTrieDB := trie.NewDatabase(serverDB)

	largeStorageRoot, _, _ := trie.GenerateTrie(t, rand, serverTrieDB, 2000, common.HashLength)
	root, _ := trie.FillAccounts(t, rand, serverTrieDB, common.Hash{}, 2000, func(t *testing.T, index int, account types.StateAccount) types.StateAccount {
		// Set the root for a single account
		if index == 10 {
			account.Root = largeStorageRoot
		}
		return account
	})
	clientDB := newClientDB()
	intercept := &interruptLeafsIntercept{
		root:           largeStorageRoot,
		interruptAfter: 1,
	}
	testSync(t, syncTest{
		prepareForTest: func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
			return clientDB, serverDB, serverTrieDB, root
		},
		expectedError:     errInterrupted,
		GetLeafsIntercept: intercept.getLeafsIntercept,
		onFinish: func(clientDB ethdb.Database) {
			intercept.numRequests = 0
		},
	})

	testSync(t, syncTest{
		prepareForTest: func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
			return clientDB, serverDB, serverTrieDB, root
		},
	})
}

func TestResumeSyncToNewRootAfterLargeStorageTrieInterrupted(t *testing.T) {
	rand := rand.New(rand.NewSource(1))
	serverDB := memorydb.New()
	serverTrieDB := trie.NewDatabase(serverDB)

	largeStorageRoot1, _, _ := trie.GenerateTrie(t, rand, serverTrieDB, 2000, common.HashLength)
	largeStorageRoot2, _, _ := trie.GenerateTrie(t, rand, serverTrieDB, 2000, common.HashLength)
	root1, _ := trie.FillAccounts(t, rand, serverTrieDB, common.Hash{}, 2000, func(t *testing.T, index int, account types.StateAccount) types.StateAccount {
		// Set the root for a single account
		if index == 10 {
			account.Root = largeStorageRoot1
		}
		return account
	})
	root2, _ := trie.FillAccounts(t, rand, serverTrieDB, root1, 100, func(t *testing.T, index int, account types.StateAccount) types.StateAccount {
		if index == 20 {
			account.Root = largeStorageRoot2
		}
		return account
	})
	clientDB := newClientDB()
	intercept := &interruptLeafsIntercept{
		root:           largeStorageRoot1,
		interruptAfter: 1,
	}
	testSync(t, syncTest{
		prepareForTest: func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
			return clientDB, serverDB, serverTrieDB, root1
		},
		expectedError:     errInterrupted,
		GetLeafsIntercept: intercept.getLeafsIntercept,
		onFinish: func(clientDB ethdb.Database) {
			intercept.numRequests = 0
			err := syncutils.ClearPartialDB(clientDB)
			require.NoError(t, err)
		},
	})

	testSync(t, syncTest{
		prepareForTest: func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
			return clientDB, serverDB, serverTrieDB, root2
		},
	})
}

func TestResumeSyncLargeStorageTrieWithConsecutiveDuplicatesInterrupted(t *testing.T) {
	rand := rand.New(rand.NewSource(1))
	serverDB := memorydb.New()
	serverTrieDB := trie.NewDatabase(serverDB)

	largeStorageRoot, _, _ := trie.GenerateTrie(t, rand, serverTrieDB, 2000, common.HashLength)
	root, _ := trie.FillAccounts(t, rand, serverTrieDB, common.Hash{}, 100, func(t *testing.T, index int, account types.StateAccount) types.StateAccount {
		// Set the root for 2 successive accounts
		if index == 10 || index == 11 {
			account.Root = largeStorageRoot
		}
		return account
	})
	clientDB := newClientDB()
	intercept := &interruptLeafsIntercept{
		root:           largeStorageRoot,
		interruptAfter: 1,
	}
	testSync(t, syncTest{
		prepareForTest: func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
			return clientDB, serverDB, serverTrieDB, root
		},
		expectedError:     errInterrupted,
		GetLeafsIntercept: intercept.getLeafsIntercept,
	})

	testSync(t, syncTest{
		prepareForTest: func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
			return clientDB, serverDB, serverTrieDB, root
		},
	})
}

func TestResumeSyncLargeStorageTrieWithSpreadOutDuplicatesInterrupted(t *testing.T) {
	rand := rand.New(rand.NewSource(1))
	serverDB := memorydb.New()
	serverTrieDB := trie.NewDatabase(serverDB)

	largeStorageRoot, _, _ := trie.GenerateTrie(t, rand, serverTrieDB, 2000, common.HashLength)
	root, _ := trie.FillAccounts(t, rand, serverTrieDB, common.Hash{}, 100, func(t *testing.T, index int, account types.StateAccount) types.StateAccount {
		if index == 10 || index == 90 {
			account.Root = largeStorageRoot
		}
		return account
	})
	clientDB := newClientDB()
	intercept := &interruptLeafsIntercept{
		root:           largeStorageRoot,
		interruptAfter: 1,
	}
	testSync(t, syncTest{
		prepareForTest: func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
			return clientDB, serverDB, serverTrieDB, root
		},
		expectedError:     errInterrupted,
		GetLeafsIntercept: intercept.getLeafsIntercept,
	})

	testSync(t, syncTest{
		prepareForTest: func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
			return clientDB, serverDB, serverTrieDB, root
		},
	})
}

func TestResyncNewRootAfterDeletes(t *testing.T) {
	for name, test := range map[string]struct {
		deleteBetweenSyncs func(*testing.T, common.Hash, ethdb.Database)
	}{
		"delete code": {
			deleteBetweenSyncs: func(t *testing.T, _ common.Hash, clientDB ethdb.Database) {
				// delete code
				it := clientDB.NewIterator(rawdb.CodePrefix, nil)
				defer it.Release()
				for it.Next() {
					if len(it.Key()) != len(rawdb.CodePrefix)+common.HashLength {
						continue
					}
					if err := clientDB.Delete(it.Key()); err != nil {
						t.Fatal(err)
					}
				}
				if err := it.Error(); err != nil {
					t.Fatal(err)
				}
			},
		},
		"delete intermediate storage nodes": {
			deleteBetweenSyncs: func(t *testing.T, root common.Hash, clientDB ethdb.Database) {
				clientTrieDB := trie.NewDatabase(clientDB)
				tr, err := trie.NewStateTrie(trie.TrieID(root), clientTrieDB)
				if err != nil {
					t.Fatal(err)
				}
				it := trie.NewIterator(tr.NodeIterator(nil))
				accountsWithStorage := 0

				// keep track of storage tries we delete trie nodes from
				// so we don't try to do it again if another account has
				// the same storage root.
				corruptedStorageRoots := make(map[common.Hash]struct{})
				for it.Next() {
					var acc types.StateAccount
					if err := rlp.DecodeBytes(it.Value, &acc); err != nil {
						t.Fatal(err)
					}
					if acc.Root == types.EmptyRootHash {
						continue
					}
					if _, found := corruptedStorageRoots[acc.Root]; found {
						// avoid trying to delete nodes from a trie we have already deleted nodes from
						continue
					}
					accountsWithStorage++
					if accountsWithStorage%2 != 0 {
						continue
					}
					corruptedStorageRoots[acc.Root] = struct{}{}
					trie.CorruptTrie(t, clientTrieDB, trie.StorageTrieID(root, common.Hash(it.Key), acc.Root), 2)
				}
				if err := it.Err; err != nil {
					t.Fatal(err)
				}
			},
		},
		"delete intermediate account trie nodes": {
			deleteBetweenSyncs: func(t *testing.T, root common.Hash, clientDB ethdb.Database) {
				clientTrieDB := trie.NewDatabase(clientDB)
				trie.CorruptTrie(t, clientTrieDB, trie.StateTrieID(root), 5)
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			testSyncerSyncsToNewRoot(t, test.deleteBetweenSyncs)
		})
	}
}

func testSyncerSyncsToNewRoot(t *testing.T, deleteBetweenSyncs func(*testing.T, common.Hash, ethdb.Database)) {
	rand := rand.New(rand.NewSource(1))
	clientDB := newClientDB()
	serverDB := memorydb.New()
	serverTrieDB := trie.NewDatabase(serverDB)

	root1, _ := FillAccountsWithOverlappingStorage(t, rand, serverTrieDB, common.Hash{}, 1000, 3)
	root2, _ := FillAccountsWithOverlappingStorage(t, rand, serverTrieDB, root1, 1000, 3)

	testSync(t, syncTest{
		prepareForTest: func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
			return clientDB, serverDB, serverTrieDB, root1
		},
		onFinish: func(clientDB ethdb.Database) {
			deleteBetweenSyncs(t, root1, clientDB)

			// delete snapshot first since this is not the responsibility of the EVM State Syncer
			err := syncutils.ClearPartialDB(clientDB)
			require.NoError(t, err)
		},
	})
	testSync(t, syncTest{
		prepareForTest: func(t *testing.T) (ethdb.Database, ethdb.Database, *trie.Database, common.Hash) {
			return clientDB, serverDB, serverTrieDB, root2
		},
	})
}
