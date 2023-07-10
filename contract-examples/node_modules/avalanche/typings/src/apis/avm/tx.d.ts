/**
 * @packageDocumentation
 * @module API-AVM-Transactions
 */
import { Buffer } from "buffer/";
import { KeyChain, KeyPair } from "./keychain";
import { StandardTx, StandardUnsignedTx } from "../../common/tx";
import { BaseTx } from "./basetx";
import { SerializedEncoding } from "../../utils/serialization";
/**
 * Takes a buffer representing the output and returns the proper [[BaseTx]] instance.
 *
 * @param txtype The id of the transaction type
 *
 * @returns An instance of an [[BaseTx]]-extended class.
 */
export declare const SelectTxClass: (txtype: number, ...args: any[]) => BaseTx;
export declare class UnsignedTx extends StandardUnsignedTx<KeyPair, KeyChain, BaseTx> {
    protected _typeName: string;
    protected _typeID: any;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    getTransaction(): BaseTx;
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Signs this [[UnsignedTx]] and returns signed [[StandardTx]]
     *
     * @param kc An [[KeyChain]] used in signing
     *
     * @returns A signed [[StandardTx]]
     */
    sign(kc: KeyChain): Tx;
}
export declare class Tx extends StandardTx<KeyPair, KeyChain, UnsignedTx> {
    protected _typeName: string;
    protected _typeID: any;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} containing an [[Tx]], parses it, populates the class, and returns the length of the Tx in bytes.
     *
     * @param bytes A {@link https://github.com/feross/buffer|Buffer} containing a raw [[Tx]]
     * @param offset A number representing the starting point of the bytes to begin parsing
     *
     * @returns The length of the raw [[Tx]]
     */
    fromBuffer(bytes: Buffer, offset?: number): number;
}
//# sourceMappingURL=tx.d.ts.map