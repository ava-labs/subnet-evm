/**
 * @packageDocumentation
 * @module API-AVM-OperationTx
 */
import { Buffer } from "buffer/";
import { TransferableOutput } from "./outputs";
import { TransferableInput } from "./inputs";
import { TransferableOperation } from "./ops";
import { KeyChain } from "./keychain";
import { Credential } from "../../common/credentials";
import { BaseTx } from "./basetx";
import { SerializedEncoding } from "../../utils/serialization";
/**
 * Class representing an unsigned Operation transaction.
 */
export declare class OperationTx extends BaseTx {
    protected _typeName: string;
    protected _codecID: number;
    protected _typeID: number;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected numOps: Buffer;
    protected ops: TransferableOperation[];
    setCodecID(codecID: number): void;
    /**
     * Returns the id of the [[OperationTx]]
     */
    getTxType(): number;
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} containing an [[OperationTx]], parses it, populates the class, and returns the length of the [[OperationTx]] in bytes.
     *
     * @param bytes A {@link https://github.com/feross/buffer|Buffer} containing a raw [[OperationTx]]
     *
     * @returns The length of the raw [[OperationTx]]
     *
     * @remarks assume not-checksummed
     */
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[OperationTx]].
     */
    toBuffer(): Buffer;
    /**
     * Returns an array of [[TransferableOperation]]s in this transaction.
     */
    getOperations(): TransferableOperation[];
    /**
     * Takes the bytes of an [[UnsignedTx]] and returns an array of [[Credential]]s
     *
     * @param msg A Buffer for the [[UnsignedTx]]
     * @param kc An [[KeyChain]] used in signing
     *
     * @returns An array of [[Credential]]s
     */
    sign(msg: Buffer, kc: KeyChain): Credential[];
    clone(): this;
    create(...args: any[]): this;
    /**
     * Class representing an unsigned Operation transaction.
     *
     * @param networkID Optional networkID, [[DefaultNetworkID]]
     * @param blockchainID Optional blockchainID, default Buffer.alloc(32, 16)
     * @param outs Optional array of the [[TransferableOutput]]s
     * @param ins Optional array of the [[TransferableInput]]s
     * @param memo Optional {@link https://github.com/feross/buffer|Buffer} for the memo field
     * @param ops Array of [[Operation]]s used in the transaction
     */
    constructor(networkID?: number, blockchainID?: Buffer, outs?: TransferableOutput[], ins?: TransferableInput[], memo?: Buffer, ops?: TransferableOperation[]);
}
//# sourceMappingURL=operationtx.d.ts.map