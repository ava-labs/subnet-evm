import { ethers } from "hardhat"
const assert = require("assert")

type MethodObject = { method: string, debug?: boolean, overrides?: any, shouldFail?: boolean }
type FnNameOrObject = string | string[] | MethodObject | MethodObject[]
type MethodWithDebugAndOverrides = MethodObject & { debug: boolean, overrides: any, shouldFail: boolean }

const testFn = (fnNameOrObject: FnNameOrObject, overrides = {}, debug = false) => {
  const fnObjects: MethodWithDebugAndOverrides[] = (Array.isArray(fnNameOrObject) ? fnNameOrObject : [fnNameOrObject]).map(fnNameOrObject => {
    fnNameOrObject = typeof fnNameOrObject === 'string' ? { method: fnNameOrObject } : fnNameOrObject
    fnNameOrObject.overrides = Object.assign({}, overrides, fnNameOrObject.overrides ?? {})
    fnNameOrObject.debug = fnNameOrObject.debug ?? debug
    fnNameOrObject.shouldFail = fnNameOrObject.shouldFail ?? false
 
    return fnNameOrObject as MethodWithDebugAndOverrides
  })

  assert(fnObjects.every(({ method }) => method.startsWith('test_')), "Solidity test functions must be prefixed with 'test_'")

  return async function() {
    return fnObjects.reduce((p: Promise<undefined>, fn) => p.then(async () => {
      const contract = fn.overrides.from
        ? this.testContract.connect(await ethers.getSigner(fn.overrides.from))
        : this.testContract
      const tx = await contract[fn.method](fn.overrides).catch(err => {
        if (fn.shouldFail) {
          if (fn.debug) console.error(`smart contract call failed with error:\n${err}\n`)

          return { failed: true }
        }
 
        console.error("smart contract call failed with error:", err)
        throw err
      })

      if (tx.failed && fn.shouldFail) return

      const txReceipt = await tx.wait().catch(err => {
        if (fn.debug) console.error(`tx failed with error:\n${err}\n`)
        return err.receipt
      })

      const failed = txReceipt.status !== 0 ? await contract.callStatic.failed() : true
      if (fn.debug || failed) {
        console.log('')

        if (!txReceipt.events) console.warn('WARNING: No parseable events found in tx-receipt\n')

        txReceipt
          .events
          ?.filter(event => fn.debug || event.event?.startsWith('log'))
          .map(event => event.args?.forEach(arg => console.log(arg)))

        console.log('')
      }

      assert(!failed, `${fn.method} failed`)
    }), Promise.resolve())
  }
}

export const test = (name, fnNameOrObject, overrides = {}) => it(name, testFn(fnNameOrObject, overrides))
test.only = (name, fnNameOrObject, overrides = {}) => it.only(name, testFn(fnNameOrObject, overrides))
test.debug = (name, fnNameOrObject, overrides = {}) => it.only(name, testFn(fnNameOrObject, overrides, true))
test.skip = (name, fnNameOrObject, overrides = {}) => it.skip(name, testFn(fnNameOrObject, overrides))
