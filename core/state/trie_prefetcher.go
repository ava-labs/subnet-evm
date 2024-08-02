// Copyright 2020 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package state

import (
	"errors"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/utils"
)

var (
	// triePrefetchMetricsPrefix is the prefix under which to publish the metrics.
	triePrefetchMetricsPrefix = "trie/prefetch/"

	// errTerminated is returned if a fetcher is attempted to be operated after it
	// has already terminated.
	errTerminated = errors.New("fetcher is already terminated")
)

// triePrefetcher is an active prefetcher, which receives accounts or storage
// items and does trie-loading of them. The goal is to get as much useful content
// into the caches as possible.
//
// Note, the prefetcher's API is not thread safe.
type triePrefetcher struct {
	db       Database               // Database to fetch trie nodes through
	root     common.Hash            // Root hash of the account trie for metrics
	fetchers map[string]*subfetcher // Subfetchers for each trie
	term     chan struct{}          // Channel to signal interruption
	noreads  bool                   // Whether to ignore state-read-only prefetch requests

	// added by avalanche
	maxConcurrency          int
	workers                 *utils.BoundedWorkers
	subfetcherWorkersMeter  metrics.Meter
	subfetcherWaitTimer     metrics.Counter
	subfetcherCopiesMeter   metrics.Meter
	storageFetchersMeter    metrics.Meter
	storageLargestLoadMeter metrics.Meter

	accountLoadReadMeter  metrics.Meter
	accountLoadWriteMeter metrics.Meter
	accountDupReadMeter   metrics.Meter
	accountDupWriteMeter  metrics.Meter
	accountDupCrossMeter  metrics.Meter
	accountWasteMeter     metrics.Meter

	storageLoadReadMeter  metrics.Meter
	storageLoadWriteMeter metrics.Meter
	storageDupReadMeter   metrics.Meter
	storageDupWriteMeter  metrics.Meter
	storageDupCrossMeter  metrics.Meter
	storageWasteMeter     metrics.Meter
}

func newTriePrefetcher(db Database, root common.Hash, namespace string, noreads bool, maxConcurrency int) *triePrefetcher {
	prefix := triePrefetchMetricsPrefix + namespace
	return &triePrefetcher{
		db:       db,
		root:     root,
		fetchers: make(map[string]*subfetcher), // Active prefetchers use the fetchers map
		term:     make(chan struct{}),
		noreads:  noreads,

		maxConcurrency: maxConcurrency,
		workers:        utils.NewBoundedWorkers(maxConcurrency), // Scale up as needed to [maxConcurrency]

		subfetcherWorkersMeter: metrics.GetOrRegisterMeter(prefix+"/subfetcher/workers", nil),
		subfetcherWaitTimer:    metrics.GetOrRegisterCounter(prefix+"/subfetcher/wait", nil),
		subfetcherCopiesMeter:  metrics.GetOrRegisterMeter(prefix+"/subfetcher/copies", nil),

		storageFetchersMeter:    metrics.GetOrRegisterMeter(prefix+"/storage/fetchers", nil),
		storageLargestLoadMeter: metrics.GetOrRegisterMeter(prefix+"/storage/lload", nil),

		accountLoadReadMeter:  metrics.GetOrRegisterMeter(prefix+"/account/load/read", nil),
		accountLoadWriteMeter: metrics.GetOrRegisterMeter(prefix+"/account/load/write", nil),
		accountDupReadMeter:   metrics.GetOrRegisterMeter(prefix+"/account/dup/read", nil),
		accountDupWriteMeter:  metrics.GetOrRegisterMeter(prefix+"/account/dup/write", nil),
		accountDupCrossMeter:  metrics.GetOrRegisterMeter(prefix+"/account/dup/cross", nil),
		accountWasteMeter:     metrics.GetOrRegisterMeter(prefix+"/account/waste", nil),

		storageLoadReadMeter:  metrics.GetOrRegisterMeter(prefix+"/storage/load/read", nil),
		storageLoadWriteMeter: metrics.GetOrRegisterMeter(prefix+"/storage/load/write", nil),
		storageDupReadMeter:   metrics.GetOrRegisterMeter(prefix+"/storage/dup/read", nil),
		storageDupWriteMeter:  metrics.GetOrRegisterMeter(prefix+"/storage/dup/write", nil),
		storageDupCrossMeter:  metrics.GetOrRegisterMeter(prefix+"/storage/dup/cross", nil),
		storageWasteMeter:     metrics.GetOrRegisterMeter(prefix+"/storage/waste", nil),
	}
}

// terminate iterates over all the subfetchers and issues a termination request
// to all of them. Depending on the async parameter, the method will either block
// until all subfetchers spin down, or return immediately.
func (p *triePrefetcher) terminate(async bool) {
	// Collect stats from all fetchers
	var (
		storageFetchers int64
		largestLoad     int64
	)

	// Short circuit if the fetcher is already closed
	select {
	case <-p.term:
		return
	default:
	}
	// Terminate all sub-fetchers, sync or async, depending on the request
	for _, fetcher := range p.fetchers {
		if metrics.Enabled {
			p.subfetcherCopiesMeter.Mark(int64(fetcher.copies()))

			if fetcher.root != p.root {
				storageFetchers++
				oseen := int64(len(fetcher.seenRead) + len(fetcher.seenWrite)) // XXX: Is this metric useful?
				if oseen > largestLoad {
					largestLoad = oseen
				}
			}
		}
		// XXX: untangle this loop from the metrics if possible
		fetcher.terminate(async)
	}
	if metrics.Enabled {
		p.storageFetchersMeter.Mark(storageFetchers)
		p.storageLargestLoadMeter.Mark(largestLoad)
	}

	// Stop all workers once fetchers are aborted (otherwise
	// could stop while waiting)
	//
	// Record number of workers that were spawned during this run
	workersUsed := int64(p.workers.Wait())
	if metrics.Enabled {
		p.subfetcherWorkersMeter.Mark(workersUsed)
	}

	close(p.term)
}

// report aggregates the pre-fetching and usage metrics and reports them.
func (p *triePrefetcher) report() {
	if !metrics.Enabled {
		return
	}
	for _, fetcher := range p.fetchers {
		fetcher.wait() // ensure the fetcher's idle before poking in its internals

		if fetcher.root == p.root {
			p.accountLoadReadMeter.Mark(int64(len(fetcher.seenRead)))
			p.accountLoadWriteMeter.Mark(int64(len(fetcher.seenWrite)))

			p.accountDupReadMeter.Mark(int64(fetcher.dupsRead))
			p.accountDupWriteMeter.Mark(int64(fetcher.dupsWrite))
			p.accountDupCrossMeter.Mark(int64(fetcher.dupsCross))

			for _, key := range fetcher.used {
				delete(fetcher.seenRead, string(key))
				delete(fetcher.seenWrite, string(key))
			}
			p.accountWasteMeter.Mark(int64(len(fetcher.seenRead) + len(fetcher.seenWrite)))
		} else {
			p.storageLoadReadMeter.Mark(int64(len(fetcher.seenRead)))
			p.storageLoadWriteMeter.Mark(int64(len(fetcher.seenWrite)))

			p.storageDupReadMeter.Mark(int64(fetcher.dupsRead))
			p.storageDupWriteMeter.Mark(int64(fetcher.dupsWrite))
			p.storageDupCrossMeter.Mark(int64(fetcher.dupsCross))

			for _, key := range fetcher.used {
				delete(fetcher.seenRead, string(key))
				delete(fetcher.seenWrite, string(key))
			}
			p.storageWasteMeter.Mark(int64(len(fetcher.seenRead) + len(fetcher.seenWrite)))
		}
	}
}

// prefetch schedules a batch of trie items to prefetch. After the prefetcher is
// closed, all the following tasks scheduled will not be executed and an error
// will be returned.
//
// prefetch is called from two locations:
//
//  1. Finalize of the state-objects storage roots. This happens at the end
//     of every transaction, meaning that if several transactions touches
//     upon the same contract, the parameters invoking this method may be
//     repeated.
//  2. Finalize of the main account trie. This happens only once per block.
func (p *triePrefetcher) prefetch(owner common.Hash, root common.Hash, addr common.Address, keys [][]byte, read bool) error {
	// If the state item is only being read, but reads are disabled, return
	if read && p.noreads {
		return nil
	}
	// Ensure the subfetcher is still alive
	select {
	case <-p.term:
		return errTerminated
	default:
	}
	id := p.trieID(owner, root)
	fetcher := p.fetchers[id]
	if fetcher == nil {
		fetcher = newSubfetcher(p, owner, root, addr)
		p.fetchers[id] = fetcher
	}
	fetcher.schedule(keys, read)
	return nil
}

// trie returns the trie matching the root hash, blocking until the fetcher of
// the given trie terminates. If no fetcher exists for the request, nil will be
// returned.
func (p *triePrefetcher) trie(owner common.Hash, root common.Hash) Trie {
	// Bail if no trie was prefetched for this root
	fetcher := p.fetchers[p.trieID(owner, root)]
	if fetcher == nil {
		log.Error("Prefetcher missed to load trie", "owner", owner, "root", root)
		return nil
	}

	// XXX: verify this code is still correct
	// Wait for the fetcher to finish and shutdown orchestrator, if it exists
	start := time.Now()
	fetcher.wait()
	if metrics.Enabled {
		p.subfetcherWaitTimer.Inc(time.Since(start).Milliseconds())
	}

	// Subfetcher exists, retrieve its trie
	return fetcher.peek()
}

// used marks a batch of state items used to allow creating statistics as to
// how useful or wasteful the fetcher is.
func (p *triePrefetcher) used(owner common.Hash, root common.Hash, used [][]byte) {
	if fetcher := p.fetchers[p.trieID(owner, root)]; fetcher != nil {
		fetcher.wait() // ensure the fetcher's idle before poking in its internals
		fetcher.used = used
	}
}

// trieID returns an unique trie identifier consists the trie owner and root hash.
func (p *triePrefetcher) trieID(owner common.Hash, root common.Hash) string {
	trieID := make([]byte, common.HashLength*2)
	copy(trieID, owner.Bytes())
	copy(trieID[common.HashLength:], root.Bytes())
	return string(trieID)
}

// subfetcher is a trie fetcher goroutine responsible for pulling entries for a
// single trie. It is spawned when a new root is encountered and lives until the
// main prefetcher is paused and either all requested items are processed or if
// the trie being worked on is retrieved from the prefetcher.
type subfetcher struct {
	p *triePrefetcher

	db    Database       // Database to load trie nodes through
	state common.Hash    // Root hash of the state to prefetch
	owner common.Hash    // Owner of the trie, usually account hash
	root  common.Hash    // Root hash of the trie to prefetch
	addr  common.Address // Address of the account that the trie belongs to

	to *trieOrchestrator // Orchestrate concurrent fetching of a single trie

	seenRead  map[string]struct{} // Tracks the entries already loaded via read operations
	seenWrite map[string]struct{} // Tracks the entries already loaded via write operations

	dupsRead  int // Number of duplicate preload tasks via reads only
	dupsWrite int // Number of duplicate preload tasks via writes only
	dupsCross int // Number of duplicate preload tasks via read-write-crosses

	used [][]byte // Tracks the entries used in the end
}

// subfetcherTask is a trie path to prefetch, tagged with whether it originates
// from a read or a write request.
type subfetcherTask struct {
	read bool
	key  []byte
}

// newSubfetcher creates a goroutine to prefetch state items belonging to a
// particular root hash.
func newSubfetcher(p *triePrefetcher, owner common.Hash, root common.Hash, addr common.Address) *subfetcher {
	sf := &subfetcher{
		p:         p,
		db:        p.db,
		state:     p.root,
		owner:     owner,
		root:      root,
		addr:      addr,
		seenRead:  make(map[string]struct{}),
		seenWrite: make(map[string]struct{}),
	}
	sf.to = newTrieOrchestrator(sf)
	if sf.to != nil {
		go sf.to.processTasks()
	}
	// We return [sf] here to ensure we don't try to re-create if
	// we aren't able to setup a [newTrieOrchestrator] the first time.
	return sf
}

// schedule adds a batch of trie keys to the queue to prefetch.
// This should never block, so an array is used instead of a channel.
//
// This is not thread-safe.
func (sf *subfetcher) schedule(keys [][]byte, read bool) {
	// Append the tasks to the current queue
	tasks := make([]*subfetcherTask, 0, len(keys))
	for _, key := range keys {
		// Check if keys already seen
		sk := string(key)
		if read {
			if _, ok := sf.seenRead[sk]; ok {
				sf.dupsRead++
				continue
			}
			if _, ok := sf.seenWrite[sk]; ok {
				sf.dupsCross++
				continue
			}
		} else {
			if _, ok := sf.seenRead[sk]; ok {
				sf.dupsCross++
				continue
			}
			if _, ok := sf.seenWrite[sk]; ok {
				sf.dupsWrite++
				continue
			}
		}
		if read {
			sf.seenRead[sk] = struct{}{}
		} else {
			sf.seenWrite[sk] = struct{}{}
		}
		key := key // closure for the append below
		tasks = append(tasks, &subfetcherTask{read: read, key: key})
	}

	// After counting keys, exit if they can't be prefetched
	if sf.to == nil {
		return
	}

	// Add tasks to queue for prefetching
	sf.to.enqueueTasks(tasks)
}

// peek retrieves the fetcher's trie, populated with any pre-fetched data. The
// returned trie will be a shallow copy, so modifying it will break subsequent
// peeks for the original data. The method will block until all the scheduled
// data has been loaded and the fethcer terminated.
func (sf *subfetcher) peek() Trie {
	if sf.to == nil {
		return nil
	}
	return sf.to.copyBase()
}

// wait must only be called if [triePrefetcher] has not been closed. If this happens,
// workers will not finish.
func (sf *subfetcher) wait() {
	if sf.to == nil {
		// Unable to open trie
		return
	}
	sf.to.wait()
}

// terminate requests the subfetcher to stop accepting new tasks and spin down
// as soon as everything is loaded. Depending on the async parameter, the method
// will either block until all disk loads finish or return immediately.
func (sf *subfetcher) terminate(async bool) {
	if sf.to == nil {
		// Unable to open trie
		return
	}
	sf.to.abort(async)
}

func (sf *subfetcher) skips() int {
	if sf.to == nil {
		// Unable to open trie
		return 0
	}
	return sf.to.skipCount()
}

func (sf *subfetcher) copies() int {
	if sf.to == nil {
		// Unable to open trie
		return 0
	}
	return sf.to.copies
}

// trieOrchestrator is not thread-safe.
type trieOrchestrator struct {
	sf *subfetcher

	// base is an unmodified Trie we keep for
	// creating copies for each worker goroutine.
	//
	// We care more about quick copies than good copies
	// because most (if not all) of the nodes that will be populated
	// in the copy will come from the underlying triedb cache. Ones
	// that don't come from this cache probably had to be fetched
	// from disk anyways.
	base     Trie
	baseLock sync.Mutex

	tasksAllowed bool
	skips        int // number of tasks skipped
	pendingTasks []*subfetcherTask
	taskLock     sync.Mutex

	processingTasks sync.WaitGroup

	wake     chan struct{}
	stop     chan struct{}
	stopOnce sync.Once
	loopTerm chan struct{}

	copies      int
	copyChan    chan Trie // XXX: seems the copy chan was removed in upstream, check if it should be removed here
	copySpawner chan struct{}
}

func newTrieOrchestrator(sf *subfetcher) *trieOrchestrator {
	// Start by opening the trie and stop processing if it fails
	var (
		base Trie
		err  error
	)
	if sf.owner == (common.Hash{}) {
		base, err = sf.db.OpenTrie(sf.root)
		if err != nil {
			log.Warn("Trie prefetcher failed opening trie", "root", sf.root, "err", err)
			return nil
		}
	} else {
		base, err = sf.db.OpenStorageTrie(sf.state, sf.addr, sf.root, nil)
		if err != nil {
			log.Warn("Trie prefetcher failed opening trie", "root", sf.root, "err", err)
			return nil
		}
	}

	// Instantiate trieOrchestrator
	to := &trieOrchestrator{
		sf:   sf,
		base: base,

		tasksAllowed: true,
		wake:         make(chan struct{}, 1),
		stop:         make(chan struct{}),
		loopTerm:     make(chan struct{}),

		copyChan:    make(chan Trie, sf.p.maxConcurrency),
		copySpawner: make(chan struct{}, sf.p.maxConcurrency),
	}

	// Create initial trie copy
	to.copies++
	to.copySpawner <- struct{}{}
	to.copyChan <- to.copyBase()
	return to
}

func (to *trieOrchestrator) copyBase() Trie {
	to.baseLock.Lock()
	defer to.baseLock.Unlock()

	return to.sf.db.CopyTrie(to.base) // XXX: does this need to be a deep copy?
}

func (to *trieOrchestrator) skipCount() int {
	to.taskLock.Lock()
	defer to.taskLock.Unlock()

	return to.skips
}

func (to *trieOrchestrator) enqueueTasks(tasks []*subfetcherTask) {
	to.taskLock.Lock()
	defer to.taskLock.Unlock()

	if len(tasks) == 0 {
		return
	}

	// Add tasks to [pendingTasks]
	if !to.tasksAllowed {
		to.skips += len(tasks)
		return
	}
	to.processingTasks.Add(len(tasks))
	to.pendingTasks = append(to.pendingTasks, tasks...)

	// Wake up processor
	select {
	case to.wake <- struct{}{}:
	default:
	}
}

func (to *trieOrchestrator) handleStop(remaining int) {
	to.taskLock.Lock()
	to.skips += remaining
	to.taskLock.Unlock()
	to.processingTasks.Add(-remaining)
}

func (to *trieOrchestrator) processTasks() {
	defer close(to.loopTerm)

	for {
		// Determine if we should process or exit
		select {
		case <-to.wake:
		case <-to.stop:
			return
		}

		// Get current tasks
		to.taskLock.Lock()
		tasks := to.pendingTasks
		to.pendingTasks = nil
		to.taskLock.Unlock()

		// Enqueue more work as soon as trie copies are available
		lt := len(tasks)
		for i := 0; i < lt; i++ {
			// Try to stop as soon as possible, if channel is closed
			remaining := lt - i
			select {
			case <-to.stop:
				to.handleStop(remaining)
				return
			default:
			}

			// Try to create to get an active copy first (select is non-deterministic,
			// so we may end up creating a new copy when we don't need to)
			var t Trie
			select {
			case t = <-to.copyChan:
			default:
				// Wait for an available copy or create one, if we weren't
				// able to get a previously created copy
				select {
				case <-to.stop:
					to.handleStop(remaining)
					return
				case t = <-to.copyChan:
				case to.copySpawner <- struct{}{}:
					to.copies++
					t = to.copyBase()
				}
			}

			// Enqueue work, unless stopped.
			fTask := tasks[i]
			f := func() {
				// XXX: double check
				var (
					task = fTask
					err  error
				)
				// Perform task
				if len(task.key) == common.AddressLength {
					_, err = t.GetAccount(common.BytesToAddress(task.key))
				} else {
					_, err = t.GetStorage(to.sf.addr, task.key)
				}
				if err != nil {
					log.Error("Trie prefetcher failed fetching", "root", to.sf.root, "err", err)
				}
				to.processingTasks.Done()

				// Return copy when we are done with it, so someone else can use it
				//
				// channel is buffered and will not block
				to.copyChan <- t
			}

			// Enqueue task for processing (may spawn new goroutine
			// if not at [maxConcurrency])
			//
			// If workers are stopped before calling [Execute], this function may
			// panic.
			to.sf.p.workers.Execute(f)
		}
	}
}

func (to *trieOrchestrator) stopAcceptingTasks() {
	to.taskLock.Lock()
	defer to.taskLock.Unlock()

	if !to.tasksAllowed {
		return
	}
	to.tasksAllowed = false

	// We don't clear [to.pendingTasks] here because
	// it will be faster to prefetch them even though we
	// are still waiting.
}

// wait stops accepting new tasks and waits for ongoing tasks to complete. If
// wait is called, it is not necessary to call [abort].
//
// It is safe to call wait multiple times.
func (to *trieOrchestrator) wait() {
	// Prevent more tasks from being enqueued
	to.stopAcceptingTasks()

	// Wait for processing tasks to complete
	to.processingTasks.Wait()

	// Stop orchestrator loop
	to.stopOnce.Do(func() {
		close(to.stop)
	})
	<-to.loopTerm
}

// abort stops any ongoing tasks and shuts down the orchestrator loop. If abort
// is called, it is not necessary to call [wait].
//
// It is safe to call abort multiple times.
func (to *trieOrchestrator) abort(async bool) {
	// Prevent more tasks from being enqueued
	to.stopAcceptingTasks()

	// Stop orchestrator loop
	to.stopOnce.Do(func() {
		close(to.stop)
	})
	if async {
		// XXX: is this necessary/correct?
		// return
	}
	<-to.loopTerm

	// Capture any dangling pending tasks (processTasks
	// may exit before enqueing all pendingTasks)
	to.taskLock.Lock()
	pendingCount := len(to.pendingTasks)
	to.skips += pendingCount
	to.pendingTasks = nil
	to.taskLock.Unlock()
	to.processingTasks.Add(-pendingCount)

	// Wait for processing tasks to complete
	to.processingTasks.Wait()
}
