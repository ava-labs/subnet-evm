/**
 * @packageDocumentation
 * @module Utils-HDNode
 */
import { Buffer } from "buffer/";
/**
 * BIP32 hierarchical deterministic keys.
 */
export default class HDNode {
    private hdkey;
    publicKey: Buffer;
    privateKey: Buffer;
    privateKeyCB58: string;
    chainCode: Buffer;
    privateExtendedKey: string;
    publicExtendedKey: string;
    /**
     * Derives the HDNode at path from the current HDNode.
     * @param path
     * @returns derived child HDNode
     */
    derive(path: string): HDNode;
    /**
     * Signs the buffer hash with the private key using secp256k1 and returns the signature as a buffer.
     * @param hash
     * @returns signature as a Buffer
     */
    sign(hash: Buffer): Buffer;
    /**
     * Verifies that the signature is valid for hash and the HDNode's public key using secp256k1.
     * @param hash
     * @param signature
     * @returns true for valid, false for invalid.
     * @throws if the hash or signature is the wrong length.
     */
    verify(hash: Buffer, signature: Buffer): boolean;
    /**
     * Wipes all record of the private key from the HDNode instance.
     * After calling this method, the instance will behave as if it was created via an xpub.
     */
    wipePrivateData(): void;
    /**
     * Creates an HDNode from a master seed or an extended public/private key
     * @param from seed or key to create HDNode from
     */
    constructor(from: string | Buffer);
}
//# sourceMappingURL=hdnode.d.ts.map