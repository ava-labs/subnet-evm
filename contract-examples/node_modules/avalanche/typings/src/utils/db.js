"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
/**
 * @packageDocumentation
 * @module Utils-DB
 */
const store2_1 = __importDefault(require("store2"));
/**
 * A class for interacting with the {@link https://github.com/nbubna/store| store2 module}
 *
 * This class should never be instantiated directly. Instead, invoke the "DB.getInstance()" static
 * function to grab the singleton instance of the database.
 *
 * ```js
 * const db = DB.getInstance();
 * const blockchaindb = db.getNamespace("mychain");
 * ```
 */
class DB {
    constructor() { }
    /**
     * Retrieves the database singleton.
     */
    static getInstance() {
        if (!DB.instance) {
            DB.instance = new DB();
        }
        return DB.instance;
    }
    /**
     * Gets a namespace from the database singleton.
     *
     * @param ns Namespace to retrieve.
     */
    static getNamespace(ns) {
        return this.store.namespace(ns);
    }
}
exports.default = DB;
DB.store = store2_1.default;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZGIuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi9zcmMvdXRpbHMvZGIudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7QUFBQTs7O0dBR0c7QUFDSCxvREFBd0M7QUFFeEM7Ozs7Ozs7Ozs7R0FVRztBQUNILE1BQXFCLEVBQUU7SUFLckIsZ0JBQXVCLENBQUM7SUFFeEI7O09BRUc7SUFDSCxNQUFNLENBQUMsV0FBVztRQUNoQixJQUFJLENBQUMsRUFBRSxDQUFDLFFBQVEsRUFBRTtZQUNoQixFQUFFLENBQUMsUUFBUSxHQUFHLElBQUksRUFBRSxFQUFFLENBQUE7U0FDdkI7UUFDRCxPQUFPLEVBQUUsQ0FBQyxRQUFRLENBQUE7SUFDcEIsQ0FBQztJQUVEOzs7O09BSUc7SUFDSCxNQUFNLENBQUMsWUFBWSxDQUFDLEVBQVU7UUFDNUIsT0FBTyxJQUFJLENBQUMsS0FBSyxDQUFDLFNBQVMsQ0FBQyxFQUFFLENBQUMsQ0FBQTtJQUNqQyxDQUFDOztBQXhCSCxxQkF5QkM7QUF0QmdCLFFBQUssR0FBRyxnQkFBSyxDQUFBIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgVXRpbHMtREJcbiAqL1xuaW1wb3J0IHN0b3JlLCB7IFN0b3JlQVBJIH0gZnJvbSBcInN0b3JlMlwiXG5cbi8qKlxuICogQSBjbGFzcyBmb3IgaW50ZXJhY3Rpbmcgd2l0aCB0aGUge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9uYnVibmEvc3RvcmV8IHN0b3JlMiBtb2R1bGV9XG4gKlxuICogVGhpcyBjbGFzcyBzaG91bGQgbmV2ZXIgYmUgaW5zdGFudGlhdGVkIGRpcmVjdGx5LiBJbnN0ZWFkLCBpbnZva2UgdGhlIFwiREIuZ2V0SW5zdGFuY2UoKVwiIHN0YXRpY1xuICogZnVuY3Rpb24gdG8gZ3JhYiB0aGUgc2luZ2xldG9uIGluc3RhbmNlIG9mIHRoZSBkYXRhYmFzZS5cbiAqXG4gKiBgYGBqc1xuICogY29uc3QgZGIgPSBEQi5nZXRJbnN0YW5jZSgpO1xuICogY29uc3QgYmxvY2tjaGFpbmRiID0gZGIuZ2V0TmFtZXNwYWNlKFwibXljaGFpblwiKTtcbiAqIGBgYFxuICovXG5leHBvcnQgZGVmYXVsdCBjbGFzcyBEQiB7XG4gIHByaXZhdGUgc3RhdGljIGluc3RhbmNlOiBEQlxuXG4gIHByaXZhdGUgc3RhdGljIHN0b3JlID0gc3RvcmVcblxuICBwcml2YXRlIGNvbnN0cnVjdG9yKCkge31cblxuICAvKipcbiAgICogUmV0cmlldmVzIHRoZSBkYXRhYmFzZSBzaW5nbGV0b24uXG4gICAqL1xuICBzdGF0aWMgZ2V0SW5zdGFuY2UoKTogREIge1xuICAgIGlmICghREIuaW5zdGFuY2UpIHtcbiAgICAgIERCLmluc3RhbmNlID0gbmV3IERCKClcbiAgICB9XG4gICAgcmV0dXJuIERCLmluc3RhbmNlXG4gIH1cblxuICAvKipcbiAgICogR2V0cyBhIG5hbWVzcGFjZSBmcm9tIHRoZSBkYXRhYmFzZSBzaW5nbGV0b24uXG4gICAqXG4gICAqIEBwYXJhbSBucyBOYW1lc3BhY2UgdG8gcmV0cmlldmUuXG4gICAqL1xuICBzdGF0aWMgZ2V0TmFtZXNwYWNlKG5zOiBzdHJpbmcpOiBTdG9yZUFQSSB7XG4gICAgcmV0dXJuIHRoaXMuc3RvcmUubmFtZXNwYWNlKG5zKVxuICB9XG59XG4iXX0=