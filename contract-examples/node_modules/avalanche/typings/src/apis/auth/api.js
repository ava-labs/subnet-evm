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
exports.AuthAPI = void 0;
const jrpcapi_1 = require("../../common/jrpcapi");
/**
 * Class for interacting with a node's AuthAPI.
 *
 * @category RPCAPIs
 *
 * @remarks This extends the [[JRPCAPI]] class. This class should not be directly called. Instead, use the [[Avalanche.addAPI]] function to register this interface with Avalanche.
 */
class AuthAPI extends jrpcapi_1.JRPCAPI {
    /**
     * This class should not be instantiated directly. Instead use the [[Avalanche.addAPI]]
     * method.
     *
     * @param core A reference to the Avalanche class
     * @param baseURL Defaults to the string "/ext/auth" as the path to rpc's baseURL
     */
    constructor(core, baseURL = "/ext/auth") {
        super(core, baseURL);
        /**
         * Creates a new authorization token that grants access to one or more API endpoints.
         *
         * @param password This node's authorization token password, set through the CLI when the node was launched.
         * @param endpoints A list of endpoints that will be accessible using the generated token. If there"s an element that is "*", this token can reach any endpoint.
         *
         * @returns Returns a Promise string containing the authorization token.
         */
        this.newToken = (password, endpoints) => __awaiter(this, void 0, void 0, function* () {
            const params = {
                password,
                endpoints
            };
            const response = yield this.callMethod("auth.newToken", params);
            return response.data.result.token
                ? response.data.result.token
                : response.data.result;
        });
        /**
         * Revokes an authorization token, removing all of its rights to access endpoints.
         *
         * @param password This node's authorization token password, set through the CLI when the node was launched.
         * @param token An authorization token whose access should be revoked.
         *
         * @returns Returns a Promise boolean indicating if a token was successfully revoked.
         */
        this.revokeToken = (password, token) => __awaiter(this, void 0, void 0, function* () {
            const params = {
                password,
                token
            };
            const response = yield this.callMethod("auth.revokeToken", params);
            return response.data.result.success;
        });
        /**
         * Change this node's authorization token password. **Any authorization tokens created under an old password will become invalid.**
         *
         * @param oldPassword This node's authorization token password, set through the CLI when the node was launched.
         * @param newPassword A new password for this node's authorization token issuance.
         *
         * @returns Returns a Promise boolean indicating if the password was successfully changed.
         */
        this.changePassword = (oldPassword, newPassword) => __awaiter(this, void 0, void 0, function* () {
            const params = {
                oldPassword,
                newPassword
            };
            const response = yield this.callMethod("auth.changePassword", params);
            return response.data.result.success;
        });
    }
}
exports.AuthAPI = AuthAPI;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYXBpLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2FwaXMvYXV0aC9hcGkudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7O0FBS0Esa0RBQThDO0FBUzlDOzs7Ozs7R0FNRztBQUNILE1BQWEsT0FBUSxTQUFRLGlCQUFPO0lBcUVsQzs7Ozs7O09BTUc7SUFDSCxZQUFZLElBQW1CLEVBQUUsVUFBa0IsV0FBVztRQUM1RCxLQUFLLENBQUMsSUFBSSxFQUFFLE9BQU8sQ0FBQyxDQUFBO1FBNUV0Qjs7Ozs7OztXQU9HO1FBQ0gsYUFBUSxHQUFHLENBQ1QsUUFBZ0IsRUFDaEIsU0FBbUIsRUFDb0IsRUFBRTtZQUN6QyxNQUFNLE1BQU0sR0FBc0I7Z0JBQ2hDLFFBQVE7Z0JBQ1IsU0FBUzthQUNWLENBQUE7WUFDRCxNQUFNLFFBQVEsR0FBd0IsTUFBTSxJQUFJLENBQUMsVUFBVSxDQUN6RCxlQUFlLEVBQ2YsTUFBTSxDQUNQLENBQUE7WUFDRCxPQUFPLFFBQVEsQ0FBQyxJQUFJLENBQUMsTUFBTSxDQUFDLEtBQUs7Z0JBQy9CLENBQUMsQ0FBQyxRQUFRLENBQUMsSUFBSSxDQUFDLE1BQU0sQ0FBQyxLQUFLO2dCQUM1QixDQUFDLENBQUMsUUFBUSxDQUFDLElBQUksQ0FBQyxNQUFNLENBQUE7UUFDMUIsQ0FBQyxDQUFBLENBQUE7UUFFRDs7Ozs7OztXQU9HO1FBQ0gsZ0JBQVcsR0FBRyxDQUFPLFFBQWdCLEVBQUUsS0FBYSxFQUFvQixFQUFFO1lBQ3hFLE1BQU0sTUFBTSxHQUF5QjtnQkFDbkMsUUFBUTtnQkFDUixLQUFLO2FBQ04sQ0FBQTtZQUNELE1BQU0sUUFBUSxHQUF3QixNQUFNLElBQUksQ0FBQyxVQUFVLENBQ3pELGtCQUFrQixFQUNsQixNQUFNLENBQ1AsQ0FBQTtZQUNELE9BQU8sUUFBUSxDQUFDLElBQUksQ0FBQyxNQUFNLENBQUMsT0FBTyxDQUFBO1FBQ3JDLENBQUMsQ0FBQSxDQUFBO1FBRUQ7Ozs7Ozs7V0FPRztRQUNILG1CQUFjLEdBQUcsQ0FDZixXQUFtQixFQUNuQixXQUFtQixFQUNELEVBQUU7WUFDcEIsTUFBTSxNQUFNLEdBQTRCO2dCQUN0QyxXQUFXO2dCQUNYLFdBQVc7YUFDWixDQUFBO1lBQ0QsTUFBTSxRQUFRLEdBQXdCLE1BQU0sSUFBSSxDQUFDLFVBQVUsQ0FDekQscUJBQXFCLEVBQ3JCLE1BQU0sQ0FDUCxDQUFBO1lBQ0QsT0FBTyxRQUFRLENBQUMsSUFBSSxDQUFDLE1BQU0sQ0FBQyxPQUFPLENBQUE7UUFDckMsQ0FBQyxDQUFBLENBQUE7SUFXRCxDQUFDO0NBQ0Y7QUEvRUQsMEJBK0VDIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgQVBJLUF1dGhcbiAqL1xuaW1wb3J0IEF2YWxhbmNoZUNvcmUgZnJvbSBcIi4uLy4uL2F2YWxhbmNoZVwiXG5pbXBvcnQgeyBKUlBDQVBJIH0gZnJvbSBcIi4uLy4uL2NvbW1vbi9qcnBjYXBpXCJcbmltcG9ydCB7IFJlcXVlc3RSZXNwb25zZURhdGEgfSBmcm9tIFwiLi4vLi4vY29tbW9uL2FwaWJhc2VcIlxuaW1wb3J0IHsgRXJyb3JSZXNwb25zZU9iamVjdCB9IGZyb20gXCIuLi8uLi91dGlscy9lcnJvcnNcIlxuaW1wb3J0IHtcbiAgQ2hhbmdlUGFzc3dvcmRJbnRlcmZhY2UsXG4gIE5ld1Rva2VuSW50ZXJmYWNlLFxuICBSZXZva2VUb2tlbkludGVyZmFjZVxufSBmcm9tIFwiLi9pbnRlcmZhY2VzXCJcblxuLyoqXG4gKiBDbGFzcyBmb3IgaW50ZXJhY3Rpbmcgd2l0aCBhIG5vZGUncyBBdXRoQVBJLlxuICpcbiAqIEBjYXRlZ29yeSBSUENBUElzXG4gKlxuICogQHJlbWFya3MgVGhpcyBleHRlbmRzIHRoZSBbW0pSUENBUEldXSBjbGFzcy4gVGhpcyBjbGFzcyBzaG91bGQgbm90IGJlIGRpcmVjdGx5IGNhbGxlZC4gSW5zdGVhZCwgdXNlIHRoZSBbW0F2YWxhbmNoZS5hZGRBUEldXSBmdW5jdGlvbiB0byByZWdpc3RlciB0aGlzIGludGVyZmFjZSB3aXRoIEF2YWxhbmNoZS5cbiAqL1xuZXhwb3J0IGNsYXNzIEF1dGhBUEkgZXh0ZW5kcyBKUlBDQVBJIHtcbiAgLyoqXG4gICAqIENyZWF0ZXMgYSBuZXcgYXV0aG9yaXphdGlvbiB0b2tlbiB0aGF0IGdyYW50cyBhY2Nlc3MgdG8gb25lIG9yIG1vcmUgQVBJIGVuZHBvaW50cy5cbiAgICpcbiAgICogQHBhcmFtIHBhc3N3b3JkIFRoaXMgbm9kZSdzIGF1dGhvcml6YXRpb24gdG9rZW4gcGFzc3dvcmQsIHNldCB0aHJvdWdoIHRoZSBDTEkgd2hlbiB0aGUgbm9kZSB3YXMgbGF1bmNoZWQuXG4gICAqIEBwYXJhbSBlbmRwb2ludHMgQSBsaXN0IG9mIGVuZHBvaW50cyB0aGF0IHdpbGwgYmUgYWNjZXNzaWJsZSB1c2luZyB0aGUgZ2VuZXJhdGVkIHRva2VuLiBJZiB0aGVyZVwicyBhbiBlbGVtZW50IHRoYXQgaXMgXCIqXCIsIHRoaXMgdG9rZW4gY2FuIHJlYWNoIGFueSBlbmRwb2ludC5cbiAgICpcbiAgICogQHJldHVybnMgUmV0dXJucyBhIFByb21pc2Ugc3RyaW5nIGNvbnRhaW5pbmcgdGhlIGF1dGhvcml6YXRpb24gdG9rZW4uXG4gICAqL1xuICBuZXdUb2tlbiA9IGFzeW5jIChcbiAgICBwYXNzd29yZDogc3RyaW5nLFxuICAgIGVuZHBvaW50czogc3RyaW5nW11cbiAgKTogUHJvbWlzZTxzdHJpbmcgfCBFcnJvclJlc3BvbnNlT2JqZWN0PiA9PiB7XG4gICAgY29uc3QgcGFyYW1zOiBOZXdUb2tlbkludGVyZmFjZSA9IHtcbiAgICAgIHBhc3N3b3JkLFxuICAgICAgZW5kcG9pbnRzXG4gICAgfVxuICAgIGNvbnN0IHJlc3BvbnNlOiBSZXF1ZXN0UmVzcG9uc2VEYXRhID0gYXdhaXQgdGhpcy5jYWxsTWV0aG9kKFxuICAgICAgXCJhdXRoLm5ld1Rva2VuXCIsXG4gICAgICBwYXJhbXNcbiAgICApXG4gICAgcmV0dXJuIHJlc3BvbnNlLmRhdGEucmVzdWx0LnRva2VuXG4gICAgICA/IHJlc3BvbnNlLmRhdGEucmVzdWx0LnRva2VuXG4gICAgICA6IHJlc3BvbnNlLmRhdGEucmVzdWx0XG4gIH1cblxuICAvKipcbiAgICogUmV2b2tlcyBhbiBhdXRob3JpemF0aW9uIHRva2VuLCByZW1vdmluZyBhbGwgb2YgaXRzIHJpZ2h0cyB0byBhY2Nlc3MgZW5kcG9pbnRzLlxuICAgKlxuICAgKiBAcGFyYW0gcGFzc3dvcmQgVGhpcyBub2RlJ3MgYXV0aG9yaXphdGlvbiB0b2tlbiBwYXNzd29yZCwgc2V0IHRocm91Z2ggdGhlIENMSSB3aGVuIHRoZSBub2RlIHdhcyBsYXVuY2hlZC5cbiAgICogQHBhcmFtIHRva2VuIEFuIGF1dGhvcml6YXRpb24gdG9rZW4gd2hvc2UgYWNjZXNzIHNob3VsZCBiZSByZXZva2VkLlxuICAgKlxuICAgKiBAcmV0dXJucyBSZXR1cm5zIGEgUHJvbWlzZSBib29sZWFuIGluZGljYXRpbmcgaWYgYSB0b2tlbiB3YXMgc3VjY2Vzc2Z1bGx5IHJldm9rZWQuXG4gICAqL1xuICByZXZva2VUb2tlbiA9IGFzeW5jIChwYXNzd29yZDogc3RyaW5nLCB0b2tlbjogc3RyaW5nKTogUHJvbWlzZTxib29sZWFuPiA9PiB7XG4gICAgY29uc3QgcGFyYW1zOiBSZXZva2VUb2tlbkludGVyZmFjZSA9IHtcbiAgICAgIHBhc3N3b3JkLFxuICAgICAgdG9rZW5cbiAgICB9XG4gICAgY29uc3QgcmVzcG9uc2U6IFJlcXVlc3RSZXNwb25zZURhdGEgPSBhd2FpdCB0aGlzLmNhbGxNZXRob2QoXG4gICAgICBcImF1dGgucmV2b2tlVG9rZW5cIixcbiAgICAgIHBhcmFtc1xuICAgIClcbiAgICByZXR1cm4gcmVzcG9uc2UuZGF0YS5yZXN1bHQuc3VjY2Vzc1xuICB9XG5cbiAgLyoqXG4gICAqIENoYW5nZSB0aGlzIG5vZGUncyBhdXRob3JpemF0aW9uIHRva2VuIHBhc3N3b3JkLiAqKkFueSBhdXRob3JpemF0aW9uIHRva2VucyBjcmVhdGVkIHVuZGVyIGFuIG9sZCBwYXNzd29yZCB3aWxsIGJlY29tZSBpbnZhbGlkLioqXG4gICAqXG4gICAqIEBwYXJhbSBvbGRQYXNzd29yZCBUaGlzIG5vZGUncyBhdXRob3JpemF0aW9uIHRva2VuIHBhc3N3b3JkLCBzZXQgdGhyb3VnaCB0aGUgQ0xJIHdoZW4gdGhlIG5vZGUgd2FzIGxhdW5jaGVkLlxuICAgKiBAcGFyYW0gbmV3UGFzc3dvcmQgQSBuZXcgcGFzc3dvcmQgZm9yIHRoaXMgbm9kZSdzIGF1dGhvcml6YXRpb24gdG9rZW4gaXNzdWFuY2UuXG4gICAqXG4gICAqIEByZXR1cm5zIFJldHVybnMgYSBQcm9taXNlIGJvb2xlYW4gaW5kaWNhdGluZyBpZiB0aGUgcGFzc3dvcmQgd2FzIHN1Y2Nlc3NmdWxseSBjaGFuZ2VkLlxuICAgKi9cbiAgY2hhbmdlUGFzc3dvcmQgPSBhc3luYyAoXG4gICAgb2xkUGFzc3dvcmQ6IHN0cmluZyxcbiAgICBuZXdQYXNzd29yZDogc3RyaW5nXG4gICk6IFByb21pc2U8Ym9vbGVhbj4gPT4ge1xuICAgIGNvbnN0IHBhcmFtczogQ2hhbmdlUGFzc3dvcmRJbnRlcmZhY2UgPSB7XG4gICAgICBvbGRQYXNzd29yZCxcbiAgICAgIG5ld1Bhc3N3b3JkXG4gICAgfVxuICAgIGNvbnN0IHJlc3BvbnNlOiBSZXF1ZXN0UmVzcG9uc2VEYXRhID0gYXdhaXQgdGhpcy5jYWxsTWV0aG9kKFxuICAgICAgXCJhdXRoLmNoYW5nZVBhc3N3b3JkXCIsXG4gICAgICBwYXJhbXNcbiAgICApXG4gICAgcmV0dXJuIHJlc3BvbnNlLmRhdGEucmVzdWx0LnN1Y2Nlc3NcbiAgfVxuXG4gIC8qKlxuICAgKiBUaGlzIGNsYXNzIHNob3VsZCBub3QgYmUgaW5zdGFudGlhdGVkIGRpcmVjdGx5LiBJbnN0ZWFkIHVzZSB0aGUgW1tBdmFsYW5jaGUuYWRkQVBJXV1cbiAgICogbWV0aG9kLlxuICAgKlxuICAgKiBAcGFyYW0gY29yZSBBIHJlZmVyZW5jZSB0byB0aGUgQXZhbGFuY2hlIGNsYXNzXG4gICAqIEBwYXJhbSBiYXNlVVJMIERlZmF1bHRzIHRvIHRoZSBzdHJpbmcgXCIvZXh0L2F1dGhcIiBhcyB0aGUgcGF0aCB0byBycGMncyBiYXNlVVJMXG4gICAqL1xuICBjb25zdHJ1Y3Rvcihjb3JlOiBBdmFsYW5jaGVDb3JlLCBiYXNlVVJMOiBzdHJpbmcgPSBcIi9leHQvYXV0aFwiKSB7XG4gICAgc3VwZXIoY29yZSwgYmFzZVVSTClcbiAgfVxufVxuIl19