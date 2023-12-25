// (c) 2020-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package state

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/trie/trienode"
	"github.com/ethereum/go-ethereum/rlp"
)

type op struct {
	receiver string
	method   string
	args     []interface{}
	expected []interface{}

	r *recording
}

func display(obj interface{}) interface{} {
	if obj == nil {
		return obj
	}
	if acc, ok := obj.(*types.StateAccount); ok {
		bytes, err := rlp.EncodeToBytes(acc)
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf("%x", bytes)
	}
	if bytes, ok := obj.([]byte); ok {
		return fmt.Sprintf("%x", bytes)
	}
	if ns, ok := obj.(*trienode.NodeSet); ok {
		return fmt.Sprintf("ns <%s>", ns.Owner)
	}
	if ns, ok := obj.(*trienode.MergedNodeSet); ok {
		return fmt.Sprintf("mns <%d>", len(ns.Sets))
	}
	return obj
}

func (op *op) String() string {
	showArgs := make([]interface{}, len(op.args))
	for i, arg := range op.args {
		showArgs[i] = display(arg)
	}
	showExpected := make([]interface{}, len(op.expected))
	for i, arg := range op.expected {
		showExpected[i] = display(arg)
	}
	return fmt.Sprintf("%s.%s(%v) -> %v", op.receiver, op.method, showArgs, showExpected)
}

type knownType struct {
	rtype reflect.Type
	objs  []interface{}
}

type recording struct {
	knownTypes map[string]*knownType

	Ops []*op
}

func NewRecording() *recording {
	return &recording{
		knownTypes: make(map[string]*knownType),
	}
}

func (r *recording) Replay() error {
	for _, op := range r.Ops {
		bits := strings.Split(op.receiver, ":")
		id, err := strconv.Atoi(bits[1])
		if err != nil {
			panic(err)
		}
		receiver := r.knownTypes[bits[0]].objs[id]
		args := make([]reflect.Value, len(op.args))
		for i, arg := range op.args {
			args[i] = reflect.ValueOf(arg)
		}
		result := reflect.ValueOf(receiver).MethodByName(op.method).Call(args)
		for i, res := range result {
			fmt.Println(res, res.Type())
			// we should register this object
			if res.IsNil() && op.expected[i] == nil {
				fmt.Printf("return value %d is nil as expected\n", i)
				continue
			}
		}
	}
	return nil
}

func (r *recording) RegisterType(in interface{}) {
	i := reflect.TypeOf(in).Elem()
	r.knownTypes[i.String()] = &knownType{rtype: i}
}

func (r *recording) Register(obj interface{}) (string, bool) {
	for _, iface := range r.knownTypes {
		if reflect.TypeOf(obj).AssignableTo(iface.rtype) {
			for i, o := range iface.objs {
				if o == obj {
					return fmt.Sprintf("%s:%d", iface.rtype.String(), i), true
				}
			}
			iface.objs = append(iface.objs, obj)
			return fmt.Sprintf("%s:%d", iface.rtype.String(), len(iface.objs)-1), true
		}
	}
	return "", false
}

func (r *recording) RegisterAll(objs ...interface{}) []interface{} {
	out := make([]interface{}, 0, len(objs))
	for _, obj := range objs {
		if obj == nil {
			out = append(out, nil)
		} else if id, ok := r.Register(obj); ok {
			out = append(out, id)
		} else {
			out = append(out, obj)
		}
	}
	return out
}

func (r *recording) Record(receiver interface{}, method string, args ...interface{}) *op {
	obj, ok := r.Register(receiver)
	if !ok {
		panic(fmt.Errorf("receiver %v not registered", receiver))
	}
	op := &op{
		r:        r,
		receiver: obj,
		method:   method,
		args:     args,
	}
	if obj == "state.Trie:4" && method == "Hash" {
		fmt.Println("here")
	}
	r.Ops = append(r.Ops, op)
	return op
}

func (op *op) Returns(expected ...interface{}) {
	op.expected = op.r.RegisterAll(expected...)
}
