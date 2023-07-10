/**
 * @packageDocumentation
 * @module API-PlatformVM-CreateChainTx
 */
import { Buffer } from "buffer/";
import { TransferableOutput } from "./outputs";
import { TransferableInput } from "./inputs";
import { Credential, SigIdx } from "../../common/credentials";
import { BaseTx } from "./basetx";
import { SerializedEncoding } from "../../utils/serialization";
import { GenesisData } from "../avm";
import { SubnetAuth } from ".";
import { KeyChain } from "./keychain";
/**
 * Class representing an unsigned CreateChainTx transaction.
 */
export declare class CreateChainTx extends BaseTx {
    protected _typeName: string;
    protected _typeID: number;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected subnetID: Buffer;
    protected chainName: string;
    protected vmID: Buffer;
    protected numFXIDs: Buffer;
    protected fxIDs: Buffer[];
    protected genesisData: Buffer;
    protected subnetAuth: SubnetAuth;
    protected sigCount: Buffer;
    protected sigIdxs: SigIdx[];
    /**
     * Returns the id of the [[CreateChainTx]]
     */
    getTxType(): number;
    /**
     * Returns the subnetAuth
     */
    getSubnetAuth(): SubnetAuth;
    /**
     * Returns the subnetID as a string
     */
    getSubnetID(): string;
    /**
     * Returns a string of the chainName
     */
    getChainName(): string;
    /**
     * Returns a Buffer of the vmID
     */
    getVMID(): Buffer;
    /**
     * Returns an array of fxIDs as Buffers
     */
    getFXIDs(): Buffer[];
    /**
     * Returns a string of the genesisData
     */
    getGenesisData(): string;
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} containing an [[CreateChainTx]], parses it, populates the class, and returns the length of the [[CreateChainTx]] in bytes.
     *
     * @param bytes A {@link https://github.com/feross/buffer|Buffer} containing a raw [[CreateChainTx]]
     *
     * @returns The length of the raw [[CreateChainTx]]
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
     * Class representing an unsigned CreateChain transaction.
     *
     * @param networkID Optional networkID, [[DefaultNetworkID]]
     * @param blockchainID Optional blockchainID, default Buffer.alloc(32, 16)
     * @param outs Optional array of the [[TransferableOutput]]s
     * @param ins Optional array of the [[TransferableInput]]s
     * @param memo Optional {@link https://github.com/feross/buffer|Buffer} for the memo field
     * @param subnetID Optional ID of the Subnet that validates this blockchain.
     * @param chainName Optional A human readable name for the chain; need not be unique
     * @param vmID Optional ID of the VM running on the new chain
     * @param fxIDs Optional IDs of the feature extensions running on the new chain
     * @param genesisData Optional Byte representation of genesis state of the new chain
     */
    constructor(networkID?: number, blockchainID?: Buffer, outs?: TransferableOutput[], ins?: TransferableInput[], memo?: Buffer, subnetID?: string | Buffer, chainName?: string, vmID?: string, fxIDs?: string[], genesisData?: string | GenesisData);
}
//# sourceMappingURL=createchaintx.d.ts.map