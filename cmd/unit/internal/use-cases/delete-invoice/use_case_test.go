package delete_invoice_test

import (
	"context"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	mocks "github.com/devdammit/shekel/cmd/unit/internal/mocks/use-cases/delete-invoice"
	delete_invoice "github.com/devdammit/shekel/cmd/unit/internal/use-cases/delete-invoice"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestUseCase_Execute(t *testing.T) {
	t.Run("should return error when we want delete single invoice and it is paid", func(t *testing.T) {
		var (
			ctx       = context.Background()
			ctrl      = gomock.NewController(t)
			invoices  = mocks.NewMockInvoicesRepository(ctrl)
			templates = mocks.NewMockInvoicesTemplatesRepository(ctrl)
		)

		useCase := delete_invoice.NewUseCase(invoices, templates)

		invoices.EXPECT().GetByID(ctx, uint64(1)).Return(&entities.Invoice{
			ID:     1,
			Status: entities.InvoiceStatusPaid,
		}, nil)

		err := useCase.Execute(ctx, uint64(1), true)

		assert.EqualError(t, err, "cannot delete paid invoice")
	})

	t.Run("should delete single invoice", func(t *testing.T) {
		var (
			ctx       = context.Background()
			ctrl      = gomock.NewController(t)
			invoices  = mocks.NewMockInvoicesRepository(ctrl)
			templates = mocks.NewMockInvoicesTemplatesRepository(ctrl)
		)

		useCase := delete_invoice.NewUseCase(invoices, templates)

		invoices.EXPECT().GetByID(ctx, uint64(1)).Return(&entities.Invoice{
			ID:     1,
			Status: entities.InvoiceStatusPending,
		}, nil)

		invoices.EXPECT().Delete(ctx, uint64(1)).Return(nil)

		err := useCase.Execute(ctx, uint64(1), true)

		assert.NoError(t, err)
	})

	t.Run("should delete all not paid invoices", func(t *testing.T) {
		var (
			ctx       = context.Background()
			ctrl      = gomock.NewController(t)
			invoices  = mocks.NewMockInvoicesRepository(ctrl)
			templates = mocks.NewMockInvoicesTemplatesRepository(ctrl)
		)

		useCase := delete_invoice.NewUseCase(invoices, templates)

		invoices.EXPECT().GetByID(ctx, uint64(1)).Return(&entities.Invoice{
			ID:     1,
			Status: entities.InvoiceStatusPending,
			Template: &entities.InvoiceTemplate{
				ID: 1,
			},
			Date: datetime.MustParseDateTime("2024-01-01 00:01"),
		}, nil)

		invoices.EXPECT().GetByTemplateID(ctx, uint64(1)).Return([]entities.Invoice{
			{
				ID:     1,
				Status: entities.InvoiceStatusPending,
				Template: &entities.InvoiceTemplate{
					ID: 1,
				},
				Date: datetime.MustParseDateTime("2024-01-01 00:01"),
			},
			{
				ID:     2,
				Status: entities.InvoiceStatusPaid,
				Template: &entities.InvoiceTemplate{
					ID: 1,
				},
				Date: datetime.MustParseDateTime("2024-01-02 00:01"),
			},
			{
				ID:     3,
				Status: entities.InvoiceStatusPending,
				Template: &entities.InvoiceTemplate{
					ID: 1,
				},
				Date: datetime.MustParseDateTime("2024-01-02 00:01"),
			},
		}, nil)

		invoices.EXPECT().Delete(ctx, uint64(1)).Return(nil)
		invoices.EXPECT().Delete(ctx, uint64(3)).Return(nil)
		templates.EXPECT().Delete(ctx, uint64(1)).Return(nil)

		err := useCase.Execute(ctx, uint64(1), false)

		assert.NoError(t, err)
	})
}
