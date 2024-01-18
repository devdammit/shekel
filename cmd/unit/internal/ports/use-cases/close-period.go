package use_cases

import "errors"

var (
	ErrorPeriodAlreadyClosed = errors.New("period already closed")
	ErrorInvoiceNotPaid      = errors.New("invoice not paid")
)
