// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements precompile upgrade tests for Tx Allow List, requires network-runner cluster.
package precompile_upgrade

import (
	"context"
	"encoding/json"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/plugin/evm"
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/tests/e2e/utils"

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

type chainConfig struct {
	evm.Config
}

type chainUpgradeConfig struct {
	// ref. https://docs.avax.network/subnets/subnet-upgrade#changing-subnet-configuration
	PrecompileUpgrades []*params.PrecompileUpgrade `json:"precompileUpgrades,omitempty"`
}

var _ = utils.DescribeLocal("[Precompile Upgrade]", func() {
	ginkgo.It("can upgrade for precompile", ginkgo.Label("precompile-upgrade"), func() {
		runnerCli := utils.GetClient()
		gomega.Expect(runnerCli).ShouldNot(gomega.BeNil())

		defaultChainCfg := chainConfig{}
		defaultChainCfg.Config.SetDefaults()
		defaultChainCfg.Config.LogJSONFormat = true
		defaultChainCfgBytes, err := json.Marshal(defaultChainCfg)
		gomega.Expect(err).Should(gomega.BeNil())

		// "ewoq" key with "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
		ewoqPrivKey := "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027"
		ewoqKey, err := ethcrypto.HexToECDSA(ewoqPrivKey)
		gomega.Expect(err).Should(gomega.BeNil())
		ewoqAddr := ethcrypto.PubkeyToAddress(ewoqKey.PublicKey)

		// "0x4fDdc14F51e0FE9651fcaf081F9ECFA725Ee9af2"
		privKey2 := "8ac3855ff600db43cf4e9c3a97df1d8ca35d478b8191ecb15c8ed29def82e063"
		key2, err := ethcrypto.HexToECDSA(privKey2)
		gomega.Expect(err).Should(gomega.BeNil())
		addr2 := ethcrypto.PubkeyToAddress(key2.PublicKey)

		blkChainID := ""
		ginkgo.By("upgrades all nodes with default chain config", func() {
			utils.Outf("{{magenta}}getting cluster status{{/}}\n")
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			sresp, err := runnerCli.Status(ctx)
			cancel()
			gomega.Expect(err).Should(gomega.BeNil())

			for blkChainID = range sresp.ClusterInfo.CustomChains {
				utils.Outf("{{magenta}}found block chain ID:{{/}} %q\n", blkChainID)
				break
			}

			for _, name := range sresp.ClusterInfo.NodeNames {
				configPath := filepath.Join(sresp.ClusterInfo.RootDataDir, name, "chainConfigs", blkChainID, "config.json")
				gomega.Expect(os.MkdirAll(filepath.Dir(configPath), 0777)).Should(gomega.BeNil())

				utils.Outf("{{magenta}}writing chain config for %q{{/}}: %q\n", name, configPath)
				gomega.Expect(os.WriteFile(configPath, defaultChainCfgBytes, 0777)).Should(gomega.BeNil())

				utils.Outf("{{magenta}}restarting the node %q{{/}}\n", name)
				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				resp, err := runnerCli.RestartNode(ctx, name)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())

				time.Sleep(20 * time.Second)

				ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
				_, err = runnerCli.Health(ctx)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())
				utils.Outf("{{green}}successfully upgraded %q{{/}} (current info: %+v)\n", name, resp.ClusterInfo.NodeInfos)
			}
		})

		evmCli := new(utils.EvmClient)
		ginkgo.By("makes txs after upgrade", func() {
			ci := utils.GetClusterInfo()
			gomega.Expect(len(ci.URIs) > 0).Should(gomega.BeTrue())
			utils.Outf("{{magenta}}sending txs to subnet-evm endpoints:{{/}} %q\n", ci.SubnetEVMRPCEndpoints)

			ep := ci.SubnetEVMRPCEndpoints[rand.Intn(100)%len(ci.SubnetEVMRPCEndpoints)]

			var err error
			evmCli, err = utils.NewEvmClient(ep, 25, 1)
			gomega.Expect(err).Should(gomega.BeNil())

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			prevBal, err := evmCli.FetchBalance(ctx, ewoqAddr)
			cancel()
			gomega.Expect(err).Should(gomega.BeNil())

			if prevBal.Cmp(new(big.Int).SetInt64(0)) == 0 {
				ginkgo.Skip("no balance... skipping tests...")
			}
			transferAmount := new(big.Int).Div(prevBal, new(big.Int).SetInt64(10))

			ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
			newBal, err := evmCli.TransferTx(ctx, ewoqAddr, ewoqKey, addr2, transferAmount)
			cancel()
			gomega.Expect(err).Should(gomega.BeNil())

			ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
			err = evmCli.WaitForBalance(ctx, ewoqAddr, newBal)
			cancel()
			gomega.Expect(err).Should(gomega.BeNil())

			utils.Outf("{{magenta}}successfully sent:{{/}} after balance %v (old balance %v)\n", newBal, prevBal)
		})

		// only allow txs from "ewoq"
		upgradeCfgAllowsEwoq := chainUpgradeConfig{
			PrecompileUpgrades: []*params.PrecompileUpgrade{
				{
					TxAllowListConfig: &precompile.TxAllowListConfig{
						AllowListConfig: precompile.AllowListConfig{
							AllowListAdmins: []common.Address{ewoqAddr},
						},
						UpgradeableConfig: precompile.UpgradeableConfig{
							BlockTimestamp: new(big.Int).SetInt64(time.Now().Unix()),
						},
					},
				},
			},
		}
		upgradeCfgAllowsEwoqBytes, err := json.Marshal(upgradeCfgAllowsEwoq)
		gomega.Expect(err).Should(gomega.BeNil())

		ginkgo.By("upgrades all nodes with precompile tx allow list", func() {
			utils.Outf("{{magenta}}getting cluster status{{/}}\n")
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			sresp, err := runnerCli.Status(ctx)
			cancel()
			gomega.Expect(err).Should(gomega.BeNil())

			for _, name := range sresp.ClusterInfo.NodeNames {
				upgradeJSONPath := filepath.Join(sresp.ClusterInfo.RootDataDir, name, "chainConfigs", blkChainID, "upgrade.json")
				gomega.Expect(os.MkdirAll(filepath.Dir(upgradeJSONPath), 0777)).Should(gomega.BeNil())

				utils.Outf("{{magenta}}writing chain config for %q{{/}}: %q\n", name, upgradeJSONPath)
				gomega.Expect(os.WriteFile(upgradeJSONPath, upgradeCfgAllowsEwoqBytes, 0777)).Should(gomega.BeNil())

				utils.Outf("{{magenta}}restarting the node %q{{/}}\n", name)
				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				resp, err := runnerCli.RestartNode(ctx, name)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())

				time.Sleep(20 * time.Second)

				ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
				_, err = runnerCli.Health(ctx)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())
				utils.Outf("{{green}}successfully upgraded %q{{/}} (current info: %+v)\n", name, resp.ClusterInfo.NodeInfos)
			}

			// advance block timestamps by issuing new blocks
			for i := 0; i < 5; i++ {
				time.Sleep(5 * time.Second)

				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				prevBal, err := evmCli.FetchBalance(ctx, ewoqAddr)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())

				if prevBal.Cmp(new(big.Int).SetInt64(0)) == 0 {
					ginkgo.Skip("no balance... skipping tests...")
				}
				transferAmount := new(big.Int).Div(prevBal, new(big.Int).SetInt64(10))

				ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
				_, err = evmCli.TransferTx(ctx, ewoqAddr, ewoqKey, addr2, transferAmount)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())
			}
		})

		ginkgo.By("tx allow list indeed restricts tx issuance from non-allow listed addresses", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			prevBal, err := evmCli.FetchBalance(ctx, addr2)
			cancel()
			gomega.Expect(err).Should(gomega.BeNil())

			if prevBal.Cmp(new(big.Int).SetInt64(0)) == 0 {
				ginkgo.Skip("no balance... skipping tests...")
			}
			transferAmount := new(big.Int).Div(prevBal, new(big.Int).SetInt64(10))

			// "addr2" is not allow-listed yet, so should fail!
			ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
			_, err = evmCli.TransferTx(ctx, addr2, key2, ewoqAddr, transferAmount)
			cancel()
			gomega.Expect(err.Error()).Should(gomega.ContainSubstring("non-allow listed address"))
		})

		// TODO: test other precompiles
		// e.g., contract deployer

		ginkgo.By("non-admin address should never be allowed to add it itself to the admin list", func() {
			ci := utils.GetClusterInfo()
			gomega.Expect(len(ci.URIs) > 0).Should(gomega.BeTrue())

			utils.Outf("{{green}}non-admin adding itself to admin list{{/}}\n")
			s, err := utils.RunCommand(
				2*time.Minute,
				"cast",
				"send",
				"--private-key="+privKey2,
				"--rpc-url="+ci.SubnetEVMRPCEndpoints[0],
				precompile.TxAllowListAddress.String(),
				"setAdmin(address)",
				addr2.String(),
			)
			gomega.Expect(err).Should(gomega.BeNil())
			gomega.Expect(s.Complete && s.Exit > 0).Should(gomega.BeTrue())

			// e.g., (code: -32000, message: non-admin cannot modify allow list: 0x4fDdc14F51e0FE9651fcaf081F9ECFA725Ee9af2, data: None)
			errMatched := false
			for _, em := range s.Stderr {
				errMatched = strings.Contains(em, "non-admin cannot modify allow list")
				if errMatched {
					break
				}
			}
			gomega.Expect(errMatched).Should(gomega.BeTrue())
		})

		ginkgo.By("non-admin address should never be allowed to add it itself to the allow list", func() {
			ci := utils.GetClusterInfo()
			gomega.Expect(len(ci.URIs) > 0).Should(gomega.BeTrue())

			utils.Outf("{{green}}non-admin adding itself to allow list{{/}}\n")
			s, err := utils.RunCommand(
				2*time.Minute,
				"cast",
				"send",
				"--private-key="+privKey2,
				"--rpc-url="+ci.SubnetEVMRPCEndpoints[0],
				precompile.TxAllowListAddress.String(),
				"setEnabled(address)",
				addr2.String(),
			)
			gomega.Expect(err).Should(gomega.BeNil())
			gomega.Expect(s.Complete && s.Exit > 0).Should(gomega.BeTrue())

			// e.g., (code: -32000, message: non-admin cannot modify allow list: 0x4fDdc14F51e0FE9651fcaf081F9ECFA725Ee9af2, data: None)
			errMatched := false
			for _, em := range s.Stderr {
				errMatched = strings.Contains(em, "non-admin cannot modify allow list")
				if errMatched {
					break
				}
			}
			gomega.Expect(errMatched).Should(gomega.BeTrue())
		})

		// call precompile contract to allow more addresses without restarts
		ginkgo.By("tx allow list admin can add more addresses to the allow lists", func() {
			ci := utils.GetClusterInfo()
			gomega.Expect(len(ci.URIs) > 0).Should(gomega.BeTrue())

			utils.Outf("{{green}}adding another address to allow list{{/}}\n")
			s, err := utils.RunCommand(
				2*time.Minute,
				"cast",
				"send",
				"--private-key="+ewoqPrivKey, // ewoq key
				"--rpc-url="+ci.SubnetEVMRPCEndpoints[0],
				precompile.TxAllowListAddress.String(),
				"setEnabled(address)",
				addr2.String(),
			)
			gomega.Expect(err).Should(gomega.BeNil())
			gomega.Expect(s.Complete && s.Exit == 0 && s.Error == nil).Should(gomega.BeTrue())

			utils.Outf("{{green}}reading the current allow list{{/}}\n")
			s, err = utils.RunCommand(
				2*time.Minute,
				"cast",
				"call",
				"--rpc-url="+ci.SubnetEVMRPCEndpoints[0],
				precompile.TxAllowListAddress.String(),
				"readAllowList(address)",
				addr2.String(),
			)
			gomega.Expect(err).Should(gomega.BeNil())
			gomega.Expect(s.Complete && s.Exit == 0 && s.Error == nil).Should(gomega.BeTrue())
			gomega.Expect(s.Stdout[0]).Should(gomega.Equal("0x0000000000000000000000000000000000000000000000000000000000000001"))
		})

		ginkgo.By("newly added address can now issue txs", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			prevBal, err := evmCli.FetchBalance(ctx, addr2)
			cancel()
			gomega.Expect(err).Should(gomega.BeNil())

			if prevBal.Cmp(new(big.Int).SetInt64(0)) == 0 {
				ginkgo.Skip("no balance... skipping tests...")
			}
			transferAmount := new(big.Int).Div(prevBal, new(big.Int).SetInt64(10))

			// "addr2" is now allow-listed, so should not fail!
			ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
			_, err = evmCli.TransferTx(ctx, addr2, key2, ewoqAddr, transferAmount)
			cancel()
			gomega.Expect(err).Should(gomega.BeNil())
		})
	})
})
