// Code generated by MockGen. DO NOT EDIT.
// Source: resource_manager/account_service.go

// Package mock_resource_manager is a generated GoMock package.
package mock_resource_manager

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockConsensusService is a mock of ConsensusService interface
type MockConsensusService struct {
	ctrl     *gomock.Controller
	recorder *MockConsensusServiceMockRecorder
}

// MockConsensusServiceMockRecorder is the mock recorder for MockConsensusService
type MockConsensusServiceMockRecorder struct {
	mock *MockConsensusService
}

// NewMockConsensusService creates a new mock instance
func NewMockConsensusService(ctrl *gomock.Controller) *MockConsensusService {
	mock := &MockConsensusService{ctrl: ctrl}
	mock.recorder = &MockConsensusServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConsensusService) EXPECT() *MockConsensusServiceMockRecorder {
	return m.recorder
}

// MockAccountService is a mock of AccountService interface
type MockAccountService struct {
	ctrl     *gomock.Controller
	recorder *MockAccountServiceMockRecorder
}

// MockAccountServiceMockRecorder is the mock recorder for MockAccountService
type MockAccountServiceMockRecorder struct {
	mock *MockAccountService
}

// NewMockAccountService creates a new mock instance
func NewMockAccountService(ctrl *gomock.Controller) *MockAccountService {
	mock := &MockAccountService{ctrl: ctrl}
	mock.recorder = &MockAccountServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAccountService) EXPECT() *MockAccountServiceMockRecorder {
	return m.recorder
}

// Start mocks base method
func (m *MockAccountService) Start() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start")
}

// Start indicates an expected call of Start
func (mr *MockAccountServiceMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockAccountService)(nil).Start))
}

// CreateAccount mocks base method
func (m *MockAccountService) CreateAccount(arg0 string, arg1 float64) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", arg0, arg1)
	ret0, _ := ret[0].(string)
	return ret0
}

// CreateAccount indicates an expected call of CreateAccount
func (mr *MockAccountServiceMockRecorder) CreateAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockAccountService)(nil).CreateAccount), arg0, arg1)
}

// ProcessSendPayment mocks base method
func (m *MockAccountService) ProcessSendPayment(arg0, arg1 string, arg2 float64) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessSendPayment", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	return ret0
}

// ProcessSendPayment indicates an expected call of ProcessSendPayment
func (mr *MockAccountServiceMockRecorder) ProcessSendPayment(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessSendPayment", reflect.TypeOf((*MockAccountService)(nil).ProcessSendPayment), arg0, arg1, arg2)
}

// ProcessReceivePayment mocks base method
func (m *MockAccountService) ProcessReceivePayment(arg0, arg1 string, arg2 float64) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessReceivePayment", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	return ret0
}

// ProcessReceivePayment indicates an expected call of ProcessReceivePayment
func (mr *MockAccountServiceMockRecorder) ProcessReceivePayment(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessReceivePayment", reflect.TypeOf((*MockAccountService)(nil).ProcessReceivePayment), arg0, arg1, arg2)
}

// ProcessRollbackPayment mocks base method
func (m *MockAccountService) ProcessRollbackPayment(accountNum string, amount float64) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessRollbackPayment", accountNum, amount)
	ret0, _ := ret[0].(string)
	return ret0
}

// ProcessRollbackPayment indicates an expected call of ProcessRollbackPayment
func (mr *MockAccountServiceMockRecorder) ProcessRollbackPayment(accountNum, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessRollbackPayment", reflect.TypeOf((*MockAccountService)(nil).ProcessRollbackPayment), accountNum, amount)
}

// Propose mocks base method
func (m *MockAccountService) Propose(arg0 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Propose", arg0)
}

// Propose indicates an expected call of Propose
func (mr *MockAccountServiceMockRecorder) Propose(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Propose", reflect.TypeOf((*MockAccountService)(nil).Propose), arg0)
}

// ReadCommits mocks base method
func (m *MockAccountService) ReadCommits(arg0 <-chan *string, arg1 <-chan error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReadCommits", arg0, arg1)
}

// ReadCommits indicates an expected call of ReadCommits
func (mr *MockAccountServiceMockRecorder) ReadCommits(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadCommits", reflect.TypeOf((*MockAccountService)(nil).ReadCommits), arg0, arg1)
}
