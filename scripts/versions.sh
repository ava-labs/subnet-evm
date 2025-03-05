#!/usr/bin/env bash

# Ignore warnings about variables appearing unused since this file is not the consumer of the variables it defines.
# shellcheck disable=SC2034

# Don't export them as they're used in the context of other calls

if [[ -z ${AVALANCHE_VERSION:-} ]]; then
  # Get module details from go.mod
  MODULE_DETAILS="$(go list -m "github.com/ava-labs/avalanchego" 2>/dev/null)"

  # Extract the version part
  AVALANCHE_VERSION="$(echo "${MODULE_DETAILS}" | awk '{print $2}')"

  # Check if the version matches the pattern where the last part is the module hash
  # v*YYYYMMDDHHMMSS-abcdef123456
  #
  # If not, the value is assumed to represent a tag
  if [[ "${AVALANCHE_VERSION}" =~ ^v.*[0-9]{14}-[0-9a-f]{12}$ ]]; then
    # Extract module hash from version
    MODULE_HASH="$(echo "${AVALANCHE_VERSION}" | cut -d'-' -f3)"

    # The first 8 chars of the hash is used as the tag of avalanchego images
    AVALANCHE_VERSION="${MODULE_HASH::8}"

    # Get full hash from GitHub API
    FULL_AVALANCHEGO_VERSION="$(curl -s "https://api.github.com/repos/ava-labs/avalanchego/commits/${MODULE_HASH}" | grep '"sha":' | head -n1 | cut -d'"' -f4)"
  fi
fi

# FULL_AVALANCHEGO_VERSION needs to either be a tag or the full SHA to
# be usable for custom github action references.
if [[ -z "${FULL_AVALANCHEGO_VERSION:-}" ]]; then
  # Assume AVALANCHEGO_VERSION is a tag.
  FULL_AVALANCHEGO_VERSION="${AVALANCHE_VERSION}"
fi

GINKGO_VERSION=${GINKGO_VERSION:-'v2.2.0'}
