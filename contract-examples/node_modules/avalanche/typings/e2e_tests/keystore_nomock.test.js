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
const e2etestlib_1 = require("./e2etestlib");
describe("Keystore", () => {
    const username1 = "avalancheJsUser1";
    const username2 = "avalancheJsUser2";
    const username3 = "avalancheJsUser3";
    const password = "avalancheJsP1ssw4rd";
    let exportedUser = { value: "" };
    const avalanche = (0, e2etestlib_1.getAvalanche)();
    const keystore = avalanche.NodeKeys();
    // test_name             response_promise                              resp_fn  matcher           expected_value/obtained_value
    const tests_spec = [
        [
            "createUserWeakPass",
            () => keystore.createUser(username1, "weak"),
            (x) => x,
            e2etestlib_1.Matcher.toThrow,
            () => "password is too weak"
        ],
        [
            "createUser",
            () => keystore.createUser(username1, password),
            (x) => x,
            e2etestlib_1.Matcher.toBe,
            () => true
        ],
        [
            "createRepeatedUser",
            () => keystore.createUser(username1, password),
            (x) => x,
            e2etestlib_1.Matcher.toThrow,
            () => "user already exists: " + username1
        ],
        [
            "listUsers",
            () => keystore.listUsers(),
            (x) => x,
            e2etestlib_1.Matcher.toContain,
            () => [username1]
        ],
        [
            "exportUser",
            () => keystore.exportUser(username1, password),
            (x) => x,
            e2etestlib_1.Matcher.toMatch,
            () => /\w{78}/
        ],
        [
            "getExportedUser",
            () => keystore.exportUser(username1, password),
            (x) => x,
            e2etestlib_1.Matcher.Get,
            () => exportedUser
        ],
        [
            "importUser",
            () => keystore.importUser(username2, exportedUser.value, password),
            (x) => x,
            e2etestlib_1.Matcher.toBe,
            () => true
        ],
        [
            "exportImportUser",
            () => (() => __awaiter(void 0, void 0, void 0, function* () {
                let exported = yield keystore.exportUser(username1, password);
                return yield keystore.importUser(username3, exported, password);
            }))(),
            (x) => x,
            e2etestlib_1.Matcher.toBe,
            () => true
        ],
        [
            "listUsers2",
            () => keystore.listUsers(),
            (x) => x,
            e2etestlib_1.Matcher.toContain,
            () => [username1, username2, username3]
        ],
        [
            "deleteUser1",
            () => keystore.deleteUser(username1, password),
            (x) => x,
            e2etestlib_1.Matcher.toBe,
            () => true
        ],
        [
            "deleteUser2",
            () => keystore.deleteUser(username2, password),
            (x) => x,
            e2etestlib_1.Matcher.toBe,
            () => true
        ],
        [
            "deleteUser3",
            () => keystore.deleteUser(username3, password),
            (x) => x,
            e2etestlib_1.Matcher.toBe,
            () => true
        ]
    ];
    (0, e2etestlib_1.createTests)(tests_spec);
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoia2V5c3RvcmVfbm9tb2NrLnRlc3QuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi9lMmVfdGVzdHMva2V5c3RvcmVfbm9tb2NrLnRlc3QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7QUFBQSw2Q0FBaUU7QUFJakUsUUFBUSxDQUFDLFVBQVUsRUFBRSxHQUFTLEVBQUU7SUFDOUIsTUFBTSxTQUFTLEdBQVcsa0JBQWtCLENBQUE7SUFDNUMsTUFBTSxTQUFTLEdBQVcsa0JBQWtCLENBQUE7SUFDNUMsTUFBTSxTQUFTLEdBQVcsa0JBQWtCLENBQUE7SUFDNUMsTUFBTSxRQUFRLEdBQVcscUJBQXFCLENBQUE7SUFFOUMsSUFBSSxZQUFZLEdBQUcsRUFBRSxLQUFLLEVBQUUsRUFBRSxFQUFFLENBQUE7SUFFaEMsTUFBTSxTQUFTLEdBQWMsSUFBQSx5QkFBWSxHQUFFLENBQUE7SUFDM0MsTUFBTSxRQUFRLEdBQWdCLFNBQVMsQ0FBQyxRQUFRLEVBQUUsQ0FBQTtJQUVsRCwrSEFBK0g7SUFDL0gsTUFBTSxVQUFVLEdBQVE7UUFDdEI7WUFDRSxvQkFBb0I7WUFDcEIsR0FBRyxFQUFFLENBQUMsUUFBUSxDQUFDLFVBQVUsQ0FBQyxTQUFTLEVBQUUsTUFBTSxDQUFDO1lBQzVDLENBQUMsQ0FBQyxFQUFFLEVBQUUsQ0FBQyxDQUFDO1lBQ1Isb0JBQU8sQ0FBQyxPQUFPO1lBQ2YsR0FBRyxFQUFFLENBQUMsc0JBQXNCO1NBQzdCO1FBQ0Q7WUFDRSxZQUFZO1lBQ1osR0FBRyxFQUFFLENBQUMsUUFBUSxDQUFDLFVBQVUsQ0FBQyxTQUFTLEVBQUUsUUFBUSxDQUFDO1lBQzlDLENBQUMsQ0FBQyxFQUFFLEVBQUUsQ0FBQyxDQUFDO1lBQ1Isb0JBQU8sQ0FBQyxJQUFJO1lBQ1osR0FBRyxFQUFFLENBQUMsSUFBSTtTQUNYO1FBQ0Q7WUFDRSxvQkFBb0I7WUFDcEIsR0FBRyxFQUFFLENBQUMsUUFBUSxDQUFDLFVBQVUsQ0FBQyxTQUFTLEVB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QyxDQUFDLENBQUEiLCJzb3VyY2VzQ29udGVudCI6WyJpbXBvcnQgeyBnZXRBdmFsYW5jaGUsIGNyZWF0ZVRlc3RzLCBNYXRjaGVyIH0gZnJvbSBcIi4vZTJldGVzdGxpYlwiXG5pbXBvcnQgeyBLZXlzdG9yZUFQSSB9IGZyb20gXCJzcmMvYXBpcy9rZXlzdG9yZS9hcGlcIlxuaW1wb3J0IEF2YWxhbmNoZSBmcm9tIFwic3JjXCJcblxuZGVzY3JpYmUoXCJLZXlzdG9yZVwiLCAoKTogdm9pZCA9PiB7XG4gIGNvbnN0IHVzZXJuYW1lMTogc3RyaW5nID0gXCJhdmFsYW5jaGVKc1VzZXIxXCJcbiAgY29uc3QgdXNlcm5hbWUyOiBzdHJpbmcgPSBcImF2YWxhbmNoZUpzVXNlcjJcIlxuICBjb25zdCB1c2VybmFtZTM6IHN0cmluZyA9IFwiYXZhbGFuY2hlSnNVc2VyM1wiXG4gIGNvbnN0IHBhc3N3b3JkOiBzdHJpbmcgPSBcImF2YWxhbmNoZUpzUDFzc3c0cmRcIlxuXG4gIGxldCBleHBvcnRlZFVzZXIgPSB7IHZhbHVlOiBcIlwiIH1cblxuICBjb25zdCBhdmFsYW5jaGU6IEF2YWxhbmNoZSA9IGdldEF2YWxhbmNoZSgpXG4gIGNvbnN0IGtleXN0b3JlOiBLZXlzdG9yZUFQSSA9IGF2YWxhbmNoZS5Ob2RlS2V5cygpXG5cbiAgLy8gdGVzdF9uYW1lICAgICAgICAgICAgIHJlc3BvbnNlX3Byb21pc2UgICAgICAgICAgICAgICAgICAgICAgICAgICAgICByZXNwX2ZuICBtYXRjaGVyICAgICAgICAgICBleHBlY3RlZF92YWx1ZS9vYnRhaW5lZF92YWx1ZVxuICBjb25zdCB0ZXN0c19zcGVjOiBhbnkgPSBbXG4gICAgW1xuICAgICAgXCJjcmVhdGVVc2VyV2Vha1Bhc3NcIixcbiAgICAgICgpID0+IGtleXN0b3JlLmNyZWF0ZVVzZXIodXNlcm5hbWUxLCBcIndlYWtcIiksXG4gICAgICAoeCkgPT4geCxcbiAgICAgIE1hdGNoZXIudG9UaHJvdyxcbiAgICAgICgpID0+IFwicGFzc3dvcmQgaXMgdG9vIHdlYWtcIlxuICAgIF0sXG4gICAgW1xuICAgICAgXCJjcmVhdGVVc2VyXCIsXG4gICAgICAoKSA9PiBrZXlzdG9yZS5jcmVhdGVVc2VyKHVzZXJuYW1lMSwgcGFzc3dvcmQpLFxuICAgICAgKHgpID0+IHgsXG4gICAgICBNYXRjaGVyLnRvQmUsXG4gICAgICAoKSA9PiB0cnVlXG4gICAgXSxcbiAgICBbXG4gICAgICBcImNyZWF0ZVJlcGVhdGVkVXNlclwiLFxuICAgICAgKCkgPT4ga2V5c3RvcmUuY3JlYXRlVXNlcih1c2VybmFtZTEsIHBhc3N3b3JkKSxcbiAgICAgICh4KSA9PiB4LFxuICAgICAgTWF0Y2hlci50b1Rocm93LFxuICAgICAgKCkgPT4gXCJ1c2VyIGFscmVhZHkgZXhpc3RzOiBcIiArIHVzZXJuYW1lMVxuICAgIF0sXG4gICAgW1xuICAgICAgXCJsaXN0VXNlcnNcIixcbiAgICAgICgpID0+IGtleXN0b3JlLmxpc3RVc2VycygpLFxuICAgICAgKHgpID0+IHgsXG4gICAgICBNYXRjaGVyLnRvQ29udGFpbixcbiAgICAgICgpID0+IFt1c2VybmFtZTFdXG4gICAgXSxcbiAgICBbXG4gICAgICBcImV4cG9ydFVzZXJcIixcbiAgICAgICgpID0+IGtleXN0b3JlLmV4cG9ydFVzZXIodXNlcm5hbWUxLCBwYXNzd29yZCksXG4gICAgICAoeCkgPT4geCxcbiAgICAgIE1hdGNoZXIudG9NYXRjaCxcbiAgICAgICgpID0+IC9cXHd7Nzh9L1xuICAgIF0sXG4gICAgW1xuICAgICAgXCJnZXRFeHBvcnRlZFVzZXJcIixcbiAgICAgICgpID0+IGtleXN0b3JlLmV4cG9ydFVzZXIodXNlcm5hbWUxLCBwYXNzd29yZCksXG4gICAgICAoeCkgPT4geCxcbiAgICAgIE1hdGNoZXIuR2V0LFxuICAgICAgKCkgPT4gZXhwb3J0ZWRVc2VyXG4gICAgXSxcbiAgICBbXG4gICAgICBcImltcG9ydFVzZXJcIixcbiAgICAgICgpID0+IGtleXN0b3JlLmltcG9ydFVzZXIodXNlcm5hbWUyLCBleHBvcnRlZFVzZXIudmFsdWUsIHBhc3N3b3JkKSxcbiAgICAgICh4KSA9PiB4LFxuICAgICAgTWF0Y2hlci50b0JlLFxuICAgICAgKCkgPT4gdHJ1ZVxuICAgIF0sXG4gICAgW1xuICAgICAgXCJleHBvcnRJbXBvcnRVc2VyXCIsXG4gICAgICAoKSA9PlxuICAgICAgICAoYXN5bmMgKCkgPT4ge1xuICAgICAgICAgIGxldCBleHBvcnRlZCA9IGF3YWl0IGtleXN0b3JlLmV4cG9ydFVzZXIodXNlcm5hbWUxLCBwYXNzd29yZClcbiAgICAgICAgICByZXR1cm4gYXdhaXQga2V5c3RvcmUuaW1wb3J0VXNlcih1c2VybmFtZTMsIGV4cG9ydGVkLCBwYXNzd29yZClcbiAgICAgICAgfSkoKSxcbiAgICAgICh4KSA9PiB4LFxuICAgICAgTWF0Y2hlci50b0JlLFxuICAgICAgKCkgPT4gdHJ1ZVxuICAgIF0sXG4gICAgW1xuICAgICAgXCJsaXN0VXNlcnMyXCIsXG4gICAgICAoKSA9PiBrZXlzdG9yZS5saXN0VXNlcnMoKSxcbiAgICAgICh4KSA9PiB4LFxuICAgICAgTWF0Y2hlci50b0NvbnRhaW4sXG4gICAgICAoKSA9PiBbdXNlcm5hbWUxLCB1c2VybmFtZTIsIHVzZXJuYW1lM11cbiAgICBdLFxuICAgIFtcbiAgICAgIFwiZGVsZXRlVXNlcjFcIixcbiAgICAgICgpID0+IGtleXN0b3JlLmRlbGV0ZVVzZXIodXNlcm5hbWUxLCBwYXNzd29yZCksXG4gICAgICAoeCkgPT4geCxcbiAgICAgIE1hdGNoZXIudG9CZSxcbiAgICAgICgpID0+IHRydWVcbiAgICBdLFxuICAgIFtcbiAgICAgIFwiZGVsZXRlVXNlcjJcIixcbiAgICAgICgpID0+IGtleXN0b3JlLmRlbGV0ZVVzZXIodXNlcm5hbWUyLCBwYXNzd29yZCksXG4gICAgICAoeCkgPT4geCxcbiAgICAgIE1hdGNoZXIudG9CZSxcbiAgICAgICgpID0+IHRydWVcbiAgICBdLFxuICAgIFtcbiAgICAgIFwiZGVsZXRlVXNlcjNcIixcbiAgICAgICgpID0+IGtleXN0b3JlLmRlbGV0ZVVzZXIodXNlcm5hbWUzLCBwYXNzd29yZCksXG4gICAgICAoeCkgPT4geCxcbiAgICAgIE1hdGNoZXIudG9CZSxcbiAgICAgICgpID0+IHRydWVcbiAgICBdXG4gIF1cblxuICBjcmVhdGVUZXN0cyh0ZXN0c19zcGVjKVxufSkiXX0=