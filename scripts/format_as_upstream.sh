#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set -x

script_dir=$(dirname "$0")

commit_msg_rename_packages_to_upstream="format: rename packages to coreth"

make_commit() {
  if git diff-index --cached --quiet HEAD --; then
    echo "No changes to commit."
  else
    git commit -m "$1"
  fi
}

sed_command='s!\([^/]\)github.com/ava-labs/subnet-evm!\1github.com/ava-labs/coreth!g'
find . \( -name '*.go' -o -name 'go.mod' -o -name 'build_test.sh' \) -exec sed -i '' -e "${sed_command}" {} \;
gofmt -w .
go mod tidy
git add -u .
make_commit "${commit_msg_rename_packages_to_upstream}"
