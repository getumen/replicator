// Code generated by MockGen. DO NOT EDIT.
// Source: snapshot_store.go

// Package mocksnapshotstore is a generated GoMock package.
package mocksnapshotstore

import (
	gomock "github.com/golang/mock/gomock"
	raft "github.com/hashicorp/raft"
	io "io"
	reflect "reflect"
)

// MockSnapshotStore is a mock of SnapshotStore interface
type MockSnapshotStore struct {
	ctrl     *gomock.Controller
	recorder *MockSnapshotStoreMockRecorder
}

// MockSnapshotStoreMockRecorder is the mock recorder for MockSnapshotStore
type MockSnapshotStoreMockRecorder struct {
	mock *MockSnapshotStore
}

// NewMockSnapshotStore creates a new mock instance
func NewMockSnapshotStore(ctrl *gomock.Controller) *MockSnapshotStore {
	mock := &MockSnapshotStore{ctrl: ctrl}
	mock.recorder = &MockSnapshotStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSnapshotStore) EXPECT() *MockSnapshotStoreMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockSnapshotStore) Create(version raft.SnapshotVersion, index, term uint64, configuration raft.Configuration, configurationIndex uint64, trans raft.Transport) (raft.SnapshotSink, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", version, index, term, configuration, configurationIndex, trans)
	ret0, _ := ret[0].(raft.SnapshotSink)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockSnapshotStoreMockRecorder) Create(version, index, term, configuration, configurationIndex, trans interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSnapshotStore)(nil).Create), version, index, term, configuration, configurationIndex, trans)
}

// List mocks base method
func (m *MockSnapshotStore) List() ([]*raft.SnapshotMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]*raft.SnapshotMeta)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockSnapshotStoreMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockSnapshotStore)(nil).List))
}

// Open mocks base method
func (m *MockSnapshotStore) Open(id string) (*raft.SnapshotMeta, io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open", id)
	ret0, _ := ret[0].(*raft.SnapshotMeta)
	ret1, _ := ret[1].(io.ReadCloser)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Open indicates an expected call of Open
func (mr *MockSnapshotStoreMockRecorder) Open(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockSnapshotStore)(nil).Open), id)
}