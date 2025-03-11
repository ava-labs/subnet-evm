#!/usr/bin/env bash

set -euo pipefail

# Updates all avalanchego version references to the module version

if ! [[ "$0" =~ scripts/set_avalanchego_version.sh ]]; then
  echo "must be run from repository root, but got $0"
  exit 255
fi

# TODO(marun) Reduce duplication between this script and versions.sh. The
# call to curl to get the full module hash fails on darwin runners tasked
# with running the unit test suite, so it's not an option to just include
# this in version.sh without further work, and that call further depends on
# the module hash, which depends on getting the module details from go mod.

# Get version from go.mod
MODULE_DETAILS="$(go list -m "github.com/ava-labs/avalanchego" 2>/dev/null)"

# Extract just the version part
AVALANCHE_VERSION="$(echo "${MODULE_DETAILS}" | awk '{print $2}')"

# Check if it matches pseudo-version pattern:
# v*YYYYMMDDHHMMSS-abcdef123456
#
# If not, the value is assumed to represent a tag
if [[ "${AVALANCHE_VERSION}" =~ ^v.*[0-9]{14}-[0-9a-f]{12}$ ]]; then
  # Extract module hash from version
  MODULE_HASH="$(echo "${AVALANCHE_VERSION}" | cut -d'-' -f3)"

  # The first 8 chars of the hash is used as the tag of avalanchego images
  AVALANCHE_VERSION="${MODULE_HASH::8}"

  FULL_AVALANCHEGO_VERSION="$(curl -s "https://api.github.com/repos/ava-labs/avalanchego/commits/${MODULE_HASH}" | grep '"sha":' | head -n1 | cut -d'"' -f4)"
fi

# FULL_AVALANCHEGO_VERSION needs to either be a tag or the full SHA to
# be usable for custom github action references.
if [[ -z "${FULL_AVALANCHEGO_VERSION:-}" ]]; then
  # Assume AVALANCHEGO_VERSION is a tag.
  FULL_AVALANCHEGO_VERSION="${AVALANCHE_VERSION}"
fi

# Ensure the monitoring version is the same as the avalanche version
CUSTOM_ACTION="ava-labs/avalanchego/.github/actions/run-monitored-tmpnet-cmd"
sed -i "s|\(uses: ${CUSTOM_ACTION}\)@.*|\1@${FULL_AVALANCHEGO_VERSION}|g" .github/workflows/tests.yml

# Ensure the flake version is the same as the avalanche version
FLAKE="github:ava-labs/avalanchego"
sed -i "s|\(${FLAKE}?ref=\).*|\1${FULL_AVALANCHEGO_VERSION}\";|g" flake.nix
