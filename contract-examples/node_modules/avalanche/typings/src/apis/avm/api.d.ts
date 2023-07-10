/**
 * @packageDocumentation
 * @module API-AVM
 */
import BN from "bn.js";
import { Buffer } from "buffer/";
import AvalancheCore from "../../avalanche";
import { UTXOSet } from "./utxos";
import { KeyChain } from "./keychain";
import { Tx, UnsignedTx } from "./tx";
import { PayloadBase } from "../../utils/payload";
import { SECPMintOutput } from "./outputs";
import { InitialStates } from "./initialstates";
import { JRPCAPI } from "../../common/jrpcapi";
import { MinterSet } from "./minterset";
import { PersistanceOptions } from "../../utils/persistenceoptions";
import { OutputOwners } from "../../common/output";
import { SECPTransferOutput } from "./outputs";
import { GetUTXOsResponse, GetAssetDescriptionResponse, GetBalanceResponse, SendResponse, SendMultipleResponse, GetAddressTxsResponse, IMinterSet } from "./interfaces";
/**
 * Class for interacting with a node endpoint that is using the AVM.
 *
 * @category RPCAPIs
 *
 * @remarks This extends the [[JRPCAPI]] class. This class should not be directly called. Instead, use the [[Avalanche.addAPI]] function to register this interface with Avalanche.
 */
export declare class AVMAPI extends JRPCAPI {
    /**
     * @ignore
     */
    protected keychain: KeyChain;
    protected blockchainID: string;
    protected blockchainAlias: string;
    protected AVAXAssetID: Buffer;
    protected txFee: BN;
    protected creationTxFee: BN;
    protected mintTxFee: BN;
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
    setBlockchainAlias: (alias: string) => undefined;
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
     * Gets the default mint fee for this chain.
     *
     * @returns The default mint fee as a {@link https://github.com/indutny/bn.js/|BN}
     */
    getDefaultMintTxFee: () => BN;
    /**
     * Gets the mint fee for this chain.
     *
     * @returns The mint fee as a {@link https://github.com/indutny/bn.js/|BN}
     */
    getMintTxFee: () => BN;
    /**
     * Gets the creation fee for this chain.
     *
     * @returns The creation fee as a {@link https://github.com/indutny/bn.js/|BN}
     */
    getCreationTxFee: () => BN;
    /**
     * Sets the mint fee for this chain.
     *
     * @param fee The mint fee amount to set as {@link https://github.com/indutny/bn.js/|BN}
     */
    setMintTxFee: (fee: BN) => void;
    /**
     * Sets the creation fee for this chain.
     *
     * @param fee The creation fee amount to set as {@link https://github.com/indutny/bn.js/|BN}
     */
    setCreationTxFee: (fee: BN) => void;
    /**
     * Gets a reference to the keychain for this class.
     *
     * @returns The instance of [[KeyChain]] for this class
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
     * Gets the balance of a particular asset on a blockchain.
     *
     * @param address The address to pull the asset balance from
     * @param assetID The assetID to pull the balance from
     * @param includePartial If includePartial=false, returns only the balance held solely
     *
     * @returns Promise with the balance of the assetID as a {@link https://github.com/indutny/bn.js/|BN} on the provided address for the blockchain.
     */
    getBalance: (address: string, assetID: string, includePartial?: boolean) => Promise<GetBalanceResponse>;
    /**
     * Creates an address (and associated private keys) on a user on a blockchain.
     *
     * @param username Name of the user to create the address under
     * @param password Password to unlock the user and encrypt the private key
     *
     * @returns Promise for a string representing the address created by the vm.
     */
    createAddress: (username: string, password: string) => Promise<string>;
    /**
     * Create a new fixed-cap, fungible asset. A quantity of it is created at initialization and there no more is ever created.
     *
     * @param username The user paying the transaction fee (in $AVAX) for asset creation
     * @param password The password for the user paying the transaction fee (in $AVAX) for asset creation
     * @param name The human-readable name for the asset
     * @param symbol Optional. The shorthand symbol for the asset. Between 0 and 4 characters
     * @param denomination Optional. Determines how balances of this asset are displayed by user interfaces. Default is 0
     * @param initialHolders An array of objects containing the field "address" and "amount" to establish the genesis values for the new asset
     *
     * ```js
     * Example initialHolders:
     * [
     *   {
     *     "address": "X-avax1kj06lhgx84h39snsljcey3tpc046ze68mek3g5",
     *     "amount": 10000
     *   },
     *   {
     *     "address": "X-avax1am4w6hfrvmh3akduzkjthrtgtqafalce6an8cr",
     *     "amount": 50000
     *   }
     * ]
     * ```
     *
     * @returns Returns a Promise string containing the base 58 string representation of the ID of the newly created asset.
     */
    createFixedCapAsset: (username: string, password: string, name: string, symbol: string, denomination: number, initialHolders: object[]) => Promise<string>;
    /**
     * Create a new variable-cap, fungible asset. No units of the asset exist at initialization. Minters can mint units of this asset using createMintTx, signMintTx and sendMintTx.
     *
     * @param username The user paying the transaction fee (in $AVAX) for asset creation
     * @param password The password for the user paying the transaction fee (in $AVAX) for asset creation
     * @param name The human-readable name for the asset
     * @param symbol Optional. The shorthand symbol for the asset -- between 0 and 4 characters
     * @param denomination Optional. Determines how balances of this asset are displayed by user interfaces. Default is 0
     * @param minterSets is a list where each element specifies that threshold of the addresses in minters may together mint more of the asset by signing a minting transaction
     *
     * ```js
     * Example minterSets:
     * [
     *    {
     *      "minters":[
     *        "X-avax1am4w6hfrvmh3akduzkjthrtgtqafalce6an8cr"
     *      ],
     *      "threshold": 1
     *     },
     *     {
     *      "minters": [
     *        "X-avax1am4w6hfrvmh3akduzkjthrtgtqafalce6an8cr",
     *        "X-avax1kj06lhgx84h39snsljcey3tpc046ze68mek3g5",
     *        "X-avax1yell3e4nln0m39cfpdhgqprsd87jkh4qnakklx"
     *      ],
     *      "threshold": 2
     *     }
     * ]
     * ```
     *
     * @returns Returns a Promise string containing the base 58 string representation of the ID of the newly created asset.
     */
    createVariableCapAsset: (username: string, password: string, name: string, symbol: string, denomination: number, minterSets: object[]) => Promise<string>;
    /**
     * Creates a family of NFT Asset. No units of the asset exist at initialization. Minters can mint units of this asset using createMintTx, signMintTx and sendMintTx.
     *
     * @param username The user paying the transaction fee (in $AVAX) for asset creation
     * @param password The password for the user paying the transaction fee (in $AVAX) for asset creation
     * @param from Optional. An array of addresses managed by the node's keystore for this blockchain which will fund this transaction
     * @param changeAddr Optional. An address to send the change
     * @param name The human-readable name for the asset
     * @param symbol Optional. The shorthand symbol for the asset -- between 0 and 4 characters
     * @param minterSets is a list where each element specifies that threshold of the addresses in minters may together mint more of the asset by signing a minting transaction
     *
     * @returns Returns a Promise string containing the base 58 string representation of the ID of the newly created asset.
     */
    createNFTAsset: (username: string, password: string, from: string[] | Buffer[], changeAddr: string, name: string, symbol: string, minterSet: IMinterSet) => Promise<string>;
    /**
     * Create an unsigned transaction to mint more of an asset.
     *
     * @param amount The units of the asset to mint
     * @param assetID The ID of the asset to mint
     * @param to The address to assign the units of the minted asset
     * @param minters Addresses of the minters responsible for signing the transaction
     *
     * @returns Returns a Promise string containing the base 58 string representation of the unsigned transaction.
     */
    mint: (username: string, password: string, amount: number | BN, assetID: Buffer | string, to: string, minters: string[]) => Promise<string>;
    /**
     * Mint non-fungible tokens which were created with AVMAPI.createNFTAsset
     *
     * @param username The user paying the transaction fee (in $AVAX) for asset creation
     * @param password The password for the user paying the transaction fee (in $AVAX) for asset creation
     * @param from Optional. An array of addresses managed by the node's keystore for this blockchain which will fund this transaction
     * @param changeAddr Optional. An address to send the change
     * @param assetID The asset id which is being sent
     * @param to Address on X-Chain of the account to which this NFT is being sent
     * @param encoding Optional.  is the encoding format to use for the payload argument. Can be either "cb58" or "hex". Defaults to "hex".
     *
     * @returns ID of the transaction
     */
    mintNFT: (username: string, password: string, from: string[] | Buffer[], changeAddr: string, payload: string, assetID: string | Buffer, to: string, encoding?: string) => Promise<string>;
    /**
     * Send NFT from one account to another on X-Chain
     *
     * @param username The user paying the transaction fee (in $AVAX) for asset creation
     * @param password The password for the user paying the transaction fee (in $AVAX) for asset creation
     * @param from Optional. An array of addresses managed by the node's keystore for this blockchain which will fund this transaction
     * @param changeAddr Optional. An address to send the change
     * @param assetID The asset id which is being sent
     * @param groupID The group this NFT is issued to.
     * @param to Address on X-Chain of the account to which this NFT is being sent
     *
     * @returns ID of the transaction
     */
    sendNFT: (username: string, password: string, from: string[] | Buffer[], changeAddr: string, assetID: string | Buffer, groupID: number, to: string) => Promise<string>;
    /**
     * Exports the private key for an address.
     *
     * @param username The name of the user with the private key
     * @param password The password used to decrypt the private key
     * @param address The address whose private key should be exported
     *
     * @returns Promise with the decrypted private key as store in the database
     */
    exportKey: (username: string, password: string, address: string) => Promise<string>;
    /**
     * Imports a private key into the node's keystore under an user and for a blockchain.
     *
     * @param username The name of the user to store the private key
     * @param password The password that unlocks the user
     * @param privateKey A string representing the private key in the vm's format
     *
     * @returns The address for the imported private key.
     */
    importKey: (username: string, password: string, privateKey: string) => Promise<string>;
    /**
     * Send ANT (Avalanche Native Token) assets including AVAX from the X-Chain to an account on the P-Chain or C-Chain.
     *
     * After calling this method, you must call the P-Chain's `import` or the C-Chain’s `import` method to complete the transfer.
     *
     * @param username The Keystore user that controls the P-Chain or C-Chain account specified in `to`
     * @param password The password of the Keystore user
     * @param to The account on the P-Chain or C-Chain to send the asset to.
     * @param amount Amount of asset to export as a {@link https://github.com/indutny/bn.js/|BN}
     * @param assetID The asset id which is being sent
     *
     * @returns String representing the transaction id
     */
    export: (username: string, password: string, to: string, amount: BN, assetID: string) => Promise<string>;
    /**
     * Send ANT (Avalanche Native Token) assets including AVAX from an account on the P-Chain or C-Chain to an address on the X-Chain. This transaction
     * must be signed with the key of the account that the asset is sent from and which pays
     * the transaction fee.
     *
     * @param username The Keystore user that controls the account specified in `to`
     * @param password The password of the Keystore user
     * @param to The address of the account the asset is sent to.
     * @param sourceChain The chainID where the funds are coming from. Ex: "C"
     *
     * @returns Promise for a string for the transaction, which should be sent to the network
     * by calling issueTx.
     */
    import: (username: string, password: string, to: string, sourceChain: string) => Promise<string>;
    /**
     * Lists all the addresses under a user.
     *
     * @param username The user to list addresses
     * @param password The password of the user to list the addresses
     *
     * @returns Promise of an array of address strings in the format specified by the blockchain.
     */
    listAddresses: (username: string, password: string) => Promise<string[]>;
    /**
     * Retrieves all assets for an address on a server and their associated balances.
     *
     * @param address The address to get a list of assets
     *
     * @returns Promise of an object mapping assetID strings with {@link https://github.com/indutny/bn.js/|BN} balance for the address on the blockchain.
     */
    getAllBalances: (address: string) => Promise<object[]>;
    /**
     * Retrieves an assets name and symbol.
     *
     * @param assetID Either a {@link https://github.com/feross/buffer|Buffer} or an b58 serialized string for the AssetID or its alias.
     *
     * @returns Returns a Promise object with keys "name" and "symbol".
     */
    getAssetDescription: (assetID: Buffer | string) => Promise<GetAssetDescriptionResponse>;
    /**
     * Returns the transaction data of a provided transaction ID by calling the node's `getTx` method.
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
     * @param txID The string representation of the transaction ID
     *
     * @returns Returns a Promise string containing the status retrieved from the node
     */
    getTxStatus: (txID: string) => Promise<string>;
    /**
     * Retrieves the UTXOs related to the addresses provided from the node's `getUTXOs` method.
     *
     * @param addresses An array of addresses as cb58 strings or addresses as {@link https://github.com/feross/buffer|Buffer}s
     * @param sourceChain A string for the chain to look for the UTXO's. Default is to use this chain, but if exported UTXOs exist from other chains, this can used to pull them instead.
     * @param limit Optional. Returns at most [limit] addresses. If [limit] == 0 or > [maxUTXOsToFetch], fetches up to [maxUTXOsToFetch].
     * @param startIndex Optional. [StartIndex] defines where to start fetching UTXOs (for pagination.)
     * UTXOs fetched are from addresses equal to or greater than [StartIndex.Address]
     * For address [StartIndex.Address], only UTXOs with IDs greater than [StartIndex.Utxo] will be returned.
     * @param persistOpts Options available to persist these UTXOs in local storage
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
     * Helper function which creates an unsigned transaction. For more granular control, you may create your own
     * [[UnsignedTx]] manually (with their corresponding [[TransferableInput]]s, [[TransferableOutput]]s, and [[TransferOperation]]s).
     *
     * @param utxoset A set of UTXOs that the transaction is built on
     * @param amount The amount of AssetID to be spent in its smallest denomination, represented as {@link https://github.com/indutny/bn.js/|BN}.
     * @param assetID The assetID of the value being sent
     * @param toAddresses The addresses to send the funds
     * @param fromAddresses The addresses being used to send the funds from the UTXOs provided
     * @param changeAddresses The addresses that can spend the change remaining from the spent UTXOs
     * @param memo Optional CB58 Buffer or String which contains arbitrary bytes, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     * @param locktime Optional. The locktime field created in the resulting outputs
     * @param threshold Optional. The number of signatures required to spend the funds in the resultant UTXO
     *
     * @returns An unsigned transaction ([[UnsignedTx]]) which contains a [[BaseTx]].
     *
     * @remarks
     * This helper exists because the endpoint API should be the primary point of entry for most functionality.
     */
    buildBaseTx: (utxoset: UTXOSet, amount: BN, assetID: Buffer | string, toAddresses: string[], fromAddresses: string[], changeAddresses: string[], memo?: PayloadBase | Buffer, asOf?: BN, locktime?: BN, threshold?: number) => Promise<UnsignedTx>;
    /**
     * Helper function which creates an unsigned NFT Transfer. For more granular control, you may create your own
     * [[UnsignedTx]] manually (with their corresponding [[TransferableInput]]s, [[TransferableOutput]]s, and [[TransferOperation]]s).
     *
     * @param utxoset  A set of UTXOs that the transaction is built on
     * @param toAddresses The addresses to send the NFT
     * @param fromAddresses The addresses being used to send the NFT from the utxoID provided
     * @param changeAddresses The addresses that can spend the change remaining from the spent UTXOs
     * @param utxoid A base58 utxoID or an array of base58 utxoIDs for the nfts this transaction is sending
     * @param memo Optional CB58 Buffer or String which contains arbitrary bytes, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     * @param locktime Optional. The locktime field created in the resulting outputs
     * @param threshold Optional. The number of signatures required to spend the funds in the resultant UTXO
     *
     * @returns An unsigned transaction ([[UnsignedTx]]) which contains a [[NFTTransferTx]].
     *
     * @remarks
     * This helper exists because the endpoint API should be the primary point of entry for most functionality.
     */
    buildNFTTransferTx: (utxoset: UTXOSet, toAddresses: string[], fromAddresses: string[], changeAddresses: string[], utxoid: string | string[], memo?: PayloadBase | Buffer, asOf?: BN, locktime?: BN, threshold?: number) => Promise<UnsignedTx>;
    /**
     * Helper function which creates an unsigned Import Tx. For more granular control, you may create your own
     * [[UnsignedTx]] manually (with their corresponding [[TransferableInput]]s, [[TransferableOutput]]s, and [[TransferOperation]]s).
     *
     * @param utxoset  A set of UTXOs that the transaction is built on
     * @param ownerAddresses The addresses being used to import
     * @param sourceChain The chainid for where the import is coming from
     * @param toAddresses The addresses to send the funds
     * @param fromAddresses The addresses being used to send the funds from the UTXOs provided
     * @param changeAddresses The addresses that can spend the change remaining from the spent UTXOs
     * @param memo Optional CB58 Buffer or String which contains arbitrary bytes, up to 256 bytes
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
     * @param memo Optional CB58 Buffer or String which contains arbitrary bytes, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     * @param locktime Optional. The locktime field created in the resulting outputs
     * @param threshold Optional. The number of signatures required to spend the funds in the resultant UTXO
     * @param assetID Optional. The assetID of the asset to send. Defaults to AVAX assetID.
     * Regardless of the asset which you"re exporting, all fees are paid in AVAX.
     *
     * @returns An unsigned transaction ([[UnsignedTx]]) which contains an [[ExportTx]].
     */
    buildExportTx: (utxoset: UTXOSet, amount: BN, destinationChain: Buffer | string, toAddresses: string[], fromAddresses: string[], changeAddresses?: string[], memo?: PayloadBase | Buffer, asOf?: BN, locktime?: BN, threshold?: number, assetID?: string) => Promise<UnsignedTx>;
    /**
     * Creates an unsigned transaction. For more granular control, you may create your own
     * [[UnsignedTx]] manually (with their corresponding [[TransferableInput]]s, [[TransferableOutput]]s, and [[TransferOperation]]s).
     *
     * @param utxoset A set of UTXOs that the transaction is built on
     * @param fromAddresses The addresses being used to send the funds from the UTXOs {@link https://github.com/feross/buffer|Buffer}
     * @param changeAddresses The addresses that can spend the change remaining from the spent UTXOs
     * @param initialState The [[InitialStates]] that represent the intial state of a created asset
     * @param name String for the descriptive name of the asset
     * @param symbol String for the ticker symbol of the asset
     * @param denomination Number for the denomination which is 10^D. D must be >= 0 and <= 32. Ex: $1 AVAX = 10^9 $nAVAX
     * @param mintOutputs Optional. Array of [[SECPMintOutput]]s to be included in the transaction. These outputs can be spent to mint more tokens.
     * @param memo Optional CB58 Buffer or String which contains arbitrary bytes, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     *
     * @returns An unsigned transaction ([[UnsignedTx]]) which contains a [[CreateAssetTx]].
     *
     */
    buildCreateAssetTx: (utxoset: UTXOSet, fromAddresses: string[], changeAddresses: string[], initialStates: InitialStates, name: string, symbol: string, denomination: number, mintOutputs?: SECPMintOutput[], memo?: PayloadBase | Buffer, asOf?: BN) => Promise<UnsignedTx>;
    buildSECPMintTx: (utxoset: UTXOSet, mintOwner: SECPMintOutput, transferOwner: SECPTransferOutput, fromAddresses: string[], changeAddresses: string[], mintUTXOID: string, memo?: PayloadBase | Buffer, asOf?: BN) => Promise<any>;
    /**
     * Creates an unsigned transaction. For more granular control, you may create your own
     * [[UnsignedTx]] manually (with their corresponding [[TransferableInput]]s, [[TransferableOutput]]s, and [[TransferOperation]]s).
     *
     * @param utxoset A set of UTXOs that the transaction is built on
     * @param fromAddresses The addresses being used to send the funds from the UTXOs {@link https://github.com/feross/buffer|Buffer}
     * @param changeAddresses The addresses that can spend the change remaining from the spent UTXOs
     * @param minterSets is a list where each element specifies that threshold of the addresses in minters may together mint more of the asset by signing a minting transaction
     * @param name String for the descriptive name of the asset
     * @param symbol String for the ticker symbol of the asset
     * @param memo Optional CB58 Buffer or String which contains arbitrary bytes, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     * @param locktime Optional. The locktime field created in the resulting mint output
     *
     * ```js
     * Example minterSets:
     * [
     *      {
     *          "minters":[
     *              "X-avax1ghstjukrtw8935lryqtnh643xe9a94u3tc75c7"
     *          ],
     *          "threshold": 1
     *      },
     *      {
     *          "minters": [
     *              "X-avax1yell3e4nln0m39cfpdhgqprsd87jkh4qnakklx",
     *              "X-avax1k4nr26c80jaquzm9369j5a4shmwcjn0vmemcjz",
     *              "X-avax1ztkzsrjnkn0cek5ryvhqswdtcg23nhge3nnr5e"
     *          ],
     *          "threshold": 2
     *      }
     * ]
     * ```
     *
     * @returns An unsigned transaction ([[UnsignedTx]]) which contains a [[CreateAssetTx]].
     *
     */
    buildCreateNFTAssetTx: (utxoset: UTXOSet, fromAddresses: string[], changeAddresses: string[], minterSets: MinterSet[], name: string, symbol: string, memo?: PayloadBase | Buffer, asOf?: BN, locktime?: BN) => Promise<UnsignedTx>;
    /**
     * Creates an unsigned transaction. For more granular control, you may create your own
     * [[UnsignedTx]] manually (with their corresponding [[TransferableInput]]s, [[TransferableOutput]]s, and [[TransferOperation]]s).
     *
     * @param utxoset  A set of UTXOs that the transaction is built on
     * @param owners Either a single or an array of [[OutputOwners]] to send the nft output
     * @param fromAddresses The addresses being used to send the NFT from the utxoID provided
     * @param changeAddresses The addresses that can spend the change remaining from the spent UTXOs
     * @param utxoid A base58 utxoID or an array of base58 utxoIDs for the nft mint output this transaction is sending
     * @param groupID Optional. The group this NFT is issued to.
     * @param payload Optional. Data for NFT Payload as either a [[PayloadBase]] or a {@link https://github.com/feross/buffer|Buffer}
     * @param memo Optional CB58 Buffer or String which contains arbitrary bytes, up to 256 bytes
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     *
     * @returns An unsigned transaction ([[UnsignedTx]]) which contains an [[OperationTx]].
     *
     */
    buildCreateNFTMintTx: (utxoset: UTXOSet, owners: OutputOwners[] | OutputOwners, fromAddresses: string[], changeAddresses: string[], utxoid: string | string[], groupID?: number, payload?: PayloadBase | Buffer, memo?: PayloadBase | Buffer, asOf?: BN) => Promise<any>;
    /**
     * Helper function which takes an unsigned transaction and signs it, returning the resulting [[Tx]].
     *
     * @param utx The unsigned transaction of type [[UnsignedTx]]
     *
     * @returns A signed transaction of type [[Tx]]
     */
    signTx: (utx: UnsignedTx) => Tx;
    /**
     * Calls the node's issueTx method from the API and returns the resulting transaction ID as a string.
     *
     * @param tx A string, {@link https://github.com/feross/buffer|Buffer}, or [[Tx]] representing a transaction
     *
     * @returns A Promise string representing the transaction ID of the posted transaction.
     */
    issueTx: (tx: string | Buffer | Tx) => Promise<string>;
    /**
     * Calls the node's getAddressTxs method from the API and returns transactions corresponding to the provided address and assetID
     *
     * @param address The address for which we're fetching related transactions.
     * @param cursor Page number or offset.
     * @param pageSize  Number of items to return per page. Optional. Defaults to 1024. If [pageSize] == 0 or [pageSize] > [maxPageSize], then it fetches at max [maxPageSize] transactions
     * @param assetID Only return transactions that changed the balance of this asset. Must be an ID or an alias for an asset.
     *
     * @returns A promise object representing the array of transaction IDs and page offset
     */
    getAddressTxs: (address: string, cursor: number, pageSize: number | undefined, assetID: string | Buffer) => Promise<GetAddressTxsResponse>;
    /**
     * Sends an amount of assetID to the specified address from a list of owned of addresses.
     *
     * @param username The user that owns the private keys associated with the `from` addresses
     * @param password The password unlocking the user
     * @param assetID The assetID of the asset to send
     * @param amount The amount of the asset to be sent
     * @param to The address of the recipient
     * @param from Optional. An array of addresses managed by the node's keystore for this blockchain which will fund this transaction
     * @param changeAddr Optional. An address to send the change
     * @param memo Optional. CB58 Buffer or String which contains arbitrary bytes, up to 256 bytes
     *
     * @returns Promise for the string representing the transaction's ID.
     */
    send: (username: string, password: string, assetID: string | Buffer, amount: number | BN, to: string, from?: string[] | Buffer[], changeAddr?: string, memo?: string | Buffer) => Promise<SendResponse>;
    /**
     * Sends an amount of assetID to an array of specified addresses from a list of owned of addresses.
     *
     * @param username The user that owns the private keys associated with the `from` addresses
     * @param password The password unlocking the user
     * @param sendOutputs The array of SendOutputs. A SendOutput is an object literal which contains an assetID, amount, and to.
     * @param from Optional. An array of addresses managed by the node's keystore for this blockchain which will fund this transaction
     * @param changeAddr Optional. An address to send the change
     * @param memo Optional. CB58 Buffer or String which contains arbitrary bytes, up to 256 bytes
     *
     * @returns Promise for the string representing the transaction"s ID.
     */
    sendMultiple: (username: string, password: string, sendOutputs: {
        assetID: string | Buffer;
        amount: number | BN;
        to: string;
    }[], from?: string[] | Buffer[], changeAddr?: string, memo?: string | Buffer) => Promise<SendMultipleResponse>;
    /**
     * Given a JSON representation of this Virtual Machine’s genesis state, create the byte representation of that state.
     *
     * @param genesisData The blockchain's genesis data object
     *
     * @returns Promise of a string of bytes
     */
    buildGenesis: (genesisData: object) => Promise<string>;
    /**
     * @ignore
     */
    protected _cleanAddressArray(addresses: string[] | Buffer[], caller: string): string[];
    /**
     * This class should not be instantiated directly. Instead use the [[Avalanche.addAP`${I}`]] method.
     *
     * @param core A reference to the Avalanche class
     * @param baseURL Defaults to the string "/ext/bc/X" as the path to blockchain's baseURL
     * @param blockchainID The Blockchain"s ID. Defaults to an empty string: ""
     */
    constructor(core: AvalancheCore, baseURL?: string, blockchainID?: string);
}
//# sourceMappingURL=api.d.ts.map