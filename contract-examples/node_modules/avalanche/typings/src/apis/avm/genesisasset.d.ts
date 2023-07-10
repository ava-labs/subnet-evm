/**
 * @packageDocumentation
 * @module API-AVM-GenesisAsset
 */
import { Buffer } from "buffer/";
import { InitialStates } from "./initialstates";
import { SerializedEncoding } from "../../utils/serialization";
import { CreateAssetTx } from "./createassettx";
export declare class GenesisAsset extends CreateAssetTx {
    protected _typeName: string;
    protected _codecID: any;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected assetAlias: string;
    /**
     * Returns the string representation of the assetAlias
     */
    getAssetAlias: () => string;
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
     * Returns a {@link https://github.com/feross/buffer|Buffer} representation of the [[GenesisAsset]].
     */
    toBuffer(networkID?: number): Buffer;
    /**
     * Class representing a GenesisAsset
     *
     * @param assetAlias Optional String for the asset alias
     * @param name Optional String for the descriptive name of the asset
     * @param symbol Optional String for the ticker symbol of the asset
     * @param denomination Optional number for the denomination which is 10^D. D must be >= 0 and <= 32. Ex: $1 AVAX = 10^9 $nAVAX
     * @param initialState Optional [[InitialStates]] that represent the intial state of a created asset
     * @param memo Optional {@link https://github.com/feross/buffer|Buffer} for the memo field
     */
    constructor(assetAlias?: string, name?: string, symbol?: string, denomination?: number, initialState?: InitialStates, memo?: Buffer);
}
//# sourceMappingURL=genesisasset.d.ts.map