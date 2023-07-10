/**
 * @packageDocumentation
 * @module API-PlatformVM-AddSubnetValidatorTx
 */
import { Buffer } from "buffer/";
import { TransferableOutput } from "./outputs";
import { TransferableInput } from "./inputs";
import { Credential, SigIdx } from "../../common/credentials";
import { BaseTx } from "./basetx";
import { SerializedEncoding } from "../../utils/serialization";
import { SubnetAuth } from ".";
import { KeyChain } from "./keychain";
import BN from "bn.js";
/**
 * Class representing an unsigned AddSubnetValidatorTx transaction.
 */
export declare class AddSubnetValidatorTx extends BaseTx {
    protected _typeName: string;
    protected _typeID: number;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected nodeID: Buffer;
    protected startTime: Buffer;
    protected endTime: Buffer;
    protected weight: Buffer;
    protected subnetID: Buffer;
    protected subnetAuth: SubnetAuth;
    protected sigCount: Buffer;
    protected sigIdxs: SigIdx[];
    /**
     * Returns the id of the [[AddSubnetValidatorTx]]
     */
    getTxType(): number;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} for the stake amount.
     */
    getNodeID(): Buffer;
    /**
     * Returns a string for the nodeID amount.
     */
    getNodeIDString(): string;
    /**
     * Returns a {@link https://github.com/indutny/bn.js/|BN} for the startTime.
     */
    getStartTime(): BN;
    /**
     * Returns a {@link https://github.com/indutny/bn.js/|BN} for the endTime.
     */
    getEndTime(): BN;
    /**
     * Returns a {@link https://github.com/indutny/bn.js/|BN} for the weight
     */
    getWeight(): BN;
    /**
     * Returns the subnetID as a string
     */
    getSubnetID(): string;
    /**
     * Returns the subnetAuth
     */
    getSubnetAuth(): SubnetAuth;
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} containing an [[AddSubnetValidatorTx]], parses it, populates the class, and returns the length of the [[CreateChainTx]] in bytes.
     *
     * @param bytes A {@link https://github.com/feross/buffer|Buffer} containing a raw [[AddSubnetValidatorTx]]
     *
     * @returns The length of the raw [[AddSubnetValidatorTx]]
     *
     * @remarks assume not-checksummed
     */
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[CreateChainTx]].
     */
    toBuffer(): Buffer;
    clone(): this;
    create(...args: any[]): this;
    /**
     * Creates and adds a [[SigIdx]] to the [[AddSubnetValidatorTx]].
     *
     * @param addressIdx The index of the address to reference in the signatures
     * @param address The address of the source of the signature
     */
    addSignatureIdx(addressIdx: number, address: Buffer): void;
    /**
     * Returns the array of [[SigIdx]] for this [[Input]]
     */
    getSigIdxs(): SigIdx[];
    getCredentialID(): number;
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
     * Class representing an unsigned AddSubnetValidator transaction.
     *
     * @param networkID Optional networkID, [[DefaultNetworkID]]
     * @param blockchainID Optional blockchainID, default Buffer.alloc(32, 16)
     * @param outs Optional array of the [[TransferableOutput]]s
     * @param ins Optional array of the [[TransferableInput]]s
     * @param memo Optional {@link https://github.com/feross/buffer|Buffer} for the memo field
     * @param nodeID Optional. The node ID of the validator being added.
     * @param startTime Optional. The Unix time when the validator starts validating the Primary Network.
     * @param endTime Optional. The Unix time when the validator stops validating the Primary Network (and staked AVAX is returned).
     * @param weight Optional. Weight of this validator used when sampling
     * @param subnetID Optional. ID of the subnet this validator is validating
     */
    constructor(networkID?: number, blockchainID?: Buffer, outs?: TransferableOutput[], ins?: TransferableInput[], memo?: Buffer, nodeID?: Buffer, startTime?: BN, endTime?: BN, weight?: BN, subnetID?: string | Buffer);
}
//# sourceMappingURL=addsubnetvalidatortx.d.ts.map