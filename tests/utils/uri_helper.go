// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"fmt"
	"strings"
)

func ToWebsocketURI(uri string, blockchainID string) string {
	baseURI, _ := strings.CutPrefix(uri, "http://")
	return fmt.Sprintf("ws://%s/ext/bc/%s/ws", baseURI, blockchainID)
}

func ToRPCURI(uri string, blockchainID string) string {
	return fmt.Sprintf("%s/ext/bc/%s/rpc", uri, blockchainID)
}
