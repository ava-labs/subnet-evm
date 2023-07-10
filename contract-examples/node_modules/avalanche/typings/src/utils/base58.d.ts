/**
 * @packageDocumentation
 * @module Utils-Base58
 */
import BN from "bn.js";
import { Buffer } from "buffer/";
/**
 * A Base58 class that uses the cross-platform Buffer module. Built so that Typescript
 * will accept the code.
 *
 * ```js
 * let b58:Base58 = new Base58();
 * let str:string = b58.encode(somebuffer);
 * let buff:Buffer = b58.decode(somestring);
 * ```
 */
export declare class Base58 {
    private static instance;
    private constructor();
    /**
     * Retrieves the Base58 singleton.
     */
    static getInstance(): Base58;
    protected b58alphabet: string;
    protected alphabetIdx0: string;
    protected b58: number[];
    protected big58Radix: BN;
    protected bigZero: BN;
    /**
     * Encodes a {@link https://github.com/feross/buffer|Buffer} as a base-58 string
     *
     * @param buff A {@link https://github.com/feross/buffer|Buffer} to encode
     *
     * @returns A base-58 string.
     */
    encode: (buff: Buffer) => string;
    /**
     * Decodes a base-58 into a {@link https://github.com/feross/buffer|Buffer}
     *
     * @param b A base-58 string to decode
     *
     * @returns A {@link https://github.com/feross/buffer|Buffer} from the decoded string.
     */
    decode: (b: string) => Buffer;
}
//# sourceMappingURL=base58.d.ts.map