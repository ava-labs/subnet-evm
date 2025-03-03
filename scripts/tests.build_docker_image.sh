#!/usr/bin/env bash

set -euo pipefail

# Sanity check the image build by attempting to build and run the image without error.

# Directory above this script
SUBNET_EVM_PATH=$(
  cd "$(dirname "${BASH_SOURCE[0]}")"
  cd .. && pwd
)
# Load the constants
source "$SUBNET_EVM_PATH"/scripts/constants.sh

# Load the versions
source "$SUBNET_EVM_PATH"/scripts/versions.sh

# Use the default node image
AVALANCHEGO_NODE_IMAGE="${AVALANCHEGO_IMAGE_NAME}:${AVALANCHE_VERSION}"

build_and_test() {
  local reponame="${1}"
  local vm_id="${2}"
  local multiarch_image="${3}"

  if [[ "${multiarch_image}" == true ]]; then
    # Assume a registry image is a multi-arch image
    local arches="linux/amd64,linux/arm64"
  else
    # Test only the host platform for non-registry/single arch builds
    local host_arch
    host_arch="$(go env GOARCH)"
    local arches="linux/$host_arch"
  fi

  # Build the avalanchego image if it cannot be pulled. This will usually be due to
  # AVALANCHE_VERSION being not yet merged since the image is published post-merge.
  if ! docker pull "${AVALANCHEGO_NODE_IMAGE}"; then
    # Use a image name without a repository (i.e. without 'avaplatform/' prefix ) to build a
    # local image that will not be pushed.
    export AVALANCHEGO_IMAGE_NAME="avalanchego"
    echo "Building ${AVALANCHEGO_IMAGE_NAME}:${AVALANCHE_VERSION} locally"

    source "${SUBNET_EVM_PATH}"/scripts/lib_avalanchego_clone.sh
    clone_avalanchego "${AVALANCHE_VERSION}"
    SKIP_BUILD_RACE=1 DOCKER_IMAGE="${AVALANCHEGO_IMAGE_NAME}" "${AVALANCHEGO_CLONE_PATH}"/scripts/build_image.sh
  fi

  local imgtag="testtag"

  PLATFORMS="${arches}" \
    BUILD_IMAGE_ID="${imgtag}" \
    VM_ID=$"${vm_id}" \
    DOCKERHUB_REPO="${reponame}" \
    ./scripts/build_docker_image.sh

  echo "listing images"
  docker images

  # Check all of the images expected to have been built
  local target_images=(
    "$reponame:$imgtag"
    "$reponame:$DOCKERHUB_TAG"
  )
  IFS=',' read -r -a archarray <<<"$arches"
  for arch in "${archarray[@]}"; do
    for target_image in "${target_images[@]}"; do
      echo "checking sanity of image $target_image for $arch by running '${VM_ID} version'"
      docker run -t --rm --platform "$arch" "$target_image" /avalanchego/build/plugins/"${VM_ID}" --version
    done
  done
}

VM_ID="${VM_ID:-${DEFAULT_VM_ID}}"

echo "checking build of single-arch image"
build_and_test "subnet-evm" "${VM_ID}" false

echo "starting local docker registry to allow verification of multi-arch image builds"
REGISTRY_CONTAINER_ID="$(docker run --rm -d -P registry:2)"
REGISTRY_PORT="$(docker port "$REGISTRY_CONTAINER_ID" 5000/tcp | grep -v "::" | awk -F: '{print $NF}')"

echo "starting docker builder that supports multiplatform builds"
# - creating a new builder enables multiplatform builds
# - '--driver-opt network=host' enables the builder to use the local registry
docker buildx create --use --name ci-builder --driver-opt network=host

# Ensure registry and builder cleanup on teardown
function cleanup {
  echo "stopping local docker registry"
  docker stop "${REGISTRY_CONTAINER_ID}"
  echo "removing multiplatform builder"
  docker buildx rm ci-builder
}
trap cleanup EXIT

echo "checking build of multi-arch images"
build_and_test "localhost:${REGISTRY_PORT}/subnet-evm" "${VM_ID}" true
