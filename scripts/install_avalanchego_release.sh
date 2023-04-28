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

############################
# download avalanchego
# https://github.com/ava-labs/avalanchego/releases
GOARCH=$(go env GOARCH)
GOOS=$(go env GOOS)
BASEDIR=${BASEDIR:-"/tmp/avalanchego-release"}
AVALANCHEGO_BUILD_PATH=${AVALANCHEGO_BUILD_PATH:-${BASEDIR}/avalanchego}

mkdir -p ${BASEDIR}

AVAGO_DOWNLOAD_URL=https://github.com/ava-labs/avalanchego/releases/download/${AVALANCHEGO_VERSION}/avalanchego-linux-${GOARCH}-${AVALANCHEGO_VERSION}.tar.gz
AVAGO_DOWNLOAD_PATH=${BASEDIR}/avalanchego-linux-${GOARCH}-${AVALANCHEGO_VERSION}.tar.gz

if [[ ${GOOS} == "darwin" ]]; then
  AVAGO_DOWNLOAD_URL=https://github.com/ava-labs/avalanchego/releases/download/${AVALANCHEGO_VERSION}/avalanchego-macos-${AVALANCHEGO_VERSION}.zip
  AVAGO_DOWNLOAD_PATH=${BASEDIR}/avalanchego-macos-${AVALANCHEGO_VERSION}.zip
fi

BUILD_DIR=${AVALANCHEGO_BUILD_PATH}-${AVALANCHEGO_VERSION}

extract_archive() {
  mkdir -p ${BUILD_DIR}

  if [[ ${AVAGO_DOWNLOAD_PATH} == *.tar.gz ]]; then
    tar xzvf ${AVAGO_DOWNLOAD_PATH} --directory ${BUILD_DIR} --strip-components 1
  elif [[ ${AVAGO_DOWNLOAD_PATH} == *.zip ]]; then
    unzip ${AVAGO_DOWNLOAD_PATH} -d ${BUILD_DIR}
    mv ${BUILD_DIR}/build/* ${BUILD_DIR}
    rm -rf ${BUILD_DIR}/build/
  fi
}

# first check if we already have the archive
if [[ -f ${AVAGO_DOWNLOAD_PATH} ]]; then
  # if the download path already exists, extract and exit
  echo "found avalanchego ${AVALANCHEGO_VERSION} at ${AVAGO_DOWNLOAD_PATH}"

  extract_archive
else
  # try to download the archive if it exists
  if curl -s --head --request GET ${AVAGO_DOWNLOAD_URL} | grep "302" > /dev/null; then
    echo "${AVAGO_DOWNLOAD_URL} found"
    echo "downloading to ${AVAGO_DOWNLOAD_PATH}"
    curl -L ${AVAGO_DOWNLOAD_URL} -o ${AVAGO_DOWNLOAD_PATH}

    extract_archive
  else
    # else the version is a git commitish (or it's invalid)
    GIT_CLONE_URL=https://github.com/ava-labs/avalanchego.git
    GIT_CLONE_PATH=${BASEDIR}/avalanchego-repo/
    
    # check to see if the repo already exists, if not clone it 
    if [[ ! -d ${GIT_CLONE_PATH} ]]; then
      echo "cloning ${GIT_CLONE_URL} to ${GIT_CLONE_PATH}"
      git clone --no-checkout ${GIT_CLONE_URL} ${GIT_CLONE_PATH}
    fi

    # check to see if the commitish exists in the repo
    WORKDIR=$(pwd)

    cd ${GIT_CLONE_PATH}

    git fetch

    echo "checking out ${AVALANCHEGO_VERSION}"

    set +e
    # try to checkout the branch
    git checkout origin/${AVALANCHEGO_VERSION} > /dev/null 2>&1
    CHECKOUT_STATUS=$?
    set -e

    # if it's not a branch, try to checkout the commit 
    if [[ $CHECKOUT_STATUS -ne 0 ]]; then
      set +e
      git checkout ${AVALANCHEGO_VERSION} > /dev/null 2>&1
      CHECKOUT_STATUS=$?
      set -e

      if [[ $CHECKOUT_STATUS -ne 0 ]]; then
        echo
        echo "'${VERSION}' is not a valid release tag, commit hash, or branch name"
        exit 1
      fi
    fi

    COMMIT=$(git rev-parse HEAD)

    # use the commit hash instead of the branch name or tag
    BUILD_DIR=${AVALANCHEGO_BUILD_PATH}-${COMMIT}

    # if the build-directory doesn't exist, build avalanchego
    if [[ ! -d ${BUILD_DIR} ]]; then    
      echo "building avalanchego ${COMMIT} to ${BUILD_DIR}"
      ./scripts/build.sh
      mkdir -p ${BUILD_DIR}

      mv ${GIT_CLONE_PATH}/build/* ${BUILD_DIR}/
    fi

    cd $WORKDIR
  fi
fi

AVALANCHEGO_PATH=${AVALANCHEGO_BUILD_PATH}/avalanchego
AVALANCHEGO_PLUGIN_DIR=${AVALANCHEGO_BUILD_PATH}/plugins

mkdir -p ${AVALANCHEGO_BUILD_PATH}
    
cp ${BUILD_DIR}/avalanchego ${AVALANCHEGO_PATH}


echo "Installed AvalancheGo release ${AVALANCHEGO_VERSION}"
echo "AvalancheGo Path: ${AVALANCHEGO_PATH}"
echo "Plugin Dir: ${AVALANCHEGO_PLUGIN_DIR}"
