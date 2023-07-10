"use strict";
/**
 * @packageDocumentation
 * @module Utils-Mnemonic
 */
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
const buffer_1 = require("buffer/");
const errors_1 = require("./errors");
const bip39 = require("bip39");
const randomBytes = require("randombytes");
/**
 * BIP39 Mnemonic code for generating deterministic keys.
 *
 */
class Mnemonic {
    constructor() {
        this.wordlists = bip39.wordlists;
    }
    /**
     * Retrieves the Mnemonic singleton.
     */
    static getInstance() {
        if (!Mnemonic.instance) {
            Mnemonic.instance = new Mnemonic();
        }
        return Mnemonic.instance;
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
            return this.wordlists[`${language}`];
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
    mnemonicToSeedSync(mnemonic, password = "") {
        const seed = bip39.mnemonicToSeedSync(mnemonic, password);
        return buffer_1.Buffer.from(seed);
    }
    /**
     * Asynchronously takes mnemonic and password and returns Promise {@link https://github.com/feross/buffer|Buffer}
     *
     * @param mnemonic the mnemonic as a string
     * @param password the password as a string
     *
     * @returns A {@link https://github.com/feross/buffer|Buffer}
     */
    mnemonicToSeed(mnemonic, password = "") {
        return __awaiter(this, void 0, void 0, function* () {
            const seed = yield bip39.mnemonicToSeed(mnemonic, password);
            return buffer_1.Buffer.from(seed);
        });
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
     */
    setDefaultWordlist(language) {
        bip39.setDefaultWordlist(language);
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
            throw new errors_1.InvalidEntropy("Error - Invalid entropy");
        }
        rng = rng || randomBytes;
        return bip39.generateMnemonic(strength, rng, wordlist);
    }
}
exports.default = Mnemonic;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibW5lbW9uaWMuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi9zcmMvdXRpbHMvbW5lbW9uaWMudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IjtBQUFBOzs7R0FHRzs7Ozs7Ozs7Ozs7QUFFSCxvQ0FBZ0M7QUFFaEMscUNBQXlDO0FBQ3pDLE1BQU0sS0FBSyxHQUFRLE9BQU8sQ0FBQyxPQUFPLENBQUMsQ0FBQTtBQUNuQyxNQUFNLFdBQVcsR0FBUSxPQUFPLENBQUMsYUFBYSxDQUFDLENBQUE7QUFFL0M7OztHQUdHO0FBQ0gsTUFBcUIsUUFBUTtJQUUzQjtRQUNVLGNBQVMsR0FBYSxLQUFLLENBQUMsU0FBUyxDQUFBO0lBRHhCLENBQUM7SUFHeEI7O09BRUc7SUFDSCxNQUFNLENBQUMsV0FBVztRQUNoQixJQUFJLENBQUMsUUFBUSxDQUFDLFFBQVEsRUFBRTtZQUN0QixRQUFRLENBQUMsUUFBUSxHQUFHLElBQUksUUFBUSxFQUFFLENBQUE7U0FDbkM7UUFDRCxPQUFPLFFBQVEsQ0FBQyxRQUFRLENBQUE7SUFDMUIsQ0FBQztJQUVEOzs7Ozs7T0FNRztJQUNILFlBQVksQ0FBQyxRQUFpQjtRQUM1QixJQUFJLFFBQVEsS0FBSyxTQUFTLEVBQUU7WUFDMUIsT0FBTyxJQUFJLENBQUMsU0FBUyxDQUFDLEdBQUcsUUFBUSxFQUFFLENBQUMsQ0FBQTtTQUNyQzthQUFNO1lBQ0wsT0FBTyxJQUFJLENBQUMsU0FBUyxDQUFBO1NBQ3RCO0lBQ0gsQ0FBQztJQUVEOzs7Ozs7O09BT0c7SUFDSCxrQkFBa0IsQ0FBQyxRQUFnQixFQUFFLFdBQW1CLEVBQUU7UUFDeEQsTUFBTSxJQUFJLEdBQVcsS0FBSyxDQUFDLGtCQUFrQixDQUFDLFFBQVEsRUFBRSxRQUFRLENBQUMsQ0FBQTtRQUNqRSxPQUFPLGVBQU0sQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7SUFDMUIsQ0FBQztJQUVEOzs7Ozs7O09BT0c7SUFDRyxjQUFjLENBQ2xCLFFBQWdCLEVBQ2hCLFdBQW1CLEVBQUU7O1lBRXJCLE1BQU0sSUFBSSxHQUFXLE1BQU0sS0FBSyxDQUFDLGNBQWMsQ0FBQyxRQUFRLEVBQUUsUUFBUSxDQUFDLENBQUE7WUFDbkUsT0FBTyxlQUFNLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO1FBQzFCLENBQUM7S0FBQTtJQUVEOzs7Ozs7O09BT0c7SUFDSCxpQkFBaUIsQ0FBQyxRQUFnQixFQUFFLFFBQW1CO1FBQ3JELE9BQU8sS0FBSyxDQUFDLGlCQUFpQixDQUFDLFFBQVEsRUFBRSxRQUFRLENBQUMsQ0FBQTtJQUNwRCxDQUFDO0lBRUQ7Ozs7Ozs7T0FPRztJQUNILGlCQUFpQixDQUFDLE9BQXdCLEVBQUUsUUFBbUI7UUFDN0QsT0FBTyxLQUFLLENBQUMsaUJBQWlCLENBQUMsT0FBTyxFQUFFLFFBQVEsQ0FBQyxDQUFBO0lBQ25ELENBQUM7SUFFRDs7Ozs7OztPQU9HO0lBQ0gsZ0JBQWdCLENBQUMsUUFBZ0IsRUFBRSxRQUFtQjtRQUNwRCxPQUFPLEtBQUssQ0FBQyxnQkFBZ0IsQ0FBQyxRQUFRLEVBQUUsUUFBUSxDQUFDLENBQUE7SUFDbkQsQ0FBQztJQUVEOzs7OztPQUtHO0lBQ0gsa0JBQWtCLENBQUMsUUFBZ0I7UUFDakMsS0FBSyxDQUFDLGtCQUFrQixDQUFDLFFBQVEsQ0FBQyxDQUFBO0lBQ3BDLENBQUM7SUFFRDs7OztPQUlHO0lBQ0gsa0JBQWtCO1FBQ2hCLE9BQU8sS0FBSyxDQUFDLGtCQUFrQixFQUFFLENBQUE7SUFDbkMsQ0FBQztJQUVEOzs7Ozs7O09BT0c7SUFDSCxnQkFBZ0IsQ0FDZCxRQUFpQixFQUNqQixHQUE4QixFQUM5QixRQUFtQjtRQUVuQixRQUFRLEdBQUcsUUFBUSxJQUFJLEdBQUcsQ0FBQTtRQUMxQixJQUFJLFFBQVEsR0FBRyxFQUFFLEtBQUssQ0FBQyxFQUFFO1lBQ3ZCLE1BQU0sSUFBSSx1QkFBYyxDQUFDLHlCQUF5QixDQUFDLENBQUE7U0FDcEQ7UUFDRCxHQUFHLEdBQUcsR0FBRyxJQUFJLFdBQVcsQ0FBQTtRQUN4QixPQUFPLEtBQUssQ0FBQyxnQkFBZ0IsQ0FBQyxRQUFRLEVBQUUsR0FBRyxFQUFFLFFBQVEsQ0FBQyxDQUFBO0lBQ3hELENBQUM7Q0FDRjtBQXRJRCwyQkFzSUMiLCJzb3VyY2VzQ29udGVudCI6WyIvKipcbiAqIEBwYWNrYWdlRG9jdW1lbnRhdGlvblxuICogQG1vZHVsZSBVdGlscy1NbmVtb25pY1xuICovXG5cbmltcG9ydCB7IEJ1ZmZlciB9IGZyb20gXCJidWZmZXIvXCJcbmltcG9ydCB7IFdvcmRsaXN0IH0gZnJvbSBcImV0aGVyc1wiXG5pbXBvcnQgeyBJbnZhbGlkRW50cm9weSB9IGZyb20gXCIuL2Vycm9yc1wiXG5jb25zdCBiaXAzOTogYW55ID0gcmVxdWlyZShcImJpcDM5XCIpXG5jb25zdCByYW5kb21CeXRlczogYW55ID0gcmVxdWlyZShcInJhbmRvbWJ5dGVzXCIpXG5cbi8qKlxuICogQklQMzkgTW5lbW9uaWMgY29kZSBmb3IgZ2VuZXJhdGluZyBkZXRlcm1pbmlzdGljIGtleXMuXG4gKlxuICovXG5leHBvcnQgZGVmYXVsdCBjbGFzcyBNbmVtb25pYyB7XG4gIHByaXZhdGUgc3RhdGljIGluc3RhbmNlOiBNbmVtb25pY1xuICBwcml2YXRlIGNvbnN0cnVjdG9yKCkge31cbiAgcHJvdGVjdGVkIHdvcmRsaXN0czogc3RyaW5nW10gPSBiaXAzOS53b3JkbGlzdHNcblxuICAvKipcbiAgICogUmV0cmlldmVzIHRoZSBNbmVtb25pYyBzaW5nbGV0b24uXG4gICAqL1xuICBzdGF0aWMgZ2V0SW5zdGFuY2UoKTogTW5lbW9uaWMge1xuICAgIGlmICghTW5lbW9uaWMuaW5zdGFuY2UpIHtcbiAgICAgIE1uZW1vbmljLmluc3RhbmNlID0gbmV3IE1uZW1vbmljKClcbiAgICB9XG4gICAgcmV0dXJuIE1uZW1vbmljLmluc3RhbmNlXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJuIHdvcmRsaXN0c1xuICAgKlxuICAgKiBAcGFyYW0gbGFuZ3VhZ2UgYSBzdHJpbmcgc3BlY2lmeWluZyB0aGUgbGFuZ3VhZ2VcbiAgICpcbiAgICogQHJldHVybnMgQSBbW1dvcmRsaXN0XV0gb2JqZWN0IG9yIGFycmF5IG9mIHN0cmluZ3NcbiAgICovXG4gIGdldFdvcmRsaXN0cyhsYW5ndWFnZT86IHN0cmluZyk6IHN0cmluZ1tdIHwgV29yZGxpc3Qge1xuICAgIGlmIChsYW5ndWFnZSAhPT0gdW5kZWZpbmVkKSB7XG4gICAgICByZXR1cm4gdGhpcy53b3JkbGlzdHNbYCR7bGFuZ3VhZ2V9YF1cbiAgICB9IGVsc2Uge1xuICAgICAgcmV0dXJuIHRoaXMud29yZGxpc3RzXG4gICAgfVxuICB9XG5cbiAgLyoqXG4gICAqIFN5bmNocm9ub3VzbHkgdGFrZXMgbW5lbW9uaWMgYW5kIHBhc3N3b3JkIGFuZCByZXR1cm5zIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9XG4gICAqXG4gICAqIEBwYXJhbSBtbmVtb25pYyB0aGUgbW5lbW9uaWMgYXMgYSBzdHJpbmdcbiAgICogQHBhcmFtIHBhc3N3b3JkIHRoZSBwYXNzd29yZCBhcyBhIHN0cmluZ1xuICAgKlxuICAgKiBAcmV0dXJucyBBIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9XG4gICAqL1xuICBtbmVtb25pY1RvU2VlZFN5bmMobW5lbW9uaWM6IHN0cmluZywgcGFzc3dvcmQ6IHN0cmluZyA9IFwiXCIpOiBCdWZmZXIge1xuICAgIGNvbnN0IHNlZWQ6IEJ1ZmZlciA9IGJpcDM5Lm1uZW1vbmljVG9TZWVkU3luYyhtbmVtb25pYywgcGFzc3dvcmQpXG4gICAgcmV0dXJuIEJ1ZmZlci5mcm9tKHNlZWQpXG4gIH1cblxuICAvKipcbiAgICogQXN5bmNocm9ub3VzbHkgdGFrZXMgbW5lbW9uaWMgYW5kIHBhc3N3b3JkIGFuZCByZXR1cm5zIFByb21pc2Uge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn1cbiAgICpcbiAgICogQHBhcmFtIG1uZW1vbmljIHRoZSBtbmVtb25pYyBhcyBhIHN0cmluZ1xuICAgKiBAcGFyYW0gcGFzc3dvcmQgdGhlIHBhc3N3b3JkIGFzIGEgc3RyaW5nXG4gICAqXG4gICAqIEByZXR1cm5zIEEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn1cbiAgICovXG4gIGFzeW5jIG1uZW1vbmljVG9TZWVkKFxuICAgIG1uZW1vbmljOiBzdHJpbmcsXG4gICAgcGFzc3dvcmQ6IHN0cmluZyA9IFwiXCJcbiAgKTogUHJvbWlzZTxCdWZmZXI+IHtcbiAgICBjb25zdCBzZWVkOiBCdWZmZXIgPSBhd2FpdCBiaXAzOS5tbmVtb25pY1RvU2VlZChtbmVtb25pYywgcGFzc3dvcmQpXG4gICAgcmV0dXJuIEJ1ZmZlci5mcm9tKHNlZWQpXG4gIH1cblxuICAvKipcbiAgICogVGFrZXMgbW5lbW9uaWMgYW5kIHdvcmRsaXN0IGFuZCByZXR1cm5zIGJ1ZmZlclxuICAgKlxuICAgKiBAcGFyYW0gbW5lbW9uaWMgdGhlIG1uZW1vbmljIGFzIGEgc3RyaW5nXG4gICAqIEBwYXJhbSB3b3JkbGlzdCBPcHRpb25hbCB0aGUgd29yZGxpc3QgYXMgYW4gYXJyYXkgb2Ygc3RyaW5nc1xuICAgKlxuICAgKiBAcmV0dXJucyBBIHN0cmluZ1xuICAgKi9cbiAgbW5lbW9uaWNUb0VudHJvcHkobW5lbW9uaWM6IHN0cmluZywgd29yZGxpc3Q/OiBzdHJpbmdbXSk6IHN0cmluZyB7XG4gICAgcmV0dXJuIGJpcDM5Lm1uZW1vbmljVG9FbnRyb3B5KG1uZW1vbmljLCB3b3JkbGlzdClcbiAgfVxuXG4gIC8qKlxuICAgKiBUYWtlcyBtbmVtb25pYyBhbmQgd29yZGxpc3QgYW5kIHJldHVybnMgYnVmZmVyXG4gICAqXG4gICAqIEBwYXJhbSBlbnRyb3B5IHRoZSBlbnRyb3B5IGFzIGEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gb3IgYXMgYSBzdHJpbmdcbiAgICogQHBhcmFtIHdvcmRsaXN0IE9wdGlvbmFsLCB0aGUgd29yZGxpc3QgYXMgYW4gYXJyYXkgb2Ygc3RyaW5nc1xuICAgKlxuICAgKiBAcmV0dXJucyBBIHN0cmluZ1xuICAgKi9cbiAgZW50cm9weVRvTW5lbW9uaWMoZW50cm9weTogQnVmZmVyIHwgc3RyaW5nLCB3b3JkbGlzdD86IHN0cmluZ1tdKTogc3RyaW5nIHtcbiAgICByZXR1cm4gYmlwMzkuZW50cm9weVRvTW5lbW9uaWMoZW50cm9weSwgd29yZGxpc3QpXG4gIH1cblxuICAvKipcbiAgICogVmFsaWRhdGVzIGEgbW5lbW9uaWNcbiAgIDExKlxuICAgKiBAcGFyYW0gbW5lbW9uaWMgdGhlIG1uZW1vbmljIGFzIGEgc3RyaW5nXG4gICAqIEBwYXJhbSB3b3JkbGlzdCBPcHRpb25hbCB0aGUgd29yZGxpc3QgYXMgYW4gYXJyYXkgb2Ygc3RyaW5nc1xuICAgKlxuICAgKiBAcmV0dXJucyBBIHN0cmluZ1xuICAgKi9cbiAgdmFsaWRhdGVNbmVtb25pYyhtbmVtb25pYzogc3RyaW5nLCB3b3JkbGlzdD86IHN0cmluZ1tdKTogc3RyaW5nIHtcbiAgICByZXR1cm4gYmlwMzkudmFsaWRhdGVNbmVtb25pYyhtbmVtb25pYywgd29yZGxpc3QpXG4gIH1cblxuICAvKipcbiAgICogU2V0cyB0aGUgZGVmYXVsdCB3b3JkIGxpc3RcbiAgICpcbiAgICogQHBhcmFtIGxhbmd1YWdlIHRoZSBsYW5ndWFnZSBhcyBhIHN0cmluZ1xuICAgKlxuICAgKi9cbiAgc2V0RGVmYXVsdFdvcmRsaXN0KGxhbmd1YWdlOiBzdHJpbmcpOiB2b2lkIHtcbiAgICBiaXAzOS5zZXREZWZhdWx0V29yZGxpc3QobGFuZ3VhZ2UpXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgbGFuZ3VhZ2Ugb2YgdGhlIGRlZmF1bHQgd29yZCBsaXN0XG4gICAqXG4gICAqIEByZXR1cm5zIEEgc3RyaW5nXG4gICAqL1xuICBnZXREZWZhdWx0V29yZGxpc3QoKTogc3RyaW5nIHtcbiAgICByZXR1cm4gYmlwMzkuZ2V0RGVmYXVsdFdvcmRsaXN0KClcbiAgfVxuXG4gIC8qKlxuICAgKiBHZW5lcmF0ZSBhIHJhbmRvbSBtbmVtb25pYyAodXNlcyBjcnlwdG8ucmFuZG9tQnl0ZXMgdW5kZXIgdGhlIGhvb2QpLCBkZWZhdWx0cyB0byAyNTYtYml0cyBvZiBlbnRyb3B5XG4gICAqXG4gICAqIEBwYXJhbSBzdHJlbmd0aCBPcHRpb25hbCB0aGUgc3RyZW5ndGggYXMgYSBudW1iZXJcbiAgICogQHBhcmFtIHJuZyBPcHRpb25hbCB0aGUgcmFuZG9tIG51bWJlciBnZW5lcmF0b3IuIERlZmF1bHRzIHRvIGNyeXB0by5yYW5kb21CeXRlc1xuICAgKiBAcGFyYW0gd29yZGxpc3QgT3B0aW9uYWxcbiAgICpcbiAgICovXG4gIGdlbmVyYXRlTW5lbW9uaWMoXG4gICAgc3RyZW5ndGg/OiBudW1iZXIsXG4gICAgcm5nPzogKHNpemU6IG51bWJlcikgPT4gQnVmZmVyLFxuICAgIHdvcmRsaXN0Pzogc3RyaW5nW11cbiAgKTogc3RyaW5nIHtcbiAgICBzdHJlbmd0aCA9IHN0cmVuZ3RoIHx8IDI1NlxuICAgIGlmIChzdHJlbmd0aCAlIDMyICE9PSAwKSB7XG4gICAgICB0aHJvdyBuZXcgSW52YWxpZEVudHJvcHkoXCJFcnJvciAtIEludmFsaWQgZW50cm9weVwiKVxuICAgIH1cbiAgICBybmcgPSBybmcgfHwgcmFuZG9tQnl0ZXNcbiAgICByZXR1cm4gYmlwMzkuZ2VuZXJhdGVNbmVtb25pYyhzdHJlbmd0aCwgcm5nLCB3b3JkbGlzdClcbiAgfVxufVxuIl19