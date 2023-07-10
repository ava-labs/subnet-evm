/**
 * @packageDocumentation
 * @module API-EVM
 */
import { Buffer } from "buffer/";
import BN from "bn.js";
import AvalancheCore from "../../avalanche";
import { JRPCAPI } from "../../common/jrpcapi";
import { UTXOSet } from "./utxos";
import { KeyChain } from "./keychain";
import { Tx, UnsignedTx } from "./tx";
import { Index } from "./../../common/interfaces";
/**
 * Class for interacting with a node's EVMAPI
 *
 * @category RPCAPIs
 *
 * @remarks This extends the [[JRPCAPI]] class. This class should not be directly called. Instead, use the [[Avalanche.addAPI]] function to register this interface with Avalanche.
 */
export declare class EVMAPI extends JRPCAPI {
    /**
     * @ignore
     */
    protected keychain: KeyChain;
    protected blockchainID: string;
    protected blockchainAlias: string;
    protected AVAXAssetID: Buffer;
    protected txFee: BN;
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
     * @returns A boolean if the blockchainID was successfully refreshed.
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
     * Retrieves an assets name and symbol.
     *
     * @param assetID Either a {@link https://github.com/feross/buffer|Buffer} or an b58 serialized string for the AssetID or its alias.
     *
     * @returns Returns a Promise Asset with keys "name", "symbol", "assetID" and "denomination".
     */
    getAssetDescription: (assetID: Buffer | string) => Promise<any>;
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
     * returns the amount of [assetID] for the given address in the state of the given block number.
     * "latest", "pending", and "accepted" meta block numbers are also allowed.
     *
     * @param hexAddress The hex representation of the address
     * @param blockHeight The block height
     * @param assetID The asset ID
     *
     * @returns Returns a Promise object containing the balance
     */
    getAssetBalance: (hexAddress: string, blockHeight: string, assetID: string) => Promise<object>;
    /**
     * Returns the status of a provided atomic transaction ID by calling the node's `getAtomicTxStatus` method.
     *
     * @param txID The string representation of the transaction ID
     *
     * @returns Returns a Promise string containing the status retrieved from the node
     */
    getAtomicTxStatus: (txID: string) => Promise<string>;
    /**
     * Returns the transaction data of a provided transaction ID by calling the node's `getAtomicTx` method.
     *
     * @param txID The string representation of the transaction ID
     *
     * @returns Returns a Promise string containing the bytes retrieved from the node
     */
    getAtomicTx: (txID: string) => Promise<string>;
    /**
     * Gets the tx fee for this chain.
     *
     * @returns The tx fee as a {@link https://github.com/indutny/bn.js/|BN}
     */
    getTxFee: () => BN;
    /**
     * Send ANT (Avalanche Native Token) assets including AVAX from the C-Chain to an account on the X-Chain.
     *
     * After calling this method, you must call the X-Chain’s import method to complete the transfer.
     *
     * @param username The Keystore user that controls the X-Chain account specified in `to`
     * @param password The password of the Keystore user
     * @param to The account on the X-Chain to send the AVAX to.
     * @param amount Amount of asset to export as a {@link https://github.com/indutny/bn.js/|BN}
     * @param assetID The asset id which is being sent
     *
     * @returns String representing the transaction id
     */
    export: (username: string, password: string, to: string, amount: BN, assetID: string) => Promise<string>;
    /**
     * Send AVAX from the C-Chain to an account on the X-Chain.
     *
     * After calling this method, you must call the X-Chain’s importAVAX method to complete the transfer.
     *
     * @param username The Keystore user that controls the X-Chain account specified in `to`
     * @param password The password of the Keystore user
     * @param to The account on the X-Chain to send the AVAX to.
     * @param amount Amount of AVAX to export as a {@link https://github.com/indutny/bn.js/|BN}
     *
     * @returns String representing the transaction id
     */
    exportAVAX: (username: string, password: string, to: string, amount: BN) => Promise<string>;
    /**
     * Retrieves the UTXOs related to the addresses provided from the node's `getUTXOs` method.
     *
     * @param addresses An array of addresses as cb58 strings or addresses as {@link https://github.com/feross/buffer|Buffer}s
     * @param sourceChain A string for the chain to look for the UTXO's. Default is to use this chain, but if exported UTXOs exist
     * from other chains, this can used to pull them instead.
     * @param limit Optional. Returns at most [limit] addresses. If [limit] == 0 or > [maxUTXOsToFetch], fetches up to [maxUTXOsToFetch].
     * @param startIndex Optional. [StartIndex] defines where to start fetching UTXOs (for pagination.)
     * UTXOs fetched are from addresses equal to or greater than [StartIndex.Address]
     * For address [StartIndex.Address], only UTXOs with IDs greater than [StartIndex.Utxo] will be returned.
     */
    getUTXOs: (addresses: string[] | string, sourceChain?: string, limit?: number, startIndex?: Index, encoding?: string) => Promise<{
        numFetched: number;
        utxos: any;
        endIndex: Index;
    }>;
    /**
     * Send ANT (Avalanche Native Token) assets including AVAX from an account on the X-Chain to an address on the C-Chain. This transaction
     * must be signed with the key of the account that the asset is sent from and which pays
     * the transaction fee.
     *
     * @param username The Keystore user that controls the account specified in `to`
     * @param password The password of the Keystore user
     * @param to The address of the account the asset is sent to.
     * @param sourceChain The chainID where the funds are coming from. Ex: "X"
     *
     * @returns Promise for a string for the transaction, which should be sent to the network
     * by calling issueTx.
     */
    import: (username: string, password: string, to: string, sourceChain: string) => Promise<string>;
    /**
     * Send AVAX from an account on the X-Chain to an address on the C-Chain. This transaction
     * must be signed with the key of the account that the AVAX is sent from and which pays
     * the transaction fee.
     *
     * @param username The Keystore user that controls the account specified in `to`
     * @param password The password of the Keystore user
     * @param to The address of the account the AVAX is sent to. This must be the same as the to
     * argument in the corresponding call to the X-Chain’s exportAVAX
     * @param sourceChain The chainID where the funds are coming from.
     *
     * @returns Promise for a string for the transaction, which should be sent to the network
     * by calling issueTx.
     */
    importAVAX: (username: string, password: string, to: string, sourceChain: string) => Promise<string>;
    /**
     * Give a user control over an address by providing the private key that controls the address.
     *
     * @param username The name of the user to store the private key
     * @param password The password that unlocks the user
     * @param privateKey A string representing the private key in the vm"s format
     *
     * @returns The address for the imported private key.
     */
    importKey: (username: string, password: string, privateKey: string) => Promise<string>;
    /**
     * Calls the node's issueTx method from the API and returns the resulting transaction ID as a string.
     *
     * @param tx A string, {@link https://github.com/feross/buffer|Buffer}, or [[Tx]] representing a transaction
     *
     * @returns A Promise string representing the transaction ID of the posted transaction.
     */
    issueTx: (tx: string | Buffer | Tx) => Promise<string>;
    /**
     * Exports the private key for an address.
     *
     * @param username The name of the user with the private key
     * @param password The password used to decrypt the private key
     * @param address The address whose private key should be exported
     *
     * @returns Promise with the decrypted private key and private key hex as store in the database
     */
    exportKey: (username: string, password: string, address: string) => Promise<object>;
    /**
     * Helper function which creates an unsigned Import Tx. For more granular control, you may create your own
     * [[UnsignedTx]] manually (with their corresponding [[TransferableInput]]s, [[TransferableOutput]]s).
     *
     * @param utxoset A set of UTXOs that the transaction is built on
     * @param toAddress The address to send the funds
     * @param ownerAddresses The addresses being used to import
     * @param sourceChain The chainid for where the import is coming from
     * @param fromAddresses The addresses being used to send the funds from the UTXOs provided
     *
     * @returns An unsigned transaction ([[UnsignedTx]]) which contains a [[ImportTx]].
     *
     * @remarks
     * This helper exists because the endpoint API should be the primary point of entry for most functionality.
     */
    buildImportTx: (utxoset: UTXOSet, toAddress: string, ownerAddresses: string[], sourceChain: Buffer | string, fromAddresses: string[], fee?: BN) => Promise<UnsignedTx>;
    /**
     * Helper function which creates an unsigned Export Tx. For more granular control, you may create your own
     * [[UnsignedTx]] manually (with their corresponding [[TransferableInput]]s, [[TransferableOutput]]s).
     *
     * @param amount The amount being exported as a {@link https://github.com/indutny/bn.js/|BN}
     * @param assetID The asset id which is being sent
     * @param destinationChain The chainid for where the assets will be sent.
     * @param toAddresses The addresses to send the funds
     * @param fromAddresses The addresses being used to send the funds from the UTXOs provided
     * @param changeAddresses The addresses that can spend the change remaining from the spent UTXOs
     * @param asOf Optional. The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     * @param locktime Optional. The locktime field created in the resulting outputs
     * @param threshold Optional. The number of signatures required to spend the funds in the resultant UTXO
     *
     * @returns An unsigned transaction ([[UnsignedTx]]) which contains an [[ExportTx]].
     */
    buildExportTx: (amount: BN, assetID: Buffer | string, destinationChain: Buffer | string, fromAddressHex: string, fromAddressBech: string, toAddresses: string[], nonce?: number, locktime?: BN, threshold?: number, fee?: BN) => Promise<UnsignedTx>;
    /**
     * Gets a reference to the keychain for this class.
     *
     * @returns The instance of [[KeyChain]] for this class
     */
    keyChain: () => KeyChain;
    /**
     *
     * @returns new instance of [[KeyChain]]
     */
    newKeyChain: () => KeyChain;
    /**
     * @ignore
     */
    protected _cleanAddressArray(addresses: string[] | Buffer[], caller: string): string[];
    /**
     * This class should not be instantiated directly.
     * Instead use the [[Avalanche.addAPI]] method.
     *
     * @param core A reference to the Avalanche class
     * @param baseURL Defaults to the string "/ext/bc/C/avax" as the path to blockchain's baseURL
     * @param blockchainID The Blockchain's ID. Defaults to an empty string: ""
     */
    constructor(core: AvalancheCore, baseURL?: string, blockchainID?: string);
    /**
     * @returns a Promise string containing the base fee for the next block.
     */
    getBaseFee: () => Promise<string>;
    /**
     * returns the priority fee needed to be included in a block.
     *
     * @returns Returns a Promise string containing the priority fee needed to be included in a block.
     */
    getMaxPriorityFeePerGas: () => Promise<string>;
}
//# sourceMappingURL=api.d.ts.map