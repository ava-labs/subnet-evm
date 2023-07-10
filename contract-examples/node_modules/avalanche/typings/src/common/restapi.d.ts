/**
 * @packageDocumentation
 * @module Common-RESTAPI
 */
import { AxiosRequestConfig } from "axios";
import AvalancheCore from "../avalanche";
import { APIBase, RequestResponseData } from "./apibase";
export declare class RESTAPI extends APIBase {
    protected contentType: string;
    protected acceptType: string;
    protected prepHeaders: (contentType?: string, acceptType?: string) => object;
    protected axConf: () => AxiosRequestConfig;
    get: (baseURL?: string, contentType?: string, acceptType?: string) => Promise<RequestResponseData>;
    post: (method: string, params?: object[] | object, baseURL?: string, contentType?: string, acceptType?: string) => Promise<RequestResponseData>;
    put: (method: string, params?: object[] | object, baseURL?: string, contentType?: string, acceptType?: string) => Promise<RequestResponseData>;
    delete: (method: string, params?: object[] | object, baseURL?: string, contentType?: string, acceptType?: string) => Promise<RequestResponseData>;
    patch: (method: string, params?: object[] | object, baseURL?: string, contentType?: string, acceptType?: string) => Promise<RequestResponseData>;
    /**
     * Returns the type of the entity attached to the incoming request
     */
    getContentType: () => string;
    /**
     * Returns what type of representation is desired at the client side
     */
    getAcceptType: () => string;
    /**
     *
     * @param core Reference to the Avalanche instance using this endpoint
     * @param baseURL Path of the APIs baseURL - ex: "/ext/bc/avm"
     * @param contentType Optional Determines the type of the entity attached to the
     * incoming request
     * @param acceptType Optional Determines the type of representation which is
     * desired on the client side
     */
    constructor(core: AvalancheCore, baseURL: string, contentType?: string, acceptType?: string);
}
//# sourceMappingURL=restapi.d.ts.map