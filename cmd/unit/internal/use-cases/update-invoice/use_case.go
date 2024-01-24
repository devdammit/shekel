package update_invoice

import (
	"context"
	"errors"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
	"github.com/devdammit/shekel/pkg/log"
)

type PeriodsRepository interface {
	GetLast(ctx context.Context) (*entities.Period, error)
}

type InvoicesRepository interface {
	GetByID(ctx context.Context, id uint64) (*entities.Invoice, error)
	GetByTemplateID(ctx context.Context, id uint64) ([]entities.Invoice, error)
	Update(ctx context.Context, invoice *entities.Invoice) (*entities.Invoice, error)
	Delete(ctx context.Context, id uint64) error
	BulkCreate(ctx context.Context, invoices []entities.Invoice) ([]entities.Invoice, error)
}

type InvoicesTemplateRepository interface {
	Create(ctx context.Context, template *entities.InvoiceTemplate) (*entities.InvoiceTemplate, error)
	Delete(ctx context.Context, id uint64) error
}

type ContactsRepository interface {
	GetByID(ctx context.Context, id uint64) (*entities.Contact, error)
}

type InvoicesService interface {
	GetScheduledInvoices(ctx context.Context, template entities.InvoiceTemplate) ([]entities.Invoice, error)
}

type Transactor interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type CalendarService interface {
	Sync(ctx context.Context) error
}

type Logger interface {
	Warn(ctx context.Context, fields ...log.Field)
}

type UpdateInvoiceUseCase struct {
	periods    PeriodsRepository
	invoices   InvoicesRepository
	templates  InvoicesTemplateRepository
	contacts   ContactsRepository
	transactor Transactor
	service    InvoicesService
	calendar   CalendarService
	logger     Logger
}

func NewUpdateInvoiceUseCase(periods PeriodsRepository, invoices InvoicesRepository, templates InvoicesTemplateRepository, contacts ContactsRepository, transactor Transactor, service InvoicesService, calendar CalendarService, logger Logger) *UpdateInvoiceUseCase {
	return &UpdateInvoiceUseCase{
		periods:    periods,
		invoices:   invoices,
		templates:  templates,
		contacts:   contacts,
		transactor: transactor,
		service:    service,
		calendar:   calendar,
		logger:     logger,
	}
}

func (u *UpdateInvoiceUseCase) Execute(ctx context.Context, req *port.UpdateInvoiceRequest) error {
	period, err := u.periods.GetLast(ctx)
	if err != nil {
		return err
	}

	if period.ClosedAt != nil {
		return errors.New("cannot update invoice at closed period")
	}

	if req.Date.Before(period.CreatedAt.Time) {
		return errors.New("cannot update invoice at previous period")
	}

	invoice, err := u.invoices.GetByID(ctx, req.InvoiceID)
	if err != nil {
		return err
	}

	if invoice.Date.Before(period.CreatedAt.Time) {
		return errors.New("cannot update invoice at previous period")
	}

	if invoice.Status == entities.InvoiceStatusPaid {
		return errors.New("cannot update paid invoice")
	}

	contact, err := u.contacts.GetByID(ctx, req.ContactID)
	if err != nil {
		return err
	}

	invoice.Name = req.Name
	invoice.Desc = req.Description
	invoice.Type = req.Type
	invoice.Amount = req.Amount
	invoice.Contact = contact
	invoice.Date = req.Date

	if req.Plan == nil {
		invoice.Template = nil

		_, err := u.invoices.Update(ctx, invoice)
		if err != nil {
			return err
		}

		err = u.calendar.Sync(ctx)
		if err != nil {
			u.logger.Warn(ctx, log.Err(err), log.String("message", "cannot sync calendar"))
		}

		return nil
	}

	invoices, err := u.invoices.GetByTemplateID(ctx, invoice.Template.ID)
	if err != nil {
		return err
	}

	err = u.transactor.Transaction(ctx, func(ctx context.Context) error {
		newTemplate, err := u.templates.Create(ctx, &entities.InvoiceTemplate{
			Name:      invoice.Name,
			Desc:      invoice.Desc,
			Type:      invoice.Type,
			Amount:    invoice.Amount,
			ContactID: invoice.Contact.ID,

			Date: invoice.Date,

			RepeatPlanner: &entities.RepeatPlanner{
				Interval:      req.Plan.Interval,
				IntervalCount: req.Plan.IntervalCount,
				DaysOfWeek:    req.Plan.DaysOfWeek,
				EndDate:       req.Plan.EndDate,
				EndCount:      req.Plan.EndCount,
			},
		})
		if err != nil {
			return err
		}

		err = u.templates.Delete(ctx, invoice.Template.ID)
		if err != nil {
			return err
		}

		deletingInvoices := make([]entities.Invoice, 0, len(invoices))

		for _, inv := range invoices {
			if req.Date.After(inv.Date.Time) {
				deletingInvoices = append(deletingInvoices, inv)
			}
		}

		for _, inv := range deletingInvoices {
			err = u.invoices.Delete(ctx, inv.ID)
			if err != nil {
				return err
			}
		}

		scheduledInvoices, err := u.service.GetScheduledInvoices(ctx, *newTemplate)
		if err != nil {
			return err
		}

		_, err = u.invoices.BulkCreate(ctx, scheduledInvoices)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	err = u.calendar.Sync(ctx)
	if err != nil {
		u.logger.Warn(ctx, log.Err(err), log.String("message", "cannot sync calendar"))
	}

	return nil
}
