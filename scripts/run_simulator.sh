#!/usr/bin/env bash
set -e

#################################
echo "building simulator"
pushd ./cmd/simulator
go install -v .

popd
echo "running simulator"
simulator \
--cluster-info-yaml=$SIMULATOR_CLUSTER_YAML_FILE_PATH \
--keys=./cmd/simulator/.simulator/keys \
--timeout=30s \
--concurrency=10 \
--base-fee=25 \
--priority-fee=1
