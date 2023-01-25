#!/usr/bin/env bash
set -e

echo "Beginning simualtor script"

run_simulator() {
    #################################
    echo "building simulator"
    pushd ./cmd/simulator
    go install -v .
    echo 

    popd
    echo "running simulator from $PWD"
    simulator \
        --rpc-endpoints=$RPC_ENDPOINTS \
        --keys=./cmd/simulator/.simulator/keys \
        --timeout=30s \
        --concurrency=10 \
        --base-fee=300 \
        --priority-fee=100
}

run_simulator
