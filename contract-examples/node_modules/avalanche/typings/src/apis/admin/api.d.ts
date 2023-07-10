/**
 * @packageDocumentation
 * @module API-Admin
 */
import AvalancheCore from "../../avalanche";
import { JRPCAPI } from "../../common/jrpcapi";
import { GetLoggerLevelResponse, LoadVMsResponse, SetLoggerLevelResponse } from "./interfaces";
/**
 * Class for interacting with a node's AdminAPI.
 *
 * @category RPCAPIs
 *
 * @remarks This extends the [[JRPCAPI]] class. This class should not be directly called.
 * Instead, use the [[Avalanche.addAPI]] function to register this interface with Avalanche.
 */
export declare class AdminAPI extends JRPCAPI {
    /**
     * Assign an API an alias, a different endpoint for the API. The original endpoint will still
     * work. This change only affects this node other nodes will not know about this alias.
     *
     * @param endpoint The original endpoint of the API. endpoint should only include the part of
     * the endpoint after /ext/
     * @param alias The API being aliased can now be called at ext/alias
     *
     * @returns Returns a Promise boolean containing success, true for success, false for failure.
     */
    alias: (endpoint: string, alias: string) => Promise<boolean>;
    /**
     * Give a blockchain an alias, a different name that can be used any place the blockchain’s
     * ID is used.
     *
     * @param chain The blockchain’s ID
     * @param alias Can now be used in place of the blockchain’s ID (in API endpoints, for example)
     *
     * @returns Returns a Promise boolean containing success, true for success, false for failure.
     */
    aliasChain: (chain: string, alias: string) => Promise<boolean>;
    /**
     * Get all aliases for given blockchain
     *
     * @param chain The blockchain’s ID
     *
     * @returns Returns a Promise string[] containing aliases of the blockchain.
     */
    getChainAliases: (chain: string) => Promise<string[]>;
    /**
     * Returns log and display levels of loggers
     *
     * @param loggerName the name of the logger to be returned. This is an optional argument. If not specified, it returns all possible loggers.
     *
     * @returns Returns a Promise containing logger levels
     */
    getLoggerLevel: (loggerName?: string) => Promise<GetLoggerLevelResponse>;
    /**
     * Dynamically loads any virtual machines installed on the node as plugins
     *
     * @returns Returns a Promise containing new VMs and failed VMs
     */
    loadVMs: () => Promise<LoadVMsResponse>;
    /**
     * Dump the mutex statistics of the node to the specified file.
     *
     * @returns Promise for a boolean that is true on success.
     */
    lockProfile: () => Promise<boolean>;
    /**
     * Dump the current memory footprint of the node to the specified file.
     *
     * @returns Promise for a boolean that is true on success.
     */
    memoryProfile: () => Promise<boolean>;
    /**
     * Sets log and display levels of loggers.
     *
     * @param loggerName the name of the logger to be changed. This is an optional parameter.
     * @param logLevel the log level of written logs, can be omitted.
     * @param displayLevel the log level of displayed logs, can be omitted.
     *
     * @returns Returns a Promise containing logger levels
     */
    setLoggerLevel: (loggerName?: string, logLevel?: string, displayLevel?: string) => Promise<SetLoggerLevelResponse>;
    /**
     * Start profiling the cpu utilization of the node. Will dump the profile information into
     * the specified file on stop.
     *
     * @returns Promise for a boolean that is true on success.
     */
    startCPUProfiler: () => Promise<boolean>;
    /**
     * Stop the CPU profile that was previously started.
     *
     * @returns Promise for a boolean that is true on success.
     */
    stopCPUProfiler: () => Promise<boolean>;
    /**
     * This class should not be instantiated directly. Instead use the [[Avalanche.addAPI]]
     * method.
     *
     * @param core A reference to the Avalanche class
     * @param baseURL Defaults to the string "/ext/admin" as the path to rpc's baseURL
     */
    constructor(core: AvalancheCore, baseURL?: string);
}
//# sourceMappingURL=api.d.ts.map