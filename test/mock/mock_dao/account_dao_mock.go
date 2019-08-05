// Code generated by MockGen. DO NOT EDIT.
// Source: dao/account_dao.go

// Package mock_dao is a generated GoMock package.
package mock_dao

import (
	reflect "reflect"

	model "github.com/nminhquan/mas/model"

	gomock "github.com/golang/mock/gomock"
)

// MockAccountDAO is a mock of AccountDAO interface
type MockAccountDAO struct {
	ctrl     *gomock.Controller
	recorder *MockAccountDAOMockRecorder
}

// MockAccountDAOMockRecorder is the mock recorder for MockAccountDAO
type MockAccountDAOMockRecorder struct {
	mock *MockAccountDAO
}

// NewMockAccountDAO creates a new mock instance
func NewMockAccountDAO(ctrl *gomock.Controller) *MockAccountDAO {
	mock := &MockAccountDAO{ctrl: ctrl}
	mock.recorder = &MockAccountDAOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAccountDAO) EXPECT() *MockAccountDAOMockRecorder {
	return m.recorder
}

// GetAccount mocks base method
func (m *MockAccountDAO) GetAccount(accountNumber string) *model.AccountInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", accountNumber)
	ret0, _ := ret[0].(*model.AccountInfo)
	return ret0
}

// GetAccount indicates an expected call of GetAccount
func (mr *MockAccountDAOMockRecorder) GetAccount(accountNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockAccountDAO)(nil).GetAccount), accountNumber)
}

// CreateAccount mocks base method
func (m *MockAccountDAO) CreateAccount(accInfo model.AccountInfo) int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", accInfo)
	ret0, _ := ret[0].(int64)
	return ret0
}

// CreateAccount indicates an expected call of CreateAccount
func (mr *MockAccountDAOMockRecorder) CreateAccount(accInfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockAccountDAO)(nil).CreateAccount), accInfo)
}

// DeductMoney mocks base method
func (m *MockAccountDAO) DeductMoney(accountNumber string, amount float64) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeductMoney", accountNumber, amount)
	ret0, _ := ret[0].(bool)
	return ret0
}

// DeductMoney indicates an expected call of DeductMoney
func (mr *MockAccountDAOMockRecorder) DeductMoney(accountNumber, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeductMoney", reflect.TypeOf((*MockAccountDAO)(nil).DeductMoney), accountNumber, amount)
}

// DepositMoney mocks base method
func (m *MockAccountDAO) DepositMoney(accountNumber string, amount float64) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DepositMoney", accountNumber, amount)
	ret0, _ := ret[0].(bool)
	return ret0
}

// DepositMoney indicates an expected call of DepositMoney
func (mr *MockAccountDAOMockRecorder) DepositMoney(accountNumber, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DepositMoney", reflect.TypeOf((*MockAccountDAO)(nil).DepositMoney), accountNumber, amount)
}
