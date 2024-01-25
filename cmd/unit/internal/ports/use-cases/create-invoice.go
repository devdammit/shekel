package use_cases

import (
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	"github.com/devdammit/shekel/pkg/currency"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

type CreateInvoiceRequest struct {
	Name        string
	Description *string

	Plan *entities.RepeatPlanner

	Type entities.InvoiceType

	Amount    currency.Amount
	ContactID uint64

	Date datetime.DateTime
}
