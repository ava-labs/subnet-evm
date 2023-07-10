"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const bintools_1 = __importDefault(require("../../src/utils/bintools"));
const evm_1 = require("src/apis/evm");
const bintools = bintools_1.default.getInstance();
describe("SECP256K1", () => {
    test("addressFromPublicKey", () => {
        const pubkeys = [
            "7ECaZ7TpWLq6mh3858DkR3EzEToGi8iFFxnjY5hUGePoCHqdjw",
            "5dS4sSyL4dHziqLYanMoath8dqUMe6ZkY1VbnVuQQSsCcgtVET"
        ];
        const addrs = [
            "b0c9654511ebb78d490bb0d7a54997d4a933972c",
            "d5bb99a29e09853da983be63a76f02259ceedf15"
        ];
        pubkeys.forEach((pubkey, index) => {
            const pubkeyBuf = bintools.cb58Decode(pubkey);
            const addrBuf = evm_1.KeyPair.addressFromPublicKey(pubkeyBuf);
            expect(addrBuf.toString("hex")).toBe(addrs[index]);
        });
    });
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoic2VjcDI1NmsxLnRlc3QuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi90ZXN0cy9jb21tb24vc2VjcDI1NmsxLnRlc3QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7QUFDQSx3RUFBK0M7QUFDL0Msc0NBQXNDO0FBRXRDLE1BQU0sUUFBUSxHQUFhLGtCQUFRLENBQUMsV0FBVyxFQUFFLENBQUE7QUFFakQsUUFBUSxDQUFDLFdBQVcsRUFBRSxHQUFTLEVBQUU7SUFDL0IsSUFBSSxDQUFDLHNCQUFzQixFQUFFLEdBQVMsRUFBRTtRQUN0QyxNQUFNLE9BQU8sR0FBYTtZQUN4QixvREFBb0Q7WUFDcEQsb0RBQW9EO1NBQ3JELENBQUE7UUFDRCxNQUFNLEtBQUssR0FBYTtZQUN0QiwwQ0FBMEM7WUFDMUMsMENBQTBDO1NBQzNDLENBQUE7UUFDRCxPQUFPLENBQUMsT0FBTyxDQUFDLENBQUMsTUFBYyxFQUFFLEtBQWEsRUFBUSxFQUFFO1lBQ3RELE1BQU0sU0FBUyxHQUFXLFFBQVEsQ0FBQyxVQUFVLENBQUMsTUFBTSxDQUFDLENBQUE7WUFDckQsTUFBTSxPQUFPLEdBQVcsYUFBTyxDQUFDLG9CQUFvQixDQUFDLFNBQVMsQ0FBQyxDQUFBO1lBQy9ELE1BQU0sQ0FBQyxPQUFPLENBQUMsUUFBUSxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsSUFBSSxDQUFDLEtBQUssQ0FBQyxLQUFLLENBQUMsQ0FBQyxDQUFBO1FBQ3BELENBQUMsQ0FBQyxDQUFBO0lBQ0osQ0FBQyxDQUFDLENBQUE7QUFDSixDQUFDLENBQUMsQ0FBQSIsInNvdXJjZXNDb250ZW50IjpbImltcG9ydCB7IEJ1ZmZlciB9IGZyb20gXCJidWZmZXIvXCJcbmltcG9ydCBCaW5Ub29scyBmcm9tIFwiLi4vLi4vc3JjL3V0aWxzL2JpbnRvb2xzXCJcbmltcG9ydCB7IEtleVBhaXIgfSBmcm9tIFwic3JjL2FwaXMvZXZtXCJcblxuY29uc3QgYmludG9vbHM6IEJpblRvb2xzID0gQmluVG9vbHMuZ2V0SW5zdGFuY2UoKVxuXG5kZXNjcmliZShcIlNFQ1AyNTZLMVwiLCAoKTogdm9pZCA9PiB7XG4gIHRlc3QoXCJhZGRyZXNzRnJvbVB1YmxpY0tleVwiLCAoKTogdm9pZCA9PiB7XG4gICAgY29uc3QgcHVia2V5czogc3RyaW5nW10gPSBbXG4gICAgICBcIjdFQ2FaN1RwV0xxNm1oMzg1OERrUjNFekVUb0dpOGlGRnhualk1aFVHZVBvQ0hxZGp3XCIsXG4gICAgICBcIjVkUzRzU3lMNGRIemlxTFlhbk1vYXRoOGRxVU1lNlprWTFWYm5WdVFRU3NDY2d0VkVUXCJcbiAgICBdXG4gICAgY29uc3QgYWRkcnM6IHN0cmluZ1tdID0gW1xuICAgICAgXCJiMGM5NjU0NTExZWJiNzhkNDkwYmIwZDdhNTQ5OTdkNGE5MzM5NzJjXCIsXG4gICAgICBcImQ1YmI5OWEyOWUwOTg1M2RhOTgzYmU2M2E3NmYwMjI1OWNlZWRmMTVcIlxuICAgIF1cbiAgICBwdWJrZXlzLmZvckVhY2goKHB1YmtleTogc3RyaW5nLCBpbmRleDogbnVtYmVyKTogdm9pZCA9PiB7XG4gICAgICBjb25zdCBwdWJrZXlCdWY6IEJ1ZmZlciA9IGJpbnRvb2xzLmNiNThEZWNvZGUocHVia2V5KVxuICAgICAgY29uc3QgYWRkckJ1ZjogQnVmZmVyID0gS2V5UGFpci5hZGRyZXNzRnJvbVB1YmxpY0tleShwdWJrZXlCdWYpXG4gICAgICBleHBlY3QoYWRkckJ1Zi50b1N0cmluZyhcImhleFwiKSkudG9CZShhZGRyc1tpbmRleF0pXG4gICAgfSlcbiAgfSlcbn0pXG4iXX0=