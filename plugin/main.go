// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"fmt"

	"github.com/ava-labs/avalanchego/version"
	"github.com/ethereum/go-ethereum/plugin/evm"
	"github.com/ethereum/go-ethereum/plugin/runner"
)

func main() {
	versionString := fmt.Sprintf("Subnet-EVM/%s [AvalancheGo=%s, rpcchainvm=%d]", evm.Version, version.Current, version.RPCChainVMProtocol)
	runner.Run(versionString)
}
