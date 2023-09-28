const { ethers, BigNumber } = require("ethers");
const { expect, assert } = require("chai");
const utils = require("../utils")

const {
    _1e6,
    addMargin,
    alice,
    cancelOrderFromLimitOrderV2,
    encodeLimitOrderV2,
    encodeLimitOrderV2WithType,
    getAMMContract,
    getOrderV2,
    getRandomSalt,
    getRequiredMarginForLongOrder,
    getRequiredMarginForShortOrder,
    juror,
    limitOrderBook,
    multiplySize,
    multiplyPrice,
    placeOrderFromLimitOrderV2,
    removeAllAvailableMargin,
    waitForOrdersToMatch,
} = utils

// Testing juror precompile contract 
describe("Testing validateLiquidationOrderAndDetermineFillPrice",async function () {
    let market = 0

    context("when liquidation amount is <= zero", async function () {
        it("returns error", async function () {
            let order = new Uint8Array(1024);
            let liquidationAmount = BigNumber.from(0)
            output = await juror.validateLiquidationOrderAndDetermineFillPrice(order, liquidationAmount)
            expect(output.err).to.equal("invalid fillAmount")
            expect(output.element).to.equal(2)
            expect(output.res.fillPrice.toNumber()).to.equal(0)
            expect(output.res.fillAmount.toNumber()).to.equal(0)
        })
    })
    context("when liquidation amount is > zero", async function () {
        context("when order is invalid", async function () {
            context("when order's status is not placed", async function () {
                context("when order was never placed", async function () {
                    let liquidationAmount = multiplySize(0.1)
                    let orderPrice = multiplyPrice(2000)
                    let longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether
                    let shortOrderBaseAssetQuantity = multiplySize(-0.1) // short 0.1 ether
                    let longOrder = getOrderV2(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, orderPrice, getRandomSalt())
                    let shortOrder = getOrderV2(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, orderPrice, getRandomSalt())
                    it("returns error for a longOrder", async function () {
                        output = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(longOrder), liquidationAmount)
                        expect(output.err).to.equal("invalid order")
                        expect(output.element).to.equal(0)
                        expect(output.res.fillPrice.toNumber()).to.equal(0)
                        expect(output.res.fillAmount.toNumber()).to.equal(0)
                    })
                    it("returns error for a shortOrder", async function () {
                        output = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(shortOrder), liquidationAmount)
                        expect(output.err).to.equal("invalid order")
                        expect(output.element).to.equal(0)
                        expect(output.res.fillPrice.toNumber()).to.equal(0)
                        expect(output.res.fillAmount.toNumber()).to.equal(0)
                    })
                })
                context("when order was cancelled", async function () {
                    let liquidationAmount = multiplySize(0.1)
                    let orderPrice = multiplyPrice(2000)
                    let longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether
                    let longOrder = getOrderV2(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, orderPrice, getRandomSalt())
                    let shortOrderBaseAssetQuantity = multiplySize(-0.1) // short 0.1 ether
                    let shortOrder = getOrderV2(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, orderPrice, getRandomSalt())

                    this.beforeAll(async function () {
                        requiredMarginForLongOrder = await getRequiredMarginForLongOrder(longOrder)
                        requiredMarginForShortOrder = await getRequiredMarginForShortOrder(shortOrder)
                        await addMargin(alice, requiredMarginForShortOrder.add(requiredMarginForLongOrder))
                    })
                    this.afterAll(async function () {
                        await removeAllAvailableMargin(alice)
                    })
                    it("returns error for a longOrder", async function () {
                        await placeOrderFromLimitOrderV2(longOrder, alice)
                        await cancelOrderFromLimitOrderV2(longOrder, alice)

                        output = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(longOrder), liquidationAmount)
                        expect(output.err).to.equal("invalid order")
                        expect(output.element).to.equal(0)
                        expect(output.res.fillPrice.toNumber()).to.equal(0)
                        expect(output.res.fillAmount.toNumber()).to.equal(0)
                    })
                    it("returns error for a shortOrder", async function () {
                        await placeOrderFromLimitOrderV2(shortOrder, alice)
                        await cancelOrderFromLimitOrderV2(shortOrder, alice)

                        output = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(shortOrder), liquidationAmount)
                        expect(output.err).to.equal("invalid order")
                        expect(output.element).to.equal(0)
                        expect(output.res.fillPrice.toNumber()).to.equal(0)
                        expect(output.res.fillAmount.toNumber()).to.equal(0)
                    })
                })
                context("when order was filled", async function () {
                    let liquidationAmount = multiplySize(0.1)
                    let orderPrice = multiplyPrice(2000)
                    let longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether
                    let longOrder = getOrderV2(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, orderPrice, getRandomSalt())
                    let shortOrderBaseAssetQuantity = multiplySize(-0.1) // short 0.1 ether
                    let shortOrder = getOrderV2(BigNumber.from(market), charlie.address, shortOrderBaseAssetQuantity, orderPrice, getRandomSalt())

                    this.beforeAll(async function () {
                        requiredMarginForLongOrder = await getRequiredMarginForLongOrder(longOrder)
                        requiredMarginForShortOrder = await getRequiredMarginForShortOrder(shortOrder)
                        await addMargin(alice, requiredMarginForLongOrder)
                        await addMargin(charlie, requiredMarginForShortOrder)
                        await placeOrderFromLimitOrderV2(longOrder, alice)
                        await placeOrderFromLimitOrderV2(shortOrder, charlie)
                        await waitForOrdersToMatch()
                    })
                    this.afterAll(async function () {
                        // alice should short and charlie should long to clean
                        let aliceOppositeOrder = getOrderV2(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, orderPrice, getRandomSalt())
                        let charlieOppositeOrder = getOrderV2(BigNumber.from(market), charlie.address, longOrderBaseAssetQuantity, orderPrice, getRandomSalt())
                        requiredMarginForAliceOppositeOrder = await getRequiredMarginForShortOrder(aliceOppositeOrder)
                        requiredMarginForCharlieOppositeOrder = await getRequiredMarginForLongOrder(charlieOppositeOrder)
                        await addMargin(alice, requiredMarginForAliceOppositeOrder)
                        await addMargin(charlie, requiredMarginForCharlieOppositeOrder)
                        await placeOrderFromLimitOrderV2(aliceOppositeOrder, alice)
                        await placeOrderFromLimitOrderV2(charlieOppositeOrder, charlie)
                        await waitForOrdersToMatch()
                        await removeAllAvailableMargin(alice)
                        await removeAllAvailableMargin(charlie)
                    })
                    it("returns error for a longOrder", async function () {
                        output = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(longOrder), liquidationAmount)
                        expect(output.err).to.equal("invalid order")
                        expect(output.element).to.equal(0)
                        expect(output.res.fillPrice.toNumber()).to.equal(0)
                        expect(output.res.fillAmount.toNumber()).to.equal(0)
                    })
                    it('returns error for a shortOrder', async function () {
                        output = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(shortOrder), liquidationAmount)
                        expect(output.err).to.equal("invalid order")
                        expect(output.element).to.equal(0)
                        expect(output.res.fillPrice.toNumber()).to.equal(0)
                        expect(output.res.fillAmount.toNumber()).to.equal(0)
                    })
                })
            })
            context("when order's status is placed", async function () {
                context("when order's filled amount + liquidationAmount is > order's baseAssetQuantity", async function () {
                    let liquidationAmount = multiplySize(0.2)
                    let orderPrice = multiplyPrice(2000)
                    let longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether
                    let shortOrderBaseAssetQuantity = multiplySize(-0.1) // short 0.1 ether
                    let longOrder = getOrderV2(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, orderPrice, getRandomSalt())
                    let shortOrder = getOrderV2(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, orderPrice, getRandomSalt())

                    context("for a longOrder", async function () {
                        this.beforeAll(async function () {
                            requiredMarginForLongOrder = await getRequiredMarginForLongOrder(longOrder)
                            await addMargin(alice, requiredMarginForLongOrder)
                            await placeOrderFromLimitOrderV2(longOrder, alice)
                        })
                        this.afterAll(async function () {
                            await cancelOrderFromLimitOrderV2(longOrder, alice)
                            await removeAllAvailableMargin(alice)
                        })
                        it("returns error", async function () {
                            output = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(longOrder), liquidationAmount)
                            expect(output.err).to.equal("overfill")
                            expect(output.element).to.equal(0)
                            expect(output.res.fillPrice.toNumber()).to.equal(0)
                            expect(output.res.fillAmount.toNumber()).to.equal(0)
                        })
                    })
                    context("for a shortOrder", async function () {
                        this.beforeAll(async function () {
                            requiredMarginForShortOrder = await getRequiredMarginForShortOrder(shortOrder)
                            await addMargin(alice, requiredMarginForShortOrder)
                            await placeOrderFromLimitOrderV2(shortOrder, alice)
                        })
                        this.afterAll(async function () {
                            await cancelOrderFromLimitOrderV2(shortOrder, alice)
                            await removeAllAvailableMargin(alice)
                        })
                        it("returns error", async function () {
                            output = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(shortOrder), liquidationAmount)
                            expect(output.err).to.equal("overfill")
                            expect(output.element).to.equal(0)
                            expect(output.res.fillPrice.toNumber()).to.equal(0)
                            expect(output.res.fillAmount.toNumber()).to.equal(0)
                        })
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
                    context("if liquidationAmount is greater than zero less than minSizeRequirement", async function () {
                        let shortOrderPrice = multiplyPrice(2001)
                        let longOrderPrice = multiplyPrice(1999)
                        let longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether
                        let shortOrderBaseAssetQuantity = multiplySize(-0.1) // short 0.1 ether
                        let longOrder = getOrderV2(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                        let shortOrder = getOrderV2(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), false)
                        let liquidationAmount

                        this.beforeAll(async function () {
                            const amm = await getAMMContract(market) 
                            minSizeRequirement = await amm.minSizeRequirement()
                            liquidationAmount = minSizeRequirement.div(BigNumber.from(2))
                            requiredMarginForLongOrder = await getRequiredMarginForLongOrder(longOrder)
                            requiredMarginForShortOrder = await getRequiredMarginForShortOrder(shortOrder)
                            await addMargin(alice, requiredMarginForShortOrder.add(requiredMarginForLongOrder))
                            await placeOrderFromLimitOrderV2(longOrder, alice)
                            await placeOrderFromLimitOrderV2(shortOrder, alice)
                        })
                        this.afterAll(async function () {
                            await cancelOrderFromLimitOrderV2(longOrder, alice)
                            await cancelOrderFromLimitOrderV2(shortOrder, alice)
                            await removeAllAvailableMargin(alice)
                        })

                        it("returns error for a long order", async function () {
                            output = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(longOrder), liquidationAmount)
                            expect(output.err).to.equal("not multiple")
                            expect(output.element).to.equal(2)
                            expect(output.res.fillPrice.toNumber()).to.equal(0)
                            expect(output.res.fillAmount.toNumber()).to.equal(0)
                        })
                        it("returns error for a short order", async function () {
                            output = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(shortOrder), liquidationAmount)
                            expect(output.err).to.equal("not multiple")
                            expect(output.element).to.equal(2)
                            expect(output.res.fillPrice.toNumber()).to.equal(0)
                            expect(output.res.fillAmount.toNumber()).to.equal(0)
                        })
                    })
                    context("if liquidationAmount is greater than minSizeRequirement but not a multiple", async function () {
                        let shortOrderPrice = multiplyPrice(2001)
                        let longOrderPrice = multiplyPrice(1999)
                        let longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether
                        let shortOrderBaseAssetQuantity = multiplySize(-0.1) // short 0.1 ether
                        let longOrder = getOrderV2(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                        let shortOrder = getOrderV2(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), false)
                        let liquidationAmount

                        this.beforeAll(async function () {
                            const amm = await getAMMContract(market) 
                            minSizeRequirement = await amm.minSizeRequirement()
                            liquidationAmount = minSizeRequirement.mul(3).div(2)
                            requiredMarginForLongOrder = await getRequiredMarginForLongOrder(longOrder)
                            requiredMarginForShortOrder = await getRequiredMarginForShortOrder(shortOrder)
                            await addMargin(alice, requiredMarginForShortOrder.add(requiredMarginForLongOrder))
                            await placeOrderFromLimitOrderV2(longOrder, alice)
                            await placeOrderFromLimitOrderV2(shortOrder, alice)
                        })
                        this.afterAll(async function () {
                            await cancelOrderFromLimitOrderV2(longOrder, alice)
                            await cancelOrderFromLimitOrderV2(shortOrder, alice)
                            await removeAllAvailableMargin(alice)
                        })

                        it("returns error for a long order", async function () {
                            output = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(longOrder), liquidationAmount)
                            expect(output.err).to.equal("not multiple")
                            expect(output.element).to.equal(2)
                            expect(output.res.fillPrice.toNumber()).to.equal(0)
                            expect(output.res.fillAmount.toNumber()).to.equal(0)
                        })
                        it("returns error for a short order", async function () {
                            output = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(shortOrder), liquidationAmount)
                            expect(output.err).to.equal("not multiple")
                            expect(output.element).to.equal(2)
                            expect(output.res.fillPrice.toNumber()).to.equal(0)
                            expect(output.res.fillAmount.toNumber()).to.equal(0)
                        })
                    })
                })
            })
            context("When liquidationAmount is valid", async function () {
                let liquidationAmount = multiplySize(0.2) // 0.2 ether
                let lowerBound, upperBound, liqLowerBound, liqUpperBound

                this.beforeAll(async function () {
                    const amm = await getAMMContract(market) 
                    let oraclePrice = (await amm.getUnderlyingPrice())
                    let maxLiquidationPriceSpread = await amm.maxLiquidationPriceSpread()
                    let oraclePriceSpreadThreshold = (await amm.maxOracleSpreadRatio())
                    liqLowerBound = oraclePrice.mul(_1e6.sub(maxLiquidationPriceSpread)).div(_1e6)
                    liqUpperBound = oraclePrice.mul(_1e6.add(maxLiquidationPriceSpread)).div(_1e6)
                    upperBound = oraclePrice.mul(_1e6.add(oraclePriceSpreadThreshold)).div(_1e6)
                    lowerBound = oraclePrice.mul(_1e6.sub(oraclePriceSpreadThreshold)).div(_1e6)
                })
                context("For a long order", async function () {
                    let longOrderBaseAssetQuantity = multiplySize(0.3) // long 0.3 ether
                    context("when price is less than liquidation lower bound price", async function () {
                        let longOrder
                        this.beforeEach(async function () {
                            longOrderPrice = liqLowerBound.sub(1)
                            longOrder = getOrderV2(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt())
                            requiredMargin = await getRequiredMarginForLongOrder(longOrder)
                            await addMargin(alice, requiredMargin)
                            await placeOrderFromLimitOrderV2(longOrder, alice)
                        })
                        this.afterEach(async function () {
                            await cancelOrderFromLimitOrderV2(longOrder, alice)
                            await removeAllAvailableMargin(alice)
                        })

                        it("returns error", async function () {
                            output = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(longOrder), liquidationAmount)
                            expect(output.err).to.equal("long price below lower bound")
                            expect(output.element).to.equal(0)
                            expect(output.res.fillPrice.toNumber()).to.equal(0)
                            expect(output.res.fillAmount.toNumber()).to.equal(0)
                        })
                    })
                    context("when price is more than upperBound", async function () {
                        let longOrder
                        this.beforeEach(async function () {
                            longOrderPrice = upperBound.add(BigNumber.from(1))
                            longOrder = getOrderV2(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt())
                            requiredMargin = await getRequiredMarginForLongOrder(longOrder)
                            await addMargin(alice, requiredMargin)
                            await placeOrderFromLimitOrderV2(longOrder, alice)
                        })
                        this.afterEach(async function () {
                            await cancelOrderFromLimitOrderV2(longOrder, alice)
                            await removeAllAvailableMargin(alice)
                        })
                        it("returns upperBound as fillPrice", async function () {
                            output =  await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(longOrder), liquidationAmount)
                            expect(output.err).to.equal("")
                            expect(output.element).to.equal(3)
                            expect(output.res.fillPrice.toString()).to.equal(upperBound.toString())
                            expect(output.res.fillAmount.toString()).to.equal(liquidationAmount.toString())
                            expect(output.res.instruction.mode).to.eq(1)
                            expectedOrderHash = await limitOrderBook.getOrderHash(longOrder)
                            expect(output.res.instruction.orderHash).to.eq(expectedOrderHash)
                            expect(output.res.encodedOrder).to.eq(encodeLimitOrderV2(longOrder))
                        })
                    })
                    context("if price is between liqLowerBound and upperBound", async function () {
                        let longOrder1, longOrder2, longOrder3
                        this.beforeEach(async function () {
                            longOrderPrice1 = upperBound.sub(BigNumber.from(1))
                            longOrder1 = getOrderV2(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, longOrderPrice1, getRandomSalt())
                            requiredMargin1 = await getRequiredMarginForLongOrder(longOrder1)
                            longOrderPrice2 = liqLowerBound
                            longOrder2 = getOrderV2(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, longOrderPrice2, getRandomSalt())
                            requiredMargin2 = await getRequiredMarginForLongOrder(longOrder2)
                            longOrderPrice3 = upperBound.add(liqLowerBound).div(2)
                            longOrder3 = getOrderV2(BigNumber.from(market), alice.address, longOrderBaseAssetQuantity, longOrderPrice3, getRandomSalt(), false)
                            requiredMargin3 = await getRequiredMarginForLongOrder(longOrder3)

                            await addMargin(alice, requiredMargin1.add(requiredMargin2).add(requiredMargin3))
                            await placeOrderFromLimitOrderV2(longOrder1, alice)
                            await placeOrderFromLimitOrderV2(longOrder2, alice)
                            await placeOrderFromLimitOrderV2(longOrder3, alice)
                        })
                        this.afterEach(async function () {
                            await cancelOrderFromLimitOrderV2(longOrder1, alice)
                            await cancelOrderFromLimitOrderV2(longOrder2, alice)
                            await cancelOrderFromLimitOrderV2(longOrder3, alice)
                            await removeAllAvailableMargin(alice)
                        })
                        it("returns longOrder's price as fillPrice", async function () {
                            responseLongOrder1 =  await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(longOrder1), liquidationAmount)
                            expect(responseLongOrder1.err).to.equal("")
                            expect(responseLongOrder1.element).to.equal(3)
                            expect(responseLongOrder1.res.fillPrice.toString()).to.equal(longOrder1.price.toString())
                            expect(responseLongOrder1.res.fillAmount.toString()).to.equal(liquidationAmount.toString())
                            expect(responseLongOrder1.res.instruction.mode).to.eq(1)
                            expectedOrderHash = await limitOrderBook.getOrderHash(longOrder1)
                            expect(responseLongOrder1.res.instruction.orderHash).to.eq(expectedOrderHash)
                            expect(responseLongOrder1.res.encodedOrder).to.eq(encodeLimitOrderV2(longOrder1))

                            responseLongOrder2 = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(longOrder2), liquidationAmount)
                            expect(responseLongOrder2.err).to.equal("")
                            expect(responseLongOrder2.element).to.equal(3)
                            expect(responseLongOrder2.res.fillPrice.toString()).to.equal(longOrder2.price.toString())
                            expect(responseLongOrder2.res.fillAmount.toString()).to.equal(liquidationAmount.toString())
                            expect(responseLongOrder2.res.instruction.mode).to.eq(1)
                            expectedOrderHash = await limitOrderBook.getOrderHash(longOrder2)
                            expect(responseLongOrder2.res.instruction.orderHash).to.eq(expectedOrderHash)
                            expect(responseLongOrder2.res.encodedOrder).to.eq(encodeLimitOrderV2(longOrder2))

                            responseLongOrder3 = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(longOrder3), liquidationAmount)
                            expect(responseLongOrder3.err).to.equal("")
                            expect(responseLongOrder3.element).to.equal(3)
                            expect(responseLongOrder3.res.fillPrice.toString()).to.equal(longOrder3.price.toString())
                            expect(responseLongOrder3.res.fillAmount.toString()).to.equal(liquidationAmount.toString())
                            expect(responseLongOrder3.res.instruction.mode).to.eq(1)
                            expectedOrderHash = await limitOrderBook.getOrderHash(longOrder3)
                            expect(responseLongOrder3.res.instruction.orderHash).to.eq(expectedOrderHash)
                            expect(responseLongOrder3.res.encodedOrder).to.eq(encodeLimitOrderV2(longOrder3))
                        })
                    }) 
                })
                context("For a short order", async function () {
                    let shortOrderBaseAssetQuantity = multiplySize(-0.4) // short 0.4 ether
                    context("when price is greater than liquidation upper bound price", async function () {
                        let shortOrder
                        this.beforeEach(async function () {
                            shortOrderPrice = liqUpperBound.add(1)
                            shortOrder = getOrderV2(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt())
                            requiredMargin = await getRequiredMarginForShortOrder(shortOrder)
                            await addMargin(alice, requiredMargin)
                            await placeOrderFromLimitOrderV2(shortOrder, alice)
                        })
                        this.afterEach(async function () {
                            await cancelOrderFromLimitOrderV2(shortOrder, alice)
                            await removeAllAvailableMargin(alice)
                        })

                        it("returns error if price is more than liquidation upperBound", async function () {
                            output = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(shortOrder), liquidationAmount)
                            expect(output.err).to.equal("short price above upper bound")
                            expect(output.element).to.equal(0)
                            expect(output.res.fillPrice.toNumber()).to.equal(0)
                            expect(output.res.fillAmount.toNumber()).to.equal(0)
                        })
                    })
                    context("when price is less than lowerBound", async function () {
                        this.beforeEach(async function () {
                            shortOrderPrice = lowerBound.sub(BigNumber.from(1))
                            shortOrder = getOrderV2(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt())
                            requiredMargin = await getRequiredMarginForShortOrder(shortOrder)
                            await addMargin(alice, requiredMargin)
                            await placeOrderFromLimitOrderV2(shortOrder, alice)
                        })
                        this.afterEach(async function () {
                            await cancelOrderFromLimitOrderV2(shortOrder, alice)
                            await removeAllAvailableMargin(alice)
                        })
                        it("returns lower bound as fillPrice", async function () {
                            output = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(shortOrder), liquidationAmount)
                            expect(output.err).to.equal("")
                            expect(output.element).to.equal(3)
                            expect(output.res.fillPrice.toString()).to.equal(lowerBound.toString())
                            expect(output.res.fillAmount.toString()).to.equal(liquidationAmount.mul(-1).toString())
                            expect(output.res.instruction.mode).to.eq(1)
                            expectedOrderHash = await limitOrderBook.getOrderHash(shortOrder)
                            expect(output.res.instruction.orderHash).to.eq(expectedOrderHash)
                            expect(output.res.encodedOrder).to.eq(encodeLimitOrderV2(shortOrder))
                        })
                    })
                    context("if price is between lowerBound and liqUpperBound", async function () {
                        this.beforeEach(async function () {
                            shortOrderPrice1 = lowerBound.add(BigNumber.from(1))
                            shortOrder1 = getOrderV2(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, shortOrderPrice1, getRandomSalt())
                            requiredMargin1 = await getRequiredMarginForShortOrder(shortOrder1)
                            shortOrderPrice2 = liqUpperBound
                            shortOrder2 = getOrderV2(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, shortOrderPrice2, getRandomSalt())
                            requiredMargin2 = await getRequiredMarginForShortOrder(shortOrder2)
                            shortOrderPrice3 = lowerBound.add(liqUpperBound).div(2)
                            shortOrder3 = getOrderV2(BigNumber.from(market), alice.address, shortOrderBaseAssetQuantity, shortOrderPrice3, getRandomSalt(), false)
                            requiredMargin3 = await getRequiredMarginForShortOrder(shortOrder3)

                            await addMargin(alice, requiredMargin1.add(requiredMargin2).add(requiredMargin3))
                            await placeOrderFromLimitOrderV2(shortOrder1, alice)
                            await placeOrderFromLimitOrderV2(shortOrder2, alice)
                            await placeOrderFromLimitOrderV2(shortOrder3, alice)
                        })
                        this.afterEach(async function () {
                            await cancelOrderFromLimitOrderV2(shortOrder1, alice)
                            await cancelOrderFromLimitOrderV2(shortOrder2, alice)
                            await cancelOrderFromLimitOrderV2(shortOrder3, alice)
                            await removeAllAvailableMargin(alice)
                        })
                        it("returns shortOrder's price as fillPrice if price is between lowerBound and upperBound", async function () {
                            output = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(shortOrder1), liquidationAmount)
                            expect(output.err).to.equal("")
                            expect(output.element).to.equal(3)
                            expect(output.res.fillPrice.toString()).to.equal(shortOrder1.price.toString())
                            expect(output.res.fillAmount.toString()).to.equal(liquidationAmount.mul(-1).toString())
                            expect(output.res.instruction.mode).to.eq(1)
                            expectedOrderHash = await limitOrderBook.getOrderHash(shortOrder1)
                            expect(output.res.instruction.orderHash).to.eq(expectedOrderHash)
                            expect(output.res.encodedOrder).to.eq(encodeLimitOrderV2(shortOrder1))

                            output = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(shortOrder2), liquidationAmount)
                            expect(output.err).to.equal("")
                            expect(output.element).to.equal(3)
                            expect(output.res.fillPrice.toString()).to.equal(shortOrder2.price.toString())
                            expect(output.res.fillAmount.toString()).to.equal(liquidationAmount.mul(-1).toString())
                            expect(output.res.instruction.mode).to.eq(1)
                            expectedOrderHash = await limitOrderBook.getOrderHash(shortOrder2)
                            expect(output.res.instruction.orderHash).to.eq(expectedOrderHash)
                            expect(output.res.encodedOrder).to.eq(encodeLimitOrderV2(shortOrder2))

                            output = await juror.validateLiquidationOrderAndDetermineFillPrice(encodeLimitOrderV2WithType(shortOrder3), liquidationAmount)
                            expect(output.err).to.equal("")
                            expect(output.element).to.equal(3)
                            expect(output.res.fillPrice.toString()).to.equal(shortOrder3.price.toString())
                            expect(output.res.fillAmount.toString()).to.equal(liquidationAmount.mul(-1).toString())
                            expect(output.res.instruction.mode).to.eq(1)
                            expectedOrderHash = await limitOrderBook.getOrderHash(shortOrder3)
                            expect(output.res.instruction.orderHash).to.eq(expectedOrderHash)
                            expect(output.res.encodedOrder).to.eq(encodeLimitOrderV2(shortOrder3))
                        })
                    })
                })
            })
        })
    })
})
