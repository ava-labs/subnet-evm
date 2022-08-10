// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package e2e

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ava-labs/avalanche-network-runner/api"
	"github.com/ava-labs/avalanche-network-runner/client"
	"github.com/ava-labs/avalanche-network-runner/rpcpb"
	"github.com/ava-labs/avalanche-network-runner/server"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/subnet-evm/plugin/evm"
	"github.com/ava-labs/subnet-evm/tests/e2e/helpers"
	"github.com/ava-labs/subnet-evm/tests/e2e/runner"
	"github.com/ava-labs/subnet-evm/tests/e2e/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/formatter"
	"github.com/onsi/gomega"
)

const (
	defaultStepTimeout   = 5 * time.Minute
	healthPollInterval   = 10 * time.Second
	gasLimit             = uint64(8_000_000)
	stateSyncGenesisPath = "./tests/e2e/genesis/state_sync.json"
)

// startSubnet starts a test network and launches a subnetEVM instance with the genesis file at [genesisPath]
func startSubnet(genesisPath string) error {
	_, err := runner.StartNetwork(evm.ID, vmName, genesisPath, utils.GetPluginDir())
	gomega.Expect(err).Should(gomega.BeNil())
	return utils.UpdateHardhatConfig()
}

// stopSubnet stops the test network.
func stopSubnet() {
	err := runner.StopNetwork()
	gomega.Expect(err).Should(gomega.BeNil())
}

var _ = ginkgo.Describe("[state-sync]", func() {
	var (
		cli     client.Client
		cluster *rpcpb.ClusterInfo
	)

	ginkgo.BeforeEach(func() {
		//		var err error
		//		cli, err = client.New(client.Config{
		//			LogLevel:    networkRunnerLogLevel,
		//			Endpoint:    gRPCEp,
		//			DialTimeout: 10 * time.Second,
		//		})
		//		gomega.Expect(err).Should(gomega.BeNil())
		//		outf("{{green}}sending 'start' with binary path:{{/}} %q\n", execPath)
		//		ctx, cancel := createDefaultCtx()
		//		defer cancel()
		//		commitInterval := 16
		//		cConfigJson := fmt.Sprintf(`{"log-level":"info", "commit-interval": %d, "state-sync-commit-interval": %d}`, commitInterval, commitInterval)
		//		_, err = cli.Start(
		//			ctx,
		//			execPath,
		//			client.WithChainConfigs(map[string]string{"C": cConfigJson}),
		//		)
		//
		//		gomega.Expect(err).Should(gomega.BeNil())
		//		outf("{{green}}polling for cluster health\n")
		//		cluster = awaitHealthy(cli)
		//		outf("{{green}}successfully started:{{/}} %+v\n", cluster.NodeNames)
		outf("{{green}}starting cluster{{/}}\n")
		err := startSubnet(stateSyncGenesisPath)
		gomega.Expect(err).Should(gomega.BeNil())

		// initialize test variables
		cli = runner.GetClient()
		cluster = awaitHealthy(cli)
		outf("{{green}}successfully started:{{/}} %+v\n", cluster.NodeNames)
	})
	ginkgo.It("can sync", func() {

		outf("{{blue}}generating state{{/}}\n")
		nodeInfo := cluster.NodeInfos[cluster.NodeNames[0]]

		var customChainID string
		outf("{{magenta}}custom chains:{{/}}\n")
		for chainID, chainInfo := range cluster.CustomChains {
			if chainInfo.VmId != evm.ID.String() {
				continue
			}
			customChainID = chainID
			outf("{{magenta}}ChainID: %v{{/}}\n", chainID)
		}
		ip, port := parseURI(nodeInfo.Uri)
		ethClient := api.NewEthClientWithChainID(ip, uint(port), customChainID)
		outf("{{red}}URI: %v{{/}}\n", nodeInfo.Uri)

		testSetup, err := generateBlocks(ethClient)
		gomega.Expect(err).Should(gomega.BeNil())

		cConfigJson := `{"log-level":"info", "state-sync-enabled": true, "state-sync-min-blocks": 1}`
		ctx, cancel := createDefaultCtx()
		defer cancel()
		_, err = cli.AddNode(
			ctx,
			"sync",
			execPath,
			client.WithPluginDir(utils.GetPluginDir()),
			client.WithBlockchainSpecs([]*rpcpb.BlockchainSpec{
				{
					VmName:  vmName,
					Genesis: stateSyncGenesisPath,
				},
			}),
			client.WithChainConfigs(map[string]string{customChainID: cConfigJson}),
		)
		gomega.Expect(err).Should(gomega.BeNil())
		cluster = awaitHealthy(cli)

		// check some sutff on the newly synced node
		syncedNodeIP, syncedNodePort := parseURI(cluster.NodeInfos["sync"].Uri)
		syncedNodeClient := api.NewEthClientWithChainID(syncedNodeIP, uint(syncedNodePort), customChainID)
		err = checkSyncedClient(syncedNodeClient, testSetup)
		gomega.Expect(err).Should(gomega.BeNil())

		outf("{{blue}}---------------------{{/}}\n")
		outf("{{blue}}   sync successful   {{/}}\n")
		outf("{{blue}}---------------------{{/}}\n")

	})
})

func checkSyncedClient(client api.EthClient, testSetup *testSetup) error {
	ctx, cancel := createDefaultCtx()
	defer cancel()
	block, err := client.BlockNumber(ctx)
	if err != nil {
		return err
	}
	expectedBlock := uint64(32)
	if block != expectedBlock {
		return fmt.Errorf("syncedNode block (%d) does not match expected (%d)", block, expectedBlock)
	}
	if balance, err := client.BalanceAt(ctx, testSetup.checkingAddr, nil); err != nil {
		return fmt.Errorf("error obtaining balance: %w", err)
	} else if balance.Cmp(testSetup.expectedBalance) != 0 {
		return fmt.Errorf("invalid return BalanceAt. expected %s got %s", testSetup.expectedBalance, balance)
	} else {
		outf("{{green}}checking balance (synced node) %v{{/}}\n", balance)
	}

	return nil
}

func parseURI(uri string) (string, uint16) {
	uri = strings.TrimPrefix(uri, "http://")
	parts := strings.Split(uri, ":")
	port, err := strconv.ParseUint(parts[1], 10, 16)
	gomega.Expect(err).Should(gomega.BeNil())
	return parts[0], uint16(port)
}

func awaitHealthy(cli client.Client) *rpcpb.ClusterInfo {
	for {
		time.Sleep(healthPollInterval)
		ctx, cancel := createDefaultCtx()
		resp, err := cli.Health(ctx)
		cancel()
		if errors.Is(err, server.ErrNotBootstrapped) {
			outf("{{yellow}}still waiting...{{/}}")
			continue
		}
		gomega.Expect(err).Should(gomega.BeNil())

		if !resp.ClusterInfo.Healthy {
			outf("{{yellow}}still waiting (main chains)...{{/}}")
			continue
		}
		if !resp.ClusterInfo.CustomChainsHealthy {
			outf("{{yellow}}still waiting (custom chains)...{{/}}")
			continue
		}

		return resp.ClusterInfo
	}
}

type testSetup struct {
	checkingAddr    common.Address
	expectedBalance *big.Int
	expectedXBals   map[string]uint64
}

func generateBlocks(client api.EthClient) (*testSetup, error) {
	ctx, cancel := createDefaultCtx()
	defer cancel()

	// Send funds from EWOQ addr (funded in genesis) to checkingAddr on custom chain
	key, err := crypto.HexToECDSA("56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027")
	if err != nil {
		panic(err)
	}
	senderKey := key
	outf("{{blue}}senderAddr: %v{{/}}\n", crypto.PubkeyToAddress(key.PublicKey))

	checkingAddr := common.HexToAddress(fmt.Sprintf("0x%s", ids.GenerateTestShortID().Hex()))
	nonce := uint64(0)
	transfer := big.NewInt(1_000_000)
	expectedBalance := new(big.Int)

	for i := 0; i < 32; i++ {
		gasPrice, err := client.SuggestGasPrice(ctx)
		if err != nil {
			return nil, err
		}
		gasPrice = helpers.MultiplyMaxGasPrice(gasPrice)

		err = helpers.AwaitedSendTransaction(
			ctx, client, senderKey, nonce, checkingAddr, transfer, nil, big.NewInt(99999), gasLimit,
			gasPrice,
		)
		if err != nil {
			return nil, err
		}

		expectedBalance = new(big.Int).Add(expectedBalance, transfer)
		nonce += 1

		if balance, err := client.BalanceAt(ctx, checkingAddr, nil); err != nil {
			return nil, fmt.Errorf("error obtaining balance: %w", err)
		} else if balance.Cmp(expectedBalance) != 0 {
			return nil, fmt.Errorf("invalid return BalanceAt. expected %s got %s", expectedBalance, balance)
		} else {
			outf("{{magenta}}checking balance: %v{{/}}\n", balance)
		}
	}

	outf("{{blue}}--------------------{{/}}\n")
	outf("{{blue}}  generated blocks  {{/}}\n")
	outf("{{blue}}--------------------{{/}}\n")
	return &testSetup{
		checkingAddr:    checkingAddr,
		expectedBalance: expectedBalance,
	}, nil
}

// Outputs to stdout.
//
// e.g.,
//   Out("{{green}}{{bold}}hi there %q{{/}}", "aa")
//   Out("{{magenta}}{{bold}}hi therea{{/}} {{cyan}}{{underline}}b{{/}}")
//
// ref.
// https://github.com/onsi/ginkgo/blob/v2.0.0/formatter/formatter.go#L52-L73
//
func outf(format string, args ...interface{}) {
	s := formatter.F(format, args...)
	fmt.Fprint(formatter.ColorableStdOut, s)
}

func createDefaultCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultStepTimeout)
}
