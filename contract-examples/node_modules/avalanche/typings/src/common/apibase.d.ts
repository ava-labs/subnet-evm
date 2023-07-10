/**
 * @packageDocumentation
 * @module Common-APIBase
 */
/// <reference types="node" />
import { StoreAPI } from "store2";
import { ClientRequest } from "http";
import AvalancheCore from "../avalanche";
/**
 * Response data for HTTP requests.
 */
export declare class RequestResponseData {
    data: any;
    headers: any;
    status: number;
    statusText: string;
    request: ClientRequest | XMLHttpRequest;
    constructor(data: any, headers: any, status: number, statusText: string, request: ClientRequest | XMLHttpRequest);
}
/**
 * Abstract class defining a generic endpoint that all endpoints must implement (extend).
 */
export declare abstract class APIBase {
    protected core: AvalancheCore;
    protected baseURL: string;
    protected db: StoreAPI;
    /**
     * Sets the path of the APIs baseURL.
     *
     * @param baseURL Path of the APIs baseURL - ex: "/ext/bc/X"
     */
    setBaseURL: (baseURL: string) => void;
    /**
     * Returns the baseURL's path.
     */
    getBaseURL: () => string;
    /**
     * Returns the baseURL's database.
     */
    getDB: () => StoreAPI;
    /**
     *
     * @param core Reference to the Avalanche instance using this baseURL
     * @param baseURL Path to the baseURL
     */
    constructor(core: AvalancheCore, baseURL: string);
}
//# sourceMappingURL=apibase.d.ts.map