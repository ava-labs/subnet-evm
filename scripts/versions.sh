#!/usr/bin/env bash

# Ignore warnings about variables appearing unused since this file is not the consumer of the variables it defines.
# shellcheck disable=SC2034

# Don't export them as they're used in the context of other calls
AVALANCHE_VERSION=${AVALANCHE_VERSION:-'fdadc08156dc660ef116afc4134a589af31f4bc2'}
GINKGO_VERSION=${GINKGO_VERSION:-'v2.2.0'}

# This won't be used, but it's here to make code syncs easier
LATEST_CORETH_VERSION='v0.13.7'
