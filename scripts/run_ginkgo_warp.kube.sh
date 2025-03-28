#!/usr/bin/env bash

set -euo pipefail

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

EXTRA_ARGS=()
AVALANCHEGO_BUILD_PATH="${AVALANCHEGO_BUILD_PATH:-}"
if [[ -n "${AVALANCHEGO_BUILD_PATH}" ]]; then
  EXTRA_ARGS=("--avalanchego-path=${AVALANCHEGO_BUILD_PATH}/avalanchego")
  echo "Running with extra args:" "${EXTRA_ARGS[@]}"
fi

export KUBECONFIG="${KUBECONFIG:-$HOME/.kube/config}"

# Enable collector start if credentials are set in the env
if [[ -n "${PROMETHEUS_USERNAME:-}" ]]; then
  export TMPNET_START_COLLECTORS=true
fi

"${SUBNET_EVM_PATH}"/bin/tmpnetctl start-kind-cluster

DOCKER_IMAGE="${DOCKER_IMAGE:-localhost:5001/subnet-evm}"
AVALANCHEGO_LOCAL_IMAGE_NAME="${AVALANCHEGO_LOCAL_IMAGE_NAME:-localhost:5001/avalanchego}"
if [[ -z "${SKIP_BUILD_IMAGE:-}" ]]; then
  FORCE_TAG_LATEST=1 IMAGE_NAME="${DOCKER_IMAGE}" AVALANCHEGO_LOCAL_IMAGE_NAME="${AVALANCHEGO_LOCAL_IMAGE_NAME}" bash -x "${SUBNET_EVM_PATH}"/scripts/build_docker_image.sh
fi

GINKGO_ARGS=()
# Reference: https://onsi.github.io/ginkgo/#spec-randomization
if [[ -n "${E2E_RANDOM_SEED:-}" ]]; then
  # Supply a specific seed to simplify reproduction of test failures
  GINKGO_ARGS+=(--seed="${E2E_RANDOM_SEED}")
else
  # Execute in random order to identify unwanted dependency
  GINKGO_ARGS+=(--randomize-all)
fi

"${SUBNET_EVM_PATH}"/bin/ginkgo -vv "${GINKGO_ARGS[@]}" --label-filter="${GINKGO_LABEL_FILTER:-}" ./tests/warp --\
  "${EXTRA_ARGS[@]}" --runtime=kube --image-name="${DOCKER_IMAGE}"
