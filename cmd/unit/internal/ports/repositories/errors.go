package repositories

import "errors"

var (
	ErrNotFound        = errors.New("entity not found")
	ErrAlreadyExists   = errors.New("entity already exists")
	ErrHasOpenedPeriod = errors.New("has opened period")
)
