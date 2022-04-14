#!/usr/bin/env bash
set -e

if ! [[ "$0" =~ scripts/test.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

SUBNET_EVM_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )

# Load the versions
source "$SUBNET_EVM_PATH"/scripts/versions.sh

# Load the constants
source "$SUBNET_EVM_PATH"/scripts/constants.sh

VERSION=$1
if [[ -z "${VERSION}" ]]; then
  echo "Missing version argument!"
  echo "Usage: ${0} [VERSION]" >> /dev/stderr
  exit 255
fi

GENESIS_ADDRESS="0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
GENESIS_PATH="scripts/tests/tx_allow_list_genesis.json"
CONTRACT_DIR="contract-examples"
TEST_PATH="test/ExampleTxAllowList.ts"

# download avalanchego
# https://github.com/ava-labs/avalanchego/releases
GOARCH=$(go env GOARCH)
GOOS=$(go env GOOS)
DOWNLOAD_URL=https://github.com/ava-labs/avalanchego/releases/download/v${VERSION}/avalanchego-linux-${GOARCH}-v${VERSION}.tar.gz
DOWNLOAD_PATH=/tmp/avalanchego.tar.gz
if [[ ${GOOS} == "darwin" ]]; then
  DOWNLOAD_URL=https://github.com/ava-labs/avalanchego/releases/download/v${VERSION}/avalanchego-macos-v${VERSION}.zip
  DOWNLOAD_PATH=/tmp/avalanchego.zip
fi

rm -rf /tmp/avalanchego-v${VERSION}
rm -rf /tmp/avalanchego-build
rm -f ${DOWNLOAD_PATH}

echo "downloading avalanchego ${VERSION} at ${DOWNLOAD_URL}"
curl -L ${DOWNLOAD_URL} -o ${DOWNLOAD_PATH}

echo "extracting downloaded avalanchego"
if [[ ${GOOS} == "linux" ]]; then
  tar xzvf ${DOWNLOAD_PATH} -C /tmp
elif [[ ${GOOS} == "darwin" ]]; then
  unzip ${DOWNLOAD_PATH} -d /tmp/avalanchego-build
  mv /tmp/avalanchego-build/build /tmp/avalanchego-v${VERSION}
fi
find /tmp/avalanchego-v${VERSION}

# Check if SUBNET_EVM_COMMIT is set, if not retrieve the last commit from the repo.
# This is used in the Dockerfile to allow a commit hash to be passed in without
# including the .git/ directory within the Docker image.
subnet_evm_commit=${SUBNET_EVM_COMMIT:-$( git rev-list -1 HEAD )}

# Build Subnet EVM, which is run as a subprocess
echo "Building Subnet EVM Version: $subnet_evm_version; GitCommit: $subnet_evm_commit"
go build \
-ldflags "-X github.com/ava-labs/subnet_evm/plugin/evm.GitCommit=$subnet_evm_commit -X github.com/ava-labs/subnet_evm/plugin/evm.Version=$subnet_evm_version" \
-o /tmp/avalanchego-v${VERSION}/plugins/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy \
"plugin/"*.go
find /tmp/avalanchego-v${VERSION}

# Create genesis file to use in network (make sure to add your address to
# "alloc")
export CHAIN_ID=99999

echo "building runner"
pushd ./runner
go build -v -o /tmp/runner .
popd

# first argument is genesis, second is hardhat test path
runTest () {
  echo "launching subnet in the background for genesis $1 and test $2"
  /tmp/runner \
  --avalanchego-path=/tmp/avalanchego-v${VERSION}/avalanchego \
  --vm-id=srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy \
  --vm-genesis-path=$1 \
  --output-path=/tmp/avalanchego-v${VERSION}/output.yaml 1> /dev/null &
  PID=${!}

  sleep 30
  while [[ ! -s /tmp/avalanchego-v${VERSION}/output.yaml ]]; do
    echo "waiting for local cluster on PID ${PID}"
    sleep 5
    # wait up to 5-minute
    ((c++)) && ((c==60)) && break
  done

  if [[ -f "/tmp/avalanchego-v${VERSION}/output.yaml" ]]; then
    echo "cluster is ready!"
    go run scripts/tests/main.go /tmp/avalanchego-v${VERSION}/output.yaml $CHAIN_ID $GENESIS_ADDRESS
  else
    echo "cluster is not ready in time... terminating ${PID}"
    kill ${PID}
    exit 255
  fi

  pushd ${CONTRACT_DIR}
  if yarn hardhat test $2 --network subnet; then
    echo "killing subnet"
    kill ${PID}
    echo "tests passed successfully"
    sleep 2s
  else
    echo "killing subnet"
    kill ${PID}
    echo "some tests failed"
    sleep 2s
  fi
  popd
  rm /tmp/avalanchego-v${VERSION}/output.yaml
}

runTest ${GENESIS_PATH} ${TEST_PATH}
runTest "scripts/tests/deployer_allow_list_genesis.json" "test/ExampleDeployerList.ts"