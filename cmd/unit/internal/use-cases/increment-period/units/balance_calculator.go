package units

import (
	"context"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
)

type AccountService interface {
	CalculateBalances(ctx context.Context, chunk []entities.Transaction) (map[uint64]entities.Account, error)
}

type BalanceCalculator struct {
	accountService AccountService
}

func NewBalanceCalculator(accountService AccountService) *BalanceCalculator {
	return &BalanceCalculator{
		accountService: accountService,
	}
}

func (u *BalanceCalculator) GetName() string {
	return "balance_calculator"
}

func (u *BalanceCalculator) Handle(ctx context.Context, request *Request, payload *Payload) (*Payload, error) {
	if payload == nil || payload.Transactions == nil {
		return nil, ErrPayloadCheckFailed
	}

	if len(payload.Transactions) == 0 {
		return payload, nil
	}

	mapAccounts, err := u.accountService.CalculateBalances(ctx, payload.Transactions)
	if err != nil {
		return nil, err
	}

	accounts := make([]entities.Account, 0, len(mapAccounts))

	for _, account := range mapAccounts {
		accounts = append(accounts, account)
	}

	payload.Accounts = accounts

	return payload, nil
}
