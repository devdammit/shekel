package initialize

import (
	"context"
	"errors"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	"github.com/devdammit/shekel/pkg/pointer"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

type PeriodsRepository interface {
	GetCount(ctx context.Context) (uint64, error)
}

type UnitOfWork interface {
	SetStartDate(date datetime.Date)
	CreatePeriod(period entities.Period) error
	Commit(ctx context.Context) error
}

type DateTimeProvider interface {
	Now() datetime.DateTime
}

type UseCase struct {
	periods  PeriodsRepository
	dateTime DateTimeProvider
	unit     UnitOfWork
}

func NewUseCase(
	periods PeriodsRepository,
	provider DateTimeProvider,
	unit UnitOfWork,
) *UseCase {
	return &UseCase{
		periods:  periods,
		dateTime: provider,
		unit:     unit,
	}
}

func (uc *UseCase) Execute(ctx context.Context, startDate datetime.Date) error {
	count, err := uc.periods.GetCount(ctx)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("periods already initialized")
	}

	uc.unit.SetStartDate(startDate)

	for date := startDate; date.Before(uc.dateTime.Now().Time); date = datetime.NewDate(date.AddDate(0, 1, 0)) {
		period := entities.Period{
			CreatedAt: datetime.NewDateTime(date.Time),
		}

		if date.AddDate(0, 1, 0).After(uc.dateTime.Now().Time) {
			period.ClosedAt = nil
		} else {
			period.ClosedAt = pointer.Ptr(datetime.NewDateTime(date.AddDate(0, 1, 0)))
		}

		err = uc.unit.CreatePeriod(period)
		if err != nil {
			return err
		}
	}

	err = uc.unit.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
