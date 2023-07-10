/**
 * @packageDocumentation
 * @module API-EVM-Outputs
 */
import { Buffer } from "buffer/";
import BN from "bn.js";
import { Output, StandardAmountOutput, StandardTransferableOutput } from "../../common/output";
import { SerializedEncoding } from "../../utils/serialization";
import { EVMInput } from "./inputs";
/**
 * Takes a buffer representing the output and returns the proper Output instance.
 *
 * @param outputID A number representing the outputID parsed prior to the bytes passed in
 *
 * @returns An instance of an [[Output]]-extended class.
 */
export declare const SelectOutputClass: (outputID: number, ...args: any[]) => Output;
export declare class TransferableOutput extends StandardTransferableOutput {
    protected _typeName: string;
    protected _typeID: any;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    fromBuffer(bytes: Buffer, offset?: number): number;
}
export declare abstract class AmountOutput extends StandardAmountOutput {
    protected _typeName: string;
    protected _typeID: any;
    /**
     *
     * @param assetID An assetID which is wrapped around the Buffer of the Output
     */
    makeTransferable(assetID: Buffer): TransferableOutput;
    select(id: number, ...args: any[]): Output;
}
/**
 * An [[Output]] class which specifies an Output that carries an ammount for an assetID and uses secp256k1 signature scheme.
 */
export declare class SECPTransferOutput extends AmountOutput {
    protected _typeName: string;
    protected _typeID: number;
    /**
     * Returns the outputID for this output
     */
    getOutputID(): number;
    create(...args: any[]): this;
    clone(): this;
}
export declare class EVMOutput {
    protected address: Buffer;
    protected amount: Buffer;
    protected amountValue: BN;
    protected assetID: Buffer;
    /**
     * Returns a function used to sort an array of [[EVMOutput]]s
     */
    static comparator: () => (a: EVMOutput | EVMInput, b: EVMOutput | EVMInput) => 1 | -1 | 0;
    /**
     * Returns the address of the input as {@link https://github.com/feross/buffer|Buffer}
     */
    getAddress: () => Buffer;
    /**
     * Returns the address as a bech32 encoded string.
     */
    getAddressString: () => string;
    /**
     * Returns the amount as a {@link https://github.com/indutny/bn.js/|BN}.
     */
    getAmount: () => BN;
    /**
     * Returns the assetID of the input as {@link https://github.com/feross/buffer|Buffer}
     */
    getAssetID: () => Buffer;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[EVMOutput]].
     */
    toBuffer(): Buffer;
    /**
     * Decodes the [[EVMOutput]] as a {@link https://github.com/feross/buffer|Buffer} and returns the size.
     */
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns a base-58 representation of the [[EVMOutput]].
     */
    toString(): string;
    create(...args: any[]): this;
    clone(): this;
    /**
     * An [[EVMOutput]] class which contains address, amount, and assetID.
     *
     * @param address The address recieving the asset as a {@link https://github.com/feross/buffer|Buffer} or a string.
     * @param amount A {@link https://github.com/indutny/bn.js/|BN} or number representing the amount.
     * @param assetID The assetID which is being sent as a {@link https://github.com/feross/buffer|Buffer} or a string.
     */
    constructor(address?: Buffer | string, amount?: BN | number, assetID?: Buffer | string);
}
//# sourceMappingURL=outputs.d.ts.map