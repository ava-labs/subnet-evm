#!/usr/bin/env bash

set -euo pipefail

# Directory above this script
SUBNET_EVM_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )
# Load the constants
source "$SUBNET_EVM_PATH"/scripts/constants.sh

echo "Building Workload..."
pushd "$SUBNET_EVM_PATH/avalanche" > /dev/null
  go build -o "$SUBNET_EVM_PATH/build/workload" "$SUBNET_EVM_PATH/avalanche/tests/antithesis/"*.go
popd > /dev/null
