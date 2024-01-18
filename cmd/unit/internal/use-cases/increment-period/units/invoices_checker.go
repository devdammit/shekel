package units

import (
	"context"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
)

type InvoicesRepository interface {
	GetAllByPeriodID(periodID uint64) ([]entities.Invoice, error)
}

type InvoicesChecker struct {
	repo InvoicesRepository
}

func NewInvoicesChecker(repo InvoicesRepository) *InvoicesChecker {
	return &InvoicesChecker{repo: repo}
}

func (u *InvoicesChecker) GetName() string {
	return "invoices_extender"
}

func (u *InvoicesChecker) Handle(ctx context.Context, request *Request, payload *Payload) (*Payload, error) {
	if payload == nil || payload.ActivePeriod == nil {
		return nil, nil
	}

	invoices, err := u.repo.GetAllByPeriodID(payload.ActivePeriod.ID)
	if err != nil {
		return nil, err
	}

	for _, invoice := range invoices {
		if invoice.Status == entities.InvoiceStatusPending {
			return nil, port.ErrorInvoiceNotPaid
		}
	}

	return payload, nil
}
