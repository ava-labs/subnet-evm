"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const bn_js_1 = __importDefault(require("bn.js"));
const buffer_1 = require("buffer/");
const bintools_1 = __importDefault(require("../../../src/utils/bintools"));
const utxos_1 = require("../../../src/apis/platformvm/utxos");
const helperfunctions_1 = require("../../../src/utils/helperfunctions");
const bintools = bintools_1.default.getInstance();
const display = "display";
describe("UTXO", () => {
    const utxohex = "000038d1b9f1138672da6fb6c35125539276a9acc2a668d63bea6ba3c795e2edb0f5000000013e07e38e2f23121be8756412c18db7246a16d26ee9936f3cba28be149cfd3558000000070000000000004dd500000000000000000000000100000001a36fd0c2dbcab311731dde7ef1514bd26fcdc74d";
    const outputidx = "00000001";
    const outtxid = "38d1b9f1138672da6fb6c35125539276a9acc2a668d63bea6ba3c795e2edb0f5";
    const outaid = "3e07e38e2f23121be8756412c18db7246a16d26ee9936f3cba28be149cfd3558";
    const utxobuff = buffer_1.Buffer.from(utxohex, "hex");
    // Payment
    const OPUTXOstr = bintools.cb58Encode(utxobuff);
    // "U9rFgK5jjdXmV8k5tpqeXkimzrN3o9eCCcXesyhMBBZu9MQJCDTDo5Wn5psKvzJVMJpiMbdkfDXkp7sKZddfCZdxpuDmyNy7VFka19zMW4jcz6DRQvNfA2kvJYKk96zc7uizgp3i2FYWrB8mr1sPJ8oP9Th64GQ5yHd8"
    // implies fromString and fromBuffer
    test("Creation", () => {
        const u1 = new utxos_1.UTXO();
        u1.fromBuffer(utxobuff);
        const u1hex = u1.toBuffer().toString("hex");
        expect(u1hex).toBe(utxohex);
    });
    test("Empty Creation", () => {
        const u1 = new utxos_1.UTXO();
        expect(() => {
            u1.toBuffer();
        }).toThrow();
    });
    test("Creation of Type", () => {
        const op = new utxos_1.UTXO();
        op.fromString(OPUTXOstr);
        expect(op.getOutput().getOutputID()).toBe(7);
    });
    describe("Funtionality", () => {
        const u1 = new utxos_1.UTXO();
        u1.fromBuffer(utxobuff);
        const u1hex = u1.toBuffer().toString("hex");
        test("getAssetID NonCA", () => {
            const assetID = u1.getAssetID();
            expect(assetID.toString("hex", 0, assetID.length)).toBe(outaid);
        });
        test("getTxID", () => {
            const txid = u1.getTxID();
            expect(txid.toString("hex", 0, txid.length)).toBe(outtxid);
        });
        test("getOutputIdx", () => {
            const txidx = u1.getOutputIdx();
            expect(txidx.toString("hex", 0, txidx.length)).toBe(outputidx);
        });
        test("getUTXOID", () => {
            const txid = buffer_1.Buffer.from(outtxid, "hex");
            const txidx = buffer_1.Buffer.from(outputidx, "hex");
            const utxoid = bintools.bufferToB58(buffer_1.Buffer.concat([txid, txidx]));
            expect(u1.getUTXOID()).toBe(utxoid);
        });
        test("toString", () => {
            const serialized = u1.toString();
            expect(serialized).toBe(bintools.cb58Encode(utxobuff));
        });
    });
});
const setMergeTester = (input, equal, notEqual) => {
    const instr = JSON.stringify(input.getUTXOIDs().sort());
    for (let i = 0; i < equal.length; i++) {
        if (JSON.stringify(equal[i].getUTXOIDs().sort()) != instr) {
            return false;
        }
    }
    for (let i = 0; i < notEqual.length; i++) {
        if (JSON.stringify(notEqual[i].getUTXOIDs().sort()) == instr) {
            return false;
        }
    }
    return true;
};
describe("UTXOSet", () => {
    const utxostrs = [
        bintools.cb58Encode(buffer_1.Buffer.from("000038d1b9f1138672da6fb6c35125539276a9acc2a668d63bea6ba3c795e2edb0f5000000013e07e38e2f23121be8756412c18db7246a16d26ee9936f3cba28be149cfd3558000000070000000000004dd500000000000000000000000100000001a36fd0c2dbcab311731dde7ef1514bd26fcdc74d", "hex")),
        bintools.cb58Encode(buffer_1.Buffer.from("0000c3e4823571587fe2bdfc502689f5a8238b9d0ea7f3277124d16af9de0d2d9911000000003e07e38e2f23121be8756412c18db7246a16d26ee9936f3cba28be149cfd355800000007000000000000001900000000000000000000000100000001e1b6b6a4bad94d2e3f20730379b9bcd6f176318e", "hex")),
        bintools.cb58Encode(buffer_1.Buffer.from("0000f29dba61fda8d57a911e7f8810f935bde810d3f8d495404685bdb8d9d8545e86000000003e07e38e2f23121be8756412c18db7246a16d26ee9936f3cba28be149cfd355800000007000000000000001900000000000000000000000100000001e1b6b6a4bad94d2e3f20730379b9bcd6f176318e", "hex"))
    ];
    const addrs = [
        bintools.cb58Decode("FuB6Lw2D62NuM8zpGLA4Avepq7eGsZRiG"),
        bintools.cb58Decode("MaTvKGccbYzCxzBkJpb2zHW7E1WReZqB8")
    ];
    test("Creation", () => {
        const set = new utxos_1.UTXOSet();
        set.add(utxostrs[0]);
        const utxo = new utxos_1.UTXO();
        utxo.fromString(utxostrs[0]);
        const setArray = set.getAllUTXOs();
        expect(utxo.toString()).toBe(setArray[0].toString());
    });
    test("Serialization", () => {
        const set = new utxos_1.UTXOSet();
        set.addArray([...utxostrs]);
        let setobj = set.serialize("cb58");
        let setstr = JSON.stringify(setobj);
        let set2newobj = JSON.parse(setstr);
        let set2 = new utxos_1.UTXOSet();
        set2.deserialize(set2newobj, "cb58");
        let set2obj = set2.serialize("cb58");
        let set2str = JSON.stringify(set2obj);
        expect(set2.getAllUTXOStrings().sort().join(",")).toBe(set.getAllUTXOStrings().sort().join(","));
    });
    test("Mutliple add", () => {
        const set = new utxos_1.UTXOSet();
        // first add
        for (let i = 0; i < utxostrs.length; i++) {
            set.add(utxostrs[i]);
        }
        // the verify (do these steps separate to ensure no overwrites)
        for (let i = 0; i < utxostrs.length; i++) {
            expect(set.includes(utxostrs[i])).toBe(true);
            const utxo = new utxos_1.UTXO();
            utxo.fromString(utxostrs[i]);
            const veriutxo = set.getUTXO(utxo.getUTXOID());
            expect(veriutxo.toString()).toBe(utxostrs[i]);
        }
    });
    test("addArray", () => {
        const set = new utxos_1.UTXOSet();
        set.addArray(utxostrs);
        for (let i = 0; i < utxostrs.length; i++) {
            const e1 = new utxos_1.UTXO();
            e1.fromString(utxostrs[i]);
            expect(set.includes(e1)).toBe(true);
            const utxo = new utxos_1.UTXO();
            utxo.fromString(utxostrs[i]);
            const veriutxo = set.getUTXO(utxo.getUTXOID());
            expect(veriutxo.toString()).toBe(utxostrs[i]);
        }
        set.addArray(set.getAllUTXOs());
        for (let i = 0; i < utxostrs.length; i++) {
            const utxo = new utxos_1.UTXO();
            utxo.fromString(utxostrs[i]);
            expect(set.includes(utxo)).toBe(true);
            const veriutxo = set.getUTXO(utxo.getUTXOID());
            expect(veriutxo.toString()).toBe(utxostrs[i]);
        }
        let o = set.serialize("hex");
        let s = new utxos_1.UTXOSet();
        s.deserialize(o);
        let t = set.serialize(display);
        let r = new utxos_1.UTXOSet();
        r.deserialize(t);
    });
    test("overwriting UTXO", () => {
        const set = new utxos_1.UTXOSet();
        set.addArray(utxostrs);
        const testutxo = new utxos_1.UTXO();
        testutxo.fromString(utxostrs[0]);
        expect(set.add(utxostrs[0], true).toString()).toBe(testutxo.toString());
        expect(set.add(utxostrs[0], false)).toBeUndefined();
        expect(set.addArray(utxostrs, true).length).toBe(3);
        expect(set.addArray(utxostrs, false).length).toBe(0);
    });
    describe("Functionality", () => {
        let set;
        let utxos;
        beforeEach(() => {
            set = new utxos_1.UTXOSet();
            set.addArray(utxostrs);
            utxos = set.getAllUTXOs();
        });
        test("remove", () => {
            const testutxo = new utxos_1.UTXO();
            testutxo.fromString(utxostrs[0]);
            expect(set.remove(utxostrs[0]).toString()).toBe(testutxo.toString());
            expect(set.remove(utxostrs[0])).toBeUndefined();
            expect(set.add(utxostrs[0], false).toString()).toBe(testutxo.toString());
            expect(set.remove(utxostrs[0]).toString()).toBe(testutxo.toString());
        });
        test("removeArray", () => {
            const testutxo = new utxos_1.UTXO();
            testutxo.fromString(utxostrs[0]);
            expect(set.removeArray(utxostrs).length).toBe(3);
            expect(set.removeArray(utxostrs).length).toBe(0);
            expect(set.add(utxostrs[0], false).toString()).toBe(testutxo.toString());
            expect(set.removeArray(utxostrs).length).toBe(1);
            expect(set.addArray(utxostrs, false).length).toBe(3);
            expect(set.removeArray(utxos).length).toBe(3);
        });
        test("getUTXOIDs", () => {
            const uids = set.getUTXOIDs();
            for (let i = 0; i < utxos.length; i++) {
                expect(uids.indexOf(utxos[i].getUTXOID())).not.toBe(-1);
            }
        });
        test("getAllUTXOs", () => {
            const allutxos = set.getAllUTXOs();
            const ustrs = [];
            for (let i = 0; i < allutxos.length; i++) {
                ustrs.push(allutxos[i].toString());
            }
            for (let i = 0; i < utxostrs.length; i++) {
                expect(ustrs.indexOf(utxostrs[i])).not.toBe(-1);
            }
            const uids = set.getUTXOIDs();
            const allutxos2 = set.getAllUTXOs(uids);
            const ustrs2 = [];
            for (let i = 0; i < allutxos.length; i++) {
                ustrs2.push(allutxos2[i].toString());
            }
            for (let i = 0; i < utxostrs.length; i++) {
                expect(ustrs2.indexOf(utxostrs[i])).not.toBe(-1);
            }
        });
        test("getUTXOIDs By Address", () => {
            let utxoids;
            utxoids = set.getUTXOIDs([addrs[0]]);
            expect(utxoids.length).toBe(1);
            utxoids = set.getUTXOIDs(addrs);
            expect(utxoids.length).toBe(3);
            utxoids = set.getUTXOIDs(addrs, false);
            expect(utxoids.length).toBe(3);
        });
        test("getAllUTXOStrings", () => {
            const ustrs = set.getAllUTXOStrings();
            for (let i = 0; i < utxostrs.length; i++) {
                expect(ustrs.indexOf(utxostrs[i])).not.toBe(-1);
            }
            const uids = set.getUTXOIDs();
            const ustrs2 = set.getAllUTXOStrings(uids);
            for (let i = 0; i < utxostrs.length; i++) {
                expect(ustrs2.indexOf(utxostrs[i])).not.toBe(-1);
            }
        });
        test("getAddresses", () => {
            expect(set.getAddresses().sort()).toStrictEqual(addrs.sort());
        });
        test("getBalance", () => {
            let balance1;
            let balance2;
            balance1 = new bn_js_1.default(0);
            balance2 = new bn_js_1.default(0);
            for (let i = 0; i < utxos.length; i++) {
                const assetID = utxos[i].getAssetID();
                balance1 = balance1.add(set.getBalance(addrs, assetID));
                balance2 = balance2.add(utxos[i].getOutput().getAmount());
            }
            expect(balance1.gt(new bn_js_1.default(0))).toBe(true);
            expect(balance2.gt(new bn_js_1.default(0))).toBe(true);
            balance1 = new bn_js_1.default(0);
            balance2 = new bn_js_1.default(0);
            const now = (0, helperfunctions_1.UnixNow)();
            for (let i = 0; i < utxos.length; i++) {
                const assetID = bintools.cb58Encode(utxos[i].getAssetID());
                balance1 = balance1.add(set.getBalance(addrs, assetID, now));
                balance2 = balance2.add(utxos[i].getOutput().getAmount());
            }
            expect(balance1.gt(new bn_js_1.default(0))).toBe(true);
            expect(balance2.gt(new bn_js_1.default(0))).toBe(true);
        });
        test("getAssetIDs", () => {
            const assetIDs = set.getAssetIDs();
            for (let i = 0; i < utxos.length; i++) {
                expect(assetIDs).toContain(utxos[i].getAssetID());
            }
            const addresses = set.getAddresses();
            expect(set.getAssetIDs(addresses)).toEqual(set.getAssetIDs());
        });
        describe("Merge Rules", () => {
            let setA;
            let setB;
            let setC;
            let setD;
            let setE;
            let setF;
            let setG;
            let setH;
            // Take-or-Leave
            const newutxo = bintools.cb58Encode(buffer_1.Buffer.from("0000acf88647b3fbaa9fdf4378f3a0df6a5d15d8efb018ad78f12690390e79e1687600000003acf88647b3fbaa9fdf4378f3a0df6a5d15d8efb018ad78f12690390e79e168760000000700000000000186a000000000000000000000000100000001fceda8f90fcb5d30614b99d79fc4baa293077626", "hex"));
            beforeEach(() => {
                setA = new utxos_1.UTXOSet();
                setA.addArray([utxostrs[0], utxostrs[2]]);
                setB = new utxos_1.UTXOSet();
                setB.addArray([utxostrs[1], utxostrs[2]]);
                setC = new utxos_1.UTXOSet();
                setC.addArray([utxostrs[0], utxostrs[1]]);
                setD = new utxos_1.UTXOSet();
                setD.addArray([utxostrs[1]]);
                setE = new utxos_1.UTXOSet();
                setE.addArray([]); // empty set
                setF = new utxos_1.UTXOSet();
                setF.addArray(utxostrs); // full set, separate from self
                setG = new utxos_1.UTXOSet();
                setG.addArray([newutxo, ...utxostrs]); // full set with new element
                setH = new utxos_1.UTXOSet();
                setH.addArray([newutxo]); // set with only a new element
            });
            test("unknown merge rule", () => {
                expect(() => {
                    set.mergeByRule(setA, "ERROR");
                }).toThrow();
                const setArray = setG.getAllUTXOs();
            });
            test("intersection", () => {
                let results;
                let test;
                results = set.mergeByRule(setA, "intersection");
                test = setMergeTester(results, [setA], [setB, setC, setD, setE, setF, setG, setH]);
                expect(test).toBe(true);
                results = set.mergeByRule(setF, "intersection");
                test = setMergeTester(results, [setF], [setA, setB, setC, setD, setE, setG, setH]);
                expect(test).toBe(true);
                results = set.mergeByRule(setG, "intersection");
                test = setMergeTester(results, [setF], [setA, setB, setC, setD, setE, setG, setH]);
                expect(test).toBe(true);
                results = set.mergeByRule(setH, "intersection");
                test = setMergeTester(results, [setE], [setA, setB, setC, setD, setF, setG, setH]);
                expect(test).toBe(true);
            });
            test("differenceSelf", () => {
                let results;
                let test;
                results = set.mergeByRule(setA, "differenceSelf");
                test = setMergeTester(results, [setD], [setA, setB, setC, setE, setF, setG, setH]);
                expect(test).toBe(true);
                results = set.mergeByRule(setF, "differenceSelf");
                test = setMergeTester(results, [setE], [setA, setB, setC, setD, setF, setG, setH]);
                expect(test).toBe(true);
                results = set.mergeByRule(setG, "differenceSelf");
                test = setMergeTester(results, [setE], [setA, setB, setC, setD, setF, setG, setH]);
                expect(test).toBe(true);
                results = set.mergeByRule(setH, "differenceSelf");
                test = setMergeTester(results, [setF], [setA, setB, setC, setD, setE, setG, setH]);
                expect(test).toBe(true);
            });
            test("differenceNew", () => {
                let results;
                let test;
                results = set.mergeByRule(setA, "differenceNew");
                test = setMergeTester(results, [setE], [setA, setB, setC, setD, setF, setG, setH]);
                expect(test).toBe(true);
                results = set.mergeByRule(setF, "differenceNew");
                test = setMergeTester(results, [setE], [setA, setB, setC, setD, setF, setG, setH]);
                expect(test).toBe(true);
                results = set.mergeByRule(setG, "differenceNew");
                test = setMergeTester(results, [setH], [setA, setB, setC, setD, setE, setF, setG]);
                expect(test).toBe(true);
                results = set.mergeByRule(setH, "differenceNew");
                test = setMergeTester(results, [setH], [setA, setB, setC, setD, setE, setF, setG]);
                expect(test).toBe(true);
            });
            test("symDifference", () => {
                let results;
                let test;
                results = set.mergeByRule(setA, "symDifference");
                test = setMergeTester(results, [setD], [setA, setB, setC, setE, setF, setG, setH]);
                expect(test).toBe(true);
                results = set.mergeByRule(setF, "symDifference");
                test = setMergeTester(results, [setE], [setA, setB, setC, setD, setF, setG, setH]);
                expect(test).toBe(true);
                results = set.mergeByRule(setG, "symDifference");
                test = setMergeTester(results, [setH], [setA, setB, setC, setD, setE, setF, setG]);
                expect(test).toBe(true);
                results = set.mergeByRule(setH, "symDifference");
                test = setMergeTester(results, [setG], [setA, setB, setC, setD, setE, setF, setH]);
                expect(test).toBe(true);
            });
            test("union", () => {
                let results;
                let test;
                results = set.mergeByRule(setA, "union");
                test = setMergeTester(results, [setF], [setA, setB, setC, setD, setE, setG, setH]);
                expect(test).toBe(true);
                results = set.mergeByRule(setF, "union");
                test = setMergeTester(results, [setF], [setA, setB, setC, setD, setE, setG, setH]);
                expect(test).toBe(true);
                results = set.mergeByRule(setG, "union");
                test = setMergeTester(results, [setG], [setA, setB, setC, setD, setE, setF, setH]);
                expect(test).toBe(true);
                results = set.mergeByRule(setH, "union");
                test = setMergeTester(results, [setG], [setA, setB, setC, setD, setE, setF, setH]);
                expect(test).toBe(true);
            });
            test("unionMinusNew", () => {
                let results;
                let test;
                results = set.mergeByRule(setA, "unionMinusNew");
                test = setMergeTester(results, [setD], [setA, setB, setC, setE, setF, setG, setH]);
                expect(test).toBe(true);
                results = set.mergeByRule(setF, "unionMinusNew");
                test = setMergeTester(results, [setE], [setA, setB, setC, setD, setF, setG, setH]);
                expect(test).toBe(true);
                results = set.mergeByRule(setG, "unionMinusNew");
                test = setMergeTester(results, [setE], [setA, setB, setC, setD, setF, setG, setH]);
                expect(test).toBe(true);
                results = set.mergeByRule(setH, "unionMinusNew");
                test = setMergeTester(results, [setF], [setA, setB, setC, setD, setE, setG, setH]);
                expect(test).toBe(true);
            });
            test("unionMinusSelf", () => {
                let results;
                let test;
                results = set.mergeByRule(setA, "unionMinusSelf");
                test = setMergeTester(results, [setE], [setA, setB, setC, setD, setF, setG, setH]);
                expect(test).toBe(true);
                results = set.mergeByRule(setF, "unionMinusSelf");
                test = setMergeTester(results, [setE], [setA, setB, setC, setD, setF, setG, setH]);
                expect(test).toBe(true);
                results = set.mergeByRule(setG, "unionMinusSelf");
                test = setMergeTester(results, [setH], [setA, setB, setC, setD, setE, setF, setG]);
                expect(test).toBe(true);
                results = set.mergeByRule(setH, "unionMinusSelf");
                test = setMergeTester(results, [setH], [setA, setB, setC, setD, setE, setF, setG]);
                expect(test).toBe(true);
            });
        });
    });
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidXR4b3MudGVzdC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3Rlc3RzL2FwaXMvcGxhdGZvcm12bS91dHhvcy50ZXN0LnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7O0FBQUEsa0RBQXNCO0FBQ3RCLG9DQUFnQztBQUNoQywyRUFBa0Q7QUFDbEQsOERBQWtFO0FBRWxFLHdFQUE0RDtBQUc1RCxNQUFNLFFBQVEsR0FBYSxrQkFBUSxDQUFDLFdBQVcsRUFBRSxDQUFBO0FBQ2pELE1BQU0sT0FBTyxHQUF1QixTQUFTLENBQUE7QUFFN0MsUUFBUSxDQUFDLE1BQU0sRUFBRSxHQUFTLEVBQUU7SUFDMUIsTUFBTSxPQUFPLEdBQ1gsOE9BQThPLENBQUE7SUFDaFAsTUFBTSxTQUFTLEdBQVcsVUFBVSxDQUFBO0lBQ3BDLE1BQU0sT0FBTyxHQUNYLGtFQUFrRSxDQUFBO0lBQ3BFLE1BQU0sTUFBTSxHQUNWLGtFQUFrRSxDQUFBO0lBQ3BFLE1BQU0sUUFBUSxHQUFXLGVBQU0sQ0FBQyxJQUFJLENBQUMsT0FBTyxFQUFFLEtBQUssQ0FBQyxDQUFBO0lBRXBELFVBQVU7SUFDVixNQUFNLFNBQVMsR0FBVyxRQUFRLENBQUMsVUFBVSxDQUFDLFFBQVEsQ0FBQyxDQUFBO0lBQ3ZELHlLQUF5SztJQUV6SyxvQ0FBb0M7SUFDcEMsSUFBSSxDQUFDLFVBQVUsRUFBRSxHQUFTLEVBQUU7UUFDMUIsTUFBTSxFQUFFLEdBQVMsSUFBSSxZQUFJLEVBQUUsQ0FBQTtRQUMzQixFQUFFLENBQUMsVUFBVSxDQUFDLFFBQVEsQ0FBQyxDQUFBO1FBQ3ZCLE1BQU0sS0FBSyxHQUFXLEVBQUUsQ0FBQyxRQUFRLEVBQUUsQ0FBQyxRQUFRLENBQUMsS0FBSyxDQUFDLENBQUE7UUFDbkQsTUFBTSxDQUFDLEtBQUssQ0FBQyxDQUFDLElBQUksQ0FBQyxPQUFPLENBQUMsQ0FBQTtJQUM3QixDQUFDLENBQUMsQ0FBQTtJQUVGLElBQUksQ0FBQyxnQkFBZ0IsRUFBRSxHQUFTLEVBQUU7UUFDaEMsTUFBTSxFQUFFLEdBQVMsSUFBSSxZQUFJLEVBQUUsQ0FBQTtRQUMzQixNQUFNLENBQUMsR0FBRyxFQUFFO1lBQ1YsRUFBRSxDQUFDLFFBQVEsRUFBRSxDQUFBO1FBQ2YsQ0FBQyxDQUFDLENBQUMsT0FBTyxFQUFFLENBQUE7SUFDZCxDQUFDLENBQUMsQ0FBQTtJQUVGLElBQUksQ0FBQyxrQkFBa0IsRUFBRSxHQUFTLEVBQUU7UUFDbEMsTUFBTSxFQUFFLEdBQVMsSUFBSSxZQUFJLEVBQUUsQ0FBQTtRQUMzQixFQUFFLENBQUMsVUFBVSxDQUFDLFNBQVMsQ0FBQyxDQUFBO1FBQ3hCLE1BQU0sQ0FBQyxFQUFFLENBQUMsU0FBUyxFQUFFLENBQUMsV0FBVyxFQUFFLENBQUMsQ0FBQyxJQUFJLENBQUMsQ0FBQyxDQUFDLENBQUE7SUFDOUMsQ0FBQyxDQUFDLENBQUE7SUFFRixRQUFRLENBQUMsY0FBYyxFQUFFLEdBQVMsRUFBRTtRQUNsQyxNQUFNLEVBQUUsR0FBUyxJQUFJLFlBQUksRUFBRSxDQUFBO1FBQzNCLEVBQUUsQ0FBQyxVQUFVLENBQUMsUUFBUSxDQUFDLENBQUE7UUFDdkIsTUFBTSxLQUFLLEdBQVcsRUFBRSxDQUFDLFFBQVEsRUFBRSxDQUFDLFFBQVEsQ0FBQyxLQUFLLENBQUMsQ0FBQTtRQUNuRCxJQUFJLENBQUMsa0JBQWtCLEVBQUUsR0FBUyxFQUFFO1lBQ2xDLE1BQU0sT0FBTyxHQUFXLEVBQUUsQ0FBQyxVQUFVLEVBQUUsQ0FBQTtZQUN2QyxNQUFNLENBQUMsT0FBTyxDQUFDLFFBQVEsQ0FBQyxLQUFLLEVBQUUsQ0FBQyxFQUFFLE9BQU8sQ0FBQyxNQUFNLENBQUMsQ0FBQyxDQUFDLElBQUksQ0FBQyxNQUFNLENBQUMsQ0FBQTtRQUNqRSxDQUFDLENBQUMsQ0FBQTtRQUNGLElBQUksQ0FBQyxTQUFTLEVBQUUsR0FBUyxFQUFFO1lBQ3pCLE1BQU0sSUFBSSxHQUFXLEVBQUUsQ0FBQyxPQUFPLEVBQUUsQ0FBQTtZQUNqQyxNQUFNLENBQUMsSUFBSSxDQUFDLFFBQVEsQ0FBQyxLQUFLLEVBQUUsQ0FBQyxFQUFFLElBQUksQ0FBQyxNQUFNLENBQUMsQ0FBQyxDQUFDLElBQUksQ0FBQyxPQUFPLENBQUMsQ0FBQTtRQUM1RCxDQUFDLENBQUMsQ0FBQTtRQUNGLElBQUksQ0FBQyxjQUFjLEVBQUUsR0FBUyxFQUFFO1lBQzlCLE1BQU0sS0FBSyxHQUFXLEVBQUUsQ0FBQyxZQUFZLEVBQUUsQ0FBQTtZQUN2QyxNQUFNLENBQUMsS0FBSyxDQUFDLFFBQVEsQ0FBQyxLQUFLLEVBQUUsQ0FBQyxFQUFFLEtBQUssQ0FBQyxNQUFNLENBQUMsQ0FBQyxDQUFDLElBQUksQ0FBQyxTQUFTLENBQUMsQ0FBQTtRQUNoRSxDQUFDLENBQUMsQ0FBQTtRQUNGLElBQUksQ0FBQyxXQUFXLEVBQUUsR0FBUyxFQUFFO1lBQzNCLE1BQU0sSUFBSSxHQUFXLGVBQU0sQ0FBQyxJQUFJLENBQUMsT0FBTyxFQUFFLEtBQUssQ0FBQyxDQUFBO1lBQ2hELE1BQU0sS0FBSyxHQUFXLGVBQU0sQ0FBQyxJQUFJLENBQUMsU0FBUyxFQUFFLEtBQUssQ0FBQyxDQUFBO1lBQ25ELE1BQU0sTUFBTSxHQUFXLFFBQVEsQ0FBQyxXQUFXLENBQUMsZUFBTSxDQUFDLE1BQU0sQ0FBQyxDQUFDLElBQUksRUFBRSxLQUFLLENBQUMsQ0FBQyxDQUFDLENBQUE7WUFDekUsTUFBTSxDQUFDLEVBQUUsQ0FBQyxTQUFTLEVBQUUsQ0FBQyxDQUFDLElBQUksQ0FBQyxNQUFNLENBQUMsQ0FBQTtRQUNyQyxDQUFDLENBQUMsQ0FBQTtRQUNGLElBQUksQ0FBQyxVQUFVLEVBQUUsR0FBUyxFQUFFO1lBQzFCLE1BQU0sVUFBVSxHQUFXLEVBQUUsQ0FBQyxRQUFRLEVBQUUsQ0FBQTtZQUN4QyxNQUFNLENBQUMsVUFBVSxDQUFDLENBQUMsSUFBSSxDQUFDLFFBQVEsQ0FBQyxVQUFVLENBQUMsUUFBUSxDQUFDLENBQUMsQ0FBQTtRQUN4RCxDQUFDLENBQUMsQ0FBQTtJQUNKLENBQUMsQ0FBQyxDQUFBO0FBQ0osQ0FBQyxDQUFDLENBQUE7QUFFRixNQUFNLGNBQWMsR0FBRyxDQUNyQixLQUFjLEVBQ2QsS0FBZ0IsRUFDaEIsUUFBbUIsRUFDVixFQUFFO0lBQ1gsTUFBTSxLQUFLLEdBQVcsSUFBSSxDQUFDLFNBQVMsQ0FBQyxLQUFLLENBQUMsVUFBVSxFQUFFLENBQUMsSUFBSSxFQUFFLENBQUMsQ0FBQTtJQUMvRCxLQUFLLElBQUksQ0FBQyxHQUFXLENBQUMsRUFBRSxDQUFDLEdBQUcsS0FBSyxDQUFDLE1BQU0sRUFBRSxDQUFDLEVBQUUsRUFBRTtRQUM3QyxJQUFJLElBQUksQ0FBQyxTQUFTLENBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQyxDQUFDLFVBQVUsRUFBRSxDQUFDLElBQUksRUFBRSxDQUFDLElBQUksS0FBSyxFQUFFO1lBQ3pELE9BQU8sS0FBSyxDQUFBO1NBQ2I7S0FDRjtJQUVELEtBQUssSUFBSSxDQUFDLEdBQVcsQ0FBQyxFQUFFLENBQUMsR0FBRyxRQUFRLENBQUMsTUFBTSxFQUFFLENBQUMsRUFBRSxFQUFFO1FBQ2hELElBQUksSUFBSSxDQUFDLFNBQVMsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLENBQUMsVUFBVSxFQUFFLENBQUMsSUFBSSxFQUFFLENBQUMsSUFBSSxLQUFLLEVBQUU7WUFDNUQsT0FBTyxLQUFLLENBQUE7U0FDYjtLQUNGO0lBQ0QsT0FBTyxJQUFJLENBQUE7QUFDYixDQUFDLENBQUE7QUFFRCxRQUFRLENBQUMsU0FBUyxFQUFFLEdBQVMsRUFBRTtJQUM3QixNQUFNLFFBQVEsR0FBYTtRQUN6QixRQUFRLENBQUMsVUFBVSxDQUNqQixlQUFNLENBQUMsSUFBSSxDQUNULDhPQUE4TyxFQUM5TyxLQUFLLENBQ04sQ0FDRjtRQUNELFFBQVEsQ0FBQyxVQUFVLENBQ2pCLGVBQU0sQ0FBQyxJQUFJLENBQ1QsOE9BQThPLEVBQzlPLEtBQUssQ0FDTixDQUNGO1FBQ0QsUUFBUSxDQUFDLFVBQVUsQ0FDakIsZUFBTSxDQUFDLElBQUksQ0FDVCw4T0FBOE8sRUFDOU8sS0FBSyxDQUNOLENBQ0Y7S0FDRixDQUFBO0lBQ0QsTUFBTSxLQUFLLEdBQWE7UUFDdEIsUUFBUSxDQUFDLFVBQVUsQ0FBQyxtQ0FBbUMsQ0FBQztRQUN4RCxRQUFRLENBQUMsVUFBVSxDQUFDLG1DQUFtQyxDQUFDO0tBQ3pELENBQUE7SUFDRCxJQUFJLENBQUMsVUFBVSxFQUFFLEdBQVMsRUFBRTtRQUMxQixNQUFNLEdBQUcsR0FBWSxJQUFJLGVBQU8sRUFBRSxDQUFBO1FBQ2xDLEdBQUcsQ0FBQyxHQUFHLENBQUMsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUE7UUFDcEIsTUFBTSxJQUFJLEdBQVMsSUFBSSxZQUFJLEVBQUUsQ0FBQTtRQUM3QixJQUFJLENBQUMsVUFBVSxDQUFDLFFBQVEsQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFBO1FBQzVCLE1BQU0sUUFBUSxHQUFXLEdBQUcsQ0FBQyxXQUFXLEVBQUUsQ0FBQTtRQUMxQyxNQUFNLENBQUMsSUFBSSxDQUFDLFFBQVEsRUFBRSxDQUFDLENBQUMsSUFBSSxDQUFDLFFBQVEsQ0FBQyxDQUFDLENBQUMsQ0FBQyxRQUFRLEVBQUUsQ0FBQyxDQUFBO0lBQ3RELENBQUMsQ0FBQyxDQUFBO0lBRUYsSUFBSSxDQUFDLGVBQWUsRUFBRSxHQUFTLEVBQUU7UUFDL0IsTUFBTSxHQUFHLEdBQVksSUFBSSxlQUFPLEVBQUUsQ0FBQTtRQUNsQyxHQUFHLENBQUMsUUFBUSxDQUFDLENBQUMsR0FBRyxRQUFRLENBQUMsQ0FBQyxDQUFBO1FBQzNCLElBQUksTUFBTSxHQUFXLEdBQUcsQ0FBQyxTQUFTLENBQUMsTUFBTSxDQUFDLENBQUE7UUFDMUMsSUFBSSxNQUFNLEdBQVcsSUFBSSxDQUFDLFNBQVMsQ0FBQyxNQUFNLENBQUMsQ0FBQTtRQUMzQyxJQUFJLFVBQVUsR0FBVyxJQUFJLENBQUMsS0FBSyxDQUFDLE1BQU0sQ0FBQyxDQUFBO1FBQzNDLElBQUksSUFBSSxHQUFZLElBQUksZUFBTyxFQUFFLENBQUE7UUFDakMsSUFBSSxDQUFDLFdBQVcsQ0FBQyxVQUFVLEVBQUUsTUFBTSxDQUFDLENBQUE7UUFDcEMsSUFBSSxPQUFPLEdBQVcsSUFBSSxDQUFDLFNBQVMsQ0FBQyxNQUFNLENBQUMsQ0FBQTtRQUM1QyxJQUFJLE9BQU8sR0FBVyxJQUFJLENBQUMsU0FBUyxDQUFDLE9BQU8sQ0FBQyxDQUFBO1FBQzdDLE1BQU0sQ0FBQyxJQUFJLENBQUMsaUJBQWlCLEVBQUUsQ0FBQyxJQUFJLEVBQUUsQ0FBQyxJQUFJLENBQUMsR0FBRyxDQUFDLENBQUMsQ0FBQyxJQUFJLENBQ3BELEdBQUcsQ0FBQyxpQkFBaUIsRUFBRSxDQUFDLElBQUksRUFBRSxDQUFDLElBQUksQ0FBQyxHQUFHLENBQUMsQ0FDekMsQ0FBQTtJQUNILENBQUMsQ0FBQyxDQUFBO0lBRUYsSUFBSSxDQUFDLGNBQWMsRUFBRSxHQUFTLEVBQUU7UUFDOUIsTUFBTSxHQUFHLEdBQVksSUFBSSxlQUFPLEVBQUUsQ0FBQTtRQUNsQyxZQUFZO1FBQ1osS0FBSyxJQUFJLENBQUMsR0FBVyxDQUFDLEVBQUUsQ0FBQyxHQUFHLFFBQVEsQ0FBQyxNQUFNLEVBQUUsQ0FBQyxFQUFFLEVBQUU7WUFDaEQsR0FBRyxDQUFDLEdBQUcsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTtTQUNyQjtRQUNELCtEQUErRDtRQUMvRCxLQUFLLElBQUksQ0FBQyxHQUFXLENBQUMsRUFBRSxDQUFDLEdBQUcsUUFBUSxDQUFDLE1BQU0sRUFBRSxDQUFDLEVBQUUsRUFBRTtZQUNoRCxNQUFNLENBQUMsR0FBRyxDQUFDLFFBQVEsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtZQUM1QyxNQUFNLElBQUksR0FBUyxJQUFJLFlBQUksRUFBRSxDQUFBO1lBQzdCLElBQUksQ0FBQyxVQUFVLENBQUMsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUE7WUFDNUIsTUFBTSxRQUFRLEdBQVMsR0FBRyxDQUFDLE9BQU8sQ0FBQyxJQUFJLENBQUMsU0FBUyxFQUFFLENBQVMsQ0FBQTtZQUM1RCxNQUFNLENBQUMsUUFBUSxDQUFDLFFBQVEsRUFBRSxDQUFDLENBQUMsSUFBSSxDQUFDLFFBQVEsQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFBO1NBQzlDO0lBQ0gsQ0FBQyxDQUFDLENBQUE7SUFFRixJQUFJLENBQUMsVUFBVSxFQUFFLEdBQVMsRUFBRTtRQUMxQixNQUFNLEdBQUcsR0FBWSxJQUFJLGVBQU8sRUFBRSxDQUFBO1FBQ2xDLEdBQUcsQ0FBQyxRQUFRLENBQUMsUUFBUSxDQUFDLENBQUE7UUFDdEIsS0FBSyxJQUFJLENBQUMsR0FBVyxDQUFDLEVBQUUsQ0FBQyxHQUFHLFFBQVEsQ0FBQyxNQUFNLEVBQUUsQ0FBQyxFQUFFLEVBQUU7WUFDaEQsTUFBTSxFQUFFLEdBQVMsSUFBSSxZQUFJLEVBQUUsQ0FBQTtZQUMzQixFQUFFLENBQUMsVUFBVSxDQUFDLFFBQVEsQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFBO1lBQzFCLE1BQU0sQ0FBQyxHQUFHLENBQUMsUUFBUSxDQUFDLEVBQUUsQ0FBQyxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO1lBQ25DLE1BQU0sSUFBSSxHQUFTLElBQUksWUFBSSxFQUFFLENBQUE7WUFDN0IsSUFBSSxDQUFDLFVBQVUsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTtZQUM1QixNQUFNLFFBQVEsR0FBUyxHQUFHLENBQUMsT0FBTyxDQUFDLElBQUksQ0FBQyxTQUFTLEVBQUUsQ0FBUyxDQUFBO1lBQzVELE1BQU0sQ0FBQyxRQUFRLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQyxJQUFJLENBQUMsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUE7U0FDOUM7UUFFRCxHQUFHLENBQUMsUUFBUSxDQUFDLEdBQUcsQ0FBQyxXQUFXLEVBQUUsQ0FBQyxDQUFBO1FBQy9CLEtBQUssSUFBSSxDQUFDLEdBQVcsQ0FBQyxFQUFFLENBQUMsR0FBRyxRQUFRLENBQUMsTUFBTSxFQUFFLENBQUMsRUFBRSxFQUFFO1lBQ2hELE1BQU0sSUFBSSxHQUFTLElBQUksWUFBSSxFQUFFLENBQUE7WUFDN0IsSUFBSSxDQUFDLFVBQVUsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTtZQUM1QixNQUFNLENBQUMsR0FBRyxDQUFDLFFBQVEsQ0FBQyxJQUFJLENBQUMsQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtZQUVyQyxNQUFNLFFBQVEsR0FBUyxHQUFHLENBQUMsT0FBTyxDQUFDLElBQUksQ0FBQyxTQUFTLEVBQUUsQ0FBUyxDQUFBO1lBQzVELE1BQU0sQ0FBQyxRQUFRLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQyxJQUFJLENBQUMsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUE7U0FDOUM7UUFDRCxJQUFJLENBQUMsR0FBVyxHQUFHLENBQUMsU0FBUyxDQUFDLEtBQUssQ0FBQyxDQUFBO1FBQ3BDLElBQUksQ0FBQyxHQUFZLElBQUksZUFBTyxFQUFFLENBQUE7UUFDOUIsQ0FBQyxDQUFDLFdBQVcsQ0FBQyxDQUFDLENBQUMsQ0FBQTtRQUNoQixJQUFJLENBQUMsR0FBVyxHQUFHLENBQUMsU0FBUyxDQUFDLE9BQU8sQ0FBQyxDQUFBO1FBQ3RDLElBQUksQ0FBQyxHQUFZLElBQUksZUFBTyxFQUFFLENBQUE7UUFDOUIsQ0FBQyxDQUFDLFdBQVcsQ0FBQyxDQUFDLENBQUMsQ0FBQTtJQUNsQixDQUFDLENBQUMsQ0FBQTtJQUVGLElBQUksQ0FBQyxrQkFBa0IsRUFBRSxHQUFTLEVBQUU7UUFDbEMsTUFBTSxHQUFHLEdBQVksSUFBSSxlQUFPLEVBQUUsQ0FBQTtRQUNsQyxHQUFHLENBQUMsUUFBUSxDQUFDLFFBQVEsQ0FBQyxDQUFBO1FBQ3RCLE1BQU0sUUFBUSxHQUFTLElBQUksWUFBSSxFQUFFLENBQUE7UUFDakMsUUFBUSxDQUFDLFVBQVUsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTtRQUNoQyxNQUFNLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLEVBQUUsSUFBSSxDQUFDLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQyxJQUFJLENBQUMsUUFBUSxDQUFDLFFBQVEsRUFBRSxDQUFDLENBQUE7UUFDdkUsTUFBTSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQUMsUUFBUSxDQUFDLENBQUMsQ0FBQyxFQUFFLEtBQUssQ0FBQyxDQUFDLENBQUMsYUFBYSxFQUFFLENBQUE7UUFDbkQsTUFBTSxDQUFDLEdBQUcsQ0FBQyxRQUFRLENBQUMsUUFBUSxFQUFFLElBQUksQ0FBQyxDQUFDLE1BQU0sQ0FBQyxDQUFDLElBQUksQ0FBQyxDQUFDLENBQUMsQ0FBQTtRQUNuRCxNQUFNLENBQUMsR0FBRyxDQUFDLFFBQVEsQ0FBQyxRQUFRLEVBQUUsS0FBSyxDQUFDLENBQUMsTUFBTSxDQUFDLENBQUMsSUFBSSxDQUFDLENBQUMsQ0FBQyxDQUFBO0lBQ3RELENBQUMsQ0FBQyxDQUFBO0lBRUYsUUFBUSxDQUFDLGVBQWUsRUFBRSxHQUFTLEVBQUU7UUFDbkMsSUFBSSxHQUFZLENBQUE7UUFDaEIsSUFBSSxLQUFhLENBQUE7UUFDakIsVUFBVSxDQUFDLEdBQUcsRUFBRTtZQUNkLEdBQUcsR0FBRyxJQUFJLGVBQU8sRUFBRSxDQUFBO1lBQ25CLEdBQUcsQ0FBQyxRQUFRLENBQUMsUUFBUSxDQUFDLENBQUE7WUFDdEIsS0FBSyxHQUFHLEdBQUcsQ0FBQyxXQUFXLEVBQUUsQ0FBQTtRQUMzQixDQUFDLENBQUMsQ0FBQTtRQUVGLElBQUksQ0FBQyxRQUFRLEVBQUUsR0FBUyxFQUFFO1lBQ3hCLE1BQU0sUUFBUSxHQUFTLElBQUksWUFBSSxFQUFFLENBQUE7WUFDakMsUUFBUSxDQUFDLFVBQVUsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTtZQUNoQyxNQUFNLENBQUMsR0FBRyxDQUFDLE1BQU0sQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxRQUFRLEVBQUUsQ0FBQyxDQUFDLElBQUksQ0FBQyxRQUFRLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQTtZQUNwRSxNQUFNLENBQUMsR0FBRyxDQUFDLE1BQU0sQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDLGFBQWEsRUFBRSxDQUFBO1lBQy9DLE1BQU0sQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFDLFFBQVEsQ0FBQyxDQUFDLENBQUMsRUFBRSxLQUFLLENBQUMsQ0FBQyxRQUFRLEVBQUUsQ0FBQyxDQUFDLElBQUksQ0FBQyxRQUFRLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQTtZQUN4RSxNQUFNLENBQUMsR0FBRyxDQUFDLE1BQU0sQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxRQUFRLEVBQUUsQ0FBQyxDQUFDLElBQUksQ0FBQyxRQUFRLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQTtRQUN0RSxDQUFDLENBQUMsQ0FBQTtRQUVGLElBQUksQ0FBQyxhQUFhLEVBQUUsR0FBUyxFQUFFO1lBQzdCLE1BQU0sUUFBUSxHQUFTLElBQUksWUFBSSxFQUFFLENBQUE7WUFDakMsUUFBUSxDQUFDLFVBQVUsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTtZQUNoQyxNQUFNLENBQUMsR0FBRyxDQUFDLFdBQVcsQ0FBQyxRQUFRLENBQUMsQ0FBQyxNQUFNLENBQUMsQ0FBQyxJQUFJLENBQUMsQ0FBQyxDQUFDLENBQUE7WUFDaEQsTUFBTSxDQUFDLEdBQUcsQ0FBQyxXQUFXLENBQUMsUUFBUSxDQUFDLENBQUMsTUFBTSxDQUFDLENBQUMsSUFBSSxDQUFDLENBQUMsQ0FBQyxDQUFBO1lBQ2hELE1BQU0sQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFDLFFBQVEsQ0FBQyxDQUFDLENBQUMsRUFBRSxLQUFLLENBQUMsQ0FBQyxRQUFRLEVBQUUsQ0FBQyxDQUFDLElBQUksQ0FBQyxRQUFRLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQTtZQUN4RSxNQUFNLENBQUMsR0FBRyxDQUFDLFdBQVcsQ0FBQyxRQUFRLENBQUMsQ0FBQyxNQUFNLENBQUMsQ0FBQyxJQUFJLENBQUMsQ0FBQyxDQUFDLENBQUE7WUFDaEQsTUFBTSxDQUFDLEdBQUcsQ0FBQyxRQUFRLENBQUMsUUFBUSxFQUFFLEtBQUssQ0FBQyxDQUFDLE1BQU0sQ0FBQyxDQUFDLElBQUksQ0FBQyxDQUFDLENBQUMsQ0FBQTtZQUNwRCxNQUFNLENBQUMsR0FBRyxDQUFDLFdBQVcsQ0FBQyxLQUFLLENBQUMsQ0FBQyxNQUFNLENBQUMsQ0FBQyxJQUFJLENBQUMsQ0FBQyxDQUFDLENBQUE7UUFDL0MsQ0FBQyxDQUFDLENBQUE7UUFFRixJQUFJLENBQUMsWUFBWSxFQUFFLEdBQVMsRUFBRTtZQUM1QixNQUFNLElBQUksR0FBYSxHQUFHLENBQUMsVUFBVSxFQUFFLENBQUE7WUFDdkMsS0FBSyxJQUFJLENBQUMsR0FBVyxDQUFDLEVBQUUsQ0FBQyxHQUFHLEtBQUssQ0FBQyxNQUFNLEVBQUUsQ0FBQyxFQUFFLEVBQUU7Z0JBQzdDLE1BQU0sQ0FBQyxJQUFJLENBQUMsT0FBTyxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQyxTQUFTLEVBQUUsQ0FBQyxDQUFDLENBQUMsR0FBRyxDQUFDLElBQUksQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFBO2FBQ3hEO1FBQ0gsQ0FBQyxDQUFDLENBQUE7UUFFRixJQUFJLENBQUMsYUFBYSxFQUFFLEdBQVMsRUFBRTtZQUM3QixNQUFNLFFBQVEsR0FBVyxHQUFHLENBQUMsV0FBVyxFQUFFLENBQUE7WUFDMUMsTUFBTSxLQUFLLEdBQWEsRUFBRSxDQUFBO1lBQzFCLEtBQUssSUFBSSxDQUFDLEdBQVcsQ0FBQyxFQUFFLENBQUMsR0FBRyxRQUFRLENBQUMsTUFBTSxFQUFFLENBQUMsRUFBRSxFQUFFO2dCQUNoRCxLQUFLLENBQUMsSUFBSSxDQUFDLFFBQVEsQ0FBQyxDQUFDLENBQUMsQ0FBQyxRQUFRLEVBQUUsQ0FBQyxDQUFBO2FBQ25DO1lBQ0QsS0FBSyxJQUFJLENBQUMsR0FBVyxDQUFDLEVBQUUsQ0FBQyxHQUFHLFFBQVEsQ0FBQyxNQUFNLEVBQUUsQ0FBQyxFQUFFLEVBQUU7Z0JBQ2hELE1BQU0sQ0FBQyxLQUFLLENBQUMsT0FBTyxDQUFDLFFBQVEsQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsR0FBRyxDQUFDLElBQUksQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFBO2FBQ2hEO1lBQ0QsTUFBTSxJQUFJLEdBQWEsR0FBRyxDQUFDLFVBQVUsRUFBRSxDQUFBO1lBQ3ZDLE1BQU0sU0FBUyxHQUFXLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxDQUFDLENBQUE7WUFDL0MsTUFBTSxNQUFNLEdBQWEsRUFBRSxDQUFBO1lBQzNCLEtBQUssSUFBSSxDQUFDLEdBQVcsQ0FBQyxFQUFFLENBQUMsR0FBRyxRQUFRLENBQUMsTUFBTSxFQUFFLENBQUMsRUFBRSxFQUFFO2dCQUNoRCxNQUFNLENBQUMsSUFBSSxDQUFDLFNBQVMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxRQUFRLEVBQUUsQ0FBQyxDQUFBO2FBQ3JDO1lBQ0QsS0FBSyxJQUFJLENBQUMsR0FBVyxDQUFDLEVBQUUsQ0FBQyxHQUFHLFFBQVEsQ0FBQyxNQUFNLEVBQUUsQ0FBQyxFQUFFLEVBQUU7Z0JBQ2hELE1BQU0sQ0FBQyxNQUFNLENBQUMsT0FBTyxDQUFDLFFBQVEsQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsR0FBRyxDQUFDLElBQUksQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFBO2FBQ2pEO1FBQ0gsQ0FBQyxDQUFDLENBQUE7UUFFRixJQUFJLENBQUMsdUJBQXVCLEVBQUUsR0FBUyxFQUFFO1lBQ3ZDLElBQUksT0FBaUIsQ0FBQTtZQUNyQixPQUFPLEdBQUcsR0FBRyxDQUFDLFVBQVUsQ0FBQyxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUE7WUFDcEMsTUFBTSxDQUFDLE9BQU8sQ0FBQyxNQUFNLENBQUMsQ0FBQyxJQUFJLENBQUMsQ0FBQyxDQUFDLENBQUE7WUFDOUIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxVQUFVLENBQUMsS0FBSyxDQUFDLENBQUE7WUFDL0IsTUFBTSxDQUFDLE9BQU8sQ0FBQyxNQUFNLENBQUMsQ0FBQyxJQUFJLENBQUMsQ0FBQyxDQUFDLENBQUE7WUFDOUIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxVQUFVLENBQUMsS0FBSyxFQUFFLEtBQUssQ0FBQyxDQUFBO1lBQ3RDLE1BQU0sQ0FBQyxPQUFPLENBQUMsTUFBTSxDQUFDLENBQUMsSUFBSSxDQUFDLENBQUMsQ0FBQyxDQUFBO1FBQ2hDLENBQUMsQ0FBQyxDQUFBO1FBRUYsSUFBSSxDQUFDLG1CQUFtQixFQUFFLEdBQVMsRUFBRTtZQUNuQyxNQUFNLEtBQUssR0FBYSxHQUFHLENBQUMsaUJBQWlCLEVBQUUsQ0FBQTtZQUMvQyxLQUFLLElBQUksQ0FBQyxHQUFXLENBQUMsRUFBRSxDQUFDLEdBQUcsUUFBUSxDQUFDLE1BQU0sRUFBRSxDQUFDLEVBQUUsRUFBRTtnQkFDaEQsTUFBTSxDQUFDLEtBQUssQ0FBQyxPQUFPLENBQUMsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxHQUFHLENBQUMsSUFBSSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUE7YUFDaEQ7WUFDRCxNQUFNLElBQUksR0FBYSxHQUFHLENBQUMsVUFBVSxFQUFFLENBQUE7WUFDdkMsTUFBTSxNQUFNLEdBQWEsR0FBRyxDQUFDLGlCQUFpQixDQUFDLElBQUksQ0FBQyxDQUFBO1lBQ3BELEtBQUssSUFBSSxDQUFDLEdBQVcsQ0FBQyxFQUFFLENBQUMsR0FBRyxRQUFRLENBQUMsTUFBTSxFQUFFLENBQUMsRUFBRSxFQUFFO2dCQUNoRCxNQUFNLENBQUMsTUFBTSxDQUFDLE9BQU8sQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDLEdBQUcsQ0FBQyxJQUFJLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTthQUNqRDtRQUNILENBQUMsQ0FBQyxDQUFBO1FBRUYsSUFBSSxDQUFDLGNBQWMsRUFBRSxHQUFTLEVBQUU7WUFDOUIsTUFBTSxDQUFDLEdBQUcsQ0FBQyxZQUFZLEVBQUUsQ0FBQyxJQUFJLEVBQUUsQ0FBQyxDQUFDLGFBQWEsQ0FBQyxLQUFLLENBQUMsSUFBSSxFQUFFLENBQUMsQ0FBQTtRQUMvRCxDQUFDLENBQUMsQ0FBQTtRQUVGLElBQUksQ0FBQyxZQUFZLEVBQUUsR0FBUyxFQUFFO1lBQzVCLElBQUksUUFBWSxDQUFBO1lBQ2hCLElBQUksUUFBWSxDQUFBO1lBQ2hCLFFBQVEsR0FBRyxJQUFJLGVBQUUsQ0FBQyxDQUFDLENBQUMsQ0FBQTtZQUNwQixRQUFRLEdBQUcsSUFBSSxlQUFFLENBQUMsQ0FBQyxDQUFDLENBQUE7WUFDcEIsS0FBSyxJQUFJLENBQUMsR0FBVyxDQUFDLEVBQUUsQ0FBQyxHQUFHLEtBQUssQ0FBQyxNQUFNLEVBQUUsQ0FBQyxFQUFFLEVBQUU7Z0JBQzdDLE1BQU0sT0FBTyxHQUFHLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQyxVQUFVLEVBQUUsQ0FBQTtnQkFDckMsUUFBUSxHQUFHLFFBQVEsQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFDLFVBQVUsQ0FBQyxLQUFLLEVBQUUsT0FBTyxDQUFDLENBQUMsQ0FBQTtnQkFDdkQsUUFBUSxHQUFHLFFBQVEsQ0FBQyxHQUFHLENBQ3BCLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQyxTQUFTLEVBQW1CLENBQUMsU0FBUyxFQUFFLENBQ25ELENBQUE7YUFDRjtZQUNELE1BQU0sQ0FBQyxRQUFRLENBQUMsRUFBRSxDQUFDLElBQUksZUFBRSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7WUFDekMsTUFBTSxDQUFDLFFBQVEsQ0FBQyxFQUFFLENBQUMsSUFBSSxlQUFFLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtZQUV6QyxRQUFRLEdBQUcsSUFBSSxlQUFFLENBQUMsQ0FBQyxDQUFDLENBQUE7WUFDcEIsUUFBUSxHQUFHLElBQUksZUFBRSxDQUFDLENBQUMsQ0FBQyxDQUFBO1lBQ3BCLE1BQU0sR0FBRyxHQUFPLElBQUEseUJBQU8sR0FBRSxDQUFBO1lBQ3pCLEtBQUssSUFBSSxDQUFDLEdBQVcsQ0FBQyxFQUFFLENBQUMsR0FBRyxLQUFLLENBQUMsTUFBTSxFQUFFLENBQUMsRUFBRSxFQUFFO2dCQUM3QyxNQUFNLE9BQU8sR0FBRyxRQUFRLENBQUMsVUFBVSxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQyxVQUFVLEVBQUUsQ0FBQyxDQUFBO2dCQUMxRCxRQUFRLEdBQUcsUUFBUSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQUMsVUFBVSxDQUFDLEtBQUssRUFBRSxPQUFPLEVBQUUsR0FBRyxDQUFDLENBQUMsQ0FBQTtnQkFDNUQsUUFBUSxHQUFHLFFBQVEsQ0FBQyxHQUFHLENBQ3BCLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQyxTQUFTLEVBQW1CLENBQUMsU0FBUyxFQUFFLENBQ25ELENBQUE7YUFDRjtZQUNELE1BQU0sQ0FBQyxRQUFRLENBQUMsRUFBRSxDQUFDLElBQUksZUFBRSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7WUFDekMsTUFBTSxDQUFDLFFBQVEsQ0FBQyxFQUFFLENBQUMsSUFBSSxlQUFFLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtRQUMzQyxDQUFDLENBQUMsQ0FBQTtRQUVGLElBQUksQ0FBQyxhQUFhLEVBQUUsR0FBUyxFQUFFO1lBQzdCLE1BQU0sUUFBUSxHQUFhLEdBQUcsQ0FBQyxXQUFXLEVBQUUsQ0FBQTtZQUM1QyxLQUFLLElBQUksQ0FBQyxHQUFXLENBQUMsRUFBRSxDQUFDLEdBQUcsS0FBSyxDQUFDLE1BQU0sRUFBRSxDQUFDLEVBQUUsRUFBRTtnQkFDN0MsTUFBTSxDQUFDLFFBQVEsQ0FBQyxDQUFDLFNBQVMsQ0FBQyxLQUFLLENBQUMsQ0FBQyxDQUFDLENBQUMsVUFBVSxFQUFFLENBQUMsQ0FBQTthQUNsRDtZQUNELE1BQU0sU0FBUyxHQUFhLEdBQUcsQ0FBQyxZQUFZLEVBQUUsQ0FBQTtZQUM5QyxNQUFNLENBQUMsR0FBRyxDQUFDLFdBQVcsQ0FBQyxTQUFTLENBQUMsQ0FBQyxDQUFDLE9BQU8sQ0FBQyxHQUFHLENBQUMsV0FBVyxFQUFFLENBQUMsQ0FBQTtRQUMvRCxDQUFDLENBQUMsQ0FBQTtRQUVGLFFBQVEsQ0FBQyxhQUFhLEVBQUUsR0FBUyxFQUFFO1lBQ2pDLElBQUksSUFBYSxDQUFBO1lBQ2pCLElBQUksSUFBYSxDQUFBO1lBQ2pCLElBQUksSUFBYSxDQUFBO1lBQ2pCLElBQUksSUFBYSxDQUFBO1lBQ2pCLElBQUksSUFBYSxDQUFBO1lBQ2pCLElBQUksSUFBYSxDQUFBO1lBQ2pCLElBQUksSUFBYSxDQUFBO1lBQ2pCLElBQUksSUFBYSxDQUFBO1lBQ2pCLGdCQUFnQjtZQUNoQixNQUFNLE9BQU8sR0FBVyxRQUFRLENBQUMsVUFBVSxDQUN6QyxlQUFNLENBQUMsSUFBSSxDQUNULDhPQUE4TyxFQUM5TyxLQUFLLENBQ04sQ0FDRixDQUFBO1lBRUQsVUFBVSxDQUFDLEdBQVMsRUFBRTtnQkFDcEIsSUFBSSxHQUFHLElBQUksZUFBTyxFQUFFLENBQUE7Z0JBQ3BCLElBQUksQ0FBQyxRQUFRLENBQUMsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLEVBQUUsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTtnQkFFekMsSUFBSSxHQUFHLElBQUksZUFBTyxFQUFFLENBQUE7Z0JBQ3BCLElBQUksQ0FBQyxRQUFRLENBQUMsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLEVBQUUsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTtnQkFFekMsSUFBSSxHQUFHLElBQUksZUFBTyxFQUFFLENBQUE7Z0JBQ3BCLElBQUksQ0FBQyxRQUFRLENBQUMsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLEVBQUUsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTtnQkFFekMsSUFBSSxHQUFHLElBQUksZUFBTyxFQUFFLENBQUE7Z0JBQ3BCLElBQUksQ0FBQyxRQUFRLENBQUMsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFBO2dCQUU1QixJQUFJLEdBQUcsSUFBSSxlQUFPLEVBQUUsQ0FBQTtnQkFDcEIsSUFBSSxDQUFDLFFBQVEsQ0FBQyxFQUFFLENBQUMsQ0FBQSxDQUFDLFlBQVk7Z0JBRTlCLElBQUksR0FBRyxJQUFJLGVBQU8sRUFBRSxDQUFBO2dCQUNwQixJQUFJLENBQUMsUUFBUSxDQUFDLFFBQVEsQ0FBQyxDQUFBLENBQUMsK0JBQStCO2dCQUV2RCxJQUFJLEdBQUcsSUFBSSxlQUFPLEVBQUUsQ0FBQTtnQkFDcEIsSUFBSSxDQUFDLFFBQVEsQ0FBQyxDQUFDLE9BQU8sRUFBRSxHQUFHLFFBQVEsQ0FBQyxDQUFDLENBQUEsQ0FBQyw0QkFBNEI7Z0JBRWxFLElBQUksR0FBRyxJQUFJLGVBQU8sRUFBRSxDQUFBO2dCQUNwQixJQUFJLENBQUMsUUFBUSxDQUFDLENBQUMsT0FBTyxDQUFDLENBQUMsQ0FBQSxDQUFDLDhCQUE4QjtZQUN6RCxDQUFDLENBQUMsQ0FBQTtZQUVGLElBQUksQ0FBQyxvQkFBb0IsRUFBRSxHQUFTLEVBQUU7Z0JBQ3BDLE1BQU0sQ0FBQyxHQUFTLEVBQUU7b0JBQ2hCLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLE9BQU8sQ0FBQyxDQUFBO2dCQUNoQyxDQUFDLENBQUMsQ0FBQyxPQUFPLEVBQUUsQ0FBQTtnQkFDWixNQUFNLFFBQVEsR0FBVyxJQUFJLENBQUMsV0FBVyxFQUFFLENBQUE7WUFDN0MsQ0FBQyxDQUFDLENBQUE7WUFFRixJQUFJLENBQUMsY0FBYyxFQUFFLEdBQVMsRUFBRTtnQkFDOUIsSUFBSSxPQUFnQixDQUFBO2dCQUNwQixJQUFJLElBQWEsQ0FBQTtnQkFFakIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLGNBQWMsQ0FBQyxDQUFBO2dCQUMvQyxJQUFJLEdBQUcsY0FBYyxDQUNuQixPQUFPLEVBQ1AsQ0FBQyxJQUFJLENBQUMsRUFDTixDQUFDLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksQ0FBQyxDQUMzQyxDQUFBO2dCQUNELE1BQU0sQ0FBQyxJQUFJLENBQUMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7Z0JBRXZCLE9BQU8sR0FBRyxHQUFHLENBQUMsV0FBVyxDQUFDLElBQUksRUFBRSxjQUFjLENBQUMsQ0FBQTtnQkFDL0MsSUFBSSxHQUFHLGNBQWMsQ0FDbkIsT0FBTyxFQUNQLENBQUMsSUFBSSxDQUFDLEVBQ04sQ0FBQyxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLENBQUMsQ0FDM0MsQ0FBQTtnQkFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO2dCQUV2QixPQUFPLEdBQUcsR0FBRyxDQUFDLFdBQVcsQ0FBQyxJQUFJLEVBQUUsY0FBYyxDQUFDLENBQUE7Z0JBQy9DLElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtnQkFFdkIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLGNBQWMsQ0FBQyxDQUFBO2dCQUMvQyxJQUFJLEdBQUcsY0FBYyxDQUNuQixPQUFPLEVBQ1AsQ0FBQyxJQUFJLENBQUMsRUFDTixDQUFDLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksQ0FBQyxDQUMzQyxDQUFBO2dCQUNELE1BQU0sQ0FBQyxJQUFJLENBQUMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7WUFDekIsQ0FBQyxDQUFDLENBQUE7WUFFRixJQUFJLENBQUMsZ0JBQWdCLEVBQUUsR0FBUyxFQUFFO2dCQUNoQyxJQUFJLE9BQWdCLENBQUE7Z0JBQ3BCLElBQUksSUFBYSxDQUFBO2dCQUVqQixPQUFPLEdBQUcsR0FBRyxDQUFDLFdBQVcsQ0FBQyxJQUFJLEVBQUUsZ0JBQWdCLENBQUMsQ0FBQTtnQkFDakQsSUFBSSxHQUFHLGNBQWMsQ0FDbkIsT0FBTyxFQUNQLENBQUMsSUFBSSxDQUFDLEVBQ04sQ0FBQyxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLENBQUMsQ0FDM0MsQ0FBQTtnQkFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO2dCQUV2QixPQUFPLEdBQUcsR0FBRyxDQUFDLFdBQVcsQ0FBQyxJQUFJLEVBQUUsZ0JBQWdCLENBQUMsQ0FBQTtnQkFDakQsSUFBSSxHQUFHLGNBQWMsQ0FDbkIsT0FBTyxFQUNQLENBQUMsSUFBSSxDQUFDLEVBQ04sQ0FBQyxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLENBQUMsQ0FDM0MsQ0FBQTtnQkFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO2dCQUV2QixPQUFPLEdBQUcsR0FBRyxDQUFDLFdBQVcsQ0FBQyxJQUFJLEVBQUUsZ0JBQWdCLENBQUMsQ0FBQTtnQkFDakQsSUFBSSxHQUFHLGNBQWMsQ0FDbkIsT0FBTyxFQUNQLENBQUMsSUFBSSxDQUFDLEVBQ04sQ0FBQyxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLENBQUMsQ0FDM0MsQ0FBQTtnQkFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO2dCQUV2QixPQUFPLEdBQUcsR0FBRyxDQUFDLFdBQVcsQ0FBQyxJQUFJLEVBQUUsZ0JBQWdCLENBQUMsQ0FBQTtnQkFDakQsSUFBSSxHQUFHLGNBQWMsQ0FDbkIsT0FBTyxFQUNQLENBQUMsSUFBSSxDQUFDLEVBQ04sQ0FBQyxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLENBQUMsQ0FDM0MsQ0FBQTtnQkFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO1lBQ3pCLENBQUMsQ0FBQyxDQUFBO1lBRUYsSUFBSSxDQUFDLGVBQWUsRUFBRSxHQUFTLEVBQUU7Z0JBQy9CLElBQUksT0FBZ0IsQ0FBQTtnQkFDcEIsSUFBSSxJQUFhLENBQUE7Z0JBRWpCLE9BQU8sR0FBRyxHQUFHLENBQUMsV0FBVyxDQUFDLElBQUksRUFBRSxlQUFlLENBQUMsQ0FBQTtnQkFDaEQsSUFBSSxHQUFHLGNBQWMsQ0FDbkIsT0FBTyxFQUNQLENBQUMsSUFBSSxDQUFDLEVBQ04sQ0FBQyxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLENBQUMsQ0FDM0MsQ0FBQTtnQkFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO2dCQUV2QixPQUFPLEdBQUcsR0FBRyxDQUFDLFdBQVcsQ0FBQyxJQUFJLEVBQUUsZUFBZSxDQUFDLENBQUE7Z0JBQ2hELElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtnQkFFdkIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLGVBQWUsQ0FBQyxDQUFBO2dCQUNoRCxJQUFJLEdBQUcsY0FBYyxDQUNuQixPQUFPLEVBQ1AsQ0FBQyxJQUFJLENBQUMsRUFDTixDQUFDLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksQ0FBQyxDQUMzQyxDQUFBO2dCQUNELE1BQU0sQ0FBQyxJQUFJLENBQUMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7Z0JBRXZCLE9BQU8sR0FBRyxHQUFHLENBQUMsV0FBVyxDQUFDLElBQUksRUFBRSxlQUFlLENBQUMsQ0FBQTtnQkFDaEQsSUFBSSxHQUFHLGNBQWMsQ0FDbkIsT0FBTyxFQUNQLENBQUMsSUFBSSxDQUFDLEVBQ04sQ0FBQyxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLENBQUMsQ0FDM0MsQ0FBQTtnQkFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO1lBQ3pCLENBQUMsQ0FBQyxDQUFBO1lBRUYsSUFBSSxDQUFDLGVBQWUsRUFBRSxHQUFTLEVBQUU7Z0JBQy9CLElBQUksT0FBZ0IsQ0FBQTtnQkFDcEIsSUFBSSxJQUFhLENBQUE7Z0JBRWpCLE9BQU8sR0FBRyxHQUFHLENBQUMsV0FBVyxDQUFDLElBQUksRUFBRSxlQUFlLENBQUMsQ0FBQTtnQkFDaEQsSUFBSSxHQUFHLGNBQWMsQ0FDbkIsT0FBTyxFQUNQLENBQUMsSUFBSSxDQUFDLEVBQ04sQ0FBQyxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLENBQUMsQ0FDM0MsQ0FBQTtnQkFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO2dCQUV2QixPQUFPLEdBQUcsR0FBRyxDQUFDLFdBQVcsQ0FBQyxJQUFJLEVBQUUsZUFBZSxDQUFDLENBQUE7Z0JBQ2hELElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtnQkFFdkIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLGVBQWUsQ0FBQyxDQUFBO2dCQUNoRCxJQUFJLEdBQUcsY0FBYyxDQUNuQixPQUFPLEVBQ1AsQ0FBQyxJQUFJLENBQUMsRUFDTixDQUFDLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksQ0FBQyxDQUMzQyxDQUFBO2dCQUNELE1BQU0sQ0FBQyxJQUFJLENBQUMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7Z0JBRXZCLE9BQU8sR0FBRyxHQUFHLENBQUMsV0FBVyxDQUFDLElBQUksRUFBRSxlQUFlLENBQUMsQ0FBQTtnQkFDaEQsSUFBSSxHQUFHLGNBQWMsQ0FDbkIsT0FBTyxFQUNQLENBQUMsSUFBSSxDQUFDLEVBQ04sQ0FBQyxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLENBQUMsQ0FDM0MsQ0FBQTtnQkFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO1lBQ3pCLENBQUMsQ0FBQyxDQUFBO1lBRUYsSUFBSSxDQUFDLE9BQU8sRUFBRSxHQUFTLEVBQUU7Z0JBQ3ZCLElBQUksT0FBZ0IsQ0FBQTtnQkFDcEIsSUFBSSxJQUFhLENBQUE7Z0JBRWpCLE9BQU8sR0FBRyxHQUFHLENBQUMsV0FBVyxDQUFDLElBQUksRUFBRSxPQUFPLENBQUMsQ0FBQTtnQkFDeEMsSUFBSSxHQUFHLGNBQWMsQ0FDbkIsT0FBTyxFQUNQLENBQUMsSUFBSSxDQUFDLEVBQ04sQ0FBQyxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLENBQUMsQ0FDM0MsQ0FBQTtnQkFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO2dCQUV2QixPQUFPLEdBQUcsR0FBRyxDQUFDLFdBQVcsQ0FBQyxJQUFJLEVBQUUsT0FBTyxDQUFDLENBQUE7Z0JBQ3hDLElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtnQkFFdkIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLE9BQU8sQ0FBQyxDQUFBO2dCQUN4QyxJQUFJLEdBQUcsY0FBYyxDQUNuQixPQUFPLEVBQ1AsQ0FBQyxJQUFJLENBQUMsRUFDTixDQUFDLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksQ0FBQyxDQUMzQyxDQUFBO2dCQUNELE1BQU0sQ0FBQyxJQUFJLENBQUMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7Z0JBRXZCLE9BQU8sR0FBRyxHQUFHLENBQUMsV0FBVyxDQUFDLElBQUksRUFBRSxPQUFPLENBQUMsQ0FBQTtnQkFDeEMsSUFBSSxHQUFHLGNBQWMsQ0FDbkIsT0FBTyxFQUNQLENBQUMsSUFBSSxDQUFDLEVBQ04sQ0FBQyxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLENBQUMsQ0FDM0MsQ0FBQTtnQkFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO1lBQ3pCLENBQUMsQ0FBQyxDQUFBO1lBRUYsSUFBSSxDQUFDLGVBQWUsRUFBRSxHQUFTLEVBQUU7Z0JBQy9CLElBQUksT0FBZ0IsQ0FBQTtnQkFDcEIsSUFBSSxJQUFhLENBQUE7Z0JBRWpCLE9BQU8sR0FBRyxHQUFHLENBQUMsV0FBVyxDQUFDLElBQUksRUFBRSxlQUFlLENBQUMsQ0FBQTtnQkFDaEQsSUFBSSxHQUFHLGNBQWMsQ0FDbkIsT0FBTyxFQUNQLENBQUMsSUFBSSxDQUFDLEVBQ04sQ0FBQyxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLENBQUMsQ0FDM0MsQ0FBQTtnQkFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO2dCQUV2QixPQUFPLEdBQUcsR0FBRyxDQUFDLFdBQVcsQ0FBQyxJQUFJLEVBQUUsZUFBZSxDQUFDLENBQUE7Z0JBQ2hELElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtnQkFFdkIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLGVBQWUsQ0FBQyxDQUFBO2dCQUNoRCxJQUFJLEdBQUcsY0FBYyxDQUNuQixPQUFPLEVBQ1AsQ0FBQyxJQUFJLENBQUMsRUFDTixDQUFDLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksQ0FBQyxDQUMzQyxDQUFBO2dCQUNELE1BQU0sQ0FBQyxJQUFJLENBQUMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7Z0JBRXZCLE9BQU8sR0FBRyxHQUFHLENBQUMsV0FBVyxDQUFDLElBQUksRUFBRSxlQUFlLENBQUMsQ0FBQTtnQkFDaEQsSUFBSSxHQUFHLGNBQWMsQ0FDbkIsT0FBTyxFQUNQLENBQUMsSUFBSSxDQUFDLEVBQ04sQ0FBQyxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLENBQUMsQ0FDM0MsQ0FBQTtnQkFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO1lBQ3pCLENBQUMsQ0FBQyxDQUFBO1lBRUYsSUFBSSxDQUFDLGdCQUFnQixFQUFFLEdBQVMsRUFBRTtnQkFDaEMsSUFBSSxPQUFnQixDQUFBO2dCQUNwQixJQUFJLElBQWEsQ0FBQTtnQkFFakIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLGdCQUFnQixDQUFDLENBQUE7Z0JBQ2pELElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtnQkFFdkIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLGdCQUFnQixDQUFDLENBQUE7Z0JBQ2pELElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtnQkFFdkIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLGdCQUFnQixDQUFDLENBQUE7Z0JBQ2pELElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtnQkFFdkIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLGdCQUFnQixDQUFDLENBQUE7Z0JBQ2pELElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtZQUN6QixDQUFDLENBQUMsQ0FBQTtRQUNKLENBQUMsQ0FBQyxDQUFBO0lBQ0osQ0FBQyxDQUFDLENBQUE7QUFDSixDQUFDLENBQUMsQ0FBQSIsInNvdXJjZXNDb250ZW50IjpbImltcG9ydCBCTiBmcm9tIFwiYm4uanNcIlxuaW1wb3J0IHsgQnVmZmVyIH0gZnJvbSBcImJ1ZmZlci9cIlxuaW1wb3J0IEJpblRvb2xzIGZyb20gXCIuLi8uLi8uLi9zcmMvdXRpbHMvYmludG9vbHNcIlxuaW1wb3J0IHsgVVRYTywgVVRYT1NldCB9IGZyb20gXCIuLi8uLi8uLi9zcmMvYXBpcy9wbGF0Zm9ybXZtL3V0eG9zXCJcbmltcG9ydCB7IEFtb3VudE91dHB1dCB9IGZyb20gXCIuLi8uLi8uLi9zcmMvYXBpcy9wbGF0Zm9ybXZtL291dHB1dHNcIlxuaW1wb3J0IHsgVW5peE5vdyB9IGZyb20gXCIuLi8uLi8uLi9zcmMvdXRpbHMvaGVscGVyZnVuY3Rpb25zXCJcbmltcG9ydCB7IFNlcmlhbGl6ZWRFbmNvZGluZyB9IGZyb20gXCIuLi8uLi8uLi9zcmMvdXRpbHNcIlxuXG5jb25zdCBiaW50b29sczogQmluVG9vbHMgPSBCaW5Ub29scy5nZXRJbnN0YW5jZSgpXG5jb25zdCBkaXNwbGF5OiBTZXJpYWxpemVkRW5jb2RpbmcgPSBcImRpc3BsYXlcIlxuXG5kZXNjcmliZShcIlVUWE9cIiwgKCk6IHZvaWQgPT4ge1xuICBjb25zdCB1dHhvaGV4OiBzdHJpbmcgPVxuICAgIFwiMDAwMDM4ZDFiOWYxMTM4NjcyZGE2ZmI2YzM1MTI1NTM5Mjc2YTlhY2MyYTY2OGQ2M2JlYTZiYTNjNzk1ZTJlZGIwZjUwMDAwMDAwMTNlMDdlMzhlMmYyMzEyMWJlODc1NjQxMmMxOGRiNzI0NmExNmQyNmVlOTkzNmYzY2JhMjhiZTE0OWNmZDM1NTgwMDAwMDAwNzAwMDAwMDAwMDAwMDRkZDUwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDEwMDAwMDAwMWEzNmZkMGMyZGJjYWIzMTE3MzFkZGU3ZWYxNTE0YmQyNmZjZGM3NGRcIlxuICBjb25zdCBvdXRwdXRpZHg6IHN0cmluZyA9IFwiMDAwMDAwMDFcIlxuICBjb25zdCBvdXR0eGlkOiBzdHJpbmcgPVxuICAgIFwiMzhkMWI5ZjExMzg2NzJkYTZmYjZjMzUxMjU1MzkyNzZhOWFjYzJhNjY4ZDYzYmVhNmJhM2M3OTVlMmVkYjBmNVwiXG4gIGNvbnN0IG91dGFpZDogc3RyaW5nID1cbiAgICBcIjNlMDdlMzhlMmYyMzEyMWJlODc1NjQxMmMxOGRiNzI0NmExNmQyNmVlOTkzNmYzY2JhMjhiZTE0OWNmZDM1NThcIlxuICBjb25zdCB1dHhvYnVmZjogQnVmZmVyID0gQnVmZmVyLmZyb20odXR4b2hleCwgXCJoZXhcIilcblxuICAvLyBQYXltZW50XG4gIGNvbnN0IE9QVVRYT3N0cjogc3RyaW5nID0gYmludG9vbHMuY2I1OEVuY29kZSh1dHhvYnVmZilcbiAgLy8gXCJVOXJGZ0s1ampkWG1WOGs1dHBxZVhraW16ck4zbzllQ0NjWGVzeWhNQkJadTlNUUpDRFREbzVXbjVwc0t2ekpWTUpwaU1iZGtmRFhrcDdzS1pkZGZDWmR4cHVEbXlOeTdWRmthMTl6TVc0amN6NkRSUXZOZkEya3ZKWUtrOTZ6Yzd1aXpncDNpMkZZV3JCOG1yMXNQSjhvUDlUaDY0R1E1eUhkOFwiXG5cbiAgLy8gaW1wbGllcyBmcm9tU3RyaW5nIGFuZCBmcm9tQnVmZmVyXG4gIHRlc3QoXCJDcmVhdGlvblwiLCAoKTogdm9pZCA9PiB7XG4gICAgY29uc3QgdTE6IFVUWE8gPSBuZXcgVVRYTygpXG4gICAgdTEuZnJvbUJ1ZmZlcih1dHhvYnVmZilcbiAgICBjb25zdCB1MWhleDogc3RyaW5nID0gdTEudG9CdWZmZXIoKS50b1N0cmluZyhcImhleFwiKVxuICAgIGV4cGVjdCh1MWhleCkudG9CZSh1dHhvaGV4KVxuICB9KVxuXG4gIHRlc3QoXCJFbXB0eSBDcmVhdGlvblwiLCAoKTogdm9pZCA9PiB7XG4gICAgY29uc3QgdTE6IFVUWE8gPSBuZXcgVVRYTygpXG4gICAgZXhwZWN0KCgpID0+IHtcbiAgICAgIHUxLnRvQnVmZmVyKClcbiAgICB9KS50b1Rocm93KClcbiAgfSlcblxuICB0ZXN0KFwiQ3JlYXRpb24gb2YgVHlwZVwiLCAoKTogdm9pZCA9PiB7XG4gICAgY29uc3Qgb3A6IFVUWE8gPSBuZXcgVVRYTygpXG4gICAgb3AuZnJvbVN0cmluZyhPUFVUWE9zdHIpXG4gICAgZXhwZWN0KG9wLmdldE91dHB1dCgpLmdldE91dHB1dElEKCkpLnRvQmUoNylcbiAgfSlcblxuICBkZXNjcmliZShcIkZ1bnRpb25hbGl0eVwiLCAoKTogdm9pZCA9PiB7XG4gICAgY29uc3QgdTE6IFVUWE8gPSBuZXcgVVRYTygpXG4gICAgdTEuZnJvbUJ1ZmZlcih1dHhvYnVmZilcbiAgICBjb25zdCB1MWhleDogc3RyaW5nID0gdTEudG9CdWZmZXIoKS50b1N0cmluZyhcImhleFwiKVxuICAgIHRlc3QoXCJnZXRBc3NldElEIE5vbkNBXCIsICgpOiB2b2lkID0+IHtcbiAgICAgIGNvbnN0IGFzc2V0SUQ6IEJ1ZmZlciA9IHUxLmdldEFzc2V0SUQoKVxuICAgICAgZXhwZWN0KGFzc2V0SUQudG9TdHJpbmcoXCJoZXhcIiwgMCwgYXNzZXRJRC5sZW5ndGgpKS50b0JlKG91dGFpZClcbiAgICB9KVxuICAgIHRlc3QoXCJnZXRUeElEXCIsICgpOiB2b2lkID0+IHtcbiAgICAgIGNvbnN0IHR4aWQ6IEJ1ZmZlciA9IHUxLmdldFR4SUQoKVxuICAgICAgZXhwZWN0KHR4aWQudG9TdHJpbmcoXCJoZXhcIiwgMCwgdHhpZC5sZW5ndGgpKS50b0JlKG91dHR4aWQpXG4gICAgfSlcbiAgICB0ZXN0KFwiZ2V0T3V0cHV0SWR4XCIsICgpOiB2b2lkID0+IHtcbiAgICAgIGNvbnN0IHR4aWR4OiBCdWZmZXIgPSB1MS5nZXRPdXRwdXRJZHgoKVxuICAgICAgZXhwZWN0KHR4aWR4LnRvU3RyaW5nKFwiaGV4XCIsIDAsIHR4aWR4Lmxlbmd0aCkpLnRvQmUob3V0cHV0aWR4KVxuICAgIH0pXG4gICAgdGVzdChcImdldFVUWE9JRFwiLCAoKTogdm9pZCA9PiB7XG4gICAgICBjb25zdCB0eGlkOiBCdWZmZXIgPSBCdWZmZXIuZnJvbShvdXR0eGlkLCBcImhleFwiKVxuICAgICAgY29uc3QgdHhpZHg6IEJ1ZmZlciA9IEJ1ZmZlci5mcm9tKG91dHB1dGlkeCwgXCJoZXhcIilcbiAgICAgIGNvbnN0IHV0eG9pZDogc3RyaW5nID0gYmludG9vbHMuYnVmZmVyVG9CNTgoQnVmZmVyLmNvbmNhdChbdHhpZCwgdHhpZHhdKSlcbiAgICAgIGV4cGVjdCh1MS5nZXRVVFhPSUQoKSkudG9CZSh1dHhvaWQpXG4gICAgfSlcbiAgICB0ZXN0KFwidG9TdHJpbmdcIiwgKCk6IHZvaWQgPT4ge1xuICAgICAgY29uc3Qgc2VyaWFsaXplZDogc3RyaW5nID0gdTEudG9TdHJpbmcoKVxuICAgICAgZXhwZWN0KHNlcmlhbGl6ZWQpLnRvQmUoYmludG9vbHMuY2I1OEVuY29kZSh1dHhvYnVmZikpXG4gICAgfSlcbiAgfSlcbn0pXG5cbmNvbnN0IHNldE1lcmdlVGVzdGVyID0gKFxuICBpbnB1dDogVVRYT1NldCxcbiAgZXF1YWw6IFVUWE9TZXRbXSxcbiAgbm90RXF1YWw6IFVUWE9TZXRbXVxuKTogYm9vbGVhbiA9PiB7XG4gIGNvbnN0IGluc3RyOiBzdHJpbmcgPSBKU09OLnN0cmluZ2lmeShpbnB1dC5nZXRVVFhPSURzKCkuc29ydCgpKVxuICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgZXF1YWwubGVuZ3RoOyBpKyspIHtcbiAgICBpZiAoSlNPTi5zdHJpbmdpZnkoZXF1YWxbaV0uZ2V0VVRYT0lEcygpLnNvcnQoKSkgIT0gaW5zdHIpIHtcbiAgICAgIHJldHVybiBmYWxzZVxuICAgIH1cbiAgfVxuXG4gIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCBub3RFcXVhbC5sZW5ndGg7IGkrKykge1xuICAgIGlmIChKU09OLnN0cmluZ2lmeShub3RFcXVhbFtpXS5nZXRVVFhPSURzKCkuc29ydCgpKSA9PSBpbnN0cikge1xuICAgICAgcmV0dXJuIGZhbHNlXG4gICAgfVxuICB9XG4gIHJldHVybiB0cnVlXG59XG5cbmRlc2NyaWJlKFwiVVRYT1NldFwiLCAoKTogdm9pZCA9PiB7XG4gIGNvbnN0IHV0eG9zdHJzOiBzdHJpbmdbXSA9IFtcbiAgICBiaW50b29scy5jYjU4RW5jb2RlKFxuICAgICAgQnVmZmVyLmZyb20oXG4gICAgICAgIFwiMDAwMDM4ZDFiOWYxMTM4NjcyZGE2ZmI2YzM1MTI1NTM5Mjc2YTlhY2MyYTY2OGQ2M2JlYTZiYTNjNzk1ZTJlZGIwZjUwMDAwMDAwMTNlMDdlMzhlMmYyMzEyMWJlODc1NjQxMmMxOGRiNzI0NmExNmQyNmVlOTkzNmYzY2JhMjhiZTE0OWNmZDM1NTgwMDAwMDAwNzAwMDAwMDAwMDAwMDRkZDUwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDEwMDAwMDAwMWEzNmZkMGMyZGJjYWIzMTE3MzFkZGU3ZWYxNTE0YmQyNmZjZGM3NGRcIixcbiAgICAgICAgXCJoZXhcIlxuICAgICAgKVxuICAgICksXG4gICAgYmludG9vbHMuY2I1OEVuY29kZShcbiAgICAgIEJ1ZmZlci5mcm9tKFxuICAgICAgICBcIjAwMDBjM2U0ODIzNTcxNTg3ZmUyYmRmYzUwMjY4OWY1YTgyMzhiOWQwZWE3ZjMyNzcxMjRkMTZhZjlkZTBkMmQ5OTExMDAwMDAwMDAzZTA3ZTM4ZTJmMjMxMjFiZTg3NTY0MTJjMThkYjcyNDZhMTZkMjZlZTk5MzZmM2NiYTI4YmUxNDljZmQzNTU4MDAwMDAwMDcwMDAwMDAwMDAwMDAwMDE5MDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAxMDAwMDAwMDFlMWI2YjZhNGJhZDk0ZDJlM2YyMDczMDM3OWI5YmNkNmYxNzYzMThlXCIsXG4gICAgICAgIFwiaGV4XCJcbiAgICAgIClcbiAgICApLFxuICAgIGJpbnRvb2xzLmNiNThFbmNvZGUoXG4gICAgICBCdWZmZXIuZnJvbShcbiAgICAgICAgXCIwMDAwZjI5ZGJhNjFmZGE4ZDU3YTkxMWU3Zjg4MTBmOTM1YmRlODEwZDNmOGQ0OTU0MDQ2ODViZGI4ZDlkODU0NWU4NjAwMDAwMDAwM2UwN2UzOGUyZjIzMTIxYmU4NzU2NDEyYzE4ZGI3MjQ2YTE2ZDI2ZWU5OTM2ZjNjYmEyOGJlMTQ5Y2ZkMzU1ODAwMDAwMDA3MDAwMDAwMDAwMDAwMDAxOTAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMTAwMDAwMDAxZTFiNmI2YTRiYWQ5NGQyZTNmMjA3MzAzNzliOWJjZDZmMTc2MzE4ZVwiLFxuICAgICAgICBcImhleFwiXG4gICAgICApXG4gICAgKVxuICBdXG4gIGNvbnN0IGFkZHJzOiBCdWZmZXJbXSA9IFtcbiAgICBiaW50b29scy5jYjU4RGVjb2RlKFwiRnVCNkx3MkQ2Mk51TTh6cEdMQTRBdmVwcTdlR3NaUmlHXCIpLFxuICAgIGJpbnRvb2xzLmNiNThEZWNvZGUoXCJNYVR2S0djY2JZekN4ekJrSnBiMnpIVzdFMVdSZVpxQjhcIilcbiAgXVxuICB0ZXN0KFwiQ3JlYXRpb25cIiwgKCk6IHZvaWQgPT4ge1xuICAgIGNvbnN0IHNldDogVVRYT1NldCA9IG5ldyBVVFhPU2V0KClcbiAgICBzZXQuYWRkKHV0eG9zdHJzWzBdKVxuICAgIGNvbnN0IHV0eG86IFVUWE8gPSBuZXcgVVRYTygpXG4gICAgdXR4by5mcm9tU3RyaW5nKHV0eG9zdHJzWzBdKVxuICAgIGNvbnN0IHNldEFycmF5OiBVVFhPW10gPSBzZXQuZ2V0QWxsVVRYT3MoKVxuICAgIGV4cGVjdCh1dHhvLnRvU3RyaW5nKCkpLnRvQmUoc2V0QXJyYXlbMF0udG9TdHJpbmcoKSlcbiAgfSlcblxuICB0ZXN0KFwiU2VyaWFsaXphdGlvblwiLCAoKTogdm9pZCA9PiB7XG4gICAgY29uc3Qgc2V0OiBVVFhPU2V0ID0gbmV3IFVUWE9TZXQoKVxuICAgIHNldC5hZGRBcnJheShbLi4udXR4b3N0cnNdKVxuICAgIGxldCBzZXRvYmo6IG9iamVjdCA9IHNldC5zZXJpYWxpemUoXCJjYjU4XCIpXG4gICAgbGV0IHNldHN0cjogc3RyaW5nID0gSlNPTi5zdHJpbmdpZnkoc2V0b2JqKVxuICAgIGxldCBzZXQybmV3b2JqOiBvYmplY3QgPSBKU09OLnBhcnNlKHNldHN0cilcbiAgICBsZXQgc2V0MjogVVRYT1NldCA9IG5ldyBVVFhPU2V0KClcbiAgICBzZXQyLmRlc2VyaWFsaXplKHNldDJuZXdvYmosIFwiY2I1OFwiKVxuICAgIGxldCBzZXQyb2JqOiBvYmplY3QgPSBzZXQyLnNlcmlhbGl6ZShcImNiNThcIilcbiAgICBsZXQgc2V0MnN0cjogc3RyaW5nID0gSlNPTi5zdHJpbmdpZnkoc2V0Mm9iailcbiAgICBleHBlY3Qoc2V0Mi5nZXRBbGxVVFhPU3RyaW5ncygpLnNvcnQoKS5qb2luKFwiLFwiKSkudG9CZShcbiAgICAgIHNldC5nZXRBbGxVVFhPU3RyaW5ncygpLnNvcnQoKS5qb2luKFwiLFwiKVxuICAgIClcbiAgfSlcblxuICB0ZXN0KFwiTXV0bGlwbGUgYWRkXCIsICgpOiB2b2lkID0+IHtcbiAgICBjb25zdCBzZXQ6IFVUWE9TZXQgPSBuZXcgVVRYT1NldCgpXG4gICAgLy8gZmlyc3QgYWRkXG4gICAgZm9yIChsZXQgaTogbnVtYmVyID0gMDsgaSA8IHV0eG9zdHJzLmxlbmd0aDsgaSsrKSB7XG4gICAgICBzZXQuYWRkKHV0eG9zdHJzW2ldKVxuICAgIH1cbiAgICAvLyB0aGUgdmVyaWZ5IChkbyB0aGVzZSBzdGVwcyBzZXBhcmF0ZSB0byBlbnN1cmUgbm8gb3ZlcndyaXRlcylcbiAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgdXR4b3N0cnMubGVuZ3RoOyBpKyspIHtcbiAgICAgIGV4cGVjdChzZXQuaW5jbHVkZXModXR4b3N0cnNbaV0pKS50b0JlKHRydWUpXG4gICAgICBjb25zdCB1dHhvOiBVVFhPID0gbmV3IFVUWE8oKVxuICAgICAgdXR4by5mcm9tU3RyaW5nKHV0eG9zdHJzW2ldKVxuICAgICAgY29uc3QgdmVyaXV0eG86IFVUWE8gPSBzZXQuZ2V0VVRYTyh1dHhvLmdldFVUWE9JRCgpKSBhcyBVVFhPXG4gICAgICBleHBlY3QodmVyaXV0eG8udG9TdHJpbmcoKSkudG9CZSh1dHhvc3Ryc1tpXSlcbiAgICB9XG4gIH0pXG5cbiAgdGVzdChcImFkZEFycmF5XCIsICgpOiB2b2lkID0+IHtcbiAgICBjb25zdCBzZXQ6IFVUWE9TZXQgPSBuZXcgVVRYT1NldCgpXG4gICAgc2V0LmFkZEFycmF5KHV0eG9zdHJzKVxuICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCB1dHhvc3Rycy5sZW5ndGg7IGkrKykge1xuICAgICAgY29uc3QgZTE6IFVUWE8gPSBuZXcgVVRYTygpXG4gICAgICBlMS5mcm9tU3RyaW5nKHV0eG9zdHJzW2ldKVxuICAgICAgZXhwZWN0KHNldC5pbmNsdWRlcyhlMSkpLnRvQmUodHJ1ZSlcbiAgICAgIGNvbnN0IHV0eG86IFVUWE8gPSBuZXcgVVRYTygpXG4gICAgICB1dHhvLmZyb21TdHJpbmcodXR4b3N0cnNbaV0pXG4gICAgICBjb25zdCB2ZXJpdXR4bzogVVRYTyA9IHNldC5nZXRVVFhPKHV0eG8uZ2V0VVRYT0lEKCkpIGFzIFVUWE9cbiAgICAgIGV4cGVjdCh2ZXJpdXR4by50b1N0cmluZygpKS50b0JlKHV0eG9zdHJzW2ldKVxuICAgIH1cblxuICAgIHNldC5hZGRBcnJheShzZXQuZ2V0QWxsVVRYT3MoKSlcbiAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgdXR4b3N0cnMubGVuZ3RoOyBpKyspIHtcbiAgICAgIGNvbnN0IHV0eG86IFVUWE8gPSBuZXcgVVRYTygpXG4gICAgICB1dHhvLmZyb21TdHJpbmcodXR4b3N0cnNbaV0pXG4gICAgICBleHBlY3Qoc2V0LmluY2x1ZGVzKHV0eG8pKS50b0JlKHRydWUpXG5cbiAgICAgIGNvbnN0IHZlcml1dHhvOiBVVFhPID0gc2V0LmdldFVUWE8odXR4by5nZXRVVFhPSUQoKSkgYXMgVVRYT1xuICAgICAgZXhwZWN0KHZlcml1dHhvLnRvU3RyaW5nKCkpLnRvQmUodXR4b3N0cnNbaV0pXG4gICAgfVxuICAgIGxldCBvOiBvYmplY3QgPSBzZXQuc2VyaWFsaXplKFwiaGV4XCIpXG4gICAgbGV0IHM6IFVUWE9TZXQgPSBuZXcgVVRYT1NldCgpXG4gICAgcy5kZXNlcmlhbGl6ZShvKVxuICAgIGxldCB0OiBvYmplY3QgPSBzZXQuc2VyaWFsaXplKGRpc3BsYXkpXG4gICAgbGV0IHI6IFVUWE9TZXQgPSBuZXcgVVRYT1NldCgpXG4gICAgci5kZXNlcmlhbGl6ZSh0KVxuICB9KVxuXG4gIHRlc3QoXCJvdmVyd3JpdGluZyBVVFhPXCIsICgpOiB2b2lkID0+IHtcbiAgICBjb25zdCBzZXQ6IFVUWE9TZXQgPSBuZXcgVVRYT1NldCgpXG4gICAgc2V0LmFkZEFycmF5KHV0eG9zdHJzKVxuICAgIGNvbnN0IHRlc3R1dHhvOiBVVFhPID0gbmV3IFVUWE8oKVxuICAgIHRlc3R1dHhvLmZyb21TdHJpbmcodXR4b3N0cnNbMF0pXG4gICAgZXhwZWN0KHNldC5hZGQodXR4b3N0cnNbMF0sIHRydWUpLnRvU3RyaW5nKCkpLnRvQmUodGVzdHV0eG8udG9TdHJpbmcoKSlcbiAgICBleHBlY3Qoc2V0LmFkZCh1dHhvc3Ryc1swXSwgZmFsc2UpKS50b0JlVW5kZWZpbmVkKClcbiAgICBleHBlY3Qoc2V0LmFkZEFycmF5KHV0eG9zdHJzLCB0cnVlKS5sZW5ndGgpLnRvQmUoMylcbiAgICBleHBlY3Qoc2V0LmFkZEFycmF5KHV0eG9zdHJzLCBmYWxzZSkubGVuZ3RoKS50b0JlKDApXG4gIH0pXG5cbiAgZGVzY3JpYmUoXCJGdW5jdGlvbmFsaXR5XCIsICgpOiB2b2lkID0+IHtcbiAgICBsZXQgc2V0OiBVVFhPU2V0XG4gICAgbGV0IHV0eG9zOiBVVFhPW11cbiAgICBiZWZvcmVFYWNoKCgpID0+IHtcbiAgICAgIHNldCA9IG5ldyBVVFhPU2V0KClcbiAgICAgIHNldC5hZGRBcnJheSh1dHhvc3RycylcbiAgICAgIHV0eG9zID0gc2V0LmdldEFsbFVUWE9zKClcbiAgICB9KVxuXG4gICAgdGVzdChcInJlbW92ZVwiLCAoKTogdm9pZCA9PiB7XG4gICAgICBjb25zdCB0ZXN0dXR4bzogVVRYTyA9IG5ldyBVVFhPKClcbiAgICAgIHRlc3R1dHhvLmZyb21TdHJpbmcodXR4b3N0cnNbMF0pXG4gICAgICBleHBlY3Qoc2V0LnJlbW92ZSh1dHhvc3Ryc1swXSkudG9TdHJpbmcoKSkudG9CZSh0ZXN0dXR4by50b1N0cmluZygpKVxuICAgICAgZXhwZWN0KHNldC5yZW1vdmUodXR4b3N0cnNbMF0pKS50b0JlVW5kZWZpbmVkKClcbiAgICAgIGV4cGVjdChzZXQuYWRkKHV0eG9zdHJzWzBdLCBmYWxzZSkudG9TdHJpbmcoKSkudG9CZSh0ZXN0dXR4by50b1N0cmluZygpKVxuICAgICAgZXhwZWN0KHNldC5yZW1vdmUodXR4b3N0cnNbMF0pLnRvU3RyaW5nKCkpLnRvQmUodGVzdHV0eG8udG9TdHJpbmcoKSlcbiAgICB9KVxuXG4gICAgdGVzdChcInJlbW92ZUFycmF5XCIsICgpOiB2b2lkID0+IHtcbiAgICAgIGNvbnN0IHRlc3R1dHhvOiBVVFhPID0gbmV3IFVUWE8oKVxuICAgICAgdGVzdHV0eG8uZnJvbVN0cmluZyh1dHhvc3Ryc1swXSlcbiAgICAgIGV4cGVjdChzZXQucmVtb3ZlQXJyYXkodXR4b3N0cnMpLmxlbmd0aCkudG9CZSgzKVxuICAgICAgZXhwZWN0KHNldC5yZW1vdmVBcnJheSh1dHhvc3RycykubGVuZ3RoKS50b0JlKDApXG4gICAgICBleHBlY3Qoc2V0LmFkZCh1dHhvc3Ryc1swXSwgZmFsc2UpLnRvU3RyaW5nKCkpLnRvQmUodGVzdHV0eG8udG9TdHJpbmcoKSlcbiAgICAgIGV4cGVjdChzZXQucmVtb3ZlQXJyYXkodXR4b3N0cnMpLmxlbmd0aCkudG9CZSgxKVxuICAgICAgZXhwZWN0KHNldC5hZGRBcnJheSh1dHhvc3RycywgZmFsc2UpLmxlbmd0aCkudG9CZSgzKVxuICAgICAgZXhwZWN0KHNldC5yZW1vdmVBcnJheSh1dHhvcykubGVuZ3RoKS50b0JlKDMpXG4gICAgfSlcblxuICAgIHRlc3QoXCJnZXRVVFhPSURzXCIsICgpOiB2b2lkID0+IHtcbiAgICAgIGNvbnN0IHVpZHM6IHN0cmluZ1tdID0gc2V0LmdldFVUWE9JRHMoKVxuICAgICAgZm9yIChsZXQgaTogbnVtYmVyID0gMDsgaSA8IHV0eG9zLmxlbmd0aDsgaSsrKSB7XG4gICAgICAgIGV4cGVjdCh1aWRzLmluZGV4T2YodXR4b3NbaV0uZ2V0VVRYT0lEKCkpKS5ub3QudG9CZSgtMSlcbiAgICAgIH1cbiAgICB9KVxuXG4gICAgdGVzdChcImdldEFsbFVUWE9zXCIsICgpOiB2b2lkID0+IHtcbiAgICAgIGNvbnN0IGFsbHV0eG9zOiBVVFhPW10gPSBzZXQuZ2V0QWxsVVRYT3MoKVxuICAgICAgY29uc3QgdXN0cnM6IHN0cmluZ1tdID0gW11cbiAgICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCBhbGx1dHhvcy5sZW5ndGg7IGkrKykge1xuICAgICAgICB1c3Rycy5wdXNoKGFsbHV0eG9zW2ldLnRvU3RyaW5nKCkpXG4gICAgICB9XG4gICAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgdXR4b3N0cnMubGVuZ3RoOyBpKyspIHtcbiAgICAgICAgZXhwZWN0KHVzdHJzLmluZGV4T2YodXR4b3N0cnNbaV0pKS5ub3QudG9CZSgtMSlcbiAgICAgIH1cbiAgICAgIGNvbnN0IHVpZHM6IHN0cmluZ1tdID0gc2V0LmdldFVUWE9JRHMoKVxuICAgICAgY29uc3QgYWxsdXR4b3MyOiBVVFhPW10gPSBzZXQuZ2V0QWxsVVRYT3ModWlkcylcbiAgICAgIGNvbnN0IHVzdHJzMjogc3RyaW5nW10gPSBbXVxuICAgICAgZm9yIChsZXQgaTogbnVtYmVyID0gMDsgaSA8IGFsbHV0eG9zLmxlbmd0aDsgaSsrKSB7XG4gICAgICAgIHVzdHJzMi5wdXNoKGFsbHV0eG9zMltpXS50b1N0cmluZygpKVxuICAgICAgfVxuICAgICAgZm9yIChsZXQgaTogbnVtYmVyID0gMDsgaSA8IHV0eG9zdHJzLmxlbmd0aDsgaSsrKSB7XG4gICAgICAgIGV4cGVjdCh1c3RyczIuaW5kZXhPZih1dHhvc3Ryc1tpXSkpLm5vdC50b0JlKC0xKVxuICAgICAgfVxuICAgIH0pXG5cbiAgICB0ZXN0KFwiZ2V0VVRYT0lEcyBCeSBBZGRyZXNzXCIsICgpOiB2b2lkID0+IHtcbiAgICAgIGxldCB1dHhvaWRzOiBzdHJpbmdbXVxuICAgICAgdXR4b2lkcyA9IHNldC5nZXRVVFhPSURzKFthZGRyc1swXV0pXG4gICAgICBleHBlY3QodXR4b2lkcy5sZW5ndGgpLnRvQmUoMSlcbiAgICAgIHV0eG9pZHMgPSBzZXQuZ2V0VVRYT0lEcyhhZGRycylcbiAgICAgIGV4cGVjdCh1dHhvaWRzLmxlbmd0aCkudG9CZSgzKVxuICAgICAgdXR4b2lkcyA9IHNldC5nZXRVVFhPSURzKGFkZHJzLCBmYWxzZSlcbiAgICAgIGV4cGVjdCh1dHhvaWRzLmxlbmd0aCkudG9CZSgzKVxuICAgIH0pXG5cbiAgICB0ZXN0KFwiZ2V0QWxsVVRYT1N0cmluZ3NcIiwgKCk6IHZvaWQgPT4ge1xuICAgICAgY29uc3QgdXN0cnM6IHN0cmluZ1tdID0gc2V0LmdldEFsbFVUWE9TdHJpbmdzKClcbiAgICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCB1dHhvc3Rycy5sZW5ndGg7IGkrKykge1xuICAgICAgICBleHBlY3QodXN0cnMuaW5kZXhPZih1dHhvc3Ryc1tpXSkpLm5vdC50b0JlKC0xKVxuICAgICAgfVxuICAgICAgY29uc3QgdWlkczogc3RyaW5nW10gPSBzZXQuZ2V0VVRYT0lEcygpXG4gICAgICBjb25zdCB1c3RyczI6IHN0cmluZ1tdID0gc2V0LmdldEFsbFVUWE9TdHJpbmdzKHVpZHMpXG4gICAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgdXR4b3N0cnMubGVuZ3RoOyBpKyspIHtcbiAgICAgICAgZXhwZWN0KHVzdHJzMi5pbmRleE9mKHV0eG9zdHJzW2ldKSkubm90LnRvQmUoLTEpXG4gICAgICB9XG4gICAgfSlcblxuICAgIHRlc3QoXCJnZXRBZGRyZXNzZXNcIiwgKCk6IHZvaWQgPT4ge1xuICAgICAgZXhwZWN0KHNldC5nZXRBZGRyZXNzZXMoKS5zb3J0KCkpLnRvU3RyaWN0RXF1YWwoYWRkcnMuc29ydCgpKVxuICAgIH0pXG5cbiAgICB0ZXN0KFwiZ2V0QmFsYW5jZVwiLCAoKTogdm9pZCA9PiB7XG4gICAgICBsZXQgYmFsYW5jZTE6IEJOXG4gICAgICBsZXQgYmFsYW5jZTI6IEJOXG4gICAgICBiYWxhbmNlMSA9IG5ldyBCTigwKVxuICAgICAgYmFsYW5jZTIgPSBuZXcgQk4oMClcbiAgICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCB1dHhvcy5sZW5ndGg7IGkrKykge1xuICAgICAgICBjb25zdCBhc3NldElEID0gdXR4b3NbaV0uZ2V0QXNzZXRJRCgpXG4gICAgICAgIGJhbGFuY2UxID0gYmFsYW5jZTEuYWRkKHNldC5nZXRCYWxhbmNlKGFkZHJzLCBhc3NldElEKSlcbiAgICAgICAgYmFsYW5jZTIgPSBiYWxhbmNlMi5hZGQoXG4gICAgICAgICAgKHV0eG9zW2ldLmdldE91dHB1dCgpIGFzIEFtb3VudE91dHB1dCkuZ2V0QW1vdW50KClcbiAgICAgICAgKVxuICAgICAgfVxuICAgICAgZXhwZWN0KGJhbGFuY2UxLmd0KG5ldyBCTigwKSkpLnRvQmUodHJ1ZSlcbiAgICAgIGV4cGVjdChiYWxhbmNlMi5ndChuZXcgQk4oMCkpKS50b0JlKHRydWUpXG5cbiAgICAgIGJhbGFuY2UxID0gbmV3IEJOKDApXG4gICAgICBiYWxhbmNlMiA9IG5ldyBCTigwKVxuICAgICAgY29uc3Qgbm93OiBCTiA9IFVuaXhOb3coKVxuICAgICAgZm9yIChsZXQgaTogbnVtYmVyID0gMDsgaSA8IHV0eG9zLmxlbmd0aDsgaSsrKSB7XG4gICAgICAgIGNvbnN0IGFzc2V0SUQgPSBiaW50b29scy5jYjU4RW5jb2RlKHV0eG9zW2ldLmdldEFzc2V0SUQoKSlcbiAgICAgICAgYmFsYW5jZTEgPSBiYWxhbmNlMS5hZGQoc2V0LmdldEJhbGFuY2UoYWRkcnMsIGFzc2V0SUQsIG5vdykpXG4gICAgICAgIGJhbGFuY2UyID0gYmFsYW5jZTIuYWRkKFxuICAgICAgICAgICh1dHhvc1tpXS5nZXRPdXRwdXQoKSBhcyBBbW91bnRPdXRwdXQpLmdldEFtb3VudCgpXG4gICAgICAgIClcbiAgICAgIH1cbiAgICAgIGV4cGVjdChiYWxhbmNlMS5ndChuZXcgQk4oMCkpKS50b0JlKHRydWUpXG4gICAgICBleHBlY3QoYmFsYW5jZTIuZ3QobmV3IEJOKDApKSkudG9CZSh0cnVlKVxuICAgIH0pXG5cbiAgICB0ZXN0KFwiZ2V0QXNzZXRJRHNcIiwgKCk6IHZvaWQgPT4ge1xuICAgICAgY29uc3QgYXNzZXRJRHM6IEJ1ZmZlcltdID0gc2V0LmdldEFzc2V0SURzKClcbiAgICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCB1dHhvcy5sZW5ndGg7IGkrKykge1xuICAgICAgICBleHBlY3QoYXNzZXRJRHMpLnRvQ29udGFpbih1dHhvc1tpXS5nZXRBc3NldElEKCkpXG4gICAgICB9XG4gICAgICBjb25zdCBhZGRyZXNzZXM6IEJ1ZmZlcltdID0gc2V0LmdldEFkZHJlc3NlcygpXG4gICAgICBleHBlY3Qoc2V0LmdldEFzc2V0SURzKGFkZHJlc3NlcykpLnRvRXF1YWwoc2V0LmdldEFzc2V0SURzKCkpXG4gICAgfSlcblxuICAgIGRlc2NyaWJlKFwiTWVyZ2UgUnVsZXNcIiwgKCk6IHZvaWQgPT4ge1xuICAgICAgbGV0IHNldEE6IFVUWE9TZXRcbiAgICAgIGxldCBzZXRCOiBVVFhPU2V0XG4gICAgICBsZXQgc2V0QzogVVRYT1NldFxuICAgICAgbGV0IHNldEQ6IFVUWE9TZXRcbiAgICAgIGxldCBzZXRFOiBVVFhPU2V0XG4gICAgICBsZXQgc2V0RjogVVRYT1NldFxuICAgICAgbGV0IHNldEc6IFVUWE9TZXRcbiAgICAgIGxldCBzZXRIOiBVVFhPU2V0XG4gICAgICAvLyBUYWtlLW9yLUxlYXZlXG4gICAgICBjb25zdCBuZXd1dHhvOiBzdHJpbmcgPSBiaW50b29scy5jYjU4RW5jb2RlKFxuICAgICAgICBCdWZmZXIuZnJvbShcbiAgICAgICAgICBcIjAwMDBhY2Y4ODY0N2IzZmJhYTlmZGY0Mzc4ZjNhMGRmNmE1ZDE1ZDhlZmIwMThhZDc4ZjEyNjkwMzkwZTc5ZTE2ODc2MDAwMDAwMDNhY2Y4ODY0N2IzZmJhYTlmZGY0Mzc4ZjNhMGRmNmE1ZDE1ZDhlZmIwMThhZDc4ZjEyNjkwMzkwZTc5ZTE2ODc2MDAwMDAwMDcwMDAwMDAwMDAwMDE4NmEwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAxMDAwMDAwMDFmY2VkYThmOTBmY2I1ZDMwNjE0Yjk5ZDc5ZmM0YmFhMjkzMDc3NjI2XCIsXG4gICAgICAgICAgXCJoZXhcIlxuICAgICAgICApXG4gICAgICApXG5cbiAgICAgIGJlZm9yZUVhY2goKCk6IHZvaWQgPT4ge1xuICAgICAgICBzZXRBID0gbmV3IFVUWE9TZXQoKVxuICAgICAgICBzZXRBLmFkZEFycmF5KFt1dHhvc3Ryc1swXSwgdXR4b3N0cnNbMl1dKVxuXG4gICAgICAgIHNldEIgPSBuZXcgVVRYT1NldCgpXG4gICAgICAgIHNldEIuYWRkQXJyYXkoW3V0eG9zdHJzWzFdLCB1dHhvc3Ryc1syXV0pXG5cbiAgICAgICAgc2V0QyA9IG5ldyBVVFhPU2V0KClcbiAgICAgICAgc2V0Qy5hZGRBcnJheShbdXR4b3N0cnNbMF0sIHV0eG9zdHJzWzFdXSlcblxuICAgICAgICBzZXREID0gbmV3IFVUWE9TZXQoKVxuICAgICAgICBzZXRELmFkZEFycmF5KFt1dHhvc3Ryc1sxXV0pXG5cbiAgICAgICAgc2V0RSA9IG5ldyBVVFhPU2V0KClcbiAgICAgICAgc2V0RS5hZGRBcnJheShbXSkgLy8gZW1wdHkgc2V0XG5cbiAgICAgICAgc2V0RiA9IG5ldyBVVFhPU2V0KClcbiAgICAgICAgc2V0Ri5hZGRBcnJheSh1dHhvc3RycykgLy8gZnVsbCBzZXQsIHNlcGFyYXRlIGZyb20gc2VsZlxuXG4gICAgICAgIHNldEcgPSBuZXcgVVRYT1NldCgpXG4gICAgICAgIHNldEcuYWRkQXJyYXkoW25ld3V0eG8sIC4uLnV0eG9zdHJzXSkgLy8gZnVsbCBzZXQgd2l0aCBuZXcgZWxlbWVudFxuXG4gICAgICAgIHNldEggPSBuZXcgVVRYT1NldCgpXG4gICAgICAgIHNldEguYWRkQXJyYXkoW25ld3V0eG9dKSAvLyBzZXQgd2l0aCBvbmx5IGEgbmV3IGVsZW1lbnRcbiAgICAgIH0pXG5cbiAgICAgIHRlc3QoXCJ1bmtub3duIG1lcmdlIHJ1bGVcIiwgKCk6IHZvaWQgPT4ge1xuICAgICAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgICAgIHNldC5tZXJnZUJ5UnVsZShzZXRBLCBcIkVSUk9SXCIpXG4gICAgICAgIH0pLnRvVGhyb3coKVxuICAgICAgICBjb25zdCBzZXRBcnJheTogVVRYT1tdID0gc2V0Ry5nZXRBbGxVVFhPcygpXG4gICAgICB9KVxuXG4gICAgICB0ZXN0KFwiaW50ZXJzZWN0aW9uXCIsICgpOiB2b2lkID0+IHtcbiAgICAgICAgbGV0IHJlc3VsdHM6IFVUWE9TZXRcbiAgICAgICAgbGV0IHRlc3Q6IGJvb2xlYW5cblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEEsIFwiaW50ZXJzZWN0aW9uXCIpXG4gICAgICAgIHRlc3QgPSBzZXRNZXJnZVRlc3RlcihcbiAgICAgICAgICByZXN1bHRzLFxuICAgICAgICAgIFtzZXRBXSxcbiAgICAgICAgICBbc2V0Qiwgc2V0Qywgc2V0RCwgc2V0RSwgc2V0Riwgc2V0Rywgc2V0SF1cbiAgICAgICAgKVxuICAgICAgICBleHBlY3QodGVzdCkudG9CZSh0cnVlKVxuXG4gICAgICAgIHJlc3VsdHMgPSBzZXQubWVyZ2VCeVJ1bGUoc2V0RiwgXCJpbnRlcnNlY3Rpb25cIilcbiAgICAgICAgdGVzdCA9IHNldE1lcmdlVGVzdGVyKFxuICAgICAgICAgIHJlc3VsdHMsXG4gICAgICAgICAgW3NldEZdLFxuICAgICAgICAgIFtzZXRBLCBzZXRCLCBzZXRDLCBzZXRELCBzZXRFLCBzZXRHLCBzZXRIXVxuICAgICAgICApXG4gICAgICAgIGV4cGVjdCh0ZXN0KS50b0JlKHRydWUpXG5cbiAgICAgICAgcmVzdWx0cyA9IHNldC5tZXJnZUJ5UnVsZShzZXRHLCBcImludGVyc2VjdGlvblwiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0Rl0sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEQsIHNldEUsIHNldEcsIHNldEhdXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEgsIFwiaW50ZXJzZWN0aW9uXCIpXG4gICAgICAgIHRlc3QgPSBzZXRNZXJnZVRlc3RlcihcbiAgICAgICAgICByZXN1bHRzLFxuICAgICAgICAgIFtzZXRFXSxcbiAgICAgICAgICBbc2V0QSwgc2V0Qiwgc2V0Qywgc2V0RCwgc2V0Riwgc2V0Rywgc2V0SF1cbiAgICAgICAgKVxuICAgICAgICBleHBlY3QodGVzdCkudG9CZSh0cnVlKVxuICAgICAgfSlcblxuICAgICAgdGVzdChcImRpZmZlcmVuY2VTZWxmXCIsICgpOiB2b2lkID0+IHtcbiAgICAgICAgbGV0IHJlc3VsdHM6IFVUWE9TZXRcbiAgICAgICAgbGV0IHRlc3Q6IGJvb2xlYW5cblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEEsIFwiZGlmZmVyZW5jZVNlbGZcIilcbiAgICAgICAgdGVzdCA9IHNldE1lcmdlVGVzdGVyKFxuICAgICAgICAgIHJlc3VsdHMsXG4gICAgICAgICAgW3NldERdLFxuICAgICAgICAgIFtzZXRBLCBzZXRCLCBzZXRDLCBzZXRFLCBzZXRGLCBzZXRHLCBzZXRIXVxuICAgICAgICApXG4gICAgICAgIGV4cGVjdCh0ZXN0KS50b0JlKHRydWUpXG5cbiAgICAgICAgcmVzdWx0cyA9IHNldC5tZXJnZUJ5UnVsZShzZXRGLCBcImRpZmZlcmVuY2VTZWxmXCIpXG4gICAgICAgIHRlc3QgPSBzZXRNZXJnZVRlc3RlcihcbiAgICAgICAgICByZXN1bHRzLFxuICAgICAgICAgIFtzZXRFXSxcbiAgICAgICAgICBbc2V0QSwgc2V0Qiwgc2V0Qywgc2V0RCwgc2V0Riwgc2V0Rywgc2V0SF1cbiAgICAgICAgKVxuICAgICAgICBleHBlY3QodGVzdCkudG9CZSh0cnVlKVxuXG4gICAgICAgIHJlc3VsdHMgPSBzZXQubWVyZ2VCeVJ1bGUoc2V0RywgXCJkaWZmZXJlbmNlU2VsZlwiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0RV0sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEQsIHNldEYsIHNldEcsIHNldEhdXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEgsIFwiZGlmZmVyZW5jZVNlbGZcIilcbiAgICAgICAgdGVzdCA9IHNldE1lcmdlVGVzdGVyKFxuICAgICAgICAgIHJlc3VsdHMsXG4gICAgICAgICAgW3NldEZdLFxuICAgICAgICAgIFtzZXRBLCBzZXRCLCBzZXRDLCBzZXRELCBzZXRFLCBzZXRHLCBzZXRIXVxuICAgICAgICApXG4gICAgICAgIGV4cGVjdCh0ZXN0KS50b0JlKHRydWUpXG4gICAgICB9KVxuXG4gICAgICB0ZXN0KFwiZGlmZmVyZW5jZU5ld1wiLCAoKTogdm9pZCA9PiB7XG4gICAgICAgIGxldCByZXN1bHRzOiBVVFhPU2V0XG4gICAgICAgIGxldCB0ZXN0OiBib29sZWFuXG5cbiAgICAgICAgcmVzdWx0cyA9IHNldC5tZXJnZUJ5UnVsZShzZXRBLCBcImRpZmZlcmVuY2VOZXdcIilcbiAgICAgICAgdGVzdCA9IHNldE1lcmdlVGVzdGVyKFxuICAgICAgICAgIHJlc3VsdHMsXG4gICAgICAgICAgW3NldEVdLFxuICAgICAgICAgIFtzZXRBLCBzZXRCLCBzZXRDLCBzZXRELCBzZXRGLCBzZXRHLCBzZXRIXVxuICAgICAgICApXG4gICAgICAgIGV4cGVjdCh0ZXN0KS50b0JlKHRydWUpXG5cbiAgICAgICAgcmVzdWx0cyA9IHNldC5tZXJnZUJ5UnVsZShzZXRGLCBcImRpZmZlcmVuY2VOZXdcIilcbiAgICAgICAgdGVzdCA9IHNldE1lcmdlVGVzdGVyKFxuICAgICAgICAgIHJlc3VsdHMsXG4gICAgICAgICAgW3NldEVdLFxuICAgICAgICAgIFtzZXRBLCBzZXRCLCBzZXRDLCBzZXRELCBzZXRGLCBzZXRHLCBzZXRIXVxuICAgICAgICApXG4gICAgICAgIGV4cGVjdCh0ZXN0KS50b0JlKHRydWUpXG5cbiAgICAgICAgcmVzdWx0cyA9IHNldC5tZXJnZUJ5UnVsZShzZXRHLCBcImRpZmZlcmVuY2VOZXdcIilcbiAgICAgICAgdGVzdCA9IHNldE1lcmdlVGVzdGVyKFxuICAgICAgICAgIHJlc3VsdHMsXG4gICAgICAgICAgW3NldEhdLFxuICAgICAgICAgIFtzZXRBLCBzZXRCLCBzZXRDLCBzZXRELCBzZXRFLCBzZXRGLCBzZXRHXVxuICAgICAgICApXG4gICAgICAgIGV4cGVjdCh0ZXN0KS50b0JlKHRydWUpXG5cbiAgICAgICAgcmVzdWx0cyA9IHNldC5tZXJnZUJ5UnVsZShzZXRILCBcImRpZmZlcmVuY2VOZXdcIilcbiAgICAgICAgdGVzdCA9IHNldE1lcmdlVGVzdGVyKFxuICAgICAgICAgIHJlc3VsdHMsXG4gICAgICAgICAgW3NldEhdLFxuICAgICAgICAgIFtzZXRBLCBzZXRCLCBzZXRDLCBzZXRELCBzZXRFLCBzZXRGLCBzZXRHXVxuICAgICAgICApXG4gICAgICAgIGV4cGVjdCh0ZXN0KS50b0JlKHRydWUpXG4gICAgICB9KVxuXG4gICAgICB0ZXN0KFwic3ltRGlmZmVyZW5jZVwiLCAoKTogdm9pZCA9PiB7XG4gICAgICAgIGxldCByZXN1bHRzOiBVVFhPU2V0XG4gICAgICAgIGxldCB0ZXN0OiBib29sZWFuXG5cbiAgICAgICAgcmVzdWx0cyA9IHNldC5tZXJnZUJ5UnVsZShzZXRBLCBcInN5bURpZmZlcmVuY2VcIilcbiAgICAgICAgdGVzdCA9IHNldE1lcmdlVGVzdGVyKFxuICAgICAgICAgIHJlc3VsdHMsXG4gICAgICAgICAgW3NldERdLFxuICAgICAgICAgIFtzZXRBLCBzZXRCLCBzZXRDLCBzZXRFLCBzZXRGLCBzZXRHLCBzZXRIXVxuICAgICAgICApXG4gICAgICAgIGV4cGVjdCh0ZXN0KS50b0JlKHRydWUpXG5cbiAgICAgICAgcmVzdWx0cyA9IHNldC5tZXJnZUJ5UnVsZShzZXRGLCBcInN5bURpZmZlcmVuY2VcIilcbiAgICAgICAgdGVzdCA9IHNldE1lcmdlVGVzdGVyKFxuICAgICAgICAgIHJlc3VsdHMsXG4gICAgICAgICAgW3NldEVdLFxuICAgICAgICAgIFtzZXRBLCBzZXRCLCBzZXRDLCBzZXRELCBzZXRGLCBzZXRHLCBzZXRIXVxuICAgICAgICApXG4gICAgICAgIGV4cGVjdCh0ZXN0KS50b0JlKHRydWUpXG5cbiAgICAgICAgcmVzdWx0cyA9IHNldC5tZXJnZUJ5UnVsZShzZXRHLCBcInN5bURpZmZlcmVuY2VcIilcbiAgICAgICAgdGVzdCA9IHNldE1lcmdlVGVzdGVyKFxuICAgICAgICAgIHJlc3VsdHMsXG4gICAgICAgICAgW3NldEhdLFxuICAgICAgICAgIFtzZXRBLCBzZXRCLCBzZXRDLCBzZXRELCBzZXRFLCBzZXRGLCBzZXRHXVxuICAgICAgICApXG4gICAgICAgIGV4cGVjdCh0ZXN0KS50b0JlKHRydWUpXG5cbiAgICAgICAgcmVzdWx0cyA9IHNldC5tZXJnZUJ5UnVsZShzZXRILCBcInN5bURpZmZlcmVuY2VcIilcbiAgICAgICAgdGVzdCA9IHNldE1lcmdlVGVzdGVyKFxuICAgICAgICAgIHJlc3VsdHMsXG4gICAgICAgICAgW3NldEddLFxuICAgICAgICAgIFtzZXRBLCBzZXRCLCBzZXRDLCBzZXRELCBzZXRFLCBzZXRGLCBzZXRIXVxuICAgICAgICApXG4gICAgICAgIGV4cGVjdCh0ZXN0KS50b0JlKHRydWUpXG4gICAgICB9KVxuXG4gICAgICB0ZXN0KFwidW5pb25cIiwgKCk6IHZvaWQgPT4ge1xuICAgICAgICBsZXQgcmVzdWx0czogVVRYT1NldFxuICAgICAgICBsZXQgdGVzdDogYm9vbGVhblxuXG4gICAgICAgIHJlc3VsdHMgPSBzZXQubWVyZ2VCeVJ1bGUoc2V0QSwgXCJ1bmlvblwiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0Rl0sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEQsIHNldEUsIHNldEcsIHNldEhdXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEYsIFwidW5pb25cIilcbiAgICAgICAgdGVzdCA9IHNldE1lcmdlVGVzdGVyKFxuICAgICAgICAgIHJlc3VsdHMsXG4gICAgICAgICAgW3NldEZdLFxuICAgICAgICAgIFtzZXRBLCBzZXRCLCBzZXRDLCBzZXRELCBzZXRFLCBzZXRHLCBzZXRIXVxuICAgICAgICApXG4gICAgICAgIGV4cGVjdCh0ZXN0KS50b0JlKHRydWUpXG5cbiAgICAgICAgcmVzdWx0cyA9IHNldC5tZXJnZUJ5UnVsZShzZXRHLCBcInVuaW9uXCIpXG4gICAgICAgIHRlc3QgPSBzZXRNZXJnZVRlc3RlcihcbiAgICAgICAgICByZXN1bHRzLFxuICAgICAgICAgIFtzZXRHXSxcbiAgICAgICAgICBbc2V0QSwgc2V0Qiwgc2V0Qywgc2V0RCwgc2V0RSwgc2V0Riwgc2V0SF1cbiAgICAgICAgKVxuICAgICAgICBleHBlY3QodGVzdCkudG9CZSh0cnVlKVxuXG4gICAgICAgIHJlc3VsdHMgPSBzZXQubWVyZ2VCeVJ1bGUoc2V0SCwgXCJ1bmlvblwiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0R10sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEQsIHNldEUsIHNldEYsIHNldEhdXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcbiAgICAgIH0pXG5cbiAgICAgIHRlc3QoXCJ1bmlvbk1pbnVzTmV3XCIsICgpOiB2b2lkID0+IHtcbiAgICAgICAgbGV0IHJlc3VsdHM6IFVUWE9TZXRcbiAgICAgICAgbGV0IHRlc3Q6IGJvb2xlYW5cblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEEsIFwidW5pb25NaW51c05ld1wiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0RF0sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEUsIHNldEYsIHNldEcsIHNldEhdXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEYsIFwidW5pb25NaW51c05ld1wiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0RV0sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEQsIHNldEYsIHNldEcsIHNldEhdXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEcsIFwidW5pb25NaW51c05ld1wiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0RV0sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEQsIHNldEYsIHNldEcsIHNldEhdXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEgsIFwidW5pb25NaW51c05ld1wiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0Rl0sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEQsIHNldEUsIHNldEcsIHNldEhdXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcbiAgICAgIH0pXG5cbiAgICAgIHRlc3QoXCJ1bmlvbk1pbnVzU2VsZlwiLCAoKTogdm9pZCA9PiB7XG4gICAgICAgIGxldCByZXN1bHRzOiBVVFhPU2V0XG4gICAgICAgIGxldCB0ZXN0OiBib29sZWFuXG5cbiAgICAgICAgcmVzdWx0cyA9IHNldC5tZXJnZUJ5UnVsZShzZXRBLCBcInVuaW9uTWludXNTZWxmXCIpXG4gICAgICAgIHRlc3QgPSBzZXRNZXJnZVRlc3RlcihcbiAgICAgICAgICByZXN1bHRzLFxuICAgICAgICAgIFtzZXRFXSxcbiAgICAgICAgICBbc2V0QSwgc2V0Qiwgc2V0Qywgc2V0RCwgc2V0Riwgc2V0Rywgc2V0SF1cbiAgICAgICAgKVxuICAgICAgICBleHBlY3QodGVzdCkudG9CZSh0cnVlKVxuXG4gICAgICAgIHJlc3VsdHMgPSBzZXQubWVyZ2VCeVJ1bGUoc2V0RiwgXCJ1bmlvbk1pbnVzU2VsZlwiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0RV0sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEQsIHNldEYsIHNldEcsIHNldEhdXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEcsIFwidW5pb25NaW51c1NlbGZcIilcbiAgICAgICAgdGVzdCA9IHNldE1lcmdlVGVzdGVyKFxuICAgICAgICAgIHJlc3VsdHMsXG4gICAgICAgICAgW3NldEhdLFxuICAgICAgICAgIFtzZXRBLCBzZXRCLCBzZXRDLCBzZXRELCBzZXRFLCBzZXRGLCBzZXRHXVxuICAgICAgICApXG4gICAgICAgIGV4cGVjdCh0ZXN0KS50b0JlKHRydWUpXG5cbiAgICAgICAgcmVzdWx0cyA9IHNldC5tZXJnZUJ5UnVsZShzZXRILCBcInVuaW9uTWludXNTZWxmXCIpXG4gICAgICAgIHRlc3QgPSBzZXRNZXJnZVRlc3RlcihcbiAgICAgICAgICByZXN1bHRzLFxuICAgICAgICAgIFtzZXRIXSxcbiAgICAgICAgICBbc2V0QSwgc2V0Qiwgc2V0Qywgc2V0RCwgc2V0RSwgc2V0Riwgc2V0R11cbiAgICAgICAgKVxuICAgICAgICBleHBlY3QodGVzdCkudG9CZSh0cnVlKVxuICAgICAgfSlcbiAgICB9KVxuICB9KVxufSlcbiJdfQ==