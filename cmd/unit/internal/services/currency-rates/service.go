package currency_rates

import (
	"context"
	"errors"

	port "github.com/devdammit/shekel/cmd/unit/internal/ports/repositories/currency-rates"
	"github.com/devdammit/shekel/pkg/currency"
	openexchange "github.com/devdammit/shekel/pkg/open-exchange"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

var BaseCurrency = currency.USD
var SupportedCurrencies = []currency.Code{
	currency.USD,
	currency.EUR,
	currency.THB,
	currency.RUB,
}

type RatesRepository interface {
	SetCurrencyRatesByDate(ctx context.Context, rates map[currency.Code]float64, date datetime.DateTime) error
	GetCurrencyRateByDate(ctx context.Context, code currency.Code, date datetime.DateTime) (float64, error)
}

type OpenExchangeRatesAPI interface {
	GetByDate(base currency.Code, symbols []currency.Code, date datetime.Date) (*openexchange.HistoricalRates, error)
}

type Service struct {
	rates RatesRepository
	api   OpenExchangeRatesAPI
}

func NewService(rates RatesRepository, api OpenExchangeRatesAPI) *Service {
	return &Service{
		rates: rates,
		api:   api,
	}
}

func (s *Service) Convert(ctx context.Context, amount currency.Amount, to currency.Code, date datetime.DateTime) (*currency.Amount, error) {
	sourceRate, err := s.rates.GetCurrencyRateByDate(ctx, amount.CurrencyCode, date)
	if err != nil {
		if errors.Is(err, port.ErrRateNotFound) {
			err = s.fillRates(ctx, date)
			if err != nil {
				return nil, err
			}

			sourceRate, err = s.rates.GetCurrencyRateByDate(ctx, amount.CurrencyCode, date)

			if err != nil {
				return nil, err
			}
		}
	}

	targetRate, err := s.rates.GetCurrencyRateByDate(ctx, to, date)
	if err != nil {
		return nil, err
	}

	return &currency.Amount{
		CurrencyCode: to,
		Value:        amount.Value * (targetRate / sourceRate),
	}, nil
}

func (s *Service) fillRates(ctx context.Context, date datetime.DateTime) error {
	rates, err := s.api.GetByDate(BaseCurrency, SupportedCurrencies, datetime.MustParseDate(date.Format("2006-01-02")))
	if err != nil {
		return err
	}

	err = s.rates.SetCurrencyRatesByDate(ctx, rates.Rates, date)
	if err != nil {
		return err
	}

	return nil
}
