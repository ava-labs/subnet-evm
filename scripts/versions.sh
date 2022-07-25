#!/usr/bin/env bash

# Set up the versions to be used
subnet_evm_version=${SUBNET_EVM_VERSION:-'v0.2.7'}
# Don't export them as they're used in the context of other calls
avalanche_version=${AVALANCHE_VERSION:-'v1.7.16'}
network_runner_version=${NETWORK_RUNNER_VERSION:-'v1.1.4'}
