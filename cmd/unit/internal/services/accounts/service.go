package accounts

import (
	"context"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	"github.com/devdammit/shekel/pkg/currency"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

type Repository interface {
	GetAll() ([]entities.Account, error)
}

type CurrencyRates interface {
	Convert(amount currency.Amount, to currency.Code, date datetime.DateTime) (*currency.Amount, error)
}

type Service struct {
	repo          Repository
	currencyRates CurrencyRates
}

func NewService(repo Repository, currencyRates CurrencyRates) *Service {
	return &Service{
		repo:          repo,
		currencyRates: currencyRates,
	}
}

// CalculateBalances calculates the balances of all accounts affected by the given transactions.
func (s *Service) CalculateBalances(_ context.Context, chunk []entities.Transaction) (map[uint64]entities.Account, error) {
	accounts, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	mapAccounts := make(map[uint64]currency.Amount, len(accounts))
	affectedAccounts := make(map[uint64]struct{}, len(accounts))

	for _, account := range accounts {
		mapAccounts[account.ID] = account.Balance
	}

	for _, transaction := range chunk {
		if transaction.From != nil {
			acc := mapAccounts[transaction.From.ID]

			affectedAccounts[transaction.From.ID] = struct{}{}

			var value currency.Amount

			if acc.CurrencyCode != transaction.Amount.CurrencyCode {
				convertedAmount, err := s.currencyRates.Convert(transaction.Amount, acc.CurrencyCode, transaction.CreatedAt)
				if err != nil {
					return nil, err
				}

				value = acc.Subtract(*convertedAmount)
			} else {
				value = acc.Subtract(transaction.Amount)
			}

			mapAccounts[transaction.From.ID] = value
		}

		if transaction.To != nil {
			acc := mapAccounts[transaction.To.ID]
			affectedAccounts[transaction.To.ID] = struct{}{}

			var value currency.Amount

			if acc.CurrencyCode != transaction.Amount.CurrencyCode {
				convertedAmount, err := s.currencyRates.Convert(transaction.Amount, acc.CurrencyCode, transaction.CreatedAt)
				if err != nil {
					return nil, err
				}

				value = acc.Add(*convertedAmount)
			} else {
				value = acc.Add(transaction.Amount)
			}

			mapAccounts[transaction.To.ID] = value
		}
	}

	result := make(map[uint64]entities.Account, len(accounts))

	for _, account := range accounts {
		if _, ok := affectedAccounts[account.ID]; ok {
			account.Balance = mapAccounts[account.ID]

			result[account.ID] = account
		}
	}

	return result, nil
}
