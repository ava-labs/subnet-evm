/**
 * @packageDocumentation
 * @module Common-Transactions
 */
import { Buffer } from "buffer/";
import { Credential } from "./credentials";
import BN from "bn.js";
import { StandardKeyChain, StandardKeyPair } from "./keychain";
import { Serializable, SerializedEncoding } from "../utils/serialization";
/**
 * Class representing a base for all transactions.
 */
export declare abstract class EVMStandardBaseTx<KPClass extends StandardKeyPair, KCClass extends StandardKeyChain<KPClass>> extends Serializable {
    protected _typeName: string;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected networkID: Buffer;
    protected blockchainID: Buffer;
    /**
     * Returns the id of the [[StandardBaseTx]]
     */
    abstract getTxType(): number;
    /**
     * Returns the NetworkID as a number
     */
    getNetworkID(): number;
    /**
     * Returns the Buffer representation of the BlockchainID
     */
    getBlockchainID(): Buffer;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[StandardBaseTx]].
     */
    toBuffer(): Buffer;
    /**
     * Returns a base-58 representation of the [[StandardBaseTx]].
     */
    toString(): string;
    abstract clone(): this;
    abstract create(...args: any[]): this;
    abstract select(id: number, ...args: any[]): this;
    /**
     * Class representing a StandardBaseTx which is the foundation for all transactions.
     *
     * @param networkID Optional networkID, [[DefaultNetworkID]]
     * @param blockchainID Optional blockchainID, default Buffer.alloc(32, 16)
     * @param outs Optional array of the [[TransferableOutput]]s
     * @param ins Optional array of the [[TransferableInput]]s
     */
    constructor(networkID?: number, blockchainID?: Buffer);
}
/**
 * Class representing an unsigned transaction.
 */
export declare abstract class EVMStandardUnsignedTx<KPClass extends StandardKeyPair, KCClass extends StandardKeyChain<KPClass>, SBTx extends EVMStandardBaseTx<KPClass, KCClass>> extends Serializable {
    protected _typeName: string;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected codecID: number;
    protected transaction: SBTx;
    /**
     * Returns the CodecID as a number
     */
    getCodecID(): number;
    /**
     * Returns the {@link https://github.com/feross/buffer|Buffer} representation of the CodecID
     */
    getCodecIDBuffer(): Buffer;
    /**
     * Returns the inputTotal as a BN
     */
    getInputTotal(assetID: Buffer): BN;
    /**
     * Returns the outputTotal as a BN
     */
    getOutputTotal(assetID: Buffer): BN;
    /**
     * Returns the number of burned tokens as a BN
     */
    getBurn(assetID: Buffer): BN;
    /**
     * Returns the Transaction
     */
    abstract getTransaction(): SBTx;
    abstract fromBuffer(bytes: Buffer, offset?: number): number;
    toBuffer(): Buffer;
    /**
     * Signs this [[UnsignedTx]] and returns signed [[StandardTx]]
     *
     * @param kc An [[KeyChain]] used in signing
     *
     * @returns A signed [[StandardTx]]
     */
    abstract sign(kc: KCClass): EVMStandardTx<KPClass, KCClass, EVMStandardUnsignedTx<KPClass, KCClass, SBTx>>;
    constructor(transaction?: SBTx, codecID?: number);
}
/**
 * Class representing a signed transaction.
 */
export declare abstract class EVMStandardTx<KPClass extends StandardKeyPair, KCClass extends StandardKeyChain<KPClass>, SUBTx extends EVMStandardUnsignedTx<KPClass, KCClass, EVMStandardBaseTx<KPClass, KCClass>>> extends Serializable {
    protected _typeName: string;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    protected unsignedTx: SUBTx;
    protected credentials: Credential[];
    /**
     * Returns the [[StandardUnsignedTx]]
     */
    getUnsignedTx(): SUBTx;
    abstract fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[StandardTx]].
     */
    toBuffer(): Buffer;
    /**
     * Takes a base-58 string containing an [[StandardTx]], parses it, populates the class, and returns the length of the Tx in bytes.
     *
     * @param serialized A base-58 string containing a raw [[StandardTx]]
     *
     * @returns The length of the raw [[StandardTx]]
     *
     * @remarks
     * unlike most fromStrings, it expects the string to be serialized in cb58 format
     */
    fromString(serialized: string): number;
    /**
     * Returns a cb58 representation of the [[StandardTx]].
     *
     * @remarks
     * unlike most toStrings, this returns in cb58 serialization format
     */
    toString(): string;
    /**
     * Class representing a signed transaction.
     *
     * @param unsignedTx Optional [[StandardUnsignedTx]]
     * @param signatures Optional array of [[Credential]]s
     */
    constructor(unsignedTx?: SUBTx, credentials?: Credential[]);
}
//# sourceMappingURL=evmtx.d.ts.map