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
exports.IndexAPI = void 0;
const jrpcapi_1 = require("../../common/jrpcapi");
/**
 * Class for interacting with a node's IndexAPI.
 *
 * @category RPCAPIs
 *
 * @remarks This extends the [[JRPCAPI]] class. This class should not be directly called. Instead, use the [[Avalanche.addAPI]] function to register this interface with Avalanche.
 */
class IndexAPI extends jrpcapi_1.JRPCAPI {
    /**
     * This class should not be instantiated directly. Instead use the [[Avalanche.addAPI]] method.
     *
     * @param core A reference to the Avalanche class
     * @param baseURL Defaults to the string "/ext/index/X/tx" as the path to rpc's baseURL
     */
    constructor(core, baseURL = "/ext/index/X/tx") {
        super(core, baseURL);
        /**
         * Get last accepted tx, vtx or block
         *
         * @param encoding
         * @param baseURL
         *
         * @returns Returns a Promise GetLastAcceptedResponse.
         */
        this.getLastAccepted = (encoding = "hex", baseURL = this.getBaseURL()) => __awaiter(this, void 0, void 0, function* () {
            this.setBaseURL(baseURL);
            const params = {
                encoding
            };
            try {
                const response = yield this.callMethod("index.getLastAccepted", params);
                return response.data.result;
            }
            catch (error) {
                console.log(error);
            }
        });
        /**
         * Get container by index
         *
         * @param index
         * @param encoding
         * @param baseURL
         *
         * @returns Returns a Promise GetContainerByIndexResponse.
         */
        this.getContainerByIndex = (index = "0", encoding = "hex", baseURL = this.getBaseURL()) => __awaiter(this, void 0, void 0, function* () {
            this.setBaseURL(baseURL);
            const params = {
                index,
                encoding
            };
            try {
                const response = yield this.callMethod("index.getContainerByIndex", params);
                return response.data.result;
            }
            catch (error) {
                console.log(error);
            }
        });
        /**
         * Get contrainer by ID
         *
         * @param containerID
         * @param encoding
         * @param baseURL
         *
         * @returns Returns a Promise GetContainerByIDResponse.
         */
        this.getContainerByID = (containerID = "0", encoding = "hex", baseURL = this.getBaseURL()) => __awaiter(this, void 0, void 0, function* () {
            this.setBaseURL(baseURL);
            const params = {
                containerID,
                encoding
            };
            try {
                const response = yield this.callMethod("index.getContainerByID", params);
                return response.data.result;
            }
            catch (error) {
                console.log(error);
            }
        });
        /**
         * Get container range
         *
         * @param startIndex
         * @param numToFetch
         * @param encoding
         * @param baseURL
         *
         * @returns Returns a Promise GetContainerRangeResponse.
         */
        this.getContainerRange = (startIndex = 0, numToFetch = 100, encoding = "hex", baseURL = this.getBaseURL()) => __awaiter(this, void 0, void 0, function* () {
            this.setBaseURL(baseURL);
            const params = {
                startIndex,
                numToFetch,
                encoding
            };
            try {
                const response = yield this.callMethod("index.getContainerRange", params);
                return response.data.result;
            }
            catch (error) {
                console.log(error);
            }
        });
        /**
         * Get index by containerID
         *
         * @param containerID
         * @param encoding
         * @param baseURL
         *
         * @returns Returns a Promise GetIndexResponse.
         */
        this.getIndex = (containerID = "", encoding = "hex", baseURL = this.getBaseURL()) => __awaiter(this, void 0, void 0, function* () {
            this.setBaseURL(baseURL);
            const params = {
                containerID,
                encoding
            };
            try {
                const response = yield this.callMethod("index.getIndex", params);
                return response.data.result.index;
            }
            catch (error) {
                console.log(error);
            }
        });
        /**
         * Check if container is accepted
         *
         * @param containerID
         * @param encoding
         * @param baseURL
         *
         * @returns Returns a Promise GetIsAcceptedResponse.
         */
        this.isAccepted = (containerID = "", encoding = "hex", baseURL = this.getBaseURL()) => __awaiter(this, void 0, void 0, function* () {
            this.setBaseURL(baseURL);
            const params = {
                containerID,
                encoding
            };
            try {
                const response = yield this.callMethod("index.isAccepted", params);
                return response.data.result;
            }
            catch (error) {
                console.log(error);
            }
        });
    }
}
exports.IndexAPI = IndexAPI;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYXBpLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2FwaXMvaW5kZXgvYXBpLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7Ozs7Ozs7OztBQUtBLGtEQUE4QztBQWdCOUM7Ozs7OztHQU1HO0FBQ0gsTUFBYSxRQUFTLFNBQVEsaUJBQU87SUEyTG5DOzs7OztPQUtHO0lBQ0gsWUFBWSxJQUFtQixFQUFFLFVBQWtCLGlCQUFpQjtRQUNsRSxLQUFLLENBQUMsSUFBSSxFQUFFLE9BQU8sQ0FBQyxDQUFBO1FBak10Qjs7Ozs7OztXQU9HO1FBQ0gsb0JBQWUsR0FBRyxDQUNoQixXQUFtQixLQUFLLEVBQ3hCLFVBQWtCLElBQUksQ0FBQyxVQUFVLEVBQUUsRUFDRCxFQUFFO1lBQ3BDLElBQUksQ0FBQyxVQUFVLENBQUMsT0FBTyxDQUFDLENBQUE7WUFDeEIsTUFBTSxNQUFNLEdBQTBCO2dCQUNwQyxRQUFRO2FBQ1QsQ0FBQTtZQUVELElBQUk7Z0JBQ0YsTUFBTSxRQUFRLEdBQXdCLE1BQU0sSUFBSSxDQUFDLFVBQVUsQ0FDekQsdUJBQXVCLEVBQ3ZCLE1BQU0sQ0FDUCxDQUFBO2dCQUNELE9BQU8sUUFBUSxDQUFDLElBQUksQ0FBQyxNQUFNLENBQUE7YUFDNUI7WUFBQyxPQUFPLEtBQUssRUFBRTtnQkFDZCxPQUFPLENBQUMsR0FBRyxDQUFDLEtBQUssQ0FBQyxDQUFBO2FBQ25CO1FBQ0gsQ0FBQyxDQUFBLENBQUE7UUFFRDs7Ozs7Ozs7V0FRRztRQUNILHdCQUFtQixHQUFHLENBQ3BCLFFBQWdCLEdBQUcsRUFDbkIsV0FBbUIsS0FBSyxFQUN4QixVQUFrQixJQUFJLENBQUMsVUFBVSxFQUFFLEVBQ0csRUFBRTtZQUN4QyxJQUFJLENBQUMsVUFBVSxDQUFDLE9BQU8sQ0FBQyxDQUFBO1lBQ3hCLE1BQU0sTUFBTSxHQUE4QjtnQkFDeEMsS0FBSztnQkFDTCxRQUFRO2FBQ1QsQ0FBQTtZQUVELElBQUk7Z0JBQ0YsTUFBTSxR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QTtnQkFDRCxPQUFPLFFBQVEsQ0FBQyxJQUFJLENBQUMsTUFBTSxDQUFBO2FBQzVCO1lBQUMsT0FBTyxLQUFLLEVBQUU7Z0JBQ2QsT0FBTyxDQUFDLEdBQUcsQ0FBQyxLQUFLLENBQUMsQ0FBQTthQUNuQjtRQUNILENBQUMsQ0FBQSxDQUFBO0lBVUQsQ0FBQztDQUNGO0FBcE1ELDRCQW9NQyIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICogQHBhY2thZ2VEb2N1bWVudGF0aW9uXG4gKiBAbW9kdWxlIEFQSS1JbmRleFxuICovXG5pbXBvcnQgQXZhbGFuY2hlQ29yZSBmcm9tIFwiLi4vLi4vYXZhbGFuY2hlXCJcbmltcG9ydCB7IEpSUENBUEkgfSBmcm9tIFwiLi4vLi4vY29tbW9uL2pycGNhcGlcIlxuaW1wb3J0IHsgUmVxdWVzdFJlc3BvbnNlRGF0YSB9IGZyb20gXCIuLi8uLi9jb21tb24vYXBpYmFzZVwiXG5pbXBvcnQge1xuICBHZXRMYXN0QWNjZXB0ZWRQYXJhbXMsXG4gIEdldExhc3RBY2NlcHRlZFJlc3BvbnNlLFxuICBHZXRDb250YWluZXJCeUluZGV4UGFyYW1zLFxuICBHZXRDb250YWluZXJCeUluZGV4UmVzcG9uc2UsXG4gIEdldENvbnRhaW5lckJ5SURQYXJhbXMsXG4gIEdldENvbnRhaW5lckJ5SURSZXNwb25zZSxcbiAgR2V0Q29udGFpbmVyUmFuZ2VQYXJhbXMsXG4gIEdldENvbnRhaW5lclJhbmdlUmVzcG9uc2UsXG4gIEdldEluZGV4UGFyYW1zLFxuICBHZXRJc0FjY2VwdGVkUGFyYW1zLFxuICBJc0FjY2VwdGVkUmVzcG9uc2Vcbn0gZnJvbSBcIi4vaW50ZXJmYWNlc1wiXG5cbi8qKlxuICogQ2xhc3MgZm9yIGludGVyYWN0aW5nIHdpdGggYSBub2RlJ3MgSW5kZXhBUEkuXG4gKlxuICogQGNhdGVnb3J5IFJQQ0FQSXNcbiAqXG4gKiBAcmVtYXJrcyBUaGlzIGV4dGVuZHMgdGhlIFtbSlJQQ0FQSV1dIGNsYXNzLiBUaGlzIGNsYXNzIHNob3VsZCBub3QgYmUgZGlyZWN0bHkgY2FsbGVkLiBJbnN0ZWFkLCB1c2UgdGhlIFtbQXZhbGFuY2hlLmFkZEFQSV1dIGZ1bmN0aW9uIHRvIHJlZ2lzdGVyIHRoaXMgaW50ZXJmYWNlIHdpdGggQXZhbGFuY2hlLlxuICovXG5leHBvcnQgY2xhc3MgSW5kZXhBUEkgZXh0ZW5kcyBKUlBDQVBJIHtcbiAgLyoqXG4gICAqIEdldCBsYXN0IGFjY2VwdGVkIHR4LCB2dHggb3IgYmxvY2tcbiAgICpcbiAgICogQHBhcmFtIGVuY29kaW5nXG4gICAqIEBwYXJhbSBiYXNlVVJMXG4gICAqXG4gICAqIEByZXR1cm5zIFJldHVybnMgYSBQcm9taXNlIEdldExhc3RBY2NlcHRlZFJlc3BvbnNlLlxuICAgKi9cbiAgZ2V0TGFzdEFjY2VwdGVkID0gYXN5bmMgKFxuICAgIGVuY29kaW5nOiBzdHJpbmcgPSBcImhleFwiLFxuICAgIGJhc2VVUkw6IHN0cmluZyA9IHRoaXMuZ2V0QmFzZVVSTCgpXG4gICk6IFByb21pc2U8R2V0TGFzdEFjY2VwdGVkUmVzcG9uc2U+ID0+IHtcbiAgICB0aGlzLnNldEJhc2VVUkwoYmFzZVVSTClcbiAgICBjb25zdCBwYXJhbXM6IEdldExhc3RBY2NlcHRlZFBhcmFtcyA9IHtcbiAgICAgIGVuY29kaW5nXG4gICAgfVxuXG4gICAgdHJ5IHtcbiAgICAgIGNvbnN0IHJlc3BvbnNlOiBSZXF1ZXN0UmVzcG9uc2VEYXRhID0gYXdhaXQgdGhpcy5jYWxsTWV0aG9kKFxuICAgICAgICBcImluZGV4LmdldExhc3RBY2NlcHRlZFwiLFxuICAgICAgICBwYXJhbXNcbiAgICAgIClcbiAgICAgIHJldHVybiByZXNwb25zZS5kYXRhLnJlc3VsdFxuICAgIH0gY2F0Y2ggKGVycm9yKSB7XG4gICAgICBjb25zb2xlLmxvZyhlcnJvcilcbiAgICB9XG4gIH1cblxuICAvKipcbiAgICogR2V0IGNvbnRhaW5lciBieSBpbmRleFxuICAgKlxuICAgKiBAcGFyYW0gaW5kZXhcbiAgICogQHBhcmFtIGVuY29kaW5nXG4gICAqIEBwYXJhbSBiYXNlVVJMXG4gICAqXG4gICAqIEByZXR1cm5zIFJldHVybnMgYSBQcm9taXNlIEdldENvbnRhaW5lckJ5SW5kZXhSZXNwb25zZS5cbiAgICovXG4gIGdldENvbnRhaW5lckJ5SW5kZXggPSBhc3luYyAoXG4gICAgaW5kZXg6IHN0cmluZyA9IFwiMFwiLFxuICAgIGVuY29kaW5nOiBzdHJpbmcgPSBcImhleFwiLFxuICAgIGJhc2VVUkw6IHN0cmluZyA9IHRoaXMuZ2V0QmFzZVVSTCgpXG4gICk6IFByb21pc2U8R2V0Q29udGFpbmVyQnlJbmRleFJlc3BvbnNlPiA9PiB7XG4gICAgdGhpcy5zZXRCYXNlVVJMKGJhc2VVUkwpXG4gICAgY29uc3QgcGFyYW1zOiBHZXRDb250YWluZXJCeUluZGV4UGFyYW1zID0ge1xuICAgICAgaW5kZXgsXG4gICAgICBlbmNvZGluZ1xuICAgIH1cblxuICAgIHRyeSB7XG4gICAgICBjb25zdCByZXNwb25zZTogUmVxdWVzdFJlc3BvbnNlRGF0YSA9IGF3YWl0IHRoaXMuY2FsbE1ldGhvZChcbiAgICAgICAgXCJpbmRleC5nZXRDb250YWluZXJCeUluZGV4XCIsXG4gICAgICAgIHBhcmFtc1xuICAgICAgKVxuICAgICAgcmV0dXJuIHJlc3BvbnNlLmRhdGEucmVzdWx0XG4gICAgfSBjYXRjaCAoZXJyb3IpIHtcbiAgICAgIGNvbnNvbGUubG9nKGVycm9yKVxuICAgIH1cbiAgfVxuXG4gIC8qKlxuICAgKiBHZXQgY29udHJhaW5lciBieSBJRFxuICAgKlxuICAgKiBAcGFyYW0gY29udGFpbmVySURcbiAgICogQHBhcmFtIGVuY29kaW5nXG4gICAqIEBwYXJhbSBiYXNlVVJMXG4gICAqXG4gICAqIEByZXR1cm5zIFJldHVybnMgYSBQcm9taXNlIEdldENvbnRhaW5lckJ5SURSZXNwb25zZS5cbiAgICovXG4gIGdldENvbnRhaW5lckJ5SUQgPSBhc3luYyAoXG4gICAgY29udGFpbmVySUQ6IHN0cmluZyA9IFwiMFwiLFxuICAgIGVuY29kaW5nOiBzdHJpbmcgPSBcImhleFwiLFxuICAgIGJhc2VVUkw6IHN0cmluZyA9IHRoaXMuZ2V0QmFzZVVSTCgpXG4gICk6IFByb21pc2U8R2V0Q29udGFpbmVyQnlJRFJlc3BvbnNlPiA9PiB7XG4gICAgdGhpcy5zZXRCYXNlVVJMKGJhc2VVUkwpXG4gICAgY29uc3QgcGFyYW1zOiBHZXRDb250YWluZXJCeUlEUGFyYW1zID0ge1xuICAgICAgY29udGFpbmVySUQsXG4gICAgICBlbmNvZGluZ1xuICAgIH1cblxuICAgIHRyeSB7XG4gICAgICBjb25zdCByZXNwb25zZTogUmVxdWVzdFJlc3BvbnNlRGF0YSA9IGF3YWl0IHRoaXMuY2FsbE1ldGhvZChcbiAgICAgICAgXCJpbmRleC5nZXRDb250YWluZXJCeUlEXCIsXG4gICAgICAgIHBhcmFtc1xuICAgICAgKVxuICAgICAgcmV0dXJuIHJlc3BvbnNlLmRhdGEucmVzdWx0XG4gICAgfSBjYXRjaCAoZXJyb3IpIHtcbiAgICAgIGNvbnNvbGUubG9nKGVycm9yKVxuICAgIH1cbiAgfVxuXG4gIC8qKlxuICAgKiBHZXQgY29udGFpbmVyIHJhbmdlXG4gICAqXG4gICAqIEBwYXJhbSBzdGFydEluZGV4XG4gICAqIEBwYXJhbSBudW1Ub0ZldGNoXG4gICAqIEBwYXJhbSBlbmNvZGluZ1xuICAgKiBAcGFyYW0gYmFzZVVSTFxuICAgKlxuICAgKiBAcmV0dXJucyBSZXR1cm5zIGEgUHJvbWlzZSBHZXRDb250YWluZXJSYW5nZVJlc3BvbnNlLlxuICAgKi9cbiAgZ2V0Q29udGFpbmVyUmFuZ2UgPSBhc3luYyAoXG4gICAgc3RhcnRJbmRleDogbnVtYmVyID0gMCxcbiAgICBudW1Ub0ZldGNoOiBudW1iZXIgPSAxMDAsXG4gICAgZW5jb2Rpbmc6IHN0cmluZyA9IFwiaGV4XCIsXG4gICAgYmFzZVVSTDogc3RyaW5nID0gdGhpcy5nZXRCYXNlVVJMKClcbiAgKTogUHJvbWlzZTxHZXRDb250YWluZXJSYW5nZVJlc3BvbnNlW10+ID0+IHtcbiAgICB0aGlzLnNldEJhc2VVUkwoYmFzZVVSTClcbiAgICBjb25zdCBwYXJhbXM6IEdldENvbnRhaW5lclJhbmdlUGFyYW1zID0ge1xuICAgICAgc3RhcnRJbmRleCxcbiAgICAgIG51bVRvRmV0Y2gsXG4gICAgICBlbmNvZGluZ1xuICAgIH1cblxuICAgIHRyeSB7XG4gICAgICBjb25zdCByZXNwb25zZTogUmVxdWVzdFJlc3BvbnNlRGF0YSA9IGF3YWl0IHRoaXMuY2FsbE1ldGhvZChcbiAgICAgICAgXCJpbmRleC5nZXRDb250YWluZXJSYW5nZVwiLFxuICAgICAgICBwYXJhbXNcbiAgICAgIClcbiAgICAgIHJldHVybiByZXNwb25zZS5kYXRhLnJlc3VsdFxuICAgIH0gY2F0Y2ggKGVycm9yKSB7XG4gICAgICBjb25zb2xlLmxvZyhlcnJvcilcbiAgICB9XG4gIH1cblxuICAvKipcbiAgICogR2V0IGluZGV4IGJ5IGNvbnRhaW5lcklEXG4gICAqXG4gICAqIEBwYXJhbSBjb250YWluZXJJRFxuICAgKiBAcGFyYW0gZW5jb2RpbmdcbiAgICogQHBhcmFtIGJhc2VVUkxcbiAgICpcbiAgICogQHJldHVybnMgUmV0dXJucyBhIFByb21pc2UgR2V0SW5kZXhSZXNwb25zZS5cbiAgICovXG4gIGdldEluZGV4ID0gYXN5bmMgKFxuICAgIGNvbnRhaW5lcklEOiBzdHJpbmcgPSBcIlwiLFxuICAgIGVuY29kaW5nOiBzdHJpbmcgPSBcImhleFwiLFxuICAgIGJhc2VVUkw6IHN0cmluZyA9IHRoaXMuZ2V0QmFzZVVSTCgpXG4gICk6IFByb21pc2U8c3RyaW5nPiA9PiB7XG4gICAgdGhpcy5zZXRCYXNlVVJMKGJhc2VVUkwpXG4gICAgY29uc3QgcGFyYW1zOiBHZXRJbmRleFBhcmFtcyA9IHtcbiAgICAgIGNvbnRhaW5lcklELFxuICAgICAgZW5jb2RpbmdcbiAgICB9XG5cbiAgICB0cnkge1xuICAgICAgY29uc3QgcmVzcG9uc2U6IFJlcXVlc3RSZXNwb25zZURhdGEgPSBhd2FpdCB0aGlzLmNhbGxNZXRob2QoXG4gICAgICAgIFwiaW5kZXguZ2V0SW5kZXhcIixcbiAgICAgICAgcGFyYW1zXG4gICAgICApXG4gICAgICByZXR1cm4gcmVzcG9uc2UuZGF0YS5yZXN1bHQuaW5kZXhcbiAgICB9IGNhdGNoIChlcnJvcikge1xuICAgICAgY29uc29sZS5sb2coZXJyb3IpXG4gICAgfVxuICB9XG5cbiAgLyoqXG4gICAqIENoZWNrIGlmIGNvbnRhaW5lciBpcyBhY2NlcHRlZFxuICAgKlxuICAgKiBAcGFyYW0gY29udGFpbmVySURcbiAgICogQHBhcmFtIGVuY29kaW5nXG4gICAqIEBwYXJhbSBiYXNlVVJMXG4gICAqXG4gICAqIEByZXR1cm5zIFJldHVybnMgYSBQcm9taXNlIEdldElzQWNjZXB0ZWRSZXNwb25zZS5cbiAgICovXG4gIGlzQWNjZXB0ZWQgPSBhc3luYyAoXG4gICAgY29udGFpbmVySUQ6IHN0cmluZyA9IFwiXCIsXG4gICAgZW5jb2Rpbmc6IHN0cmluZyA9IFwiaGV4XCIsXG4gICAgYmFzZVVSTDogc3RyaW5nID0gdGhpcy5nZXRCYXNlVVJMKClcbiAgKTogUHJvbWlzZTxJc0FjY2VwdGVkUmVzcG9uc2U+ID0+IHtcbiAgICB0aGlzLnNldEJhc2VVUkwoYmFzZVVSTClcbiAgICBjb25zdCBwYXJhbXM6IEdldElzQWNjZXB0ZWRQYXJhbXMgPSB7XG4gICAgICBjb250YWluZXJJRCxcbiAgICAgIGVuY29kaW5nXG4gICAgfVxuXG4gICAgdHJ5IHtcbiAgICAgIGNvbnN0IHJlc3BvbnNlOiBSZXF1ZXN0UmVzcG9uc2VEYXRhID0gYXdhaXQgdGhpcy5jYWxsTWV0aG9kKFxuICAgICAgICBcImluZGV4LmlzQWNjZXB0ZWRcIixcbiAgICAgICAgcGFyYW1zXG4gICAgICApXG4gICAgICByZXR1cm4gcmVzcG9uc2UuZGF0YS5yZXN1bHRcbiAgICB9IGNhdGNoIChlcnJvcikge1xuICAgICAgY29uc29sZS5sb2coZXJyb3IpXG4gICAgfVxuICB9XG5cbiAgLyoqXG4gICAqIFRoaXMgY2xhc3Mgc2hvdWxkIG5vdCBiZSBpbnN0YW50aWF0ZWQgZGlyZWN0bHkuIEluc3RlYWQgdXNlIHRoZSBbW0F2YWxhbmNoZS5hZGRBUEldXSBtZXRob2QuXG4gICAqXG4gICAqIEBwYXJhbSBjb3JlIEEgcmVmZXJlbmNlIHRvIHRoZSBBdmFsYW5jaGUgY2xhc3NcbiAgICogQHBhcmFtIGJhc2VVUkwgRGVmYXVsdHMgdG8gdGhlIHN0cmluZyBcIi9leHQvaW5kZXgvWC90eFwiIGFzIHRoZSBwYXRoIHRvIHJwYydzIGJhc2VVUkxcbiAgICovXG4gIGNvbnN0cnVjdG9yKGNvcmU6IEF2YWxhbmNoZUNvcmUsIGJhc2VVUkw6IHN0cmluZyA9IFwiL2V4dC9pbmRleC9YL3R4XCIpIHtcbiAgICBzdXBlcihjb3JlLCBiYXNlVVJMKVxuICB9XG59XG4iXX0=