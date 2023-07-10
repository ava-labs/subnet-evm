/**
 * @packageDocumentation
 * @module API-EVM-Inputs
 */
import { Buffer } from "buffer/";
import { Input, StandardTransferableInput, StandardAmountInput } from "../../common/input";
import { SerializedEncoding } from "../../utils/serialization";
import { EVMOutput } from "./outputs";
import BN from "bn.js";
import { SigIdx } from "../../common/credentials";
/**
 * Takes a buffer representing the output and returns the proper [[Input]] instance.
 *
 * @param inputID A number representing the inputID parsed prior to the bytes passed in
 *
 * @returns An instance of an [[Input]]-extended class.
 */
export declare const SelectInputClass: (inputID: number, ...args: any[]) => Input;
export declare class TransferableInput extends StandardTransferableInput {
    protected _typeName: string;
    protected _typeID: any;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    /**
     *
     * Assesses the amount to be paid based on the number of signatures required
     * @returns the amount to be paid
     */
    getCost: () => number;
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} containing a [[TransferableInput]], parses it, populates the class, and returns the length of the [[TransferableInput]] in bytes.
     *
     * @param bytes A {@link https://github.com/feross/buffer|Buffer} containing a raw [[TransferableInput]]
     *
     * @returns The length of the raw [[TransferableInput]]
     */
    fromBuffer(bytes: Buffer, offset?: number): number;
}
export declare abstract class AmountInput extends StandardAmountInput {
    protected _typeName: string;
    protected _typeID: any;
    select(id: number, ...args: any[]): Input;
}
export declare class SECPTransferInput extends AmountInput {
    protected _typeName: string;
    protected _typeID: number;
    /**
     * Returns the inputID for this input
     */
    getInputID(): number;
    getCredentialID: () => number;
    create(...args: any[]): this;
    clone(): this;
}
export declare class EVMInput extends EVMOutput {
    protected nonce: Buffer;
    protected nonceValue: BN;
    protected sigCount: Buffer;
    protected sigIdxs: SigIdx[];
    /**
     * Returns the array of [[SigIdx]] for this [[Input]]
     */
    getSigIdxs: () => SigIdx[];
    /**
     * Creates and adds a [[SigIdx]] to the [[Input]].
     *
     * @param addressIdx The index of the address to reference in the signatures
     * @param address The address of the source of the signature
     */
    addSignatureIdx: (addressIdx: number, address: Buffer) => void;
    /**
     * Returns the nonce as a {@link https://github.com/indutny/bn.js/|BN}.
     */
    getNonce: () => BN;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[EVMOutput]].
     */
    toBuffer(): Buffer;
    getCredentialID: () => number;
    /**
     * Decodes the [[EVMInput]] as a {@link https://github.com/feross/buffer|Buffer} and returns the size.
     *
     * @param bytes The bytes as a {@link https://github.com/feross/buffer|Buffer}.
     * @param offset An offset as a number.
     */
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns a base-58 representation of the [[EVMInput]].
     */
    toString(): string;
    create(...args: any[]): this;
    clone(): this;
    /**
     * An [[EVMInput]] class which contains address, amount, assetID, nonce.
     *
     * @param address is the EVM address from which to transfer funds.
     * @param amount is the amount of the asset to be transferred (specified in nAVAX for AVAX and the smallest denomination for all other assets).
     * @param assetID The assetID which is being sent as a {@link https://github.com/feross/buffer|Buffer} or as a string.
     * @param nonce A {@link https://github.com/indutny/bn.js/|BN} or a number representing the nonce.
     */
    constructor(address?: Buffer | string, amount?: BN | number, assetID?: Buffer | string, nonce?: BN | number);
}
//# sourceMappingURL=inputs.d.ts.map