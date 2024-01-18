package use_cases

import (
	"errors"
	"github.com/devdammit/shekel/pkg/currency"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

type UpdateTxRequest struct {
	ID        uint64
	Date      datetime.Date
	Amount    currency.Amount
	FromID    *uint64
	ToID      *uint64
	InvoiceID *uint64
}

var (
	CannotUpdateTxAtClosedPeriod = errors.New("cannot update transaction at closed period")
)
