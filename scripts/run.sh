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
AVALANCHEGO_BUILD_PATH=${AVALANCHEGO_BUILD_PATH:-"$GOPATH/src/github.com/ava-labs/avalanchego/build"}
AVALANCHEGO_PATH=${AVALANCHEGO_PATH:-"$AVALANCHEGO_BUILD_PATH/avalanchego"}
AVALANCHEGO_PLUGIN_DIR=${AVALANCHEGO_PLUGIN_DIR:-"$AVALANCHEGO_BUILD_PATH/plugins"}
DATA_DIR=${DATA_DIR:-/tmp/subnet-evm-start-node/$(date "+%Y-%m-%d%:%H:%M:%S")}

# Set the config file contents for the path passed in as the first argument
function _set_config(){
  cat <<EOF >$1
  {
    "network-id": "local",
    "staking-enabled": false,
    "health-check-frequency": "5s",
    "plugin-dir": "$AVALANCHEGO_PLUGIN_DIR"
  }
EOF
}

DATA_DIRS=("$DATA_DIR/node1" "$DATA_DIR/node2" "$DATA_DIR/node3" "$DATA_DIR/node4" "$DATA_DIR/node5")
CMDS=()
for (( i=0; i <5; i++ ))
do
  echo "Creating data directory: ${DATA_DIRS[i]}"
  mkdir -p ${DATA_DIRS[i]}
  NODE_DATA_DIR=${DATA_DIRS[i]}
  NODE_CONFIG_FILE_PATH="$NODE_DATA_DIR/config.json"
  _set_config $NODE_CONFIG_FILE_PATH
  
  CMD="$AVALANCHEGO_PATH --data-dir=$NODE_DATA_DIR"
  if [ $i -gt 0 ]; then
    echo "Adding CLI options for node$(($i+1))"
    CMD="$CMD --staking-port=$((9651+2*$i)) --http-port=$((9650+2*$i)) --bootstrap-ips=127.0.0.1:9651 --bootstrap-ids=NodeID-5Q84knz4UyG4AkUdo7r2KiUiK6FHKxA8Z"
  fi
  
  echo "Created command $CMD"
  CMDS+=("$CMD")
done


echo "Starting AvalancheGo network with the commands:"
echo ""
for (( i=0; i<5; i++ ))
do
  echo "CMD $i : ${CMDS[i]}"
  echo ""
done

# cleanup sends SIGINT to each tracked process
function cleanup_process_group(){
  echo "Terminating AvalancheGo network process group"
  kill 0
}


(trap 'cleanup_process_group' SIGINT; ${CMDS[0]} & ${CMDS[1]} & ${CMDS[2]} & ${CMDS[3]} & ${CMDS[4]} & wait)
