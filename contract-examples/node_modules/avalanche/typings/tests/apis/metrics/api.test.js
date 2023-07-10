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
const api_1 = require("../../../src/apis/metrics/api");
describe("Metrics", () => {
    const ip = "127.0.0.1";
    const port = 9650;
    const protocol = "https";
    const avalanche = new src_1.Avalanche(ip, port, protocol, 12345, undefined, undefined, undefined, true);
    let metrics;
    beforeAll(() => {
        metrics = new api_1.MetricsAPI(avalanche);
    });
    afterEach(() => {
        jest_mock_axios_1.default.reset();
    });
    test("getMetrics", () => __awaiter(void 0, void 0, void 0, function* () {
        const result = metrics.getMetrics();
        const payload = `
              gecko_timestamp_handler_get_failed_bucket{le="100"} 0
              gecko_timestamp_handler_get_failed_bucket{le="1000"} 0
              gecko_timestamp_handler_get_failed_bucket{le="10000"} 0
              gecko_timestamp_handler_get_failed_bucket{le="100000"} 0
              gecko_timestamp_handler_get_failed_bucket{le="1e+06"} 0
              gecko_timestamp_handler_get_failed_bucket{le="1e+07"} 0
              gecko_timestamp_handler_get_failed_bucket{le="1e+08"} 0
              gecko_timestamp_handler_get_failed_bucket{le="1e+09"} 0
              gecko_timestamp_handler_get_failed_bucket{le="+Inf"} 0
        `;
        const responseObj = {
            data: payload
        };
        jest_mock_axios_1.default.mockResponse(responseObj);
        const response = yield result;
        expect(jest_mock_axios_1.default.request).toHaveBeenCalledTimes(1);
        expect(response).toBe(payload);
    }));
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYXBpLnRlc3QuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi90ZXN0cy9hcGlzL21ldHJpY3MvYXBpLnRlc3QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7Ozs7QUFBQSxzRUFBdUM7QUFFdkMsNkJBQStCO0FBQy9CLHVEQUEwRDtBQUUxRCxRQUFRLENBQUMsU0FBUyxFQUFFLEdBQVMsRUFBRTtJQUM3QixNQUFNLEVBQUUsR0FBVyxXQUFXLENBQUE7SUFDOUIsTUFBTSxJQUFJLEdBQVcsSUFBSSxDQUFBO0lBQ3pCLE1BQU0sUUFBUSxHQUFXLE9BQU8sQ0FBQTtJQUVoQyxNQUFNLFNBQVMsR0FBYyxJQUFJLGVBQVMsQ0FDeEMsRUFBRSxFQUNGLElBQUksRUFDSixRQUFRLEVBQ1IsS0FBSyxFQUNMLFNBQVMsRUFDVCxTQUFTLEVBQ1QsU0FBUyxFQUNULElBQUksQ0FDTCxDQUFBO0lBQ0QsSUFBSSxPQUFtQixDQUFBO0lBRXZCLFNBQVMsQ0FBQyxHQUFTLEVBQUU7UUFDbkIsT0FBTyxHQUFHLElBQUksZ0JBQVUsQ0FBQyxTQUFTLENBQUMsQ0FBQTtJQUNyQyxDQUFDLENBQUMsQ0FBQTtJQUVGLFNBQVMsQ0FBQyxHQUFTLEVBQUU7UUFDbkIseUJBQVMsQ0FBQyxLQUFLLEVBQUUsQ0FBQTtJQUNuQixDQUFDLENBQUMsQ0FBQTtJQUVGLElBQUksQ0FBQyxZQUFZLEVBQUUsR0FBd0IsRUFBRTtRQUMzQyxNQUFNLE1BQU0sR0FBb0IsT0FBTyxDQUFDLFVBQVUsRUFBRSxDQUFBO1FBQ3BELE1BQU0sT0FBTyxHQUFXOzs7Ozs7Ozs7O1NBVW5CLENBQUE7UUFDTCxNQUFNLFdBQVcsR0FBaUI7WUFDaEMsSUFBSSxFQUFFLE9BQU87U0FDZCxDQUFBO1FBRUQseUJBQVMsQ0FBQyxZQUFZLENBQUMsV0FBVyxDQUFDLENBQUE7UUFDbkMsTUFBTSxRQUFRLEdBQVcsTUFBTSxNQUFNLENBQUE7UUFFckMsTUFBTSxDQUFDLHlCQUFTLENBQUMsT0FBTyxDQUFDLENBQUMscUJBQXFCLENBQUMsQ0FBQyxDQUFDLENBQUE7UUFDbEQsTUFBTSxDQUFDLFFBQVEsQ0FBQyxDQUFDLElBQUksQ0FBQyxPQUFPLENBQUMsQ0FBQTtJQUNoQyxDQUFDLENBQUEsQ0FBQyxDQUFBO0FBQ0osQ0FBQyxDQUFDLENBQUEiLCJzb3VyY2VzQ29udGVudCI6WyJpbXBvcnQgbW9ja0F4aW9zIGZyb20gXCJqZXN0LW1vY2stYXhpb3NcIlxuaW1wb3J0IHsgSHR0cFJlc3BvbnNlIH0gZnJvbSBcImplc3QtbW9jay1heGlvcy9kaXN0L2xpYi9tb2NrLWF4aW9zLXR5cGVzXCJcbmltcG9ydCB7IEF2YWxhbmNoZSB9IGZyb20gXCJzcmNcIlxuaW1wb3J0IHsgTWV0cmljc0FQSSB9IGZyb20gXCIuLi8uLi8uLi9zcmMvYXBpcy9tZXRyaWNzL2FwaVwiXG5cbmRlc2NyaWJlKFwiTWV0cmljc1wiLCAoKTogdm9pZCA9PiB7XG4gIGNvbnN0IGlwOiBzdHJpbmcgPSBcIjEyNy4wLjAuMVwiXG4gIGNvbnN0IHBvcnQ6IG51bWJlciA9IDk2NTBcbiAgY29uc3QgcHJvdG9jb2w6IHN0cmluZyA9IFwiaHR0cHNcIlxuXG4gIGNvbnN0IGF2YWxhbmNoZTogQXZhbGFuY2hlID0gbmV3IEF2YWxhbmNoZShcbiAgICBpcCxcbiAgICBwb3J0LFxuICAgIHByb3RvY29sLFxuICAgIDEyMzQ1LFxuICAgIHVuZGVmaW5lZCxcbiAgICB1bmRlZmluZWQsXG4gICAgdW5kZWZpbmVkLFxuICAgIHRydWVcbiAgKVxuICBsZXQgbWV0cmljczogTWV0cmljc0FQSVxuXG4gIGJlZm9yZUFsbCgoKTogdm9pZCA9PiB7XG4gICAgbWV0cmljcyA9IG5ldyBNZXRyaWNzQVBJKGF2YWxhbmNoZSlcbiAgfSlcblxuICBhZnRlckVhY2goKCk6IHZvaWQgPT4ge1xuICAgIG1vY2tBeGlvcy5yZXNldCgpXG4gIH0pXG5cbiAgdGVzdChcImdldE1ldHJpY3NcIiwgYXN5bmMgKCk6IFByb21pc2U8dm9pZD4gPT4ge1xuICAgIGNvbnN0IHJlc3VsdDogUHJvbWlzZTxzdHJpbmc+ID0gbWV0cmljcy5nZXRNZXRyaWNzKClcbiAgICBjb25zdCBwYXlsb2FkOiBzdHJpbmcgPSBgXG4gICAgICAgICAgICAgIGdlY2tvX3RpbWVzdGFtcF9oYW5kbGVyX2dldF9mYWlsZWRfYnVja2V0e2xlPVwiMTAwXCJ9IDBcbiAgICAgICAgICAgICAgZ2Vja29fdGltZXN0YW1wX2hhbmRsZXJfZ2V0X2ZhaWxlZF9idWNrZXR7bGU9XCIxMDAwXCJ9IDBcbiAgICAgICAgICAgICAgZ2Vja29fdGltZXN0YW1wX2hhbmRsZXJfZ2V0X2ZhaWxlZF9idWNrZXR7bGU9XCIxMDAwMFwifSAwXG4gICAgICAgICAgICAgIGdlY2tvX3RpbWVzdGFtcF9oYW5kbGVyX2dldF9mYWlsZWRfYnVja2V0e2xlPVwiMTAwMDAwXCJ9IDBcbiAgICAgICAgICAgICAgZ2Vja29fdGltZXN0YW1wX2hhbmRsZXJfZ2V0X2ZhaWxlZF9idWNrZXR7bGU9XCIxZSswNlwifSAwXG4gICAgICAgICAgICAgIGdlY2tvX3RpbWVzdGFtcF9oYW5kbGVyX2dldF9mYWlsZWRfYnVja2V0e2xlPVwiMWUrMDdcIn0gMFxuICAgICAgICAgICAgICBnZWNrb190aW1lc3RhbXBfaGFuZGxlcl9nZXRfZmFpbGVkX2J1Y2tldHtsZT1cIjFlKzA4XCJ9IDBcbiAgICAgICAgICAgICAgZ2Vja29fdGltZXN0YW1wX2hhbmRsZXJfZ2V0X2ZhaWxlZF9idWNrZXR7bGU9XCIxZSswOVwifSAwXG4gICAgICAgICAgICAgIGdlY2tvX3RpbWVzdGFtcF9oYW5kbGVyX2dldF9mYWlsZWRfYnVja2V0e2xlPVwiK0luZlwifSAwXG4gICAgICAgIGBcbiAgICBjb25zdCByZXNwb25zZU9iajogSHR0cFJlc3BvbnNlID0ge1xuICAgICAgZGF0YTogcGF5bG9hZFxuICAgIH1cblxuICAgIG1vY2tBeGlvcy5tb2NrUmVzcG9uc2UocmVzcG9uc2VPYmopXG4gICAgY29uc3QgcmVzcG9uc2U6IHN0cmluZyA9IGF3YWl0IHJlc3VsdFxuXG4gICAgZXhwZWN0KG1vY2tBeGlvcy5yZXF1ZXN0KS50b0hhdmVCZWVuQ2FsbGVkVGltZXMoMSlcbiAgICBleHBlY3QocmVzcG9uc2UpLnRvQmUocGF5bG9hZClcbiAgfSlcbn0pXG4iXX0=