#!/usr/bin/env bash
set -e

# This script assumes that an AvalancheGo and Subnet-EVM binaries are available in the standard location
# within the $GOPATH
# The AvalancheGo and PluginDir paths can be specified via the environment variables used in ./scripts/run.sh.

# Load the versions
SUBNET_EVM_PATH=$(
  cd "$(dirname "${BASH_SOURCE[0]}")"
  cd .. && pwd
)

source "$SUBNET_EVM_PATH"/scripts/constants.sh

source "$SUBNET_EVM_PATH"/scripts/versions.sh

echo "compiling Warp contracts"
solc --evm-version paris "./tests/warp/contracts/ExampleWarp.sol"  --abi --bin -o "./tests/warp/contracts/" --overwrite

# Build ginkgo
echo "building precompile.test"
# to install the ginkgo binary (required for test build and run)
go install -v github.com/onsi/ginkgo/v2/ginkgo@${GINKGO_VERSION}

ACK_GINKGO_RC=true ginkgo build ./tests/precompile ./tests/load ./tests/warp

# By default, it runs all e2e test cases!
# Use "--ginkgo.skip" to skip tests.
# Use "--ginkgo.focus" to select tests.
./tests/precompile/precompile.test \
  --ginkgo.vv \
  --ginkgo.label-filter=${GINKGO_LABEL_FILTER:-""}\
  --ginkgo.focus="none"

./tests/load/load.test \
  --ginkgo.vv \
  --ginkgo.label-filter=${GINKGO_LABEL_FILTER:-""}\
  --ginkgo.focus="none"

./tests/warp/warp.test \
  --ginkgo.vv \
  --ginkgo.label-filter=${GINKGO_LABEL_FILTER:-""}
