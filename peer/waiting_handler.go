// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package peer

import (
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/subnet-evm/plugin/evm/message"
)

var _ message.ResponseHandler = &waitingResponseHandler{}
var _ message.CrossChainResponseHandler = &waitingCrossChainResponseHandler{}

// waitingResponseHandler implements the ResponseHandler interface
// Internally used to wait for response after making a request synchronously
// responseChan may contain response bytes if the original request has not failed
// responseChan is closed in either fail or success scenario
type waitingResponseHandler struct {
	responseChan chan []byte // blocking channel with response bytes
	failed       bool        // whether the original request is failed
}

// OnResponse passes the response bytes to the responseChan and closes the channel
func (w *waitingResponseHandler) OnResponse(_ ids.NodeID, _ uint32, response []byte) error {
	w.responseChan <- response
	close(w.responseChan)
	return nil
}

// OnFailure sets the failed flag to true and closes the channel
func (w *waitingResponseHandler) OnFailure(ids.NodeID, uint32) error {
	w.failed = true
	close(w.responseChan)
	return nil
}

// newWaitingResponseHandler returns new instance of the waitingResponseHandler
func newWaitingResponseHandler() *waitingResponseHandler {
	return &waitingResponseHandler{responseChan: make(chan []byte)}
}

// waitingCrossChainResponseHandler implements the CrossChainResponseHandler interface.
// This is used internally to wait for a response after making a synchronous CrossChainRequest.
// responseChan may contain response bytes if the original request has not failed.
// responseChan is closed in either fail or success scenario.
type waitingCrossChainResponseHandler struct {
	responseChan chan []byte // blocking channel with response bytes
	failed       bool        // whether the original request is failed
}

// OnResponse passes the response bytes to the responseChan and closes the channel
func (w *waitingCrossChainResponseHandler) OnResponse(_ ids.ID, _ uint32, response []byte) error {
	w.responseChan <- response
	close(w.responseChan)
	return nil
}

// OnFailure sets the failed flag to true and closes the channel
func (w *waitingCrossChainResponseHandler) OnFailure(ids.ID, uint32) error {
	w.failed = true
	close(w.responseChan)
	return nil
}

// newWaitingCrossChainResponseHandler returns new instance of the waitingCrossChainResponseHandler
func newWaitingCrossChainResponseHandler() *waitingCrossChainResponseHandler {
	return &waitingCrossChainResponseHandler{responseChan: make(chan []byte)}
}
