#!/usr/bin/env bash

# Ignore warnings about variables appearing unused since this file is not the consumer of the variables it defines.
# shellcheck disable=SC2034

# Don't export them as they're used in the context of other calls
AVALANCHE_VERSION=${AVALANCHE_VERSION:-'15c496b09f92cc5ac23b5aa2937d17a258d9a14f'}
GINKGO_VERSION=${GINKGO_VERSION:-'v2.2.0'}
