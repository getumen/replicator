// Code generated by MockGen. DO NOT EDIT.
// Source: transport.go

// Package mocktransport is a generated GoMock package.
package mocktransport

import (
	gomock "github.com/golang/mock/gomock"
	raft "github.com/hashicorp/raft"
	io "io"
	reflect "reflect"
)

// MockTransport is a mock of Transport interface
type MockTransport struct {
	ctrl     *gomock.Controller
	recorder *MockTransportMockRecorder
}

// MockTransportMockRecorder is the mock recorder for MockTransport
type MockTransportMockRecorder struct {
	mock *MockTransport
}

// NewMockTransport creates a new mock instance
func NewMockTransport(ctrl *gomock.Controller) *MockTransport {
	mock := &MockTransport{ctrl: ctrl}
	mock.recorder = &MockTransportMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTransport) EXPECT() *MockTransportMockRecorder {
	return m.recorder
}

// Consumer mocks base method
func (m *MockTransport) Consumer() <-chan raft.RPC {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Consumer")
	ret0, _ := ret[0].(<-chan raft.RPC)
	return ret0
}

// Consumer indicates an expected call of Consumer
func (mr *MockTransportMockRecorder) Consumer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Consumer", reflect.TypeOf((*MockTransport)(nil).Consumer))
}

// LocalAddr mocks base method
func (m *MockTransport) LocalAddr() raft.ServerAddress {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LocalAddr")
	ret0, _ := ret[0].(raft.ServerAddress)
	return ret0
}

// LocalAddr indicates an expected call of LocalAddr
func (mr *MockTransportMockRecorder) LocalAddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LocalAddr", reflect.TypeOf((*MockTransport)(nil).LocalAddr))
}

// AppendEntriesPipeline mocks base method
func (m *MockTransport) AppendEntriesPipeline(id raft.ServerID, target raft.ServerAddress) (raft.AppendPipeline, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppendEntriesPipeline", id, target)
	ret0, _ := ret[0].(raft.AppendPipeline)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AppendEntriesPipeline indicates an expected call of AppendEntriesPipeline
func (mr *MockTransportMockRecorder) AppendEntriesPipeline(id, target interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppendEntriesPipeline", reflect.TypeOf((*MockTransport)(nil).AppendEntriesPipeline), id, target)
}

// AppendEntries mocks base method
func (m *MockTransport) AppendEntries(id raft.ServerID, target raft.ServerAddress, args *raft.AppendEntriesRequest, resp *raft.AppendEntriesResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppendEntries", id, target, args, resp)
	ret0, _ := ret[0].(error)
	return ret0
}

// AppendEntries indicates an expected call of AppendEntries
func (mr *MockTransportMockRecorder) AppendEntries(id, target, args, resp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppendEntries", reflect.TypeOf((*MockTransport)(nil).AppendEntries), id, target, args, resp)
}

// RequestVote mocks base method
func (m *MockTransport) RequestVote(id raft.ServerID, target raft.ServerAddress, args *raft.RequestVoteRequest, resp *raft.RequestVoteResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RequestVote", id, target, args, resp)
	ret0, _ := ret[0].(error)
	return ret0
}

// RequestVote indicates an expected call of RequestVote
func (mr *MockTransportMockRecorder) RequestVote(id, target, args, resp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestVote", reflect.TypeOf((*MockTransport)(nil).RequestVote), id, target, args, resp)
}

// InstallSnapshot mocks base method
func (m *MockTransport) InstallSnapshot(id raft.ServerID, target raft.ServerAddress, args *raft.InstallSnapshotRequest, resp *raft.InstallSnapshotResponse, data io.Reader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstallSnapshot", id, target, args, resp, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// InstallSnapshot indicates an expected call of InstallSnapshot
func (mr *MockTransportMockRecorder) InstallSnapshot(id, target, args, resp, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstallSnapshot", reflect.TypeOf((*MockTransport)(nil).InstallSnapshot), id, target, args, resp, data)
}

// EncodePeer mocks base method
func (m *MockTransport) EncodePeer(id raft.ServerID, addr raft.ServerAddress) []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncodePeer", id, addr)
	ret0, _ := ret[0].([]byte)
	return ret0
}

// EncodePeer indicates an expected call of EncodePeer
func (mr *MockTransportMockRecorder) EncodePeer(id, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncodePeer", reflect.TypeOf((*MockTransport)(nil).EncodePeer), id, addr)
}

// DecodePeer mocks base method
func (m *MockTransport) DecodePeer(arg0 []byte) raft.ServerAddress {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecodePeer", arg0)
	ret0, _ := ret[0].(raft.ServerAddress)
	return ret0
}

// DecodePeer indicates an expected call of DecodePeer
func (mr *MockTransportMockRecorder) DecodePeer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecodePeer", reflect.TypeOf((*MockTransport)(nil).DecodePeer), arg0)
}

// SetHeartbeatHandler mocks base method
func (m *MockTransport) SetHeartbeatHandler(cb func(raft.RPC)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetHeartbeatHandler", cb)
}

// SetHeartbeatHandler indicates an expected call of SetHeartbeatHandler
func (mr *MockTransportMockRecorder) SetHeartbeatHandler(cb interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeartbeatHandler", reflect.TypeOf((*MockTransport)(nil).SetHeartbeatHandler), cb)
}

// TimeoutNow mocks base method
func (m *MockTransport) TimeoutNow(id raft.ServerID, target raft.ServerAddress, args *raft.TimeoutNowRequest, resp *raft.TimeoutNowResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TimeoutNow", id, target, args, resp)
	ret0, _ := ret[0].(error)
	return ret0
}

// TimeoutNow indicates an expected call of TimeoutNow
func (mr *MockTransportMockRecorder) TimeoutNow(id, target, args, resp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TimeoutNow", reflect.TypeOf((*MockTransport)(nil).TimeoutNow), id, target, args, resp)
}
