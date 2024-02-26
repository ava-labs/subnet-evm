package orderbook

import (
	"github.com/ava-labs/subnet-evm/metrics"
)

var (
	transactionsPerBlockHistogram = metrics.NewRegisteredHistogram("transactions/rate", nil, metrics.ResettingSample(metrics.NewExpDecaySample(1028, 0.015)))

	gasUsedPerBlockHistogram      = metrics.NewRegisteredHistogram("gas_used_per_block", nil, metrics.ResettingSample(metrics.NewExpDecaySample(1028, 0.015)))
	blockGasCostPerBlockHistogram = metrics.NewRegisteredHistogram("block_gas_cost", nil, metrics.ResettingSample(metrics.NewExpDecaySample(1028, 0.015)))

	ordersPlacedPerBlock    = metrics.NewRegisteredHistogram("orders_placed_per_block", nil, metrics.ResettingSample(metrics.NewExpDecaySample(1028, 0.015)))
	ordersCancelledPerBlock = metrics.NewRegisteredHistogram("orders_cancelled_per_block", nil, metrics.ResettingSample(metrics.NewExpDecaySample(1028, 0.015)))

	// only valid for OrderBook transactions send by this validator
	orderBookTransactionsSuccessTotalCounter = metrics.NewRegisteredCounter("orderbooktxs/total/success", nil)
	orderBookTransactionsFailureTotalCounter = metrics.NewRegisteredCounter("orderbooktxs/total/failure", nil)

	// panics are recovered but monitored
	AllPanicsCounter                         = metrics.NewRegisteredCounter("all_panics", nil)
	RunMatchingPipelinePanicsCounter         = metrics.NewRegisteredCounter("matching_pipeline_panics", nil)
	RunSanitaryPipelinePanicsCounter         = metrics.NewRegisteredCounter("sanitary_pipeline_panics", nil)
	HandleHubbleFeedLogsPanicsCounter        = metrics.NewRegisteredCounter("handle_hubble_feed_logs_panics", nil)
	HandleChainAcceptedLogsPanicsCounter     = metrics.NewRegisteredCounter("handle_chain_accepted_logs_panics", nil)
	HandleChainAcceptedEventPanicsCounter    = metrics.NewRegisteredCounter("handle_chain_accepted_event_panics", nil)
	HandleMatchingPipelineTimerPanicsCounter = metrics.NewRegisteredCounter("handle_matching_pipeline_timer_panics", nil)
	RPCPanicsCounter                         = metrics.NewRegisteredCounter("rpc_panic", nil)
	AwaitSignedOrdersGossipPanicsCounter     = metrics.NewRegisteredCounter("await_signed_orders_gossip_panics", nil)

	BuildBlockFailedWithLowBlockGasCounter = metrics.NewRegisteredCounter("build_block_failed_low_block_gas", nil)

	// lag between head and accepted block
	headBlockLagHistogram = metrics.NewRegisteredHistogram("head_block_lag", nil, metrics.ResettingSample(metrics.NewExpDecaySample(1028, 0.015)))

	// order id not found while deleting
	deleteOrderIdNotFoundCounter = metrics.NewRegisteredCounter("delete_order_id_not_found", nil)

	// unquenched liquidations
	unquenchedLiquidationsCounter = metrics.NewRegisteredCounter("unquenched_liquidations", nil)
	placeSignedOrderCounter       = metrics.NewRegisteredCounter("place_signed_order", nil)

	// makerbook write failures
	makerBookWriteFailuresCounter = metrics.NewRegisteredCounter("makerbook_write_failures", nil)

	// snapshot write failures
	SnapshotWriteFailuresCounter = metrics.NewRegisteredCounter("snapshot_write_failures", nil)
)
