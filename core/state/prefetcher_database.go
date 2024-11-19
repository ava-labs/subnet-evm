// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package state

import (
	"sync"

	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

// PrefetcherDB is an interface that extends Database with additional methods
// used in trie_prefetcher.  This includes specific methods for prefetching
// accounts and storage slots, (which may be non-blocking and/or parallelized)
// and methods to wait for pending prefetches.
type PrefetcherDB interface {
	// From Database
	OpenTrie(root common.Hash) (Trie, error)
	OpenStorageTrie(stateRoot common.Hash, address common.Address, root common.Hash, trie Trie) (Trie, error)
	CopyTrie(t Trie) Trie

	// Additional methods
	PrefetchAccount(t Trie, address common.Address)
	PrefetchStorage(t Trie, address common.Address, key []byte)
	CanPrefetchDuringShutdown() bool
	WaitTrie(t Trie)
	Close()
}

// withPrefetcher is an optional interface that a Database can implement to
// signal PrefetcherDB() should be called to get a Database for use in
// trie_prefetcher.  Each call to PrefetcherDB() should return a new
// PrefetcherDB instance.
type withPrefetcherDB interface {
	PrefetcherDB() PrefetcherDB
}

type withPrefetcher struct {
	Database
	maxConcurrency int
}

func (db *withPrefetcher) PrefetcherDB() PrefetcherDB {
	return newPrefetcherDatabase(db.Database, db.maxConcurrency)
}

func WithPrefetcher(db Database, maxConcurrency int) Database {
	return &withPrefetcher{db, maxConcurrency}
}

// withPrefetcherDefaults extends Database and implements PrefetcherDB by adding
// default implementations for PrefetchAccount and PrefetchStorage that read the
// account and storage slot from the trie.
type withPrefetcherDefaults struct {
	Database
}

func (withPrefetcherDefaults) PrefetchAccount(t Trie, address common.Address) {
	_, _ = t.GetAccount(address)
}

func (withPrefetcherDefaults) PrefetchStorage(t Trie, address common.Address, key []byte) {
	_, _ = t.GetStorage(address, key)
}

func (withPrefetcherDefaults) CanPrefetchDuringShutdown() bool { return false }
func (withPrefetcherDefaults) WaitTrie(Trie)                   {}
func (withPrefetcherDefaults) Close()                          {}

type prefetcherDatabase struct {
	Database

	maxConcurrency int
	workers        *utils.BoundedWorkers
}

func newPrefetcherDatabase(db Database, maxConcurrency int) *prefetcherDatabase {
	return &prefetcherDatabase{
		Database:       db,
		maxConcurrency: maxConcurrency,
		workers:        utils.NewBoundedWorkers(maxConcurrency),
	}
}

func (p *prefetcherDatabase) OpenTrie(root common.Hash) (Trie, error) {
	trie, err := p.Database.OpenTrie(root)
	return newPrefetcherTrie(p, trie), err
}

func (p *prefetcherDatabase) OpenStorageTrie(stateRoot common.Hash, address common.Address, root common.Hash, trie Trie) (Trie, error) {
	storageTrie, err := p.Database.OpenStorageTrie(stateRoot, address, root, trie)
	return newPrefetcherTrie(p, storageTrie), err
}

func (p *prefetcherDatabase) CopyTrie(t Trie) Trie {
	switch t := t.(type) {
	case *prefetcherTrie:
		return t.getCopy()
	default:
		return p.Database.CopyTrie(t)
	}
}

// PrefetchAccount should only be called on a trie returned from OpenTrie or OpenStorageTrie
func (*prefetcherDatabase) PrefetchAccount(t Trie, address common.Address) {
	t.(*prefetcherTrie).PrefetchAccount(address)
}

// PrefetchStorage should only be called on a trie returned from OpenTrie or OpenStorageTrie
func (*prefetcherDatabase) PrefetchStorage(t Trie, address common.Address, key []byte) {
	t.(*prefetcherTrie).PrefetchStorage(address, key)
}

// WaitTrie should only be called on a trie returned from OpenTrie or OpenStorageTrie
func (*prefetcherDatabase) WaitTrie(t Trie) {
	t.(*prefetcherTrie).Wait()
}

func (p *prefetcherDatabase) Close() {
	p.workers.Wait()
}

func (p *prefetcherDatabase) CanPrefetchDuringShutdown() bool {
	return true
}

type prefetcherTrie struct {
	p *prefetcherDatabase

	Trie
	copyLock sync.Mutex

	copies chan Trie
	wg     sync.WaitGroup
}

// newPrefetcherTrie returns a new prefetcherTrie that wraps the given trie.
// prefetcherTrie prefetches accounts and storage slots in parallel, using
// bounded workers from the prefetcherDatabase.  As Trie is not safe for
// concurrent access, each prefetch operation uses a copy. The copy is kept in
// a buffered channel for reuse.
func newPrefetcherTrie(p *prefetcherDatabase, t Trie) *prefetcherTrie {
	prefetcher := &prefetcherTrie{
		p:      p,
		Trie:   t,
		copies: make(chan Trie, p.maxConcurrency),
	}
	prefetcher.copies <- prefetcher.getCopy()
	return prefetcher
}

func (p *prefetcherTrie) Wait() {
	p.wg.Wait()
}

// getCopy returns a copy of the trie. The copy is taken from the copies channel
// if available, otherwise a new copy is created.
func (p *prefetcherTrie) getCopy() Trie {
	select {
	case copy := <-p.copies:
		return copy
	default:
		p.copyLock.Lock()
		defer p.copyLock.Unlock()
		return p.p.Database.CopyTrie(p.Trie)
	}
}

// putCopy keeps the copy for future use.  If the buffer is full, the copy is
// discarded.
func (p *prefetcherTrie) putCopy(copy Trie) {
	select {
	case p.copies <- copy:
	default:
	}
}

func (p *prefetcherTrie) PrefetchAccount(address common.Address) {
	p.wg.Add(1)
	f := func() {
		defer p.wg.Done()

		tr := p.getCopy()
		_, err := tr.GetAccount(address)
		if err != nil {
			log.Error("GetAccount failed in prefetcher", "err", err)
		}
		p.putCopy(tr)
	}
	p.p.workers.Execute(f)
}

func (p *prefetcherTrie) PrefetchStorage(address common.Address, key []byte) {
	p.wg.Add(1)
	f := func() {
		defer p.wg.Done()

		tr := p.getCopy()
		_, err := tr.GetStorage(address, key)
		if err != nil {
			log.Error("GetStorage failed in prefetcher", "err", err)
		}
		p.putCopy(tr)
	}
	p.p.workers.Execute(f)
}
