const assert = require("assert")

const testFn = (fnNameOrNames: string | string[], overrides = {}, debug = false) => async function () {
  const fnNames: string[] = Array.isArray(fnNameOrNames) ? fnNameOrNames : [fnNameOrNames]
  assert(fnNames.every(name => name.startsWith('test_')))

  return fnNames.reduce((p: Promise<undefined>, fnName) => p.then(async () => {
    const tx = await this.testContract[fnName](overrides)
    const txReceipt = await tx.wait().catch(err => err.receipt)

    const failed = txReceipt.status !== 0 ? await this.testContract.callStatic.failed() : true
    
    if (debug || failed) {
      console.log('')

      if (!txReceipt.events) console.warn('WARNING: No parseable events found in tx-receipt\n')

      txReceipt
        .events
        ?.filter(event => debug || event.event?.startsWith('log'))
        .map(event => event.args?.forEach(arg => console.log(arg)))

      console.log('')
    }

    assert(!failed, `${fnName} failed`)
  }), Promise.resolve());
}

export const test = (name, fnName, overrides = {}) => it(name, testFn(fnName, overrides));
test.only = (name, fnName, overrides = {}) => it.only(name, testFn(fnName, overrides));
test.debug = (name, fnName, overrides = {}) => it.only(name, testFn(fnName, overrides, true));
test.skip = (name, fnName, overrides = {}) => it.skip(name, testFn(fnName, overrides));