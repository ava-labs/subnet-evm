/**
 * @packageDocumentation
 * @module API-AVM-CreateAssetTx
 */
import { Buffer } from 'buffer/';
import { Serializable, SerializedEncoding } from '../../utils/serialization';
export declare class GenesisState extends Serializable {
    protected _typeName: string;
    protected _codecID: number;
    protected networkid: Buffer;
    protected blockchainid: Buffer;
    serialize(encoding?: SerializedEncoding): object;
    /**
    * Class representing a GenesisState
    *
    * @param networkid Optional networkid, [[DefaultNetworkID]]
    * @param blockchainid Optional blockchainid, default Buffer.alloc(32, 16)
    */
    constructor(networkid?: number, blockchainid?: Buffer);
}
//# sourceMappingURL=genesisstate.d.ts.map