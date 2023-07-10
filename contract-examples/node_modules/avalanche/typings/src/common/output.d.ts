/**
 * @packageDocumentation
 * @module Common-Output
 */
import { Buffer } from "buffer/";
import BN from "bn.js";
import { NBytes } from "./nbytes";
import { Serializable, SerializedEncoding } from "../utils/serialization";
/**
 * Class for representing an address used in [[Output]] types
 */
export declare class Address extends NBytes {
    protected _typeName: string;
    protected _typeID: any;
    protected bytes: Buffer;
    protected bsize: number;
    /**
     * Returns a function used to sort an array of [[Address]]es
     */
    static comparator: () => (a: Address, b: Address) => 1 | -1 | 0;
    /**
     * Returns a base-58 representation of the [[Address]].
     */
    toString(): string;
    /**
     * Takes a base-58 string containing an [[Address]], parses it, populates the class, and returns the length of the Address in bytes.
     *
     * @param bytes A base-58 string containing a raw [[Address]]
     *
     * @returns The length of the raw [[Address]]
     */
    fromString(addr: string): number;
    clone(): this;
    create(...args: any[]): this;
    /**
     * Class for representing an address used in [[Output]] types
     */
    constructor();
}
/**
 * Defines the most basic values for output ownership. Mostly inherited from, but can be used in population of NFT Owner data.
 */
export declare class OutputOwners extends Serializable {
    protected _typeName: string;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected locktime: Buffer;
    protected threshold: Buffer;
    protected numaddrs: Buffer;
    protected addresses: Address[];
    /**
     * Returns the threshold of signers required to spend this output.
     */
    getThreshold: () => number;
    /**
     * Returns the a {@link https://github.com/indutny/bn.js/|BN} repersenting the UNIX Timestamp when the lock is made available.
     */
    getLocktime: () => BN;
    /**
     * Returns an array of {@link https://github.com/feross/buffer|Buffer}s for the addresses.
     */
    getAddresses: () => Buffer[];
    /**
     * Returns the index of the address.
     *
     * @param address A {@link https://github.com/feross/buffer|Buffer} of the address to look up to return its index.
     *
     * @returns The index of the address.
     */
    getAddressIdx: (address: Buffer) => number;
    /**
     * Returns the address from the index provided.
     *
     * @param idx The index of the address.
     *
     * @returns Returns the string representing the address.
     */
    getAddress: (idx: number) => Buffer;
    /**
     * Given an array of address {@link https://github.com/feross/buffer|Buffer}s and an optional timestamp, returns true if the addresses meet the threshold required to spend the output.
     */
    meetsThreshold: (addresses: Buffer[], asOf?: BN) => boolean;
    /**
     * Given an array of addresses and an optional timestamp, select an array of address {@link https://github.com/feross/buffer|Buffer}s of qualified spenders for the output.
     */
    getSpenders: (addresses: Buffer[], asOf?: BN) => Buffer[];
    /**
     * Returns a base-58 string representing the [[Output]].
     */
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns the buffer representing the [[Output]] instance.
     */
    toBuffer(): Buffer;
    /**
     * Returns a base-58 string representing the [[Output]].
     */
    toString(): string;
    static comparator: () => (a: Output, b: Output) => 1 | -1 | 0;
    /**
     * An [[Output]] class which contains addresses, locktimes, and thresholds.
     *
     * @param addresses An array of {@link https://github.com/feross/buffer|Buffer}s representing output owner's addresses
     * @param locktime A {@link https://github.com/indutny/bn.js/|BN} representing the locktime
     * @param threshold A number representing the the threshold number of signers required to sign the transaction
     */
    constructor(addresses?: Buffer[], locktime?: BN, threshold?: number);
}
export declare abstract class Output extends OutputOwners {
    protected _typeName: string;
    protected _typeID: any;
    /**
     * Returns the outputID for the output which tells parsers what type it is
     */
    abstract getOutputID(): number;
    abstract clone(): this;
    abstract create(...args: any[]): this;
    abstract select(id: number, ...args: any[]): Output;
    /**
     *
     * @param assetID An assetID which is wrapped around the Buffer of the Output
     *
     * Must be implemented to use the appropriate TransferableOutput for the VM.
     */
    abstract makeTransferable(assetID: Buffer): StandardTransferableOutput;
}
export declare abstract class StandardParseableOutput extends Serializable {
    protected _typeName: string;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    protected output: Output;
    /**
     * Returns a function used to sort an array of [[ParseableOutput]]s
     */
    static comparator: () => (a: StandardParseableOutput, b: StandardParseableOutput) => 1 | -1 | 0;
    getOutput: () => Output;
    abstract fromBuffer(bytes: Buffer, offset?: number): number;
    toBuffer(): Buffer;
    /**
     * Class representing an [[ParseableOutput]] for a transaction.
     *
     * @param output A number representing the InputID of the [[ParseableOutput]]
     */
    constructor(output?: Output);
}
export declare abstract class StandardTransferableOutput extends StandardParseableOutput {
    protected _typeName: string;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected assetID: Buffer;
    getAssetID: () => Buffer;
    abstract fromBuffer(bytes: Buffer, offset?: number): number;
    toBuffer(): Buffer;
    /**
     * Class representing an [[StandardTransferableOutput]] for a transaction.
     *
     * @param assetID A {@link https://github.com/feross/buffer|Buffer} representing the assetID of the [[Output]]
     * @param output A number representing the InputID of the [[StandardTransferableOutput]]
     */
    constructor(assetID?: Buffer, output?: Output);
}
/**
 * An [[Output]] class which specifies a token amount .
 */
export declare abstract class StandardAmountOutput extends Output {
    protected _typeName: string;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected amount: Buffer;
    protected amountValue: BN;
    /**
     * Returns the amount as a {@link https://github.com/indutny/bn.js/|BN}.
     */
    getAmount(): BN;
    /**
     * Popuates the instance from a {@link https://github.com/feross/buffer|Buffer} representing the [[StandardAmountOutput]] and returns the size of the output.
     */
    fromBuffer(outbuff: Buffer, offset?: number): number;
    /**
     * Returns the buffer representing the [[StandardAmountOutput]] instance.
     */
    toBuffer(): Buffer;
    /**
     * A [[StandardAmountOutput]] class which issues a payment on an assetID.
     *
     * @param amount A {@link https://github.com/indutny/bn.js/|BN} representing the amount in the output
     * @param addresses An array of {@link https://github.com/feross/buffer|Buffer}s representing addresses
     * @param locktime A {@link https://github.com/indutny/bn.js/|BN} representing the locktime
     * @param threshold A number representing the the threshold number of signers required to sign the transaction
     */
    constructor(amount?: BN, addresses?: Buffer[], locktime?: BN, threshold?: number);
}
/**
 * An [[Output]] class which specifies an NFT.
 */
export declare abstract class BaseNFTOutput extends Output {
    protected _typeName: string;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected groupID: Buffer;
    /**
     * Returns the groupID as a number.
     */
    getGroupID: () => number;
}
//# sourceMappingURL=output.d.ts.map