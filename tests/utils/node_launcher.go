// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
)

const (
	isPartitioningEnabled    = false
	enclaveIdPrefix          = "test"
	avalancheStarlarkPackage = "github.com/kurtosis-tech/avalanche-package"
	// forces the node to launch on 9650 instead of ephemeral ports
	forceExposeOn9650        = `{"test_mode": true}`
	defaultParallelism       = 4
	firstNodeId              = "node-0"
	validationErrorDelimiter = ", "
)

func SpinupAvalancheNode() (string, func(), error) {
	ctx := context.Background()

	kurtosisCtx, err := kurtosis_context.NewKurtosisContextFromLocalEngine()
	if err != nil {
		return "", nil, err
	}

	enclaveId := fmt.Sprintf("%s-%d", enclaveIdPrefix, time.Now().Unix())
	enclaveCtx, err := kurtosisCtx.CreateEnclave(ctx, enclaveId, isPartitioningEnabled)
	if err != nil {
		return "", nil, err
	}

	_, err = enclaveCtx.RunStarlarkRemotePackageBlocking(ctx, avalancheStarlarkPackage, forceExposeOn9650, false, defaultParallelism)
	if err != nil {
		return "", nil, fmt.Errorf("an error occurred while running Starlark Package: %v", err)
	}

	serviceCtx, err := enclaveCtx.GetServiceContext(firstNodeId)
	if err != nil {
		return "", nil, err
	}

	publicRpcPorts := serviceCtx.GetPublicPorts()
	rpcPortSpec, found := publicRpcPorts["rpc"]
	if !found {
		return "", nil, fmt.Errorf("couldn't find RPC port in the node '%v' that was spun up", firstNodeId)
	}

	rpcPortNumber := rpcPortSpec.GetNumber()

	tearDownFunction := func() {
		fmt.Println(fmt.Printf("Destroying enclave with id '%v'", enclaveId))
		if err = kurtosisCtx.StopEnclave(ctx, enclaveId); err != nil {
			fmt.Printf("An error occurred while stopping the enclave with id '%v'\n", enclaveId)
		}
		if err = kurtosisCtx.DestroyEnclave(ctx, enclaveId); err != nil {
			fmt.Printf("An error occurred while cleaning up the enclave with id '%v'\n", enclaveId)
		}
	}

	return fmt.Sprintf("http://127.0.0.1:%d", rpcPortNumber), tearDownFunction, nil
}
