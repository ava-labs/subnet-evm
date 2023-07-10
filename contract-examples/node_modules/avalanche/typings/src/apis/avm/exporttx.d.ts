/**
 * @packageDocumentation
 * @module API-AVM-ExportTx
 */
import { Buffer } from "buffer/";
import { TransferableOutput } from "./outputs";
import { TransferableInput } from "./inputs";
import { BaseTx } from "./basetx";
import BN from "bn.js";
import { SerializedEncoding } from "../../utils/serialization";
/**
 * Class representing an unsigned Export transaction.
 */
export declare class ExportTx extends BaseTx {
    protected _typeName: string;
    protected _codecID: number;
    protected _typeID: number;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected destinationChain: Buffer;
    protected numOuts: Buffer;
    protected exportOuts: TransferableOutput[];
    /**
     * Set the codecID
     *
     * @param codecID The codecID to set
     */
    setCodecID(codecID: number): void;
    /**
     * Returns the id of the [[ExportTx]]
     */
    getTxType(): number;
    /**
     * Returns an array of [[TransferableOutput]]s in this transaction.
     */
    getExportOutputs(): TransferableOutput[];
    /**
     * Returns the totall exported amount as a {@link https://github.com/indutny/bn.js/|BN}.
     */
    getExportTotal(): BN;
    getTotalOuts(): TransferableOutput[];
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} for the destination chainid.
     */
    getDestinationChain(): Buffer;
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} containing an [[ExportTx]], parses it, populates the class, and returns the length of the [[ExportTx]] in bytes.
     *
     * @param bytes A {@link https://github.com/feross/buffer|Buffer} containing a raw [[ExportTx]]
     *
     * @returns The length of the raw [[ExportTx]]
     *
     * @remarks assume not-checksummed
     */
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[ExportTx]].
     */
    toBuffer(): Buffer;
    clone(): this;
    create(...args: any[]): this;
    /**
     * Class representing an unsigned Export transaction.
     *
     * @param networkID Optional networkID, [[DefaultNetworkID]]
     * @param blockchainID Optional blockchainID, default Buffer.alloc(32, 16)
     * @param outs Optional array of the [[TransferableOutput]]s
     * @param ins Optional array of the [[TransferableInput]]s
     * @param memo Optional {@link https://github.com/feross/buffer|Buffer} for the memo field
     * @param destinationChain Optional chainid which identifies where the funds will sent to
     * @param exportOuts Array of [[TransferableOutputs]]s used in the transaction
     */
    constructor(networkID?: number, blockchainID?: Buffer, outs?: TransferableOutput[], ins?: TransferableInput[], memo?: Buffer, destinationChain?: Buffer, exportOuts?: TransferableOutput[]);
}
//# sourceMappingURL=exporttx.d.ts.map