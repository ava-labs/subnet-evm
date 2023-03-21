// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package runner

import (
	"context"
	"fmt"
	"os"
	"time"

	runner_sdk "github.com/ava-labs/avalanche-network-runner/client"
	"github.com/ava-labs/avalanche-network-runner/rpcpb"
	runner_server "github.com/ava-labs/avalanche-network-runner/server"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/ava-labs/avalanchego/utils/wrappers"
	"github.com/ava-labs/subnet-evm/plugin/evm"
	"github.com/ethereum/go-ethereum/log"
)

type Subnet struct {
	// SubnetID is the txID of the transaction that created the subnet
	SubnetID ids.ID
	// Current ANR assumes one blockchain per subnet, so we have a single blockchainID here
	BlockchainID ids.ID
	// ValidatorURIs is the base URIs for each participant of the Subnet
	ValidatorURIs []string
}

type ANRConfig struct {
	LogLevel            string
	AvalancheGoExecPath string
	PluginDir           string
	GlobalNodeConfig    string
}

type NetworkManager struct {
	ANRConfig ANRConfig

	// Map SubnetID to Subnet details
	subnets []*Subnet

	logFactory      logging.Factory
	anrClient       runner_sdk.Client
	anrServer       runner_server.Server
	done            chan struct{}
	serverCtxCancel context.CancelFunc
}

func NewDefaultANRConfig() ANRConfig {
	defaultConfig := ANRConfig{
		LogLevel:            "info",
		AvalancheGoExecPath: os.ExpandEnv("$GOPATH/src/github.com/ava-labs/avalanchego/build/avalanchego"),
		PluginDir:           os.ExpandEnv("$GOPATH/src/github.com/ava-labs/avalanchego/build/plugins"),
		GlobalNodeConfig: `{
			"log-display-level":"info",
			"proposervm-use-current-height":true,
			"throttler-inbound-validator-alloc-size":"107374182",
			"throttler-inbound-node-max-processing-msgs":"100000",
			"throttler-inbound-bandwidth-refill-rate":"1073741824",
			"throttler-inbound-bandwidth-max-burst-size":"1073741824",
			"throttler-inbound-cpu-validator-alloc":"100000",
			"throttler-inbound-disk-validator-alloc":"10737418240000",
			"throttler-outbound-validator-alloc-size":"107374182"
		}`,
	}
	// If AVALANCHEGO_BUILD_PATH is populated, override location set by GOPATH
	if envBuildPath, exists := os.LookupEnv("AVALANCHEGO_BUILD_PATH"); exists {
		defaultConfig.AvalancheGoExecPath = fmt.Sprintf("%s/avalanchego", envBuildPath)
		defaultConfig.PluginDir = fmt.Sprintf("%s/plugins", envBuildPath)
	}
	return defaultConfig
}

// NewNetworkManager constructs a new instance of a network manager
func NewNetworkManager(config ANRConfig) *NetworkManager {
	manager := &NetworkManager{
		ANRConfig: config,
	}

	logLevel, err := logging.ToLevel(config.LogLevel)
	if err != nil {
		panic(fmt.Errorf("invalid ANR log level: %w", err))
	}
	manager.logFactory = logging.NewFactory(logging.Config{
		DisplayLevel: logLevel,
		LogLevel:     logLevel,
	})

	return manager
}

func (n *NetworkManager) startServer(ctx context.Context) (<-chan struct{}, error) {
	done := make(chan struct{})
	zapServerLog, err := n.logFactory.Make("server")
	if err != nil {
		return nil, fmt.Errorf("failed to make server log: %w", err)
	}

	n.anrServer, err = runner_server.New(
		runner_server.Config{
			Port:                ":12352",
			GwPort:              ":12353",
			GwDisabled:          false,
			DialTimeout:         10 * time.Second,
			RedirectNodesOutput: true,
			SnapshotsDir:        "",
		},
		zapServerLog,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start ANR server: %w", err)
	}
	n.done = done

	// Use a separate background context here, since the server should only be canceled by explicit shutdown
	serverCtx, serverCtxCancel := context.WithCancel(context.Background())
	n.serverCtxCancel = serverCtxCancel
	go func() {
		if err := n.anrServer.Run(serverCtx); err != nil {
			log.Error("Error shutting down ANR server", "err", err)
		} else {
			log.Info("Terminating ANR Server")
		}
		close(done)
	}()

	return done, nil
}

// This is an ugly hack to redial the server and create a new client connection.
// This is used to support tearing down the network from an external command.
func (n *NetworkManager) startClient() error {
	logLevel, err := logging.ToLevel(n.ANRConfig.LogLevel)
	if err != nil {
		return fmt.Errorf("failed to parse ANR log level: %w", err)
	}
	logFactory := logging.NewFactory(logging.Config{
		DisplayLevel: logLevel,
		LogLevel:     logLevel,
	})
	zapLog, err := logFactory.Make("main")
	if err != nil {
		return fmt.Errorf("failed to make client log: %w", err)
	}

	n.anrClient, err = runner_sdk.New(runner_sdk.Config{
		Endpoint:    "0.0.0.0:12352",
		DialTimeout: 10 * time.Second,
	}, zapLog)
	if err != nil {
		return fmt.Errorf("failed to start ANR client: %w", err)
	}

	return nil
}

func (n *NetworkManager) initServer() error {
	if n.anrServer != nil {
		return nil
	}

	_, err := n.startServer(context.Background())
	return err
}

func (n *NetworkManager) initClient() error {
	if n.anrClient != nil {
		return nil
	}

	return n.startClient()
}

func (n *NetworkManager) init() error {
	if err := n.initServer(); err != nil {
		return err
	}
	return n.initClient()
}

// StartDefaultNetwork constructs a default 5 node network.
func (n *NetworkManager) StartDefaultNetwork(ctx context.Context) (<-chan struct{}, error) {
	if err := n.init(); err != nil {
		return nil, err
	}

	log.Info("Sending 'start'", "AvalancheGoExecPath", n.ANRConfig.AvalancheGoExecPath)

	// Start cluster
	resp, err := n.anrClient.Start(
		ctx,
		n.ANRConfig.AvalancheGoExecPath,
		runner_sdk.WithPluginDir(n.ANRConfig.PluginDir),
		runner_sdk.WithGlobalNodeConfig(n.ANRConfig.GlobalNodeConfig),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start ANR network: %w", err)
	}
	log.Info("successfully started cluster", "RootDataDir", resp.ClusterInfo.RootDataDir, "Subnets", resp.GetClusterInfo().GetSubnets())
	return n.done, nil
}

// SetupNetwork constructs blockchains with the given [blockchainSpecs] and adds them to the network manager.
// Uses [execPath] as the AvalancheGo binary execution path for any started nodes.
func (n *NetworkManager) SetupNetwork(ctx context.Context, execPath string, blockchainSpecs []*rpcpb.BlockchainSpec) error {
	if err := n.init(); err != nil {
		return err
	}
	sresp, err := n.anrClient.CreateSpecificBlockchains(
		ctx,
		execPath,
		blockchainSpecs,
	)
	if err != nil {
		return fmt.Errorf("failed to create blockchains: %w", err)
	}

	// TODO: network runner health should imply custom VM healthiness
	// or provide a separate API for custom VM healthiness
	// "start" is async, so wait some time for cluster health
	log.Info("waiting for all VMs to report healthy", "VMID", evm.ID)
	for {
		v, err := n.anrClient.Health(ctx)
		log.Info("Pinged CLI Health", "result", v, "err", err)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		// TODO: clean this up
		if err != nil {
			return fmt.Errorf("failed to ping ANR health: %w", err)
		}
		break
	}

	cctx, ccancel := context.WithTimeout(ctx, 2*time.Minute)
	status, err := n.anrClient.Status(cctx)
	ccancel()
	if err != nil {
		return fmt.Errorf("failed to get ANR status: %w", err)
	}
	nodeInfos := status.GetClusterInfo().GetNodeInfos()

	for i, chainSpec := range blockchainSpecs {
		blockchainIDStr := sresp.Chains[i]
		blockchainID, err := ids.FromString(blockchainIDStr)
		if err != nil {
			panic(err)
		}
		subnetIDStr := sresp.ClusterInfo.CustomChains[blockchainIDStr].SubnetId
		subnetID, err := ids.FromString(subnetIDStr)
		if err != nil {
			panic(err)
		}
		subnet := &Subnet{
			SubnetID:     subnetID,
			BlockchainID: blockchainID,
		}
		for _, nodeName := range chainSpec.Participants {
			subnet.ValidatorURIs = append(subnet.ValidatorURIs, nodeInfos[nodeName].Uri)
		}
		n.subnets = append(n.subnets, subnet)
	}

	return nil
}

// TeardownNetwork tears down the network constructed by the network manager and cleans up
// everything associated with it.
func (n *NetworkManager) TeardownNetwork() error {
	if err := n.initClient(); err != nil {
		return err
	}
	errs := wrappers.Errs{}
	log.Info("Shutting down cluster")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	_, err := n.anrClient.Stop(ctx)
	cancel()
	errs.Add(err)
	errs.Add(n.anrClient.Close())
	if n.serverCtxCancel != nil {
		n.serverCtxCancel()
	}
	return errs.Err
}

// CloseClient closes the connection between the ANR client and server without terminating the
// running network.
func (n *NetworkManager) CloseClient() error {
	if n.anrClient == nil {
		return nil
	}
	err := n.anrClient.Close()
	n.anrClient = nil
	return err
}

// GetSubnets returns the IDs of the currently running subnets
func (n *NetworkManager) GetSubnets() []ids.ID {
	subnetIDs := make([]ids.ID, 0, len(n.subnets))
	for _, subnet := range n.subnets {
		subnetIDs = append(subnetIDs, subnet.SubnetID)
	}
	return subnetIDs
}

// GetSubnet retrieves the subnet details for the requested subnetID
func (n *NetworkManager) GetSubnet(subnetID ids.ID) (*Subnet, bool) {
	for _, subnet := range n.subnets {
		if subnet.SubnetID == subnetID {
			return subnet, true
		}
	}
	return nil, false
}
