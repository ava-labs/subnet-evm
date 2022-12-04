// (c) 2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package core

import (
	"sync"
)

// FIFOCache evicts the oldest element added to it after [limit] items are
// added.
type FIFOCache[K comparable, V any] struct {
	limit int

	l sync.RWMutex

	buffer *BoundedBuffer[K]
	m      map[K]V
}

func NewFIFOCache[K comparable, V any](limit int) *FIFOCache[K, V] {
	c := &FIFOCache[K, V]{
		limit: limit,
		m:     map[K]V{},
	}
	c.buffer = NewBoundedBuffer[K](limit, c.remove)
	return c
}

// remove is used as the callback in [BoundedBuffer]. It is assumed that the
// [WriteLock] is held when this is accessed.
func (f *FIFOCache[K, V]) remove(key K) {
	delete(f.m, key)
}

func (f *FIFOCache[K, V]) Put(key K, val V) {
	f.l.Lock()
	defer f.l.Unlock()

	f.buffer.Insert(key) // Insert will remove the oldest [K] if we are at the [limit]
	f.m[key] = val
}

func (f *FIFOCache[K, V]) Get(key K) (V, bool) {
	f.l.RLock()
	defer f.l.RUnlock()

	v, ok := f.m[key]
	return v, ok
}
