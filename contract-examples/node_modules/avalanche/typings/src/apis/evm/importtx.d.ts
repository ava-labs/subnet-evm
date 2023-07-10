/**
 * @packageDocumentation
 * @module API-EVM-ImportTx
 */
import { Buffer } from "buffer/";
import BN from "bn.js";
import { EVMOutput } from "./outputs";
import { TransferableInput } from "./inputs";
import { EVMBaseTx } from "./basetx";
import { Credential } from "../../common/credentials";
import { KeyChain } from "./keychain";
import { SerializedEncoding } from "../../utils/serialization";
/**
 * Class representing an unsigned Import transaction.
 */
export declare class ImportTx extends EVMBaseTx {
    protected _typeName: string;
    protected _typeID: number;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected sourceChain: Buffer;
    protected numIns: Buffer;
    protected importIns: TransferableInput[];
    protected numOuts: Buffer;
    protected outs: EVMOutput[];
    /**
     * Returns the id of the [[ImportTx]]
     */
    getTxType(): number;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} for the source chainid.
     */
    getSourceChain(): Buffer;
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} containing an [[ImportTx]], parses it,
     * populates the class, and returns the length of the [[ImportTx]] in bytes.
     *
     * @param bytes A {@link https://github.com/feross/buffer|Buffer} containing a raw [[ImportTx]]
     * @param offset A number representing the byte offset. Defaults to 0.
     *
     * @returns The length of the raw [[ImportTx]]
     *
     * @remarks assume not-checksummed
     */
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[ImportTx]].
     */
    toBuffer(): Buffer;
    /**
     * Returns an array of [[TransferableInput]]s in this transaction.
     */
    getImportInputs(): TransferableInput[];
    /**
     * Returns an array of [[EVMOutput]]s in this transaction.
     */
    getOuts(): EVMOutput[];
    clone(): this;
    create(...args: any[]): this;
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
     * Class representing an unsigned Import transaction.
     *
     * @param networkID Optional networkID, [[DefaultNetworkID]]
     * @param blockchainID Optional blockchainID, default Buffer.alloc(32, 16)
     * @param sourceChainID Optional chainID for the source inputs to import. Default Buffer.alloc(32, 16)
     * @param importIns Optional array of [[TransferableInput]]s used in the transaction
     * @param outs Optional array of the [[EVMOutput]]s
     * @param fee Optional the fee as a BN
     */
    constructor(networkID?: number, blockchainID?: Buffer, sourceChainID?: Buffer, importIns?: TransferableInput[], outs?: EVMOutput[], fee?: BN);
    private validateOuts;
}
//# sourceMappingURL=importtx.d.ts.map