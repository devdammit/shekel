package services

import "errors"

var (
	ErrUnknownCurrency = errors.New("unknown currency")
	ErrRatesNotFound   = errors.New("rates not found")
)
