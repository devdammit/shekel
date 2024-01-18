package create_tx

import (
	"context"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

type TransactionsRepository interface {
	Create(ctx context.Context, tx entities.Transaction) (*entities.Transaction, error)
}

type AccountsRepository interface {
	GetByID(ctx context.Context, ID uint64) (*entities.Account, error)
}

type InvoicesRepository interface {
	GetByID(ctx context.Context, ID uint64) (*entities.Invoice, error)
}

type PeriodsRepository interface {
	GetLast() (*entities.Period, error)
}

type Service struct {
	periods      PeriodsRepository
	transactions TransactionsRepository
	accounts     AccountsRepository
	invoices     InvoicesRepository
}

func NewService(periodsRepo PeriodsRepository, transactionsRepo TransactionsRepository, accountsRepo AccountsRepository, invoicesRepo InvoicesRepository) *Service {
	return &Service{
		periods:      periodsRepo,
		transactions: transactionsRepo,
		accounts:     accountsRepo,
		invoices:     invoicesRepo,
	}
}

func (s *Service) Execute(ctx context.Context, request port.CreateTxRequest) error {
	period, err := s.periods.GetLast()
	if err != nil {
		return err
	}

	if period.ClosedAt != nil {
		return port.ErrorPeriodAlreadyClosed
	}

	if request.Date.Before(period.CreatedAt.Time) {
		return port.ErrorTxDateOutOfRange
	}

	if request.FromID == nil && request.ToID == nil {
		return port.ErrorTxNoAccounts
	}

	var fromAcc, toAcc *entities.Account

	if request.FromID != nil {
		fromAcc, err = s.accounts.GetByID(ctx, *request.FromID)
		if err != nil {
			return err
		}
	}

	if request.ToID != nil {
		toAcc, err = s.accounts.GetByID(ctx, *request.ToID)
		if err != nil {
			return err
		}
	}

	var invoice *entities.Invoice

	if request.InvoiceID != nil {
		invoice, err = s.invoices.GetByID(ctx, *request.InvoiceID)
		if err != nil {
			return err
		}
	}

	_, err = s.transactions.Create(ctx, entities.Transaction{
		Amount:    request.Amount,
		From:      fromAcc,
		To:        toAcc,
		Invoice:   invoice,
		CreatedAt: datetime.DateTime(request.Date),
	})

	if err != nil {
		return err
	}

	return nil
}
