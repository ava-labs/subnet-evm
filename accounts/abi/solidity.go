package abi

import (
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/log"
)

type SolidityJSON struct {
	ContractName string `json:"contractName"`
	SourceName   string `json:"sourceName"`
	Abi          []Abi  `json:"abi"`
}

type Abi struct {
	Inputs          []Input `json:"inputs"`
	StateMutability string  `json:"stateMutability,omitempty"`
	Type            string  `json:"type"`
	Anonymous       bool    `json:"anonymous,omitempty"`
	Name            string  `json:"name,omitempty"`
	Outputs         []Input `json:"outputs,omitempty"`
}

type Input struct {
	Components   []Input `json:"components"`
	InternalType string  `json:"internalType"`
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	Indexed      bool    `json:"indexed"`
}

func getFunctionType(type_ string) FunctionType {
	typeMap := map[string]FunctionType{
		"function":    Function,
		"receive":     Receive,
		"fallback":    Fallback,
		"constructor": Constructor,
	}
	return typeMap[type_]
}

func FromSolidityJson(abiJsonInput string) (ABI, error) {
	solidityJson := SolidityJSON{}
	err := json.Unmarshal([]byte(abiJsonInput), &solidityJson)
	if err != nil {
		log.Error("Error in decoding ABI json", "err", err)
		return ABI{}, err
	}

	var constructor Method
	var receive Method
	var fallback Method

	methods := map[string]Method{}
	events := map[string]Event{}
	errors := map[string]Error{}

	for _, method := range solidityJson.Abi {
		inputs := []Argument{}
		for _, input := range method.Inputs {
			components := []ArgumentMarshaling{}
			if strings.HasPrefix(input.Type, "tuple") { // covers "tuple", "tuple[2]", "tuple[]"
				for _, component := range input.Components {
					components = append(components, ArgumentMarshaling{
						Name:         component.Name,
						Type:         component.Type,
						InternalType: component.InternalType,
					})
				}
			}

			type_, _ := NewType(input.Type, input.InternalType, components)
			inputs = append(inputs, Argument{
				Name:    input.Name,
				Type:    type_,
				Indexed: input.Indexed,
			})
		}

		if method.Type == "event" {
			abiEvent := NewEvent(method.Name, method.Name, method.Anonymous, inputs)
			events[method.Name] = abiEvent
			continue
		}

		if method.Type == "error" {
			abiError := NewError(method.Name, inputs)
			errors[method.Name] = abiError
			continue
		}

		outputs := []Argument{}
		for _, output := range method.Outputs {
			components := []ArgumentMarshaling{}
			if output.Type == "tuple" || output.Type == "tuple[2]" {
				for _, component := range output.Components {
					components = append(components, ArgumentMarshaling{
						Name:         component.Name,
						Type:         component.Type,
						InternalType: component.InternalType,
					})
				}
			}
			type_, _ := NewType(output.Type, output.InternalType, components)
			outputs = append(outputs, Argument{
				Name:    output.Name,
				Type:    type_,
				Indexed: output.Indexed,
			})
		}

		methodType := getFunctionType(method.Type)
		abiMethod := NewMethod(method.Name, method.Name, methodType, method.StateMutability, false, method.StateMutability == "payable", inputs, outputs)

		// don't include the method in the list if it's a constructor
		if methodType == Constructor {
			constructor = abiMethod
			continue
		}
		if methodType == Fallback {
			fallback = abiMethod
			continue
		}
		if methodType == Receive {
			receive = abiMethod
			continue
		}
		methods[method.Name] = abiMethod
	}

	Abi := ABI{
		Constructor: constructor,
		Methods:     methods,
		Events:      events,
		Errors:      errors,

		Receive:  receive,
		Fallback: fallback,
	}

	return Abi, nil
}
