const { expect } = require("chai");
const { BigNumber } = require("ethers");
const utils = require("../utils")

const {
    addMargin,
    alice,
    cancelOrderFromLimitOrder,
    getOrderV2,
    getRandomSalt,
    juror,
    multiplyPrice,
    multiplySize,
    placeOrderFromLimitOrder,
    removeAllAvailableMargin,
} = utils

describe("Testing ValidateCancelLimitOrder", async function() {
    market = BigNumber.from(0)
    longBaseAssetQuantity = multiplySize(0.1)
    shortBaseAssetQuantity = multiplySize("-0.1")
    price = multiplyPrice(1800)
    salt = getRandomSalt()
    initialMargin = multiplyPrice(500000)

    context("when order's status is not placed", async function() {
        context("when order's status is invalid", async function() {
            it("should return error", async function() {
                assertLowMargin = false
                longOrder = getOrderV2(market, alice.address, longBaseAssetQuantity, price, salt)
                let { err, orderHash } = await juror.validateCancelLimitOrder(longOrder, alice.address, assertLowMargin)
                expect(err).to.equal("Invalid")
                expect(orderHash).to.equal(await utils.orderBook.getOrderHashV2(longOrder))

                shortOrder = getOrderV2(market, alice.address, shortBaseAssetQuantity, price, salt, true)
                ;({ err, orderHash } = await juror.validateCancelLimitOrder(shortOrder, alice.address, assertLowMargin))
                expect(err).to.equal("Invalid")
                expect(orderHash).to.equal(await utils.orderBook.getOrderHashV2(shortOrder))
            })
        })
        context("when order's status is cancelled", async function() {
            this.beforeEach(async function() {
                await addMargin(alice, initialMargin)
            })
            this.afterEach(async function() {
                await removeAllAvailableMargin(alice)
            })

            it("should return error", async function() {
                longOrder = getOrderV2(market, alice.address, longBaseAssetQuantity, price, salt)
                await placeOrderFromLimitOrder(longOrder, alice)
                await cancelOrderFromLimitOrder(longOrder, alice)
                let { err, orderHash } = await juror.validateCancelLimitOrder(longOrder, alice.address, assertLowMargin)
                expect(err).to.equal("Cancelled")
                expect(orderHash).to.equal(await utils.orderBook.getOrderHashV2(longOrder))
            })
        })
        it("should return error when order's status is filled", async function() {
        })
    })
})
