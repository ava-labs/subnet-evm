// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Defines the stateless interface for unmarshalling an arbitrary config of a precompile
package execution

import (
	"fmt"
)

var registry = make(map[string]Execution)

func RegisterExecution(name string, execution Execution) error {
	_, exists := registry[name]
	if exists {
		return fmt.Errorf("cannot register duplicate precompile execution with the name: %s", name)
	}
	registry[name] = execution
	return nil
}

// TODO: add lookups as needed
