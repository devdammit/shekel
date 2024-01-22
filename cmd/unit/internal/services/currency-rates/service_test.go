package currency_rates_test

import (
	"context"
	mocks "github.com/devdammit/shekel/cmd/unit/internal/mocks/services/currency-rates"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/repositories/currency-rates"
	currency_rates "github.com/devdammit/shekel/cmd/unit/internal/services/currency-rates"
	"github.com/devdammit/shekel/pkg/currency"
	openexchange "github.com/devdammit/shekel/pkg/open-exchange"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestService_Convert(t *testing.T) {
	t.Run("should convert currency", func(t *testing.T) {
		var (
			ctrl = gomock.NewController(t)
			repo = mocks.NewMockRatesRepository(ctrl)
			api  = mocks.NewMockOpenExchangeRatesAPI(ctrl)
		)

		repo.EXPECT().GetCurrencyRateByDate(context.Background(), currency.USD, datetime.MustParseDateTime("2020-01-01 14:00")).Times(1).Return(float64(0), port.ErrRateNotFound)

		api.EXPECT().GetByDate(currency.USD, []currency.Code{currency.USD, currency.EUR, currency.THB, currency.RUB}, datetime.MustParseDate("2020-01-01")).Times(1).Return(&openexchange.HistoricalRates{
			Base: currency.USD,
			Rates: map[currency.Code]float64{
				currency.USD: 1,
				currency.EUR: 0.9,
				currency.THB: 30,
				currency.RUB: 60,
			},
		}, nil)

		repo.EXPECT().SetCurrencyRatesByDate(context.Background(), map[currency.Code]float64{
			currency.USD: 1,
			currency.EUR: 0.9,
			currency.THB: 30,
			currency.RUB: 60,
		}, datetime.MustParseDateTime("2020-01-01 14:00")).Times(1).Return(nil)

		repo.EXPECT().GetCurrencyRateByDate(context.Background(), currency.USD, datetime.MustParseDateTime("2020-01-01 14:00")).Times(1).Return(float64(1), nil)
		repo.EXPECT().GetCurrencyRateByDate(context.Background(), currency.RUB, datetime.MustParseDateTime("2020-01-01 14:00")).Times(1).Return(float64(60), nil)

		ctx := context.Background()
		service := currency_rates.NewService(repo, api)

		rate, err := service.Convert(ctx, currency.Amount{
			CurrencyCode: currency.USD,
			Value:        1,
		}, currency.RUB, datetime.MustParseDateTime("2020-01-01 14:00"))

		assert.NoError(t, err)
		assert.Equal(t, currency.Amount{
			CurrencyCode: currency.RUB,
			Value:        60,
		}, *rate)
	})
}
