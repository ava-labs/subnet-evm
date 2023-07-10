"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.KeyChain = exports.KeyPair = void 0;
const bintools_1 = __importDefault(require("../../utils/bintools"));
const secp256k1_1 = require("../../common/secp256k1");
const utils_1 = require("../../utils");
/**
 * @ignore
 */
const bintools = bintools_1.default.getInstance();
const serialization = utils_1.Serialization.getInstance();
/**
 * Class for representing a private and public keypair on the Platform Chain.
 */
class KeyPair extends secp256k1_1.SECP256k1KeyPair {
    clone() {
        let newkp = new KeyPair(this.hrp, this.chainID);
        newkp.importKey(bintools.copyFrom(this.getPrivateKey()));
        return newkp;
    }
    create(...args) {
        if (args.length == 2) {
            return new KeyPair(args[0], args[1]);
        }
        return new KeyPair(this.hrp, this.chainID);
    }
}
exports.KeyPair = KeyPair;
/**
 * Class for representing a key chain in Avalanche.
 *
 * @typeparam KeyPair Class extending [[KeyPair]] which is used as the key in [[KeyChain]]
 */
class KeyChain extends secp256k1_1.SECP256k1KeyChain {
    /**
     * Returns instance of KeyChain.
     */
    constructor(hrp, chainID) {
        super();
        this.hrp = "";
        this.chainID = "";
        /**
         * Makes a new key pair, returns the address.
         *
         * @returns The new key pair
         */
        this.makeKey = () => {
            let keypair = new KeyPair(this.hrp, this.chainID);
            this.addKey(keypair);
            return keypair;
        };
        this.addKey = (newKey) => {
            newKey.setChainID(this.chainID);
            super.addKey(newKey);
        };
        /**
         * Given a private key, makes a new key pair, returns the address.
         *
         * @param privk A {@link https://github.com/feross/buffer|Buffer} or cb58 serialized string representing the private key
         *
         * @returns The new key pair
         */
        this.importKey = (privk) => {
            let keypair = new KeyPair(this.hrp, this.chainID);
            let pk;
            if (typeof privk === "string") {
                pk = bintools.cb58Decode(privk.split("-")[1]);
            }
            else {
                pk = bintools.copyFrom(privk);
            }
            keypair.importKey(pk);
            if (!(keypair.getAddress().toString("hex") in this.keys)) {
                this.addKey(keypair);
            }
            return keypair;
        };
        this.hrp = hrp;
        this.chainID = chainID;
    }
    create(...args) {
        if (args.length == 2) {
            return new KeyChain(args[0], args[1]);
        }
        return new KeyChain(this.hrp, this.chainID);
    }
    clone() {
        const newkc = new KeyChain(this.hrp, this.chainID);
        for (let k in this.keys) {
            newkc.addKey(this.keys[`${k}`].clone());
        }
        return newkc;
    }
    union(kc) {
        let newkc = kc.clone();
        for (let k in this.keys) {
            newkc.addKey(this.keys[`${k}`].clone());
        }
        return newkc;
    }
}
exports.KeyChain = KeyChain;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoia2V5Y2hhaW4uanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi9zcmMvYXBpcy9wbGF0Zm9ybXZtL2tleWNoYWluLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7OztBQUtBLG9FQUEyQztBQUMzQyxzREFBNEU7QUFDNUUsdUNBQTJEO0FBRTNEOztHQUVHO0FBQ0gsTUFBTSxRQUFRLEdBQWEsa0JBQVEsQ0FBQyxXQUFXLEVBQUUsQ0FBQTtBQUNqRCxNQUFNLGFBQWEsR0FBa0IscUJBQWEsQ0FBQyxXQUFXLEVBQUUsQ0FBQTtBQUVoRTs7R0FFRztBQUNILE1BQWEsT0FBUSxTQUFRLDRCQUFnQjtJQUMzQyxLQUFLO1FBQ0gsSUFBSSxLQUFLLEdBQVksSUFBSSxPQUFPLENBQUMsSUFBSSxDQUFDLEdBQUcsRUFBRSxJQUFJLENBQUMsT0FBTyxDQUFDLENBQUE7UUFDeEQsS0FBSyxDQUFDLFNBQVMsQ0FBQyxRQUFRLENBQUMsUUFBUSxDQUFDLElBQUksQ0FBQyxhQUFhLEVBQUUsQ0FBQyxDQUFDLENBQUE7UUFDeEQsT0FBTyxLQUFhLENBQUE7SUFDdEIsQ0FBQztJQUVELE1BQU0sQ0FBQyxHQUFHLElBQVc7UUFDbkIsSUFBSSxJQUFJLENBQUMsTUFBTSxJQUFJLENBQUMsRUFBRTtZQUNwQixPQUFPLElBQUksT0FBTyxDQUFDLElBQUksQ0FBQyxDQUFDLENBQUMsRUFBRSxJQUFJLENBQUMsQ0FBQyxDQUFDLENBQVMsQ0FBQTtTQUM3QztRQUNELE9BQU8sSUFBSSxPQUFPLENBQUMsSUFBSSxDQUFDLEdBQUcsRUFBRSxJQUFJLENBQUMsT0FBTyxDQUFTLENBQUE7SUFDcEQsQ0FBQztDQUNGO0FBYkQsMEJBYUM7QUFFRDs7OztHQUlHO0FBQ0gsTUFBYSxRQUFTLFNBQVEsNkJBQTBCO0lBaUV0RDs7T0FFRztJQUNILFlBQVksR0FBVyxFQUFFLE9BQWU7UUFDdEMsS0FBSyxFQUFFLENBQUE7UUFwRVQsUUFBRyxHQUFXLEVBQUUsQ0FBQTtR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IEBwYWNrYWdlRG9jdW1lbnRhdGlvblxuICogQG1vZHVsZSBBUEktUGxhdGZvcm1WTS1LZXlDaGFpblxuICovXG5pbXBvcnQgeyBCdWZmZXIgfSBmcm9tIFwiYnVmZmVyL1wiXG5pbXBvcnQgQmluVG9vbHMgZnJvbSBcIi4uLy4uL3V0aWxzL2JpbnRvb2xzXCJcbmltcG9ydCB7IFNFQ1AyNTZrMUtleUNoYWluLCBTRUNQMjU2azFLZXlQYWlyIH0gZnJvbSBcIi4uLy4uL2NvbW1vbi9zZWNwMjU2azFcIlxuaW1wb3J0IHsgU2VyaWFsaXphdGlvbiwgU2VyaWFsaXplZFR5cGUgfSBmcm9tIFwiLi4vLi4vdXRpbHNcIlxuXG4vKipcbiAqIEBpZ25vcmVcbiAqL1xuY29uc3QgYmludG9vbHM6IEJpblRvb2xzID0gQmluVG9vbHMuZ2V0SW5zdGFuY2UoKVxuY29uc3Qgc2VyaWFsaXphdGlvbjogU2VyaWFsaXphdGlvbiA9IFNlcmlhbGl6YXRpb24uZ2V0SW5zdGFuY2UoKVxuXG4vKipcbiAqIENsYXNzIGZvciByZXByZXNlbnRpbmcgYSBwcml2YXRlIGFuZCBwdWJsaWMga2V5cGFpciBvbiB0aGUgUGxhdGZvcm0gQ2hhaW4uXG4gKi9cbmV4cG9ydCBjbGFzcyBLZXlQYWlyIGV4dGVuZHMgU0VDUDI1NmsxS2V5UGFpciB7XG4gIGNsb25lKCk6IHRoaXMge1xuICAgIGxldCBuZXdrcDogS2V5UGFpciA9IG5ldyBLZXlQYWlyKHRoaXMuaHJwLCB0aGlzLmNoYWluSUQpXG4gICAgbmV3a3AuaW1wb3J0S2V5KGJpbnRvb2xzLmNvcHlGcm9tKHRoaXMuZ2V0UHJpdmF0ZUtleSgpKSlcbiAgICByZXR1cm4gbmV3a3AgYXMgdGhpc1xuICB9XG5cbiAgY3JlYXRlKC4uLmFyZ3M6IGFueVtdKTogdGhpcyB7XG4gICAgaWYgKGFyZ3MubGVuZ3RoID09IDIpIHtcbiAgICAgIHJldHVybiBuZXcgS2V5UGFpcihhcmdzWzBdLCBhcmdzWzFdKSBhcyB0aGlzXG4gICAgfVxuICAgIHJldHVybiBuZXcgS2V5UGFpcih0aGlzLmhycCwgdGhpcy5jaGFpbklEKSBhcyB0aGlzXG4gIH1cbn1cblxuLyoqXG4gKiBDbGFzcyBmb3IgcmVwcmVzZW50aW5nIGEga2V5IGNoYWluIGluIEF2YWxhbmNoZS5cbiAqXG4gKiBAdHlwZXBhcmFtIEtleVBhaXIgQ2xhc3MgZXh0ZW5kaW5nIFtbS2V5UGFpcl1dIHdoaWNoIGlzIHVzZWQgYXMgdGhlIGtleSBpbiBbW0tleUNoYWluXV1cbiAqL1xuZXhwb3J0IGNsYXNzIEtleUNoYWluIGV4dGVuZHMgU0VDUDI1NmsxS2V5Q2hhaW48S2V5UGFpcj4ge1xuICBocnA6IHN0cmluZyA9IFwiXCJcbiAgY2hhaW5JRDogc3RyaW5nID0gXCJcIlxuXG4gIC8qKlxuICAgKiBNYWtlcyBhIG5ldyBrZXkgcGFpciwgcmV0dXJucyB0aGUgYWRkcmVzcy5cbiAgICpcbiAgICogQHJldHVybnMgVGhlIG5ldyBrZXkgcGFpclxuICAgKi9cbiAgbWFrZUtleSA9ICgpOiBLZXlQYWlyID0+IHtcbiAgICBsZXQga2V5cGFpcjogS2V5UGFpciA9IG5ldyBLZXlQYWlyKHRoaXMuaHJwLCB0aGlzLmNoYWluSUQpXG4gICAgdGhpcy5hZGRLZXkoa2V5cGFpcilcbiAgICByZXR1cm4ga2V5cGFpclxuICB9XG5cbiAgYWRkS2V5ID0gKG5ld0tleTogS2V5UGFpcikgPT4ge1xuICAgIG5ld0tleS5zZXRDaGFpbklEKHRoaXMuY2hhaW5JRClcbiAgICBzdXBlci5hZGRLZXkobmV3S2V5KVxuICB9XG5cbiAgLyoqXG4gICAqIEdpdmVuIGEgcHJpdmF0ZSBrZXksIG1ha2VzIGEgbmV3IGtleSBwYWlyLCByZXR1cm5zIHRoZSBhZGRyZXNzLlxuICAgKlxuICAgKiBAcGFyYW0gcHJpdmsgQSB7QGxpbmsgaHR0cHM6Ly9naXRodWIuY29tL2Zlcm9zcy9idWZmZXJ8QnVmZmVyfSBvciBjYjU4IHNlcmlhbGl6ZWQgc3RyaW5nIHJlcHJlc2VudGluZyB0aGUgcHJpdmF0ZSBrZXlcbiAgICpcbiAgICogQHJldHVybnMgVGhlIG5ldyBrZXkgcGFpclxuICAgKi9cbiAgaW1wb3J0S2V5ID0gKHByaXZrOiBCdWZmZXIgfCBzdHJpbmcpOiBLZXlQYWlyID0+IHtcbiAgICBsZXQga2V5cGFpcjogS2V5UGFpciA9IG5ldyBLZXlQYWlyKHRoaXMuaHJwLCB0aGlzLmNoYWluSUQpXG4gICAgbGV0IHBrOiBCdWZmZXJcbiAgICBpZiAodHlwZW9mIHByaXZrID09PSBcInN0cmluZ1wiKSB7XG4gICAgICBwayA9IGJpbnRvb2xzLmNiNThEZWNvZGUocHJpdmsuc3BsaXQoXCItXCIpWzFdKVxuICAgIH0gZWxzZSB7XG4gICAgICBwayA9IGJpbnRvb2xzLmNvcHlGcm9tKHByaXZrKVxuICAgIH1cbiAgICBrZXlwYWlyLmltcG9ydEtleShwaylcbiAgICBpZiAoIShrZXlwYWlyLmdldEFkZHJlc3MoKS50b1N0cmluZyhcImhleFwiKSBpbiB0aGlzLmtleXMpKSB7XG4gICAgICB0aGlzLmFkZEtleShrZXlwYWlyKVxuICAgIH1cbiAgICByZXR1cm4ga2V5cGFpclxuICB9XG5cbiAgY3JlYXRlKC4uLmFyZ3M6IGFueVtdKTogdGhpcyB7XG4gICAgaWYgKGFyZ3MubGVuZ3RoID09IDIpIHtcbiAgICAgIHJldHVybiBuZXcgS2V5Q2hhaW4oYXJnc1swXSwgYXJnc1sxXSkgYXMgdGhpc1xuICAgIH1cbiAgICByZXR1cm4gbmV3IEtleUNoYWluKHRoaXMuaHJwLCB0aGlzLmNoYWluSUQpIGFzIHRoaXNcbiAgfVxuXG4gIGNsb25lKCk6IHRoaXMge1xuICAgIGNvbnN0IG5ld2tjOiBLZXlDaGFpbiA9IG5ldyBLZXlDaGFpbih0aGlzLmhycCwgdGhpcy5jaGFpbklEKVxuICAgIGZvciAobGV0IGsgaW4gdGhpcy5rZXlzKSB7XG4gICAgICBuZXdrYy5hZGRLZXkodGhpcy5rZXlzW2Ake2t9YF0uY2xvbmUoKSlcbiAgICB9XG4gICAgcmV0dXJuIG5ld2tjIGFzIHRoaXNcbiAgfVxuXG4gIHVuaW9uKGtjOiB0aGlzKTogdGhpcyB7XG4gICAgbGV0IG5ld2tjOiBLZXlDaGFpbiA9IGtjLmNsb25lKClcbiAgICBmb3IgKGxldCBrIGluIHRoaXMua2V5cykge1xuICAgICAgbmV3a2MuYWRkS2V5KHRoaXMua2V5c1tgJHtrfWBdLmNsb25lKCkpXG4gICAgfVxuICAgIHJldHVybiBuZXdrYyBhcyB0aGlzXG4gIH1cblxuICAvKipcbiAgICogUmV0dXJucyBpbnN0YW5jZSBvZiBLZXlDaGFpbi5cbiAgICovXG4gIGNvbnN0cnVjdG9yKGhycDogc3RyaW5nLCBjaGFpbklEOiBzdHJpbmcpIHtcbiAgICBzdXBlcigpXG4gICAgdGhpcy5ocnAgPSBocnBcbiAgICB0aGlzLmNoYWluSUQgPSBjaGFpbklEXG4gIH1cbn1cbiJdfQ==