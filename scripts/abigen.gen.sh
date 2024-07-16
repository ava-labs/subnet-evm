#!/usr/bin/env bash

set -eu;

if ! [[ "$0" =~ scripts/abigen.gen.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

go install github.com/ethereum/go-ethereum/cmd/abigen;

# TODO(arr4n) There are go:generate directives in geth that are out of sync with
# our checked-in code, so we limit the scope here. This needs to be changed to
# ./... once we have a proper geth fork.
go generate ./contracts/... ./testing/...;