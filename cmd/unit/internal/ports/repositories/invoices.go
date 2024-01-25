package repositories

import "github.com/devdammit/shekel/pkg/types/datetime"

type GetAllInvoicesRequest struct {
	StartedAt  *datetime.Date
	Limit      *uint64
	Offset     *uint64
	TemplateID *uint64
}

type InvoicesOrderBy string

const (
	OrderByDateAsc  InvoicesOrderBy = "date_asc"
	OrderByDateDesc InvoicesOrderBy = "date_desc"
)

type FindByDatesRequest struct {
	StartedAt datetime.DateTime
	EndedAt   datetime.DateTime

	Limit  *uint64
	Offset *uint64

	OrderBy *InvoicesOrderBy
}
