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
exports.createTests = exports.Matcher = exports.getAvalanche = void 0;
const src_1 = require("src");
const getAvalanche = () => {
    if (typeof process.env.AVALANCHEGO_IP === "undefined") {
        throw "Undefined environment variable: AVALANCHEGO_IP";
    }
    if (typeof process.env.AVALANCHEGO_PORT === "undefined") {
        throw "Undefined environment variable: AVALANCHEGO_PORT";
    }
    const avalanche = new src_1.Avalanche(process.env.AVALANCHEGO_IP, parseInt(process.env.AVALANCHEGO_PORT));
    return avalanche;
};
exports.getAvalanche = getAvalanche;
var Matcher;
(function (Matcher) {
    Matcher[Matcher["toBe"] = 0] = "toBe";
    Matcher[Matcher["toEqual"] = 1] = "toEqual";
    Matcher[Matcher["toContain"] = 2] = "toContain";
    Matcher[Matcher["toMatch"] = 3] = "toMatch";
    Matcher[Matcher["toThrow"] = 4] = "toThrow";
    Matcher[Matcher["Get"] = 5] = "Get";
})(Matcher = exports.Matcher || (exports.Matcher = {}));
const createTests = (tests_spec) => {
    for (const [testName, promise, preprocess, matcher, expected] of tests_spec) {
        test(testName, () => __awaiter(void 0, void 0, void 0, function* () {
            if (matcher == Matcher.toBe) {
                expect(preprocess(yield promise())).toBe(expected());
            }
            if (matcher == Matcher.toEqual) {
                expect(preprocess(yield promise())).toEqual(expected());
            }
            if (matcher == Matcher.toContain) {
                expect(preprocess(yield promise())).toEqual(expect.arrayContaining(expected()));
            }
            if (matcher == Matcher.toMatch) {
                expect(preprocess(yield promise())).toMatch(expected());
            }
            if (matcher == Matcher.toThrow) {
                yield expect(preprocess(promise())).rejects.toThrow(expected());
            }
            if (matcher == Matcher.Get) {
                expected().value = preprocess(yield promise());
                expect(true).toBe(true);
            }
        }));
    }
};
exports.createTests = createTests;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZTJldGVzdGxpYi5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uL2UyZV90ZXN0cy9lMmV0ZXN0bGliLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7Ozs7Ozs7OztBQUFBLDZCQUErQjtBQUV4QixNQUFNLFlBQVksR0FBRyxHQUFjLEVBQUU7SUFDMUMsSUFBSSxPQUFPLE9BQU8sQ0FBQyxHQUFHLENBQUMsY0FBYyxLQUFLLFdBQVcsRUFBRTtRQUNyRCxNQUFNLGdEQUFnRCxDQUFBO0tBQ3ZEO0lBQ0QsSUFBSSxPQUFPLE9BQU8sQ0FBQyxHQUFHLENBQUMsZ0JBQWdCLEtBQUssV0FBVyxFQUFFO1FBQ3ZELE1BQU0sa0RBQWtELENBQUE7S0FDekQ7SUFDRCxNQUFNLFNBQVMsR0FBYyxJQUFJLGVBQVMsQ0FDeEMsT0FBTyxDQUFDLEdBQUcsQ0FBQyxjQUFjLEVBQzFCLFFBQVEsQ0FBQyxPQUFPLENBQUMsR0FBRyxDQUFDLGdCQUFnQixDQUFDLENBQ3ZDLENBQUE7SUFDRCxPQUFPLFNBQVMsQ0FBQTtBQUNsQixDQUFDLENBQUE7QUFaWSxRQUFBLFlBQVksZ0JBWXhCO0FBRUQsSUFBWSxPQU9YO0FBUEQsV0FBWSxPQUFPO0lBQ2pCLHFDQUFJLENBQUE7SUFDSiwyQ0FBTyxDQUFBO0lBQ1AsK0NBQVMsQ0FBQTtJQUNULDJDQUFPLENBQUE7SUFDUCwyQ0FBTyxDQUFBO0lBQ1AsbUNBQUcsQ0FBQTtBQUNMLENBQUMsRUFQVyxPQUFPLEdBQVAsZUFBTyxLQUFQLGVBQU8sUUFPbEI7QUFFTSxNQUFNLFdBQVcsR0FBRyxDQUFDLFVBQWlCLEVBQVEsRUFBRTtJQUNyRCxLQUFLLE1BQU0sQ0FBQyxRQUFRLEVBQUUsT0FBTyxFQUFFLFVBQVUsRUFBRSxPQUFPLEVBQUUsUUFBUSxDQUFDLElBQUksVUFBVSxFQUFFO1FBQzNFLElBQUksQ0FBQyxRQUFRLEVBQUUsR0FBd0IsRUFBRTtZQUN2QyxJQUFJLE9BQU8sSUFBSSxPQUFPLENBQUMsSUFBSSxFQUFFO2dCQUMzQixNQUFNLENBQUMsVUFBVSxDQUFDLE1BQU0sT0FBTyxFQUFFLENBQUMsQ0FBQyxDQUFDLElBQUksQ0FBQyxRQUFRLEVBQUUsQ0FBQyxDQUFBO2FBQ3JEO1lBQ0QsSUFBSSxPQUFPLElBQUksT0FBTyxDQUFDLE9BQU8sRUFBRTtnQkFDOUIsTUFBTSxDQUFDLFVBQVUsQ0FBQyxNQUFNLE9BQU8sRUFBRSxDQUFDLENBQUMsQ0FBQyxPQUFPLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQTthQUN4RDtZQUNELElBQUksT0FBTyxJQUFJLE9BQU8sQ0FBQyxTQUFTLEVBQUU7Z0JBQ2hDLE1BQU0sQ0FBQyxVQUFVLENBQUMsTUFBTSxPQUFPLEVBQUUsQ0FBQyxDQUFDLENBQUMsT0FBTyxDQUFDLE1BQU0sQ0FBQyxlQUFlLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQyxDQUFBO2FBQ2hGO1lBQ0QsSUFBSSxPQUFPLElBQUksT0FBTyxDQUFDLE9BQU8sRUFBRTtnQkFDOUIsTUFBTSxDQUFDLFVBQVUsQ0FBQyxNQUFNLE9BQU8sRUFBRSxDQUFDLENBQUMsQ0FBQyxPQUFPLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQTthQUN4RDtZQUNELElBQUksT0FBTyxJQUFJLE9BQU8sQ0FBQyxPQUFPLEVBQUU7Z0JBQzlCLE1BQU0sTUFBTSxDQUFDLFVBQVUsQ0FBQyxPQUFPLEVBQUUsQ0FBQyxDQUFDLENBQUMsT0FBTyxDQUFDLE9BQU8sQ0FBQyxRQUFRLEVBQUUsQ0FBQyxDQUFBO2FBQ2hFO1lBQ0QsSUFBSSxPQUFPLElBQUksT0FBTyxDQUFDLEdBQUcsRUFBRTtnQkFDMUIsUUFBUSxFQUFFLENBQUMsS0FBSyxHQUFHLFVBQVUsQ0FBQyxNQUFNLE9BQU8sRUFBRSxDQUFDLENBQUE7Z0JBQzlDLE1BQU0sQ0FBQyxJQUFJLENBQUMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7YUFDeEI7UUFDSCxDQUFDLENBQUEsQ0FBQyxDQUFBO0tBQ0g7QUFDSCxDQUFDLENBQUE7QUF4QlksUUFBQSxXQUFXLGVBd0J2QiIsInNvdXJjZXNDb250ZW50IjpbImltcG9ydCB7IEF2YWxhbmNoZSB9IGZyb20gXCJzcmNcIlxuXG5leHBvcnQgY29uc3QgZ2V0QXZhbGFuY2hlID0gKCk6IEF2YWxhbmNoZSA9PiB7XG4gIGlmICh0eXBlb2YgcHJvY2Vzcy5lbnYuQVZBTEFOQ0hFR09fSVAgPT09IFwidW5kZWZpbmVkXCIpIHtcbiAgICB0aHJvdyBcIlVuZGVmaW5lZCBlbnZpcm9ubWVudCB2YXJpYWJsZTogQVZBTEFOQ0hFR09fSVBcIlxuICB9XG4gIGlmICh0eXBlb2YgcHJvY2Vzcy5lbnYuQVZBTEFOQ0hFR09fUE9SVCA9PT0gXCJ1bmRlZmluZWRcIikge1xuICAgIHRocm93IFwiVW5kZWZpbmVkIGVudmlyb25tZW50IHZhcmlhYmxlOiBBVkFMQU5DSEVHT19QT1JUXCJcbiAgfVxuICBjb25zdCBhdmFsYW5jaGU6IEF2YWxhbmNoZSA9IG5ldyBBdmFsYW5jaGUoXG4gICAgcHJvY2Vzcy5lbnYuQVZBTEFOQ0hFR09fSVAsXG4gICAgcGFyc2VJbnQocHJvY2Vzcy5lbnYuQVZBTEFOQ0hFR09fUE9SVClcbiAgKVxuICByZXR1cm4gYXZhbGFuY2hlXG59XG5cbmV4cG9ydCBlbnVtIE1hdGNoZXIge1xuICB0b0JlLFxuICB0b0VxdWFsLFxuICB0b0NvbnRhaW4sXG4gIHRvTWF0Y2gsXG4gIHRvVGhyb3csXG4gIEdldFxufVxuXG5leHBvcnQgY29uc3QgY3JlYXRlVGVzdHMgPSAodGVzdHNfc3BlYzogYW55W10pOiB2b2lkID0+IHtcbiAgZm9yIChjb25zdCBbdGVzdE5hbWUsIHByb21pc2UsIHByZXByb2Nlc3MsIG1hdGNoZXIsIGV4cGVjdGVkXSBvZiB0ZXN0c19zcGVjKSB7XG4gICAgdGVzdCh0ZXN0TmFtZSwgYXN5bmMgKCk6IFByb21pc2U8dm9pZD4gPT4ge1xuICAgICAgaWYgKG1hdGNoZXIgPT0gTWF0Y2hlci50b0JlKSB7XG4gICAgICAgIGV4cGVjdChwcmVwcm9jZXNzKGF3YWl0IHByb21pc2UoKSkpLnRvQmUoZXhwZWN0ZWQoKSlcbiAgICAgIH1cbiAgICAgIGlmIChtYXRjaGVyID09IE1hdGNoZXIudG9FcXVhbCkge1xuICAgICAgICBleHBlY3QocHJlcHJvY2Vzcyhhd2FpdCBwcm9taXNlKCkpKS50b0VxdWFsKGV4cGVjdGVkKCkpXG4gICAgICB9XG4gICAgICBpZiAobWF0Y2hlciA9PSBNYXRjaGVyLnRvQ29udGFpbikge1xuICAgICAgICBleHBlY3QocHJlcHJvY2Vzcyhhd2FpdCBwcm9taXNlKCkpKS50b0VxdWFsKGV4cGVjdC5hcnJheUNvbnRhaW5pbmcoZXhwZWN0ZWQoKSkpXG4gICAgICB9XG4gICAgICBpZiAobWF0Y2hlciA9PSBNYXRjaGVyLnRvTWF0Y2gpIHtcbiAgICAgICAgZXhwZWN0KHByZXByb2Nlc3MoYXdhaXQgcHJvbWlzZSgpKSkudG9NYXRjaChleHBlY3RlZCgpKVxuICAgICAgfVxuICAgICAgaWYgKG1hdGNoZXIgPT0gTWF0Y2hlci50b1Rocm93KSB7XG4gICAgICAgIGF3YWl0IGV4cGVjdChwcmVwcm9jZXNzKHByb21pc2UoKSkpLnJlamVjdHMudG9UaHJvdyhleHBlY3RlZCgpKVxuICAgICAgfVxuICAgICAgaWYgKG1hdGNoZXIgPT0gTWF0Y2hlci5HZXQpIHtcbiAgICAgICAgZXhwZWN0ZWQoKS52YWx1ZSA9IHByZXByb2Nlc3MoYXdhaXQgcHJvbWlzZSgpKVxuICAgICAgICBleHBlY3QodHJ1ZSkudG9CZSh0cnVlKVxuICAgICAgfVxuICAgIH0pXG4gIH1cbn1cblxuIl19