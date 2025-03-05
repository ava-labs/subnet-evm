#!/usr/bin/env bash

set -euo pipefail

# Updates all avalanchego version references to the module version

if ! [[ "$0" =~ scripts/set_avalanchego_version.sh ]]; then
  echo "must be run from repository root, but got $0"
  exit 255
fi

# Load the avalanche version
source ./scripts/versions.sh

# Ensure the monitoring version is the same as the avalanche version
CUSTOM_ACTION="ava-labs/avalanchego/.github/actions/run-monitored-tmpnet-cmd"
sed -i "s|\(uses: ${CUSTOM_ACTION}\)@.*|\1@${FULL_AVALANCHEGO_VERSION}|g" .github/workflows/tests.yml

# Ensure the flake version is the same as the avalanche version
FLAKE="github:ava-labs/avalanchego"
sed -i "s|\(${FLAKE}?ref=\).*|\1${FULL_AVALANCHEGO_VERSION}\";|g" flake.nix
