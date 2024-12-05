// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ava-labs/avalanchego/cache"
	"github.com/ava-labs/avalanchego/cache/metercacher"
	"github.com/ava-labs/avalanchego/network/p2p"
	"github.com/ava-labs/avalanchego/network/p2p/acp118"
	"github.com/ava-labs/avalanchego/network/p2p/gossip"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/consensus/dummy"
	"github.com/ava-labs/subnet-evm/constants"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/txpool"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/eth"
	"github.com/ava-labs/subnet-evm/eth/ethconfig"
	"github.com/ava-labs/subnet-evm/metrics"
	subnetEVMPrometheus "github.com/ava-labs/subnet-evm/metrics/prometheus"
	"github.com/ava-labs/subnet-evm/miner"
	"github.com/ava-labs/subnet-evm/node"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/peer"
	"github.com/ava-labs/subnet-evm/plugin/evm/message"
	"github.com/ava-labs/subnet-evm/plugin/evm/validators"
	"github.com/ava-labs/subnet-evm/plugin/evm/validators/interfaces"
	"github.com/ava-labs/subnet-evm/triedb"
	"github.com/ava-labs/subnet-evm/triedb/hashdb"

	warpcontract "github.com/ava-labs/subnet-evm/precompile/contracts/warp"
	"github.com/ava-labs/subnet-evm/rpc"
	statesyncclient "github.com/ava-labs/subnet-evm/sync/client"
	"github.com/ava-labs/subnet-evm/sync/client/stats"
	"github.com/ava-labs/subnet-evm/warp"

	// Force-load tracer engine to trigger registration
	//
	// We must import this package (not referenced elsewhere) so that the native "callTracer"
	// is added to a map of client-accessible tracers. In geth, this is done
	// inside of cmd/geth.
	_ "github.com/ava-labs/subnet-evm/eth/tracers/js"
	_ "github.com/ava-labs/subnet-evm/eth/tracers/native"

	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	// Force-load precompiles to trigger registration
	_ "github.com/ava-labs/subnet-evm/precompile/registry"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"

	avalancheRPC "github.com/gorilla/rpc/v2"

	"github.com/ava-labs/avalanchego/codec"
	"github.com/ava-labs/avalanchego/database/versiondb"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/consensus/snowman"
	"github.com/ava-labs/avalanchego/snow/engine/snowman/block"
	"github.com/ava-labs/avalanchego/utils/perms"
	"github.com/ava-labs/avalanchego/utils/profiler"
	"github.com/ava-labs/avalanchego/utils/timer/mockable"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/version"
	"github.com/ava-labs/avalanchego/vms/components/chain"

	commonEng "github.com/ava-labs/avalanchego/snow/engine/common"

	"github.com/ava-labs/avalanchego/database"
	avalancheUtils "github.com/ava-labs/avalanchego/utils"
	avalancheJSON "github.com/ava-labs/avalanchego/utils/json"
)

var (
	_ block.ChainVM                      = &VM{}
	_ block.BuildBlockWithContextChainVM = &VM{}
	_ block.StateSyncableVM              = &VM{}
	_ statesyncclient.EthBlockParser     = &VM{}
)

const (
	// Max time from current time allowed for blocks, before they're considered future blocks
	// and fail verification
	maxFutureBlockTime     = 10 * time.Second
	decidedCacheSize       = 10 * units.MiB
	missingCacheSize       = 50
	unverifiedCacheSize    = 5 * units.MiB
	bytesToIDCacheSize     = 5 * units.MiB
	warpSignatureCacheSize = 500

	// Prefixes for metrics gatherers
	ethMetricsPrefix        = "eth"
	sdkMetricsPrefix        = "sdk"
	chainStateMetricsPrefix = "chain_state"

	// gossip constants
	pushGossipDiscardedElements          = 16_384
	txGossipBloomMinTargetElements       = 8 * 1024
	txGossipBloomTargetFalsePositiveRate = 0.01
	txGossipBloomResetFalsePositiveRate  = 0.05
	txGossipBloomChurnMultiplier         = 3
	txGossipTargetMessageSize            = 20 * units.KiB
	maxValidatorSetStaleness             = time.Minute
	txGossipThrottlingPeriod             = 10 * time.Second
	txGossipThrottlingLimit              = 2
	txGossipPollSize                     = 1
)

// Define the API endpoints for the VM
const (
	adminEndpoint        = "/admin"
	ethRPCEndpoint       = "/rpc"
	ethWSEndpoint        = "/ws"
	validatorsEndpoint   = "/validators"
	ethTxGossipNamespace = "eth_tx_gossip"
)

var (
	// Set last accepted key to be longer than the keys used to store accepted block IDs.
	lastAcceptedKey    = []byte("last_accepted_key")
	acceptedPrefix     = []byte("snowman_accepted")
	metadataPrefix     = []byte("metadata")
	warpPrefix         = []byte("warp")
	ethDBPrefix        = []byte("ethdb")
	validatorsDBPrefix = []byte("validators")
)

var (
	errEmptyBlock                    = errors.New("empty block")
	errUnsupportedFXs                = errors.New("unsupported feature extensions")
	errInvalidBlock                  = errors.New("invalid block")
	errInvalidNonce                  = errors.New("invalid nonce")
	errUnclesUnsupported             = errors.New("uncles unsupported")
	errNilBaseFeeSubnetEVM           = errors.New("nil base fee is invalid after subnetEVM")
	errNilBlockGasCostSubnetEVM      = errors.New("nil blockGasCost is invalid after subnetEVM")
	errInvalidHeaderPredicateResults = errors.New("invalid header predicate results")
)

// legacyApiNames maps pre geth v1.10.20 api names to their updated counterparts.
// used in attachEthService for backward configuration compatibility.
var legacyApiNames = map[string]string{
	"internal-public-eth":              "internal-eth",
	"internal-public-blockchain":       "internal-blockchain",
	"internal-public-transaction-pool": "internal-transaction",
	"internal-public-tx-pool":          "internal-tx-pool",
	"internal-public-debug":            "internal-debug",
	"internal-private-debug":           "internal-debug",
	"internal-public-account":          "internal-account",
	"internal-private-personal":        "internal-personal",

	"public-eth":        "eth",
	"public-eth-filter": "eth-filter",
	"private-admin":     "admin",
	"public-debug":      "debug",
	"private-debug":     "debug",
}

// VM implements the snowman.ChainVM interface
type VM struct {
	ctx *snow.Context
	// [cancel] may be nil until [snow.NormalOp] starts
	cancel context.CancelFunc
	// *chain.State helps to implement the VM interface by wrapping blocks
	// with an efficient caching layer.
	*chain.State

	config Config

	networkID   uint64
	genesisHash common.Hash
	chainConfig *params.ChainConfig
	ethConfig   ethconfig.Config

	// pointers to eth constructs
	eth        *eth.Ethereum
	txPool     *txpool.TxPool
	blockChain *core.BlockChain
	miner      *miner.Miner

	// [versiondb] is the VM's current versioned database
	versiondb *versiondb.Database

	// [db] is the VM's current database
	db database.Database

	// metadataDB is used to store one off keys.
	metadataDB database.Database

	// [chaindb] is the database supplied to the Ethereum backend
	chaindb ethdb.Database

	// [acceptedBlockDB] is the database to store the last accepted
	// block.
	acceptedBlockDB database.Database
	// [warpDB] is used to store warp message signatures
	// set to a prefixDB with the prefix [warpPrefix]
	warpDB database.Database

	validatorsDB database.Database

	toEngine chan<- commonEng.Message

	syntacticBlockValidator BlockValidator

	builder *blockBuilder

	clock mockable.Clock

	shutdownChan chan struct{}
	shutdownWg   sync.WaitGroup

	// Continuous Profiler
	profiler profiler.ContinuousProfiler

	peer.Network
	client       peer.NetworkClient
	networkCodec codec.Manager

	p2pValidators *p2p.Validators

	// Metrics
	sdkMetrics *prometheus.Registry

	bootstrapped avalancheUtils.Atomic[bool]

	logger SubnetEVMLogger
	// State sync server and client
	StateSyncServer
	StateSyncClient

	// Avalanche Warp Messaging backend
	// Used to serve BLS signatures of warp messages over RPC
	warpBackend warp.Backend

	// Initialize only sets these if nil so they can be overridden in tests
	p2pSender          commonEng.AppSender
	ethTxGossipHandler p2p.Handler
	ethTxPushGossiper  avalancheUtils.Atomic[*gossip.PushGossiper[*GossipEthTx]]
	ethTxPullGossiper  gossip.Gossiper

	validatorsManager interfaces.Manager

	chainAlias string
	// RPC handlers (should be stopped before closing chaindb)
	rpcHandlers []interface{ Stop() }
}

// Initialize implements the snowman.ChainVM interface
func (vm *VM) Initialize(
	_ context.Context,
	chainCtx *snow.Context,
	db database.Database,
	genesisBytes []byte,
	upgradeBytes []byte,
	configBytes []byte,
	toEngine chan<- commonEng.Message,
	fxs []*commonEng.Fx,
	appSender commonEng.AppSender,
) error {
	vm.config.SetDefaults()
	if len(configBytes) > 0 {
		if err := json.Unmarshal(configBytes, &vm.config); err != nil {
			return fmt.Errorf("failed to unmarshal config %s: %w", string(configBytes), err)
		}
	}
	if err := vm.config.Validate(); err != nil {
		return err
	}
	// We should deprecate config flags as the first thing, before we do anything else
	// because this can set old flags to new flags. log the message after we have
	// initialized the logger.
	deprecateMsg := vm.config.Deprecate()

	vm.ctx = chainCtx

	// Create logger
	alias, err := vm.ctx.BCLookup.PrimaryAlias(vm.ctx.ChainID)
	if err != nil {
		// fallback to ChainID string instead of erroring
		alias = vm.ctx.ChainID.String()
	}
	vm.chainAlias = alias

	subnetEVMLogger, err := InitLogger(vm.chainAlias, vm.config.LogLevel, vm.config.LogJSONFormat, vm.ctx.Log)
	if err != nil {
		return fmt.Errorf("failed to initialize logger due to: %w ", err)
	}
	vm.logger = subnetEVMLogger

	log.Info("Initializing Subnet EVM VM", "Version", Version, "Config", vm.config)

	if deprecateMsg != "" {
		log.Warn("Deprecation Warning", "msg", deprecateMsg)
	}

	if len(fxs) > 0 {
		return errUnsupportedFXs
	}

	// Enable debug-level metrics that might impact runtime performance
	metrics.EnabledExpensive = vm.config.MetricsExpensiveEnabled

	vm.toEngine = toEngine
	vm.shutdownChan = make(chan struct{}, 1)

	if err := vm.initializeMetrics(); err != nil {
		return fmt.Errorf("failed to initialize metrics: %w", err)
	}

	// Initialize the database
	if err := vm.initializeDBs(db); err != nil {
		return fmt.Errorf("failed to initialize databases: %w", err)
	}

	if vm.config.InspectDatabase {
		if err := vm.inspectDatabases(); err != nil {
			return err
		}
	}

	g := new(core.Genesis)
	if err := json.Unmarshal(genesisBytes, g); err != nil {
		return err
	}

	if g.Config == nil {
		g.Config = params.SubnetEVMDefaultChainConfig
	}

	// Set the Avalanche Context on the ChainConfig
	g.Config.AvalancheContext = params.AvalancheContext{
		SnowCtx: chainCtx,
	}

	g.Config.SetNetworkUpgradeDefaults()

	// Load airdrop file if provided
	if vm.config.AirdropFile != "" {
		g.AirdropData, err = os.ReadFile(vm.config.AirdropFile)
		if err != nil {
			return fmt.Errorf("could not read airdrop file '%s': %w", vm.config.AirdropFile, err)
		}
	}

	vm.syntacticBlockValidator = NewBlockValidator()

	if g.Config.FeeConfig == commontype.EmptyFeeConfig {
		log.Info("No fee config given in genesis, setting default fee config", "DefaultFeeConfig", params.DefaultFeeConfig)
		g.Config.FeeConfig = params.DefaultFeeConfig
	}

	// Apply upgradeBytes (if any) by unmarshalling them into [chainConfig.UpgradeConfig].
	// Initializing the chain will verify upgradeBytes are compatible with existing values.
	// This should be called before g.Verify().
	if len(upgradeBytes) > 0 {
		var upgradeConfig params.UpgradeConfig
		if err := json.Unmarshal(upgradeBytes, &upgradeConfig); err != nil {
			return fmt.Errorf("failed to parse upgrade bytes: %w", err)
		}
		g.Config.UpgradeConfig = upgradeConfig
	}

	if g.Config.UpgradeConfig.NetworkUpgradeOverrides != nil {
		overrides := g.Config.UpgradeConfig.NetworkUpgradeOverrides
		marshaled, err := json.Marshal(overrides)
		if err != nil {
			log.Warn("Failed to marshal network upgrade overrides", "error", err, "overrides", overrides)
		} else {
			log.Info("Applying network upgrade overrides", "overrides", string(marshaled))
		}
		g.Config.Override(overrides)
	}

	g.Config.SetEthUpgrades(g.Config.NetworkUpgrades)

	if err := g.Verify(); err != nil {
		return fmt.Errorf("failed to verify genesis: %w", err)
	}

	vm.ethConfig = ethconfig.NewDefaultConfig()
	vm.ethConfig.Genesis = g
	// NetworkID here is different than Avalanche's NetworkID.
	// Avalanche's NetworkID represents the Avalanche network is running on
	// like Fuji, Mainnet, Local, etc.
	// The NetworkId here is kept same as ChainID to be compatible with
	// Ethereum tooling.
	vm.ethConfig.NetworkId = g.Config.ChainID.Uint64()

	// Set minimum price for mining and default gas price oracle value to the min
	// gas price to prevent so transactions and blocks all use the correct fees
	vm.ethConfig.RPCGasCap = vm.config.RPCGasCap
	vm.ethConfig.RPCEVMTimeout = vm.config.APIMaxDuration.Duration
	vm.ethConfig.RPCTxFeeCap = vm.config.RPCTxFeeCap

	vm.ethConfig.TxPool.Locals = vm.config.PriorityRegossipAddresses
	vm.ethConfig.TxPool.NoLocals = !vm.config.LocalTxsEnabled
	vm.ethConfig.TxPool.PriceLimit = vm.config.TxPoolPriceLimit
	vm.ethConfig.TxPool.PriceBump = vm.config.TxPoolPriceBump
	vm.ethConfig.TxPool.AccountSlots = vm.config.TxPoolAccountSlots
	vm.ethConfig.TxPool.GlobalSlots = vm.config.TxPoolGlobalSlots
	vm.ethConfig.TxPool.AccountQueue = vm.config.TxPoolAccountQueue
	vm.ethConfig.TxPool.GlobalQueue = vm.config.TxPoolGlobalQueue
	vm.ethConfig.TxPool.Lifetime = vm.config.TxPoolLifetime.Duration

	vm.ethConfig.AllowUnfinalizedQueries = vm.config.AllowUnfinalizedQueries
	vm.ethConfig.AllowUnprotectedTxs = vm.config.AllowUnprotectedTxs
	vm.ethConfig.AllowUnprotectedTxHashes = vm.config.AllowUnprotectedTxHashes
	vm.ethConfig.Preimages = vm.config.Preimages
	vm.ethConfig.Pruning = vm.config.Pruning
	vm.ethConfig.TrieCleanCache = vm.config.TrieCleanCache
	vm.ethConfig.TrieDirtyCache = vm.config.TrieDirtyCache
	vm.ethConfig.TrieDirtyCommitTarget = vm.config.TrieDirtyCommitTarget
	vm.ethConfig.TriePrefetcherParallelism = vm.config.TriePrefetcherParallelism
	vm.ethConfig.SnapshotCache = vm.config.SnapshotCache
	vm.ethConfig.AcceptorQueueLimit = vm.config.AcceptorQueueLimit
	vm.ethConfig.PopulateMissingTries = vm.config.PopulateMissingTries
	vm.ethConfig.PopulateMissingTriesParallelism = vm.config.PopulateMissingTriesParallelism
	vm.ethConfig.AllowMissingTries = vm.config.AllowMissingTries
	vm.ethConfig.SnapshotDelayInit = vm.config.StateSyncEnabled
	vm.ethConfig.SnapshotWait = vm.config.SnapshotWait
	vm.ethConfig.SnapshotVerify = vm.config.SnapshotVerify
	vm.ethConfig.OfflinePruning = vm.config.OfflinePruning
	vm.ethConfig.OfflinePruningBloomFilterSize = vm.config.OfflinePruningBloomFilterSize
	vm.ethConfig.OfflinePruningDataDirectory = vm.config.OfflinePruningDataDirectory
	vm.ethConfig.CommitInterval = vm.config.CommitInterval
	vm.ethConfig.SkipUpgradeCheck = vm.config.SkipUpgradeCheck
	vm.ethConfig.AcceptedCacheSize = vm.config.AcceptedCacheSize
	vm.ethConfig.TransactionHistory = vm.config.TransactionHistory
	vm.ethConfig.SkipTxIndexing = vm.config.SkipTxIndexing

	// Create directory for offline pruning
	if len(vm.ethConfig.OfflinePruningDataDirectory) != 0 {
		if err := os.MkdirAll(vm.ethConfig.OfflinePruningDataDirectory, perms.ReadWriteExecute); err != nil {
			log.Error("failed to create offline pruning data directory", "error", err)
			return err
		}
	}

	// Handle custom fee recipient
	if common.IsHexAddress(vm.config.FeeRecipient) {
		address := common.HexToAddress(vm.config.FeeRecipient)
		log.Info("Setting fee recipient", "address", address)
		vm.ethConfig.Miner.Etherbase = address
	} else {
		log.Info("Config has not specified any coinbase address. Defaulting to the blackhole address.")
		vm.ethConfig.Miner.Etherbase = constants.BlackholeAddr
	}

	vm.chainConfig = g.Config
	vm.networkID = vm.ethConfig.NetworkId

	// create genesisHash after applying upgradeBytes in case
	// upgradeBytes modifies genesis.
	vm.genesisHash = vm.ethConfig.Genesis.ToBlock().Hash() // must create genesis hash before [vm.readLastAccepted]
	lastAcceptedHash, lastAcceptedHeight, err := vm.readLastAccepted()
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("lastAccepted = %s", lastAcceptedHash))

	// initialize peer network
	if vm.p2pSender == nil {
		vm.p2pSender = appSender
	}

	p2pNetwork, err := p2p.NewNetwork(vm.ctx.Log, vm.p2pSender, vm.sdkMetrics, "p2p")
	if err != nil {
		return fmt.Errorf("failed to initialize p2p network: %w", err)
	}
	vm.p2pValidators = p2p.NewValidators(p2pNetwork.Peers, vm.ctx.Log, vm.ctx.SubnetID, vm.ctx.ValidatorState, maxValidatorSetStaleness)
	vm.networkCodec = message.Codec
	vm.Network = peer.NewNetwork(p2pNetwork, appSender, vm.networkCodec, chainCtx.NodeID, vm.config.MaxOutboundActiveRequests)
	vm.client = peer.NewNetworkClient(vm.Network)

	vm.validatorsManager, err = validators.NewManager(vm.ctx, vm.validatorsDB, &vm.clock)
	if err != nil {
		return fmt.Errorf("failed to initialize validators manager: %w", err)
	}

	// Initialize warp backend
	offchainWarpMessages := make([][]byte, len(vm.config.WarpOffChainMessages))
	for i, hexMsg := range vm.config.WarpOffChainMessages {
		offchainWarpMessages[i] = []byte(hexMsg)
	}
	warpSignatureCache := &cache.LRU[ids.ID, []byte]{Size: warpSignatureCacheSize}
	meteredCache, err := metercacher.New("warp_signature_cache", vm.sdkMetrics, warpSignatureCache)
	if err != nil {
		return fmt.Errorf("failed to create warp signature cache: %w", err)
	}

	// clear warpdb on initialization if config enabled
	if vm.config.PruneWarpDB {
		if err := database.Clear(vm.warpDB, ethdb.IdealBatchSize); err != nil {
			return fmt.Errorf("failed to prune warpDB: %w", err)
		}
	}

	vm.warpBackend, err = warp.NewBackend(
		vm.ctx.NetworkID,
		vm.ctx.ChainID,
		vm.ctx.WarpSigner,
		vm,
		vm.validatorsManager,
		vm.warpDB,
		meteredCache,
		offchainWarpMessages,
	)
	if err != nil {
		return err
	}

	if err := vm.initializeChain(lastAcceptedHash, vm.ethConfig); err != nil {
		return err
	}

	go vm.ctx.Log.RecoverAndPanic(vm.startContinuousProfiler)

	// Add p2p warp message warpHandler
	warpHandler := acp118.NewCachedHandler(meteredCache, vm.warpBackend, vm.ctx.WarpSigner)
	vm.Network.AddHandler(p2p.SignatureRequestHandlerID, warpHandler)

	vm.setAppRequestHandlers()

	vm.StateSyncServer = NewStateSyncServer(&stateSyncServerConfig{
		Chain:            vm.blockChain,
		SyncableInterval: vm.config.StateSyncCommitInterval,
	})
	return vm.initializeStateSyncClient(lastAcceptedHeight)
}

func (vm *VM) initializeMetrics() error {
	vm.sdkMetrics = prometheus.NewRegistry()
	// If metrics are enabled, register the default metrics registry
	if !metrics.Enabled {
		return nil
	}

	gatherer := subnetEVMPrometheus.Gatherer(metrics.DefaultRegistry)
	if err := vm.ctx.Metrics.Register(ethMetricsPrefix, gatherer); err != nil {
		return err
	}
	return vm.ctx.Metrics.Register(sdkMetricsPrefix, vm.sdkMetrics)
}

func (vm *VM) initializeChain(lastAcceptedHash common.Hash, ethConfig ethconfig.Config) error {
	nodecfg := &node.Config{
		SubnetEVMVersion:      Version,
		KeyStoreDir:           vm.config.KeystoreDirectory,
		ExternalSigner:        vm.config.KeystoreExternalSigner,
		InsecureUnlockAllowed: vm.config.KeystoreInsecureUnlockAllowed,
	}
	node, err := node.New(nodecfg)
	if err != nil {
		return err
	}
	vm.eth, err = eth.New(
		node,
		&vm.ethConfig,
		&EthPushGossiper{vm: vm},
		vm.chaindb,
		vm.config.EthBackendSettings(),
		lastAcceptedHash,
		dummy.NewFakerWithClock(&vm.clock),
		&vm.clock,
	)
	if err != nil {
		return err
	}
	vm.eth.SetEtherbase(ethConfig.Miner.Etherbase)
	vm.txPool = vm.eth.TxPool()
	vm.txPool.SetMinFee(vm.chainConfig.FeeConfig.MinBaseFee)
	vm.txPool.SetGasTip(big.NewInt(0))
	vm.blockChain = vm.eth.BlockChain()
	vm.miner = vm.eth.Miner()

	vm.eth.Start()
	return vm.initChainState(vm.blockChain.LastAcceptedBlock())
}

// initializeStateSyncClient initializes the client for performing state sync.
// If state sync is disabled, this function will wipe any ongoing summary from
// disk to ensure that we do not continue syncing from an invalid snapshot.
func (vm *VM) initializeStateSyncClient(lastAcceptedHeight uint64) error {
	// parse nodeIDs from state sync IDs in vm config
	var stateSyncIDs []ids.NodeID
	if vm.config.StateSyncEnabled && len(vm.config.StateSyncIDs) > 0 {
		nodeIDs := strings.Split(vm.config.StateSyncIDs, ",")
		stateSyncIDs = make([]ids.NodeID, len(nodeIDs))
		for i, nodeIDString := range nodeIDs {
			nodeID, err := ids.NodeIDFromString(nodeIDString)
			if err != nil {
				return fmt.Errorf("failed to parse %s as NodeID: %w", nodeIDString, err)
			}
			stateSyncIDs[i] = nodeID
		}
	}

	vm.StateSyncClient = NewStateSyncClient(&stateSyncClientConfig{
		chain: vm.eth,
		state: vm.State,
		client: statesyncclient.NewClient(
			&statesyncclient.ClientConfig{
				NetworkClient:    vm.client,
				Codec:            vm.networkCodec,
				Stats:            stats.NewClientSyncerStats(),
				StateSyncNodeIDs: stateSyncIDs,
				BlockParser:      vm,
			},
		),
		enabled:              vm.config.StateSyncEnabled,
		skipResume:           vm.config.StateSyncSkipResume,
		stateSyncMinBlocks:   vm.config.StateSyncMinBlocks,
		stateSyncRequestSize: vm.config.StateSyncRequestSize,
		lastAcceptedHeight:   lastAcceptedHeight, // TODO clean up how this is passed around
		chaindb:              vm.chaindb,
		metadataDB:           vm.metadataDB,
		acceptedBlockDB:      vm.acceptedBlockDB,
		db:                   vm.versiondb,
		toEngine:             vm.toEngine,
	})

	// If StateSync is disabled, clear any ongoing summary so that we will not attempt to resume
	// sync using a snapshot that has been modified by the node running normal operations.
	if !vm.config.StateSyncEnabled {
		return vm.StateSyncClient.ClearOngoingSummary()
	}

	return nil
}

func (vm *VM) initChainState(lastAcceptedBlock *types.Block) error {
	block := vm.newBlock(lastAcceptedBlock)

	config := &chain.Config{
		DecidedCacheSize:      decidedCacheSize,
		MissingCacheSize:      missingCacheSize,
		UnverifiedCacheSize:   unverifiedCacheSize,
		BytesToIDCacheSize:    bytesToIDCacheSize,
		GetBlock:              vm.getBlock,
		UnmarshalBlock:        vm.parseBlock,
		BuildBlock:            vm.buildBlock,
		BuildBlockWithContext: vm.buildBlockWithContext,
		LastAcceptedBlock:     block,
	}

	// Register chain state metrics
	chainStateRegisterer := prometheus.NewRegistry()
	state, err := chain.NewMeteredState(chainStateRegisterer, config)
	if err != nil {
		return fmt.Errorf("could not create metered state: %w", err)
	}
	vm.State = state

	if !metrics.Enabled {
		return nil
	}

	return vm.ctx.Metrics.Register(chainStateMetricsPrefix, chainStateRegisterer)
}

func (vm *VM) SetState(_ context.Context, state snow.State) error {
	switch state {
	case snow.StateSyncing:
		vm.bootstrapped.Set(false)
		return nil
	case snow.Bootstrapping:
		return vm.onBootstrapStarted()
	case snow.NormalOp:
		return vm.onNormalOperationsStarted()
	default:
		return snow.ErrUnknownState
	}
}

// onBootstrapStarted marks this VM as bootstrapping
func (vm *VM) onBootstrapStarted() error {
	vm.bootstrapped.Set(false)
	if err := vm.StateSyncClient.Error(); err != nil {
		return err
	}
	// After starting bootstrapping, do not attempt to resume a previous state sync.
	if err := vm.StateSyncClient.ClearOngoingSummary(); err != nil {
		return err
	}
	// Ensure snapshots are initialized before bootstrapping (i.e., if state sync is skipped).
	// Note calling this function has no effect if snapshots are already initialized.
	vm.blockChain.InitializeSnapshots()

	return nil
}

// onNormalOperationsStarted marks this VM as bootstrapped
func (vm *VM) onNormalOperationsStarted() error {
	if vm.bootstrapped.Get() {
		return nil
	}
	vm.bootstrapped.Set(true)

	ctx, cancel := context.WithCancel(context.TODO())
	vm.cancel = cancel

	// sync validators first
	if err := vm.validatorsManager.Sync(ctx); err != nil {
		return fmt.Errorf("failed to update validators: %w", err)
	}
	vdrIDs := vm.validatorsManager.GetNodeIDs().List()
	// Then start tracking with updated validators
	// StartTracking initializes the uptime tracking with the known validators
	// and update their uptime to account for the time we were being offline.
	if err := vm.validatorsManager.StartTracking(vdrIDs); err != nil {
		return fmt.Errorf("failed to start tracking uptime: %w", err)
	}
	// dispatch validator set update
	vm.shutdownWg.Add(1)
	go func() {
		vm.validatorsManager.DispatchSync(ctx)
		vm.shutdownWg.Done()
	}()

	// Initialize goroutines related to block building
	// once we enter normal operation as there is no need to handle mempool gossip before this point.
	ethTxGossipMarshaller := GossipEthTxMarshaller{}
	ethTxGossipClient := vm.Network.NewClient(p2p.TxGossipHandlerID, p2p.WithValidatorSampling(vm.p2pValidators))
	ethTxGossipMetrics, err := gossip.NewMetrics(vm.sdkMetrics, ethTxGossipNamespace)
	if err != nil {
		return fmt.Errorf("failed to initialize eth tx gossip metrics: %w", err)
	}
	ethTxPool, err := NewGossipEthTxPool(vm.txPool, vm.sdkMetrics)
	if err != nil {
		return fmt.Errorf("failed to initialize gossip eth tx pool: %w", err)
	}
	vm.shutdownWg.Add(1)
	go func() {
		ethTxPool.Subscribe(ctx)
		vm.shutdownWg.Done()
	}()

	pushGossipParams := gossip.BranchingFactor{
		StakePercentage: vm.config.PushGossipPercentStake,
		Validators:      vm.config.PushGossipNumValidators,
		Peers:           vm.config.PushGossipNumPeers,
	}
	pushRegossipParams := gossip.BranchingFactor{
		Validators: vm.config.PushRegossipNumValidators,
		Peers:      vm.config.PushRegossipNumPeers,
	}

	ethTxPushGossiper := vm.ethTxPushGossiper.Get()
	if ethTxPushGossiper == nil {
		ethTxPushGossiper, err = gossip.NewPushGossiper[*GossipEthTx](
			ethTxGossipMarshaller,
			ethTxPool,
			vm.p2pValidators,
			ethTxGossipClient,
			ethTxGossipMetrics,
			pushGossipParams,
			pushRegossipParams,
			pushGossipDiscardedElements,
			txGossipTargetMessageSize,
			vm.config.RegossipFrequency.Duration,
		)
		if err != nil {
			return fmt.Errorf("failed to initialize eth tx push gossiper: %w", err)
		}
		vm.ethTxPushGossiper.Set(ethTxPushGossiper)
	}

	// NOTE: gossip network must be initialized first otherwise ETH tx gossip will not work.
	gossipStats := NewGossipStats()
	vm.builder = vm.NewBlockBuilder(vm.toEngine)
	vm.builder.awaitSubmittedTxs()
	vm.Network.SetGossipHandler(NewGossipHandler(vm, gossipStats))

	if vm.ethTxGossipHandler == nil {
		vm.ethTxGossipHandler = newTxGossipHandler[*GossipEthTx](
			vm.ctx.Log,
			ethTxGossipMarshaller,
			ethTxPool,
			ethTxGossipMetrics,
			txGossipTargetMessageSize,
			txGossipThrottlingPeriod,
			txGossipThrottlingLimit,
			vm.p2pValidators,
		)
	}

	if err := vm.Network.AddHandler(p2p.TxGossipHandlerID, vm.ethTxGossipHandler); err != nil {
		return fmt.Errorf("failed to add eth tx gossip handler: %w", err)
	}

	if vm.ethTxPullGossiper == nil {
		ethTxPullGossiper := gossip.NewPullGossiper[*GossipEthTx](
			vm.ctx.Log,
			ethTxGossipMarshaller,
			ethTxPool,
			ethTxGossipClient,
			ethTxGossipMetrics,
			txGossipPollSize,
		)

		vm.ethTxPullGossiper = gossip.ValidatorGossiper{
			Gossiper:   ethTxPullGossiper,
			NodeID:     vm.ctx.NodeID,
			Validators: vm.p2pValidators,
		}
	}

	vm.shutdownWg.Add(2)
	go func() {
		gossip.Every(ctx, vm.ctx.Log, ethTxPushGossiper, vm.config.PushGossipFrequency.Duration)
		vm.shutdownWg.Done()
	}()
	go func() {
		gossip.Every(ctx, vm.ctx.Log, vm.ethTxPullGossiper, vm.config.PullGossipFrequency.Duration)
		vm.shutdownWg.Done()
	}()

	return nil
}

// setAppRequestHandlers sets the request handlers for the VM to serve state sync
// requests.
func (vm *VM) setAppRequestHandlers() {
	// Create standalone EVM TrieDB (read only) for serving leafs requests.
	// We create a standalone TrieDB here, so that it has a standalone cache from the one
	// used by the node when processing blocks.
	evmTrieDB := triedb.NewDatabase(
		vm.chaindb,
		&triedb.Config{
			HashDB: &hashdb.Config{
				CleanCacheSize: vm.config.StateSyncServerTrieCache * units.MiB,
			},
		},
	)

	networkHandler := newNetworkHandler(vm.blockChain, vm.chaindb, evmTrieDB, vm.warpBackend, vm.networkCodec)
	vm.Network.SetRequestHandler(networkHandler)
}

// Shutdown implements the snowman.ChainVM interface
func (vm *VM) Shutdown(context.Context) error {
	if vm.ctx == nil {
		return nil
	}
	if vm.cancel != nil {
		vm.cancel()
	}
	if vm.bootstrapped.Get() {
		vdrIDs := vm.validatorsManager.GetNodeIDs().List()
		if err := vm.validatorsManager.StopTracking(vdrIDs); err != nil {
			return fmt.Errorf("failed to stop tracking uptime: %w", err)
		}
		if err := vm.validatorsManager.WriteState(); err != nil {
			return fmt.Errorf("failed to write validator: %w", err)
		}
	}
	vm.Network.Shutdown()
	if err := vm.StateSyncClient.Shutdown(); err != nil {
		log.Error("error stopping state syncer", "err", err)
	}
	close(vm.shutdownChan)
	// Stop RPC handlers before eth.Stop which will close the database
	for _, handler := range vm.rpcHandlers {
		handler.Stop()
	}
	vm.eth.Stop()
	log.Info("Ethereum backend stop completed")
	if err := vm.db.Close(); err != nil {
		log.Error("failed to close database: %w", err)
	} else {
		log.Info("Database closed")
	}
	vm.shutdownWg.Wait()
	log.Info("Subnet-EVM Shutdown completed")
	return nil
}

// buildBlock builds a block to be wrapped by ChainState
func (vm *VM) buildBlock(ctx context.Context) (snowman.Block, error) {
	return vm.buildBlockWithContext(ctx, nil)
}

func (vm *VM) buildBlockWithContext(ctx context.Context, proposerVMBlockCtx *block.Context) (snowman.Block, error) {
	if proposerVMBlockCtx != nil {
		log.Debug("Building block with context", "pChainBlockHeight", proposerVMBlockCtx.PChainHeight)
	} else {
		log.Debug("Building block without context")
	}
	predicateCtx := &precompileconfig.PredicateContext{
		SnowCtx:            vm.ctx,
		ProposerVMBlockCtx: proposerVMBlockCtx,
	}

	block, err := vm.miner.GenerateBlock(predicateCtx)
	vm.builder.handleGenerateBlock()
	if err != nil {
		return nil, err
	}

	// Note: the status of block is set by ChainState
	blk := vm.newBlock(block)

	// Verify is called on a non-wrapped block here, such that this
	// does not add [blk] to the processing blocks map in ChainState.
	//
	// TODO cache verification since Verify() will be called by the
	// consensus engine as well.
	//
	// Note: this is only called when building a new block, so caching
	// verification will only be a significant optimization for nodes
	// that produce a large number of blocks.
	// We call verify without writes here to avoid generating a reference
	// to the blk state root in the triedb when we are going to call verify
	// again from the consensus engine with writes enabled.
	if err := blk.verify(predicateCtx, false /*=writes*/); err != nil {
		return nil, fmt.Errorf("block failed verification due to: %w", err)
	}

	log.Debug(fmt.Sprintf("Built block %s", blk.ID()))
	// Marks the current transactions from the mempool as being successfully issued
	// into a block.
	return blk, nil
}

// parseBlock parses [b] into a block to be wrapped by ChainState.
func (vm *VM) parseBlock(_ context.Context, b []byte) (snowman.Block, error) {
	ethBlock := new(types.Block)
	if err := rlp.DecodeBytes(b, ethBlock); err != nil {
		return nil, err
	}

	// Note: the status of block is set by ChainState
	block := vm.newBlock(ethBlock)
	// Performing syntactic verification in ParseBlock allows for
	// short-circuiting bad blocks before they are processed by the VM.
	if err := block.syntacticVerify(); err != nil {
		return nil, fmt.Errorf("syntactic block verification failed: %w", err)
	}
	return block, nil
}

func (vm *VM) ParseEthBlock(b []byte) (*types.Block, error) {
	block, err := vm.parseBlock(context.TODO(), b)
	if err != nil {
		return nil, err
	}

	return block.(*Block).ethBlock, nil
}

// getBlock attempts to retrieve block [id] from the VM to be wrapped
// by ChainState.
func (vm *VM) getBlock(_ context.Context, id ids.ID) (snowman.Block, error) {
	ethBlock := vm.blockChain.GetBlockByHash(common.Hash(id))
	// If [ethBlock] is nil, return [database.ErrNotFound] here
	// so that the miss is considered cacheable.
	if ethBlock == nil {
		return nil, database.ErrNotFound
	}
	// Note: the status of block is set by ChainState
	return vm.newBlock(ethBlock), nil
}

// GetAcceptedBlock attempts to retrieve block [blkID] from the VM. This method
// only returns accepted blocks.
func (vm *VM) GetAcceptedBlock(ctx context.Context, blkID ids.ID) (snowman.Block, error) {
	blk, err := vm.GetBlock(ctx, blkID)
	if err != nil {
		return nil, err
	}

	height := blk.Height()
	acceptedBlkID, err := vm.GetBlockIDAtHeight(ctx, height)
	if err != nil {
		return nil, err
	}

	if acceptedBlkID != blkID {
		// The provided block is not accepted.
		return nil, database.ErrNotFound
	}
	return blk, nil
}

// SetPreference sets what the current tail of the chain is
func (vm *VM) SetPreference(ctx context.Context, blkID ids.ID) error {
	// Since each internal handler used by [vm.State] always returns a block
	// with non-nil ethBlock value, GetBlockInternal should never return a
	// (*Block) with a nil ethBlock value.
	block, err := vm.GetBlockInternal(ctx, blkID)
	if err != nil {
		return fmt.Errorf("failed to set preference to %s: %w", blkID, err)
	}

	return vm.blockChain.SetPreference(block.(*Block).ethBlock)
}

// VerifyHeightIndex always returns a nil error since the index is maintained by
// vm.blockChain.
func (vm *VM) VerifyHeightIndex(context.Context) error {
	return nil
}

// GetBlockIDAtHeight returns the canonical block at [height].
// Note: the engine assumes that if a block is not found at [height], then
// [database.ErrNotFound] will be returned. This indicates that the VM has state
// synced and does not have all historical blocks available.
func (vm *VM) GetBlockIDAtHeight(_ context.Context, height uint64) (ids.ID, error) {
	lastAcceptedBlock := vm.LastAcceptedBlock()
	if lastAcceptedBlock.Height() < height {
		return ids.ID{}, database.ErrNotFound
	}

	hash := vm.blockChain.GetCanonicalHash(height)
	if hash == (common.Hash{}) {
		return ids.ID{}, database.ErrNotFound
	}
	return ids.ID(hash), nil
}

func (vm *VM) Version(context.Context) (string, error) {
	return Version, nil
}

// NewHandler returns a new Handler for a service where:
//   - The handler's functionality is defined by [service]
//     [service] should be a gorilla RPC service (see https://www.gorillatoolkit.org/pkg/rpc/v2)
//   - The name of the service is [name]
func newHandler(name string, service interface{}) (http.Handler, error) {
	server := avalancheRPC.NewServer()
	server.RegisterCodec(avalancheJSON.NewCodec(), "application/json")
	server.RegisterCodec(avalancheJSON.NewCodec(), "application/json;charset=UTF-8")
	return server, server.RegisterService(service, name)
}

// CreateHandlers makes new http handlers that can handle API calls
func (vm *VM) CreateHandlers(context.Context) (map[string]http.Handler, error) {
	handler := rpc.NewServer(vm.config.APIMaxDuration.Duration)
	if vm.config.HttpBodyLimit > 0 {
		handler.SetHTTPBodyLimit(int(vm.config.HttpBodyLimit))
	}

	enabledAPIs := vm.config.EthAPIs()
	if err := attachEthService(handler, vm.eth.APIs(), enabledAPIs); err != nil {
		return nil, err
	}

	apis := make(map[string]http.Handler)
	if vm.config.AdminAPIEnabled {
		adminAPI, err := newHandler("admin", NewAdminService(vm, os.ExpandEnv(fmt.Sprintf("%s_subnet_evm_performance_%s", vm.config.AdminAPIDir, vm.chainAlias))))
		if err != nil {
			return nil, fmt.Errorf("failed to register service for admin API due to %w", err)
		}
		apis[adminEndpoint] = adminAPI
		enabledAPIs = append(enabledAPIs, "subnet-evm-admin")
	}

	if vm.config.ValidatorsAPIEnabled {
		validatorsAPI, err := newHandler("validators", &ValidatorsAPI{vm})
		if err != nil {
			return nil, fmt.Errorf("failed to register service for validators API due to %w", err)
		}
		apis[validatorsEndpoint] = validatorsAPI
		enabledAPIs = append(enabledAPIs, "validators")
	}

	// RPC APIs
	if vm.config.SnowmanAPIEnabled {
		if err := handler.RegisterName("snowman", &SnowmanAPI{vm}); err != nil {
			return nil, err
		}
		enabledAPIs = append(enabledAPIs, "snowman")
	}

	if vm.config.WarpAPIEnabled {
		if err := handler.RegisterName("warp", warp.NewAPI(vm.ctx.NetworkID, vm.ctx.SubnetID, vm.ctx.ChainID, vm.ctx.ValidatorState, vm.warpBackend, vm.client, vm.requirePrimaryNetworkSigners)); err != nil {
			return nil, err
		}
		enabledAPIs = append(enabledAPIs, "warp")
	}

	log.Info(fmt.Sprintf("Enabled APIs: %s", strings.Join(enabledAPIs, ", ")))
	apis[ethRPCEndpoint] = handler
	apis[ethWSEndpoint] = handler.WebsocketHandlerWithDuration(
		[]string{"*"},
		vm.config.APIMaxDuration.Duration,
		vm.config.WSCPURefillRate.Duration,
		vm.config.WSCPUMaxStored.Duration,
	)

	vm.rpcHandlers = append(vm.rpcHandlers, handler)
	return apis, nil
}

// CreateStaticHandlers makes new http handlers that can handle API calls
func (vm *VM) CreateStaticHandlers(context.Context) (map[string]http.Handler, error) {
	handler := rpc.NewServer(0)
	if vm.config.HttpBodyLimit > 0 {
		handler.SetHTTPBodyLimit(int(vm.config.HttpBodyLimit))
	}
	if err := handler.RegisterName("static", &StaticService{}); err != nil {
		return nil, err
	}

	vm.rpcHandlers = append(vm.rpcHandlers, handler)
	return map[string]http.Handler{
		"/rpc": handler,
	}, nil
}

/*
 ******************************************************************************
 *********************************** Helpers **********************************
 ******************************************************************************
 */

// GetCurrentNonce returns the nonce associated with the address at the
// preferred block
func (vm *VM) GetCurrentNonce(address common.Address) (uint64, error) {
	// Note: current state uses the state of the preferred block.
	state, err := vm.blockChain.State()
	if err != nil {
		return 0, err
	}
	return state.GetNonce(address), nil
}

// currentRules returns the chain rules for the current block.
func (vm *VM) currentRules() params.Rules {
	header := vm.eth.APIBackend.CurrentHeader()
	return vm.chainConfig.Rules(header.Number, header.Time)
}

// requirePrimaryNetworkSigners returns true if warp messages from the primary
// network must be signed by the primary network validators.
// This is necessary when the subnet is not validating the primary network.
func (vm *VM) requirePrimaryNetworkSigners() bool {
	switch c := vm.currentRules().ActivePrecompiles[warpcontract.ContractAddress].(type) {
	case *warpcontract.Config:
		return c.RequirePrimaryNetworkSigners
	default: // includes nil due to non-presence
		return false
	}
}

func (vm *VM) startContinuousProfiler() {
	// If the profiler directory is empty, return immediately
	// without creating or starting a continuous profiler.
	if vm.config.ContinuousProfilerDir == "" {
		return
	}
	vm.profiler = profiler.NewContinuous(
		filepath.Join(vm.config.ContinuousProfilerDir),
		vm.config.ContinuousProfilerFrequency.Duration,
		vm.config.ContinuousProfilerMaxFiles,
	)
	defer vm.profiler.Shutdown()

	vm.shutdownWg.Add(1)
	go func() {
		defer vm.shutdownWg.Done()
		log.Info("Dispatching continuous profiler", "dir", vm.config.ContinuousProfilerDir, "freq", vm.config.ContinuousProfilerFrequency, "maxFiles", vm.config.ContinuousProfilerMaxFiles)
		err := vm.profiler.Dispatch()
		if err != nil {
			log.Error("continuous profiler failed", "err", err)
		}
	}()
	// Wait for shutdownChan to be closed
	<-vm.shutdownChan
}

// readLastAccepted reads the last accepted hash from [acceptedBlockDB] and returns the
// last accepted block hash and height by reading directly from [vm.chaindb] instead of relying
// on [chain].
// Note: assumes [vm.chaindb] and [vm.genesisHash] have been initialized.
func (vm *VM) readLastAccepted() (common.Hash, uint64, error) {
	// Attempt to load last accepted block to determine if it is necessary to
	// initialize state with the genesis block.
	lastAcceptedBytes, lastAcceptedErr := vm.acceptedBlockDB.Get(lastAcceptedKey)
	switch {
	case lastAcceptedErr == database.ErrNotFound:
		// If there is nothing in the database, return the genesis block hash and height
		return vm.genesisHash, 0, nil
	case lastAcceptedErr != nil:
		return common.Hash{}, 0, fmt.Errorf("failed to get last accepted block ID due to: %w", lastAcceptedErr)
	case len(lastAcceptedBytes) != common.HashLength:
		return common.Hash{}, 0, fmt.Errorf("last accepted bytes should have been length %d, but found %d", common.HashLength, len(lastAcceptedBytes))
	default:
		lastAcceptedHash := common.BytesToHash(lastAcceptedBytes)
		height := rawdb.ReadHeaderNumber(vm.chaindb, lastAcceptedHash)
		if height == nil {
			return common.Hash{}, 0, fmt.Errorf("failed to retrieve header number of last accepted block: %s", lastAcceptedHash)
		}
		return lastAcceptedHash, *height, nil
	}
}

// attachEthService registers the backend RPC services provided by Ethereum
// to the provided handler under their assigned namespaces.
func attachEthService(handler *rpc.Server, apis []rpc.API, names []string) error {
	enabledServicesSet := make(map[string]struct{})
	for _, ns := range names {
		// handle pre geth v1.10.20 api names as aliases for their updated values
		// to allow configurations to be backwards compatible.
		if newName, isLegacy := legacyApiNames[ns]; isLegacy {
			log.Info("deprecated api name referenced in configuration.", "deprecated", ns, "new", newName)
			enabledServicesSet[newName] = struct{}{}
			continue
		}

		enabledServicesSet[ns] = struct{}{}
	}

	apiSet := make(map[string]rpc.API)
	for _, api := range apis {
		if existingAPI, exists := apiSet[api.Name]; exists {
			return fmt.Errorf("duplicated API name: %s, namespaces %s and %s", api.Name, api.Namespace, existingAPI.Namespace)
		}
		apiSet[api.Name] = api
	}

	for name := range enabledServicesSet {
		api, exists := apiSet[name]
		if !exists {
			return fmt.Errorf("API service %s not found", name)
		}
		if err := handler.RegisterName(api.Namespace, api.Service); err != nil {
			return err
		}
	}

	return nil
}

func (vm *VM) Connected(ctx context.Context, nodeID ids.NodeID, version *version.Application) error {
	if err := vm.validatorsManager.Connect(nodeID); err != nil {
		return fmt.Errorf("uptime manager failed to connect node %s: %w", nodeID, err)
	}
	return vm.Network.Connected(ctx, nodeID, version)
}

func (vm *VM) Disconnected(ctx context.Context, nodeID ids.NodeID) error {
	if err := vm.validatorsManager.Disconnect(nodeID); err != nil {
		return fmt.Errorf("uptime manager failed to disconnect node %s: %w", nodeID, err)
	}

	return vm.Network.Disconnected(ctx, nodeID)
}
