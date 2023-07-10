/**
 * @packageDocumentation
 * @module API-AVM-Outputs
 */
import { Buffer } from "buffer/";
import BN from "bn.js";
import { Output, StandardAmountOutput, StandardTransferableOutput, BaseNFTOutput } from "../../common/output";
import { SerializedEncoding } from "../../utils/serialization";
/**
 * Takes a buffer representing the output and returns the proper Output instance.
 *
 * @param outputid A number representing the inputID parsed prior to the bytes passed in
 *
 * @returns An instance of an [[Output]]-extended class.
 */
export declare const SelectOutputClass: (outputid: number, ...args: any[]) => Output;
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
export declare abstract class NFTOutput extends BaseNFTOutput {
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
    protected _codecID: number;
    protected _typeID: number;
    /**
     * Set the codecID
     *
     * @param codecID The codecID to set
     */
    setCodecID(codecID: number): void;
    /**
     * Returns the outputID for this output
     */
    getOutputID(): number;
    create(...args: any[]): this;
    clone(): this;
}
/**
 * An [[Output]] class which specifies an Output that carries an ammount for an assetID and uses secp256k1 signature scheme.
 */
export declare class SECPMintOutput extends Output {
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
     * Returns the outputID for this output
     */
    getOutputID(): number;
    /**
     *
     * @param assetID An assetID which is wrapped around the Buffer of the Output
     */
    makeTransferable(assetID: Buffer): TransferableOutput;
    create(...args: any[]): this;
    clone(): this;
    select(id: number, ...args: any[]): Output;
}
/**
 * An [[Output]] class which specifies an Output that carries an NFT Mint and uses secp256k1 signature scheme.
 */
export declare class NFTMintOutput extends NFTOutput {
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
     * Returns the outputID for this output
     */
    getOutputID(): number;
    /**
     * Popuates the instance from a {@link https://github.com/feross/buffer|Buffer} representing the [[NFTMintOutput]] and returns the size of the output.
     */
    fromBuffer(utxobuff: Buffer, offset?: number): number;
    /**
     * Returns the buffer representing the [[NFTMintOutput]] instance.
     */
    toBuffer(): Buffer;
    create(...args: any[]): this;
    clone(): this;
    /**
     * An [[Output]] class which contains an NFT mint for an assetID.
     *
     * @param groupID A number specifies the group this NFT is issued to
     * @param addresses An array of {@link https://github.com/feross/buffer|Buffer}s representing  addresses
     * @param locktime A {@link https://github.com/indutny/bn.js/|BN} representing the locktime
     * @param threshold A number representing the the threshold number of signers required to sign the transaction
  
     */
    constructor(groupID?: number, addresses?: Buffer[], locktime?: BN, threshold?: number);
}
/**
 * An [[Output]] class which specifies an Output that carries an NFT and uses secp256k1 signature scheme.
 */
export declare class NFTTransferOutput extends NFTOutput {
    protected _typeName: string;
    protected _codecID: number;
    protected _typeID: number;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected sizePayload: Buffer;
    protected payload: Buffer;
    /**
     * Set the codecID
     *
     * @param codecID The codecID to set
     */
    setCodecID(codecID: number): void;
    /**
     * Returns the outputID for this output
     */
    getOutputID(): number;
    /**
     * Returns the payload as a {@link https://github.com/feross/buffer|Buffer} with content only.
     */
    getPayload: () => Buffer;
    /**
     * Returns the payload as a {@link https://github.com/feross/buffer|Buffer} with length of payload prepended.
     */
    getPayloadBuffer: () => Buffer;
    /**
     * Popuates the instance from a {@link https://github.com/feross/buffer|Buffer} representing the [[NFTTransferOutput]] and returns the size of the output.
     */
    fromBuffer(utxobuff: Buffer, offset?: number): number;
    /**
     * Returns the buffer representing the [[NFTTransferOutput]] instance.
     */
    toBuffer(): Buffer;
    create(...args: any[]): this;
    clone(): this;
    /**
       * An [[Output]] class which contains an NFT on an assetID.
       *
       * @param groupID A number representing the amount in the output
       * @param payload A {@link https://github.com/feross/buffer|Buffer} of max length 1024
       * @param addresses An array of {@link https://github.com/feross/buffer|Buffer}s representing addresses
       * @param locktime A {@link https://github.com/indutny/bn.js/|BN} representing the locktime
       * @param threshold A number representing the the threshold number of signers required to sign the transaction
  
       */
    constructor(groupID?: number, payload?: Buffer, addresses?: Buffer[], locktime?: BN, threshold?: number);
}
//# sourceMappingURL=outputs.d.ts.map