#!/usr/bin/env bash

set -euo pipefail

# Builds docker images for antithesis testing.

# e.g.,
# ./scripts/build_antithesis_images.sh                                           # Build local images
# IMAGE_PREFIX=<registry>/<repo> TAG=latest ./scripts/build_antithesis_images.sh # Specify a prefix to enable image push and use a specific tag

# Directory above this script
SUBNET_EVM_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )

# Allow configuring the clone path to point to a shared and/or existing clone of the avalanchego repo
AVALANCHEGO_CLONE_PATH="${AVALANCHEGO_CLONE_PATH:-${SUBNET_EVM_PATH}/avalanchego}"

# Specifying an image prefix will ensure the image is pushed after build
IMAGE_PREFIX="${IMAGE_PREFIX:-}"

TAG="${TAG:-}"
if [[ -z "${TAG}" ]]; then
  # Default to tagging with the commit hash
  source "${SUBNET_EVM_PATH}"/scripts/constants.sh
  TAG="${SUBNET_EVM_COMMIT::8}"
fi

# The dockerfiles don't specify the golang version to minimize the changes required to bump
# the version. Instead, the golang version is provided as an argument.
GO_VERSION="$(go list -m -f '{{.GoVersion}}')"

function build_images {
  local base_image_name=$1
  local uninstrumented_node_dockerfile=$2
  local avalanche_node_image=$3

  # Define image names
  if [[ -n "${IMAGE_PREFIX}" ]]; then
    base_image_name="${IMAGE_PREFIX}/${base_image_name}"
  fi
  local node_image_name="${base_image_name}-node:${TAG}"
  local workload_image_name="${base_image_name}-workload:${TAG}"
  local config_image_name="${base_image_name}-config:${TAG}"

  # Define dockerfiles
  local base_dockerfile="${SUBNET_EVM_PATH}/tests/antithesis/Dockerfile"
  local node_dockerfile="${base_dockerfile}.node"
  if [[ "$(go env GOARCH)" == "arm64" ]]; then
    # Antithesis instrumentation is only supported on amd64. On apple silicon (arm64), the
    # uninstrumented Dockerfile will be used to build the node image to enable local test
    # development.
    node_dockerfile="${uninstrumented_node_dockerfile}"
  fi

  # Define default build command
  local docker_cmd="docker buildx build --build-arg GO_VERSION=${GO_VERSION} --build-arg NODE_IMAGE=${node_image_name}"

  # Build node image first to allow the workload image to be based on it.
  ${docker_cmd} --build-arg AVALANCHEGO_NODE_IMAGE="${avalanche_node_image}" -t "${node_image_name}" \
                -f "${node_dockerfile}" "${SUBNET_EVM_PATH}"
  TARGET_PATH="${SUBNET_EVM_PATH}/build/antithesis"
  if [[ -d "${TARGET_PATH}" ]]; then
    # Ensure the target path is empty before generating the compose config
    rm -r "${TARGET_PATH}"
  fi

  # Ensure avalanchego and subnet-evm binaries are available to create an initial db state that includes subnets.
  "${AVALANCHEGO_CLONE_PATH}"/scripts/build.sh
  PLUGIN_PATH="${TARGET_PATH}"/plugins
  "${SUBNET_EVM_PATH}"/scripts/build.sh "${PLUGIN_PATH}"/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy

  # Generate compose config and db state for the config image
  TARGET_PATH="${TARGET_PATH}"\
    IMAGE_TAG="${TAG}"\
    AVALANCHEGO_PATH="${AVALANCHEGO_CLONE_PATH}/build/avalanchego"\
    AVALANCHEGO_PLUGIN_DIR="${PLUGIN_PATH}"\
    go run "${SUBNET_EVM_PATH}/tests/antithesis/gencomposeconfig"

  # Build config image
  ${docker_cmd} -t "${config_image_name}" -f "${base_dockerfile}.config" "${SUBNET_EVM_PATH}"

  # Build workload image
  ${docker_cmd} -t "${workload_image_name}" -f "${base_dockerfile}.workload" "${SUBNET_EVM_PATH}"
}

# Assume it's necessary to build the avalanchego node image from source
# TODO(marun) Support use of a released node image if using a release version of avalanchego

source "${SUBNET_EVM_PATH}"/scripts/versions.sh

echo "checking out target avalanchego version ${AVALANCHE_VERSION}"
if [[ -d "${AVALANCHEGO_CLONE_PATH}" ]]; then
  echo "updating existing clone"
  cd "${AVALANCHEGO_CLONE_PATH}"
  git fetch
else
  echo "creating new clone"
  git clone https://github.com/ava-labs/avalanchego.git "${AVALANCHEGO_CLONE_PATH}"
  cd "${AVALANCHEGO_CLONE_PATH}"
fi
# Branch will be reset to $AVALANCHE_VERSION if it already exists
git checkout -B "test-${AVALANCHE_VERSION}" "${AVALANCHE_VERSION}"
cd "${SUBNET_EVM_PATH}"

# Build avalanchego node image. Supply an empty tag so the tag can be discovered from the hash of the avalanchego repo.
NODE_ONLY=1 TEST_SETUP=avalanchego IMAGE_PREFIX="${IMAGE_PREFIX}" TAG='' bash -x "${AVALANCHEGO_CLONE_PATH}"/scripts/build_antithesis_images.sh

build_images antithesis-subnet-evm "${SUBNET_EVM_PATH}/Dockerfile" "antithesis-avalanchego-node:${AVALANCHE_VERSION::8}"
