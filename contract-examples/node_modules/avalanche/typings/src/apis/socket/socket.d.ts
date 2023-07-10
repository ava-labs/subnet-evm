/// <reference types="node" />
/**
 * @packageDocumentation
 * @module API-Socket
 */
import { ClientRequestArgs } from "http";
import WebSocket from "isomorphic-ws";
export declare class Socket extends WebSocket {
    onopen: any;
    onmessage: any;
    onclose: any;
    onerror: any;
    /**
     * Send a message to the server
     *
     * @param data
     * @param cb Optional
     */
    send(data: any, cb?: any): void;
    /**
     * Terminates the connection completely
     *
     * @param mcode Optional
     * @param data Optional
     */
    close(mcode?: number, data?: string): void;
    /**
     * Provides the API for creating and managing a WebSocket connection to a server, as well as for sending and receiving data on the connection.
     *
     * @param url Defaults to [[MainnetAPI]]
     * @param options Optional
     */
    constructor(url?: string | import("url").URL, options?: WebSocket.ClientOptions | ClientRequestArgs);
}
//# sourceMappingURL=socket.d.ts.map