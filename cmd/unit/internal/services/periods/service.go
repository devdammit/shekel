package periods

import (
	"context"
	"errors"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	"github.com/devdammit/shekel/pkg/pointer"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

type Repository interface {
	Create(ctx context.Context, period entities.Period) (*entities.Period, error)
	GetCount(ctx context.Context) (uint64, error)
}

type AppConfig interface {
	GetStartYear() (datetime.DateTime, error)
}

type DateTimeProvider interface {
	Now() datetime.DateTime
}

type Service struct {
	periods   Repository
	appConfig AppConfig
	dateTime  DateTimeProvider
}

func NewService(periods Repository, appConfig AppConfig, provider DateTimeProvider) *Service {
	return &Service{
		periods:   periods,
		appConfig: appConfig,
		dateTime:  provider,
	}
}

func (s *Service) InitPeriods(ctx context.Context) error {
	count, err := s.periods.GetCount(ctx)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("periods already initialized")
	}

	startYearMonthFrom, err := s.appConfig.GetStartYear()
	if err != nil {
		return err
	}

	for date := startYearMonthFrom; date.Before(s.dateTime.Now().Time); date = datetime.NewDateTime(date.AddDate(0, 1, 0)) {
		period := entities.Period{
			CreatedAt: date,
		}

		if date.AddDate(0, 1, 0).After(s.dateTime.Now().Time) {
			period.ClosedAt = nil
		} else {
			period.ClosedAt = pointer.Ptr(datetime.NewDateTime(date.AddDate(0, 1, 0)))
		}

		_, err := s.periods.Create(ctx, period)
		if err != nil {
			return err
		}
	}

	return nil
}
