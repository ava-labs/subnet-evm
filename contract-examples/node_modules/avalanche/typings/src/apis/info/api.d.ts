/**
 * @packageDocumentation
 * @module API-Info
 */
import AvalancheCore from "../../avalanche";
import { JRPCAPI } from "../../common/jrpcapi";
import { GetTxFeeResponse, PeersResponse, UptimeResponse } from "./interfaces";
/**
 * Class for interacting with a node's InfoAPI.
 *
 * @category RPCAPIs
 *
 * @remarks This extends the [[JRPCAPI]] class. This class should not be directly called. Instead, use the [[Avalanche.addAPI]] function to register this interface with Avalanche.
 */
export declare class InfoAPI extends JRPCAPI {
    /**
     * Fetches the blockchainID from the node for a given alias.
     *
     * @param alias The blockchain alias to get the blockchainID
     *
     * @returns Returns a Promise string containing the base 58 string representation of the blockchainID.
     */
    getBlockchainID: (alias: string) => Promise<string>;
    /**
     * Fetches the IP address from the node.
     *
     * @returns Returns a Promise string of the node IP address.
     */
    getNodeIP: () => Promise<string>;
    /**
     * Fetches the networkID from the node.
     *
     * @returns Returns a Promise number of the networkID.
     */
    getNetworkID: () => Promise<number>;
    /**
     * Fetches the network name this node is running on
     *
     * @returns Returns a Promise string containing the network name.
     */
    getNetworkName: () => Promise<string>;
    /**
     * Fetches the nodeID from the node.
     *
     * @returns Returns a Promise string of the nodeID.
     */
    getNodeID: () => Promise<string>;
    /**
     * Fetches the version of Gecko this node is running
     *
     * @returns Returns a Promise string containing the version of Gecko.
     */
    getNodeVersion: () => Promise<string>;
    /**
     * Fetches the transaction fee from the node.
     *
     * @returns Returns a Promise object of the transaction fee in nAVAX.
     */
    getTxFee: () => Promise<GetTxFeeResponse>;
    /**
     * Check whether a given chain is done bootstrapping
     * @param chain The ID or alias of a chain.
     *
     * @returns Returns a Promise boolean of whether the chain has completed bootstrapping.
     */
    isBootstrapped: (chain: string) => Promise<boolean>;
    /**
     * Returns the peers connected to the node.
     * @param nodeIDs an optional parameter to specify what nodeID's descriptions should be returned.
     * If this parameter is left empty, descriptions for all active connections will be returned.
     * If the node is not connected to a specified nodeID, it will be omitted from the response.
     *
     * @returns Promise for the list of connected peers in PeersResponse format.
     */
    peers: (nodeIDs?: string[]) => Promise<PeersResponse[]>;
    /**
     * Returns the network's observed uptime of this node.
     *
     * @returns Returns a Promise UptimeResponse which contains rewardingStakePercentage and weightedAveragePercentage.
     */
    uptime: () => Promise<UptimeResponse>;
    /**
     * This class should not be instantiated directly. Instead use the [[Avalanche.addAPI]] method.
     *
     * @param core A reference to the Avalanche class
     * @param baseURL Defaults to the string "/ext/info" as the path to rpc's baseURL
     */
    constructor(core: AvalancheCore, baseURL?: string);
}
//# sourceMappingURL=api.d.ts.map