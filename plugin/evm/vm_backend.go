package evm

import (
	"github.com/ava-labs/coreth/consensus/dummy"
	"github.com/ava-labs/coreth/core"
	"github.com/ava-labs/coreth/core/rawdb"
	"github.com/ava-labs/coreth/eth"
	"github.com/ava-labs/coreth/node"
	"github.com/ava-labs/coreth/v2/chain"
	"github.com/ethereum/go-ethereum/common"
	corevm "github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
)

func (vm *VM) createBackend(lastAcceptedHash common.Hash) (Backend, error) {
	if vm.config.UseV1Backend {
		return vm.createBackendV1(lastAcceptedHash)
	}
	return vm.createBackendV2(lastAcceptedHash)
}

func (vm *VM) createBackendV1(lastAcceptedHash common.Hash) (Backend, error) {
	nodecfg := &node.Config{
		SubnetEVMVersion:      Version,
		KeyStoreDir:           vm.config.KeystoreDirectory,
		ExternalSigner:        vm.config.KeystoreExternalSigner,
		InsecureUnlockAllowed: vm.config.KeystoreInsecureUnlockAllowed,
	}
	node, err := node.New(nodecfg)
	if err != nil {
		return nil, err
	}
	eth, err := eth.New(
		node,
		&vm.ethConfig,
		&EthPushGossiper{vm: vm},
		vm.chaindb,
		vm.config.EthBackendSettings(),
		lastAcceptedHash,
		&vm.clock,
	)
	if err != nil {
		return nil, err
	}
	return &ethBackender{eth}, nil
}

func (vm *VM) createBackendV2(lastAcceptedHash common.Hash) (Backend, error) {
	log.Info("Creating v2 backend")
	config := &vm.ethConfig

	// round TrieCleanCache and SnapshotCache up to nearest 64MB, since fastcache will mmap
	// memory in 64MBs chunks.
	config.TrieCleanCache = roundUpCacheSize(config.TrieCleanCache, 64)
	vmConfig := corevm.Config{
		EnablePreimageRecording: config.EnablePreimageRecording,
	}
	cacheConfig := &core.CacheConfig{
		TrieCleanLimit:            config.TrieCleanCache,
		TrieDirtyLimit:            config.TrieDirtyCache,
		TrieDirtyCommitTarget:     config.TrieDirtyCommitTarget,
		TriePrefetcherParallelism: config.TriePrefetcherParallelism,
		Pruning:                   config.Pruning,
		CommitInterval:            config.CommitInterval,
		Preimages:                 config.Preimages,
		AcceptedCacheSize:         config.AcceptedCacheSize,
		StateHistory:              config.StateHistory,
		StateScheme:               rawdb.PathScheme, // XXX: hardcoded to pathdb
	}
	engine := dummy.NewFakerWithClock(&vm.clock)
	// TODO: add support for separate trie db disk storage.
	// TODO: add support for chainDir
	tdb := chain.NewTrieDB(vm.chaindb, cacheConfig)
	committable := core.AsCommittable(vm.chaindb, tdb)
	blockchain, err := chain.NewBlockChain(
		vm.chaindb,
		committable,
		cacheConfig,
		config.Genesis,
		engine,
		vmConfig,
		lastAcceptedHash,
		vm.config.SkipUpgradeCheck,
	)
	if err != nil {
		return nil, err
	}
	backend, err := chain.NewLegacyBackend(
		blockchain,
		config.TxPool,
		&config.Miner,
		&vm.clock,
		config.GPO,
		config.AllowUnfinalizedQueries,
	)
	if err != nil {
		return nil, err
	}
	return &v2Backender{backend}, nil
}

// roundUpCacheSize returns [input] rounded up to the next multiple of [allocSize]
func roundUpCacheSize(input int, allocSize int) int {
	cacheChunks := (input + allocSize - 1) / allocSize
	return cacheChunks * allocSize
}
