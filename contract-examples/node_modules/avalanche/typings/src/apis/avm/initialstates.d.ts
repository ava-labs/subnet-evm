/**
 * @packageDocumentation
 * @module API-AVM-InitialStates
 */
import { Buffer } from "buffer/";
import { Output } from "../../common/output";
import { Serializable, SerializedEncoding } from "../../utils/serialization";
/**
 * Class for creating initial output states used in asset creation
 */
export declare class InitialStates extends Serializable {
    protected _typeName: string;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected fxs: {
        [fxid: number]: Output[];
    };
    /**
     *
     * @param out The output state to add to the collection
     * @param fxid The FxID that will be used for this output, default AVMConstants.SECPFXID
     */
    addOutput(out: Output, fxid?: number): void;
    fromBuffer(bytes: Buffer, offset?: number): number;
    toBuffer(): Buffer;
}
//# sourceMappingURL=initialstates.d.ts.map