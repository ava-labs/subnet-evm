#!/usr/bin/env bash
set -e
source ./scripts/utils.sh

avalanche network stop --snapshot-name snap1

./scripts/build.sh custom_evm.bin

avalanche subnet upgrade vm hubblenet --binary custom_evm.bin --local

# utse tee to keep showing outut while storing in a var
OUTPUT=$(avalanche network start --avalanchego-version v1.10.0 --snapshot-name snap1 --config .avalanche-cli.json | tee /dev/fd/2)

setStatus
