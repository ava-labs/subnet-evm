#!/usr/bin/env bash
set -e

# Load the versions
SUBNET_EVM_PATH=$(
  cd "$(dirname "${BASH_SOURCE[0]}")"
  cd .. && pwd
)
source "$SUBNET_EVM_PATH"/scripts/versions.sh

# Load the constants
source "$SUBNET_EVM_PATH"/scripts/constants.sh

VERSION=$AVALANCHEGO_VERSION

############################
# download avalanchego
# https://github.com/ava-labs/avalanchego/releases
GOARCH=$(go env GOARCH)
GOOS=$(go env GOOS)
BASEDIR=${BASE_DIR:-"/tmp/avalanchego-release"}
mkdir -p ${BASEDIR}
AVAGO_DOWNLOAD_URL=https://github.com/ava-labs/avalanchego/releases/download/${VERSION}/avalanchego-linux-${GOARCH}-${VERSION}.tar.gz
AVAGO_DOWNLOAD_PATH=${BASEDIR}/avalanchego-linux-${GOARCH}-${VERSION}.tar.gz
if [[ ${GOOS} == "darwin" ]]; then
  AVAGO_DOWNLOAD_URL=https://github.com/ava-labs/avalanchego/releases/download/${VERSION}/avalanchego-macos-${VERSION}.zip
  AVAGO_DOWNLOAD_PATH=${BASEDIR}/avalanchego-macos-${VERSION}.zip
fi

# TODO: Please read:
# not sure why the version isn't included in CI currently... open to suggestions on how to handle this better
AVALANCHEGO_BUILD_PATH=${AVALANCHEGO_BUILD_PATH-${VERSION}:-${BASEDIR}/avalanchego-${VERSION}}
mkdir -p $AVALANCHEGO_BUILD_PATH

if [[ ! -f ${AVAGO_DOWNLOAD_PATH} ]]; then
  echo "downloading avalanchego ${VERSION} at ${AVAGO_DOWNLOAD_URL} to ${AVAGO_DOWNLOAD_PATH}"

  # test if the tarball is already available for download
  if curl -s --head --request GET ${AVAGO_DOWNLOAD_URL} | grep "302" > /dev/null; then
    curl -L ${AVAGO_DOWNLOAD_URL} -o ${AVAGO_DOWNLOAD_PATH}
  else
    GIT_CLONE_URL=https://github.com/ava-labs/avalanchego.git
    GIT_CLONE_PATH=${BASEDIR}/avalanchego-git/
    mkdir -p $GIT_CLONE_PATH

    WORKDIR=$PWD
    cd $GIT_CLONE_PATH

    # if the git repo already exists, fetch, otherwise clone
    if [[ -d .git ]]; then
      git fetch
    else
      git clone ${GIT_CLONE_URL} .
    fi

    set +e
    git checkout ${VERSION}
    CHECKOUT_STATUS=$?
    set -e

    # if the previous command failed, exit
    if [[ $CHECKOUT_STATUS -ne 0 ]]; then
      echo
      echo "'${VERSION}' is not a valid release tag, commit hash, or branch name"
      exit 1
    fi

    # build avalanchego
    echo "building avalanchego ${VERSION}"
    ./scripts/build.sh

    # copy the build to the download path
    cp build/avalanchego ${AVALANCHEGO_BUILD_PATH}
    cd $WORKDIR
  fi
fi

if [[ ! -f ${AVALANCHEGO_BUILD_PATH}/avalanchego ]]; then
  echo "extracting downloaded avalanchego to ${AVALANCHEGO_BUILD_PATH}"

  if [[ ${GOOS} == "linux" ]]; then
    mkdir -p ${AVALANCHEGO_BUILD_PATH} && tar xzvf ${AVAGO_DOWNLOAD_PATH} --directory ${AVALANCHEGO_BUILD_PATH} --strip-components 1
  elif [[ ${GOOS} == "darwin" ]]; then
    unzip ${AVAGO_DOWNLOAD_PATH} -d ${AVALANCHEGO_BUILD_PATH}
    mv ${AVALANCHEGO_BUILD_PATH}/build/* ${AVALANCHEGO_BUILD_PATH}
    rm -rf ${AVALANCHEGO_BUILD_PATH}/build/
  fi
fi

AVALANCHEGO_PATH=${AVALANCHEGO_BUILD_PATH}/avalanchego
AVALANCHEGO_PLUGIN_DIR=${AVALANCHEGO_BUILD_PATH}/plugins

echo "Installed AvalancheGo release ${VERSION}"
echo "AvalancheGo Path: ${AVALANCHEGO_PATH}"
echo "Plugin Dir: ${AVALANCHEGO_PLUGIN_DIR}"
