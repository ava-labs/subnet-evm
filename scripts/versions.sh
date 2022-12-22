#!/usr/bin/env bash

# Set up the versions to be used
subnet_evm_version=${SUBNET_EVM_VERSION:-'v0.4.7'}
# Don't export them as they're used in the context of other calls
avalanche_version=${AVALANCHE_VERSION:-'v1.9.5'}
network_runner_version=${NETWORK_RUNNER_VERSION:-'35be10cd3867a94fbe960a1c14a455f179de60d9'}
ginkgo_version=${GINKGO_VERSION:-'v2.2.0'}

# This won't be used, but it's here to make code syncs easier
latest_coreth_version=0.11.3
