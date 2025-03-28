# Releasing

## When to release

- When [AvalancheGo](https://github.com/ava-labs/avalanchego/releases) increases its RPC chain VM protocol version, which you can also check in [its `version/compatibility.json`](https://github.com/ava-labs/avalanchego/blob/master/version/compatibility.json)
- When Subnet-EVM needs to release a new feature or bug fix.

## Version semantics

## Procedure

In this document, we create a release `v0.7.3-rc.0` and the releaser Github username is `myusername`. We therefore assign these environment variables to simplify copying instructions:

```bash
export VERSION=v0.7.3-rc.0
export SEMVER_VERSION=v0.7.3
```

1. Create your branch, usually from the tip of the `master` branch:

    ```bash
    git fetch origin master:master
    git checkout master
    git checkout -b "releases/$VERSION"
    ```

1. Modify [plugin/evm/version.go](../../plugin/evm/version.go)'s `Version` global string variable and set it to the desired `$SEMVER_VERSION`.
1. Check the RPC chain VM versions are compatible by running:

    ```bash
    go test -run ^TestCompatibility$ github.com/ava-labs/subnet-evm/plugin/evm
    ```

    If the test fails:

    - First, check the [compatiblity.json](../../compatibility.json): if `$SEMVER_VERSION` is missing, you need to add it to the `"rpcChainVMProtocolVersion"` JSON object. In our example, we add it as

        ```json
        {
            "rpcChainVMProtocolVersion": {
                "v0.7.3": 39,
                "v0.7.2": 39,
                "v0.7.1": 39,
                "v0.7.0": 38
            }
        }
        ```

        💁 If you are unsure about the RPC chain VM protocol version:

        1. Check [go.mod](../../go.mod) and spot the version used for `github.com/ava-labs/avalanchego`. For example `v1.12.3-0.20250321175346-50f1601bf39a`.
        1. Refer to the [Avalanchego repository `version/compatibility.json` file](https://github.com/ava-labs/avalanchego/blob/master/version/compatibility.json) to find the RPC chain VM protocol version matching the subnet-evm AvalancheGo version. In our case, we use an AvalancheGo version based on top of `v1.12.2`, so the RPC chain VM protocol version is `39`:

            ```json
            {
                "39": [
                    "v1.12.2",
                    "v1.13.0"
                ],
            }
            ```

        That should allow the test to pass.
    - If the target `$SEMVER_VERSION` is already present in [compatiblity.json](../../compatibility.json), you might then have to change the AvalancheGo dependency to be compatible with subnet-evm's RPC chain VM protocol version:
        1. Refer to the [Avalanchego repository `version/compatibility.json` file](https://github.com/ava-labs/avalanchego/blob/master/version/compatibility.json) to find the AvalancheGo version supporting the subnet-evm RPC chain VM protocol version. For example, if the subnet-evm RPC chain VM protocol version is `39`, we can spot the AvalancheGo version `v1.12.2` supporting such version:

            ```json
            {
                "39": [
                    "v1.12.2",
                    "v1.13.0"
                ],
            }
            ```

        1. Change the subnet-evm avalanchego dependency using the version previously found with:

            ```bash
            go get github.com/ava-labs/avalanchego@v1.13.0
            ```

        That should allow the test to pass.
1. Specify the AvalancheGo compatibility in the [README.md relevant section](../../README.md#avalanchego-compatibility). For example we would add:

    ```text
    ...
    [v0.7.3] AvalancheGo@v1.12.2/1.13.0-fuji (Protocol Version: 39)
    ```

1. Commit your changes and push the branch

    ```bash
    git add plugin/evm/version.go compatibility.json go.mod go.sum README.md
    git commit -S -m "chore: release $VERSION"
    git push -u origin "releases/$VERSION"
    ```

1. Create a pull request (PR) from your branch targeting master, for example using [`gh`](https://cli.github.com/):

    ```bash
    gh pr create --repo github.com/ava-labs/subnet-evm --base master --title "chore: release $VERSION"
    ```

1. Once the PR checks pass, squash and merge it
1. There are two cases:
    - You are creating a release candidate (`-rc.`): update your master branch, tag it and push the tag:

        ```bash
        git checkout master
        git fetch origin master:master
        git tag "$VERSION"
        git push -u origin "$VERSION"
        ```

    - You are creating a release (i.e. `v0.7.4`): create a new release through the [Github web interface](https://github.com/ava-labs/subnet-evm/releases/new), targeting the master branch and creating the tag there.

We first deploy RC to a local node (I prefer to bootstrap canaries echo/dispatch from scratch in Fuji)
If all good then we update echo/dispatch with RC under <https://github.com/ava-labs/external-plugins-builder>
Confirm goreleaser job has successfully generated binaries by checking the releases page
