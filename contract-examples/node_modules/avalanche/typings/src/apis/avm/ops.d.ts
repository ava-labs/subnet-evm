/**
 * @packageDocumentation
 * @module API-AVM-Operations
 */
import { Buffer } from "buffer/";
import { NFTTransferOutput, SECPMintOutput, SECPTransferOutput } from "./outputs";
import { NBytes } from "../../common/nbytes";
import { SigIdx } from "../../common/credentials";
import { OutputOwners } from "../../common/output";
import { Serializable, SerializedEncoding } from "../../utils/serialization";
/**
 * Takes a buffer representing the output and returns the proper [[Operation]] instance.
 *
 * @param opid A number representing the operation ID parsed prior to the bytes passed in
 *
 * @returns An instance of an [[Operation]]-extended class.
 */
export declare const SelectOperationClass: (opid: number, ...args: any[]) => Operation;
/**
 * A class representing an operation. All operation types must extend on this class.
 */
export declare abstract class Operation extends Serializable {
    protected _typeName: string;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected sigCount: Buffer;
    protected sigIdxs: SigIdx[];
    static comparator: () => (a: Operation, b: Operation) => 1 | -1 | 0;
    abstract getOperationID(): number;
    /**
     * Returns the array of [[SigIdx]] for this [[Operation]]
     */
    getSigIdxs: () => SigIdx[];
    /**
     * Returns the credential ID.
     */
    abstract getCredentialID(): number;
    /**
     * Creates and adds a [[SigIdx]] to the [[Operation]].
     *
     * @param addressIdx The index of the address to reference in the signatures
     * @param address The address of the source of the signature
     */
    addSignatureIdx: (addressIdx: number, address: Buffer) => void;
    fromBuffer(bytes: Buffer, offset?: number): number;
    toBuffer(): Buffer;
    /**
     * Returns a base-58 string representing the [[NFTMintOperation]].
     */
    toString(): string;
}
/**
 * A class which contains an [[Operation]] for transfers.
 *
 */
export declare class TransferableOperation extends Serializable {
    protected _typeName: string;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected assetID: Buffer;
    protected utxoIDs: UTXOID[];
    protected operation: Operation;
    /**
     * Returns a function used to sort an array of [[TransferableOperation]]s
     */
    static comparator: () => (a: TransferableOperation, b: TransferableOperation) => 1 | -1 | 0;
    /**
     * Returns the assetID as a {@link https://github.com/feross/buffer|Buffer}.
     */
    getAssetID: () => Buffer;
    /**
     * Returns an array of UTXOIDs in this operation.
     */
    getUTXOIDs: () => UTXOID[];
    /**
     * Returns the operation
     */
    getOperation: () => Operation;
    fromBuffer(bytes: Buffer, offset?: number): number;
    toBuffer(): Buffer;
    constructor(assetID?: Buffer, utxoids?: UTXOID[] | string[] | Buffer[], operation?: Operation);
}
/**
 * An [[Operation]] class which specifies a SECP256k1 Mint Op.
 */
export declare class SECPMintOperation extends Operation {
    protected _typeName: string;
    protected _codecID: number;
    protected _typeID: number;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected mintOutput: SECPMintOutput;
    protected transferOutput: SECPTransferOutput;
    /**
     * Set the codecID
     *
     * @param codecID The codecID to set
     */
    setCodecID(codecID: number): void;
    /**
     * Returns the operation ID.
     */
    getOperationID(): number;
    /**
     * Returns the credential ID.
     */
    getCredentialID(): number;
    /**
     * Returns the [[SECPMintOutput]] to be produced by this operation.
     */
    getMintOutput(): SECPMintOutput;
    /**
     * Returns [[SECPTransferOutput]] to be produced by this operation.
     */
    getTransferOutput(): SECPTransferOutput;
    /**
     * Popuates the instance from a {@link https://github.com/feross/buffer|Buffer} representing the [[SECPMintOperation]] and returns the updated offset.
     */
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns the buffer representing the [[SECPMintOperation]] instance.
     */
    toBuffer(): Buffer;
    /**
     * An [[Operation]] class which mints new tokens on an assetID.
     *
     * @param mintOutput The [[SECPMintOutput]] that will be produced by this transaction.
     * @param transferOutput A [[SECPTransferOutput]] that will be produced from this minting operation.
     */
    constructor(mintOutput?: SECPMintOutput, transferOutput?: SECPTransferOutput);
}
/**
 * An [[Operation]] class which specifies a NFT Mint Op.
 */
export declare class NFTMintOperation extends Operation {
    protected _typeName: string;
    protected _codecID: number;
    protected _typeID: number;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected groupID: Buffer;
    protected payload: Buffer;
    protected outputOwners: OutputOwners[];
    /**
     * Set the codecID
     *
     * @param codecID The codecID to set
     */
    setCodecID(codecID: number): void;
    /**
     * Returns the operation ID.
     */
    getOperationID(): number;
    /**
     * Returns the credential ID.
     */
    getCredentialID: () => number;
    /**
     * Returns the payload.
     */
    getGroupID: () => Buffer;
    /**
     * Returns the payload.
     */
    getPayload: () => Buffer;
    /**
     * Returns the payload's raw {@link https://github.com/feross/buffer|Buffer} with length prepended, for use with [[PayloadBase]]'s fromBuffer
     */
    getPayloadBuffer: () => Buffer;
    /**
     * Returns the outputOwners.
     */
    getOutputOwners: () => OutputOwners[];
    /**
     * Popuates the instance from a {@link https://github.com/feross/buffer|Buffer} representing the [[NFTMintOperation]] and returns the updated offset.
     */
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns the buffer representing the [[NFTMintOperation]] instance.
     */
    toBuffer(): Buffer;
    /**
     * Returns a base-58 string representing the [[NFTMintOperation]].
     */
    toString(): string;
    /**
     * An [[Operation]] class which contains an NFT on an assetID.
     *
     * @param groupID The group to which to issue the NFT Output
     * @param payload A {@link https://github.com/feross/buffer|Buffer} of the NFT payload
     * @param outputOwners An array of outputOwners
     */
    constructor(groupID?: number, payload?: Buffer, outputOwners?: OutputOwners[]);
}
/**
 * A [[Operation]] class which specifies a NFT Transfer Op.
 */
export declare class NFTTransferOperation extends Operation {
    protected _typeName: string;
    protected _codecID: number;
    protected _typeID: number;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected output: NFTTransferOutput;
    /**
     * Set the codecID
     *
     * @param codecID The codecID to set
     */
    setCodecID(codecID: number): void;
    /**
     * Returns the operation ID.
     */
    getOperationID(): number;
    /**
     * Returns the credential ID.
     */
    getCredentialID(): number;
    getOutput: () => NFTTransferOutput;
    /**
     * Popuates the instance from a {@link https://github.com/feross/buffer|Buffer} representing the [[NFTTransferOperation]] and returns the updated offset.
     */
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns the buffer representing the [[NFTTransferOperation]] instance.
     */
    toBuffer(): Buffer;
    /**
     * Returns a base-58 string representing the [[NFTTransferOperation]].
     */
    toString(): string;
    /**
     * An [[Operation]] class which contains an NFT on an assetID.
     *
     * @param output An [[NFTTransferOutput]]
     */
    constructor(output?: NFTTransferOutput);
}
/**
 * Class for representing a UTXOID used in [[TransferableOp]] types
 */
export declare class UTXOID extends NBytes {
    protected _typeName: string;
    protected _typeID: any;
    protected bytes: Buffer;
    protected bsize: number;
    /**
     * Returns a function used to sort an array of [[UTXOID]]s
     */
    static comparator: () => (a: UTXOID, b: UTXOID) => 1 | -1 | 0;
    /**
     * Returns a base-58 representation of the [[UTXOID]].
     */
    toString(): string;
    /**
     * Takes a base-58 string containing an [[UTXOID]], parses it, populates the class, and returns the length of the UTXOID in bytes.
     *
     * @param bytes A base-58 string containing a raw [[UTXOID]]
     *
     * @returns The length of the raw [[UTXOID]]
     */
    fromString(utxoid: string): number;
    clone(): this;
    create(...args: any[]): this;
    /**
     * Class for representing a UTXOID used in [[TransferableOp]] types
     */
    constructor();
}
//# sourceMappingURL=ops.d.ts.map