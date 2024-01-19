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

type InvoicesRepository interface {
	CreateTemplate(ctx context.Context, template entities.InvoiceTemplate) (*entities.InvoiceTemplate, error)
	BulkCreate(ctx context.Context, invoices []entities.Invoice) ([]entities.Invoice, error)
}

type PeriodsRepository interface {
	GetLast(ctx context.Context) (*entities.Period, error)
}

type UseCase struct {
	invoices InvoicesRepository
	service  InvoicesService
	periods  PeriodsRepository
}

func NewUseCase(invoices InvoicesRepository, service InvoicesService, periods PeriodsRepository) *UseCase {
	return &UseCase{
		invoices: invoices,
		service:  service,
		periods:  periods,
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
	} else if len(invoices) > 1 {
		entity, err := u.invoices.CreateTemplate(ctx, template)
		if err != nil {
			return err
		}

		template = *entity
	}

	for _, invoice := range invoices {
		invoice.Template = &template
	}

	if len(invoices) == 1 {
		invoices[0].Template = nil
	}

	_, err = u.invoices.BulkCreate(ctx, invoices)
	if err != nil {
		return err
	}

	return nil
}
