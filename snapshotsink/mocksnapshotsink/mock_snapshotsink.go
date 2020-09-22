// Code generated by MockGen. DO NOT EDIT.
// Source: snapshotsink.go

// Package mocksnapshotsink is a generated GoMock package.
package mocksnapshotsink

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSnapshotSink is a mock of SnapshotSink interface
type MockSnapshotSink struct {
	ctrl     *gomock.Controller
	recorder *MockSnapshotSinkMockRecorder
}

// MockSnapshotSinkMockRecorder is the mock recorder for MockSnapshotSink
type MockSnapshotSinkMockRecorder struct {
	mock *MockSnapshotSink
}

// NewMockSnapshotSink creates a new mock instance
func NewMockSnapshotSink(ctrl *gomock.Controller) *MockSnapshotSink {
	mock := &MockSnapshotSink{ctrl: ctrl}
	mock.recorder = &MockSnapshotSinkMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSnapshotSink) EXPECT() *MockSnapshotSinkMockRecorder {
	return m.recorder
}

// Write mocks base method
func (m *MockSnapshotSink) Write(p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write
func (mr *MockSnapshotSinkMockRecorder) Write(p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockSnapshotSink)(nil).Write), p)
}

// Close mocks base method
func (m *MockSnapshotSink) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockSnapshotSinkMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockSnapshotSink)(nil).Close))
}

// ID mocks base method
func (m *MockSnapshotSink) ID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID
func (mr *MockSnapshotSinkMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockSnapshotSink)(nil).ID))
}

// Cancel mocks base method
func (m *MockSnapshotSink) Cancel() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cancel")
	ret0, _ := ret[0].(error)
	return ret0
}

// Cancel indicates an expected call of Cancel
func (mr *MockSnapshotSinkMockRecorder) Cancel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cancel", reflect.TypeOf((*MockSnapshotSink)(nil).Cancel))
}