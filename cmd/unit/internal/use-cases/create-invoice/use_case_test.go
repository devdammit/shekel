package create_invoice_test

import (
	"context"
	"testing"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	mocks "github.com/devdammit/shekel/cmd/unit/internal/mocks/use-cases/create-invoice"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
	create_invoice "github.com/devdammit/shekel/cmd/unit/internal/use-cases/create-invoice"
	"github.com/devdammit/shekel/pkg/currency"
	"github.com/devdammit/shekel/pkg/planner"
	"github.com/devdammit/shekel/pkg/pointer"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Execute(t *testing.T) {
	t.Run("should create one invoice without template", func(t *testing.T) {
		var (
			mockController = gomock.NewController(t)
			service        = mocks.NewMockInvoicesService(mockController)
			periods        = mocks.NewMockPeriodsRepository(mockController)
			calendar       = mocks.NewMockCalendarService(mockController)
			uow            = mocks.NewMockUnitOfWork(mockController)
		)

		uow.EXPECT().CreateInvoice(gomock.Any()).Times(1).Do(func(invoice entities.Invoice) {
			assert.Equal(t, "Invoice 1", invoice.Name)
		})

		periods.EXPECT().GetLast(gomock.Any()).Times(1).Return(&entities.Period{
			CreatedAt: datetime.MustParseDateTime("2024-01-01 00:01"),
		}, nil)

		service.EXPECT().GetScheduledInvoices(gomock.Any(), gomock.Any()).Times(1).Return([]entities.Invoice{
			{
				Name: "Invoice 1",
			},
		}, nil)

		uow.EXPECT().Commit(gomock.Any()).Times(1).Return(nil)

		calendar.EXPECT().Sync(gomock.Any()).Times(1).Return(nil)

		useCase := create_invoice.NewUseCase(service, periods, calendar, uow)

		err := useCase.Execute(context.Background(), port.CreateInvoiceRequest{
			Name:        "Invoice 1",
			Description: pointer.Ptr("Description"),
			Type:        entities.InvoiceTypeExpense,
			Amount: currency.Amount{
				CurrencyCode: currency.USD,
				Value:        100,
			},

			Date: datetime.MustParseDateTime("2024-01-01 04:20"),
		})

		assert.NoError(t, err)
	})

	t.Run("should create two invoices with template", func(t *testing.T) {
		var (
			mockController = gomock.NewController(t)
			service        = mocks.NewMockInvoicesService(mockController)
			periods        = mocks.NewMockPeriodsRepository(mockController)
			calendar       = mocks.NewMockCalendarService(mockController)
			uow            = mocks.NewMockUnitOfWork(mockController)
		)

		periods.EXPECT().GetLast(gomock.Any()).Times(1).Return(&entities.Period{
			CreatedAt: datetime.MustParseDateTime("2024-01-01 00:01"),
		}, nil)

		uow.EXPECT().CreateInvoices(gomock.Any(), gomock.Any()).Times(1).Do(func(invoices []entities.Invoice, template entities.InvoiceTemplate) {
			assert.Equal(t, 2, len(invoices))
		})

		uow.EXPECT().Commit(gomock.Any()).Times(1).Return(nil)

		service.EXPECT().GetScheduledInvoices(gomock.Any(), gomock.Any()).Times(1).Return([]entities.Invoice{
			{
				Name: "Invoice 1",
				Template: &entities.InvoiceTemplate{
					ID: 1,
				},
			},
			{
				Name: "Invoice 2",
				Template: &entities.InvoiceTemplate{
					ID: 1,
				},
			},
		}, nil)
		calendar.EXPECT().Sync(gomock.Any()).Times(1).Return(nil)

		useCase := create_invoice.NewUseCase(service, periods, calendar, uow)

		err := useCase.Execute(context.Background(), port.CreateInvoiceRequest{
			Name:        "Invoice 1",
			Description: pointer.Ptr("Description"),
			Type:        entities.InvoiceTypeExpense,
			Amount: currency.Amount{
				CurrencyCode: currency.USD,
				Value:        100,
			},

			Date: datetime.MustParseDateTime("2024-01-01 04:20"),
			Plan: &entities.RepeatPlanner{
				Interval:      planner.PlanRepeatIntervalMonthly,
				IntervalCount: 1,
				EndCount:      pointer.Ptr(uint32(2)),
			},
		})

		assert.NoError(t, err)
	})

	t.Run("should return error if no invoices to create", func(t *testing.T) {
		var (
			mockController = gomock.NewController(t)
			service        = mocks.NewMockInvoicesService(mockController)
			periods        = mocks.NewMockPeriodsRepository(mockController)
			calendar       = mocks.NewMockCalendarService(mockController)
			uow            = mocks.NewMockUnitOfWork(mockController)
		)

		periods.EXPECT().GetLast(gomock.Any()).Times(1).Return(&entities.Period{
			CreatedAt: datetime.MustParseDateTime("2024-01-01 00:01"),
		}, nil)

		service.EXPECT().GetScheduledInvoices(gomock.Any(), gomock.Any()).Times(1).Return(nil, nil)

		useCase := create_invoice.NewUseCase(service, periods, calendar, uow)

		err := useCase.Execute(context.Background(), port.CreateInvoiceRequest{
			Name:        "Invoice 1",
			Description: pointer.Ptr("Description"),
			Type:        entities.InvoiceTypeExpense,
			Amount: currency.Amount{
				CurrencyCode: currency.USD,
				Value:        100,
			},

			Date: datetime.MustParseDateTime("2024-01-01 04:20"),
		})

		assert.Error(t, err)
	})

	t.Run("should error if user try create invoice before current period", func(t *testing.T) {
		var (
			mockController = gomock.NewController(t)
			service        = mocks.NewMockInvoicesService(mockController)
			periods        = mocks.NewMockPeriodsRepository(mockController)
			calendar       = mocks.NewMockCalendarService(mockController)
			uow            = mocks.NewMockUnitOfWork(mockController)
		)

		periods.EXPECT().GetLast(gomock.Any()).Times(1).Return(&entities.Period{
			CreatedAt: datetime.MustParseDateTime("2024-01-01 00:01"),
		}, nil)

		service.EXPECT().GetScheduledInvoices(gomock.Any(), gomock.Any()).Times(0)

		useCase := create_invoice.NewUseCase(service, periods, calendar, uow)

		err := useCase.Execute(context.Background(), port.CreateInvoiceRequest{
			Name:        "Invoice 1",
			Description: pointer.Ptr("Description"),
			Type:        entities.InvoiceTypeExpense,
			Amount: currency.Amount{
				CurrencyCode: currency.USD,
				Value:        100,
			},

			Date: datetime.MustParseDateTime("2023-11-01 04:20"),
		})

		assert.EqualError(t, err, "cannot create invoice before current period")
	})

	t.Run("should error if user try create invoice at closed period", func(t *testing.T) {
		var (
			mockController = gomock.NewController(t)
			service        = mocks.NewMockInvoicesService(mockController)
			periods        = mocks.NewMockPeriodsRepository(mockController)
			calendar       = mocks.NewMockCalendarService(mockController)
			uow            = mocks.NewMockUnitOfWork(mockController)
		)

		periods.EXPECT().GetLast(gomock.Any()).Times(1).Return(&entities.Period{
			CreatedAt: datetime.MustParseDateTime("2024-01-01 00:01"),
			ClosedAt:  pointer.Ptr(datetime.MustParseDateTime("2024-01-02 00:01")),
		}, nil)

		service.EXPECT().GetScheduledInvoices(gomock.Any(), gomock.Any()).Times(0)

		useCase := create_invoice.NewUseCase(service, periods, calendar, uow)

		err := useCase.Execute(context.Background(), port.CreateInvoiceRequest{
			Name:        "Invoice 1",
			Description: pointer.Ptr("Description"),
			Type:        entities.InvoiceTypeExpense,
			Amount: currency.Amount{
				CurrencyCode: currency.USD,
				Value:        100,
			},

			Date: datetime.MustParseDateTime("2024-01-01 04:20"),
		})

		assert.EqualError(t, err, "cannot create invoice at closed period")
	})
}
