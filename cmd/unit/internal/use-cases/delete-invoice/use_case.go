package delete_invoice

import (
	"context"
	"errors"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
)

type InvoicesRepository interface {
	GetByID(ctx context.Context, id uint64) (*entities.Invoice, error)
	Delete(ctx context.Context, id uint64) error
	GetByTemplateID(ctx context.Context, templateID uint64) ([]entities.Invoice, error)
}

type InvoicesTemplatesRepository interface {
	Delete(ctx context.Context, id uint64) error
}

type UseCase struct {
	invoices  InvoicesRepository
	templates InvoicesTemplatesRepository
}

func NewUseCase(invoices InvoicesRepository, templates InvoicesTemplatesRepository) *UseCase {
	return &UseCase{
		invoices:  invoices,
		templates: templates,
	}
}

func (u *UseCase) Execute(ctx context.Context, id uint64, single bool) error {
	invoice, err := u.invoices.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if invoice.Status == entities.InvoiceStatusPaid {
		return errors.New("cannot delete paid invoice")
	}

	if single {
		return u.invoices.Delete(ctx, id)
	}

	invoices, err := u.invoices.GetByTemplateID(ctx, invoice.Template.ID)
	if err != nil {
		return err
	}

	deletingInvoices := make([]uint64, 0, len(invoices))

	for _, inv := range invoices {
		if inv.Date.Equal(invoice.Date.Time) || inv.Date.After(invoice.Date.Time) {
			if inv.Status == entities.InvoiceStatusPending {
				deletingInvoices = append(deletingInvoices, inv.ID)
			}
		}
	}

	for _, id := range deletingInvoices {
		err = u.invoices.Delete(ctx, id)
		if err != nil {
			return err
		}
	}

	err = u.templates.Delete(ctx, invoice.Template.ID)
	if err != nil {
		return err
	}

	return nil
}
