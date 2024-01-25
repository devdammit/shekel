package create_invoice

import (
	"context"
	"errors"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	"github.com/devdammit/shekel/internal/resources"
	"go.etcd.io/bbolt"
)

type invoicesRepository interface {
	CreateTx(tx *bbolt.Tx, invoice entities.Invoice) (*entities.Invoice, error)
}

type templatesRepository interface {
	CreateTx(tx *bbolt.Tx, template entities.InvoiceTemplate) (*entities.InvoiceTemplate, error)
}

type Repositories struct {
	invoices  invoicesRepository
	templates templatesRepository
}

type CreateInvoiceUoW struct {
	db *resources.Bolt

	invoices []entities.Invoice
	template *entities.InvoiceTemplate

	repositories Repositories
}

func NewUoW(db *resources.Bolt, repository invoicesRepository, templatesRepository templatesRepository) *CreateInvoiceUoW {
	return &CreateInvoiceUoW{
		db:       db,
		invoices: make([]entities.Invoice, 0),

		repositories: Repositories{
			invoices:  repository,
			templates: templatesRepository,
		},
	}
}

func (u *CreateInvoiceUoW) CreateInvoices(invoices []entities.Invoice, template entities.InvoiceTemplate) {
	u.template = &template
	u.invoices = invoices
}

func (u *CreateInvoiceUoW) CreateInvoice(invoice entities.Invoice) {
	u.invoices = append(u.invoices, invoice)
}

func (u *CreateInvoiceUoW) Commit(ctx context.Context) error {
	if u.invoices == nil {
		return errors.New("missing data")
	}

	return u.db.Update(func(tx *bbolt.Tx) error {
		if len(u.invoices) == 1 {
			_, err := u.repositories.invoices.CreateTx(tx, u.invoices[0])
			if err != nil {
				return err
			}

			return nil
		}

		if u.template == nil {
			return errors.New("missing data")
		}

		template, err := u.repositories.templates.CreateTx(tx, *u.template)
		if err != nil {
			return err
		}

		for _, invoice := range u.invoices {
			invoice.Template = template

			_, err = u.repositories.invoices.CreateTx(tx, invoice)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
