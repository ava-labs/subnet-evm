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
Object.defineProperty(exports, "__esModule", { value: true });
exports.AdminAPI = void 0;
const jrpcapi_1 = require("../../common/jrpcapi");
/**
 * Class for interacting with a node's AdminAPI.
 *
 * @category RPCAPIs
 *
 * @remarks This extends the [[JRPCAPI]] class. This class should not be directly called.
 * Instead, use the [[Avalanche.addAPI]] function to register this interface with Avalanche.
 */
class AdminAPI extends jrpcapi_1.JRPCAPI {
    /**
     * This class should not be instantiated directly. Instead use the [[Avalanche.addAPI]]
     * method.
     *
     * @param core A reference to the Avalanche class
     * @param baseURL Defaults to the string "/ext/admin" as the path to rpc's baseURL
     */
    constructor(core, baseURL = "/ext/admin") {
        super(core, baseURL);
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
        this.alias = (endpoint, alias) => __awaiter(this, void 0, void 0, function* () {
            const params = {
                endpoint,
                alias
            };
            const response = yield this.callMethod("admin.alias", params);
            return response.data.result.success
                ? response.data.result.success
                : response.data.result;
        });
        /**
         * Give a blockchain an alias, a different name that can be used any place the blockchain’s
         * ID is used.
         *
         * @param chain The blockchain’s ID
         * @param alias Can now be used in place of the blockchain’s ID (in API endpoints, for example)
         *
         * @returns Returns a Promise boolean containing success, true for success, false for failure.
         */
        this.aliasChain = (chain, alias) => __awaiter(this, void 0, void 0, function* () {
            const params = {
                chain,
                alias
            };
            const response = yield this.callMethod("admin.aliasChain", params);
            return response.data.result.success
                ? response.data.result.success
                : response.data.result;
        });
        /**
         * Get all aliases for given blockchain
         *
         * @param chain The blockchain’s ID
         *
         * @returns Returns a Promise string[] containing aliases of the blockchain.
         */
        this.getChainAliases = (chain) => __awaiter(this, void 0, void 0, function* () {
            const params = {
                chain
            };
            const response = yield this.callMethod("admin.getChainAliases", params);
            return response.data.result.aliases
                ? response.data.result.aliases
                : response.data.result;
        });
        /**
         * Returns log and display levels of loggers
         *
         * @param loggerName the name of the logger to be returned. This is an optional argument. If not specified, it returns all possible loggers.
         *
         * @returns Returns a Promise containing logger levels
         */
        this.getLoggerLevel = (loggerName) => __awaiter(this, void 0, void 0, function* () {
            const params = {};
            if (typeof loggerName !== "undefined") {
                params.loggerName = loggerName;
            }
            const response = yield this.callMethod("admin.getLoggerLevel", params);
            return response.data.result;
        });
        /**
         * Dynamically loads any virtual machines installed on the node as plugins
         *
         * @returns Returns a Promise containing new VMs and failed VMs
         */
        this.loadVMs = () => __awaiter(this, void 0, void 0, function* () {
            const response = yield this.callMethod("admin.loadVMs");
            return response.data.result.aliases
                ? response.data.result.aliases
                : response.data.result;
        });
        /**
         * Dump the mutex statistics of the node to the specified file.
         *
         * @returns Promise for a boolean that is true on success.
         */
        this.lockProfile = () => __awaiter(this, void 0, void 0, function* () {
            const response = yield this.callMethod("admin.lockProfile");
            return response.data.result.success
                ? response.data.result.success
                : response.data.result;
        });
        /**
         * Dump the current memory footprint of the node to the specified file.
         *
         * @returns Promise for a boolean that is true on success.
         */
        this.memoryProfile = () => __awaiter(this, void 0, void 0, function* () {
            const response = yield this.callMethod("admin.memoryProfile");
            return response.data.result.success
                ? response.data.result.success
                : response.data.result;
        });
        /**
         * Sets log and display levels of loggers.
         *
         * @param loggerName the name of the logger to be changed. This is an optional parameter.
         * @param logLevel the log level of written logs, can be omitted.
         * @param displayLevel the log level of displayed logs, can be omitted.
         *
         * @returns Returns a Promise containing logger levels
         */
        this.setLoggerLevel = (loggerName, logLevel, displayLevel) => __awaiter(this, void 0, void 0, function* () {
            const params = {};
            if (typeof loggerName !== "undefined") {
                params.loggerName = loggerName;
            }
            if (typeof logLevel !== "undefined") {
                params.logLevel = logLevel;
            }
            if (typeof displayLevel !== "undefined") {
                params.displayLevel = displayLevel;
            }
            const response = yield this.callMethod("admin.setLoggerLevel", params);
            return response.data.result;
        });
        /**
         * Start profiling the cpu utilization of the node. Will dump the profile information into
         * the specified file on stop.
         *
         * @returns Promise for a boolean that is true on success.
         */
        this.startCPUProfiler = () => __awaiter(this, void 0, void 0, function* () {
            const response = yield this.callMethod("admin.startCPUProfiler");
            return response.data.result.success
                ? response.data.result.success
                : response.data.result;
        });
        /**
         * Stop the CPU profile that was previously started.
         *
         * @returns Promise for a boolean that is true on success.
         */
        this.stopCPUProfiler = () => __awaiter(this, void 0, void 0, function* () {
            const response = yield this.callMethod("admin.stopCPUProfiler");
            return response.data.result.success
                ? response.data.result.success
                : response.data.result;
        });
    }
}
exports.AdminAPI = AdminAPI;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYXBpLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2FwaXMvYWRtaW4vYXBpLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7Ozs7Ozs7OztBQUtBLGtEQUE4QztBQWE5Qzs7Ozs7OztHQU9HO0FBRUgsTUFBYSxRQUFTLFNBQVEsaUJBQU87SUE2TG5DOzs7Ozs7T0FNRztJQUNILFlBQVksSUFBbUIsRUFBRSxVQUFrQixZQUFZO1FBQzdELEtBQUssQ0FBQyxJQUFJLEVBQUUsT0FBTyxDQUFDLENBQUE7UUFwTXRCOzs7Ozs7Ozs7V0FTRztRQUNILFVBQUssR0FBRyxDQUFPLFFBQWdCLEVBQUUsS0FBYSxFQUFvQixFQUFFO1lBQ2xFLE1BQU0sTUFBTSxHQUFnQjtnQkFDMUIsUUFBUTtnQkFDUixLQUFLO2FBQ04sQ0FBQTtZQUNELE1BQU0sUUFBUSxHQUF3QixNQUFNLElBQUksQ0FBQyxVQUFVLENBQ3pELGFBQWEsRUFDYixNQUFNLENBQ1AsQ0FBQTtZQUNELE9BQU8sUUFBUSxDQUFDLElBQUksQ0FBQyxNQUFNLENBQUMsT0FBTztnQkFDakMsQ0FBQyxDQUFDLFFBQVEsQ0FBQyxJQUFJLENBQUMsTUFBTSxDQUFDLE9BQU87Z0JBQzlCLENBQUMsQ0FBQyxRQUFRLENBQUMsSUFBSSxDQUFDLE1BQU0sQ0FBQTtRQUMxQixDQUFDLENBQUEsQ0FBQTtRQUVEOzs7Ozs7OztXQVFHO1FBQ0gsZUFBVSxHQUFHLENBQU8sS0FBYSxFQUFFLEtBQWEsRUFBb0IsRUFBRTtZQUNwRSxNQUFNLE1BQU0sR0FBcUI7Z0JBQy9CLEtBQUs7Z0JBQ0wsS0FBSzthQUNOLENBQUE7WUFDRCxNQUFNLFFBQVEsR0FBd0IsTUFBTSxJQUFJLENBQUMsVUFBVSxDQUN6RCxrQkFBa0IsRUFDbEIsTUFBTSxDQUNQLENBQUE7WUFDRCxPQUFPLFFBQVEsQ0FBQyxJQUFJLENBQUMsTUFBTSxDQUFDLE9BQU87Z0JBQ2pDLENBQUMsQ0FBQyxRQUFRLENBQUMsSUFB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O1lBQ0QsSUFBSSxPQUFPLFlBQVksS0FBSyxXQUFXLEVBQUU7Z0JBQ3ZDLE1BQU0sQ0FBQyxZQUFZLEdBQUcsWUFBWSxDQUFBO2FBQ25DO1lBQ0QsTUFBTSxRQUFRLEdBQXdCLE1BQU0sSUFBSSxDQUFDLFVBQVUsQ0FDekQsc0JBQXNCLEVBQ3RCLE1BQU0sQ0FDUCxDQUFBO1lBQ0QsT0FBTyxRQUFRLENBQUMsSUFBSSxDQUFDLE1BQU0sQ0FBQTtRQUM3QixDQUFDLENBQUEsQ0FBQTtRQUVEOzs7OztXQUtHO1FBQ0gscUJBQWdCLEdBQUcsR0FBMkIsRUFBRTtZQUM5QyxNQUFNLFFBQVEsR0FBd0IsTUFBTSxJQUFJLENBQUMsVUFBVSxDQUN6RCx3QkFBd0IsQ0FDekIsQ0FBQTtZQUNELE9BQU8sUUFBUSxDQUFDLElBQUksQ0FBQyxNQUFNLENBQUMsT0FBTztnQkFDakMsQ0FBQyxDQUFDLFFBQVEsQ0FBQyxJQUFJLENBQUMsTUFBTSxDQUFDLE9BQU87Z0JBQzlCLENBQUMsQ0FBQyxRQUFRLENBQUMsSUFBSSxDQUFDLE1BQU0sQ0FBQTtRQUMxQixDQUFDLENBQUEsQ0FBQTtRQUVEOzs7O1dBSUc7UUFDSCxvQkFBZSxHQUFHLEdBQTJCLEVBQUU7WUFDN0MsTUFBTSxRQUFRLEdBQXdCLE1BQU0sSUFBSSxDQUFDLFVBQVUsQ0FDekQsdUJBQXVCLENBQ3hCLENBQUE7WUFDRCxPQUFPLFFBQVEsQ0FBQyxJQUFJLENBQUMsTUFBTSxDQUFDLE9BQU87Z0JBQ2pDLENBQUMsQ0FBQyxRQUFRLENBQUMsSUFBSSxDQUFDLE1BQU0sQ0FBQyxPQUFPO2dCQUM5QixDQUFDLENBQUMsUUFBUSxDQUFDLElBQUksQ0FBQyxNQUFNLENBQUE7UUFDMUIsQ0FBQyxDQUFBLENBQUE7SUFXRCxDQUFDO0NBQ0Y7QUF2TUQsNEJBdU1DIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgQVBJLUFkbWluXG4gKi9cbmltcG9ydCBBdmFsYW5jaGVDb3JlIGZyb20gXCIuLi8uLi9hdmFsYW5jaGVcIlxuaW1wb3J0IHsgSlJQQ0FQSSB9IGZyb20gXCIuLi8uLi9jb21tb24vanJwY2FwaVwiXG5pbXBvcnQgeyBSZXF1ZXN0UmVzcG9uc2VEYXRhIH0gZnJvbSBcIi4uLy4uL2NvbW1vbi9hcGliYXNlXCJcbmltcG9ydCB7XG4gIEFsaWFzQ2hhaW5QYXJhbXMsXG4gIEFsaWFzUGFyYW1zLFxuICBHZXRDaGFpbkFsaWFzZXNQYXJhbXMsXG4gIEdldExvZ2dlckxldmVsUGFyYW1zLFxuICBHZXRMb2dnZXJMZXZlbFJlc3BvbnNlLFxuICBMb2FkVk1zUmVzcG9uc2UsXG4gIFNldExvZ2dlckxldmVsUGFyYW1zLFxuICBTZXRMb2dnZXJMZXZlbFJlc3BvbnNlXG59IGZyb20gXCIuL2ludGVyZmFjZXNcIlxuXG4vKipcbiAqIENsYXNzIGZvciBpbnRlcmFjdGluZyB3aXRoIGEgbm9kZSdzIEFkbWluQVBJLlxuICpcbiAqIEBjYXRlZ29yeSBSUENBUElzXG4gKlxuICogQHJlbWFya3MgVGhpcyBleHRlbmRzIHRoZSBbW0pSUENBUEldXSBjbGFzcy4gVGhpcyBjbGFzcyBzaG91bGQgbm90IGJlIGRpcmVjdGx5IGNhbGxlZC5cbiAqIEluc3RlYWQsIHVzZSB0aGUgW1tBdmFsYW5jaGUuYWRkQVBJXV0gZnVuY3Rpb24gdG8gcmVnaXN0ZXIgdGhpcyBpbnRlcmZhY2Ugd2l0aCBBdmFsYW5jaGUuXG4gKi9cblxuZXhwb3J0IGNsYXNzIEFkbWluQVBJIGV4dGVuZHMgSlJQQ0FQSSB7XG4gIC8qKlxuICAgKiBBc3NpZ24gYW4gQVBJIGFuIGFsaWFzLCBhIGRpZmZlcmVudCBlbmRwb2ludCBmb3IgdGhlIEFQSS4gVGhlIG9yaWdpbmFsIGVuZHBvaW50IHdpbGwgc3RpbGxcbiAgICogd29yay4gVGhpcyBjaGFuZ2Ugb25seSBhZmZlY3RzIHRoaXMgbm9kZSBvdGhlciBub2RlcyB3aWxsIG5vdCBrbm93IGFib3V0IHRoaXMgYWxpYXMuXG4gICAqXG4gICAqIEBwYXJhbSBlbmRwb2ludCBUaGUgb3JpZ2luYWwgZW5kcG9pbnQgb2YgdGhlIEFQSS4gZW5kcG9pbnQgc2hvdWxkIG9ubHkgaW5jbHVkZSB0aGUgcGFydCBvZlxuICAgKiB0aGUgZW5kcG9pbnQgYWZ0ZXIgL2V4dC9cbiAgICogQHBhcmFtIGFsaWFzIFRoZSBBUEkgYmVpbmcgYWxpYXNlZCBjYW4gbm93IGJlIGNhbGxlZCBhdCBleHQvYWxpYXNcbiAgICpcbiAgICogQHJldHVybnMgUmV0dXJucyBhIFByb21pc2UgYm9vbGVhbiBjb250YWluaW5nIHN1Y2Nlc3MsIHRydWUgZm9yIHN1Y2Nlc3MsIGZhbHNlIGZvciBmYWlsdXJlLlxuICAgKi9cbiAgYWxpYXMgPSBhc3luYyAoZW5kcG9pbnQ6IHN0cmluZywgYWxpYXM6IHN0cmluZyk6IFByb21pc2U8Ym9vbGVhbj4gPT4ge1xuICAgIGNvbnN0IHBhcmFtczogQWxpYXNQYXJhbXMgPSB7XG4gICAgICBlbmRwb2ludCxcbiAgICAgIGFsaWFzXG4gICAgfVxuICAgIGNvbnN0IHJlc3BvbnNlOiBSZXF1ZXN0UmVzcG9uc2VEYXRhID0gYXdhaXQgdGhpcy5jYWxsTWV0aG9kKFxuICAgICAgXCJhZG1pbi5hbGlhc1wiLFxuICAgICAgcGFyYW1zXG4gICAgKVxuICAgIHJldHVybiByZXNwb25zZS5kYXRhLnJlc3VsdC5zdWNjZXNzXG4gICAgICA/IHJlc3BvbnNlLmRhdGEucmVzdWx0LnN1Y2Nlc3NcbiAgICAgIDogcmVzcG9uc2UuZGF0YS5yZXN1bHRcbiAgfVxuXG4gIC8qKlxuICAgKiBHaXZlIGEgYmxvY2tjaGFpbiBhbiBhbGlhcywgYSBkaWZmZXJlbnQgbmFtZSB0aGF0IGNhbiBiZSB1c2VkIGFueSBwbGFjZSB0aGUgYmxvY2tjaGFpbuKAmXNcbiAgICogSUQgaXMgdXNlZC5cbiAgICpcbiAgICogQHBhcmFtIGNoYWluIFRoZSBibG9ja2NoYWlu4oCZcyBJRFxuICAgKiBAcGFyYW0gYWxpYXMgQ2FuIG5vdyBiZSB1c2VkIGluIHBsYWNlIG9mIHRoZSBibG9ja2NoYWlu4oCZcyBJRCAoaW4gQVBJIGVuZHBvaW50cywgZm9yIGV4YW1wbGUpXG4gICAqXG4gICAqIEByZXR1cm5zIFJldHVybnMgYSBQcm9taXNlIGJvb2xlYW4gY29udGFpbmluZyBzdWNjZXNzLCB0cnVlIGZvciBzdWNjZXNzLCBmYWxzZSBmb3IgZmFpbHVyZS5cbiAgICovXG4gIGFsaWFzQ2hhaW4gPSBhc3luYyAoY2hhaW46IHN0cmluZywgYWxpYXM6IHN0cmluZyk6IFByb21pc2U8Ym9vbGVhbj4gPT4ge1xuICAgIGNvbnN0IHBhcmFtczogQWxpYXNDaGFpblBhcmFtcyA9IHtcbiAgICAgIGNoYWluLFxuICAgICAgYWxpYXNcbiAgICB9XG4gICAgY29uc3QgcmVzcG9uc2U6IFJlcXVlc3RSZXNwb25zZURhdGEgPSBhd2FpdCB0aGlzLmNhbGxNZXRob2QoXG4gICAgICBcImFkbWluLmFsaWFzQ2hhaW5cIixcbiAgICAgIHBhcmFtc1xuICAgIClcbiAgICByZXR1cm4gcmVzcG9uc2UuZGF0YS5yZXN1bHQuc3VjY2Vzc1xuICAgICAgPyByZXNwb25zZS5kYXRhLnJlc3VsdC5zdWNjZXNzXG4gICAgICA6IHJlc3BvbnNlLmRhdGEucmVzdWx0XG4gIH1cblxuICAvKipcbiAgICogR2V0IGFsbCBhbGlhc2VzIGZvciBnaXZlbiBibG9ja2NoYWluXG4gICAqXG4gICAqIEBwYXJhbSBjaGFpbiBUaGUgYmxvY2tjaGFpbuKAmXMgSURcbiAgICpcbiAgICogQHJldHVybnMgUmV0dXJucyBhIFByb21pc2Ugc3RyaW5nW10gY29udGFpbmluZyBhbGlhc2VzIG9mIHRoZSBibG9ja2NoYWluLlxuICAgKi9cbiAgZ2V0Q2hhaW5BbGlhc2VzID0gYXN5bmMgKGNoYWluOiBzdHJpbmcpOiBQcm9taXNlPHN0cmluZ1tdPiA9PiB7XG4gICAgY29uc3QgcGFyYW1zOiBHZXRDaGFpbkFsaWFzZXNQYXJhbXMgPSB7XG4gICAgICBjaGFpblxuICAgIH1cbiAgICBjb25zdCByZXNwb25zZTogUmVxdWVzdFJlc3BvbnNlRGF0YSA9IGF3YWl0IHRoaXMuY2FsbE1ldGhvZChcbiAgICAgIFwiYWRtaW4uZ2V0Q2hhaW5BbGlhc2VzXCIsXG4gICAgICBwYXJhbXNcbiAgICApXG4gICAgcmV0dXJuIHJlc3BvbnNlLmRhdGEucmVzdWx0LmFsaWFzZXNcbiAgICAgID8gcmVzcG9uc2UuZGF0YS5yZXN1bHQuYWxpYXNlc1xuICAgICAgOiByZXNwb25zZS5kYXRhLnJlc3VsdFxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgbG9nIGFuZCBkaXNwbGF5IGxldmVscyBvZiBsb2dnZXJzXG4gICAqXG4gICAqIEBwYXJhbSBsb2dnZXJOYW1lIHRoZSBuYW1lIG9mIHRoZSBsb2dnZXIgdG8gYmUgcmV0dXJuZWQuIFRoaXMgaXMgYW4gb3B0aW9uYWwgYXJndW1lbnQuIElmIG5vdCBzcGVjaWZpZWQsIGl0IHJldHVybnMgYWxsIHBvc3NpYmxlIGxvZ2dlcnMuXG4gICAqXG4gICAqIEByZXR1cm5zIFJldHVybnMgYSBQcm9taXNlIGNvbnRhaW5pbmcgbG9nZ2VyIGxldmVsc1xuICAgKi9cbiAgZ2V0TG9nZ2VyTGV2ZWwgPSBhc3luYyAoXG4gICAgbG9nZ2VyTmFtZT86IHN0cmluZ1xuICApOiBQcm9taXNlPEdldExvZ2dlckxldmVsUmVzcG9uc2U+ID0+IHtcbiAgICBjb25zdCBwYXJhbXM6IEdldExvZ2dlckxldmVsUGFyYW1zID0ge31cbiAgICBpZiAodHlwZW9mIGxvZ2dlck5hbWUgIT09IFwidW5kZWZpbmVkXCIpIHtcbiAgICAgIHBhcmFtcy5sb2dnZXJOYW1lID0gbG9nZ2VyTmFtZVxuICAgIH1cbiAgICBjb25zdCByZXNwb25zZTogUmVxdWVzdFJlc3BvbnNlRGF0YSA9IGF3YWl0IHRoaXMuY2FsbE1ldGhvZChcbiAgICAgIFwiYWRtaW4uZ2V0TG9nZ2VyTGV2ZWxcIixcbiAgICAgIHBhcmFtc1xuICAgIClcbiAgICByZXR1cm4gcmVzcG9uc2UuZGF0YS5yZXN1bHRcbiAgfVxuXG4gIC8qKlxuICAgKiBEeW5hbWljYWxseSBsb2FkcyBhbnkgdmlydHVhbCBtYWNoaW5lcyBpbnN0YWxsZWQgb24gdGhlIG5vZGUgYXMgcGx1Z2luc1xuICAgKlxuICAgKiBAcmV0dXJucyBSZXR1cm5zIGEgUHJvbWlzZSBjb250YWluaW5nIG5ldyBWTXMgYW5kIGZhaWxlZCBWTXNcbiAgICovXG4gIGxvYWRWTXMgPSBhc3luYyAoKTogUHJvbWlzZTxMb2FkVk1zUmVzcG9uc2U+ID0+IHtcbiAgICBjb25zdCByZXNwb25zZTogUmVxdWVzdFJlc3BvbnNlRGF0YSA9IGF3YWl0IHRoaXMuY2FsbE1ldGhvZChcImFkbWluLmxvYWRWTXNcIilcbiAgICByZXR1cm4gcmVzcG9uc2UuZGF0YS5yZXN1bHQuYWxpYXNlc1xuICAgICAgPyByZXNwb25zZS5kYXRhLnJlc3VsdC5hbGlhc2VzXG4gICAgICA6IHJlc3BvbnNlLmRhdGEucmVzdWx0XG4gIH1cblxuICAvKipcbiAgICogRHVtcCB0aGUgbXV0ZXggc3RhdGlzdGljcyBvZiB0aGUgbm9kZSB0byB0aGUgc3BlY2lmaWVkIGZpbGUuXG4gICAqXG4gICAqIEByZXR1cm5zIFByb21pc2UgZm9yIGEgYm9vbGVhbiB0aGF0IGlzIHRydWUgb24gc3VjY2Vzcy5cbiAgICovXG4gIGxvY2tQcm9maWxlID0gYXN5bmMgKCk6IFByb21pc2U8Ym9vbGVhbj4gPT4ge1xuICAgIGNvbnN0IHJlc3BvbnNlOiBSZXF1ZXN0UmVzcG9uc2VEYXRhID0gYXdhaXQgdGhpcy5jYWxsTWV0aG9kKFxuICAgICAgXCJhZG1pbi5sb2NrUHJvZmlsZVwiXG4gICAgKVxuICAgIHJldHVybiByZXNwb25zZS5kYXRhLnJlc3VsdC5zdWNjZXNzXG4gICAgICA/IHJlc3BvbnNlLmRhdGEucmVzdWx0LnN1Y2Nlc3NcbiAgICAgIDogcmVzcG9uc2UuZGF0YS5yZXN1bHRcbiAgfVxuXG4gIC8qKlxuICAgKiBEdW1wIHRoZSBjdXJyZW50IG1lbW9yeSBmb290cHJpbnQgb2YgdGhlIG5vZGUgdG8gdGhlIHNwZWNpZmllZCBmaWxlLlxuICAgKlxuICAgKiBAcmV0dXJucyBQcm9taXNlIGZvciBhIGJvb2xlYW4gdGhhdCBpcyB0cnVlIG9uIHN1Y2Nlc3MuXG4gICAqL1xuICBtZW1vcnlQcm9maWxlID0gYXN5bmMgKCk6IFByb21pc2U8Ym9vbGVhbj4gPT4ge1xuICAgIGNvbnN0IHJlc3BvbnNlOiBSZXF1ZXN0UmVzcG9uc2VEYXRhID0gYXdhaXQgdGhpcy5jYWxsTWV0aG9kKFxuICAgICAgXCJhZG1pbi5tZW1vcnlQcm9maWxlXCJcbiAgICApXG4gICAgcmV0dXJuIHJlc3BvbnNlLmRhdGEucmVzdWx0LnN1Y2Nlc3NcbiAgICAgID8gcmVzcG9uc2UuZGF0YS5yZXN1bHQuc3VjY2Vzc1xuICAgICAgOiByZXNwb25zZS5kYXRhLnJlc3VsdFxuICB9XG5cbiAgLyoqXG4gICAqIFNldHMgbG9nIGFuZCBkaXNwbGF5IGxldmVscyBvZiBsb2dnZXJzLlxuICAgKlxuICAgKiBAcGFyYW0gbG9nZ2VyTmFtZSB0aGUgbmFtZSBvZiB0aGUgbG9nZ2VyIHRvIGJlIGNoYW5nZWQuIFRoaXMgaXMgYW4gb3B0aW9uYWwgcGFyYW1ldGVyLlxuICAgKiBAcGFyYW0gbG9nTGV2ZWwgdGhlIGxvZyBsZXZlbCBvZiB3cml0dGVuIGxvZ3MsIGNhbiBiZSBvbWl0dGVkLlxuICAgKiBAcGFyYW0gZGlzcGxheUxldmVsIHRoZSBsb2cgbGV2ZWwgb2YgZGlzcGxheWVkIGxvZ3MsIGNhbiBiZSBvbWl0dGVkLlxuICAgKlxuICAgKiBAcmV0dXJucyBSZXR1cm5zIGEgUHJvbWlzZSBjb250YWluaW5nIGxvZ2dlciBsZXZlbHNcbiAgICovXG4gIHNldExvZ2dlckxldmVsID0gYXN5bmMgKFxuICAgIGxvZ2dlck5hbWU/OiBzdHJpbmcsXG4gICAgbG9nTGV2ZWw/OiBzdHJpbmcsXG4gICAgZGlzcGxheUxldmVsPzogc3RyaW5nXG4gICk6IFByb21pc2U8U2V0TG9nZ2VyTGV2ZWxSZXNwb25zZT4gPT4ge1xuICAgIGNvbnN0IHBhcmFtczogU2V0TG9nZ2VyTGV2ZWxQYXJhbXMgPSB7fVxuICAgIGlmICh0eXBlb2YgbG9nZ2VyTmFtZSAhPT0gXCJ1bmRlZmluZWRcIikge1xuICAgICAgcGFyYW1zLmxvZ2dlck5hbWUgPSBsb2dnZXJOYW1lXG4gICAgfVxuICAgIGlmICh0eXBlb2YgbG9nTGV2ZWwgIT09IFwidW5kZWZpbmVkXCIpIHtcbiAgICAgIHBhcmFtcy5sb2dMZXZlbCA9IGxvZ0xldmVsXG4gICAgfVxuICAgIGlmICh0eXBlb2YgZGlzcGxheUxldmVsICE9PSBcInVuZGVmaW5lZFwiKSB7XG4gICAgICBwYXJhbXMuZGlzcGxheUxldmVsID0gZGlzcGxheUxldmVsXG4gICAgfVxuICAgIGNvbnN0IHJlc3BvbnNlOiBSZXF1ZXN0UmVzcG9uc2VEYXRhID0gYXdhaXQgdGhpcy5jYWxsTWV0aG9kKFxuICAgICAgXCJhZG1pbi5zZXRMb2dnZXJMZXZlbFwiLFxuICAgICAgcGFyYW1zXG4gICAgKVxuICAgIHJldHVybiByZXNwb25zZS5kYXRhLnJlc3VsdFxuICB9XG5cbiAgLyoqXG4gICAqIFN0YXJ0IHByb2ZpbGluZyB0aGUgY3B1IHV0aWxpemF0aW9uIG9mIHRoZSBub2RlLiBXaWxsIGR1bXAgdGhlIHByb2ZpbGUgaW5mb3JtYXRpb24gaW50b1xuICAgKiB0aGUgc3BlY2lmaWVkIGZpbGUgb24gc3RvcC5cbiAgICpcbiAgICogQHJldHVybnMgUHJvbWlzZSBmb3IgYSBib29sZWFuIHRoYXQgaXMgdHJ1ZSBvbiBzdWNjZXNzLlxuICAgKi9cbiAgc3RhcnRDUFVQcm9maWxlciA9IGFzeW5jICgpOiBQcm9taXNlPGJvb2xlYW4+ID0+IHtcbiAgICBjb25zdCByZXNwb25zZTogUmVxdWVzdFJlc3BvbnNlRGF0YSA9IGF3YWl0IHRoaXMuY2FsbE1ldGhvZChcbiAgICAgIFwiYWRtaW4uc3RhcnRDUFVQcm9maWxlclwiXG4gICAgKVxuICAgIHJldHVybiByZXNwb25zZS5kYXRhLnJlc3VsdC5zdWNjZXNzXG4gICAgICA/IHJlc3BvbnNlLmRhdGEucmVzdWx0LnN1Y2Nlc3NcbiAgICAgIDogcmVzcG9uc2UuZGF0YS5yZXN1bHRcbiAgfVxuXG4gIC8qKlxuICAgKiBTdG9wIHRoZSBDUFUgcHJvZmlsZSB0aGF0IHdhcyBwcmV2aW91c2x5IHN0YXJ0ZWQuXG4gICAqXG4gICAqIEByZXR1cm5zIFByb21pc2UgZm9yIGEgYm9vbGVhbiB0aGF0IGlzIHRydWUgb24gc3VjY2Vzcy5cbiAgICovXG4gIHN0b3BDUFVQcm9maWxlciA9IGFzeW5jICgpOiBQcm9taXNlPGJvb2xlYW4+ID0+IHtcbiAgICBjb25zdCByZXNwb25zZTogUmVxdWVzdFJlc3BvbnNlRGF0YSA9IGF3YWl0IHRoaXMuY2FsbE1ldGhvZChcbiAgICAgIFwiYWRtaW4uc3RvcENQVVByb2ZpbGVyXCJcbiAgICApXG4gICAgcmV0dXJuIHJlc3BvbnNlLmRhdGEucmVzdWx0LnN1Y2Nlc3NcbiAgICAgID8gcmVzcG9uc2UuZGF0YS5yZXN1bHQuc3VjY2Vzc1xuICAgICAgOiByZXNwb25zZS5kYXRhLnJlc3VsdFxuICB9XG5cbiAgLyoqXG4gICAqIFRoaXMgY2xhc3Mgc2hvdWxkIG5vdCBiZSBpbnN0YW50aWF0ZWQgZGlyZWN0bHkuIEluc3RlYWQgdXNlIHRoZSBbW0F2YWxhbmNoZS5hZGRBUEldXVxuICAgKiBtZXRob2QuXG4gICAqXG4gICAqIEBwYXJhbSBjb3JlIEEgcmVmZXJlbmNlIHRvIHRoZSBBdmFsYW5jaGUgY2xhc3NcbiAgICogQHBhcmFtIGJhc2VVUkwgRGVmYXVsdHMgdG8gdGhlIHN0cmluZyBcIi9leHQvYWRtaW5cIiBhcyB0aGUgcGF0aCB0byBycGMncyBiYXNlVVJMXG4gICAqL1xuICBjb25zdHJ1Y3Rvcihjb3JlOiBBdmFsYW5jaGVDb3JlLCBiYXNlVVJMOiBzdHJpbmcgPSBcIi9leHQvYWRtaW5cIikge1xuICAgIHN1cGVyKGNvcmUsIGJhc2VVUkwpXG4gIH1cbn1cbiJdfQ==