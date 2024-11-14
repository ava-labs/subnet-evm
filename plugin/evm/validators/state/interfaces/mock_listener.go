// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ava-labs/subnet-evm/plugin/evm/validators/state/interfaces (interfaces: StateCallbackListener)
//
// Generated by this command:
//
//	mockgen -package=interfaces -destination=plugin/evm/validators/state/interfaces/mock_listener.go github.com/ava-labs/subnet-evm/plugin/evm/validators/state/interfaces StateCallbackListener
//

// Package interfaces is a generated GoMock package.
package interfaces

import (
	reflect "reflect"

	ids "github.com/ava-labs/avalanchego/ids"
	gomock "go.uber.org/mock/gomock"
)

// MockStateCallbackListener is a mock of StateCallbackListener interface.
type MockStateCallbackListener struct {
	ctrl     *gomock.Controller
	recorder *MockStateCallbackListenerMockRecorder
}

// MockStateCallbackListenerMockRecorder is the mock recorder for MockStateCallbackListener.
type MockStateCallbackListenerMockRecorder struct {
	mock *MockStateCallbackListener
}

// NewMockStateCallbackListener creates a new mock instance.
func NewMockStateCallbackListener(ctrl *gomock.Controller) *MockStateCallbackListener {
	mock := &MockStateCallbackListener{ctrl: ctrl}
	mock.recorder = &MockStateCallbackListenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStateCallbackListener) EXPECT() *MockStateCallbackListenerMockRecorder {
	return m.recorder
}

// OnValidatorAdded mocks base method.
func (m *MockStateCallbackListener) OnValidatorAdded(arg0 ids.ID, arg1 ids.NodeID, arg2 uint64, arg3 bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnValidatorAdded", arg0, arg1, arg2, arg3)
}

// OnValidatorAdded indicates an expected call of OnValidatorAdded.
func (mr *MockStateCallbackListenerMockRecorder) OnValidatorAdded(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnValidatorAdded", reflect.TypeOf((*MockStateCallbackListener)(nil).OnValidatorAdded), arg0, arg1, arg2, arg3)
}

// OnValidatorRemoved mocks base method.
func (m *MockStateCallbackListener) OnValidatorRemoved(arg0 ids.ID, arg1 ids.NodeID) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnValidatorRemoved", arg0, arg1)
}

// OnValidatorRemoved indicates an expected call of OnValidatorRemoved.
func (mr *MockStateCallbackListenerMockRecorder) OnValidatorRemoved(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnValidatorRemoved", reflect.TypeOf((*MockStateCallbackListener)(nil).OnValidatorRemoved), arg0, arg1)
}

// OnValidatorStatusUpdated mocks base method.
func (m *MockStateCallbackListener) OnValidatorStatusUpdated(arg0 ids.ID, arg1 ids.NodeID, arg2 bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnValidatorStatusUpdated", arg0, arg1, arg2)
}

// OnValidatorStatusUpdated indicates an expected call of OnValidatorStatusUpdated.
func (mr *MockStateCallbackListenerMockRecorder) OnValidatorStatusUpdated(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnValidatorStatusUpdated", reflect.TypeOf((*MockStateCallbackListener)(nil).OnValidatorStatusUpdated), arg0, arg1, arg2)
}
