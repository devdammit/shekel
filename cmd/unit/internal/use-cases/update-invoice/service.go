package update_invoice

import (
	"context"
	"errors"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
)

type PeriodsRepository interface {
	GetLast(ctx context.Context) (*entities.Period, error)
}

type InvoicesRepository interface {
	GetByID(ctx context.Context, id uint64) (*entities.Invoice, error)
}

type InvoicesTemplateRepository interface {
	GetByID(ctx context.Context, id uint64) (*entities.InvoiceTemplate, error)
	Update(ctx context.Context, template *entities.InvoiceTemplate) (*entities.InvoiceTemplate, error)
}

type ContactsRepository interface {
	GetByID(ctx context.Context, id uint64) (*entities.Contact, error)
}

type UpdateInvoiceUseCase struct {
	periods   PeriodsRepository
	invoices  InvoicesRepository
	templates InvoicesTemplateRepository
	contacts  ContactsRepository
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
	} else {
		template, err := u.templates.GetByID(ctx, invoice.Template.ID)
		if err != nil {
			return err
		}

		template.Name = req.Name
		template.Desc = req.Description
		template.Type = req.Type
		template.Amount = req.Amount
		template.ContactID = contact.ID
		template.Date = req.Date
		template.RepeatPlanner = &entities.RepeatPlanner{
			Interval:      req.Plan.Interval,
			IntervalCount: req.Plan.IntervalCount,
			DaysOfWeek:    req.Plan.DaysOfWeek,
			EndDate:       req.Plan.EndDate,
			EndCount:      req.Plan.EndCount,
		}

		invoice.Template = template

	}

	return nil
}
