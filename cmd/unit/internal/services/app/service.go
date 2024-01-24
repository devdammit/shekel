package app

import (
	"context"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	serviceport "github.com/devdammit/shekel/cmd/unit/internal/ports/services"
)

type PeriodsRepository interface {
	GetCount(ctx context.Context) (uint64, error)
	GetLast(ctx context.Context) (*entities.Period, error)
}

type Service struct {
	periods PeriodsRepository
}

func NewService(periodsRepo PeriodsRepository) *Service {
	return &Service{
		periods: periodsRepo,
	}
}

func (s *Service) GetInfo(ctx context.Context) (*serviceport.AppInfo, error) {
	res := &serviceport.AppInfo{
		Version: "0.0.1",
	}

	count, err := s.periods.GetCount(ctx)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		res.Initialized = true
	}

	if res.Initialized {
		res.ActivePeriod, err = s.periods.GetLast(ctx)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
