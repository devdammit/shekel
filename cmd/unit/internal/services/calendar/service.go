package calendar

import (
	"context"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Start() {
}

func (s *Service) Sync(ctx context.Context) error {
	// 1. get all invoices from current period
	// 2. remove all events from calendar beginning from current period
	// 3. add all events from invoices

	// event payload
	// 1. title: invoice title
	// 2. description: invoice description
	// 3. status: invoice status
	// 4. type: invoice type
	// 5. contact: invoice contact
	// 6. amount: invoice amount

	return nil
}
