/**
 * @packageDocumentation
 * @module AvalancheCore
 */
import { AxiosRequestConfig, AxiosRequestHeaders, Method } from "axios";
import { APIBase, RequestResponseData } from "./common/apibase";
/**
 * AvalancheCore is middleware for interacting with Avalanche node RPC APIs.
 *
 * Example usage:
 * ```js
 * let avalanche = new AvalancheCore("127.0.0.1", 9650, "https")
 * ```
 *
 *
 */
export default class AvalancheCore {
    protected networkID: number;
    protected hrp: string;
    protected protocol: string;
    protected ip: string;
    protected host: string;
    protected port: number;
    protected baseEndpoint: string;
    protected url: string;
    protected auth: string;
    protected headers: {
        [k: string]: string;
    };
    protected requestConfig: AxiosRequestConfig;
    protected apis: {
        [k: string]: APIBase;
    };
    /**
     * Sets the address and port of the main Avalanche Client.
     *
     * @param host The hostname to resolve to reach the Avalanche Client RPC APIs.
     * @param port The port to resolve to reach the Avalanche Client RPC APIs.
     * @param protocol The protocol string to use before a "://" in a request,
     * ex: "http", "https", etc. Defaults to http
     * @param baseEndpoint the base endpoint to reach the Avalanche Client RPC APIs,
     * ex: "/rpc". Defaults to "/"
     * The following special characters are removed from host and protocol
     * &#,@+()$~%'":*?{} also less than and greater than signs
     */
    setAddress: (host: string, port: number, protocol?: string, baseEndpoint?: string) => void;
    /**
     * Returns the protocol such as "http", "https", "git", "ws", etc.
     */
    getProtocol: () => string;
    /**
     * Returns the host for the Avalanche node.
     */
    getHost: () => string;
    /**
     * Returns the IP for the Avalanche node.
     */
    getIP: () => string;
    /**
     * Returns the port for the Avalanche node.
     */
    getPort: () => number;
    /**
     * Returns the base endpoint for the Avalanche node.
     */
    getBaseEndpoint: () => string;
    /**
     * Returns the URL of the Avalanche node (ip + port)
     */
    getURL: () => string;
    /**
     * Returns the custom headers
     */
    getHeaders: () => object;
    /**
     * Returns the custom request config
     */
    getRequestConfig: () => AxiosRequestConfig;
    /**
     * Returns the networkID
     */
    getNetworkID: () => number;
    /**
     * Sets the networkID
     */
    setNetworkID: (netID: number) => void;
    /**
     * Returns the Human-Readable-Part of the network associated with this key.
     *
     * @returns The [[KeyPair]]'s Human-Readable-Part of the network's Bech32 addressing scheme
     */
    getHRP: () => string;
    /**
     * Sets the the Human-Readable-Part of the network associated with this key.
     *
     * @param hrp String for the Human-Readable-Part of Bech32 addresses
     */
    setHRP: (hrp: string) => void;
    /**
     * Adds a new custom header to be included with all requests.
     *
     * @param key Header name
     * @param value Header value
     */
    setHeader: (key: string, value: string) => void;
    /**
     * Removes a previously added custom header.
     *
     * @param key Header name
     */
    removeHeader: (key: string) => void;
    /**
     * Removes all headers.
     */
    removeAllHeaders: () => void;
    /**
     * Adds a new custom config value to be included with all requests.
     *
     * @param key Config name
     * @param value Config value
     */
    setRequestConfig: (key: string, value: string | boolean) => void;
    /**
     * Removes a previously added request config.
     *
     * @param key Header name
     */
    removeRequestConfig: (key: string) => void;
    /**
     * Removes all request configs.
     */
    removeAllRequestConfigs: () => void;
    /**
     * Sets the temporary auth token used for communicating with the node.
     *
     * @param auth A temporary token provided by the node enabling access to the endpoints on the node.
     */
    setAuthToken: (auth: string) => void;
    protected _setHeaders: (headers: any) => AxiosRequestHeaders;
    /**
     * Adds an API to the middleware. The API resolves to a registered blockchain's RPC.
     *
     * In TypeScript:
     * ```js
     * avalanche.addAPI<MyVMClass>("mychain", MyVMClass, "/ext/bc/mychain")
     * ```
     *
     * In Javascript:
     * ```js
     * avalanche.addAPI("mychain", MyVMClass, "/ext/bc/mychain")
     * ```
     *
     * @typeparam GA Class of the API being added
     * @param apiName A label for referencing the API in the future
     * @param ConstructorFN A reference to the class which instantiates the API
     * @param baseurl Path to resolve to reach the API
     *
     */
    addAPI: <GA extends APIBase>(apiName: string, ConstructorFN: new (avax: AvalancheCore, baseurl?: string, ...args: any[]) => GA, baseurl?: string, ...args: any[]) => void;
    /**
     * Retrieves a reference to an API by its apiName label.
     *
     * @param apiName Name of the API to return
     */
    api: <GA extends APIBase>(apiName: string) => GA;
    /**
     * @ignore
     */
    protected _request: (xhrmethod: Method, baseurl: string, getdata: object, postdata: string | object | ArrayBuffer | ArrayBufferView, headers?: AxiosRequestHeaders, axiosConfig?: AxiosRequestConfig) => Promise<RequestResponseData>;
    /**
     * Makes a GET call to an API.
     *
     * @param baseurl Path to the api
     * @param getdata Object containing the key value pairs sent in GET
     * @param headers An array HTTP Request Headers
     * @param axiosConfig Configuration for the axios javascript library that will be the
     * foundation for the rest of the parameters
     *
     * @returns A promise for [[RequestResponseData]]
     */
    get: (baseurl: string, getdata: object, headers?: object, axiosConfig?: AxiosRequestConfig) => Promise<RequestResponseData>;
    /**
     * Makes a DELETE call to an API.
     *
     * @param baseurl Path to the API
     * @param getdata Object containing the key value pairs sent in DELETE
     * @param headers An array HTTP Request Headers
     * @param axiosConfig Configuration for the axios javascript library that will be the
     * foundation for the rest of the parameters
     *
     * @returns A promise for [[RequestResponseData]]
     */
    delete: (baseurl: string, getdata: object, headers?: object, axiosConfig?: AxiosRequestConfig) => Promise<RequestResponseData>;
    /**
     * Makes a POST call to an API.
     *
     * @param baseurl Path to the API
     * @param getdata Object containing the key value pairs sent in POST
     * @param postdata Object containing the key value pairs sent in POST
     * @param headers An array HTTP Request Headers
     * @param axiosConfig Configuration for the axios javascript library that will be the
     * foundation for the rest of the parameters
     *
     * @returns A promise for [[RequestResponseData]]
     */
    post: (baseurl: string, getdata: object, postdata: string | object | ArrayBuffer | ArrayBufferView, headers?: object, axiosConfig?: AxiosRequestConfig) => Promise<RequestResponseData>;
    /**
     * Makes a PUT call to an API.
     *
     * @param baseurl Path to the baseurl
     * @param getdata Object containing the key value pairs sent in PUT
     * @param postdata Object containing the key value pairs sent in PUT
     * @param headers An array HTTP Request Headers
     * @param axiosConfig Configuration for the axios javascript library that will be the
     * foundation for the rest of the parameters
     *
     * @returns A promise for [[RequestResponseData]]
     */
    put: (baseurl: string, getdata: object, postdata: string | object | ArrayBuffer | ArrayBufferView, headers?: object, axiosConfig?: AxiosRequestConfig) => Promise<RequestResponseData>;
    /**
     * Makes a PATCH call to an API.
     *
     * @param baseurl Path to the baseurl
     * @param getdata Object containing the key value pairs sent in PATCH
     * @param postdata Object containing the key value pairs sent in PATCH
     * @param parameters Object containing the parameters of the API call
     * @param headers An array HTTP Request Headers
     * @param axiosConfig Configuration for the axios javascript library that will be the
     * foundation for the rest of the parameters
     *
     * @returns A promise for [[RequestResponseData]]
     */
    patch: (baseurl: string, getdata: object, postdata: string | object | ArrayBuffer | ArrayBufferView, headers?: object, axiosConfig?: AxiosRequestConfig) => Promise<RequestResponseData>;
    /**
     * Creates a new Avalanche instance. Sets the address and port of the main Avalanche Client.
     *
     * @param host The hostname to resolve to reach the Avalanche Client APIs
     * @param port The port to resolve to reach the Avalanche Client APIs
     * @param protocol The protocol string to use before a "://" in a request, ex: "http", "https", "git", "ws", etc ...
     */
    constructor(host?: string, port?: number, protocol?: string);
}
//# sourceMappingURL=avalanche.d.ts.map