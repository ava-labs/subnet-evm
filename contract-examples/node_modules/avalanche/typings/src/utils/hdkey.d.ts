/**
 * @packageDocumentation
 * @module Utils-HDKey
 */
import { Buffer } from 'buffer/';
import { HDKeyJSON } from 'src/common';
import { default as hdkey } from 'hdkey';
/**
 * BIP32 hierarchical deterministic keys.
 *
 */
export default class HDKey {
    private static instance;
    private constructor();
    /**
     * Retrieves the HDKey singleton.
     */
    static getInstance(): HDKey;
    /**
     * Creates an HDNode from a master seed buffer
     *
     * @param seedBuffer Buffer
     *
     * @returns HDNode
     */
    fromMasterSeed(seed: Buffer): hdkey;
    /**
     * Creates an HDNode from a xprv or xpub extended key string. Accepts an optional versions object.
     *
     * @param xpriv string
     *
     * @returns HDNode
     */
    fromExtendedKey(xpriv: string): hdkey;
    /**
     * Creates an HDNode from an object created via hdkey.toJSON().
     *
     * @param obj HDKeyJSON
     *
     * @returns HDNode
     */
    fromJSON(obj: HDKeyJSON): hdkey;
}
//# sourceMappingURL=hdkey.d.ts.map