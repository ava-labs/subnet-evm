#!/usr/bin/env bash

set -euo pipefail

# On Windows, ensure task uses bash for script execution
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || "$(uname -s)" == "MINGW"* ]]; then
  # Windows detected - set environment variable to force bash execution
  export TASK_SHELL=bash
fi

# Assume the system-installed task is compatible with the taskfile version
if command -v task > /dev/null 2>&1; then
  exec task "${@}"
else
  go run github.com/go-task/task/v3/cmd/task@v3.39.2 "${@}"
fi
