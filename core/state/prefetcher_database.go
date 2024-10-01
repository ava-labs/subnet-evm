// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package state

import (
	"fmt"
	"sync"

	"github.com/ava-labs/coreth/core/types"
	"github.com/ava-labs/coreth/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

// ForPrefetchingOnly returns a new database that is only suitable for prefetching
// operations. It will not be safe to use for any other operations.
// Close must be called on the returned database when it is no longer needed
// to wait on all spawned goroutines.
func (*cachingDB) ForPrefetchingOnly(db Database, maxConcurrency int) Database {
	return newPrefetcherDatabase(db, maxConcurrency)
}

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
	t, err := p.Database.OpenTrie(root)
	return newPrefetcherTrie(p, t), err
}

func (p *prefetcherDatabase) OpenStorageTrie(stateRoot common.Hash, address common.Address, root common.Hash, trie Trie) (Trie, error) {
	t, err := p.Database.OpenStorageTrie(stateRoot, address, root, trie)
	return newPrefetcherTrie(p, t), err
}

func (p *prefetcherDatabase) CopyTrie(t Trie) Trie {
	switch t := t.(type) {
	case *prefetcherTrie:
		return newPrefetcherTrie(p, t.getCopy())
	default:
		panic(fmt.Errorf("unknown trie type %T", t))
	}
}

func (p *prefetcherDatabase) Close() {
	p.workers.Wait()
}

type prefetcherTrie struct {
	p *prefetcherDatabase

	Trie
	copyLock sync.Mutex

	copies chan Trie
	wg     sync.WaitGroup
}

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

func (p *prefetcherTrie) putCopy(copy Trie) {
	select {
	case p.copies <- copy:
	default:
	}
}

func (p *prefetcherTrie) GetAccount(address common.Address) (*types.StateAccount, error) {
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
	return nil, nil // Note this result is never used by the prefetcher
}

func (p *prefetcherTrie) GetStorage(address common.Address, key []byte) ([]byte, error) {
	p.wg.Add(1)
	f := func() {
		defer p.wg.Done()

		tr := p.getCopy()
		_, err := tr.GetStorage(address, key)
		if err != nil {
			log.Error("GetAccount failed in prefetcher", "err", err)
		}
		p.putCopy(tr)
	}
	p.p.workers.Execute(f)
	return nil, nil // Note this result is never used by the prefetcher
}
