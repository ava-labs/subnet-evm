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
exports.InfoAPI = void 0;
const jrpcapi_1 = require("../../common/jrpcapi");
const bn_js_1 = __importDefault(require("bn.js"));
/**
 * Class for interacting with a node's InfoAPI.
 *
 * @category RPCAPIs
 *
 * @remarks This extends the [[JRPCAPI]] class. This class should not be directly called. Instead, use the [[Avalanche.addAPI]] function to register this interface with Avalanche.
 */
class InfoAPI extends jrpcapi_1.JRPCAPI {
    /**
     * This class should not be instantiated directly. Instead use the [[Avalanche.addAPI]] method.
     *
     * @param core A reference to the Avalanche class
     * @param baseURL Defaults to the string "/ext/info" as the path to rpc's baseURL
     */
    constructor(core, baseURL = "/ext/info") {
        super(core, baseURL);
        /**
         * Fetches the blockchainID from the node for a given alias.
         *
         * @param alias The blockchain alias to get the blockchainID
         *
         * @returns Returns a Promise string containing the base 58 string representation of the blockchainID.
         */
        this.getBlockchainID = (alias) => __awaiter(this, void 0, void 0, function* () {
            const params = {
                alias
            };
            const response = yield this.callMethod("info.getBlockchainID", params);
            return response.data.result.blockchainID;
        });
        /**
         * Fetches the IP address from the node.
         *
         * @returns Returns a Promise string of the node IP address.
         */
        this.getNodeIP = () => __awaiter(this, void 0, void 0, function* () {
            const response = yield this.callMethod("info.getBlockchainID");
            return response.data.result.ip;
        });
        /**
         * Fetches the networkID from the node.
         *
         * @returns Returns a Promise number of the networkID.
         */
        this.getNetworkID = () => __awaiter(this, void 0, void 0, function* () {
            const response = yield this.callMethod("info.getNetworkID");
            return response.data.result.networkID;
        });
        /**
         * Fetches the network name this node is running on
         *
         * @returns Returns a Promise string containing the network name.
         */
        this.getNetworkName = () => __awaiter(this, void 0, void 0, function* () {
            const response = yield this.callMethod("info.getNetworkName");
            return response.data.result.networkName;
        });
        /**
         * Fetches the nodeID from the node.
         *
         * @returns Returns a Promise string of the nodeID.
         */
        this.getNodeID = () => __awaiter(this, void 0, void 0, function* () {
            const response = yield this.callMethod("info.getNodeID");
            return response.data.result.nodeID;
        });
        /**
         * Fetches the version of Gecko this node is running
         *
         * @returns Returns a Promise string containing the version of Gecko.
         */
        this.getNodeVersion = () => __awaiter(this, void 0, void 0, function* () {
            const response = yield this.callMethod("info.getNodeVersion");
            return response.data.result.version;
        });
        /**
         * Fetches the transaction fee from the node.
         *
         * @returns Returns a Promise object of the transaction fee in nAVAX.
         */
        this.getTxFee = () => __awaiter(this, void 0, void 0, function* () {
            const response = yield this.callMethod("info.getTxFee");
            return {
                txFee: new bn_js_1.default(response.data.result.txFee, 10),
                creationTxFee: new bn_js_1.default(response.data.result.creationTxFee, 10)
            };
        });
        /**
         * Check whether a given chain is done bootstrapping
         * @param chain The ID or alias of a chain.
         *
         * @returns Returns a Promise boolean of whether the chain has completed bootstrapping.
         */
        this.isBootstrapped = (chain) => __awaiter(this, void 0, void 0, function* () {
            const params = {
                chain
            };
            const response = yield this.callMethod("info.isBootstrapped", params);
            return response.data.result.isBootstrapped;
        });
        /**
         * Returns the peers connected to the node.
         * @param nodeIDs an optional parameter to specify what nodeID's descriptions should be returned.
         * If this parameter is left empty, descriptions for all active connections will be returned.
         * If the node is not connected to a specified nodeID, it will be omitted from the response.
         *
         * @returns Promise for the list of connected peers in PeersResponse format.
         */
        this.peers = (nodeIDs = []) => __awaiter(this, void 0, void 0, function* () {
            const params = {
                nodeIDs
            };
            const response = yield this.callMethod("info.peers", params);
            return response.data.result.peers;
        });
        /**
         * Returns the network's observed uptime of this node.
         *
         * @returns Returns a Promise UptimeResponse which contains rewardingStakePercentage and weightedAveragePercentage.
         */
        this.uptime = () => __awaiter(this, void 0, void 0, function* () {
            const response = yield this.callMethod("info.uptime");
            return response.data.result;
        });
    }
}
exports.InfoAPI = InfoAPI;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYXBpLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2FwaXMvaW5mby9hcGkudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7Ozs7O0FBS0Esa0RBQThDO0FBRTlDLGtEQUFzQjtBQVV0Qjs7Ozs7O0dBTUc7QUFDSCxNQUFhLE9BQVEsU0FBUSxpQkFBTztJQTJJbEM7Ozs7O09BS0c7SUFDSCxZQUFZLElBQW1CLEVBQUUsVUFBa0IsV0FBVztRQUM1RCxLQUFLLENBQUMsSUFBSSxFQUFFLE9BQU8sQ0FBQyxDQUFBO1FBakp0Qjs7Ozs7O1dBTUc7UUFDSCxvQkFBZSxHQUFHLENBQU8sS0FBYSxFQUFtQixFQUFFO1lBQ3pELE1BQU0sTUFBTSxHQUEwQjtnQkFDcEMsS0FBSzthQUNOLENBQUE7WUFFRCxNQUFNLFFBQVEsR0FBd0IsTUFBTSxJQUFJLENBQUMsVUFBVSxDQUN6RCxzQkFBc0IsRUFDdEIsTUFBTSxDQUNQLENBQUE7WUFDRCxPQUFPLFFBQVEsQ0FBQyxJQUFJLENBQUMsTUFBTSxDQUFDLFlBQVksQ0FBQTtRQUMxQyxDQUFDLENBQUEsQ0FBQTtRQUVEOzs7O1dBSUc7UUFDSCxjQUFTLEdBQUcsR0FBMEIsRUFBRTtZQUN0QyxNQUFNLFFBQVEsR0FBd0IsTUFBTSxJQUFJLENBQUMsVUFBVSxDQUN6RCxzQkFBc0IsQ0FDdkIsQ0FBQTtZQUNELE9BQU8sUUFBUSxDQUFDLElBQUksQ0FBQyxNQUFNLENBQUMsRUFBRSxDQUFBO1FBQ2hDLENBQUMsQ0FBQSxDQUFBO1FBRUQ7Ozs7V0FJRztRQUNILGlCQUFZLEdBQUcsR0FBMEIsRUFBRTtZQUN6QyxNQUFNLFFBQVEsR0FBd0IsTUFBTSxJQUFJLENBQUMsVUFBVSxDQUN6RCxtQkFBbUIsQ0FDcEIsQ0FBQTtZQUNELE9BQU8sUUFBUSxDQUFDLElBQUksQ0FBQyxNQUFNLENBQUMsU0FBUyxDQUFBO1FBQ3ZDLENBQUMsQ0FBQSxDQUFBO1FBRUQ7Ozs7V0FJRztRQUNILG1CQUFjLEdBQUcsR0FBMEIsRUFBRTtZQUMzQyxNQUFNLFFBQVEsR0FBd0IsTUFBTSxJQUFJLENBQUMsVUFBVSxDQUN6RCxxQkFBcUIsQ0FDdEIsQ0FBQTtZQUNELE9BQU8sUUFBUSxDQUFDLElBQUksQ0FBQyxNQUFNLENBQUMsV0FBVyxDQUFBO1FBQ3pDLENBQUMsQ0FBQSxDQUFBO1FBRUQ7Ozs7V0FJRztRQUNILGNBQVMsR0FBRyxHQUEwQixFQUFFO1lBQ3RDLE1BQU0sUUFBUSxHQUF3QixNQUFNLElBQUksQ0FBQyxVQUFVLENBQ3pELGdCQUFnQixDQUNqQixDQUFBO1lBQ0QsT0FBTyxRQUFRLENBQUMsSUFBSSxDQUFDLE1BQU0sQ0FBQyxNQUFNLENBQUE7UUFDcEMsQ0FBQyxDQUFBLENBQUE7UUFFRDs7OztXQUlHO1FBQ0gsbUJBQWMsR0FBRyxHQUEwQixFQUFFO1lBQzNDLE1BQU0sUUFBUSxHQUF3QixNQUFNLElBQUksQ0FBQyxVQUFVLENBQ3pELHFCQUFxQixDQUN0QixDQUFBO1lBQ0QsT0FBTyxRQUFRLENBQUMsSUFBSSxDQUFDLE1BQU0sQ0FBQyxPQUFPLENBQUE7UUFDckMsQ0FBQyxDQUFBLENBQUE7UUFFRDs7OztXQUlHO1FBQ0gsYUFBUSxHQUFHLEdBQW9DLEVBQUU7WUFDL0MsTUFBTSxRQUFRLEdBQXdCLE1BQU0sSUFBSSxDQUFDLFVBQVUsQ0FBQyxlQUFlLENBQUMsQ0FBQTtZQUM1RSxPQUFPO2dCQUNMLEtBQUssRUFBRSxJQUFJLGVBQUUsQ0FBQyxRQUFRLENBQUMsSUFBSSxDQUFDLE1BQU0sQ0FBQyxLQUFLLEVBQUUsRUFBRSxDQUFDO2dCQUM3QyxhQUFhLEVBQUUsSUFBSSxlQUFFLENBQUMsUUFBUSxDQUFDLElBQUksQ0FBQyxNQUFNLENBQUMsYUFBYSxFQUFFLEVBQUUsQ0FBQzthQUM5RCxDQUFBO1FBQ0gsQ0FBQyxDQUFBLENBQUE7UUFFRDs7Ozs7V0FLRztRQUNILG1CQUFjLEdBQUcsQ0FBTyxLQUFhLEVBQW9CLEVBQUU7WUFDekQsTUFBTSxNQUFNLEdBQXlCO2dCQUNuQyxLQUFLO2FBQ04sQ0FBQTtZQUNELE1BQU0sUUFBUSxHQUF3QixNQUFNLElBQUksQ0FBQyxVQUFVLENBQ3pELHFCQUFxQixFQUNyQixNQUFNLENBQ1AsQ0FBQTtZQUNELE9BQU8sUUFBUSxDQUFDLElBQUksQ0FBQyxNQUFNLENBQUMsY0FBYyxDQUFBO1FBQzVDLENBQUMsQ0FBQSxDQUFBO1FBRUQ7Ozs7Ozs7V0FPRztRQUNILFVBQUssR0FBRyxDQUFPLFVBQW9CLEVBQUUsRUFBNEIsRUFBRTtZQUNqRSxNQUFNLE1BQU0sR0FBZ0I7Z0JBQzFCLE9BQU87YUFDUixDQUFBO1lBQ0QsTUFBTSxRQUFRLEdBQXdCLE1BQU0sSUFBSSxDQUFDLFVBQVUsQ0FDekQsWUFBWSxFQUNaLE1BQU0sQ0FDUCxDQUFBO1lBQ0QsT0FBTyxRQUFRLENBQUMsSUFBSSxDQUFDLE1BQU0sQ0FBQyxLQUFLLENBQUE7UUFDbkMsQ0FBQyxDQUFBLENBQUE7UUFFRDs7OztXQUlHO1FBQ0gsV0FBTSxHQUFHLEdBQWtDLEVBQUU7WUFDM0MsTUFBTSxRQUFRLEdBQXdCLE1BQU0sSUFBSSxDQUFDLFVBQVUsQ0FBQyxhQUFhLENBQUMsQ0FBQTtZQUMxRSxPQUFPLFFBQVEsQ0FBQyxJQUFJLENBQUMsTUFBTSxDQUFBO1FBQzdCLENBQUMsQ0FBQSxDQUFBO0lBVUQsQ0FBQztDQUNGO0FBcEpELDBCQW9KQyIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICogQHBhY2thZ2VEb2N1bWVudGF0aW9uXG4gKiBAbW9kdWxlIEFQSS1JbmZvXG4gKi9cbmltcG9ydCBBdmFsYW5jaGVDb3JlIGZyb20gXCIuLi8uLi9hdmFsYW5jaGVcIlxuaW1wb3J0IHsgSlJQQ0FQSSB9IGZyb20gXCIuLi8uLi9jb21tb24vanJwY2FwaVwiXG5pbXBvcnQgeyBSZXF1ZXN0UmVzcG9uc2VEYXRhIH0gZnJvbSBcIi4uLy4uL2NvbW1vbi9hcGliYXNlXCJcbmltcG9ydCBCTiBmcm9tIFwiYm4uanNcIlxuaW1wb3J0IHtcbiAgR2V0QmxvY2tjaGFpbklEUGFyYW1zLFxuICBHZXRUeEZlZVJlc3BvbnNlLFxuICBJc0Jvb3RzdHJhcHBlZFBhcmFtcyxcbiAgUGVlcnNQYXJhbXMsXG4gIFBlZXJzUmVzcG9uc2UsXG4gIFVwdGltZVJlc3BvbnNlXG59IGZyb20gXCIuL2ludGVyZmFjZXNcIlxuXG4vKipcbiAqIENsYXNzIGZvciBpbnRlcmFjdGluZyB3aXRoIGEgbm9kZSdzIEluZm9BUEkuXG4gKlxuICogQGNhdGVnb3J5IFJQQ0FQSXNcbiAqXG4gKiBAcmVtYXJrcyBUaGlzIGV4dGVuZHMgdGhlIFtbSlJQQ0FQSV1dIGNsYXNzLiBUaGlzIGNsYXNzIHNob3VsZCBub3QgYmUgZGlyZWN0bHkgY2FsbGVkLiBJbnN0ZWFkLCB1c2UgdGhlIFtbQXZhbGFuY2hlLmFkZEFQSV1dIGZ1bmN0aW9uIHRvIHJlZ2lzdGVyIHRoaXMgaW50ZXJmYWNlIHdpdGggQXZhbGFuY2hlLlxuICovXG5leHBvcnQgY2xhc3MgSW5mb0FQSSBleHRlbmRzIEpSUENBUEkge1xuICAvKipcbiAgICogRmV0Y2hlcyB0aGUgYmxvY2tjaGFpbklEIGZyb20gdGhlIG5vZGUgZm9yIGEgZ2l2ZW4gYWxpYXMuXG4gICAqXG4gICAqIEBwYXJhbSBhbGlhcyBUaGUgYmxvY2tjaGFpbiBhbGlhcyB0byBnZXQgdGhlIGJsb2NrY2hhaW5JRFxuICAgKlxuICAgKiBAcmV0dXJucyBSZXR1cm5zIGEgUHJvbWlzZSBzdHJpbmcgY29udGFpbmluZyB0aGUgYmFzZSA1OCBzdHJpbmcgcmVwcmVzZW50YXRpb24gb2YgdGhlIGJsb2NrY2hhaW5JRC5cbiAgICovXG4gIGdldEJsb2NrY2hhaW5JRCA9IGFzeW5jIChhbGlhczogc3RyaW5nKTogUHJvbWlzZTxzdHJpbmc+ID0+IHtcbiAgICBjb25zdCBwYXJhbXM6IEdldEJsb2NrY2hhaW5JRFBhcmFtcyA9IHtcbiAgICAgIGFsaWFzXG4gICAgfVxuXG4gICAgY29uc3QgcmVzcG9uc2U6IFJlcXVlc3RSZXNwb25zZURhdGEgPSBhd2FpdCB0aGlzLmNhbGxNZXRob2QoXG4gICAgICBcImluZm8uZ2V0QmxvY2tjaGFpbklEXCIsXG4gICAgICBwYXJhbXNcbiAgICApXG4gICAgcmV0dXJuIHJlc3BvbnNlLmRhdGEucmVzdWx0LmJsb2NrY2hhaW5JRFxuICB9XG5cbiAgLyoqXG4gICAqIEZldGNoZXMgdGhlIElQIGFkZHJlc3MgZnJvbSB0aGUgbm9kZS5cbiAgICpcbiAgICogQHJldHVybnMgUmV0dXJucyBhIFByb21pc2Ugc3RyaW5nIG9mIHRoZSBub2RlIElQIGFkZHJlc3MuXG4gICAqL1xuICBnZXROb2RlSVAgPSBhc3luYyAoKTogUHJvbWlzZTxzdHJpbmc+ID0+IHtcbiAgICBjb25zdCByZXNwb25zZTogUmVxdWVzdFJlc3BvbnNlRGF0YSA9IGF3YWl0IHRoaXMuY2FsbE1ldGhvZChcbiAgICAgIFwiaW5mby5nZXRCbG9ja2NoYWluSURcIlxuICAgIClcbiAgICByZXR1cm4gcmVzcG9uc2UuZGF0YS5yZXN1bHQuaXBcbiAgfVxuXG4gIC8qKlxuICAgKiBGZXRjaGVzIHRoZSBuZXR3b3JrSUQgZnJvbSB0aGUgbm9kZS5cbiAgICpcbiAgICogQHJldHVybnMgUmV0dXJucyBhIFByb21pc2UgbnVtYmVyIG9mIHRoZSBuZXR3b3JrSUQuXG4gICAqL1xuICBnZXROZXR3b3JrSUQgPSBhc3luYyAoKTogUHJvbWlzZTxudW1iZXI+ID0+IHtcbiAgICBjb25zdCByZXNwb25zZTogUmVxdWVzdFJlc3BvbnNlRGF0YSA9IGF3YWl0IHRoaXMuY2FsbE1ldGhvZChcbiAgICAgIFwiaW5mby5nZXROZXR3b3JrSURcIlxuICAgIClcbiAgICByZXR1cm4gcmVzcG9uc2UuZGF0YS5yZXN1bHQubmV0d29ya0lEXG4gIH1cblxuICAvKipcbiAgICogRmV0Y2hlcyB0aGUgbmV0d29yayBuYW1lIHRoaXMgbm9kZSBpcyBydW5uaW5nIG9uXG4gICAqXG4gICAqIEByZXR1cm5zIFJldHVybnMgYSBQcm9taXNlIHN0cmluZyBjb250YWluaW5nIHRoZSBuZXR3b3JrIG5hbWUuXG4gICAqL1xuICBnZXROZXR3b3JrTmFtZSA9IGFzeW5jICgpOiBQcm9taXNlPHN0cmluZz4gPT4ge1xuICAgIGNvbnN0IHJlc3BvbnNlOiBSZXF1ZXN0UmVzcG9uc2VEYXRhID0gYXdhaXQgdGhpcy5jYWxsTWV0aG9kKFxuICAgICAgXCJpbmZvLmdldE5ldHdvcmtOYW1lXCJcbiAgICApXG4gICAgcmV0dXJuIHJlc3BvbnNlLmRhdGEucmVzdWx0Lm5ldHdvcmtOYW1lXG4gIH1cblxuICAvKipcbiAgICogRmV0Y2hlcyB0aGUgbm9kZUlEIGZyb20gdGhlIG5vZGUuXG4gICAqXG4gICAqIEByZXR1cm5zIFJldHVybnMgYSBQcm9taXNlIHN0cmluZyBvZiB0aGUgbm9kZUlELlxuICAgKi9cbiAgZ2V0Tm9kZUlEID0gYXN5bmMgKCk6IFByb21pc2U8c3RyaW5nPiA9PiB7XG4gICAgY29uc3QgcmVzcG9uc2U6IFJlcXVlc3RSZXNwb25zZURhdGEgPSBhd2FpdCB0aGlzLmNhbGxNZXRob2QoXG4gICAgICBcImluZm8uZ2V0Tm9kZUlEXCJcbiAgICApXG4gICAgcmV0dXJuIHJlc3BvbnNlLmRhdGEucmVzdWx0Lm5vZGVJRFxuICB9XG5cbiAgLyoqXG4gICAqIEZldGNoZXMgdGhlIHZlcnNpb24gb2YgR2Vja28gdGhpcyBub2RlIGlzIHJ1bm5pbmdcbiAgICpcbiAgICogQHJldHVybnMgUmV0dXJucyBhIFByb21pc2Ugc3RyaW5nIGNvbnRhaW5pbmcgdGhlIHZlcnNpb24gb2YgR2Vja28uXG4gICAqL1xuICBnZXROb2RlVmVyc2lvbiA9IGFzeW5jICgpOiBQcm9taXNlPHN0cmluZz4gPT4ge1xuICAgIGNvbnN0IHJlc3BvbnNlOiBSZXF1ZXN0UmVzcG9uc2VEYXRhID0gYXdhaXQgdGhpcy5jYWxsTWV0aG9kKFxuICAgICAgXCJpbmZvLmdldE5vZGVWZXJzaW9uXCJcbiAgICApXG4gICAgcmV0dXJuIHJlc3BvbnNlLmRhdGEucmVzdWx0LnZlcnNpb25cbiAgfVxuXG4gIC8qKlxuICAgKiBGZXRjaGVzIHRoZSB0cmFuc2FjdGlvbiBmZWUgZnJvbSB0aGUgbm9kZS5cbiAgICpcbiAgICogQHJldHVybnMgUmV0dXJucyBhIFByb21pc2Ugb2JqZWN0IG9mIHRoZSB0cmFuc2FjdGlvbiBmZWUgaW4gbkFWQVguXG4gICAqL1xuICBnZXRUeEZlZSA9IGFzeW5jICgpOiBQcm9taXNlPEdldFR4RmVlUmVzcG9uc2U+ID0+IHtcbiAgICBjb25zdCByZXNwb25zZTogUmVxdWVzdFJlc3BvbnNlRGF0YSA9IGF3YWl0IHRoaXMuY2FsbE1ldGhvZChcImluZm8uZ2V0VHhGZWVcIilcbiAgICByZXR1cm4ge1xuICAgICAgdHhGZWU6IG5ldyBCTihyZXNwb25zZS5kYXRhLnJlc3VsdC50eEZlZSwgMTApLFxuICAgICAgY3JlYXRpb25UeEZlZTogbmV3IEJOKHJlc3BvbnNlLmRhdGEucmVzdWx0LmNyZWF0aW9uVHhGZWUsIDEwKVxuICAgIH1cbiAgfVxuXG4gIC8qKlxuICAgKiBDaGVjayB3aGV0aGVyIGEgZ2l2ZW4gY2hhaW4gaXMgZG9uZSBib290c3RyYXBwaW5nXG4gICAqIEBwYXJhbSBjaGFpbiBUaGUgSUQgb3IgYWxpYXMgb2YgYSBjaGFpbi5cbiAgICpcbiAgICogQHJldHVybnMgUmV0dXJucyBhIFByb21pc2UgYm9vbGVhbiBvZiB3aGV0aGVyIHRoZSBjaGFpbiBoYXMgY29tcGxldGVkIGJvb3RzdHJhcHBpbmcuXG4gICAqL1xuICBpc0Jvb3RzdHJhcHBlZCA9IGFzeW5jIChjaGFpbjogc3RyaW5nKTogUHJvbWlzZTxib29sZWFuPiA9PiB7XG4gICAgY29uc3QgcGFyYW1zOiBJc0Jvb3RzdHJhcHBlZFBhcmFtcyA9IHtcbiAgICAgIGNoYWluXG4gICAgfVxuICAgIGNvbnN0IHJlc3BvbnNlOiBSZXF1ZXN0UmVzcG9uc2VEYXRhID0gYXdhaXQgdGhpcy5jYWxsTWV0aG9kKFxuICAgICAgXCJpbmZvLmlzQm9vdHN0cmFwcGVkXCIsXG4gICAgICBwYXJhbXNcbiAgICApXG4gICAgcmV0dXJuIHJlc3BvbnNlLmRhdGEucmVzdWx0LmlzQm9vdHN0cmFwcGVkXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgcGVlcnMgY29ubmVjdGVkIHRvIHRoZSBub2RlLlxuICAgKiBAcGFyYW0gbm9kZUlEcyBhbiBvcHRpb25hbCBwYXJhbWV0ZXIgdG8gc3BlY2lmeSB3aGF0IG5vZGVJRCdzIGRlc2NyaXB0aW9ucyBzaG91bGQgYmUgcmV0dXJuZWQuXG4gICAqIElmIHRoaXMgcGFyYW1ldGVyIGlzIGxlZnQgZW1wdHksIGRlc2NyaXB0aW9ucyBmb3IgYWxsIGFjdGl2ZSBjb25uZWN0aW9ucyB3aWxsIGJlIHJldHVybmVkLlxuICAgKiBJZiB0aGUgbm9kZSBpcyBub3QgY29ubmVjdGVkIHRvIGEgc3BlY2lmaWVkIG5vZGVJRCwgaXQgd2lsbCBiZSBvbWl0dGVkIGZyb20gdGhlIHJlc3BvbnNlLlxuICAgKlxuICAgKiBAcmV0dXJucyBQcm9taXNlIGZvciB0aGUgbGlzdCBvZiBjb25uZWN0ZWQgcGVlcnMgaW4gUGVlcnNSZXNwb25zZSBmb3JtYXQuXG4gICAqL1xuICBwZWVycyA9IGFzeW5jIChub2RlSURzOiBzdHJpbmdbXSA9IFtdKTogUHJvbWlzZTxQZWVyc1Jlc3BvbnNlW10+ID0+IHtcbiAgICBjb25zdCBwYXJhbXM6IFBlZXJzUGFyYW1zID0ge1xuICAgICAgbm9kZUlEc1xuICAgIH1cbiAgICBjb25zdCByZXNwb25zZTogUmVxdWVzdFJlc3BvbnNlRGF0YSA9IGF3YWl0IHRoaXMuY2FsbE1ldGhvZChcbiAgICAgIFwiaW5mby5wZWVyc1wiLFxuICAgICAgcGFyYW1zXG4gICAgKVxuICAgIHJldHVybiByZXNwb25zZS5kYXRhLnJlc3VsdC5wZWVyc1xuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIG5ldHdvcmsncyBvYnNlcnZlZCB1cHRpbWUgb2YgdGhpcyBub2RlLlxuICAgKlxuICAgKiBAcmV0dXJucyBSZXR1cm5zIGEgUHJvbWlzZSBVcHRpbWVSZXNwb25zZSB3aGljaCBjb250YWlucyByZXdhcmRpbmdTdGFrZVBlcmNlbnRhZ2UgYW5kIHdlaWdodGVkQXZlcmFnZVBlcmNlbnRhZ2UuXG4gICAqL1xuICB1cHRpbWUgPSBhc3luYyAoKTogUHJvbWlzZTxVcHRpbWVSZXNwb25zZT4gPT4ge1xuICAgIGNvbnN0IHJlc3BvbnNlOiBSZXF1ZXN0UmVzcG9uc2VEYXRhID0gYXdhaXQgdGhpcy5jYWxsTWV0aG9kKFwiaW5mby51cHRpbWVcIilcbiAgICByZXR1cm4gcmVzcG9uc2UuZGF0YS5yZXN1bHRcbiAgfVxuXG4gIC8qKlxuICAgKiBUaGlzIGNsYXNzIHNob3VsZCBub3QgYmUgaW5zdGFudGlhdGVkIGRpcmVjdGx5LiBJbnN0ZWFkIHVzZSB0aGUgW1tBdmFsYW5jaGUuYWRkQVBJXV0gbWV0aG9kLlxuICAgKlxuICAgKiBAcGFyYW0gY29yZSBBIHJlZmVyZW5jZSB0byB0aGUgQXZhbGFuY2hlIGNsYXNzXG4gICAqIEBwYXJhbSBiYXNlVVJMIERlZmF1bHRzIHRvIHRoZSBzdHJpbmcgXCIvZXh0L2luZm9cIiBhcyB0aGUgcGF0aCB0byBycGMncyBiYXNlVVJMXG4gICAqL1xuICBjb25zdHJ1Y3Rvcihjb3JlOiBBdmFsYW5jaGVDb3JlLCBiYXNlVVJMOiBzdHJpbmcgPSBcIi9leHQvaW5mb1wiKSB7XG4gICAgc3VwZXIoY29yZSwgYmFzZVVSTClcbiAgfVxufVxuIl19