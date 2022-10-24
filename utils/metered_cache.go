// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/ava-labs/subnet-evm/metrics"
)

// MeteredCache wraps *fastcache.Cache and periodically pulls stats from it.
type MeteredCache struct {
	*fastcache.Cache
	namespace string

	// stats to be surfaced
	entriesCount metrics.Gauge
	bytesSize    metrics.Gauge
	collisions   metrics.Gauge
	gets         metrics.Gauge
	sets         metrics.Gauge
	misses       metrics.Gauge

	// synchronization
	quitCh chan struct{}
	waitWg sync.WaitGroup
}

// NewMeteredCache returns a new MeteredCache that will update stats to the
// provided namespace at the given updateFrequency.
func NewMeteredCache(size int, journal string, namespace string, updateFrequency time.Duration) *MeteredCache {
	var cache *fastcache.Cache
	if journal == "" {
		cache = fastcache.New(size)
	} else {
		cache = fastcache.LoadFromFileOrNew(journal, size)
	}
	mc := &MeteredCache{
		Cache:     cache,
		namespace: namespace,
		quitCh:    make(chan struct{}),
	}
	if namespace != "" {
		// only register stats if a namespace is provided.
		mc.entriesCount = metrics.GetOrRegisterGauge(fmt.Sprintf("%s/entriesCount", namespace), nil)
		mc.bytesSize = metrics.GetOrRegisterGauge(fmt.Sprintf("%s/bytesSize", namespace), nil)
		mc.collisions = metrics.GetOrRegisterGauge(fmt.Sprintf("%s/collisions", namespace), nil)
		mc.gets = metrics.GetOrRegisterGauge(fmt.Sprintf("%s/gets", namespace), nil)
		mc.sets = metrics.GetOrRegisterGauge(fmt.Sprintf("%s/sets", namespace), nil)
		mc.misses = metrics.GetOrRegisterGauge(fmt.Sprintf("%s/misses", namespace), nil)
	}

	if updateFrequency > 0 {
		// spawn a goroutine to periodically update stats from the cache
		mc.waitWg.Add(1)
		go func() {
			defer mc.waitWg.Done()
			ticker := time.NewTicker(updateFrequency)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					mc.updateStats()
				case <-mc.quitCh:
					return
				}
			}
		}()
		// Note: clean up the goroutine once the object is ready for gc.
		runtime.SetFinalizer(mc, func(mc *MeteredCache) { mc.Shutdown() })
	}
	return mc
}

// updateStats updates metrics from fastcache
func (mc *MeteredCache) updateStats() {
	if mc.namespace == "" {
		return
	}

	s := fastcache.Stats{}
	mc.UpdateStats(&s)
	mc.entriesCount.Update(int64(s.EntriesCount))
	mc.bytesSize.Update(int64(s.BytesSize))
	mc.collisions.Update(int64(s.Collisions))
	mc.gets.Update(int64(s.GetCalls))
	mc.sets.Update(int64(s.SetCalls))
	mc.misses.Update(int64(s.Misses))
}

func (mc *MeteredCache) Shutdown() {
	close(mc.quitCh)
	mc.waitWg.Wait()
}
