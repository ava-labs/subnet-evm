#!/usr/bin/env bash

set -euo pipefail

# https://github.com/uber-go/mock
go install -v go.uber.org/mock/mockgen@v0.4.0

go generate -run "mockgen.+license_header" ./...
