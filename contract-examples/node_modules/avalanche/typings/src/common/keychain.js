"use strict";
/**
 * @packageDocumentation
 * @module Common-KeyChain
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.StandardKeyChain = exports.StandardKeyPair = void 0;
const buffer_1 = require("buffer/");
/**
 * Class for representing a private and public keypair in Avalanche.
 * All APIs that need key pairs should extend on this class.
 */
class StandardKeyPair {
    /**
     * Returns a reference to the private key.
     *
     * @returns A {@link https://github.com/feross/buffer|Buffer} containing the private key
     */
    getPrivateKey() {
        return this.privk;
    }
    /**
     * Returns a reference to the public key.
     *
     * @returns A {@link https://github.com/feross/buffer|Buffer} containing the public key
     */
    getPublicKey() {
        return this.pubk;
    }
}
exports.StandardKeyPair = StandardKeyPair;
/**
 * Class for representing a key chain in Avalanche.
 * All endpoints that need key chains should extend on this class.
 *
 * @typeparam KPClass extending [[StandardKeyPair]] which is used as the key in [[StandardKeyChain]]
 */
class StandardKeyChain {
    constructor() {
        this.keys = {};
        /**
         * Gets an array of addresses stored in the [[StandardKeyChain]].
         *
         * @returns An array of {@link https://github.com/feross/buffer|Buffer}  representations
         * of the addresses
         */
        this.getAddresses = () => Object.values(this.keys).map((kp) => kp.getAddress());
        /**
         * Gets an array of addresses stored in the [[StandardKeyChain]].
         *
         * @returns An array of string representations of the addresses
         */
        this.getAddressStrings = () => Object.values(this.keys).map((kp) => kp.getAddressString());
        /**
         * Removes the key pair from the list of they keys managed in the [[StandardKeyChain]].
         *
         * @param key A {@link https://github.com/feross/buffer|Buffer} for the address or
         * KPClass to remove
         *
         * @returns The boolean true if a key was removed.
         */
        this.removeKey = (key) => {
            let kaddr;
            if (key instanceof buffer_1.Buffer) {
                kaddr = key.toString("hex");
            }
            else {
                kaddr = key.getAddress().toString("hex");
            }
            if (kaddr in this.keys) {
                delete this.keys[`${kaddr}`];
                return true;
            }
            return false;
        };
        /**
         * Checks if there is a key associated with the provided address.
         *
         * @param address The address to check for existence in the keys database
         *
         * @returns True on success, false if not found
         */
        this.hasKey = (address) => address.toString("hex") in this.keys;
        /**
         * Returns the [[StandardKeyPair]] listed under the provided address
         *
         * @param address The {@link https://github.com/feross/buffer|Buffer} of the address to
         * retrieve from the keys database
         *
         * @returns A reference to the [[StandardKeyPair]] in the keys database
         */
        this.getKey = (address) => this.keys[address.toString("hex")];
    }
    /**
     * Adds the key pair to the list of the keys managed in the [[StandardKeyChain]].
     *
     * @param newKey A key pair of the appropriate class to be added to the [[StandardKeyChain]]
     */
    addKey(newKey) {
        this.keys[newKey.getAddress().toString("hex")] = newKey;
    }
}
exports.StandardKeyChain = StandardKeyChain;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoia2V5Y2hhaW4uanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi9zcmMvY29tbW9uL2tleWNoYWluLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7QUFBQTs7O0dBR0c7OztBQUVILG9DQUFnQztBQUVoQzs7O0dBR0c7QUFDSCxNQUFzQixlQUFlO0lBb0RuQzs7OztPQUlHO0lBQ0gsYUFBYTtRQUNYLE9BQU8sSUFBSSxDQUFDLEtBQUssQ0FBQTtJQUNuQixDQUFDO0lBRUQ7Ozs7T0FJRztJQUNILFlBQVk7UUFDVixPQUFPLElBQUksQ0FBQyxJQUFJLENBQUE7SUFDbEIsQ0FBQztDQWlDRjtBQXJHRCwwQ0FxR0M7QUFFRDs7Ozs7R0FLRztBQUNILE1BQXNCLGdCQUFnQjtJQUF0QztRQUNZLFNBQUksR0FBbUMsRUFBRSxDQUFBO1FBa0JuRDs7Ozs7V0FLRztRQUNILGlCQUFZLEdBQUcsR0FBYSxFQUFFLENBQzVCLE1BQU0sQ0FBQyxNQUFNLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDLEdBQUcsQ0FBQyxDQUFDLEVBQUUsRUFBRSxFQUFFLENBQUMsRUFBRSxDQUFDLFVBQVUsRUFBRSxDQUFDLENBQUE7UUFFdkQ7Ozs7V0FJRztRQUNILHNCQUFpQixHQUFHLEdBQWEsRUFBRSxDQUNqQyxNQUFNLENBQUMsTUFBTSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQyxHQUFHLENBQUMsQ0FBQyxFQUFFLEVBQUUsRUFBRSxDQUFDLEVBQUUsQ0FBQyxnQkFBZ0IsRUFBRSxDQUFDLENBQUE7UUFXN0Q7Ozs7Ozs7V0FPRztRQUNILGNBQVMsR0FBRyxDQUFDLEdBQXFCLEVBQUUsRUFBRTtZQUNwQyxJQUFJLEtBQWEsQ0FBQTtZQUNqQixJQUFJLEdBQUcsWUFBWSxlQUFNLEVBQUU7Z0JBQ3pCLEtBQUssR0FBRyxHQUFHLENBQUMsUUFBUSxDQUFDLEtBQUssQ0FBQyxDQUFBO2FBQzVCO2lCQUFNO2dCQUNMLEtBQUssR0FBRyxHQUFHLENBQUMsVUFBVSxFQUFFLENBQUMsUUFBUSxDQUFDLEtBQUssQ0FBQyxDQUFBO2FBQ3pDO1lBQ0QsSUFBSSxLQUFLLElBQUksSUFBSSxDQUFDLElBQUksRUFBRTtnQkFDdEIsT0FBTyxJQUFJLENBQUMsSUFBSSxDQUFDLEdBQUcsS0FBSyxFQUFFLENBQUMsQ0FBQTtnQkFDNUIsT0FBTyxJQUFJLENBQUE7YUFDWjtZQUNELE9BQU8sS0FBSyxDQUFBO1FBQ2QsQ0FBQyxDQUFBO1FBRUQ7Ozs7OztXQU1HO1FBQ0gsV0FBTSxHQUFHLENBQUMsT0FBZSxFQUFXLEVBQUUsQ0FBQyxPQUFPLENBQUMsUUFBUSxDQUFDLEtBQUssQ0FBQyxJQUFJLElBQUksQ0FBQyxJQUFJLENBQUE7UUFFM0U7Ozs7Ozs7V0FPRztRQUNILFdBQU0sR0FBRyxDQUFDLE9BQWUsRUFBVyxFQUFFLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxPQUFPLENBQUMsUUFBUSxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUE7SUFPM0UsQ0FBQztJQXZEQzs7OztPQUlHO0lBQ0gsTUFBTSxDQUFDLE1BQWU7UUFDcEIsSUFBSSxDQUFDLElBQUksQ0FBQyxNQUFNLENBQUMsVUFBVSxFQUFFLENBQUMsUUFBUSxDQUFDLEtBQUssQ0FBQyxDQUFDLEdBQUcsTUFBTSxDQUFBO0lBQ3pELENBQUM7Q0FnREY7QUEzRkQsNENBMkZDIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgQ29tbW9uLUtleUNoYWluXG4gKi9cblxuaW1wb3J0IHsgQnVmZmVyIH0gZnJvbSBcImJ1ZmZlci9cIlxuXG4vKipcbiAqIENsYXNzIGZvciByZXByZXNlbnRpbmcgYSBwcml2YXRlIGFuZCBwdWJsaWMga2V5cGFpciBpbiBBdmFsYW5jaGUuXG4gKiBBbGwgQVBJcyB0aGF0IG5lZWQga2V5IHBhaXJzIHNob3VsZCBleHRlbmQgb24gdGhpcyBjbGFzcy5cbiAqL1xuZXhwb3J0IGFic3RyYWN0IGNsYXNzIFN0YW5kYXJkS2V5UGFpciB7XG4gIHByb3RlY3RlZCBwdWJrOiBCdWZmZXJcbiAgcHJvdGVjdGVkIHByaXZrOiBCdWZmZXJcblxuICAvKipcbiAgICogR2VuZXJhdGVzIGEgbmV3IGtleXBhaXIuXG4gICAqXG4gICAqIEBwYXJhbSBlbnRyb3B5IE9wdGlvbmFsIHBhcmFtZXRlciB0aGF0IG1heSBiZSBuZWNlc3NhcnkgdG8gcHJvZHVjZSBzZWN1cmUga2V5c1xuICAgKi9cbiAgYWJzdHJhY3QgZ2VuZXJhdGVLZXkoZW50cm9weT86IEJ1ZmZlcik6IHZvaWRcblxuICAvKipcbiAgICogSW1wb3J0cyBhIHByaXZhdGUga2V5IGFuZCBnZW5lcmF0ZXMgdGhlIGFwcHJvcHJpYXRlIHB1YmxpYyBrZXkuXG4gICAqXG4gICAqIEBwYXJhbSBwcml2ayBBIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IHJlcHJlc2VudGluZyB0aGUgcHJpdmF0ZSBrZXlcbiAgICpcbiAgICogQHJldHVybnMgdHJ1ZSBvbiBzdWNjZXNzLCBmYWxzZSBvbiBmYWlsdXJlXG4gICAqL1xuICBhYnN0cmFjdCBpbXBvcnRLZXkocHJpdms6IEJ1ZmZlcik6IGJvb2xlYW5cblxuICAvKipcbiAgICogVGFrZXMgYSBtZXNzYWdlLCBzaWducyBpdCwgYW5kIHJldHVybnMgdGhlIHNpZ25hdHVyZS5cbiAgICpcbiAgICogQHBhcmFtIG1zZyBUaGUgbWVzc2FnZSB0byBzaWduXG4gICAqXG4gICAqIEByZXR1cm5zIEEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gY29udGFpbmluZyB0aGUgc2lnbmF0dXJlXG4gICAqL1xuICBhYnN0cmFjdCBzaWduKG1zZzogQnVmZmVyKTogQnVmZmVyXG5cbiAgLyoqXG4gICAqIFJlY292ZXJzIHRoZSBwdWJsaWMga2V5IG9mIGEgbWVzc2FnZSBzaWduZXIgZnJvbSBhIG1lc3NhZ2UgYW5kIGl0cyBhc3NvY2lhdGVkIHNpZ25hdHVyZS5cbiAgICpcbiAgICogQHBhcmFtIG1zZyBUaGUgbWVzc2FnZSB0aGF0J3Mgc2lnbmVkXG4gICAqIEBwYXJhbSBzaWcgVGhlIHNpZ25hdHVyZSB0aGF0J3Mgc2lnbmVkIG9uIHRoZSBtZXNzYWdlXG4gICAqXG4gICAqIEByZXR1cm5zIEEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gY29udGFpbmluZyB0aGUgcHVibGljXG4gICAqIGtleSBvZiB0aGUgc2lnbmVyXG4gICAqL1xuICBhYnN0cmFjdCByZWNvdmVyKG1zZzogQnVmZmVyLCBzaWc6IEJ1ZmZlcik6IEJ1ZmZlclxuXG4gIC8qKlxuICAgKiBWZXJpZmllcyB0aGF0IHRoZSBwcml2YXRlIGtleSBhc3NvY2lhdGVkIHdpdGggdGhlIHByb3ZpZGVkIHB1YmxpYyBrZXkgcHJvZHVjZXMgdGhlXG4gICAqIHNpZ25hdHVyZSBhc3NvY2lhdGVkIHdpdGggdGhlIGdpdmVuIG1lc3NhZ2UuXG4gICAqXG4gICAqIEBwYXJhbSBtc2cgVGhlIG1lc3NhZ2UgYXNzb2NpYXRlZCB3aXRoIHRoZSBzaWduYXR1cmVcbiAgICogQHBhcmFtIHNpZyBUaGUgc2lnbmF0dXJlIG9mIHRoZSBzaWduZWQgbWVzc2FnZVxuICAgKiBAcGFyYW0gcHViayBUaGUgcHVibGljIGtleSBhc3NvY2lhdGVkIHdpdGggdGhlIG1lc3NhZ2Ugc2lnbmF0dXJlXG4gICAqXG4gICAqIEByZXR1cm5zIFRydWUgb24gc3VjY2VzcywgZmFsc2Ugb24gZmFpbHVyZVxuICAgKi9cbiAgYWJzdHJhY3QgdmVyaWZ5KG1zZzogQnVmZmVyLCBzaWc6IEJ1ZmZlciwgcHViazogQnVmZmVyKTogYm9vbGVhblxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIGEgcmVmZXJlbmNlIHRvIHRoZSBwcml2YXRlIGtleS5cbiAgICpcbiAgICogQHJldHVybnMgQSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSBjb250YWluaW5nIHRoZSBwcml2YXRlIGtleVxuICAgKi9cbiAgZ2V0UHJpdmF0ZUtleSgpOiBCdWZmZXIge1xuICAgIHJldHVybiB0aGlzLnByaXZrXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyBhIHJlZmVyZW5jZSB0byB0aGUgcHVibGljIGtleS5cbiAgICpcbiAgICogQHJldHVybnMgQSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSBjb250YWluaW5nIHRoZSBwdWJsaWMga2V5XG4gICAqL1xuICBnZXRQdWJsaWNLZXkoKTogQnVmZmVyIHtcbiAgICByZXR1cm4gdGhpcy5wdWJrXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyBhIHN0cmluZyByZXByZXNlbnRhdGlvbiBvZiB0aGUgcHJpdmF0ZSBrZXkuXG4gICAqXG4gICAqIEByZXR1cm5zIEEgc3RyaW5nIHJlcHJlc2VudGF0aW9uIG9mIHRoZSBwdWJsaWMga2V5XG4gICAqL1xuICBhYnN0cmFjdCBnZXRQcml2YXRlS2V5U3RyaW5nKCk6IHN0cmluZ1xuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBwdWJsaWMga2V5LlxuICAgKlxuICAgKiBAcmV0dXJucyBBIHN0cmluZyByZXByZXNlbnRhdGlvbiBvZiB0aGUgcHVibGljIGtleVxuICAgKi9cbiAgYWJzdHJhY3QgZ2V0UHVibGljS2V5U3RyaW5nKCk6IHN0cmluZ1xuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBhZGRyZXNzLlxuICAgKlxuICAgKiBAcmV0dXJucyBBIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9ICByZXByZXNlbnRhdGlvbiBvZiB0aGUgYWRkcmVzc1xuICAgKi9cbiAgYWJzdHJhY3QgZ2V0QWRkcmVzcygpOiBCdWZmZXJcblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgYWRkcmVzcydzIHN0cmluZyByZXByZXNlbnRhdGlvbi5cbiAgICpcbiAgICogQHJldHVybnMgQSBzdHJpbmcgcmVwcmVzZW50YXRpb24gb2YgdGhlIGFkZHJlc3NcbiAgICovXG4gIGFic3RyYWN0IGdldEFkZHJlc3NTdHJpbmcoKTogc3RyaW5nXG5cbiAgYWJzdHJhY3QgY3JlYXRlKC4uLmFyZ3M6IGFueVtdKTogdGhpc1xuXG4gIGFic3RyYWN0IGNsb25lKCk6IHRoaXNcbn1cblxuLyoqXG4gKiBDbGFzcyBmb3IgcmVwcmVzZW50aW5nIGEga2V5IGNoYWluIGluIEF2YWxhbmNoZS5cbiAqIEFsbCBlbmRwb2ludHMgdGhhdCBuZWVkIGtleSBjaGFpbnMgc2hvdWxkIGV4dGVuZCBvbiB0aGlzIGNsYXNzLlxuICpcbiAqIEB0eXBlcGFyYW0gS1BDbGFzcyBleHRlbmRpbmcgW1tTdGFuZGFyZEtleVBhaXJdXSB3aGljaCBpcyB1c2VkIGFzIHRoZSBrZXkgaW4gW1tTdGFuZGFyZEtleUNoYWluXV1cbiAqL1xuZXhwb3J0IGFic3RyYWN0IGNsYXNzIFN0YW5kYXJkS2V5Q2hhaW48S1BDbGFzcyBleHRlbmRzIFN0YW5kYXJkS2V5UGFpcj4ge1xuICBwcm90ZWN0ZWQga2V5czogeyBbYWRkcmVzczogc3RyaW5nXTogS1BDbGFzcyB9ID0ge31cblxuICAvKipcbiAgICogTWFrZXMgYSBuZXcgW1tTdGFuZGFyZEtleVBhaXJdXSwgcmV0dXJucyB0aGUgYWRkcmVzcy5cbiAgICpcbiAgICogQHJldHVybnMgQWRkcmVzcyBvZiB0aGUgbmV3IFtbU3RhbmRhcmRLZXlQYWlyXV1cbiAgICovXG4gIG1ha2VLZXk6ICgpID0+IEtQQ2xhc3NcblxuICAvKipcbiAgICogR2l2ZW4gYSBwcml2YXRlIGtleSwgbWFrZXMgYSBuZXcgW1tTdGFuZGFyZEtleVBhaXJdXSwgcmV0dXJucyB0aGUgYWRkcmVzcy5cbiAgICpcbiAgICogQHBhcmFtIHByaXZrIEEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gcmVwcmVzZW50aW5nIHRoZSBwcml2YXRlIGtleVxuICAgKlxuICAgKiBAcmV0dXJucyBBIG5ldyBbW1N0YW5kYXJkS2V5UGFpcl1dXG4gICAqL1xuICBpbXBvcnRLZXk6IChwcml2azogQnVmZmVyKSA9PiBLUENsYXNzXG5cbiAgLyoqXG4gICAqIEdldHMgYW4gYXJyYXkgb2YgYWRkcmVzc2VzIHN0b3JlZCBpbiB0aGUgW1tTdGFuZGFyZEtleUNoYWluXV0uXG4gICAqXG4gICAqIEByZXR1cm5zIEFuIGFycmF5IG9mIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9ICByZXByZXNlbnRhdGlvbnNcbiAgICogb2YgdGhlIGFkZHJlc3Nlc1xuICAgKi9cbiAgZ2V0QWRkcmVzc2VzID0gKCk6IEJ1ZmZlcltdID0+XG4gICAgT2JqZWN0LnZhbHVlcyh0aGlzLmtleXMpLm1hcCgoa3ApID0+IGtwLmdldEFkZHJlc3MoKSlcblxuICAvKipcbiAgICogR2V0cyBhbiBhcnJheSBvZiBhZGRyZXNzZXMgc3RvcmVkIGluIHRoZSBbW1N0YW5kYXJkS2V5Q2hhaW5dXS5cbiAgICpcbiAgICogQHJldHVybnMgQW4gYXJyYXkgb2Ygc3RyaW5nIHJlcHJlc2VudGF0aW9ucyBvZiB0aGUgYWRkcmVzc2VzXG4gICAqL1xuICBnZXRBZGRyZXNzU3RyaW5ncyA9ICgpOiBzdHJpbmdbXSA9PlxuICAgIE9iamVjdC52YWx1ZXModGhpcy5rZXlzKS5tYXAoKGtwKSA9PiBrcC5nZXRBZGRyZXNzU3RyaW5nKCkpXG5cbiAgLyoqXG4gICAqIEFkZHMgdGhlIGtleSBwYWlyIHRvIHRoZSBsaXN0IG9mIHRoZSBrZXlzIG1hbmFnZWQgaW4gdGhlIFtbU3RhbmRhcmRLZXlDaGFpbl1dLlxuICAgKlxuICAgKiBAcGFyYW0gbmV3S2V5IEEga2V5IHBhaXIgb2YgdGhlIGFwcHJvcHJpYXRlIGNsYXNzIHRvIGJlIGFkZGVkIHRvIHRoZSBbW1N0YW5kYXJkS2V5Q2hhaW5dXVxuICAgKi9cbiAgYWRkS2V5KG5ld0tleTogS1BDbGFzcykge1xuICAgIHRoaXMua2V5c1tuZXdLZXkuZ2V0QWRkcmVzcygpLnRvU3RyaW5nKFwiaGV4XCIpXSA9IG5ld0tleVxuICB9XG5cbiAgLyoqXG4gICAqIFJlbW92ZXMgdGhlIGtleSBwYWlyIGZyb20gdGhlIGxpc3Qgb2YgdGhleSBrZXlzIG1hbmFnZWQgaW4gdGhlIFtbU3RhbmRhcmRLZXlDaGFpbl1dLlxuICAgKlxuICAgKiBAcGFyYW0ga2V5IEEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gZm9yIHRoZSBhZGRyZXNzIG9yXG4gICAqIEtQQ2xhc3MgdG8gcmVtb3ZlXG4gICAqXG4gICAqIEByZXR1cm5zIFRoZSBib29sZWFuIHRydWUgaWYgYSBrZXkgd2FzIHJlbW92ZWQuXG4gICAqL1xuICByZW1vdmVLZXkgPSAoa2V5OiBLUENsYXNzIHwgQnVmZmVyKSA9PiB7XG4gICAgbGV0IGthZGRyOiBzdHJpbmdcbiAgICBpZiAoa2V5IGluc3RhbmNlb2YgQnVmZmVyKSB7XG4gICAgICBrYWRkciA9IGtleS50b1N0cmluZyhcImhleFwiKVxuICAgIH0gZWxzZSB7XG4gICAgICBrYWRkciA9IGtleS5nZXRBZGRyZXNzKCkudG9TdHJpbmcoXCJoZXhcIilcbiAgICB9XG4gICAgaWYgKGthZGRyIGluIHRoaXMua2V5cykge1xuICAgICAgZGVsZXRlIHRoaXMua2V5c1tgJHtrYWRkcn1gXVxuICAgICAgcmV0dXJuIHRydWVcbiAgICB9XG4gICAgcmV0dXJuIGZhbHNlXG4gIH1cblxuICAvKipcbiAgICogQ2hlY2tzIGlmIHRoZXJlIGlzIGEga2V5IGFzc29jaWF0ZWQgd2l0aCB0aGUgcHJvdmlkZWQgYWRkcmVzcy5cbiAgICpcbiAgICogQHBhcmFtIGFkZHJlc3MgVGhlIGFkZHJlc3MgdG8gY2hlY2sgZm9yIGV4aXN0ZW5jZSBpbiB0aGUga2V5cyBkYXRhYmFzZVxuICAgKlxuICAgKiBAcmV0dXJucyBUcnVlIG9uIHN1Y2Nlc3MsIGZhbHNlIGlmIG5vdCBmb3VuZFxuICAgKi9cbiAgaGFzS2V5ID0gKGFkZHJlc3M6IEJ1ZmZlcik6IGJvb2xlYW4gPT4gYWRkcmVzcy50b1N0cmluZyhcImhleFwiKSBpbiB0aGlzLmtleXNcblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgW1tTdGFuZGFyZEtleVBhaXJdXSBsaXN0ZWQgdW5kZXIgdGhlIHByb3ZpZGVkIGFkZHJlc3NcbiAgICpcbiAgICogQHBhcmFtIGFkZHJlc3MgVGhlIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IG9mIHRoZSBhZGRyZXNzIHRvXG4gICAqIHJldHJpZXZlIGZyb20gdGhlIGtleXMgZGF0YWJhc2VcbiAgICpcbiAgICogQHJldHVybnMgQSByZWZlcmVuY2UgdG8gdGhlIFtbU3RhbmRhcmRLZXlQYWlyXV0gaW4gdGhlIGtleXMgZGF0YWJhc2VcbiAgICovXG4gIGdldEtleSA9IChhZGRyZXNzOiBCdWZmZXIpOiBLUENsYXNzID0+IHRoaXMua2V5c1thZGRyZXNzLnRvU3RyaW5nKFwiaGV4XCIpXVxuXG4gIGFic3RyYWN0IGNyZWF0ZSguLi5hcmdzOiBhbnlbXSk6IHRoaXNcblxuICBhYnN0cmFjdCBjbG9uZSgpOiB0aGlzXG5cbiAgYWJzdHJhY3QgdW5pb24oa2M6IHRoaXMpOiB0aGlzXG59XG4iXX0=