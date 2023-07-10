"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const utils_1 = require("src/utils");
describe("Errors", () => {
    test("AvalancheError", () => {
        try {
            throw new utils_1.AvalancheError("Testing AvalancheError", "0");
        }
        catch (error) {
            expect(error.getCode()).toBe("0");
        }
        expect(() => {
            throw new utils_1.AvalancheError("Testing AvalancheError", "0");
        }).toThrow("Testing AvalancheError");
        expect(() => {
            throw new utils_1.AvalancheError("Testing AvalancheError", "0");
        }).toThrowError();
    });
    test("AddressError", () => {
        try {
            throw new utils_1.AddressError("Testing AddressError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1000");
        }
        expect(() => {
            throw new utils_1.AddressError("Testing AddressError");
        }).toThrow("Testing AddressError");
        expect(() => {
            throw new utils_1.AddressError("Testing AddressError");
        }).toThrowError();
    });
    test("GooseEggCheckError", () => {
        try {
            throw new utils_1.GooseEggCheckError("Testing GooseEggCheckError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1001");
        }
        expect(() => {
            throw new utils_1.GooseEggCheckError("Testing GooseEggCheckError");
        }).toThrow("Testing GooseEggCheckError");
        expect(() => {
            throw new utils_1.GooseEggCheckError("Testing GooseEggCheckError");
        }).toThrowError();
    });
    test("ChainIdError", () => {
        try {
            throw new utils_1.ChainIdError("Testing ChainIdError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1002");
        }
        expect(() => {
            throw new utils_1.ChainIdError("Testing ChainIdError");
        }).toThrow("Testing ChainIdError");
        expect(() => {
            throw new utils_1.ChainIdError("Testing ChainIdError");
        }).toThrowError();
    });
    test("NoAtomicUTXOsError", () => {
        try {
            throw new utils_1.NoAtomicUTXOsError("Testing NoAtomicUTXOsError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1003");
        }
        expect(() => {
            throw new utils_1.NoAtomicUTXOsError("Testing NoAtomicUTXOsError");
        }).toThrow("Testing NoAtomicUTXOsError");
        expect(() => {
            throw new utils_1.NoAtomicUTXOsError("Testing NoAtomicUTXOsError");
        }).toThrowError();
    });
    test("SymbolError", () => {
        try {
            throw new utils_1.SymbolError("Testing SymbolError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1004");
        }
        expect(() => {
            throw new utils_1.SymbolError("Testing SymbolError");
        }).toThrow("Testing SymbolError");
        expect(() => {
            throw new utils_1.SymbolError("Testing SymbolError");
        }).toThrowError();
    });
    test("NameError", () => {
        try {
            throw new utils_1.NameError("Testing NameError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1005");
        }
        expect(() => {
            throw new utils_1.NameError("Testing NameError");
        }).toThrow("Testing NameError");
        expect(() => {
            throw new utils_1.NameError("Testing NameError");
        }).toThrowError();
    });
    test("TransactionError", () => {
        try {
            throw new utils_1.TransactionError("Testing TransactionError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1006");
        }
        expect(() => {
            throw new utils_1.TransactionError("Testing TransactionError");
        }).toThrow("Testing TransactionError");
        expect(() => {
            throw new utils_1.TransactionError("Testing TransactionError");
        }).toThrowError();
    });
    test("CodecIdError", () => {
        try {
            throw new utils_1.CodecIdError("Testing CodecIdError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1007");
        }
        expect(() => {
            throw new utils_1.CodecIdError("Testing CodecIdError");
        }).toThrow("Testing CodecIdError");
        expect(() => {
            throw new utils_1.CodecIdError("Testing CodecIdError");
        }).toThrowError();
    });
    test("CredIdError", () => {
        try {
            throw new utils_1.CredIdError("Testing CredIdError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1008");
        }
        expect(() => {
            throw new utils_1.CredIdError("Testing CredIdError");
        }).toThrow("Testing CredIdError");
        expect(() => {
            throw new utils_1.CredIdError("Testing CredIdError");
        }).toThrowError();
    });
    test("TransferableOutputError", () => {
        try {
            throw new utils_1.TransferableOutputError("Testing TransferableOutputError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1009");
        }
        expect(() => {
            throw new utils_1.TransferableOutputError("Testing TransferableOutputError");
        }).toThrow("Testing TransferableOutputError");
        expect(() => {
            throw new utils_1.TransferableOutputError("Testing TransferableOutputError");
        }).toThrowError();
    });
    test("TransferableInputError", () => {
        try {
            throw new utils_1.TransferableInputError("Testing TransferableInputError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1010");
        }
        expect(() => {
            throw new utils_1.TransferableInputError("Testing TransferableInputError");
        }).toThrow("Testing TransferableInputError");
        expect(() => {
            throw new utils_1.TransferableInputError("Testing TransferableInputError");
        }).toThrowError();
    });
    test("InputIdError", () => {
        try {
            throw new utils_1.InputIdError("Testing InputIdError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1011");
        }
        expect(() => {
            throw new utils_1.InputIdError("Testing InputIdError");
        }).toThrow("Testing InputIdError");
        expect(() => {
            throw new utils_1.InputIdError("Testing InputIdError");
        }).toThrowError();
    });
    test("OperationError", () => {
        try {
            throw new utils_1.OperationError("Testing OperationError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1012");
        }
        expect(() => {
            throw new utils_1.OperationError("Testing OperationError");
        }).toThrow("Testing OperationError");
        expect(() => {
            throw new utils_1.OperationError("Testing OperationError");
        }).toThrowError();
    });
    test("InvalidOperationIdError", () => {
        try {
            throw new utils_1.InvalidOperationIdError("Testing InvalidOperationIdError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1013");
        }
        expect(() => {
            throw new utils_1.InvalidOperationIdError("Testing InvalidOperationIdError");
        }).toThrow("Testing InvalidOperationIdError");
        expect(() => {
            throw new utils_1.InvalidOperationIdError("Testing InvalidOperationIdError");
        }).toThrowError();
    });
    test("ChecksumError", () => {
        try {
            throw new utils_1.ChecksumError("Testing ChecksumError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1014");
        }
        expect(() => {
            throw new utils_1.ChecksumError("Testing ChecksumError");
        }).toThrow("Testing ChecksumError");
        expect(() => {
            throw new utils_1.ChecksumError("Testing ChecksumError");
        }).toThrowError();
    });
    test("OutputIdError", () => {
        try {
            throw new utils_1.OutputIdError("Testing OutputIdError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1015");
        }
        expect(() => {
            throw new utils_1.OutputIdError("Testing OutputIdError");
        }).toThrow("Testing OutputIdError");
        expect(() => {
            throw new utils_1.OutputIdError("Testing OutputIdError");
        }).toThrowError();
    });
    test("UTXOError", () => {
        try {
            throw new utils_1.UTXOError("Testing UTXOError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1016");
        }
        expect(() => {
            throw new utils_1.UTXOError("Testing UTXOError");
        }).toThrow("Testing UTXOError");
        expect(() => {
            throw new utils_1.UTXOError("Testing UTXOError");
        }).toThrowError();
    });
    test("InsufficientFundsError", () => {
        try {
            throw new utils_1.InsufficientFundsError("Testing InsufficientFundsError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1017");
        }
        expect(() => {
            throw new utils_1.InsufficientFundsError("Testing InsufficientFundsError");
        }).toThrow("Testing InsufficientFundsError");
        expect(() => {
            throw new utils_1.InsufficientFundsError("Testing InsufficientFundsError");
        }).toThrowError();
    });
    test("ThresholdError", () => {
        try {
            throw new utils_1.ThresholdError("Testing ThresholdError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1018");
        }
        expect(() => {
            throw new utils_1.ThresholdError("Testing ThresholdError");
        }).toThrow("Testing ThresholdError");
        expect(() => {
            throw new utils_1.ThresholdError("Testing ThresholdError");
        }).toThrowError();
    });
    test("SECPMintOutputError", () => {
        try {
            throw new utils_1.SECPMintOutputError("Testing SECPMintOutputError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1019");
        }
        expect(() => {
            throw new utils_1.SECPMintOutputError("Testing SECPMintOutputError");
        }).toThrow("Testing SECPMintOutputError");
        expect(() => {
            throw new utils_1.SECPMintOutputError("Testing SECPMintOutputError");
        }).toThrowError();
    });
    test("EVMInputError", () => {
        try {
            throw new utils_1.EVMInputError("Testing EVMInputError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1020");
        }
        expect(() => {
            throw new utils_1.EVMInputError("Testing EVMInputError");
        }).toThrow("Testing EVMInputError");
        expect(() => {
            throw new utils_1.EVMInputError("Testing EVMInputError");
        }).toThrowError();
    });
    test("EVMOutputError", () => {
        try {
            throw new utils_1.EVMOutputError("Testing EVMOutputError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1021");
        }
        expect(() => {
            throw new utils_1.EVMOutputError("Testing EVMOutputError");
        }).toThrow("Testing EVMOutputError");
        expect(() => {
            throw new utils_1.EVMOutputError("Testing EVMOutputError");
        }).toThrowError();
    });
    test("FeeAssetError", () => {
        try {
            throw new utils_1.FeeAssetError("Testing FeeAssetError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1022");
        }
        expect(() => {
            throw new utils_1.FeeAssetError("Testing FeeAssetError");
        }).toThrow("Testing FeeAssetError");
        expect(() => {
            throw new utils_1.FeeAssetError("Testing FeeAssetError");
        }).toThrowError();
    });
    test("StakeError", () => {
        try {
            throw new utils_1.StakeError("Testing StakeError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1023");
        }
        expect(() => {
            throw new utils_1.StakeError("Testing StakeError");
        }).toThrow("Testing StakeError");
        expect(() => {
            throw new utils_1.StakeError("Testing StakeError");
        }).toThrowError();
    });
    test("TimeError", () => {
        try {
            throw new utils_1.TimeError("Testing TimeError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1024");
        }
        expect(() => {
            throw new utils_1.TimeError("Testing TimeError");
        }).toThrow("Testing TimeError");
        expect(() => {
            throw new utils_1.TimeError("Testing TimeError");
        }).toThrowError();
    });
    test("DelegationFeeError", () => {
        try {
            throw new utils_1.DelegationFeeError("Testing DelegationFeeError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1025");
        }
        expect(() => {
            throw new utils_1.DelegationFeeError("Testing DelegationFeeError");
        }).toThrow("Testing DelegationFeeError");
        expect(() => {
            throw new utils_1.DelegationFeeError("Testing DelegationFeeError");
        }).toThrowError();
    });
    test("SubnetOwnerError", () => {
        try {
            throw new utils_1.SubnetOwnerError("Testing SubnetOwnerError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1026");
        }
        expect(() => {
            throw new utils_1.SubnetOwnerError("Testing SubnetOwnerError");
        }).toThrow("Testing SubnetOwnerError");
        expect(() => {
            throw new utils_1.SubnetOwnerError("Testing SubnetOwnerError");
        }).toThrowError();
    });
    test("BufferSizeError", () => {
        try {
            throw new utils_1.BufferSizeError("Testing BufferSizeError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1027");
        }
        expect(() => {
            throw new utils_1.BufferSizeError("Testing BufferSizeError");
        }).toThrow("Testing BufferSizeError");
        expect(() => {
            throw new utils_1.BufferSizeError("Testing BufferSizeError");
        }).toThrowError();
    });
    test("AddressIndexError", () => {
        try {
            throw new utils_1.AddressIndexError("Testing AddressIndexError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1028");
        }
        expect(() => {
            throw new utils_1.AddressIndexError("Testing AddressIndexError");
        }).toThrow("Testing AddressIndexError");
        expect(() => {
            throw new utils_1.AddressIndexError("Testing AddressIndexError");
        }).toThrowError();
    });
    test("PublicKeyError", () => {
        try {
            throw new utils_1.PublicKeyError("Testing PublicKeyError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1029");
        }
        expect(() => {
            throw new utils_1.PublicKeyError("Testing PublicKeyError");
        }).toThrow("Testing PublicKeyError");
        expect(() => {
            throw new utils_1.PublicKeyError("Testing PublicKeyError");
        }).toThrowError();
    });
    test("MergeRuleError", () => {
        try {
            throw new utils_1.MergeRuleError("Testing MergeRuleError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1030");
        }
        expect(() => {
            throw new utils_1.MergeRuleError("Testing MergeRuleError");
        }).toThrow("Testing MergeRuleError");
        expect(() => {
            throw new utils_1.MergeRuleError("Testing MergeRuleError");
        }).toThrowError();
    });
    test("Base58Error", () => {
        try {
            throw new utils_1.Base58Error("Testing Base58Error");
        }
        catch (error) {
            expect(error.getCode()).toBe("1031");
        }
        expect(() => {
            throw new utils_1.Base58Error("Testing Base58Error");
        }).toThrow("Testing Base58Error");
        expect(() => {
            throw new utils_1.Base58Error("Testing Base58Error");
        }).toThrowError();
    });
    test("PrivateKeyError", () => {
        try {
            throw new utils_1.PrivateKeyError("Testing PrivateKeyError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1032");
        }
        expect(() => {
            throw new utils_1.PrivateKeyError("Testing PrivateKeyError");
        }).toThrow("Testing PrivateKeyError");
        expect(() => {
            throw new utils_1.PrivateKeyError("Testing PrivateKeyError");
        }).toThrowError();
    });
    test("NodeIdError", () => {
        try {
            throw new utils_1.NodeIdError("Testing NodeIdError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1033");
        }
        expect(() => {
            throw new utils_1.NodeIdError("Testing NodeIdError");
        }).toThrow("Testing NodeIdError");
        expect(() => {
            throw new utils_1.NodeIdError("Testing NodeIdError");
        }).toThrowError();
    });
    test("HexError", () => {
        try {
            throw new utils_1.HexError("Testing HexError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1034");
        }
        expect(() => {
            throw new utils_1.HexError("Testing HexError");
        }).toThrow("Testing HexError");
        expect(() => {
            throw new utils_1.HexError("Testing HexError");
        }).toThrowError();
    });
    test("TypeIdError", () => {
        try {
            throw new utils_1.TypeIdError("Testing TypeIdError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1035");
        }
        expect(() => {
            throw new utils_1.TypeIdError("Testing TypeIdError");
        }).toThrow("Testing TypeIdError");
        expect(() => {
            throw new utils_1.TypeIdError("Testing TypeIdError");
        }).toThrowError();
    });
    test("TypeNameError", () => {
        try {
            throw new utils_1.TypeNameError("Testing TypeNameError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1042");
        }
        expect(() => {
            throw new utils_1.TypeNameError("Testing TypeNameError");
        }).toThrow("Testing TypeNameError");
        expect(() => {
            throw new utils_1.TypeNameError("Testing TypeNameError");
        }).toThrowError();
    });
    test("UnknownTypeError", () => {
        try {
            throw new utils_1.UnknownTypeError("Testing UnknownTypeError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1036");
        }
        expect(() => {
            throw new utils_1.UnknownTypeError("Testing UnknownTypeError");
        }).toThrow("Testing UnknownTypeError");
        expect(() => {
            throw new utils_1.UnknownTypeError("Testing UnknownTypeError");
        }).toThrowError();
    });
    test("Bech32Error", () => {
        try {
            throw new utils_1.Bech32Error("Testing Bech32Error");
        }
        catch (error) {
            expect(error.getCode()).toBe("1037");
        }
        expect(() => {
            throw new utils_1.Bech32Error("Testing Bech32Error");
        }).toThrow("Testing Bech32Error");
        expect(() => {
            throw new utils_1.Bech32Error("Testing Bech32Error");
        }).toThrowError();
    });
    test("EVMFeeError", () => {
        try {
            throw new utils_1.EVMFeeError("Testing EVMFeeError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1038");
        }
        expect(() => {
            throw new utils_1.EVMFeeError("Testing EVMFeeError");
        }).toThrow("Testing EVMFeeError");
        expect(() => {
            throw new utils_1.EVMFeeError("Testing EVMFeeError");
        }).toThrowError();
    });
    test("InvalidEntropy", () => {
        try {
            throw new utils_1.InvalidEntropy("Testing InvalidEntropy");
        }
        catch (error) {
            expect(error.getCode()).toBe("1039");
        }
        expect(() => {
            throw new utils_1.InvalidEntropy("Testing InvalidEntropy");
        }).toThrow("Testing InvalidEntropy");
        expect(() => {
            throw new utils_1.InvalidEntropy("Testing InvalidEntropy");
        }).toThrowError();
    });
    test("ProtocolError", () => {
        try {
            throw new utils_1.ProtocolError("Testing ProtocolError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1040");
        }
        expect(() => {
            throw new utils_1.ProtocolError("Testing ProtocolError");
        }).toThrow("Testing ProtocolError");
        expect(() => {
            throw new utils_1.ProtocolError("Testing ProtocolError");
        }).toThrowError();
    });
    test("SubnetIdError", () => {
        try {
            throw new utils_1.SubnetIdError("Testing SubnetIdError");
        }
        catch (error) {
            expect(error.getCode()).toBe("1041");
        }
        expect(() => {
            throw new utils_1.SubnetIdError("Testing SubnetIdError");
        }).toThrow("Testing SubnetIdError");
        expect(() => {
            throw new utils_1.SubnetIdError("Testing SubnetIdError");
        }).toThrowError();
    });
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZXJyb3JzLnRlc3QuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi90ZXN0cy91dGlscy9lcnJvcnMudGVzdC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOztBQUNBLHFDQTZDa0I7QUFFbEIsUUFBUSxDQUFDLFFBQVEsRUFBRSxHQUFTLEVBQUU7SUFDNUIsSUFBSSxDQUFDLGdCQUFnQixFQUFFLEdBQVMsRUFBRTtRQUNoQyxJQUFJO1lBQ0YsTUFBTSxJQUFJLHNCQUFjLENBQUMsd0JBQXdCLEVBQUUsR0FBRyxDQUFDLENBQUE7U0FDeEQ7UUFBQyxPQUFPLEtBQVUsRUFBRTtZQUNuQixNQUFNLENBQUMsS0FBSyxDQUFDLE9BQU8sRUFBRSxDQUFDLENBQUMsSUFBSSxDQUFDLEdBQUcsQ0FBQyxDQUFBO1NBQ2xDO1FBQ0QsTUFBTSxDQUFDLEdBQVMsRUFBRTtZQUNoQixNQUFNLElBQUksc0JBQWMsQ0FBQyx3QkFBd0IsRUFBRSxHQUFHLENBQUMsQ0FBQTtRQUN6RCxDQUFDLENBQUMsQ0FBQyxPQUFPLENBQUMsd0JBQXdCLENBQUMsQ0FBQTtRQUNwQyxNQUFNLENBQUMsR0FBUyxFQUFFO1lBQ2hCLE1BQU0sSUFBSSxzQkFBYyxDQUFDLHdCQUF3QixFQUFFLEdBQUcsQ0FBQyxDQUFBO1FBQ3pELENBQUMsQ0FBQyxDQUFDLFlBQVksRUFBRSxDQUFBO0lBQ25CLENBQUMsQ0FBQyxDQUFBO0lBRUYsSUFBSSxDQUFDLGNBQWMsRUFBRSxHQUFTLEVBQUU7UUFDOUIsSUFBSTtZQUNGLE1BQU0sSUFBSSxvQkFBWSxDQUFDLHNCQUFzQixDQUFDLENBQUE7U0FDL0M7UUFBQyxPQUFPLEtBQVUsRUFBRTtZQUNuQixNQUFNLENBQUMsS0FBSyxDQUFDLE9BQU8sRUFBRSxDQUFDLENBQUMsSUFBSSxDQUFDLE1BQU0sQ0FBQyxDQUFBO1NBQ3JDO1FBQ0QsTUFBTSxDQUFDLEdBQVMsRUFBRTtZQUNoQixNQUFNLElBQUksb0JBQVksQ0FBQyxzQkFBc0IsQ0FB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QyxDQUFBO1FBQ25DLE1BQU0sQ0FBQyxHQUFTLEVBQUU7WUFDaEIsTUFBTSxJQUFJLHFCQUFhLENBQUMsdUJBQXVCLENBQUMsQ0FBQTtRQUNsRCxDQUFDLENBQUMsQ0FBQyxZQUFZLEVBQUUsQ0FBQTtJQUNuQixDQUFDLENBQUMsQ0FBQTtJQUVGLElBQUksQ0FBQyxlQUFlLEVBQUUsR0FBUyxFQUFFO1FBQy9CLElBQUk7WUFDRixNQUFNLElBQUkscUJBQWEsQ0FBQyx1QkFBdUIsQ0FBQyxDQUFBO1NBQ2pEO1FBQUMsT0FBTyxLQUFVLEVBQUU7WUFDbkIsTUFBTSxDQUFDLEtBQUssQ0FBQyxPQUFPLEVBQUUsQ0FBQyxDQUFDLElBQUksQ0FBQyxNQUFNLENBQUMsQ0FBQTtTQUNyQztRQUNELE1BQU0sQ0FBQyxHQUFTLEVBQUU7WUFDaEIsTUFBTSxJQUFJLHFCQUFhLENBQUMsdUJBQXVCLENBQUMsQ0FBQTtRQUNsRCxDQUFDLENBQUMsQ0FBQyxPQUFPLENBQUMsdUJBQXVCLENBQUMsQ0FBQTtRQUNuQyxNQUFNLENBQUMsR0FBUyxFQUFFO1lBQ2hCLE1BQU0sSUFBSSxxQkFBYSxDQUFDLHVCQUF1QixDQUFDLENBQUE7UUFDbEQsQ0FBQyxDQUFDLENBQUMsWUFBWSxFQUFFLENBQUE7SUFDbkIsQ0FBQyxDQUFDLENBQUE7QUFDSixDQUFDLENBQUMsQ0FBQSIsInNvdXJjZXNDb250ZW50IjpbImltcG9ydCBCaW5Ub29scyBmcm9tIFwiLi4vLi4vc3JjL3V0aWxzL2JpbnRvb2xzXCJcbmltcG9ydCB7XG4gIEF2YWxhbmNoZUVycm9yLFxuICBBZGRyZXNzRXJyb3IsXG4gIEdvb3NlRWdnQ2hlY2tFcnJvcixcbiAgQ2hhaW5JZEVycm9yLFxuICBOb0F0b21pY1VUWE9zRXJyb3IsXG4gIFN5bWJvbEVycm9yLFxuICBOYW1lRXJyb3IsXG4gIFRyYW5zYWN0aW9uRXJyb3IsXG4gIENvZGVjSWRFcnJvcixcbiAgQ3JlZElkRXJyb3IsXG4gIFRyYW5zZmVyYWJsZU91dHB1dEVycm9yLFxuICBUcmFuc2ZlcmFibGVJbnB1dEVycm9yLFxuICBJbnB1dElkRXJyb3IsXG4gIE9wZXJhdGlvbkVycm9yLFxuICBJbnZhbGlkT3BlcmF0aW9uSWRFcnJvcixcbiAgQ2hlY2tzdW1FcnJvcixcbiAgT3V0cHV0SWRFcnJvcixcbiAgVVRYT0Vycm9yLFxuICBJbnN1ZmZpY2llbnRGdW5kc0Vycm9yLFxuICBUaHJlc2hvbGRFcnJvcixcbiAgU0VDUE1pbnRPdXRwdXRFcnJvcixcbiAgRVZNSW5wdXRFcnJvcixcbiAgRVZNT3V0cHV0RXJyb3IsXG4gIEZlZUFzc2V0RXJyb3IsXG4gIFN0YWtlRXJyb3IsXG4gIFRpbWVFcnJvcixcbiAgRGVsZWdhdGlvbkZlZUVycm9yLFxuICBTdWJuZXRPd25lckVycm9yLFxuICBCdWZmZXJTaXplRXJyb3IsXG4gIEFkZHJlc3NJbmRleEVycm9yLFxuICBQdWJsaWNLZXlFcnJvcixcbiAgTWVyZ2VSdWxlRXJyb3IsXG4gIEJhc2U1OEVycm9yLFxuICBQcml2YXRlS2V5RXJyb3IsXG4gIE5vZGVJZEVycm9yLFxuICBIZXhFcnJvcixcbiAgVHlwZUlkRXJyb3IsXG4gIFR5cGVOYW1lRXJyb3IsXG4gIFVua25vd25UeXBlRXJyb3IsXG4gIEJlY2gzMkVycm9yLFxuICBFVk1GZWVFcnJvcixcbiAgSW52YWxpZEVudHJvcHksXG4gIFByb3RvY29sRXJyb3IsXG4gIFN1Ym5ldElkRXJyb3Jcbn0gZnJvbSBcInNyYy91dGlsc1wiXG5cbmRlc2NyaWJlKFwiRXJyb3JzXCIsICgpOiB2b2lkID0+IHtcbiAgdGVzdChcIkF2YWxhbmNoZUVycm9yXCIsICgpOiB2b2lkID0+IHtcbiAgICB0cnkge1xuICAgICAgdGhyb3cgbmV3IEF2YWxhbmNoZUVycm9yKFwiVGVzdGluZyBBdmFsYW5jaGVFcnJvclwiLCBcIjBcIilcbiAgICB9IGNhdGNoIChlcnJvcjogYW55KSB7XG4gICAgICBleHBlY3QoZXJyb3IuZ2V0Q29kZSgpKS50b0JlKFwiMFwiKVxuICAgIH1cbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IEF2YWxhbmNoZUVycm9yKFwiVGVzdGluZyBBdmFsYW5jaGVFcnJvclwiLCBcIjBcIilcbiAgICB9KS50b1Rocm93KFwiVGVzdGluZyBBdmFsYW5jaGVFcnJvclwiKVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgQXZhbGFuY2hlRXJyb3IoXCJUZXN0aW5nIEF2YWxhbmNoZUVycm9yXCIsIFwiMFwiKVxuICAgIH0pLnRvVGhyb3dFcnJvcigpXG4gIH0pXG5cbiAgdGVzdChcIkFkZHJlc3NFcnJvclwiLCAoKTogdm9pZCA9PiB7XG4gICAgdHJ5IHtcbiAgICAgIHRocm93IG5ldyBBZGRyZXNzRXJyb3IoXCJUZXN0aW5nIEFkZHJlc3NFcnJvclwiKVxuICAgIH0gY2F0Y2ggKGVycm9yOiBhbnkpIHtcbiAgICAgIGV4cGVjdChlcnJvci5nZXRDb2RlKCkpLnRvQmUoXCIxMDAwXCIpXG4gICAgfVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgQWRkcmVzc0Vycm9yKFwiVGVzdGluZyBBZGRyZXNzRXJyb3JcIilcbiAgICB9KS50b1Rocm93KFwiVGVzdGluZyBBZGRyZXNzRXJyb3JcIilcbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IEFkZHJlc3NFcnJvcihcIlRlc3RpbmcgQWRkcmVzc0Vycm9yXCIpXG4gICAgfSkudG9UaHJvd0Vycm9yKClcbiAgfSlcblxuICB0ZXN0KFwiR29vc2VFZ2dDaGVja0Vycm9yXCIsICgpOiB2b2lkID0+IHtcbiAgICB0cnkge1xuICAgICAgdGhyb3cgbmV3IEdvb3NlRWdnQ2hlY2tFcnJvcihcIlRlc3RpbmcgR29vc2VFZ2dDaGVja0Vycm9yXCIpXG4gICAgfSBjYXRjaCAoZXJyb3I6IGFueSkge1xuICAgICAgZXhwZWN0KGVycm9yLmdldENvZGUoKSkudG9CZShcIjEwMDFcIilcbiAgICB9XG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBHb29zZUVnZ0NoZWNrRXJyb3IoXCJUZXN0aW5nIEdvb3NlRWdnQ2hlY2tFcnJvclwiKVxuICAgIH0pLnRvVGhyb3coXCJUZXN0aW5nIEdvb3NlRWdnQ2hlY2tFcnJvclwiKVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgR29vc2VFZ2dDaGVja0Vycm9yKFwiVGVzdGluZyBHb29zZUVnZ0NoZWNrRXJyb3JcIilcbiAgICB9KS50b1Rocm93RXJyb3IoKVxuICB9KVxuXG4gIHRlc3QoXCJDaGFpbklkRXJyb3JcIiwgKCk6IHZvaWQgPT4ge1xuICAgIHRyeSB7XG4gICAgICB0aHJvdyBuZXcgQ2hhaW5JZEVycm9yKFwiVGVzdGluZyBDaGFpbklkRXJyb3JcIilcbiAgICB9IGNhdGNoIChlcnJvcjogYW55KSB7XG4gICAgICBleHBlY3QoZXJyb3IuZ2V0Q29kZSgpKS50b0JlKFwiMTAwMlwiKVxuICAgIH1cbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IENoYWluSWRFcnJvcihcIlRlc3RpbmcgQ2hhaW5JZEVycm9yXCIpXG4gICAgfSkudG9UaHJvdyhcIlRlc3RpbmcgQ2hhaW5JZEVycm9yXCIpXG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBDaGFpbklkRXJyb3IoXCJUZXN0aW5nIENoYWluSWRFcnJvclwiKVxuICAgIH0pLnRvVGhyb3dFcnJvcigpXG4gIH0pXG5cbiAgdGVzdChcIk5vQXRvbWljVVRYT3NFcnJvclwiLCAoKTogdm9pZCA9PiB7XG4gICAgdHJ5IHtcbiAgICAgIHRocm93IG5ldyBOb0F0b21pY1VUWE9zRXJyb3IoXCJUZXN0aW5nIE5vQXRvbWljVVRYT3NFcnJvclwiKVxuICAgIH0gY2F0Y2ggKGVycm9yOiBhbnkpIHtcbiAgICAgIGV4cGVjdChlcnJvci5nZXRDb2RlKCkpLnRvQmUoXCIxMDAzXCIpXG4gICAgfVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgTm9BdG9taWNVVFhPc0Vycm9yKFwiVGVzdGluZyBOb0F0b21pY1VUWE9zRXJyb3JcIilcbiAgICB9KS50b1Rocm93KFwiVGVzdGluZyBOb0F0b21pY1VUWE9zRXJyb3JcIilcbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IE5vQXRvbWljVVRYT3NFcnJvcihcIlRlc3RpbmcgTm9BdG9taWNVVFhPc0Vycm9yXCIpXG4gICAgfSkudG9UaHJvd0Vycm9yKClcbiAgfSlcblxuICB0ZXN0KFwiU3ltYm9sRXJyb3JcIiwgKCk6IHZvaWQgPT4ge1xuICAgIHRyeSB7XG4gICAgICB0aHJvdyBuZXcgU3ltYm9sRXJyb3IoXCJUZXN0aW5nIFN5bWJvbEVycm9yXCIpXG4gICAgfSBjYXRjaCAoZXJyb3I6IGFueSkge1xuICAgICAgZXhwZWN0KGVycm9yLmdldENvZGUoKSkudG9CZShcIjEwMDRcIilcbiAgICB9XG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBTeW1ib2xFcnJvcihcIlRlc3RpbmcgU3ltYm9sRXJyb3JcIilcbiAgICB9KS50b1Rocm93KFwiVGVzdGluZyBTeW1ib2xFcnJvclwiKVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgU3ltYm9sRXJyb3IoXCJUZXN0aW5nIFN5bWJvbEVycm9yXCIpXG4gICAgfSkudG9UaHJvd0Vycm9yKClcbiAgfSlcblxuICB0ZXN0KFwiTmFtZUVycm9yXCIsICgpOiB2b2lkID0+IHtcbiAgICB0cnkge1xuICAgICAgdGhyb3cgbmV3IE5hbWVFcnJvcihcIlRlc3RpbmcgTmFtZUVycm9yXCIpXG4gICAgfSBjYXRjaCAoZXJyb3I6IGFueSkge1xuICAgICAgZXhwZWN0KGVycm9yLmdldENvZGUoKSkudG9CZShcIjEwMDVcIilcbiAgICB9XG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBOYW1lRXJyb3IoXCJUZXN0aW5nIE5hbWVFcnJvclwiKVxuICAgIH0pLnRvVGhyb3coXCJUZXN0aW5nIE5hbWVFcnJvclwiKVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgTmFtZUVycm9yKFwiVGVzdGluZyBOYW1lRXJyb3JcIilcbiAgICB9KS50b1Rocm93RXJyb3IoKVxuICB9KVxuXG4gIHRlc3QoXCJUcmFuc2FjdGlvbkVycm9yXCIsICgpOiB2b2lkID0+IHtcbiAgICB0cnkge1xuICAgICAgdGhyb3cgbmV3IFRyYW5zYWN0aW9uRXJyb3IoXCJUZXN0aW5nIFRyYW5zYWN0aW9uRXJyb3JcIilcbiAgICB9IGNhdGNoIChlcnJvcjogYW55KSB7XG4gICAgICBleHBlY3QoZXJyb3IuZ2V0Q29kZSgpKS50b0JlKFwiMTAwNlwiKVxuICAgIH1cbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IFRyYW5zYWN0aW9uRXJyb3IoXCJUZXN0aW5nIFRyYW5zYWN0aW9uRXJyb3JcIilcbiAgICB9KS50b1Rocm93KFwiVGVzdGluZyBUcmFuc2FjdGlvbkVycm9yXCIpXG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBUcmFuc2FjdGlvbkVycm9yKFwiVGVzdGluZyBUcmFuc2FjdGlvbkVycm9yXCIpXG4gICAgfSkudG9UaHJvd0Vycm9yKClcbiAgfSlcblxuICB0ZXN0KFwiQ29kZWNJZEVycm9yXCIsICgpOiB2b2lkID0+IHtcbiAgICB0cnkge1xuICAgICAgdGhyb3cgbmV3IENvZGVjSWRFcnJvcihcIlRlc3RpbmcgQ29kZWNJZEVycm9yXCIpXG4gICAgfSBjYXRjaCAoZXJyb3I6IGFueSkge1xuICAgICAgZXhwZWN0KGVycm9yLmdldENvZGUoKSkudG9CZShcIjEwMDdcIilcbiAgICB9XG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBDb2RlY0lkRXJyb3IoXCJUZXN0aW5nIENvZGVjSWRFcnJvclwiKVxuICAgIH0pLnRvVGhyb3coXCJUZXN0aW5nIENvZGVjSWRFcnJvclwiKVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgQ29kZWNJZEVycm9yKFwiVGVzdGluZyBDb2RlY0lkRXJyb3JcIilcbiAgICB9KS50b1Rocm93RXJyb3IoKVxuICB9KVxuXG4gIHRlc3QoXCJDcmVkSWRFcnJvclwiLCAoKTogdm9pZCA9PiB7XG4gICAgdHJ5IHtcbiAgICAgIHRocm93IG5ldyBDcmVkSWRFcnJvcihcIlRlc3RpbmcgQ3JlZElkRXJyb3JcIilcbiAgICB9IGNhdGNoIChlcnJvcjogYW55KSB7XG4gICAgICBleHBlY3QoZXJyb3IuZ2V0Q29kZSgpKS50b0JlKFwiMTAwOFwiKVxuICAgIH1cbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IENyZWRJZEVycm9yKFwiVGVzdGluZyBDcmVkSWRFcnJvclwiKVxuICAgIH0pLnRvVGhyb3coXCJUZXN0aW5nIENyZWRJZEVycm9yXCIpXG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBDcmVkSWRFcnJvcihcIlRlc3RpbmcgQ3JlZElkRXJyb3JcIilcbiAgICB9KS50b1Rocm93RXJyb3IoKVxuICB9KVxuXG4gIHRlc3QoXCJUcmFuc2ZlcmFibGVPdXRwdXRFcnJvclwiLCAoKTogdm9pZCA9PiB7XG4gICAgdHJ5IHtcbiAgICAgIHRocm93IG5ldyBUcmFuc2ZlcmFibGVPdXRwdXRFcnJvcihcIlRlc3RpbmcgVHJhbnNmZXJhYmxlT3V0cHV0RXJyb3JcIilcbiAgICB9IGNhdGNoIChlcnJvcjogYW55KSB7XG4gICAgICBleHBlY3QoZXJyb3IuZ2V0Q29kZSgpKS50b0JlKFwiMTAwOVwiKVxuICAgIH1cbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IFRyYW5zZmVyYWJsZU91dHB1dEVycm9yKFwiVGVzdGluZyBUcmFuc2ZlcmFibGVPdXRwdXRFcnJvclwiKVxuICAgIH0pLnRvVGhyb3coXCJUZXN0aW5nIFRyYW5zZmVyYWJsZU91dHB1dEVycm9yXCIpXG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBUcmFuc2ZlcmFibGVPdXRwdXRFcnJvcihcIlRlc3RpbmcgVHJhbnNmZXJhYmxlT3V0cHV0RXJyb3JcIilcbiAgICB9KS50b1Rocm93RXJyb3IoKVxuICB9KVxuXG4gIHRlc3QoXCJUcmFuc2ZlcmFibGVJbnB1dEVycm9yXCIsICgpOiB2b2lkID0+IHtcbiAgICB0cnkge1xuICAgICAgdGhyb3cgbmV3IFRyYW5zZmVyYWJsZUlucHV0RXJyb3IoXCJUZXN0aW5nIFRyYW5zZmVyYWJsZUlucHV0RXJyb3JcIilcbiAgICB9IGNhdGNoIChlcnJvcjogYW55KSB7XG4gICAgICBleHBlY3QoZXJyb3IuZ2V0Q29kZSgpKS50b0JlKFwiMTAxMFwiKVxuICAgIH1cbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IFRyYW5zZmVyYWJsZUlucHV0RXJyb3IoXCJUZXN0aW5nIFRyYW5zZmVyYWJsZUlucHV0RXJyb3JcIilcbiAgICB9KS50b1Rocm93KFwiVGVzdGluZyBUcmFuc2ZlcmFibGVJbnB1dEVycm9yXCIpXG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBUcmFuc2ZlcmFibGVJbnB1dEVycm9yKFwiVGVzdGluZyBUcmFuc2ZlcmFibGVJbnB1dEVycm9yXCIpXG4gICAgfSkudG9UaHJvd0Vycm9yKClcbiAgfSlcblxuICB0ZXN0KFwiSW5wdXRJZEVycm9yXCIsICgpOiB2b2lkID0+IHtcbiAgICB0cnkge1xuICAgICAgdGhyb3cgbmV3IElucHV0SWRFcnJvcihcIlRlc3RpbmcgSW5wdXRJZEVycm9yXCIpXG4gICAgfSBjYXRjaCAoZXJyb3I6IGFueSkge1xuICAgICAgZXhwZWN0KGVycm9yLmdldENvZGUoKSkudG9CZShcIjEwMTFcIilcbiAgICB9XG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBJbnB1dElkRXJyb3IoXCJUZXN0aW5nIElucHV0SWRFcnJvclwiKVxuICAgIH0pLnRvVGhyb3coXCJUZXN0aW5nIElucHV0SWRFcnJvclwiKVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgSW5wdXRJZEVycm9yKFwiVGVzdGluZyBJbnB1dElkRXJyb3JcIilcbiAgICB9KS50b1Rocm93RXJyb3IoKVxuICB9KVxuXG4gIHRlc3QoXCJPcGVyYXRpb25FcnJvclwiLCAoKTogdm9pZCA9PiB7XG4gICAgdHJ5IHtcbiAgICAgIHRocm93IG5ldyBPcGVyYXRpb25FcnJvcihcIlRlc3RpbmcgT3BlcmF0aW9uRXJyb3JcIilcbiAgICB9IGNhdGNoIChlcnJvcjogYW55KSB7XG4gICAgICBleHBlY3QoZXJyb3IuZ2V0Q29kZSgpKS50b0JlKFwiMTAxMlwiKVxuICAgIH1cbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IE9wZXJhdGlvbkVycm9yKFwiVGVzdGluZyBPcGVyYXRpb25FcnJvclwiKVxuICAgIH0pLnRvVGhyb3coXCJUZXN0aW5nIE9wZXJhdGlvbkVycm9yXCIpXG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBPcGVyYXRpb25FcnJvcihcIlRlc3RpbmcgT3BlcmF0aW9uRXJyb3JcIilcbiAgICB9KS50b1Rocm93RXJyb3IoKVxuICB9KVxuXG4gIHRlc3QoXCJJbnZhbGlkT3BlcmF0aW9uSWRFcnJvclwiLCAoKTogdm9pZCA9PiB7XG4gICAgdHJ5IHtcbiAgICAgIHRocm93IG5ldyBJbnZhbGlkT3BlcmF0aW9uSWRFcnJvcihcIlRlc3RpbmcgSW52YWxpZE9wZXJhdGlvbklkRXJyb3JcIilcbiAgICB9IGNhdGNoIChlcnJvcjogYW55KSB7XG4gICAgICBleHBlY3QoZXJyb3IuZ2V0Q29kZSgpKS50b0JlKFwiMTAxM1wiKVxuICAgIH1cbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IEludmFsaWRPcGVyYXRpb25JZEVycm9yKFwiVGVzdGluZyBJbnZhbGlkT3BlcmF0aW9uSWRFcnJvclwiKVxuICAgIH0pLnRvVGhyb3coXCJUZXN0aW5nIEludmFsaWRPcGVyYXRpb25JZEVycm9yXCIpXG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBJbnZhbGlkT3BlcmF0aW9uSWRFcnJvcihcIlRlc3RpbmcgSW52YWxpZE9wZXJhdGlvbklkRXJyb3JcIilcbiAgICB9KS50b1Rocm93RXJyb3IoKVxuICB9KVxuXG4gIHRlc3QoXCJDaGVja3N1bUVycm9yXCIsICgpOiB2b2lkID0+IHtcbiAgICB0cnkge1xuICAgICAgdGhyb3cgbmV3IENoZWNrc3VtRXJyb3IoXCJUZXN0aW5nIENoZWNrc3VtRXJyb3JcIilcbiAgICB9IGNhdGNoIChlcnJvcjogYW55KSB7XG4gICAgICBleHBlY3QoZXJyb3IuZ2V0Q29kZSgpKS50b0JlKFwiMTAxNFwiKVxuICAgIH1cbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IENoZWNrc3VtRXJyb3IoXCJUZXN0aW5nIENoZWNrc3VtRXJyb3JcIilcbiAgICB9KS50b1Rocm93KFwiVGVzdGluZyBDaGVja3N1bUVycm9yXCIpXG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBDaGVja3N1bUVycm9yKFwiVGVzdGluZyBDaGVja3N1bUVycm9yXCIpXG4gICAgfSkudG9UaHJvd0Vycm9yKClcbiAgfSlcblxuICB0ZXN0KFwiT3V0cHV0SWRFcnJvclwiLCAoKTogdm9pZCA9PiB7XG4gICAgdHJ5IHtcbiAgICAgIHRocm93IG5ldyBPdXRwdXRJZEVycm9yKFwiVGVzdGluZyBPdXRwdXRJZEVycm9yXCIpXG4gICAgfSBjYXRjaCAoZXJyb3I6IGFueSkge1xuICAgICAgZXhwZWN0KGVycm9yLmdldENvZGUoKSkudG9CZShcIjEwMTVcIilcbiAgICB9XG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBPdXRwdXRJZEVycm9yKFwiVGVzdGluZyBPdXRwdXRJZEVycm9yXCIpXG4gICAgfSkudG9UaHJvdyhcIlRlc3RpbmcgT3V0cHV0SWRFcnJvclwiKVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgT3V0cHV0SWRFcnJvcihcIlRlc3RpbmcgT3V0cHV0SWRFcnJvclwiKVxuICAgIH0pLnRvVGhyb3dFcnJvcigpXG4gIH0pXG5cbiAgdGVzdChcIlVUWE9FcnJvclwiLCAoKTogdm9pZCA9PiB7XG4gICAgdHJ5IHtcbiAgICAgIHRocm93IG5ldyBVVFhPRXJyb3IoXCJUZXN0aW5nIFVUWE9FcnJvclwiKVxuICAgIH0gY2F0Y2ggKGVycm9yOiBhbnkpIHtcbiAgICAgIGV4cGVjdChlcnJvci5nZXRDb2RlKCkpLnRvQmUoXCIxMDE2XCIpXG4gICAgfVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgVVRYT0Vycm9yKFwiVGVzdGluZyBVVFhPRXJyb3JcIilcbiAgICB9KS50b1Rocm93KFwiVGVzdGluZyBVVFhPRXJyb3JcIilcbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IFVUWE9FcnJvcihcIlRlc3RpbmcgVVRYT0Vycm9yXCIpXG4gICAgfSkudG9UaHJvd0Vycm9yKClcbiAgfSlcblxuICB0ZXN0KFwiSW5zdWZmaWNpZW50RnVuZHNFcnJvclwiLCAoKTogdm9pZCA9PiB7XG4gICAgdHJ5IHtcbiAgICAgIHRocm93IG5ldyBJbnN1ZmZpY2llbnRGdW5kc0Vycm9yKFwiVGVzdGluZyBJbnN1ZmZpY2llbnRGdW5kc0Vycm9yXCIpXG4gICAgfSBjYXRjaCAoZXJyb3I6IGFueSkge1xuICAgICAgZXhwZWN0KGVycm9yLmdldENvZGUoKSkudG9CZShcIjEwMTdcIilcbiAgICB9XG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBJbnN1ZmZpY2llbnRGdW5kc0Vycm9yKFwiVGVzdGluZyBJbnN1ZmZpY2llbnRGdW5kc0Vycm9yXCIpXG4gICAgfSkudG9UaHJvdyhcIlRlc3RpbmcgSW5zdWZmaWNpZW50RnVuZHNFcnJvclwiKVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgSW5zdWZmaWNpZW50RnVuZHNFcnJvcihcIlRlc3RpbmcgSW5zdWZmaWNpZW50RnVuZHNFcnJvclwiKVxuICAgIH0pLnRvVGhyb3dFcnJvcigpXG4gIH0pXG5cbiAgdGVzdChcIlRocmVzaG9sZEVycm9yXCIsICgpOiB2b2lkID0+IHtcbiAgICB0cnkge1xuICAgICAgdGhyb3cgbmV3IFRocmVzaG9sZEVycm9yKFwiVGVzdGluZyBUaHJlc2hvbGRFcnJvclwiKVxuICAgIH0gY2F0Y2ggKGVycm9yOiBhbnkpIHtcbiAgICAgIGV4cGVjdChlcnJvci5nZXRDb2RlKCkpLnRvQmUoXCIxMDE4XCIpXG4gICAgfVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgVGhyZXNob2xkRXJyb3IoXCJUZXN0aW5nIFRocmVzaG9sZEVycm9yXCIpXG4gICAgfSkudG9UaHJvdyhcIlRlc3RpbmcgVGhyZXNob2xkRXJyb3JcIilcbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IFRocmVzaG9sZEVycm9yKFwiVGVzdGluZyBUaHJlc2hvbGRFcnJvclwiKVxuICAgIH0pLnRvVGhyb3dFcnJvcigpXG4gIH0pXG5cbiAgdGVzdChcIlNFQ1BNaW50T3V0cHV0RXJyb3JcIiwgKCk6IHZvaWQgPT4ge1xuICAgIHRyeSB7XG4gICAgICB0aHJvdyBuZXcgU0VDUE1pbnRPdXRwdXRFcnJvcihcIlRlc3RpbmcgU0VDUE1pbnRPdXRwdXRFcnJvclwiKVxuICAgIH0gY2F0Y2ggKGVycm9yOiBhbnkpIHtcbiAgICAgIGV4cGVjdChlcnJvci5nZXRDb2RlKCkpLnRvQmUoXCIxMDE5XCIpXG4gICAgfVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgU0VDUE1pbnRPdXRwdXRFcnJvcihcIlRlc3RpbmcgU0VDUE1pbnRPdXRwdXRFcnJvclwiKVxuICAgIH0pLnRvVGhyb3coXCJUZXN0aW5nIFNFQ1BNaW50T3V0cHV0RXJyb3JcIilcbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IFNFQ1BNaW50T3V0cHV0RXJyb3IoXCJUZXN0aW5nIFNFQ1BNaW50T3V0cHV0RXJyb3JcIilcbiAgICB9KS50b1Rocm93RXJyb3IoKVxuICB9KVxuXG4gIHRlc3QoXCJFVk1JbnB1dEVycm9yXCIsICgpOiB2b2lkID0+IHtcbiAgICB0cnkge1xuICAgICAgdGhyb3cgbmV3IEVWTUlucHV0RXJyb3IoXCJUZXN0aW5nIEVWTUlucHV0RXJyb3JcIilcbiAgICB9IGNhdGNoIChlcnJvcjogYW55KSB7XG4gICAgICBleHBlY3QoZXJyb3IuZ2V0Q29kZSgpKS50b0JlKFwiMTAyMFwiKVxuICAgIH1cbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IEVWTUlucHV0RXJyb3IoXCJUZXN0aW5nIEVWTUlucHV0RXJyb3JcIilcbiAgICB9KS50b1Rocm93KFwiVGVzdGluZyBFVk1JbnB1dEVycm9yXCIpXG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBFVk1JbnB1dEVycm9yKFwiVGVzdGluZyBFVk1JbnB1dEVycm9yXCIpXG4gICAgfSkudG9UaHJvd0Vycm9yKClcbiAgfSlcblxuICB0ZXN0KFwiRVZNT3V0cHV0RXJyb3JcIiwgKCk6IHZvaWQgPT4ge1xuICAgIHRyeSB7XG4gICAgICB0aHJvdyBuZXcgRVZNT3V0cHV0RXJyb3IoXCJUZXN0aW5nIEVWTU91dHB1dEVycm9yXCIpXG4gICAgfSBjYXRjaCAoZXJyb3I6IGFueSkge1xuICAgICAgZXhwZWN0KGVycm9yLmdldENvZGUoKSkudG9CZShcIjEwMjFcIilcbiAgICB9XG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBFVk1PdXRwdXRFcnJvcihcIlRlc3RpbmcgRVZNT3V0cHV0RXJyb3JcIilcbiAgICB9KS50b1Rocm93KFwiVGVzdGluZyBFVk1PdXRwdXRFcnJvclwiKVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgRVZNT3V0cHV0RXJyb3IoXCJUZXN0aW5nIEVWTU91dHB1dEVycm9yXCIpXG4gICAgfSkudG9UaHJvd0Vycm9yKClcbiAgfSlcblxuICB0ZXN0KFwiRmVlQXNzZXRFcnJvclwiLCAoKTogdm9pZCA9PiB7XG4gICAgdHJ5IHtcbiAgICAgIHRocm93IG5ldyBGZWVBc3NldEVycm9yKFwiVGVzdGluZyBGZWVBc3NldEVycm9yXCIpXG4gICAgfSBjYXRjaCAoZXJyb3I6IGFueSkge1xuICAgICAgZXhwZWN0KGVycm9yLmdldENvZGUoKSkudG9CZShcIjEwMjJcIilcbiAgICB9XG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBGZWVBc3NldEVycm9yKFwiVGVzdGluZyBGZWVBc3NldEVycm9yXCIpXG4gICAgfSkudG9UaHJvdyhcIlRlc3RpbmcgRmVlQXNzZXRFcnJvclwiKVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgRmVlQXNzZXRFcnJvcihcIlRlc3RpbmcgRmVlQXNzZXRFcnJvclwiKVxuICAgIH0pLnRvVGhyb3dFcnJvcigpXG4gIH0pXG5cbiAgdGVzdChcIlN0YWtlRXJyb3JcIiwgKCk6IHZvaWQgPT4ge1xuICAgIHRyeSB7XG4gICAgICB0aHJvdyBuZXcgU3Rha2VFcnJvcihcIlRlc3RpbmcgU3Rha2VFcnJvclwiKVxuICAgIH0gY2F0Y2ggKGVycm9yOiBhbnkpIHtcbiAgICAgIGV4cGVjdChlcnJvci5nZXRDb2RlKCkpLnRvQmUoXCIxMDIzXCIpXG4gICAgfVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgU3Rha2VFcnJvcihcIlRlc3RpbmcgU3Rha2VFcnJvclwiKVxuICAgIH0pLnRvVGhyb3coXCJUZXN0aW5nIFN0YWtlRXJyb3JcIilcbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IFN0YWtlRXJyb3IoXCJUZXN0aW5nIFN0YWtlRXJyb3JcIilcbiAgICB9KS50b1Rocm93RXJyb3IoKVxuICB9KVxuXG4gIHRlc3QoXCJUaW1lRXJyb3JcIiwgKCk6IHZvaWQgPT4ge1xuICAgIHRyeSB7XG4gICAgICB0aHJvdyBuZXcgVGltZUVycm9yKFwiVGVzdGluZyBUaW1lRXJyb3JcIilcbiAgICB9IGNhdGNoIChlcnJvcjogYW55KSB7XG4gICAgICBleHBlY3QoZXJyb3IuZ2V0Q29kZSgpKS50b0JlKFwiMTAyNFwiKVxuICAgIH1cbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IFRpbWVFcnJvcihcIlRlc3RpbmcgVGltZUVycm9yXCIpXG4gICAgfSkudG9UaHJvdyhcIlRlc3RpbmcgVGltZUVycm9yXCIpXG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBUaW1lRXJyb3IoXCJUZXN0aW5nIFRpbWVFcnJvclwiKVxuICAgIH0pLnRvVGhyb3dFcnJvcigpXG4gIH0pXG5cbiAgdGVzdChcIkRlbGVnYXRpb25GZWVFcnJvclwiLCAoKTogdm9pZCA9PiB7XG4gICAgdHJ5IHtcbiAgICAgIHRocm93IG5ldyBEZWxlZ2F0aW9uRmVlRXJyb3IoXCJUZXN0aW5nIERlbGVnYXRpb25GZWVFcnJvclwiKVxuICAgIH0gY2F0Y2ggKGVycm9yOiBhbnkpIHtcbiAgICAgIGV4cGVjdChlcnJvci5nZXRDb2RlKCkpLnRvQmUoXCIxMDI1XCIpXG4gICAgfVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgRGVsZWdhdGlvbkZlZUVycm9yKFwiVGVzdGluZyBEZWxlZ2F0aW9uRmVlRXJyb3JcIilcbiAgICB9KS50b1Rocm93KFwiVGVzdGluZyBEZWxlZ2F0aW9uRmVlRXJyb3JcIilcbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IERlbGVnYXRpb25GZWVFcnJvcihcIlRlc3RpbmcgRGVsZWdhdGlvbkZlZUVycm9yXCIpXG4gICAgfSkudG9UaHJvd0Vycm9yKClcbiAgfSlcblxuICB0ZXN0KFwiU3VibmV0T3duZXJFcnJvclwiLCAoKTogdm9pZCA9PiB7XG4gICAgdHJ5IHtcbiAgICAgIHRocm93IG5ldyBTdWJuZXRPd25lckVycm9yKFwiVGVzdGluZyBTdWJuZXRPd25lckVycm9yXCIpXG4gICAgfSBjYXRjaCAoZXJyb3I6IGFueSkge1xuICAgICAgZXhwZWN0KGVycm9yLmdldENvZGUoKSkudG9CZShcIjEwMjZcIilcbiAgICB9XG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBTdWJuZXRPd25lckVycm9yKFwiVGVzdGluZyBTdWJuZXRPd25lckVycm9yXCIpXG4gICAgfSkudG9UaHJvdyhcIlRlc3RpbmcgU3VibmV0T3duZXJFcnJvclwiKVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgU3VibmV0T3duZXJFcnJvcihcIlRlc3RpbmcgU3VibmV0T3duZXJFcnJvclwiKVxuICAgIH0pLnRvVGhyb3dFcnJvcigpXG4gIH0pXG5cbiAgdGVzdChcIkJ1ZmZlclNpemVFcnJvclwiLCAoKTogdm9pZCA9PiB7XG4gICAgdHJ5IHtcbiAgICAgIHRocm93IG5ldyBCdWZmZXJTaXplRXJyb3IoXCJUZXN0aW5nIEJ1ZmZlclNpemVFcnJvclwiKVxuICAgIH0gY2F0Y2ggKGVycm9yOiBhbnkpIHtcbiAgICAgIGV4cGVjdChlcnJvci5nZXRDb2RlKCkpLnRvQmUoXCIxMDI3XCIpXG4gICAgfVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgQnVmZmVyU2l6ZUVycm9yKFwiVGVzdGluZyBCdWZmZXJTaXplRXJyb3JcIilcbiAgICB9KS50b1Rocm93KFwiVGVzdGluZyBCdWZmZXJTaXplRXJyb3JcIilcbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IEJ1ZmZlclNpemVFcnJvcihcIlRlc3RpbmcgQnVmZmVyU2l6ZUVycm9yXCIpXG4gICAgfSkudG9UaHJvd0Vycm9yKClcbiAgfSlcblxuICB0ZXN0KFwiQWRkcmVzc0luZGV4RXJyb3JcIiwgKCk6IHZvaWQgPT4ge1xuICAgIHRyeSB7XG4gICAgICB0aHJvdyBuZXcgQWRkcmVzc0luZGV4RXJyb3IoXCJUZXN0aW5nIEFkZHJlc3NJbmRleEVycm9yXCIpXG4gICAgfSBjYXRjaCAoZXJyb3I6IGFueSkge1xuICAgICAgZXhwZWN0KGVycm9yLmdldENvZGUoKSkudG9CZShcIjEwMjhcIilcbiAgICB9XG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBBZGRyZXNzSW5kZXhFcnJvcihcIlRlc3RpbmcgQWRkcmVzc0luZGV4RXJyb3JcIilcbiAgICB9KS50b1Rocm93KFwiVGVzdGluZyBBZGRyZXNzSW5kZXhFcnJvclwiKVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgQWRkcmVzc0luZGV4RXJyb3IoXCJUZXN0aW5nIEFkZHJlc3NJbmRleEVycm9yXCIpXG4gICAgfSkudG9UaHJvd0Vycm9yKClcbiAgfSlcblxuICB0ZXN0KFwiUHVibGljS2V5RXJyb3JcIiwgKCk6IHZvaWQgPT4ge1xuICAgIHRyeSB7XG4gICAgICB0aHJvdyBuZXcgUHVibGljS2V5RXJyb3IoXCJUZXN0aW5nIFB1YmxpY0tleUVycm9yXCIpXG4gICAgfSBjYXRjaCAoZXJyb3I6IGFueSkge1xuICAgICAgZXhwZWN0KGVycm9yLmdldENvZGUoKSkudG9CZShcIjEwMjlcIilcbiAgICB9XG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBQdWJsaWNLZXlFcnJvcihcIlRlc3RpbmcgUHVibGljS2V5RXJyb3JcIilcbiAgICB9KS50b1Rocm93KFwiVGVzdGluZyBQdWJsaWNLZXlFcnJvclwiKVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgUHVibGljS2V5RXJyb3IoXCJUZXN0aW5nIFB1YmxpY0tleUVycm9yXCIpXG4gICAgfSkudG9UaHJvd0Vycm9yKClcbiAgfSlcblxuICB0ZXN0KFwiTWVyZ2VSdWxlRXJyb3JcIiwgKCk6IHZvaWQgPT4ge1xuICAgIHRyeSB7XG4gICAgICB0aHJvdyBuZXcgTWVyZ2VSdWxlRXJyb3IoXCJUZXN0aW5nIE1lcmdlUnVsZUVycm9yXCIpXG4gICAgfSBjYXRjaCAoZXJyb3I6IGFueSkge1xuICAgICAgZXhwZWN0KGVycm9yLmdldENvZGUoKSkudG9CZShcIjEwMzBcIilcbiAgICB9XG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBNZXJnZVJ1bGVFcnJvcihcIlRlc3RpbmcgTWVyZ2VSdWxlRXJyb3JcIilcbiAgICB9KS50b1Rocm93KFwiVGVzdGluZyBNZXJnZVJ1bGVFcnJvclwiKVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgTWVyZ2VSdWxlRXJyb3IoXCJUZXN0aW5nIE1lcmdlUnVsZUVycm9yXCIpXG4gICAgfSkudG9UaHJvd0Vycm9yKClcbiAgfSlcblxuICB0ZXN0KFwiQmFzZTU4RXJyb3JcIiwgKCk6IHZvaWQgPT4ge1xuICAgIHRyeSB7XG4gICAgICB0aHJvdyBuZXcgQmFzZTU4RXJyb3IoXCJUZXN0aW5nIEJhc2U1OEVycm9yXCIpXG4gICAgfSBjYXRjaCAoZXJyb3I6IGFueSkge1xuICAgICAgZXhwZWN0KGVycm9yLmdldENvZGUoKSkudG9CZShcIjEwMzFcIilcbiAgICB9XG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBCYXNlNThFcnJvcihcIlRlc3RpbmcgQmFzZTU4RXJyb3JcIilcbiAgICB9KS50b1Rocm93KFwiVGVzdGluZyBCYXNlNThFcnJvclwiKVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgQmFzZTU4RXJyb3IoXCJUZXN0aW5nIEJhc2U1OEVycm9yXCIpXG4gICAgfSkudG9UaHJvd0Vycm9yKClcbiAgfSlcblxuICB0ZXN0KFwiUHJpdmF0ZUtleUVycm9yXCIsICgpOiB2b2lkID0+IHtcbiAgICB0cnkge1xuICAgICAgdGhyb3cgbmV3IFByaXZhdGVLZXlFcnJvcihcIlRlc3RpbmcgUHJpdmF0ZUtleUVycm9yXCIpXG4gICAgfSBjYXRjaCAoZXJyb3I6IGFueSkge1xuICAgICAgZXhwZWN0KGVycm9yLmdldENvZGUoKSkudG9CZShcIjEwMzJcIilcbiAgICB9XG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBQcml2YXRlS2V5RXJyb3IoXCJUZXN0aW5nIFByaXZhdGVLZXlFcnJvclwiKVxuICAgIH0pLnRvVGhyb3coXCJUZXN0aW5nIFByaXZhdGVLZXlFcnJvclwiKVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgUHJpdmF0ZUtleUVycm9yKFwiVGVzdGluZyBQcml2YXRlS2V5RXJyb3JcIilcbiAgICB9KS50b1Rocm93RXJyb3IoKVxuICB9KVxuXG4gIHRlc3QoXCJOb2RlSWRFcnJvclwiLCAoKTogdm9pZCA9PiB7XG4gICAgdHJ5IHtcbiAgICAgIHRocm93IG5ldyBOb2RlSWRFcnJvcihcIlRlc3RpbmcgTm9kZUlkRXJyb3JcIilcbiAgICB9IGNhdGNoIChlcnJvcjogYW55KSB7XG4gICAgICBleHBlY3QoZXJyb3IuZ2V0Q29kZSgpKS50b0JlKFwiMTAzM1wiKVxuICAgIH1cbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IE5vZGVJZEVycm9yKFwiVGVzdGluZyBOb2RlSWRFcnJvclwiKVxuICAgIH0pLnRvVGhyb3coXCJUZXN0aW5nIE5vZGVJZEVycm9yXCIpXG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBOb2RlSWRFcnJvcihcIlRlc3RpbmcgTm9kZUlkRXJyb3JcIilcbiAgICB9KS50b1Rocm93RXJyb3IoKVxuICB9KVxuXG4gIHRlc3QoXCJIZXhFcnJvclwiLCAoKTogdm9pZCA9PiB7XG4gICAgdHJ5IHtcbiAgICAgIHRocm93IG5ldyBIZXhFcnJvcihcIlRlc3RpbmcgSGV4RXJyb3JcIilcbiAgICB9IGNhdGNoIChlcnJvcjogYW55KSB7XG4gICAgICBleHBlY3QoZXJyb3IuZ2V0Q29kZSgpKS50b0JlKFwiMTAzNFwiKVxuICAgIH1cbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IEhleEVycm9yKFwiVGVzdGluZyBIZXhFcnJvclwiKVxuICAgIH0pLnRvVGhyb3coXCJUZXN0aW5nIEhleEVycm9yXCIpXG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBIZXhFcnJvcihcIlRlc3RpbmcgSGV4RXJyb3JcIilcbiAgICB9KS50b1Rocm93RXJyb3IoKVxuICB9KVxuXG4gIHRlc3QoXCJUeXBlSWRFcnJvclwiLCAoKTogdm9pZCA9PiB7XG4gICAgdHJ5IHtcbiAgICAgIHRocm93IG5ldyBUeXBlSWRFcnJvcihcIlRlc3RpbmcgVHlwZUlkRXJyb3JcIilcbiAgICB9IGNhdGNoIChlcnJvcjogYW55KSB7XG4gICAgICBleHBlY3QoZXJyb3IuZ2V0Q29kZSgpKS50b0JlKFwiMTAzNVwiKVxuICAgIH1cbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IFR5cGVJZEVycm9yKFwiVGVzdGluZyBUeXBlSWRFcnJvclwiKVxuICAgIH0pLnRvVGhyb3coXCJUZXN0aW5nIFR5cGVJZEVycm9yXCIpXG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBUeXBlSWRFcnJvcihcIlRlc3RpbmcgVHlwZUlkRXJyb3JcIilcbiAgICB9KS50b1Rocm93RXJyb3IoKVxuICB9KVxuXG4gIHRlc3QoXCJUeXBlTmFtZUVycm9yXCIsICgpOiB2b2lkID0+IHtcbiAgICB0cnkge1xuICAgICAgdGhyb3cgbmV3IFR5cGVOYW1lRXJyb3IoXCJUZXN0aW5nIFR5cGVOYW1lRXJyb3JcIilcbiAgICB9IGNhdGNoIChlcnJvcjogYW55KSB7XG4gICAgICBleHBlY3QoZXJyb3IuZ2V0Q29kZSgpKS50b0JlKFwiMTA0MlwiKVxuICAgIH1cbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IFR5cGVOYW1lRXJyb3IoXCJUZXN0aW5nIFR5cGVOYW1lRXJyb3JcIilcbiAgICB9KS50b1Rocm93KFwiVGVzdGluZyBUeXBlTmFtZUVycm9yXCIpXG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBUeXBlTmFtZUVycm9yKFwiVGVzdGluZyBUeXBlTmFtZUVycm9yXCIpXG4gICAgfSkudG9UaHJvd0Vycm9yKClcbiAgfSlcblxuICB0ZXN0KFwiVW5rbm93blR5cGVFcnJvclwiLCAoKTogdm9pZCA9PiB7XG4gICAgdHJ5IHtcbiAgICAgIHRocm93IG5ldyBVbmtub3duVHlwZUVycm9yKFwiVGVzdGluZyBVbmtub3duVHlwZUVycm9yXCIpXG4gICAgfSBjYXRjaCAoZXJyb3I6IGFueSkge1xuICAgICAgZXhwZWN0KGVycm9yLmdldENvZGUoKSkudG9CZShcIjEwMzZcIilcbiAgICB9XG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBVbmtub3duVHlwZUVycm9yKFwiVGVzdGluZyBVbmtub3duVHlwZUVycm9yXCIpXG4gICAgfSkudG9UaHJvdyhcIlRlc3RpbmcgVW5rbm93blR5cGVFcnJvclwiKVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgVW5rbm93blR5cGVFcnJvcihcIlRlc3RpbmcgVW5rbm93blR5cGVFcnJvclwiKVxuICAgIH0pLnRvVGhyb3dFcnJvcigpXG4gIH0pXG5cbiAgdGVzdChcIkJlY2gzMkVycm9yXCIsICgpOiB2b2lkID0+IHtcbiAgICB0cnkge1xuICAgICAgdGhyb3cgbmV3IEJlY2gzMkVycm9yKFwiVGVzdGluZyBCZWNoMzJFcnJvclwiKVxuICAgIH0gY2F0Y2ggKGVycm9yOiBhbnkpIHtcbiAgICAgIGV4cGVjdChlcnJvci5nZXRDb2RlKCkpLnRvQmUoXCIxMDM3XCIpXG4gICAgfVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgQmVjaDMyRXJyb3IoXCJUZXN0aW5nIEJlY2gzMkVycm9yXCIpXG4gICAgfSkudG9UaHJvdyhcIlRlc3RpbmcgQmVjaDMyRXJyb3JcIilcbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IEJlY2gzMkVycm9yKFwiVGVzdGluZyBCZWNoMzJFcnJvclwiKVxuICAgIH0pLnRvVGhyb3dFcnJvcigpXG4gIH0pXG5cbiAgdGVzdChcIkVWTUZlZUVycm9yXCIsICgpOiB2b2lkID0+IHtcbiAgICB0cnkge1xuICAgICAgdGhyb3cgbmV3IEVWTUZlZUVycm9yKFwiVGVzdGluZyBFVk1GZWVFcnJvclwiKVxuICAgIH0gY2F0Y2ggKGVycm9yOiBhbnkpIHtcbiAgICAgIGV4cGVjdChlcnJvci5nZXRDb2RlKCkpLnRvQmUoXCIxMDM4XCIpXG4gICAgfVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgRVZNRmVlRXJyb3IoXCJUZXN0aW5nIEVWTUZlZUVycm9yXCIpXG4gICAgfSkudG9UaHJvdyhcIlRlc3RpbmcgRVZNRmVlRXJyb3JcIilcbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IEVWTUZlZUVycm9yKFwiVGVzdGluZyBFVk1GZWVFcnJvclwiKVxuICAgIH0pLnRvVGhyb3dFcnJvcigpXG4gIH0pXG5cbiAgdGVzdChcIkludmFsaWRFbnRyb3B5XCIsICgpOiB2b2lkID0+IHtcbiAgICB0cnkge1xuICAgICAgdGhyb3cgbmV3IEludmFsaWRFbnRyb3B5KFwiVGVzdGluZyBJbnZhbGlkRW50cm9weVwiKVxuICAgIH0gY2F0Y2ggKGVycm9yOiBhbnkpIHtcbiAgICAgIGV4cGVjdChlcnJvci5nZXRDb2RlKCkpLnRvQmUoXCIxMDM5XCIpXG4gICAgfVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgSW52YWxpZEVudHJvcHkoXCJUZXN0aW5nIEludmFsaWRFbnRyb3B5XCIpXG4gICAgfSkudG9UaHJvdyhcIlRlc3RpbmcgSW52YWxpZEVudHJvcHlcIilcbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IEludmFsaWRFbnRyb3B5KFwiVGVzdGluZyBJbnZhbGlkRW50cm9weVwiKVxuICAgIH0pLnRvVGhyb3dFcnJvcigpXG4gIH0pXG5cbiAgdGVzdChcIlByb3RvY29sRXJyb3JcIiwgKCk6IHZvaWQgPT4ge1xuICAgIHRyeSB7XG4gICAgICB0aHJvdyBuZXcgUHJvdG9jb2xFcnJvcihcIlRlc3RpbmcgUHJvdG9jb2xFcnJvclwiKVxuICAgIH0gY2F0Y2ggKGVycm9yOiBhbnkpIHtcbiAgICAgIGV4cGVjdChlcnJvci5nZXRDb2RlKCkpLnRvQmUoXCIxMDQwXCIpXG4gICAgfVxuICAgIGV4cGVjdCgoKTogdm9pZCA9PiB7XG4gICAgICB0aHJvdyBuZXcgUHJvdG9jb2xFcnJvcihcIlRlc3RpbmcgUHJvdG9jb2xFcnJvclwiKVxuICAgIH0pLnRvVGhyb3coXCJUZXN0aW5nIFByb3RvY29sRXJyb3JcIilcbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IFByb3RvY29sRXJyb3IoXCJUZXN0aW5nIFByb3RvY29sRXJyb3JcIilcbiAgICB9KS50b1Rocm93RXJyb3IoKVxuICB9KVxuXG4gIHRlc3QoXCJTdWJuZXRJZEVycm9yXCIsICgpOiB2b2lkID0+IHtcbiAgICB0cnkge1xuICAgICAgdGhyb3cgbmV3IFN1Ym5ldElkRXJyb3IoXCJUZXN0aW5nIFN1Ym5ldElkRXJyb3JcIilcbiAgICB9IGNhdGNoIChlcnJvcjogYW55KSB7XG4gICAgICBleHBlY3QoZXJyb3IuZ2V0Q29kZSgpKS50b0JlKFwiMTA0MVwiKVxuICAgIH1cbiAgICBleHBlY3QoKCk6IHZvaWQgPT4ge1xuICAgICAgdGhyb3cgbmV3IFN1Ym5ldElkRXJyb3IoXCJUZXN0aW5nIFN1Ym5ldElkRXJyb3JcIilcbiAgICB9KS50b1Rocm93KFwiVGVzdGluZyBTdWJuZXRJZEVycm9yXCIpXG4gICAgZXhwZWN0KCgpOiB2b2lkID0+IHtcbiAgICAgIHRocm93IG5ldyBTdWJuZXRJZEVycm9yKFwiVGVzdGluZyBTdWJuZXRJZEVycm9yXCIpXG4gICAgfSkudG9UaHJvd0Vycm9yKClcbiAgfSlcbn0pXG4iXX0=