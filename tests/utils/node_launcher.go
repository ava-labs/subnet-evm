// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
)

const (
	isPartitioningEnabled    = false
	enclaveIdPrefix          = "test"
	avalancheStarlarkPackage = "github.com/kurtosis-tech/avalanche-package"
	defaultParallelism       = 4
	nodePrefix               = "node-"
	testImageId              = "avaplatform/avalanchego:test"
)

func SpinupAvalancheNodes(nodeCount int) ([]string, func(), error) {
	ctx := context.Background()

	packageArgumentsToStartNNodeTestNet := `{
		"test_mode": true,
		"node_count": ` + strconv.Itoa(nodeCount) + `,
		"avalanchego_image": "` + testImageId + `"
	}`

	kurtosisCtx, err := kurtosis_context.NewKurtosisContextFromLocalEngine()
	if err != nil {
		return nil, nil, err
	}

	enclaveId := fmt.Sprintf("%s-%d", enclaveIdPrefix, time.Now().Unix())
	enclaveCtx, err := kurtosisCtx.CreateEnclave(ctx, enclaveId, isPartitioningEnabled)
	if err != nil {
		return nil, nil, err
	}

	_, err = enclaveCtx.RunStarlarkRemotePackageBlocking(ctx, avalancheStarlarkPackage, packageArgumentsToStartNNodeTestNet, false, defaultParallelism)
	if err != nil {
		return nil, nil, fmt.Errorf("an error occurred while running Starlark Package: %v", err)
	}

	var nodeRpcUris []string

	for nodeIdx := 0; nodeIdx < nodeCount; nodeIdx++ {
		nodeId := fmt.Sprintf("%v%d", nodePrefix, nodeIdx)
		serviceCtx, err := enclaveCtx.GetServiceContext(nodeId)
		if err != nil {
			return nil, nil, err
		}

		publicRpcPorts := serviceCtx.GetPublicPorts()
		rpcPortSpec, found := publicRpcPorts["rpc"]
		if !found {
			return nil, nil, fmt.Errorf("couldn't find RPC port in the node '%v' that was spun up", nodeId)
		}

		rpcPortNumber := rpcPortSpec.GetNumber()
		nodeRpcUris = append(nodeRpcUris, fmt.Sprintf("http://127.0.0.1:%d", rpcPortNumber))
	}

	tearDownFunction := func() {
		fmt.Println(fmt.Printf("Destroying enclave with id '%v'", enclaveId))
		if err = kurtosisCtx.StopEnclave(ctx, enclaveId); err != nil {
			fmt.Printf("An error occurred while stopping the enclave with id '%v'\n", enclaveId)
		}
		if err = kurtosisCtx.DestroyEnclave(ctx, enclaveId); err != nil {
			fmt.Printf("An error occurred while cleaning up the enclave with id '%v'\n", enclaveId)
		}
	}

	return nodeRpcUris, tearDownFunction, nil
}
