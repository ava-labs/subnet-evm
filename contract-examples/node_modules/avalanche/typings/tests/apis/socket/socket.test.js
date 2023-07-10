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
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const jest_websocket_mock_1 = __importDefault(require("jest-websocket-mock"));
// const server = new WS("ws://localhost:1234")
// const client = new Socket("ws://localhost:1234")
describe("Socket", () => __awaiter(void 0, void 0, void 0, function* () {
    // await server.connected // wait for the server to have established the connection
    // the mock websocket server will record all the messages it receives
    // client.send("hello")
    // the mock websocket server can also send messages to all connected clients
    // server.send("hello everyone")
    // ...simulate an error and close the connection
    // server.error()
    // ...or gracefully close the connection
    // server.close()
    // The WS class also has a static "clean" method to gracefully close all open connections,
    // particularly useful to reset the environment between test runs.
    // WS.clean()
    test("foobar", () => __awaiter(void 0, void 0, void 0, function* () {
        const server = new jest_websocket_mock_1.default("ws://localhost:1234");
        // console.log(server)
        const cient = new WebSocket("ws://localhost:1234");
        // const client = new Socket("ws://localhost:1234/")
        console.log(cient);
        // await server.connected
        // client.send("hello")
        // await expect(server).toReceiveMessage("hello")
        // expect(server).toHaveReceivedMessages(["hello"])
    }));
}));
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoic29ja2V0LnRlc3QuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi90ZXN0cy9hcGlzL3NvY2tldC9zb2NrZXQudGVzdC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7Ozs7OztBQUtBLDhFQUFvQztBQUVwQywrQ0FBK0M7QUFFL0MsbURBQW1EO0FBQ25ELFFBQVEsQ0FBQyxRQUFRLEVBQUUsR0FBd0IsRUFBRTtJQUMzQyxtRkFBbUY7SUFFbkYscUVBQXFFO0lBQ3JFLHVCQUF1QjtJQUV2Qiw0RUFBNEU7SUFDNUUsZ0NBQWdDO0lBRWhDLGdEQUFnRDtJQUNoRCxpQkFBaUI7SUFFakIsd0NBQXdDO0lBQ3hDLGlCQUFpQjtJQUVqQiwwRkFBMEY7SUFDMUYsa0VBQWtFO0lBQ2xFLGFBQWE7SUFFYixJQUFJLENBQUMsUUFBUSxFQUFFLEdBQXdCLEVBQUU7UUFDdkMsTUFBTSxNQUFNLEdBQU8sSUFBSSw2QkFBRSxDQUFDLHFCQUFxQixDQUFDLENBQUE7UUFDaEQsc0JBQXNCO1FBQ3RCLE1BQU0sS0FBSyxHQUFHLElBQUksU0FBUyxDQUFDLHFCQUFxQixDQUFDLENBQUE7UUFDbEQsb0RBQW9EO1FBQ3BELE9BQU8sQ0FBQyxHQUFHLENBQUMsS0FBSyxDQUFDLENBQUE7UUFFbEIseUJBQXlCO1FBQ3pCLHVCQUF1QjtRQUN2QixpREFBaUQ7UUFDakQsbURBQW1EO0lBQ3JELENBQUMsQ0FBQSxDQUFDLENBQUE7QUFDSixDQUFDLENBQUEsQ0FBQyxDQUFBIiwic291cmNlc0NvbnRlbnQiOlsiaW1wb3J0IG1vY2tBeGlvcyBmcm9tIFwiamVzdC1tb2NrLWF4aW9zXCJcbmltcG9ydCB7IEF2YWxhbmNoZSwgU29ja2V0IH0gZnJvbSBcInNyY1wiXG5pbXBvcnQgeyBJbmZvQVBJIH0gZnJvbSBcInNyYy9hcGlzL2luZm8vYXBpXCJcbmltcG9ydCBCTiBmcm9tIFwiYm4uanNcIlxuaW1wb3J0IHsgUGVlcnNQYXJhbXMsIFBlZXJzUmVzcG9uc2UgfSBmcm9tIFwic3JjL2NvbW1vblwiXG5pbXBvcnQgV1MgZnJvbSBcImplc3Qtd2Vic29ja2V0LW1vY2tcIlxuXG4vLyBjb25zdCBzZXJ2ZXIgPSBuZXcgV1MoXCJ3czovL2xvY2FsaG9zdDoxMjM0XCIpXG5cbi8vIGNvbnN0IGNsaWVudCA9IG5ldyBTb2NrZXQoXCJ3czovL2xvY2FsaG9zdDoxMjM0XCIpXG5kZXNjcmliZShcIlNvY2tldFwiLCBhc3luYyAoKTogUHJvbWlzZTx2b2lkPiA9PiB7XG4gIC8vIGF3YWl0IHNlcnZlci5jb25uZWN0ZWQgLy8gd2FpdCBmb3IgdGhlIHNlcnZlciB0byBoYXZlIGVzdGFibGlzaGVkIHRoZSBjb25uZWN0aW9uXG5cbiAgLy8gdGhlIG1vY2sgd2Vic29ja2V0IHNlcnZlciB3aWxsIHJlY29yZCBhbGwgdGhlIG1lc3NhZ2VzIGl0IHJlY2VpdmVzXG4gIC8vIGNsaWVudC5zZW5kKFwiaGVsbG9cIilcblxuICAvLyB0aGUgbW9jayB3ZWJzb2NrZXQgc2VydmVyIGNhbiBhbHNvIHNlbmQgbWVzc2FnZXMgdG8gYWxsIGNvbm5lY3RlZCBjbGllbnRzXG4gIC8vIHNlcnZlci5zZW5kKFwiaGVsbG8gZXZlcnlvbmVcIilcblxuICAvLyAuLi5zaW11bGF0ZSBhbiBlcnJvciBhbmQgY2xvc2UgdGhlIGNvbm5lY3Rpb25cbiAgLy8gc2VydmVyLmVycm9yKClcblxuICAvLyAuLi5vciBncmFjZWZ1bGx5IGNsb3NlIHRoZSBjb25uZWN0aW9uXG4gIC8vIHNlcnZlci5jbG9zZSgpXG5cbiAgLy8gVGhlIFdTIGNsYXNzIGFsc28gaGFzIGEgc3RhdGljIFwiY2xlYW5cIiBtZXRob2QgdG8gZ3JhY2VmdWxseSBjbG9zZSBhbGwgb3BlbiBjb25uZWN0aW9ucyxcbiAgLy8gcGFydGljdWxhcmx5IHVzZWZ1bCB0byByZXNldCB0aGUgZW52aXJvbm1lbnQgYmV0d2VlbiB0ZXN0IHJ1bnMuXG4gIC8vIFdTLmNsZWFuKClcblxuICB0ZXN0KFwiZm9vYmFyXCIsIGFzeW5jICgpOiBQcm9taXNlPHZvaWQ+ID0+IHtcbiAgICBjb25zdCBzZXJ2ZXI6IFdTID0gbmV3IFdTKFwid3M6Ly9sb2NhbGhvc3Q6MTIzNFwiKVxuICAgIC8vIGNvbnNvbGUubG9nKHNlcnZlcilcbiAgICBjb25zdCBjaWVudCA9IG5ldyBXZWJTb2NrZXQoXCJ3czovL2xvY2FsaG9zdDoxMjM0XCIpXG4gICAgLy8gY29uc3QgY2xpZW50ID0gbmV3IFNvY2tldChcIndzOi8vbG9jYWxob3N0OjEyMzQvXCIpXG4gICAgY29uc29sZS5sb2coY2llbnQpXG5cbiAgICAvLyBhd2FpdCBzZXJ2ZXIuY29ubmVjdGVkXG4gICAgLy8gY2xpZW50LnNlbmQoXCJoZWxsb1wiKVxuICAgIC8vIGF3YWl0IGV4cGVjdChzZXJ2ZXIpLnRvUmVjZWl2ZU1lc3NhZ2UoXCJoZWxsb1wiKVxuICAgIC8vIGV4cGVjdChzZXJ2ZXIpLnRvSGF2ZVJlY2VpdmVkTWVzc2FnZXMoW1wiaGVsbG9cIl0pXG4gIH0pXG59KVxuIl19