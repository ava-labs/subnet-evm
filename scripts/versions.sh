#!/usr/bin/env bash

# Ignore warnings about variables appearing unused since this file is not the consumer of the variables it defines.
# shellcheck disable=SC2034

# Don't export them as they're used in the context of other calls
<<<<<<< HEAD
AVALANCHE_VERSION=${AVALANCHE_VERSION:-'v1.11.9'}
GINKGO_VERSION=${GINKGO_VERSION:-'v2.2.0'}

# This won't be used, but it's here to make code syncs easier
LATEST_CORETH_VERSION='7684836'
=======
AVALANCHE_VERSION=${AVALANCHE_VERSION:-'87a2b4f'}
>>>>>>> 34a7752258 (Update to latest p2p API (#594))
