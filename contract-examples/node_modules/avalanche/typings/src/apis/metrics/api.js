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
exports.MetricsAPI = void 0;
const restapi_1 = require("../../common/restapi");
/**
 * Class for interacting with a node API that is using the node's MetricsApi.
 *
 * @category RPCAPIs
 *
 * @remarks This extends the [[RESTAPI]] class. This class should not be directly called. Instead, use the [[Avalanche.addAPI]] function to register this interface with Avalanche.
 */
class MetricsAPI extends restapi_1.RESTAPI {
    /**
     * This class should not be instantiated directly. Instead use the [[Avalanche.addAPI]] method.
     *
     * @param core A reference to the Avalanche class
     * @param baseURL Defaults to the string "/ext/metrics" as the path to rpc's baseurl
     */
    constructor(core, baseURL = "/ext/metrics") {
        super(core, baseURL);
        this.axConf = () => {
            return {
                baseURL: `${this.core.getProtocol()}://${this.core.getHost()}:${this.core.getPort()}`,
                responseType: "text"
            };
        };
        /**
         *
         * @returns Promise for an object containing the metrics response
         */
        this.getMetrics = () => __awaiter(this, void 0, void 0, function* () {
            const response = yield this.post("");
            return response.data;
        });
    }
}
exports.MetricsAPI = MetricsAPI;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYXBpLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2FwaXMvbWV0cmljcy9hcGkudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7O0FBS0Esa0RBQThDO0FBSTlDOzs7Ozs7R0FNRztBQUNILE1BQWEsVUFBVyxTQUFRLGlCQUFPO0lBaUJyQzs7Ozs7T0FLRztJQUNILFlBQVksSUFBbUIsRUFBRSxVQUFrQixjQUFjO1FBQy9ELEtBQUssQ0FBQyxJQUFJLEVBQUUsT0FBTyxDQUFDLENBQUE7UUF2QlosV0FBTSxHQUFHLEdBQXVCLEVBQUU7WUFDMUMsT0FBTztnQkFDTCxPQUFPLEVBQUUsR0FBRyxJQUFJLENBQUMsSUFBSSxDQUFDLFdBQVcsRUFBRSxNQUFNLElBQUksQ0FBQyxJQUFJLENBQUMsT0FBTyxFQUFFLElBQUksSUFBSSxDQUFDLElBQUksQ0FBQyxPQUFPLEVBQUUsRUFBRTtnQkFDckYsWUFBWSxFQUFFLE1BQU07YUFDckIsQ0FBQTtRQUNILENBQUMsQ0FBQTtRQUVEOzs7V0FHRztRQUNILGVBQVUsR0FBRyxHQUEwQixFQUFFO1lBQ3ZDLE1BQU0sUUFBUSxHQUF3QixNQUFNLElBQUksQ0FBQyxJQUFJLENBQUMsRUFBRSxDQUFDLENBQUE7WUFDekQsT0FBTyxRQUFRLENBQUMsSUFBYyxDQUFBO1FBQ2hDLENBQUMsQ0FBQSxDQUFBO0lBVUQsQ0FBQztDQUNGO0FBMUJELGdDQTBCQyIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICogQHBhY2thZ2VEb2N1bWVudGF0aW9uXG4gKiBAbW9kdWxlIEFQSS1NZXRyaWNzXG4gKi9cbmltcG9ydCBBdmFsYW5jaGVDb3JlIGZyb20gXCIuLi8uLi9hdmFsYW5jaGVcIlxuaW1wb3J0IHsgUkVTVEFQSSB9IGZyb20gXCIuLi8uLi9jb21tb24vcmVzdGFwaVwiXG5pbXBvcnQgeyBSZXF1ZXN0UmVzcG9uc2VEYXRhIH0gZnJvbSBcIi4uLy4uL2NvbW1vbi9hcGliYXNlXCJcbmltcG9ydCB7IEF4aW9zUmVxdWVzdENvbmZpZyB9IGZyb20gXCJheGlvc1wiXG5cbi8qKlxuICogQ2xhc3MgZm9yIGludGVyYWN0aW5nIHdpdGggYSBub2RlIEFQSSB0aGF0IGlzIHVzaW5nIHRoZSBub2RlJ3MgTWV0cmljc0FwaS5cbiAqXG4gKiBAY2F0ZWdvcnkgUlBDQVBJc1xuICpcbiAqIEByZW1hcmtzIFRoaXMgZXh0ZW5kcyB0aGUgW1tSRVNUQVBJXV0gY2xhc3MuIFRoaXMgY2xhc3Mgc2hvdWxkIG5vdCBiZSBkaXJlY3RseSBjYWxsZWQuIEluc3RlYWQsIHVzZSB0aGUgW1tBdmFsYW5jaGUuYWRkQVBJXV0gZnVuY3Rpb24gdG8gcmVnaXN0ZXIgdGhpcyBpbnRlcmZhY2Ugd2l0aCBBdmFsYW5jaGUuXG4gKi9cbmV4cG9ydCBjbGFzcyBNZXRyaWNzQVBJIGV4dGVuZHMgUkVTVEFQSSB7XG4gIHByb3RlY3RlZCBheENvbmYgPSAoKTogQXhpb3NSZXF1ZXN0Q29uZmlnID0+IHtcbiAgICByZXR1cm4ge1xuICAgICAgYmFzZVVSTDogYCR7dGhpcy5jb3JlLmdldFByb3RvY29sKCl9Oi8vJHt0aGlzLmNvcmUuZ2V0SG9zdCgpfToke3RoaXMuY29yZS5nZXRQb3J0KCl9YCxcbiAgICAgIHJlc3BvbnNlVHlwZTogXCJ0ZXh0XCJcbiAgICB9XG4gIH1cblxuICAvKipcbiAgICpcbiAgICogQHJldHVybnMgUHJvbWlzZSBmb3IgYW4gb2JqZWN0IGNvbnRhaW5pbmcgdGhlIG1ldHJpY3MgcmVzcG9uc2VcbiAgICovXG4gIGdldE1ldHJpY3MgPSBhc3luYyAoKTogUHJvbWlzZTxzdHJpbmc+ID0+IHtcbiAgICBjb25zdCByZXNwb25zZTogUmVxdWVzdFJlc3BvbnNlRGF0YSA9IGF3YWl0IHRoaXMucG9zdChcIlwiKVxuICAgIHJldHVybiByZXNwb25zZS5kYXRhIGFzIHN0cmluZ1xuICB9XG5cbiAgLyoqXG4gICAqIFRoaXMgY2xhc3Mgc2hvdWxkIG5vdCBiZSBpbnN0YW50aWF0ZWQgZGlyZWN0bHkuIEluc3RlYWQgdXNlIHRoZSBbW0F2YWxhbmNoZS5hZGRBUEldXSBtZXRob2QuXG4gICAqXG4gICAqIEBwYXJhbSBjb3JlIEEgcmVmZXJlbmNlIHRvIHRoZSBBdmFsYW5jaGUgY2xhc3NcbiAgICogQHBhcmFtIGJhc2VVUkwgRGVmYXVsdHMgdG8gdGhlIHN0cmluZyBcIi9leHQvbWV0cmljc1wiIGFzIHRoZSBwYXRoIHRvIHJwYydzIGJhc2V1cmxcbiAgICovXG4gIGNvbnN0cnVjdG9yKGNvcmU6IEF2YWxhbmNoZUNvcmUsIGJhc2VVUkw6IHN0cmluZyA9IFwiL2V4dC9tZXRyaWNzXCIpIHtcbiAgICBzdXBlcihjb3JlLCBiYXNlVVJMKVxuICB9XG59XG4iXX0=