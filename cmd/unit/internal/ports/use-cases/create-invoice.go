package use_cases

import (
	"time"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	"github.com/devdammit/shekel/pkg/currency"
	"github.com/devdammit/shekel/pkg/planner"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

type InvoicePlan struct {
	Interval      planner.PlanRepeatInterval
	IntervalCount uint32
	DaysOfWeek    []time.Weekday
	EndDate       *datetime.Date
	EndCount      *uint32
}

type CreateInvoiceRequest struct {
	Name        string
	Description *string

	Plan *InvoicePlan

	Type entities.InvoiceType

	Amount    currency.Amount
	ContactID uint64

	Date datetime.DateTime
}
