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
}

type ContactsRepository interface {
	GetByID(ctx context.Context, id uint64) (*entities.Contact, error)
}

type InvoicesService interface {
	GetScheduledInvoices(ctx context.Context, template entities.InvoiceTemplate) ([]entities.Invoice, error)
}

type CalendarService interface {
	Sync(ctx context.Context) error
}

type Logger interface {
	Warn(msg string, fields ...log.Field)
}

type UnitOfWork interface {
	CreateInvoices(invoices []entities.Invoice, template entities.InvoiceTemplate)
	DeleteInvoice(id uint64)
	DeleteTemplate(id uint64)
	UpdateInvoice(invoice entities.Invoice)
	Commit(ctx context.Context) error
}

type UseCase struct {
	periods  PeriodsRepository
	invoices InvoicesRepository
	contacts ContactsRepository
	service  InvoicesService
	calendar CalendarService
	logger   Logger
	uow      UnitOfWork
}

func NewUseCase(
	periods PeriodsRepository,
	invoices InvoicesRepository,
	contacts ContactsRepository,
	service InvoicesService,
	calendar CalendarService,
	logger Logger,
	uow UnitOfWork,
) *UseCase {
	return &UseCase{
		periods:  periods,
		invoices: invoices,
		contacts: contacts,
		service:  service,
		calendar: calendar,
		logger:   logger,
		uow:      uow,
	}
}

func (u *UseCase) Execute(ctx context.Context, req port.UpdateInvoiceRequest) error {
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

		u.uow.UpdateInvoice(*invoice)

		err = u.uow.Commit(ctx)
		if err != nil {
			return err
		}

		err = u.calendar.Sync(ctx)
		if err != nil {
			u.logger.Warn("cannot sync calendar", log.Err(err))
		}

		return nil
	}

	invoices, err := u.invoices.GetByTemplateID(ctx, invoice.Template.ID)
	if err != nil {
		return err
	}

	newTemplate := entities.InvoiceTemplate{
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
	}

	u.uow.DeleteTemplate(invoice.Template.ID)

	deletingInvoices := make([]entities.Invoice, 0, len(invoices))

	for _, inv := range invoices {
		if req.Date.After(inv.Date.Time) {
			deletingInvoices = append(deletingInvoices, inv)
		}
	}

	for _, inv := range deletingInvoices {
		u.uow.DeleteInvoice(inv.ID)
	}

	scheduledInvoices, err := u.service.GetScheduledInvoices(ctx, newTemplate)
	if err != nil {
		return err
	}

	u.uow.CreateInvoices(scheduledInvoices, newTemplate)

	err = u.uow.Commit(ctx)
	if err != nil {
		return err
	}

	err = u.calendar.Sync(ctx)
	if err != nil {
		u.logger.Warn("cannot sync calendar", log.Err(err))
	}

	return nil
}
