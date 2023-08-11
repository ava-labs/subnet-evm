const { ethers, BigNumber } = require("ethers");
const { expect, assert } = require("chai");
const utils = require("../utils")

const {
    _1e6,
    addMargin,
    alice,
    cancelOrderFromLimitOrder,
    encodeLimitOrderWithType,
    getAMMContract,
    getRandomSalt,
    getOrder,
    juror,
    multiplySize,
    multiplyPrice,
    placeOrderFromLimitOrder,
    removeAllAvailableMargin,
    waitForOrdersToMatch,
} = utils

// Testing juror precompile contract 
describe("Testing validateLiquidationOrderAndDetermineFillPrice",async function () {
    market = 0

    context("when liquidation amount is <= zero", async function () {
        it("returns error", async function () {
            let order = new Uint8Array(1024);
            let liquidationAmount = BigNumber.from(0)
            try {
                await juror.validateLiquidationOrderAndDetermineFillPrice(order, liquidationAmount)
            } catch (error) {
                error_message = JSON.parse(error.error.body).error.message
                expect(error_message).to.equal("invalid fillAmount")
                return
            }
            expect.fail("Expected throw not received");
        })
    })
    context("when liquidation amount is > zero", async function () {
        context("when order is invalid", async function () {
            context("when order's status is not placed", async function () {
                it("returns error when order was never placed", async function () {
                    let liquidationAmount = multiplySize(0.1)
                    orderPrice = multiplyPrice(1800)
                    longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether
                    shortOrderBaseAssetQuantity = multiplySize(-0.1) // short 0.1 ether
                    longOrder = getOrder(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, orderPrice, getRandomSalt(), false)

                    // try long order
                    try {
                        await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(longOrder), liquidationAmount)
                    } catch (error) {
                        error_message = JSON.parse(error.error.body).error.message
                        expect(error_message).to.equal("invalid order")
                    }

                    // try short order
                    shortOrder = getOrder(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, orderPrice, getRandomSalt(), false)
                    try {
                        await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(shortOrder), liquidationAmount)
                    } catch (error) {
                        error_message = JSON.parse(error.error.body).error.message
                        expect(error_message).to.equal("invalid order")
                        return
                    }
                    expect.fail("Expected throw not received");
                })
                it("returns error when order was cancelled", async function () {
                    margin = multiplyPrice(150000)
                    await addMargin(alice, margin)

                    let liquidationAmount = multiplySize(0.1)
                    orderPrice = multiplyPrice(1800)
                    longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether
                    longOrder = getOrder(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, orderPrice, getRandomSalt(), false)
                    await placeOrderFromLimitOrder(longOrder, alice)
                    await cancelOrderFromLimitOrder(longOrder, alice)

                    // try long order
                    try {
                        await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(longOrder), liquidationAmount)
                    } catch (error) {
                        error_message = JSON.parse(error.error.body).error.message
                        expect(error_message).to.equal("invalid order")
                    }

                    // try short order
                    shortOrderBaseAssetQuantity = multiplySize(-0.1) // short 0.1 ether
                    shortOrder = getOrder(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, orderPrice, getRandomSalt(), false)
                    await placeOrderFromLimitOrder(shortOrder, alice)
                    await cancelOrderFromLimitOrder(shortOrder, alice)
                    try {
                        await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(shortOrder), liquidationAmount)
                    } catch (error) {
                        await removeAllAvailableMargin(alice)
                        error_message = JSON.parse(error.error.body).error.message
                        expect(error_message).to.equal("invalid order")
                        return
                    }
                    expect.fail("Expected throw not received");
                })
                it("returns error when order was filled", async function () {
                    margin = multiplyPrice(150000)
                    await addMargin(alice, margin)
                    await addMargin(charlie, margin)

                    let liquidationAmount = multiplySize(0.1)
                    orderPrice = multiplyPrice(1800)
                    longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether

                    longOrder = getOrder(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, orderPrice, getRandomSalt(), false)
                    await placeOrderFromLimitOrder(longOrder, alice)

                    shortOrderBaseAssetQuantity = multiplySize(-0.1) // short 0.1 ether
                    shortOrder = getOrder(BigNumber.from(market), charlie.address, shortOrderBaseAssetQuantity, orderPrice, getRandomSalt(), false)
                    await placeOrderFromLimitOrder(shortOrder, charlie)

                    await waitForOrdersToMatch()

                    // try long order
                    try {
                        await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(longOrder), liquidationAmount)
                    } catch (error) {
                        error_message = JSON.parse(error.error.body).error.message
                        expect(error_message).to.equal("invalid order")
                    }

                    try {
                        await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(shortOrder), liquidationAmount)
                    } catch (error) {
                        error_message = JSON.parse(error.error.body).error.message
                        expect(error_message).to.equal("invalid order")
                        // cleanup
                        longOrder = getOrder(BigNumber.from(market), charlie.address, longOrderBaseAssetQuantity, orderPrice, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(longOrder, charlie)
                        shortOrder = getOrder(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, orderPrice, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(shortOrder, alice)
                        await waitForOrdersToMatch()
                        await removeAllAvailableMargin(alice)
                        await removeAllAvailableMargin(charlie)
                        return
                    }
                    expect.fail("Expected throw not received");
                })
            })
            context("when order's status is placed", async function () {
                context("when order's filled amount + liquidationAmount is > order's baseAssetQuantity", async function () {
                    it("returns error", async function () {
                        margin = multiplyPrice(150000)
                        await addMargin(alice, margin)

                        let liquidationAmount = multiplySize(0.2)
                        orderPrice = multiplyPrice(1800)
                        longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether
                        shortOrderBaseAssetQuantity = multiplySize(-0.1) // short 0.1 ether
                        longOrder = getOrder(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, orderPrice, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(longOrder, alice)

                        // try long order
                        try {
                            await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(longOrder), liquidationAmount)
                        } catch (error) {
                            error_message = JSON.parse(error.error.body).error.message
                            expect(error_message).to.equal("overfill")
                        }

                        await cancelOrderFromLimitOrder(longOrder, alice)

                        // try short order
                        shortOrder = getOrder(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, orderPrice, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(shortOrder, alice)

                        try {
                            await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(shortOrder), liquidationAmount)
                        } catch (error) {
                            error_message = JSON.parse(error.error.body).error.message
                            expect(error_message).to.equal("overfill")
                            await cancelOrderFromLimitOrder(shortOrder, alice)
                            await removeAllAvailableMargin(alice)
                            return
                        }
                        expect.fail("Expected throw not received");
                    })
                })
                context("when order's filled amount + liquidationAmount is <= order's baseAssetQuantity", async function () {
                    it.skip("returns error if order is reduceOnly and liquidationAmount > currentPosition", async function () {
                    })
                })
            })
        })
        context("when order is valid", async function () {
            context("when liquidationAmount is invalid", async function () {
                context("When liquidation amount is not multiple of minSizeRequirement", async function () {
                    it("returns error if liquidationAmount is greater than zero less than minSizeRequirement", async function () {
                        margin = multiplyPrice(150000)
                        await addMargin(alice, margin)

                        const amm = await getAMMContract(market) 
                        minSizeRequirement = await amm.minSizeRequirement()
                        liquidationAmount = minSizeRequirement.div(BigNumber.from(2))

                        orderPrice = multiplyPrice(1800)
                        longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether
                        shortOrderBaseAssetQuantity = multiplySize(-0.1) // short 0.1 ether
                        longOrder = getOrder(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, orderPrice, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(longOrder, alice)

                        // try long order
                        try {
                            await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(longOrder), liquidationAmount)
                        } catch (error) {
                            error_message = JSON.parse(error.error.body).error.message
                            expect(error_message).to.equal("not multiple")
                        }

                        await cancelOrderFromLimitOrder(longOrder, alice)

                        // try short order
                        shortOrder = getOrder(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, orderPrice, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(shortOrder, alice)
                        try {
                            await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(shortOrder), liquidationAmount)
                        } catch (error) {
                            error_message = JSON.parse(error.error.body).error.message
                            expect(error_message).to.equal("not multiple")
                            await cancelOrderFromLimitOrder(shortOrder, alice)
                            await removeAllAvailableMargin(alice)
                            return
                        }
                        expect.fail("Expected throw not received");
                    })
                    it("returns error if liquidationAmount is greater than minSizeRequirement but not a multiple", async function () {
                        margin = multiplyPrice(150000)
                        await addMargin(alice, margin)

                        const amm = await getAMMContract(market) 
                        minSizeRequirement = await amm.minSizeRequirement()
                        liquidationAmount = minSizeRequirement.mul(BigNumber.from(3)).div(BigNumber.from(2))

                        orderPrice = multiplyPrice(1800)
                        longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether
                        shortOrderBaseAssetQuantity = multiplySize(-0.1) // short 0.1 ether
                        longOrder = getOrder(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, orderPrice, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(longOrder, alice)

                        // try long order
                        try {
                            response = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(longOrder), liquidationAmount)
                        } catch (error) {
                            error_message = JSON.parse(error.error.body).error.message
                            expect(error_message).to.equal("not multiple")
                        }
                        await cancelOrderFromLimitOrder(longOrder, alice)
                        await waitForOrdersToMatch()

                        // try short order
                        shortOrder = getOrder(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, orderPrice, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(shortOrder, alice)

                        try {
                            await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(shortOrder), liquidationAmount)
                        } catch (error) {
                            error_message = JSON.parse(error.error.body).error.message
                            expect(error_message).to.equal("not multiple")
                            await cancelOrderFromLimitOrder(shortOrder, alice)
                            await removeAllAvailableMargin(alice)
                            return
                        }
                        expect.fail("Expected throw not received");
                    })
                })
            })
            context("When liquidationAmount is valid", async function () {
                context("For a long order", async function () {
                    it("returns error if price is less than liquidation lower bound price", async function () {
                        margin = multiplyPrice(150000)
                        await addMargin(alice, margin)

                        longOrderBaseAssetQuantity = multiplySize(0.3) // long 0.3 ether
                        liquidationAmount = multiplySize(0.2) // 0.2 ether
                        const amm = await getAMMContract(market) 
                        oraclePrice = (await amm.getUnderlyingPrice())
                        maxLiquidationPriceSpread = await amm.maxLiquidationPriceSpread()
                        // liqLowerBound = oraclePrice*(1e6 - liquidationPriceSpread)/1e6
                        liqLowerBound = oraclePrice.mul(_1e6.sub(maxLiquidationPriceSpread)).div(_1e6)
                        longOrderPrice = liqLowerBound.sub(1)

                        longOrder = getOrder(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(longOrder, alice)

                        try {
                            await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(longOrder), liquidationAmount)
                        } catch (error) {
                            expect(error.error.body).to.match(/OB_long_order_price_too_low/)
                            await cancelOrderFromLimitOrder(longOrder, alice)
                            await removeAllAvailableMargin(alice)
                            return
                        }
                        expect.fail("Expected throw not received");
                    })
                    it("returns upperBound as fillPrice if price is more than upperBound", async function () {
                        margin = multiplyPrice(150000)
                        await addMargin(alice, margin)

                        longOrderBaseAssetQuantity = multiplySize(0.3) // long 0.3 ether
                        liquidationAmount = multiplySize(0.2) // 0.2 ether
                        const amm = await getAMMContract(market) 
                        oraclePrice = (await amm.getUnderlyingPrice())
                        oraclePriceSpreadThreshold = (await amm.maxOracleSpreadRatio())
                        // upperBound = (oraclePrice*(1e6 + oraclePriceSpreadThreshold))/1e6
                        upperBound = oraclePrice.mul(_1e6.add(oraclePriceSpreadThreshold)).div(_1e6)

                        longOrderPrice1 = upperBound.add(BigNumber.from(1))
                        longOrder1 = getOrder(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, longOrderPrice1, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(longOrder1, alice)
                        responseLongOrder1 =  await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(longOrder1), liquidationAmount)
                        expect(responseLongOrder1.fillPrice.toString()).to.equal(upperBound.toString())
                        await cancelOrderFromLimitOrder(longOrder1, alice)

                        longOrderPrice2 = upperBound.add(BigNumber.from(1000))
                        longOrder2 = getOrder(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, longOrderPrice2, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(longOrder2, alice)
                        responseLongOrder2 = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(longOrder2), liquidationAmount)
                        expect(responseLongOrder2.fillPrice.toString()).to.equal(upperBound.toString())
                        await cancelOrderFromLimitOrder(longOrder2, alice)

                        longOrderPrice3 = upperBound.add(BigNumber.from(_1e6))
                        longOrder3 = getOrder(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, longOrderPrice3, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(longOrder3, alice)
                        responseLongOrder3 = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(longOrder3), liquidationAmount)
                        expect(responseLongOrder3.fillPrice.toString()).to.equal(upperBound.toString())
                        await cancelOrderFromLimitOrder(longOrder3, alice)

                        //cleanup
                        await removeAllAvailableMargin(alice)
                    })
                    it("returns longOrder's price as fillPrice if price is between lowerBound and upperBound", async function () {
                        margin = multiplyPrice(150000)
                        await addMargin(alice, margin)

                        longOrderBaseAssetQuantity = multiplySize(0.3) // long 0.3 ether
                        liquidationAmount = multiplySize(0.2) // 0.2 ether
                        const amm = await getAMMContract(market) 
                        oraclePrice = (await amm.getUnderlyingPrice())
                        oraclePriceSpreadThreshold = (await amm.maxOracleSpreadRatio())
                        // upperBound = (oraclePrice*(1e6 + oraclePriceSpreadThreshold))/1e6
                        upperBound = oraclePrice.mul(_1e6.add(oraclePriceSpreadThreshold)).div(_1e6)
                        lowerBound = oraclePrice.mul(_1e6.sub(oraclePriceSpreadThreshold)).div(_1e6)

                        longOrderPrice1 = upperBound.sub(BigNumber.from(1))
                        longOrder1 = getOrder(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, longOrderPrice1, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(longOrder1, alice)
                        responseLongOrder1 =  await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(longOrder1), liquidationAmount)
                        expect(responseLongOrder1.fillPrice.toString()).to.equal(longOrderPrice1.toString())
                        await cancelOrderFromLimitOrder(longOrder1, alice)

                        longOrderPrice2 = lowerBound
                        longOrder2 = getOrder(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, longOrderPrice2, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(longOrder2, alice)
                        responseLongOrder2 = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(longOrder2), liquidationAmount)
                        expect(responseLongOrder2.fillPrice.toString()).to.equal(longOrderPrice2.toString())
                        await cancelOrderFromLimitOrder(longOrder2, alice)

                        longOrderPrice3 = upperBound.add(lowerBound).div(2)
                        longOrder3 = getOrder(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, longOrderPrice3, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(longOrder3, alice)
                        responseLongOrder3 = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(longOrder3), liquidationAmount)
                        expect(responseLongOrder3.fillPrice.toString()).to.equal(longOrderPrice3.toString())
                        await cancelOrderFromLimitOrder(longOrder3, alice)

                        await removeAllAvailableMargin(alice)
                    }) 
                })
                context("For a short order", async function () {
                    it("returns lower bound as fillPrice if shortPrice is less than lowerBound", async function () {
                        margin = multiplyPrice(150000)
                        await addMargin(alice, margin)

                        shortOrderBaseAssetQuantity = multiplySize(-0.4) // short 0.4 ether
                        liquidationAmount = multiplySize(0.2) // 0.2 ether
                        const amm = await getAMMContract(market) 
                        oraclePrice = (await amm.getUnderlyingPrice())
                        oraclePriceSpreadThreshold = (await amm.maxOracleSpreadRatio())
                        // lowerBound = (oraclePrice*(1e6 - oraclePriceSpreadThreshold))/1e6
                        lowerBound = oraclePrice.mul(_1e6.sub(oraclePriceSpreadThreshold)).div(_1e6)

                        shortOrderPrice1 = lowerBound.sub(BigNumber.from(1))
                        shortOrder1 = getOrder(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, shortOrderPrice1, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(shortOrder1, alice)
                        responseShortOrder1 = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(shortOrder1), liquidationAmount)
                        expect(responseShortOrder1.fillPrice.toString()).to.equal(lowerBound.toString())
                        await cancelOrderFromLimitOrder(shortOrder1, alice)

                        shortOrderPrice2 = lowerBound.sub(BigNumber.from(1000))
                        shortOrder2 = getOrder(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, shortOrderPrice2, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(shortOrder2, alice)
                        responseShortOrder2 = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(shortOrder2), liquidationAmount)
                        expect(responseShortOrder2.fillPrice.toString()).to.equal(lowerBound.toString())
                        await cancelOrderFromLimitOrder(shortOrder2, alice)

                        shortOrderPrice3 = lowerBound.sub(BigNumber.from(_1e6))
                        shortOrder3 = getOrder(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, shortOrderPrice3, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(shortOrder3, alice)
                        responseShortOrder3 = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(shortOrder3), liquidationAmount)
                        expect(responseShortOrder3.fillPrice.toString()).to.equal(lowerBound.toString())
                        await cancelOrderFromLimitOrder(shortOrder3, alice)
                        await removeAllAvailableMargin(alice)
                    })
                    it("returns error if price is more than liquidation upperBound", async function () {
                        await removeAllAvailableMargin(alice)
                        margin = multiplyPrice(150000)
                        await addMargin(alice, margin)

                        const amm = await getAMMContract(market) 
                        oraclePrice = (await amm.getUnderlyingPrice())
                        maxLiquidationPriceSpread = (await amm.maxLiquidationPriceSpread())
                        // liqUpperBound = oraclePrice*(1e6 + maxLiquidationPriceSpread))
                        liqUpperBound = oraclePrice.mul(_1e6.add(maxLiquidationPriceSpread)).div(_1e6)
                        shortOrderPrice = liqUpperBound.add(BigNumber.from(1))
                        shortOrderBaseAssetQuantity = multiplySize(-0.4) // short 0.4 ether
                        liquidationAmount = multiplySize(0.2) // 0.2 ether

                        shortOrder = getOrder(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(shortOrder, alice)

                        try {
                            await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(shortOrder), liquidationAmount)
                        } catch (error) {
                            error_message = JSON.parse(error.error.body).error.message
                            expect(error_message).to.equal("OB_short_order_price_too_high")
                            await cancelOrderFromLimitOrder(shortOrder, alice)
                            await removeAllAvailableMargin(alice)
                            return
                        }
                        expect.fail("Expected throw not received");
                    })
                    it("returns shortOrder's price as fillPrice if price is between lowerBound and upperBound", async function () {
                        margin = multiplyPrice(150000)
                        await addMargin(alice, margin)

                        shortOrderBaseAssetQuantity = multiplySize(-0.4) // short 0.4 ether
                        liquidationAmount = multiplySize(0.2) // 0.2 ether
                        const amm = await getAMMContract(market) 
                        oraclePrice = (await amm.getUnderlyingPrice())
                        oraclePriceSpreadThreshold = (await amm.maxOracleSpreadRatio())
                        lowerBound = oraclePrice.mul(_1e6.sub(oraclePriceSpreadThreshold)).div(_1e6)
                        upperBound = oraclePrice.mul(_1e6.add(oraclePriceSpreadThreshold)).div(_1e6)

                        shortOrderPrice1 = upperBound.sub(BigNumber.from(1))
                        shortOrder1 = getOrder(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, shortOrderPrice1, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(shortOrder1, alice)
                        responseShortOrder1 = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(shortOrder1), liquidationAmount)
                        expect(responseShortOrder1.fillPrice.toString()).to.equal(shortOrderPrice1.toString())
                        await cancelOrderFromLimitOrder(shortOrder1, alice)

                        shortOrderPrice2 = lowerBound
                        shortOrder2 = getOrder(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, shortOrderPrice2, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(shortOrder2, alice)
                        responseShortOrder2 = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(shortOrder2), liquidationAmount)
                        expect(responseShortOrder2.fillPrice.toString()).to.equal(shortOrderPrice2.toString())
                        await cancelOrderFromLimitOrder(shortOrder2, alice)

                        shortOrderPrice3 = lowerBound.add(upperBound).div(2)
                        shortOrder3 = getOrder(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, shortOrderPrice3, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(shortOrder3, alice)
                        responseShortOrder3 = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderWithType(shortOrder3), liquidationAmount)
                        expect(responseShortOrder3.fillPrice.toString()).to.equal(shortOrderPrice3.toString())
                        await cancelOrderFromLimitOrder(shortOrder3, alice)
                        await removeAllAvailableMargin(alice)
                    })
                })
            })
        })
    })
})
