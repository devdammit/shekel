package entities

import (
	"github.com/devdammit/shekel/pkg/currency"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

type Transaction struct {
	ID     uint64          `json:"id"`
	Amount currency.Amount `json:"amount"`

	From *Account `json:"from"`
	To   *Account `json:"to"`

	Invoice   *Invoice          `json:"invoice"`
	CreatedAt datetime.DateTime `json:"created_at"`
}
