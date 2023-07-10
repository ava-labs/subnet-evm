/**
 * @packageDocumentation
 * @module API-EVM-ExportTx
 */
import { Buffer } from "buffer/";
import { KeyChain } from "./keychain";
import { EVMBaseTx } from "./basetx";
import { Credential } from "../../common/credentials";
import { EVMInput } from "./inputs";
import { SerializedEncoding } from "../../utils/serialization";
import { TransferableOutput } from "./outputs";
export declare class ExportTx extends EVMBaseTx {
    protected _typeName: string;
    protected _typeID: number;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected destinationChain: Buffer;
    protected numInputs: Buffer;
    protected inputs: EVMInput[];
    protected numExportedOutputs: Buffer;
    protected exportedOutputs: TransferableOutput[];
    /**
     * Returns the destinationChain as a {@link https://github.com/feross/buffer|Buffer}
     */
    getDestinationChain(): Buffer;
    /**
     * Returns the inputs as an array of [[EVMInputs]]
     */
    getInputs(): EVMInput[];
    /**
     * Returns the outs as an array of [[EVMOutputs]]
     */
    getExportedOutputs(): TransferableOutput[];
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[ExportTx]].
     */
    toBuffer(): Buffer;
    /**
     * Decodes the [[ExportTx]] as a {@link https://github.com/feross/buffer|Buffer} and returns the size.
     */
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns a base-58 representation of the [[ExportTx]].
     */
    toString(): string;
    /**
     * Takes the bytes of an [[UnsignedTx]] and returns an array of [[Credential]]s
     *
     * @param msg A Buffer for the [[UnsignedTx]]
     * @param kc An [[KeyChain]] used in signing
     *
     * @returns An array of [[Credential]]s
     */
    sign(msg: Buffer, kc: KeyChain): Credential[];
    /**
     * Class representing a ExportTx.
     *
     * @param networkID Optional networkID
     * @param blockchainID Optional blockchainID, default Buffer.alloc(32, 16)
     * @param destinationChain Optional destinationChain, default Buffer.alloc(32, 16)
     * @param inputs Optional array of the [[EVMInputs]]s
     * @param exportedOutputs Optional array of the [[EVMOutputs]]s
     */
    constructor(networkID?: number, blockchainID?: Buffer, destinationChain?: Buffer, inputs?: EVMInput[], exportedOutputs?: TransferableOutput[]);
}
//# sourceMappingURL=exporttx.d.ts.map