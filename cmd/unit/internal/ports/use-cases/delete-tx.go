package use_cases

import "errors"

var (
	ErrCannotDeleteTxAtClosedPeriod = errors.New("cannot delete transaction at closed period")
)
