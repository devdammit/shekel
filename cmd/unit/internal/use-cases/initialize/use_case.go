package initialize

import (
	"context"
	"errors"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	"github.com/devdammit/shekel/pkg/pointer"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

type PeriodsRepository interface {
	Create(ctx context.Context, period entities.Period) (*entities.Period, error)
	GetCount(ctx context.Context) (uint64, error)
}

type AppConfig interface {
	SetStartDate(ctx context.Context, date datetime.Date) error
}

type DateTimeProvider interface {
	Now() datetime.DateTime
}

type UseCase struct {
	periods   PeriodsRepository
	appConfig AppConfig
	dateTime  DateTimeProvider
}

func NewUseCase(periods PeriodsRepository, appConfig AppConfig, provider DateTimeProvider) *UseCase {
	return &UseCase{
		periods:   periods,
		appConfig: appConfig,
		dateTime:  provider,
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

	err = uc.appConfig.SetStartDate(ctx, startDate)
	if err != nil {
		return err
	}

	for date := startDate; date.Before(uc.dateTime.Now().Time); date = datetime.NewDate(date.AddDate(0, 1, 0)) {
		period := entities.Period{
			CreatedAt: datetime.NewDateTime(date.Time),
		}

		if date.AddDate(0, 1, 0).After(uc.dateTime.Now().Time) {
			period.ClosedAt = nil
		} else {
			period.ClosedAt = pointer.Ptr(datetime.NewDateTime(date.AddDate(0, 1, 0)))
		}

		_, err := uc.periods.Create(ctx, period)
		if err != nil {
			return err
		}
	}

	return nil
}
