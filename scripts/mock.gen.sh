#!/usr/bin/env bash

set -euo pipefail

go generate -run "mockgen" ./...
