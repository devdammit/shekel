// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repositories/bbolt/currency_rates.go
//
// Generated by this command:
//
//	mockgen -source=internal/repositories/bbolt/currency_rates.go -destination=internal/mocks/repositories/bbolt/currency_rates.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	currency "github.com/devdammit/shekel/pkg/currency"
	open_exchange "github.com/devdammit/shekel/pkg/open-exchange"
	datetime "github.com/devdammit/shekel/pkg/types/datetime"
	gomock "go.uber.org/mock/gomock"
)

// MockOpenExchangeRatesAPI is a mock of OpenExchangeRatesAPI interface.
type MockOpenExchangeRatesAPI struct {
	ctrl     *gomock.Controller
	recorder *MockOpenExchangeRatesAPIMockRecorder
}

// MockOpenExchangeRatesAPIMockRecorder is the mock recorder for MockOpenExchangeRatesAPI.
type MockOpenExchangeRatesAPIMockRecorder struct {
	mock *MockOpenExchangeRatesAPI
}

// NewMockOpenExchangeRatesAPI creates a new mock instance.
func NewMockOpenExchangeRatesAPI(ctrl *gomock.Controller) *MockOpenExchangeRatesAPI {
	mock := &MockOpenExchangeRatesAPI{ctrl: ctrl}
	mock.recorder = &MockOpenExchangeRatesAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOpenExchangeRatesAPI) EXPECT() *MockOpenExchangeRatesAPIMockRecorder {
	return m.recorder
}

// GetByDate mocks base method.
func (m *MockOpenExchangeRatesAPI) GetByDate(base currency.Code, symbols []currency.Code, date datetime.Date) (*open_exchange.HistoricalRates, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByDate", base, symbols, date)
	ret0, _ := ret[0].(*open_exchange.HistoricalRates)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByDate indicates an expected call of GetByDate.
func (mr *MockOpenExchangeRatesAPIMockRecorder) GetByDate(base, symbols, date any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByDate", reflect.TypeOf((*MockOpenExchangeRatesAPI)(nil).GetByDate), base, symbols, date)
}
