"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.Socket = void 0;
const isomorphic_ws_1 = __importDefault(require("isomorphic-ws"));
const utils_1 = require("../../utils");
class Socket extends isomorphic_ws_1.default {
    /**
     * Provides the API for creating and managing a WebSocket connection to a server, as well as for sending and receiving data on the connection.
     *
     * @param url Defaults to [[MainnetAPI]]
     * @param options Optional
     */
    constructor(url = `wss://${utils_1.MainnetAPI}:443/ext/bc/X/events`, options) {
        super(url, options);
    }
    /**
     * Send a message to the server
     *
     * @param data
     * @param cb Optional
     */
    send(data, cb) {
        super.send(data, cb);
    }
    /**
     * Terminates the connection completely
     *
     * @param mcode Optional
     * @param data Optional
     */
    close(mcode, data) {
        super.close(mcode, data);
    }
}
exports.Socket = Socket;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoic29ja2V0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2FwaXMvc29ja2V0L3NvY2tldC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7QUFLQSxrRUFBcUM7QUFDckMsdUNBQXdDO0FBQ3hDLE1BQWEsTUFBTyxTQUFRLHVCQUFTO0lBOEJuQzs7Ozs7T0FLRztJQUNILFlBQ0UsTUFBa0MsU0FBUyxrQkFBVSxzQkFBc0IsRUFDM0UsT0FBcUQ7UUFFckQsS0FBSyxDQUFDLEdBQUcsRUFBRSxPQUFPLENBQUMsQ0FBQTtJQUNyQixDQUFDO0lBL0JEOzs7OztPQUtHO0lBQ0gsSUFBSSxDQUFDLElBQVMsRUFBRSxFQUFRO1FBQ3RCLEtBQUssQ0FBQyxJQUFJLENBQUMsSUFBSSxFQUFFLEVBQUUsQ0FBQyxDQUFBO0lBQ3RCLENBQUM7SUFFRDs7Ozs7T0FLRztJQUNILEtBQUssQ0FBQyxLQUFjLEVBQUUsSUFBYTtRQUNqQyxLQUFLLENBQUMsS0FBSyxDQUFDLEtBQUssRUFBRSxJQUFJLENBQUMsQ0FBQTtJQUMxQixDQUFDO0NBY0Y7QUExQ0Qsd0JBMENDIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBAcGFja2FnZURvY3VtZW50YXRpb25cbiAqIEBtb2R1bGUgQVBJLVNvY2tldFxuICovXG5pbXBvcnQgeyBDbGllbnRSZXF1ZXN0QXJncyB9IGZyb20gXCJodHRwXCJcbmltcG9ydCBXZWJTb2NrZXQgZnJvbSBcImlzb21vcnBoaWMtd3NcIlxuaW1wb3J0IHsgTWFpbm5ldEFQSSB9IGZyb20gXCIuLi8uLi91dGlsc1wiXG5leHBvcnQgY2xhc3MgU29ja2V0IGV4dGVuZHMgV2ViU29ja2V0IHtcbiAgLy8gRmlyZXMgb25jZSB0aGUgY29ubmVjdGlvbiBoYXMgYmVlbiBlc3RhYmxpc2hlZCBiZXR3ZWVuIHRoZSBjbGllbnQgYW5kIHRoZSBzZXJ2ZXJcbiAgb25vcGVuOiBhbnlcbiAgLy8gRmlyZXMgd2hlbiB0aGUgc2VydmVyIHNlbmRzIHNvbWUgZGF0YVxuICBvbm1lc3NhZ2U6IGFueVxuICAvLyBGaXJlcyBhZnRlciBlbmQgb2YgdGhlIGNvbW11bmljYXRpb24gYmV0d2VlbiBzZXJ2ZXIgYW5kIHRoZSBjbGllbnRcbiAgb25jbG9zZTogYW55XG4gIC8vIEZpcmVzIGZvciBzb21lIG1pc3Rha2UsIHdoaWNoIGhhcHBlbnMgZHVyaW5nIHRoZSBjb21tdW5pY2F0aW9uXG4gIG9uZXJyb3I6IGFueVxuXG4gIC8qKlxuICAgKiBTZW5kIGEgbWVzc2FnZSB0byB0aGUgc2VydmVyXG4gICAqXG4gICAqIEBwYXJhbSBkYXRhXG4gICAqIEBwYXJhbSBjYiBPcHRpb25hbFxuICAgKi9cbiAgc2VuZChkYXRhOiBhbnksIGNiPzogYW55KTogdm9pZCB7XG4gICAgc3VwZXIuc2VuZChkYXRhLCBjYilcbiAgfVxuXG4gIC8qKlxuICAgKiBUZXJtaW5hdGVzIHRoZSBjb25uZWN0aW9uIGNvbXBsZXRlbHlcbiAgICpcbiAgICogQHBhcmFtIG1jb2RlIE9wdGlvbmFsXG4gICAqIEBwYXJhbSBkYXRhIE9wdGlvbmFsXG4gICAqL1xuICBjbG9zZShtY29kZT86IG51bWJlciwgZGF0YT86IHN0cmluZyk6IHZvaWQge1xuICAgIHN1cGVyLmNsb3NlKG1jb2RlLCBkYXRhKVxuICB9XG5cbiAgLyoqXG4gICAqIFByb3ZpZGVzIHRoZSBBUEkgZm9yIGNyZWF0aW5nIGFuZCBtYW5hZ2luZyBhIFdlYlNvY2tldCBjb25uZWN0aW9uIHRvIGEgc2VydmVyLCBhcyB3ZWxsIGFzIGZvciBzZW5kaW5nIGFuZCByZWNlaXZpbmcgZGF0YSBvbiB0aGUgY29ubmVjdGlvbi5cbiAgICpcbiAgICogQHBhcmFtIHVybCBEZWZhdWx0cyB0byBbW01haW5uZXRBUEldXVxuICAgKiBAcGFyYW0gb3B0aW9ucyBPcHRpb25hbFxuICAgKi9cbiAgY29uc3RydWN0b3IoXG4gICAgdXJsOiBzdHJpbmcgfCBpbXBvcnQoXCJ1cmxcIikuVVJMID0gYHdzczovLyR7TWFpbm5ldEFQSX06NDQzL2V4dC9iYy9YL2V2ZW50c2AsXG4gICAgb3B0aW9ucz86IFdlYlNvY2tldC5DbGllbnRPcHRpb25zIHwgQ2xpZW50UmVxdWVzdEFyZ3NcbiAgKSB7XG4gICAgc3VwZXIodXJsLCBvcHRpb25zKVxuICB9XG59XG4iXX0=