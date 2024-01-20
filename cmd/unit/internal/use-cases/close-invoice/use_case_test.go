package close_invoice_test

import (
	"context"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	mocks "github.com/devdammit/shekel/cmd/unit/internal/mocks/use-cases/close-invoice"
	close_invoice "github.com/devdammit/shekel/cmd/unit/internal/use-cases/close-invoice"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestUseCase_Execute(t *testing.T) {
	t.Run("should return error if invoice already paid", func(t *testing.T) {
		var (
			mockController = gomock.NewController(t)
			invoices       = mocks.NewMockinvoicesRepository(mockController)
			periods        = mocks.NewMockperiodsRepository(mockController)
		)

		invoices.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(1).Return(&entities.Invoice{
			Status: entities.InvoiceStatusPaid,
		}, nil)

		useCase := close_invoice.NewUseCase(invoices, periods)

		err := useCase.Execute(context.Background(), 1)

		assert.EqualError(t, err, "invoice already paid")
	})

	t.Run("should return error if invoice date is before current period", func(t *testing.T) {
		var (
			mockController = gomock.NewController(t)
			invoices       = mocks.NewMockinvoicesRepository(mockController)
			periods        = mocks.NewMockperiodsRepository(mockController)
		)

		invoices.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(1).Return(&entities.Invoice{
			Status: entities.InvoiceStatusPending,
			Date:   datetime.MustParseDateTime("2024-01-01 00:01"),
		}, nil)

		periods.EXPECT().GetLast(gomock.Any()).Times(1).Return(&entities.Period{
			CreatedAt: datetime.MustParseDateTime("2024-01-01 00:02"),
		}, nil)

		useCase := close_invoice.NewUseCase(invoices, periods)

		err := useCase.Execute(context.Background(), 1)

		assert.EqualError(t, err, "cannot close invoice before current period")
	})

	t.Run("should mark invoice as paid", func(t *testing.T) {
		var (
			mockController = gomock.NewController(t)
			invoices       = mocks.NewMockinvoicesRepository(mockController)
			periods        = mocks.NewMockperiodsRepository(mockController)
		)

		invoices.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(1).Return(&entities.Invoice{
			Status: entities.InvoiceStatusPending,
			Date:   datetime.MustParseDateTime("2024-01-01 00:03"),
		}, nil)

		periods.EXPECT().GetLast(gomock.Any()).Times(1).Return(&entities.Period{
			CreatedAt: datetime.MustParseDateTime("2024-01-01 00:02"),
		}, nil)

		invoices.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).Return(nil, nil).Do(func(ctx context.Context, invoice *entities.Invoice) {
			assert.Equal(t, entities.InvoiceStatusPaid, invoice.Status)
		})

		useCase := close_invoice.NewUseCase(invoices, periods)

		err := useCase.Execute(context.Background(), 1)

		assert.NoError(t, err)
	})
}
