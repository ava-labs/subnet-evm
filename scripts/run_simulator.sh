#!/usr/bin/env bash
# This script runs a load simulation when the following options are provided:
# --test-type (load or warp)
# --endpoints (comma separated list of RPC endpoints) or --endpoints-file (path to file containing RPC endpoints)

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

run_simulator() {
    echo "building simulator"
    pushd ./cmd/simulator
    go build -o ./simulator main/*.go
    echo 

    popd
    echo "running simulator from $PWD"
    ./cmd/simulator/simulator \
        "$@" \
        --key-dir=./cmd/simulator/.simulator/keys \
        --timeout=30s \
        --workers=1 \
        --max-fee-cap=300 \
        --max-tip-cap=100
}

run_simulator "$@"
