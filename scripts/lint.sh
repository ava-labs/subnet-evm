#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

golangci-lint run --path-prefix=. --skip-dirs=coreth --timeout 3m
