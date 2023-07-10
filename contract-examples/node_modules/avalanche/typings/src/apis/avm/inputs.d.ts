/**
 * @packageDocumentation
 * @module API-AVM-Inputs
 */
import { Buffer } from "buffer/";
import { Input, StandardTransferableInput, StandardAmountInput } from "../../common/input";
import { SerializedEncoding } from "../../utils/serialization";
/**
 * Takes a buffer representing the output and returns the proper [[Input]] instance.
 *
 * @param inputid A number representing the inputID parsed prior to the bytes passed in
 *
 * @returns An instance of an [[Input]]-extended class.
 */
export declare const SelectInputClass: (inputid: number, ...args: any[]) => Input;
export declare class TransferableInput extends StandardTransferableInput {
    protected _typeName: string;
    protected _typeID: any;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
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
    protected _codecID: number;
    protected _typeID: number;
    /**
     * Set the codecID
     *
     * @param codecID The codecID to set
     */
    setCodecID(codecID: number): void;
    /**
     * Returns the inputID for this input
     */
    getInputID(): number;
    getCredentialID(): number;
    create(...args: any[]): this;
    clone(): this;
}
//# sourceMappingURL=inputs.d.ts.map