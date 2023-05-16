
const { ethers } = require('ethers');
const { BigNumber } = require('ethers')
const axios = require('axios');
const { expect } = require('chai');
const { randomInt } = require('crypto');

const OrderBookContractAddress = "0x0300000000000000000000000000000000000069"
const MarginAccountContractAddress = "0x0300000000000000000000000000000000000070"
const MarginAccountHelperContractAddress = "0x610178dA211FEF7D417bC0e6FeD39F05609AD788"
const ClearingHouseContractAddress = "0x0300000000000000000000000000000000000071"

let provider, domain, orderType, orderBook, marginAccount, marginAccountHelper, clearingHouse
let alice, bob, charlie, aliceAddress, bobAddress, charlieAddress
let governance
let alicePartialMatchedLongOrder, bobHighPriceShortOrder

const ZERO = BigNumber.from(0)
const _1e6 = BigNumber.from(10).pow(6)
const _1e8 = BigNumber.from(10).pow(8)
const _1e12 = BigNumber.from(10).pow(12)
const _1e18 = ethers.constants.WeiPerEther
const maxLeverage = 5
const tradeFeeRatio = 0.0025

const homedir = require('os').homedir()
let conf = require(`${homedir}/.hubblenet.json`)
const url = `http://127.0.0.1:9650/ext/bc/${conf.chain_id}/rpc`

describe('Submit transaction and compare with EVM state', function () {
    before('', async function () {
        provider = new ethers.providers.JsonRpcProvider(url);

        // Set up signer
        governance = new ethers.Wallet('0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80', provider) // governance
        alice = new ethers.Wallet('0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d', provider); // 0x70997970c51812dc3a010c7d01b50e0d17dc79c8
        bob = new ethers.Wallet('0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a', provider); // 0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc
        charlie = new ethers.Wallet('15614556be13730e9e8d6eacc1603143e7b96987429df8726384c2ec4502ef6e', provider); // 0x55ee05df718f1a5c1441e76190eb1a19ee2c9430
        aliceAddress = alice.address.toLowerCase()
        bobAddress = bob.address.toLowerCase()
        charlieAddress = charlie.address.toLowerCase()
        console.log({ alice: aliceAddress, bob: bobAddress, charlie: charlieAddress });

        // Set up contract interface
        orderBook = new ethers.Contract(OrderBookContractAddress, require('./abi/OrderBook.json'), provider);
        clearingHouse = new ethers.Contract(ClearingHouseContractAddress, require('./abi/ClearingHouse.json'), provider);
        marginAccount = new ethers.Contract(MarginAccountContractAddress, require('./abi/MarginAccount.json'), provider);
        marginAccountHelper = new ethers.Contract(MarginAccountHelperContractAddress, require('./abi/MarginAccountHelper.json'));
        domain = {
            name: 'Hubble',
            version: '2.0',
            chainId: (await provider.getNetwork()).chainId,
            verifyingContract: orderBook.address
        }

        orderType = {
            Order: [
                // field ordering must be the same as LIMIT_ORDER_TYPEHASH
                { name: "ammIndex", type: "uint256" },
                { name: "trader", type: "address" },
                { name: "baseAssetQuantity", type: "int256" },
                { name: "price", type: "uint256" },
                { name: "salt", type: "uint256" },
                { name: "reduceOnly", type: "bool" },
            ]
        }

    })

    let aliceMargin = _1e6 * 150
    let bobMargin = _1e6 * 150
    let charlieMargin = 0

    let aliceOrderSize = 0.1
    let aliceOrderPrice = 1800
    let aliceReserved = getReservedMargin(aliceOrderSize * aliceOrderPrice)
    let aliceTradeFee = getTradeFee(aliceOrderSize * aliceOrderPrice)
    let aliceOpenNotional = Math.abs(aliceOrderSize * aliceOrderPrice) * 1e6
    let aliceLiquidationThreshold = getLiquidationThreshold(aliceOrderSize)


    let bobOrderSize = -0.1
    let bobOrderPrice = 1800
    let bobOpenNotional = Math.abs(bobOrderSize * bobOrderPrice) * 1e6
    let bobLiquidationThreshold = getLiquidationThreshold(bobOrderSize)
    let bobTradeFee = getTradeFee(bobOrderSize * bobOrderPrice)

    it('Add margin', async function () {
        tx = await addMargin(alice, aliceMargin)
        await tx.wait();

        tx = await addMargin(bob, bobMargin)
        await tx.wait();

        const expectedState = {
            "order_map": {},
            "trader_map": {
                [bobAddress]: {
                    "positions": {},
                    "margin": {
                        "reserved": 0,
                        "deposited": {
                            "0": bobMargin
                        }
                    }
                },
                [aliceAddress]: {
                    "positions": {},
                    "margin": {
                        "reserved": 0,
                        "deposited": {
                            "0": aliceMargin
                        }
                    }
                }
            },
            "next_funding_time": await getNextFundingTime(),
            "last_price": {
                "0": 0
            }
        }
        const evmState = await getEVMState()
        // console.log(JSON.stringify(evmState, null, 2))
        expect(evmState).to.deep.contain(expectedState)
    });

    it('Remove margin', async function () {
        const aliceMarginRemoved = _1e6 * 1
        const tx = await marginAccount.connect(alice).removeMargin(0, aliceMarginRemoved)
        aliceMargin = aliceMargin - aliceMarginRemoved
        await tx.wait();

        const expectedState = {
            "order_map": {},
            "trader_map": {
                [bobAddress]: {
                    "positions": {},
                    "margin": {
                        "reserved": 0,
                        "deposited": {
                            "0": bobMargin
                        }
                    }
                },
                [aliceAddress]: {
                    "positions": {},
                    "margin": {
                        "reserved": 0,
                        "deposited": {
                            "0": aliceMargin
                        }
                    }
                }
            },
            "next_funding_time": await getNextFundingTime(),
            "last_price": {
                "0": 0
            }
        }
        const evmState = await getEVMState()
        // console.log(JSON.stringify(evmState, null, 2))
        expect(evmState).to.deep.contain(expectedState)
    });

    it('Place order', async function () {
        const { hash, order, signature, tx, txReceipt } = await placeOrder(alice, aliceOrderSize, aliceOrderPrice, 101)

        const expectedState = {
            "order_map": {
                [hash]: {
                    "market": 0,
                    "position_type": "long",
                    "user_address": aliceAddress,
                    "base_asset_quantity": _1e18 * aliceOrderSize,
                    "filled_base_asset_quantity": 0,
                    "salt": 101,
                    "price": _1e6 * aliceOrderPrice,
                    "lifecycle_list": [
                        {
                            "BlockNumber": txReceipt.blockNumber,
                            "Status": 0
                        }
                    ],
                    "signature": signature,
                    "block_number": txReceipt.blockNumber,
                    "reduce_only": false
                }
            },
            "trader_map": {
                [bobAddress]: {
                    "positions": {},
                    "margin": {
                        "reserved": 0,
                        "deposited": {
                            "0": bobMargin
                        }
                    }
                },
                [aliceAddress]: {
                    "positions": {},
                    "margin": {
                        "reserved": aliceReserved,
                        "deposited": {
                            "0": aliceMargin
                        }
                    }
                }
            },
            "next_funding_time": await getNextFundingTime(),
            "last_price": {
                "0": 0
            }
        }

        const evmState = await getEVMState()
        // console.log(JSON.stringify(evmState, null, 2))
        expect(evmState).to.deep.contain(expectedState)
    });

    it('Match order', async function () {
        const {tx} = await placeOrder(bob, bobOrderSize, bobOrderPrice, 201)
        
        await sleep(2)
        
        const expectedState = {
            "order_map": {},
            "trader_map": {
                [bobAddress]: {
                    "positions": {
                        "0": {
                            "open_notional": bobOpenNotional,
                            "size": _1e18 * bobOrderSize,
                            "unrealised_funding": 0,
                            "last_premium_fraction": 0,
                            "liquidation_threshold": bobLiquidationThreshold
                        }
                    },
                    "margin": {
                        "reserved": 0,
                        "deposited": {
                            "0": bobMargin - bobTradeFee
                        }
                    }
                },
                [aliceAddress]: {
                    "positions": {
                        "0": {
                            "open_notional": aliceOpenNotional,
                            "size": _1e18 * aliceOrderSize,
                            "unrealised_funding": 0,
                            "last_premium_fraction": 0,
                            "liquidation_threshold": aliceLiquidationThreshold
                        }
                    },
                    "margin": {
                        "reserved": 0,
                        "deposited": {
                            "0": aliceMargin - aliceTradeFee
                        }
                    }
                }
            },
            "next_funding_time": await getNextFundingTime(),
            "last_price": {
                "0": _1e6 * aliceOrderPrice
            }
        }
        const evmState = await getEVMState()
        // console.log(JSON.stringify(evmState, null, 2))
        expect(evmState).to.deep.contain(expectedState)
    });

    it('Order cancel', async function () {
        const { hash, order } = await placeOrder(alice, 0.1, 1800, 401)

        tx = await orderBook.connect(alice).cancelOrder(hash)
        await tx.wait()

        // same as last test scenario
        const expectedState = {
            "order_map": {},
            "trader_map": {
                [bobAddress]: {
                    "positions": {
                        "0": {
                            "open_notional": bobOpenNotional,
                            "size": _1e18 * bobOrderSize,
                            "unrealised_funding": 0,
                            "last_premium_fraction": 0,
                            "liquidation_threshold": bobLiquidationThreshold
                        }
                    },
                    "margin": {
                        "reserved": 0,
                        "deposited": {
                            "0": bobMargin - bobTradeFee
                        }
                    }
                },
                [aliceAddress]: {
                    "positions": {
                        "0": {
                            "open_notional": aliceOpenNotional,
                            "size": _1e18 * aliceOrderSize,
                            "unrealised_funding": 0,
                            "last_premium_fraction": 0,
                            "liquidation_threshold": aliceLiquidationThreshold
                        }
                    },
                    "margin": {
                        "reserved": 0,
                        "deposited": {
                            "0": aliceMargin - aliceTradeFee
                        }
                    }
                }
            },
            "next_funding_time": await getNextFundingTime(),
            "last_price": {
                "0": _1e6 * aliceOrderPrice
            }
        }

        const evmState = await getEVMState()
        // console.log(JSON.stringify(evmState, null, 2))
        expect(evmState).to.deep.contain(expectedState)
    });

    it('Partially match order', async function () {
        alicePartialMatchedLongOrder = await placeOrder(alice, 0.2, aliceOrderPrice, 301)
        const { tx, hash: bobShortOrderHash } = await placeOrder(bob, -0.1, bobOrderPrice, 302)

        await tx.wait()

        aliceOrderSize = aliceOrderSize + 0.1
        bobOrderSize = bobOrderSize - 0.1
        bobOpenNotional = Math.abs(bobOrderSize * bobOrderPrice) * 1e6
        aliceOpenNotional = Math.abs(aliceOrderSize * aliceOrderPrice) * 1e6
        bobLiquidationThreshold = getLiquidationThreshold(bobOrderSize)
        aliceLiquidationThreshold = getLiquidationThreshold(aliceOrderSize)

        aliceTradeFee = getTradeFee(0.2 * aliceOrderPrice)
        bobTradeFee = getTradeFee(0.2 * bobOrderPrice)
        aliceReserved = getReservedMargin(0.1 * aliceOrderPrice) // reserved only for 0.1 size
        const expectedState = {
            "order_map": {
                [alicePartialMatchedLongOrder.hash]: {
                    "market": 0,
                    "position_type": "long",
                    "user_address": aliceAddress,
                    "base_asset_quantity": 200000000000000000,
                    "filled_base_asset_quantity": 100000000000000000,
                    "salt": 301,
                    "price": 1800000000,
                    "lifecycle_list": [
                        {
                            "BlockNumber": alicePartialMatchedLongOrder.txReceipt.blockNumber,
                            "Status": 0
                        }
                    ],
                    "signature": alicePartialMatchedLongOrder.signature,
                    "block_number": alicePartialMatchedLongOrder.txReceipt.blockNumber,
                    "reduce_only": false
                }
            },
            "trader_map": {
                [bobAddress]: {
                    "positions": {
                        "0": {
                            "open_notional": bobOpenNotional,
                            "size": _1e18 * bobOrderSize,
                            "unrealised_funding": 0,
                            "last_premium_fraction": 0,
                            "liquidation_threshold": bobLiquidationThreshold
                        }
                    },
                    "margin": {
                        "reserved": 0,
                        "deposited": {
                            "0": bobMargin - bobTradeFee
                        }
                    }
                },
                [aliceAddress]: {
                    "positions": {
                        "0": {
                            "open_notional": aliceOpenNotional,
                            "size": _1e18 * aliceOrderSize,
                            "unrealised_funding": 0,
                            "last_premium_fraction": 0,
                            "liquidation_threshold": aliceLiquidationThreshold
                        }
                    },
                    "margin": {
                        "reserved": aliceReserved,
                        "deposited": {
                            "0": aliceMargin - aliceTradeFee
                        }
                    }
                }
            },
            "next_funding_time": await getNextFundingTime(),
            "last_price": {
                "0": _1e6 * aliceOrderPrice
            }
        }
        const evmState = await getEVMState()
        // console.log(JSON.stringify(evmState, null, 2))
        expect(evmState).to.deep.contain(expectedState)
    });


    // it.skip('Order match error', async function () {
    //     // place an order with reduceOnly which should fail
    //     // const { hash: charlieHash } = await placeOrder(charlie, 50, 12, 501)
    //     await orderBook.connect(alice).cancelOrder(alicePartialMatchedLongOrder.hash)
    //     bobReverseOrder = await placeOrder(bob, 6.95, 10, 501)
    //     aliceReverseOrder = await placeOrder(alice, -7, 10, 502)
    //     // bobHighPriceShortOrder = await placeOrder(bob, -2, 10, 502) // reduceOnly; this should fail while matching

    //     await sleep(2)
    //     const expectedState = {
    //         "order_map": {
    //             [aliceReverseOrder.hash]: {
    //                 "market": 0,
    //                 "position_type": "short",
    //                 "user_address": aliceAddress,
    //                 "base_asset_quantity": -7000000000000000000,
    //                 "filled_base_asset_quantity": -6950000000000000000,
    //                 "salt": 502,
    //                 "price": 10000000,
    //                 "lifecycle_list": [
    //                     {
    //                         "BlockNumber": aliceReverseOrder.txReceipt.blockNumber,
    //                         "Status": 0
    //                     }
    //                 ],
    //                 "signature": aliceReverseOrder.signature,
    //                 "block_number": aliceReverseOrder.txReceipt.blockNumber,
    //                 "reduce_only": false
    //             }
    //         },
    //         "trader_map": {
    //             [bobAddress]: {
    //                 "positions": {
    //                     "0": {
    //                         "open_notional": 70000000,
    //                         "size": -7000000000000000000,
    //                         "unrealised_funding": 0,
    //                         "last_premium_fraction": 0,
    //                         "liquidation_threshold": -5000000000000000000
    //                     }
    //                 },
    //                 "margin": {
    //                     "reserved": 0,
    //                     "deposited": {
    //                         "0": 39965000
    //                     }
    //                 }
    //             },
    //             [aliceAddress]: {
    //                 "positions": {
    //                     "0": {
    //                         "open_notional": 70000000,
    //                         "size": 7000000000000000000,
    //                         "unrealised_funding": 0,
    //                         "last_premium_fraction": 0,
    //                         "liquidation_threshold": 5000000000000000000
    //                     }
    //                 },
    //                 "margin": {
    //                     "reserved": 0,
    //                     "deposited": {
    //                         "0": 29965000
    //                     }
    //                 }
    //             }
    //         },
    //         "next_funding_time": await getNextFundingTime(),
    //         "last_price": {
    //             "0": 10000000
    //         }
    //     }
    //     const evmState = await getEVMState()
    //     console.log(JSON.stringify(evmState, null, 2))
    //     expect(evmState).to.deep.contain(expectedState)
    // });

    it('Liquidate trader', async function () {
        await addMargin(charlie, _1e6.mul(100))
        await addMargin(alice, _1e6.mul(300))
        await addMargin(bob, _1e6.mul(200))

        await orderBook.connect(alice).cancelOrder(alicePartialMatchedLongOrder.hash)
        
        aliceMargin = aliceMargin + (_1e6 * 300)
        bobMargin = bobMargin + (_1e6 * 200)
        charlieMargin = _1e6 * 100

        // large position by charlie
        let charlieOrderSize = 0.25
        let charliePrice = 1800
        await placeOrder(charlie, charlieOrderSize, charliePrice, 601)
        await placeOrder(bob, -charlieOrderSize, charliePrice, 602)

        bobOrderSize -= charlieOrderSize
        bobOpenNotional = bobOpenNotional + Math.abs(charlieOrderSize * charliePrice * _1e6)
        let charlieOpenNotional = Math.abs(charlieOrderSize * charliePrice * _1e6)

        // reduce the price
        let reducedPrice = 1400
        await setOraclePrice(0, reducedPrice * _1e6)
        const { hash: aliceHash } = await placeOrder(alice, 0.01, reducedPrice, 603)
        const { hash: bobHash2 } = await placeOrder(bob, -0.01, reducedPrice, 604)

        bobOpenNotional = bobOpenNotional + Math.abs(0.01 * reducedPrice * _1e6)
        aliceOpenNotional = aliceOpenNotional + Math.abs(0.01 * reducedPrice * _1e6)
        aliceOrderSize += 0.01
        bobOrderSize -= 0.01
        bobLiquidationThreshold = getLiquidationThreshold(bobOrderSize)
        aliceLiquidationThreshold = getLiquidationThreshold(aliceOrderSize)
        let charlieLiquidationThreshold = getLiquidationThreshold(charlieOrderSize)
        aliceTradeFee =  getTradeFee(aliceOpenNotional/_1e6)
        bobTradeFee = getTradeFee(bobOpenNotional/_1e6)
        let charlieTradeFee = getTradeFee(charlieOpenNotional/_1e6)

        // 1 long order so that liquidation can run
        // const aliceNewPrice = 1800
        // let increasedPrice = 1800
        // await setOraclePrice(0, increasedPrice * _1e6)
        const aliceLongOrderForLiquidation = await placeOrder(alice, charlieOrderSize, reducedPrice, 605)
        aliceOrderSize += charlieOrderSize
        aliceOpenNotional = aliceOpenNotional + Math.abs(charlieOrderSize * reducedPrice * _1e6)
        aliceLiquidationThreshold = getLiquidationThreshold(aliceOrderSize)
        aliceTradeFee = aliceTradeFee + getTradeFee(charlieOrderSize * reducedPrice)

        charlieMargin = (_1e6 * 100
            - charlieTradeFee // tradeFee for initial order
            - (reducedPrice * charlieOrderSize * 0.05 * _1e6) // 5% liquidation penalty
            - ((charliePrice - reducedPrice) *  charlieOrderSize * _1e6)) // negative pnl for liquidated position


        const expectedState = {
            "order_map": {},
            "trader_map": {
                [bobAddress]: {
                    "positions": {
                        "0": {
                            "open_notional": bobOpenNotional,
                            "size": _1e18 * bobOrderSize,
                            "unrealised_funding": 0,
                            "last_premium_fraction": 0,
                            "liquidation_threshold": bobLiquidationThreshold
                        }
                    },
                    "margin": {
                        "reserved": 0,
                        "deposited": {
                            "0": bobMargin - bobTradeFee
                        }
                    }
                },
                [aliceAddress]: {
                    "positions": {
                        "0": {
                            "open_notional": aliceOpenNotional,
                            "size": _1e18 * aliceOrderSize,
                            "unrealised_funding": 0,
                            "last_premium_fraction": 0,
                            "liquidation_threshold": aliceLiquidationThreshold
                        }
                    },
                    "margin": {
                        "reserved": 0,
                        "deposited": {
                            "0": aliceMargin - aliceTradeFee
                        }
                    }
                },
                [charlieAddress]: {
                    "positions": {
                        "0": {
                            "open_notional": 0,
                            "size": 0,
                            "unrealised_funding": 0,
                            "last_premium_fraction": 0,
                            "liquidation_threshold": 0
                        }
                    },
                    "margin": {
                        "reserved": 0,
                        "deposited": {
                            "0": charlieMargin
                        }
                    }
                }
            },
            "next_funding_time": await getNextFundingTime(),
            "last_price": {
                "0": _1e6 * reducedPrice
            }
        }

        await sleep(5)
        const evmState = await getEVMState()
        // console.log(JSON.stringify(evmState, null, 2))
        // console.log(JSON.stringify(expectedState, null, 2))
        expect(evmState).to.deep.contain(expectedState)
    });
});

async function placeOrder(trader, size, price, salt, reduceOnly=false) {
    const order = {
        ammIndex: ZERO,
        trader: trader.address,
        baseAssetQuantity: ethers.utils.parseEther(size.toString()),
        price: ethers.utils.parseUnits(price.toString(), 6),
        salt: BigNumber.from(salt),
        reduceOnly: reduceOnly,
    }
    const signature = await trader._signTypedData(domain, orderType, order)
    const hash = await orderBook.connect(trader).getOrderHash(order)

    const tx = await orderBook.connect(trader).placeOrder(order, signature)
    const txReceipt = await tx.wait()
    return { tx, txReceipt, hash, order, signature: signature.slice(2) }
}

function addMargin(trader, amount) {
    const hgtAmount = _1e12.mul(amount)
    return marginAccountHelper.connect(trader).addVUSDMarginWithReserve(amount, { value: hgtAmount })
}

async function getNextFundingTime() {
    const fundingEvents = await clearingHouse.queryFilter('FundingRateUpdated')
    const latestFundingEvent = fundingEvents.pop()
    return latestFundingEvent.args.nextFundingTime.toNumber()
}

function getLiquidationThreshold(size) {
    const absSize = Math.abs(size)
    let liquidationThreshold = Math.max(absSize / 4, 0.01)
    return size >= 0 ? _1e18 * liquidationThreshold : _1e18 * -liquidationThreshold;
}

function getReservedMargin(notional) {
    const leveraged = Math.abs(notional / maxLeverage)
    let tradeFee = leveraged * tradeFeeRatio
    let reserved = leveraged + tradeFee
    return _1e6 * reserved
}

function getTradeFee(notional) {
    const leveraged = Math.abs(notional / maxLeverage)
    let tradeFee = leveraged * tradeFeeRatio
    return _1e6 * tradeFee
}

async function setOraclePrice(market, price) {
    const ammAddress = await clearingHouse.amms(market)
    const amm = new ethers.Contract(ammAddress, require('./abi/AMM.json'), provider);
    const underlying = await amm.underlyingAsset()
    const oracleAddress = await marginAccount.oracle()
    const oracle = new ethers.Contract(oracleAddress, require('./abi/Oracle.json'), provider);

    await oracle.connect(governance).setStablePrice(underlying, price)
}

async function getEVMState() {
    const response = await axios.post(url, {
        jsonrpc: '2.0',
        id: 1,
        method: 'orderbook_getDetailedOrderBookData',
        params: []
    }, {
        headers: {
            'Content-Type': 'application/json'
        }
    });

    return response.data.result
}

function sleep(s) {
    console.log(`Requested a sleep of ${s} seconds...`)
    return new Promise(resolve => setTimeout(resolve, s * 1000));
}
