package update_tx

import (
	"context"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

type PeriodsRepository interface {
	GetLast() (*entities.Period, error)
}

type AccountsRepository interface {
	GetByID(ctx context.Context, ID uint64) (*entities.Account, error)
}

type InvoicesRepository interface {
	GetByID(ctx context.Context, ID uint64) (*entities.Invoice, error)
}

type TransactionsRepository interface {
	Update(ctx context.Context, tx entities.Transaction) (*entities.Transaction, error)
}

type UpdateTxUseCase struct {
	periods      PeriodsRepository
	accounts     AccountsRepository
	invoices     InvoicesRepository
	transactions TransactionsRepository
}

func NewUpdateTxUseCase(periodsRepo PeriodsRepository, accountsRepo AccountsRepository, invoicesRepo InvoicesRepository, transactionsRepo TransactionsRepository) *UpdateTxUseCase {
	return &UpdateTxUseCase{
		periods:      periodsRepo,
		accounts:     accountsRepo,
		invoices:     invoicesRepo,
		transactions: transactionsRepo,
	}
}

func (uc *UpdateTxUseCase) Execute(ctx context.Context, request port.UpdateTxRequest) error {
	period, err := uc.periods.GetLast()
	if err != nil {
		return err
	}

	if period.ClosedAt != nil {
		return port.CannotUpdateTxAtClosedPeriod
	}

	if request.Date.Before(period.CreatedAt.Time) {
		return port.ErrorTxDateOutOfRange
	}

	if request.FromID == nil && request.ToID == nil {
		return port.ErrorTxNoAccounts
	}

	var fromAcc, toAcc *entities.Account

	if request.FromID != nil {
		fromAcc, err = uc.accounts.GetByID(ctx, *request.FromID)
		if err != nil {
			return err
		}
	}

	if request.ToID != nil {
		toAcc, err = uc.accounts.GetByID(ctx, *request.ToID)
		if err != nil {
			return err
		}
	}

	var invoice *entities.Invoice

	if request.InvoiceID != nil {
		invoice, err = uc.invoices.GetByID(ctx, *request.InvoiceID)
		if err != nil {
			return err
		}
	}

	_, err = uc.transactions.Update(ctx, entities.Transaction{
		ID:        request.ID,
		CreatedAt: datetime.DateTime(request.Date),
		Amount:    request.Amount,
		From:      fromAcc,
		To:        toAcc,
		Invoice:   invoice,
	})

	if err != nil {
		return err
	}

	return nil
}
