package evm

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/ava-labs/avalanchego/api/keystore"
	"github.com/ava-labs/avalanchego/chains/atomic"
	"github.com/ava-labs/avalanchego/database/manager"
	"github.com/ava-labs/avalanchego/database/prefixdb"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	commonEng "github.com/ava-labs/avalanchego/snow/engine/common"
	"github.com/ava-labs/avalanchego/snow/validators"
	avalancheConstants "github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/ava-labs/avalanchego/version"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/stretchr/testify/require"
)

var (
	testNetworkID uint32 = 10
	testCChainID         = ids.ID{'c', 'c', 'h', 'a', 'i', 'n', 't', 'e', 's', 't'}
	testXChainID         = ids.ID{'t', 'e', 's', 't', 'x'}

	testMinGasPrice int64 = 225_000_000_000
	testAvaxAssetID       = ids.ID{1, 2, 3}
	username              = "Johns"
	password              = "CjasdjhiPeirbSenfeI13" // #nosec G101

	genesisJSONSubnetEVM    = "{\"config\":{\"chainId\":43111,\"homesteadBlock\":0,\"eip150Block\":0,\"eip150Hash\":\"0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0\",\"eip155Block\":0,\"eip158Block\":0,\"byzantiumBlock\":0,\"constantinopleBlock\":0,\"petersburgBlock\":0,\"istanbulBlock\":0,\"muirGlacierBlock\":0,\"subnetEVMTimestamp\":0},\"nonce\":\"0x0\",\"timestamp\":\"0x0\",\"extraData\":\"0x00\",\"gasLimit\":\"0x7A1200\",\"difficulty\":\"0x0\",\"mixHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"coinbase\":\"0x0000000000000000000000000000000000000000\",\"alloc\":{\"0x71562b71999873DB5b286dF957af199Ec94617F7\": {\"balance\":\"0x4192927743b88000\"}, \"0x703c4b2bD70c169f5717101CaeE543299Fc946C7\": {\"balance\":\"0x4192927743b88000\"}},\"number\":\"0x0\",\"gasUsed\":\"0x0\",\"parentHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\"}"
	genesisJSONDUpgrade     = "{\"config\":{\"chainId\":43111,\"homesteadBlock\":0,\"eip150Block\":0,\"eip150Hash\":\"0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0\",\"eip155Block\":0,\"eip158Block\":0,\"byzantiumBlock\":0,\"constantinopleBlock\":0,\"petersburgBlock\":0,\"istanbulBlock\":0,\"muirGlacierBlock\":0,\"subnetEVMTimestamp\":0,\"dUpgradeTimestamp\":0},\"nonce\":\"0x0\",\"timestamp\":\"0x0\",\"extraData\":\"0x00\",\"gasLimit\":\"0x7A1200\",\"difficulty\":\"0x0\",\"mixHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"coinbase\":\"0x0000000000000000000000000000000000000000\",\"alloc\":{\"0x71562b71999873DB5b286dF957af199Ec94617F7\": {\"balance\":\"0x4192927743b88000\"}, \"0x703c4b2bD70c169f5717101CaeE543299Fc946C7\": {\"balance\":\"0x4192927743b88000\"}},\"number\":\"0x0\",\"gasUsed\":\"0x0\",\"parentHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\"}"
	genesisJSONPreSubnetEVM = "{\"config\":{\"chainId\":43111,\"homesteadBlock\":0,\"eip150Block\":0,\"eip150Hash\":\"0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0\",\"eip155Block\":0,\"eip158Block\":0,\"byzantiumBlock\":0,\"constantinopleBlock\":0,\"petersburgBlock\":0,\"istanbulBlock\":0,\"muirGlacierBlock\":0},\"nonce\":\"0x0\",\"timestamp\":\"0x0\",\"extraData\":\"0x00\",\"gasLimit\":\"0x7A1200\",\"difficulty\":\"0x0\",\"mixHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"coinbase\":\"0x0000000000000000000000000000000000000000\",\"alloc\":{\"0x71562b71999873DB5b286dF957af199Ec94617F7\": {\"balance\":\"0x4192927743b88000\"}, \"0x703c4b2bD70c169f5717101CaeE543299Fc946C7\": {\"balance\":\"0x4192927743b88000\"}},\"number\":\"0x0\",\"gasUsed\":\"0x0\",\"parentHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\"}"
	genesisJSONLatest       = genesisJSONDUpgrade

	defaultMaxWorkers = 5
)

var (
	noTestVar = errors.New("testing variable nil")
)

type VmTestManagerConfig struct {
	IsE2E bool `json:"IsE2E"`
	maxWorkers int
}

type GenesisVMConfig struct {
	FinishBootstrapping bool   `json:"finishBootstrapping"`
	GenesisJSON         string `json:"GenesisJSON"`
	ConfigJSON          string `json:"ConfigJSON"`
	UpgradeJSON         string `json:"UpgradeJSON"`
}

type testVector struct {
	test func(t *testing.T, c VmTestManager)
	expectedResult error
}

type VmTestManager interface {
	RunVectors(vecs []testVector) error
	Create(t *testing.T, config GenesisVMConfig,) (VMWorker, error)
}

type vmTestManager struct {
	config VmTestManagerConfig
}

func NewVmTestManager(config VmTestManagerConfig) (VmTestManager, error) {
	manager := &vmTestManager{
		config: config,
	}

	if config.IsE2E {
		//launch avalanche network runner
	}
	return manager, nil
}

func (v *vmTestManager) Create(t *testing.T, config GenesisVMConfig) (VMWorker, error) {
	worker := &vmWorker{
		workerConfig: config,
	}

	if v.config.IsE2E {
		//
	} else {
		if t == nil {
			return nil, noTestVar
		}

		issuer, vm, dbmanager, appsender := GenesisVM(
			t, config.FinishBootstrapping, config.GenesisJSON, config.ConfigJSON, config.UpgradeJSON)
		worker.issuer = issuer
		worker.vm = vm
		worker.dbManager = dbmanager
		worker.appsender = appsender
	}

	return worker, nil
}

func (v *vmTestManager) RunVectors(vecs []testVector) error {
	return nil
}

type VMWorker interface {
	GetVMConfig() (Config, error)
	IssueTxs(txs []*types.Transaction) error
	ConfirmTxs(txs []*types.Transaction) error
	Shutdown(ctx context.Context) error
}

type vmWorker struct {
	workerConfig GenesisVMConfig
	isE2E bool

	//only e2e

	//only unit
	issuer    chan commonEng.Message
	vm        *VM
	dbManager manager.Manager
	appsender *commonEng.SenderTest
}

func (v *vmWorker) GetVMConfig() (Config, error) {
	var (
		config Config
		err error = nil
	)
	if v.isE2E {
	} else {
		config = v.vm.config
	}
	return config, err
}

func (v *vmWorker) IssueTxs(txs []*types.Transaction) error {
	if v.isE2E {
	} else {

	}
	return nil
}
func (v *vmWorker) ConfirmTxs(txs []*types.Transaction) error {
	if v.isE2E {

	} else {
			
	}
	return nil
}
func (v *vmWorker) Shutdown(ctx context.Context) error {
	if v.isE2E {

	} else {
			
	}
	return nil
}

// GenesisVM creates a VM instance with the genesis test bytes and returns
// the channel use to send messages to the engine, the VM, database manager,
// and sender.
// If [genesisJSON] is empty, defaults to using [genesisJSONLatest]
func GenesisVM(t *testing.T,
	finishBootstrapping bool,
	genesisJSON string,
	configJSON string,
	upgradeJSON string,
) (chan commonEng.Message,
	*VM,
	manager.Manager,
	*commonEng.SenderTest,
) {
	vm := &VM{}
	ctx, dbManager, genesisBytes, issuer, _ := setupGenesis(t, genesisJSON)
	appSender := &commonEng.SenderTest{T: t}
	appSender.CantSendAppGossip = true
	appSender.SendAppGossipF = func(context.Context, []byte) error { return nil }
	err := vm.Initialize(
		context.Background(),
		ctx,
		dbManager,
		genesisBytes,
		[]byte(upgradeJSON),
		[]byte(configJSON),
		issuer,
		[]*commonEng.Fx{},
		appSender,
	)
	require.NoError(t, err, "error initializing GenesisVM")

	if finishBootstrapping {
		require.NoError(t, vm.SetState(context.Background(), snow.Bootstrapping))
		require.NoError(t, vm.SetState(context.Background(), snow.NormalOp))
	}

	return issuer, vm, dbManager, appSender
}

// If [genesisJSON] is empty, defaults to using [genesisJSONLatest]
func setupGenesis(
	t *testing.T,
	genesisJSON string,
) (*snow.Context,
	manager.Manager,
	[]byte,
	chan commonEng.Message,
	*atomic.Memory,
) {
	if len(genesisJSON) == 0 {
		genesisJSON = genesisJSONLatest
	}
	genesisBytes := buildGenesisTest(t, genesisJSON)
	ctx := NewContext()

	baseDBManager := manager.NewMemDB(&version.Semantic{
		Major: 1,
		Minor: 4,
		Patch: 5,
	})

	// initialize the atomic memory
	atomicMemory := atomic.NewMemory(prefixdb.New([]byte{0}, baseDBManager.Current().Database))
	ctx.SharedMemory = atomicMemory.NewSharedMemory(ctx.ChainID)

	// NB: this lock is intentionally left locked when this function returns.
	// The caller of this function is responsible for unlocking.
	ctx.Lock.Lock()

	userKeystore := keystore.New(logging.NoLog{}, manager.NewMemDB(&version.Semantic{
		Major: 1,
		Minor: 4,
		Patch: 5,
	}))
	if err := userKeystore.CreateUser(username, password); err != nil {
		t.Fatal(err)
	}
	ctx.Keystore = userKeystore.NewBlockchainKeyStore(ctx.ChainID)

	issuer := make(chan commonEng.Message, 1)
	prefixedDBManager := baseDBManager.NewPrefixDBManager([]byte{1})
	return ctx, prefixedDBManager, genesisBytes, issuer, atomicMemory
}

// BuildGenesisTest returns the genesis bytes for Subnet EVM VM to be used in testing
func buildGenesisTest(t *testing.T, genesisJSON string) []byte {
	ss := CreateStaticService()

	genesis := &core.Genesis{}
	if err := json.Unmarshal([]byte(genesisJSON), genesis); err != nil {
		t.Fatalf("Problem unmarshaling genesis JSON: %s", err)
	}
	args := &BuildGenesisArgs{GenesisData: genesis}
	reply := &BuildGenesisReply{}
	err := ss.BuildGenesis(nil, args, reply)
	if err != nil {
		t.Fatalf("Failed to create test genesis")
	}
	genesisBytes, err := formatting.Decode(reply.Encoding, reply.GenesisBytes)
	if err != nil {
		t.Fatalf("Failed to decode genesis bytes: %s", err)
	}
	return genesisBytes
}

func NewContext() *snow.Context {
	ctx := snow.DefaultContextTest()
	ctx.NetworkID = testNetworkID
	ctx.NodeID = ids.GenerateTestNodeID()
	ctx.ChainID = testCChainID
	ctx.AVAXAssetID = testAvaxAssetID
	ctx.XChainID = testXChainID
	aliaser := ctx.BCLookup.(ids.Aliaser)
	_ = aliaser.Alias(testCChainID, "C")
	_ = aliaser.Alias(testCChainID, testCChainID.String())
	_ = aliaser.Alias(testXChainID, "X")
	_ = aliaser.Alias(testXChainID, testXChainID.String())
	ctx.ValidatorState = &validators.TestState{
		GetSubnetIDF: func(_ context.Context, chainID ids.ID) (ids.ID, error) {
			subnetID, ok := map[ids.ID]ids.ID{
				avalancheConstants.PlatformChainID: avalancheConstants.PrimaryNetworkID,
				testXChainID:                       avalancheConstants.PrimaryNetworkID,
				testCChainID:                       avalancheConstants.PrimaryNetworkID,
			}[chainID]
			if !ok {
				return ids.Empty, errors.New("unknown chain")
			}
			return subnetID, nil
		},
	}
	blsSecretKey, err := bls.NewSecretKey()
	if err != nil {
		panic(err)
	}
	ctx.WarpSigner = avalancheWarp.NewSigner(blsSecretKey, ctx.NetworkID, ctx.ChainID)
	ctx.PublicKey = bls.PublicFromSecretKey(blsSecretKey)
	return ctx
}
