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
const jest_mock_axios_1 = __importDefault(require("jest-mock-axios"));
const src_1 = require("src");
const api_1 = require("../../../src/apis/health/api");
describe("Health", () => {
    const ip = "127.0.0.1";
    const port = 9650;
    const protocol = "https";
    const avalanche = new src_1.Avalanche(ip, port, protocol, 12345, undefined, undefined, undefined, true);
    let health;
    beforeAll(() => {
        health = new api_1.HealthAPI(avalanche);
    });
    afterEach(() => {
        jest_mock_axios_1.default.reset();
    });
    test("health", () => __awaiter(void 0, void 0, void 0, function* () {
        const result = health.health();
        const payload = {
            result: {
                checks: {
                    C: {
                        message: [Object],
                        timestamp: "2021-09-29T15:31:20.274427-07:00",
                        duration: 275539,
                        contiguousFailures: 0,
                        timeOfFirstFailure: null
                    },
                    P: {
                        message: [Object],
                        timestamp: "2021-09-29T15:31:20.274508-07:00",
                        duration: 14576,
                        contiguousFailures: 0,
                        timeOfFirstFailure: null
                    },
                    X: {
                        message: [Object],
                        timestamp: "2021-09-29T15:31:20.274529-07:00",
                        duration: 4563,
                        contiguousFailures: 0,
                        timeOfFirstFailure: null
                    },
                    isBootstrapped: {
                        timestamp: "2021-09-29T15:31:19.448314-07:00",
                        duration: 392,
                        contiguousFailures: 0,
                        timeOfFirstFailure: null
                    },
                    network: {
                        message: [Object],
                        timestamp: "2021-09-29T15:31:19.448311-07:00",
                        duration: 4866,
                        contiguousFailures: 0,
                        timeOfFirstFailure: null
                    },
                    router: {
                        message: [Object],
                        timestamp: "2021-09-29T15:31:19.448452-07:00",
                        duration: 3932,
                        contiguousFailures: 0,
                        timeOfFirstFailure: null
                    }
                },
                healthy: true
            }
        };
        const responseObj = {
            data: payload
        };
        jest_mock_axios_1.default.mockResponse(responseObj);
        const response = yield result;
        expect(jest_mock_axios_1.default.request).toHaveBeenCalledTimes(1);
        expect(response).toBe(payload.result);
    }));
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYXBpLnRlc3QuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi90ZXN0cy9hcGlzL2hlYWx0aC9hcGkudGVzdC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7Ozs7OztBQUFBLHNFQUF1QztBQUV2Qyw2QkFBK0I7QUFDL0Isc0RBQXdEO0FBSXhELFFBQVEsQ0FBQyxRQUFRLEVBQUUsR0FBUyxFQUFFO0lBQzVCLE1BQU0sRUFBRSxHQUFXLFdBQVcsQ0FBQTtJQUM5QixNQUFNLElBQUksR0FBVyxJQUFJLENBQUE7SUFDekIsTUFBTSxRQUFRLEdBQVcsT0FBTyxDQUFBO0lBQ2hDLE1BQU0sU0FBUyxHQUFjLElBQUksZUFBUyxDQUN4QyxFQUFFLEVBQ0YsSUFBSSxFQUNKLFFBQVEsRUFDUixLQUFLLEVBQ0wsU0FBUyxFQUNULFNBQVMsRUFDVCxTQUFTLEVBQ1QsSUFBSSxDQUNMLENBQUE7SUFDRCxJQUFJLE1BQWlCLENBQUE7SUFFckIsU0FBUyxDQUFDLEdBQVMsRUFBRTtRQUNuQixNQUFNLEdBQUcsSUFBSSxlQUFTLENBQUMsU0FBUyxDQUFDLENBQUE7SUFDbkMsQ0FBQyxDQUFDLENBQUE7SUFFRixTQUFTLENBQUMsR0FBUyxFQUFFO1FBQ25CLHlCQUFTLENBQUMsS0FBSyxFQUFFLENBQUE7SUFDbkIsQ0FBQyxDQUFDLENBQUE7SUFFRixJQUFJLENBQUMsUUFBUSxFQUFFLEdBQXdCLEVBQUU7UUFDdkMsTUFBTSxNQUFNLEdBQTRCLE1BQU0sQ0FBQyxNQUFNLEVBQUUsQ0FBQTtRQUN2RCxNQUFNLE9BQU8sR0FBUTtZQUNuQixNQUFNLEVBQUU7Z0JBQ04sTUFBTSxFQUFFO29CQUNOLENBQUMsRUFBRTt3QkFDRCxPQUFPLEVBQUUsQ0FBQyxNQUFNLENBQUM7d0JBQ2pCLFNBQVMsRUFBRSxrQ0FBa0M7d0JBQzdDLFFBQVEsRUFBRSxNQUFNO3dCQUNoQixrQkFBa0IsRUFBRSxDQUFDO3dCQUNyQixrQkFBa0IsRUFBRSxJQUFJO3FCQUN6QjtvQkFDRCxDQUFDLEVBQUU7d0JBQ0QsT0FBTyxFQUFFLENBQUMsTUFBTSxDQUFDO3dCQUNqQixTQUFTLEVBQUUsa0NBQWtDO3dCQUM3QyxRQUFRLEVBQUUsS0FBSzt3QkFDZixrQkFBa0IsRUFBRSxDQUFDO3dCQUNyQixrQkFBa0IsRUFBRSxJQUFJO3FCQUN6QjtvQkFDRCxDQUFDLEVBQUU7d0JBQ0QsT0FBTyxFQUFFLENBQUMsTUFBTSxDQUFDO3dCQUNqQixTQUFTLEVBQUUsa0NBQWtDO3dCQUM3QyxRQUFRLEVBQUUsSUFBSTt3QkFDZCxrQkFBa0IsRUFBRSxDQUFDO3dCQUNyQixrQkFBa0IsRUFBRSxJQUFJO3FCQUN6QjtvQkFDRCxjQUFjLEVBQUU7d0JBQ2QsU0FBUyxFQUFFLGtDQUFrQzt3QkFDN0MsUUFBUSxFQUFFLEdBQUc7d0JBQ2Isa0JBQWtCLEVBQUUsQ0FBQzt3QkFDckIsa0JBQWtCLEVBQUUsSUFBSTtxQkFDekI7b0JBQ0QsT0FBTyxFQUFFO3dCQUNQLE9BQU8sRUFBRSxDQUFDLE1BQU0sQ0FBQzt3QkFDakIsU0FBUyxFQUFFLGtDQUFrQzt3QkFDN0MsUUFBUSxFQUFFLElBQUk7d0JBQ2Qsa0JBQWtCLEVBQUUsQ0FBQzt3QkFDckIsa0JBQWtCLEVBQUUsSUFBSTtxQkFDekI7b0JBQ0QsTUFBTSxFQUFFO3dCQUNOLE9BQU8sRUFBRSxDQUFDLE1BQU0sQ0FBQzt3QkFDakIsU0FBUyxFQUFFLGtDQUFrQzt3QkFDN0MsUUFBUSxFQUFFLElBQUk7d0JBQ2Qsa0JBQWtCLEVBQUUsQ0FBQzt3QkFDckIsa0JBQWtCLEVBQUUsSUFBSTtxQkFDekI7aUJBQ0Y7Z0JBQ0QsT0FBTyxFQUFFLElBQUk7YUFDZDtTQUNGLENBQUE7UUFDRCxNQUFNLFdBQVcsR0FBaUI7WUFDaEMsSUFBSSxFQUFFLE9BQU87U0FDZCxDQUFBO1FBRUQseUJBQVMsQ0FBQyxZQUFZLENBQUMsV0FBVyxDQUFDLENBQUE7UUFDbkMsTUFBTSxRQUFRLEdBQVEsTUFBTSxNQUFNLENBQUE7UUFFbEMsTUFBTSxDQUFDLHlCQUFTLENBQUMsT0FBTyxDQUFDLENBQUMscUJBQXFCLENBQUMsQ0FBQyxDQUFDLENBQUE7UUFDbEQsTUFBTSxDQUFDLFFBQVEsQ0FBQyxDQUFDLElBQUksQ0FBQyxPQUFPLENBQUMsTUFBTSxDQUFDLENBQUE7SUFDdkMsQ0FBQyxDQUFBLENBQUMsQ0FBQTtBQUNKLENBQUMsQ0FBQyxDQUFBIiwic291cmNlc0NvbnRlbnQiOlsiaW1wb3J0IG1vY2tBeGlvcyBmcm9tIFwiamVzdC1tb2NrLWF4aW9zXCJcblxuaW1wb3J0IHsgQXZhbGFuY2hlIH0gZnJvbSBcInNyY1wiXG5pbXBvcnQgeyBIZWFsdGhBUEkgfSBmcm9tIFwiLi4vLi4vLi4vc3JjL2FwaXMvaGVhbHRoL2FwaVwiXG5pbXBvcnQgeyBIZWFsdGhSZXNwb25zZSB9IGZyb20gXCIuLi8uLi8uLi9zcmMvYXBpcy9oZWFsdGgvaW50ZXJmYWNlc1wiXG5pbXBvcnQgeyBIdHRwUmVzcG9uc2UgfSBmcm9tIFwiamVzdC1tb2NrLWF4aW9zL2Rpc3QvbGliL21vY2stYXhpb3MtdHlwZXNcIlxuXG5kZXNjcmliZShcIkhlYWx0aFwiLCAoKTogdm9pZCA9PiB7XG4gIGNvbnN0IGlwOiBzdHJpbmcgPSBcIjEyNy4wLjAuMVwiXG4gIGNvbnN0IHBvcnQ6IG51bWJlciA9IDk2NTBcbiAgY29uc3QgcHJvdG9jb2w6IHN0cmluZyA9IFwiaHR0cHNcIlxuICBjb25zdCBhdmFsYW5jaGU6IEF2YWxhbmNoZSA9IG5ldyBBdmFsYW5jaGUoXG4gICAgaXAsXG4gICAgcG9ydCxcbiAgICBwcm90b2NvbCxcbiAgICAxMjM0NSxcbiAgICB1bmRlZmluZWQsXG4gICAgdW5kZWZpbmVkLFxuICAgIHVuZGVmaW5lZCxcbiAgICB0cnVlXG4gIClcbiAgbGV0IGhlYWx0aDogSGVhbHRoQVBJXG5cbiAgYmVmb3JlQWxsKCgpOiB2b2lkID0+IHtcbiAgICBoZWFsdGggPSBuZXcgSGVhbHRoQVBJKGF2YWxhbmNoZSlcbiAgfSlcblxuICBhZnRlckVhY2goKCk6IHZvaWQgPT4ge1xuICAgIG1vY2tBeGlvcy5yZXNldCgpXG4gIH0pXG5cbiAgdGVzdChcImhlYWx0aFwiLCBhc3luYyAoKTogUHJvbWlzZTx2b2lkPiA9PiB7XG4gICAgY29uc3QgcmVzdWx0OiBQcm9taXNlPEhlYWx0aFJlc3BvbnNlPiA9IGhlYWx0aC5oZWFsdGgoKVxuICAgIGNvbnN0IHBheWxvYWQ6IGFueSA9IHtcbiAgICAgIHJlc3VsdDoge1xuICAgICAgICBjaGVja3M6IHtcbiAgICAgICAgICBDOiB7XG4gICAgICAgICAgICBtZXNzYWdlOiBbT2JqZWN0XSxcbiAgICAgICAgICAgIHRpbWVzdGFtcDogXCIyMDIxLTA5LTI5VDE1OjMxOjIwLjI3NDQyNy0wNzowMFwiLFxuICAgICAgICAgICAgZHVyYXRpb246IDI3NTUzOSxcbiAgICAgICAgICAgIGNvbnRpZ3VvdXNGYWlsdXJlczogMCxcbiAgICAgICAgICAgIHRpbWVPZkZpcnN0RmFpbHVyZTogbnVsbFxuICAgICAgICAgIH0sXG4gICAgICAgICAgUDoge1xuICAgICAgICAgICAgbWVzc2FnZTogW09iamVjdF0sXG4gICAgICAgICAgICB0aW1lc3RhbXA6IFwiMjAyMS0wOS0yOVQxNTozMToyMC4yNzQ1MDgtMDc6MDBcIixcbiAgICAgICAgICAgIGR1cmF0aW9uOiAxNDU3NixcbiAgICAgICAgICAgIGNvbnRpZ3VvdXNGYWlsdXJlczogMCxcbiAgICAgICAgICAgIHRpbWVPZkZpcnN0RmFpbHVyZTogbnVsbFxuICAgICAgICAgIH0sXG4gICAgICAgICAgWDoge1xuICAgICAgICAgICAgbWVzc2FnZTogW09iamVjdF0sXG4gICAgICAgICAgICB0aW1lc3RhbXA6IFwiMjAyMS0wOS0yOVQxNTozMToyMC4yNzQ1MjktMDc6MDBcIixcbiAgICAgICAgICAgIGR1cmF0aW9uOiA0NTYzLFxuICAgICAgICAgICAgY29udGlndW91c0ZhaWx1cmVzOiAwLFxuICAgICAgICAgICAgdGltZU9mRmlyc3RGYWlsdXJlOiBudWxsXG4gICAgICAgICAgfSxcbiAgICAgICAgICBpc0Jvb3RzdHJhcHBlZDoge1xuICAgICAgICAgICAgdGltZXN0YW1wOiBcIjIwMjEtMDktMjlUMTU6MzE6MTkuNDQ4MzE0LTA3OjAwXCIsXG4gICAgICAgICAgICBkdXJhdGlvbjogMzkyLFxuICAgICAgICAgICAgY29udGlndW91c0ZhaWx1cmVzOiAwLFxuICAgICAgICAgICAgdGltZU9mRmlyc3RGYWlsdXJlOiBudWxsXG4gICAgICAgICAgfSxcbiAgICAgICAgICBuZXR3b3JrOiB7XG4gICAgICAgICAgICBtZXNzYWdlOiBbT2JqZWN0XSxcbiAgICAgICAgICAgIHRpbWVzdGFtcDogXCIyMDIxLTA5LTI5VDE1OjMxOjE5LjQ0ODMxMS0wNzowMFwiLFxuICAgICAgICAgICAgZHVyYXRpb246IDQ4NjYsXG4gICAgICAgICAgICBjb250aWd1b3VzRmFpbHVyZXM6IDAsXG4gICAgICAgICAgICB0aW1lT2ZGaXJzdEZhaWx1cmU6IG51bGxcbiAgICAgICAgICB9LFxuICAgICAgICAgIHJvdXRlcjoge1xuICAgICAgICAgICAgbWVzc2FnZTogW09iamVjdF0sXG4gICAgICAgICAgICB0aW1lc3RhbXA6IFwiMjAyMS0wOS0yOVQxNTozMToxOS40NDg0NTItMDc6MDBcIixcbiAgICAgICAgICAgIGR1cmF0aW9uOiAzOTMyLFxuICAgICAgICAgICAgY29udGlndW91c0ZhaWx1cmVzOiAwLFxuICAgICAgICAgICAgdGltZU9mRmlyc3RGYWlsdXJlOiBudWxsXG4gICAgICAgICAgfVxuICAgICAgICB9LFxuICAgICAgICBoZWFsdGh5OiB0cnVlXG4gICAgICB9XG4gICAgfVxuICAgIGNvbnN0IHJlc3BvbnNlT2JqOiBIdHRwUmVzcG9uc2UgPSB7XG4gICAgICBkYXRhOiBwYXlsb2FkXG4gICAgfVxuXG4gICAgbW9ja0F4aW9zLm1vY2tSZXNwb25zZShyZXNwb25zZU9iailcbiAgICBjb25zdCByZXNwb25zZTogYW55ID0gYXdhaXQgcmVzdWx0XG5cbiAgICBleHBlY3QobW9ja0F4aW9zLnJlcXVlc3QpLnRvSGF2ZUJlZW5DYWxsZWRUaW1lcygxKVxuICAgIGV4cGVjdChyZXNwb25zZSkudG9CZShwYXlsb2FkLnJlc3VsdClcbiAgfSlcbn0pXG4iXX0=