package limitorders

import (
	"github.com/ava-labs/avalanchego/snow/engine/common"
	"github.com/ethereum/go-ethereum/log"
)

var toEngine chan<- common.Message

func SetToEngine(toEngineChan chan<- common.Message) {
	toEngine = toEngineChan
}

func SendTxReadySignal() {
	select {
	case toEngine <- common.PendingTxs:
		log.Info("SendTxReadySignal - notified the consensus engine")
	default:
		log.Error("SendTxReadySignal - Failed to push PendingTxs notification to the consensus engine.")
	}
}
