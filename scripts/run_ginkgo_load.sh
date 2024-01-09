#!/usr/bin/env bash

set -euo pipefail

# e.g.,
# ./scripts/run_ginkgo_load.sh
# ./scripts/run_ginkgo_load.sh --ginkgo.label-filter=x  # All arguments are supplied to ginkgo
if ! [[ "$0" =~ scripts/run_ginkgo_load.sh ]]; then
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

# Run tests in random order to avoid dependency
ginkgo --vv --randomize-all ./tests/load -- "${@}"
