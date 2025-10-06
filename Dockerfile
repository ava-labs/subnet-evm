# ============= Setting up base Stage ================
# AVALANCHEGO_NODE_IMAGE needs to identify an existing node image and should include the tag
# This value is not intended to be used but silences a warning
ARG AVALANCHEGO_NODE_IMAGE="invalid-image"

# ============= Compilation Stage ================
FROM --platform=$BUILDPLATFORM golang:1.24.7-bookworm AS builder

WORKDIR /build

# Copy module files first (improves Docker layer caching)
COPY go.mod go.sum ./
# Download module dependencies
RUN go mod download

# Copy the code into the container
COPY . .

# If a local avalanchego module is present, move it out of the module tree and
# point the replace directive at it to avoid flattening into this module's packages.
RUN if [ -f ./avalanchego/go.mod ]; then \
  mkdir -p /third_party && \
  mv ./avalanchego /third_party/avalanchego && \
  go mod edit -replace github.com/ava-labs/avalanchego=./third_party/avalanchego && \
  go mod tidy; \
fi

# Ensure pre-existing builds are not available for inclusion in the final image
RUN [ -d ./build ] && rm -rf ./build/* || true

ARG TARGETPLATFORM
ARG BUILDPLATFORM

# Configure a cross-compiler if the target platform differs from the build platform.
#
# build_env.sh is used to capture the environmental changes required by the build step since RUN
# environment state is not otherwise persistent.
RUN if [ "$TARGETPLATFORM" = "linux/arm64" ] && [ "$BUILDPLATFORM" != "linux/arm64" ]; then \
  apt-get update && apt-get install -y gcc-aarch64-linux-gnu && \
  echo "export CC=aarch64-linux-gnu-gcc" > ./build_env.sh \
  ; elif [ "$TARGETPLATFORM" = "linux/amd64" ] && [ "$BUILDPLATFORM" != "linux/amd64" ]; then \
  apt-get update && apt-get install -y gcc-x86-64-linux-gnu && \
  echo "export CC=x86_64-linux-gnu-gcc" > ./build_env.sh \
  ; else \
  echo "export CC=gcc" > ./build_env.sh \
  ; fi

# Pass in SUBNET_EVM_COMMIT as an arg to allow the build script to set this externally
ARG SUBNET_EVM_COMMIT
ARG CURRENT_BRANCH

RUN . ./build_env.sh && \
  echo "{CC=$CC, TARGETPLATFORM=$TARGETPLATFORM, BUILDPLATFORM=$BUILDPLATFORM}" && \
  export GOARCH=$(echo ${TARGETPLATFORM} | cut -d / -f2) && \
  export CURRENT_BRANCH=$CURRENT_BRANCH && \
  export SUBNET_EVM_COMMIT=$SUBNET_EVM_COMMIT && \
  ./scripts/build.sh build/subnet-evm

# ============= Cleanup Stage ================
FROM $AVALANCHEGO_NODE_IMAGE AS execution

# Copy the evm binary into the correct location in the container
ARG VM_ID=srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy
ENV AVAGO_PLUGIN_DIR="/avalanchego/build/plugins"
COPY --from=builder /build/build/subnet-evm $AVAGO_PLUGIN_DIR/$VM_ID
