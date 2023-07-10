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
exports.fetchAdapter = void 0;
function createRequest(config) {
    const headers = new Headers(config.headers);
    if (config.auth) {
        const username = config.auth.username || "";
        const password = config.auth.password
            ? encodeURIComponent(config.auth.password)
            : "";
        headers.set("Authorization", `Basic ${Buffer.from(`${username}:${password}`).toString("base64")}`);
    }
    const method = config.method.toUpperCase();
    const options = {
        headers: headers,
        method
    };
    if (method !== "GET" && method !== "HEAD") {
        options.body = config.data;
    }
    if (!!config.withCredentials) {
        options.credentials = config.withCredentials ? "include" : "omit";
    }
    const fullPath = new URL(config.url, config.baseURL);
    const params = new URLSearchParams(config.params);
    const url = `${fullPath}${params}`;
    return new Request(url, options);
}
function getResponse(request, config) {
    return __awaiter(this, void 0, void 0, function* () {
        let stageOne;
        try {
            stageOne = yield fetch(request);
        }
        catch (e) {
            const error = Object.assign(Object.assign({}, new Error("Network Error")), { config,
                request, isAxiosError: true, toJSON: () => error });
            return Promise.reject(error);
        }
        const response = {
            status: stageOne.status,
            statusText: stageOne.statusText,
            headers: Object.assign({}, stageOne.headers),
            config: config,
            request,
            data: undefined // we set it below
        };
        if (stageOne.status >= 200 && stageOne.status !== 204) {
            switch (config.responseType) {
                case "arraybuffer":
                    response.data = yield stageOne.arrayBuffer();
                    break;
                case "blob":
                    response.data = yield stageOne.blob();
                    break;
                case "json":
                    response.data = yield stageOne.json();
                    break;
                case "formData":
                    response.data = yield stageOne.formData();
                    break;
                default:
                    response.data = yield stageOne.text();
                    break;
            }
        }
        return Promise.resolve(response);
    });
}
function fetchAdapter(config) {
    return __awaiter(this, void 0, void 0, function* () {
        const request = createRequest(config);
        const promiseChain = [getResponse(request, config)];
        if (config.timeout && config.timeout > 0) {
            promiseChain.push(new Promise((res, reject) => {
                setTimeout(() => {
                    const message = config.timeoutErrorMessage
                        ? config.timeoutErrorMessage
                        : "timeout of " + config.timeout + "ms exceeded";
                    const error = Object.assign(Object.assign({}, new Error(message)), { config,
                        request, code: "ECONNABORTED", isAxiosError: true, toJSON: () => error });
                    reject(error);
                }, config.timeout);
            }));
        }
        const response = yield Promise.race(promiseChain);
        return new Promise((resolve, reject) => {
            if (response instanceof Error) {
                reject(response);
            }
            else {
                if (!response.status ||
                    !response.config.validateStatus ||
                    response.config.validateStatus(response.status)) {
                    resolve(response);
                }
                else {
                    const error = Object.assign(Object.assign({}, new Error("Request failed with status code " + response.status)), { config,
                        request, code: response.status >= 500 ? "ERR_BAD_RESPONSE" : "ERR_BAD_REQUEST", isAxiosError: true, toJSON: () => error });
                    reject(error);
                }
            }
        });
    });
}
exports.fetchAdapter = fetchAdapter;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZmV0Y2hhZGFwdGVyLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vc3JjL3V0aWxzL2ZldGNoYWRhcHRlci50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7Ozs7QUFFQSxTQUFTLGFBQWEsQ0FBQyxNQUEwQjtJQUMvQyxNQUFNLE9BQU8sR0FBRyxJQUFJLE9BQU8sQ0FBQyxNQUFNLENBQUMsT0FBaUMsQ0FBQyxDQUFBO0lBRXJFLElBQUksTUFBTSxDQUFDLElBQUksRUFBRTtRQUNmLE1BQU0sUUFBUSxHQUFHLE1BQU0sQ0FBQyxJQUFJLENBQUMsUUFBUSxJQUFJLEVBQUUsQ0FBQTtRQUMzQyxNQUFNLFFBQVEsR0FBRyxNQUFNLENBQUMsSUFBSSxDQUFDLFFBQVE7WUFDbkMsQ0FBQyxDQUFDLGtCQUFrQixDQUFDLE1BQU0sQ0FBQyxJQUFJLENBQUMsUUFBUSxDQUFDO1lBQzFDLENBQUMsQ0FBQyxFQUFFLENBQUE7UUFDTixPQUFPLENBQUMsR0FBRyxDQUNULGVBQWUsRUFDZixTQUFTLE1BQU0sQ0FBQyxJQUFJLENBQUMsR0FBRyxRQUFRLElBQUksUUFBUSxFQUFFLENBQUMsQ0FBQyxRQUFRLENBQUMsUUFBUSxDQUFDLEVBQUUsQ0FDckUsQ0FBQTtLQUNGO0lBRUQsTUFBTSxNQUFNLEdBQUcsTUFBTSxDQUFDLE1BQU0sQ0FBQyxXQUFXLEVBQUUsQ0FBQTtJQUMxQyxNQUFNLE9BQU8sR0FBZ0I7UUFDM0IsT0FBTyxFQUFFLE9BQU87UUFDaEIsTUFBTTtLQUNQLENBQUE7SUFDRCxJQUFJLE1BQU0sS0FBSyxLQUFLLElBQUksTUFBTSxLQUFLLE1BQU0sRUFBRTtRQUN6QyxPQUFPLENBQUMsSUFBSSxHQUFHLE1BQU0sQ0FBQyxJQUFJLENBQUE7S0FDM0I7SUFFRCxJQUFJLENBQUMsQ0FBQyxNQUFNLENBQUMsZUFBZSxFQUFFO1FBQzVCLE9BQU8sQ0FBQyxXQUFXLEdBQUcsTUFBTSxDQUFDLGVBQWUsQ0FBQyxDQUFDLENBQUMs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TSxDQUFDLEtBQUssQ0FBQyxDQUFBO2lCQUNkO2FBQ0Y7UUFDSCxDQUFDLENBQUMsQ0FBQTtJQUNKLENBQUM7Q0FBQTtBQXBERCxvQ0FvREMiLCJzb3VyY2VzQ29udGVudCI6WyJpbXBvcnQgeyBBeGlvc1JlcXVlc3RDb25maWcsIEF4aW9zUmVzcG9uc2UsIEF4aW9zRXJyb3IgfSBmcm9tIFwiYXhpb3NcIlxuXG5mdW5jdGlvbiBjcmVhdGVSZXF1ZXN0KGNvbmZpZzogQXhpb3NSZXF1ZXN0Q29uZmlnKTogUmVxdWVzdCB7XG4gIGNvbnN0IGhlYWRlcnMgPSBuZXcgSGVhZGVycyhjb25maWcuaGVhZGVycyBhcyBSZWNvcmQ8c3RyaW5nLCBzdHJpbmc+KVxuXG4gIGlmIChjb25maWcuYXV0aCkge1xuICAgIGNvbnN0IHVzZXJuYW1lID0gY29uZmlnLmF1dGgudXNlcm5hbWUgfHwgXCJcIlxuICAgIGNvbnN0IHBhc3N3b3JkID0gY29uZmlnLmF1dGgucGFzc3dvcmRcbiAgICAgID8gZW5jb2RlVVJJQ29tcG9uZW50KGNvbmZpZy5hdXRoLnBhc3N3b3JkKVxuICAgICAgOiBcIlwiXG4gICAgaGVhZGVycy5zZXQoXG4gICAgICBcIkF1dGhvcml6YXRpb25cIixcbiAgICAgIGBCYXNpYyAke0J1ZmZlci5mcm9tKGAke3VzZXJuYW1lfToke3Bhc3N3b3JkfWApLnRvU3RyaW5nKFwiYmFzZTY0XCIpfWBcbiAgICApXG4gIH1cblxuICBjb25zdCBtZXRob2QgPSBjb25maWcubWV0aG9kLnRvVXBwZXJDYXNlKClcbiAgY29uc3Qgb3B0aW9uczogUmVxdWVzdEluaXQgPSB7XG4gICAgaGVhZGVyczogaGVhZGVycyxcbiAgICBtZXRob2RcbiAgfVxuICBpZiAobWV0aG9kICE9PSBcIkdFVFwiICYmIG1ldGhvZCAhPT0gXCJIRUFEXCIpIHtcbiAgICBvcHRpb25zLmJvZHkgPSBjb25maWcuZGF0YVxuICB9XG5cbiAgaWYgKCEhY29uZmlnLndpdGhDcmVkZW50aWFscykge1xuICAgIG9wdGlvbnMuY3JlZGVudGlhbHMgPSBjb25maWcud2l0aENyZWRlbnRpYWxzID8gXCJpbmNsdWRlXCIgOiBcIm9taXRcIlxuICB9XG5cbiAgY29uc3QgZnVsbFBhdGggPSBuZXcgVVJMKGNvbmZpZy51cmwsIGNvbmZpZy5iYXNlVVJMKVxuICBjb25zdCBwYXJhbXMgPSBuZXcgVVJMU2VhcmNoUGFyYW1zKGNvbmZpZy5wYXJhbXMpXG5cbiAgY29uc3QgdXJsID0gYCR7ZnVsbFBhdGh9JHtwYXJhbXN9YFxuXG4gIHJldHVybiBuZXcgUmVxdWVzdCh1cmwsIG9wdGlvbnMpXG59XG5cbmFzeW5jIGZ1bmN0aW9uIGdldFJlc3BvbnNlKHJlcXVlc3QsIGNvbmZpZyk6IFByb21pc2U8QXhpb3NSZXNwb25zZT4ge1xuICBsZXQgc3RhZ2VPbmVcbiAgdHJ5IHtcbiAgICBzdGFnZU9uZSA9IGF3YWl0IGZldGNoKHJlcXVlc3QpXG4gIH0gY2F0Y2ggKGUpIHtcbiAgICBjb25zdCBlcnJvcjogQXhpb3NFcnJvciA9IHtcbiAgICAgIC4uLm5ldyBFcnJvcihcIk5ldHdvcmsgRXJyb3JcIiksXG4gICAgICBjb25maWcsXG4gICAgICByZXF1ZXN0LFxuICAgICAgaXNBeGlvc0Vycm9yOiB0cnVlLFxuICAgICAgdG9KU09OOiAoKSA9PiBlcnJvclxuICAgIH1cbiAgICByZXR1cm4gUHJvbWlzZS5yZWplY3QoZXJyb3IpXG4gIH1cblxuICBjb25zdCByZXNwb25zZTogQXhpb3NSZXNwb25zZSA9IHtcbiAgICBzdGF0dXM6IHN0YWdlT25lLnN0YXR1cyxcbiAgICBzdGF0dXNUZXh0OiBzdGFnZU9uZS5zdGF0dXNUZXh0LFxuICAgIGhlYWRlcnM6IHsgLi4uc3RhZ2VPbmUuaGVhZGVycyB9LCAvLyBtYWtlIGEgY29weSBvZiB0aGUgaGVhZGVyc1xuICAgIGNvbmZpZzogY29uZmlnLFxuICAgIHJlcXVlc3QsXG4gICAgZGF0YTogdW5kZWZpbmVkIC8vIHdlIHNldCBpdCBiZWxvd1xuICB9XG5cbiAgaWYgKHN0YWdlT25lLnN0YXR1cyA+PSAyMDAgJiYgc3RhZ2VPbmUuc3RhdHVzICE9PSAyMDQpIHtcbiAgICBzd2l0Y2ggKGNvbmZpZy5yZXNwb25zZVR5cGUpIHtcbiAgICAgIGNhc2UgXCJhcnJheWJ1ZmZlclwiOlxuICAgICAgICByZXNwb25zZS5kYXRhID0gYXdhaXQgc3RhZ2VPbmUuYXJyYXlCdWZmZXIoKVxuICAgICAgICBicmVha1xuICAgICAgY2FzZSBcImJsb2JcIjpcbiAgICAgICAgcmVzcG9uc2UuZGF0YSA9IGF3YWl0IHN0YWdlT25lLmJsb2IoKVxuICAgICAgICBicmVha1xuICAgICAgY2FzZSBcImpzb25cIjpcbiAgICAgICAgcmVzcG9uc2UuZGF0YSA9IGF3YWl0IHN0YWdlT25lLmpzb24oKVxuICAgICAgICBicmVha1xuICAgICAgY2FzZSBcImZvcm1EYXRhXCI6XG4gICAgICAgIHJlc3BvbnNlLmRhdGEgPSBhd2FpdCBzdGFnZU9uZS5mb3JtRGF0YSgpXG4gICAgICAgIGJyZWFrXG4gICAgICBkZWZhdWx0OlxuICAgICAgICByZXNwb25zZS5kYXRhID0gYXdhaXQgc3RhZ2VPbmUudGV4dCgpXG4gICAgICAgIGJyZWFrXG4gICAgfVxuICB9XG5cbiAgcmV0dXJuIFByb21pc2UucmVzb2x2ZShyZXNwb25zZSlcbn1cblxuZXhwb3J0IGFzeW5jIGZ1bmN0aW9uIGZldGNoQWRhcHRlcihcbiAgY29uZmlnOiBBeGlvc1JlcXVlc3RDb25maWdcbik6IFByb21pc2U8QXhpb3NSZXNwb25zZT4ge1xuICBjb25zdCByZXF1ZXN0ID0gY3JlYXRlUmVxdWVzdChjb25maWcpXG5cbiAgY29uc3QgcHJvbWlzZUNoYWluID0gW2dldFJlc3BvbnNlKHJlcXVlc3QsIGNvbmZpZyldXG5cbiAgaWYgKGNvbmZpZy50aW1lb3V0ICYmIGNvbmZpZy50aW1lb3V0ID4gMCkge1xuICAgIHByb21pc2VDaGFpbi5wdXNoKFxuICAgICAgbmV3IFByb21pc2UoKHJlcywgcmVqZWN0KSA9PiB7XG4gICAgICAgIHNldFRpbWVvdXQoKCkgPT4ge1xuICAgICAgICAgIGNvbnN0IG1lc3NhZ2UgPSBjb25maWcudGltZW91dEVycm9yTWVzc2FnZVxuICAgICAgICAgICAgPyBjb25maWcudGltZW91dEVycm9yTWVzc2FnZVxuICAgICAgICAgICAgOiBcInRpbWVvdXQgb2YgXCIgKyBjb25maWcudGltZW91dCArIFwibXMgZXhjZWVkZWRcIlxuICAgICAgICAgIGNvbnN0IGVycm9yOiBBeGlvc0Vycm9yID0ge1xuICAgICAgICAgICAgLi4ubmV3IEVycm9yKG1lc3NhZ2UpLFxuICAgICAgICAgICAgY29uZmlnLFxuICAgICAgICAgICAgcmVxdWVzdCxcbiAgICAgICAgICAgIGNvZGU6IFwiRUNPTk5BQk9SVEVEXCIsXG4gICAgICAgICAgICBpc0F4aW9zRXJyb3I6IHRydWUsXG4gICAgICAgICAgICB0b0pTT046ICgpID0+IGVycm9yXG4gICAgICAgICAgfVxuICAgICAgICAgIHJlamVjdChlcnJvcilcbiAgICAgICAgfSwgY29uZmlnLnRpbWVvdXQpXG4gICAgICB9KVxuICAgIClcbiAgfVxuXG4gIGNvbnN0IHJlc3BvbnNlID0gYXdhaXQgUHJvbWlzZS5yYWNlKHByb21pc2VDaGFpbilcbiAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICBpZiAocmVzcG9uc2UgaW5zdGFuY2VvZiBFcnJvcikge1xuICAgICAgcmVqZWN0KHJlc3BvbnNlKVxuICAgIH0gZWxzZSB7XG4gICAgICBpZiAoXG4gICAgICAgICFyZXNwb25zZS5zdGF0dXMgfHxcbiAgICAgICAgIXJlc3BvbnNlLmNvbmZpZy52YWxpZGF0ZVN0YXR1cyB8fFxuICAgICAgICByZXNwb25zZS5jb25maWcudmFsaWRhdGVTdGF0dXMocmVzcG9uc2Uuc3RhdHVzKVxuICAgICAgKSB7XG4gICAgICAgIHJlc29sdmUocmVzcG9uc2UpXG4gICAgICB9IGVsc2Uge1xuICAgICAgICBjb25zdCBlcnJvcjogQXhpb3NFcnJvciA9IHtcbiAgICAgICAgICAuLi5uZXcgRXJyb3IoXCJSZXF1ZXN0IGZhaWxlZCB3aXRoIHN0YXR1cyBjb2RlIFwiICsgcmVzcG9uc2Uuc3RhdHVzKSxcbiAgICAgICAgICBjb25maWcsXG4gICAgICAgICAgcmVxdWVzdCxcbiAgICAgICAgICBjb2RlOiByZXNwb25zZS5zdGF0dXMgPj0gNTAwID8gXCJFUlJfQkFEX1JFU1BPTlNFXCIgOiBcIkVSUl9CQURfUkVRVUVTVFwiLFxuICAgICAgICAgIGlzQXhpb3NFcnJvcjogdHJ1ZSxcbiAgICAgICAgICB0b0pTT046ICgpID0+IGVycm9yXG4gICAgICAgIH1cbiAgICAgICAgcmVqZWN0KGVycm9yKVxuICAgICAgfVxuICAgIH1cbiAgfSlcbn1cbiJdfQ==