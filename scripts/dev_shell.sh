#!/usr/bin/env bash

# Load AVALANCHE_VERSION
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck source=/scripts/versions.sh
source "$SCRIPT_DIR"/versions.sh

# Start a dev shell with the avalanchego flake
FLAKE="github:ava-labs/avalanchego?ref=${AVALANCHE_VERSION}"
echo "Starting nix shell for ${FLAKE}"
nix develop "${FLAKE}" "${@}"
