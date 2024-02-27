package orderbook

const (
	HandleChainAcceptedEventPanicMessage  = "panic while processing chainAcceptedEvent"
	HandleChainAcceptedLogsPanicMessage   = "panic while processing chainAcceptedLogs"
	HandleHubbleFeedLogsPanicMessage      = "panic while processing hubbleFeedLogs"
	RunMatchingPipelinePanicMessage       = "panic while running matching pipeline"
	RunSanitaryPipelinePanicMessage       = "panic while running sanitary pipeline"
	MakerBookFileWriteChannelPanicMessage = "panic while sending to makerbook file write channel"
)
