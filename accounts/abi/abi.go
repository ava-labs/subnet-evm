// (c) 2019-2020, Ava Labs, Inc.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package abi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// The ABI holds information about a contract's context and available
// invokable methods. It will allow you to type check function calls and
// packs data accordingly.
type ABI struct {
	Constructor Method
	Methods     map[string]Method
	Events      map[string]Event
	Errors      map[string]Error

	// Additional "special" functions introduced in solidity v0.6.0.
	// It's separated from the original default fallback. Each contract
	// can only define one fallback and receive function.
	Fallback Method // Note it's also used to represent legacy fallback before v0.6.0
	Receive  Method
}

// JSON returns a parsed ABI interface and error if it failed.
func JSON(reader io.Reader) (ABI, error) {
	dec := json.NewDecoder(reader)

	var abi ABI
	if err := dec.Decode(&abi); err != nil {
		return ABI{}, err
	}
	return abi, nil
}

// Pack the given method name to conform the ABI. Method call's data
// will consist of method_id, args0, arg1, ... argN. Method id consists
// of 4 bytes and arguments are all 32 bytes.
// Method ids are created from the first 4 bytes of the hash of the
// methods string signature. (signature = baz(uint32,string32))
func (abi ABI) Pack(name string, args ...interface{}) ([]byte, error) {
	// Fetch the ABI of the requested method
	if name == "" {
		// constructor
		arguments, err := abi.Constructor.Inputs.Pack(args...)
		if err != nil {
			return nil, err
		}
		return arguments, nil
	}
	method, exist := abi.Methods[name]
	if !exist {
		return nil, fmt.Errorf("method '%s' not found", name)
	}
	arguments, err := method.Inputs.Pack(args...)
	if err != nil {
		return nil, err
	}
	// Pack up the method ID too if not a constructor and return
	return append(method.ID, arguments...), nil
}

// PackEvent packs the given event name and arguments to conform the ABI.
// Returns the topics for the event including the event signature (if non-anonymous event) and
// hashes derived from indexed arguments and the packed data of non-indexed args according to
// the event ABI specification.
// The order of arguments must match the order of the event definition.
// https://docs.soliditylang.org/en/v0.8.17/abi-spec.html#indexed-event-encoding.
// Note: PackEvent does not support array (fixed or dynamic-size) or struct types.
func (abi ABI) PackEvent(name string, args ...interface{}) ([]common.Hash, []byte, error) {
	event, exist := abi.Events[name]
	if !exist {
		return nil, nil, fmt.Errorf("event '%s' not found", name)
	}
	if len(args) != len(event.Inputs) {
		return nil, nil, fmt.Errorf("event '%s' unexpected number of inputs %d", name, len(args))
	}

	var (
		nonIndexedInputs = make([]interface{}, 0)
		indexedInputs    = make([]interface{}, 0)
		nonIndexedArgs   Arguments
		indexedArgs      Arguments
	)

	for i, arg := range event.Inputs {
		if arg.Indexed {
			indexedArgs = append(indexedArgs, arg)
			indexedInputs = append(indexedInputs, args[i])
		} else {
			nonIndexedArgs = append(nonIndexedArgs, arg)
			nonIndexedInputs = append(nonIndexedInputs, args[i])
		}
	}

	packedArguments, err := nonIndexedArgs.Pack(nonIndexedInputs...)
	if err != nil {
		return nil, nil, err
	}
	topics := make([]common.Hash, 0, len(indexedArgs)+1)
	if !event.Anonymous {
		topics = append(topics, event.ID)
	}
	indexedTopics, err := PackTopics(indexedInputs)
	if err != nil {
		return nil, nil, err
	}

	return append(topics, indexedTopics...), packedArguments, nil
}

// PackOutput packs the given [args] as the output of given method [name] to conform the ABI.
// This does not include method ID.
func (abi ABI) PackOutput(name string, args ...interface{}) ([]byte, error) {
	// Fetch the ABI of the requested method
	method, exist := abi.Methods[name]
	if !exist {
		return nil, fmt.Errorf("method '%s' not found", name)
	}
	arguments, err := method.Outputs.Pack(args...)
	if err != nil {
		return nil, err
	}
	return arguments, nil
}

// getInputs gets input arguments of the given [name] method.
// useStrictMode indicates whether to check the input data length strictly.
func (abi ABI) getInputs(name string, data []byte, useStrictMode bool) (Arguments, error) {
	// since there can't be naming collisions with contracts and events,
	// we need to decide whether we're calling a method or an event
	var args Arguments
	if method, ok := abi.Methods[name]; ok {
		if useStrictMode && len(data)%32 != 0 {
			return nil, fmt.Errorf("abi: improperly formatted input: %s - Bytes: [%+v]", string(data), data)
		}
		args = method.Inputs
	}
	if event, ok := abi.Events[name]; ok {
		args = event.Inputs
	}
	if args == nil {
		return nil, fmt.Errorf("abi: could not locate named method or event: %s", name)
	}
	return args, nil
}

// getArguments gets output arguments of the given [name] method.
func (abi ABI) getArguments(name string, data []byte) (Arguments, error) {
	// since there can't be naming collisions with contracts and events,
	// we need to decide whether we're calling a method or an event
	var args Arguments
	if method, ok := abi.Methods[name]; ok {
		if len(data)%32 != 0 {
			return nil, fmt.Errorf("abi: improperly formatted output: %q - Bytes: %+v", data, data)
		}
		args = method.Outputs
	}
	if event, ok := abi.Events[name]; ok {
		args = event.Inputs
	}
	if args == nil {
		return nil, fmt.Errorf("abi: could not locate named method or event: %s", name)
	}
	return args, nil
}

// UnpackInput unpacks the input according to the ABI specification.
// useStrictMode indicates whether to check the input data length strictly.
// By default it was set to true. In order to support the general EVM tool compatibility this
// should be set to false. This transition (true -> false) should be done with a network upgrade.
func (abi ABI) UnpackInput(name string, data []byte, useStrictMode bool) ([]interface{}, error) {
	args, err := abi.getInputs(name, data, useStrictMode)
	if err != nil {
		return nil, err
	}
	return args.Unpack(data)
}

// Unpack unpacks the output according to the abi specification.
func (abi ABI) Unpack(name string, data []byte) ([]interface{}, error) {
	args, err := abi.getArguments(name, data)
	if err != nil {
		return nil, err
	}
	return args.Unpack(data)
}

// UnpackInputIntoInterface unpacks the input in v according to the ABI specification.
// It performs an additional copy. Please only use, if you want to unpack into a
// structure that does not strictly conform to the ABI structure (e.g. has additional arguments)
// useStrictMode indicates whether to check the input data length strictly.
// By default it was set to true. In order to support the general EVM tool compatibility this
// should be set to false. This transition (true -> false) should be done with a network upgrade.
func (abi ABI) UnpackInputIntoInterface(v interface{}, name string, data []byte, useStrictMode bool) error {
	args, err := abi.getInputs(name, data, useStrictMode)
	if err != nil {
		return err
	}
	unpacked, err := args.Unpack(data)
	if err != nil {
		return err
	}
	return args.Copy(v, unpacked)
}

// UnpackIntoInterface unpacks the output in v according to the abi specification.
// It performs an additional copy. Please only use, if you want to unpack into a
// structure that does not strictly conform to the abi structure (e.g. has additional arguments)
func (abi ABI) UnpackIntoInterface(v interface{}, name string, data []byte) error {
	args, err := abi.getArguments(name, data)
	if err != nil {
		return err
	}
	unpacked, err := args.Unpack(data)
	if err != nil {
		return err
	}
	return args.Copy(v, unpacked)
}

// UnpackIntoMap unpacks a log into the provided map[string]interface{}.
func (abi ABI) UnpackIntoMap(v map[string]interface{}, name string, data []byte) (err error) {
	args, err := abi.getArguments(name, data)
	if err != nil {
		return err
	}
	return args.UnpackIntoMap(v, data)
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (abi *ABI) UnmarshalJSON(data []byte) error {
	var fields []struct {
		Type    string
		Name    string
		Inputs  []Argument
		Outputs []Argument

		// Status indicator which can be: "pure", "view",
		// "nonpayable" or "payable".
		StateMutability string

		// Deprecated Status indicators, but removed in v0.6.0.
		Constant bool // True if function is either pure or view
		Payable  bool // True if function is payable

		// Event relevant indicator represents the event is
		// declared as anonymous.
		Anonymous bool
	}
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	abi.Methods = make(map[string]Method)
	abi.Events = make(map[string]Event)
	abi.Errors = make(map[string]Error)
	for _, field := range fields {
		switch field.Type {
		case "constructor":
			abi.Constructor = NewMethod("", "", Constructor, field.StateMutability, field.Constant, field.Payable, field.Inputs, nil)
		case "function":
			name := ResolveNameConflict(field.Name, func(s string) bool { _, ok := abi.Methods[s]; return ok })
			abi.Methods[name] = NewMethod(name, field.Name, Function, field.StateMutability, field.Constant, field.Payable, field.Inputs, field.Outputs)
		case "fallback":
			// New introduced function type in v0.6.0, check more detail
			// here https://solidity.readthedocs.io/en/v0.6.0/contracts.html#fallback-function
			if abi.HasFallback() {
				return errors.New("only single fallback is allowed")
			}
			abi.Fallback = NewMethod("", "", Fallback, field.StateMutability, field.Constant, field.Payable, nil, nil)
		case "receive":
			// New introduced function type in v0.6.0, check more detail
			// here https://solidity.readthedocs.io/en/v0.6.0/contracts.html#fallback-function
			if abi.HasReceive() {
				return errors.New("only single receive is allowed")
			}
			if field.StateMutability != "payable" {
				return errors.New("the statemutability of receive can only be payable")
			}
			abi.Receive = NewMethod("", "", Receive, field.StateMutability, field.Constant, field.Payable, nil, nil)
		case "event":
			name := ResolveNameConflict(field.Name, func(s string) bool { _, ok := abi.Events[s]; return ok })
			abi.Events[name] = NewEvent(name, field.Name, field.Anonymous, field.Inputs)
		case "error":
			// Errors cannot be overloaded or overridden but are inherited,
			// no need to resolve the name conflict here.
			abi.Errors[field.Name] = NewError(field.Name, field.Inputs)
		default:
			return fmt.Errorf("abi: could not recognize type %v of field %v", field.Type, field.Name)
		}
	}
	return nil
}

// MethodById looks up a method by the 4-byte id,
// returns nil if none found.
func (abi *ABI) MethodById(sigdata []byte) (*Method, error) {
	if len(sigdata) < 4 {
		return nil, fmt.Errorf("data too short (%d bytes) for abi method lookup", len(sigdata))
	}
	for _, method := range abi.Methods {
		if bytes.Equal(method.ID, sigdata[:4]) {
			return &method, nil
		}
	}
	return nil, fmt.Errorf("no method with id: %#x", sigdata[:4])
}

// EventByID looks an event up by its topic hash in the
// ABI and returns nil if none found.
func (abi *ABI) EventByID(topic common.Hash) (*Event, error) {
	for _, event := range abi.Events {
		if bytes.Equal(event.ID.Bytes(), topic.Bytes()) {
			return &event, nil
		}
	}
	return nil, fmt.Errorf("no event with id: %#x", topic.Hex())
}

// ErrorByID looks up an error by the 4-byte id,
// returns nil if none found.
func (abi *ABI) ErrorByID(sigdata [4]byte) (*Error, error) {
	for _, errABI := range abi.Errors {
		if bytes.Equal(errABI.ID[:4], sigdata[:]) {
			return &errABI, nil
		}
	}
	return nil, fmt.Errorf("no error with id: %#x", sigdata[:])
}

// HasFallback returns an indicator whether a fallback function is included.
func (abi *ABI) HasFallback() bool {
	return abi.Fallback.Type == Fallback
}

// HasReceive returns an indicator whether a receive function is included.
func (abi *ABI) HasReceive() bool {
	return abi.Receive.Type == Receive
}

// revertSelector is a special function selector for revert reason unpacking.
var revertSelector = crypto.Keccak256([]byte("Error(string)"))[:4]

// UnpackRevert resolves the abi-encoded revert reason. According to the solidity
// spec https://solidity.readthedocs.io/en/latest/control-structures.html#revert,
// the provided revert reason is abi-encoded as if it were a call to a function
// `Error(string)`. So it's a special tool for it.
func UnpackRevert(data []byte) (string, error) {
	if len(data) < 4 {
		return "", errors.New("invalid data for unpacking")
	}
	if !bytes.Equal(data[:4], revertSelector) {
		return "", errors.New("invalid data for unpacking")
	}
	typ, err := NewType("string", "", nil)
	if err != nil {
		return "", err
	}
	unpacked, err := (Arguments{{Type: typ}}).Unpack(data[4:])
	if err != nil {
		return "", err
	}
	return unpacked[0].(string), nil
}
