#!/usr/bin/env bash

set -euo pipefail

# Mostly taken from https://github.com/golang/go/issues/46312#issuecomment-1153345129

# Directory above this script
SUBNET_EVM_PATH=$(
  cd "$(dirname "${BASH_SOURCE[0]}")"
  cd .. && pwd
)

source "$SUBNET_EVM_PATH"/scripts/constants.sh

fuzzTime=${1:-1}
files=$(grep -r --include='**_test.go' --files-with-matches 'func Fuzz' .)
failed=false
for file in ${files}; do
  funcs=$(grep -oP 'func \K(Fuzz\w*)' $file)
  for func in ${funcs}; do
    echo "Fuzzing $func in $file"
    parentDir=$(dirname $file)
    go test $parentDir -run=$func -fuzz=$func -fuzztime=${fuzzTime}s -tags fuzz
    # If any of the fuzz tests fail, return exit code 1
    if [ $? -ne 0 ]; then
      failed=true
    fi
  done
done

if $failed; then
  exit 1
fi
