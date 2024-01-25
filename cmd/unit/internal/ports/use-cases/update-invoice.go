package use_cases

import (
	"errors"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	"github.com/devdammit/shekel/pkg/currency"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

var (
	ErrorCannotUpdateInvoiceAtClosedPeriod = errors.New("cannot update invoice at closed period")
)

type UpdateInvoiceRequest struct {
	InvoiceID   uint64
	Name        string
	Description *string

	Plan *entities.RepeatPlanner

	Type      entities.InvoiceType
	Amount    currency.Amount
	ContactID uint64
	Date      datetime.DateTime
}
