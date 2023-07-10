/**
 * @packageDocumentation
 * @module API-AVM-Vertex
 */
import { Buffer } from "buffer/";
import { Tx } from "./tx";
import { Serializable } from "../../utils";
import BN from "bn.js";
/**
 * Class representing a Vertex
 */
export declare class Vertex extends Serializable {
    protected _typeName: string;
    protected _codecID: number;
    protected networkID: number;
    protected blockchainID: Buffer;
    protected height: BN;
    protected epoch: number;
    protected parentIDs: Buffer[];
    protected numParentIDs: number;
    protected txs: Tx[];
    protected numTxs: number;
    protected restrictions: Buffer[];
    protected numRestrictions: number;
    /**
     * Returns the NetworkID as a number
     */
    getNetworkID(): number;
    /**
     * Returns the BlockchainID as a CB58 string
     */
    getBlockchainID(): string;
    /**
     * Returns the Height as a {@link https://github.com/indutny/bn.js/|BN}.
     */
    getHeight(): BN;
    /**
     * Returns the Epoch as a number.
     */
    getEpoch(): number;
    /**
     * @returns An array of Buffers
     */
    getParentIDs(): Buffer[];
    /**
     * Returns array of UnsignedTxs.
     */
    getTxs(): Tx[];
    /**
     * @returns An array of Buffers
     */
    getRestrictions(): Buffer[];
    /**
     * Set the codecID
     *
     * @param codecID The codecID to set
     */
    setCodecID(codecID: number): void;
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} containing an [[Vertex]], parses it, populates the class, and returns the length of the Vertex in bytes.
     *
     * @param bytes A {@link https://github.com/feross/buffer|Buffer} containing a raw [[Vertex]]
     *
     * @returns The length of the raw [[Vertex]]
     *
     * @remarks assume not-checksummed
     */
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[Vertex]].
     */
    toBuffer(): Buffer;
    clone(): this;
    /**
     * Class representing a Vertex which is a container for AVM Transactions.
     *
     * @param networkID Optional, [[DefaultNetworkID]]
     * @param blockchainID Optional, default "2oYMBNV4eNHyqk2fjjV5nVQLDbtmNJzq5s3qs3Lo6ftnC6FByM"
     * @param height Optional, default new BN(0)
     * @param epoch Optional, default new BN(0)
     * @param parentIDs Optional, default []
     * @param txs Optional, default []
     * @param restrictions Optional, default []
     */
    constructor(networkID?: number, blockchainID?: string, height?: BN, epoch?: number, parentIDs?: Buffer[], txs?: Tx[], restrictions?: Buffer[]);
}
//# sourceMappingURL=vertex.d.ts.map