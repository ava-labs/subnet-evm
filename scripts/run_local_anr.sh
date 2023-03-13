#!/usr/bin/env bash
set -e

if ! [[ "$0" =~ scripts/run_local_anr.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

VERSION='v1.9.7'
ANR_VERSION='8438e423db523743c48bd178bc20642f9c3ba049'

# Load the versions
SUBNET_EVM_PATH=$(
  cd "$(dirname "${BASH_SOURCE[0]}")"
  cd .. && pwd
)

# Load the constants
source "$SUBNET_EVM_PATH"/scripts/constants.sh


############################
# download avalanchego
# https://github.com/ava-labs/avalanchego/releases
GOARCH=$(go env GOARCH)
GOOS=$(go env GOOS)
BASEDIR=/tmp/subnet-evm-runner
mkdir -p ${BASEDIR}
AVAGO_DOWNLOAD_URL=https://github.com/ava-labs/avalanchego/releases/download/${VERSION}/avalanchego-linux-${GOARCH}-${VERSION}.tar.gz
AVAGO_DOWNLOAD_PATH=${BASEDIR}/avalanchego-linux-${GOARCH}-${VERSION}.tar.gz
if [[ ${GOOS} == "darwin" ]]; then
  AVAGO_DOWNLOAD_URL=https://github.com/ava-labs/avalanchego/releases/download/${VERSION}/avalanchego-macos-${VERSION}.zip
  AVAGO_DOWNLOAD_PATH=${BASEDIR}/avalanchego-macos-${VERSION}.zip
fi

AVAGO_FILEPATH=${BASEDIR}/avalanchego-${VERSION}
if [[ ! -d ${AVAGO_FILEPATH} ]]; then
  if [[ ! -f ${AVAGO_DOWNLOAD_PATH} ]]; then
    echo "downloading avalanchego ${VERSION} at ${AVAGO_DOWNLOAD_URL} to ${AVAGO_DOWNLOAD_PATH}"
    curl -L ${AVAGO_DOWNLOAD_URL} -o ${AVAGO_DOWNLOAD_PATH}
  fi
  echo "extracting downloaded avalanchego to ${AVAGO_FILEPATH}"
  if [[ ${GOOS} == "linux" ]]; then
    mkdir -p ${AVAGO_FILEPATH} && tar xzvf ${AVAGO_DOWNLOAD_PATH} --directory ${AVAGO_FILEPATH} --strip-components 1
  elif [[ ${GOOS} == "darwin" ]]; then
    unzip ${AVAGO_DOWNLOAD_PATH} -d ${AVAGO_FILEPATH}
    mv ${AVAGO_FILEPATH}/build/* ${AVAGO_FILEPATH}
    rm -rf ${AVAGO_FILEPATH}/build/
  fi
  find ${BASEDIR}/avalanchego-${VERSION}
fi

AVALANCHEGO_PATH=${AVAGO_FILEPATH}/avalanchego
AVALANCHEGO_PLUGIN_DIR=${AVAGO_FILEPATH}/plugins


#################################
# compile subnet-evm
# Check if SUBNET_EVM_COMMIT is set, if not retrieve the last commit from the repo.
# This is used in the Dockerfile to allow a commit hash to be passed in without
# including the .git/ directory within the Docker image.
subnet_evm_commit=${SUBNET_EVM_COMMIT:-$(git rev-list -1 HEAD)}

# Build Subnet EVM, which is run as a subprocess
echo "Building Subnet EVM Version: $subnet_evm_version; GitCommit: $subnet_evm_commit"
go build \
  -ldflags "-X github.com/ava-labs/subnet_evm/plugin/evm.GitCommit=$subnet_evm_commit -X github.com/ava-labs/subnet_evm/plugin/evm.Version=$subnet_evm_version" \
  -o $AVALANCHEGO_PLUGIN_DIR/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy \
  "plugin/"*.go


export CHAIN_ID=99999
echo "creating genesis"

cp genesis.json $BASEDIR/genesis.json

#################################
# download avalanche-network-runner
# https://github.com/ava-labs/avalanche-network-runner
ANR_REPO_PATH=github.com/ava-labs/avalanche-network-runner
# version set
go install -v ${ANR_REPO_PATH}@${ANR_VERSION}

#################################
# run "avalanche-network-runner" server
GOPATH=$(go env GOPATH)
if [[ -z ${GOBIN+x} ]]; then
  # no gobin set
  BIN=${GOPATH}/bin/avalanche-network-runner
else
  # gobin set
  BIN=${GOBIN}/avalanche-network-runner
fi
echo "launch avalanche-network-runner in the background"
$BIN server \
  --log-level debug \
  --port=":12342" \
  --grpc-gateway-port=":12343" &
PID=${!}

CHAIN_CONFIG_PATH=${BASEDIR}/chain_config.json

cat <<EOF >$CHAIN_CONFIG_PATH
{
  "local-txs-enabled": true,
  "priority-regossip-frequency": "1s",
  "tx-regossip-max-size": 32,
  "priority-regossip-max-txs": 500,
  "priority-regossip-txs-per-address": 200,
  "priority-regossip-addresses": ["0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", "0x70997970C51812dc3A010C7d01b50e0d17dc79C8", "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC", "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"]
}
EOF

$BIN control start \
  --log-level debug \
  --endpoint="0.0.0.0:12342" \
  --number-of-nodes=5 \
  --dial-timeout 30s \
  --avalanchego-path ${AVALANCHEGO_PATH} \
  --plugin-dir ${AVALANCHEGO_PLUGIN_DIR} \
  --blockchain-specs '[{"vm_name": "subnetevm", "genesis": "/tmp/subnet-evm-runner/genesis.json", "chain_config": "'$CHAIN_CONFIG_PATH'"}]'
  # --blockchain-specs '[{"vm_name": "subnetevm", "genesis": "/tmp/subnet-evm.genesis.json", "chain_config": "'$CHAIN_CONFIG_PATH'", "network_upgrade": "'$NETWORK_UPGRADE_PATH'", "subnet_config": "'$SUBNET_CONFIG_PATH'"}]'



echo "pkill -P ${PID} && kill -2 ${PID} && pkill -9 -f srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy" > kill.sh
