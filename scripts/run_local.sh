
#!/usr/bin/env bash
set -e
source ./scripts/utils.sh

if ! [[ "$0" =~ scripts/run_local.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

if [[ -z "${VALIDATOR_PRIVATE_KEY}" ]]; then
  echo "VALIDATOR_PRIVATE_KEY must be set"
  exit 255
fi

avalanche network clean

./scripts/build.sh custom_evm.bin

avalanche subnet create hubblenet --force --custom --genesis genesis.json --vm custom_evm.bin --config .avalanche-cli.json

# configure and add chain.json
avalanche subnet configure hubblenet --chain-config chain.json --config .avalanche-cli.json
# avalanche subnet configure hubblenet --per-node-chain-config node_config.json --config .avalanche-cli.json

# use the same avalanchego version as the one used in subnet-evm
# use tee to keep showing outut while storing in a var
export OUTPUT=$(avalanche subnet deploy hubblenet -l --avalanchego-version v1.9.7 --config .avalanche-cli.json | tee /dev/fd/2)

setStatus
