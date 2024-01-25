package create_invoice

import (
	"context"
	"errors"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
)

type InvoicesService interface {
	GetScheduledInvoices(ctx context.Context, template entities.InvoiceTemplate) ([]entities.Invoice, error)
}

type PeriodsRepository interface {
	GetLast(ctx context.Context) (*entities.Period, error)
}

type CalendarService interface {
	Sync(ctx context.Context) error
}

type UnitOfWork interface {
	CreateInvoices(invoices []entities.Invoice, template entities.InvoiceTemplate)
	CreateInvoice(invoice entities.Invoice)
	Commit(ctx context.Context) error
}

type UseCase struct {
	service  InvoicesService
	periods  PeriodsRepository
	calendar CalendarService
	uow      UnitOfWork
}

func NewUseCase(service InvoicesService, periods PeriodsRepository, calendar CalendarService, work UnitOfWork) *UseCase {
	return &UseCase{
		service:  service,
		periods:  periods,
		calendar: calendar,
		uow:      work,
	}
}

func (u *UseCase) Execute(ctx context.Context, request port.CreateInvoiceRequest) error {
	period, err := u.periods.GetLast(ctx)
	if err != nil {
		return err
	}

	if period.ClosedAt != nil {
		return errors.New("cannot create invoice at closed period")
	}

	if request.Date.Time.Before(period.CreatedAt.Time) {
		return errors.New("cannot create invoice before current period")
	}

	template := entities.InvoiceTemplate{
		Name:      request.Name,
		Desc:      request.Description,
		Type:      request.Type,
		Amount:    request.Amount,
		ContactID: request.ContactID,

		Date: request.Date,
	}

	if request.Plan != nil {
		template.RepeatPlanner = &entities.RepeatPlanner{
			Interval:      request.Plan.Interval,
			IntervalCount: request.Plan.IntervalCount,
			DaysOfWeek:    request.Plan.DaysOfWeek,
			EndDate:       request.Plan.EndDate,
			EndCount:      request.Plan.EndCount,
		}
	}

	invoices, err := u.service.GetScheduledInvoices(ctx, template)
	if err != nil {
		return err
	}

	if len(invoices) == 0 {
		return errors.New("no invoices to create")
	} else if len(invoices) == 1 {
		u.uow.CreateInvoice(invoices[0])
	} else {
		u.uow.CreateInvoices(invoices, template)
	}

	err = u.uow.Commit(ctx)
	if err != nil {
		return err
	}

	err = u.calendar.Sync(ctx)
	if err != nil {
		return err
	}

	return nil
}
