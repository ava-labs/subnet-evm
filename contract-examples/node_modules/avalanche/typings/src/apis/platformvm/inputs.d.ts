/**
 * @packageDocumentation
 * @module API-PlatformVM-Inputs
 */
import { Buffer } from "buffer/";
import { Input, StandardTransferableInput, StandardAmountInput, StandardParseableInput } from "../../common/input";
import { SerializedEncoding } from "../../utils/serialization";
import BN from "bn.js";
/**
 * Takes a buffer representing the output and returns the proper [[Input]] instance.
 *
 * @param inputid A number representing the inputID parsed prior to the bytes passed in
 *
 * @returns An instance of an [[Input]]-extended class.
 */
export declare const SelectInputClass: (inputid: number, ...args: any[]) => Input;
export declare class ParseableInput extends StandardParseableInput {
    protected _typeName: string;
    protected _typeID: any;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    fromBuffer(bytes: Buffer, offset?: number): number;
}
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
    protected _typeID: number;
    /**
     * Returns the inputID for this input
     */
    getInputID(): number;
    getCredentialID: () => number;
    create(...args: any[]): this;
    clone(): this;
}
/**
 * An [[Input]] class which specifies an input that has a locktime which can also enable staking of the value held, preventing transfers but not validation.
 */
export declare class StakeableLockIn extends AmountInput {
    protected _typeName: string;
    protected _typeID: number;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected stakeableLocktime: Buffer;
    protected transferableInput: ParseableInput;
    private synchronize;
    getStakeableLocktime(): BN;
    getTransferablInput(): ParseableInput;
    /**
     * Returns the inputID for this input
     */
    getInputID(): number;
    getCredentialID: () => number;
    /**
     * Popuates the instance from a {@link https://github.com/feross/buffer|Buffer} representing the [[StakeableLockIn]] and returns the size of the output.
     */
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns the buffer representing the [[StakeableLockIn]] instance.
     */
    toBuffer(): Buffer;
    create(...args: any[]): this;
    clone(): this;
    select(id: number, ...args: any[]): Input;
    /**
     * A [[Output]] class which specifies an [[Input]] that has a locktime which can also enable staking of the value held, preventing transfers but not validation.
     *
     * @param amount A {@link https://github.com/indutny/bn.js/|BN} representing the amount in the input
     * @param stakeableLocktime A {@link https://github.com/indutny/bn.js/|BN} representing the stakeable locktime
     * @param transferableInput A [[ParseableInput]] which is embedded into this input.
     */
    constructor(amount?: BN, stakeableLocktime?: BN, transferableInput?: ParseableInput);
}
//# sourceMappingURL=inputs.d.ts.map