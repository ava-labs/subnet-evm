#!/usr/bin/env bash
set -e

if ! [[ "$0" =~ scripts/install_cli.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

# Load the versions
SUBNET_EVM_PATH=$(
  cd "$(dirname "${BASH_SOURCE[0]}")"
  cd .. && pwd
)
source "$SUBNET_EVM_PATH"/scripts/versions.sh

# Load the constants
source "$SUBNET_EVM_PATH"/scripts/constants.sh

# Place avalanche-cli in 
BUILD_DIR=${AVALANCHE_CLI_BIN:-"$SUBNET_EVM_PATH/build"}
AVALANCHE_CLI_BIN="$BUILD_DIR/avalanche" # No override - since install script places avalanche binary in specified BUILD_DIR

curl -sSfL https://raw.githubusercontent.com/ava-labs/avalanche-cli/main/scripts/install.sh | sh -s -- -b $BUILD_DIR $AVALANCHE_CLI_VERSION
