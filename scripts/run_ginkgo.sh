#!/usr/bin/env bash
set -e

# Load the versions
SUBNET_EVM_PATH=$(
  cd "$(dirname "${BASH_SOURCE[0]}")"
  cd .. && pwd
)

source "$SUBNET_EVM_PATH"/scripts/versions.sh

run_ginkgo() {
  echo "building e2e.test"
  # to install the ginkgo binary (required for test build and run)
  go install -v github.com/onsi/ginkgo/v2/ginkgo@${ginkgo_version}

  ACK_GINKGO_RC=true ginkgo build ./tests/e2e

  # By default, it runs all e2e test cases!
  # Use "--ginkgo.skip" to skip tests.
  # Use "--ginkgo.focus" to select tests.
  ./tests/e2e/e2e.test \
    --ginkgo.vv \
    --ginkgo.label-filter=${GINKGO_LABEL_FILTER:-""}
}

run_ginkgo