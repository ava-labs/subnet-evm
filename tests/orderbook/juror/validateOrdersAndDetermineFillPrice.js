const { ethers, BigNumber } = require("ethers");
const { expect, assert } = require("chai");
const utils = require("../utils")

const {
    _1e6,
    addMargin,
    alice,
    bob,
    cancelOrderFromLimitOrder,
    disableValidatorMatching,
    enableValidatorMatching,
    encodeLimitOrder,
    encodeLimitOrderWithType,
    getAMMContract,
    getOrder,
    getRandomSalt,
    juror,
    multiplySize,
    multiplyPrice,
    orderBook,
    placeOrderFromLimitOrder,
    removeAllAvailableMargin,
    waitForOrdersToMatch,
} = utils

// Testing juror precompile contract 
describe("Test validateOrdersAndDetermineFillPrice", function () {
    beforeEach(async function () {
        market = BigNumber.from(0)
        longOrderBaseAssetQuantity = multiplySize(0.1) // 0.1 ether
        shortOrderBaseAssetQuantity = multiplySize(-0.1) // 0.1 ether
        longOrderPrice = multiplyPrice(1800)
        shortOrderPrice = multiplyPrice(1800)
        initialMargin = multiplyPrice(150000)
    });

    context("when fillAmount is <= 0", async function () {
        it("returns error when fillAmount=0", async function () {
            try {
                await juror.validateOrdersAndDetermineFillPrice([1,1], 0)
            } catch (error) {
                error_message = JSON.parse(error.error.body).error.message
                expect(error_message).to.equal("invalid fillAmount")
            }
        })
        it("returns error when fillAmount<0", async function () {
            let fillAmount = BigNumber.from("-1")
            try {
                await juror.validateOrdersAndDetermineFillPrice([1,1], fillAmount)
            } catch (error) {
                error_message = JSON.parse(error.error.body).error.message
                expect(error_message).to.equal("invalid fillAmount")
            }
        })
    })
    context("when fillAmount is > 0", async function () {
        context("when either longOrder or shortOrder is invalid", async function () {
            context("when longOrder is invalid", async function () {
                context("when longOrder's status is not placed", async function () {
                    it("returns error if longOrder was never placed", async function () {
                        longOrder = getOrder(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                        fillAmount = longOrderBaseAssetQuantity
                        try {
                            await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderWithType(longOrder), encodeLimitOrderWithType(longOrder)], fillAmount)
                        } catch (error) {
                            error_message = JSON.parse(error.error.body).error.message
                            expect(error_message).to.equal("invalid order")
                            return
                        }
                        expect.fail('Expected throw not received');
                    });
                    it("returns error if longOrder's status is cancelled", async function () {
                        await addMargin(alice, initialMargin)
                        longOrder = getOrder(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(longOrder, alice)
                        await cancelOrderFromLimitOrder(longOrder, alice)
                        fillAmount = longOrderBaseAssetQuantity

                        try {
                            await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderWithType(longOrder), encodeLimitOrderWithType(longOrder)], fillAmount)
                        } catch (error) {
                            // cleanup
                            await removeAllAvailableMargin(alice)

                            error_message = JSON.parse(error.error.body).error.message
                            expect(error_message).to.equal("invalid order")
                            return
                        }
                        expect.fail('Expected throw not received');
                    });
                    it("returns error if longOrder's status is filled", async function () {
                        await addMargin(alice, initialMargin)
                        await addMargin(bob, initialMargin)
                        longOrder = getOrder(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(longOrder, alice)
                        shortOrder = getOrder(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(shortOrder, bob)
                        fillAmount = longOrderBaseAssetQuantity

                        await waitForOrdersToMatch()

                        try {
                            await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderWithType(longOrder), encodeLimitOrderWithType(longOrder)], fillAmount)
                        } catch (error) {
                            //cleanup
                            aliceOppositeOrder = getOrder(market, alice.address, shortOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), true)
                            bobOppositeOrder = getOrder(market, bob.address, longOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), true)
                            await placeOrderFromLimitOrder(aliceOppositeOrder, alice)
                            await placeOrderFromLimitOrder(bobOppositeOrder, bob)
                            await waitForOrdersToMatch()
                            await removeAllAvailableMargin(alice)
                            await removeAllAvailableMargin(bob)

                            error_message = JSON.parse(error.error.body).error.message
                            expect(error_message).to.equal("invalid order")
                            return
                        }
                        expect.fail('Expected throw not received');
                    });
                });
                context("when longOrder's status is placed", async function () {
                    context("when longOrder's baseAssetQuantity is negative", async function () {
                        it("returns error", async function () {
                            await addMargin(bob, initialMargin)
                            shortOrder = getOrder(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), false)
                            await placeOrderFromLimitOrder(shortOrder, bob)
                            fillAmount = longOrderBaseAssetQuantity

                            try {
                                await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderWithType(shortOrder), encodeLimitOrderWithType(shortOrder)], fillAmount)
                            } catch (error) {
                                //cleanup
                                await cancelOrderFromLimitOrder(shortOrder, bob)
                                await removeAllAvailableMargin(bob)

                                error_message = JSON.parse(error.error.body).error.message
                                expect(error_message).to.equal("not long")
                                return
                            }
                            expect.fail('Expected throw not received');
                        })
                    })
                    context("when longOrder's baseAssetQuantity is positive", async function () {
                        context("when longOrder's unfilled < fillAmount", async function () {
                            it("returns error", async function () {
                                await addMargin(alice, initialMargin)
                                longOrder = getOrder(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                                await placeOrderFromLimitOrder(longOrder, alice)
                                fillAmount = longOrderBaseAssetQuantity.mul(2)

                                try {
                                    await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderWithType(longOrder), encodeLimitOrderWithType(longOrder)], fillAmount)
                                } catch (error) {
                                    //cleanup
                                    await cancelOrderFromLimitOrder(longOrder, alice)
                                    await removeAllAvailableMargin(alice)

                                    error_message = JSON.parse(error.error.body).error.message
                                    expect(error_message).to.equal("overfill")
                                    return
                                }
                                expect.fail('Expected throw not received');
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
                        it("returns error if shortOrder was never placed", async function () {
                            await addMargin(alice, initialMargin)
                            longOrder = getOrder(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                            await placeOrderFromLimitOrder(longOrder, alice)
                            shortOrder = getOrder(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), false)
                            fillAmount = longOrderBaseAssetQuantity
                            try {
                                await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderWithType(longOrder), encodeLimitOrderWithType(shortOrder)], fillAmount)
                            } catch (error) {
                                // cleanup
                                await cancelOrderFromLimitOrder(longOrder, alice)
                                await removeAllAvailableMargin(alice)

                                error_message = JSON.parse(error.error.body).error.message
                                expect(error_message).to.equal("invalid order")
                                return
                            }
                            expect.fail('Expected throw not received');
                        });
                        it("returns error if shortOrder's status is cancelled", async function () {
                            await addMargin(bob, initialMargin)
                            shortOrder = getOrder(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), false)
                            await placeOrderFromLimitOrder(shortOrder, bob)
                            await cancelOrderFromLimitOrder(shortOrder, bob)
                            await addMargin(alice, initialMargin)
                            longOrder = getOrder(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                            await placeOrderFromLimitOrder(longOrder, alice)
                            fillAmount = longOrderBaseAssetQuantity
                            try {
                                await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderWithType(longOrder), encodeLimitOrderWithType(shortOrder)], fillAmount)
                            } catch (error) {
                                // cleanup
                                await cancelOrderFromLimitOrder(longOrder, alice)
                                await removeAllAvailableMargin(alice)
                                await removeAllAvailableMargin(bob)

                                error_message = JSON.parse(error.error.body).error.message
                                expect(error_message).to.equal("invalid order")
                                return
                            }
                            expect.fail('Expected throw not received');
                        });
                        it("returns error if shortOrder's status is filled", async function () {
                            await addMargin(bob, initialMargin)
                            shortOrder = getOrder(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), false)
                            await placeOrderFromLimitOrder(shortOrder, bob)
                            await addMargin(alice, initialMargin)
                            longOrder = getOrder(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                            await placeOrderFromLimitOrder(longOrder, alice)
                            fillAmount = longOrderBaseAssetQuantity
                            await waitForOrdersToMatch()

                            longOrder2 = getOrder(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                            await placeOrderFromLimitOrder(longOrder2, alice)

                            try {
                                await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderWithType(longOrder2), encodeLimitOrderWithType(shortOrder)], fillAmount)
                            } catch (error) {
                                // cleanup
                                await cancelOrderFromLimitOrder(longOrder2, alice)
                                shortOrder = getOrder(market, alice.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), false)
                                longOrder = getOrder(market, bob.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                                await placeOrderFromLimitOrder(shortOrder, alice)
                                await placeOrderFromLimitOrder(longOrder, bob)
                                await waitForOrdersToMatch()
                                await removeAllAvailableMargin(alice)
                                await removeAllAvailableMargin(bob)

                                error_message = JSON.parse(error.error.body).error.message
                                expect(error_message).to.equal("invalid order")
                                return
                            }
                            expect.fail('Expected throw not received');
                        });
                    });
                    context("when shortOrder's status is placed", async function () {
                        context("when shortOrder's baseAssetQuantity is positive", async function () {
                            it("returns error", async function () {
                                await addMargin(alice, initialMargin)
                                longOrder = getOrder(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                                await placeOrderFromLimitOrder(longOrder, alice)

                                fillAmount = longOrderBaseAssetQuantity

                                try {
                                    await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderWithType(longOrder), encodeLimitOrderWithType(longOrder)], fillAmount)
                                } catch (error) {
                                    // cleanup
                                    await cancelOrderFromLimitOrder(longOrder, alice)
                                    await removeAllAvailableMargin(alice)

                                    error_message = JSON.parse(error.error.body).error.message
                                    expect(error_message).to.equal("not short")
                                    return
                                }
                                expect.fail('Expected throw not received');
                            })
                        })
                        context("when shortOrder's baseAssetQuantity is negative", async function () {
                            context("when shortOrder's unfilled < fillAmount", async function () {
                                it("returns error", async function () {
                                    await disableValidatorMatching()
                                    await addMargin(alice, initialMargin)
                                    longOrder = getOrder(market, alice.address, longOrderBaseAssetQuantity.mul(3), longOrderPrice, getRandomSalt(), false)
                                    await placeOrderFromLimitOrder(longOrder, alice)
                                    await addMargin(bob, initialMargin)
                                    shortOrder = getOrder(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), false)
                                    await placeOrderFromLimitOrder(shortOrder, bob)

                                    fillAmount = shortOrderBaseAssetQuantity.abs().mul(2)

    
                                    try {
                                        await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderWithType(longOrder), encodeLimitOrderWithType(shortOrder)], fillAmount)
                                    } catch (error) {
                                        //cleanup
                                        await cancelOrderFromLimitOrder(longOrder, alice)
                                        await cancelOrderFromLimitOrder(shortOrder, bob)
                                        await enableValidatorMatching()
                                        await removeAllAvailableMargin(alice)
                                        await removeAllAvailableMargin(bob)
    
                                        error_message = JSON.parse(error.error.body).error.message
                                        expect(error_message).to.equal("overfill")
                                        return
                                    }
                                    expect.fail('Expected throw not received');
                                })
                            })
                            context("when shortOrder's unfilled > fillAmount", async function () {
                                context.skip("when order is reduceOnly", async function () {
                                    it("returns error if fillAmount > currentPosition of shortOrder's trader", async function () {
                                        console.log("stuff")
                                    })
                                })
                            })
                        })
                    })
                })
            })
        })
        context("when both orders are valid", async function () {
            this.beforeEach(async function () {
                await disableValidatorMatching()
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
                    it("returns error ", async function () {
                        await addMargin(alice, initialMargin)
                        longOrder = getOrder(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(longOrder, alice)
                        await addMargin(bob, initialMargin)
                        shortOrderPrice = longOrderPrice.add(1)
                        shortOrder = getOrder(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), false)
                        await placeOrderFromLimitOrder(shortOrder, bob)

                        fillAmount = longOrderBaseAssetQuantity

                        try {
                            await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderWithType(longOrder), encodeLimitOrderWithType(shortOrder)], fillAmount)
                        } catch (error) {
                            // cleanup
                            await cancelOrderFromLimitOrder(longOrder, alice)
                            await cancelOrderFromLimitOrder(shortOrder, bob)
                            await removeAllAvailableMargin(alice)
                            await removeAllAvailableMargin(bob)

                            error_message = JSON.parse(error.error.body).error.message
                            expect(error_message).to.equal("OB_orders_do_not_match")
                            return
                        }
                        expect.fail('Expected throw not received');
                    })
                })
                context("when longOrder's price is greater than shortOrder's price", async function () {
                    context("when fillAmount is not a multiple of minSizeRequirement", async function () {
                        it("returns error if fillAmount < minSizeRequirement", async function () {
                            amm = await getAMMContract(market)
                            minSizeRequirement = await amm.minSizeRequirement()
                            fillAmount = minSizeRequirement.div(2)

                            await addMargin(alice, initialMargin)
                            longOrder = getOrder(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                            await placeOrderFromLimitOrder(longOrder, alice)
                            await addMargin(bob, initialMargin)
                            shortOrder = getOrder(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), false)
                            await placeOrderFromLimitOrder(shortOrder, bob)

                            try {
                                await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderWithType(longOrder), encodeLimitOrderWithType(shortOrder)], fillAmount)
                            } catch (error) {
                                // cleanup
                                await cancelOrderFromLimitOrder(longOrder, alice)
                                await cancelOrderFromLimitOrder(shortOrder, bob)
                                await removeAllAvailableMargin(alice)
                                await removeAllAvailableMargin(bob)

                                error_message = JSON.parse(error.error.body).error.message
                                expect(error_message).to.equal("not multiple")
                                return
                            }
                            expect.fail('Expected throw not received');
                        })              
                        it("returns error if fillAmount > minSizeRequirement", async function () {
                            amm = await getAMMContract(market)
                            minSizeRequirement = await amm.minSizeRequirement()
                            fillAmount = minSizeRequirement.mul(3).div(2)

                            await addMargin(alice, initialMargin)
                            longOrder = getOrder(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                            await placeOrderFromLimitOrder(longOrder, alice)
                            await addMargin(bob, initialMargin)
                            shortOrder = getOrder(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), false)
                            await placeOrderFromLimitOrder(shortOrder, bob)

                            try {
                                await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderWithType(longOrder), encodeLimitOrderWithType(shortOrder)], fillAmount)
                            } catch (error) {
                                // cleanup
                                await cancelOrderFromLimitOrder(longOrder, alice)
                                await cancelOrderFromLimitOrder(shortOrder, bob)
                                await removeAllAvailableMargin(alice)
                                await removeAllAvailableMargin(bob)

                                error_message = JSON.parse(error.error.body).error.message
                                expect(error_message).to.equal("not multiple")
                                return
                            }
                            expect.fail('Expected throw not received');
                        })              
                    })
                    context("when fillAmount is a multiple of minSizeRequirement", async function () {
                        context("when longOrder price is less than lowerBoundPrice", async function () {
                            it("returns error", async function () {
                                amm = await getAMMContract(market)
                                minSizeRequirement = await amm.minSizeRequirement()
                                fillAmount = minSizeRequirement.mul(3)
                                maxOracleSpreadRatio = await amm.maxOracleSpreadRatio()
                                oraclePrice = await amm.getUnderlyingPrice()
                                lowerBoundPrice = oraclePrice.mul(_1e6.sub(maxOracleSpreadRatio)).div(_1e6)

                                await addMargin(alice, initialMargin)
                                longOrderPrice = lowerBoundPrice.sub(1)
                                longOrder = getOrder(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                                await placeOrderFromLimitOrder(longOrder, alice)
                                await addMargin(bob, initialMargin)
                                shortOrderPrice = longOrderPrice
                                shortOrder = getOrder(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), false)
                                await placeOrderFromLimitOrder(shortOrder, bob)

                                try {
                                    await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderWithType(longOrder), encodeLimitOrderWithType(shortOrder)], fillAmount)
                                } catch (error) {
                                    // cleanup
                                    await cancelOrderFromLimitOrder(longOrder, alice)
                                    await cancelOrderFromLimitOrder(shortOrder, bob)
                                    await removeAllAvailableMargin(alice)
                                    await removeAllAvailableMargin(bob)

                                    error_message = JSON.parse(error.error.body).error.message
                                    expect(error_message).to.equal("OB_long_order_price_too_low")
                                    return
                                }
                                expect.fail('Expected throw not received');
                            });
                        })
                        context("when longOrder price is >= lowerBoundPrice", async function () {
                            context("when shortOrder price is greater than upperBoundPrice", async function () {
                                it("returns error", async function () {
                                    amm = await getAMMContract(market)
                                    minSizeRequirement = await amm.minSizeRequirement()
                                    fillAmount = minSizeRequirement.mul(3)
                                    maxOracleSpreadRatio = await amm.maxOracleSpreadRatio()
                                    oraclePrice = await amm.getUnderlyingPrice()
                                    lowerBoundPrice = oraclePrice.mul(_1e6.sub(maxOracleSpreadRatio)).div(_1e6)
                                    upperBoundPrice = oraclePrice.mul(_1e6.add(maxOracleSpreadRatio)).div(_1e6)

                                    await addMargin(alice, initialMargin)
                                    longOrderPrice = upperBoundPrice.add(1)
                                    longOrder = getOrder(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                                    await placeOrderFromLimitOrder(longOrder, alice)
                                    await addMargin(bob, initialMargin)
                                    shortOrderPrice = upperBoundPrice.add(1)
                                    shortOrder = getOrder(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), false)
                                    await placeOrderFromLimitOrder(shortOrder, bob)

                                    try {
                                        await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderWithType(longOrder), encodeLimitOrderWithType(shortOrder)], fillAmount)
                                    } catch (error) {
                                        // cleanup
                                        await cancelOrderFromLimitOrder(longOrder, alice)
                                        await cancelOrderFromLimitOrder(shortOrder, bob)
                                        await removeAllAvailableMargin(alice)
                                        await removeAllAvailableMargin(bob)

                                        error_message = JSON.parse(error.error.body).error.message
                                        expect(error_message).to.equal("OB_short_order_price_too_high")
                                        return
                                    }
                                    expect.fail('Expected throw not received');
                                });
                            })
                            context("when shortOrder price is <= upperBoundPrice", async function () {
                                context("When longOrder was placed in earlier block than shortOrder", async function () {
                                    it("returns longOrder's price as fillPrice if longOrder price is greater than lowerBoundPrice but less than upperBoundPrice", async function () {
                                        amm = await getAMMContract(market)
                                        minSizeRequirement = await amm.minSizeRequirement()
                                        fillAmount = minSizeRequirement.mul(3)
                                        maxOracleSpreadRatio = await amm.maxOracleSpreadRatio()
                                        oraclePrice = await amm.getUnderlyingPrice()
                                        lowerBoundPrice = oraclePrice.mul(_1e6.sub(maxOracleSpreadRatio)).div(_1e6)
                                        upperBoundPrice = oraclePrice.mul(_1e6.add(maxOracleSpreadRatio)).div(_1e6)
                                        longOrderPrice = upperBoundPrice.sub(1)

                                        await addMargin(alice, initialMargin)
                                        longOrder = getOrder(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                                        await placeOrderFromLimitOrder(longOrder, alice)
                                        await addMargin(bob, initialMargin)
                                        shortOrderPrice = longOrderPrice
                                        shortOrder = getOrder(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), false)
                                        await placeOrderFromLimitOrder(shortOrder, bob)

                                        response = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderWithType(longOrder), encodeLimitOrderWithType(shortOrder)], fillAmount)
                                        // cleanup
                                        await cancelOrderFromLimitOrder(longOrder, alice)
                                        await cancelOrderFromLimitOrder(shortOrder, bob)
                                        await removeAllAvailableMargin(alice)
                                        await removeAllAvailableMargin(bob)

                                        expect(response.fillPrice.toString()).to.equal(longOrderPrice.toString())
                                        expect(response.instructions.length).to.equal(2)
                                        //longOrder
                                        expect(response.instructions[0].ammIndex.toNumber()).to.equal(0)
                                        expect(response.instructions[0].trader).to.equal(alice.address)
                                        expect(response.instructions[0].mode).to.equal(1)
                                        longOrderHash = await orderBook.getOrderHash(longOrder)
                                        expect(response.instructions[0].orderHash).to.equal(longOrderHash)

                                        //shortOrder
                                        expect(response.instructions[1].ammIndex.toNumber()).to.equal(0)
                                        expect(response.instructions[1].trader).to.equal(bob.address)
                                        expect(response.instructions[1].mode).to.equal(0)
                                        shortOrderHash = await orderBook.getOrderHash(shortOrder)
                                        expect(response.instructions[1].orderHash).to.equal(shortOrderHash)

                                        expect(response.orderTypes.length).to.equal(2)
                                        expect(response.orderTypes[0]).to.equal(0)
                                        expect(response.orderTypes[1]).to.equal(0)

                                        expect(response.encodedOrders[0]).to.equal(encodeLimitOrder(longOrder))
                                        expect(response.encodedOrders[1]).to.equal(encodeLimitOrder(shortOrder))
                                    });
                                    it("returns upperBound as fillPrice if longOrder price is greater than upperBoundPrice", async function () {
                                        amm = await getAMMContract(market)
                                        minSizeRequirement = await amm.minSizeRequirement()
                                        fillAmount = minSizeRequirement.mul(3)
                                        maxOracleSpreadRatio = await amm.maxOracleSpreadRatio()
                                        oraclePrice = await amm.getUnderlyingPrice()
                                        lowerBoundPrice = oraclePrice.mul(_1e6.sub(maxOracleSpreadRatio)).div(_1e6)
                                        upperBoundPrice = oraclePrice.mul(_1e6.add(maxOracleSpreadRatio)).div(_1e6)

                                        await addMargin(alice, initialMargin)
                                        longOrderPrice = upperBoundPrice.add(1)
                                        longOrder = getOrder(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                                        await placeOrderFromLimitOrder(longOrder, alice)
                                        await addMargin(bob, initialMargin)
                                        shortOrderPrice = upperBoundPrice.sub(1)
                                        shortOrder = getOrder(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), false)
                                        await placeOrderFromLimitOrder(shortOrder, bob)

                                        response = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderWithType(longOrder), encodeLimitOrderWithType(shortOrder)], fillAmount)
                                        // cleanup
                                        await cancelOrderFromLimitOrder(longOrder, alice)
                                        await cancelOrderFromLimitOrder(shortOrder, bob)
                                        await removeAllAvailableMargin(alice)
                                        await removeAllAvailableMargin(bob)

                                        expect(response.fillPrice.toString()).to.equal(upperBoundPrice.toString())
                                        expect(response.instructions.length).to.equal(2)
                                        //longOrder
                                        expect(response.instructions[0].ammIndex.toNumber()).to.equal(0)
                                        expect(response.instructions[0].trader).to.equal(alice.address)
                                        expect(response.instructions[0].mode).to.equal(1)
                                        longOrderHash = await orderBook.getOrderHash(longOrder)
                                        expect(response.instructions[0].orderHash).to.equal(longOrderHash)

                                        //shortOrder
                                        expect(response.instructions[1].ammIndex.toNumber()).to.equal(0)
                                        expect(response.instructions[1].trader).to.equal(bob.address)
                                        expect(response.instructions[1].mode).to.equal(0)
                                        shortOrderHash = await orderBook.getOrderHash(shortOrder)
                                        expect(response.instructions[1].orderHash).to.equal(shortOrderHash)

                                        expect(response.orderTypes.length).to.equal(2)
                                        expect(response.orderTypes[0]).to.equal(0)
                                        expect(response.orderTypes[1]).to.equal(0)

                                        expect(response.encodedOrders[0]).to.equal(encodeLimitOrder(longOrder))
                                        expect(response.encodedOrders[1]).to.equal(encodeLimitOrder(shortOrder))
                                    })
                                });
                                context("When shortOrder was placed in same or earlier block than longOrder", async function () {
                                    it("returns shortOrder's price as fillPrice if shortOrder price is less than upperBoundPrice greater than lowerBoundPrice", async function () {
                                        amm = await getAMMContract(market)
                                        minSizeRequirement = await amm.minSizeRequirement()
                                        fillAmount = minSizeRequirement.mul(3)
                                        maxOracleSpreadRatio = await amm.maxOracleSpreadRatio()
                                        oraclePrice = await amm.getUnderlyingPrice()
                                        lowerBoundPrice = oraclePrice.mul(_1e6.sub(maxOracleSpreadRatio)).div(_1e6)
                                        upperBoundPrice = oraclePrice.mul(_1e6.add(maxOracleSpreadRatio)).div(_1e6)

                                        await addMargin(bob, initialMargin)
                                        shortOrderPrice = upperBoundPrice.sub(1)
                                        shortOrder = getOrder(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), false)
                                        await placeOrderFromLimitOrder(shortOrder, bob)
                                        await addMargin(alice, initialMargin)
                                        longOrderPrice = shortOrderPrice
                                        longOrder = getOrder(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                                        await placeOrderFromLimitOrder(longOrder, alice)
                                        
                                        response = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderWithType(longOrder), encodeLimitOrderWithType(shortOrder)], fillAmount)
                                        // cleanup
                                        await cancelOrderFromLimitOrder(longOrder, alice)
                                        await cancelOrderFromLimitOrder(shortOrder, bob)
                                        await removeAllAvailableMargin(alice)
                                        await removeAllAvailableMargin(bob)

                                        expect(response.fillPrice.toString()).to.equal(shortOrderPrice.toString())
                                        expect(response.instructions.length).to.equal(2)
                                        //longOrder
                                        expect(response.instructions[0].ammIndex.toNumber()).to.equal(0)
                                        expect(response.instructions[0].trader).to.equal(alice.address)
                                        expect(response.instructions[0].mode).to.equal(0)
                                        longOrderHash = await orderBook.getOrderHash(longOrder)
                                        expect(response.instructions[0].orderHash).to.equal(longOrderHash)

                                        //shortOrder
                                        expect(response.instructions[1].ammIndex.toNumber()).to.equal(0)
                                        expect(response.instructions[1].trader).to.equal(bob.address)
                                        expect(response.instructions[1].mode).to.equal(1)
                                        shortOrderHash = await orderBook.getOrderHash(shortOrder)
                                        expect(response.instructions[1].orderHash).to.equal(shortOrderHash)

                                        expect(response.orderTypes.length).to.equal(2)
                                        expect(response.orderTypes[0]).to.equal(0)
                                        expect(response.orderTypes[1]).to.equal(0)

                                        expect(response.encodedOrders[0]).to.equal(encodeLimitOrder(longOrder))
                                        expect(response.encodedOrders[1]).to.equal(encodeLimitOrder(shortOrder))
                                    });
                                    it("returns lowerBoundPrice price as fillPrice if shortOrder's price is less than lowerBoundPrice", async function () {
                                        amm = await getAMMContract(market)
                                        minSizeRequirement = await amm.minSizeRequirement()
                                        fillAmount = minSizeRequirement.mul(3)
                                        maxOracleSpreadRatio = await amm.maxOracleSpreadRatio()
                                        oraclePrice = await amm.getUnderlyingPrice()
                                        lowerBoundPrice = oraclePrice.mul(_1e6.sub(maxOracleSpreadRatio)).div(_1e6)
                                        upperBoundPrice = oraclePrice.mul(_1e6.add(maxOracleSpreadRatio)).div(_1e6)

                                        await addMargin(bob, initialMargin)
                                        shortOrderPrice = lowerBoundPrice.sub(1)
                                        shortOrder = getOrder(market, bob.address, shortOrderBaseAssetQuantity, shortOrderPrice, getRandomSalt(), false)
                                        await placeOrderFromLimitOrder(shortOrder, bob)
                                        await addMargin(alice, initialMargin)
                                        longOrderPrice = upperBoundPrice
                                        longOrder = getOrder(market, alice.address, longOrderBaseAssetQuantity, longOrderPrice, getRandomSalt(), false)
                                        await placeOrderFromLimitOrder(longOrder, alice)
                                        
                                        response = await juror.validateOrdersAndDetermineFillPrice([encodeLimitOrderWithType(longOrder), encodeLimitOrderWithType(shortOrder)], fillAmount)
                                        // cleanup
                                        await cancelOrderFromLimitOrder(longOrder, alice)
                                        await cancelOrderFromLimitOrder(shortOrder, bob)
                                        await removeAllAvailableMargin(alice)
                                        await removeAllAvailableMargin(bob)

                                        expect(response.fillPrice.toString()).to.equal(lowerBoundPrice.toString())
                                        expect(response.instructions.length).to.equal(2)
                                        //longOrder
                                        expect(response.instructions[0].ammIndex.toNumber()).to.equal(0)
                                        expect(response.instructions[0].trader).to.equal(alice.address)
                                        expect(response.instructions[0].mode).to.equal(0)
                                        longOrderHash = await orderBook.getOrderHash(longOrder)
                                        expect(response.instructions[0].orderHash).to.equal(longOrderHash)

                                        //shortOrder
                                        expect(response.instructions[1].ammIndex.toNumber()).to.equal(0)
                                        expect(response.instructions[1].trader).to.equal(bob.address)
                                        expect(response.instructions[1].mode).to.equal(1)
                                        shortOrderHash = await orderBook.getOrderHash(shortOrder)
                                        expect(response.instructions[1].orderHash).to.equal(shortOrderHash)

                                        expect(response.orderTypes.length).to.equal(2)
                                        expect(response.orderTypes[0]).to.equal(0)
                                        expect(response.orderTypes[1]).to.equal(0)

                                        expect(response.encodedOrders[0]).to.equal(encodeLimitOrder(longOrder))
                                        expect(response.encodedOrders[1]).to.equal(encodeLimitOrder(shortOrder))
                                    });
                                });
                            })
                        });
                    })
                })
            })
        })
    })
});
