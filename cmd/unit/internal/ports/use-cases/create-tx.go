package use_cases

import (
	"errors"
	"github.com/devdammit/shekel/pkg/currency"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

type CreateTxRequest struct {
	Amount    currency.Amount
	FromID    *uint64
	ToID      *uint64
	Date      datetime.Date
	InvoiceID *uint64
}

var (
	ErrorTxDateOutOfRange = errors.New("transaction date out of range")
	ErrorTxNoAccounts     = errors.New("transaction must have at least one account")
)
