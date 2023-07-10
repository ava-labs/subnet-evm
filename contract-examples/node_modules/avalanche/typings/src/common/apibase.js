"use strict";
/**
 * @packageDocumentation
 * @module Common-APIBase
 */
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.APIBase = exports.RequestResponseData = void 0;
const db_1 = __importDefault(require("../utils/db"));
/**
 * Response data for HTTP requests.
 */
class RequestResponseData {
    constructor(data, headers, status, statusText, request) {
        this.data = data;
        this.headers = headers;
        this.status = status;
        this.statusText = statusText;
        this.request = request;
    }
}
exports.RequestResponseData = RequestResponseData;
/**
 * Abstract class defining a generic endpoint that all endpoints must implement (extend).
 */
class APIBase {
    /**
     *
     * @param core Reference to the Avalanche instance using this baseURL
     * @param baseURL Path to the baseURL
     */
    constructor(core, baseURL) {
        /**
         * Sets the path of the APIs baseURL.
         *
         * @param baseURL Path of the APIs baseURL - ex: "/ext/bc/X"
         */
        this.setBaseURL = (baseURL) => {
            if (this.db && this.baseURL !== baseURL) {
                const backup = this.db.getAll();
                this.db.clearAll();
                this.baseURL = baseURL;
                this.db = db_1.default.getNamespace(baseURL);
                this.db.setAll(backup, true);
            }
            else {
                this.baseURL = baseURL;
                this.db = db_1.default.getNamespace(baseURL);
            }
        };
        /**
         * Returns the baseURL's path.
         */
        this.getBaseURL = () => this.baseURL;
        /**
         * Returns the baseURL's database.
         */
        this.getDB = () => this.db;
        this.core = core;
        this.setBaseURL(baseURL);
    }
}
exports.APIBase = APIBase;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYXBpYmFzZS5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3NyYy9jb21tb24vYXBpYmFzZS50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiO0FBQUE7OztHQUdHOzs7Ozs7QUFJSCxxREFBNEI7QUFHNUI7O0dBRUc7QUFDSCxNQUFhLG1CQUFtQjtJQUM5QixZQUNTLElBQVMsRUFDVCxPQUFZLEVBQ1osTUFBYyxFQUNkLFVBQWtCLEVBQ2xCLE9BQXVDO1FBSnZDLFNBQUksR0FBSixJQUFJLENBQUs7UUFDVCxZQUFPLEdBQVAsT0FBTyxDQUFLO1FBQ1osV0FBTSxHQUFOLE1BQU0sQ0FBUTtRQUNkLGVBQVUsR0FBVixVQUFVLENBQVE7UUFDbEIsWUFBTyxHQUFQLE9BQU8sQ0FBZ0M7SUFDN0MsQ0FBQztDQUNMO0FBUkQsa0RBUUM7QUFFRDs7R0FFRztBQUNILE1BQXNCLE9BQU87SUFpQzNCOzs7O09BSUc7SUFDSCxZQUFZLElBQW1CLEVBQUUsT0FBZTtRQWpDaEQ7Ozs7V0FJRztRQUNILGVBQVUsR0FBRyxDQUFDLE9BQWUsRUFBRSxFQUFFO1lBQy9CLElBQUksSUFBSSxDQUFDLEVBQUUsSUFBSSxJQUFJLENBQUMsT0FBTyxLQUFLLE9BQU8sRUFBRTtnQkFDdkMsTUFBTSxNQUFNLEdBQUcsSUFBSSxDQUFDLEVBQUUsQ0FBQyxNQUFNLEVBQUUsQ0FBQTtnQkFDL0IsSUFBSSxDQUFDLEVBQUUsQ0FBQyxRQUFRLEVBQUUsQ0FBQTtnQkFDbEIsSUFBSSxDQUFDLE9BQU8sR0FBRyxPQUFPLENBQUE7Z0JBQ3RCLElBQUksQ0FBQyxFQUFFLEdBQUcsWUFBRSxDQUFDLFlBQVksQ0FBQyxPQUFPLENBQUMsQ0FBQTtnQkFDbEMsSUFBSSxDQUFDLEVBQUUsQ0FBQyxNQUFNLENBQUMsTUFBTSxFQUFFLElBQUksQ0FBQyxDQUFBO2FBQzdCO2lCQUFNO2dCQUNMLElBQUksQ0FBQyxPQUFPLEdBQUcsT0FBTyxDQUFBO2dCQUN0QixJQUFJLENBQUMsRUFBRSxHQUFHLFlBQUUsQ0FBQyxZQUFZLENBQUMsT0FBTyxDQUFDLENBQUE7YUFDbkM7UUFDSCxDQUFDLENBQUE7UUFFRDs7V0FFRztRQUNILGVBQVUsR0FBRyxHQUFXLEVBQUUsQ0FBQyxJQUFJLENBQUMsT0FBTyxDQUFBO1FBRXZDOztXQUVHO1FBQ0gsVUFBSyxHQUFHLEdBQWEsRUFBRSxDQUFDLElBQUksQ0FBQyxFQUFFLENBQUE7UUFRN0IsSUFBSSxDQUFDLElBQUksR0FBRyxJQUFJLENBQUE7UUFDaEIsSUFBSSxDQUFDLFVBQVUsQ0FBQyxPQUFPLENBQUMsQ0FBQTtJQUMxQixDQUFDO0NBQ0Y7QUExQ0QsMEJBMENDIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgQ29tbW9uLUFQSUJhc2VcbiAqL1xuXG5pbXBvcnQgeyBTdG9yZUFQSSB9IGZyb20gXCJzdG9yZTJcIlxuaW1wb3J0IHsgQ2xpZW50UmVxdWVzdCB9IGZyb20gXCJodHRwXCJcbmltcG9ydCBEQiBmcm9tIFwiLi4vdXRpbHMvZGJcIlxuaW1wb3J0IEF2YWxhbmNoZUNvcmUgZnJvbSBcIi4uL2F2YWxhbmNoZVwiXG5cbi8qKlxuICogUmVzcG9uc2UgZGF0YSBmb3IgSFRUUCByZXF1ZXN0cy5cbiAqL1xuZXhwb3J0IGNsYXNzIFJlcXVlc3RSZXNwb25zZURhdGEge1xuICBjb25zdHJ1Y3RvcihcbiAgICBwdWJsaWMgZGF0YTogYW55LFxuICAgIHB1YmxpYyBoZWFkZXJzOiBhbnksXG4gICAgcHVibGljIHN0YXR1czogbnVtYmVyLFxuICAgIHB1YmxpYyBzdGF0dXNUZXh0OiBzdHJpbmcsXG4gICAgcHVibGljIHJlcXVlc3Q6IENsaWVudFJlcXVlc3QgfCBYTUxIdHRwUmVxdWVzdFxuICApIHt9XG59XG5cbi8qKlxuICogQWJzdHJhY3QgY2xhc3MgZGVmaW5pbmcgYSBnZW5lcmljIGVuZHBvaW50IHRoYXQgYWxsIGVuZHBvaW50cyBtdXN0IGltcGxlbWVudCAoZXh0ZW5kKS5cbiAqL1xuZXhwb3J0IGFic3RyYWN0IGNsYXNzIEFQSUJhc2Uge1xuICBwcm90ZWN0ZWQgY29yZTogQXZhbGFuY2hlQ29yZVxuICBwcm90ZWN0ZWQgYmFzZVVSTDogc3RyaW5nXG4gIHByb3RlY3RlZCBkYjogU3RvcmVBUElcblxuICAvKipcbiAgICogU2V0cyB0aGUgcGF0aCBvZiB0aGUgQVBJcyBiYXNlVVJMLlxuICAgKlxuICAgKiBAcGFyYW0gYmFzZVVSTCBQYXRoIG9mIHRoZSBBUElzIGJhc2VVUkwgLSBleDogXCIvZXh0L2JjL1hcIlxuICAgKi9cbiAgc2V0QmFzZVVSTCA9IChiYXNlVVJMOiBzdHJpbmcpID0+IHtcbiAgICBpZiAodGhpcy5kYiAmJiB0aGlzLmJhc2VVUkwgIT09IGJhc2VVUkwpIHtcbiAgICAgIGNvbnN0IGJhY2t1cCA9IHRoaXMuZGIuZ2V0QWxsKClcbiAgICAgIHRoaXMuZGIuY2xlYXJBbGwoKVxuICAgICAgdGhpcy5iYXNlVVJMID0gYmFzZVVSTFxuICAgICAgdGhpcy5kYiA9IERCLmdldE5hbWVzcGFjZShiYXNlVVJMKVxuICAgICAgdGhpcy5kYi5zZXRBbGwoYmFja3VwLCB0cnVlKVxuICAgIH0gZWxzZSB7XG4gICAgICB0aGlzLmJhc2VVUkwgPSBiYXNlVVJMXG4gICAgICB0aGlzLmRiID0gREIuZ2V0TmFtZXNwYWNlKGJhc2VVUkwpXG4gICAgfVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIGJhc2VVUkwncyBwYXRoLlxuICAgKi9cbiAgZ2V0QmFzZVVSTCA9ICgpOiBzdHJpbmcgPT4gdGhpcy5iYXNlVVJMXG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIGJhc2VVUkwncyBkYXRhYmFzZS5cbiAgICovXG4gIGdldERCID0gKCk6IFN0b3JlQVBJID0+IHRoaXMuZGJcblxuICAvKipcbiAgICpcbiAgICogQHBhcmFtIGNvcmUgUmVmZXJlbmNlIHRvIHRoZSBBdmFsYW5jaGUgaW5zdGFuY2UgdXNpbmcgdGhpcyBiYXNlVVJMXG4gICAqIEBwYXJhbSBiYXNlVVJMIFBhdGggdG8gdGhlIGJhc2VVUkxcbiAgICovXG4gIGNvbnN0cnVjdG9yKGNvcmU6IEF2YWxhbmNoZUNvcmUsIGJhc2VVUkw6IHN0cmluZykge1xuICAgIHRoaXMuY29yZSA9IGNvcmVcbiAgICB0aGlzLnNldEJhc2VVUkwoYmFzZVVSTClcbiAgfVxufVxuIl19