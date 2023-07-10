/**
 * @packageDocumentation
 * @module Utils-DB
 */
import { StoreAPI } from "store2";
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
export default class DB {
    private static instance;
    private static store;
    private constructor();
    /**
     * Retrieves the database singleton.
     */
    static getInstance(): DB;
    /**
     * Gets a namespace from the database singleton.
     *
     * @param ns Namespace to retrieve.
     */
    static getNamespace(ns: string): StoreAPI;
}
//# sourceMappingURL=db.d.ts.map