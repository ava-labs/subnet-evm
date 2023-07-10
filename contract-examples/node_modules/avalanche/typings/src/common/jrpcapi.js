"use strict";
/**
 * @packageDocumentation
 * @module Common-JRPCAPI
 */
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.JRPCAPI = void 0;
const utils_1 = require("../utils");
const apibase_1 = require("./apibase");
class JRPCAPI extends apibase_1.APIBase {
    /**
     *
     * @param core Reference to the Avalanche instance using this endpoint
     * @param baseURL Path of the APIs baseURL - ex: "/ext/bc/avm"
     * @param jrpcVersion The jrpc version to use, default "2.0".
     */
    constructor(core, baseURL, jrpcVersion = "2.0") {
        super(core, baseURL);
        this.jrpcVersion = "2.0";
        this.rpcID = 1;
        this.callMethod = (method, params, baseURL, headers) => __awaiter(this, void 0, void 0, function* () {
            const ep = baseURL || this.baseURL;
            const rpc = {};
            rpc.id = this.rpcID;
            rpc.method = method;
            // Set parameters if exists
            if (params) {
                rpc.params = params;
            }
            else if (this.jrpcVersion === "1.0") {
                rpc.params = [];
            }
            if (this.jrpcVersion !== "1.0") {
                rpc.jsonrpc = this.jrpcVersion;
            }
            let headrs = { "Content-Type": "application/json;charset=UTF-8" };
            if (headers) {
                headrs = Object.assign(Object.assign({}, headrs), headers);
            }
            baseURL = this.core.getURL();
            const axConf = {
                baseURL: baseURL,
                responseType: "json",
                // use the fetch adapter if fetch is available e.g. non Node<17 env
                adapter: typeof fetch !== "undefined" ? utils_1.fetchAdapter : undefined
            };
            const resp = yield this.core.post(ep, {}, JSON.stringify(rpc), headrs, axConf);
            if (resp.status >= 200 && resp.status < 300) {
                this.rpcID += 1;
                if (typeof resp.data === "string") {
                    resp.data = JSON.parse(resp.data);
                }
                if (typeof resp.data === "object" &&
                    (resp.data === null || "error" in resp.data)) {
                    throw new Error(resp.data.error.message);
                }
            }
            return resp;
        });
        /**
         * Returns the rpcid, a strictly-increasing number, starting from 1, indicating the next
         * request ID that will be sent.
         */
        this.getRPCID = () => this.rpcID;
        this.jrpcVersion = jrpcVersion;
        this.rpcID = 1;
    }
}
exports.JRPCAPI = JRPCAPI;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoianJwY2FwaS5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3NyYy9jb21tb24vanJwY2FwaS50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiO0FBQUE7OztHQUdHOzs7Ozs7Ozs7Ozs7QUFHSCxvQ0FBdUM7QUFFdkMsdUNBQXdEO0FBRXhELE1BQWEsT0FBUSxTQUFRLGlCQUFPO0lBb0VsQzs7Ozs7T0FLRztJQUNILFlBQ0UsSUFBbUIsRUFDbkIsT0FBZSxFQUNmLGNBQXNCLEtBQUs7UUFFM0IsS0FBSyxDQUFDLElBQUksRUFBRSxPQUFPLENBQUMsQ0FBQTtRQTlFWixnQkFBVyxHQUFXLEtBQUssQ0FBQTtRQUMzQixVQUFLLEdBQUcsQ0FBQyxDQUFBO1FBRW5CLGVBQVUsR0FBRyxDQUNYLE1BQWMsRUFDZCxNQUEwQixFQUMxQixPQUFnQixFQUNoQixPQUFnQixFQUNjLEVBQUU7WUFDaEMsTUFBTSxFQUFFLEdBQVcsT0FBTyxJQUFJLElBQUksQ0FBQyxPQUFPLENBQUE7WUFDMUMsTUFBTSxHQUFHLEdBQVEsRUFBRSxDQUFBO1lBQ25CLEdBQUcsQ0FBQyxFQUFFLEdBQUcsSUFBSSxDQUFDLEtBQUssQ0FBQTtZQUNuQixHQUFHLENBQUMsTUFBTSxHQUFHLE1BQU0sQ0FBQTtZQUVuQiwyQkFBMkI7WUFDM0IsSUFBSSxNQUFNLEVBQUU7Z0JBQ1YsR0FBRyxDQUFDLE1BQU0sR0FBRyxNQUFNLENBQUE7YUFDcEI7aUJBQU0sSUFBSSxJQUFJLENBQUMsV0FBVyxLQUFLLEtBQUssRUFBRTtnQkFDckMsR0FBRyxDQUFDLE1BQU0sR0FBRyxFQUFFLENBQUE7YUFDaEI7WUFFRCxJQUFJLElBQUksQ0FBQyxXQUFXLEtBQUssS0FBSyxFQUFFO2dCQUM5QixHQUFHLENBQUMsT0FBTyxHQUFHLElBQUksQ0FBQyxXQUFXLENBQUE7YUFDL0I7WUFFRCxJQUFJLE1BQU0sR0FBVyxFQUFFLGNBQWMsRUFBRSxnQ0FBZ0MsRUFBRSxDQUFBO1lBQ3pFLElBQUksT0FBTyxFQUFFO2dCQUNYLE1BQU0sbUNBQVEsTUFBTSxHQUFLLE9BQU8sQ0FBRSxDQUFBO2FBQ25DO1lBRUQsT0FBTyxHQUFHLElBQUksQ0FBQyxJQUFJLENBQUMsTUFBTSxFQUFFLENBQUE7WUFFNUIsTUFBTSxNQUFNLEdBQXVCO2dCQUNqQyxPQUFPLEVBQUUsT0FBTztnQkFDaEIsWUFBWSxFQUFFLE1BQU07Z0JBQ3BCLG1FQUFtRTtnQkFDbkUsT0FBTyxFQUFFLE9BQU8sS0FBSyxLQUFLLFdBQVcsQ0FBQyxDQUFDLENBQUMsb0JBQVksQ0FBQyxDQUFDLENBQUMsU0FBUzthQUNqRSxDQUFBO1lBRUQsTUFBTSxJQUFJLEdBQXdCLE1BQU0sSUFBSSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQ3BELEVBQUUsRUFDRixFQUFFLEVBQ0YsSUFBSSxDQUFDLFNBQVMsQ0FBQyxHQUFHLENBQUMsRUFDbkIsTUFBTSxFQUNOLE1BQU0sQ0FDUCxDQUFBO1lBQ0QsSUFBSSxJQUFJLENBQUMsTUFBTSxJQUFJLEdBQUcsSUFBSSxJQUFJLENBQUMsTUFBTSxHQUFHLEdBQUcsRUFBRTtnQkFDM0MsSUFBSSxDQUFDLEtBQUssSUFBSSxDQUFDLENBQUE7Z0JBQ2YsSUFBSSxPQUFPLElBQUksQ0FBQyxJQUFJLEtBQUssUUFBUSxFQUFFO29CQUNqQyxJQUFJLENBQUMsSUFBSSxHQUFHLElBQUksQ0FBQyxLQUFLLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO2lCQUNsQztnQkFDRCxJQUNFLE9BQU8sSUFBSSxDQUFDLElBQUksS0FBSyxRQUFRO29CQUM3QixDQUFDLElBQUksQ0FBQyxJQUFJLEtBQUssSUFBSSxJQUFJLE9BQU8sSUFBSSxJQUFJLENBQUMsSUFBSSxDQUFDLEVBQzVDO29CQUNBLE1BQU0sSUFBSSxLQUFLLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxLQUFLLENBQUMsT0FBTyxDQUFDLENBQUE7aUJBQ3pDO2FBQ0Y7WUFDRCxPQUFPLElBQUksQ0FBQTtRQUNiLENBQUMsQ0FBQSxDQUFBO1FBRUQ7OztXQUdHO1FBQ0gsYUFBUSxHQUFHLEdBQVcsRUFBRSxDQUFDLElBQUksQ0FBQyxLQUFLLENBQUE7UUFjakMsSUFBSSxDQUFDLFdBQVcsR0FBRyxXQUFXLENBQUE7UUFDOUIsSUFBSSxDQUFDLEtBQUssR0FBRyxDQUFDLENBQUE7SUFDaEIsQ0FBQztDQUNGO0FBbkZELDBCQW1GQyIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICogQHBhY2thZ2VEb2N1bWVudGF0aW9uXG4gKiBAbW9kdWxlIENvbW1vbi1KUlBDQVBJXG4gKi9cblxuaW1wb3J0IHsgQXhpb3NSZXF1ZXN0Q29uZmlnIH0gZnJvbSBcImF4aW9zXCJcbmltcG9ydCB7IGZldGNoQWRhcHRlciB9IGZyb20gXCIuLi91dGlsc1wiXG5pbXBvcnQgQXZhbGFuY2hlQ29yZSBmcm9tIFwiLi4vYXZhbGFuY2hlXCJcbmltcG9ydCB7IEFQSUJhc2UsIFJlcXVlc3RSZXNwb25zZURhdGEgfSBmcm9tIFwiLi9hcGliYXNlXCJcblxuZXhwb3J0IGNsYXNzIEpSUENBUEkgZXh0ZW5kcyBBUElCYXNlIHtcbiAgcHJvdGVjdGVkIGpycGNWZXJzaW9uOiBzdHJpbmcgPSBcIjIuMFwiXG4gIHByb3RlY3RlZCBycGNJRCA9IDFcblxuICBjYWxsTWV0aG9kID0gYXN5bmMgKFxuICAgIG1ldGhvZDogc3RyaW5nLFxuICAgIHBhcmFtcz86IG9iamVjdFtdIHwgb2JqZWN0LFxuICAgIGJhc2VVUkw/OiBzdHJpbmcsXG4gICAgaGVhZGVycz86IG9iamVjdFxuICApOiBQcm9taXNlPFJlcXVlc3RSZXNwb25zZURhdGE+ID0+IHtcbiAgICBjb25zdCBlcDogc3RyaW5nID0gYmFzZVVSTCB8fCB0aGlzLmJhc2VVUkxcbiAgICBjb25zdCBycGM6IGFueSA9IHt9XG4gICAgcnBjLmlkID0gdGhpcy5ycGNJRFxuICAgIHJwYy5tZXRob2QgPSBtZXRob2RcblxuICAgIC8vIFNldCBwYXJhbWV0ZXJzIGlmIGV4aXN0c1xuICAgIGlmIChwYXJhbXMpIHtcbiAgICAgIHJwYy5wYXJhbXMgPSBwYXJhbXNcbiAgICB9IGVsc2UgaWYgKHRoaXMuanJwY1ZlcnNpb24gPT09IFwiMS4wXCIpIHtcbiAgICAgIHJwYy5wYXJhbXMgPSBbXVxuICAgIH1cblxuICAgIGlmICh0aGlzLmpycGNWZXJzaW9uICE9PSBcIjEuMFwiKSB7XG4gICAgICBycGMuanNvbnJwYyA9IHRoaXMuanJwY1ZlcnNpb25cbiAgICB9XG5cbiAgICBsZXQgaGVhZHJzOiBvYmplY3QgPSB7IFwiQ29udGVudC1UeXBlXCI6IFwiYXBwbGljYXRpb24vanNvbjtjaGFyc2V0PVVURi04XCIgfVxuICAgIGlmIChoZWFkZXJzKSB7XG4gICAgICBoZWFkcnMgPSB7IC4uLmhlYWRycywgLi4uaGVhZGVycyB9XG4gICAgfVxuXG4gICAgYmFzZVVSTCA9IHRoaXMuY29yZS5nZXRVUkwoKVxuXG4gICAgY29uc3QgYXhDb25mOiBBeGlvc1JlcXVlc3RDb25maWcgPSB7XG4gICAgICBiYXNlVVJMOiBiYXNlVVJMLFxuICAgICAgcmVzcG9uc2VUeXBlOiBcImpzb25cIixcbiAgICAgIC8vIHVzZSB0aGUgZmV0Y2ggYWRhcHRlciBpZiBmZXRjaCBpcyBhdmFpbGFibGUgZS5nLiBub24gTm9kZTwxNyBlbnZcbiAgICAgIGFkYXB0ZXI6IHR5cGVvZiBmZXRjaCAhPT0gXCJ1bmRlZmluZWRcIiA/IGZldGNoQWRhcHRlciA6IHVuZGVmaW5lZFxuICAgIH1cblxuICAgIGNvbnN0IHJlc3A6IFJlcXVlc3RSZXNwb25zZURhdGEgPSBhd2FpdCB0aGlzLmNvcmUucG9zdChcbiAgICAgIGVwLFxuICAgICAge30sXG4gICAgICBKU09OLnN0cmluZ2lmeShycGMpLFxuICAgICAgaGVhZHJzLFxuICAgICAgYXhDb25mXG4gICAgKVxuICAgIGlmIChyZXNwLnN0YXR1cyA+PSAyMDAgJiYgcmVzcC5zdGF0dXMgPCAzMDApIHtcbiAgICAgIHRoaXMucnBjSUQgKz0gMVxuICAgICAgaWYgKHR5cGVvZiByZXNwLmRhdGEgPT09IFwic3RyaW5nXCIpIHtcbiAgICAgICAgcmVzcC5kYXRhID0gSlNPTi5wYXJzZShyZXNwLmRhdGEpXG4gICAgICB9XG4gICAgICBpZiAoXG4gICAgICAgIHR5cGVvZiByZXNwLmRhdGEgPT09IFwib2JqZWN0XCIgJiZcbiAgICAgICAgKHJlc3AuZGF0YSA9PT0gbnVsbCB8fCBcImVycm9yXCIgaW4gcmVzcC5kYXRhKVxuICAgICAgKSB7XG4gICAgICAgIHRocm93IG5ldyBFcnJvcihyZXNwLmRhdGEuZXJyb3IubWVzc2FnZSlcbiAgICAgIH1cbiAgICB9XG4gICAgcmV0dXJuIHJlc3BcbiAgfVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBycGNpZCwgYSBzdHJpY3RseS1pbmNyZWFzaW5nIG51bWJlciwgc3RhcnRpbmcgZnJvbSAxLCBpbmRpY2F0aW5nIHRoZSBuZXh0XG4gICAqIHJlcXVlc3QgSUQgdGhhdCB3aWxsIGJlIHNlbnQuXG4gICAqL1xuICBnZXRSUENJRCA9ICgpOiBudW1iZXIgPT4gdGhpcy5ycGNJRFxuXG4gIC8qKlxuICAgKlxuICAgKiBAcGFyYW0gY29yZSBSZWZlcmVuY2UgdG8gdGhlIEF2YWxhbmNoZSBpbnN0YW5jZSB1c2luZyB0aGlzIGVuZHBvaW50XG4gICAqIEBwYXJhbSBiYXNlVVJMIFBhdGggb2YgdGhlIEFQSXMgYmFzZVVSTCAtIGV4OiBcIi9leHQvYmMvYXZtXCJcbiAgICogQHBhcmFtIGpycGNWZXJzaW9uIFRoZSBqcnBjIHZlcnNpb24gdG8gdXNlLCBkZWZhdWx0IFwiMi4wXCIuXG4gICAqL1xuICBjb25zdHJ1Y3RvcihcbiAgICBjb3JlOiBBdmFsYW5jaGVDb3JlLFxuICAgIGJhc2VVUkw6IHN0cmluZyxcbiAgICBqcnBjVmVyc2lvbjogc3RyaW5nID0gXCIyLjBcIlxuICApIHtcbiAgICBzdXBlcihjb3JlLCBiYXNlVVJMKVxuICAgIHRoaXMuanJwY1ZlcnNpb24gPSBqcnBjVmVyc2lvblxuICAgIHRoaXMucnBjSUQgPSAxXG4gIH1cbn1cbiJdfQ==