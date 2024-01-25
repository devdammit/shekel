package initialize

import (
	"context"
	"errors"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	"github.com/devdammit/shekel/internal/resources"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"go.etcd.io/bbolt"
)

type appConfigRepository interface {
	SetStartDateTx(ctx context.Context, tx *bbolt.Tx, date datetime.Date) error
}

type periodsRepository interface {
	CreateTx(ctx context.Context, tx *bbolt.Tx, period entities.Period) (*entities.Period, error)
}

type UoW struct {
	startDate *datetime.Date
	periods   []entities.Period

	db *resources.Bolt

	periodsRepo periodsRepository
	appConfig   appConfigRepository
}

func NewUoW(db *resources.Bolt, repository appConfigRepository, periodsRepository periodsRepository) *UoW {
	return &UoW{
		appConfig:   repository,
		db:          db,
		periodsRepo: periodsRepository,
		periods:     make([]entities.Period, 0),
	}
}

func (u *UoW) SetStartDate(date datetime.Date) {
	u.startDate = &date
}

func (u *UoW) CreatePeriod(period entities.Period) error {
	u.periods = append(u.periods, period)

	return nil
}

func (u *UoW) Commit(ctx context.Context) error {
	if u.startDate == nil || u.periods == nil {
		return errors.New("missing data")
	}

	return u.db.Update(func(tx *bbolt.Tx) error {
		err := u.appConfig.SetStartDateTx(ctx, tx, *u.startDate)
		if err != nil {
			return err
		}

		for _, period := range u.periods {
			_, err = u.periodsRepo.CreateTx(ctx, tx, period)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
