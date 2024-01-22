// Code generated by MockGen. DO NOT EDIT.
// Source: internal/services/currency-rates/service.go
//
// Generated by this command:
//
//	mockgen -source=internal/services/currency-rates/service.go -destination=internal/mocks/services/currency-rates/mocks.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	currency "github.com/devdammit/shekel/pkg/currency"
	open_exchange "github.com/devdammit/shekel/pkg/open-exchange"
	datetime "github.com/devdammit/shekel/pkg/types/datetime"
	gomock "go.uber.org/mock/gomock"
)

// MockRatesRepository is a mock of RatesRepository interface.
type MockRatesRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRatesRepositoryMockRecorder
}

// MockRatesRepositoryMockRecorder is the mock recorder for MockRatesRepository.
type MockRatesRepositoryMockRecorder struct {
	mock *MockRatesRepository
}

// NewMockRatesRepository creates a new mock instance.
func NewMockRatesRepository(ctrl *gomock.Controller) *MockRatesRepository {
	mock := &MockRatesRepository{ctrl: ctrl}
	mock.recorder = &MockRatesRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRatesRepository) EXPECT() *MockRatesRepositoryMockRecorder {
	return m.recorder
}

// GetCurrencyRateByDate mocks base method.
func (m *MockRatesRepository) GetCurrencyRateByDate(ctx context.Context, code currency.Code, date datetime.DateTime) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrencyRateByDate", ctx, code, date)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrencyRateByDate indicates an expected call of GetCurrencyRateByDate.
func (mr *MockRatesRepositoryMockRecorder) GetCurrencyRateByDate(ctx, code, date any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrencyRateByDate", reflect.TypeOf((*MockRatesRepository)(nil).GetCurrencyRateByDate), ctx, code, date)
}

// SetCurrencyRatesByDate mocks base method.
func (m *MockRatesRepository) SetCurrencyRatesByDate(ctx context.Context, rates map[currency.Code]float64, date datetime.DateTime) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCurrencyRatesByDate", ctx, rates, date)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCurrencyRatesByDate indicates an expected call of SetCurrencyRatesByDate.
func (mr *MockRatesRepositoryMockRecorder) SetCurrencyRatesByDate(ctx, rates, date any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCurrencyRatesByDate", reflect.TypeOf((*MockRatesRepository)(nil).SetCurrencyRatesByDate), ctx, rates, date)
}

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
