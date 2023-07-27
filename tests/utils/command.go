// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ava-labs/avalanchego/api/health"
	"github.com/ethereum/go-ethereum/log"
	"github.com/go-cmd/cmd"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

const (
	// Timeout to boot the AvalancheGo node
	bootAvalancheNodeTimeout = 5 * time.Minute

	// Timeout for the health API to check the AvalancheGo is ready
	healthCheckTimeout = 5 * time.Second
)

// GetSynchronizedBlockchainID is a function that returns the blockchain ID for a given alias.
// This function is set by CreateSubnetsSynchronized() and is used by the tests to get the blockchain ID
// for a given test case.
var GetSynchronizedBlockchainID func(alias string) string

// RunCommand starts the command [bin] with the given [args] and returns the command to the caller
// TODO cmd package mentions we can do this more efficiently with cmd.NewCmdOptions rather than looping
// and calling Status().
func RunCommand(bin string, args ...string) (*cmd.Cmd, error) {
	log.Info("Executing", "cmd", fmt.Sprintf("%s %s", bin, strings.Join(args, " ")))

	curCmd := cmd.NewCmd(bin, args...)
	_ = curCmd.Start()

	// to stream outputs
	ticker := time.NewTicker(10 * time.Millisecond)
	go func() {
		prevLine := ""
		for range ticker.C {
			status := curCmd.Status()
			n := len(status.Stdout)
			if n == 0 {
				continue
			}

			line := status.Stdout[n-1]
			if prevLine != line && line != "" {
				fmt.Println("[streaming output]", line)
			}

			prevLine = line
		}
	}()

	return curCmd, nil
}

func RegisterPingTest() {
	ginkgo.It("ping the network", ginkgo.Label("ping"), func() {
		client := health.NewClient(DefaultLocalNodeURI)
		healthy, err := client.Readiness(context.Background(), nil)
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(healthy.Healthy).Should(gomega.BeTrue())
	})
}

// RegisterNodeRun registers a before suite that starts an AvalancheGo process to use for the e2e tests
// and an after suite that stops the AvalancheGo process
func RegisterNodeRun() {
	// BeforeSuite starts an AvalancheGo process to use for the e2e tests
	var startCmd *cmd.Cmd
	_ = ginkgo.BeforeSuite(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		wd, err := os.Getwd()
		gomega.Expect(err).Should(gomega.BeNil())
		log.Info("Starting AvalancheGo node", "wd", wd)
		cmd, err := RunCommand("./scripts/run.sh")
		startCmd = cmd
		gomega.Expect(err).Should(gomega.BeNil())

		// Assumes that startCmd will launch a node with HTTP Port at [utils.DefaultLocalNodeURI]
		healthClient := health.NewClient(DefaultLocalNodeURI)
		healthy, err := health.AwaitReady(ctx, healthClient, 5*time.Second, nil)
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(healthy).Should(gomega.BeTrue())
		log.Info("AvalancheGo node is healthy")
	})

	ginkgo.AfterSuite(func() {
		gomega.Expect(startCmd).ShouldNot(gomega.BeNil())
		gomega.Expect(startCmd.Stop()).Should(gomega.BeNil())
		// TODO add a new node to bootstrap off of the existing node and ensure it can bootstrap all subnets
		// created during the test
	})
}

// CreateSubnetsSynchronized creates subnets for given [genesisFiles], and registers a before suite that starts an AvalancheGo process to use for the e2e tests.
// genesisFiles is a map of test aliases to genesis file paths.
func CreateSubnetsSynchronized(genesisFiles map[string]string) {
	// Keep track of the AvalancheGo external bash script, it is null for most
	// processes except the first process that starts AvalancheGo
	var startCmd *cmd.Cmd

	// Our test suite runs in separate processes, ginkgo has
	// SynchronizedBeforeSuite() which runs once, and its return value is passed
	// over to each worker.
	//
	// Here an AvalancheGo node instance is started, and subnets are created for
	// each test case. Each test case has its own subnet, therefore all tests
	// can run in parallel without any issue.
	//
	var _ = ginkgo.SynchronizedBeforeSuite(func() []byte {
		ctx, cancel := context.WithTimeout(context.Background(), bootAvalancheNodeTimeout)
		defer cancel()

		wd, err := os.Getwd()
		gomega.Expect(err).Should(gomega.BeNil())
		log.Info("Starting AvalancheGo node", "wd", wd)
		cmd, err := RunCommand("./scripts/run.sh")
		startCmd = cmd
		gomega.Expect(err).Should(gomega.BeNil())

		// Assumes that startCmd will launch a node with HTTP Port at [utils.DefaultLocalNodeURI]
		healthClient := health.NewClient(DefaultLocalNodeURI)
		healthy, err := health.AwaitReady(ctx, healthClient, healthCheckTimeout, nil)
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(healthy).Should(gomega.BeTrue())
		log.Info("AvalancheGo node is healthy")

		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		blockchainIDs := make(map[string]string)
		for alias, file := range genesisFiles {
			blockchainIDs[alias] = CreateNewSubnet(ctx, file)
		}

		blockchainIDsBytes, err := json.Marshal(blockchainIDs)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		return blockchainIDsBytes
	}, func(ctx ginkgo.SpecContext, data []byte) {
		blockchainIDs := make(map[string]string)
		err := json.Unmarshal(data, &blockchainIDs)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		GetSynchronizedBlockchainID = func(alias string) string {
			return blockchainIDs[alias]
		}
	})

	// SynchronizedAfterSuite() takes two functions, the first runs after each test suite is done and the second
	// function is executed once when all the tests are done. This function is used
	// to gracefully shutdown the AvalancheGo node.
	var _ = ginkgo.SynchronizedAfterSuite(func() {}, func() {
		gomega.Expect(startCmd).ShouldNot(gomega.BeNil())
		gomega.Expect(startCmd.Stop()).Should(gomega.BeNil())
	})
}

// GetFilesAndAliases returns a map of aliases to file paths in given [dir].
func GetFilesAndAliases(dir string) (map[string]string, error) {
	files, err := filepath.Glob(dir)
	if err != nil {
		return nil, err
	}
	aliasesToFiles := make(map[string]string)
	for _, file := range files {
		alias := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		aliasesToFiles[alias] = file
	}
	return aliasesToFiles, nil
}
