package invoices

import (
	"context"
	"errors"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	"github.com/devdammit/shekel/pkg/planner"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"time"
)

type AppConfig interface {
	GetFinancialYearStart() time.Time
}

type Services struct {
}

type Service struct {
	Services  Services
	AppConfig AppConfig
}

func (s *Service) GetScheduledInvoices(_ context.Context, startAt datetime.Date, template entities.InvoiceTemplate) ([]entities.Invoice, error) {
	if template.Date.Before(startAt.Time) && template.RepeatPlanner == nil {
		return nil, errors.New("no invoices to create")
	}

	if template.RepeatPlanner == nil {
		return []entities.Invoice{
			{
				Name: template.Name,
				Desc: template.Desc,

				Status: entities.InvoiceStatusPending,
				Type:   template.Type,
				Contact: &entities.Contact{
					ID: template.ContactID,
				},
				Template: &template,
				Amount:   template.Amount,
				Date:     template.Date,
			},
		}, nil
	}

	if template.RepeatPlanner.EndDate == nil {
		template.RepeatPlanner.EndDate = &datetime.Date{
			Time: s.AppConfig.GetFinancialYearStart().AddDate(1, 1, 0),
		}
	}

	p := planner.NewPlanner(
		template.RepeatPlanner.IntervalCount,
		template.RepeatPlanner.Interval,
		template.RepeatPlanner.EndDate.Time,
		template.Date.Time,
		template.RepeatPlanner.EndCount,
		template.RepeatPlanner.DaysOfWeek,
	)

	dates, err := p.GetDates()
	if err != nil {
		return nil, err
	}

	invoices := make([]entities.Invoice, 0, len(dates))

	for _, date := range dates {
		invoices = append(invoices, entities.Invoice{
			Name: template.Name,
			Desc: template.Desc,

			Status: entities.InvoiceStatusPending,
			Type:   template.Type,
			Contact: &entities.Contact{
				ID: template.ContactID,
			},
			Template: &template,
			Amount:   template.Amount,
			Date:     datetime.NewDateTime(date),
		})
	}

	return invoices, nil
}
