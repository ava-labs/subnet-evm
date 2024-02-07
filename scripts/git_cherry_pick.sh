#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

while true; do
    git cherry-pick "$@" && exit 0 || echo "cherry-pick has conflicts, attempting to resolve..."
    exit 1
done

