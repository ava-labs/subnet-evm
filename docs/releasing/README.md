# Releasing

## When to release

- When [AvalancheGo](https://github.com/ava-labs/avalanchego/releases) increases its RPC chain VM protocol version, which you can also check in [its `version/compatibility.json`](https://github.com/ava-labs/avalanchego/blob/master/version/compatibility.json)
- When Subnet-EVM needs to release a new feature or bug fix.

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

1. Modify the [plugin/evm/version.go](../../plugin/evm/version.go) `Version` global string variable and set it to the desired `$VERSION`.
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

Q: what's the use to setup an AWS create platform account through okta, using the platform sandbox account.

1. Once the PR checks pass, squash and merge it
1. Update your master branch, create a tag and push it:

    ```bash
    git checkout master
    git fetch origin master:master
    git tag "$VERSION_RC"
    git push -u origin "$VERSION_RC"
    ```

Once the tag is created, you need to test it on the Fuji testnet both locally and then as canaries, using the Dispatch and Echo subnets.

#### Local deployment

1. Find the Dispatch and Echo L1s blockchain ID and subnet ID:
    - [Dispath L1 details](https://subnets-test.avax.network/dispatch/details). Its subnet id is `7WtoAMPhrmh5KosDUsFL9yTcvw7YSxiKHPpdfs4JsgW47oZT5`.
    - [Echo L1 details](https://subnets-test.avax.network/echo/details). Its subnet id is `i9gFpZQHPLcGfZaQLiwFAStddQD7iTKBpFfurPFJsXm1CkTZK`.
1. Get the VM IDs of the Echo and Dispatch L1s with:
    - Dispatch:

        ```bash
        curl -X POST --silent -H 'content-type:application/json' --data '{
            "jsonrpc": "2.0",
            "method": "platform.getBlockchains",
            "params": {},
            "id": 1
        }'  https://api.avax-test.network/ext/bc/P | \
        jq -r '.result.blockchains[] | select(.subnetID=="7WtoAMPhrmh5KosDUsFL9yTcvw7YSxiKHPpdfs4JsgW47oZT5") |  "\(.name)\nBlockchain id: \(.id)\nVM id: \(.vmID)\n"'
        ```

        Which as the time of this writing returns:

        ```text
        dispatch
        Blockchain id: 2D8RG4UpSXbPbvPCAWppNJyqTG2i2CAXSkTgmTBBvs7GKNZjsY
        VM id: mDtV8ES8wRL1j2m6Kvc1qRFAvnpq4kufhueAY1bwbzVhk336o
        ```

    - Echo:

        ```bash
        curl -X POST --silent -H 'content-type:application/json' --data '{
            "jsonrpc": "2.0",
            "method": "platform.getBlockchains",
            "params": {},
            "id": 1
        }'  https://api.avax-test.network/ext/bc/P | \
        jq -r '.result.blockchains[] | select(.subnetID=="i9gFpZQHPLcGfZaQLiwFAStddQD7iTKBpFfurPFJsXm1CkTZK") |  "\(.name)\nBlockchain id: \(.id)\nVM id: \(.vmID)\n"'
        ```

        Which as the time of this writing returns:

        ```text
        echo
        Blockchain id: 98qnjenm7MBd8G2cPZoRvZrgJC33JGSAAKghsQ6eojbLCeRNp
        VM id: meq3bv7qCMZZ69L8xZRLwyKnWp6chRwyscq8VPtHWignRQVVF
        ```

1. Clone [AvalancheGo](https://github.com/ava-labs/avalanchego):

    ```bash
    git clone git@github.com:ava-labs/avalanchego.git
    ```

1. Build AvalancheGo using those VM ids:

    ```bash
    cd avalanchego
    ./scripts/build.sh ~/.avalanchego/plugins/mDtV8ES8wRL1j2m6Kvc1qRFAvnpq4kufhueAY1bwbzVhk336o
    ./scripts/build.sh ~/.avalanchego/plugins/meq3bv7qCMZZ69L8xZRLwyKnWp6chRwyscq8VPtHWignRQVVF
    ```

1. Get upgrades for each L1 and write them out to `~/.avalanchego/configs/chains/<blockchain-id>/upgrade.json`:

    ```bash
    mkdir -p ~/.avalanchego/configs/chains/2D8RG4UpSXbPbvPCAWppNJyqTG2i2CAXSkTgmTBBvs7GKNZjsY
    curl -X POST --silent --header 'Content-Type: application/json' --data '{
        "jsonrpc": "2.0",
        "method": "eth_getChainConfig",
        "params": [],
        "id": 1
    }' https://subnets.avax.network/dispatch/testnet/rpc | \
    jq -r '.result.upgrades' > ~/.avalanchego/configs/chains/2D8RG4UpSXbPbvPCAWppNJyqTG2i2CAXSkTgmTBBvs7GKNZjsY/upgrade.json
    ```

    Note it's possible there is no upgrades so the upgrade.json might just be `{}`.

    ```bash
    mkdir -p ~/.avalanchego/configs/chains/98qnjenm7MBd8G2cPZoRvZrgJC33JGSAAKghsQ6eojbLCeRNp
    curl -X POST --silent --header 'Content-Type: application/json' --data '{
        "jsonrpc": "2.0",
        "method": "eth_getChainConfig",
        "params": [],
        "id": 1
    }' https://subnets.avax.network/echo/testnet/rpc | \
    jq -r '.result.upgrades' > ~/.avalanchego/configs/chains/98qnjenm7MBd8G2cPZoRvZrgJC33JGSAAKghsQ6eojbLCeRNp/upgrade.json
    ```

1. (Optional) You can tweak the `config.json` for each L1 if you want to test a particular feature for example.
    - Dispatch: `~/.avalanchego/configs/chains/2D8RG4UpSXbPbvPCAWppNJyqTG2i2CAXSkTgmTBBvs7GKNZjsY/config.json`
    - Echo: `~/.avalanchego/configs/chains/98qnjenm7MBd8G2cPZoRvZrgJC33JGSAAKghsQ6eojbLCeRNp/config.json`
1. (Optional) If you want to reboostrap completely the chain, you can remove `~/.avalanchego/chainData/<blockchain-id>/db/pebbledb`, for example:
    - Dispatch: `rm -r ~/.avalanchego/chainData/2D8RG4UpSXbPbvPCAWppNJyqTG2i2CAXSkTgmTBBvs7GKNZjsY/db/pebbledb`
    - Echo: `rm -r ~/.avalanchego/chainData/98qnjenm7MBd8G2cPZoRvZrgJC33JGSAAKghsQ6eojbLCeRNp/db/pebbledb`

    AvalancheGo keeps its database in `~/.avalanchego/db/fuji/v1.4.5/*.ldb` which you should not delete.
1. Build AvalancheGo:

    ```bash
    ./scripts/build.sh
    ```

1. Run AvalancheGo tracking the Dispatch and Echo VM IDs:

    ```bash
    ./build/avalanchego --network-id=fuji --partial-sync-primary-network --public-ip=127.0.0.1 \
    --track-subnets=7WtoAMPhrmh5KosDUsFL9yTcvw7YSxiKHPpdfs4JsgW47oZT5,i9gFpZQHPLcGfZaQLiwFAStddQD7iTKBpFfurPFJsXm1CkTZK
    ```

1. Follow the logs and wait until you see line stating the health `check started passing`, for example:

    ```text
    [04-02|12:03:09.830] INFO health/worker.go:261 check started passing {"name": "health", "name": "network", "tags": ["application"]}
    [04-02|12:03:09.831] INFO health/worker.go:261 check started passing {"name": "health", "name": "validation", "tags": ["application"]}
    [04-02|12:03:09.831] INFO health/worker.go:261 check started passing {"name": "health", "name": "P", "tags": ["11111111111111111111111111111111LpoYY"]}
    ```

1. In another terminal, check you can obtain the current block number for both chains:

    - Dispatch:

        ```bash
        curl -X POST --silent --header 'Content-Type: application/json' --data '{
            "jsonrpc": "2.0",
            "method": "eth_blockNumber",
            "params": [],
            "id": 1
        }' localhost:9650/ext/bc/2D8RG4UpSXbPbvPCAWppNJyqTG2i2CAXSkTgmTBBvs7GKNZjsY/rpc
        ```

    - Echo:

        ```bash
        curl -X POST --silent --header 'Content-Type: application/json' --data '{
            "jsonrpc": "2.0",
            "method": "eth_blockNumber",
            "params": [],
            "id": 1
        }' localhost:9650/ext/bc/98qnjenm7MBd8G2cPZoRvZrgJC33JGSAAKghsQ6eojbLCeRNp/rpc
        ```

### Release

If a successful release candidate was created, you can now create a release.

Following the previous example in the [Release candidate section](#release-candidate), we will create a release `v0.7.3` indicated by the `$VERSION` variable.

1. Head to the last release candidate pull request, and **restore** the deleted branch at the bottom of the Github page.
1. Create a new release through the [Github web interface](https://github.com/ava-labs/subnet-evm/releases/new)
    1. In the "Choose a tag" box, enter `$VERSION` (`v0.7.3`)
    1. In the "Target", pick the previously restored last release candidate branch `releases/${VERSION_RC}`, for example `releases/v0.7.3-rc.0`.
    Do not select `master` as the target branch to prevent adding new master branch commits to the release.
    1. Pick the previous release, for example as `v0.7.2` in our case, since the default would be the last release candidate.
    1. Set the "Release title" to `$VERSION` (`v0.7.3`)
    1. Set the description (breaking changes, features, fixes, documentation)
    1. Only tick the box "Set as the latest release"
    1. Click on the "Create release" button
1. Monitor the [release Github workflow](https://github.com/ava-labs/subnet-evm/actions/workflows/release.yml) to ensure the GoReleaser step succeeds and check the binaries are then published to [the releases page](https://github.com/ava-labs/subnet-evm/releases)
1. Monitor the [Publish Docker image workflow](https://github.com/ava-labs/subnet-evm/actions/workflows/publish_docker.yml) succeeds.
