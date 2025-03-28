#!/usr/bin/env bash
# Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
# See the file LICENSE for licensing terms.

# Ignore warnings about variables appearing unused since this file is not the consumer of the variables it defines.
# shellcheck disable=SC2034

set -euo pipefail

SUBNET_EVM_PATH=$(
  cd "$(dirname "${BASH_SOURCE[0]}")"
  cd .. && pwd
)

# WARNING: this will use the most recent commit even if there are un-committed changes present
GIT_COMMIT="${SUBNET_EVM_COMMIT:-$(git --git-dir="${SUBNET_EVM_PATH}/.git" rev-parse HEAD)}"
COMMIT_HASH="${GIT_COMMIT::8}"
