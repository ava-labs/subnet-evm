// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/ava-labs/subnet-evm/metrics"
)

var (
	statsMap = map[string]func(s *fastcache.Stats) uint64{
		"entries":    func(s *fastcache.Stats) uint64 { return s.EntriesCount },
		"bytesSize":  func(s *fastcache.Stats) uint64 { return s.BytesSize },
		"collisions": func(s *fastcache.Stats) uint64 { return s.Collisions },
		"gets":       func(s *fastcache.Stats) uint64 { return s.GetCalls },
		"sets":       func(s *fastcache.Stats) uint64 { return s.SetCalls },
		"misses":     func(s *fastcache.Stats) uint64 { return s.Misses },
	}
)

// MeteredCache wraps *fastcache.Cache and periodically pulls stats from it.
type MeteredCache struct {
	*fastcache.Cache
	stats  map[string]metrics.Gauge
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

	stats := make(map[string]metrics.Gauge, len(statsMap))
	if namespace != "" {
		// avoid registering stats if a namespace was not provided.
		for statName := range statsMap {
			stats[statName] = metrics.GetOrRegisterGauge(fmt.Sprintf("%s/%s", namespace, statName), nil)
		}
	}
	mc := &MeteredCache{
		Cache:  cache,
		stats:  stats,
		quitCh: make(chan struct{}),
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
					break
				}
			}
		}()
	}
	return mc
}

// updateStats updates metrics from fastcache
func (mc *MeteredCache) updateStats() {
	s := fastcache.Stats{}
	mc.UpdateStats(&s)
	for statName, stat := range mc.stats {
		stat.Update(int64(statsMap[statName](&s)))
	}
}

func (mc *MeteredCache) Shutdown() {
	close(mc.quitCh)
	mc.waitWg.Wait()
}
