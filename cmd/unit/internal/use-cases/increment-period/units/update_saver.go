package units

import (
	"context"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
)

type AccountsRepository interface {
	ButchUpdate(ctx context.Context, accounts []entities.Account) error
}

type Transactor interface {
	Transaction(fn func() error) error
}

type ResultSaver struct {
	accountsRepo AccountsRepository
	periodsRepo  PeriodsRepository
	db           Transactor
}

func NewResultSaver(accountsRepo AccountsRepository, periodsRepo PeriodsRepository) *ResultSaver {
	return &ResultSaver{
		accountsRepo: accountsRepo,
		periodsRepo:  periodsRepo,
	}
}

func (u *ResultSaver) GetName() string {
	return "period_closer"
}

func (u *ResultSaver) Handle(ctx context.Context, _ *Request, payload *Payload) (*Payload, error) {
	if payload == nil || payload.ActivePeriod == nil || payload.Accounts == nil {
		return nil, ErrPayloadCheckFailed
	}

	err := u.db.Transaction(func() error {
		err := u.accountsRepo.ButchUpdate(ctx, payload.Accounts)
		if err != nil {
			return err
		}

		payload.ActivePeriod.Close()

		err = u.periodsRepo.Update(ctx, payload.ActivePeriod)
		if err != nil {
			return err
		}

		_, err = u.periodsRepo.Create(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return payload, nil
}
