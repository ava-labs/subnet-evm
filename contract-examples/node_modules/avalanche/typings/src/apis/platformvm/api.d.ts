/**
 * @packageDocumentation
 * @module API-PlatformVM
 */
import { Buffer } from "buffer/";
import BN from "bn.js";
import AvalancheCore from "../../avalanche";
import { JRPCAPI } from "../../common/jrpcapi";
import { ErrorResponseObject } from "../../utils/errors";
import { KeyChain } from "./keychain";
import { UnsignedTx, Tx } from "./tx";
import { PayloadBase } from "../../utils/payload";
import { UTXOSet } from "../platformvm/utxos";
import { PersistanceOptions } from "../../utils/persistenceoptions";
import { GetRewardUTXOsResponse, GetStakeResponse, Subnet, GetValidatorsAtResponse, GetBalanceResponse, GetUTXOsResponse, Blockchain, GetTxStatusResponse, GetMinStakeResponse } from "./interfaces";
import { GenesisData } from "../avm";
/**
 * Class for interacting with a node's PlatformVMAPI
 *
 * @category RPCAPIs
 *
 * @remarks This extends the [[JRPCAPI]] class. This class should not be directly called. Instead, use the [[Avalanche.addAPI]] function to register this interface with Avalanche.
 */
export declare class PlatformVMAPI extends JRPCAPI {
    /**
     * @ignore
     */
    protected keychain: KeyChain;
    protected blockchainID: string;
    protected blockchainAlias: string;
    protected AVAXAssetID: Buffer;
    protected txFee: BN;
    protected creationTxFee: BN;
    protected minValidatorStake: BN;
    protected minDelegatorStake: BN;
    /**
     * Gets the alias for the blockchainID if it exists, otherwise returns `undefined`.
     *
     * @returns The alias for the blockchainID
     */
    getBlockchainAlias: () => string;
    /**
     * Sets the alias for the blockchainID.
     *
     * @param alias The alias for the blockchainID.
     *
     */
    setBlockchainAlias: (alias: string) => string;
    /**
     * Gets the blockchainID and returns it.
     *
     * @returns The blockchainID
     */
    getBlockchainID: () => string;
    /**
     * Refresh blockchainID, and if a blockchainID is passed in, use that.
     *
     * @param Optional. BlockchainID to assign, if none, uses the default based on networkID.
     *
     * @returns The blockchainID
     */
    refreshBlockchainID: (blockchainID?: string) => boolean;
    /**
     * Takes an address string and returns its {@link https://github.com/feross/buffer|Buffer} representation if valid.
     *
     * @returns A {@link https://github.com/feross/buffer|Buffer} for the address if valid, undefined if not valid.
     */
    parseAddress: (addr: string) => Buffer;
    addressFromBuffer: (address: Buffer) => string;
    /**
     * Fetches the AVAX AssetID and returns it in a Promise.
     *
     * @param refresh This function caches the response. Refresh = true will bust the cache.
     *
     * @returns The the provided string representing the AVAX AssetID
     */
    getAVAXAssetID: (refresh?: boolean) => Promise<Buffer>;
    /**
     * Overrides the defaults and sets the cache to a specific AVAX AssetID
     *
     * @param avaxAssetID A cb58 string or Buffer representing the AVAX AssetID
     *
     * @returns The the provided string representing the AVAX AssetID
     */
    setAVAXAssetID: (avaxAssetID: string | Buffer) => void;
    /**
     * Gets the default tx fee for this chain.
     *
     * @returns The default tx fee as a {@link https://github.com/indutny/bn.js/|BN}
     */
    getDefaultTxFee: () => BN;
    /**
     * Gets the tx fee for this chain.
     *
     * @returns The tx fee as a {@link https://github.com/indutny/bn.js/|BN}
     */
    getTxFee: () => BN;
    /**
     * Gets the CreateSubnetTx fee.
     *
     * @returns The CreateSubnetTx fee as a {@link https://github.com/indutny/bn.js/|BN}
     */
    getCreateSubnetTxFee: () => BN;
    /**
     * Gets the CreateChainTx fee.
     *
     * @returns The CreateChainTx fee as a {@link https://github.com/indutny/bn.js/|BN}
     */
    getCreateChainTxFee: () => BN;
    /**
     * Sets the tx fee for this chain.
     *
     * @param fee The tx fee amount to set as {@link https://github.com/indutny/bn.js/|BN}
     */
    setTxFee: (fee: BN) => void;
    /**
     * Gets the default creation fee for this chain.
     *
     * @returns The default creation fee as a {@link https://github.com/indutny/bn.js/|BN}
     */
    getDefaultCreationTxFee: () => BN;
    /**
     * Gets the creation fee for this chain.
     *
     * @returns The creation fee as a {@link https://github.com/indutny/bn.js/|BN}
     */
    getCreationTxFee: () => BN;
    /**
     * Sets the creation fee for this chain.
     *
     * @param fee The creation fee amount to set as {@link https://github.com/indutny/bn.js/|BN}
     */
    setCreationTxFee: (fee: BN) => void;
    /**
     * Gets a reference to the keychain for this class.
     *
     * @returns The instance of [[]] for this class
     */
    keyChain: () => KeyChain;
    /**
     * @ignore
     */
    newKeyChain: () => KeyChain;
    /**
     * Helper function which determines if a tx is a goose egg transaction.
     *
     * @param utx An UnsignedTx
     *
     * @returns boolean true if passes goose egg test and false if fails.
     *
     * @remarks
     * A "Goose Egg Transaction" is when the fee far exceeds a reasonable amount
     */
    checkGooseEgg: (utx: UnsignedTx, outTotal?: BN) => Promise<boolean>;
    /**
     * Retrieves an assetID for a subnet"s staking assset.
     *
     * @returns Returns a Promise string with cb58 encoded value of the assetID.
     */
    getStakingAssetID: () => Promise<string>;
    /**
     * Creates a new blockchain.
     *
     * @param username The username of the Keystore user that controls the new account
     * @param password The password of the Keystore user that controls the new account
     * @param subnetID Optional. Either a {@link https://github.com/feross/buffer|Buffer} or an cb58 serialized string for the SubnetID or its alias.
     * @param vmID The ID of the Virtual Machine the blockchain runs. Can also be an alias of the Virtual Machine.
     * @param fxIDs The ids of the FXs the VM is running.
     * @param name A human-readable name for the new blockchain
     * @param genesis The base 58 (with checksum) representation of the genesis state of the new blockchain. Virtual Machines should have a static API method named buildGenesis that can be used to generate genesisData.
     *
     * @returns Promise for the unsigned transaction to create this blockchain. Must be signed by a sufficient number of the Subnet’s control keys and by the account paying the transaction fee.
     */
    createBlockchain: (username: string, password: string, subnetID: Buffer | string, vmID: string, fxIDs: number[], name: string, genesis: string) => Promise<string>;
    /**
     * Gets the status of a blockchain.
     *
     * @param blockchainID The blockchainID requesting a status update
     *
     * @returns Promise for a string of one of: "Validating", "Created", "Preferred", "Unknown".
     */
    getBlockchainStatus: (blockchainID: string) => Promise<string>;
    /**
     * Get the validators and their weights of a subnet or the Primary Network at a given P-Chain height.
     *
     * @param height The P-Chain height to get the validator set at.
     * @param subnetID Optional. A cb58 serialized string for the SubnetID or its alias.
     *
     * @returns Promise GetValidatorsAtResponse
     */
    getValidatorsAt: (height: number, subnetID?: string) => Promise<GetValidatorsAtResponse>;
    /**
     * Create an address in the node's keystore.
     *
     * @param username The username of the Keystore user that controls the new account
     * @param password The password of the Keystore user that controls the new account
     *
     * @returns Promise for a string of the newly created account address.
     */
    createAddress: (username: string, password: string) => Promise<string>;
    /**
     * Gets the balance of a particular asset.
     *
     * @param address The address to pull the asset balance from
     *
     * @returns Promise with the balance as a {@link https://github.com/indutny/bn.js/|BN} on the provided address.
     */
    getBalance: (address: string) => Promise<GetBalanceResponse>;
    /**
     * List the addresses controlled by the user.
     *
     * @param username The username of the Keystore user
     * @param password The password of the Keystore user
     *
     * @returns Promise for an array of addresses.
     */
    listAddresses: (username: string, password: string) => Promise<string[]>;
    /**
     * Lists the set of current validators.
     *
     * @param subnetID Optional. Either a {@link https://github.com/feross/buffer|Buffer} or an
     * cb58 serialized string for the SubnetID or its alias.
     * @param nodeIDs Optional. An array of strings
     *
     * @returns Promise for an array of validators that are currently staking, see: {@link https://docs.avax.network/v1.0/en/api/platform/#platformgetcurrentvalidators|platform.getCurrentValidators documentation}.
     *
     */
    getCurrentValidators: (subnetID?: Buffer | string, nodeIDs?: string[]) => Promise<object>;
    /**
     * Lists the set of pending validators.
     *
     * @param subnetID Optional. Either a {@link https://github.com/feross/buffer|Buffer}
     * or a cb58 serialized string for the SubnetID or its alias.
     * @param nodeIDs Optional. An array of strings
     *
     * @returns Promise for an array of validators that are pending staking, see: {@link https://docs.avax.network/v1.0/en/api/platform/#platformgetpendingvalidators|platform.getPendingValidators documentation}.
     *
     */
    getPendingValidators: (subnetID?: Buffer | string, nodeIDs?: string[]) => Promise<object>;
    /**
     * Samples `Size` validators from the current validator set.
     *
     * @param sampleSize Of the total universe of validators, select this many at random
     * @param subnetID Optional. Either a {@link https://github.com/feross/buffer|Buffer} or an
     * cb58 serialized string for the SubnetID or its alias.
     *
     * @returns Promise for an array of validator"s stakingIDs.
     */
    sampleValidators: (sampleSize: number, subnetID?: Buffer | string) => Promise<string[]>;
    /**
     * Add a validator to the Primary Network.
     *
     * @param username The username of the Keystore user
     * @param password The password of the Keystore user
     * @param nodeID The node ID of the validator
     * @param startTime Javascript Date object for the start time to validate
     * @param endTime Javascript Date object for the end time to validate
     * @param stakeAmount The amount of nAVAX the validator is staking as
     * a {@link https://github.com/indutny/bn.js/|BN}
     * @param rewardAddress The address the validator reward will go to, if there is one.
     * @param delegationFeeRate Optional. A {@link https://github.com/indutny/bn.js/|BN} for the percent fee this validator
     * charges when others delegate stake to them. Up to 4 decimal places allowed additional decimal places are ignored.
     * Must be between 0 and 100, inclusive. For example, if delegationFeeRate is 1.2345 and someone delegates to this
     * validator, then when the delegation period is over, 1.2345% of the reward goes to the validator and the rest goes
     * to the delegator.
     *
     * @returns Promise for a base58 string of the unsigned transaction.
     */
    addValidator: (username: string, password: string, nodeID: string, startTime: Date, endTime: Date, stakeAmount: BN, rewardAddress: string, delegationFeeRate?: BN) => Promise<string>;
    /**
     * Add a validator to a Subnet other than the Primary Network. The validator must validate the Primary Network for the entire duration they validate this Subnet.
     *
     * @param username The username of the Keystore user
     * @param password The password of the Keystore user
     * @param nodeID The node ID of the validator
     * @param subnetID Either a {@link https://github.com/feross/buffer|Buffer} or a cb58 serialized string for the SubnetID or its alias.
     * @param startTime Javascript Date object for the start time to validate
     * @param endTime Javascript Date object for the end time to validate
     * @param weight The validator’s weight used for sampling
     *
     * @returns Promise for the unsigned transaction. It must be signed (using sign) by the proper number of the Subnet’s control keys and by the key of the account paying the transaction fee before it can be issued.
     */
    addSubnetValidator: (username: string, password: string, nodeID: string, subnetID: Buffer | string, startTime: Date, endTime: Date, weight: number) => Promise<string>;
    /**
     * Add a delegator to the Primary Network.
     *
     * @param username The username of the Keystore user
     * @param password The password of the Keystore user
     * @param nodeID The node ID of the delegatee
     * @param startTime Javascript Date object for when the delegator starts delegating
     * @param endTime Javascript Date object for when the delegator starts delegating
     * @param stakeAmount The amount of nAVAX the delegator is staking as
     * a {@link https://github.com/indutny/bn.js/|BN}
     * @param rewardAddress The address of the account the staked AVAX and validation reward
     * (if applicable) are sent to at endTime
     *
     * @returns Promise for an array of validator"s stakingIDs.
     */
    addDelegator: (username: string, password: string, nodeID: string, startTime: Date, endTime: Date, stakeAmount: BN, rewardAddress: string) => Promise<string>;
    /**
     * Create an unsigned transaction to create a new Subnet. The unsigned transaction must be
     * signed with the key of the account paying the transaction fee. The Subnet’s ID is the ID of the transaction that creates it (ie the response from issueTx when issuing the signed transaction).
     *
     * @param username The username of the Keystore user
     * @param password The password of the Keystore user
     * @param controlKeys Array of platform addresses as strings
     * @param threshold To add a validator to this Subnet, a transaction must have threshold
     * signatures, where each signature is from a key whose address is an element of `controlKeys`
     *
     * @returns Promise for a string with the unsigned transaction encoded as base58.
     */
    createSubnet: (username: string, password: string, controlKeys: string[], threshold: number) => Promise<string | ErrorResponseObject>;
    /**
     * Get the Subnet that validates a given blockchain.
     *
     * @param blockchainID Either a {@link https://github.com/feross/buffer|Buffer} or a cb58
     * encoded string for the blockchainID or its alias.
     *
     * @returns Promise for a string of the subnetID that validates the blockchain.
     */
    validatedBy: (blockchainID: string) => Promise<string>;
    /**
     * Get the IDs of the blockchains a Subnet validates.
     *
     * @param subnetID Either a {@link https://github.com/feross/buffer|Buffer} or an AVAX
     * serialized string for the SubnetID or its alias.
     *
     * @returns Promise for an array of blockchainIDs the subnet validates.
     */
    validates: (subnetID: Buffer | string) => Promise<string[]>;
    /**
     * Get all the blockchains that exist (excluding the P-Chain).
     *
     * @returns Promise for an array of objects containing fields "id", "subnetID", and "vmID".
     */
    getBlockchains: () => Promise<Blockchain[]>;
    /**
     * Send AVAX from an account on the P-Chain to an address on the X-Chain. This transaction
     * must be signed with the key of the account that the AVAX is sent from and which pays the
     * transaction fee. After issuing this transaction, you must call the X-Chain’s importAVAX
     * method to complete the transfer.
     *
     * @param username The Keystore user that controls the account specified in `to`
     * @param password The password of the Keystore user
     * @param to The address on the X-Chain to send the AVAX to. Do not include X- in the address
     * @param amount Amount of AVAX to export as a {@link https://github.com/indutny/bn.js/|BN}
     *
     * @returns Promise for an unsigned transaction to be signed by the account the the AVAX is
     * sent from and pays the transaction fee.
     */
    exportAVAX: (username: string, password: string, amount: BN, to: string) => Promise<string | ErrorResponseObject>;
    /**
     * Send AVAX from an account on the P-Chain to an address on the X-Chain. This transaction
     * must be signed with the key of the account that the AVAX is sent from and which pays
     * the transaction fee. After issuing this transaction, you must call the X-Chain’s
     * importAVAX method to complete the transfer.
     *
     * @param username The Keystore user that controls the account specified in `to`
     * @param password The password of the Keystore user
     * @param to The ID of the account the AVAX is sent to. This must be the same as the to
     * argument in the corresponding call to the X-Chain’s exportAVAX
     * @param sourceChain The chainID where the funds are coming from.
     *
     * @returns Promise for a string for the transaction, which should be sent to the network
     * by calling issueTx.
     */
    importAVAX: (username: string, password: string, to: string, sourceChain: string) => Promise<string | ErrorResponseObject>;
    /**
     * Calls the node's issueTx method from the API and returns the resulting transaction ID as a string.
     *
     * @param tx A string, {@link https://github.com/feross/buffer|Buffer}, or [[Tx]] representing a transaction
     *
     * @returns A Promise string representing the transaction ID of the posted transaction.
     */
    issueTx: (tx: string | Buffer | Tx) => Promise<string>;
    /**
     * Returns an upper bound on the amount of tokens that exist. Not monotonically increasing because this number can go down if a staker"s reward is denied.
     */
    getCurrentSupply: () => Promise<BN>;
    /**
     * Returns the height of the platform chain.
     */
    getHeight: () => Promise<BN>;
    /**
     * Gets the minimum staking amount.
     *
     * @param refresh A boolean to bypass the local cached value of Minimum Stake Amount, polling the node instead.
     */
    getMinStake: (refresh?: boolean) => Promise<GetMinStakeResponse>;
    /**
     * getTotalStake() returns the total amount staked on the Primary Network
     *
     * @returns A big number representing total staked by validators on the primary network
     */
    getTotalStake: () => Promise<BN>;
    /**
     * getMaxStakeAmount() returns the maximum amount of nAVAX staking to the named node during the time period.
     *
     * @param subnetID A Buffer or cb58 string representing subnet
     * @param nodeID A string representing ID of the node whose stake amount is required during the given duration
     * @param startTime A big number denoting start time of the duration during which stake amount of the node is required.
     * @param endTime A big number denoting end time of the duration during which stake amount of the node is required.
     * @returns A big number representing total staked by validators on the primary network
     */
    getMaxStakeAmount: (subnetID: string | Buffer, nodeID: string, startTime: BN, endTime: BN) => Promise<BN>;
    /**
     * Sets the minimum stake cached in this class.
     * @param minValidatorStake A {@link https://github.com/indutny/bn.js/|BN} to set the minimum stake amount cached in this class.
     * @param minDelegatorStake A {@link https://github.com/indutny/bn.js/|BN} to set the minimum delegation amount cached in this class.
     */
    setMinStake: (minValidatorStake?: BN, minDelegatorStake?: BN) => void;
    /**
     * Gets the total amount staked for an array of addresses.
     */
    getStake: (addresses: string[], encoding?: string) => Promise<GetStakeResponse>;
    /**
     * Get all the subnets that exist.
     *
     * @param ids IDs of the subnets to retrieve information about. If omitted, gets all subnets
     *
     * @returns Promise for an array of objects containing fields "id",
     * "controlKeys", and "threshold".
     */
    getSubnets: (ids?: string[]) => Promise<Subnet[]>;
    /**
     * Exports the private key for an address.
     *
     * @param username The name of the user with the private key
     * @param password The password used to decrypt the private key
     * @param address The address whose private key should be exported
     *
     * @returns Promise with the decrypted private key as store in the database
     */
    exportKey: (username: string, password: string, address: string) => Promise<string | ErrorResponseObject>;
    /**
     * Give a user control over an address by providing the private key that controls the address.
     *
     * @param username The name of the user to store the private key
     * @param password The password that unlocks the user
     * @param privateKey A string representing the private key in the vm"s format
     *
     * @returns The address for the imported private key.
     */
    importKey: (username: string, password: string, privateKey: string) => Promise<string | ErrorResponseObject>;
    /**
     * Returns the treansaction data of a provided transaction ID by calling the node's `getTx` method.
     *
     * @param txID The string representation of the transaction ID
     * @param encoding sets the format of the returned transaction. Can be, "cb58", "hex" or "json". Defaults to "cb58".
     *
     * @returns Returns a Promise string or object containing the bytes retrieved from the node
     */
    getTx: (txID: string, encoding?: string) => Promise<string | object>;
    /**
     * Returns the status of a provided transaction ID by calling the node's `getTxStatus` method.
     *
     * @param txid The string representation of the transaction ID
     * @param includeReason Return the reason tx was dropped, if applicable. Defaults to true
     *
     * @returns Returns a Promise string containing the status retrieved from the node and the reason a tx was dropped, if applicable.
     */
    getTxStatus: (txid: string, includeReason?: boolean) => Promise<string | GetTxStatusResponse>;
    /**
     * Retrieves the UTXOs related to the addresses provided from the node's `getUTXOs` method.
     *
     * @param addresses An array of addresses as cb58 strings or addresses as {@link https://github.com/feross/buffer|Buffer}s
     * @param sourceChain A string for the chain to look for the UTXO"s. Default is to use this chain, but if exported UTXOs exist from other chains, this can used to pull them instead.
     * @param limit Optional. Returns at most [limit] addresses. If [limit] == 0 or > [maxUTXOsToFetch], fetches up to [maxUTXOsToFetch].
     * @param startIndex Optional. [StartIndex] defines where to start fetching UTXOs (for pagination.)
     * UTXOs fetched are from addresses equal to or greater than [StartIndex.Address]
     * For address [StartIndex.Address], only UTXOs with IDs greater than [StartIndex.Utxo] will be returned.
     * @param persistOpts Options available to persist these UTXOs in local storage
     * @param encoding Optional.  is the encoding format to use for the payload argument. Can be either "cb58" or "hex". Defaults to "hex".
     *
     * @remarks
     * persistOpts is optional and must be of type [[PersistanceOptions]]
     *
     */
    getUTXOs: (addresses: string[] | string, sourceChain?: string, limit?: number, startIndex?: {
        address: string;
        utxo: string;
    }, persistOpts?: PersistanceOptions, encoding?: string) => Promise<GetUTXOsResponse>;
    /**
     * Helper function which creates an unsigned Import Tx. For more granular control, you may create your own
     * [[UnsignedTx]] manually (with their corresponding [[TransferableInput]]s, [[TransferableOutput]]s, and [[TransferOperation]]s).
     *
     * @param utxoset A set of UTXOs that the transaction is built on
     * @param ownerAddresses The addresses being used to import
     * @param sourceChain The chainid for where the import is coming from.
     * @param toAddresses The addresses to send the funds
     * @param fromAddresses The addresses being used to send the funds from the UTXOs provided
     * @param changeAddresses The addresses that can spend the change remaining from the spent UTXOs
     * @param memo Optional contains arbitrary bytes, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     * @param locktime Optional. The locktime field created in the resulting outputs
     * @param threshold Optional. The number of signatures required to spend the funds in the resultant UTXO
     *
     * @returns An unsigned transaction ([[UnsignedTx]]) which contains a [[ImportTx]].
     *
     * @remarks
     * This helper exists because the endpoint API should be the primary point of entry for most functionality.
     */
    buildImportTx: (utxoset: UTXOSet, ownerAddresses: string[], sourceChain: Buffer | string, toAddresses: string[], fromAddresses: string[], changeAddresses?: string[], memo?: PayloadBase | Buffer, asOf?: BN, locktime?: BN, threshold?: number) => Promise<UnsignedTx>;
    /**
     * Helper function which creates an unsigned Export Tx. For more granular control, you may create your own
     * [[UnsignedTx]] manually (with their corresponding [[TransferableInput]]s, [[TransferableOutput]]s, and [[TransferOperation]]s).
     *
     * @param utxoset A set of UTXOs that the transaction is built on
     * @param amount The amount being exported as a {@link https://github.com/indutny/bn.js/|BN}
     * @param destinationChain The chainid for where the assets will be sent.
     * @param toAddresses The addresses to send the funds
     * @param fromAddresses The addresses being used to send the funds from the UTXOs provided
     * @param changeAddresses The addresses that can spend the change remaining from the spent UTXOs
     * @param memo Optional contains arbitrary bytes, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     * @param locktime Optional. The locktime field created in the resulting outputs
     * @param threshold Optional. The number of signatures required to spend the funds in the resultant UTXO
     *
     * @returns An unsigned transaction ([[UnsignedTx]]) which contains an [[ExportTx]].
     */
    buildExportTx: (utxoset: UTXOSet, amount: BN, destinationChain: Buffer | string, toAddresses: string[], fromAddresses: string[], changeAddresses?: string[], memo?: PayloadBase | Buffer, asOf?: BN, locktime?: BN, threshold?: number) => Promise<UnsignedTx>;
    /**
     * Helper function which creates an unsigned [[AddSubnetValidatorTx]]. For more granular control, you may create your own
     * [[UnsignedTx]] manually and import the [[AddSubnetValidatorTx]] class directly.
     *
     * @param utxoset A set of UTXOs that the transaction is built on.
     * @param fromAddresses An array of addresses as {@link https://github.com/feross/buffer|Buffer} who pays the fees in AVAX
     * @param changeAddresses An array of addresses as {@link https://github.com/feross/buffer|Buffer} who gets the change leftover from the fee payment
     * @param nodeID The node ID of the validator being added.
     * @param startTime The Unix time when the validator starts validating the Primary Network.
     * @param endTime The Unix time when the validator stops validating the Primary Network (and staked AVAX is returned).
     * @param weight The amount of weight for this subnet validator.
     * @param memo Optional contains arbitrary bytes, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     * @param subnetAuthCredentials Optional. An array of index and address to sign for each SubnetAuth.
     *
     * @returns An unsigned transaction created from the passed in parameters.
     */
    buildAddSubnetValidatorTx: (utxoset: UTXOSet, fromAddresses: string[], changeAddresses: string[], nodeID: string, startTime: BN, endTime: BN, weight: BN, subnetID: string, memo?: PayloadBase | Buffer, asOf?: BN, subnetAuthCredentials?: [number, Buffer][]) => Promise<UnsignedTx>;
    /**
     * Helper function which creates an unsigned [[AddDelegatorTx]]. For more granular control, you may create your own
     * [[UnsignedTx]] manually and import the [[AddDelegatorTx]] class directly.
     *
     * @param utxoset A set of UTXOs that the transaction is built on
     * @param toAddresses An array of addresses as {@link https://github.com/feross/buffer|Buffer} who received the staked tokens at the end of the staking period
     * @param fromAddresses An array of addresses as {@link https://github.com/feross/buffer|Buffer} who own the staking UTXOs the fees in AVAX
     * @param changeAddresses An array of addresses as {@link https://github.com/feross/buffer|Buffer} who gets the change leftover from the fee payment
     * @param nodeID The node ID of the validator being added.
     * @param startTime The Unix time when the validator starts validating the Primary Network.
     * @param endTime The Unix time when the validator stops validating the Primary Network (and staked AVAX is returned).
     * @param stakeAmount The amount being delegated as a {@link https://github.com/indutny/bn.js/|BN}
     * @param rewardAddresses The addresses which will recieve the rewards from the delegated stake.
     * @param rewardLocktime Optional. The locktime field created in the resulting reward outputs
     * @param rewardThreshold Opional. The number of signatures required to spend the funds in the resultant reward UTXO. Default 1.
     * @param memo Optional contains arbitrary bytes, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     *
     * @returns An unsigned transaction created from the passed in parameters.
     */
    buildAddDelegatorTx: (utxoset: UTXOSet, toAddresses: string[], fromAddresses: string[], changeAddresses: string[], nodeID: string, startTime: BN, endTime: BN, stakeAmount: BN, rewardAddresses: string[], rewardLocktime?: BN, rewardThreshold?: number, memo?: PayloadBase | Buffer, asOf?: BN) => Promise<UnsignedTx>;
    /**
     * Helper function which creates an unsigned [[AddValidatorTx]]. For more granular control, you may create your own
     * [[UnsignedTx]] manually and import the [[AddValidatorTx]] class directly.
     *
     * @param utxoset A set of UTXOs that the transaction is built on
     * @param toAddresses An array of addresses as {@link https://github.com/feross/buffer|Buffer} who received the staked tokens at the end of the staking period
     * @param fromAddresses An array of addresses as {@link https://github.com/feross/buffer|Buffer} who own the staking UTXOs the fees in AVAX
     * @param changeAddresses An array of addresses as {@link https://github.com/feross/buffer|Buffer} who gets the change leftover from the fee payment
     * @param nodeID The node ID of the validator being added.
     * @param startTime The Unix time when the validator starts validating the Primary Network.
     * @param endTime The Unix time when the validator stops validating the Primary Network (and staked AVAX is returned).
     * @param stakeAmount The amount being delegated as a {@link https://github.com/indutny/bn.js/|BN}
     * @param rewardAddresses The addresses which will recieve the rewards from the delegated stake.
     * @param delegationFee A number for the percentage of reward to be given to the validator when someone delegates to them. Must be between 0 and 100.
     * @param rewardLocktime Optional. The locktime field created in the resulting reward outputs
     * @param rewardThreshold Opional. The number of signatures required to spend the funds in the resultant reward UTXO. Default 1.
     * @param memo Optional contains arbitrary bytes, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     *
     * @returns An unsigned transaction created from the passed in parameters.
     */
    buildAddValidatorTx: (utxoset: UTXOSet, toAddresses: string[], fromAddresses: string[], changeAddresses: string[], nodeID: string, startTime: BN, endTime: BN, stakeAmount: BN, rewardAddresses: string[], delegationFee: number, rewardLocktime?: BN, rewardThreshold?: number, memo?: PayloadBase | Buffer, asOf?: BN) => Promise<UnsignedTx>;
    /**
     * Class representing an unsigned [[CreateSubnetTx]] transaction.
     *
     * @param utxoset A set of UTXOs that the transaction is built on
     * @param fromAddresses The addresses being used to send the funds from the UTXOs {@link https://github.com/feross/buffer|Buffer}
     * @param changeAddresses The addresses that can spend the change remaining from the spent UTXOs
     * @param subnetOwnerAddresses An array of addresses for owners of the new subnet
     * @param subnetOwnerThreshold A number indicating the amount of signatures required to add validators to a subnet
     * @param memo Optional contains arbitrary bytes, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     *
     * @returns An unsigned transaction created from the passed in parameters.
     */
    buildCreateSubnetTx: (utxoset: UTXOSet, fromAddresses: string[], changeAddresses: string[], subnetOwnerAddresses: string[], subnetOwnerThreshold: number, memo?: PayloadBase | Buffer, asOf?: BN) => Promise<UnsignedTx>;
    /**
     * Build an unsigned [[CreateChainTx]].
     *
     * @param utxoset A set of UTXOs that the transaction is built on
     * @param fromAddresses The addresses being used to send the funds from the UTXOs {@link https://github.com/feross/buffer|Buffer}
     * @param changeAddresses The addresses that can spend the change remaining from the spent UTXOs
     * @param subnetID Optional ID of the Subnet that validates this blockchain
     * @param chainName Optional A human readable name for the chain; need not be unique
     * @param vmID Optional ID of the VM running on the new chain
     * @param fxIDs Optional IDs of the feature extensions running on the new chain
     * @param genesisData Optional Byte representation of genesis state of the new chain
     * @param memo Optional contains arbitrary bytes, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     * @param subnetAuthCredentials Optional. An array of index and address to sign for each SubnetAuth.
     *
     * @returns An unsigned transaction created from the passed in parameters.
     */
    buildCreateChainTx: (utxoset: UTXOSet, fromAddresses: string[], changeAddresses: string[], subnetID?: string | Buffer, chainName?: string, vmID?: string, fxIDs?: string[], genesisData?: string | GenesisData, memo?: PayloadBase | Buffer, asOf?: BN, subnetAuthCredentials?: [number, Buffer][]) => Promise<UnsignedTx>;
    /**
     * @ignore
     */
    protected _cleanAddressArray(addresses: string[] | Buffer[], caller: string): string[];
    /**
     * This class should not be instantiated directly.
     * Instead use the [[Avalanche.addAPI]] method.
     *
     * @param core A reference to the Avalanche class
     * @param baseURL Defaults to the string "/ext/P" as the path to blockchain's baseURL
     */
    constructor(core: AvalancheCore, baseURL?: string);
    /**
     * @returns the current timestamp on chain.
     */
    getTimestamp: () => Promise<number>;
    /**
     * @returns the UTXOs that were rewarded after the provided transaction"s staking or delegation period ended.
     */
    getRewardUTXOs: (txID: string, encoding?: string) => Promise<GetRewardUTXOsResponse>;
}
//# sourceMappingURL=api.d.ts.map