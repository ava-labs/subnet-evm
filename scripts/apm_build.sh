#!/usr/bin/env bash

if ! [[ "$0" =~ scripts/apm_build.sh ]]; then
    echo "must be run from repository root"
    exit 255
fi

# Create output build directory
mkdir -p ./build

source ./scripts/build.sh ./build/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy
