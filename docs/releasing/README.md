# Releasing

## When to release

- When [AvalancheGo](https://github.com/ava-labs/avalanchego/releases) increases its RPC chain VM protocol version, which you can also check in [its `version/compatibility.json`](https://github.com/ava-labs/avalanchego/blob/master/version/compatibility.json)
- When Subnet-EVM needs to release a new feature or bug fix.

## Version semantics

## Procedure

### Release candidate

ℹ️ you should always create a release candidate first, and only if everything is fine, you can create a release.

In this section, we create a release candidate `v0.7.3-rc.0`. We therefore assign these environment variables to simplify copying instructions:

```bash
export VERSION_RC=v0.7.3-rc.0
export VERSION=v0.7.3
```

1. Create your branch, usually from the tip of the `master` branch:

    ```bash
    git fetch origin master:master
    git checkout master
    git checkout -b "releases/$VERSION_RC"
    ```

1. Modify [plugin/evm/version.go](../../plugin/evm/version.go)'s `Version` global string variable and set it to the desired `$VERSION`.
1. Ensure the AvalancheGo version used in [go.mod](../../go.mod) is [its last release](https://github.com/ava-labs/avalanchego/releases). If not, upgrade it with, for example:

    ```bash
    go get github.com/ava-labs/avalanchego@v1.13.0
    go mod tidy
    ```

    And fix any errors that may arise from the upgrade. If it requires significant changes, you may want to create a separate PR for the upgrade and wait for it to be merged before continuing with this procedure.
1. Modify [compatiblity.json](../../compatibility.json) by adding `$VERSION` to the `"rpcChainVMProtocolVersion"` JSON object. In our example, we add it as

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

    1. Check [go.mod](../../go.mod) and spot the version used for `github.com/ava-labs/avalanchego`. For example `v1.13.0`.
    1. Refer to the [Avalanchego repository `version/compatibility.json` file](https://github.com/ava-labs/avalanchego/blob/master/version/compatibility.json) to find the RPC chain VM protocol version matching the subnet-evm AvalancheGo version. In our case, we use an AvalancheGo version `v1.13.0`, so the RPC chain VM protocol version is `39`:

        ```json
        {
            "39": [
                "v1.12.2",
                "v1.13.0"
            ],
        }
        ```

    Finally, check the RPC chain VM protocol version compatibility is setup properly by running:

    ```bash
    go test -run ^TestCompatibility$ github.com/ava-labs/subnet-evm/plugin/evm
    ```

1. Specify the AvalancheGo compatibility in the [README.md relevant section](../../README.md#avalanchego-compatibility). For example we would add:

    ```text
    ...
    [v0.7.3] AvalancheGo@v1.12.2/1.13.0-fuji/1.13.0 (Protocol Version: 39)
    ```

1. Commit your changes and push the branch

    ```bash
    git add .
    git commit -S -m "chore: release $VERSION_RC"
    git push -u origin "releases/$VERSION_RC"
    ```

1. Create a pull request (PR) from your branch targeting master, for example using [`gh`](https://cli.github.com/):

    ```bash
    gh pr create --repo github.com/ava-labs/subnet-evm --base master --title "chore: release $VERSION_RC"
    ```

1. Once the PR checks pass, squash and merge it
1. Update your master branch, create a tag and push it:

    ```bash
    git checkout master
    git fetch origin master:master
    git tag "$VERSION_RC"
    git push -u origin "$VERSION_RC"
    ```

1. Deploy the release candidate tagged to a local node (bootstrap canaries echo/dispatch from scratch in Fuji)
1. Update echo/dispatch with RC under <https://github.com/ava-labs/external-plugins-builder>

### Release

If you have create a successful release candidate, you can now create a release.

Following the previous example in the [Release candidate section](#release-candidate), we will create a release `v0.7.3` indicated by the `$VERSION` variable.

1. Create a new release through the [Github web interface](https://github.com/ava-labs/subnet-evm/releases/new)
    1. In the "Choose a tag" box, enter `$VERSION` (`v0.7.3`)
    1. In the "Target", pick the last release candidate branch `releases/${VERSION_RC}`, for example `releases/v0.7.3-rc.0`.
    Do not select `master` as the target branch to prevent adding new master branch commits to the release.
    1. Set the "Release title" to `$VERSION` (`v0.7.3`)
    1. Set the description
    1. Only tick the box "Set as the latest release"
    1. Click on the "Create release" button
1. Monitor the [release Github workflow](https://github.com/ava-labs/subnet-evm/actions/workflows/release.yml) to ensure the GoReleaser step succeeds and check the binaries are then published to [the releases page](https://github.com/ava-labs/subnet-evm/releases)
1. Monitor the [Publish Docker image workflow](https://github.com/ava-labs/subnet-evm/actions/workflows/publish_docker.yml) succeeds.
