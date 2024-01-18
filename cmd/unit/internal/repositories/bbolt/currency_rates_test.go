package bbolt_test

import (
	"errors"
	mocks "github.com/devdammit/shekel/cmd/unit/internal/mocks/repositories/bbolt"
	apirepo "github.com/devdammit/shekel/cmd/unit/internal/ports/apis/open-exchange"
	"github.com/devdammit/shekel/cmd/unit/internal/repositories/bbolt"
	"github.com/devdammit/shekel/pkg/currency"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestCurrencyRatesRepository_GetCurrencyRateByDate(t *testing.T) {
	t.Run("should return error if bucket does not exist", func(t *testing.T) {
		var (
			mockController = gomock.NewController(t)
			openExchange   = mocks.NewMockOpenExchangeRatesAPI(mockController)
		)

		openExchange.EXPECT().GetByDate(currency.USD, []currency.Code{currency.USD, currency.EUR, currency.THB, currency.RUB}, datetime.MustParseDate("2020-01-01")).Times(1).Return(nil, errors.New("error"))

		repo := bbolt.NewCurrencyRatesRepository(nil, openExchange)

		_, err := repo.GetCurrencyRateByDate(currency.USD, currency.RUB, datetime.MustParseDate("2020-01-01"))

		assert.Error(t, err)
	})

	t.Run("should return rates", func(t *testing.T) {
		var (
			mockController = gomock.NewController(t)
			openExchange   = mocks.NewMockOpenExchangeRatesAPI(mockController)
		)

		openExchange.EXPECT().GetByDate(currency.USD, []currency.Code{currency.USD, currency.EUR, currency.THB, currency.RUB}, datetime.MustParseDate("2020-01-01")).Times(1).Return(&apirepo.HistoricalRates{
			Base: currency.USD,
			Rates: map[currency.Code]float64{
				currency.USD: 1,
				currency.EUR: 0.9,
				currency.THB: 30,
				currency.RUB: 60,
			},
		}, nil)

		repo := bbolt.NewCurrencyRatesRepository(nil, openExchange)

		rates, err := repo.GetCurrencyRateByDate(currency.USD, currency.RUB, datetime.MustParseDate("2020-01-01"))

		assert.NoError(t, err)
		assert.Equal(t, 60.0, rates)
	})

	t.Run("should memoize rates", func(t *testing.T) {
		var (
			mockController = gomock.NewController(t)
			openExchange   = mocks.NewMockOpenExchangeRatesAPI(mockController)
		)

		openExchange.EXPECT().GetByDate(currency.USD, []currency.Code{currency.USD, currency.EUR, currency.THB, currency.RUB}, datetime.MustParseDate("2020-01-01")).Times(1).Return(&apirepo.HistoricalRates{
			Base: currency.USD,
			Rates: map[currency.Code]float64{
				currency.USD: 1,
				currency.EUR: 0.9,
				currency.THB: 30,
				currency.RUB: 60,
			},
		}, nil)

		repo := bbolt.NewCurrencyRatesRepository(nil, openExchange)

		_, err := repo.GetCurrencyRateByDate(currency.USD, currency.RUB, datetime.MustParseDate("2020-01-01"))
		assert.NoError(t, err)

		rates, err := repo.GetCurrencyRateByDate(currency.USD, currency.RUB, datetime.MustParseDate("2020-01-01"))

		assert.NoError(t, err)
		assert.Equal(t, 60.0, rates)
	})
}
