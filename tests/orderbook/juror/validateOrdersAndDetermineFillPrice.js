const { ethers, BigNumber } = require("ethers")
const { expect, assert } = require("chai")
const utils = require("../utils")

const {
    _1e6,
    addMargin,
    alice,
    bob,
    cancelOrderFromLimitOrderV2,
    disableValidatorMatching,
    enableValidatorMatching,
    encodeLimitOrderV2,
    encodeLimitOrderV2WithType,
    getAMMContract,
    getOrderV2,
    getRandomSalt,
    getRequiredMarginForLongOrder,
    getRequiredMarginForShortOrder,
    juror,
    multiplySize,
    multiplyPrice,
    orderBook,
    placeOrderFromLimitOrderV2,
    removeAllAvailableMargin,
    waitForOrdersToMatch,
} = utils

// Testing juror precompile contract 
describe("Test validateOrdersAndDetermineFillPrice", function () {
    let market = BigNumber.from(0)
    let longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether
    let shortOrderBaseAssetQuantity = multiplySize(-0.1) // 0.1 ether
    let longOrderPrice = multiplyPrice(2000)
    let shortOrderPrice = multiplyPrice(2000)

    context("when fillAmount is <= 0", async function () {
        it("returns error when fillAmount=0", async function () {
            output = await juror.validateOrdersAndDetermineFillPrice([1,1], 0)
            expect(output.err).to.equal("invalid fillAmount")
            expect(output.element).to.equal(2)
            expect(output.res.fillPrice.toNumber()).to.equal(0)
        })
        it("returns error when fillAmount<0", async function () {
            let fillAmount = BigNumber.from("-1")
            output = await juror.validateOrdersAndDetermineFillPrice([1,1], fillAmount)
            expect(output.err).to.equal("invalid fillAmount")
            expect(output.element).to.equal(2)
            expect(output.res.fillPrice.toNumber()).to.equal(0)
        })
    })
    context("when fillAmount is > 0", async function () {
        context("when either longOrder or shortOrder is invalid", async function () {
            context("when longOrder is invalid", async function () {
                context("when longOrder's status is not placed", async function () {
                    context("when longOrder was never placed", async function () {
                        it("returns error", async function () {
                            let fillAmount = longOrderBaseAssetQuantity
                            let longOrder = getOrderV2(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt())
                            let output = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderV2WithType(longOrder), encodeLimitOrderV2WithType(longOrder)], fillAmount)
                            expect(output.err).to.equal("invalid order")
                            expect(output.element).to.equal(0)
                            expect(output.res.fillPrice.toNumber()).to.equal(0)
                        })
                    })
                    context("if longOrder's status is cancelled", async function () {
                        let longOrder = getOrderV2(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt())
                        this.beforeAll(async function () {
                            requiredMargin = await getRequiredMarginForLongOrder(longOrder)
                            await addMargin(alice, requiredMargin)
                            await placeOrderFromLimitOrderV2(longOrder, alice)
                            await cancelOrderFromLimitOrderV2(longOrder, alice)
                        })
                        this.afterAll(async function () {
                            await removeAllAvailableMargin(alice)
                        })
                        it("returns error", async function () {
                            let fillAmount = longOrderBaseAssetQuantity
                            let output = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderV2WithType(longOrder), encodeLimitOrderV2WithType(longOrder)], fillAmount)
                            expect(output.err).to.equal("invalid order")
                            expect(output.element).to.equal(0)
                            expect(output.res.fillPrice.toNumber()).to.equal(0)
                        })
                    })
                    context("if longOrder's status is filled", async function () {
                        let longOrder = getOrderV2(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt())
                        let shortOrder = getOrderV2(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt())

                        this.beforeAll(async function () {
                            requiredMarginForLongOrder = await getRequiredMarginForLongOrder(longOrder)
                            await addMargin(alice, requiredMarginForLongOrder)
                            await placeOrderFromLimitOrderV2(longOrder, alice)
                            requiredMarginForShortOrder = await getRequiredMarginForShortOrder(shortOrder)
                            await addMargin(bob, requiredMarginForShortOrder)
                            await placeOrderFromLimitOrderV2(shortOrder, bob)
                            await waitForOrdersToMatch()
                        })
                        this.afterAll(async function () {
                            aliceOppositeOrder = getOrderV2(market, alice.address, shortOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), true)
                            bobOppositeOrder = getOrderV2(market, bob.address, longOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), true)
                            await placeOrderFromLimitOrderV2(aliceOppositeOrder, alice)
                            await placeOrderFromLimitOrderV2(bobOppositeOrder, bob)
                            await waitForOrdersToMatch()
                            await removeAllAvailableMargin(alice)
                            await removeAllAvailableMargin(bob)
                        })
                        it("returns error", async function () {
                            let fillAmount = longOrderBaseAssetQuantity
                            let output = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderV2WithType(longOrder), encodeLimitOrderV2WithType(shortOrder)], fillAmount)
                            expect(output.err).to.equal("invalid order")
                            expect(output.element).to.equal(0)
                            expect(output.res.fillPrice.toNumber()).to.equal(0)
                        })
                    })
                })
                context("when longOrder's status is placed", async function () {
                    context("when longOrder's baseAssetQuantity is negative", async function () {
                        let shortOrder = getOrderV2(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt())
                        this.beforeAll(async function () {
                            requiredMargin = await getRequiredMarginForShortOrder(shortOrder)
                            await addMargin(bob, requiredMargin)
                            await placeOrderFromLimitOrderV2(shortOrder, bob)
                        })
                        this.afterAll(async function () {
                            await cancelOrderFromLimitOrderV2(shortOrder, bob)
                            await removeAllAvailableMargin(bob)
                        })

                        it("returns error", async function () {
                            fillAmount = longOrderBaseAssetQuantity
                            let output = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderV2WithType(shortOrder), encodeLimitOrderV2WithType(shortOrder)], fillAmount)
                            expect(output.err).to.equal("not long")
                            expect(output.element).to.equal(0)
                            expect(output.res.fillPrice.toNumber()).to.equal(0)
                        })
                    })
                    context("when longOrder's baseAssetQuantity is positive", async function () {
                        context("when longOrder's unfilled < fillAmount", async function () {
                            let longOrder = getOrderV2(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt())
                            this.beforeAll(async function () {
                                requiredMargin = await getRequiredMarginForLongOrder(longOrder)
                                await addMargin(alice, requiredMargin)
                                await placeOrderFromLimitOrderV2(longOrder, alice)
                            })
                            this.afterAll(async function () {
                                await cancelOrderFromLimitOrderV2(longOrder, alice)
                                await removeAllAvailableMargin(alice)
                            })

                            it("returns error", async function () {
                                fillAmount = longOrderBaseAssetQuantity.mul(2)
                                let output = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderV2WithType(longOrder), encodeLimitOrderV2WithType(longOrder)], fillAmount)
                                expect(output.err).to.equal("overfill")
                                expect(output.element).to.equal(0)
                                expect(output.res.fillPrice.toNumber()).to.equal(0)
                            })
                        })
                        context("when longOrder's unfilled > fillAmount", async function () {
                            context.skip("when order is reduceOnly", async function () {
                                it("returns error if fillAmount > currentPosition of longOrder trader", async function () {
                                })
                            })
                        })
                    })
                })
            })
            context("when longOrder is valid", async function () {
                context("when shortOrder is invalid", async function () {
                    context("when shortOrder's status is not placed", async function () {
                        let fillAmount = longOrderBaseAssetQuantity
                        context("if shortOrder was never placed", async function () {
                            let longOrder = getOrderV2(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt())
                            let shortOrder = getOrderV2(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt())

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
                                output = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderV2WithType(longOrder), encodeLimitOrderV2WithType(shortOrder)], fillAmount)
                                expect(output.err).to.equal("invalid order")
                                expect(output.element).to.equal(1)
                                expect(output.res.fillPrice.toNumber()).to.equal(0)
                            })
                        })
                        context("if shortOrder's status is cancelled", async function () {
                            let longOrder = getOrderV2(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt())
                            let shortOrder = getOrderV2(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt())
                            this.beforeAll(async function () {
                                //placing short order first to avoid matching. We can use disableValidatorMatching() also
                                requiredMarginForShortOrder = await getRequiredMarginForShortOrder(shortOrder)
                                await addMargin(bob, requiredMarginForShortOrder)
                                await placeOrderFromLimitOrderV2(shortOrder, bob)
                                await cancelOrderFromLimitOrderV2(shortOrder, bob)
                                requiredMarginForLongOrder = await getRequiredMarginForLongOrder(longOrder)
                                await addMargin(alice, requiredMarginForLongOrder)
                                await placeOrderFromLimitOrderV2(longOrder, alice)
                            })
                            this.afterAll(async function () {
                                await cancelOrderFromLimitOrderV2(longOrder, alice)
                                await removeAllAvailableMargin(alice)
                                await removeAllAvailableMargin(bob)
                            })
                            it("returns error", async function () {
                                output = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderV2WithType(longOrder), encodeLimitOrderV2WithType(shortOrder)], fillAmount)
                                expect(output.err).to.equal("invalid order")
                                expect(output.element).to.equal(1)
                                expect(output.res.fillPrice.toNumber()).to.equal(0)
                            })
                        })
                        context("if shortOrder's status is filled", async function () {
                            let longOrder = getOrderV2(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt())
                            let shortOrder = getOrderV2(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt())
                            let longOrder2 = getOrderV2(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt())
                            this.beforeAll(async function () {
                                requiredMarginForLongOrder = await getRequiredMarginForLongOrder(longOrder)
                                await addMargin(alice, requiredMarginForLongOrder)
                                await placeOrderFromLimitOrderV2(longOrder, alice)
                                requiredMarginForShortOrder = await getRequiredMarginForShortOrder(shortOrder)
                                await addMargin(bob, requiredMarginForShortOrder)
                                await placeOrderFromLimitOrderV2(shortOrder, bob)
                                await waitForOrdersToMatch()
                                requiredMarginForLongOrder2 = await getRequiredMarginForLongOrder(longOrder2)
                                await addMargin(alice, requiredMarginForLongOrder2)
                                await placeOrderFromLimitOrderV2(longOrder2, alice)
                            })
                            this.afterAll(async function () {
                                await cancelOrderFromLimitOrderV2(longOrder2, alice)
                                aliceOppositeOrder = getOrderV2(market, alice.address, shortOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), true)
                                requiredMarginForAliceOppositeOrder = await getRequiredMarginForShortOrder(aliceOppositeOrder)
                                bobOppositeOrder = getOrderV2(market, bob.address, longOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), true)
                                requiredMarginForBobOppositeOrder = await getRequiredMarginForLongOrder(bobOppositeOrder)
                                await addMargin(alice, requiredMarginForAliceOppositeOrder)
                                await addMargin(bob, requiredMarginForBobOppositeOrder)
                                await placeOrderFromLimitOrderV2(aliceOppositeOrder, alice)
                                await placeOrderFromLimitOrderV2(bobOppositeOrder, bob)
                                await waitForOrdersToMatch()
                                await removeAllAvailableMargin(alice)
                                await removeAllAvailableMargin(bob)
                            })
                            it("returns error", async function () {
                                output = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderV2WithType(longOrder2), encodeLimitOrderV2WithType(shortOrder)], fillAmount)
                                expect(output.err).to.equal("invalid order")
                                expect(output.element).to.equal(1)
                                expect(output.res.fillPrice.toNumber()).to.equal(0)
                            })
                        })
                    })
                    context("when shortOrder's status is placed", async function () {
                        context("when shortOrder's baseAssetQuantity is positive", async function () {
                            let longOrder = getOrderV2(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt())
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
                                fillAmount = longOrderBaseAssetQuantity
                                output = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderV2WithType(longOrder), encodeLimitOrderV2WithType(longOrder)], fillAmount)
                                expect(output.err).to.equal("not short")
                                expect(output.element).to.equal(1)
                                expect(output.res.fillPrice.toNumber()).to.equal(0)
                            })
                        })
                        context("when shortOrder's baseAssetQuantity is negative", async function () {
                            context("when shortOrder's unfilled < fillAmount", async function () {
                                let newLongOrderPrice = multiplyPrice(1999)
                                let newShortOrderPrice = multiplyPrice(2001)
                                let longOrder = getOrderV2(market, alice.address, longOrderBaseAssetQuantity.mul(3), newLongOrderPrice, getRandomSalt())
                                let shortOrder = getOrderV2(market, bob.address, shortOrderBaseAssetQuantity, newShortOrderPrice, getRandomSalt())

                                this.beforeAll(async function () {
                                    requiredMarginForLongOrder = await getRequiredMarginForLongOrder(longOrder)
                                    await addMargin(alice, requiredMarginForLongOrder)
                                    await placeOrderFromLimitOrderV2(longOrder, alice)
                                    requiredMarginForShortOrder = await getRequiredMarginForShortOrder(shortOrder)
                                    await addMargin(bob, requiredMarginForShortOrder)
                                    await placeOrderFromLimitOrderV2(shortOrder, bob)
                                })
                                this.afterAll(async function () {
                                    await cancelOrderFromLimitOrderV2(longOrder, alice)
                                    await cancelOrderFromLimitOrderV2(shortOrder, bob)
                                    await removeAllAvailableMargin(alice)
                                    await removeAllAvailableMargin(bob)
                                })

                                it("returns error", async function () {
                                    fillAmount = shortOrderBaseAssetQuantity.abs().mul(2)
                                    output = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderV2WithType(longOrder), encodeLimitOrderV2WithType(shortOrder)], fillAmount)
                                    expect(output.err).to.equal("overfill")
                                    expect(output.element).to.equal(1)
                                    expect(output.res.fillPrice.toNumber()).to.equal(0)
                                })
                            })
                            context("when shortOrder's unfilled > fillAmount", async function () {
                                context.skip("when order is reduceOnly", async function () {
                                    it("returns error if fillAmount > currentPosition of shortOrder's trader", async function () {
                                    })
                                })
                            })
                        })
                    })
                })
            })
        })
        context("when both orders are valid", async function () {
            let amm, minSizeRequirement, lowerBoundPrice, upperBoundPrice
            this.beforeEach(async function () {
                await disableValidatorMatching()
                amm = await getAMMContract(market)
                minSizeRequirement = await amm.minSizeRequirement()
                let maxOracleSpreadRatio = await amm.maxOracleSpreadRatio()
                let oraclePrice = await amm.getUnderlyingPrice()
                upperBoundPrice = oraclePrice.mul(_1e6.add(maxOracleSpreadRatio)).div(_1e6)
                lowerBoundPrice = oraclePrice.mul(_1e6.sub(maxOracleSpreadRatio)).div(_1e6)
            })

            this.afterEach(async function () {
                await enableValidatorMatching()
            })

            context("when amm is different for long and short orders", async function () {
                it("returns error", async function () {
                    // needs deploying another amm
                })
            })
            context("when amm is same for long and short orders", async function () {
                context("when longOrder's price is less than shortOrder's price", async function () {
                    let longOrder = getOrderV2(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt())
                    let newShortOrderPrice = longOrderPrice.add(1)
                    let shortOrder = getOrderV2(market, bob.address, shortOrderBaseAssetQuantity, newShortOrderPrice, getRandomSalt())
                    this.beforeEach(async function () {
                        requiredMarginForLongOrder = await getRequiredMarginForLongOrder(longOrder)
                        await addMargin(alice, requiredMarginForLongOrder)
                        await placeOrderFromLimitOrderV2(longOrder, alice)
                        requiredMarginForShortOrder = await getRequiredMarginForShortOrder(shortOrder)
                        await addMargin(bob, requiredMarginForShortOrder)
                        await placeOrderFromLimitOrderV2(shortOrder, bob)
                    })
                    this.afterEach(async function () {
                        await cancelOrderFromLimitOrderV2(longOrder, alice)
                        await cancelOrderFromLimitOrderV2(shortOrder, bob)
                        await removeAllAvailableMargin(alice)
                        await removeAllAvailableMargin(bob)
                    })

                    it("returns error ", async function () {
                        fillAmount = longOrderBaseAssetQuantity
                        output = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderV2WithType(longOrder), encodeLimitOrderV2WithType(shortOrder)], fillAmount)
                        expect(output.err).to.equal("OB_orders_do_not_match")
                        expect(output.element).to.equal(2)
                        expect(output.res.fillPrice.toNumber()).to.equal(0)
                    })
                })
                context("when longOrder's price is greater than shortOrder's price", async function () {
                    context("when fillAmount is not a multiple of minSizeRequirement", async function () {
                        context("when fillAmount < minSizeRequirement", async function () {
                            let longOrder, shortOrder
                            this.beforeEach(async function () {
                                longOrder = getOrderV2(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt())
                                let requiredMarginForLongOrder = await getRequiredMarginForLongOrder(longOrder)
                                shortOrder = getOrderV2(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt())
                                let requiredMarginForShortOrder = await getRequiredMarginForShortOrder(shortOrder)
                                await addMargin(alice, requiredMarginForLongOrder)
                                await addMargin(bob, requiredMarginForShortOrder)
                                await placeOrderFromLimitOrderV2(longOrder, alice)
                                await placeOrderFromLimitOrderV2(shortOrder, bob)
                            })
                            this.afterEach(async function () {
                                await cancelOrderFromLimitOrderV2(longOrder, alice)
                                await cancelOrderFromLimitOrderV2(shortOrder, bob)
                                await removeAllAvailableMargin(alice)
                                await removeAllAvailableMargin(bob)
                            })

                            it("returns error", async function () {
                                let fillAmount = minSizeRequirement.div(2)
                                output = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderV2WithType(longOrder), encodeLimitOrderV2WithType(shortOrder)], fillAmount)
                                expect(output.err).to.equal("not multiple")
                                expect(output.element).to.equal(2)
                                expect(output.res.fillPrice.toNumber()).to.equal(0)
                            })
                        })              
                        context("when fillAmount > minSizeRequirement", async function () {
                            let longOrder, shortOrder
                            this.beforeEach(async function () {
                                longOrder = getOrderV2(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt())
                                let requiredMarginForLongOrder = await getRequiredMarginForLongOrder(longOrder)
                                shortOrder = getOrderV2(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt())
                                let requiredMarginForShortOrder = await getRequiredMarginForShortOrder(shortOrder)
                                await addMargin(alice, requiredMarginForLongOrder)
                                await addMargin(bob, requiredMarginForShortOrder)
                                await placeOrderFromLimitOrderV2(longOrder, alice)
                                await placeOrderFromLimitOrderV2(shortOrder, bob)
                            })
                            this.afterEach(async function () {
                                await cancelOrderFromLimitOrderV2(longOrder, alice)
                                await cancelOrderFromLimitOrderV2(shortOrder, bob)
                                await removeAllAvailableMargin(alice)
                                await removeAllAvailableMargin(bob)
                            })

                            it("returns error", async function () {
                                let fillAmount = minSizeRequirement.mul(3).div(2)
                                output = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderV2WithType(longOrder), encodeLimitOrderV2WithType(shortOrder)], fillAmount)
                                expect(output.err).to.equal("not multiple")
                                expect(output.element).to.equal(2)
                                expect(output.res.fillPrice.toNumber()).to.equal(0)
                            })
                        })
                    })
                    context("when fillAmount is a multiple of minSizeRequirement", async function () {
                        context("when longOrder price is less than lowerBoundPrice", async function () {
                            let longOrder, shortOrder
                            this.beforeEach(async function () {
                                let longOrderPrice = lowerBoundPrice.sub(1)
                                longOrder = getOrderV2(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt())
                                let requiredMarginForLongOrder = await getRequiredMarginForLongOrder(longOrder)
                                await addMargin(alice, requiredMarginForLongOrder)
                                let shortOrderPrice = longOrderPrice
                                shortOrder = getOrderV2(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt())
                                let requiredMarginForShortOrder = await getRequiredMarginForShortOrder(shortOrder)
                                await addMargin(bob, requiredMarginForShortOrder)
                                await placeOrderFromLimitOrderV2(longOrder, alice)
                                await placeOrderFromLimitOrderV2(shortOrder, bob)
                            })
                            this.afterEach(async function () {
                                await cancelOrderFromLimitOrderV2(longOrder, alice)
                                await cancelOrderFromLimitOrderV2(shortOrder, bob)
                                await removeAllAvailableMargin(alice)
                                await removeAllAvailableMargin(bob)
                            })

                            it("returns error", async function () {
                                fillAmount = minSizeRequirement.mul(3)
                                output = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderV2WithType(longOrder), encodeLimitOrderV2WithType(shortOrder)], fillAmount)
                                expect(output.err).to.equal("long price below lower bound")
                                expect(output.element).to.equal(0)
                                expect(output.res.fillPrice.toNumber()).to.equal(0)
                            })
                        })
                        context("when longOrder price is >= lowerBoundPrice", async function () {
                            context("when shortOrder price is greater than upperBoundPrice", async function () {
                                let longOrder, shortOrder
                                this.beforeEach(async function () {
                                    longOrderPrice = upperBoundPrice.add(1)
                                    longOrder = getOrderV2(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt())
                                    let requiredMarginForLongOrder = await getRequiredMarginForLongOrder(longOrder)
                                    await addMargin(alice, requiredMarginForLongOrder)
                                    shortOrderPrice = upperBoundPrice.add(1)
                                    shortOrder = getOrderV2(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt())
                                    let requiredMarginForShortOrder = await getRequiredMarginForShortOrder(shortOrder)
                                    await addMargin(bob, requiredMarginForShortOrder)
                                    await placeOrderFromLimitOrderV2(longOrder, alice)
                                    await placeOrderFromLimitOrderV2(shortOrder, bob)
                                })
                                this.afterEach(async function () {
                                    await cancelOrderFromLimitOrderV2(longOrder, alice)
                                    await cancelOrderFromLimitOrderV2(shortOrder, bob)
                                    await removeAllAvailableMargin(alice)
                                    await removeAllAvailableMargin(bob)
                                })

                                it("returns error", async function () {
                                    fillAmount = minSizeRequirement.mul(3)
                                    output = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderV2WithType(longOrder), encodeLimitOrderV2WithType(shortOrder)], fillAmount)
                                    expect(output.err).to.equal("short price above upper bound")
                                    expect(output.element).to.equal(1)
                                    expect(output.res.fillPrice.toNumber()).to.equal(0)
                                })
                            })
                            context("when shortOrder price is <= upperBoundPrice", async function () {
                                context("When longOrder was placed in earlier block than shortOrder", async function () {
                                    context("if longOrder price is greater than lowerBoundPrice but less than upperBoundPrice", async function () {
                                        let longOrder, shortOrder
                                        this.beforeEach(async function () {
                                            let longOrderPrice = lowerBoundPrice.add(1)
                                            longOrder = getOrderV2(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt())
                                            let requiredMarginForLongOrder = await getRequiredMarginForLongOrder(longOrder)
                                            await addMargin(alice, requiredMarginForLongOrder)
                                            let shortOrderPrice = longOrderPrice
                                            shortOrder = getOrderV2(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt())
                                            let requiredMarginForShortOrder = await getRequiredMarginForShortOrder(shortOrder)
                                            await addMargin(bob, requiredMarginForShortOrder)
                                            await placeOrderFromLimitOrderV2(longOrder, alice)
                                            await placeOrderFromLimitOrderV2(shortOrder, bob)
                                        })
                                        this.afterEach(async function () {
                                            await cancelOrderFromLimitOrderV2(longOrder, alice)
                                            await cancelOrderFromLimitOrderV2(shortOrder, bob)
                                            await removeAllAvailableMargin(alice)
                                            await removeAllAvailableMargin(bob)
                                        })
                                        it("returns longOrder's price as fillPrice", async function () {
                                            let fillAmount = minSizeRequirement.mul(3)
                                            response = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderV2WithType(longOrder), encodeLimitOrderV2WithType(shortOrder)], fillAmount)
                                            expect(response.err).to.equal("")
                                            expect(response.element).to.equal(3)
                                            expect(response.res.fillPrice.toString()).to.equal(longOrder.price.toString())
                                            expect(response.res.instructions.length).to.equal(2)
                                            //longOrder
                                            expect(response.res.instructions[0].ammIndex.toNumber()).to.equal(0)
                                            expect(response.res.instructions[0].trader).to.equal(alice.address)
                                            expect(response.res.instructions[0].mode).to.equal(1)
                                            longOrderHash = await limitOrderBook.getOrderHash(longOrder)
                                            expect(response.res.instructions[0].orderHash).to.equal(longOrderHash)

                                            //shortOrder
                                            expect(response.res.instructions[1].ammIndex.toNumber()).to.equal(0)
                                            expect(response.res.instructions[1].trader).to.equal(bob.address)
                                            expect(response.res.instructions[1].mode).to.equal(0)
                                            shortOrderHash = await limitOrderBook.getOrderHash(shortOrder)
                                            expect(response.res.instructions[1].orderHash).to.equal(shortOrderHash)

                                            expect(response.res.orderTypes.length).to.equal(2)
                                            expect(response.res.orderTypes[0]).to.equal(0)
                                            expect(response.res.orderTypes[1]).to.equal(0)

                                            expect(response.res.encodedOrders[0]).to.equal(encodeLimitOrderV2(longOrder))
                                            expect(response.res.encodedOrders[1]).to.equal(encodeLimitOrderV2(shortOrder))
                                        })
                                    })
                                    context("if longOrder price is greater than upperBoundPrice", async function () {
                                        let longOrder, shortOrder
                                        this.beforeEach(async function () {
                                            let longOrderPrice = upperBoundPrice.add(1)
                                            longOrder = getOrderV2(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt())
                                            let requiredMarginForLongOrder = await getRequiredMarginForLongOrder(longOrder)
                                            await addMargin(alice, requiredMarginForLongOrder)
                                            let shortOrderPrice = upperBoundPrice.sub(1)
                                            shortOrder = getOrderV2(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt())
                                            let requiredMarginForShortOrder = await getRequiredMarginForShortOrder(shortOrder)
                                            await addMargin(bob, requiredMarginForShortOrder)
                                            await placeOrderFromLimitOrderV2(longOrder, alice)
                                            await placeOrderFromLimitOrderV2(shortOrder, bob)
                                        })
                                        this.afterEach(async function () {
                                            await cancelOrderFromLimitOrderV2(longOrder, alice)
                                            await cancelOrderFromLimitOrderV2(shortOrder, bob)
                                            await removeAllAvailableMargin(alice)
                                            await removeAllAvailableMargin(bob)
                                        })
                                        it("returns upperBound as fillPrice", async function () {
                                            let fillAmount = minSizeRequirement.mul(3)
                                            response = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderV2WithType(longOrder), encodeLimitOrderV2WithType(shortOrder)], fillAmount)
                                            expect(response.err).to.equal("")
                                            expect(response.element).to.equal(3)
                                            expect(response.res.fillPrice.toString()).to.equal(upperBoundPrice.toString())
                                            expect(response.res.instructions.length).to.equal(2)
                                            //longOrder
                                            expect(response.res.instructions[0].ammIndex.toNumber()).to.equal(0)
                                            expect(response.res.instructions[0].trader).to.equal(alice.address)
                                            expect(response.res.instructions[0].mode).to.equal(1)
                                            longOrderHash = await limitOrderBook.getOrderHash(longOrder)
                                            expect(response.res.instructions[0].orderHash).to.equal(longOrderHash)

                                            //shortOrder
                                            expect(response.res.instructions[1].ammIndex.toNumber()).to.equal(0)
                                            expect(response.res.instructions[1].trader).to.equal(bob.address)
                                            expect(response.res.instructions[1].mode).to.equal(0)
                                            shortOrderHash = await limitOrderBook.getOrderHash(shortOrder)
                                            expect(response.res.instructions[1].orderHash).to.equal(shortOrderHash)

                                            expect(response.res.orderTypes.length).to.equal(2)
                                            expect(response.res.orderTypes[0]).to.equal(0)
                                            expect(response.res.orderTypes[1]).to.equal(0)

                                            expect(response.res.encodedOrders[0]).to.equal(encodeLimitOrderV2(longOrder))
                                            expect(response.res.encodedOrders[1]).to.equal(encodeLimitOrderV2(shortOrder))
                                        })
                                    })
                                })
                                context("When shortOrder was placed in same or earlier block than longOrder", async function () {
                                    context("if shortOrder price is less than upperBoundPrice greater than lowerBoundPrice", async function () {
                                        let longOrder, shortOrder
                                        this.beforeEach(async function () {
                                            let shortOrderPrice = upperBoundPrice.sub(1)
                                            shortOrder = getOrderV2(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt())
                                            let requiredMarginForShortOrder = await getRequiredMarginForShortOrder(shortOrder)
                                            await addMargin(bob, requiredMarginForShortOrder)
                                            let longOrderPrice = shortOrderPrice.add(2)
                                            longOrder = getOrderV2(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt())
                                            let requiredMarginForLongOrder = await getRequiredMarginForLongOrder(longOrder)
                                            await addMargin(alice, requiredMarginForLongOrder)
                                            await placeOrderFromLimitOrderV2(shortOrder, bob)
                                            await placeOrderFromLimitOrderV2(longOrder, alice)
                                        })
                                        this.afterEach(async function () {
                                            await cancelOrderFromLimitOrderV2(longOrder, alice)
                                            await cancelOrderFromLimitOrderV2(shortOrder, bob)
                                            await removeAllAvailableMargin(alice)
                                            await removeAllAvailableMargin(bob)
                                        })

                                        it("returns shortOrder's price as fillPrice", async function () {
                                            fillAmount = minSizeRequirement.mul(3)
                                            response = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderV2WithType(longOrder), encodeLimitOrderV2WithType(shortOrder)], fillAmount)
                                            expect(response.res.fillPrice.toString()).to.equal(shortOrder.price.toString())
                                            expect(response.res.instructions.length).to.equal(2)
                                            //longOrder
                                            expect(response.res.instructions[0].ammIndex.toNumber()).to.equal(0)
                                            expect(response.res.instructions[0].trader).to.equal(alice.address)
                                            expect(response.res.instructions[0].mode).to.equal(0)
                                            longOrderHash = await limitOrderBook.getOrderHash(longOrder)
                                            expect(response.res.instructions[0].orderHash).to.equal(longOrderHash)

                                            //shortOrder
                                            expect(response.res.instructions[1].ammIndex.toNumber()).to.equal(0)
                                            expect(response.res.instructions[1].trader).to.equal(bob.address)
                                            expect(response.res.instructions[1].mode).to.equal(1)
                                            shortOrderHash = await limitOrderBook.getOrderHash(shortOrder)
                                            expect(response.res.instructions[1].orderHash).to.equal(shortOrderHash)

                                            expect(response.res.orderTypes.length).to.equal(2)
                                            expect(response.res.orderTypes[0]).to.equal(0)
                                            expect(response.res.orderTypes[1]).to.equal(0)

                                            expect(response.res.encodedOrders[0]).to.equal(encodeLimitOrderV2(longOrder))
                                            expect(response.res.encodedOrders[1]).to.equal(encodeLimitOrderV2(shortOrder))
                                        })
                                    })
                                    context("returns lowerBoundPrice price as fillPrice if shortOrder's price is less than lowerBoundPrice", async function () {
                                        let longOrder, shortOrder
                                        this.beforeEach(async function () {
                                            let shortOrderPrice = lowerBoundPrice.sub(1)
                                            shortOrder = getOrderV2(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt())
                                            let requiredMarginForShortOrder = await getRequiredMarginForShortOrder(shortOrder)
                                            await addMargin(bob, requiredMarginForShortOrder)
                                            let longOrderPrice = shortOrderPrice.add(2)
                                            longOrder = getOrderV2(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt())
                                            let requiredMarginForLongOrder = await getRequiredMarginForLongOrder(longOrder)
                                            await addMargin(alice, requiredMarginForLongOrder)
                                            await placeOrderFromLimitOrderV2(shortOrder, bob)
                                            await placeOrderFromLimitOrderV2(longOrder, alice)
                                        })
                                        this.afterEach(async function () {
                                            await cancelOrderFromLimitOrderV2(longOrder, alice)
                                            await cancelOrderFromLimitOrderV2(shortOrder, bob)
                                            await removeAllAvailableMargin(alice)
                                            await removeAllAvailableMargin(bob)
                                        })
                                        it("returns lowerBoundPrice price as fillPrice", async function () {
                                            fillAmount = minSizeRequirement.mul(3)
                                            response = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderV2WithType(longOrder), encodeLimitOrderV2WithType(shortOrder)], fillAmount)
                                            expect(response.res.fillPrice.toString()).to.equal(lowerBoundPrice.toString())
                                            expect(response.res.instructions.length).to.equal(2)
                                            //longOrder
                                            expect(response.res.instructions[0].ammIndex.toNumber()).to.equal(0)
                                            expect(response.res.instructions[0].trader).to.equal(alice.address)
                                            expect(response.res.instructions[0].mode).to.equal(0)
                                            longOrderHash = await limitOrderBook.getOrderHash(longOrder)
                                            expect(response.res.instructions[0].orderHash).to.equal(longOrderHash)

                                            //shortOrder
                                            expect(response.res.instructions[1].ammIndex.toNumber()).to.equal(0)
                                            expect(response.res.instructions[1].trader).to.equal(bob.address)
                                            expect(response.res.instructions[1].mode).to.equal(1)
                                            shortOrderHash = await limitOrderBook.getOrderHash(shortOrder)
                                            expect(response.res.instructions[1].orderHash).to.equal(shortOrderHash)

                                            expect(response.res.orderTypes.length).to.equal(2)
                                            expect(response.res.orderTypes[0]).to.equal(0)
                                            expect(response.res.orderTypes[1]).to.equal(0)

                                            expect(response.res.encodedOrders[0]).to.equal(encodeLimitOrderV2(longOrder))
                                            expect(response.res.encodedOrders[1]).to.equal(encodeLimitOrderV2(shortOrder))
                                        })
                                    })
                                })
                            })
                        })
                    })
                })
            })
        })
    })
})
