/**
 * @packageDocumentation
 * @module API-AVM-Credentials
 */
import { Credential } from "../../common/credentials";
/**
 * Takes a buffer representing the credential and returns the proper [[Credential]] instance.
 *
 * @param credid A number representing the credential ID parsed prior to the bytes passed in
 *
 * @returns An instance of an [[Credential]]-extended class.
 */
export declare const SelectCredentialClass: (credid: number, ...args: any[]) => Credential;
export declare class SECPCredential extends Credential {
    protected _typeName: string;
    protected _codecID: number;
    protected _typeID: number;
    /**
     * Set the codecID
     *
     * @param codecID The codecID to set
     */
    setCodecID(codecID: number): void;
    getCredentialID(): number;
    clone(): this;
    create(...args: any[]): this;
    select(id: number, ...args: any[]): Credential;
}
export declare class NFTCredential extends Credential {
    protected _typeName: string;
    protected _codecID: number;
    protected _typeID: number;
    /**
     * Set the codecID
     *
     * @param codecID The codecID to set
     */
    setCodecID(codecID: number): void;
    getCredentialID(): number;
    clone(): this;
    create(...args: any[]): this;
    select(id: number, ...args: any[]): Credential;
}
//# sourceMappingURL=credentials.d.ts.map