
#!/usr/bin/env bash
set -e
source ./scripts/utils.sh

if ! [[ "$0" =~ scripts/run_local.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

avalanche network clean

./scripts/build.sh custom_evm.bin

FILE=/tmp/validator.pk
if [ ! -f "$FILE" ]
then
    echo "$FILE does not exist; creating"
    echo "31b571bf6894a248831ff937bb49f7754509fe93bbd2517c9c73c4144c0e97dc" > $FILE
fi

avalanche subnet create hubblenet --force --custom --genesis genesis.json --vm custom_evm.bin --config .avalanche-cli.json

# configure and add chain.json
avalanche subnet configure hubblenet --chain-config chain.json --config .avalanche-cli.json
avalanche subnet configure hubblenet --subnet-config subnet.json --config .avalanche-cli.json
# avalanche subnet configure hubblenet --per-node-chain-config node_config.json --config .avalanche-cli.json

# use the same avalanchego version as the one used in subnet-evm
# use tee to keep showing outut while storing in a var
OUTPUT=$(avalanche subnet deploy hubblenet -l --avalanchego-version v1.10.13 --config .avalanche-cli.json | tee /dev/fd/2)

setStatus
