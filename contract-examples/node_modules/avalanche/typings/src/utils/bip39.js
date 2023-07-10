"use strict";
/**
 * @packageDocumentation
 * @module Utils-BIP39
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.BIP39 = void 0;
const errors_1 = require("./errors");
const bip39 = require('bip39');
const randomBytes = require("randombytes");
/**
 * Implementation of Mnemonic. Mnemonic code for generating deterministic keys.
 *
 */
class BIP39 {
    constructor() {
        this.wordlists = bip39.wordlists;
    }
    /**
     * Retrieves the Mnemonic singleton.
     */
    static getInstance() {
        if (!BIP39.instance) {
            BIP39.instance = new BIP39();
        }
        return BIP39.instance;
    }
    /**
     * Return wordlists
     *
     * @param language a string specifying the language
     *
     * @returns A [[Wordlist]] object or array of strings
     */
    getWordlists(language) {
        if (language !== undefined) {
            return this.wordlists[language];
        }
        else {
            return this.wordlists;
        }
    }
    /**
     * Synchronously takes mnemonic and password and returns {@link https://github.com/feross/buffer|Buffer}
     *
     * @param mnemonic the mnemonic as a string
     * @param password the password as a string
     *
     * @returns A {@link https://github.com/feross/buffer|Buffer}
     */
    mnemonicToSeedSync(mnemonic, password) {
        return bip39.mnemonicToSeedSync(mnemonic, password);
    }
    /**
     * Asynchronously takes mnemonic and password and returns Promise<{@link https://github.com/feross/buffer|Buffer}>
     *
     * @param mnemonic the mnemonic as a string
     * @param password the password as a string
     *
     * @returns A {@link https://github.com/feross/buffer|Buffer}
     */
    mnemonicToSeed(mnemonic, password) {
        return bip39.mnemonicToSeed(mnemonic, password);
    }
    /**
     * Takes mnemonic and wordlist and returns buffer
     *
     * @param mnemonic the mnemonic as a string
     * @param wordlist Optional the wordlist as an array of strings
     *
     * @returns A string
     */
    mnemonicToEntropy(mnemonic, wordlist) {
        return bip39.mnemonicToEntropy(mnemonic, wordlist);
    }
    /**
     * Takes mnemonic and wordlist and returns buffer
     *
     * @param entropy the entropy as a {@link https://github.com/feross/buffer|Buffer} or as a string
     * @param wordlist Optional, the wordlist as an array of strings
     *
     * @returns A string
     */
    entropyToMnemonic(entropy, wordlist) {
        return bip39.entropyToMnemonic(entropy, wordlist);
    }
    /**
     * Validates a mnemonic
     11*
     * @param mnemonic the mnemonic as a string
     * @param wordlist Optional the wordlist as an array of strings
     *
     * @returns A string
     */
    validateMnemonic(mnemonic, wordlist) {
        return bip39.validateMnemonic(mnemonic, wordlist);
    }
    /**
     * Sets the default word list
     *
     * @param language the language as a string
     *
     * @returns A string
     */
    setDefaultWordlist(language) {
        return bip39.setDefaultWordlist(language);
    }
    /**
     * Returns the language of the default word list
     *
     * @returns A string
     */
    getDefaultWordlist() {
        return bip39.getDefaultWordlist();
    }
    /**
     * Generate a random mnemonic (uses crypto.randomBytes under the hood), defaults to 256-bits of entropy
     *
     * @param strength Optional the strength as a number
     * @param rng Optional the random number generator. Defaults to crypto.randomBytes
     * @param wordlist Optional
     *
     */
    generateMnemonic(strength, rng, wordlist) {
        strength = strength || 256;
        if (strength % 32 !== 0) {
            throw new errors_1.InvalidEntropy('Error - Invalid entropy');
        }
        rng = rng || randomBytes;
        return bip39.generateMnemonic(strength, rng, wordlist);
    }
}
exports.BIP39 = BIP39;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYmlwMzkuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi9zcmMvdXRpbHMvYmlwMzkudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IjtBQUFBOzs7R0FHRzs7O0FBSUgscUNBQXlDO0FBQ3pDLE1BQU0sS0FBSyxHQUFRLE9BQU8sQ0FBQyxPQUFPLENBQUMsQ0FBQTtBQUNuQyxNQUFNLFdBQVcsR0FBUSxPQUFPLENBQUMsYUFBYSxDQUFDLENBQUE7QUFFL0M7OztHQUdHO0FBQ0gsTUFBYSxLQUFLO0lBRWhCO1FBQ1UsY0FBUyxHQUFhLEtBQUssQ0FBQyxTQUFTLENBQUE7SUFEdkIsQ0FBQztJQUd6Qjs7T0FFRztJQUNILE1BQU0sQ0FBQyxXQUFXO1FBQ2hCLElBQUksQ0FBQyxLQUFLLENBQUMsUUFBUSxFQUFFO1lBQ25CLEtBQUssQ0FBQyxRQUFRLEdBQUcsSUFBSSxLQUFLLEVBQUUsQ0FBQTtTQUM3QjtRQUNELE9BQU8sS0FBSyxDQUFDLFFBQVEsQ0FBQTtJQUN2QixDQUFDO0lBRUQ7Ozs7OztPQU1HO0lBQ0gsWUFBWSxDQUFDLFFBQWlCO1FBQzVCLElBQUksUUFBUSxLQUFLLFNBQVMsRUFBRTtZQUMxQixPQUFPLElBQUksQ0FBQyxTQUFTLENBQUMsUUFBUSxDQUFDLENBQUE7U0FDaEM7YUFBTTtZQUNMLE9BQU8sSUFBSSxDQUFDLFNBQVMsQ0FBQTtTQUN0QjtJQUNILENBQUM7SUFFRDs7Ozs7OztPQU9HO0lBQ0gsa0JBQWtCLENBQUMsUUFBZ0IsRUFBRSxRQUFnQjtRQUNuRCxPQUFPLEtBQUssQ0FBQyxrQkFBa0IsQ0FBQyxRQUFRLEVBQUUsUUFBUSxDQUFDLENBQUE7SUFDckQsQ0FBQztJQUVEOzs7Ozs7O09BT0c7SUFDSCxjQUFjLENBQUMsUUFBZ0IsRUFBRSxRQUFnQjtRQUMvQyxPQUFPLEtBQUssQ0FBQyxjQUFjLENBQUMsUUFBUSxFQUFFLFFBQVEsQ0FBQyxDQUFBO0lBQ2pELENBQUM7SUFFRDs7Ozs7OztPQU9HO0lBQ0gsaUJBQWlCLENBQ2YsUUFBZ0IsRUFDaEIsUUFBbUI7UUFFbkIsT0FBTyxLQUFLLENBQUMsaUJBQWlCLENBQUMsUUFBUSxFQUFFLFFBQVEsQ0FBQyxDQUFBO0lBQ3BELENBQUM7SUFFRDs7Ozs7OztPQU9HO0lBQ0gsaUJBQWlCLENBQ2YsT0FBd0IsRUFDeEIsUUFBbUI7UUFFbkIsT0FBTyxLQUFLLENBQUMsaUJBQWlCLENBQUMsT0FBTyxFQUFFLFFBQVEsQ0FBQyxDQUFBO0lBQ25ELENBQUM7SUFFRDs7Ozs7OztPQU9HO0lBQ0gsZ0JBQWdCLENBQ2QsUUFBZ0IsRUFDaEIsUUFBbUI7UUFFbkIsT0FBTyxLQUFLLENBQUMsZ0JBQWdCLENBQUMsUUFBUSxFQUFFLFFBQVEsQ0FBQyxDQUFBO0lBQ25ELENBQUM7SUFFRDs7Ozs7O09BTUc7SUFDSCxrQkFBa0IsQ0FBQyxRQUFnQjtRQUNqQyxPQUFPLEtBQUssQ0FBQyxrQkFBa0IsQ0FBQyxRQUFRLENBQUMsQ0FBQTtJQUMzQyxDQUFDO0lBRUQ7Ozs7T0FJRztJQUNILGtCQUFrQjtRQUNoQixPQUFPLEtBQUssQ0FBQyxrQkFBa0IsRUFBRSxDQUFBO0lBQ25DLENBQUM7SUFFRDs7Ozs7OztPQU9HO0lBQ0gsZ0JBQWdCLENBQ2QsUUFBaUIsRUFDakIsR0FBOEIsRUFDOUIsUUFBbUI7UUFFbkIsUUFBUSxHQUFHLFFBQVEsSUFBSSxHQUFHLENBQUE7UUFDMUIsSUFBSSxRQUFRLEdBQUcsRUFBRSxLQUFLLENBQUMsRUFBRTtZQUN2QixNQUFNLElBQUksdUJBQWMsQ0FBQyx5QkFBeUIsQ0FBQyxDQUFBO1NBQ3BEO1FBQ0QsR0FBRyxHQUFHLEdBQUcsSUFBSSxXQUFXLENBQUE7UUFDeEIsT0FBTyxLQUFLLENBQUMsZ0JBQWdCLENBQUMsUUFBUSxFQUFFLEdBQUcsRUFBRSxRQUFRLENBQUMsQ0FBQTtJQUN4RCxDQUFDO0NBQ0Y7QUEzSUQsc0JBMklDIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgVXRpbHMtQklQMzlcbiAqL1xuXG5pbXBvcnQgeyBCdWZmZXIgfSBmcm9tICdidWZmZXIvJ1xuaW1wb3J0IHsgV29yZGxpc3QgfSBmcm9tICdldGhlcnMnXG5pbXBvcnQgeyBJbnZhbGlkRW50cm9weSB9IGZyb20gJy4vZXJyb3JzJ1xuY29uc3QgYmlwMzk6IGFueSA9IHJlcXVpcmUoJ2JpcDM5JylcbmNvbnN0IHJhbmRvbUJ5dGVzOiBhbnkgPSByZXF1aXJlKFwicmFuZG9tYnl0ZXNcIilcblxuLyoqXG4gKiBJbXBsZW1lbnRhdGlvbiBvZiBNbmVtb25pYy4gTW5lbW9uaWMgY29kZSBmb3IgZ2VuZXJhdGluZyBkZXRlcm1pbmlzdGljIGtleXMuXG4gKlxuICovXG5leHBvcnQgY2xhc3MgQklQMzkge1xuICBwcml2YXRlIHN0YXRpYyBpbnN0YW5jZTogQklQMzlcbiAgcHJpdmF0ZSBjb25zdHJ1Y3RvcigpIHsgfVxuICBwcm90ZWN0ZWQgd29yZGxpc3RzOiBzdHJpbmdbXSA9IGJpcDM5LndvcmRsaXN0c1xuXG4gIC8qKlxuICAgKiBSZXRyaWV2ZXMgdGhlIE1uZW1vbmljIHNpbmdsZXRvbi5cbiAgICovXG4gIHN0YXRpYyBnZXRJbnN0YW5jZSgpOiBCSVAzOSB7XG4gICAgaWYgKCFCSVAzOS5pbnN0YW5jZSkge1xuICAgICAgQklQMzkuaW5zdGFuY2UgPSBuZXcgQklQMzkoKVxuICAgIH1cbiAgICByZXR1cm4gQklQMzkuaW5zdGFuY2VcbiAgfVxuXG4gIC8qKlxuICAgKiBSZXR1cm4gd29yZGxpc3RzXG4gICAqXG4gICAqIEBwYXJhbSBsYW5ndWFnZSBhIHN0cmluZyBzcGVjaWZ5aW5nIHRoZSBsYW5ndWFnZVxuICAgKlxuICAgKiBAcmV0dXJucyBBIFtbV29yZGxpc3RdXSBvYmplY3Qgb3IgYXJyYXkgb2Ygc3RyaW5nc1xuICAgKi9cbiAgZ2V0V29yZGxpc3RzKGxhbmd1YWdlPzogc3RyaW5nKTogc3RyaW5nW10gfCBXb3JkbGlzdCB7XG4gICAgaWYgKGxhbmd1YWdlICE9PSB1bmRlZmluZWQpIHtcbiAgICAgIHJldHVybiB0aGlzLndvcmRsaXN0c1tsYW5ndWFnZV1cbiAgICB9IGVsc2Uge1xuICAgICAgcmV0dXJuIHRoaXMud29yZGxpc3RzXG4gICAgfVxuICB9XG5cbiAgLyoqXG4gICAqIFN5bmNocm9ub3VzbHkgdGFrZXMgbW5lbW9uaWMgYW5kIHBhc3N3b3JkIGFuZCByZXR1cm5zIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9XG4gICAqXG4gICAqIEBwYXJhbSBtbmVtb25pYyB0aGUgbW5lbW9uaWMgYXMgYSBzdHJpbmdcbiAgICogQHBhcmFtIHBhc3N3b3JkIHRoZSBwYXNzd29yZCBhcyBhIHN0cmluZ1xuICAgKlxuICAgKiBAcmV0dXJucyBBIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9XG4gICAqL1xuICBtbmVtb25pY1RvU2VlZFN5bmMobW5lbW9uaWM6IHN0cmluZywgcGFzc3dvcmQ6IHN0cmluZyk6IEJ1ZmZlciB7XG4gICAgcmV0dXJuIGJpcDM5Lm1uZW1vbmljVG9TZWVkU3luYyhtbmVtb25pYywgcGFzc3dvcmQpXG4gIH1cblxuICAvKipcbiAgICogQXN5bmNocm9ub3VzbHkgdGFrZXMgbW5lbW9uaWMgYW5kIHBhc3N3b3JkIGFuZCByZXR1cm5zIFByb21pc2U8e0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0+XG4gICAqXG4gICAqIEBwYXJhbSBtbmVtb25pYyB0aGUgbW5lbW9uaWMgYXMgYSBzdHJpbmdcbiAgICogQHBhcmFtIHBhc3N3b3JkIHRoZSBwYXNzd29yZCBhcyBhIHN0cmluZ1xuICAgKlxuICAgKiBAcmV0dXJucyBBIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9XG4gICAqL1xuICBtbmVtb25pY1RvU2VlZChtbmVtb25pYzogc3RyaW5nLCBwYXNzd29yZDogc3RyaW5nKTogQnVmZmVyIHtcbiAgICByZXR1cm4gYmlwMzkubW5lbW9uaWNUb1NlZWQobW5lbW9uaWMsIHBhc3N3b3JkKVxuICB9XG5cbiAgLyoqXG4gICAqIFRha2VzIG1uZW1vbmljIGFuZCB3b3JkbGlzdCBhbmQgcmV0dXJucyBidWZmZXJcbiAgICpcbiAgICogQHBhcmFtIG1uZW1vbmljIHRoZSBtbmVtb25pYyBhcyBhIHN0cmluZ1xuICAgKiBAcGFyYW0gd29yZGxpc3QgT3B0aW9uYWwgdGhlIHdvcmRsaXN0IGFzIGFuIGFycmF5IG9mIHN0cmluZ3NcbiAgICpcbiAgICogQHJldHVybnMgQSBzdHJpbmdcbiAgICovXG4gIG1uZW1vbmljVG9FbnRyb3B5KFxuICAgIG1uZW1vbmljOiBzdHJpbmcsXG4gICAgd29yZGxpc3Q/OiBzdHJpbmdbXVxuICApOiBzdHJpbmcge1xuICAgIHJldHVybiBiaXAzOS5tbmVtb25pY1RvRW50cm9weShtbmVtb25pYywgd29yZGxpc3QpXG4gIH1cblxuICAvKipcbiAgICogVGFrZXMgbW5lbW9uaWMgYW5kIHdvcmRsaXN0IGFuZCByZXR1cm5zIGJ1ZmZlclxuICAgKlxuICAgKiBAcGFyYW0gZW50cm9weSB0aGUgZW50cm9weSBhcyBhIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IG9yIGFzIGEgc3RyaW5nXG4gICAqIEBwYXJhbSB3b3JkbGlzdCBPcHRpb25hbCwgdGhlIHdvcmRsaXN0IGFzIGFuIGFycmF5IG9mIHN0cmluZ3NcbiAgICpcbiAgICogQHJldHVybnMgQSBzdHJpbmdcbiAgICovXG4gIGVudHJvcHlUb01uZW1vbmljKFxuICAgIGVudHJvcHk6IEJ1ZmZlciB8IHN0cmluZyxcbiAgICB3b3JkbGlzdD86IHN0cmluZ1tdXG4gICk6IHN0cmluZyB7XG4gICAgcmV0dXJuIGJpcDM5LmVudHJvcHlUb01uZW1vbmljKGVudHJvcHksIHdvcmRsaXN0KVxuICB9XG5cbiAgLyoqXG4gICAqIFZhbGlkYXRlcyBhIG1uZW1vbmljXG4gICAxMSpcbiAgICogQHBhcmFtIG1uZW1vbmljIHRoZSBtbmVtb25pYyBhcyBhIHN0cmluZ1xuICAgKiBAcGFyYW0gd29yZGxpc3QgT3B0aW9uYWwgdGhlIHdvcmRsaXN0IGFzIGFuIGFycmF5IG9mIHN0cmluZ3NcbiAgICpcbiAgICogQHJldHVybnMgQSBzdHJpbmdcbiAgICovXG4gIHZhbGlkYXRlTW5lbW9uaWMoXG4gICAgbW5lbW9uaWM6IHN0cmluZyxcbiAgICB3b3JkbGlzdD86IHN0cmluZ1tdXG4gICk6IHN0cmluZyB7XG4gICAgcmV0dXJuIGJpcDM5LnZhbGlkYXRlTW5lbW9uaWMobW5lbW9uaWMsIHdvcmRsaXN0KVxuICB9XG5cbiAgLyoqXG4gICAqIFNldHMgdGhlIGRlZmF1bHQgd29yZCBsaXN0XG4gICAqXG4gICAqIEBwYXJhbSBsYW5ndWFnZSB0aGUgbGFuZ3VhZ2UgYXMgYSBzdHJpbmdcbiAgICpcbiAgICogQHJldHVybnMgQSBzdHJpbmdcbiAgICovXG4gIHNldERlZmF1bHRXb3JkbGlzdChsYW5ndWFnZTogc3RyaW5nKTogc3RyaW5nIHtcbiAgICByZXR1cm4gYmlwMzkuc2V0RGVmYXVsdFdvcmRsaXN0KGxhbmd1YWdlKVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIGxhbmd1YWdlIG9mIHRoZSBkZWZhdWx0IHdvcmQgbGlzdFxuICAgKiBcbiAgICogQHJldHVybnMgQSBzdHJpbmdcbiAgICovXG4gIGdldERlZmF1bHRXb3JkbGlzdCgpOiBzdHJpbmcge1xuICAgIHJldHVybiBiaXAzOS5nZXREZWZhdWx0V29yZGxpc3QoKVxuICB9XG5cbiAgLyoqXG4gICAqIEdlbmVyYXRlIGEgcmFuZG9tIG1uZW1vbmljICh1c2VzIGNyeXB0by5yYW5kb21CeXRlcyB1bmRlciB0aGUgaG9vZCksIGRlZmF1bHRzIHRvIDI1Ni1iaXRzIG9mIGVudHJvcHlcbiAgICogXG4gICAqIEBwYXJhbSBzdHJlbmd0aCBPcHRpb25hbCB0aGUgc3RyZW5ndGggYXMgYSBudW1iZXJcbiAgICogQHBhcmFtIHJuZyBPcHRpb25hbCB0aGUgcmFuZG9tIG51bWJlciBnZW5lcmF0b3IuIERlZmF1bHRzIHRvIGNyeXB0by5yYW5kb21CeXRlc1xuICAgKiBAcGFyYW0gd29yZGxpc3QgT3B0aW9uYWxcbiAgICogXG4gICAqL1xuICBnZW5lcmF0ZU1uZW1vbmljKFxuICAgIHN0cmVuZ3RoPzogbnVtYmVyLFxuICAgIHJuZz86IChzaXplOiBudW1iZXIpID0+IEJ1ZmZlcixcbiAgICB3b3JkbGlzdD86IHN0cmluZ1tdLFxuICApOiBzdHJpbmcge1xuICAgIHN0cmVuZ3RoID0gc3RyZW5ndGggfHwgMjU2XG4gICAgaWYgKHN0cmVuZ3RoICUgMzIgIT09IDApIHtcbiAgICAgIHRocm93IG5ldyBJbnZhbGlkRW50cm9weSgnRXJyb3IgLSBJbnZhbGlkIGVudHJvcHknKVxuICAgIH1cbiAgICBybmcgPSBybmcgfHwgcmFuZG9tQnl0ZXNcbiAgICByZXR1cm4gYmlwMzkuZ2VuZXJhdGVNbmVtb25pYyhzdHJlbmd0aCwgcm5nLCB3b3JkbGlzdClcbiAgfVxufSJdfQ==