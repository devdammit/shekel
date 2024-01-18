package entities

import (
	"errors"
	"github.com/devdammit/shekel/pkg/currency"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

type AccountType string

const (
	AccountTypeCash   AccountType = "cash"
	AccountTypeCredit AccountType = "credit"
	AccountTypeDebit  AccountType = "debit"
)

var (
	ErrorAccountExists = errors.New("account already exists")
)

type Account struct {
	ID          uint64          `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	Type        AccountType     `json:"type"`
	Balance     currency.Amount `json:"amount"`

	IsArchived bool `json:"deleted"`

	CreatedAt datetime.Time `json:"created_at"`
	UpdateAt  datetime.Time `json:"updated_at"`
}
