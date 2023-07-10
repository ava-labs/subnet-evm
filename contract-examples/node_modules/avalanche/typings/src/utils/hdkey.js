"use strict";
/**
 * @packageDocumentation
 * @module Utils-HDKey
 */
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const hdkey_1 = __importDefault(require("hdkey"));
/**
 * BIP32 hierarchical deterministic keys.
 *
 */
class HDKey {
    constructor() { }
    /**
     * Retrieves the HDKey singleton.
     */
    static getInstance() {
        if (!HDKey.instance) {
            HDKey.instance = new HDKey();
        }
        return HDKey.instance;
    }
    /**
     * Creates an HDNode from a master seed buffer
     *
     * @param seedBuffer Buffer
     *
     * @returns HDNode
     */
    fromMasterSeed(seed) {
        return hdkey_1.default.fromMasterSeed(seed);
    }
    /**
     * Creates an HDNode from a xprv or xpub extended key string. Accepts an optional versions object.
     *
     * @param xpriv string
     *
     * @returns HDNode
     */
    fromExtendedKey(xpriv) {
        return hdkey_1.default.fromExtendedKey(xpriv);
    }
    /**
     * Creates an HDNode from an object created via hdkey.toJSON().
     *
     * @param obj HDKeyJSON
     *
     * @returns HDNode
     */
    fromJSON(obj) {
        return hdkey_1.default.fromJSON(obj);
    }
}
exports.default = HDKey;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiaGRrZXkuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi9zcmMvdXRpbHMvaGRrZXkudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IjtBQUFBOzs7R0FHRzs7Ozs7QUFJSCxrREFBd0M7QUFHeEM7OztHQUdHO0FBQ0gsTUFBcUIsS0FBSztJQUV4QixnQkFBd0IsQ0FBQztJQUV6Qjs7T0FFRztJQUNILE1BQU0sQ0FBQyxXQUFXO1FBQ2hCLElBQUksQ0FBQyxLQUFLLENBQUMsUUFBUSxFQUFFO1lBQ25CLEtBQUssQ0FBQyxRQUFRLEdBQUcsSUFBSSxLQUFLLEVBQUUsQ0FBQTtTQUM3QjtRQUNELE9BQU8sS0FBSyxDQUFDLFFBQVEsQ0FBQTtJQUN2QixDQUFDO0lBRUQ7Ozs7OztPQU1HO0lBQ0gsY0FBYyxDQUFDLElBQVk7UUFDekIsT0FBTyxlQUFLLENBQUMsY0FBYyxDQUFDLElBQW9DLENBQUMsQ0FBQTtJQUNuRSxDQUFDO0lBRUQ7Ozs7OztPQU1HO0lBQ0gsZUFBZSxDQUFDLEtBQWE7UUFDM0IsT0FBTyxlQUFLLENBQUMsZUFBZSxDQUFDLEtBQUssQ0FBQyxDQUFBO0lBQ3JDLENBQUM7SUFFRDs7Ozs7O09BTUc7SUFDSCxRQUFRLENBQUMsR0FBYztRQUNyQixPQUFPLGVBQUssQ0FBQyxRQUFRLENBQUMsR0FBRyxDQUFDLENBQUE7SUFDNUIsQ0FBQztDQUNGO0FBOUNELHdCQThDQyIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICogQHBhY2thZ2VEb2N1bWVudGF0aW9uXG4gKiBAbW9kdWxlIFV0aWxzLUhES2V5XG4gKi9cblxuaW1wb3J0IHsgQnVmZmVyIH0gZnJvbSAnYnVmZmVyLydcbmltcG9ydCB7IEhES2V5SlNPTiB9IGZyb20gJ3NyYy9jb21tb24nXG5pbXBvcnQgeyBkZWZhdWx0IGFzIGhka2V5IH0gZnJvbSAnaGRrZXknXG5pbXBvcnQgSEROb2RlIGZyb20gJy4vaGRub2RlJ1xuXG4vKipcbiAqIEJJUDMyIGhpZXJhcmNoaWNhbCBkZXRlcm1pbmlzdGljIGtleXMuXG4gKlxuICovXG5leHBvcnQgZGVmYXVsdCBjbGFzcyBIREtleSB7XG4gIHByaXZhdGUgc3RhdGljIGluc3RhbmNlOiBIREtleVxuICBwcml2YXRlIGNvbnN0cnVjdG9yKCkgeyB9XG5cbiAgLyoqXG4gICAqIFJldHJpZXZlcyB0aGUgSERLZXkgc2luZ2xldG9uLlxuICAgKi9cbiAgc3RhdGljIGdldEluc3RhbmNlKCk6IEhES2V5IHtcbiAgICBpZiAoIUhES2V5Lmluc3RhbmNlKSB7XG4gICAgICBIREtleS5pbnN0YW5jZSA9IG5ldyBIREtleSgpXG4gICAgfVxuICAgIHJldHVybiBIREtleS5pbnN0YW5jZVxuICB9XG5cbiAgLyoqXG4gICAqIENyZWF0ZXMgYW4gSEROb2RlIGZyb20gYSBtYXN0ZXIgc2VlZCBidWZmZXJcbiAgICpcbiAgICogQHBhcmFtIHNlZWRCdWZmZXIgQnVmZmVyXG4gICAqXG4gICAqIEByZXR1cm5zIEhETm9kZVxuICAgKi9cbiAgZnJvbU1hc3RlclNlZWQoc2VlZDogQnVmZmVyKTogaGRrZXkge1xuICAgIHJldHVybiBoZGtleS5mcm9tTWFzdGVyU2VlZChzZWVkIGFzIHVua25vd24gYXMgZ2xvYmFsVGhpcy5CdWZmZXIpXG4gIH1cblxuICAvKipcbiAgICogQ3JlYXRlcyBhbiBIRE5vZGUgZnJvbSBhIHhwcnYgb3IgeHB1YiBleHRlbmRlZCBrZXkgc3RyaW5nLiBBY2NlcHRzIGFuIG9wdGlvbmFsIHZlcnNpb25zIG9iamVjdC5cbiAgICpcbiAgICogQHBhcmFtIHhwcml2IHN0cmluZ1xuICAgKlxuICAgKiBAcmV0dXJucyBIRE5vZGVcbiAgICovXG4gIGZyb21FeHRlbmRlZEtleSh4cHJpdjogc3RyaW5nKTogaGRrZXkge1xuICAgIHJldHVybiBoZGtleS5mcm9tRXh0ZW5kZWRLZXkoeHByaXYpXG4gIH1cblxuICAvKipcbiAgICogQ3JlYXRlcyBhbiBIRE5vZGUgZnJvbSBhbiBvYmplY3QgY3JlYXRlZCB2aWEgaGRrZXkudG9KU09OKCkuXG4gICAqXG4gICAqIEBwYXJhbSBvYmogSERLZXlKU09OXG4gICAqIFxuICAgKiBAcmV0dXJucyBIRE5vZGVcbiAgICovXG4gIGZyb21KU09OKG9iajogSERLZXlKU09OKTogaGRrZXkge1xuICAgIHJldHVybiBoZGtleS5mcm9tSlNPTihvYmopXG4gIH1cbn0iXX0=