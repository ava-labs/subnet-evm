// (c) 2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package core

import "github.com/ava-labs/avalanchego/utils/linkedhashmap"

// FIFOCache evicts the oldest element added to it after [limit] items are
// added.
//
// TODO: Move this to AvalancheGo
type FIFOCache[K, V any] struct {
	limit int

	m linkedhashmap.LinkedHashmap[K, V]

	// [Linkedhashmap] is thread-safe, so no additional locking is required
}

func NewFIFOCache[K comparable, V any](limit int) *FIFOCache[K, V] {
	return &FIFOCache[K, V]{
		limit: limit,
		m:     linkedhashmap.New[K, V](),
	}
}

func (f *FIFOCache[K, V]) Put(key K, val V) {
	if f.m.Len() >= f.limit {
		oldest, _, exists := f.m.Oldest()
		if exists {
			f.m.Delete(oldest)
		}
	}
	f.m.Put(key, val)
}

func (f *FIFOCache[K, V]) Get(key K) (V, bool) {
	return f.m.Get(key)
}
