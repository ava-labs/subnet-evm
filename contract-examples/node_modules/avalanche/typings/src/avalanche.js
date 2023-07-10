"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
/**
 * @packageDocumentation
 * @module AvalancheCore
 */
const axios_1 = __importDefault(require("axios"));
const apibase_1 = require("./common/apibase");
const errors_1 = require("./utils/errors");
const fetchadapter_1 = require("./utils/fetchadapter");
const helperfunctions_1 = require("./utils/helperfunctions");
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
class AvalancheCore {
    /**
     * Creates a new Avalanche instance. Sets the address and port of the main Avalanche Client.
     *
     * @param host The hostname to resolve to reach the Avalanche Client APIs
     * @param port The port to resolve to reach the Avalanche Client APIs
     * @param protocol The protocol string to use before a "://" in a request, ex: "http", "https", "git", "ws", etc ...
     */
    constructor(host, port, protocol = "http") {
        this.networkID = 0;
        this.hrp = "";
        this.auth = undefined;
        this.headers = {};
        this.requestConfig = {};
        this.apis = {};
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
        this.setAddress = (host, port, protocol = "http", baseEndpoint = "") => {
            host = host.replace(/[&#,@+()$~%'":*?<>{}]/g, "");
            protocol = protocol.replace(/[&#,@+()$~%'":*?<>{}]/g, "");
            const protocols = ["http", "https"];
            if (!protocols.includes(protocol)) {
                /* istanbul ignore next */
                throw new errors_1.ProtocolError("Error - AvalancheCore.setAddress: Invalid protocol");
            }
            this.host = host;
            this.port = port;
            this.protocol = protocol;
            this.baseEndpoint = baseEndpoint;
            let url = `${protocol}://${host}`;
            if (port != undefined && typeof port === "number" && port >= 0) {
                url = `${url}:${port}`;
            }
            if (baseEndpoint != undefined &&
                typeof baseEndpoint == "string" &&
                baseEndpoint.length > 0) {
                if (baseEndpoint[0] != "/") {
                    baseEndpoint = `/${baseEndpoint}`;
                }
                url = `${url}${baseEndpoint}`;
            }
            this.url = url;
        };
        /**
         * Returns the protocol such as "http", "https", "git", "ws", etc.
         */
        this.getProtocol = () => this.protocol;
        /**
         * Returns the host for the Avalanche node.
         */
        this.getHost = () => this.host;
        /**
         * Returns the IP for the Avalanche node.
         */
        this.getIP = () => this.host;
        /**
         * Returns the port for the Avalanche node.
         */
        this.getPort = () => this.port;
        /**
         * Returns the base endpoint for the Avalanche node.
         */
        this.getBaseEndpoint = () => this.baseEndpoint;
        /**
         * Returns the URL of the Avalanche node (ip + port)
         */
        this.getURL = () => this.url;
        /**
         * Returns the custom headers
         */
        this.getHeaders = () => this.headers;
        /**
         * Returns the custom request config
         */
        this.getRequestConfig = () => this.requestConfig;
        /**
         * Returns the networkID
         */
        this.getNetworkID = () => this.networkID;
        /**
         * Sets the networkID
         */
        this.setNetworkID = (netID) => {
            this.networkID = netID;
            this.hrp = (0, helperfunctions_1.getPreferredHRP)(this.networkID);
        };
        /**
         * Returns the Human-Readable-Part of the network associated with this key.
         *
         * @returns The [[KeyPair]]'s Human-Readable-Part of the network's Bech32 addressing scheme
         */
        this.getHRP = () => this.hrp;
        /**
         * Sets the the Human-Readable-Part of the network associated with this key.
         *
         * @param hrp String for the Human-Readable-Part of Bech32 addresses
         */
        this.setHRP = (hrp) => {
            this.hrp = hrp;
        };
        /**
         * Adds a new custom header to be included with all requests.
         *
         * @param key Header name
         * @param value Header value
         */
        this.setHeader = (key, value) => {
            this.headers[`${key}`] = value;
        };
        /**
         * Removes a previously added custom header.
         *
         * @param key Header name
         */
        this.removeHeader = (key) => {
            delete this.headers[`${key}`];
        };
        /**
         * Removes all headers.
         */
        this.removeAllHeaders = () => {
            for (const prop in this.headers) {
                if (Object.prototype.hasOwnProperty.call(this.headers, prop)) {
                    delete this.headers[`${prop}`];
                }
            }
        };
        /**
         * Adds a new custom config value to be included with all requests.
         *
         * @param key Config name
         * @param value Config value
         */
        this.setRequestConfig = (key, value) => {
            this.requestConfig[`${key}`] = value;
        };
        /**
         * Removes a previously added request config.
         *
         * @param key Header name
         */
        this.removeRequestConfig = (key) => {
            delete this.requestConfig[`${key}`];
        };
        /**
         * Removes all request configs.
         */
        this.removeAllRequestConfigs = () => {
            for (const prop in this.requestConfig) {
                if (Object.prototype.hasOwnProperty.call(this.requestConfig, prop)) {
                    delete this.requestConfig[`${prop}`];
                }
            }
        };
        /**
         * Sets the temporary auth token used for communicating with the node.
         *
         * @param auth A temporary token provided by the node enabling access to the endpoints on the node.
         */
        this.setAuthToken = (auth) => {
            this.auth = auth;
        };
        this._setHeaders = (headers) => {
            if (typeof this.headers === "object") {
                for (const [key, value] of Object.entries(this.headers)) {
                    headers[`${key}`] = value;
                }
            }
            if (typeof this.auth === "string") {
                headers.Authorization = `Bearer ${this.auth}`;
            }
            return headers;
        };
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
        this.addAPI = (apiName, ConstructorFN, baseurl = undefined, ...args) => {
            if (typeof baseurl === "undefined") {
                this.apis[`${apiName}`] = new ConstructorFN(this, undefined, ...args);
            }
            else {
                this.apis[`${apiName}`] = new ConstructorFN(this, baseurl, ...args);
            }
        };
        /**
         * Retrieves a reference to an API by its apiName label.
         *
         * @param apiName Name of the API to return
         */
        this.api = (apiName) => this.apis[`${apiName}`];
        /**
         * @ignore
         */
        this._request = (xhrmethod, baseurl, getdata, postdata, headers = {}, axiosConfig = undefined) => __awaiter(this, void 0, void 0, function* () {
            let config;
            if (axiosConfig) {
                config = Object.assign(Object.assign({}, axiosConfig), this.requestConfig);
            }
            else {
                config = Object.assign({ baseURL: this.url, responseType: "text" }, this.requestConfig);
            }
            config.url = baseurl;
            config.method = xhrmethod;
            config.headers = headers;
            config.data = postdata;
            config.params = getdata;
            // use the fetch adapter if fetch is available e.g. non Node<17 env
            if (typeof fetch !== "undefined") {
                config.adapter = fetchadapter_1.fetchAdapter;
            }
            const resp = yield axios_1.default.request(config);
            // purging all that is axios
            const xhrdata = new apibase_1.RequestResponseData(resp.data, resp.headers, resp.status, resp.statusText, resp.request);
            return xhrdata;
        });
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
        this.get = (baseurl, getdata, headers = {}, axiosConfig = undefined) => this._request("GET", baseurl, getdata, {}, this._setHeaders(headers), axiosConfig);
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
        this.delete = (baseurl, getdata, headers = {}, axiosConfig = undefined) => this._request("DELETE", baseurl, getdata, {}, this._setHeaders(headers), axiosConfig);
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
        this.post = (baseurl, getdata, postdata, headers = {}, axiosConfig = undefined) => this._request("POST", baseurl, getdata, postdata, this._setHeaders(headers), axiosConfig);
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
        this.put = (baseurl, getdata, postdata, headers = {}, axiosConfig = undefined) => this._request("PUT", baseurl, getdata, postdata, this._setHeaders(headers), axiosConfig);
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
        this.patch = (baseurl, getdata, postdata, headers = {}, axiosConfig = undefined) => this._request("PATCH", baseurl, getdata, postdata, this._setHeaders(headers), axiosConfig);
        if (host != undefined) {
            this.setAddress(host, port, protocol);
        }
    }
}
exports.default = AvalancheCore;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYXZhbGFuY2hlLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vc3JjL2F2YWxhbmNoZS50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7Ozs7OztBQUFBOzs7R0FHRztBQUNILGtEQUtjO0FBQ2QsOENBQStEO0FBQy9ELDJDQUE4QztBQUM5Qyx1REFBbUQ7QUFDbkQsNkRBQXlEO0FBRXpEOzs7Ozs7Ozs7R0FTRztBQUNILE1BQXFCLGFBQWE7SUF3YmhDOzs7Ozs7T0FNRztJQUNILFlBQVksSUFBYSxFQUFFLElBQWEsRUFBRSxXQUFtQixNQUFNO1FBOWJ6RCxjQUFTLEdBQVcsQ0FBQyxDQUFBO1FBQ3JCLFFBQUcsR0FBVyxFQUFFLENBQUE7UUFPaEIsU0FBSSxHQUFXLFNBQVMsQ0FBQTtRQUN4QixZQUFPLEdBQTRCLEVBQUUsQ0FBQTtRQUNyQyxrQkFBYSxHQUF1QixFQUFFLENBQUE7UUFDdEMsU0FBSSxHQUE2QixFQUFFLENBQUE7UUFFN0M7Ozs7Ozs7Ozs7O1dBV0c7UUFDSCxlQUFVLEdBQUcsQ0FDWCxJQUFZLEVBQ1osSUFBWSxFQUNaLFdBQW1CLE1BQU0sRUFDekIsZUFBdUIsRUFBRSxFQUNuQixFQUFFO1lBQ1IsSUFBSSxHQUFHLElBQUksQ0FBQyxPQUFPLENBQUMsd0JBQXdCLEVBQUUsRUFBRSxDQUFDLENBQUE7WUFDakQsUUFBUSxHQUFHLFFBQVEsQ0FBQyxPQUFPLENBQUMsd0JBQXdCLEVBQUUsRUFBRSxDQUFDLENBQUE7WUFDekQsTUFBTSxTQUFTLEdBQWEsQ0FBQyxNQUFNLEVBQUUsT0FBTyxDQUFDLENBQUE7WUFDN0MsSUFBSSxDQUFDLFNBQVMsQ0FBQyxRQUFRLENBQUMsUUFBUSxDQUFDLEVBQUU7Z0JBQ2pDLDBCQUEwQjtnQkFDMUIsTUFBTSxJQUFJLHNCQUFhLENBQ3JCLG9EQUFvRCxDQUNyRCxDQUFBO2FBQ0Y7WUFFRCxJQUFJLENBQUMsSUFBSSxHQUFHLElBQUksQ0FBQTtZQUNoQixJQUFJLENBQUMsSUFB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Q1osQ0FBQTtRQVVELElBQUksSUFBSSxJQUFJLFNBQVMsRUFBRTtZQUNyQixJQUFJLENBQUMsVUFBVSxDQUFDLElBQUksRUFBRSxJQUFJLEVBQUUsUUFBUSxDQUFDLENBQUE7U0FDdEM7SUFDSCxDQUFDO0NBQ0Y7QUFwY0QsZ0NBb2NDIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgQXZhbGFuY2hlQ29yZVxuICovXG5pbXBvcnQgYXhpb3MsIHtcbiAgQXhpb3NSZXF1ZXN0Q29uZmlnLFxuICBBeGlvc1JlcXVlc3RIZWFkZXJzLFxuICBBeGlvc1Jlc3BvbnNlLFxuICBNZXRob2Rcbn0gZnJvbSBcImF4aW9zXCJcbmltcG9ydCB7IEFQSUJhc2UsIFJlcXVlc3RSZXNwb25zZURhdGEgfSBmcm9tIFwiLi9jb21tb24vYXBpYmFzZVwiXG5pbXBvcnQgeyBQcm90b2NvbEVycm9yIH0gZnJvbSBcIi4vdXRpbHMvZXJyb3JzXCJcbmltcG9ydCB7IGZldGNoQWRhcHRlciB9IGZyb20gXCIuL3V0aWxzL2ZldGNoYWRhcHRlclwiXG5pbXBvcnQgeyBnZXRQcmVmZXJyZWRIUlAgfSBmcm9tIFwiLi91dGlscy9oZWxwZXJmdW5jdGlvbnNcIlxuXG4vKipcbiAqIEF2YWxhbmNoZUNvcmUgaXMgbWlkZGxld2FyZSBmb3IgaW50ZXJhY3Rpbmcgd2l0aCBBdmFsYW5jaGUgbm9kZSBSUEMgQVBJcy5cbiAqXG4gKiBFeGFtcGxlIHVzYWdlOlxuICogYGBganNcbiAqIGxldCBhdmFsYW5jaGUgPSBuZXcgQXZhbGFuY2hlQ29yZShcIjEyNy4wLjAuMVwiLCA5NjUwLCBcImh0dHBzXCIpXG4gKiBgYGBcbiAqXG4gKlxuICovXG5leHBvcnQgZGVmYXVsdCBjbGFzcyBBdmFsYW5jaGVDb3JlIHtcbiAgcHJvdGVjdGVkIG5ldHdvcmtJRDogbnVtYmVyID0gMFxuICBwcm90ZWN0ZWQgaHJwOiBzdHJpbmcgPSBcIlwiXG4gIHByb3RlY3RlZCBwcm90b2NvbDogc3RyaW5nXG4gIHByb3RlY3RlZCBpcDogc3RyaW5nXG4gIHByb3RlY3RlZCBob3N0OiBzdHJpbmdcbiAgcHJvdGVjdGVkIHBvcnQ6IG51bWJlclxuICBwcm90ZWN0ZWQgYmFzZUVuZHBvaW50OiBzdHJpbmdcbiAgcHJvdGVjdGVkIHVybDogc3RyaW5nXG4gIHByb3RlY3RlZCBhdXRoOiBzdHJpbmcgPSB1bmRlZmluZWRcbiAgcHJvdGVjdGVkIGhlYWRlcnM6IHsgW2s6IHN0cmluZ106IHN0cmluZyB9ID0ge31cbiAgcHJvdGVjdGVkIHJlcXVlc3RDb25maWc6IEF4aW9zUmVxdWVzdENvbmZpZyA9IHt9XG4gIHByb3RlY3RlZCBhcGlzOiB7IFtrOiBzdHJpbmddOiBBUElCYXNlIH0gPSB7fVxuXG4gIC8qKlxuICAgKiBTZXRzIHRoZSBhZGRyZXNzIGFuZCBwb3J0IG9mIHRoZSBtYWluIEF2YWxhbmNoZSBDbGllbnQuXG4gICAqXG4gICAqIEBwYXJhbSBob3N0IFRoZSBob3N0bmFtZSB0byByZXNvbHZlIHRvIHJlYWNoIHRoZSBBdmFsYW5jaGUgQ2xpZW50IFJQQyBBUElzLlxuICAgKiBAcGFyYW0gcG9ydCBUaGUgcG9ydCB0byByZXNvbHZlIHRvIHJlYWNoIHRoZSBBdmFsYW5jaGUgQ2xpZW50IFJQQyBBUElzLlxuICAgKiBAcGFyYW0gcHJvdG9jb2wgVGhlIHByb3RvY29sIHN0cmluZyB0byB1c2UgYmVmb3JlIGEgXCI6Ly9cIiBpbiBhIHJlcXVlc3QsXG4gICAqIGV4OiBcImh0dHBcIiwgXCJodHRwc1wiLCBldGMuIERlZmF1bHRzIHRvIGh0dHBcbiAgICogQHBhcmFtIGJhc2VFbmRwb2ludCB0aGUgYmFzZSBlbmRwb2ludCB0byByZWFjaCB0aGUgQXZhbGFuY2hlIENsaWVudCBSUEMgQVBJcyxcbiAgICogZXg6IFwiL3JwY1wiLiBEZWZhdWx0cyB0byBcIi9cIlxuICAgKiBUaGUgZm9sbG93aW5nIHNwZWNpYWwgY2hhcmFjdGVycyBhcmUgcmVtb3ZlZCBmcm9tIGhvc3QgYW5kIHByb3RvY29sXG4gICAqICYjLEArKCkkfiUnXCI6Kj97fSBhbHNvIGxlc3MgdGhhbiBhbmQgZ3JlYXRlciB0aGFuIHNpZ25zXG4gICAqL1xuICBzZXRBZGRyZXNzID0gKFxuICAgIGhvc3Q6IHN0cmluZyxcbiAgICBwb3J0OiBudW1iZXIsXG4gICAgcHJvdG9jb2w6IHN0cmluZyA9IFwiaHR0cFwiLFxuICAgIGJhc2VFbmRwb2ludDogc3RyaW5nID0gXCJcIlxuICApOiB2b2lkID0+IHtcbiAgICBob3N0ID0gaG9zdC5yZXBsYWNlKC9bJiMsQCsoKSR+JSdcIjoqPzw+e31dL2csIFwiXCIpXG4gICAgcHJvdG9jb2wgPSBwcm90b2NvbC5yZXBsYWNlKC9bJiMsQCsoKSR+JSdcIjoqPzw+e31dL2csIFwiXCIpXG4gICAgY29uc3QgcHJvdG9jb2xzOiBzdHJpbmdbXSA9IFtcImh0dHBcIiwgXCJodHRwc1wiXVxuICAgIGlmICghcHJvdG9jb2xzLmluY2x1ZGVzKHByb3RvY29sKSkge1xuICAgICAgLyogaXN0YW5idWwgaWdub3JlIG5leHQgKi9cbiAgICAgIHRocm93IG5ldyBQcm90b2NvbEVycm9yKFxuICAgICAgICBcIkVycm9yIC0gQXZhbGFuY2hlQ29yZS5zZXRBZGRyZXNzOiBJbnZhbGlkIHByb3RvY29sXCJcbiAgICAgIClcbiAgICB9XG5cbiAgICB0aGlzLmhvc3QgPSBob3N0XG4gICAgdGhpcy5wb3J0ID0gcG9ydFxuICAgIHRoaXMucHJvdG9jb2wgPSBwcm90b2NvbFxuICAgIHRoaXMuYmFzZUVuZHBvaW50ID0gYmFzZUVuZHBvaW50XG4gICAgbGV0IHVybDogc3RyaW5nID0gYCR7cHJvdG9jb2x9Oi8vJHtob3N0fWBcbiAgICBpZiAocG9ydCAhPSB1bmRlZmluZWQgJiYgdHlwZW9mIHBvcnQgPT09IFwibnVtYmVyXCIgJiYgcG9ydCA+PSAwKSB7XG4gICAgICB1cmwgPSBgJHt1cmx9OiR7cG9ydH1gXG4gICAgfVxuICAgIGlmIChcbiAgICAgIGJhc2VFbmRwb2ludCAhPSB1bmRlZmluZWQgJiZcbiAgICAgIHR5cGVvZiBiYXNlRW5kcG9pbnQgPT0gXCJzdHJpbmdcIiAmJlxuICAgICAgYmFzZUVuZHBvaW50Lmxlbmd0aCA+IDBcbiAgICApIHtcbiAgICAgIGlmIChiYXNlRW5kcG9pbnRbMF0gIT0gXCIvXCIpIHtcbiAgICAgICAgYmFzZUVuZHBvaW50ID0gYC8ke2Jhc2VFbmRwb2ludH1gXG4gICAgICB9XG4gICAgICB1cmwgPSBgJHt1cmx9JHtiYXNlRW5kcG9pbnR9YFxuICAgIH1cbiAgICB0aGlzLnVybCA9IHVybFxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIHByb3RvY29sIHN1Y2ggYXMgXCJodHRwXCIsIFwiaHR0cHNcIiwgXCJnaXRcIiwgXCJ3c1wiLCBldGMuXG4gICAqL1xuICBnZXRQcm90b2NvbCA9ICgpOiBzdHJpbmcgPT4gdGhpcy5wcm90b2NvbFxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBob3N0IGZvciB0aGUgQXZhbGFuY2hlIG5vZGUuXG4gICAqL1xuICBnZXRIb3N0ID0gKCk6IHN0cmluZyA9PiB0aGlzLmhvc3RcblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgSVAgZm9yIHRoZSBBdmFsYW5jaGUgbm9kZS5cbiAgICovXG4gIGdldElQID0gKCk6IHN0cmluZyA9PiB0aGlzLmhvc3RcblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgcG9ydCBmb3IgdGhlIEF2YWxhbmNoZSBub2RlLlxuICAgKi9cbiAgZ2V0UG9ydCA9ICgpOiBudW1iZXIgPT4gdGhpcy5wb3J0XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIGJhc2UgZW5kcG9pbnQgZm9yIHRoZSBBdmFsYW5jaGUgbm9kZS5cbiAgICovXG4gIGdldEJhc2VFbmRwb2ludCA9ICgpOiBzdHJpbmcgPT4gdGhpcy5iYXNlRW5kcG9pbnRcblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgVVJMIG9mIHRoZSBBdmFsYW5jaGUgbm9kZSAoaXAgKyBwb3J0KVxuICAgKi9cbiAgZ2V0VVJMID0gKCk6IHN0cmluZyA9PiB0aGlzLnVybFxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBjdXN0b20gaGVhZGVyc1xuICAgKi9cbiAgZ2V0SGVhZGVycyA9ICgpOiBvYmplY3QgPT4gdGhpcy5oZWFkZXJzXG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIGN1c3RvbSByZXF1ZXN0IGNvbmZpZ1xuICAgKi9cbiAgZ2V0UmVxdWVzdENvbmZpZyA9ICgpOiBBeGlvc1JlcXVlc3RDb25maWcgPT4gdGhpcy5yZXF1ZXN0Q29uZmlnXG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIG5ldHdvcmtJRFxuICAgKi9cbiAgZ2V0TmV0d29ya0lEID0gKCk6IG51bWJlciA9PiB0aGlzLm5ldHdvcmtJRFxuXG4gIC8qKlxuICAgKiBTZXRzIHRoZSBuZXR3b3JrSURcbiAgICovXG4gIHNldE5ldHdvcmtJRCA9IChuZXRJRDogbnVtYmVyKTogdm9pZCA9PiB7XG4gICAgdGhpcy5uZXR3b3JrSUQgPSBuZXRJRFxuICAgIHRoaXMuaHJwID0gZ2V0UHJlZmVycmVkSFJQKHRoaXMubmV0d29ya0lEKVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIEh1bWFuLVJlYWRhYmxlLVBhcnQgb2YgdGhlIG5ldHdvcmsgYXNzb2NpYXRlZCB3aXRoIHRoaXMga2V5LlxuICAgKlxuICAgKiBAcmV0dXJucyBUaGUgW1tLZXlQYWlyXV0ncyBIdW1hbi1SZWFkYWJsZS1QYXJ0IG9mIHRoZSBuZXR3b3JrJ3MgQmVjaDMyIGFkZHJlc3Npbmcgc2NoZW1lXG4gICAqL1xuICBnZXRIUlAgPSAoKTogc3RyaW5nID0+IHRoaXMuaHJwXG5cbiAgLyoqXG4gICAqIFNldHMgdGhlIHRoZSBIdW1hbi1SZWFkYWJsZS1QYXJ0IG9mIHRoZSBuZXR3b3JrIGFzc29jaWF0ZWQgd2l0aCB0aGlzIGtleS5cbiAgICpcbiAgICogQHBhcmFtIGhycCBTdHJpbmcgZm9yIHRoZSBIdW1hbi1SZWFkYWJsZS1QYXJ0IG9mIEJlY2gzMiBhZGRyZXNzZXNcbiAgICovXG4gIHNldEhSUCA9IChocnA6IHN0cmluZyk6IHZvaWQgPT4ge1xuICAgIHRoaXMuaHJwID0gaHJwXG4gIH1cblxuICAvKipcbiAgICogQWRkcyBhIG5ldyBjdXN0b20gaGVhZGVyIHRvIGJlIGluY2x1ZGVkIHdpdGggYWxsIHJlcXVlc3RzLlxuICAgKlxuICAgKiBAcGFyYW0ga2V5IEhlYWRlciBuYW1lXG4gICAqIEBwYXJhbSB2YWx1ZSBIZWFkZXIgdmFsdWVcbiAgICovXG4gIHNldEhlYWRlciA9IChrZXk6IHN0cmluZywgdmFsdWU6IHN0cmluZyk6IHZvaWQgPT4ge1xuICAgIHRoaXMuaGVhZGVyc1tgJHtrZXl9YF0gPSB2YWx1ZVxuICB9XG5cbiAgLyoqXG4gICAqIFJlbW92ZXMgYSBwcmV2aW91c2x5IGFkZGVkIGN1c3RvbSBoZWFkZXIuXG4gICAqXG4gICAqIEBwYXJhbSBrZXkgSGVhZGVyIG5hbWVcbiAgICovXG4gIHJlbW92ZUhlYWRlciA9IChrZXk6IHN0cmluZyk6IHZvaWQgPT4ge1xuICAgIGRlbGV0ZSB0aGlzLmhlYWRlcnNbYCR7a2V5fWBdXG4gIH1cblxuICAvKipcbiAgICogUmVtb3ZlcyBhbGwgaGVhZGVycy5cbiAgICovXG4gIHJlbW92ZUFsbEhlYWRlcnMgPSAoKTogdm9pZCA9PiB7XG4gICAgZm9yIChjb25zdCBwcm9wIGluIHRoaXMuaGVhZGVycykge1xuICAgICAgaWYgKE9iamVjdC5wcm90b3R5cGUuaGFzT3duUHJvcGVydHkuY2FsbCh0aGlzLmhlYWRlcnMsIHByb3ApKSB7XG4gICAgICAgIGRlbGV0ZSB0aGlzLmhlYWRlcnNbYCR7cHJvcH1gXVxuICAgICAgfVxuICAgIH1cbiAgfVxuXG4gIC8qKlxuICAgKiBBZGRzIGEgbmV3IGN1c3RvbSBjb25maWcgdmFsdWUgdG8gYmUgaW5jbHVkZWQgd2l0aCBhbGwgcmVxdWVzdHMuXG4gICAqXG4gICAqIEBwYXJhbSBrZXkgQ29uZmlnIG5hbWVcbiAgICogQHBhcmFtIHZhbHVlIENvbmZpZyB2YWx1ZVxuICAgKi9cbiAgc2V0UmVxdWVzdENvbmZpZyA9IChrZXk6IHN0cmluZywgdmFsdWU6IHN0cmluZyB8IGJvb2xlYW4pOiB2b2lkID0+IHtcbiAgICB0aGlzLnJlcXVlc3RDb25maWdbYCR7a2V5fWBdID0gdmFsdWVcbiAgfVxuXG4gIC8qKlxuICAgKiBSZW1vdmVzIGEgcHJldmlvdXNseSBhZGRlZCByZXF1ZXN0IGNvbmZpZy5cbiAgICpcbiAgICogQHBhcmFtIGtleSBIZWFkZXIgbmFtZVxuICAgKi9cbiAgcmVtb3ZlUmVxdWVzdENvbmZpZyA9IChrZXk6IHN0cmluZyk6IHZvaWQgPT4ge1xuICAgIGRlbGV0ZSB0aGlzLnJlcXVlc3RDb25maWdbYCR7a2V5fWBdXG4gIH1cblxuICAvKipcbiAgICogUmVtb3ZlcyBhbGwgcmVxdWVzdCBjb25maWdzLlxuICAgKi9cbiAgcmVtb3ZlQWxsUmVxdWVzdENvbmZpZ3MgPSAoKTogdm9pZCA9PiB7XG4gICAgZm9yIChjb25zdCBwcm9wIGluIHRoaXMucmVxdWVzdENvbmZpZykge1xuICAgICAgaWYgKE9iamVjdC5wcm90b3R5cGUuaGFzT3duUHJvcGVydHkuY2FsbCh0aGlzLnJlcXVlc3RDb25maWcsIHByb3ApKSB7XG4gICAgICAgIGRlbGV0ZSB0aGlzLnJlcXVlc3RDb25maWdbYCR7cHJvcH1gXVxuICAgICAgfVxuICAgIH1cbiAgfVxuXG4gIC8qKlxuICAgKiBTZXRzIHRoZSB0ZW1wb3JhcnkgYXV0aCB0b2tlbiB1c2VkIGZvciBjb21tdW5pY2F0aW5nIHdpdGggdGhlIG5vZGUuXG4gICAqXG4gICAqIEBwYXJhbSBhdXRoIEEgdGVtcG9yYXJ5IHRva2VuIHByb3ZpZGVkIGJ5IHRoZSBub2RlIGVuYWJsaW5nIGFjY2VzcyB0byB0aGUgZW5kcG9pbnRzIG9uIHRoZSBub2RlLlxuICAgKi9cbiAgc2V0QXV0aFRva2VuID0gKGF1dGg6IHN0cmluZyk6IHZvaWQgPT4ge1xuICAgIHRoaXMuYXV0aCA9IGF1dGhcbiAgfVxuXG4gIHByb3RlY3RlZCBfc2V0SGVhZGVycyA9IChoZWFkZXJzOiBhbnkpOiBBeGlvc1JlcXVlc3RIZWFkZXJzID0+IHtcbiAgICBpZiAodHlwZW9mIHRoaXMuaGVhZGVycyA9PT0gXCJvYmplY3RcIikge1xuICAgICAgZm9yIChjb25zdCBba2V5LCB2YWx1ZV0gb2YgT2JqZWN0LmVudHJpZXModGhpcy5oZWFkZXJzKSkge1xuICAgICAgICBoZWFkZXJzW2Ake2tleX1gXSA9IHZhbHVlXG4gICAgICB9XG4gICAgfVxuXG4gICAgaWYgKHR5cGVvZiB0aGlzLmF1dGggPT09IFwic3RyaW5nXCIpIHtcbiAgICAgIGhlYWRlcnMuQXV0aG9yaXphdGlvbiA9IGBCZWFyZXIgJHt0aGlzLmF1dGh9YFxuICAgIH1cbiAgICByZXR1cm4gaGVhZGVyc1xuICB9XG5cbiAgLyoqXG4gICAqIEFkZHMgYW4gQVBJIHRvIHRoZSBtaWRkbGV3YXJlLiBUaGUgQVBJIHJlc29sdmVzIHRvIGEgcmVnaXN0ZXJlZCBibG9ja2NoYWluJ3MgUlBDLlxuICAgKlxuICAgKiBJbiBUeXBlU2NyaXB0OlxuICAgKiBgYGBqc1xuICAgKiBhdmFsYW5jaGUuYWRkQVBJPE15Vk1DbGFzcz4oXCJteWNoYWluXCIsIE15Vk1DbGFzcywgXCIvZXh0L2JjL215Y2hhaW5cIilcbiAgICogYGBgXG4gICAqXG4gICAqIEluIEphdmFzY3JpcHQ6XG4gICAqIGBgYGpzXG4gICAqIGF2YWxhbmNoZS5hZGRBUEkoXCJteWNoYWluXCIsIE15Vk1DbGFzcywgXCIvZXh0L2JjL215Y2hhaW5cIilcbiAgICogYGBgXG4gICAqXG4gICAqIEB0eXBlcGFyYW0gR0EgQ2xhc3Mgb2YgdGhlIEFQSSBiZWluZyBhZGRlZFxuICAgKiBAcGFyYW0gYXBpTmFtZSBBIGxhYmVsIGZvciByZWZlcmVuY2luZyB0aGUgQVBJIGluIHRoZSBmdXR1cmVcbiAgICogQHBhcmFtIENvbnN0cnVjdG9yRk4gQSByZWZlcmVuY2UgdG8gdGhlIGNsYXNzIHdoaWNoIGluc3RhbnRpYXRlcyB0aGUgQVBJXG4gICAqIEBwYXJhbSBiYXNldXJsIFBhdGggdG8gcmVzb2x2ZSB0byByZWFjaCB0aGUgQVBJXG4gICAqXG4gICAqL1xuICBhZGRBUEkgPSA8R0EgZXh0ZW5kcyBBUElCYXNlPihcbiAgICBhcGlOYW1lOiBzdHJpbmcsXG4gICAgQ29uc3RydWN0b3JGTjogbmV3IChcbiAgICAgIGF2YXg6IEF2YWxhbmNoZUNvcmUsXG4gICAgICBiYXNldXJsPzogc3RyaW5nLFxuICAgICAgLi4uYXJnczogYW55W11cbiAgICApID0+IEdBLFxuICAgIGJhc2V1cmw6IHN0cmluZyA9IHVuZGVmaW5lZCxcbiAgICAuLi5hcmdzOiBhbnlbXVxuICApID0+IHtcbiAgICBpZiAodHlwZW9mIGJhc2V1cmwgPT09IFwidW5kZWZpbmVkXCIpIHtcbiAgICAgIHRoaXMuYXBpc1tgJHthcGlOYW1lfWBdID0gbmV3IENvbnN0cnVjdG9yRk4odGhpcywgdW5kZWZpbmVkLCAuLi5hcmdzKVxuICAgIH0gZWxzZSB7XG4gICAgICB0aGlzLmFwaXNbYCR7YXBpTmFtZX1gXSA9IG5ldyBDb25zdHJ1Y3RvckZOKHRoaXMsIGJhc2V1cmwsIC4uLmFyZ3MpXG4gICAgfVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHJpZXZlcyBhIHJlZmVyZW5jZSB0byBhbiBBUEkgYnkgaXRzIGFwaU5hbWUgbGFiZWwuXG4gICAqXG4gICAqIEBwYXJhbSBhcGlOYW1lIE5hbWUgb2YgdGhlIEFQSSB0byByZXR1cm5cbiAgICovXG4gIGFwaSA9IDxHQSBleHRlbmRzIEFQSUJhc2U+KGFwaU5hbWU6IHN0cmluZyk6IEdBID0+XG4gICAgdGhpcy5hcGlzW2Ake2FwaU5hbWV9YF0gYXMgR0FcblxuICAvKipcbiAgICogQGlnbm9yZVxuICAgKi9cbiAgcHJvdGVjdGVkIF9yZXF1ZXN0ID0gYXN5bmMgKFxuICAgIHhocm1ldGhvZDogTWV0aG9kLFxuICAgIGJhc2V1cmw6IHN0cmluZyxcbiAgICBnZXRkYXRhOiBvYmplY3QsXG4gICAgcG9zdGRhdGE6IHN0cmluZyB8IG9iamVjdCB8IEFycmF5QnVmZmVyIHwgQXJyYXlCdWZmZXJWaWV3LFxuICAgIGhlYWRlcnM6IEF4aW9zUmVxdWVzdEhlYWRlcnMgPSB7fSxcbiAgICBheGlvc0NvbmZpZzogQXhpb3NSZXF1ZXN0Q29uZmlnID0gdW5kZWZpbmVkXG4gICk6IFByb21pc2U8UmVxdWVzdFJlc3BvbnNlRGF0YT4gPT4ge1xuICAgIGxldCBjb25maWc6IEF4aW9zUmVxdWVzdENvbmZpZ1xuICAgIGlmIChheGlvc0NvbmZpZykge1xuICAgICAgY29uZmlnID0ge1xuICAgICAgICAuLi5heGlvc0NvbmZpZyxcbiAgICAgICAgLi4udGhpcy5yZXF1ZXN0Q29uZmlnXG4gICAgICB9XG4gICAgfSBlbHNlIHtcbiAgICAgIGNvbmZpZyA9IHtcbiAgICAgICAgYmFzZVVSTDogdGhpcy51cmwsXG4gICAgICAgIHJlc3BvbnNlVHlwZTogXCJ0ZXh0XCIsXG4gICAgICAgIC4uLnRoaXMucmVxdWVzdENvbmZpZ1xuICAgICAgfVxuICAgIH1cbiAgICBjb25maWcudXJsID0gYmFzZXVybFxuICAgIGNvbmZpZy5tZXRob2QgPSB4aHJtZXRob2RcbiAgICBjb25maWcuaGVhZGVycyA9IGhlYWRlcnNcbiAgICBjb25maWcuZGF0YSA9IHBvc3RkYXRhXG4gICAgY29uZmlnLnBhcmFtcyA9IGdldGRhdGFcbiAgICAvLyB1c2UgdGhlIGZldGNoIGFkYXB0ZXIgaWYgZmV0Y2ggaXMgYXZhaWxhYmxlIGUuZy4gbm9uIE5vZGU8MTcgZW52XG4gICAgaWYgKHR5cGVvZiBmZXRjaCAhPT0gXCJ1bmRlZmluZWRcIikge1xuICAgICAgY29uZmlnLmFkYXB0ZXIgPSBmZXRjaEFkYXB0ZXJcbiAgICB9XG4gICAgY29uc3QgcmVzcDogQXhpb3NSZXNwb25zZTxhbnk+ID0gYXdhaXQgYXhpb3MucmVxdWVzdChjb25maWcpXG4gICAgLy8gcHVyZ2luZyBhbGwgdGhhdCBpcyBheGlvc1xuICAgIGNvbnN0IHhocmRhdGE6IFJlcXVlc3RSZXNwb25zZURhdGEgPSBuZXcgUmVxdWVzdFJlc3BvbnNlRGF0YShcbiAgICAgIHJlc3AuZGF0YSxcbiAgICAgIHJlc3AuaGVhZGVycyxcbiAgICAgIHJlc3Auc3RhdHVzLFxuICAgICAgcmVzcC5zdGF0dXNUZXh0LFxuICAgICAgcmVzcC5yZXF1ZXN0XG4gICAgKVxuICAgIHJldHVybiB4aHJkYXRhXG4gIH1cblxuICAvKipcbiAgICogTWFrZXMgYSBHRVQgY2FsbCB0byBhbiBBUEkuXG4gICAqXG4gICAqIEBwYXJhbSBiYXNldXJsIFBhdGggdG8gdGhlIGFwaVxuICAgKiBAcGFyYW0gZ2V0ZGF0YSBPYmplY3QgY29udGFpbmluZyB0aGUga2V5IHZhbHVlIHBhaXJzIHNlbnQgaW4gR0VUXG4gICAqIEBwYXJhbSBoZWFkZXJzIEFuIGFycmF5IEhUVFAgUmVxdWVzdCBIZWFkZXJzXG4gICAqIEBwYXJhbSBheGlvc0NvbmZpZyBDb25maWd1cmF0aW9uIGZvciB0aGUgYXhpb3MgamF2YXNjcmlwdCBsaWJyYXJ5IHRoYXQgd2lsbCBiZSB0aGVcbiAgICogZm91bmRhdGlvbiBmb3IgdGhlIHJlc3Qgb2YgdGhlIHBhcmFtZXRlcnNcbiAgICpcbiAgICogQHJldHVybnMgQSBwcm9taXNlIGZvciBbW1JlcXVlc3RSZXNwb25zZURhdGFdXVxuICAgKi9cbiAgZ2V0ID0gKFxuICAgIGJhc2V1cmw6IHN0cmluZyxcbiAgICBnZXRkYXRhOiBvYmplY3QsXG4gICAgaGVhZGVyczogb2JqZWN0ID0ge30sXG4gICAgYXhpb3NDb25maWc6IEF4aW9zUmVxdWVzdENvbmZpZyA9IHVuZGVmaW5lZFxuICApOiBQcm9taXNlPFJlcXVlc3RSZXNwb25zZURhdGE+ID0+XG4gICAgdGhpcy5fcmVxdWVzdChcbiAgICAgIFwiR0VUXCIsXG4gICAgICBiYXNldXJsLFxuICAgICAgZ2V0ZGF0YSxcbiAgICAgIHt9LFxuICAgICAgdGhpcy5fc2V0SGVhZGVycyhoZWFkZXJzKSxcbiAgICAgIGF4aW9zQ29uZmlnXG4gICAgKVxuXG4gIC8qKlxuICAgKiBNYWtlcyBhIERFTEVURSBjYWxsIHRvIGFuIEFQSS5cbiAgICpcbiAgICogQHBhcmFtIGJhc2V1cmwgUGF0aCB0byB0aGUgQVBJXG4gICAqIEBwYXJhbSBnZXRkYXRhIE9iamVjdCBjb250YWluaW5nIHRoZSBrZXkgdmFsdWUgcGFpcnMgc2VudCBpbiBERUxFVEVcbiAgICogQHBhcmFtIGhlYWRlcnMgQW4gYXJyYXkgSFRUUCBSZXF1ZXN0IEhlYWRlcnNcbiAgICogQHBhcmFtIGF4aW9zQ29uZmlnIENvbmZpZ3VyYXRpb24gZm9yIHRoZSBheGlvcyBqYXZhc2NyaXB0IGxpYnJhcnkgdGhhdCB3aWxsIGJlIHRoZVxuICAgKiBmb3VuZGF0aW9uIGZvciB0aGUgcmVzdCBvZiB0aGUgcGFyYW1ldGVyc1xuICAgKlxuICAgKiBAcmV0dXJucyBBIHByb21pc2UgZm9yIFtbUmVxdWVzdFJlc3BvbnNlRGF0YV1dXG4gICAqL1xuICBkZWxldGUgPSAoXG4gICAgYmFzZXVybDogc3RyaW5nLFxuICAgIGdldGRhdGE6IG9iamVjdCxcbiAgICBoZWFkZXJzOiBvYmplY3QgPSB7fSxcbiAgICBheGlvc0NvbmZpZzogQXhpb3NSZXF1ZXN0Q29uZmlnID0gdW5kZWZpbmVkXG4gICk6IFByb21pc2U8UmVxdWVzdFJlc3BvbnNlRGF0YT4gPT5cbiAgICB0aGlzLl9yZXF1ZXN0KFxuICAgICAgXCJERUxFVEVcIixcbiAgICAgIGJhc2V1cmwsXG4gICAgICBnZXRkYXRhLFxuICAgICAge30sXG4gICAgICB0aGlzLl9zZXRIZWFkZXJzKGhlYWRlcnMpLFxuICAgICAgYXhpb3NDb25maWdcbiAgICApXG5cbiAgLyoqXG4gICAqIE1ha2VzIGEgUE9TVCBjYWxsIHRvIGFuIEFQSS5cbiAgICpcbiAgICogQHBhcmFtIGJhc2V1cmwgUGF0aCB0byB0aGUgQVBJXG4gICAqIEBwYXJhbSBnZXRkYXRhIE9iamVjdCBjb250YWluaW5nIHRoZSBrZXkgdmFsdWUgcGFpcnMgc2VudCBpbiBQT1NUXG4gICAqIEBwYXJhbSBwb3N0ZGF0YSBPYmplY3QgY29udGFpbmluZyB0aGUga2V5IHZhbHVlIHBhaXJzIHNlbnQgaW4gUE9TVFxuICAgKiBAcGFyYW0gaGVhZGVycyBBbiBhcnJheSBIVFRQIFJlcXVlc3QgSGVhZGVyc1xuICAgKiBAcGFyYW0gYXhpb3NDb25maWcgQ29uZmlndXJhdGlvbiBmb3IgdGhlIGF4aW9zIGphdmFzY3JpcHQgbGlicmFyeSB0aGF0IHdpbGwgYmUgdGhlXG4gICAqIGZvdW5kYXRpb24gZm9yIHRoZSByZXN0IG9mIHRoZSBwYXJhbWV0ZXJzXG4gICAqXG4gICAqIEByZXR1cm5zIEEgcHJvbWlzZSBmb3IgW1tSZXF1ZXN0UmVzcG9uc2VEYXRhXV1cbiAgICovXG4gIHBvc3QgPSAoXG4gICAgYmFzZXVybDogc3RyaW5nLFxuICAgIGdldGRhdGE6IG9iamVjdCxcbiAgICBwb3N0ZGF0YTogc3RyaW5nIHwgb2JqZWN0IHwgQXJyYXlCdWZmZXIgfCBBcnJheUJ1ZmZlclZpZXcsXG4gICAgaGVhZGVyczogb2JqZWN0ID0ge30sXG4gICAgYXhpb3NDb25maWc6IEF4aW9zUmVxdWVzdENvbmZpZyA9IHVuZGVmaW5lZFxuICApOiBQcm9taXNlPFJlcXVlc3RSZXNwb25zZURhdGE+ID0+XG4gICAgdGhpcy5fcmVxdWVzdChcbiAgICAgIFwiUE9TVFwiLFxuICAgICAgYmFzZXVybCxcbiAgICAgIGdldGRhdGEsXG4gICAgICBwb3N0ZGF0YSxcbiAgICAgIHRoaXMuX3NldEhlYWRlcnMoaGVhZGVycyksXG4gICAgICBheGlvc0NvbmZpZ1xuICAgIClcblxuICAvKipcbiAgICogTWFrZXMgYSBQVVQgY2FsbCB0byBhbiBBUEkuXG4gICAqXG4gICAqIEBwYXJhbSBiYXNldXJsIFBhdGggdG8gdGhlIGJhc2V1cmxcbiAgICogQHBhcmFtIGdldGRhdGEgT2JqZWN0IGNvbnRhaW5pbmcgdGhlIGtleSB2YWx1ZSBwYWlycyBzZW50IGluIFBVVFxuICAgKiBAcGFyYW0gcG9zdGRhdGEgT2JqZWN0IGNvbnRhaW5pbmcgdGhlIGtleSB2YWx1ZSBwYWlycyBzZW50IGluIFBVVFxuICAgKiBAcGFyYW0gaGVhZGVycyBBbiBhcnJheSBIVFRQIFJlcXVlc3QgSGVhZGVyc1xuICAgKiBAcGFyYW0gYXhpb3NDb25maWcgQ29uZmlndXJhdGlvbiBmb3IgdGhlIGF4aW9zIGphdmFzY3JpcHQgbGlicmFyeSB0aGF0IHdpbGwgYmUgdGhlXG4gICAqIGZvdW5kYXRpb24gZm9yIHRoZSByZXN0IG9mIHRoZSBwYXJhbWV0ZXJzXG4gICAqXG4gICAqIEByZXR1cm5zIEEgcHJvbWlzZSBmb3IgW1tSZXF1ZXN0UmVzcG9uc2VEYXRhXV1cbiAgICovXG4gIHB1dCA9IChcbiAgICBiYXNldXJsOiBzdHJpbmcsXG4gICAgZ2V0ZGF0YTogb2JqZWN0LFxuICAgIHBvc3RkYXRhOiBzdHJpbmcgfCBvYmplY3QgfCBBcnJheUJ1ZmZlciB8IEFycmF5QnVmZmVyVmlldyxcbiAgICBoZWFkZXJzOiBvYmplY3QgPSB7fSxcbiAgICBheGlvc0NvbmZpZzogQXhpb3NSZXF1ZXN0Q29uZmlnID0gdW5kZWZpbmVkXG4gICk6IFByb21pc2U8UmVxdWVzdFJlc3BvbnNlRGF0YT4gPT5cbiAgICB0aGlzLl9yZXF1ZXN0KFxuICAgICAgXCJQVVRcIixcbiAgICAgIGJhc2V1cmwsXG4gICAgICBnZXRkYXRhLFxuICAgICAgcG9zdGRhdGEsXG4gICAgICB0aGlzLl9zZXRIZWFkZXJzKGhlYWRlcnMpLFxuICAgICAgYXhpb3NDb25maWdcbiAgICApXG5cbiAgLyoqXG4gICAqIE1ha2VzIGEgUEFUQ0ggY2FsbCB0byBhbiBBUEkuXG4gICAqXG4gICAqIEBwYXJhbSBiYXNldXJsIFBhdGggdG8gdGhlIGJhc2V1cmxcbiAgICogQHBhcmFtIGdldGRhdGEgT2JqZWN0IGNvbnRhaW5pbmcgdGhlIGtleSB2YWx1ZSBwYWlycyBzZW50IGluIFBBVENIXG4gICAqIEBwYXJhbSBwb3N0ZGF0YSBPYmplY3QgY29udGFpbmluZyB0aGUga2V5IHZhbHVlIHBhaXJzIHNlbnQgaW4gUEFUQ0hcbiAgICogQHBhcmFtIHBhcmFtZXRlcnMgT2JqZWN0IGNvbnRhaW5pbmcgdGhlIHBhcmFtZXRlcnMgb2YgdGhlIEFQSSBjYWxsXG4gICAqIEBwYXJhbSBoZWFkZXJzIEFuIGFycmF5IEhUVFAgUmVxdWVzdCBIZWFkZXJzXG4gICAqIEBwYXJhbSBheGlvc0NvbmZpZyBDb25maWd1cmF0aW9uIGZvciB0aGUgYXhpb3MgamF2YXNjcmlwdCBsaWJyYXJ5IHRoYXQgd2lsbCBiZSB0aGVcbiAgICogZm91bmRhdGlvbiBmb3IgdGhlIHJlc3Qgb2YgdGhlIHBhcmFtZXRlcnNcbiAgICpcbiAgICogQHJldHVybnMgQSBwcm9taXNlIGZvciBbW1JlcXVlc3RSZXNwb25zZURhdGFdXVxuICAgKi9cbiAgcGF0Y2ggPSAoXG4gICAgYmFzZXVybDogc3RyaW5nLFxuICAgIGdldGRhdGE6IG9iamVjdCxcbiAgICBwb3N0ZGF0YTogc3RyaW5nIHwgb2JqZWN0IHwgQXJyYXlCdWZmZXIgfCBBcnJheUJ1ZmZlclZpZXcsXG4gICAgaGVhZGVyczogb2JqZWN0ID0ge30sXG4gICAgYXhpb3NDb25maWc6IEF4aW9zUmVxdWVzdENvbmZpZyA9IHVuZGVmaW5lZFxuICApOiBQcm9taXNlPFJlcXVlc3RSZXNwb25zZURhdGE+ID0+XG4gICAgdGhpcy5fcmVxdWVzdChcbiAgICAgIFwiUEFUQ0hcIixcbiAgICAgIGJhc2V1cmwsXG4gICAgICBnZXRkYXRhLFxuICAgICAgcG9zdGRhdGEsXG4gICAgICB0aGlzLl9zZXRIZWFkZXJzKGhlYWRlcnMpLFxuICAgICAgYXhpb3NDb25maWdcbiAgICApXG5cbiAgLyoqXG4gICAqIENyZWF0ZXMgYSBuZXcgQXZhbGFuY2hlIGluc3RhbmNlLiBTZXRzIHRoZSBhZGRyZXNzIGFuZCBwb3J0IG9mIHRoZSBtYWluIEF2YWxhbmNoZSBDbGllbnQuXG4gICAqXG4gICAqIEBwYXJhbSBob3N0IFRoZSBob3N0bmFtZSB0byByZXNvbHZlIHRvIHJlYWNoIHRoZSBBdmFsYW5jaGUgQ2xpZW50IEFQSXNcbiAgICogQHBhcmFtIHBvcnQgVGhlIHBvcnQgdG8gcmVzb2x2ZSB0byByZWFjaCB0aGUgQXZhbGFuY2hlIENsaWVudCBBUElzXG4gICAqIEBwYXJhbSBwcm90b2NvbCBUaGUgcHJvdG9jb2wgc3RyaW5nIHRvIHVzZSBiZWZvcmUgYSBcIjovL1wiIGluIGEgcmVxdWVzdCwgZXg6IFwiaHR0cFwiLCBcImh0dHBzXCIsIFwiZ2l0XCIsIFwid3NcIiwgZXRjIC4uLlxuICAgKi9cbiAgY29uc3RydWN0b3IoaG9zdD86IHN0cmluZywgcG9ydD86IG51bWJlciwgcHJvdG9jb2w6IHN0cmluZyA9IFwiaHR0cFwiKSB7XG4gICAgaWYgKGhvc3QgIT0gdW5kZWZpbmVkKSB7XG4gICAgICB0aGlzLnNldEFkZHJlc3MoaG9zdCwgcG9ydCwgcHJvdG9jb2wpXG4gICAgfVxuICB9XG59XG4iXX0=