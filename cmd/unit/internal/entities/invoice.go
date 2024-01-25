package entities

import (
	"github.com/devdammit/shekel/pkg/currency"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

type InvoiceType string

const (
	InvoiceTypeIncome  InvoiceType = "Income"
	InvoiceTypeExpense InvoiceType = "Expense"
)

type InvoiceStatus string

const (
	InvoiceStatusPending InvoiceStatus = "Pending"
	InvoiceStatusPaid    InvoiceStatus = "Paid"
)

type Invoice struct {
	ID   uint64  `json:"id"`
	Name string  `json:"name"`
	Desc *string `json:"description"`

	Status InvoiceStatus `json:"status"`
	Type   InvoiceType   `json:"type"`

	Template     *InvoiceTemplate `json:"template"`
	Contact      *Contact         `json:"contact"`
	Transactions []Transaction    `json:"transactions"`

	Amount currency.Amount `json:"amount"`

	Date      datetime.DateTime `json:"date"`
	CreatedAt datetime.DateTime `json:"created_at"`
	UpdateAt  datetime.DateTime `json:"updated_at"`
}
