"use strict";
/**
 * @packageDocumentation
 * @module API-EVM-Credentials
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.SECPCredential = exports.SelectCredentialClass = void 0;
const constants_1 = require("./constants");
const credentials_1 = require("../../common/credentials");
const errors_1 = require("../../utils/errors");
/**
 * Takes a buffer representing the credential and returns the proper [[Credential]] instance.
 *
 * @param credid A number representing the credential ID parsed prior to the bytes passed in
 *
 * @returns An instance of an [[Credential]]-extended class.
 */
const SelectCredentialClass = (credid, ...args) => {
    if (credid === constants_1.EVMConstants.SECPCREDENTIAL) {
        return new SECPCredential(...args);
    }
    /* istanbul ignore next */
    throw new errors_1.CredIdError("Error - SelectCredentialClass: unknown credid");
};
exports.SelectCredentialClass = SelectCredentialClass;
class SECPCredential extends credentials_1.Credential {
    constructor() {
        super(...arguments);
        this._typeName = "SECPCredential";
        this._typeID = constants_1.EVMConstants.SECPCREDENTIAL;
    }
    //serialize and deserialize both are inherited
    getCredentialID() {
        return this._typeID;
    }
    clone() {
        let newbase = new SECPCredential();
        newbase.fromBuffer(this.toBuffer());
        return newbase;
    }
    create(...args) {
        return new SECPCredential(...args);
    }
    select(id, ...args) {
        let credential = (0, exports.SelectCredentialClass)(id, ...args);
        return credential;
    }
}
exports.SECPCredential = SECPCredential;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY3JlZGVudGlhbHMuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi9zcmMvYXBpcy9ldm0vY3JlZGVudGlhbHMudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IjtBQUFBOzs7R0FHRzs7O0FBRUgsMkNBQTBDO0FBQzFDLDBEQUFxRDtBQUNyRCwrQ0FBZ0Q7QUFFaEQ7Ozs7OztHQU1HO0FBQ0ksTUFBTSxxQkFBcUIsR0FBRyxDQUNuQyxNQUFjLEVBQ2QsR0FBRyxJQUFXLEVBQ0YsRUFBRTtJQUNkLElBQUksTUFBTSxLQUFLLHdCQUFZLENBQUMsY0FBYyxFQUFFO1FBQzFDLE9BQU8sSUFBSSxjQUFjLENBQUMsR0FBRyxJQUFJLENBQUMsQ0FBQTtLQUNuQztJQUNELDBCQUEwQjtJQUMxQixNQUFNLElBQUksb0JBQVcsQ0FBQywrQ0FBK0MsQ0FBQyxDQUFBO0FBQ3hFLENBQUMsQ0FBQTtBQVRZLFFBQUEscUJBQXFCLHlCQVNqQztBQUVELE1BQWEsY0FBZSxTQUFRLHdCQUFVO0lBQTlDOztRQUNZLGNBQVMsR0FBVyxnQkFBZ0IsQ0FBQTtRQUNwQyxZQUFPLEdBQVcsd0JBQVksQ0FBQyxjQUFjLENBQUE7SUFzQnpELENBQUM7SUFwQkMsOENBQThDO0lBRTlDLGVBQWU7UUFDYixPQUFPLElBQUksQ0FBQyxPQUFPLENBQUE7SUFDckIsQ0FBQztJQUVELEtBQUs7UUFDSCxJQUFJLE9BQU8sR0FBbUIsSUFBSSxjQUFjLEVBQUUsQ0FBQTtRQUNsRCxPQUFPLENBQUMsVUFBVSxDQUFDLElBQUksQ0FBQyxRQUFRLEVBQUUsQ0FBQyxDQUFBO1FBQ25DLE9BQU8sT0FBZSxDQUFBO0lBQ3hCLENBQUM7SUFFRCxNQUFNLENBQUMsR0FBRyxJQUFXO1FBQ25CLE9BQU8sSUFBSSxjQUFjLENBQUMsR0FBRyxJQUFJLENBQVMsQ0FBQTtJQUM1QyxDQUFDO0lBRUQsTUFBTSxDQUFDLEVBQVUsRUFBRSxHQUFHLElBQVc7UUFDL0IsSUFBSSxVQUFVLEdBQWUsSUFBQSw2QkFBcUIsRUFBQyxFQUFFLEVBQUUsR0FBRyxJQUFJLENBQUMsQ0FBQTtRQUMvRCxPQUFPLFVBQVUsQ0FBQTtJQUNuQixDQUFDO0NBQ0Y7QUF4QkQsd0NBd0JDIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgQVBJLUVWTS1DcmVkZW50aWFsc1xuICovXG5cbmltcG9ydCB7IEVWTUNvbnN0YW50cyB9IGZyb20gXCIuL2NvbnN0YW50c1wiXG5pbXBvcnQgeyBDcmVkZW50aWFsIH0gZnJvbSBcIi4uLy4uL2NvbW1vbi9jcmVkZW50aWFsc1wiXG5pbXBvcnQgeyBDcmVkSWRFcnJvciB9IGZyb20gXCIuLi8uLi91dGlscy9lcnJvcnNcIlxuXG4vKipcbiAqIFRha2VzIGEgYnVmZmVyIHJlcHJlc2VudGluZyB0aGUgY3JlZGVudGlhbCBhbmQgcmV0dXJucyB0aGUgcHJvcGVyIFtbQ3JlZGVudGlhbF1dIGluc3RhbmNlLlxuICpcbiAqIEBwYXJhbSBjcmVkaWQgQSBudW1iZXIgcmVwcmVzZW50aW5nIHRoZSBjcmVkZW50aWFsIElEIHBhcnNlZCBwcmlvciB0byB0aGUgYnl0ZXMgcGFzc2VkIGluXG4gKlxuICogQHJldHVybnMgQW4gaW5zdGFuY2Ugb2YgYW4gW1tDcmVkZW50aWFsXV0tZXh0ZW5kZWQgY2xhc3MuXG4gKi9cbmV4cG9ydCBjb25zdCBTZWxlY3RDcmVkZW50aWFsQ2xhc3MgPSAoXG4gIGNyZWRpZDogbnVtYmVyLFxuICAuLi5hcmdzOiBhbnlbXVxuKTogQ3JlZGVudGlhbCA9PiB7XG4gIGlmIChjcmVkaWQgPT09IEVWTUNvbnN0YW50cy5TRUNQQ1JFREVOVElBTCkge1xuICAgIHJldHVybiBuZXcgU0VDUENyZWRlbnRpYWwoLi4uYXJncylcbiAgfVxuICAvKiBpc3RhbmJ1bCBpZ25vcmUgbmV4dCAqL1xuICB0aHJvdyBuZXcgQ3JlZElkRXJyb3IoXCJFcnJvciAtIFNlbGVjdENyZWRlbnRpYWxDbGFzczogdW5rbm93biBjcmVkaWRcIilcbn1cblxuZXhwb3J0IGNsYXNzIFNFQ1BDcmVkZW50aWFsIGV4dGVuZHMgQ3JlZGVudGlhbCB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWU6IHN0cmluZyA9IFwiU0VDUENyZWRlbnRpYWxcIlxuICBwcm90ZWN0ZWQgX3R5cGVJRDogbnVtYmVyID0gRVZNQ29uc3RhbnRzLlNFQ1BDUkVERU5USUFMXG5cbiAgLy9zZXJpYWxpemUgYW5kIGRlc2VyaWFsaXplIGJvdGggYXJlIGluaGVyaXRlZFxuXG4gIGdldENyZWRlbnRpYWxJRCgpOiBudW1iZXIge1xuICAgIHJldHVybiB0aGlzLl90eXBlSURcbiAgfVxuXG4gIGNsb25lKCk6IHRoaXMge1xuICAgIGxldCBuZXdiYXNlOiBTRUNQQ3JlZGVudGlhbCA9IG5ldyBTRUNQQ3JlZGVudGlhbCgpXG4gICAgbmV3YmFzZS5mcm9tQnVmZmVyKHRoaXMudG9CdWZmZXIoKSlcbiAgICByZXR1cm4gbmV3YmFzZSBhcyB0aGlzXG4gIH1cblxuICBjcmVhdGUoLi4uYXJnczogYW55W10pOiB0aGlzIHtcbiAgICByZXR1cm4gbmV3IFNFQ1BDcmVkZW50aWFsKC4uLmFyZ3MpIGFzIHRoaXNcbiAgfVxuXG4gIHNlbGVjdChpZDogbnVtYmVyLCAuLi5hcmdzOiBhbnlbXSk6IENyZWRlbnRpYWwge1xuICAgIGxldCBjcmVkZW50aWFsOiBDcmVkZW50aWFsID0gU2VsZWN0Q3JlZGVudGlhbENsYXNzKGlkLCAuLi5hcmdzKVxuICAgIHJldHVybiBjcmVkZW50aWFsXG4gIH1cbn1cbiJdfQ==