/**
 * @packageDocumentation
 * @module Utils-BinTools
 */
import BN from "bn.js";
import { Buffer } from "buffer/";
/**
 * A class containing tools useful in interacting with binary data cross-platform using
 * nodejs & javascript.
 *
 * This class should never be instantiated directly. Instead,
 * invoke the "BinTools.getInstance()" static * function to grab the singleton
 * instance of the tools.
 *
 * Everything in this library uses
 * the {@link https://github.com/feross/buffer|feross's Buffer class}.
 *
 * ```js
 * const bintools: BinTools = BinTools.getInstance();
 * const b58str:  = bintools.bufferToB58(Buffer.from("Wubalubadubdub!"));
 * ```
 */
export default class BinTools {
    private static instance;
    private constructor();
    private b58;
    /**
     * Retrieves the BinTools singleton.
     */
    static getInstance(): BinTools;
    /**
     * Returns true if base64, otherwise false
     * @param str the string to verify is Base64
     */
    isBase64(str: string): boolean;
    /**
     * Returns true if cb58, otherwise false
     * @param cb58 the string to verify is cb58
     */
    isCB58(cb58: string): boolean;
    /**
     * Returns true if base58, otherwise false
     * @param base58 the string to verify is base58
     */
    isBase58(base58: string): boolean;
    /**
     * Returns true if hexidecimal, otherwise false
     * @param hex the string to verify is hexidecimal
     */
    isHex(hex: string): boolean;
    /**
     * Returns true if decimal, otherwise false
     * @param str the string to verify is hexidecimal
     */
    isDecimal(str: string): boolean;
    /**
     * Returns true if meets requirements to parse as an address as Bech32 on X-Chain or P-Chain, otherwise false
     * @param address the string to verify is address
     */
    isPrimaryBechAddress: (address: string) => boolean;
    /**
     * Produces a string from a {@link https://github.com/feross/buffer|Buffer}
     * representing a string. ONLY USED IN TRANSACTION FORMATTING, ASSUMED LENGTH IS PREPENDED.
     *
     * @param buff The {@link https://github.com/feross/buffer|Buffer} to convert to a string
     */
    bufferToString: (buff: Buffer) => string;
    /**
     * Produces a {@link https://github.com/feross/buffer|Buffer} from a string. ONLY USED IN TRANSACTION FORMATTING, LENGTH IS PREPENDED.
     *
     * @param str The string to convert to a {@link https://github.com/feross/buffer|Buffer}
     */
    stringToBuffer: (str: string) => Buffer;
    /**
     * Makes a copy (no reference) of a {@link https://github.com/feross/buffer|Buffer}
     * over provided indecies.
     *
     * @param buff The {@link https://github.com/feross/buffer|Buffer} to copy
     * @param start The index to start the copy
     * @param end The index to end the copy
     */
    copyFrom: (buff: Buffer, start?: number, end?: number) => Buffer;
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} and returns a base-58 string of
     * the {@link https://github.com/feross/buffer|Buffer}.
     *
     * @param buff The {@link https://github.com/feross/buffer|Buffer} to convert to base-58
     */
    bufferToB58: (buff: Buffer) => string;
    /**
     * Takes a base-58 string and returns a {@link https://github.com/feross/buffer|Buffer}.
     *
     * @param b58str The base-58 string to convert
     * to a {@link https://github.com/feross/buffer|Buffer}
     */
    b58ToBuffer: (b58str: string) => Buffer;
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} and returns an ArrayBuffer.
     *
     * @param buff The {@link https://github.com/feross/buffer|Buffer} to
     * convert to an ArrayBuffer
     */
    fromBufferToArrayBuffer: (buff: Buffer) => ArrayBuffer;
    /**
     * Takes an ArrayBuffer and converts it to a {@link https://github.com/feross/buffer|Buffer}.
     *
     * @param ab The ArrayBuffer to convert to a {@link https://github.com/feross/buffer|Buffer}
     */
    fromArrayBufferToBuffer: (ab: ArrayBuffer) => Buffer;
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} and converts it
     * to a {@link https://github.com/indutny/bn.js/|BN}.
     *
     * @param buff The {@link https://github.com/feross/buffer|Buffer} to convert
     * to a {@link https://github.com/indutny/bn.js/|BN}
     */
    fromBufferToBN: (buff: Buffer) => BN;
    /**
     * Takes a {@link https://github.com/indutny/bn.js/|BN} and converts it
     * to a {@link https://github.com/feross/buffer|Buffer}.
     *
     * @param bn The {@link https://github.com/indutny/bn.js/|BN} to convert
     * to a {@link https://github.com/feross/buffer|Buffer}
     * @param length The zero-padded length of the {@link https://github.com/feross/buffer|Buffer}
     */
    fromBNToBuffer: (bn: BN, length?: number) => Buffer;
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} and adds a checksum, returning
     * a {@link https://github.com/feross/buffer|Buffer} with the 4-byte checksum appended.
     *
     * @param buff The {@link https://github.com/feross/buffer|Buffer} to append a checksum
     */
    addChecksum: (buff: Buffer) => Buffer;
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} with an appended 4-byte checksum
     * and returns true if the checksum is valid, otherwise false.
     *
     * @param b The {@link https://github.com/feross/buffer|Buffer} to validate the checksum
     */
    validateChecksum: (buff: Buffer) => boolean;
    /**
     * Takes a {@link https://github.com/feross/buffer|Buffer} and returns a base-58 string with
     * checksum as per the cb58 standard.
     *
     * @param bytes A {@link https://github.com/feross/buffer|Buffer} to serialize
     *
     * @returns A serialized base-58 string of the Buffer.
     */
    cb58Encode: (bytes: Buffer) => string;
    /**
     * Takes a cb58 serialized {@link https://github.com/feross/buffer|Buffer} or base-58 string
     * and returns a {@link https://github.com/feross/buffer|Buffer} of the original data. Throws on error.
     *
     * @param bytes A cb58 serialized {@link https://github.com/feross/buffer|Buffer} or base-58 string
     */
    cb58Decode: (bytes: Buffer | string) => Buffer;
    cb58DecodeWithChecksum: (bytes: Buffer | string) => string;
    addressToString: (hrp: string, chainid: string, bytes: Buffer) => string;
    stringToAddress: (address: string, hrp?: string) => Buffer;
    /**
     * Takes an address and returns its {@link https://github.com/feross/buffer|Buffer}
     * representation if valid. A more strict version of stringToAddress.
     *
     * @param addr A string representation of the address
     * @param blockchainID A cb58 encoded string representation of the blockchainID
     * @param alias A chainID alias, if any, that the address can also parse from.
     * @param addrlen VMs can use any addressing scheme that they like, so this is the appropriate number of address bytes. Default 20.
     *
     * @returns A {@link https://github.com/feross/buffer|Buffer} for the address if valid,
     * undefined if not valid.
     */
    parseAddress: (addr: string, blockchainID: string, alias?: string, addrlen?: number) => Buffer;
}
//# sourceMappingURL=bintools.d.ts.map