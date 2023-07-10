/**
 * @packageDocumentation
 * @module API-EVM-BaseTx
 */
import { Buffer } from "buffer/";
import { KeyChain, KeyPair } from "./keychain";
import { EVMStandardBaseTx } from "../../common/evmtx";
import { Credential } from "../../common/credentials";
import { SerializedEncoding } from "../../utils/serialization";
/**
 * Class representing a base for all transactions.
 */
export declare class EVMBaseTx extends EVMStandardBaseTx<KeyPair, KeyChain> {
    protected _typeName: string;
    protected _typeID: any;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    /**
     * Returns the id of the [[BaseTx]]
     */
    getTxType(): number;
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} containing an [[BaseTx]], parses it, populates the class, and returns the length of the BaseTx in bytes.
     *
     * @param bytes A {@link https://github.com/feross/buffer|Buffer} containing a raw [[BaseTx]]
     *
     * @returns The length of the raw [[BaseTx]]
     *
     * @remarks assume not-checksummed
     */
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Takes the bytes of an [[UnsignedTx]] and returns an array of [[Credential]]s
     *
     * @param msg A Buffer for the [[UnsignedTx]]
     * @param kc An [[KeyChain]] used in signing
     *
     * @returns An array of [[Credential]]s
     */
    sign(msg: Buffer, kc: KeyChain): Credential[];
    clone(): this;
    create(...args: any[]): this;
    select(id: number, ...args: any[]): this;
    /**
     * Class representing an EVMBaseTx which is the foundation for all EVM transactions.
     *
     * @param networkID Optional networkID, [[DefaultNetworkID]]
     * @param blockchainID Optional blockchainID, default Buffer.alloc(32, 16)
     */
    constructor(networkID?: number, blockchainID?: Buffer);
}
//# sourceMappingURL=basetx.d.ts.map