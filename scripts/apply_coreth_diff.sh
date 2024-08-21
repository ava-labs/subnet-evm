#!/bin/bash

set -e;
set -u;

ROOT=$(git rev-parse --show-toplevel);
cd "${ROOT}";

BASE="${1}";

git diff .."${BASE}" --binary | git apply --whitespace=nowarn

sed_command='s!github.com/ava-labs/coreth!github.com/ava-labs/subnet-evm!g'

# TODO: improve this command that finds all the "coreth" references and replaces them with "subnet-evm"
LANG=C find . -type f \! -name 'apply_coreth_diff.sh' \! -path './.git/*' \! -path './contracts/node_modules/*' -exec sed -i '' -e "${sed_command}" {} \;
gofmt -w .
go mod tidy

# Restore contracts/.gitignore
git checkout -- contracts/.gitignore