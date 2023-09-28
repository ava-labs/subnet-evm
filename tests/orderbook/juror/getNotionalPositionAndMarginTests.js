const { BigNumber } = require('ethers');
const { expect } = require('chai');

const utils = require('../utils')

const {
    _1e6,
    _1e18,
    addMargin,
    alice,
    charlie,
    clearingHouse,
    getOrderV2,
    getMakerFee,
    getRandomSalt,
    getTakerFee,
    juror,
    multiplyPrice,
    multiplySize,
    placeOrder,
    placeOrderFromLimitOrderV2,
    removeAllAvailableMargin,
    waitForOrdersToMatch
} = utils

// Testing juror precompile contract

describe('Testing getNotionalPositionAndMargin',async function () {
    aliceInitialMargin = multiplyPrice(BigNumber.from(600000))
    charlieInitialMargin = multiplyPrice(BigNumber.from(600000))
    aliceOrderPrice = multiplyPrice(1800)
    charlieOrderPrice = multiplyPrice(1800)
    aliceOrderSize = multiplySize(0.1)
    charlieOrderSize = multiplySize(-0.1)
    market = BigNumber.from(0)

    context('When position and margin are 0', async function () {
        it('should return 0 as notionalPosition and 0 as margin', async function () {
            await removeAllAvailableMargin(alice)
            result = await juror.getNotionalPositionAndMargin(alice.address, false, 0)
            expect(result.notionalPosition.toString()).to.equal("0")
            expect(result.margin.toString()).to.equal("0")
        })
    })

    context('When position is zero but margin is non zero', async function () {
        context("when user never opened a position", async function () {
            this.afterAll(async function () {
                await removeAllAvailableMargin(alice)
            })
            it('should return 0 as notionalPosition and amount deposited as margin for trader', async function () {
                await addMargin(alice, aliceInitialMargin)

                result = await juror.getNotionalPositionAndMargin(alice.address, false, 0)
                expect(result.notionalPosition.toString()).to.equal("0")
                expect(result.margin.toString()).to.equal(aliceInitialMargin.toString())
            })
        })
        context('when user opens and closes whole position', async function () {
            this.afterAll(async function () {
                await removeAllAvailableMargin(alice)
                await removeAllAvailableMargin(charlie)
            })

            it('returns 0 as position and amountDeposited - ordersFee as margin', async function () {
                await addMargin(alice, aliceInitialMargin)
                await addMargin(charlie, charlieInitialMargin)
                //create position

                longOrder = getOrderV2(market, alice.address, aliceOrderSize, aliceOrderPrice, getRandomSalt())
                await placeOrderFromLimitOrderV2(longOrder, alice)
                shortOrder = getOrderV2(market, charlie.address, charlieOrderSize, charlieOrderPrice, getRandomSalt())
                await placeOrderFromLimitOrderV2(shortOrder, charlie)
                await waitForOrdersToMatch()
                // close position; charlie is taker for 2nd order
                oppositeLongOrder = getOrderV2(market, charlie.address, aliceOrderSize, aliceOrderPrice, getRandomSalt())
                await placeOrderFromLimitOrderV2(oppositeLongOrder, charlie)
                oppositeShortOrder = getOrderV2(market, alice.address, charlieOrderSize, charlieOrderPrice, getRandomSalt())
                await placeOrderFromLimitOrderV2(oppositeShortOrder, alice)
                await waitForOrdersToMatch()

                makerFee = await getMakerFee() 
                takerFee = await getTakerFee() 

                resultCharlie = await juror.getNotionalPositionAndMargin(charlie.address, false, 0)
                charlieOrder1Fee = makerFee.mul(charlieOrderSize.abs()).mul(charlieOrderPrice).div(_1e18).div(_1e6)
                charlieOrder2Fee = takerFee.mul(charlieOrderSize.abs()).mul(charlieOrderPrice).div(_1e18).div(_1e6)
                expectedCharlieMargin = charlieInitialMargin.sub(charlieOrder1Fee).sub(charlieOrder2Fee)
                expect(resultCharlie.notionalPosition.toString()).to.equal("0")
                expect(resultCharlie.margin.toString()).to.equal(expectedCharlieMargin.toString())

                resultAlice = await juror.getNotionalPositionAndMargin(alice.address, false, 0)
                aliceOrder1Fee = takerFee.mul(aliceOrderSize.abs()).mul(aliceOrderPrice).div(_1e18).div(_1e6)
                aliceOrder2Fee = makerFee.mul(aliceOrderSize.abs()).mul(aliceOrderPrice).div(_1e18).div(_1e6)
                expectedAliceMargin = aliceInitialMargin.sub(aliceOrder1Fee).sub(aliceOrder2Fee)
                expect(resultAlice.notionalPosition.toString()).to.equal("0")
                expect(resultAlice.margin.toString()).to.equal(expectedAliceMargin.toString())
            })
        })
    })

    context('When position and margin are both non zero', async function () {
        //create position
        let aliceOrder1 = getOrderV2(market, alice.address, aliceOrderSize, aliceOrderPrice, getRandomSalt())
        let charlieOrder1 = getOrderV2(market, charlie.address, charlieOrderSize, charlieOrderPrice, getRandomSalt())
        let oppositeAliceOrder1 = getOrderV2(market, alice.address, charlieOrderSize, charlieOrderPrice, getRandomSalt())
        let oppositeCharlieOrder1 = getOrderV2(market, charlie.address, aliceOrderSize, aliceOrderPrice, getRandomSalt())
        // increase position
        let aliceOrder2Size = multiplySize(0.2)
        let charlieOrder2Size = multiplySize(-0.2)
        let aliceOrder2 = getOrderV2(market, alice.address, aliceOrder2Size, aliceOrderPrice, getRandomSalt())
        let charlieOrder2 = getOrderV2(market, charlie.address, charlieOrder2Size, charlieOrderPrice, getRandomSalt())
        // decrease position
        let aliceOrder3Size = multiplySize(-0.4)
        let charlieOrder3Size = multiplySize(0.4)
        let aliceOrder3 = getOrderV2(market, alice.address, aliceOrder3Size, aliceOrderPrice, getRandomSalt())
        let charlieOrder3 = getOrderV2(market, charlie.address, charlieOrder3Size, charlieOrderPrice, getRandomSalt())

        let makerFee, takerFee

        this.beforeAll(async function () {
            makerFee = await getMakerFee() 
            takerFee = await getTakerFee() 
            await addMargin(alice, aliceInitialMargin)
            await addMargin(charlie, charlieInitialMargin)
            // charlie places a short order and alice places a long order
            await placeOrderFromLimitOrderV2(aliceOrder1, alice)
            await placeOrderFromLimitOrderV2(charlieOrder1, charlie)
            await waitForOrdersToMatch()
        })

        this.afterAll(async function () {
            let resultCharlie = await juror.getNotionalPositionAndMargin(charlie.address, false, 0)
            let resultAlice = await juror.getNotionalPositionAndMargin(alice.address, false, 0)
            // charlie places a long order and alice places a short order
            charlieTotalSize = charlieOrder1.baseAssetQuantity.add(charlieOrder2Size).add(charlieOrder3Size)
            aliceTotalSize = aliceOrder1.baseAssetQuantity.add(aliceOrder2Size).add(aliceOrder3Size)
            aliceCleanupOrder = getOrderV2(market, alice.address, charlieTotalSize, charlieOrderPrice, getRandomSalt())
            charlieCleanupOrder = getOrderV2(market, charlie.address, aliceTotalSize, aliceOrderPrice, getRandomSalt())
            aliceCleanupOrderMargin = await utils.getRequiredMarginForShortOrder(aliceCleanupOrder)
            charlieCleanupOrderMargin = await utils.getRequiredMarginForShortOrder(charlieCleanupOrder)
            await addMargin(alice, aliceCleanupOrderMargin)
            await addMargin(charlie, charlieCleanupOrderMargin)
            await placeOrderFromLimitOrderV2(aliceCleanupOrder, alice)
            await placeOrderFromLimitOrderV2(charlieCleanupOrder, charlie)
            await waitForOrdersToMatch()
            await removeAllAvailableMargin(alice)
            await removeAllAvailableMargin(charlie)
        })

        context('when user creates a position', async function () {
            it('should return correct notional position and margin', async function () {
                let resultCharlie = await juror.getNotionalPositionAndMargin(charlie.address, false, 0)
                let charlieOrderFee = takerFee.mul(charlieOrderSize.abs()).mul(charlieOrderPrice).div(_1e18).div(_1e6)
                let expectedCharlieMargin = charlieInitialMargin.sub(charlieOrderFee)
                let expectedCharlieNotionalPosition = charlieOrderSize.abs().mul(charlieOrderPrice).div(_1e18)
                expect(resultCharlie.notionalPosition.toString()).to.equal(expectedCharlieNotionalPosition.toString())
                expect(resultCharlie.margin.toString()).to.equal(expectedCharlieMargin.toString())

                let resultAlice = await juror.getNotionalPositionAndMargin(alice.address, false, 0)
                let aliceOrderFee = takerFee.mul(aliceOrderSize).mul(aliceOrderPrice).div(_1e18).div(_1e6)
                let expectedAliceMargin = aliceInitialMargin.sub(aliceOrderFee)
                let expectedAliceNotionalPosition = aliceOrderSize.mul(aliceOrderPrice).div(_1e18)
                expect(resultAlice.notionalPosition.toString()).to.equal(expectedAliceNotionalPosition.toString())
                expect(resultAlice.margin.toString()).to.equal(expectedAliceMargin.toString())
            })
        })

        context('when user increases the position', async function () {
            it('should return increased notional position and correct margin', async function () {
                // increase position , charlie is taker for 2nd order
                await placeOrderFromLimitOrderV2(aliceOrder2, alice)
                await placeOrderFromLimitOrderV2(charlieOrder2, charlie)
                await waitForOrdersToMatch()
                // tests
                let resultCharlie = await juror.getNotionalPositionAndMargin(charlie.address, false, 0)
                let charlieOrder1Fee = makerFee.mul(charlieOrderSize.abs()).mul(charlieOrderPrice).div(_1e18).div(_1e6)
                let charlieOrder2Fee = takerFee.mul(charlieOrder2Size.abs()).mul(charlieOrderPrice).div(_1e18).div(_1e6)
                let expectedCharlieMargin = charlieInitialMargin.sub(charlieOrder1Fee).sub(charlieOrder2Fee)
                let charlieOrder1Notional = charlieOrderSize.mul(charlieOrderPrice).div(_1e18).abs()
                let charlieOrder2Notional = charlieOrder2Size.mul(charlieOrderPrice).div(_1e18).abs()
                let expectedCharlieNotionalPosition = charlieOrder1Notional.add(charlieOrder2Notional) 
                expect(resultCharlie.notionalPosition.toString()).to.equal(expectedCharlieNotionalPosition.toString())
                expect(resultCharlie.margin.toString()).to.equal(expectedCharlieMargin.toString())

                let resultAlice = await juror.getNotionalPositionAndMargin(alice.address, false, 0)
                let aliceOrder1Fee = takerFee.mul(aliceOrderSize).mul(aliceOrderPrice).div(_1e18).div(_1e6)
                let aliceOrder2Fee = makerFee.mul(aliceOrder2Size).mul(aliceOrderPrice).div(_1e18).div(_1e6)
                let expectedAliceMargin = aliceInitialMargin.sub(aliceOrder1Fee).sub(aliceOrder2Fee)
                let aliceOrder1Notional = aliceOrderSize.mul(aliceOrderPrice).div(_1e18)
                let aliceOrder2Notional = aliceOrder2Size.mul(aliceOrderPrice).div(_1e18)
                let expectedAliceNotionalPosition = aliceOrder1Notional.add(aliceOrder2Notional)
                expect(resultAlice.notionalPosition.toString()).to.equal(expectedAliceNotionalPosition.toString())
                expect(resultAlice.margin.toString()).to.equal(expectedAliceMargin.toString())
            })
        })

        context('when user decreases the position', async function () {
            it('should returns decreased notional position and margin', async function () {
                // increase position and charlie is maker for 3rd order
                await placeOrderFromLimitOrderV2(charlieOrder3, charlie)
                await placeOrderFromLimitOrderV2(aliceOrder3, alice)
                await waitForOrdersToMatch()
                let resultCharlie = await juror.getNotionalPositionAndMargin(charlie.address, false, 0)
                let charlieOrder1Fee = makerFee.mul(charlieOrderSize.abs()).mul(charlieOrderPrice).div(_1e18).div(_1e6)
                let charlieOrder2Fee = takerFee.mul(charlieOrder2Size.abs()).mul(charlieOrderPrice).div(_1e18).div(_1e6)
                let charlieOrder3Fee = makerFee.mul(charlieOrder3Size.abs()).mul(charlieOrderPrice).div(_1e18).div(_1e6)
                let expectedCharlieMargin = charlieInitialMargin.sub(charlieOrder1Fee).sub(charlieOrder2Fee).sub(charlieOrder3Fee)
                let charlieOrder1Notional = charlieOrderSize.mul(charlieOrderPrice).div(_1e18)
                let charlieOrder2Notional = charlieOrder2Size.mul(charlieOrderPrice).div(_1e18)
                let charlieOrder3Notional = charlieOrder3Size.mul(charlieOrderPrice).div(_1e18)
                let expectedCharlieNotionalPosition = charlieOrder1Notional.add(charlieOrder2Notional).add(charlieOrder3Notional).abs()
                expect(resultCharlie.notionalPosition.toString()).to.equal(expectedCharlieNotionalPosition.toString())
                expect(resultCharlie.margin.toString()).to.equal(expectedCharlieMargin.toString())

                let resultAlice = await juror.getNotionalPositionAndMargin(alice.address, false, 0)
                let aliceOrder1Fee = takerFee.mul(aliceOrderSize.abs()).mul(aliceOrderPrice).div(_1e18).div(_1e6)
                let aliceOrder2Fee = makerFee.mul(aliceOrder2Size.abs()).mul(aliceOrderPrice).div(_1e18).div(_1e6)
                let aliceOrder3Fee = takerFee.mul(aliceOrder3Size.abs()).mul(aliceOrderPrice).div(_1e18).div(_1e6)
                let expectedAliceMargin = aliceInitialMargin.sub(aliceOrder1Fee).sub(aliceOrder2Fee).sub(aliceOrder3Fee)
                let aliceOrder1Notional = aliceOrderSize.mul(aliceOrderPrice).div(_1e18)
                let aliceOrder2Notional = aliceOrder2Size.mul(aliceOrderPrice).div(_1e18)
                let aliceOrder3Notional = aliceOrder3Size.mul(aliceOrderPrice).div(_1e18)
                let expectedAliceNotionalPosition = aliceOrder1Notional.add(aliceOrder2Notional).add(aliceOrder3Notional).abs()
                expect(resultAlice.notionalPosition.toString()).to.equal(expectedAliceNotionalPosition.toString())
                expect(resultAlice.margin.toString()).to.equal(expectedAliceMargin.toString())
            })
        })
    })
})
