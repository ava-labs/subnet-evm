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
describe("Index", () => {
    const ip = "127.0.0.1";
    const port = 9650;
    const protocol = "https";
    const avalanche = new src_1.Avalanche(ip, port, protocol, 12345);
    let index;
    const id = "6fXf5hncR8LXvwtM8iezFQBpK5cubV6y1dWgpJCcNyzGB1EzY";
    const bytes = "111115HRzXVDSeonLBcv6QdJkQFjPzPEobMWy7PyGuoheggsMCx73MVXZo2hJMEXXvR5gFFasTRJH36aVkLiWHtTTFcghyFTqjaHnBhdXTRiLaYcro3jpseqLAFVn3ngnAB47nebQiBBKmg3nFWKzQUDxMuE6uDGXgnGouDSaEKZxfKreoLHYNUxH56rgi5c8gKFYSDi8AWBgy26siwAWj6V8EgFnPVgm9pmKCfXio6BP7Bua4vrupoX8jRGqdrdkN12dqGAibJ78Rf44SSUXhEvJtPxAzjEGfiTyAm5BWFqPdheKN72HyrBBtwC6y7wG6suHngZ1PMBh93Ubkbt8jjjGoEgs5NjpasJpE8YA9ZMLTPeNZ6ELFxV99zA46wvkjAwYHGzegBXvzGU5pGPbg28iW3iKhLoYAnReysY4x3fBhjPBsags37Z9P3SqioVifVX4wwzxYqbV72u1AWZ4JNmsnhVDP196Gu99QTzmySGTVGP5ABNdZrngTRfmGTFCRbt9CHsgNbhgetkxbsEG7tySi3gFxMzGuJ2Npk2gnSr68LgtYdSHf48Ns";
    const timestamp = "2021-04-02T15:34:00.262979-07:00";
    const idx = "0";
    beforeAll(() => {
        index = avalanche.Index();
    });
    afterEach(() => {
        jest_mock_axios_1.default.reset();
    });
    test("getLastAccepted", () => __awaiter(void 0, void 0, void 0, function* () {
        const encoding = "hex";
        const baseurl = "/ext/index/X/tx";
        const respobj = {
            id,
            bytes,
            timestamp,
            encoding,
            idx
        };
        const result = index.getLastAccepted(encoding, baseurl);
        const payload = {
            result: respobj
        };
        const responseObj = {
            data: payload
        };
        jest_mock_axios_1.default.mockResponse(responseObj);
        const response = yield result;
        expect(jest_mock_axios_1.default.request).toHaveBeenCalledTimes(1);
        expect(JSON.stringify(response)).toBe(JSON.stringify(respobj));
    }));
    test("getContainerByIndex", () => __awaiter(void 0, void 0, void 0, function* () {
        const encoding = "hex";
        const baseurl = "/ext/index/X/tx";
        const respobj = {
            id,
            bytes,
            timestamp,
            encoding,
            idx
        };
        const result = index.getContainerByIndex(idx, encoding, baseurl);
        const payload = {
            result: respobj
        };
        const responseObj = {
            data: payload
        };
        jest_mock_axios_1.default.mockResponse(responseObj);
        const response = yield result;
        expect(jest_mock_axios_1.default.request).toHaveBeenCalledTimes(1);
        expect(JSON.stringify(response)).toBe(JSON.stringify(respobj));
    }));
    test("getContainerByID", () => __awaiter(void 0, void 0, void 0, function* () {
        const encoding = "hex";
        const baseurl = "/ext/index/X/tx";
        const respobj = {
            id,
            bytes,
            timestamp,
            encoding,
            idx
        };
        const result = index.getContainerByIndex(id, encoding, baseurl);
        const payload = {
            result: respobj
        };
        const responseObj = {
            data: payload
        };
        jest_mock_axios_1.default.mockResponse(responseObj);
        const response = yield result;
        expect(jest_mock_axios_1.default.request).toHaveBeenCalledTimes(1);
        expect(JSON.stringify(response)).toBe(JSON.stringify(respobj));
    }));
    test("getContainerRange", () => __awaiter(void 0, void 0, void 0, function* () {
        const startIndex = 0;
        const numToFetch = 100;
        const encoding = "hex";
        const baseurl = "/ext/index/X/tx";
        const respobj = {
            id,
            bytes,
            timestamp,
            encoding,
            idx
        };
        const result = index.getContainerRange(startIndex, numToFetch, encoding, baseurl);
        const payload = {
            result: respobj
        };
        const responseObj = {
            data: payload
        };
        jest_mock_axios_1.default.mockResponse(responseObj);
        const response = yield result;
        expect(jest_mock_axios_1.default.request).toHaveBeenCalledTimes(1);
        expect(JSON.stringify(response)).toBe(JSON.stringify(respobj));
    }));
    test("getIndex", () => __awaiter(void 0, void 0, void 0, function* () {
        const encoding = "hex";
        const baseurl = "/ext/index/X/tx";
        const result = index.getIndex(id, encoding, baseurl);
        const payload = {
            result: {
                index: "0"
            }
        };
        const responseObj = {
            data: payload
        };
        jest_mock_axios_1.default.mockResponse(responseObj);
        const response = yield result;
        expect(jest_mock_axios_1.default.request).toHaveBeenCalledTimes(1);
        expect(response).toBe("0");
    }));
    test("isAccepted", () => __awaiter(void 0, void 0, void 0, function* () {
        const encoding = "hex";
        const baseurl = "/ext/index/X/tx";
        const result = index.isAccepted(id, encoding, baseurl);
        const payload = {
            result: true
        };
        const responseObj = {
            data: payload
        };
        jest_mock_axios_1.default.mockResponse(responseObj);
        const response = yield result;
        expect(jest_mock_axios_1.default.request).toHaveBeenCalledTimes(1);
        expect(response).toBe(true);
    }));
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYXBpLnRlc3QuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi90ZXN0cy9hcGlzL2luZGV4L2FwaS50ZXN0LnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7Ozs7Ozs7Ozs7O0FBQUEsc0VBQXVDO0FBRXZDLDZCQUErQjtBQUkvQixRQUFRLENBQUMsT0FBTyxFQUFFLEdBQUcsRUFBRTtJQUNyQixNQUFNLEVBQUUsR0FBVyxXQUFXLENBQUE7SUFDOUIsTUFBTSxJQUFJLEdBQVcsSUFBSSxDQUFBO0lBQ3pCLE1BQU0sUUFBUSxHQUFXLE9BQU8sQ0FBQTtJQUVoQyxNQUFNLFNBQVMsR0FBYyxJQUFJLGVBQVMsQ0FBQyxFQUFFLEVBQUUsSUFBSSxFQUFFLFFBQVEsRUFBRSxLQUFLLENBQUMsQ0FBQTtJQUNyRSxJQUFJLEtBQWUsQ0FBQTtJQUVuQixNQUFNLEVBQUUsR0FBVyxtREFBbUQsQ0FBQTtJQUN0RSxNQUFNLEtBQUssR0FDVCw0akJBQTRqQixDQUFBO0lBQzlqQixNQUFNLFNBQVMsR0FBVyxrQ0FBa0MsQ0FBQTtJQUM1RCxNQUFNLEdBQUcsR0FBVyxHQUFHLENBQUE7SUFFdkIsU0FBUyxDQUFDLEdBQUcsRUFBRTtRQUNiLEtBQUssR0FBRyxTQUFTLENBQUMsS0FBSyxFQUFFLENBQUE7SUFDM0IsQ0FBQyxDQUFDLENBQUE7SUFFRixTQUFTLENBQUMsR0FBRyxFQUFFO1FBQ2IseUJBQVMsQ0FBQyxLQUFLLEVBQUUsQ0FBQTtJQUNuQixDQUFDLENBQUMsQ0FBQTtJQUVGLElBQUksQ0FBQyxpQkFBaUIsRUFBRSxHQUFTLEVBQUU7UUFDakMsTUFBTSxRQUFRLEdBQVcsS0FBSyxDQUFBO1FBQzlCLE1BQU0sT0FBTyxHQUFXLGlCQUFpQixDQUFBO1FBQ3pDLE1BQU0sT0FBTyxHQUFHO1lBQ2QsRUFBRTtZQUNGLEtBQUs7WUFDTCxTQUFTO1lBQ1QsUUFBUTtZQUNSLEdBQUc7U0FDSixDQUFBO1FBQ0QsTUFBTSxNQUFNLEdBQW9CLEtB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aEMsSUFBSSxFQUFFLE9BQU87U0FDZCxDQUFBO1FBRUQseUJBQVMsQ0FBQyxZQUFZLENBQUMsV0FBVyxDQUFDLENBQUE7UUFDbkMsTUFBTSxRQUFRLEdBQVcsTUFBTSxNQUFNLENBQUE7UUFFckMsTUFBTSxDQUFDLHlCQUFTLENBQUMsT0FBTyxDQUFDLENBQUMscUJBQXFCLENBQUMsQ0FBQyxDQUFDLENBQUE7UUFDbEQsTUFBTSxDQUFDLFFBQVEsQ0FBQyxDQUFDLElBQUksQ0FBQyxHQUFHLENBQUMsQ0FBQTtJQUM1QixDQUFDLENBQUEsQ0FBQyxDQUFBO0lBRUYsSUFBSSxDQUFDLFlBQVksRUFBRSxHQUFTLEVBQUU7UUFDNUIsTUFBTSxRQUFRLEdBQVcsS0FBSyxDQUFBO1FBQzlCLE1BQU0sT0FBTyxHQUFXLGlCQUFpQixDQUFBO1FBQ3pDLE1BQU0sTUFBTSxHQUFnQyxLQUFLLENBQUMsVUFBVSxDQUMxRCxFQUFFLEVBQ0YsUUFBUSxFQUNSLE9BQU8sQ0FDUixDQUFBO1FBRUQsTUFBTSxPQUFPLEdBQVc7WUFDdEIsTUFBTSxFQUFFLElBQUk7U0FDYixDQUFBO1FBQ0QsTUFBTSxXQUFXLEdBQWlCO1lBQ2hDLElBQUksRUFBRSxPQUFPO1NBQ2QsQ0FBQTtRQUVELHlCQUFTLENBQUMsWUFBWSxDQUFDLFdBQVcsQ0FBQyxDQUFBO1FBQ25DLE1BQU0sUUFBUSxHQUF1QixNQUFNLE1BQU0sQ0FBQTtRQUVqRCxNQUFNLENBQUMseUJBQVMsQ0FBQyxPQUFPLENBQUMsQ0FBQyxxQkFBcUIsQ0FBQyxDQUFDLENBQUMsQ0FBQTtRQUNsRCxNQUFNLENBQUMsUUFBUSxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO0lBQzdCLENBQUMsQ0FBQSxDQUFDLENBQUE7QUFDSixDQUFDLENBQUMsQ0FBQSIsInNvdXJjZXNDb250ZW50IjpbImltcG9ydCBtb2NrQXhpb3MgZnJvbSBcImplc3QtbW9jay1heGlvc1wiXG5pbXBvcnQgeyBIdHRwUmVzcG9uc2UgfSBmcm9tIFwiamVzdC1tb2NrLWF4aW9zL2Rpc3QvbGliL21vY2stYXhpb3MtdHlwZXNcIlxuaW1wb3J0IHsgQXZhbGFuY2hlIH0gZnJvbSBcInNyY1wiXG5pbXBvcnQgeyBJc0FjY2VwdGVkUmVzcG9uc2UgfSBmcm9tIFwic3JjL2FwaXMvaW5kZXgvaW50ZXJmYWNlc1wiXG5pbXBvcnQgeyBJbmRleEFQSSB9IGZyb20gXCIuLi8uLi8uLi9zcmMvYXBpcy9pbmRleC9hcGlcIlxuXG5kZXNjcmliZShcIkluZGV4XCIsICgpID0+IHtcbiAgY29uc3QgaXA6IHN0cmluZyA9IFwiMTI3LjAuMC4xXCJcbiAgY29uc3QgcG9ydDogbnVtYmVyID0gOTY1MFxuICBjb25zdCBwcm90b2NvbDogc3RyaW5nID0gXCJodHRwc1wiXG5cbiAgY29uc3QgYXZhbGFuY2hlOiBBdmFsYW5jaGUgPSBuZXcgQXZhbGFuY2hlKGlwLCBwb3J0LCBwcm90b2NvbCwgMTIzNDUpXG4gIGxldCBpbmRleDogSW5kZXhBUElcblxuICBjb25zdCBpZDogc3RyaW5nID0gXCI2ZlhmNWhuY1I4TFh2d3RNOGllekZRQnBLNWN1YlY2eTFkV2dwSkNjTnl6R0IxRXpZXCJcbiAgY29uc3QgYnl0ZXM6IHN0cmluZyA9XG4gICAgXCIxMTExMTVIUnpYVkRTZW9uTEJjdjZRZEprUUZqUHpQRW9iTVd5N1B5R3VvaGVnZ3NNQ3g3M01WWFpvMmhKTUVYWHZSNWdGRmFzVFJKSDM2YVZrTGlXSHRUVEZjZ2h5RlRxamFIbkJoZFhUUmlMYVljcm8zanBzZXFMQUZWbjNuZ25BQjQ3bmViUWlCQkttZzNuRldLelFVRHhNdUU2dURHWGduR291RFNhRUtaeGZLcmVvTEhZTlV4SDU2cmdpNWM4Z0tGWVNEaThBV0JneTI2c2l3QVdqNlY4RWdGblBWZ205cG1LQ2ZYaW82QlA3QnVhNHZydXBvWDhqUkdxZHJka04xMmRxR0FpYko3OFJmNDRTU1VYaEV2SnRQeEF6akVHZmlUeUFtNUJXRnFQZGhlS043Mkh5ckJCdHdDNnk3d0c2c3VIbmdaMVBNQmg5M1Via2J0OGpqakdvRWdzNU5qcGFzSnBFOFlBOVpNTFRQZU5aNkVMRnhWOTl6QTQ2d3ZrakF3WUhHemVnQlh2ekdVNXBHUGJnMjhpVzNpS2hMb1lBblJleXNZNHgzZkJoalBCc2FnczM3WjlQM1NxaW9WaWZWWDR3d3p4WXFiVjcydTFBV1o0Sk5tc25oVkRQMTk2R3U5OVFUem15U0dUVkdQNUFCTmRacm5nVFJmbUdURkNSYnQ5Q0hzZ05iaGdldGt4YnNFRzd0eVNpM2dGeE16R3VKMk5wazJnblNyNjhMZ3RZZFNIZjQ4TnNcIlxuICBjb25zdCB0aW1lc3RhbXA6IHN0cmluZyA9IFwiMjAyMS0wNC0wMlQxNTozNDowMC4yNjI5NzktMDc6MDBcIlxuICBjb25zdCBpZHg6IHN0cmluZyA9IFwiMFwiXG5cbiAgYmVmb3JlQWxsKCgpID0+IHtcbiAgICBpbmRleCA9IGF2YWxhbmNoZS5JbmRleCgpXG4gIH0pXG5cbiAgYWZ0ZXJFYWNoKCgpID0+IHtcbiAgICBtb2NrQXhpb3MucmVzZXQoKVxuICB9KVxuXG4gIHRlc3QoXCJnZXRMYXN0QWNjZXB0ZWRcIiwgYXN5bmMgKCkgPT4ge1xuICAgIGNvbnN0IGVuY29kaW5nOiBzdHJpbmcgPSBcImhleFwiXG4gICAgY29uc3QgYmFzZXVybDogc3RyaW5nID0gXCIvZXh0L2luZGV4L1gvdHhcIlxuICAgIGNvbnN0IHJlc3BvYmogPSB7XG4gICAgICBpZCxcbiAgICAgIGJ5dGVzLFxuICAgICAgdGltZXN0YW1wLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBpZHhcbiAgICB9XG4gICAgY29uc3QgcmVzdWx0OiBQcm9taXNlPG9iamVjdD4gPSBpbmRleC5nZXRMYXN0QWNjZXB0ZWQoZW5jb2RpbmcsIGJhc2V1cmwpXG5cbiAgICBjb25zdCBwYXlsb2FkOiBvYmplY3QgPSB7XG4gICAgICByZXN1bHQ6IHJlc3BvYmpcbiAgICB9XG4gICAgY29uc3QgcmVzcG9uc2VPYmo6IEh0dHBSZXNwb25zZSA9IHtcbiAgICAgIGRhdGE6IHBheWxvYWRcbiAgICB9XG5cbiAgICBtb2NrQXhpb3MubW9ja1Jlc3BvbnNlKHJlc3BvbnNlT2JqKVxuICAgIGNvbnN0IHJlc3BvbnNlOiBvYmplY3QgPSBhd2FpdCByZXN1bHRcbiAgICBleHBlY3QobW9ja0F4aW9zLnJlcXVlc3QpLnRvSGF2ZUJlZW5DYWxsZWRUaW1lcygxKVxuICAgIGV4cGVjdChKU09OLnN0cmluZ2lmeShyZXNwb25zZSkpLnRvQmUoSlNPTi5zdHJpbmdpZnkocmVzcG9iaikpXG4gIH0pXG5cbiAgdGVzdChcImdldENvbnRhaW5lckJ5SW5kZXhcIiwgYXN5bmMgKCkgPT4ge1xuICAgIGNvbnN0IGVuY29kaW5nOiBzdHJpbmcgPSBcImhleFwiXG4gICAgY29uc3QgYmFzZXVybDogc3RyaW5nID0gXCIvZXh0L2luZGV4L1gvdHhcIlxuICAgIGNvbnN0IHJlc3BvYmogPSB7XG4gICAgICBpZCxcbiAgICAgIGJ5dGVzLFxuICAgICAgdGltZXN0YW1wLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBpZHhcbiAgICB9XG4gICAgY29uc3QgcmVzdWx0OiBQcm9taXNlPG9iamVjdD4gPSBpbmRleC5nZXRDb250YWluZXJCeUluZGV4KFxuICAgICAgaWR4LFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBiYXNldXJsXG4gICAgKVxuXG4gICAgY29uc3QgcGF5bG9hZDogb2JqZWN0ID0ge1xuICAgICAgcmVzdWx0OiByZXNwb2JqXG4gICAgfVxuICAgIGNvbnN0IHJlc3BvbnNlT2JqOiBIdHRwUmVzcG9uc2UgPSB7XG4gICAgICBkYXRhOiBwYXlsb2FkXG4gICAgfVxuXG4gICAgbW9ja0F4aW9zLm1vY2tSZXNwb25zZShyZXNwb25zZU9iailcbiAgICBjb25zdCByZXNwb25zZTogb2JqZWN0ID0gYXdhaXQgcmVzdWx0XG4gICAgZXhwZWN0KG1vY2tBeGlvcy5yZXF1ZXN0KS50b0hhdmVCZWVuQ2FsbGVkVGltZXMoMSlcbiAgICBleHBlY3QoSlNPTi5zdHJpbmdpZnkocmVzcG9uc2UpKS50b0JlKEpTT04uc3RyaW5naWZ5KHJlc3BvYmopKVxuICB9KVxuXG4gIHRlc3QoXCJnZXRDb250YWluZXJCeUlEXCIsIGFzeW5jICgpID0+IHtcbiAgICBjb25zdCBlbmNvZGluZzogc3RyaW5nID0gXCJoZXhcIlxuICAgIGNvbnN0IGJhc2V1cmw6IHN0cmluZyA9IFwiL2V4dC9pbmRleC9YL3R4XCJcbiAgICBjb25zdCByZXNwb2JqID0ge1xuICAgICAgaWQsXG4gICAgICBieXRlcyxcbiAgICAgIHRpbWVzdGFtcCxcbiAgICAgIGVuY29kaW5nLFxuICAgICAgaWR4XG4gICAgfVxuICAgIGNvbnN0IHJlc3VsdDogUHJvbWlzZTxvYmplY3Q+ID0gaW5kZXguZ2V0Q29udGFpbmVyQnlJbmRleChcbiAgICAgIGlkLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBiYXNldXJsXG4gICAgKVxuXG4gICAgY29uc3QgcGF5bG9hZDogb2JqZWN0ID0ge1xuICAgICAgcmVzdWx0OiByZXNwb2JqXG4gICAgfVxuICAgIGNvbnN0IHJlc3BvbnNlT2JqOiBIdHRwUmVzcG9uc2UgPSB7XG4gICAgICBkYXRhOiBwYXlsb2FkXG4gICAgfVxuXG4gICAgbW9ja0F4aW9zLm1vY2tSZXNwb25zZShyZXNwb25zZU9iailcbiAgICBjb25zdCByZXNwb25zZTogb2JqZWN0ID0gYXdhaXQgcmVzdWx0XG4gICAgZXhwZWN0KG1vY2tBeGlvcy5yZXF1ZXN0KS50b0hhdmVCZWVuQ2FsbGVkVGltZXMoMSlcbiAgICBleHBlY3QoSlNPTi5zdHJpbmdpZnkocmVzcG9uc2UpKS50b0JlKEpTT04uc3RyaW5naWZ5KHJlc3BvYmopKVxuICB9KVxuXG4gIHRlc3QoXCJnZXRDb250YWluZXJSYW5nZVwiLCBhc3luYyAoKSA9PiB7XG4gICAgY29uc3Qgc3RhcnRJbmRleDogbnVtYmVyID0gMFxuICAgIGNvbnN0IG51bVRvRmV0Y2g6IG51bWJlciA9IDEwMFxuICAgIGNvbnN0IGVuY29kaW5nOiBzdHJpbmcgPSBcImhleFwiXG4gICAgY29uc3QgYmFzZXVybDogc3RyaW5nID0gXCIvZXh0L2luZGV4L1gvdHhcIlxuICAgIGNvbnN0IHJlc3BvYmogPSB7XG4gICAgICBpZCxcbiAgICAgIGJ5dGVzLFxuICAgICAgdGltZXN0YW1wLFxuICAgICAgZW5jb2RpbmcsXG4gICAgICBpZHhcbiAgICB9XG4gICAgY29uc3QgcmVzdWx0OiBQcm9taXNlPG9iamVjdFtdPiA9IGluZGV4LmdldENvbnRhaW5lclJhbmdlKFxuICAgICAgc3RhcnRJbmRleCxcbiAgICAgIG51bVRvRmV0Y2gsXG4gICAgICBlbmNvZGluZyxcbiAgICAgIGJhc2V1cmxcbiAgICApXG5cbiAgICBjb25zdCBwYXlsb2FkOiBvYmplY3QgPSB7XG4gICAgICByZXN1bHQ6IHJlc3BvYmpcbiAgICB9XG4gICAgY29uc3QgcmVzcG9uc2VPYmo6IEh0dHBSZXNwb25zZSA9IHtcbiAgICAgIGRhdGE6IHBheWxvYWRcbiAgICB9XG5cbiAgICBtb2NrQXhpb3MubW9ja1Jlc3BvbnNlKHJlc3BvbnNlT2JqKVxuICAgIGNvbnN0IHJlc3BvbnNlOiBvYmplY3QgPSBhd2FpdCByZXN1bHRcbiAgICBleHBlY3QobW9ja0F4aW9zLnJlcXVlc3QpLnRvSGF2ZUJlZW5DYWxsZWRUaW1lcygxKVxuICAgIGV4cGVjdChKU09OLnN0cmluZ2lmeShyZXNwb25zZSkpLnRvQmUoSlNPTi5zdHJpbmdpZnkocmVzcG9iaikpXG4gIH0pXG5cbiAgdGVzdChcImdldEluZGV4XCIsIGFzeW5jICgpID0+IHtcbiAgICBjb25zdCBlbmNvZGluZzogc3RyaW5nID0gXCJoZXhcIlxuICAgIGNvbnN0IGJhc2V1cmw6IHN0cmluZyA9IFwiL2V4dC9pbmRleC9YL3R4XCJcbiAgICBjb25zdCByZXN1bHQ6IFByb21pc2U8c3RyaW5nPiA9IGluZGV4LmdldEluZGV4KGlkLCBlbmNvZGluZywgYmFzZXVybClcblxuICAgIGNvbnN0IHBheWxvYWQ6IG9iamVjdCA9IHtcbiAgICAgIHJlc3VsdDoge1xuICAgICAgICBpbmRleDogXCIwXCJcbiAgICAgIH1cbiAgICB9XG4gICAgY29uc3QgcmVzcG9uc2VPYmo6IEh0dHBSZXNwb25zZSA9IHtcbiAgICAgIGRhdGE6IHBheWxvYWRcbiAgICB9XG5cbiAgICBtb2NrQXhpb3MubW9ja1Jlc3BvbnNlKHJlc3BvbnNlT2JqKVxuICAgIGNvbnN0IHJlc3BvbnNlOiBzdHJpbmcgPSBhd2FpdCByZXN1bHRcblxuICAgIGV4cGVjdChtb2NrQXhpb3MucmVxdWVzdCkudG9IYXZlQmVlbkNhbGxlZFRpbWVzKDEpXG4gICAgZXhwZWN0KHJlc3BvbnNlKS50b0JlKFwiMFwiKVxuICB9KVxuXG4gIHRlc3QoXCJpc0FjY2VwdGVkXCIsIGFzeW5jICgpID0+IHtcbiAgICBjb25zdCBlbmNvZGluZzogc3RyaW5nID0gXCJoZXhcIlxuICAgIGNvbnN0IGJhc2V1cmw6IHN0cmluZyA9IFwiL2V4dC9pbmRleC9YL3R4XCJcbiAgICBjb25zdCByZXN1bHQ6IFByb21pc2U8SXNBY2NlcHRlZFJlc3BvbnNlPiA9IGluZGV4LmlzQWNjZXB0ZWQoXG4gICAgICBpZCxcbiAgICAgIGVuY29kaW5nLFxuICAgICAgYmFzZXVybFxuICAgIClcblxuICAgIGNvbnN0IHBheWxvYWQ6IG9iamVjdCA9IHtcbiAgICAgIHJlc3VsdDogdHJ1ZVxuICAgIH1cbiAgICBjb25zdCByZXNwb25zZU9iajogSHR0cFJlc3BvbnNlID0ge1xuICAgICAgZGF0YTogcGF5bG9hZFxuICAgIH1cblxuICAgIG1vY2tBeGlvcy5tb2NrUmVzcG9uc2UocmVzcG9uc2VPYmopXG4gICAgY29uc3QgcmVzcG9uc2U6IElzQWNjZXB0ZWRSZXNwb25zZSA9IGF3YWl0IHJlc3VsdFxuXG4gICAgZXhwZWN0KG1vY2tBeGlvcy5yZXF1ZXN0KS50b0hhdmVCZWVuQ2FsbGVkVGltZXMoMSlcbiAgICBleHBlY3QocmVzcG9uc2UpLnRvQmUodHJ1ZSlcbiAgfSlcbn0pXG4iXX0=