// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"fmt"

	"github.com/ava-labs/avalanchego/cache"
	"github.com/ava-labs/subnet-evm/metrics"
	"github.com/ethereum/go-ethereum/log"
)

const (
	avgObjSize = 512 // 512B
)

// MeteredCache wraps *fastcache.Cache and periodically pulls stats from it.
type MeteredCache struct {
	cache     *cache.LRU
	namespace string

	// stats to be surfaced
	puts   metrics.Counter
	hits   metrics.Counter
	misses metrics.Counter
}

// NewMeteredCache returns a new MeteredCache that will update stats to the
// provided namespace once per each [updateFrequency] operations.
// Note: if [updateFrequency] is passed as 0, it will be treated as 1.
func NewMeteredCache(size int, namespace string) *MeteredCache {
	estSize := size / avgObjSize
	if estSize < 4096 {
		estSize = 4096
	}
	log.Info("Creating cache", "namespace", namespace, "entries", estSize)
	mc := &MeteredCache{
		cache:     &cache.LRU{Size: estSize},
		namespace: namespace,
	}
	if namespace != "" {
		// only register stats if a namespace is provided.
		mc.puts = metrics.NewRegisteredCounter(fmt.Sprintf("%s/puts", namespace), nil)
		mc.hits = metrics.NewRegisteredCounter(fmt.Sprintf("%s/hits", namespace), nil)
		mc.misses = metrics.NewRegisteredCounter(fmt.Sprintf("%s/misses", namespace), nil)
	}
	return mc
}

func (mc *MeteredCache) Del(k []byte) {
	mc.cache.Evict(string(k))
}

func (mc *MeteredCache) incMiss() {
	if mc.misses == nil {
		return
	}
	mc.misses.Inc(1)
}

func (mc *MeteredCache) incHit() {
	if mc.hits == nil {
		return
	}
	mc.hits.Inc(1)
}

func (mc *MeteredCache) incPuts() {
	if mc.puts == nil {
		return
	}
	mc.puts.Inc(1)
}

func (mc *MeteredCache) Get(k []byte) []byte {
	v, ok := mc.cache.Get(string(k))
	if !ok {
		mc.incMiss()
		return nil
	}
	mc.incHit()
	return v.([]byte)
}

func (mc *MeteredCache) Set(k, v []byte) {
	mc.incPuts()
	mc.cache.Put(string(k), v)
}
