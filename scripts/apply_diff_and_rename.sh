#!/bin/bash

# Add other repo as a remote: `git remote add -f coreth git@github.com:ava-labs/coreth.git`.
# Usage: ./scripts/apply_diff_and_rename.sh coreth/master (or subnet-evm/master)

set -e;
set -u;

ROOT=$(git rev-parse --show-toplevel);
cd "${ROOT}";

BASE="${1}";

git diff .."${BASE}" --binary | git apply --whitespace=nowarn

if [[ "${BASE}" == coreth* ]]; then
    echo "Replacing coreth with subnet-evm"
    sed_command='s!github.com/ava-labs/coreth!github.com/ava-labs/subnet-evm!g'
else
    echo "Replacing subnet-evm with coreth"
    sed_command='s!github.com/ava-labs/subnet-evm!github.com/ava-labs/coreth!g'
fi

# TODO: improve this command that finds all the "coreth" references and replaces them with "subnet-evm"
LANG=C find . -type f \! -name 'apply_diff_and_rename.sh' \! -path './.git/*' \! -path './contracts/node_modules/*' -exec sed -i '' -e "${sed_command}" {} \;
gofmt -w .
go mod tidy

# Restore contracts/.gitignore
git checkout -- contracts/.gitignore