#!/usr/bin/env bash
set -e

# This script assumes that an AvalancheGo and Subnet-EVM binaries are available in the standard location
# within the $GOPATH
# The AvalancheGo and PluginDir paths can be specified via the environment variables used in ./scripts/run.sh.

# e.g.,
# ./scripts/run_ginkgo_warp.sh
# ./scripts/run_ginkgo_warp.sh --ginkgo.label-filter=x  # All arguments are supplied to ginkgo
if ! [[ "$0" =~ scripts/run_ginkgo_warp.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

# Ensure avalanchego release is available
./scripts/install_avalanchego_release.sh

# Build subnet-evm
./scripts/build.sh

# Ensure the ginkgo version is available
source ./scripts/versions.sh

# Install the ginkgo binary
go install -v github.com/onsi/ginkgo/v2/ginkgo@${GINKGO_VERSION}

ginkgo --vv --randomize-all ./tests/warp -- "${@}"
