// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	"github.com/devdammit/shekel/pkg/gql"
)

type AddContactInput struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type Amount struct {
	Amount   float64      `json:"amount"`
	Currency gql.Currency `json:"currency"`
}

type AmountInput struct {
	Amount   float64      `json:"amount"`
	Currency gql.Currency `json:"currency"`
}

type App struct {
	Initialized  bool             `json:"initialized"`
	ActivePeriod *entities.Period `json:"activePeriod,omitempty"`
	Version      string           `json:"version"`
}

type CreateAccountInput struct {
	Name        string       `json:"name"`
	Description *string      `json:"description,omitempty"`
	Type        AccountType  `json:"type"`
	Balance     *AmountInput `json:"balance"`
}

type CreateInvoiceInput struct {
	Name        string              `json:"name"`
	Description *string             `json:"description,omitempty"`
	Plan        *RepeatPlannerInput `json:"plan,omitempty"`
	Type        InvoiceType         `json:"type"`
	Amount      *AmountInput        `json:"amount"`
	ContactID   uint64              `json:"contactId"`
	Date        gql.DateTime        `json:"date"`
}

type FindInvoiceByPeriod struct {
	PeriodID    uint64           `json:"periodId"`
	OnlyPending *bool            `json:"onlyPending,omitempty"`
	OnlyPaid    *bool            `json:"onlyPaid,omitempty"`
	Limit       *uint64          `json:"limit,omitempty"`
	Offset      *uint64          `json:"offset,omitempty"`
	OrderBy     *InvoicesOrderBy `json:"orderBy,omitempty"`
}

type Mutation struct {
}

type QRCodeInput struct {
	File graphql.Upload `json:"file"`
	Bank string         `json:"bank"`
}

type Query struct {
}

type RepeatPlannerInput struct {
	IntervalCount uint32             `json:"intervalCount"`
	Interval      PlanRepeatInterval `json:"interval"`
	DaysOfWeek    []uint32           `json:"daysOfWeek,omitempty"`
	EndDate       *gql.Date          `json:"endDate,omitempty"`
	EndCount      *uint32            `json:"endCount,omitempty"`
}

type UpdateAccountInput struct {
	ID          uint64       `json:"id"`
	Name        string       `json:"name"`
	Description *string      `json:"description,omitempty"`
	Balance     *AmountInput `json:"balance"`
}

type UpdateInvoiceInput struct {
	ID          uint64              `json:"id"`
	Name        string              `json:"name"`
	Description *string             `json:"description,omitempty"`
	Plan        *RepeatPlannerInput `json:"plan,omitempty"`
	Type        InvoiceType         `json:"type"`
	Amount      *AmountInput        `json:"amount"`
	ContactID   uint64              `json:"contactId"`
	Date        gql.DateTime        `json:"date"`
}

type AccountType string

const (
	AccountTypeCash   AccountType = "CASH"
	AccountTypeCredit AccountType = "CREDIT"
	AccountTypeDebit  AccountType = "DEBIT"
)

var AllAccountType = []AccountType{
	AccountTypeCash,
	AccountTypeCredit,
	AccountTypeDebit,
}

func (e AccountType) IsValid() bool {
	switch e {
	case AccountTypeCash, AccountTypeCredit, AccountTypeDebit:
		return true
	}
	return false
}

func (e AccountType) String() string {
	return string(e)
}

func (e *AccountType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AccountType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AccountType", str)
	}
	return nil
}

func (e AccountType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type InvoiceStatus string

const (
	InvoiceStatusPending InvoiceStatus = "Pending"
	InvoiceStatusPaid    InvoiceStatus = "Paid"
)

var AllInvoiceStatus = []InvoiceStatus{
	InvoiceStatusPending,
	InvoiceStatusPaid,
}

func (e InvoiceStatus) IsValid() bool {
	switch e {
	case InvoiceStatusPending, InvoiceStatusPaid:
		return true
	}
	return false
}

func (e InvoiceStatus) String() string {
	return string(e)
}

func (e *InvoiceStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = InvoiceStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid InvoiceStatus", str)
	}
	return nil
}

func (e InvoiceStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type InvoiceType string

const (
	InvoiceTypeIncome  InvoiceType = "Income"
	InvoiceTypeExpense InvoiceType = "Expense"
)

var AllInvoiceType = []InvoiceType{
	InvoiceTypeIncome,
	InvoiceTypeExpense,
}

func (e InvoiceType) IsValid() bool {
	switch e {
	case InvoiceTypeIncome, InvoiceTypeExpense:
		return true
	}
	return false
}

func (e InvoiceType) String() string {
	return string(e)
}

func (e *InvoiceType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = InvoiceType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid InvoiceType", str)
	}
	return nil
}

func (e InvoiceType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type InvoicesOrderBy string

const (
	InvoicesOrderByDateAsc  InvoicesOrderBy = "Date_ASC"
	InvoicesOrderByDateDesc InvoicesOrderBy = "Date_DESC"
)

var AllInvoicesOrderBy = []InvoicesOrderBy{
	InvoicesOrderByDateAsc,
	InvoicesOrderByDateDesc,
}

func (e InvoicesOrderBy) IsValid() bool {
	switch e {
	case InvoicesOrderByDateAsc, InvoicesOrderByDateDesc:
		return true
	}
	return false
}

func (e InvoicesOrderBy) String() string {
	return string(e)
}

func (e *InvoicesOrderBy) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = InvoicesOrderBy(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid InvoicesOrderBy", str)
	}
	return nil
}

func (e InvoicesOrderBy) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type OrderBy string

const (
	OrderByAsc  OrderBy = "ASC"
	OrderByDesc OrderBy = "DESC"
)

var AllOrderBy = []OrderBy{
	OrderByAsc,
	OrderByDesc,
}

func (e OrderBy) IsValid() bool {
	switch e {
	case OrderByAsc, OrderByDesc:
		return true
	}
	return false
}

func (e OrderBy) String() string {
	return string(e)
}

func (e *OrderBy) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = OrderBy(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid OrderBy", str)
	}
	return nil
}

func (e OrderBy) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type PlanRepeatInterval string

const (
	PlanRepeatIntervalDaily   PlanRepeatInterval = "Daily"
	PlanRepeatIntervalWeekly  PlanRepeatInterval = "Weekly"
	PlanRepeatIntervalMonthly PlanRepeatInterval = "Monthly"
	PlanRepeatIntervalYearly  PlanRepeatInterval = "Yearly"
)

var AllPlanRepeatInterval = []PlanRepeatInterval{
	PlanRepeatIntervalDaily,
	PlanRepeatIntervalWeekly,
	PlanRepeatIntervalMonthly,
	PlanRepeatIntervalYearly,
}

func (e PlanRepeatInterval) IsValid() bool {
	switch e {
	case PlanRepeatIntervalDaily, PlanRepeatIntervalWeekly, PlanRepeatIntervalMonthly, PlanRepeatIntervalYearly:
		return true
	}
	return false
}

func (e PlanRepeatInterval) String() string {
	return string(e)
}

func (e *PlanRepeatInterval) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PlanRepeatInterval(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PlanRepeatInterval", str)
	}
	return nil
}

func (e PlanRepeatInterval) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
