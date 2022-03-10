# syntax=docker/dockerfile:experimental

# ============= Setting up base Stage ================
# Set required AVALANCHE_VERSION parameter in build image script
ARG AVALANCHE_VERSION=v1.7.6

# VMID generated (https://github.com/ava-labs/subnet-cli#subnet-cli-create-vmid) 
ARG VMID=srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy

# ============= Compilation Stage ================
FROM golang:1.17.4-buster AS builder

# Declare ARGs in order to inherit global values above
# Pass in SUBNET_EVM_COMMIT as an arg to allow the build script to set this externally. If ommited the latest commit ID will be used
ARG SUBNET_EVM_COMMIT
ARG VMID

RUN apt-get update && apt-get install -y --no-install-recommends bash=5.0-4 git=1:2.20.1-2+deb10u3 make=4.2.1-1.2 gcc=4:8.3.0-1 musl-dev=1.1.21-2 ca-certificates=20200601~deb10u2

WORKDIR /build

# Copy avalanche dependencies first (intermediate docker image caching)
# Copy avalanchego directory if present (for manual CI case, which uses local dependency)
COPY go.mod go.sum avalanchego* ./

# Download avalanche dependencies using go mod
RUN go mod download

# Copy the code into the container
COPY . .

RUN export SUBNET_EVM_COMMIT=$SUBNET_EVM_COMMIT && ./scripts/build.sh /build/$VMID

# ============= Cleanup Stage ================
FROM avaplatform/avalanchego:$AVALANCHE_VERSION AS builtImage

# Declare ARGs in order to inherit global values above
ARG VMID

# Copy the evm binary into the correct location in the container
COPY --from=builder /build/$VMID /avalanchego/build/plugins/$VMID
