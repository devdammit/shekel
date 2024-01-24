package entities

import (
	"time"

	"github.com/devdammit/shekel/pkg/currency"
	"github.com/devdammit/shekel/pkg/planner"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

type RepeatPlanner struct {
	IntervalCount uint32                     `json:"interval_count"`
	Interval      planner.PlanRepeatInterval `json:"interval"`
	DaysOfWeek    []time.Weekday             `json:"day_of_week"`
	EndDate       *datetime.Date             `json:"end_date"`
	EndCount      *uint32                    `json:"end_count"`
}

type InvoiceTemplate struct {
	ID   uint64  `json:"id"`
	Name string  `json:"name"`
	Desc *string `json:"description"`

	Type InvoiceType `json:"type"`

	Amount currency.Amount `json:"amount"`

	RepeatPlanner *RepeatPlanner `json:"repeat_planner"`

	ContactID uint64 `json:"contact_id"`

	Date      datetime.DateTime  `json:"date"`
	DeletedAt *datetime.DateTime `json:"delete_at"`
	CreatedAt datetime.DateTime  `json:"created_at"`
	UpdateAt  datetime.DateTime  `json:"updated_at"`
}
