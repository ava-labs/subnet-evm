#!/usr/bin/env bash
set -e

# This script starts a single node running the local network with staking disabled and expects the caller to take care of cleaning up the created AvalancheGo process.
# This script uses the data directory to use node configuration and populate data into a predictable location.
if ! [[ "$0" =~ scripts/run_single_node.sh ]]; then
  echo "must be run from repository root, but got $0"
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

# Set up avalanche binary path and assume build directory is set
AVALANCHEGO_PATH=${AVALANCHEGO_PATH:-"$GOPATH/src/github.com/ava-labs/avalanchego/build/avalanchego"}
AVALANCHEGO_PLUGIN_DIR=${AVALANCHEGO_PATH:-"$AVALANCHEGO_PATH/plugins"}
DATA_DIR=${DATA_DIR:-/tmp/subnet-evm-start-node/$(date "+%Y-%m-%d%:%H:%M:%S")}

# Create node config in DATA_DIR
echo "creating node config"
mkdir -p $DATA_DIR
CONFIG_FILE_PATH=$DATA_DIR/config.json

  cat <<EOF >$CONFIG_FILE_PATH
{
  "network-id": "local",
  "staking-enabled": false,
  "network-health-min-conn-peers": 0,
  "network-health-max-time-since-msg-received": 4611686018427387904,
  "network-health-max-time-since-msg-sent": 4611686018427387904,
  "health-check-frequency": "5s",
  "plugin-dir": "$AVALANCHEGO_PLUGIN_DIR"
}
EOF

echo "Starting AvalancheGo node"
echo "AvalancheGo Binary Path: ${AVALANCHEGO_PATH}"
echo "Data directory: ${DATA_DIR}"
echo "Config file: ${CONFIG_FILE_PATH}"

# Run the node
$AVALANCHEGO_PATH --data-dir=$DATA_DIR --config-file=$CONFIG_FILE_PATH
