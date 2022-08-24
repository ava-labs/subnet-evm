#!/usr/bin/env bash

# Set up the versions to be used
subnet_evm_version=${SUBNET_EVM_VERSION:-'v0.2.9'}
# Don't export them as they're used in the context of other calls
avalanche_version=${AVALANCHE_VERSION:-'v1.7.18'}
network_runner_version=${NETWORK_RUNNER_VERSION:-'v1.2.0'}
ginkgo_version=${GINKGO_VERSION:-'v2.1.4'}

