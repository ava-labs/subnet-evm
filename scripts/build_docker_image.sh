#!/usr/bin/env bash

set -euo pipefail

# If set to non-empty, prompts the building of a multi-arch image when the image
# name indicates use of a registry.
#
# A registry is required to build a multi-arch image since a multi-arch image is
# not really an image at all. A multi-arch image (also called a manifest) is
# basically a list of arch-specific images available from the same registry that
# hosts the manifest. Manifests are not supported for local images.
#
# Reference: https://docs.docker.com/build/building/multi-platform/
PLATFORMS="${PLATFORMS:-}"

# If set to non-empty, the image will be published to the registry.
PUBLISH="${PUBLISH:-}"

# Directory above this script
SUBNET_EVM_PATH=$(
  cd "$(dirname "${BASH_SOURCE[0]}")"
  cd .. && pwd
)

# Load the constants
source "$SUBNET_EVM_PATH"/scripts/constants.sh

# Load the versions
source "$SUBNET_EVM_PATH"/scripts/versions.sh

# WARNING: this will use the most recent commit even if there are un-committed changes present
BUILD_IMAGE_ID=${BUILD_IMAGE_ID:-"${CURRENT_BRANCH}"}

# buildx (BuildKit) improves the speed and UI of builds over the legacy builder and
# simplifies creation of multi-arch images.
#
# Reference: https://docs.docker.com/build/buildkit/
DOCKER_CMD="docker buildx build"

if [[ -n "${PUBLISH}" ]]; then
  echo "Pushing $DOCKERHUB_REPO:$BUILD_IMAGE_ID"

  # A populated DOCKER_USERNAME env var triggers login
  if [[ -n "${DOCKER_USERNAME:-}" ]]; then
    echo "$DOCKER_PASS" | docker login --username "$DOCKER_USERNAME" --password-stdin
  fi
fi

# Build a specified platform image if requested
if [[ -n "${PLATFORMS}" ]]; then
  DOCKER_CMD="${DOCKER_CMD} --platform=${PLATFORMS}"
  if [[ "$PLATFORMS" == *,* ]]; then ## Multi-arch
    DOCKER_CMD="${DOCKER_CMD} --push"
  else ## Single arch
    # Build a single-arch image since the image name does not include a slash which
    # indicates that a registry is not available.
    #
    # Building a single-arch image with buildx and having the resulting image show up
    # in the local store of docker images (ala 'docker build') requires explicitly
    # loading it from the buildx store with '--load'.
    DOCKER_CMD="${DOCKER_CMD} --load"
  fi
else ## Single arch no explicit platform
  DOCKER_CMD="${DOCKER_CMD} --load"
fi

VM_ID=${VM_ID:-"${DEFAULT_VM_ID}"}
if [[ "${VM_ID}" != "${DEFAULT_VM_ID}" ]]; then
  DOCKERHUB_TAG="${VM_ID}-${DOCKERHUB_TAG}"
fi

# Default to the release image. Will need to be overridden when testing against unreleased versions.
AVALANCHEGO_NODE_IMAGE="${AVALANCHEGO_NODE_IMAGE:-${AVALANCHEGO_IMAGE_NAME}:${AVALANCHE_VERSION}}"

echo "docker cmd is $DOCKER_CMD"
echo "Building Docker Image: $DOCKERHUB_REPO:$BUILD_IMAGE_ID based of AvalancheGo@$AVALANCHE_VERSION"
${DOCKER_CMD} -t "$DOCKERHUB_REPO:$BUILD_IMAGE_ID" -t "$DOCKERHUB_REPO:${DOCKERHUB_TAG}" \
  "$SUBNET_EVM_PATH" -f "$SUBNET_EVM_PATH/Dockerfile" \
  --build-arg AVALANCHEGO_NODE_IMAGE="$AVALANCHEGO_NODE_IMAGE" \
  --build-arg SUBNET_EVM_COMMIT="$SUBNET_EVM_COMMIT" \
  --build-arg CURRENT_BRANCH="$CURRENT_BRANCH" \
  --build-arg VM_ID="$VM_ID"

if [[ -n "${PUBLISH}" && $CURRENT_BRANCH == "master" ]]; then
  echo "Tagging current image as $DOCKERHUB_REPO:latest"
  docker buildx imagetools create -t "$DOCKERHUB_REPO:latest" "$DOCKERHUB_REPO:$BUILD_IMAGE_ID"
fi
