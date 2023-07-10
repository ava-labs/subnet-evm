/**
 * @packageDocumentation
 * @module Common-NBytes
 */
import { Buffer } from "buffer/";
import { Serializable, SerializedEncoding } from "../utils/serialization";
/**
 * Abstract class that implements basic functionality for managing a
 * {@link https://github.com/feross/buffer|Buffer} of an exact length.
 *
 * Create a class that extends this one and override bsize to make it validate for exactly
 * the correct length.
 */
export declare abstract class NBytes extends Serializable {
    protected _typeName: string;
    protected _typeID: any;
    serialize(encoding?: SerializedEncoding): object;
    deserialize(fields: object, encoding?: SerializedEncoding): void;
    protected bytes: Buffer;
    protected bsize: number;
    /**
     * Returns the length of the {@link https://github.com/feross/buffer|Buffer}.
     *
     * @returns The exact length requirement of this class
     */
    getSize: () => number;
    /**
     * Takes a base-58 encoded string, verifies its length, and stores it.
     *
     * @returns The size of the {@link https://github.com/feross/buffer|Buffer}
     */
    fromString(b58str: string): number;
    /**
     * Takes a [[Buffer]], verifies its length, and stores it.
     *
     * @returns The size of the {@link https://github.com/feross/buffer|Buffer}
     */
    fromBuffer(buff: Buffer, offset?: number): number;
    /**
     * @returns A reference to the stored {@link https://github.com/feross/buffer|Buffer}
     */
    toBuffer(): Buffer;
    /**
     * @returns A base-58 string of the stored {@link https://github.com/feross/buffer|Buffer}
     */
    toString(): string;
    abstract clone(): this;
    abstract create(...args: any[]): this;
}
//# sourceMappingURL=nbytes.d.ts.map