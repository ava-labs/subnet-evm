"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.utils = exports.platformvm = exports.metrics = exports.keystore = exports.info = exports.index = exports.health = exports.evm = exports.common = exports.avm = exports.auth = exports.admin = exports.Socket = exports.PubSub = exports.Mnemonic = exports.GenesisData = exports.GenesisAsset = exports.HDNode = exports.DB = exports.Buffer = exports.BN = exports.BinTools = exports.AvalancheCore = exports.Avalanche = void 0;
/**
 * @packageDocumentation
 * @module Avalanche
 */
const avalanche_1 = __importDefault(require("./avalanche"));
exports.AvalancheCore = avalanche_1.default;
const api_1 = require("./apis/admin/api");
const api_2 = require("./apis/auth/api");
const api_3 = require("./apis/avm/api");
const api_4 = require("./apis/evm/api");
const genesisasset_1 = require("./apis/avm/genesisasset");
Object.defineProperty(exports, "GenesisAsset", { enumerable: true, get: function () { return genesisasset_1.GenesisAsset; } });
const genesisdata_1 = require("./apis/avm/genesisdata");
Object.defineProperty(exports, "GenesisData", { enumerable: true, get: function () { return genesisdata_1.GenesisData; } });
const api_5 = require("./apis/health/api");
const api_6 = require("./apis/index/api");
const api_7 = require("./apis/info/api");
const api_8 = require("./apis/keystore/api");
const api_9 = require("./apis/metrics/api");
const api_10 = require("./apis/platformvm/api");
const socket_1 = require("./apis/socket/socket");
Object.defineProperty(exports, "Socket", { enumerable: true, get: function () { return socket_1.Socket; } });
const constants_1 = require("./utils/constants");
const helperfunctions_1 = require("./utils/helperfunctions");
const bintools_1 = __importDefault(require("./utils/bintools"));
exports.BinTools = bintools_1.default;
const db_1 = __importDefault(require("./utils/db"));
exports.DB = db_1.default;
const mnemonic_1 = __importDefault(require("./utils/mnemonic"));
exports.Mnemonic = mnemonic_1.default;
const pubsub_1 = __importDefault(require("./utils/pubsub"));
exports.PubSub = pubsub_1.default;
const hdnode_1 = __importDefault(require("./utils/hdnode"));
exports.HDNode = hdnode_1.default;
const bn_js_1 = __importDefault(require("bn.js"));
exports.BN = bn_js_1.default;
const buffer_1 = require("buffer/");
Object.defineProperty(exports, "Buffer", { enumerable: true, get: function () { return buffer_1.Buffer; } });
/**
 * AvalancheJS is middleware for interacting with Avalanche node RPC APIs.
 *
 * Example usage:
 * ```js
 * const avalanche: Avalanche = new Avalanche("127.0.0.1", 9650, "https")
 * ```
 *
 */
class Avalanche extends avalanche_1.default {
    /**
     * Creates a new Avalanche instance. Sets the address and port of the main Avalanche Client.
     *
     * @param host The hostname to resolve to reach the Avalanche Client RPC APIs
     * @param port The port to resolve to reach the Avalanche Client RPC APIs
     * @param protocol The protocol string to use before a "://" in a request,
     * ex: "http", "https", "git", "ws", etc. Defaults to http
     * @param networkID Sets the NetworkID of the class. Default [[DefaultNetworkID]]
     * @param XChainID Sets the blockchainID for the AVM. Will try to auto-detect,
     * otherwise default "2eNy1mUFdmaxXNj1eQHUe7Np4gju9sJsEtWQ4MX3ToiNKuADed"
     * @param CChainID Sets the blockchainID for the EVM. Will try to auto-detect,
     * otherwise default "2CA6j5zYzasynPsFeNoqWkmTCt3VScMvXUZHbfDJ8k3oGzAPtU"
     * @param hrp The human-readable part of the bech32 addresses
     * @param skipinit Skips creating the APIs. Defaults to false
     */
    constructor(host, port, protocol = "http", networkID = constants_1.DefaultNetworkID, XChainID = undefined, CChainID = undefined, hrp = undefined, skipinit = false) {
        super(host, port, protocol);
        /**
         * Returns a reference to the Admin RPC.
         */
        this.Admin = () => this.apis.admin;
        /**
         * Returns a reference to the Auth RPC.
         */
        this.Auth = () => this.apis.auth;
        /**
         * Returns a reference to the EVMAPI RPC pointed at the C-Chain.
         */
        this.CChain = () => this.apis.cchain;
        /**
         * Returns a reference to the AVM RPC pointed at the X-Chain.
         */
        this.XChain = () => this.apis.xchain;
        /**
         * Returns a reference to the Health RPC for a node.
         */
        this.Health = () => this.apis.health;
        /**
         * Returns a reference to the Index RPC for a node.
         */
        this.Index = () => this.apis.index;
        /**
         * Returns a reference to the Info RPC for a node.
         */
        this.Info = () => this.apis.info;
        /**
         * Returns a reference to the Metrics RPC.
         */
        this.Metrics = () => this.apis.metrics;
        /**
         * Returns a reference to the Keystore RPC for a node. We label it "NodeKeys" to reduce
         * confusion about what it's accessing.
         */
        this.NodeKeys = () => this.apis.keystore;
        /**
         * Returns a reference to the PlatformVM RPC pointed at the P-Chain.
         */
        this.PChain = () => this.apis.pchain;
        let xchainid = XChainID;
        let cchainid = CChainID;
        if (typeof XChainID === "undefined" ||
            !XChainID ||
            XChainID.toLowerCase() === "x") {
            if (networkID.toString() in constants_1.Defaults.network) {
                xchainid = constants_1.Defaults.network[`${networkID}`].X.blockchainID;
            }
            else {
                xchainid = constants_1.Defaults.network[12345].X.blockchainID;
            }
        }
        if (typeof CChainID === "undefined" ||
            !CChainID ||
            CChainID.toLowerCase() === "c") {
            if (networkID.toString() in constants_1.Defaults.network) {
                cchainid = constants_1.Defaults.network[`${networkID}`].C.blockchainID;
            }
            else {
                cchainid = constants_1.Defaults.network[12345].C.blockchainID;
            }
        }
        if (typeof networkID === "number" && networkID >= 0) {
            this.networkID = networkID;
        }
        else if (typeof networkID === "undefined") {
            networkID = constants_1.DefaultNetworkID;
        }
        if (typeof hrp !== "undefined") {
            this.hrp = hrp;
        }
        else {
            this.hrp = (0, helperfunctions_1.getPreferredHRP)(this.networkID);
        }
        if (!skipinit) {
            this.addAPI("admin", api_1.AdminAPI);
            this.addAPI("auth", api_2.AuthAPI);
            this.addAPI("xchain", api_3.AVMAPI, "/ext/bc/X", xchainid);
            this.addAPI("cchain", api_4.EVMAPI, "/ext/bc/C/avax", cchainid);
            this.addAPI("health", api_5.HealthAPI);
            this.addAPI("info", api_7.InfoAPI);
            this.addAPI("index", api_6.IndexAPI);
            this.addAPI("keystore", api_8.KeystoreAPI);
            this.addAPI("metrics", api_9.MetricsAPI);
            this.addAPI("pchain", api_10.PlatformVMAPI);
        }
    }
}
exports.default = Avalanche;
exports.Avalanche = Avalanche;
exports.admin = __importStar(require("./apis/admin"));
exports.auth = __importStar(require("./apis/auth"));
exports.avm = __importStar(require("./apis/avm"));
exports.common = __importStar(require("./common"));
exports.evm = __importStar(require("./apis/evm"));
exports.health = __importStar(require("./apis/health"));
exports.index = __importStar(require("./apis/index"));
exports.info = __importStar(require("./apis/info"));
exports.keystore = __importStar(require("./apis/keystore"));
exports.metrics = __importStar(require("./apis/metrics"));
exports.platformvm = __importStar(require("./apis/platformvm"));
exports.utils = __importStar(require("./utils"));
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiaW5kZXguanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi9zcmMvaW5kZXgudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7QUFBQTs7O0dBR0c7QUFDSCw0REFBdUM7QUFtSzlCLHdCQW5LRixtQkFBYSxDQW1LRTtBQWxLdEIsMENBQTJDO0FBQzNDLHlDQUF5QztBQUN6Qyx3Q0FBdUM7QUFDdkMsd0NBQXVDO0FBQ3ZDLDBEQUFzRDtBQW9LN0MsNkZBcEtBLDJCQUFZLE9Bb0tBO0FBbktyQix3REFBb0Q7QUFvSzNDLDRGQXBLQSx5QkFBVyxPQW9LQTtBQW5LcEIsMkNBQTZDO0FBQzdDLDBDQUEyQztBQUMzQyx5Q0FBeUM7QUFDekMsNkNBQWlEO0FBQ2pELDRDQUErQztBQUMvQyxnREFBcUQ7QUFDckQsaURBQTZDO0FBZ0twQyx1RkFoS0EsZUFBTSxPQWdLQTtBQS9KZixpREFBOEQ7QUFDOUQsNkRBQXlEO0FBQ3pELGdFQUF1QztBQW9KOUIsbUJBcEpGLGtCQUFRLENBb0pFO0FBbkpqQixvREFBMkI7QUFzSmxCLGFBdEpGLFlBQUUsQ0FzSkU7QUFySlgsZ0VBQXVDO0FBeUo5QixtQkF6SkYsa0JBQVEsQ0F5SkU7QUF4SmpCLDREQUFtQztBQXlKMUIsaUJBekpGLGdCQUFNLENBeUpFO0FBeEpmLDREQUFtQztBQW9KMUIsaUJBcEpGLGdCQUFNLENBb0pFO0FBbkpmLGtEQUFzQjtBQWdKYixhQWhKRixlQUFFLENBZ0pFO0FBL0lYLG9DQUFnQztBQWdKdkIsdUZBaEpBLGVBQU0sT0FnSkE7QUE5SWY7Ozs7Ozs7O0dBUUc7QUFDSCxNQUFxQixTQUFVLFNBQVEsbUJBQWE7SUFvRGxEOzs7Ozs7Ozs7Ozs7OztPQWNHO0lBQ0gsWUFDRSxJQUFhLEVBQ2IsSUFBYSxFQUNiLFdBQW1CLE1BQU0sRUFDekIsWUFBb0IsNEJBQWdCLEVBQ3BDLFdBQW1CLFNBQVMsRUFDNUIsV0FBbUIsU0FBUyxFQUM1QixNQUFj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QUksQ0FBQyxTQUFTLENBQUMsQ0FBQTtTQUMzQztRQUVELElBQUksQ0FBQyxRQUFRLEVBQUU7WUFDYixJQUFJLENBQUMsTUFBTSxDQUFDLE9BQU8sRUFBRSxjQUFRLENBQUMsQ0FBQTtZQUM5QixJQUFJLENBQUMsTUFBTSxDQUFDLE1BQU0sRUFBRSxhQUFPLENBQUMsQ0FBQTtZQUM1QixJQUFJLENBQUMsTUFBTSxDQUFDLFFBQVEsRUFBRSxZQUFNLEVBQUUsV0FBVyxFQUFFLFFBQVEsQ0FBQyxDQUFBO1lBQ3BELElBQUksQ0FBQyxNQUFNLENBQUMsUUFBUSxFQUFFLFlBQU0sRUFBRSxnQkFBZ0IsRUFBRSxRQUFRLENBQUMsQ0FBQTtZQUN6RCxJQUFJLENBQUMsTUFBTSxDQUFDLFFBQVEsRUFBRSxlQUFTLENBQUMsQ0FBQTtZQUNoQyxJQUFJLENBQUMsTUFBTSxDQUFDLE1BQU0sRUFBRSxhQUFPLENBQUMsQ0FBQTtZQUM1QixJQUFJLENBQUMsTUFBTSxDQUFDLE9BQU8sRUFBRSxjQUFRLENBQUMsQ0FBQTtZQUM5QixJQUFJLENBQUMsTUFBTSxDQUFDLFVBQVUsRUFBRSxpQkFBVyxDQUFDLENBQUE7WUFDcEMsSUFBSSxDQUFDLE1BQU0sQ0FBQyxTQUFTLEVBQUUsZ0JBQVUsQ0FBQyxDQUFBO1lBQ2xDLElBQUksQ0FBQyxNQUFNLENBQUMsUUFBUSxFQUFFLG9CQUFhLENBQUMsQ0FBQTtTQUNyQztJQUNILENBQUM7Q0FDRjtBQS9IRCw0QkErSEM7QUFFUSw4QkFBUztBQWFsQixzREFBcUM7QUFDckMsb0RBQW1DO0FBQ25DLGtEQUFpQztBQUNqQyxtREFBa0M7QUFDbEMsa0RBQWlDO0FBQ2pDLHdEQUF1QztBQUN2QyxzREFBcUM7QUFDckMsb0RBQW1DO0FBQ25DLDREQUEyQztBQUMzQywwREFBeUM7QUFDekMsZ0VBQStDO0FBQy9DLGlEQUFnQyIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICogQHBhY2thZ2VEb2N1bWVudGF0aW9uXG4gKiBAbW9kdWxlIEF2YWxhbmNoZVxuICovXG5pbXBvcnQgQXZhbGFuY2hlQ29yZSBmcm9tIFwiLi9hdmFsYW5jaGVcIlxuaW1wb3J0IHsgQWRtaW5BUEkgfSBmcm9tIFwiLi9hcGlzL2FkbWluL2FwaVwiXG5pbXBvcnQgeyBBdXRoQVBJIH0gZnJvbSBcIi4vYXBpcy9hdXRoL2FwaVwiXG5pbXBvcnQgeyBBVk1BUEkgfSBmcm9tIFwiLi9hcGlzL2F2bS9hcGlcIlxuaW1wb3J0IHsgRVZNQVBJIH0gZnJvbSBcIi4vYXBpcy9ldm0vYXBpXCJcbmltcG9ydCB7IEdlbmVzaXNBc3NldCB9IGZyb20gXCIuL2FwaXMvYXZtL2dlbmVzaXNhc3NldFwiXG5pbXBvcnQgeyBHZW5lc2lzRGF0YSB9IGZyb20gXCIuL2FwaXMvYXZtL2dlbmVzaXNkYXRhXCJcbmltcG9ydCB7IEhlYWx0aEFQSSB9IGZyb20gXCIuL2FwaXMvaGVhbHRoL2FwaVwiXG5pbXBvcnQgeyBJbmRleEFQSSB9IGZyb20gXCIuL2FwaXMvaW5kZXgvYXBpXCJcbmltcG9ydCB7IEluZm9BUEkgfSBmcm9tIFwiLi9hcGlzL2luZm8vYXBpXCJcbmltcG9ydCB7IEtleXN0b3JlQVBJIH0gZnJvbSBcIi4vYXBpcy9rZXlzdG9yZS9hcGlcIlxuaW1wb3J0IHsgTWV0cmljc0FQSSB9IGZyb20gXCIuL2FwaXMvbWV0cmljcy9hcGlcIlxuaW1wb3J0IHsgUGxhdGZvcm1WTUFQSSB9IGZyb20gXCIuL2FwaXMvcGxhdGZvcm12bS9hcGlcIlxuaW1wb3J0IHsgU29ja2V0IH0gZnJvbSBcIi4vYXBpcy9zb2NrZXQvc29ja2V0XCJcbmltcG9ydCB7IERlZmF1bHROZXR3b3JrSUQsIERlZmF1bHRzIH0gZnJvbSBcIi4vdXRpbHMvY29uc3RhbnRzXCJcbmltcG9ydCB7IGdldFByZWZlcnJlZEhSUCB9IGZyb20gXCIuL3V0aWxzL2hlbHBlcmZ1bmN0aW9uc1wiXG5pbXBvcnQgQmluVG9vbHMgZnJvbSBcIi4vdXRpbHMvYmludG9vbHNcIlxuaW1wb3J0IERCIGZyb20gXCIuL3V0aWxzL2RiXCJcbmltcG9ydCBNbmVtb25pYyBmcm9tIFwiLi91dGlscy9tbmVtb25pY1wiXG5pbXBvcnQgUHViU3ViIGZyb20gXCIuL3V0aWxzL3B1YnN1YlwiXG5pbXBvcnQgSEROb2RlIGZyb20gXCIuL3V0aWxzL2hkbm9kZVwiXG5pbXBvcnQgQk4gZnJvbSBcImJuLmpzXCJcbmltcG9ydCB7IEJ1ZmZlciB9IGZyb20gXCJidWZmZXIvXCJcblxuLyoqXG4gKiBBdmFsYW5jaGVKUyBpcyBtaWRkbGV3YXJlIGZvciBpbnRlcmFjdGluZyB3aXRoIEF2YWxhbmNoZSBub2RlIFJQQyBBUElzLlxuICpcbiAqIEV4YW1wbGUgdXNhZ2U6XG4gKiBgYGBqc1xuICogY29uc3QgYXZhbGFuY2hlOiBBdmFsYW5jaGUgPSBuZXcgQXZhbGFuY2hlKFwiMTI3LjAuMC4xXCIsIDk2NTAsIFwiaHR0cHNcIilcbiAqIGBgYFxuICpcbiAqL1xuZXhwb3J0IGRlZmF1bHQgY2xhc3MgQXZhbGFuY2hlIGV4dGVuZHMgQXZhbGFuY2hlQ29yZSB7XG4gIC8qKlxuICAgKiBSZXR1cm5zIGEgcmVmZXJlbmNlIHRvIHRoZSBBZG1pbiBSUEMuXG4gICAqL1xuICBBZG1pbiA9ICgpID0+IHRoaXMuYXBpcy5hZG1pbiBhcyBBZG1pbkFQSVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIGEgcmVmZXJlbmNlIHRvIHRoZSBBdXRoIFJQQy5cbiAgICovXG4gIEF1dGggPSAoKSA9PiB0aGlzLmFwaXMuYXV0aCBhcyBBdXRoQVBJXG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSByZWZlcmVuY2UgdG8gdGhlIEVWTUFQSSBSUEMgcG9pbnRlZCBhdCB0aGUgQy1DaGFpbi5cbiAgICovXG4gIENDaGFpbiA9ICgpID0+IHRoaXMuYXBpcy5jY2hhaW4gYXMgRVZNQVBJXG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSByZWZlcmVuY2UgdG8gdGhlIEFWTSBSUEMgcG9pbnRlZCBhdCB0aGUgWC1DaGFpbi5cbiAgICovXG4gIFhDaGFpbiA9ICgpID0+IHRoaXMuYXBpcy54Y2hhaW4gYXMgQVZNQVBJXG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSByZWZlcmVuY2UgdG8gdGhlIEhlYWx0aCBSUEMgZm9yIGEgbm9kZS5cbiAgICovXG4gIEhlYWx0aCA9ICgpID0+IHRoaXMuYXBpcy5oZWFsdGggYXMgSGVhbHRoQVBJXG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSByZWZlcmVuY2UgdG8gdGhlIEluZGV4IFJQQyBmb3IgYSBub2RlLlxuICAgKi9cbiAgSW5kZXggPSAoKSA9PiB0aGlzLmFwaXMuaW5kZXggYXMgSW5kZXhBUElcblxuICAvKipcbiAgICogUmV0dXJucyBhIHJlZmVyZW5jZSB0byB0aGUgSW5mbyBSUEMgZm9yIGEgbm9kZS5cbiAgICovXG4gIEluZm8gPSAoKSA9PiB0aGlzLmFwaXMuaW5mbyBhcyBJbmZvQVBJXG5cbiAgLyoqXG4gICAqIFJldHVybnMgYSByZWZlcmVuY2UgdG8gdGhlIE1ldHJpY3MgUlBDLlxuICAgKi9cbiAgTWV0cmljcyA9ICgpID0+IHRoaXMuYXBpcy5tZXRyaWNzIGFzIE1ldHJpY3NBUElcblxuICAvKipcbiAgICogUmV0dXJucyBhIHJlZmVyZW5jZSB0byB0aGUgS2V5c3RvcmUgUlBDIGZvciBhIG5vZGUuIFdlIGxhYmVsIGl0IFwiTm9kZUtleXNcIiB0byByZWR1Y2VcbiAgICogY29uZnVzaW9uIGFib3V0IHdoYXQgaXQncyBhY2Nlc3NpbmcuXG4gICAqL1xuICBOb2RlS2V5cyA9ICgpID0+IHRoaXMuYXBpcy5rZXlzdG9yZSBhcyBLZXlzdG9yZUFQSVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIGEgcmVmZXJlbmNlIHRvIHRoZSBQbGF0Zm9ybVZNIFJQQyBwb2ludGVkIGF0IHRoZSBQLUNoYWluLlxuICAgKi9cbiAgUENoYWluID0gKCkgPT4gdGhpcy5hcGlzLnBjaGFpbiBhcyBQbGF0Zm9ybVZNQVBJXG5cbiAgLyoqXG4gICAqIENyZWF0ZXMgYSBuZXcgQXZhbGFuY2hlIGluc3RhbmNlLiBTZXRzIHRoZSBhZGRyZXNzIGFuZCBwb3J0IG9mIHRoZSBtYWluIEF2YWxhbmNoZSBDbGllbnQuXG4gICAqXG4gICAqIEBwYXJhbSBob3N0IFRoZSBob3N0bmFtZSB0byByZXNvbHZlIHRvIHJlYWNoIHRoZSBBdmFsYW5jaGUgQ2xpZW50IFJQQyBBUElzXG4gICAqIEBwYXJhbSBwb3J0IFRoZSBwb3J0IHRvIHJlc29sdmUgdG8gcmVhY2ggdGhlIEF2YWxhbmNoZSBDbGllbnQgUlBDIEFQSXNcbiAgICogQHBhcmFtIHByb3RvY29sIFRoZSBwcm90b2NvbCBzdHJpbmcgdG8gdXNlIGJlZm9yZSBhIFwiOi8vXCIgaW4gYSByZXF1ZXN0LFxuICAgKiBleDogXCJodHRwXCIsIFwiaHR0cHNcIiwgXCJnaXRcIiwgXCJ3c1wiLCBldGMuIERlZmF1bHRzIHRvIGh0dHBcbiAgICogQHBhcmFtIG5ldHdvcmtJRCBTZXRzIHRoZSBOZXR3b3JrSUQgb2YgdGhlIGNsYXNzLiBEZWZhdWx0IFtbRGVmYXVsdE5ldHdvcmtJRF1dXG4gICAqIEBwYXJhbSBYQ2hhaW5JRCBTZXRzIHRoZSBibG9ja2NoYWluSUQgZm9yIHRoZSBBVk0uIFdpbGwgdHJ5IHRvIGF1dG8tZGV0ZWN0LFxuICAgKiBvdGhlcndpc2UgZGVmYXVsdCBcIjJlTnkxbVVGZG1heFhOajFlUUhVZTdOcDRnanU5c0pzRXRXUTRNWDNUb2lOS3VBRGVkXCJcbiAgICogQHBhcmFtIENDaGFpbklEIFNldHMgdGhlIGJsb2NrY2hhaW5JRCBmb3IgdGhlIEVWTS4gV2lsbCB0cnkgdG8gYXV0by1kZXRlY3QsXG4gICAqIG90aGVyd2lzZSBkZWZhdWx0IFwiMkNBNmo1ell6YXN5blBzRmVOb3FXa21UQ3QzVlNjTXZYVVpIYmZESjhrM29HekFQdFVcIlxuICAgKiBAcGFyYW0gaHJwIFRoZSBodW1hbi1yZWFkYWJsZSBwYXJ0IG9mIHRoZSBiZWNoMzIgYWRkcmVzc2VzXG4gICAqIEBwYXJhbSBza2lwaW5pdCBTa2lwcyBjcmVhdGluZyB0aGUgQVBJcy4gRGVmYXVsdHMgdG8gZmFsc2VcbiAgICovXG4gIGNvbnN0cnVjdG9yKFxuICAgIGhvc3Q/OiBzdHJpbmcsXG4gICAgcG9ydD86IG51bWJlcixcbiAgICBwcm90b2NvbDogc3RyaW5nID0gXCJodHRwXCIsXG4gICAgbmV0d29ya0lEOiBudW1iZXIgPSBEZWZhdWx0TmV0d29ya0lELFxuICAgIFhDaGFpbklEOiBzdHJpbmcgPSB1bmRlZmluZWQsXG4gICAgQ0NoYWluSUQ6IHN0cmluZyA9IHVuZGVmaW5lZCxcbiAgICBocnA6IHN0cmluZyA9IHVuZGVmaW5lZCxcbiAgICBza2lwaW5pdDogYm9vbGVhbiA9IGZhbHNlXG4gICkge1xuICAgIHN1cGVyKGhvc3QsIHBvcnQsIHByb3RvY29sKVxuICAgIGxldCB4Y2hhaW5pZDogc3RyaW5nID0gWENoYWluSURcbiAgICBsZXQgY2NoYWluaWQ6IHN0cmluZyA9IENDaGFpbklEXG5cbiAgICBpZiAoXG4gICAgICB0eXBlb2YgWENoYWluSUQgPT09IFwidW5kZWZpbmVkXCIgfHxcbiAgICAgICFYQ2hhaW5JRCB8fFxuICAgICAgWENoYWluSUQudG9Mb3dlckNhc2UoKSA9PT0gXCJ4XCJcbiAgICApIHtcbiAgICAgIGlmIChuZXR3b3JrSUQudG9TdHJpbmcoKSBpbiBEZWZhdWx0cy5uZXR3b3JrKSB7XG4gICAgICAgIHhjaGFpbmlkID0gRGVmYXVsdHMubmV0d29ya1tgJHtuZXR3b3JrSUR9YF0uWC5ibG9ja2NoYWluSURcbiAgICAgIH0gZWxzZSB7XG4gICAgICAgIHhjaGFpbmlkID0gRGVmYXVsdHMubmV0d29ya1sxMjM0NV0uWC5ibG9ja2NoYWluSURcbiAgICAgIH1cbiAgICB9XG4gICAgaWYgKFxuICAgICAgdHlwZW9mIENDaGFpbklEID09PSBcInVuZGVmaW5lZFwiIHx8XG4gICAgICAhQ0NoYWluSUQgfHxcbiAgICAgIENDaGFpbklELnRvTG93ZXJDYXNlKCkgPT09IFwiY1wiXG4gICAgKSB7XG4gICAgICBpZiAobmV0d29ya0lELnRvU3RyaW5nKCkgaW4gRGVmYXVsdHMubmV0d29yaykge1xuICAgICAgICBjY2hhaW5pZCA9IERlZmF1bHRzLm5ldHdvcmtbYCR7bmV0d29ya0lEfWBdLkMuYmxvY2tjaGFpbklEXG4gICAgICB9IGVsc2Uge1xuICAgICAgICBjY2hhaW5pZCA9IERlZmF1bHRzLm5ldHdvcmtbMTIzNDVdLkMuYmxvY2tjaGFpbklEXG4gICAgICB9XG4gICAgfVxuICAgIGlmICh0eXBlb2YgbmV0d29ya0lEID09PSBcIm51bWJlclwiICYmIG5ldHdvcmtJRCA+PSAwKSB7XG4gICAgICB0aGlzLm5ldHdvcmtJRCA9IG5ldHdvcmtJRFxuICAgIH0gZWxzZSBpZiAodHlwZW9mIG5ldHdvcmtJRCA9PT0gXCJ1bmRlZmluZWRcIikge1xuICAgICAgbmV0d29ya0lEID0gRGVmYXVsdE5ldHdvcmtJRFxuICAgIH1cbiAgICBpZiAodHlwZW9mIGhycCAhPT0gXCJ1bmRlZmluZWRcIikge1xuICAgICAgdGhpcy5ocnAgPSBocnBcbiAgICB9IGVsc2Uge1xuICAgICAgdGhpcy5ocnAgPSBnZXRQcmVmZXJyZWRIUlAodGhpcy5uZXR3b3JrSUQpXG4gICAgfVxuXG4gICAgaWYgKCFza2lwaW5pdCkge1xuICAgICAgdGhpcy5hZGRBUEkoXCJhZG1pblwiLCBBZG1pbkFQSSlcbiAgICAgIHRoaXMuYWRkQVBJKFwiYXV0aFwiLCBBdXRoQVBJKVxuICAgICAgdGhpcy5hZGRBUEkoXCJ4Y2hhaW5cIiwgQVZNQVBJLCBcIi9leHQvYmMvWFwiLCB4Y2hhaW5pZClcbiAgICAgIHRoaXMuYWRkQVBJKFwiY2NoYWluXCIsIEVWTUFQSSwgXCIvZXh0L2JjL0MvYXZheFwiLCBjY2hhaW5pZClcbiAgICAgIHRoaXMuYWRkQVBJKFwiaGVhbHRoXCIsIEhlYWx0aEFQSSlcbiAgICAgIHRoaXMuYWRkQVBJKFwiaW5mb1wiLCBJbmZvQVBJKVxuICAgICAgdGhpcy5hZGRBUEkoXCJpbmRleFwiLCBJbmRleEFQSSlcbiAgICAgIHRoaXMuYWRkQVBJKFwia2V5c3RvcmVcIiwgS2V5c3RvcmVBUEkpXG4gICAgICB0aGlzLmFkZEFQSShcIm1ldHJpY3NcIiwgTWV0cmljc0FQSSlcbiAgICAgIHRoaXMuYWRkQVBJKFwicGNoYWluXCIsIFBsYXRmb3JtVk1BUEkpXG4gICAgfVxuICB9XG59XG5cbmV4cG9ydCB7IEF2YWxhbmNoZSB9XG5leHBvcnQgeyBBdmFsYW5jaGVDb3JlIH1cbmV4cG9ydCB7IEJpblRvb2xzIH1cbmV4cG9ydCB7IEJOIH1cbmV4cG9ydCB7IEJ1ZmZlciB9XG5leHBvcnQgeyBEQiB9XG5leHBvcnQgeyBIRE5vZGUgfVxuZXhwb3J0IHsgR2VuZXNpc0Fzc2V0IH1cbmV4cG9ydCB7IEdlbmVzaXNEYXRhIH1cbmV4cG9ydCB7IE1uZW1vbmljIH1cbmV4cG9ydCB7IFB1YlN1YiB9XG5leHBvcnQgeyBTb2NrZXQgfVxuXG5leHBvcnQgKiBhcyBhZG1pbiBmcm9tIFwiLi9hcGlzL2FkbWluXCJcbmV4cG9ydCAqIGFzIGF1dGggZnJvbSBcIi4vYXBpcy9hdXRoXCJcbmV4cG9ydCAqIGFzIGF2bSBmcm9tIFwiLi9hcGlzL2F2bVwiXG5leHBvcnQgKiBhcyBjb21tb24gZnJvbSBcIi4vY29tbW9uXCJcbmV4cG9ydCAqIGFzIGV2bSBmcm9tIFwiLi9hcGlzL2V2bVwiXG5leHBvcnQgKiBhcyBoZWFsdGggZnJvbSBcIi4vYXBpcy9oZWFsdGhcIlxuZXhwb3J0ICogYXMgaW5kZXggZnJvbSBcIi4vYXBpcy9pbmRleFwiXG5leHBvcnQgKiBhcyBpbmZvIGZyb20gXCIuL2FwaXMvaW5mb1wiXG5leHBvcnQgKiBhcyBrZXlzdG9yZSBmcm9tIFwiLi9hcGlzL2tleXN0b3JlXCJcbmV4cG9ydCAqIGFzIG1ldHJpY3MgZnJvbSBcIi4vYXBpcy9tZXRyaWNzXCJcbmV4cG9ydCAqIGFzIHBsYXRmb3Jtdm0gZnJvbSBcIi4vYXBpcy9wbGF0Zm9ybXZtXCJcbmV4cG9ydCAqIGFzIHV0aWxzIGZyb20gXCIuL3V0aWxzXCJcbiJdfQ==