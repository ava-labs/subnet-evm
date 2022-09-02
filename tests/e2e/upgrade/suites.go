// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements upgrade tests, requires network-runner cluster.
package upgrade

import (
	"context"
	"encoding/json"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/ava-labs/avalanchego/tests"
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

var _ = utils.DescribeLocal("[Upgrade]", func() {
	ginkgo.It("can upgrade with precompile", ginkgo.Label("upgrade"), func() {
		runnerCli := utils.GetClient()
		gomega.Expect(runnerCli).ShouldNot(gomega.BeNil())

		chainCfg := chainConfig{}
		chainCfg.Config.SetDefaults()
		chainCfg.Config.LogJSONFormat = true

		chainUpgradeCfg := chainUpgradeConfig{}

		// "ewoq" key with "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
		key1, err := ethcrypto.HexToECDSA("56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027")
		gomega.Expect(err).Should(gomega.BeNil())
		addr1 := ethcrypto.PubkeyToAddress(key1.PublicKey)

		// "0x4fDdc14F51e0FE9651fcaf081F9ECFA725Ee9af2"
		key2, err := ethcrypto.HexToECDSA("8ac3855ff600db43cf4e9c3a97df1d8ca35d478b8191ecb15c8ed29def82e063")
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

			chainCfgBytes, err := json.Marshal(chainCfg)
			gomega.Expect(err).Should(gomega.BeNil())

			for _, name := range sresp.ClusterInfo.NodeNames {
				configPath := filepath.Join(sresp.ClusterInfo.RootDataDir, name, "chainConfigs", blkChainID, "config.json")
				gomega.Expect(os.MkdirAll(filepath.Dir(configPath), 0777)).Should(gomega.BeNil())

				utils.Outf("{{magenta}}writing chain config for %q{{/}}: %q\n", name, configPath)
				gomega.Expect(os.WriteFile(configPath, chainCfgBytes, 0777)).Should(gomega.BeNil())

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
				tests.Outf("{{green}}successfully upgraded %q{{/}} (current info: %+v)\n", name, resp.ClusterInfo.NodeInfos)
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
			prevBal, err := evmCli.FetchBalance(ctx, addr1)
			cancel()
			gomega.Expect(err).Should(gomega.BeNil())

			if prevBal.Cmp(new(big.Int).SetInt64(0)) == 0 {
				ginkgo.Skip("no balance... skipping tests...")
			}
			transferAmount := new(big.Int).Div(prevBal, new(big.Int).SetInt64(10))

			ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
			newBal, err := evmCli.TransferTx(ctx, addr1, key1, addr2, transferAmount)
			cancel()
			gomega.Expect(err).Should(gomega.BeNil())

			ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
			err = evmCli.WaitForBalance(ctx, addr1, newBal)
			cancel()
			gomega.Expect(err).Should(gomega.BeNil())

			utils.Outf("{{magenta}}successfully sent:{{/}} after balance %v (old balance %v)\n", newBal, prevBal)
		})

		ginkgo.By("upgrades all nodes with precompile tx allow list", func() {
			utils.Outf("{{magenta}}getting cluster status{{/}}\n")
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			sresp, err := runnerCli.Status(ctx)
			cancel()
			gomega.Expect(err).Should(gomega.BeNil())

			// only allow txs from "ewoq"
			chainUpgradeCfg.PrecompileUpgrades = []*params.PrecompileUpgrade{
				{
					TxAllowListConfig: &precompile.TxAllowListConfig{
						AllowListConfig: precompile.AllowListConfig{
							AllowListAdmins: []common.Address{addr1},
						},
						UpgradeableConfig: precompile.UpgradeableConfig{
							BlockTimestamp: new(big.Int).SetInt64(time.Now().Unix()),
						},
					},
				},
			}

			chainUpgradeCfgBytes, err := json.Marshal(chainUpgradeCfg)
			gomega.Expect(err).Should(gomega.BeNil())

			for _, name := range sresp.ClusterInfo.NodeNames {
				configPath := filepath.Join(sresp.ClusterInfo.RootDataDir, name, "chainConfigs", blkChainID, "upgrade.json")
				gomega.Expect(os.MkdirAll(filepath.Dir(configPath), 0777)).Should(gomega.BeNil())

				utils.Outf("{{magenta}}writing chain config for %q{{/}}: %q\n", name, configPath)
				gomega.Expect(os.WriteFile(configPath, chainUpgradeCfgBytes, 0777)).Should(gomega.BeNil())

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
				tests.Outf("{{green}}successfully upgraded %q{{/}} (current info: %+v)\n", name, resp.ClusterInfo.NodeInfos)
			}

			// advance block timestamps by issuing new blocks
			for i := 0; i < 5; i++ {
				time.Sleep(5 * time.Second)

				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				prevBal, err := evmCli.FetchBalance(ctx, addr1)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())

				if prevBal.Cmp(new(big.Int).SetInt64(0)) == 0 {
					ginkgo.Skip("no balance... skipping tests...")
				}
				transferAmount := new(big.Int).Div(prevBal, new(big.Int).SetInt64(10))

				ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
				_, err = evmCli.TransferTx(ctx, addr1, key1, addr2, transferAmount)
				cancel()
				gomega.Expect(err).Should(gomega.BeNil())
			}
		})

		ginkgo.By("checks whether tx allow list indeed limits tx issuers", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			prevBal, err := evmCli.FetchBalance(ctx, addr2)
			cancel()
			gomega.Expect(err).Should(gomega.BeNil())

			if prevBal.Cmp(new(big.Int).SetInt64(0)) == 0 {
				ginkgo.Skip("no balance... skipping tests...")
			}
			transferAmount := new(big.Int).Div(prevBal, new(big.Int).SetInt64(10))

			ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
			_, err = evmCli.TransferTx(ctx, addr2, key2, addr1, transferAmount)
			cancel()
			gomega.Expect(err.Error()).Should(gomega.ContainSubstring("non-allow listed address"))
		})
	})
})
