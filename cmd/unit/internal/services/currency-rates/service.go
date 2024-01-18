package currency_rates

import (
	"github.com/devdammit/shekel/pkg/currency"
	openexchange "github.com/devdammit/shekel/pkg/open-exchange"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

type RatesRepository interface {
	GetCurrencyRateByDate(source, target currency.Code, date datetime.DateTime) (float64, error)
}

type Service struct {
	rates RatesRepository
	api   openexchange.Client
}

func (s *Service) Convert(amount currency.Amount, to currency.Code, date datetime.DateTime) (*currency.Amount, error) {
	rates, err := s.rates.GetCurrencyRateByDate(amount.CurrencyCode, to, date)
	if err != nil {
		return nil, err
	}

	return &currency.Amount{
		CurrencyCode: to,
		Value:        amount.Value * rates,
	}, nil
}
