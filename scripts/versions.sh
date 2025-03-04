#!/usr/bin/env bash

# Ignore warnings about variables appearing unused since this file is not the consumer of the variables it defines.
# shellcheck disable=SC2034

# Don't export them as they're used in the context of other calls
# When updating this version, make sure to also update:
# - avalanchego ref in flake.nix
# - run-monitored-tmpnet-cmd version in .github/workflows/tests.yml
AVALANCHE_VERSION=${AVALANCHE_VERSION:-'26eb5c858f7f26162efaf03a2e83487b8b750447'}
GINKGO_VERSION=${GINKGO_VERSION:-'v2.2.0'}
