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
    getAMMContract,
    getMakerFee,
    getTakerFee,
    hubblebibliophile,
    multiplyPrice,
    multiplySize,
    placeOrder,
    removeAllAvailableMargin,
    waitForOrdersToMatch
} = utils

// Testing hubblebibliophile precompile contract

describe('Testing getNotionalPositionAndMargin and getPositionSizesAndUpperBoundsForMarkets',async function () {
    aliceInitialMargin = multiplyPrice(BigNumber.from(600000))
    charlieInitialMargin = multiplyPrice(BigNumber.from(600000))
    aliceOrderPrice = multiplyPrice(1800)
    charlieOrderPrice = multiplyPrice(1800)
    aliceOrderSize = multiplySize(0.1)
    charlieOrderSize = multiplySize(-0.1)
    market = BigNumber.from(0)

    context('When position and margin are 0', async function () {
        it('should returns the upperBound, 0 as positions, 0 as notionalPosition and 0 as margin', async function () {
            await removeAllAvailableMargin(alice)
            result = await hubblebibliophile.getPositionSizesAndUpperBoundsForMarkets(alice.address)
            expect(result.posSizes[0].toString()).to.equal("0")
            expectedUpperBound = await getUpperBoundForMarket(market)
            expect(result.upperBounds[0].toString()).to.equal(expectedUpperBound.toString())

            result = await hubblebibliophile.getNotionalPositionAndMargin(alice.address, false, 0)
            expect(result.notionalPosition.toString()).to.equal("0")
            expect(result.margin.toString()).to.equal("0")

        })
    })

    context('When position is zero but margin is non zero', async function () {
        context("when user never opened a position", async function () {
            this.afterAll(async function () {
                await removeAllAvailableMargin(alice)
            })
            it('should returns the upperBound, 0 as position, 0 as notionalPosition and amount deposited as margin for trader', async function () {
                await addMargin(alice, aliceInitialMargin)

                result = await hubblebibliophile.getPositionSizesAndUpperBoundsForMarkets(alice.address)
                expect(result.posSizes[0].toString()).to.equal("0")
                expectedUpperBound = await getUpperBoundForMarket(market)
                expect(result.upperBounds[0].toString()).to.equal(expectedUpperBound.toString())

                result = await hubblebibliophile.getNotionalPositionAndMargin(alice.address, false, 0)
                expect(result.notionalPosition.toString()).to.equal("0")
                expect(result.margin.toString()).to.equal(aliceInitialMargin.toString())
            })
        })
        context('when user closes whole position', async function () {
            this.afterAll(async function () {
                await removeAllAvailableMargin(alice)
                await removeAllAvailableMargin(charlie)
            })

            it('returns the upperBound, 0 as positions, 0 as position and amountDeposited - ordersFee as margin', async function () {
                await addMargin(alice, aliceInitialMargin)
                await addMargin(charlie, charlieInitialMargin)
                //create position
                await placeOrder(market, alice, aliceOrderSize, aliceOrderPrice)
                await placeOrder(market, charlie, charlieOrderSize, charlieOrderPrice)
                await waitForOrdersToMatch()
                // close position; charlie is taker for 2nd order
                await placeOrder(market, alice, charlieOrderSize, aliceOrderPrice)
                await placeOrder(market, charlie, aliceOrderSize, charlieOrderPrice)
                await waitForOrdersToMatch()
                makerFee = await getMakerFee() 
                takerFee = await getTakerFee() 

                expectedUpperBound = await getUpperBoundForMarket(market)
                result = await hubblebibliophile.getPositionSizesAndUpperBoundsForMarkets(alice.address)
                expect(result.posSizes[0].toString()).to.equal("0")
                expect(result.upperBounds[0].toString()).to.equal(expectedUpperBound.toString())
                result = await hubblebibliophile.getPositionSizesAndUpperBoundsForMarkets(charlie.address)
                expect(result.posSizes[0].toString()).to.equal("0")
                expect(result.upperBounds[0].toString()).to.equal(expectedUpperBound.toString())

                resultCharlie = await hubblebibliophile.getNotionalPositionAndMargin(charlie.address, false, 0)
                charlieOrder1Fee = makerFee.mul(charlieOrderSize.abs()).mul(charlieOrderPrice).div(_1e18).div(_1e6)
                charlieOrder2Fee = takerFee.mul(charlieOrderSize.abs()).mul(charlieOrderPrice).div(_1e18).div(_1e6)
                expectedCharlieMargin = charlieInitialMargin.sub(charlieOrder1Fee).sub(charlieOrder2Fee)
                expect(resultCharlie.notionalPosition.toString()).to.equal("0")
                expect(resultCharlie.margin.toString()).to.equal(expectedCharlieMargin.toString())

                resultAlice = await hubblebibliophile.getNotionalPositionAndMargin(alice.address, false, 0)
                aliceOrder1Fee = takerFee.mul(aliceOrderSize.abs()).mul(aliceOrderPrice).div(_1e18).div(_1e6)
                aliceOrder2Fee = makerFee.mul(aliceOrderSize.abs()).mul(aliceOrderPrice).div(_1e18).div(_1e6)
                expectedAliceMargin = aliceInitialMargin.sub(aliceOrder1Fee).sub(aliceOrder2Fee)
                expect(resultAlice.notionalPosition.toString()).to.equal("0")
                expect(resultAlice.margin.toString()).to.equal(expectedAliceMargin.toString())
            })
        })
    })

    context('When position and margin are both non zero', async function () {
        this.beforeEach(async function () {
            await addMargin(alice, aliceInitialMargin)
            await addMargin(charlie, charlieInitialMargin)
            // charlie places a short order and alice places a long order
            await placeOrder(market, charlie, charlieOrderSize, charlieOrderPrice)
            await placeOrder(market, alice, aliceOrderSize, aliceOrderPrice)
            await waitForOrdersToMatch()
        })

        this.afterEach(async function () {
            // charlie places a long order and alice places a short order
            await placeOrder(market, charlie, aliceOrderSize, aliceOrderPrice)
            await placeOrder(market, alice, charlieOrderSize, charlieOrderPrice)
            await waitForOrdersToMatch()
            await removeAllAvailableMargin(alice)
            await removeAllAvailableMargin(charlie)
        })

        context('when user creates a position', async function () {
            it('returns the positions, notional position and margin', async function () {
                expectedUpperBound = await getUpperBoundForMarket(market)
                result = await hubblebibliophile.getPositionSizesAndUpperBoundsForMarkets(alice.address)
                expect(result.posSizes[0].toString()).to.equal(aliceOrderSize.toString())
                expect(result.upperBounds[0].toString()).to.equal(expectedUpperBound.toString())
                result = await hubblebibliophile.getPositionSizesAndUpperBoundsForMarkets(charlie.address)
                expect(result.posSizes[0].toString()).to.equal(charlieOrderSize.toString())
                expect(result.upperBounds[0].toString()).to.equal(expectedUpperBound.toString())

                resultCharlie = await hubblebibliophile.getNotionalPositionAndMargin(charlie.address, false, 0)
                takerFee = await clearingHouse.takerFee() // in 1e6 units
                charlieOrderFee = takerFee.mul(charlieOrderSize.abs()).mul(charlieOrderPrice).div(_1e18).div(_1e6)
                expectedCharlieMargin = charlieInitialMargin.sub(charlieOrderFee)
                expectedCharlieNotionalPosition = charlieOrderSize.abs().mul(charlieOrderPrice).div(_1e18)
                expect(resultCharlie.notionalPosition.toString()).to.equal(expectedCharlieNotionalPosition.toString())
                expect(resultCharlie.margin.toString()).to.equal(expectedCharlieMargin.toString())

                resultAlice = await hubblebibliophile.getNotionalPositionAndMargin(alice.address, false, 0)
                aliceOrderFee = takerFee.mul(aliceOrderSize).mul(aliceOrderPrice).div(_1e18).div(_1e6)
                expectedAliceMargin = aliceInitialMargin.sub(aliceOrderFee)
                expectedAliceNotionalPosition = aliceOrderSize.mul(aliceOrderPrice).div(_1e18)
                expect(resultAlice.notionalPosition.toString()).to.equal(expectedAliceNotionalPosition.toString())
                expect(resultAlice.margin.toString()).to.equal(expectedAliceMargin.toString())
            })
        })

        context('when user increases the position', async function () {
            let aliceOrder2Size = multiplySize(0.2)
            let charlieOrder2Size = multiplySize(-0.2)

            this.afterEach(async function () {
                await placeOrder(market, charlie, aliceOrder2Size, charlieOrderPrice)
                await placeOrder(market, alice, charlieOrder2Size, aliceOrderPrice)
                await waitForOrdersToMatch()
            })
            it('returns the upperBound, positions, notional position and margin', async function () {
                // increase position , charlie is taker for 2nd order
                await placeOrder(market, alice, aliceOrder2Size, aliceOrderPrice)
                await placeOrder(market, charlie, charlieOrder2Size, charlieOrderPrice)
                await waitForOrdersToMatch()

                expectedUpperBound = await getUpperBoundForMarket(market)
                result = await hubblebibliophile.getPositionSizesAndUpperBoundsForMarkets(alice.address)
                totalAliceOrderSize = aliceOrderSize.add(aliceOrder2Size)
                expect(result.posSizes[0].toString()).to.equal(totalAliceOrderSize.toString())
                expect(result.upperBounds[0].toString()).to.equal(expectedUpperBound.toString())
                result = await hubblebibliophile.getPositionSizesAndUpperBoundsForMarkets(charlie.address)
                totalCharlieOrderSize = charlieOrderSize.add(charlieOrder2Size)
                expect(result.posSizes[0].toString()).to.equal(totalCharlieOrderSize.toString())
                expect(result.upperBounds[0].toString()).to.equal(expectedUpperBound.toString())

                makerFee = await getMakerFee() 
                takerFee = await getTakerFee() 

                // tests
                resultCharlie = await hubblebibliophile.getNotionalPositionAndMargin(charlie.address, false, 0)
                charlieOrder1Fee = makerFee.mul(charlieOrderSize.abs()).mul(charlieOrderPrice).div(_1e18).div(_1e6)
                charlieOrder2Fee = takerFee.mul(charlieOrder2Size.abs()).mul(charlieOrderPrice).div(_1e18).div(_1e6)
                expectedCharlieMargin = charlieInitialMargin.sub(charlieOrder1Fee).sub(charlieOrder2Fee)
                charlieOrder1Notional = charlieOrderSize.mul(charlieOrderPrice).div(_1e18).abs()
                charlieOrder2Notional = charlieOrder2Size.mul(charlieOrderPrice).div(_1e18).abs()
                expectedCharlieNotionalPosition = charlieOrder1Notional.add(charlieOrder2Notional) 
                expect(resultCharlie.notionalPosition.toString()).to.equal(expectedCharlieNotionalPosition.toString())
                expect(resultCharlie.margin.toString()).to.equal(expectedCharlieMargin.toString())

                resultAlice = await hubblebibliophile.getNotionalPositionAndMargin(alice.address, false, 0)
                aliceOrder1Fee = takerFee.mul(aliceOrderSize).mul(aliceOrderPrice).div(_1e18).div(_1e6)
                aliceOrder2Fee = makerFee.mul(aliceOrder2Size).mul(aliceOrderPrice).div(_1e18).div(_1e6)
                expectedAliceMargin = aliceInitialMargin.sub(aliceOrder1Fee).sub(aliceOrder2Fee)
                aliceOrder1Notional = aliceOrderSize.mul(aliceOrderPrice).div(_1e18)
                aliceOrder2Notional = aliceOrder2Size.mul(aliceOrderPrice).div(_1e18)
                expectedAliceNotionalPosition = aliceOrder1Notional.add(aliceOrder2Notional)
                expect(resultAlice.notionalPosition.toString()).to.equal(expectedAliceNotionalPosition.toString())
                expect(resultAlice.margin.toString()).to.equal(expectedAliceMargin.toString())
            })
        })

        context('when user decreases the position', async function () {
            let aliceOrder2Size = multiplySize(-0.2)
            let charlieOrder2Size = multiplySize(0.2)

            this.afterEach(async function () {
                await placeOrder(market, charlie, aliceOrder2Size, charlieOrderPrice)
                await placeOrder(market, alice, charlieOrder2Size, aliceOrderPrice)
                await waitForOrdersToMatch()
            })
            it('returns the upperBound, position, notional position and margin', async function () {
                // increase position and charlie is taker for 2nd order
                await placeOrder(market, alice, aliceOrder2Size, aliceOrderPrice)
                await placeOrder(market, charlie, charlieOrder2Size, charlieOrderPrice)
                await waitForOrdersToMatch()

                expectedUpperBound = await getUpperBoundForMarket(market)
                result = await hubblebibliophile.getPositionSizesAndUpperBoundsForMarkets(alice.address)
                totalAliceOrderSize = aliceOrderSize.add(aliceOrder2Size)
                expect(result.posSizes[0].toString()).to.equal(totalAliceOrderSize.toString())
                expect(result.upperBounds[0].toString()).to.equal(expectedUpperBound.toString())
                result = await hubblebibliophile.getPositionSizesAndUpperBoundsForMarkets(charlie.address)
                totalCharlieOrderSize = charlieOrderSize.add(charlieOrder2Size)
                expect(result.posSizes[0].toString()).to.equal(totalCharlieOrderSize.toString())
                expect(result.upperBounds[0].toString()).to.equal(expectedUpperBound.toString())

                makerFee = await getMakerFee() 
                takerFee = await getTakerFee() 

                resultCharlie = await hubblebibliophile.getNotionalPositionAndMargin(charlie.address, false, 0)
                charlieOrder1Fee = makerFee.mul(charlieOrderSize.abs()).mul(charlieOrderPrice).div(_1e18).div(_1e6)
                charlieOrder2Fee = takerFee.mul(charlieOrder2Size.abs()).mul(charlieOrderPrice).div(_1e18).div(_1e6)
                expectedCharlieMargin = charlieInitialMargin.sub(charlieOrder1Fee).sub(charlieOrder2Fee)
                charlieOrder1Notional = charlieOrderSize.mul(charlieOrderPrice).div(_1e18)
                charlieOrder2Notional = charlieOrder2Size.mul(charlieOrderPrice).div(_1e18)
                expectedNotionalPosition = charlieOrder1Notional.add(charlieOrder2Notional).abs()
                expect(resultCharlie.notionalPosition.toString()).to.equal(expectedNotionalPosition.toString())
                expect(resultCharlie.margin.toString()).to.equal(expectedCharlieMargin.toString())

                resultAlice = await hubblebibliophile.getNotionalPositionAndMargin(alice.address, false, 0)
                aliceOrder1Fee = takerFee.mul(aliceOrderSize.abs()).mul(aliceOrderPrice).div(_1e18).div(_1e6)
                aliceOrder2Fee = makerFee.mul(aliceOrder2Size.abs()).mul(aliceOrderPrice).div(_1e18).div(_1e6)
                expectedAliceMargin = aliceInitialMargin.sub(aliceOrder1Fee).sub(aliceOrder2Fee)
                aliceOrder1Notional = aliceOrderSize.mul(aliceOrderPrice).div(_1e18)
                aliceOrder2Notional = aliceOrder2Size.mul(aliceOrderPrice).div(_1e18)
                expectedNotionalPosition = aliceOrder1Notional.add(aliceOrder2Notional).abs()
                expect(resultAlice.notionalPosition.toString()).to.equal(expectedNotionalPosition.toString())
                expect(resultAlice.margin.toString()).to.equal(expectedAliceMargin.toString())
            })
        })
    })
})


async function getUpperBoundForMarket(market) {
    amm = await getAMMContract(market)
    maxOraclePriceSpread = await amm.maxOracleSpreadRatio()
    underlyingPrice = await amm.getUnderlyingPrice()
    upperBound = underlyingPrice.mul(_1e6.add(maxOraclePriceSpread)).div(_1e6)
    return upperBound
}

