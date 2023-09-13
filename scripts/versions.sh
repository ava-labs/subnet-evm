#!/usr/bin/env bash

# Don't export them as they're used in the context of other calls
<<<<<<< HEAD
AVALANCHE_VERSION=${AVALANCHE_VERSION:-'v1.10.10-rc.2'}
AVALANCHEGO_VERSION=${AVALANCHEGO_VERSION:-$AVALANCHE_VERSION}
GINKGO_VERSION=${GINKGO_VERSION:-'v2.2.0'}

# This won't be used, but it's here to make code syncs easier
LATEST_CORETH_VERSION='0.12.4-rc.0'
=======
avalanche_version=${AVALANCHE_VERSION:-'v1.10.10-rc.4'}
>>>>>>> 5034cf341 (Drop outbound gossip requests for non-validators (#334))
