"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.SECP256k1KeyChain = exports.SECP256k1KeyPair = void 0;
/**
 * @packageDocumentation
 * @module Common-SECP256k1KeyChain
 */
const buffer_1 = require("buffer/");
const elliptic = __importStar(require("elliptic"));
const create_hash_1 = __importDefault(require("create-hash"));
const bintools_1 = __importDefault(require("../utils/bintools"));
const keychain_1 = require("./keychain");
const errors_1 = require("../utils/errors");
const utils_1 = require("../utils");
/**
 * @ignore
 */
const EC = elliptic.ec;
/**
 * @ignore
 */
const ec = new EC("secp256k1");
/**
 * @ignore
 */
const ecparams = ec.curve;
/**
 * @ignore
 */
const BN = ecparams.n.constructor;
/**
 * @ignore
 */
const bintools = bintools_1.default.getInstance();
const serialization = utils_1.Serialization.getInstance();
/**
 * Class for representing a private and public keypair on the Platform Chain.
 */
class SECP256k1KeyPair extends keychain_1.StandardKeyPair {
    constructor(hrp, chainID) {
        super();
        this.chainID = "";
        this.hrp = "";
        this.chainID = chainID;
        this.hrp = hrp;
        this.generateKey();
    }
    /**
     * @ignore
     */
    _sigFromSigBuffer(sig) {
        const r = new BN(bintools.copyFrom(sig, 0, 32));
        const s = new BN(bintools.copyFrom(sig, 32, 64));
        const recoveryParam = bintools
            .copyFrom(sig, 64, 65)
            .readUIntBE(0, 1);
        const sigOpt = {
            r: r,
            s: s,
            recoveryParam: recoveryParam
        };
        return sigOpt;
    }
    /**
     * Generates a new keypair.
     */
    generateKey() {
        this.keypair = ec.genKeyPair();
        // doing hex translation to get Buffer class
        this.privk = buffer_1.Buffer.from(this.keypair.getPrivate("hex").padStart(64, "0"), "hex");
        this.pubk = buffer_1.Buffer.from(this.keypair.getPublic(true, "hex").padStart(66, "0"), "hex");
    }
    /**
     * Imports a private key and generates the appropriate public key.
     *
     * @param privk A {@link https://github.com/feross/buffer|Buffer} representing the private key
     *
     * @returns true on success, false on failure
     */
    importKey(privk) {
        this.keypair = ec.keyFromPrivate(privk.toString("hex"), "hex");
        // doing hex translation to get Buffer class
        try {
            this.privk = buffer_1.Buffer.from(this.keypair.getPrivate("hex").padStart(64, "0"), "hex");
            this.pubk = buffer_1.Buffer.from(this.keypair.getPublic(true, "hex").padStart(66, "0"), "hex");
            return true; // silly I know, but the interface requires so it returns true on success, so if Buffer fails validation...
        }
        catch (error) {
            return false;
        }
    }
    /**
     * Returns the address as a {@link https://github.com/feross/buffer|Buffer}.
     *
     * @returns A {@link https://github.com/feross/buffer|Buffer} representation of the address
     */
    getAddress() {
        return SECP256k1KeyPair.addressFromPublicKey(this.pubk);
    }
    /**
     * Returns the address's string representation.
     *
     * @returns A string representation of the address
     */
    getAddressString() {
        const addr = SECP256k1KeyPair.addressFromPublicKey(this.pubk);
        const type = "bech32";
        return serialization.bufferToType(addr, type, this.hrp, this.chainID);
    }
    /**
     * Returns an address given a public key.
     *
     * @param pubk A {@link https://github.com/feross/buffer|Buffer} representing the public key
     *
     * @returns A {@link https://github.com/feross/buffer|Buffer} for the address of the public key.
     */
    static addressFromPublicKey(pubk) {
        if (pubk.length === 65) {
            /* istanbul ignore next */
            pubk = buffer_1.Buffer.from(ec.keyFromPublic(pubk).getPublic(true, "hex").padStart(66, "0"), "hex"); // make compact, stick back into buffer
        }
        if (pubk.length === 33) {
            const sha256 = buffer_1.Buffer.from((0, create_hash_1.default)("sha256").update(pubk).digest());
            const ripesha = buffer_1.Buffer.from((0, create_hash_1.default)("ripemd160").update(sha256).digest());
            return ripesha;
        }
        /* istanbul ignore next */
        throw new errors_1.PublicKeyError("Unable to make address.");
    }
    /**
     * Returns a string representation of the private key.
     *
     * @returns A cb58 serialized string representation of the private key
     */
    getPrivateKeyString() {
        return `PrivateKey-${bintools.cb58Encode(this.privk)}`;
    }
    /**
     * Returns the public key.
     *
     * @returns A cb58 serialized string representation of the public key
     */
    getPublicKeyString() {
        return bintools.cb58Encode(this.pubk);
    }
    /**
     * Takes a message, signs it, and returns the signature.
     *
     * @param msg The message to sign, be sure to hash first if expected
     *
     * @returns A {@link https://github.com/feross/buffer|Buffer} containing the signature
     */
    sign(msg) {
        const sigObj = this.keypair.sign(msg, undefined, {
            canonical: true
        });
        const recovery = buffer_1.Buffer.alloc(1);
        recovery.writeUInt8(sigObj.recoveryParam, 0);
        const r = buffer_1.Buffer.from(sigObj.r.toArray("be", 32)); //we have to skip native Buffer class, so this is the way
        const s = buffer_1.Buffer.from(sigObj.s.toArray("be", 32)); //we have to skip native Buffer class, so this is the way
        const result = buffer_1.Buffer.concat([r, s, recovery], 65);
        return result;
    }
    /**
     * Verifies that the private key associated with the provided public key produces the signature associated with the given message.
     *
     * @param msg The message associated with the signature
     * @param sig The signature of the signed message
     *
     * @returns True on success, false on failure
     */
    verify(msg, sig) {
        const sigObj = this._sigFromSigBuffer(sig);
        return ec.verify(msg, sigObj, this.keypair);
    }
    /**
     * Recovers the public key of a message signer from a message and its associated signature.
     *
     * @param msg The message that's signed
     * @param sig The signature that's signed on the message
     *
     * @returns A {@link https://github.com/feross/buffer|Buffer} containing the public key of the signer
     */
    recover(msg, sig) {
        const sigObj = this._sigFromSigBuffer(sig);
        const pubk = ec.recoverPubKey(msg, sigObj, sigObj.recoveryParam);
        return buffer_1.Buffer.from(pubk.encodeCompressed());
    }
    /**
     * Returns the chainID associated with this key.
     *
     * @returns The [[KeyPair]]'s chainID
     */
    getChainID() {
        return this.chainID;
    }
    /**
     * Sets the the chainID associated with this key.
     *
     * @param chainID String for the chainID
     */
    setChainID(chainID) {
        this.chainID = chainID;
    }
    /**
     * Returns the Human-Readable-Part of the network associated with this key.
     *
     * @returns The [[KeyPair]]'s Human-Readable-Part of the network's Bech32 addressing scheme
     */
    getHRP() {
        return this.hrp;
    }
    /**
     * Sets the the Human-Readable-Part of the network associated with this key.
     *
     * @param hrp String for the Human-Readable-Part of Bech32 addresses
     */
    setHRP(hrp) {
        this.hrp = hrp;
    }
}
exports.SECP256k1KeyPair = SECP256k1KeyPair;
/**
 * Class for representing a key chain in Avalanche.
 *
 * @typeparam SECP256k1KeyPair Class extending [[StandardKeyPair]] which is used as the key in [[SECP256k1KeyChain]]
 */
class SECP256k1KeyChain extends keychain_1.StandardKeyChain {
    addKey(newKey) {
        super.addKey(newKey);
    }
}
exports.SECP256k1KeyChain = SECP256k1KeyChain;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoic2VjcDI1NmsxLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vc3JjL2NvbW1vbi9zZWNwMjU2azEudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7QUFBQTs7O0dBR0c7QUFDSCxvQ0FBZ0M7QUFDaEMsbURBQW9DO0FBQ3BDLDhEQUFvQztBQUNwQyxpRUFBd0M7QUFDeEMseUNBQThEO0FBQzlELDRDQUFnRDtBQUVoRCxvQ0FBd0Q7QUFFeEQ7O0dBRUc7QUFDSCxNQUFNLEVBQUUsR0FBdUIsUUFBUSxDQUFDLEVBQUUsQ0FBQTtBQUUxQzs7R0FFRztBQUNILE1BQU0sRUFBRSxHQUFnQixJQUFJLEVBQUUsQ0FBQyxXQUFXLENBQUMsQ0FBQTtBQUUzQzs7R0FFRztBQUNILE1BQU0sUUFBUSxHQUFRLEVBQUUsQ0FBQyxLQUFLLENBQUE7QUFFOUI7O0dBRUc7QUFDSCxNQUFNLEVBQUUsR0FBUSxRQUFRLENBQUMsQ0FBQyxDQUFDLFdBQVcsQ0FBQTtBQUV0Qzs7R0FFRztBQUNILE1BQU0sUUFBUSxHQUFhLGtCQUFRLENBQUMsV0FBVyxFQUFFLENBQUE7QUFDakQsTUFBTSxhQUFhLEdBQWtCLHFCQUFhLENBQUMsV0FBVyxFQUFFLENBQUE7QUFFaEU7O0dBRUc7QUFDSCxNQUFzQixnQkFBaUIsU0FBUSwwQkFBZTtJQW9ONUQsWUFBWSxHQUFXLEVBQUUsT0FBZTtRQUN0QyxLQUFLLEVBQUUsQ0FBQTtRQW5OQyxZQUFPLEdBQVcsRUFBRSxDQUFBO1FBQ3BCLFFBQUcsR0FBVyxFQUFFLENBQUE7UUFtTnhCLElBQUksQ0FBQyxPQUFPLEdBQUcsT0FBTyxDQUFBO1FBQ3RCLElBQUksQ0FBQyxHQUFHLEdBQUcsR0FBRyxDQUFBO1FBQ2QsSUFBSSxDQUFDLFdBQVcsRUFBRSxDQUFBO0lBQ3BCLENBQUM7SUFwTkQ7O09BRUc7SUFDTyxpQkFBaUIsQ0FBQyxHQUFXO1FBQ3JDLE1BQU0sQ0FBQyxHQUFZLElB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QU8sSUFBSSxDQUFDLE9BQU8sQ0FBQTtJQUNyQixDQUFDO0lBRUQ7Ozs7T0FJRztJQUNILFVBQVUsQ0FBQyxPQUFlO1FBQ3hCLElBQUksQ0FBQyxPQUFPLEdBQUcsT0FBTyxDQUFBO0lBQ3hCLENBQUM7SUFFRDs7OztPQUlHO0lBQ0gsTUFBTTtRQUNKLE9BQU8sSUFBSSxDQUFDLEdBQUcsQ0FBQTtJQUNqQixDQUFDO0lBRUQ7Ozs7T0FJRztJQUNILE1BQU0sQ0FBQyxHQUFXO1FBQ2hCLElBQUksQ0FBQyxHQUFHLEdBQUcsR0FBRyxDQUFBO0lBQ2hCLENBQUM7Q0FRRjtBQTFORCw0Q0EwTkM7QUFFRDs7OztHQUlHO0FBQ0gsTUFBc0IsaUJBRXBCLFNBQVEsMkJBQTZCO0lBUXJDLE1BQU0sQ0FBQyxNQUFtQjtRQUN4QixLQUFLLENBQUMsTUFBTSxDQUFDLE1BQU0sQ0FBQyxDQUFBO0lBQ3RCLENBQUM7Q0FVRjtBQXRCRCw4Q0FzQkMiLCJzb3VyY2VzQ29udGVudCI6WyIvKipcbiAqIEBwYWNrYWdlRG9jdW1lbnRhdGlvblxuICogQG1vZHVsZSBDb21tb24tU0VDUDI1NmsxS2V5Q2hhaW5cbiAqL1xuaW1wb3J0IHsgQnVmZmVyIH0gZnJvbSBcImJ1ZmZlci9cIlxuaW1wb3J0ICogYXMgZWxsaXB0aWMgZnJvbSBcImVsbGlwdGljXCJcbmltcG9ydCBjcmVhdGVIYXNoIGZyb20gXCJjcmVhdGUtaGFzaFwiXG5pbXBvcnQgQmluVG9vbHMgZnJvbSBcIi4uL3V0aWxzL2JpbnRvb2xzXCJcbmltcG9ydCB7IFN0YW5kYXJkS2V5UGFpciwgU3RhbmRhcmRLZXlDaGFpbiB9IGZyb20gXCIuL2tleWNoYWluXCJcbmltcG9ydCB7IFB1YmxpY0tleUVycm9yIH0gZnJvbSBcIi4uL3V0aWxzL2Vycm9yc1wiXG5pbXBvcnQgeyBCTklucHV0IH0gZnJvbSBcImVsbGlwdGljXCJcbmltcG9ydCB7IFNlcmlhbGl6YXRpb24sIFNlcmlhbGl6ZWRUeXBlIH0gZnJvbSBcIi4uL3V0aWxzXCJcblxuLyoqXG4gKiBAaWdub3JlXG4gKi9cbmNvbnN0IEVDOiB0eXBlb2YgZWxsaXB0aWMuZWMgPSBlbGxpcHRpYy5lY1xuXG4vKipcbiAqIEBpZ25vcmVcbiAqL1xuY29uc3QgZWM6IGVsbGlwdGljLmVjID0gbmV3IEVDKFwic2VjcDI1NmsxXCIpXG5cbi8qKlxuICogQGlnbm9yZVxuICovXG5jb25zdCBlY3BhcmFtczogYW55ID0gZWMuY3VydmVcblxuLyoqXG4gKiBAaWdub3JlXG4gKi9cbmNvbnN0IEJOOiBhbnkgPSBlY3BhcmFtcy5uLmNvbnN0cnVjdG9yXG5cbi8qKlxuICogQGlnbm9yZVxuICovXG5jb25zdCBiaW50b29sczogQmluVG9vbHMgPSBCaW5Ub29scy5nZXRJbnN0YW5jZSgpXG5jb25zdCBzZXJpYWxpemF0aW9uOiBTZXJpYWxpemF0aW9uID0gU2VyaWFsaXphdGlvbi5nZXRJbnN0YW5jZSgpXG5cbi8qKlxuICogQ2xhc3MgZm9yIHJlcHJlc2VudGluZyBhIHByaXZhdGUgYW5kIHB1YmxpYyBrZXlwYWlyIG9uIHRoZSBQbGF0Zm9ybSBDaGFpbi5cbiAqL1xuZXhwb3J0IGFic3RyYWN0IGNsYXNzIFNFQ1AyNTZrMUtleVBhaXIgZXh0ZW5kcyBTdGFuZGFyZEtleVBhaXIge1xuICBwcm90ZWN0ZWQga2V5cGFpcjogZWxsaXB0aWMuZWMuS2V5UGFpclxuICBwcm90ZWN0ZWQgY2hhaW5JRDogc3RyaW5nID0gXCJcIlxuICBwcm90ZWN0ZWQgaHJwOiBzdHJpbmcgPSBcIlwiXG5cbiAgLyoqXG4gICAqIEBpZ25vcmVcbiAgICovXG4gIHByb3RlY3RlZCBfc2lnRnJvbVNpZ0J1ZmZlcihzaWc6IEJ1ZmZlcik6IGVsbGlwdGljLmVjLlNpZ25hdHVyZU9wdGlvbnMge1xuICAgIGNvbnN0IHI6IEJOSW5wdXQgPSBuZXcgQk4oYmludG9vbHMuY29weUZyb20oc2lnLCAwLCAzMikpXG4gICAgY29uc3QgczogQk5JbnB1dCA9IG5ldyBCTihiaW50b29scy5jb3B5RnJvbShzaWcsIDMyLCA2NCkpXG4gICAgY29uc3QgcmVjb3ZlcnlQYXJhbTogbnVtYmVyID0gYmludG9vbHNcbiAgICAgIC5jb3B5RnJvbShzaWcsIDY0LCA2NSlcbiAgICAgIC5yZWFkVUludEJFKDAsIDEpXG4gICAgY29uc3Qgc2lnT3B0ID0ge1xuICAgICAgcjogcixcbiAgICAgIHM6IHMsXG4gICAgICByZWNvdmVyeVBhcmFtOiByZWNvdmVyeVBhcmFtXG4gICAgfVxuICAgIHJldHVybiBzaWdPcHRcbiAgfVxuXG4gIC8qKlxuICAgKiBHZW5lcmF0ZXMgYSBuZXcga2V5cGFpci5cbiAgICovXG4gIGdlbmVyYXRlS2V5KCkge1xuICAgIHRoaXMua2V5cGFpciA9IGVjLmdlbktleVBhaXIoKVxuXG4gICAgLy8gZG9pbmcgaGV4IHRyYW5zbGF0aW9uIHRvIGdldCBCdWZmZXIgY2xhc3NcbiAgICB0aGlzLnByaXZrID0gQnVmZmVyLmZyb20oXG4gICAgICB0aGlzLmtleXBhaXIuZ2V0UHJpdmF0ZShcImhleFwiKS5wYWRTdGFydCg2NCwgXCIwXCIpLFxuICAgICAgXCJoZXhcIlxuICAgIClcbiAgICB0aGlzLnB1YmsgPSBCdWZmZXIuZnJvbShcbiAgICAgIHRoaXMua2V5cGFpci5nZXRQdWJsaWModHJ1ZSwgXCJoZXhcIikucGFkU3RhcnQoNjYsIFwiMFwiKSxcbiAgICAgIFwiaGV4XCJcbiAgICApXG4gIH1cblxuICAvKipcbiAgICogSW1wb3J0cyBhIHByaXZhdGUga2V5IGFuZCBnZW5lcmF0ZXMgdGhlIGFwcHJvcHJpYXRlIHB1YmxpYyBrZXkuXG4gICAqXG4gICAqIEBwYXJhbSBwcml2ayBBIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IHJlcHJlc2VudGluZyB0aGUgcHJpdmF0ZSBrZXlcbiAgICpcbiAgICogQHJldHVybnMgdHJ1ZSBvbiBzdWNjZXNzLCBmYWxzZSBvbiBmYWlsdXJlXG4gICAqL1xuICBpbXBvcnRLZXkocHJpdms6IEJ1ZmZlcik6IGJvb2xlYW4ge1xuICAgIHRoaXMua2V5cGFpciA9IGVjLmtleUZyb21Qcml2YXRlKHByaXZrLnRvU3RyaW5nKFwiaGV4XCIpLCBcImhleFwiKVxuICAgIC8vIGRvaW5nIGhleCB0cmFuc2xhdGlvbiB0byBnZXQgQnVmZmVyIGNsYXNzXG4gICAgdHJ5IHtcbiAgICAgIHRoaXMucHJpdmsgPSBCdWZmZXIuZnJvbShcbiAgICAgICAgdGhpcy5rZXlwYWlyLmdldFByaXZhdGUoXCJoZXhcIikucGFkU3RhcnQoNjQsIFwiMFwiKSxcbiAgICAgICAgXCJoZXhcIlxuICAgICAgKVxuICAgICAgdGhpcy5wdWJrID0gQnVmZmVyLmZyb20oXG4gICAgICAgIHRoaXMua2V5cGFpci5nZXRQdWJsaWModHJ1ZSwgXCJoZXhcIikucGFkU3RhcnQoNjYsIFwiMFwiKSxcbiAgICAgICAgXCJoZXhcIlxuICAgICAgKVxuICAgICAgcmV0dXJuIHRydWUgLy8gc2lsbHkgSSBrbm93LCBidXQgdGhlIGludGVyZmFjZSByZXF1aXJlcyBzbyBpdCByZXR1cm5zIHRydWUgb24gc3VjY2Vzcywgc28gaWYgQnVmZmVyIGZhaWxzIHZhbGlkYXRpb24uLi5cbiAgICB9IGNhdGNoIChlcnJvcikge1xuICAgICAgcmV0dXJuIGZhbHNlXG4gICAgfVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgdGhlIGFkZHJlc3MgYXMgYSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfS5cbiAgICpcbiAgICogQHJldHVybnMgQSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSByZXByZXNlbnRhdGlvbiBvZiB0aGUgYWRkcmVzc1xuICAgKi9cbiAgZ2V0QWRkcmVzcygpOiBCdWZmZXIge1xuICAgIHJldHVybiBTRUNQMjU2azFLZXlQYWlyLmFkZHJlc3NGcm9tUHVibGljS2V5KHRoaXMucHViaylcbiAgfVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBhZGRyZXNzJ3Mgc3RyaW5nIHJlcHJlc2VudGF0aW9uLlxuICAgKlxuICAgKiBAcmV0dXJucyBBIHN0cmluZyByZXByZXNlbnRhdGlvbiBvZiB0aGUgYWRkcmVzc1xuICAgKi9cbiAgZ2V0QWRkcmVzc1N0cmluZygpOiBzdHJpbmcge1xuICAgIGNvbnN0IGFkZHI6IEJ1ZmZlciA9IFNFQ1AyNTZrMUtleVBhaXIuYWRkcmVzc0Zyb21QdWJsaWNLZXkodGhpcy5wdWJrKVxuICAgIGNvbnN0IHR5cGU6IFNlcmlhbGl6ZWRUeXBlID0gXCJiZWNoMzJcIlxuICAgIHJldHVybiBzZXJpYWxpemF0aW9uLmJ1ZmZlclRvVHlwZShhZGRyLCB0eXBlLCB0aGlzLmhycCwgdGhpcy5jaGFpbklEKVxuICB9XG5cbiAgLyoqXG4gICAqIFJldHVybnMgYW4gYWRkcmVzcyBnaXZlbiBhIHB1YmxpYyBrZXkuXG4gICAqXG4gICAqIEBwYXJhbSBwdWJrIEEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gcmVwcmVzZW50aW5nIHRoZSBwdWJsaWMga2V5XG4gICAqXG4gICAqIEByZXR1cm5zIEEge0BsaW5rIGh0dHBzOi8vZ2l0aHViLmNvbS9mZXJvc3MvYnVmZmVyfEJ1ZmZlcn0gZm9yIHRoZSBhZGRyZXNzIG9mIHRoZSBwdWJsaWMga2V5LlxuICAgKi9cbiAgc3RhdGljIGFkZHJlc3NGcm9tUHVibGljS2V5KHB1Yms6IEJ1ZmZlcik6IEJ1ZmZlciB7XG4gICAgaWYgKHB1YmsubGVuZ3RoID09PSA2NSkge1xuICAgICAgLyogaXN0YW5idWwgaWdub3JlIG5leHQgKi9cbiAgICAgIHB1YmsgPSBCdWZmZXIuZnJvbShcbiAgICAgICAgZWMua2V5RnJvbVB1YmxpYyhwdWJrKS5nZXRQdWJsaWModHJ1ZSwgXCJoZXhcIikucGFkU3RhcnQoNjYsIFwiMFwiKSxcbiAgICAgICAgXCJoZXhcIlxuICAgICAgKSAvLyBtYWtlIGNvbXBhY3QsIHN0aWNrIGJhY2sgaW50byBidWZmZXJcbiAgICB9XG4gICAgaWYgKHB1YmsubGVuZ3RoID09PSAzMykge1xuICAgICAgY29uc3Qgc2hhMjU2OiBCdWZmZXIgPSBCdWZmZXIuZnJvbShcbiAgICAgICAgY3JlYXRlSGFzaChcInNoYTI1NlwiKS51cGRhdGUocHViaykuZGlnZXN0KClcbiAgICAgIClcbiAgICAgIGNvbnN0IHJpcGVzaGE6IEJ1ZmZlciA9IEJ1ZmZlci5mcm9tKFxuICAgICAgICBjcmVhdGVIYXNoKFwicmlwZW1kMTYwXCIpLnVwZGF0ZShzaGEyNTYpLmRpZ2VzdCgpXG4gICAgICApXG4gICAgICByZXR1cm4gcmlwZXNoYVxuICAgIH1cbiAgICAvKiBpc3RhbmJ1bCBpZ25vcmUgbmV4dCAqL1xuICAgIHRocm93IG5ldyBQdWJsaWNLZXlFcnJvcihcIlVuYWJsZSB0byBtYWtlIGFkZHJlc3MuXCIpXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyBhIHN0cmluZyByZXByZXNlbnRhdGlvbiBvZiB0aGUgcHJpdmF0ZSBrZXkuXG4gICAqXG4gICAqIEByZXR1cm5zIEEgY2I1OCBzZXJpYWxpemVkIHN0cmluZyByZXByZXNlbnRhdGlvbiBvZiB0aGUgcHJpdmF0ZSBrZXlcbiAgICovXG4gIGdldFByaXZhdGVLZXlTdHJpbmcoKTogc3RyaW5nIHtcbiAgICByZXR1cm4gYFByaXZhdGVLZXktJHtiaW50b29scy5jYjU4RW5jb2RlKHRoaXMucHJpdmspfWBcbiAgfVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBwdWJsaWMga2V5LlxuICAgKlxuICAgKiBAcmV0dXJucyBBIGNiNTggc2VyaWFsaXplZCBzdHJpbmcgcmVwcmVzZW50YXRpb24gb2YgdGhlIHB1YmxpYyBrZXlcbiAgICovXG4gIGdldFB1YmxpY0tleVN0cmluZygpOiBzdHJpbmcge1xuICAgIHJldHVybiBiaW50b29scy5jYjU4RW5jb2RlKHRoaXMucHViaylcbiAgfVxuXG4gIC8qKlxuICAgKiBUYWtlcyBhIG1lc3NhZ2UsIHNpZ25zIGl0LCBhbmQgcmV0dXJucyB0aGUgc2lnbmF0dXJlLlxuICAgKlxuICAgKiBAcGFyYW0gbXNnIFRoZSBtZXNzYWdlIHRvIHNpZ24sIGJlIHN1cmUgdG8gaGFzaCBmaXJzdCBpZiBleHBlY3RlZFxuICAgKlxuICAgKiBAcmV0dXJucyBBIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IGNvbnRhaW5pbmcgdGhlIHNpZ25hdHVyZVxuICAgKi9cbiAgc2lnbihtc2c6IEJ1ZmZlcik6IEJ1ZmZlciB7XG4gICAgY29uc3Qgc2lnT2JqOiBlbGxpcHRpYy5lYy5TaWduYXR1cmUgPSB0aGlzLmtleXBhaXIuc2lnbihtc2csIHVuZGVmaW5lZCwge1xuICAgICAgY2Fub25pY2FsOiB0cnVlXG4gICAgfSlcbiAgICBjb25zdCByZWNvdmVyeTogQnVmZmVyID0gQnVmZmVyLmFsbG9jKDEpXG4gICAgcmVjb3Zlcnkud3JpdGVVSW50OChzaWdPYmoucmVjb3ZlcnlQYXJhbSwgMClcbiAgICBjb25zdCByOiBCdWZmZXIgPSBCdWZmZXIuZnJvbShzaWdPYmouci50b0FycmF5KFwiYmVcIiwgMzIpKSAvL3dlIGhhdmUgdG8gc2tpcCBuYXRpdmUgQnVmZmVyIGNsYXNzLCBzbyB0aGlzIGlzIHRoZSB3YXlcbiAgICBjb25zdCBzOiBCdWZmZXIgPSBCdWZmZXIuZnJvbShzaWdPYmoucy50b0FycmF5KFwiYmVcIiwgMzIpKSAvL3dlIGhhdmUgdG8gc2tpcCBuYXRpdmUgQnVmZmVyIGNsYXNzLCBzbyB0aGlzIGlzIHRoZSB3YXlcbiAgICBjb25zdCByZXN1bHQ6IEJ1ZmZlciA9IEJ1ZmZlci5jb25jYXQoW3IsIHMsIHJlY292ZXJ5XSwgNjUpXG4gICAgcmV0dXJuIHJlc3VsdFxuICB9XG5cbiAgLyoqXG4gICAqIFZlcmlmaWVzIHRoYXQgdGhlIHByaXZhdGUga2V5IGFzc29jaWF0ZWQgd2l0aCB0aGUgcHJvdmlkZWQgcHVibGljIGtleSBwcm9kdWNlcyB0aGUgc2lnbmF0dXJlIGFzc29jaWF0ZWQgd2l0aCB0aGUgZ2l2ZW4gbWVzc2FnZS5cbiAgICpcbiAgICogQHBhcmFtIG1zZyBUaGUgbWVzc2FnZSBhc3NvY2lhdGVkIHdpdGggdGhlIHNpZ25hdHVyZVxuICAgKiBAcGFyYW0gc2lnIFRoZSBzaWduYXR1cmUgb2YgdGhlIHNpZ25lZCBtZXNzYWdlXG4gICAqXG4gICAqIEByZXR1cm5zIFRydWUgb24gc3VjY2VzcywgZmFsc2Ugb24gZmFpbHVyZVxuICAgKi9cbiAgdmVyaWZ5KG1zZzogQnVmZmVyLCBzaWc6IEJ1ZmZlcik6IGJvb2xlYW4ge1xuICAgIGNvbnN0IHNpZ09iajogZWxsaXB0aWMuZWMuU2lnbmF0dXJlT3B0aW9ucyA9IHRoaXMuX3NpZ0Zyb21TaWdCdWZmZXIoc2lnKVxuICAgIHJldHVybiBlYy52ZXJpZnkobXNnLCBzaWdPYmosIHRoaXMua2V5cGFpcilcbiAgfVxuXG4gIC8qKlxuICAgKiBSZWNvdmVycyB0aGUgcHVibGljIGtleSBvZiBhIG1lc3NhZ2Ugc2lnbmVyIGZyb20gYSBtZXNzYWdlIGFuZCBpdHMgYXNzb2NpYXRlZCBzaWduYXR1cmUuXG4gICAqXG4gICAqIEBwYXJhbSBtc2cgVGhlIG1lc3NhZ2UgdGhhdCdzIHNpZ25lZFxuICAgKiBAcGFyYW0gc2lnIFRoZSBzaWduYXR1cmUgdGhhdCdzIHNpZ25lZCBvbiB0aGUgbWVzc2FnZVxuICAgKlxuICAgKiBAcmV0dXJucyBBIHtAbGluayBodHRwczovL2dpdGh1Yi5jb20vZmVyb3NzL2J1ZmZlcnxCdWZmZXJ9IGNvbnRhaW5pbmcgdGhlIHB1YmxpYyBrZXkgb2YgdGhlIHNpZ25lclxuICAgKi9cbiAgcmVjb3Zlcihtc2c6IEJ1ZmZlciwgc2lnOiBCdWZmZXIpOiBCdWZmZXIge1xuICAgIGNvbnN0IHNpZ09iajogZWxsaXB0aWMuZWMuU2lnbmF0dXJlT3B0aW9ucyA9IHRoaXMuX3NpZ0Zyb21TaWdCdWZmZXIoc2lnKVxuICAgIGNvbnN0IHB1YmsgPSBlYy5yZWNvdmVyUHViS2V5KG1zZywgc2lnT2JqLCBzaWdPYmoucmVjb3ZlcnlQYXJhbSlcbiAgICByZXR1cm4gQnVmZmVyLmZyb20ocHViay5lbmNvZGVDb21wcmVzc2VkKCkpXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyB0aGUgY2hhaW5JRCBhc3NvY2lhdGVkIHdpdGggdGhpcyBrZXkuXG4gICAqXG4gICAqIEByZXR1cm5zIFRoZSBbW0tleVBhaXJdXSdzIGNoYWluSURcbiAgICovXG4gIGdldENoYWluSUQoKTogc3RyaW5nIHtcbiAgICByZXR1cm4gdGhpcy5jaGFpbklEXG4gIH1cblxuICAvKipcbiAgICogU2V0cyB0aGUgdGhlIGNoYWluSUQgYXNzb2NpYXRlZCB3aXRoIHRoaXMga2V5LlxuICAgKlxuICAgKiBAcGFyYW0gY2hhaW5JRCBTdHJpbmcgZm9yIHRoZSBjaGFpbklEXG4gICAqL1xuICBzZXRDaGFpbklEKGNoYWluSUQ6IHN0cmluZyk6IHZvaWQge1xuICAgIHRoaXMuY2hhaW5JRCA9IGNoYWluSURcbiAgfVxuXG4gIC8qKlxuICAgKiBSZXR1cm5zIHRoZSBIdW1hbi1SZWFkYWJsZS1QYXJ0IG9mIHRoZSBuZXR3b3JrIGFzc29jaWF0ZWQgd2l0aCB0aGlzIGtleS5cbiAgICpcbiAgICogQHJldHVybnMgVGhlIFtbS2V5UGFpcl1dJ3MgSHVtYW4tUmVhZGFibGUtUGFydCBvZiB0aGUgbmV0d29yaydzIEJlY2gzMiBhZGRyZXNzaW5nIHNjaGVtZVxuICAgKi9cbiAgZ2V0SFJQKCk6IHN0cmluZyB7XG4gICAgcmV0dXJuIHRoaXMuaHJwXG4gIH1cblxuICAvKipcbiAgICogU2V0cyB0aGUgdGhlIEh1bWFuLVJlYWRhYmxlLVBhcnQgb2YgdGhlIG5ldHdvcmsgYXNzb2NpYXRlZCB3aXRoIHRoaXMga2V5LlxuICAgKlxuICAgKiBAcGFyYW0gaHJwIFN0cmluZyBmb3IgdGhlIEh1bWFuLVJlYWRhYmxlLVBhcnQgb2YgQmVjaDMyIGFkZHJlc3Nlc1xuICAgKi9cbiAgc2V0SFJQKGhycDogc3RyaW5nKTogdm9pZCB7XG4gICAgdGhpcy5ocnAgPSBocnBcbiAgfVxuXG4gIGNvbnN0cnVjdG9yKGhycDogc3RyaW5nLCBjaGFpbklEOiBzdHJpbmcpIHtcbiAgICBzdXBlcigpXG4gICAgdGhpcy5jaGFpbklEID0gY2hhaW5JRFxuICAgIHRoaXMuaHJwID0gaHJwXG4gICAgdGhpcy5nZW5lcmF0ZUtleSgpXG4gIH1cbn1cblxuLyoqXG4gKiBDbGFzcyBmb3IgcmVwcmVzZW50aW5nIGEga2V5IGNoYWluIGluIEF2YWxhbmNoZS5cbiAqXG4gKiBAdHlwZXBhcmFtIFNFQ1AyNTZrMUtleVBhaXIgQ2xhc3MgZXh0ZW5kaW5nIFtbU3RhbmRhcmRLZXlQYWlyXV0gd2hpY2ggaXMgdXNlZCBhcyB0aGUga2V5IGluIFtbU0VDUDI1NmsxS2V5Q2hhaW5dXVxuICovXG5leHBvcnQgYWJzdHJhY3QgY2xhc3MgU0VDUDI1NmsxS2V5Q2hhaW48XG4gIFNFQ1BLUENsYXNzIGV4dGVuZHMgU0VDUDI1NmsxS2V5UGFpclxuPiBleHRlbmRzIFN0YW5kYXJkS2V5Q2hhaW48U0VDUEtQQ2xhc3M+IHtcbiAgLyoqXG4gICAqIE1ha2VzIGEgbmV3IGtleSBwYWlyLCByZXR1cm5zIHRoZSBhZGRyZXNzLlxuICAgKlxuICAgKiBAcmV0dXJucyBBZGRyZXNzIG9mIHRoZSBuZXcga2V5IHBhaXJcbiAgICovXG4gIG1ha2VLZXk6ICgpID0+IFNFQ1BLUENsYXNzXG5cbiAgYWRkS2V5KG5ld0tleTogU0VDUEtQQ2xhc3MpOiB2b2lkIHtcbiAgICBzdXBlci5hZGRLZXkobmV3S2V5KVxuICB9XG5cbiAgLyoqXG4gICAqIEdpdmVuIGEgcHJpdmF0ZSBrZXksIG1ha2VzIGEgbmV3IGtleSBwYWlyLCByZXR1cm5zIHRoZSBhZGRyZXNzLlxuICAgKlxuICAgKiBAcGFyYW0gcHJpdmsgQSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSBvciBjYjU4IHNlcmlhbGl6ZWQgc3RyaW5nIHJlcHJlc2VudGluZyB0aGUgcHJpdmF0ZSBrZXlcbiAgICpcbiAgICogQHJldHVybnMgQWRkcmVzcyBvZiB0aGUgbmV3IGtleSBwYWlyXG4gICAqL1xuICBpbXBvcnRLZXk6IChwcml2azogQnVmZmVyIHwgc3RyaW5nKSA9PiBTRUNQS1BDbGFzc1xufVxuIl19