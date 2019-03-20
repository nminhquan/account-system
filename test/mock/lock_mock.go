package mytest_mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockLock is a mock of Lock interface
type MockLock struct {
	ctrl     *gomock.Controller
	recorder *MockLockMockRecorder
}

// MockLockMockRecorder is the mock recorder for MockLock
type MockLockMockRecorder struct {
	mock *MockLock
}

// NewMockLock creates a new mock instance
func NewMockLock(ctrl *gomock.Controller) *MockLock {
	mock := &MockLock{ctrl: ctrl}
	mock.recorder = &MockLockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLock) EXPECT() *MockLockMockRecorder {
	return m.recorder
}

// On mocks base method
func (m *MockLock) On() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "On")
	ret0, _ := ret[0].(bool)
	return ret0
}

// On indicates an expected call of On
func (mr *MockLockMockRecorder) On() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "On", reflect.TypeOf((*MockLock)(nil).On))
}

// Off mocks base method
func (m *MockLock) Off() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Off")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Off indicates an expected call of Off
func (mr *MockLockMockRecorder) Off() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Off", reflect.TypeOf((*MockLock)(nil).Off))
}

// Status mocks base method
func (m *MockLock) Status() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Status indicates an expected call of Status
func (mr *MockLockMockRecorder) Status() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockLock)(nil).Status))
}
