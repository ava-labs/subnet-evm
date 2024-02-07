#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

while true; do
    git cherry-pick "$@" && exit 0 || echo "attempting to resolve delete/update conflicts..."
    delete_update_conflicts=$(git status --porcelain | grep ^DU | cut -d' ' -f2 | xargs)
    git rm ${delete_update_conflicts}
    git cherry-pick --continue || exit 1
done