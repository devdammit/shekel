package close_invoice

import (
	"context"
	"errors"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
)

type invoicesRepository interface {
	GetByID(ctx context.Context, id uint64) (*entities.Invoice, error)
	Update(ctx context.Context, invoice *entities.Invoice) (*entities.Invoice, error)
}

type periodsRepository interface {
	GetLast(ctx context.Context) (*entities.Period, error)
}

type UseCase struct {
	invoices invoicesRepository
	periods  periodsRepository
}

func NewUseCase(invoices invoicesRepository, repository periodsRepository) *UseCase {
	return &UseCase{
		invoices: invoices,
		periods:  repository,
	}
}

func (u *UseCase) Execute(ctx context.Context, ID uint64) error {
	invoice, err := u.invoices.GetByID(ctx, ID)
	if err != nil {
		return err
	}

	if invoice.Status == entities.InvoiceStatusPaid {
		return errors.New("invoice already paid")
	}

	period, err := u.periods.GetLast(ctx)
	if err != nil {
		return err
	}

	if period.CreatedAt.After(invoice.Date.Time) {
		return errors.New("cannot close invoice before current period")
	}

	invoice.Status = entities.InvoiceStatusPaid

	_, err = u.invoices.Update(ctx, invoice)
	if err != nil {
		return err
	}

	return nil
}
