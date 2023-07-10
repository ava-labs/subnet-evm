/**
 * @packageDocumentation
 * @module Common-SECP256k1KeyChain
 */
import { Buffer } from "buffer/";
import * as elliptic from "elliptic";
import { StandardKeyPair, StandardKeyChain } from "./keychain";
/**
 * Class for representing a private and public keypair on the Platform Chain.
 */
export declare abstract class SECP256k1KeyPair extends StandardKeyPair {
    protected keypair: elliptic.ec.KeyPair;
    protected chainID: string;
    protected hrp: string;
    /**
     * @ignore
     */
    protected _sigFromSigBuffer(sig: Buffer): elliptic.ec.SignatureOptions;
    /**
     * Generates a new keypair.
     */
    generateKey(): void;
    /**
     * Imports a private key and generates the appropriate public key.
     *
     * @param privk A {@link https://github.com/feross/buffer|Buffer} representing the private key
     *
     * @returns true on success, false on failure
     */
    importKey(privk: Buffer): boolean;
    /**
     * Returns the address as a {@link https://github.com/feross/buffer|Buffer}.
     *
     * @returns A {@link https://github.com/feross/buffer|Buffer} representation of the address
     */
    getAddress(): Buffer;
    /**
     * Returns the address's string representation.
     *
     * @returns A string representation of the address
     */
    getAddressString(): string;
    /**
     * Returns an address given a public key.
     *
     * @param pubk A {@link https://github.com/feross/buffer|Buffer} representing the public key
     *
     * @returns A {@link https://github.com/feross/buffer|Buffer} for the address of the public key.
     */
    static addressFromPublicKey(pubk: Buffer): Buffer;
    /**
     * Returns a string representation of the private key.
     *
     * @returns A cb58 serialized string representation of the private key
     */
    getPrivateKeyString(): string;
    /**
     * Returns the public key.
     *
     * @returns A cb58 serialized string representation of the public key
     */
    getPublicKeyString(): string;
    /**
     * Takes a message, signs it, and returns the signature.
     *
     * @param msg The message to sign, be sure to hash first if expected
     *
     * @returns A {@link https://github.com/feross/buffer|Buffer} containing the signature
     */
    sign(msg: Buffer): Buffer;
    /**
     * Verifies that the private key associated with the provided public key produces the signature associated with the given message.
     *
     * @param msg The message associated with the signature
     * @param sig The signature of the signed message
     *
     * @returns True on success, false on failure
     */
    verify(msg: Buffer, sig: Buffer): boolean;
    /**
     * Recovers the public key of a message signer from a message and its associated signature.
     *
     * @param msg The message that's signed
     * @param sig The signature that's signed on the message
     *
     * @returns A {@link https://github.com/feross/buffer|Buffer} containing the public key of the signer
     */
    recover(msg: Buffer, sig: Buffer): Buffer;
    /**
     * Returns the chainID associated with this key.
     *
     * @returns The [[KeyPair]]'s chainID
     */
    getChainID(): string;
    /**
     * Sets the the chainID associated with this key.
     *
     * @param chainID String for the chainID
     */
    setChainID(chainID: string): void;
    /**
     * Returns the Human-Readable-Part of the network associated with this key.
     *
     * @returns The [[KeyPair]]'s Human-Readable-Part of the network's Bech32 addressing scheme
     */
    getHRP(): string;
    /**
     * Sets the the Human-Readable-Part of the network associated with this key.
     *
     * @param hrp String for the Human-Readable-Part of Bech32 addresses
     */
    setHRP(hrp: string): void;
    constructor(hrp: string, chainID: string);
}
/**
 * Class for representing a key chain in Avalanche.
 *
 * @typeparam SECP256k1KeyPair Class extending [[StandardKeyPair]] which is used as the key in [[SECP256k1KeyChain]]
 */
export declare abstract class SECP256k1KeyChain<SECPKPClass extends SECP256k1KeyPair> extends StandardKeyChain<SECPKPClass> {
    /**
     * Makes a new key pair, returns the address.
     *
     * @returns Address of the new key pair
     */
    makeKey: () => SECPKPClass;
    addKey(newKey: SECPKPClass): void;
    /**
     * Given a private key, makes a new key pair, returns the address.
     *
     * @param privk A {@link https://github.com/feross/buffer|Buffer} or cb58 serialized string representing the private key
     *
     * @returns Address of the new key pair
     */
    importKey: (privk: Buffer | string) => SECPKPClass;
}
//# sourceMappingURL=secp256k1.d.ts.map