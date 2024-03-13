#!/usr/bin/env bash

# Ignore warnings about variables appearing unused since this file is not the consumer of the variables it defines.
# shellcheck disable=SC2034

# Don't export them as they're used in the context of other calls
<<<<<<< HEAD
AVALANCHE_VERSION=${AVALANCHE_VERSION:-'v1.11.2'}
GINKGO_VERSION=${GINKGO_VERSION:-'v2.2.0'}

# This won't be used, but it's here to make code syncs easier
LATEST_CORETH_VERSION='0.13.1-rc.5'
=======
avalanche_version=${AVALANCHE_VERSION:-'v1.11.3-stake-weighted-gossip.2'}
>>>>>>> 16cf2556ea (Integrate stake weighted gossip selection (#511))
