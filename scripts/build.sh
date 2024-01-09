#!/usr/bin/env bash

set -euo pipefail

go_version_minimum="1.20.12"

go_version() {
    go version | sed -nE -e 's/[^0-9.]+([0-9.]+).+/\1/p'
}

version_lt() {
    # Return true if $1 is a lower version than than $2,
    local ver1=$1
    local ver2=$2
    # Reverse sort the versions, if the 1st item != ver1 then ver1 < ver2
    if [[ $(echo -e -n "$ver1\n$ver2\n" | sort -rV | head -n1) != "$ver1" ]]; then
        return 0
    else
        return 1
    fi
}

if version_lt "$(go_version)" "$go_version_minimum"; then
    echo "SubnetEVM requires Go >= $go_version_minimum, Go $(go_version) found." >&2
    exit 1
fi

# Root directory
SUBNET_EVM_PATH=$(
    cd "$(dirname "${BASH_SOURCE[0]}")"
    cd .. && pwd
)

# Load the versions
source "$SUBNET_EVM_PATH"/scripts/versions.sh

# Load the constants
source "$SUBNET_EVM_PATH"/scripts/constants.sh

BINARY_PATH="./build/subnet-evm"
echo "Building Subnet EVM @ GitCommit: $SUBNET_EVM_COMMIT at $BINARY_PATH"
go build -ldflags "-X github.com/ava-labs/subnet-evm/plugin/evm.GitCommit=${SUBNET_EVM_COMMIT} ${STATIC_LD_FLAGS}" \
   -o "${BINARY_PATH}" ./plugin/

PLUGIN_DIR="$HOME/.avalanchego/plugins"
PLUGIN_PATH="${PLUGIN_DIR}/${PLUGIN_FILENAME}"
echo "Symlinking ${BINARY_PATH} to ${PLUGIN_PATH}"
mkdir -p "${PLUGIN_DIR}"
ln -sf "${PWD}/${BINARY_PATH}" "${PLUGIN_PATH}"
