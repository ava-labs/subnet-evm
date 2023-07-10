/**
 * @packageDocumentation
 * @module API-PlatformVM-CreateSubnetTx
 */
import { Buffer } from "buffer/";
import { BaseTx } from "./basetx";
import { TransferableOutput, SECPOwnerOutput } from "./outputs";
import { TransferableInput } from "./inputs";
import { SerializedEncoding } from "../../utils/serialization";
export declare class CreateSubnetTx extends BaseTx {
    protected _typeName: string;
    protected _typeID: number;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected subnetOwners: SECPOwnerOutput;
    /**
     * Returns the id of the [[CreateSubnetTx]]
     */
    getTxType(): number;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} for the reward address.
     */
    getSubnetOwners(): SECPOwnerOutput;
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} containing an [[CreateSubnetTx]], parses it, populates the class, and returns the length of the [[CreateSubnetTx]] in bytes.
     *
     * @param bytes A {@link https://github.com/feross/buffer|Buffer} containing a raw [[CreateSubnetTx]]
     * @param offset A number for the starting position in the bytes.
     *
     * @returns The length of the raw [[CreateSubnetTx]]
     *
     * @remarks assume not-checksummed
     */
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[CreateSubnetTx]].
     */
    toBuffer(): Buffer;
    /**
     * Class representing an unsigned Create Subnet transaction.
     *
     * @param networkID Optional networkID, [[DefaultNetworkID]]
     * @param blockchainID Optional blockchainID, default Buffer.alloc(32, 16)
     * @param outs Optional array of the [[TransferableOutput]]s
     * @param ins Optional array of the [[TransferableInput]]s
     * @param memo Optional {@link https://github.com/feross/buffer|Buffer} for the memo field
     * @param subnetOwners Optional [[SECPOwnerOutput]] class for specifying who owns the subnet.
     */
    constructor(networkID?: number, blockchainID?: Buffer, outs?: TransferableOutput[], ins?: TransferableInput[], memo?: Buffer, subnetOwners?: SECPOwnerOutput);
}
//# sourceMappingURL=createsubnettx.d.ts.map