// Copyright 2024 the libevm authors.
//
// The libevm additions to go-ethereum are free software: you can redistribute
// them and/or modify them under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.
//
// The libevm additions are distributed in the hope that they will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Lesser
// General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see
// <http://www.gnu.org/licenses/>.

package state

import (
	"sync"

	"github.com/ava-labs/subnet-evm/libevm/options"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

// A PrefetcherOption configures behaviour of trie prefetching.
type PrefetcherOption = options.Option[prefetcherConfig]

type prefetcherConfig struct {
	newWorkers func() WorkerPool
}

// A WorkerPool is responsible for executing functions, possibly asynchronously.
type WorkerPool interface {
	Execute(func())
	Wait()
}

// WithWorkerPools configures trie prefetching to execute asynchronously. The
// provided constructor is called once for each trie being fetched and it MAY
// return the same pool.
func WithWorkerPools(ctor func() WorkerPool) PrefetcherOption {
	return options.Func[prefetcherConfig](func(c *prefetcherConfig) {
		c.newWorkers = ctor
	})
}

type subfetcherPool struct {
	workers WorkerPool
	tries   sync.Pool
}

// applyTo configures the [subfetcher] to use a [WorkerPool] if one was provided
// with a [PrefetcherOption].
func (c *prefetcherConfig) applyTo(sf *subfetcher) {
	sf.pool = &subfetcherPool{
		tries: sync.Pool{
			// Although the workers may be shared between all subfetchers, each
			// MUST have its own Trie pool.
			New: func() any {
				return sf.db.CopyTrie(sf.trie)
			},
		},
	}
	if c.newWorkers != nil {
		sf.pool.workers = c.newWorkers()
	}
}

func (sf *subfetcher) wait() {
	if w := sf.pool.workers; w != nil {
		w.Wait()
	}
}

// execute runs the provided function with a copy of the subfetcher's Trie.
// Copies are stored in a [sync.Pool] to reduce creation overhead. If sf was
// configured with a [WorkerPool] then it is used for function execution,
// otherwise `fn` is just called directly.
func (sf *subfetcher) execute(fn func(Trie)) {
	if w := sf.pool.workers; w != nil {
		w.Execute(func() {
			trie := sf.pool.tries.Get().(Trie)
			fn(trie)
			sf.pool.tries.Put(trie)
		})
	} else {
		trie := sf.pool.tries.Get().(Trie)
		fn(trie)
		sf.pool.tries.Put(trie)
	}
}

// GetAccount optimistically pre-fetches an account, dropping the returned value
// and logging errors. See [subfetcher.execute] re worker pools.
func (sf *subfetcher) GetAccount(addr common.Address) {
	sf.execute(func(t Trie) {
		if _, err := t.GetAccount(addr); err != nil {
			log.Error("account prefetching failed", "address", addr, "err", err)
		}
	})
}

// GetStorage is the storage equivalent of [subfetcher.GetAccount].
func (sf *subfetcher) GetStorage(addr common.Address, key []byte) {
	sf.execute(func(t Trie) {
		if _, err := t.GetStorage(addr, key); err != nil {
			log.Error("storage prefetching failed", "address", addr, "key", key, "err", err)
		}
	})
}
