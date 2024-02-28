# Developing with tmpnet

The `load/` and `warp/` paths contain end-to-end (e2e) tests that use
the [tmpnet
fixture](https://github.com/ava-labs/avalanchego/blob/dev/tests/fixture/tmpnet/README.md). By
default both test suites use the tmpnet fixture to create a temporary
network that exists for only the duration of their execution.

It is possible to create a temporary network that can be reused across
test runs to minimize the setup cost involved:

```bash
# From the root of a clone of avalanchego, build the tmpnetctl cli
$ ./scripts/build_tmpnetctl.sh

# Start a new temporary network configured with subnet-evm's default plugin path
$ ./build/tmpnetctl start-network \
  --avalanche-path=./build/avalanchego
  --plugin-dir=$GOPATH/src/github.com/ava-labs/avalanchego/build/plugins

# From the root of a clone of subnet-evm, execute the warp test suite against the existing network
$ ginkgo -vv ./tests/warp -- --use-existing-network --network-dir=$HOME/.tmpnet/networks/latest
```

The network started by `tmpnetctl` won't come with subnets configured,
so the test suite will add them to the network the first time it
runs. Subsequent test runs will be able to reuse those subnets without
having to set them up.
