"use strict";
/**
 * @packageDocumentation
 * @module API-PlatformVM-Credentials
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
    if (credid === constants_1.PlatformVMConstants.SECPCREDENTIAL) {
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
        this._typeID = constants_1.PlatformVMConstants.SECPCREDENTIAL;
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
        let newbasetx = (0, exports.SelectCredentialClass)(id, ...args);
        return newbasetx;
    }
}
exports.SECPCredential = SECPCredential;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY3JlZGVudGlhbHMuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi9zcmMvYXBpcy9wbGF0Zm9ybXZtL2NyZWRlbnRpYWxzLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7QUFBQTs7O0dBR0c7OztBQUVILDJDQUFpRDtBQUNqRCwwREFBcUQ7QUFDckQsK0NBQWdEO0FBRWhEOzs7Ozs7R0FNRztBQUNJLE1BQU0scUJBQXFCLEdBQUcsQ0FDbkMsTUFBYyxFQUNkLEdBQUcsSUFBVyxFQUNGLEVBQUU7SUFDZCxJQUFJLE1BQU0sS0FBSywrQkFBbUIsQ0FBQyxjQUFjLEVBQUU7UUFDakQsT0FBTyxJQUFJLGNBQWMsQ0FBQyxHQUFHLElBQUksQ0FBQyxDQUFBO0tBQ25DO0lBQ0QsMEJBQTBCO0lBQzFCLE1BQU0sSUFBSSxvQkFBVyxDQUFDLCtDQUErQyxDQUFDLENBQUE7QUFDeEUsQ0FBQyxDQUFBO0FBVFksUUFBQSxxQkFBcUIseUJBU2pDO0FBRUQsTUFBYSxjQUFlLFNBQVEsd0JBQVU7SUFBOUM7O1FBQ1ksY0FBUyxHQUFHLGdCQUFnQixDQUFBO1FBQzVCLFlBQU8sR0FBRywrQkFBbUIsQ0FBQyxjQUFjLENBQUE7SUFzQnhELENBQUM7SUFwQkMsOENBQThDO0lBRTlDLGVBQWU7UUFDYixPQUFPLElBQUksQ0FBQyxPQUFPLENBQUE7SUFDckIsQ0FBQztJQUVELEtBQUs7UUFDSCxJQUFJLE9BQU8sR0FBbUIsSUFBSSxjQUFjLEVBQUUsQ0FBQTtRQUNsRCxPQUFPLENBQUMsVUFBVSxDQUFDLElBQUksQ0FBQyxRQUFRLEVBQUUsQ0FBQyxDQUFBO1FBQ25DLE9BQU8sT0FBZSxDQUFBO0lBQ3hCLENBQUM7SUFFRCxNQUFNLENBQUMsR0FBRyxJQUFXO1FBQ25CLE9BQU8sSUFBSSxjQUFjLENBQUMsR0FBRyxJQUFJLENBQVMsQ0FBQTtJQUM1QyxDQUFDO0lBRUQsTUFBTSxDQUFDLEVBQVUsRUFBRSxHQUFHLElBQVc7UUFDL0IsSUFBSSxTQUFTLEdBQWUsSUFBQSw2QkFBcUIsRUFBQyxFQUFFLEVBQUUsR0FBRyxJQUFJLENBQUMsQ0FBQTtRQUM5RCxPQUFPLFNBQVMsQ0FBQTtJQUNsQixDQUFDO0NBQ0Y7QUF4QkQsd0NBd0JDIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgQVBJLVBsYXRmb3JtVk0tQ3JlZGVudGlhbHNcbiAqL1xuXG5pbXBvcnQgeyBQbGF0Zm9ybVZNQ29uc3RhbnRzIH0gZnJvbSBcIi4vY29uc3RhbnRzXCJcbmltcG9ydCB7IENyZWRlbnRpYWwgfSBmcm9tIFwiLi4vLi4vY29tbW9uL2NyZWRlbnRpYWxzXCJcbmltcG9ydCB7IENyZWRJZEVycm9yIH0gZnJvbSBcIi4uLy4uL3V0aWxzL2Vycm9yc1wiXG5cbi8qKlxuICogVGFrZXMgYSBidWZmZXIgcmVwcmVzZW50aW5nIHRoZSBjcmVkZW50aWFsIGFuZCByZXR1cm5zIHRoZSBwcm9wZXIgW1tDcmVkZW50aWFsXV0gaW5zdGFuY2UuXG4gKlxuICogQHBhcmFtIGNyZWRpZCBBIG51bWJlciByZXByZXNlbnRpbmcgdGhlIGNyZWRlbnRpYWwgSUQgcGFyc2VkIHByaW9yIHRvIHRoZSBieXRlcyBwYXNzZWQgaW5cbiAqXG4gKiBAcmV0dXJucyBBbiBpbnN0YW5jZSBvZiBhbiBbW0NyZWRlbnRpYWxdXS1leHRlbmRlZCBjbGFzcy5cbiAqL1xuZXhwb3J0IGNvbnN0IFNlbGVjdENyZWRlbnRpYWxDbGFzcyA9IChcbiAgY3JlZGlkOiBudW1iZXIsXG4gIC4uLmFyZ3M6IGFueVtdXG4pOiBDcmVkZW50aWFsID0+IHtcbiAgaWYgKGNyZWRpZCA9PT0gUGxhdGZvcm1WTUNvbnN0YW50cy5TRUNQQ1JFREVOVElBTCkge1xuICAgIHJldHVybiBuZXcgU0VDUENyZWRlbnRpYWwoLi4uYXJncylcbiAgfVxuICAvKiBpc3RhbmJ1bCBpZ25vcmUgbmV4dCAqL1xuICB0aHJvdyBuZXcgQ3JlZElkRXJyb3IoXCJFcnJvciAtIFNlbGVjdENyZWRlbnRpYWxDbGFzczogdW5rbm93biBjcmVkaWRcIilcbn1cblxuZXhwb3J0IGNsYXNzIFNFQ1BDcmVkZW50aWFsIGV4dGVuZHMgQ3JlZGVudGlhbCB7XG4gIHByb3RlY3RlZCBfdHlwZU5hbWUgPSBcIlNFQ1BDcmVkZW50aWFsXCJcbiAgcHJvdGVjdGVkIF90eXBlSUQgPSBQbGF0Zm9ybVZNQ29uc3RhbnRzLlNFQ1BDUkVERU5USUFMXG5cbiAgLy9zZXJpYWxpemUgYW5kIGRlc2VyaWFsaXplIGJvdGggYXJlIGluaGVyaXRlZFxuXG4gIGdldENyZWRlbnRpYWxJRCgpOiBudW1iZXIge1xuICAgIHJldHVybiB0aGlzLl90eXBlSURcbiAgfVxuXG4gIGNsb25lKCk6IHRoaXMge1xuICAgIGxldCBuZXdiYXNlOiBTRUNQQ3JlZGVudGlhbCA9IG5ldyBTRUNQQ3JlZGVudGlhbCgpXG4gICAgbmV3YmFzZS5mcm9tQnVmZmVyKHRoaXMudG9CdWZmZXIoKSlcbiAgICByZXR1cm4gbmV3YmFzZSBhcyB0aGlzXG4gIH1cblxuICBjcmVhdGUoLi4uYXJnczogYW55W10pOiB0aGlzIHtcbiAgICByZXR1cm4gbmV3IFNFQ1BDcmVkZW50aWFsKC4uLmFyZ3MpIGFzIHRoaXNcbiAgfVxuXG4gIHNlbGVjdChpZDogbnVtYmVyLCAuLi5hcmdzOiBhbnlbXSk6IENyZWRlbnRpYWwge1xuICAgIGxldCBuZXdiYXNldHg6IENyZWRlbnRpYWwgPSBTZWxlY3RDcmVkZW50aWFsQ2xhc3MoaWQsIC4uLmFyZ3MpXG4gICAgcmV0dXJuIG5ld2Jhc2V0eFxuICB9XG59XG4iXX0=