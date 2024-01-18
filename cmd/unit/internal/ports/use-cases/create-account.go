package use_cases

import (
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	"github.com/devdammit/shekel/pkg/currency"
)

type CreateAccountParams struct {
	Name        string
	Description *string
	Type        entities.AccountType
	Balance     currency.Amount
}
