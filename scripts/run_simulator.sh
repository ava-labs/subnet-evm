#!/usr/bin/env bash
# This script runs a load simulation against endpoints specified either as:
# - Comma separated list in the RPC_ENDPOINTS environment variable.
# - In the file specified by the RPC_ENDPOINTS_FILE environment variable.

set -e

echo "Beginning simulator script"

if ! [[ "$0" =~ scripts/run_simulator.sh ]]; then
  echo "must be run from repository root, but got $0"
  exit 255
fi

# Load the versions
SUBNET_EVM_PATH=$(
  cd "$(dirname "${BASH_SOURCE[0]}")"
  cd .. && pwd
)
source "$SUBNET_EVM_PATH"/scripts/versions.sh

# Load the constants
source "$SUBNET_EVM_PATH"/scripts/constants.sh

# Pass either --endpoints-file or --endpoints
if [[ -n "$RPC_ENDPOINTS_FILE" ]]; then
  ENDPOINT_OPTS="--endpoints-file=$RPC_ENDPOINTS_FILE"
else 
  ENDPOINT_OPTS="--endpoints=$RPC_ENDPOINTS"
fi

run_simulator() {
    #################################
    echo "building simulator"
    pushd ./cmd/simulator
    go build -o ./simulator main/*.go
    echo 

    popd
    echo "running simulator from $PWD"
    ./cmd/simulator/simulator \
        ${ENDPOINT_OPTS} \
        --key-dir=./cmd/simulator/.simulator/keys \
        --timeout=30s \
        --workers=1 \
        --max-fee-cap=300 \
        --max-tip-cap=100
}

run_simulator
