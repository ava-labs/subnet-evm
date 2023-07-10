"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const bn_js_1 = __importDefault(require("bn.js"));
const buffer_1 = require("buffer/");
const bintools_1 = __importDefault(require("../../../src/utils/bintools"));
const utxos_1 = require("../../../src/apis/avm/utxos");
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
    test("bad creation", () => {
        const set = new utxos_1.UTXOSet();
        const bad = bintools.cb58Encode(buffer_1.Buffer.from("aasdfasd", "hex"));
        set.add(bad);
        const utxo = new utxos_1.UTXO();
        expect(() => {
            utxo.fromString(bad);
        }).toThrow();
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
                balance1.add(set.getBalance(addrs, assetID));
                balance2.add(utxos[i].getOutput().getAmount());
            }
            expect(balance1.toString()).toBe(balance2.toString());
            balance1 = new bn_js_1.default(0);
            balance2 = new bn_js_1.default(0);
            const now = (0, helperfunctions_1.UnixNow)();
            for (let i = 0; i < utxos.length; i++) {
                const assetID = bintools.cb58Encode(utxos[i].getAssetID());
                balance1.add(set.getBalance(addrs, assetID, now));
                balance2.add(utxos[i].getOutput().getAmount());
            }
            expect(balance1.toString()).toBe(balance2.toString());
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
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidXR4b3MudGVzdC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3Rlc3RzL2FwaXMvYXZtL3V0eG9zLnRlc3QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7QUFBQSxrREFBc0I7QUFDdEIsb0NBQWdDO0FBQ2hDLDJFQUFrRDtBQUNsRCx1REFBMkQ7QUFFM0Qsd0VBQTREO0FBRzVELE1BQU0sUUFBUSxHQUFhLGtCQUFRLENBQUMsV0FBVyxFQUFFLENBQUE7QUFDakQsTUFBTSxPQUFPLEdBQXVCLFNBQVMsQ0FBQTtBQUU3QyxRQUFRLENBQUMsTUFBTSxFQUFFLEdBQVMsRUFBRTtJQUMxQixNQUFNLE9BQU8sR0FDWCw4T0FBOE8sQ0FBQTtJQUNoUCxNQUFNLFNBQVMsR0FBVyxVQUFVLENBQUE7SUFDcEMsTUFBTSxPQUFPLEdBQ1gsa0VBQWtFLENBQUE7SUFDcEUsTUFBTSxNQUFNLEdBQ1Ysa0VBQWtFLENBQUE7SUFDcEUsTUFBTSxRQUFRLEdBQVcsZUFBTSxDQUFDLElBQUksQ0FBQyxPQUFPLEVBQUUsS0FBSyxDQUFDLENBQUE7SUFFcEQsVUFBVTtJQUNWLE1BQU0sU0FBUyxHQUFXLFFBQVEsQ0FBQyxVQUFVLENBQUMsUUFBUSxDQUFDLENBQUE7SUFDdkQseUtBQXlLO0lBRXpLLG9DQUFvQztJQUNwQyxJQUFJLENBQUMsVUFBVSxFQUFFLEdBQVMsRUFBRTtRQUMxQixNQUFNLEVBQUUsR0FBUyxJQUFJLFlBQUksRUFBRSxDQUFBO1FBQzNCLEVBQUUsQ0FBQyxVQUFVLENBQUMsUUFBUSxDQUFDLENBQUE7UUFDdkIsTUFBTSxLQUFLLEdBQVcsRUFBRSxDQUFDLFFBQVEsRUFBRSxDQUFDLFFBQVEsQ0FBQyxLQUFLLENBQUMsQ0FBQTtRQUNuRCxNQUFNLENBQUMsS0FBSyxDQUFDLENBQUMsSUFBSSxDQUFDLE9BQU8sQ0FBQyxDQUFBO0lBQzdCLENBQUMsQ0FBQyxDQUFBO0lBRUYsSUFBSSxDQUFDLGdCQUFnQixFQUFFLEdBQVMs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QU0sQ0FBQyxHQUFHLENBQUMsUUFBUSxDQUFDLEVBQUUsQ0FBQyxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO1lBQ25DLE1BQU0sSUFBSSxHQUFTLElBQUksWUFBSSxFQUFFLENBQUE7WUFDN0IsSUFBSSxDQUFDLFVBQVUsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTtZQUM1QixNQUFNLFFBQVEsR0FBUyxHQUFHLENBQUMsT0FBTyxDQUFDLElBQUksQ0FBQyxTQUFTLEVBQUUsQ0FBUyxDQUFBO1lBQzVELE1BQU0sQ0FBQyxRQUFRLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQyxJQUFJLENBQUMsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUE7U0FDOUM7UUFFRCxHQUFHLENBQUMsUUFBUSxDQUFDLEdBQUcsQ0FBQyxXQUFXLEVBQUUsQ0FBQyxDQUFBO1FBQy9CLEtBQUssSUFBSSxDQUFDLEdBQVcsQ0FBQyxFQUFFLENBQUMsR0FBRyxRQUFRLENBQUMsTUFBTSxFQUFFLENBQUMsRUFBRSxFQUFFO1lBQ2hELE1BQU0sSUFBSSxHQUFTLElBQUksWUFBSSxFQUFFLENBQUE7WUFDN0IsSUFBSSxDQUFDLFVBQVUsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTtZQUM1QixNQUFNLENBQUMsR0FBRyxDQUFDLFFBQVEsQ0FBQyxJQUFJLENBQUMsQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtZQUVyQyxNQUFNLFFBQVEsR0FBUyxHQUFHLENBQUMsT0FBTyxDQUFDLElBQUksQ0FBQyxTQUFTLEVBQUUsQ0FBUyxDQUFBO1lBQzVELE1BQU0sQ0FBQyxRQUFRLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQyxJQUFJLENBQUMsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUE7U0FDOUM7UUFFRCxJQUFJLENBQUMsR0FBVyxHQUFHLENBQUMsU0FBUyxDQUFDLEtBQUssQ0FBQyxDQUFBO1FBQ3BDLElBQUksQ0FBQyxHQUFZLElBQUksZUFBTyxFQUFFLENBQUE7UUFDOUIsQ0FBQyxDQUFDLFdBQVcsQ0FBQyxDQUFDLENBQUMsQ0FBQTtRQUNoQixJQUFJLENBQUMsR0FBVyxHQUFHLENBQUMsU0FBUyxDQUFDLE9BQU8sQ0FBQyxDQUFBO1FBQ3RDLElBQUksQ0FBQyxHQUFZLElBQUksZUFBTyxFQUFFLENBQUE7UUFDOUIsQ0FBQyxDQUFDLFdBQVcsQ0FBQyxDQUFDLENBQUMsQ0FBQTtJQUNsQixDQUFDLENBQUMsQ0FBQTtJQUVGLElBQUksQ0FBQyxrQkFBa0IsRUFBRSxHQUFTLEVBQUU7UUFDbEMsTUFBTSxHQUFHLEdBQVksSUFBSSxlQUFPLEVBQUUsQ0FBQTtRQUNsQyxHQUFHLENBQUMsUUFBUSxDQUFDLFFBQVEsQ0FBQyxDQUFBO1FBQ3RCLE1BQU0sUUFBUSxHQUFTLElBQUksWUFBSSxFQUFFLENBQUE7UUFDakMsUUFBUSxDQUFDLFVBQVUsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTtRQUNoQyxNQUFNLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLEVBQUUsSUFBSSxDQUFDLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQyxJQUFJLENBQUMsUUFBUSxDQUFDLFFBQVEsRUFBRSxDQUFDLENBQUE7UUFDdkUsTUFBTSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQUMsUUFBUSxDQUFDLENBQUMsQ0FBQyxFQUFFLEtBQUssQ0FBQyxDQUFDLENBQUMsYUFBYSxFQUFFLENBQUE7UUFDbkQsTUFBTSxDQUFDLEdBQUcsQ0FBQyxRQUFRLENBQUMsUUFBUSxFQUFFLElBQUksQ0FBQyxDQUFDLE1BQU0sQ0FBQyxDQUFDLElBQUksQ0FBQyxDQUFDLENBQUMsQ0FBQTtRQUNuRCxNQUFNLENBQUMsR0FBRyxDQUFDLFFBQVEsQ0FBQyxRQUFRLEVBQUUsS0FBSyxDQUFDLENBQUMsTUFBTSxDQUFDLENBQUMsSUFBSSxDQUFDLENBQUMsQ0FBQyxDQUFBO0lBQ3RELENBQUMsQ0FBQyxDQUFBO0lBRUYsUUFBUSxDQUFDLGVBQWUsRUFBRSxHQUFTLEVBQUU7UUFDbkMsSUFBSSxHQUFZLENBQUE7UUFDaEIsSUFBSSxLQUFhLENBQUE7UUFDakIsVUFBVSxDQUFDLEdBQVMsRUFBRTtZQUNwQixHQUFHLEdBQUcsSUFBSSxlQUFPLEVBQUUsQ0FBQTtZQUNuQixHQUFHLENBQUMsUUFBUSxDQUFDLFFBQVEsQ0FBQyxDQUFBO1lBQ3RCLEtBQUssR0FBRyxHQUFHLENBQUMsV0FBVyxFQUFFLENBQUE7UUFDM0IsQ0FBQyxDQUFDLENBQUE7UUFFRixJQUFJLENBQUMsUUFBUSxFQUFFLEdBQVMsRUFBRTtZQUN4QixNQUFNLFFBQVEsR0FBUyxJQUFJLFlBQUksRUFBRSxDQUFBO1lBQ2pDLFFBQVEsQ0FBQyxVQUFVLENBQUMsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUE7WUFDaEMsTUFBTSxDQUFDLEdBQUcsQ0FBQyxNQUFNLENBQUMsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQyxJQUFJLENBQUMsUUFBUSxDQUFDLFFBQVEsRUFBRSxDQUFDLENBQUE7WUFDcEUsTUFBTSxDQUFDLEdBQUcsQ0FBQyxNQUFNLENBQUMsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxhQUFhLEVBQUUsQ0FBQTtZQUMvQyxNQUFNLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLEVBQUUsS0FBSyxDQUFDLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQyxJQUFJLENBQUMsUUFBUSxDQUFDLFFBQVEsRUFBRSxDQUFDLENBQUE7WUFDeEUsTUFBTSxDQUFDLEdBQUcsQ0FBQyxNQUFNLENBQUMsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQyxJQUFJLENBQUMsUUFBUSxDQUFDLFFBQVEsRUFBRSxDQUFDLENBQUE7UUFDdEUsQ0FBQyxDQUFDLENBQUE7UUFFRixJQUFJLENBQUMsYUFBYSxFQUFFLEdBQVMsRUFBRTtZQUM3QixNQUFNLFFBQVEsR0FBUyxJQUFJLFlBQUksRUFBRSxDQUFBO1lBQ2pDLFFBQVEsQ0FBQyxVQUFVLENBQUMsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUE7WUFDaEMsTUFBTSxDQUFDLEdBQUcsQ0FBQyxXQUFXLENBQUMsUUFBUSxDQUFDLENBQUMsTUFBTSxDQUFDLENBQUMsSUFBSSxDQUFDLENBQUMsQ0FBQyxDQUFBO1lBQ2hELE1BQU0sQ0FBQyxHQUFHLENBQUMsV0FBVyxDQUFDLFFBQVEsQ0FBQyxDQUFDLE1BQU0sQ0FBQyxDQUFDLElBQUksQ0FBQyxDQUFDLENBQUMsQ0FBQTtZQUNoRCxNQUFNLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLEVBQUUsS0FBSyxDQUFDLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQyxJQUFJLENBQUMsUUFBUSxDQUFDLFFBQVEsRUFBRSxDQUFDLENBQUE7WUFDeEUsTUFBTSxDQUFDLEdBQUcsQ0FBQyxXQUFXLENBQUMsUUFBUSxDQUFDLENBQUMsTUFBTSxDQUFDLENBQUMsSUFBSSxDQUFDLENBQUMsQ0FBQyxDQUFBO1lBQ2hELE1BQU0sQ0FBQyxHQUFHLENBQUMsUUFBUSxDQUFDLFFBQVEsRUFBRSxLQUFLLENBQUMsQ0FBQyxNQUFNLENBQUMsQ0FBQyxJQUFJLENBQUMsQ0FBQyxDQUFDLENBQUE7WUFDcEQsTUFBTSxDQUFDLEdBQUcsQ0FBQyxXQUFXLENBQUMsS0FBSyxDQUFDLENBQUMsTUFBTSxDQUFDLENBQUMsSUFBSSxDQUFDLENBQUMsQ0FBQyxDQUFBO1FBQy9DLENBQUMsQ0FBQyxDQUFBO1FBRUYsSUFBSSxDQUFDLFlBQVksRUFBRSxHQUFTLEVBQUU7WUFDNUIsTUFBTSxJQUFJLEdBQWEsR0FBRyxDQUFDLFVBQVUsRUFBRSxDQUFBO1lBQ3ZDLEtBQUssSUFBSSxDQUFDLEdBQVcsQ0FBQyxFQUFFLENBQUMsR0FBRyxLQUFLLENBQUMsTUFBTSxFQUFFLENBQUMsRUFBRSxFQUFFO2dCQUM3QyxNQUFNLENBQUMsSUFBSSxDQUFDLE9BQU8sQ0FBQyxLQUFLLENBQUMsQ0FBQyxDQUFDLENBQUMsU0FBUyxFQUFFLENBQUMsQ0FBQyxDQUFDLEdBQUcsQ0FBQyxJQUFJLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTthQUN4RDtRQUNILENBQUMsQ0FBQyxDQUFBO1FBRUYsSUFBSSxDQUFDLGFBQWEsRUFBRSxHQUFTLEVBQUU7WUFDN0IsTUFBTSxRQUFRLEdBQVcsR0FBRyxDQUFDLFdBQVcsRUFBRSxDQUFBO1lBQzFDLE1BQU0sS0FBSyxHQUFhLEVBQUUsQ0FBQTtZQUMxQixLQUFLLElBQUksQ0FBQyxHQUFXLENBQUMsRUFBRSxDQUFDLEdBQUcsUUFBUSxDQUFDLE1BQU0sRUFBRSxDQUFDLEVBQUUsRUFBRTtnQkFDaEQsS0FBSyxDQUFDLElBQUksQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQTthQUNuQztZQUNELEtBQUssSUFBSSxDQUFDLEdBQVcsQ0FBQyxFQUFFLENBQUMsR0FBRyxRQUFRLENBQUMsTUFBTSxFQUFFLENBQUMsRUFBRSxFQUFFO2dCQUNoRCxNQUFNLENBQUMsS0FBSyxDQUFDLE9BQU8sQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDLEdBQUcsQ0FBQyxJQUFJLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTthQUNoRDtZQUNELE1BQU0sSUFBSSxHQUFhLEdBQUcsQ0FBQyxVQUFVLEVBQUUsQ0FBQTtZQUN2QyxNQUFNLFNBQVMsR0FBVyxHQUFHLENBQUMsV0FBVyxDQUFDLElBQUksQ0FBQyxDQUFBO1lBQy9DLE1BQU0sTUFBTSxHQUFhLEVBQUUsQ0FBQTtZQUMzQixLQUFLLElBQUksQ0FBQyxHQUFXLENBQUMsRUFBRSxDQUFDLEdBQUcsUUFBUSxDQUFDLE1BQU0sRUFBRSxDQUFDLEVBQUUsRUFBRTtnQkFDaEQsTUFBTSxDQUFDLElBQUksQ0FBQyxTQUFTLENBQUMsQ0FBQyxDQUFDLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQTthQUNyQztZQUNELEtBQUssSUFBSSxDQUFDLEdBQVcsQ0FBQyxFQUFFLENBQUMsR0FBRyxRQUFRLENBQUMsTUFBTSxFQUFFLENBQUMsRUFBRSxFQUFFO2dCQUNoRCxNQUFNLENBQUMsTUFBTSxDQUFDLE9BQU8sQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDLEdBQUcsQ0FBQyxJQUFJLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTthQUNqRDtRQUNILENBQUMsQ0FBQyxDQUFBO1FBRUYsSUFBSSxDQUFDLHVCQUF1QixFQUFFLEdBQVMsRUFBRTtZQUN2QyxJQUFJLE9BQWlCLENBQUE7WUFDckIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxVQUFVLENBQUMsQ0FBQyxLQUFLLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFBO1lBQ3BDLE1BQU0sQ0FBQyxPQUFPLENBQUMsTUFBTSxDQUFDLENBQUMsSUFBSSxDQUFDLENBQUMsQ0FBQyxDQUFBO1lBQzlCLE9BQU8sR0FBRyxHQUFHLENBQUMsVUFBVSxDQUFDLEtBQUssQ0FBQyxDQUFBO1lBQy9CLE1BQU0sQ0FBQyxPQUFPLENBQUMsTUFBTSxDQUFDLENBQUMsSUFBSSxDQUFDLENBQUMsQ0FBQyxDQUFBO1lBQzlCLE9BQU8sR0FBRyxHQUFHLENBQUMsVUFBVSxDQUFDLEtBQUssRUFBRSxLQUFLLENBQUMsQ0FBQTtZQUN0QyxNQUFNLENBQUMsT0FBTyxDQUFDLE1BQU0sQ0FBQyxDQUFDLElBQUksQ0FBQyxDQUFDLENBQUMsQ0FBQTtRQUNoQyxDQUFDLENBQUMsQ0FBQTtRQUVGLElBQUksQ0FBQyxtQkFBbUIsRUFBRSxHQUFTLEVBQUU7WUFDbkMsTUFBTSxLQUFLLEdBQWEsR0FBRyxDQUFDLGlCQUFpQixFQUFFLENBQUE7WUFDL0MsS0FBSyxJQUFJLENBQUMsR0FBVyxDQUFDLEVBQUUsQ0FBQyxHQUFHLFFBQVEsQ0FBQyxNQUFNLEVBQUUsQ0FBQyxFQUFFLEVBQUU7Z0JBQ2hELE1BQU0sQ0FBQyxLQUFLLENBQUMsT0FBTyxDQUFDLFFBQVEsQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsR0FBRyxDQUFDLElBQUksQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFBO2FBQ2hEO1lBQ0QsTUFBTSxJQUFJLEdBQWEsR0FBRyxDQUFDLFVBQVUsRUFBRSxDQUFBO1lBQ3ZDLE1BQU0sTUFBTSxHQUFhLEdBQUcsQ0FBQyxpQkFBaUIsQ0FBQyxJQUFJLENBQUMsQ0FBQTtZQUNwRCxLQUFLLElBQUksQ0FBQyxHQUFXLENBQUMsRUFBRSxDQUFDLEdBQUcsUUFBUSxDQUFDLE1BQU0sRUFBRSxDQUFDLEVBQUUsRUFBRTtnQkFDaEQsTUFBTSxDQUFDLE1BQU0sQ0FBQyxPQUFPLENBQUMsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxHQUFHLENBQUMsSUFBSSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUE7YUFDakQ7UUFDSCxDQUFDLENBQUMsQ0FBQTtRQUVGLElBQUksQ0FBQyxjQUFjLEVBQUUsR0FBUyxFQUFFO1lBQzlCLE1BQU0sQ0FBQyxHQUFHLENBQUMsWUFBWSxFQUFFLENBQUMsSUFBSSxFQUFFLENBQUMsQ0FBQyxhQUFhLENBQUMsS0FBSyxDQUFDLElBQUksRUFBRSxDQUFDLENBQUE7UUFDL0QsQ0FBQyxDQUFDLENBQUE7UUFFRixJQUFJLENBQUMsWUFBWSxFQUFFLEdBQVMsRUFBRTtZQUM1QixJQUFJLFFBQVksQ0FBQTtZQUNoQixJQUFJLFFBQVksQ0FBQTtZQUNoQixRQUFRLEdBQUcsSUFBSSxlQUFFLENBQUMsQ0FBQyxDQUFDLENBQUE7WUFDcEIsUUFBUSxHQUFHLElBQUksZUFBRSxDQUFDLENBQUMsQ0FBQyxDQUFBO1lBQ3BCLEtBQUssSUFBSSxDQUFDLEdBQVcsQ0FBQyxFQUFFLENBQUMsR0FBRyxLQUFLLENBQUMsTUFBTSxFQUFFLENBQUMsRUFBRSxFQUFFO2dCQUM3QyxNQUFNLE9BQU8sR0FBRyxLQUFLLENBQUMsQ0FBQyxDQUFDLENBQUMsVUFBVSxFQUFFLENBQUE7Z0JBQ3JDLFFBQVEsQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFDLFVBQVUsQ0FBQyxLQUFLLEVBQUUsT0FBTyxDQUFDLENBQUMsQ0FBQTtnQkFDNUMsUUFBUSxDQUFDLEdBQUcsQ0FBRSxLQUFLLENBQUMsQ0FBQyxDQUFDLENBQUMsU0FBUyxFQUFtQixDQUFDLFNBQVMsRUFBRSxDQUFDLENBQUE7YUFDakU7WUFDRCxNQUFNLENBQUMsUUFBUSxDQUFDLFFBQVEsRUFBRSxDQUFDLENBQUMsSUFBSSxDQUFDLFFBQVEsQ0FBQyxRQUFRLEVBQUUsQ0FBQyxDQUFBO1lBRXJELFFBQVEsR0FBRyxJQUFJLGVBQUUsQ0FBQyxDQUFDLENBQUMsQ0FBQTtZQUNwQixRQUFRLEdBQUcsSUFBSSxlQUFFLENBQUMsQ0FBQyxDQUFDLENBQUE7WUFDcEIsTUFBTSxHQUFHLEdBQU8sSUFBQSx5QkFBTyxHQUFFLENBQUE7WUFDekIsS0FBSyxJQUFJLENBQUMsR0FBVyxDQUFDLEVBQUUsQ0FBQyxHQUFHLEtBQUssQ0FBQyxNQUFNLEVBQUUsQ0FBQyxFQUFFLEVBQUU7Z0JBQzdDLE1BQU0sT0FBTyxHQUFHLFFBQVEsQ0FBQyxVQUFVLENBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQyxDQUFDLFVBQVUsRUFBRSxDQUFDLENBQUE7Z0JBQzFELFFBQVEsQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFDLFVBQVUsQ0FBQyxLQUFLLEVBQUUsT0FBTyxFQUFFLEdBQUcsQ0FBQyxDQUFDLENBQUE7Z0JBQ2pELFFBQVEsQ0FBQyxHQUFHLENBQUUsS0FBSyxDQUFDLENBQUMsQ0FBQyxDQUFDLFNBQVMsRUFBbUIsQ0FBQyxTQUFTLEVBQUUsQ0FBQyxDQUFBO2FBQ2pFO1lBQ0QsTUFBTSxDQUFDLFFBQVEsQ0FBQyxRQUFRLEVBQUUsQ0FBQyxDQUFDLElBQUksQ0FBQyxRQUFRLENBQUMsUUFBUSxFQUFFLENBQUMsQ0FBQTtRQUN2RCxDQUFDLENBQUMsQ0FBQTtRQUVGLElBQUksQ0FBQyxhQUFhLEVBQUUsR0FBUyxFQUFFO1lBQzdCLE1BQU0sUUFBUSxHQUFhLEdBQUcsQ0FBQyxXQUFXLEVBQUUsQ0FBQTtZQUM1QyxLQUFLLElBQUksQ0FBQyxHQUFXLENBQUMsRUFBRSxDQUFDLEdBQUcsS0FBSyxDQUFDLE1BQU0sRUFBRSxDQUFDLEVBQUUsRUFBRTtnQkFDN0MsTUFBTSxDQUFDLFFBQVEsQ0FBQyxDQUFDLFNBQVMsQ0FBQyxLQUFLLENBQUMsQ0FBQyxDQUFDLENBQUMsVUFBVSxFQUFFLENBQUMsQ0FBQTthQUNsRDtZQUNELE1BQU0sU0FBUyxHQUFhLEdBQUcsQ0FBQyxZQUFZLEVBQUUsQ0FBQTtZQUM5QyxNQUFNLENBQUMsR0FBRyxDQUFDLFdBQVcsQ0FBQyxTQUFTLENBQUMsQ0FBQyxDQUFDLE9BQU8sQ0FBQyxHQUFHLENBQUMsV0FBVyxFQUFFLENBQUMsQ0FBQTtRQUMvRCxDQUFDLENBQUMsQ0FBQTtRQUVGLFFBQVEsQ0FBQyxhQUFhLEVBQUUsR0FBUyxFQUFFO1lBQ2pDLElBQUksSUFBYSxDQUFBO1lBQ2pCLElBQUksSUFBYSxDQUFBO1lBQ2pCLElBQUksSUFBYSxDQUFBO1lBQ2pCLElBQUksSUFBYSxDQUFBO1lBQ2pCLElBQUksSUFBYSxDQUFBO1lBQ2pCLElBQUksSUFBYSxDQUFBO1lBQ2pCLElBQUksSUFBYSxDQUFBO1lBQ2pCLElBQUksSUFBYSxDQUFBO1lBQ2pCLGdCQUFnQjtZQUNoQixNQUFNLE9BQU8sR0FBVyxRQUFRLENBQUMsVUFBVSxDQUN6QyxlQUFNLENBQUMsSUFBSSxDQUNULDhPQUE4TyxFQUM5TyxLQUFLLENBQ04sQ0FDRixDQUFBO1lBRUQsVUFBVSxDQUFDLEdBQVMsRUFBRTtnQkFDcEIsSUFBSSxHQUFHLElBQUksZUFBTyxFQUFFLENBQUE7Z0JBQ3BCLElBQUksQ0FBQyxRQUFRLENBQUMsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLEVBQUUsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTtnQkFFekMsSUFBSSxHQUFHLElBQUksZUFBTyxFQUFFLENBQUE7Z0JBQ3BCLElBQUksQ0FBQyxRQUFRLENBQUMsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLEVBQUUsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTtnQkFFekMsSUFBSSxHQUFHLElBQUksZUFBTyxFQUFFLENBQUE7Z0JBQ3BCLElBQUksQ0FBQyxRQUFRLENBQUMsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLEVBQUUsUUFBUSxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTtnQkFFekMsSUFBSSxHQUFHLElBQUksZUFBTyxFQUFFLENBQUE7Z0JBQ3BCLElBQUksQ0FBQyxRQUFRLENBQUMsQ0FBQyxRQUFRLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFBO2dCQUU1QixJQUFJLEdBQUcsSUFBSSxlQUFPLEVBQUUsQ0FBQTtnQkFDcEIsSUFBSSxDQUFDLFFBQVEsQ0FBQyxFQUFFLENBQUMsQ0FBQSxDQUFDLFlBQVk7Z0JBRTlCLElBQUksR0FBRyxJQUFJLGVBQU8sRUFBRSxDQUFBO2dCQUNwQixJQUFJLENBQUMsUUFBUSxDQUFDLFFBQVEsQ0FBQyxDQUFBLENBQUMsK0JBQStCO2dCQUV2RCxJQUFJLEdBQUcsSUFBSSxlQUFPLEVBQUUsQ0FBQTtnQkFDcEIsSUFBSSxDQUFDLFFBQVEsQ0FBQyxDQUFDLE9BQU8sRUFBRSxHQUFHLFFBQVEsQ0FBQyxDQUFDLENBQUEsQ0FBQyw0QkFBNEI7Z0JBRWxFLElBQUksR0FBRyxJQUFJLGVBQU8sRUFBRSxDQUFBO2dCQUNwQixJQUFJLENBQUMsUUFBUSxDQUFDLENBQUMsT0FBTyxDQUFDLENBQUMsQ0FBQSxDQUFDLDhCQUE4QjtZQUN6RCxDQUFDLENBQUMsQ0FBQTtZQUVGLElBQUksQ0FBQyxvQkFBb0IsRUFBRSxHQUFTLEVBQUU7Z0JBQ3BDLE1BQU0sQ0FBQyxHQUFTLEVBQUU7b0JBQ2hCLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLE9BQU8sQ0FBQyxDQUFBO2dCQUNoQyxDQUFDLENBQUMsQ0FBQyxPQUFPLEVBQUUsQ0FBQTtZQUNkLENBQUMsQ0FBQyxDQUFBO1lBRUYsSUFBSSxDQUFDLGNBQWMsRUFBRSxHQUFTLEVBQUU7Z0JBQzlCLElBQUksT0FBZ0IsQ0FBQTtnQkFDcEIsSUFBSSxJQUFhLENBQUE7Z0JBRWpCLE9BQU8sR0FBRyxHQUFHLENBQUMsV0FBVyxDQUFDLElBQUksRUFBRSxjQUFjLENBQUMsQ0FBQTtnQkFDL0MsSUFBSSxHQUFHLGNBQWMsQ0FDbkIsT0FBTyxFQUNQLENBQUMsSUFBSSxDQUFDLEVBQ04sQ0FBQyxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLENBQUMsQ0FDM0MsQ0FBQTtnQkFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO2dCQUV2QixPQUFPLEdBQUcsR0FBRyxDQUFDLFdBQVcsQ0FBQyxJQUFJLEVBQUUsY0FBYyxDQUFDLENBQUE7Z0JBQy9DLElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtnQkFFdkIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLGNBQWMsQ0FBQyxDQUFBO2dCQUMvQyxJQUFJLEdBQUcsY0FBYyxDQUNuQixPQUFPLEVBQ1AsQ0FBQyxJQUFJLENBQUMsRUFDTixDQUFDLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksQ0FBQyxDQUMzQyxDQUFBO2dCQUNELE1BQU0sQ0FBQyxJQUFJLENBQUMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7Z0JBRXZCLE9BQU8sR0FBRyxHQUFHLENBQUMsV0FBVyxDQUFDLElBQUksRUFBRSxjQUFjLENBQUMsQ0FBQTtnQkFDL0MsSUFBSSxHQUFHLGNBQWMsQ0FDbkIsT0FBTyxFQUNQLENBQUMsSUFBSSxDQUFDLEVBQ04sQ0FBQyxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLENBQUMsQ0FDM0MsQ0FBQTtnQkFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO1lBQ3pCLENBQUMsQ0FBQyxDQUFBO1lBRUYsSUFBSSxDQUFDLGdCQUFnQixFQUFFLEdBQVMsRUFBRTtnQkFDaEMsSUFBSSxPQUFnQixDQUFBO2dCQUNwQixJQUFJLElBQWEsQ0FBQTtnQkFFakIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLGdCQUFnQixDQUFDLENBQUE7Z0JBQ2pELElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtnQkFFdkIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLGdCQUFnQixDQUFDLENBQUE7Z0JBQ2pELElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtnQkFFdkIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLGdCQUFnQixDQUFDLENBQUE7Z0JBQ2pELElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtnQkFFdkIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLGdCQUFnQixDQUFDLENBQUE7Z0JBQ2pELElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtZQUN6QixDQUFDLENBQUMsQ0FBQTtZQUVGLElBQUksQ0FBQyxlQUFlLEVBQUUsR0FBUyxFQUFFO2dCQUMvQixJQUFJLE9BQWdCLENBQUE7Z0JBQ3BCLElBQUksSUFBYSxDQUFBO2dCQUVqQixPQUFPLEdBQUcsR0FBRyxDQUFDLFdBQVcsQ0FBQyxJQUFJLEVBQUUsZUFBZSxDQUFDLENBQUE7Z0JBQ2hELElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtnQkFFdkIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLGVBQWUsQ0FBQyxDQUFBO2dCQUNoRCxJQUFJLEdBQUcsY0FBYyxDQUNuQixPQUFPLEVBQ1AsQ0FBQyxJQUFJLENBQUMsRUFDTixDQUFDLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksQ0FBQyxDQUMzQyxDQUFBO2dCQUNELE1BQU0sQ0FBQyxJQUFJLENBQUMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7Z0JBRXZCLE9BQU8sR0FBRyxHQUFHLENBQUMsV0FBVyxDQUFDLElBQUksRUFBRSxlQUFlLENBQUMsQ0FBQTtnQkFDaEQsSUFBSSxHQUFHLGNBQWMsQ0FDbkIsT0FBTyxFQUNQLENBQUMsSUFBSSxDQUFDLEVBQ04sQ0FBQyxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLENBQUMsQ0FDM0MsQ0FBQTtnQkFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO2dCQUV2QixPQUFPLEdBQUcsR0FBRyxDQUFDLFdBQVcsQ0FBQyxJQUFJLEVBQUUsZUFBZSxDQUFDLENBQUE7Z0JBQ2hELElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtZQUN6QixDQUFDLENBQUMsQ0FBQTtZQUVGLElBQUksQ0FBQyxlQUFlLEVBQUUsR0FBUyxFQUFFO2dCQUMvQixJQUFJLE9BQWdCLENBQUE7Z0JBQ3BCLElBQUksSUFBYSxDQUFBO2dCQUVqQixPQUFPLEdBQUcsR0FBRyxDQUFDLFdBQVcsQ0FBQyxJQUFJLEVBQUUsZUFBZSxDQUFDLENBQUE7Z0JBQ2hELElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtnQkFFdkIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLGVBQWUsQ0FBQyxDQUFBO2dCQUNoRCxJQUFJLEdBQUcsY0FBYyxDQUNuQixPQUFPLEVBQ1AsQ0FBQyxJQUFJLENBQUMsRUFDTixDQUFDLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksQ0FBQyxDQUMzQyxDQUFBO2dCQUNELE1BQU0sQ0FBQyxJQUFJLENBQUMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7Z0JBRXZCLE9BQU8sR0FBRyxHQUFHLENBQUMsV0FBVyxDQUFDLElBQUksRUFBRSxlQUFlLENBQUMsQ0FBQTtnQkFDaEQsSUFBSSxHQUFHLGNBQWMsQ0FDbkIsT0FBTyxFQUNQLENBQUMsSUFBSSxDQUFDLEVBQ04sQ0FBQyxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLENBQUMsQ0FDM0MsQ0FBQTtnQkFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO2dCQUV2QixPQUFPLEdBQUcsR0FBRyxDQUFDLFdBQVcsQ0FBQyxJQUFJLEVBQUUsZUFBZSxDQUFDLENBQUE7Z0JBQ2hELElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtZQUN6QixDQUFDLENBQUMsQ0FBQTtZQUVGLElBQUksQ0FBQyxPQUFPLEVBQUUsR0FBUyxFQUFFO2dCQUN2QixJQUFJLE9BQWdCLENBQUE7Z0JBQ3BCLElBQUksSUFBYSxDQUFBO2dCQUVqQixPQUFPLEdBQUcsR0FBRyxDQUFDLFdBQVcsQ0FBQyxJQUFJLEVBQUUsT0FBTyxDQUFDLENBQUE7Z0JBQ3hDLElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtnQkFFdkIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLE9BQU8sQ0FBQyxDQUFBO2dCQUN4QyxJQUFJLEdBQUcsY0FBYyxDQUNuQixPQUFPLEVBQ1AsQ0FBQyxJQUFJLENBQUMsRUFDTixDQUFDLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksQ0FBQyxDQUMzQyxDQUFBO2dCQUNELE1BQU0sQ0FBQyxJQUFJLENBQUMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7Z0JBRXZCLE9BQU8sR0FBRyxHQUFHLENBQUMsV0FBVyxDQUFDLElBQUksRUFBRSxPQUFPLENBQUMsQ0FBQTtnQkFDeEMsSUFBSSxHQUFHLGNBQWMsQ0FDbkIsT0FBTyxFQUNQLENBQUMsSUFBSSxDQUFDLEVBQ04sQ0FBQyxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLENBQUMsQ0FDM0MsQ0FBQTtnQkFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO2dCQUV2QixPQUFPLEdBQUcsR0FBRyxDQUFDLFdBQVcsQ0FBQyxJQUFJLEVBQUUsT0FBTyxDQUFDLENBQUE7Z0JBQ3hDLElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtZQUN6QixDQUFDLENBQUMsQ0FBQTtZQUVGLElBQUksQ0FBQyxlQUFlLEVBQUUsR0FBUyxFQUFFO2dCQUMvQixJQUFJLE9BQWdCLENBQUE7Z0JBQ3BCLElBQUksSUFBYSxDQUFBO2dCQUVqQixPQUFPLEdBQUcsR0FBRyxDQUFDLFdBQVcsQ0FBQyxJQUFJLEVBQUUsZUFBZSxDQUFDLENBQUE7Z0JBQ2hELElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtnQkFFdkIsT0FBTyxHQUFHLEdBQUcsQ0FBQyxXQUFXLENBQUMsSUFBSSxFQUFFLGVBQWUsQ0FBQyxDQUFBO2dCQUNoRCxJQUFJLEdBQUcsY0FBYyxDQUNuQixPQUFPLEVBQ1AsQ0FBQyxJQUFJLENBQUMsRUFDTixDQUFDLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksQ0FBQyxDQUMzQyxDQUFBO2dCQUNELE1BQU0sQ0FBQyxJQUFJLENBQUMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7Z0JBRXZCLE9BQU8sR0FBRyxHQUFHLENBQUMsV0FBVyxDQUFDLElBQUksRUFBRSxlQUFlLENBQUMsQ0FBQTtnQkFDaEQsSUFBSSxHQUFHLGNBQWMsQ0FDbkIsT0FBTyxFQUNQLENBQUMsSUFBSSxDQUFDLEVBQ04sQ0FBQyxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLENBQUMsQ0FDM0MsQ0FBQTtnQkFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFBO2dCQUV2QixPQUFPLEdBQUcsR0FBRyxDQUFDLFdBQVcsQ0FBQyxJQUFJLEVBQUUsZUFBZSxDQUFDLENBQUE7Z0JBQ2hELElBQUksR0FBRyxjQUFjLENBQ25CLE9BQU8sRUFDUCxDQUFDLElBQUksQ0FBQyxFQUNOLENBQUMsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxDQUFDLENBQzNDLENBQUE7Z0JBQ0QsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtZQUN6QixDQUFDLENBQUMsQ0FBQTtZQUVGLElBQUksQ0FBQyxnQkFBZ0IsRUFBRSxHQUFTLEVBQUU7Z0JBQ2hDLElBQUksT0FBZ0IsQ0FBQTtnQkFDcEIsSUFBSSxJQUFhLENBQUE7Z0JBRWpCLE9BQU8sR0FBRyxHQUFHLENBQUMsV0FBVyxDQUFDLElBQUksRUFBRSxnQkFBZ0IsQ0FBQyxDQUFBO2dCQUNqRCxJQUFJLEdBQUcsY0FBYyxDQUNuQixPQUFPLEVBQ1AsQ0FBQyxJQUFJLENBQUMsRUFDTixDQUFDLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksQ0FBQyxDQUMzQyxDQUFBO2dCQUNELE1BQU0sQ0FBQyxJQUFJLENBQUMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7Z0JBRXZCLE9BQU8sR0FBRyxHQUFHLENBQUMsV0FBVyxDQUFDLElBQUksRUFBRSxnQkFBZ0IsQ0FBQyxDQUFBO2dCQUNqRCxJQUFJLEdBQUcsY0FBYyxDQUNuQixPQUFPLEVBQ1AsQ0FBQyxJQUFJLENBQUMsRUFDTixDQUFDLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksQ0FBQyxDQUMzQyxDQUFBO2dCQUNELE1BQU0sQ0FBQyxJQUFJLENBQUMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7Z0JBRXZCLE9BQU8sR0FBRyxHQUFHLENBQUMsV0FBVyxDQUFDLElBQUksRUFBRSxnQkFBZ0IsQ0FBQyxDQUFBO2dCQUNqRCxJQUFJLEdBQUcsY0FBYyxDQUNuQixPQUFPLEVBQ1AsQ0FBQyxJQUFJLENBQUMsRUFDTixDQUFDLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksQ0FBQyxDQUMzQyxDQUFBO2dCQUNELE1BQU0sQ0FBQyxJQUFJLENBQUMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7Z0JBRXZCLE9BQU8sR0FBRyxHQUFHLENBQUMsV0FBVyxDQUFDLElBQUksRUFBRSxnQkFBZ0IsQ0FBQyxDQUFBO2dCQUNqRCxJQUFJLEdBQUcsY0FBYyxDQUNuQixPQUFPLEVBQ1AsQ0FBQyxJQUFJLENBQUMsRUFDTixDQUFDLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksQ0FBQyxDQUMzQyxDQUFBO2dCQUNELE1BQU0sQ0FBQyxJQUFJLENBQUMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUE7WUFDekIsQ0FBQyxDQUFDLENBQUE7UUFDSixDQUFDLENBQUMsQ0FBQTtJQUNKLENBQUMsQ0FBQyxDQUFBO0FBQ0osQ0FBQyxDQUFDLENBQUEiLCJzb3VyY2VzQ29udGVudCI6WyJpbXBvcnQgQk4gZnJvbSBcImJuLmpzXCJcbmltcG9ydCB7IEJ1ZmZlciB9IGZyb20gXCJidWZmZXIvXCJcbmltcG9ydCBCaW5Ub29scyBmcm9tIFwiLi4vLi4vLi4vc3JjL3V0aWxzL2JpbnRvb2xzXCJcbmltcG9ydCB7IFVUWE8sIFVUWE9TZXQgfSBmcm9tIFwiLi4vLi4vLi4vc3JjL2FwaXMvYXZtL3V0eG9zXCJcbmltcG9ydCB7IEFtb3VudE91dHB1dCB9IGZyb20gXCIuLi8uLi8uLi9zcmMvYXBpcy9hdm0vb3V0cHV0c1wiXG5pbXBvcnQgeyBVbml4Tm93IH0gZnJvbSBcIi4uLy4uLy4uL3NyYy91dGlscy9oZWxwZXJmdW5jdGlvbnNcIlxuaW1wb3J0IHsgU2VyaWFsaXplZEVuY29kaW5nIH0gZnJvbSBcIi4uLy4uLy4uL3NyYy91dGlsc1wiXG5cbmNvbnN0IGJpbnRvb2xzOiBCaW5Ub29scyA9IEJpblRvb2xzLmdldEluc3RhbmNlKClcbmNvbnN0IGRpc3BsYXk6IFNlcmlhbGl6ZWRFbmNvZGluZyA9IFwiZGlzcGxheVwiXG5cbmRlc2NyaWJlKFwiVVRYT1wiLCAoKTogdm9pZCA9PiB7XG4gIGNvbnN0IHV0eG9oZXg6IHN0cmluZyA9XG4gICAgXCIwMDAwMzhkMWI5ZjExMzg2NzJkYTZmYjZjMzUxMjU1MzkyNzZhOWFjYzJhNjY4ZDYzYmVhNmJhM2M3OTVlMmVkYjBmNTAwMDAwMDAxM2UwN2UzOGUyZjIzMTIxYmU4NzU2NDEyYzE4ZGI3MjQ2YTE2ZDI2ZWU5OTM2ZjNjYmEyOGJlMTQ5Y2ZkMzU1ODAwMDAwMDA3MDAwMDAwMDAwMDAwNGRkNTAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMTAwMDAwMDAxYTM2ZmQwYzJkYmNhYjMxMTczMWRkZTdlZjE1MTRiZDI2ZmNkYzc0ZFwiXG4gIGNvbnN0IG91dHB1dGlkeDogc3RyaW5nID0gXCIwMDAwMDAwMVwiXG4gIGNvbnN0IG91dHR4aWQ6IHN0cmluZyA9XG4gICAgXCIzOGQxYjlmMTEzODY3MmRhNmZiNmMzNTEyNTUzOTI3NmE5YWNjMmE2NjhkNjNiZWE2YmEzYzc5NWUyZWRiMGY1XCJcbiAgY29uc3Qgb3V0YWlkOiBzdHJpbmcgPVxuICAgIFwiM2UwN2UzOGUyZjIzMTIxYmU4NzU2NDEyYzE4ZGI3MjQ2YTE2ZDI2ZWU5OTM2ZjNjYmEyOGJlMTQ5Y2ZkMzU1OFwiXG4gIGNvbnN0IHV0eG9idWZmOiBCdWZmZXIgPSBCdWZmZXIuZnJvbSh1dHhvaGV4LCBcImhleFwiKVxuXG4gIC8vIFBheW1lbnRcbiAgY29uc3QgT1BVVFhPc3RyOiBzdHJpbmcgPSBiaW50b29scy5jYjU4RW5jb2RlKHV0eG9idWZmKVxuICAvLyBcIlU5ckZnSzVqamRYbVY4azV0cHFlWGtpbXpyTjNvOWVDQ2NYZXN5aE1CQlp1OU1RSkNEVERvNVduNXBzS3Z6SlZNSnBpTWJka2ZEWGtwN3NLWmRkZkNaZHhwdURteU55N1ZGa2ExOXpNVzRqY3o2RFJRdk5mQTJrdkpZS2s5NnpjN3VpemdwM2kyRllXckI4bXIxc1BKOG9QOVRoNjRHUTV5SGQ4XCJcblxuICAvLyBpbXBsaWVzIGZyb21TdHJpbmcgYW5kIGZyb21CdWZmZXJcbiAgdGVzdChcIkNyZWF0aW9uXCIsICgpOiB2b2lkID0+IHtcbiAgICBjb25zdCB1MTogVVRYTyA9IG5ldyBVVFhPKClcbiAgICB1MS5mcm9tQnVmZmVyKHV0eG9idWZmKVxuICAgIGNvbnN0IHUxaGV4OiBzdHJpbmcgPSB1MS50b0J1ZmZlcigpLnRvU3RyaW5nKFwiaGV4XCIpXG4gICAgZXhwZWN0KHUxaGV4KS50b0JlKHV0eG9oZXgpXG4gIH0pXG5cbiAgdGVzdChcIkVtcHR5IENyZWF0aW9uXCIsICgpOiB2b2lkID0+IHtcbiAgICBjb25zdCB1MTogVVRYTyA9IG5ldyBVVFhPKClcbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdTEudG9CdWZmZXIoKVxuICAgIH0pLnRvVGhyb3coKVxuICB9KVxuXG4gIHRlc3QoXCJDcmVhdGlvbiBvZiBUeXBlXCIsICgpOiB2b2lkID0+IHtcbiAgICBjb25zdCBvcDogVVRYTyA9IG5ldyBVVFhPKClcbiAgICBvcC5mcm9tU3RyaW5nKE9QVVRYT3N0cilcbiAgICBleHBlY3Qob3AuZ2V0T3V0cHV0KCkuZ2V0T3V0cHV0SUQoKSkudG9CZSg3KVxuICB9KVxuXG4gIGRlc2NyaWJlKFwiRnVudGlvbmFsaXR5XCIsICgpOiB2b2lkID0+IHtcbiAgICBjb25zdCB1MTogVVRYTyA9IG5ldyBVVFhPKClcbiAgICB1MS5mcm9tQnVmZmVyKHV0eG9idWZmKVxuICAgIHRlc3QoXCJnZXRBc3NldElEIE5vbkNBXCIsICgpOiB2b2lkID0+IHtcbiAgICAgIGNvbnN0IGFzc2V0SUQ6IEJ1ZmZlciA9IHUxLmdldEFzc2V0SUQoKVxuICAgICAgZXhwZWN0KGFzc2V0SUQudG9TdHJpbmcoXCJoZXhcIiwgMCwgYXNzZXRJRC5sZW5ndGgpKS50b0JlKG91dGFpZClcbiAgICB9KVxuICAgIHRlc3QoXCJnZXRUeElEXCIsICgpOiB2b2lkID0+IHtcbiAgICAgIGNvbnN0IHR4aWQ6IEJ1ZmZlciA9IHUxLmdldFR4SUQoKVxuICAgICAgZXhwZWN0KHR4aWQudG9TdHJpbmcoXCJoZXhcIiwgMCwgdHhpZC5sZW5ndGgpKS50b0JlKG91dHR4aWQpXG4gICAgfSlcbiAgICB0ZXN0KFwiZ2V0T3V0cHV0SWR4XCIsICgpOiB2b2lkID0+IHtcbiAgICAgIGNvbnN0IHR4aWR4OiBCdWZmZXIgPSB1MS5nZXRPdXRwdXRJZHgoKVxuICAgICAgZXhwZWN0KHR4aWR4LnRvU3RyaW5nKFwiaGV4XCIsIDAsIHR4aWR4Lmxlbmd0aCkpLnRvQmUob3V0cHV0aWR4KVxuICAgIH0pXG4gICAgdGVzdChcImdldFVUWE9JRFwiLCAoKTogdm9pZCA9PiB7XG4gICAgICBjb25zdCB0eGlkOiBCdWZmZXIgPSBCdWZmZXIuZnJvbShvdXR0eGlkLCBcImhleFwiKVxuICAgICAgY29uc3QgdHhpZHg6IEJ1ZmZlciA9IEJ1ZmZlci5mcm9tKG91dHB1dGlkeCwgXCJoZXhcIilcbiAgICAgIGNvbnN0IHV0eG9pZDogc3RyaW5nID0gYmludG9vbHMuYnVmZmVyVG9CNTgoQnVmZmVyLmNvbmNhdChbdHhpZCwgdHhpZHhdKSlcbiAgICAgIGV4cGVjdCh1MS5nZXRVVFhPSUQoKSkudG9CZSh1dHhvaWQpXG4gICAgfSlcbiAgICB0ZXN0KFwidG9TdHJpbmdcIiwgKCk6IHZvaWQgPT4ge1xuICAgICAgY29uc3Qgc2VyaWFsaXplZDogc3RyaW5nID0gdTEudG9TdHJpbmcoKVxuICAgICAgZXhwZWN0KHNlcmlhbGl6ZWQpLnRvQmUoYmludG9vbHMuY2I1OEVuY29kZSh1dHhvYnVmZikpXG4gICAgfSlcbiAgfSlcbn0pXG5cbmNvbnN0IHNldE1lcmdlVGVzdGVyID0gKFxuICBpbnB1dDogVVRYT1NldCxcbiAgZXF1YWw6IFVUWE9TZXRbXSxcbiAgbm90RXF1YWw6IFVUWE9TZXRbXVxuKTogYm9vbGVhbiA9PiB7XG4gIGNvbnN0IGluc3RyOiBzdHJpbmcgPSBKU09OLnN0cmluZ2lmeShpbnB1dC5nZXRVVFhPSURzKCkuc29ydCgpKVxuICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgZXF1YWwubGVuZ3RoOyBpKyspIHtcbiAgICBpZiAoSlNPTi5zdHJpbmdpZnkoZXF1YWxbaV0uZ2V0VVRYT0lEcygpLnNvcnQoKSkgIT0gaW5zdHIpIHtcbiAgICAgIHJldHVybiBmYWxzZVxuICAgIH1cbiAgfVxuXG4gIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCBub3RFcXVhbC5sZW5ndGg7IGkrKykge1xuICAgIGlmIChKU09OLnN0cmluZ2lmeShub3RFcXVhbFtpXS5nZXRVVFhPSURzKCkuc29ydCgpKSA9PSBpbnN0cikge1xuICAgICAgcmV0dXJuIGZhbHNlXG4gICAgfVxuICB9XG4gIHJldHVybiB0cnVlXG59XG5cbmRlc2NyaWJlKFwiVVRYT1NldFwiLCAoKTogdm9pZCA9PiB7XG4gIGNvbnN0IHV0eG9zdHJzOiBzdHJpbmdbXSA9IFtcbiAgICBiaW50b29scy5jYjU4RW5jb2RlKFxuICAgICAgQnVmZmVyLmZyb20oXG4gICAgICAgIFwiMDAwMDM4ZDFiOWYxMTM4NjcyZGE2ZmI2YzM1MTI1NTM5Mjc2YTlhY2MyYTY2OGQ2M2JlYTZiYTNjNzk1ZTJlZGIwZjUwMDAwMDAwMTNlMDdlMzhlMmYyMzEyMWJlODc1NjQxMmMxOGRiNzI0NmExNmQyNmVlOTkzNmYzY2JhMjhiZTE0OWNmZDM1NTgwMDAwMDAwNzAwMDAwMDAwMDAwMDRkZDUwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDEwMDAwMDAwMWEzNmZkMGMyZGJjYWIzMTE3MzFkZGU3ZWYxNTE0YmQyNmZjZGM3NGRcIixcbiAgICAgICAgXCJoZXhcIlxuICAgICAgKVxuICAgICksXG4gICAgYmludG9vbHMuY2I1OEVuY29kZShcbiAgICAgIEJ1ZmZlci5mcm9tKFxuICAgICAgICBcIjAwMDBjM2U0ODIzNTcxNTg3ZmUyYmRmYzUwMjY4OWY1YTgyMzhiOWQwZWE3ZjMyNzcxMjRkMTZhZjlkZTBkMmQ5OTExMDAwMDAwMDAzZTA3ZTM4ZTJmMjMxMjFiZTg3NTY0MTJjMThkYjcyNDZhMTZkMjZlZTk5MzZmM2NiYTI4YmUxNDljZmQzNTU4MDAwMDAwMDcwMDAwMDAwMDAwMDAwMDE5MDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAxMDAwMDAwMDFlMWI2YjZhNGJhZDk0ZDJlM2YyMDczMDM3OWI5YmNkNmYxNzYzMThlXCIsXG4gICAgICAgIFwiaGV4XCJcbiAgICAgIClcbiAgICApLFxuICAgIGJpbnRvb2xzLmNiNThFbmNvZGUoXG4gICAgICBCdWZmZXIuZnJvbShcbiAgICAgICAgXCIwMDAwZjI5ZGJhNjFmZGE4ZDU3YTkxMWU3Zjg4MTBmOTM1YmRlODEwZDNmOGQ0OTU0MDQ2ODViZGI4ZDlkODU0NWU4NjAwMDAwMDAwM2UwN2UzOGUyZjIzMTIxYmU4NzU2NDEyYzE4ZGI3MjQ2YTE2ZDI2ZWU5OTM2ZjNjYmEyOGJlMTQ5Y2ZkMzU1ODAwMDAwMDA3MDAwMDAwMDAwMDAwMDAxOTAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMTAwMDAwMDAxZTFiNmI2YTRiYWQ5NGQyZTNmMjA3MzAzNzliOWJjZDZmMTc2MzE4ZVwiLFxuICAgICAgICBcImhleFwiXG4gICAgICApXG4gICAgKVxuICBdXG4gIGNvbnN0IGFkZHJzOiBCdWZmZXJbXSA9IFtcbiAgICBiaW50b29scy5jYjU4RGVjb2RlKFwiRnVCNkx3MkQ2Mk51TTh6cEdMQTRBdmVwcTdlR3NaUmlHXCIpLFxuICAgIGJpbnRvb2xzLmNiNThEZWNvZGUoXCJNYVR2S0djY2JZekN4ekJrSnBiMnpIVzdFMVdSZVpxQjhcIilcbiAgXVxuICB0ZXN0KFwiQ3JlYXRpb25cIiwgKCk6IHZvaWQgPT4ge1xuICAgIGNvbnN0IHNldDogVVRYT1NldCA9IG5ldyBVVFhPU2V0KClcbiAgICBzZXQuYWRkKHV0eG9zdHJzWzBdKVxuICAgIGNvbnN0IHV0eG86IFVUWE8gPSBuZXcgVVRYTygpXG4gICAgdXR4by5mcm9tU3RyaW5nKHV0eG9zdHJzWzBdKVxuICAgIGNvbnN0IHNldEFycmF5OiBVVFhPW10gPSBzZXQuZ2V0QWxsVVRYT3MoKVxuICAgIGV4cGVjdCh1dHhvLnRvU3RyaW5nKCkpLnRvQmUoc2V0QXJyYXlbMF0udG9TdHJpbmcoKSlcbiAgfSlcblxuICB0ZXN0KFwiYmFkIGNyZWF0aW9uXCIsICgpOiB2b2lkID0+IHtcbiAgICBjb25zdCBzZXQ6IFVUWE9TZXQgPSBuZXcgVVRYT1NldCgpXG4gICAgY29uc3QgYmFkOiBzdHJpbmcgPSBiaW50b29scy5jYjU4RW5jb2RlKEJ1ZmZlci5mcm9tKFwiYWFzZGZhc2RcIiwgXCJoZXhcIikpXG4gICAgc2V0LmFkZChiYWQpXG4gICAgY29uc3QgdXR4bzogVVRYTyA9IG5ldyBVVFhPKClcblxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB1dHhvLmZyb21TdHJpbmcoYmFkKVxuICAgIH0pLnRvVGhyb3coKVxuICB9KVxuXG4gIHRlc3QoXCJNdXRsaXBsZSBhZGRcIiwgKCk6IHZvaWQgPT4ge1xuICAgIGNvbnN0IHNldDogVVRYT1NldCA9IG5ldyBVVFhPU2V0KClcbiAgICAvLyBmaXJzdCBhZGRcbiAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgdXR4b3N0cnMubGVuZ3RoOyBpKyspIHtcbiAgICAgIHNldC5hZGQodXR4b3N0cnNbaV0pXG4gICAgfVxuICAgIC8vIHRoZSB2ZXJpZnkgKGRvIHRoZXNlIHN0ZXBzIHNlcGFyYXRlIHRvIGVuc3VyZSBubyBvdmVyd3JpdGVzKVxuICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCB1dHhvc3Rycy5sZW5ndGg7IGkrKykge1xuICAgICAgZXhwZWN0KHNldC5pbmNsdWRlcyh1dHhvc3Ryc1tpXSkpLnRvQmUodHJ1ZSlcbiAgICAgIGNvbnN0IHV0eG86IFVUWE8gPSBuZXcgVVRYTygpXG4gICAgICB1dHhvLmZyb21TdHJpbmcodXR4b3N0cnNbaV0pXG4gICAgICBjb25zdCB2ZXJpdXR4bzogVVRYTyA9IHNldC5nZXRVVFhPKHV0eG8uZ2V0VVRYT0lEKCkpIGFzIFVUWE9cbiAgICAgIGV4cGVjdCh2ZXJpdXR4by50b1N0cmluZygpKS50b0JlKHV0eG9zdHJzW2ldKVxuICAgIH1cbiAgfSlcblxuICB0ZXN0KFwiYWRkQXJyYXlcIiwgKCk6IHZvaWQgPT4ge1xuICAgIGNvbnN0IHNldDogVVRYT1NldCA9IG5ldyBVVFhPU2V0KClcbiAgICBzZXQuYWRkQXJyYXkodXR4b3N0cnMpXG4gICAgZm9yIChsZXQgaTogbnVtYmVyID0gMDsgaSA8IHV0eG9zdHJzLmxlbmd0aDsgaSsrKSB7XG4gICAgICBjb25zdCBlMTogVVRYTyA9IG5ldyBVVFhPKClcbiAgICAgIGUxLmZyb21TdHJpbmcodXR4b3N0cnNbaV0pXG4gICAgICBleHBlY3Qoc2V0LmluY2x1ZGVzKGUxKSkudG9CZSh0cnVlKVxuICAgICAgY29uc3QgdXR4bzogVVRYTyA9IG5ldyBVVFhPKClcbiAgICAgIHV0eG8uZnJvbVN0cmluZyh1dHhvc3Ryc1tpXSlcbiAgICAgIGNvbnN0IHZlcml1dHhvOiBVVFhPID0gc2V0LmdldFVUWE8odXR4by5nZXRVVFhPSUQoKSkgYXMgVVRYT1xuICAgICAgZXhwZWN0KHZlcml1dHhvLnRvU3RyaW5nKCkpLnRvQmUodXR4b3N0cnNbaV0pXG4gICAgfVxuXG4gICAgc2V0LmFkZEFycmF5KHNldC5nZXRBbGxVVFhPcygpKVxuICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCB1dHhvc3Rycy5sZW5ndGg7IGkrKykge1xuICAgICAgY29uc3QgdXR4bzogVVRYTyA9IG5ldyBVVFhPKClcbiAgICAgIHV0eG8uZnJvbVN0cmluZyh1dHhvc3Ryc1tpXSlcbiAgICAgIGV4cGVjdChzZXQuaW5jbHVkZXModXR4bykpLnRvQmUodHJ1ZSlcblxuICAgICAgY29uc3QgdmVyaXV0eG86IFVUWE8gPSBzZXQuZ2V0VVRYTyh1dHhvLmdldFVUWE9JRCgpKSBhcyBVVFhPXG4gICAgICBleHBlY3QodmVyaXV0eG8udG9TdHJpbmcoKSkudG9CZSh1dHhvc3Ryc1tpXSlcbiAgICB9XG5cbiAgICBsZXQgbzogb2JqZWN0ID0gc2V0LnNlcmlhbGl6ZShcImhleFwiKVxuICAgIGxldCBzOiBVVFhPU2V0ID0gbmV3IFVUWE9TZXQoKVxuICAgIHMuZGVzZXJpYWxpemUobylcbiAgICBsZXQgdDogb2JqZWN0ID0gc2V0LnNlcmlhbGl6ZShkaXNwbGF5KVxuICAgIGxldCByOiBVVFhPU2V0ID0gbmV3IFVUWE9TZXQoKVxuICAgIHIuZGVzZXJpYWxpemUodClcbiAgfSlcblxuICB0ZXN0KFwib3ZlcndyaXRpbmcgVVRYT1wiLCAoKTogdm9pZCA9PiB7XG4gICAgY29uc3Qgc2V0OiBVVFhPU2V0ID0gbmV3IFVUWE9TZXQoKVxuICAgIHNldC5hZGRBcnJheSh1dHhvc3RycylcbiAgICBjb25zdCB0ZXN0dXR4bzogVVRYTyA9IG5ldyBVVFhPKClcbiAgICB0ZXN0dXR4by5mcm9tU3RyaW5nKHV0eG9zdHJzWzBdKVxuICAgIGV4cGVjdChzZXQuYWRkKHV0eG9zdHJzWzBdLCB0cnVlKS50b1N0cmluZygpKS50b0JlKHRlc3R1dHhvLnRvU3RyaW5nKCkpXG4gICAgZXhwZWN0KHNldC5hZGQodXR4b3N0cnNbMF0sIGZhbHNlKSkudG9CZVVuZGVmaW5lZCgpXG4gICAgZXhwZWN0KHNldC5hZGRBcnJheSh1dHhvc3RycywgdHJ1ZSkubGVuZ3RoKS50b0JlKDMpXG4gICAgZXhwZWN0KHNldC5hZGRBcnJheSh1dHhvc3RycywgZmFsc2UpLmxlbmd0aCkudG9CZSgwKVxuICB9KVxuXG4gIGRlc2NyaWJlKFwiRnVuY3Rpb25hbGl0eVwiLCAoKTogdm9pZCA9PiB7XG4gICAgbGV0IHNldDogVVRYT1NldFxuICAgIGxldCB1dHhvczogVVRYT1tdXG4gICAgYmVmb3JlRWFjaCgoKTogdm9pZCA9PiB7XG4gICAgICBzZXQgPSBuZXcgVVRYT1NldCgpXG4gICAgICBzZXQuYWRkQXJyYXkodXR4b3N0cnMpXG4gICAgICB1dHhvcyA9IHNldC5nZXRBbGxVVFhPcygpXG4gICAgfSlcblxuICAgIHRlc3QoXCJyZW1vdmVcIiwgKCk6IHZvaWQgPT4ge1xuICAgICAgY29uc3QgdGVzdHV0eG86IFVUWE8gPSBuZXcgVVRYTygpXG4gICAgICB0ZXN0dXR4by5mcm9tU3RyaW5nKHV0eG9zdHJzWzBdKVxuICAgICAgZXhwZWN0KHNldC5yZW1vdmUodXR4b3N0cnNbMF0pLnRvU3RyaW5nKCkpLnRvQmUodGVzdHV0eG8udG9TdHJpbmcoKSlcbiAgICAgIGV4cGVjdChzZXQucmVtb3ZlKHV0eG9zdHJzWzBdKSkudG9CZVVuZGVmaW5lZCgpXG4gICAgICBleHBlY3Qoc2V0LmFkZCh1dHhvc3Ryc1swXSwgZmFsc2UpLnRvU3RyaW5nKCkpLnRvQmUodGVzdHV0eG8udG9TdHJpbmcoKSlcbiAgICAgIGV4cGVjdChzZXQucmVtb3ZlKHV0eG9zdHJzWzBdKS50b1N0cmluZygpKS50b0JlKHRlc3R1dHhvLnRvU3RyaW5nKCkpXG4gICAgfSlcblxuICAgIHRlc3QoXCJyZW1vdmVBcnJheVwiLCAoKTogdm9pZCA9PiB7XG4gICAgICBjb25zdCB0ZXN0dXR4bzogVVRYTyA9IG5ldyBVVFhPKClcbiAgICAgIHRlc3R1dHhvLmZyb21TdHJpbmcodXR4b3N0cnNbMF0pXG4gICAgICBleHBlY3Qoc2V0LnJlbW92ZUFycmF5KHV0eG9zdHJzKS5sZW5ndGgpLnRvQmUoMylcbiAgICAgIGV4cGVjdChzZXQucmVtb3ZlQXJyYXkodXR4b3N0cnMpLmxlbmd0aCkudG9CZSgwKVxuICAgICAgZXhwZWN0KHNldC5hZGQodXR4b3N0cnNbMF0sIGZhbHNlKS50b1N0cmluZygpKS50b0JlKHRlc3R1dHhvLnRvU3RyaW5nKCkpXG4gICAgICBleHBlY3Qoc2V0LnJlbW92ZUFycmF5KHV0eG9zdHJzKS5sZW5ndGgpLnRvQmUoMSlcbiAgICAgIGV4cGVjdChzZXQuYWRkQXJyYXkodXR4b3N0cnMsIGZhbHNlKS5sZW5ndGgpLnRvQmUoMylcbiAgICAgIGV4cGVjdChzZXQucmVtb3ZlQXJyYXkodXR4b3MpLmxlbmd0aCkudG9CZSgzKVxuICAgIH0pXG5cbiAgICB0ZXN0KFwiZ2V0VVRYT0lEc1wiLCAoKTogdm9pZCA9PiB7XG4gICAgICBjb25zdCB1aWRzOiBzdHJpbmdbXSA9IHNldC5nZXRVVFhPSURzKClcbiAgICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCB1dHhvcy5sZW5ndGg7IGkrKykge1xuICAgICAgICBleHBlY3QodWlkcy5pbmRleE9mKHV0eG9zW2ldLmdldFVUWE9JRCgpKSkubm90LnRvQmUoLTEpXG4gICAgICB9XG4gICAgfSlcblxuICAgIHRlc3QoXCJnZXRBbGxVVFhPc1wiLCAoKTogdm9pZCA9PiB7XG4gICAgICBjb25zdCBhbGx1dHhvczogVVRYT1tdID0gc2V0LmdldEFsbFVUWE9zKClcbiAgICAgIGNvbnN0IHVzdHJzOiBzdHJpbmdbXSA9IFtdXG4gICAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgYWxsdXR4b3MubGVuZ3RoOyBpKyspIHtcbiAgICAgICAgdXN0cnMucHVzaChhbGx1dHhvc1tpXS50b1N0cmluZygpKVxuICAgICAgfVxuICAgICAgZm9yIChsZXQgaTogbnVtYmVyID0gMDsgaSA8IHV0eG9zdHJzLmxlbmd0aDsgaSsrKSB7XG4gICAgICAgIGV4cGVjdCh1c3Rycy5pbmRleE9mKHV0eG9zdHJzW2ldKSkubm90LnRvQmUoLTEpXG4gICAgICB9XG4gICAgICBjb25zdCB1aWRzOiBzdHJpbmdbXSA9IHNldC5nZXRVVFhPSURzKClcbiAgICAgIGNvbnN0IGFsbHV0eG9zMjogVVRYT1tdID0gc2V0LmdldEFsbFVUWE9zKHVpZHMpXG4gICAgICBjb25zdCB1c3RyczI6IHN0cmluZ1tdID0gW11cbiAgICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCBhbGx1dHhvcy5sZW5ndGg7IGkrKykge1xuICAgICAgICB1c3RyczIucHVzaChhbGx1dHhvczJbaV0udG9TdHJpbmcoKSlcbiAgICAgIH1cbiAgICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCB1dHhvc3Rycy5sZW5ndGg7IGkrKykge1xuICAgICAgICBleHBlY3QodXN0cnMyLmluZGV4T2YodXR4b3N0cnNbaV0pKS5ub3QudG9CZSgtMSlcbiAgICAgIH1cbiAgICB9KVxuXG4gICAgdGVzdChcImdldFVUWE9JRHMgQnkgQWRkcmVzc1wiLCAoKTogdm9pZCA9PiB7XG4gICAgICBsZXQgdXR4b2lkczogc3RyaW5nW11cbiAgICAgIHV0eG9pZHMgPSBzZXQuZ2V0VVRYT0lEcyhbYWRkcnNbMF1dKVxuICAgICAgZXhwZWN0KHV0eG9pZHMubGVuZ3RoKS50b0JlKDEpXG4gICAgICB1dHhvaWRzID0gc2V0LmdldFVUWE9JRHMoYWRkcnMpXG4gICAgICBleHBlY3QodXR4b2lkcy5sZW5ndGgpLnRvQmUoMylcbiAgICAgIHV0eG9pZHMgPSBzZXQuZ2V0VVRYT0lEcyhhZGRycywgZmFsc2UpXG4gICAgICBleHBlY3QodXR4b2lkcy5sZW5ndGgpLnRvQmUoMylcbiAgICB9KVxuXG4gICAgdGVzdChcImdldEFsbFVUWE9TdHJpbmdzXCIsICgpOiB2b2lkID0+IHtcbiAgICAgIGNvbnN0IHVzdHJzOiBzdHJpbmdbXSA9IHNldC5nZXRBbGxVVFhPU3RyaW5ncygpXG4gICAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgdXR4b3N0cnMubGVuZ3RoOyBpKyspIHtcbiAgICAgICAgZXhwZWN0KHVzdHJzLmluZGV4T2YodXR4b3N0cnNbaV0pKS5ub3QudG9CZSgtMSlcbiAgICAgIH1cbiAgICAgIGNvbnN0IHVpZHM6IHN0cmluZ1tdID0gc2V0LmdldFVUWE9JRHMoKVxuICAgICAgY29uc3QgdXN0cnMyOiBzdHJpbmdbXSA9IHNldC5nZXRBbGxVVFhPU3RyaW5ncyh1aWRzKVxuICAgICAgZm9yIChsZXQgaTogbnVtYmVyID0gMDsgaSA8IHV0eG9zdHJzLmxlbmd0aDsgaSsrKSB7XG4gICAgICAgIGV4cGVjdCh1c3RyczIuaW5kZXhPZih1dHhvc3Ryc1tpXSkpLm5vdC50b0JlKC0xKVxuICAgICAgfVxuICAgIH0pXG5cbiAgICB0ZXN0KFwiZ2V0QWRkcmVzc2VzXCIsICgpOiB2b2lkID0+IHtcbiAgICAgIGV4cGVjdChzZXQuZ2V0QWRkcmVzc2VzKCkuc29ydCgpKS50b1N0cmljdEVxdWFsKGFkZHJzLnNvcnQoKSlcbiAgICB9KVxuXG4gICAgdGVzdChcImdldEJhbGFuY2VcIiwgKCk6IHZvaWQgPT4ge1xuICAgICAgbGV0IGJhbGFuY2UxOiBCTlxuICAgICAgbGV0IGJhbGFuY2UyOiBCTlxuICAgICAgYmFsYW5jZTEgPSBuZXcgQk4oMClcbiAgICAgIGJhbGFuY2UyID0gbmV3IEJOKDApXG4gICAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgdXR4b3MubGVuZ3RoOyBpKyspIHtcbiAgICAgICAgY29uc3QgYXNzZXRJRCA9IHV0eG9zW2ldLmdldEFzc2V0SUQoKVxuICAgICAgICBiYWxhbmNlMS5hZGQoc2V0LmdldEJhbGFuY2UoYWRkcnMsIGFzc2V0SUQpKVxuICAgICAgICBiYWxhbmNlMi5hZGQoKHV0eG9zW2ldLmdldE91dHB1dCgpIGFzIEFtb3VudE91dHB1dCkuZ2V0QW1vdW50KCkpXG4gICAgICB9XG4gICAgICBleHBlY3QoYmFsYW5jZTEudG9TdHJpbmcoKSkudG9CZShiYWxhbmNlMi50b1N0cmluZygpKVxuXG4gICAgICBiYWxhbmNlMSA9IG5ldyBCTigwKVxuICAgICAgYmFsYW5jZTIgPSBuZXcgQk4oMClcbiAgICAgIGNvbnN0IG5vdzogQk4gPSBVbml4Tm93KClcbiAgICAgIGZvciAobGV0IGk6IG51bWJlciA9IDA7IGkgPCB1dHhvcy5sZW5ndGg7IGkrKykge1xuICAgICAgICBjb25zdCBhc3NldElEID0gYmludG9vbHMuY2I1OEVuY29kZSh1dHhvc1tpXS5nZXRBc3NldElEKCkpXG4gICAgICAgIGJhbGFuY2UxLmFkZChzZXQuZ2V0QmFsYW5jZShhZGRycywgYXNzZXRJRCwgbm93KSlcbiAgICAgICAgYmFsYW5jZTIuYWRkKCh1dHhvc1tpXS5nZXRPdXRwdXQoKSBhcyBBbW91bnRPdXRwdXQpLmdldEFtb3VudCgpKVxuICAgICAgfVxuICAgICAgZXhwZWN0KGJhbGFuY2UxLnRvU3RyaW5nKCkpLnRvQmUoYmFsYW5jZTIudG9TdHJpbmcoKSlcbiAgICB9KVxuXG4gICAgdGVzdChcImdldEFzc2V0SURzXCIsICgpOiB2b2lkID0+IHtcbiAgICAgIGNvbnN0IGFzc2V0SURzOiBCdWZmZXJbXSA9IHNldC5nZXRBc3NldElEcygpXG4gICAgICBmb3IgKGxldCBpOiBudW1iZXIgPSAwOyBpIDwgdXR4b3MubGVuZ3RoOyBpKyspIHtcbiAgICAgICAgZXhwZWN0KGFzc2V0SURzKS50b0NvbnRhaW4odXR4b3NbaV0uZ2V0QXNzZXRJRCgpKVxuICAgICAgfVxuICAgICAgY29uc3QgYWRkcmVzc2VzOiBCdWZmZXJbXSA9IHNldC5nZXRBZGRyZXNzZXMoKVxuICAgICAgZXhwZWN0KHNldC5nZXRBc3NldElEcyhhZGRyZXNzZXMpKS50b0VxdWFsKHNldC5nZXRBc3NldElEcygpKVxuICAgIH0pXG5cbiAgICBkZXNjcmliZShcIk1lcmdlIFJ1bGVzXCIsICgpOiB2b2lkID0+IHtcbiAgICAgIGxldCBzZXRBOiBVVFhPU2V0XG4gICAgICBsZXQgc2V0QjogVVRYT1NldFxuICAgICAgbGV0IHNldEM6IFVUWE9TZXRcbiAgICAgIGxldCBzZXREOiBVVFhPU2V0XG4gICAgICBsZXQgc2V0RTogVVRYT1NldFxuICAgICAgbGV0IHNldEY6IFVUWE9TZXRcbiAgICAgIGxldCBzZXRHOiBVVFhPU2V0XG4gICAgICBsZXQgc2V0SDogVVRYT1NldFxuICAgICAgLy8gVGFrZS1vci1MZWF2ZVxuICAgICAgY29uc3QgbmV3dXR4bzogc3RyaW5nID0gYmludG9vbHMuY2I1OEVuY29kZShcbiAgICAgICAgQnVmZmVyLmZyb20oXG4gICAgICAgICAgXCIwMDAwYWNmODg2NDdiM2ZiYWE5ZmRmNDM3OGYzYTBkZjZhNWQxNWQ4ZWZiMDE4YWQ3OGYxMjY5MDM5MGU3OWUxNjg3NjAwMDAwMDAzYWNmODg2NDdiM2ZiYWE5ZmRmNDM3OGYzYTBkZjZhNWQxNWQ4ZWZiMDE4YWQ3OGYxMjY5MDM5MGU3OWUxNjg3NjAwMDAwMDA3MDAwMDAwMDAwMDAxODZhMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMTAwMDAwMDAxZmNlZGE4ZjkwZmNiNWQzMDYxNGI5OWQ3OWZjNGJhYTI5MzA3NzYyNlwiLFxuICAgICAgICAgIFwiaGV4XCJcbiAgICAgICAgKVxuICAgICAgKVxuXG4gICAgICBiZWZvcmVFYWNoKCgpOiB2b2lkID0+IHtcbiAgICAgICAgc2V0QSA9IG5ldyBVVFhPU2V0KClcbiAgICAgICAgc2V0QS5hZGRBcnJheShbdXR4b3N0cnNbMF0sIHV0eG9zdHJzWzJdXSlcblxuICAgICAgICBzZXRCID0gbmV3IFVUWE9TZXQoKVxuICAgICAgICBzZXRCLmFkZEFycmF5KFt1dHhvc3Ryc1sxXSwgdXR4b3N0cnNbMl1dKVxuXG4gICAgICAgIHNldEMgPSBuZXcgVVRYT1NldCgpXG4gICAgICAgIHNldEMuYWRkQXJyYXkoW3V0eG9zdHJzWzBdLCB1dHhvc3Ryc1sxXV0pXG5cbiAgICAgICAgc2V0RCA9IG5ldyBVVFhPU2V0KClcbiAgICAgICAgc2V0RC5hZGRBcnJheShbdXR4b3N0cnNbMV1dKVxuXG4gICAgICAgIHNldEUgPSBuZXcgVVRYT1NldCgpXG4gICAgICAgIHNldEUuYWRkQXJyYXkoW10pIC8vIGVtcHR5IHNldFxuXG4gICAgICAgIHNldEYgPSBuZXcgVVRYT1NldCgpXG4gICAgICAgIHNldEYuYWRkQXJyYXkodXR4b3N0cnMpIC8vIGZ1bGwgc2V0LCBzZXBhcmF0ZSBmcm9tIHNlbGZcblxuICAgICAgICBzZXRHID0gbmV3IFVUWE9TZXQoKVxuICAgICAgICBzZXRHLmFkZEFycmF5KFtuZXd1dHhvLCAuLi51dHhvc3Ryc10pIC8vIGZ1bGwgc2V0IHdpdGggbmV3IGVsZW1lbnRcblxuICAgICAgICBzZXRIID0gbmV3IFVUWE9TZXQoKVxuICAgICAgICBzZXRILmFkZEFycmF5KFtuZXd1dHhvXSkgLy8gc2V0IHdpdGggb25seSBhIG5ldyBlbGVtZW50XG4gICAgICB9KVxuXG4gICAgICB0ZXN0KFwidW5rbm93biBtZXJnZSBydWxlXCIsICgpOiB2b2lkID0+IHtcbiAgICAgICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgICAgICBzZXQubWVyZ2VCeVJ1bGUoc2V0QSwgXCJFUlJPUlwiKVxuICAgICAgICB9KS50b1Rocm93KClcbiAgICAgIH0pXG5cbiAgICAgIHRlc3QoXCJpbnRlcnNlY3Rpb25cIiwgKCk6IHZvaWQgPT4ge1xuICAgICAgICBsZXQgcmVzdWx0czogVVRYT1NldFxuICAgICAgICBsZXQgdGVzdDogYm9vbGVhblxuXG4gICAgICAgIHJlc3VsdHMgPSBzZXQubWVyZ2VCeVJ1bGUoc2V0QSwgXCJpbnRlcnNlY3Rpb25cIilcbiAgICAgICAgdGVzdCA9IHNldE1lcmdlVGVzdGVyKFxuICAgICAgICAgIHJlc3VsdHMsXG4gICAgICAgICAgW3NldEFdLFxuICAgICAgICAgIFtzZXRCLCBzZXRDLCBzZXRELCBzZXRFLCBzZXRGLCBzZXRHLCBzZXRIXVxuICAgICAgICApXG4gICAgICAgIGV4cGVjdCh0ZXN0KS50b0JlKHRydWUpXG5cbiAgICAgICAgcmVzdWx0cyA9IHNldC5tZXJnZUJ5UnVsZShzZXRGLCBcImludGVyc2VjdGlvblwiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0Rl0sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEQsIHNldEUsIHNldEcsIHNldEhdXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEcsIFwiaW50ZXJzZWN0aW9uXCIpXG4gICAgICAgIHRlc3QgPSBzZXRNZXJnZVRlc3RlcihcbiAgICAgICAgICByZXN1bHRzLFxuICAgICAgICAgIFtzZXRGXSxcbiAgICAgICAgICBbc2V0QSwgc2V0Qiwgc2V0Qywgc2V0RCwgc2V0RSwgc2V0Rywgc2V0SF1cbiAgICAgICAgKVxuICAgICAgICBleHBlY3QodGVzdCkudG9CZSh0cnVlKVxuXG4gICAgICAgIHJlc3VsdHMgPSBzZXQubWVyZ2VCeVJ1bGUoc2V0SCwgXCJpbnRlcnNlY3Rpb25cIilcbiAgICAgICAgdGVzdCA9IHNldE1lcmdlVGVzdGVyKFxuICAgICAgICAgIHJlc3VsdHMsXG4gICAgICAgICAgW3NldEVdLFxuICAgICAgICAgIFtzZXRBLCBzZXRCLCBzZXRDLCBzZXRELCBzZXRGLCBzZXRHLCBzZXRIXVxuICAgICAgICApXG4gICAgICAgIGV4cGVjdCh0ZXN0KS50b0JlKHRydWUpXG4gICAgICB9KVxuXG4gICAgICB0ZXN0KFwiZGlmZmVyZW5jZVNlbGZcIiwgKCk6IHZvaWQgPT4ge1xuICAgICAgICBsZXQgcmVzdWx0czogVVRYT1NldFxuICAgICAgICBsZXQgdGVzdDogYm9vbGVhblxuXG4gICAgICAgIHJlc3VsdHMgPSBzZXQubWVyZ2VCeVJ1bGUoc2V0QSwgXCJkaWZmZXJlbmNlU2VsZlwiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0RF0sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEUsIHNldEYsIHNldEcsIHNldEhdXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEYsIFwiZGlmZmVyZW5jZVNlbGZcIilcbiAgICAgICAgdGVzdCA9IHNldE1lcmdlVGVzdGVyKFxuICAgICAgICAgIHJlc3VsdHMsXG4gICAgICAgICAgW3NldEVdLFxuICAgICAgICAgIFtzZXRBLCBzZXRCLCBzZXRDLCBzZXRELCBzZXRGLCBzZXRHLCBzZXRIXVxuICAgICAgICApXG4gICAgICAgIGV4cGVjdCh0ZXN0KS50b0JlKHRydWUpXG5cbiAgICAgICAgcmVzdWx0cyA9IHNldC5tZXJnZUJ5UnVsZShzZXRHLCBcImRpZmZlcmVuY2VTZWxmXCIpXG4gICAgICAgIHRlc3QgPSBzZXRNZXJnZVRlc3RlcihcbiAgICAgICAgICByZXN1bHRzLFxuICAgICAgICAgIFtzZXRFXSxcbiAgICAgICAgICBbc2V0QSwgc2V0Qiwgc2V0Qywgc2V0RCwgc2V0Riwgc2V0Rywgc2V0SF1cbiAgICAgICAgKVxuICAgICAgICBleHBlY3QodGVzdCkudG9CZSh0cnVlKVxuXG4gICAgICAgIHJlc3VsdHMgPSBzZXQubWVyZ2VCeVJ1bGUoc2V0SCwgXCJkaWZmZXJlbmNlU2VsZlwiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0Rl0sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEQsIHNldEUsIHNldEcsIHNldEhdXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcbiAgICAgIH0pXG5cbiAgICAgIHRlc3QoXCJkaWZmZXJlbmNlTmV3XCIsICgpOiB2b2lkID0+IHtcbiAgICAgICAgbGV0IHJlc3VsdHM6IFVUWE9TZXRcbiAgICAgICAgbGV0IHRlc3Q6IGJvb2xlYW5cblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEEsIFwiZGlmZmVyZW5jZU5ld1wiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0RV0sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEQsIHNldEYsIHNldEcsIHNldEhdXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEYsIFwiZGlmZmVyZW5jZU5ld1wiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0RV0sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEQsIHNldEYsIHNldEcsIHNldEhdXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEcsIFwiZGlmZmVyZW5jZU5ld1wiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0SF0sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEQsIHNldEUsIHNldEYsIHNldEddXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEgsIFwiZGlmZmVyZW5jZU5ld1wiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0SF0sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEQsIHNldEUsIHNldEYsIHNldEddXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcbiAgICAgIH0pXG5cbiAgICAgIHRlc3QoXCJzeW1EaWZmZXJlbmNlXCIsICgpOiB2b2lkID0+IHtcbiAgICAgICAgbGV0IHJlc3VsdHM6IFVUWE9TZXRcbiAgICAgICAgbGV0IHRlc3Q6IGJvb2xlYW5cblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEEsIFwic3ltRGlmZmVyZW5jZVwiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0RF0sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEUsIHNldEYsIHNldEcsIHNldEhdXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEYsIFwic3ltRGlmZmVyZW5jZVwiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0RV0sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEQsIHNldEYsIHNldEcsIHNldEhdXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEcsIFwic3ltRGlmZmVyZW5jZVwiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0SF0sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEQsIHNldEUsIHNldEYsIHNldEddXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEgsIFwic3ltRGlmZmVyZW5jZVwiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0R10sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEQsIHNldEUsIHNldEYsIHNldEhdXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcbiAgICAgIH0pXG5cbiAgICAgIHRlc3QoXCJ1bmlvblwiLCAoKTogdm9pZCA9PiB7XG4gICAgICAgIGxldCByZXN1bHRzOiBVVFhPU2V0XG4gICAgICAgIGxldCB0ZXN0OiBib29sZWFuXG5cbiAgICAgICAgcmVzdWx0cyA9IHNldC5tZXJnZUJ5UnVsZShzZXRBLCBcInVuaW9uXCIpXG4gICAgICAgIHRlc3QgPSBzZXRNZXJnZVRlc3RlcihcbiAgICAgICAgICByZXN1bHRzLFxuICAgICAgICAgIFtzZXRGXSxcbiAgICAgICAgICBbc2V0QSwgc2V0Qiwgc2V0Qywgc2V0RCwgc2V0RSwgc2V0Rywgc2V0SF1cbiAgICAgICAgKVxuICAgICAgICBleHBlY3QodGVzdCkudG9CZSh0cnVlKVxuXG4gICAgICAgIHJlc3VsdHMgPSBzZXQubWVyZ2VCeVJ1bGUoc2V0RiwgXCJ1bmlvblwiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0Rl0sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEQsIHNldEUsIHNldEcsIHNldEhdXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEcsIFwidW5pb25cIilcbiAgICAgICAgdGVzdCA9IHNldE1lcmdlVGVzdGVyKFxuICAgICAgICAgIHJlc3VsdHMsXG4gICAgICAgICAgW3NldEddLFxuICAgICAgICAgIFtzZXRBLCBzZXRCLCBzZXRDLCBzZXRELCBzZXRFLCBzZXRGLCBzZXRIXVxuICAgICAgICApXG4gICAgICAgIGV4cGVjdCh0ZXN0KS50b0JlKHRydWUpXG5cbiAgICAgICAgcmVzdWx0cyA9IHNldC5tZXJnZUJ5UnVsZShzZXRILCBcInVuaW9uXCIpXG4gICAgICAgIHRlc3QgPSBzZXRNZXJnZVRlc3RlcihcbiAgICAgICAgICByZXN1bHRzLFxuICAgICAgICAgIFtzZXRHXSxcbiAgICAgICAgICBbc2V0QSwgc2V0Qiwgc2V0Qywgc2V0RCwgc2V0RSwgc2V0Riwgc2V0SF1cbiAgICAgICAgKVxuICAgICAgICBleHBlY3QodGVzdCkudG9CZSh0cnVlKVxuICAgICAgfSlcblxuICAgICAgdGVzdChcInVuaW9uTWludXNOZXdcIiwgKCk6IHZvaWQgPT4ge1xuICAgICAgICBsZXQgcmVzdWx0czogVVRYT1NldFxuICAgICAgICBsZXQgdGVzdDogYm9vbGVhblxuXG4gICAgICAgIHJlc3VsdHMgPSBzZXQubWVyZ2VCeVJ1bGUoc2V0QSwgXCJ1bmlvbk1pbnVzTmV3XCIpXG4gICAgICAgIHRlc3QgPSBzZXRNZXJnZVRlc3RlcihcbiAgICAgICAgICByZXN1bHRzLFxuICAgICAgICAgIFtzZXREXSxcbiAgICAgICAgICBbc2V0QSwgc2V0Qiwgc2V0Qywgc2V0RSwgc2V0Riwgc2V0Rywgc2V0SF1cbiAgICAgICAgKVxuICAgICAgICBleHBlY3QodGVzdCkudG9CZSh0cnVlKVxuXG4gICAgICAgIHJlc3VsdHMgPSBzZXQubWVyZ2VCeVJ1bGUoc2V0RiwgXCJ1bmlvbk1pbnVzTmV3XCIpXG4gICAgICAgIHRlc3QgPSBzZXRNZXJnZVRlc3RlcihcbiAgICAgICAgICByZXN1bHRzLFxuICAgICAgICAgIFtzZXRFXSxcbiAgICAgICAgICBbc2V0QSwgc2V0Qiwgc2V0Qywgc2V0RCwgc2V0Riwgc2V0Rywgc2V0SF1cbiAgICAgICAgKVxuICAgICAgICBleHBlY3QodGVzdCkudG9CZSh0cnVlKVxuXG4gICAgICAgIHJlc3VsdHMgPSBzZXQubWVyZ2VCeVJ1bGUoc2V0RywgXCJ1bmlvbk1pbnVzTmV3XCIpXG4gICAgICAgIHRlc3QgPSBzZXRNZXJnZVRlc3RlcihcbiAgICAgICAgICByZXN1bHRzLFxuICAgICAgICAgIFtzZXRFXSxcbiAgICAgICAgICBbc2V0QSwgc2V0Qiwgc2V0Qywgc2V0RCwgc2V0Riwgc2V0Rywgc2V0SF1cbiAgICAgICAgKVxuICAgICAgICBleHBlY3QodGVzdCkudG9CZSh0cnVlKVxuXG4gICAgICAgIHJlc3VsdHMgPSBzZXQubWVyZ2VCeVJ1bGUoc2V0SCwgXCJ1bmlvbk1pbnVzTmV3XCIpXG4gICAgICAgIHRlc3QgPSBzZXRNZXJnZVRlc3RlcihcbiAgICAgICAgICByZXN1bHRzLFxuICAgICAgICAgIFtzZXRGXSxcbiAgICAgICAgICBbc2V0QSwgc2V0Qiwgc2V0Qywgc2V0RCwgc2V0RSwgc2V0Rywgc2V0SF1cbiAgICAgICAgKVxuICAgICAgICBleHBlY3QodGVzdCkudG9CZSh0cnVlKVxuICAgICAgfSlcblxuICAgICAgdGVzdChcInVuaW9uTWludXNTZWxmXCIsICgpOiB2b2lkID0+IHtcbiAgICAgICAgbGV0IHJlc3VsdHM6IFVUWE9TZXRcbiAgICAgICAgbGV0IHRlc3Q6IGJvb2xlYW5cblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEEsIFwidW5pb25NaW51c1NlbGZcIilcbiAgICAgICAgdGVzdCA9IHNldE1lcmdlVGVzdGVyKFxuICAgICAgICAgIHJlc3VsdHMsXG4gICAgICAgICAgW3NldEVdLFxuICAgICAgICAgIFtzZXRBLCBzZXRCLCBzZXRDLCBzZXRELCBzZXRGLCBzZXRHLCBzZXRIXVxuICAgICAgICApXG4gICAgICAgIGV4cGVjdCh0ZXN0KS50b0JlKHRydWUpXG5cbiAgICAgICAgcmVzdWx0cyA9IHNldC5tZXJnZUJ5UnVsZShzZXRGLCBcInVuaW9uTWludXNTZWxmXCIpXG4gICAgICAgIHRlc3QgPSBzZXRNZXJnZVRlc3RlcihcbiAgICAgICAgICByZXN1bHRzLFxuICAgICAgICAgIFtzZXRFXSxcbiAgICAgICAgICBbc2V0QSwgc2V0Qiwgc2V0Qywgc2V0RCwgc2V0Riwgc2V0Rywgc2V0SF1cbiAgICAgICAgKVxuICAgICAgICBleHBlY3QodGVzdCkudG9CZSh0cnVlKVxuXG4gICAgICAgIHJlc3VsdHMgPSBzZXQubWVyZ2VCeVJ1bGUoc2V0RywgXCJ1bmlvbk1pbnVzU2VsZlwiKVxuICAgICAgICB0ZXN0ID0gc2V0TWVyZ2VUZXN0ZXIoXG4gICAgICAgICAgcmVzdWx0cyxcbiAgICAgICAgICBbc2V0SF0sXG4gICAgICAgICAgW3NldEEsIHNldEIsIHNldEMsIHNldEQsIHNldEUsIHNldEYsIHNldEddXG4gICAgICAgIClcbiAgICAgICAgZXhwZWN0KHRlc3QpLnRvQmUodHJ1ZSlcblxuICAgICAgICByZXN1bHRzID0gc2V0Lm1lcmdlQnlSdWxlKHNldEgsIFwidW5pb25NaW51c1NlbGZcIilcbiAgICAgICAgdGVzdCA9IHNldE1lcmdlVGVzdGVyKFxuICAgICAgICAgIHJlc3VsdHMsXG4gICAgICAgICAgW3NldEhdLFxuICAgICAgICAgIFtzZXRBLCBzZXRCLCBzZXRDLCBzZXRELCBzZXRFLCBzZXRGLCBzZXRHXVxuICAgICAgICApXG4gICAgICAgIGV4cGVjdCh0ZXN0KS50b0JlKHRydWUpXG4gICAgICB9KVxuICAgIH0pXG4gIH0pXG59KVxuIl19