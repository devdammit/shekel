package delete_tx

import (
	"context"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
)

type PeriodsRepository interface {
	GetLast(ctx context.Context) (*entities.Period, error)
}

type TransactionsRepository interface {
	Delete(ctx context.Context, ID uint64) error
}

type DeleteTxUseCase struct {
	periods      PeriodsRepository
	transactions TransactionsRepository
}

func NewDeleteTxUseCase(periodsRepo PeriodsRepository, transactionsRepo TransactionsRepository) *DeleteTxUseCase {
	return &DeleteTxUseCase{
		periods:      periodsRepo,
		transactions: transactionsRepo,
	}
}

func (uc *DeleteTxUseCase) Execute(ctx context.Context, ID uint64) error {
	period, err := uc.periods.GetLast(ctx)
	if err != nil {
		return err
	}

	if period.ClosedAt != nil {
		return port.ErrCannotDeleteTxAtClosedPeriod
	}

	return uc.transactions.Delete(ctx, ID)
}
