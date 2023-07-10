"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const db_1 = __importDefault(require("src/utils/db"));
describe("DB", () => {
    test("instantiate singletone", () => {
        const db1 = db_1.default.getInstance();
        const db2 = db_1.default.getInstance();
        expect(db1).toEqual(db2);
    });
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZGIudGVzdC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3Rlc3RzL3V0aWxzL2RiLnRlc3QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7QUFBQSxzREFBNkI7QUFFN0IsUUFBUSxDQUFDLElBQUksRUFBRSxHQUFTLEVBQUU7SUFDeEIsSUFBSSxDQUFDLHdCQUF3QixFQUFFLEdBQVMsRUFBRTtRQUN4QyxNQUFNLEdBQUcsR0FBTyxZQUFFLENBQUMsV0FBVyxFQUFFLENBQUE7UUFDaEMsTUFBTSxHQUFHLEdBQU8sWUFBRSxDQUFDLFdBQVcsRUFBRSxDQUFBO1FBQ2hDLE1BQU0sQ0FBQyxHQUFHLENBQUMsQ0FBQyxPQUFPLENBQUMsR0FBRyxDQUFDLENBQUE7SUFDMUIsQ0FBQyxDQUFDLENBQUE7QUFDSixDQUFDLENBQUMsQ0FBQSIsInNvdXJjZXNDb250ZW50IjpbImltcG9ydCBEQiBmcm9tIFwic3JjL3V0aWxzL2RiXCJcblxuZGVzY3JpYmUoXCJEQlwiLCAoKTogdm9pZCA9PiB7XG4gIHRlc3QoXCJpbnN0YW50aWF0ZSBzaW5nbGV0b25lXCIsICgpOiB2b2lkID0+IHtcbiAgICBjb25zdCBkYjE6IERCID0gREIuZ2V0SW5zdGFuY2UoKVxuICAgIGNvbnN0IGRiMjogREIgPSBEQi5nZXRJbnN0YW5jZSgpXG4gICAgZXhwZWN0KGRiMSkudG9FcXVhbChkYjIpXG4gIH0pXG59KVxuIl19