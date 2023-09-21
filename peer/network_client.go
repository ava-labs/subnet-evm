// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package peer

import (
	"errors"

	"github.com/ava-labs/avalanchego/ids"

	"github.com/ava-labs/avalanchego/version"
)

var (
	_ NetworkClient = &networkClient{}

	ErrRequestFailed = errors.New("request failed")
)

// NetworkClient defines ability to send request / response through the Network
type NetworkClient interface {
	// SendAppRequestAny synchronously sends request to an arbitrary peer with a
	// node version greater than or equal to minVersion.
	// Returns response bytes, the ID of the chosen peer, and ErrRequestFailed if
	// the request should be retried.
	SendAppRequestAny(minVersion *version.Application, request []byte) ([]byte, ids.NodeID, error)

	// SendAppRequest synchronously sends request to the selected nodeID
	// Returns response bytes, and ErrRequestFailed if the request should be retried.
	SendAppRequest(nodeID ids.NodeID, request []byte) ([]byte, error)

	// SendCrossChainRequest sends a request to a specific blockchain running on this node.
	// Returns response bytes, and ErrRequestFailed if the request failed.
	SendCrossChainRequest(chainID ids.ID, request []byte) ([]byte, error)

	// Gossip sends given gossip message to peers
	Gossip(gossip []byte) error

	// TrackBandwidth should be called for each valid request with the bandwidth
	// (length of response divided by request time), and with 0 if the response is invalid.
	TrackBandwidth(nodeID ids.NodeID, bandwidth float64)
}

// networkClient implements NetworkClient interface
// provides ability to send request / responses through the Network and wait for a response
// so that the caller gets the result synchronously.
type networkClient struct {
	network Network
}

// NewNetworkClient returns Client for a given network
func NewNetworkClient(network Network) NetworkClient {
	return &networkClient{
		network: network,
	}
}

// SendAppRequestAny synchronously sends request to an arbitrary peer with a
// node version greater than or equal to minVersion.
// Returns response bytes, the ID of the chosen peer, and ErrRequestFailed if
// the request should be retried.
func (c *networkClient) SendAppRequestAny(minVersion *version.Application, request []byte) ([]byte, ids.NodeID, error) {
	waitingHandler := newWaitingResponseHandler()
	nodeID, err := c.network.SendAppRequestAny(minVersion, request, waitingHandler)
	if err != nil {
		return nil, nodeID, err
	}
	response := <-waitingHandler.responseChan
	if waitingHandler.failed {
		return nil, nodeID, ErrRequestFailed
	}
	return response, nodeID, nil
}

// SendAppRequest synchronously sends request to the specified nodeID
// Returns response bytes and ErrRequestFailed if the request should be retried.
func (c *networkClient) SendAppRequest(nodeID ids.NodeID, request []byte) ([]byte, error) {
	waitingHandler := newWaitingResponseHandler()
	if err := c.network.SendAppRequest(nodeID, request, waitingHandler); err != nil {
		return nil, err
	}
	response := <-waitingHandler.responseChan
	if waitingHandler.failed {
		return nil, ErrRequestFailed
	}
	return response, nil
}

// SendCrossChainRequest synchronously sends request to the specified chainID
// Returns response bytes and ErrRequestFailed if the request should be retried.
func (c *networkClient) SendCrossChainRequest(chainID ids.ID, request []byte) ([]byte, error) {
	waitingHandler := newWaitingResponseHandler()
	if err := c.network.SendCrossChainRequest(chainID, request, waitingHandler); err != nil {
		return nil, err
	}
	response := <-waitingHandler.responseChan
	if waitingHandler.failed {
		return nil, ErrRequestFailed
	}
	return response, nil
}

func (c *networkClient) Gossip(gossip []byte) error {
	return c.network.Gossip(gossip)
}

func (c *networkClient) TrackBandwidth(nodeID ids.NodeID, bandwidth float64) {
	c.network.TrackBandwidth(nodeID, bandwidth)
}
