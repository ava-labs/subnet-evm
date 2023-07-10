/**
 * @packageDocumentation
 * @module Common-Inputs
 */
import { Buffer } from "buffer/";
import BN from "bn.js";
import { SigIdx } from "./credentials";
import { Serializable, SerializedEncoding } from "../utils/serialization";
export declare abstract class Input extends Serializable {
    protected _typeName: string;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected sigCount: Buffer;
    protected sigIdxs: SigIdx[];
    static comparator: () => (a: Input, b: Input) => 1 | -1 | 0;
    abstract getInputID(): number;
    /**
     * Returns the array of [[SigIdx]] for this [[Input]]
     */
    getSigIdxs: () => SigIdx[];
    abstract getCredentialID(): number;
    /**
     * Creates and adds a [[SigIdx]] to the [[Input]].
     *
     * @param addressIdx The index of the address to reference in the signatures
     * @param address The address of the source of the signature
     */
    addSignatureIdx: (addressIdx: number, address: Buffer) => void;
    fromBuffer(bytes: Buffer, offset?: number): number;
    toBuffer(): Buffer;
    /**
     * Returns a base-58 representation of the [[Input]].
     */
    toString(): string;
    abstract clone(): this;
    abstract create(...args: any[]): this;
    abstract select(id: number, ...args: any[]): Input;
}
export declare abstract class StandardParseableInput extends Serializable {
    protected _typeName: string;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    protected input: Input;
    /**
     * Returns a function used to sort an array of [[StandardParseableInput]]s
     */
    static comparator: () => (a: StandardParseableInput, b: StandardParseableInput) => 1 | -1 | 0;
    getInput: () => Input;
    abstract fromBuffer(bytes: Buffer, offset?: number): number;
    toBuffer(): Buffer;
    /**
     * Class representing an [[StandardParseableInput]] for a transaction.
     *
     * @param input A number representing the InputID of the [[StandardParseableInput]]
     */
    constructor(input?: Input);
}
export declare abstract class StandardTransferableInput extends StandardParseableInput {
    protected _typeName: string;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected txid: Buffer;
    protected outputidx: Buffer;
    protected assetID: Buffer;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} of the TxID.
     */
    getTxID: () => Buffer;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer}  of the OutputIdx.
     */
    getOutputIdx: () => Buffer;
    /**
     * Returns a base-58 string representation of the UTXOID this [[StandardTransferableInput]] references.
     */
    getUTXOID: () => string;
    /**
     * Returns the input.
     */
    getInput: () => Input;
    /**
     * Returns the assetID of the input.
     */
    getAssetID: () => Buffer;
    abstract fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[StandardTransferableInput]].
     */
    toBuffer(): Buffer;
    /**
     * Returns a base-58 representation of the [[StandardTransferableInput]].
     */
    toString(): string;
    /**
     * Class representing an [[StandardTransferableInput]] for a transaction.
     *
     * @param txid A {@link https://github.com/feross/buffer|Buffer} containing the transaction ID of the referenced UTXO
     * @param outputidx A {@link https://github.com/feross/buffer|Buffer} containing the index of the output in the transaction consumed in the [[StandardTransferableInput]]
     * @param assetID A {@link https://github.com/feross/buffer|Buffer} representing the assetID of the [[Input]]
     * @param input An [[Input]] to be made transferable
     */
    constructor(txid?: Buffer, outputidx?: Buffer, assetID?: Buffer, input?: Input);
}
/**
 * An [[Input]] class which specifies a token amount .
 */
export declare abstract class StandardAmountInput extends Input {
    protected _typeName: string;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected amount: Buffer;
    protected amountValue: BN;
    /**
     * Returns the amount as a {@link https://github.com/indutny/bn.js/|BN}.
     */
    getAmount: () => BN;
    /**
     * Popuates the instance from a {@link https://github.com/feross/buffer|Buffer} representing the [[AmountInput]] and returns the size of the input.
     */
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns the buffer representing the [[AmountInput]] instance.
     */
    toBuffer(): Buffer;
    /**
     * An [[AmountInput]] class which issues a payment on an assetID.
     *
     * @param amount A {@link https://github.com/indutny/bn.js/|BN} representing the amount in the input
     */
    constructor(amount?: BN);
}
//# sourceMappingURL=input.d.ts.map