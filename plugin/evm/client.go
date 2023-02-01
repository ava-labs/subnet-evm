// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"fmt"

	"github.com/ava-labs/avalanchego/api"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/rpc"
	"github.com/ava-labs/subnet-evm/plugin/evm/message"
	"github.com/ethereum/go-ethereum/log"
)

// Interface compliance
var _ Client = (*client)(nil)

// Client interface for interacting with EVM [chain]
type Client interface {
	StartCPUProfiler(ctx context.Context) error
	StopCPUProfiler(ctx context.Context) error
	MemoryProfile(ctx context.Context) error
	LockProfile(ctx context.Context) error
	SetLogLevel(ctx context.Context, level log.Lvl) error
	GetVMConfig(ctx context.Context) (*Config, error)
	GetSignature(ctx context.Context, signatureRequest message.SignatureRequest) (*[bls.SignatureLen]byte, error)
}

// Client implementation for interacting with EVM [chain]
type client struct {
	requester rpc.EndpointRequester
}

// NewClient returns a Client for interacting with EVM [chain]
func NewClient(uri, chain, api string) Client {
	return &client{
		requester: rpc.NewEndpointRequester(fmt.Sprintf("%s/ext/bc/%s/%s", uri, chain, api)),
	}
}

// NewCChainClient returns a Client for interacting with the C Chain
func NewCChainClient(uri string) Client {
	// TODO: Update for Subnet-EVM compatibility
	return NewClient(uri, "C", "admin")
}

func (c *client) StartCPUProfiler(ctx context.Context) error {
	return c.requester.SendRequest(ctx, "admin.startCPUProfiler", struct{}{}, &api.EmptyReply{})
}

func (c *client) StopCPUProfiler(ctx context.Context) error {
	return c.requester.SendRequest(ctx, "admin.stopCPUProfiler", struct{}{}, &api.EmptyReply{})
}

func (c *client) MemoryProfile(ctx context.Context) error {
	return c.requester.SendRequest(ctx, "admin.memoryProfile", struct{}{}, &api.EmptyReply{})
}

func (c *client) LockProfile(ctx context.Context) error {
	return c.requester.SendRequest(ctx, "admin.lockProfile", struct{}{}, &api.EmptyReply{})
}

// SetLogLevel dynamically sets the log level for the C Chain
func (c *client) SetLogLevel(ctx context.Context, level log.Lvl) error {
	return c.requester.SendRequest(ctx, "admin.setLogLevel", &SetLogLevelArgs{
		Level: level.String(),
	}, &api.EmptyReply{})
}

// GetVMConfig returns the current config of the VM
func (c *client) GetVMConfig(ctx context.Context) (*Config, error) {
	res := &ConfigReply{}
	err := c.requester.SendRequest(ctx, "admin.getVMConfig", struct{}{}, res)
	return res.Config, err
}

func (c *client) GetSignature(ctx context.Context, signatureRequest message.SignatureRequest) (*[bls.SignatureLen]byte, error) {
	res := &message.SignatureResponse{}
	err := c.requester.SendRequest(ctx, "snowman.getSignature", &signatureRequest, res)
	return &res.Signature, err
}
