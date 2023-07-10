/**
 * @packageDocumentation
 * @module API-AVM-GenesisData
 */
import { Buffer } from "buffer/";
import { Serializable, SerializedEncoding } from "../../utils/serialization";
import { GenesisAsset } from ".";
export declare class GenesisData extends Serializable {
    protected _typeName: string;
    protected _codecID: number;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected genesisAssets: GenesisAsset[];
    protected networkID: Buffer;
    /**
     * Returns the GenesisAssets[]
     */
    getGenesisAssets: () => GenesisAsset[];
    /**
     * Returns the NetworkID as a number
     */
    getNetworkID: () => number;
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} containing an [[GenesisAsset]], parses it, populates the class, and returns the length of the [[GenesisAsset]] in bytes.
     *
     * @param bytes A {@link https://github.com/feross/buffer|Buffer} containing a raw [[GenesisAsset]]
     *
     * @returns The length of the raw [[GenesisAsset]]
     *
     * @remarks assume not-checksummed
     */
    fromBuffer(bytes: Buffer, offset?: number): number;
    /**
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[GenesisData]].
     */
    toBuffer(): Buffer;
    /**
     * Class representing AVM GenesisData
     *
     * @param genesisAssets Optional GenesisAsset[]
     * @param networkID Optional DefaultNetworkID
     */
    constructor(genesisAssets?: GenesisAsset[], networkID?: number);
}
//# sourceMappingURL=genesisdata.d.ts.map