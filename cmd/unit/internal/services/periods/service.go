package periods

import (
	"context"
	"errors"
	"github.com/devdammit/shekel/pkg/log"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"time"
)

type AppConfig interface {
	GetFinancialYearStart() time.Time
}

type Service struct {
	config AppConfig
}

func NewService(c AppConfig) *Service {
	return &Service{
		config: c,
	}
}

func (s *Service) GetEstimatedEndDate(_ context.Context) (datetime.DateTime, error) {
	startedAt := s.config.GetFinancialYearStart()

	log.With(log.String("started_at", startedAt.String())).Info("started at")

	for d := startedAt; d.Before(time.Now()); d = d.AddDate(0, 1, 0) {
		if time.Now().Before(d.AddDate(0, 1, 0)) {
			return datetime.NewDateTime(d.AddDate(0, 1, 0)), nil
		}
	}

	return datetime.NewDateTime(startedAt), errors.New("unable to calculate end date")
}
