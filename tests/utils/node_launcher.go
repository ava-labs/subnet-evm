package utils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
)

const (
	isPartitioningEnabled    = false
	enclaveIdPrefix          = "test"
	avalancheStarlarkPackage = "github.com/kurtosis-tech/avalanche-package"
	// forces the node to launch on 9650 instead of ephemeral ports
	forceExposeOn9650  = `{"test_mode": true}`
	defaultParallelism = 4
	firstNodeId        = "node-0"
)

func SpinupAvalancheNode() (error, string, func()) {
	ctx := context.Background()

	kurtosisCtx, err := kurtosis_context.NewKurtosisContextFromLocalEngine()
	if err != nil {
		return err, "", nil
	}

	enclaveId := fmt.Sprintf("%s-%d", enclaveIdPrefix, time.Now().Unix())
	enclaveCtx, err := kurtosisCtx.CreateEnclave(ctx, enclaveId, isPartitioningEnabled)
	if err != nil {
		return err, "", nil
	}

	runResult, err := enclaveCtx.RunStarlarkRemotePackageBlocking(ctx, avalancheStarlarkPackage, forceExposeOn9650, false, defaultParallelism)
	if err != nil {
		return err, "", nil
	}

	if runResult.InterpretationError != nil {
		return errors.New("error interpreting Starlark code"), "", nil
	}
	if len(runResult.ValidationErrors) != 0 {
		return errors.New("error validating Starlark code"), "", nil
	}
	if runResult.ExecutionError != nil {
		return errors.New("error executing Starlark code"), "", nil
	}

	serviceCtx, err := enclaveCtx.GetServiceContext(firstNodeId)
	if err != nil {
		return nil, "", nil
	}

	publicRpcPorts := serviceCtx.GetPublicPorts()
	rpcPortSpec, found := publicRpcPorts["rpc"]
	if !found {
		return fmt.Errorf("couldn't find RPC port in the node '%v' that was spun up", firstNodeId), "", nil
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

	return nil, fmt.Sprintf("http://127.0.0.1:%d", rpcPortNumber), tearDownFunction
}
