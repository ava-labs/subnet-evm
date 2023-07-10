/**
 * @packageDocumentation
 * @module Common-UTXOs
 */
import { Buffer } from "buffer/";
import BN from "bn.js";
import { Output } from "./output";
import { MergeRule } from "../utils/constants";
import { Serializable, SerializedEncoding } from "../utils/serialization";
/**
 * Class for representing a single StandardUTXO.
 */
export declare abstract class StandardUTXO extends Serializable {
    protected _typeName: string;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected codecID: Buffer;
    protected txid: Buffer;
    protected outputidx: Buffer;
    protected assetID: Buffer;
    protected output: Output;
    /**
     * Returns the numeric representation of the CodecID.
     */
    getCodecID: () => number;
    /**
     * Returns the {@link https://github.com/feross/buffer|Buffer} representation of the CodecID
     */
    getCodecIDBuffer: () => Buffer;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} of the TxID.
     */
    getTxID: () => Buffer;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer}  of the OutputIdx.
     */
    getOutputIdx: () => Buffer;
    /**
     * Returns the assetID as a {@link https://github.com/feross/buffer|Buffer}.
     */
    getAssetID: () => Buffer;
    /**
     * Returns the UTXOID as a base-58 string (UTXOID is a string )
     */
    getUTXOID: () => string;
    /**
     * Returns a reference to the output
     */
    getOutput: () => Output;
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} containing an [[StandardUTXO]], parses it, populates the class, and returns the length of the StandardUTXO in bytes.
     *
     * @param bytes A {@link https://github.com/feross/buffer|Buffer} containing a raw [[StandardUTXO]]
     */
    abstract fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[StandardUTXO]].
     */
    toBuffer(): Buffer;
    abstract fromString(serialized: string): number;
    abstract toString(): string;
    abstract clone(): this;
    abstract create(codecID?: number, txid?: Buffer, outputidx?: Buffer | number, assetID?: Buffer, output?: Output): this;
    /**
     * Class for representing a single StandardUTXO.
     *
     * @param codecID Optional number which specifies the codeID of the UTXO. Default 0
     * @param txID Optional {@link https://github.com/feross/buffer|Buffer} of transaction ID for the StandardUTXO
     * @param txidx Optional {@link https://github.com/feross/buffer|Buffer} or number for the index of the transaction's [[Output]]
     * @param assetID Optional {@link https://github.com/feross/buffer|Buffer} of the asset ID for the StandardUTXO
     * @param outputid Optional {@link https://github.com/feross/buffer|Buffer} or number of the output ID for the StandardUTXO
     */
    constructor(codecID?: number, txID?: Buffer, outputidx?: Buffer | number, assetID?: Buffer, output?: Output);
}
/**
 * Class representing a set of [[StandardUTXO]]s.
 */
export declare abstract class StandardUTXOSet<UTXOClass extends StandardUTXO> extends Serializable {
    protected _typeName: string;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    protected utxos: {
        [utxoid: string]: UTXOClass;
    };
    protected addressUTXOs: {
        [address: string]: {
            [utxoid: string]: BN;
        };
    };
    abstract parseUTXO(utxo: UTXOClass | string): UTXOClass;
    /**
     * Returns true if the [[StandardUTXO]] is in the StandardUTXOSet.
     *
     * @param utxo Either a [[StandardUTXO]] a cb58 serialized string representing a StandardUTXO
     */
    includes: (utxo: UTXOClass | string) => boolean;
    /**
     * Adds a [[StandardUTXO]] to the StandardUTXOSet.
     *
     * @param utxo Either a [[StandardUTXO]] an cb58 serialized string representing a StandardUTXO
     * @param overwrite If true, if the UTXOID already exists, overwrite it... default false
     *
     * @returns A [[StandardUTXO]] if one was added and undefined if nothing was added.
     */
    add(utxo: UTXOClass | string, overwrite?: boolean): UTXOClass;
    /**
     * Adds an array of [[StandardUTXO]]s to the [[StandardUTXOSet]].
     *
     * @param utxo Either a [[StandardUTXO]] an cb58 serialized string representing a StandardUTXO
     * @param overwrite If true, if the UTXOID already exists, overwrite it... default false
     *
     * @returns An array of StandardUTXOs which were added.
     */
    addArray(utxos: string[] | UTXOClass[], overwrite?: boolean): StandardUTXO[];
    /**
     * Removes a [[StandardUTXO]] from the [[StandardUTXOSet]] if it exists.
     *
     * @param utxo Either a [[StandardUTXO]] an cb58 serialized string representing a StandardUTXO
     *
     * @returns A [[StandardUTXO]] if it was removed and undefined if nothing was removed.
     */
    remove: (utxo: UTXOClass | string) => UTXOClass;
    /**
     * Removes an array of [[StandardUTXO]]s to the [[StandardUTXOSet]].
     *
     * @param utxo Either a [[StandardUTXO]] an cb58 serialized string representing a StandardUTXO
     * @param overwrite If true, if the UTXOID already exists, overwrite it... default false
     *
     * @returns An array of UTXOs which were removed.
     */
    removeArray: (utxos: string[] | UTXOClass[]) => UTXOClass[];
    /**
     * Gets a [[StandardUTXO]] from the [[StandardUTXOSet]] by its UTXOID.
     *
     * @param utxoid String representing the UTXOID
     *
     * @returns A [[StandardUTXO]] if it exists in the set.
     */
    getUTXO: (utxoid: string) => UTXOClass;
    /**
     * Gets all the [[StandardUTXO]]s, optionally that match with UTXOIDs in an array
     *
     * @param utxoids An optional array of UTXOIDs, returns all [[StandardUTXO]]s if not provided
     *
     * @returns An array of [[StandardUTXO]]s.
     */
    getAllUTXOs: (utxoids?: string[]) => UTXOClass[];
    /**
     * Gets all the [[StandardUTXO]]s as strings, optionally that match with UTXOIDs in an array.
     *
     * @param utxoids An optional array of UTXOIDs, returns all [[StandardUTXO]]s if not provided
     *
     * @returns An array of [[StandardUTXO]]s as cb58 serialized strings.
     */
    getAllUTXOStrings: (utxoids?: string[]) => string[];
    /**
     * Given an address or array of addresses, returns all the UTXOIDs for those addresses
     *
     * @param address An array of address {@link https://github.com/feross/buffer|Buffer}s
     * @param spendable If true, only retrieves UTXOIDs whose locktime has passed
     *
     * @returns An array of addresses.
     */
    getUTXOIDs: (addresses?: Buffer[], spendable?: boolean) => string[];
    /**
     * Gets the addresses in the [[StandardUTXOSet]] and returns an array of {@link https://github.com/feross/buffer|Buffer}.
     */
    getAddresses: () => Buffer[];
    /**
     * Returns the balance of a set of addresses in the StandardUTXOSet.
     *
     * @param addresses An array of addresses
     * @param assetID Either a {@link https://github.com/feross/buffer|Buffer} or an cb58 serialized representation of an AssetID
     * @param asOf The timestamp to verify the transaction against as a {@link https://github.com/indutny/bn.js/|BN}
     *
     * @returns Returns the total balance as a {@link https://github.com/indutny/bn.js/|BN}.
     */
    getBalance: (addresses: Buffer[], assetID: Buffer | string, asOf?: BN) => BN;
    /**
     * Gets all the Asset IDs, optionally that match with Asset IDs in an array
     *
     * @param utxoids An optional array of Addresses as string or Buffer, returns all Asset IDs if not provided
     *
     * @returns An array of {@link https://github.com/feross/buffer|Buffer} representing the Asset IDs.
     */
    getAssetIDs: (addresses?: Buffer[]) => Buffer[];
    abstract clone(): this;
    abstract create(...args: any[]): this;
    filter(args: any[], lambda: (utxo: UTXOClass, ...largs: any[]) => boolean): this;
    /**
     * Returns a new set with copy of UTXOs in this and set parameter.
     *
     * @param utxoset The [[StandardUTXOSet]] to merge with this one
     * @param hasUTXOIDs Will subselect a set of [[StandardUTXO]]s which have the UTXOIDs provided in this array, defults to all UTXOs
     *
     * @returns A new StandardUTXOSet that contains all the filtered elements.
     */
    merge: (utxoset: this, hasUTXOIDs?: string[]) => this;
    /**
     * Set intersetion between this set and a parameter.
     *
     * @param utxoset The set to intersect
     *
     * @returns A new StandardUTXOSet containing the intersection
     */
    intersection: (utxoset: this) => this;
    /**
     * Set difference between this set and a parameter.
     *
     * @param utxoset The set to difference
     *
     * @returns A new StandardUTXOSet containing the difference
     */
    difference: (utxoset: this) => this;
    /**
     * Set symmetrical difference between this set and a parameter.
     *
     * @param utxoset The set to symmetrical difference
     *
     * @returns A new StandardUTXOSet containing the symmetrical difference
     */
    symDifference: (utxoset: this) => this;
    /**
     * Set union between this set and a parameter.
     *
     * @param utxoset The set to union
     *
     * @returns A new StandardUTXOSet containing the union
     */
    union: (utxoset: this) => this;
    /**
     * Merges a set by the rule provided.
     *
     * @param utxoset The set to merge by the MergeRule
     * @param mergeRule The [[MergeRule]] to apply
     *
     * @returns A new StandardUTXOSet containing the merged data
     *
     * @remarks
     * The merge rules are as follows:
     *   * "intersection" - the intersection of the set
     *   * "differenceSelf" - the difference between the existing data and new set
     *   * "differenceNew" - the difference between the new data and the existing set
     *   * "symDifference" - the union of the differences between both sets of data
     *   * "union" - the unique set of all elements contained in both sets
     *   * "unionMinusNew" - the unique set of all elements contained in both sets, excluding values only found in the new set
     *   * "unionMinusSelf" - the unique set of all elements contained in both sets, excluding values only found in the existing set
     */
    mergeByRule: (utxoset: this, mergeRule: MergeRule) => this;
}
//# sourceMappingURL=utxos.d.ts.map