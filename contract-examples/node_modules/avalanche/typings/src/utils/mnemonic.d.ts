/**
 * @packageDocumentation
 * @module Utils-Mnemonic
 */
import { Buffer } from "buffer/";
import { Wordlist } from "ethers";
/**
 * BIP39 Mnemonic code for generating deterministic keys.
 *
 */
export default class Mnemonic {
    private static instance;
    private constructor();
    protected wordlists: string[];
    /**
     * Retrieves the Mnemonic singleton.
     */
    static getInstance(): Mnemonic;
    /**
     * Return wordlists
     *
     * @param language a string specifying the language
     *
     * @returns A [[Wordlist]] object or array of strings
     */
    getWordlists(language?: string): string[] | Wordlist;
    /**
     * Synchronously takes mnemonic and password and returns {@link https://github.com/feross/buffer|Buffer}
     *
     * @param mnemonic the mnemonic as a string
     * @param password the password as a string
     *
     * @returns A {@link https://github.com/feross/buffer|Buffer}
     */
    mnemonicToSeedSync(mnemonic: string, password?: string): Buffer;
    /**
     * Asynchronously takes mnemonic and password and returns Promise {@link https://github.com/feross/buffer|Buffer}
     *
     * @param mnemonic the mnemonic as a string
     * @param password the password as a string
     *
     * @returns A {@link https://github.com/feross/buffer|Buffer}
     */
    mnemonicToSeed(mnemonic: string, password?: string): Promise<Buffer>;
    /**
     * Takes mnemonic and wordlist and returns buffer
     *
     * @param mnemonic the mnemonic as a string
     * @param wordlist Optional the wordlist as an array of strings
     *
     * @returns A string
     */
    mnemonicToEntropy(mnemonic: string, wordlist?: string[]): string;
    /**
     * Takes mnemonic and wordlist and returns buffer
     *
     * @param entropy the entropy as a {@link https://github.com/feross/buffer|Buffer} or as a string
     * @param wordlist Optional, the wordlist as an array of strings
     *
     * @returns A string
     */
    entropyToMnemonic(entropy: Buffer | string, wordlist?: string[]): string;
    /**
     * Validates a mnemonic
     11*
     * @param mnemonic the mnemonic as a string
     * @param wordlist Optional the wordlist as an array of strings
     *
     * @returns A string
     */
    validateMnemonic(mnemonic: string, wordlist?: string[]): string;
    /**
     * Sets the default word list
     *
     * @param language the language as a string
     *
     */
    setDefaultWordlist(language: string): void;
    /**
     * Returns the language of the default word list
     *
     * @returns A string
     */
    getDefaultWordlist(): string;
    /**
     * Generate a random mnemonic (uses crypto.randomBytes under the hood), defaults to 256-bits of entropy
     *
     * @param strength Optional the strength as a number
     * @param rng Optional the random number generator. Defaults to crypto.randomBytes
     * @param wordlist Optional
     *
     */
    generateMnemonic(strength?: number, rng?: (size: number) => Buffer, wordlist?: string[]): string;
}
//# sourceMappingURL=mnemonic.d.ts.map