import { Buffer } from "buffer/";
import { Serialized } from "../common";
export declare const SERIALIZATIONVERSION: number;
export declare type SerializedType = "hex" | "BN" | "Buffer" | "bech32" | "nodeID" | "privateKey" | "cb58" | "base58" | "base64" | "decimalString" | "number" | "utf8";
export declare type SerializedEncoding = "hex" | "cb58" | "base58" | "base64" | "decimalString" | "number" | "utf8" | "display";
export declare abstract class Serializable {
    protected _typeName: string;
    protected _typeID: number;
    protected _codecID: number;
    /**
     * Used in serialization. TypeName is a string name for the type of object being output.
     */
    getTypeName(): string;
    /**
     * Used in serialization. Optional. TypeID is a number for the typeID of object being output.
     */
    getTypeID(): number;
    /**
     * Used in serialization. Optional. TypeID is a number for the typeID of object being output.
     */
    getCodecID(): number;
    /**
     * Sanitize to prevent cross scripting attacks.
     */
    sanitizeObject(obj: object): object;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
}
export declare class Serialization {
    private static instance;
    private constructor();
    private bintools;
    /**
     * Retrieves the Serialization singleton.
     */
    static getInstance(): Serialization;
    /**
     * Convert {@link https://github.com/feross/buffer|Buffer} to [[SerializedType]]
     *
     * @param vb {@link https://github.com/feross/buffer|Buffer}
     * @param type [[SerializedType]]
     * @param ...args remaining arguments
     * @returns type of [[SerializedType]]
     */
    bufferToType(vb: Buffer, type: SerializedType, ...args: any[]): any;
    /**
     * Convert [[SerializedType]] to {@link https://github.com/feross/buffer|Buffer}
     *
     * @param v type of [[SerializedType]]
     * @param type [[SerializedType]]
     * @param ...args remaining arguments
     * @returns {@link https://github.com/feross/buffer|Buffer}
     */
    typeToBuffer(v: any, type: SerializedType, ...args: any[]): Buffer;
    /**
     * Convert value to type of [[SerializedType]] or [[SerializedEncoding]]
     *
     * @param value
     * @param encoding [[SerializedEncoding]]
     * @param intype [[SerializedType]]
     * @param outtype [[SerializedType]]
     * @param ...args remaining arguments
     * @returns type of [[SerializedType]] or [[SerializedEncoding]]
     */
    encoder(value: any, encoding: SerializedEncoding, intype: SerializedType, outtype: SerializedType, ...args: any[]): any;
    /**
     * Convert value to type of [[SerializedType]] or [[SerializedEncoding]]
     *
     * @param value
     * @param encoding [[SerializedEncoding]]
     * @param intype [[SerializedType]]
     * @param outtype [[SerializedType]]
     * @param ...args remaining arguments
     * @returns type of [[SerializedType]] or [[SerializedEncoding]]
     */
    decoder(value: string, encoding: SerializedEncoding, intype: SerializedType, outtype: SerializedType, ...args: any[]): any;
    serialize(serialize: Serializable, vm: string, encoding?: SerializedEncoding, notes?: string): Serialized;
    deserialize(input: Serialized, output: Serializable): void;
}
//# sourceMappingURL=serialization.d.ts.map