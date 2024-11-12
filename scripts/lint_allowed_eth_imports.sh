#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Ensure that there are no eth imports that are not marked as explicitly allowed via ./scripts/eth-allowed-packages.txt
# 1. Recursively search through all go files for any lines that include a direct import from libevm
# 2. Ignore lines that import libevm with a named import starting with "eth"
# 3. Sort the unique results
# 4. Print out the difference between the search results and the list of specified allowed package imports from libevm.
libevm_regexp='"github.com/ava-labs/libevm/.*"'
allow_named_imports='eth\w\+ "'
extra_imports=$(grep -r --include='*.go' --exclude=mocks.go --exclude-dir='simulator' "${libevm_regexp}" -h | grep -v "${allow_named_imports}" | grep -o "${libevm_regexp}" | sort -u | comm -23 - ./scripts/eth-allowed-packages.txt)
if [ -n "${extra_imports}" ]; then
    echo "new ethereum imports should be added to ./scripts/eth-allowed-packages.txt to prevent accidental imports:"
    echo "${extra_imports}"
    exit 1
fi

extra_imports=$(grep -r --include='*.go' '"github.com/ava-labs/coreth/.*"' -o -h || true | sort -u)
if [ -n "${extra_imports}" ]; then
    echo "subnet-evm should not import packages from coreth:"
    echo "${extra_imports}"
    exit 1
fi
