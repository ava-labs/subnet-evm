#!/usr/bin/env bash
set -e

# This script assumes that Kurtosis is installed and an engine is running
# Head over to https://docs.kurtosis.com/install/#ii-install-the-cli to see how to install Kurtosis
# In case an engine isn't running use kurtosis engine restart
# This assumes that the node image avaplatform/avalanchego with tag :test exists;
# you can create it with BUILD_IMAGE_ID=test ./scripts/build_image.sh

# Load the versions
SUBNET_EVM_PATH=$(
  cd "$(dirname "${BASH_SOURCE[0]}")"
  cd .. && pwd
)

source "$SUBNET_EVM_PATH"/scripts/constants.sh

source "$SUBNET_EVM_PATH"/scripts/versions.sh

# Build ginkgo
echo "building precompile.test"
# to install the ginkgo binary (required for test build and run)
go install -v github.com/onsi/ginkgo/v2/ginkgo@${GINKGO_VERSION}

ACK_GINKGO_RC=true ginkgo build ./tests/precompile ./tests/load

# By default, it runs all e2e test cases!
# Use "--ginkgo.skip" to skip tests.
# Use "--ginkgo.focus" to select tests.
./tests/precompile/precompile.test \
  --ginkgo.vv \
  --ginkgo.label-filter=${GINKGO_LABEL_FILTER:-""}

./tests/load/load.test \
  --ginkgo.vv \
  --ginkgo.label-filter=${GINKGO_LABEL_FILTER:-""}
