// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package customethclient

import (
	"context"

	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/rpc"
)

// Client wraps the ethclient.Client interface to provide extra data types (in header, block body).
// If you want to use the standardized Ethereum RPC functionality without extra types, use [ethclient.Client] instead.
type Client struct {
	c      *rpc.Client
	client ethclient.Client
}

// New creates a client that uses the given RPC client.
func New(c *rpc.Client) *Client {
	return &Client{c: c, client: ethclient.NewClientWithHook(c, nil)}
}

// Dial connects a client to the given URL.
func Dial(rawurl string) (*Client, error) {
	return DialContext(context.Background(), rawurl)
}

// DialContext connects a client to the given URL with context.
func DialContext(ctx context.Context, rawurl string) (*Client, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return New(c), nil
}
