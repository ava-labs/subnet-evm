
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
let alicePartialMatchedLongOrder, bobHighPriceShortOrder

const ZERO = BigNumber.from(0)
const _1e6 = BigNumber.from(10).pow(6)
const _1e8 = BigNumber.from(10).pow(8)
const _1e12 = BigNumber.from(10).pow(12)
const _1e18 = ethers.constants.WeiPerEther

const homedir = require('os').homedir()
let conf = require(`${homedir}/.hubblenet.json`)
const url = `http://127.0.0.1:9650/ext/bc/${conf.chain_id}/rpc`

describe('Submit transaction and compare with EVM state', function () {
    before('', async function () {
        provider = new ethers.providers.JsonRpcProvider(url);

        // Set up signer
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
            ]
        }

    })

    it('Add margin', async function () {
        tx = await addMargin(alice, _1e6.mul(40))
        await tx.wait();

        tx = await addMargin(bob, _1e6.mul(40))
        await tx.wait();

        const expectedState = {
            "order_map": {},
            "trader_map": {
                [bobAddress]: {
                    "positions": {},
                    "margins": {
                        "0": 40000000
                    }
                },
                [aliceAddress]: {
                    "positions": {},
                    "margins": {
                        "0": 40000000
                    }
                }
            },
            "next_funding_time": await getNextFundingTime(),
            "last_price": {
                "0": 0
            }
        }
        const evmState = await getEVMState()
        expect(evmState).to.deep.contain(expectedState)
    });

    it('Remove margin', async function () {
        const tx = await marginAccount.connect(alice).removeMargin(0, _1e6.mul(10))
        await tx.wait();

        const expectedState = {
            "order_map": {},
            "trader_map": {
                [bobAddress]: {
                    "positions": {},
                    "margins": {
                        "0": 40000000
                    }
                },
                [aliceAddress]: {
                    "positions": {},
                    "margins": {
                        "0": 30000000
                    }
                }
            },
            "next_funding_time": await getNextFundingTime(),
            "last_price": {
                "0": 0
            }
        }
        const evmState = await getEVMState()
        expect(evmState).to.deep.contain(expectedState)
    });

    it('Place order', async function () {
        const { hash, order, signature, tx, txReceipt } = await placeOrder(alice, 5, 10, 101)

        const expectedState = {
            "order_map": {
                [hash]: {
                    "market": 0,
                    "position_type": "long",
                    "user_address": aliceAddress,
                    "base_asset_quantity": 5000000000000000000,
                    "filled_base_asset_quantity": 0,
                    "salt": 101,
                    "price": 10000000,
                    "lifecycle_list": [
                        {
                            "BlockNumber": txReceipt.blockNumber,
                            "Status": 0
                        }
                    ],
                    "signature": signature,
                    "block_number": txReceipt.blockNumber
                }
            },
            "trader_map": {
                [bobAddress]: {
                    "positions": {},
                    "margins": {
                        "0": 40000000
                    }
                },
                [aliceAddress]: {
                    "positions": {},
                    "margins": {
                        "0": 30000000
                    }
                }
            },
            "next_funding_time": await getNextFundingTime(),
            "last_price": {
                "0": 0
            }
        }

        const evmState = await getEVMState()
        expect(evmState).to.deep.contain(expectedState)
    });

    it('Match order', async function () {
        await placeOrder(bob, -5, 10, 201)

        const expectedState = {
            "order_map": {},
            "trader_map": {
                [bobAddress]: {
                    "positions": {
                        "0": {
                            "open_notional": 50000000,
                            "size": -5000000000000000000,
                            "unrealised_funding": null,
                            "last_premium_fraction": 0,
                            "liquidation_threshold": -5000000000000000000
                        }
                    },
                    "margins": {
                        "0": 39975000
                    }
                },
                [aliceAddress]: {
                    "positions": {
                        "0": {
                            "open_notional": 50000000,
                            "size": 5000000000000000000,
                            "unrealised_funding": null,
                            "last_premium_fraction": 0,
                            "liquidation_threshold": 5000000000000000000
                        }
                    },
                    "margins": {
                        "0": 29975000
                    }
                }
            },
            "next_funding_time": await getNextFundingTime(),
            "last_price": {
                "0": 10000000
            }
        }
        const evmState = await getEVMState()
        expect(evmState).to.deep.contain(expectedState)
    });

    it('Partially match order', async function () {
        alicePartialMatchedLongOrder = await placeOrder(alice, 5, 10, 301)
        const { hash: bobShortOrderHash } = await placeOrder(bob, -2, 10, 302)

        const expectedState = {
            "order_map": {
                [alicePartialMatchedLongOrder.hash]: {
                    "market": 0,
                    "position_type": "long",
                    "user_address": aliceAddress,
                    "base_asset_quantity": 5000000000000000000,
                    "filled_base_asset_quantity": 2000000000000000000,
                    "salt": 301,
                    "price": 10000000,
                    "lifecycle_list": [
                        {
                            "BlockNumber": alicePartialMatchedLongOrder.txReceipt.blockNumber,
                            "Status": 0
                        }
                    ],
                    "signature": alicePartialMatchedLongOrder.signature,
                    "block_number": alicePartialMatchedLongOrder.txReceipt.blockNumber
                }
            },
            "trader_map": {
                [bobAddress]: {
                    "positions": {
                        "0": {
                            "open_notional": 70000000,
                            "size": -7000000000000000000,
                            "unrealised_funding": null,
                            "last_premium_fraction": 0,
                            "liquidation_threshold": -5000000000000000000
                        }
                    },
                    "margins": {
                        "0": 39965000
                    }
                },
                [aliceAddress]: {
                    "positions": {
                        "0": {
                            "open_notional": 70000000,
                            "size": 7000000000000000000,
                            "unrealised_funding": null,
                            "last_premium_fraction": 0,
                            "liquidation_threshold": 5000000000000000000
                        }
                    },
                    "margins": {
                        "0": 29965000
                    }
                }
            },
            "next_funding_time": await getNextFundingTime(),
            "last_price": {
                "0": 10000000
            }
        }
        const evmState = await getEVMState()
        expect(evmState).to.deep.contain(expectedState)
    });

    it('Order cancel', async function () {
        const { hash, order } = await placeOrder(alice, 2, 14, 401)

        tx = await orderBook.connect(alice).cancelOrder(order)
        await tx.wait()

        const expectedState = {
            "order_map": {
                [alicePartialMatchedLongOrder.hash]: {
                    "market": 0,
                    "position_type": "long",
                    "user_address": aliceAddress,
                    "base_asset_quantity": 5000000000000000000,
                    "filled_base_asset_quantity": 2000000000000000000,
                    "salt": 301,
                    "price": 10000000,
                    "lifecycle_list": [
                        {
                            "BlockNumber": alicePartialMatchedLongOrder.txReceipt.blockNumber,
                            "Status": 0
                        }
                    ],
                    "signature": alicePartialMatchedLongOrder.signature,
                    "block_number": alicePartialMatchedLongOrder.txReceipt.blockNumber
                }
            },
            "trader_map": {
                [bobAddress]: {
                    "positions": {
                        "0": {
                            "open_notional": 70000000,
                            "size": -7000000000000000000,
                            "unrealised_funding": null,
                            "last_premium_fraction": 0,
                            "liquidation_threshold": -5000000000000000000
                        }
                    },
                    "margins": {
                        "0": 39965000
                    }
                },
                [aliceAddress]: {
                    "positions": {
                        "0": {
                            "open_notional": 70000000,
                            "size": 7000000000000000000,
                            "unrealised_funding": null,
                            "last_premium_fraction": 0,
                            "liquidation_threshold": 5000000000000000000
                        }
                    },
                    "margins": {
                        "0": 29965000
                    }
                }
            },
            "next_funding_time": await getNextFundingTime(),
            "last_price": {
                "0": 10000000
            }
        }

        const evmState = await getEVMState()
        expect(evmState).to.deep.contain(expectedState)
    });

    it('Order match error', async function () {
        const { hash: charlieHash } = await placeOrder(charlie, 50, 12, 501)
        bobHighPriceShortOrder = await placeOrder(bob, -10, 12, 502)

        const expectedState = {
            "order_map": {
                [alicePartialMatchedLongOrder.hash]: {
                    "market": 0,
                    "position_type": "long",
                    "user_address": aliceAddress,
                    "base_asset_quantity": 5000000000000000000,
                    "filled_base_asset_quantity": 2000000000000000000,
                    "salt": 301,
                    "price": 10000000,
                    "lifecycle_list": [
                        {
                            "BlockNumber": alicePartialMatchedLongOrder.txReceipt.blockNumber,
                            "Status": 0
                        }
                    ],
                    "signature": alicePartialMatchedLongOrder.signature,
                    "block_number": alicePartialMatchedLongOrder.txReceipt.blockNumber
                },
                [bobHighPriceShortOrder.hash]: {
                    "market": 0,
                    "position_type": "short",
                    "user_address": bobAddress,
                    "base_asset_quantity": -10000000000000000000,
                    "filled_base_asset_quantity": 0,
                    "salt": 502,
                    "price": 12000000,
                    "lifecycle_list": [
                        {
                            "BlockNumber": bobHighPriceShortOrder.txReceipt.blockNumber,
                            "Status": 0
                        }
                    ],
                    "signature": bobHighPriceShortOrder.signature,
                    "block_number": bobHighPriceShortOrder.txReceipt.blockNumber
                }
            },
            "trader_map": {
                [bobAddress]: {
                    "positions": {
                        "0": {
                            "open_notional": 70000000,
                            "size": -7000000000000000000,
                            "unrealised_funding": null,
                            "last_premium_fraction": 0,
                            "liquidation_threshold": -5000000000000000000
                        }
                    },
                    "margins": {
                        "0": 39965000
                    }
                },
                [aliceAddress]: {
                    "positions": {
                        "0": {
                            "open_notional": 70000000,
                            "size": 7000000000000000000,
                            "unrealised_funding": null,
                            "last_premium_fraction": 0,
                            "liquidation_threshold": 5000000000000000000
                        }
                    },
                    "margins": {
                        "0": 29965000
                    }
                }
            },
            "next_funding_time": await getNextFundingTime(),
            "last_price": {
                "0": 10000000
            }
        }
        const evmState = await getEVMState()
        expect(evmState).to.deep.contain(expectedState)
    });

    it('Liquidate trader', async function () {
        await addMargin(charlie, _1e6.mul(100))
        await addMargin(alice, _1e6.mul(200))
        await addMargin(bob, _1e6.mul(200))

        await sleep(3)

        // large position by charlie
        const { hash: charlieHash } = await placeOrder(charlie, 49, 10, 601) // 46 + 3 is fulfilled
        const { hash: bobHash1 } = await placeOrder(bob, -49, 10, 602) // 46 + 3

        // reduce the price
        const { hash: aliceHash } = await placeOrder(alice, 10, 8, 603) // 7 matched; 3 used for liquidation
        const { hash: bobHash2 } = await placeOrder(bob, -10, 8, 604) // 3 + 7

        // long order so that liquidation can run
        const { hash } = await placeOrder(alice, 10, 8, 605) // 10 used for liquidation

        const expectedState = {
            "order_map": {
                [bobHighPriceShortOrder.hash]: {
                    "market": 0,
                    "position_type": "short",
                    "user_address": bobAddress,
                    "base_asset_quantity": -10000000000000000000,
                    "filled_base_asset_quantity": 0,
                    "salt": 502,
                    "price": 12000000,
                    "lifecycle_list": [
                        {
                            "BlockNumber": bobHighPriceShortOrder.txReceipt.blockNumber,
                            "Status": 0
                        }
                    ],
                    "signature": bobHighPriceShortOrder.signature,
                    "block_number": bobHighPriceShortOrder.txReceipt.blockNumber
                }
            },
            "trader_map": {
                [bobAddress]: {
                    "positions": {
                        "0": {
                            "open_notional": 646000000,
                            "size": -66000000000000000000,
                            "unrealised_funding": null,
                            "last_premium_fraction": 0,
                            "liquidation_threshold": -16500000000000000000
                        }
                    },
                    "margins": {
                        "0": 239677000
                    }
                },
                [charlieAddress]: {
                    "positions": {
                        "0": {
                            "open_notional": 360000000,
                            "size": 36000000000000000000,
                            "unrealised_funding": null,
                            "last_premium_fraction": 0,
                            "liquidation_threshold": 12250000000000000000
                        }
                    },
                    "margins": {
                        "0": 68555000
                    }
                },
                [aliceAddress]: {
                    "positions": {
                        "0": {
                            "open_notional": 260000000,
                            "size": 30000000000000000000,
                            "unrealised_funding": null,
                            "last_premium_fraction": 0,
                            "liquidation_threshold": 7500000000000000000
                        }
                    },
                    "margins": {
                        "0": 229870000
                    }
                }
            },
            "next_funding_time": await getNextFundingTime(),
            "last_price": {
                "0": 8000000
            }
        }

        const evmState = await getEVMState()
        expect(evmState).to.deep.contain(expectedState)
    });
});

async function placeOrder(trader, size, price, salt) {
    const order = {
        ammIndex: ZERO,
        trader: trader.address,
        baseAssetQuantity: ethers.utils.parseEther(size.toString()),
        price: ethers.utils.parseUnits(price.toString(), 6),
        salt: BigNumber.from(salt)
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
