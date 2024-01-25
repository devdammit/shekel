package update_invoice

import (
	"context"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	"github.com/devdammit/shekel/internal/resources"
	"go.etcd.io/bbolt"
)

type templatesRepository interface {
	CreateTx(tx *bbolt.Tx, template entities.InvoiceTemplate) (*entities.InvoiceTemplate, error)
	DeleteTx(tx *bbolt.Tx, id uint64) error
}

type invoicesRepository interface {
	CreateTx(tx *bbolt.Tx, invoice entities.Invoice) (*entities.Invoice, error)
	UpdateTx(tx *bbolt.Tx, invoice entities.Invoice) (*entities.Invoice, error)
	DeleteTx(tx *bbolt.Tx, id uint64) error
}

type repositories struct {
	templates templatesRepository
	invoices  invoicesRepository
}

type UoW struct {
	db *resources.Bolt

	createdInvoices  []entities.Invoice
	createdTemplates []entities.InvoiceTemplate
	updatedInvoices  []entities.Invoice
	deletedTemplates []uint64
	deleteInvoices   []uint64

	repositories repositories
}

func NewUoW(db *resources.Bolt, templates templatesRepository, invoices invoicesRepository) *UoW {
	return &UoW{
		db: db,

		createdInvoices:  make([]entities.Invoice, 0),
		createdTemplates: make([]entities.InvoiceTemplate, 0),
		updatedInvoices:  make([]entities.Invoice, 0),
		deletedTemplates: make([]uint64, 0),
		deleteInvoices:   make([]uint64, 0),

		repositories: repositories{
			templates: templates,
			invoices:  invoices,
		},
	}
}

func (u *UoW) CreateTemplate(template entities.InvoiceTemplate) {
	u.createdTemplates = append(u.createdTemplates, template)
}

func (u *UoW) DeleteTemplate(id uint64) {
	u.deletedTemplates = append(u.deletedTemplates, id)
}

func (u *UoW) DeleteInvoice(id uint64) {
	u.deleteInvoices = append(u.deleteInvoices, id)
}

func (u *UoW) UpdateInvoice(invoice entities.Invoice) {
	u.updatedInvoices = append(u.updatedInvoices, invoice)
}

func (u *UoW) CreateInvoices(invoices []entities.Invoice, template entities.InvoiceTemplate) {
	u.createdTemplates = append(u.createdTemplates, template)
	u.createdInvoices = append(u.updatedInvoices, invoices...)
}

func (u *UoW) Commit(_ context.Context) error {
	return u.db.Update(func(tx *bbolt.Tx) error {
		for _, template := range u.createdTemplates {
			_, err := u.repositories.templates.CreateTx(tx, template)
			if err != nil {
				return err
			}
		}

		for _, id := range u.deletedTemplates {
			err := u.repositories.templates.DeleteTx(tx, id)
			if err != nil {
				return err
			}
		}

		for _, id := range u.deleteInvoices {
			err := u.repositories.invoices.DeleteTx(tx, id)
			if err != nil {
				return err
			}
		}

		for _, invoice := range u.updatedInvoices {
			_, err := u.repositories.invoices.UpdateTx(tx, invoice)
			if err != nil {
				return err
			}
		}

		for _, invoice := range u.createdInvoices {
			_, err := u.repositories.invoices.CreateTx(tx, invoice)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
