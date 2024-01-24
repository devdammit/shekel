package units

import (
	"context"

	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
)

type PeriodExtender struct {
	repo PeriodsRepository
}

func NewPeriodExtender(repo PeriodsRepository) *PeriodExtender {
	return &PeriodExtender{
		repo: repo,
	}
}

func (u *PeriodExtender) GetName() string {
	return "period_extender"
}

func (u *PeriodExtender) Handle(ctx context.Context, request *Request, payload *Payload) (*Payload, error) {
	period, err := u.repo.GetLast(ctx)
	if err != nil {
		return nil, err
	}

	if period.ClosedAt != nil {
		return nil, port.ErrorPeriodAlreadyClosed
	}

	payload.ActivePeriod = period

	return payload, nil
}
