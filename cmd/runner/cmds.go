// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ava-labs/avalanche-network-runner/rpcpb"
	"github.com/ava-labs/subnet-evm/plugin/evm"
	"github.com/ava-labs/subnet-evm/tests/utils/runner"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "subnet-evm",
		Short:   "avalanche-network-runner wrapper for launching Subnet-EVM",
		Version: evm.Version,
	}
	config  = runner.NewDefaultANRConfig()
	manager = runner.NewNetworkManager(config)
)

func init() {
	rootCmd.AddCommand(commands()...)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Subnet-EVM Run failed %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func commands() []*cobra.Command {
	return []*cobra.Command{
		{
			Use:   "default",
			Short: "Start an empty default 5 node network.",
			RunE:  startDefaultNetworkFunc,
			Args:  cobra.ExactArgs(0),
		},
		{
			Use:   "stop",
			Short: "stop the running network",
			RunE:  stopNetwork,
			Args:  cobra.ExactArgs(0),
		},
		{
			Use:   "two",
			Short: "Start a network with 10 additional nodes with reigtered BLS keys 5 running two subnets",
			RunE:  startTwoSubnetNetwork,
			Args:  cobra.ExactArgs(0),
		},
	}
}

func startDefaultNetworkFunc(*cobra.Command, []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done, err := manager.StartDefaultNetwork(ctx)
	if err != nil {
		return err
	}

	if err := manager.CloseClient(); err != nil {
		return err
	}

	terminatedChan := awaitShutdown()
	select {
	case <-done:
	case <-terminatedChan:
	}
	return nil
}

func awaitShutdown() <-chan struct{} {
	done := make(chan struct{})
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	signal.Notify(signals, syscall.SIGTERM)

	go func() {
		<-signals
		close(done)
	}()
	return done
}

func startTwoSubnetNetwork(*cobra.Command, []string) error {
	// Name 10 new validators (which should have BLS key registered)
	subnetA := []string{}
	subnetB := []string{}
	for i := 1; i <= 10; i++ {
		n := fmt.Sprintf("node%d-bls", i)
		if i <= 5 {
			subnetA = append(subnetA, n)
		} else {
			subnetB = append(subnetB, n)
		}
	}

	ctx := context.Background()

	done, err := manager.StartDefaultNetwork(ctx)
	if err != nil {
		return err
	}
	err = manager.SetupNetwork(
		ctx,
		config.AvalancheGoExecPath,
		[]*rpcpb.BlockchainSpec{
			{
				VmName:       evm.IDStr,
				Genesis:      os.ExpandEnv("$GOPATH/src/github.com/ava-labs/subnet-evm/tests/precompile/genesis/fee_manager.json"),
				ChainConfig:  "",
				SubnetConfig: "",
				Participants: subnetA,
			},
			{
				VmName:       evm.IDStr,
				Genesis:      os.ExpandEnv("$GOPATH/src/github.com/ava-labs/subnet-evm/tests/precompile/genesis/fee_manager.json"),
				ChainConfig:  "",
				SubnetConfig: "",
				Participants: subnetB,
			},
		},
	)
	if err != nil {
		return err
	}

	if err := manager.CloseClient(); err != nil {
		return err
	}

	terminateChan := awaitShutdown()
	select {
	case <-done:
	case <-terminateChan:
	}
	return nil
}

func stopNetwork(*cobra.Command, []string) error {
	return manager.TeardownNetwork()
}
