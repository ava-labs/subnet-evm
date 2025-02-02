#!/usr/bin/env bash

# Ignore warnings about variables appearing unused since this file is not the consumer of the variables it defines.
# shellcheck disable=SC2034

# Don't export them as they're used in the context of other calls
# When updating this version, make sure to also update:
# - avalanchego ref in flake.nix
# - run-monitored-tmpnet-cmd version in .github/workflows/tests.yml
AVALANCHE_VERSION=${AVALANCHE_VERSION:-'198b68f0a850fbfa12e50735bed56b14d99fe0f1'}
GINKGO_VERSION=${GINKGO_VERSION:-'v2.2.0'}
