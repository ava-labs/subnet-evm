/**
 * @packageDocumentation
 * @module API-Auth
 */
import AvalancheCore from "../../avalanche";
import { JRPCAPI } from "../../common/jrpcapi";
import { ErrorResponseObject } from "../../utils/errors";
/**
 * Class for interacting with a node's AuthAPI.
 *
 * @category RPCAPIs
 *
 * @remarks This extends the [[JRPCAPI]] class. This class should not be directly called. Instead, use the [[Avalanche.addAPI]] function to register this interface with Avalanche.
 */
export declare class AuthAPI extends JRPCAPI {
    /**
     * Creates a new authorization token that grants access to one or more API endpoints.
     *
     * @param password This node's authorization token password, set through the CLI when the node was launched.
     * @param endpoints A list of endpoints that will be accessible using the generated token. If there"s an element that is "*", this token can reach any endpoint.
     *
     * @returns Returns a Promise string containing the authorization token.
     */
    newToken: (password: string, endpoints: string[]) => Promise<string | ErrorResponseObject>;
    /**
     * Revokes an authorization token, removing all of its rights to access endpoints.
     *
     * @param password This node's authorization token password, set through the CLI when the node was launched.
     * @param token An authorization token whose access should be revoked.
     *
     * @returns Returns a Promise boolean indicating if a token was successfully revoked.
     */
    revokeToken: (password: string, token: string) => Promise<boolean>;
    /**
     * Change this node's authorization token password. **Any authorization tokens created under an old password will become invalid.**
     *
     * @param oldPassword This node's authorization token password, set through the CLI when the node was launched.
     * @param newPassword A new password for this node's authorization token issuance.
     *
     * @returns Returns a Promise boolean indicating if the password was successfully changed.
     */
    changePassword: (oldPassword: string, newPassword: string) => Promise<boolean>;
    /**
     * This class should not be instantiated directly. Instead use the [[Avalanche.addAPI]]
     * method.
     *
     * @param core A reference to the Avalanche class
     * @param baseURL Defaults to the string "/ext/auth" as the path to rpc's baseURL
     */
    constructor(core: AvalancheCore, baseURL?: string);
}
//# sourceMappingURL=api.d.ts.map