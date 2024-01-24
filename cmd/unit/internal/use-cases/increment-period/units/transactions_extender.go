package units

import (
	"context"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
)

type TransactionsRepository interface {
	GetAllByPeriodID(periodID uint64) ([]entities.Transaction, error)
}

type TransactionsExtender struct {
	repo TransactionsRepository
}

func NewTransactionsExtender(repo TransactionsRepository) *TransactionsExtender {
	return &TransactionsExtender{repo: repo}
}

func (u *TransactionsExtender) GetName() string {
	return "transactions_extender"
}

func (u *TransactionsExtender) Handle(ctx context.Context, request *Request, payload *Payload) (*Payload, error) {
	if payload == nil || payload.ActivePeriod == nil {
		return nil, ErrPayloadCheckFailed
	}

	transactions, err := u.repo.GetAllByPeriodID(payload.ActivePeriod.ID)
	if err != nil {
		return nil, err
	}

	payload.Transactions = transactions

	return payload, nil
}
