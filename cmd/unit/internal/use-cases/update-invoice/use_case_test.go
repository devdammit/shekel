package update_invoice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	mocks "github.com/devdammit/shekel/cmd/unit/internal/mocks/use-cases/update-invoice"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
	update_invoice "github.com/devdammit/shekel/cmd/unit/internal/use-cases/update-invoice"
	"github.com/devdammit/shekel/pkg/log"
	"github.com/devdammit/shekel/pkg/planner"
	"github.com/devdammit/shekel/pkg/pointer"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpdateInvoiceUseCase_Execute(t *testing.T) {
	t.Run("should return error when period is closed", func(t *testing.T) {
		var (
			mockController = gomock.NewController(t)
			invoices       = mocks.NewMockInvoicesRepository(mockController)
			contacts       = mocks.NewMockContactsRepository(mockController)
			periods        = mocks.NewMockPeriodsRepository(mockController)
			service        = mocks.NewMockInvoicesService(mockController)
			calendar       = mocks.NewMockCalendarService(mockController)
			logger         = mocks.NewMockLogger(mockController)
			uow            = mocks.NewMockUnitOfWork(mockController)
		)

		periods.EXPECT().GetLast(gomock.Any()).Times(1).Return(&entities.Period{
			ClosedAt: pointer.Ptr(datetime.MustParseDateTime("2024-01-01 00:01")),
		}, nil)

		useCase := update_invoice.NewUseCase(
			periods,
			invoices,
			contacts,
			service,
			calendar,
			logger,
			uow,
		)

		err := useCase.Execute(context.Background(), port.UpdateInvoiceRequest{})

		assert.EqualError(t, err, "cannot update invoice at closed period")
	})

	t.Run("should return error when updating invoice at previous period", func(t *testing.T) {
		var (
			mockController = gomock.NewController(t)
			invoices       = mocks.NewMockInvoicesRepository(mockController)
			contacts       = mocks.NewMockContactsRepository(mockController)
			periods        = mocks.NewMockPeriodsRepository(mockController)
			service        = mocks.NewMockInvoicesService(mockController)
			calendar       = mocks.NewMockCalendarService(mockController)
			logger         = mocks.NewMockLogger(mockController)
			uow            = mocks.NewMockUnitOfWork(mockController)
		)

		periods.EXPECT().GetLast(gomock.Any()).Times(1).Return(&entities.Period{
			CreatedAt: datetime.MustParseDateTime("2024-01-01 00:01"),
		}, nil)

		useCase := update_invoice.NewUseCase(
			periods,
			invoices,
			contacts,
			service,
			calendar,
			logger,
			uow,
		)

		err := useCase.Execute(context.Background(), port.UpdateInvoiceRequest{
			Date: datetime.MustParseDateTime("2023-01-01 00:01"),
		})

		assert.EqualError(t, err, "cannot update invoice at previous period")
	})

	t.Run("should return error when updating invoice from not current period", func(t *testing.T) {
		var (
			mockController = gomock.NewController(t)
			invoices       = mocks.NewMockInvoicesRepository(mockController)
			contacts       = mocks.NewMockContactsRepository(mockController)
			periods        = mocks.NewMockPeriodsRepository(mockController)
			service        = mocks.NewMockInvoicesService(mockController)
			calendar       = mocks.NewMockCalendarService(mockController)
			logger         = mocks.NewMockLogger(mockController)
			uow            = mocks.NewMockUnitOfWork(mockController)
		)

		periods.EXPECT().GetLast(gomock.Any()).Times(1).Return(&entities.Period{
			CreatedAt: datetime.MustParseDateTime("2024-01-01 00:01"),
		}, nil)
		invoices.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(1).Return(&entities.Invoice{
			Date: datetime.MustParseDateTime("2023-01-01 00:01"),
		}, nil)

		useCase := update_invoice.NewUseCase(
			periods,
			invoices,
			contacts,
			service,
			calendar,
			logger,
			uow,
		)

		err := useCase.Execute(context.Background(), port.UpdateInvoiceRequest{
			Date: datetime.MustParseDateTime("2024-01-01 00:01"),
		})

		assert.EqualError(t, err, "cannot update invoice at previous period")
	})

	t.Run("should return error when updating paid invoice", func(t *testing.T) {
		var (
			mockController = gomock.NewController(t)
			invoices       = mocks.NewMockInvoicesRepository(mockController)
			contacts       = mocks.NewMockContactsRepository(mockController)
			periods        = mocks.NewMockPeriodsRepository(mockController)
			service        = mocks.NewMockInvoicesService(mockController)
			calendar       = mocks.NewMockCalendarService(mockController)
			logger         = mocks.NewMockLogger(mockController)
			uow            = mocks.NewMockUnitOfWork(mockController)
		)

		periods.EXPECT().GetLast(gomock.Any()).Times(1).Return(&entities.Period{
			CreatedAt: datetime.MustParseDateTime("2024-01-01 00:00"),
		}, nil)
		invoices.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(1).Return(&entities.Invoice{
			Date:   datetime.MustParseDateTime("2024-01-01 00:01"),
			Status: entities.InvoiceStatusPaid,
		}, nil)

		useCase := update_invoice.NewUseCase(periods, invoices, contacts, service, calendar, logger, uow)

		err := useCase.Execute(context.Background(), port.UpdateInvoiceRequest{
			Date: datetime.MustParseDateTime("2024-01-01 00:01"),
		})

		assert.EqualError(t, err, "cannot update paid invoice")
	})

	t.Run("should update only one invoice", func(t *testing.T) {
		var (
			mockController = gomock.NewController(t)
			invoices       = mocks.NewMockInvoicesRepository(mockController)
			contacts       = mocks.NewMockContactsRepository(mockController)
			periods        = mocks.NewMockPeriodsRepository(mockController)
			service        = mocks.NewMockInvoicesService(mockController)
			calendar       = mocks.NewMockCalendarService(mockController)
			logger         = mocks.NewMockLogger(mockController)
			uow            = mocks.NewMockUnitOfWork(mockController)
		)

		periods.EXPECT().GetLast(gomock.Any()).Times(1).Return(&entities.Period{
			CreatedAt: datetime.MustParseDateTime("2024-01-01 00:00"),
		}, nil)

		invoices.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(1).Return(&entities.Invoice{
			ID:     1,
			Date:   datetime.MustParseDateTime("2024-01-02 00:02"),
			Status: entities.InvoiceStatusPending,
			Template: &entities.InvoiceTemplate{
				ID: 1,
			},
		}, nil)

		uow.EXPECT().UpdateInvoice(gomock.Any()).Times(1).Do(func(invoice entities.Invoice) {
			assert.Equal(t, uint64(1), invoice.ID)
			assert.Equal(t, datetime.MustParseDateTime("2024-01-01 00:01"), invoice.Date)
			assert.Equal(t, entities.InvoiceStatusPending, invoice.Status)
			assert.Nil(t, invoice.Template)
		})

		contacts.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(1).Return(&entities.Contact{
			ID: 1,
		}, nil)

		uow.EXPECT().Commit(gomock.Any()).Times(1).Return(nil)

		calendar.EXPECT().Sync(gomock.Any()).Times(1).Return(nil)

		useCase := update_invoice.NewUseCase(periods, invoices, contacts, service, calendar, logger, uow)

		err := useCase.Execute(context.Background(), port.UpdateInvoiceRequest{
			InvoiceID: 1,
			ContactID: 1,
			Date:      datetime.MustParseDateTime("2024-01-01 00:01"),
		})

		assert.NoError(t, err)
	})

	t.Run("should update all invoices after date with change templates", func(t *testing.T) {
		var (
			mockController = gomock.NewController(t)
			invoices       = mocks.NewMockInvoicesRepository(mockController)
			contacts       = mocks.NewMockContactsRepository(mockController)
			periods        = mocks.NewMockPeriodsRepository(mockController)
			service        = mocks.NewMockInvoicesService(mockController)
			calendar       = mocks.NewMockCalendarService(mockController)
			logger         = mocks.NewMockLogger(mockController)
			uow            = mocks.NewMockUnitOfWork(mockController)
		)

		periods.EXPECT().GetLast(gomock.Any()).Times(1).Return(&entities.Period{
			CreatedAt: datetime.MustParseDateTime("2024-01-01 00:00"),
		}, nil)

		invoices.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(1).Return(&entities.Invoice{
			ID:     3,
			Date:   datetime.MustParseDateTime("2024-01-02 00:02"),
			Status: entities.InvoiceStatusPending,
			Template: &entities.InvoiceTemplate{
				ID: 1,
			},
		}, nil)

		contacts.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(1).Return(&entities.Contact{
			ID: 1,
		}, nil)

		invoices.EXPECT().GetByTemplateID(gomock.Any(), gomock.Any()).Times(1).Return([]entities.Invoice{
			{
				ID:   1,
				Date: datetime.MustParseDateTime("2024-01-01 13:02"),
				Template: &entities.InvoiceTemplate{
					ID: 1,
				},
				Status: entities.InvoiceStatusPaid,
			},
			{
				ID:   2,
				Date: datetime.MustParseDateTime("2024-01-02 13:02"),
				Template: &entities.InvoiceTemplate{
					ID: 1,
				},
				Status: entities.InvoiceStatusPending,
			},
			{
				ID:   3,
				Date: datetime.MustParseDateTime("2024-01-03 13:02"),
				Template: &entities.InvoiceTemplate{
					ID: 1,
				},
				Status: entities.InvoiceStatusPending,
			},
			{
				ID:   4,
				Date: datetime.MustParseDateTime("2024-01-04 13:02"),
				Template: &entities.InvoiceTemplate{
					ID: 1,
				},
				Status: entities.InvoiceStatusPending,
			},
			{
				ID:   5,
				Date: datetime.MustParseDateTime("2024-01-05 13:02"),
				Template: &entities.InvoiceTemplate{
					ID: 1,
				},
				Status: entities.InvoiceStatusPending,
			},
		}, nil)

		uow.EXPECT().DeleteTemplate(gomock.Any()).Times(1).Do(func(id uint64) {
			assert.Equal(t, uint64(1), id)
		})

		uow.EXPECT().DeleteInvoice(gomock.Any()).Times(3).Do(func(id uint64) {
			assert.GreaterOrEqual(t, uint64(3), id)
		})

		uow.EXPECT().Commit(gomock.Any()).Times(1).Return(nil)

		service.EXPECT().GetScheduledInvoices(gomock.Any(), gomock.Any()).Times(1).Return([]entities.Invoice{
			{
				ID: 3,
				Template: &entities.InvoiceTemplate{
					ID: 2,
				},
				Status: entities.InvoiceStatusPending,
			},
			{
				ID: 4,
				Template: &entities.InvoiceTemplate{
					ID: 2,
				},
				Status: entities.InvoiceStatusPending,
			},
			{
				ID: 5,
				Template: &entities.InvoiceTemplate{
					ID: 2,
				},
				Status: entities.InvoiceStatusPending,
			},
		}, nil)

		uow.EXPECT().CreateInvoices(gomock.Any(), gomock.Any()).Times(1).Do(func(invoices []entities.Invoice, template entities.InvoiceTemplate) {
			var startedID uint64 = 3
			for _, invoice := range invoices {
				assert.Equal(t, startedID, invoice.ID)
				assert.Equal(t, entities.InvoiceStatusPending, invoice.Status)

				startedID++
			}
		})

		calendar.EXPECT().Sync(gomock.Any()).Times(1).Return(nil)

		useCase := update_invoice.NewUseCase(periods, invoices, contacts, service, calendar, logger, uow)

		err := useCase.Execute(context.Background(), port.UpdateInvoiceRequest{
			InvoiceID: 1,
			ContactID: 1,
			Plan: &entities.RepeatPlanner{
				Interval:      planner.PlanRepeatIntervalMonthly,
				IntervalCount: 1,
				EndCount:      pointer.Ptr(uint32(3)),
			},
			Date: datetime.MustParseDateTime("2024-01-03 15:00"),
		})

		assert.NoError(t, err)
	})

	t.Run("should log warn when calendar sync fails and updating one invoice", func(t *testing.T) {
		var (
			mockController = gomock.NewController(t)
			invoices       = mocks.NewMockInvoicesRepository(mockController)
			contacts       = mocks.NewMockContactsRepository(mockController)
			periods        = mocks.NewMockPeriodsRepository(mockController)
			service        = mocks.NewMockInvoicesService(mockController)
			calendar       = mocks.NewMockCalendarService(mockController)
			logger         = mocks.NewMockLogger(mockController)
			uow            = mocks.NewMockUnitOfWork(mockController)
		)

		periods.EXPECT().GetLast(gomock.Any()).Times(1).Return(&entities.Period{
			CreatedAt: datetime.MustParseDateTime("2024-01-01 00:00"),
		}, nil)

		invoices.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(1).Return(&entities.Invoice{
			ID:     1,
			Date:   datetime.MustParseDateTime("2024-01-02 00:02"),
			Status: entities.InvoiceStatusPending,
			Template: &entities.InvoiceTemplate{
				ID: 1,
			},
		}, nil)

		uow.EXPECT().Commit(gomock.Any()).Times(1).Return(nil)

		uow.EXPECT().UpdateInvoice(gomock.Any()).Times(1).Do(func(invoice entities.Invoice) {
			assert.Equal(t, uint64(1), invoice.ID)
			assert.Equal(t, datetime.MustParseDateTime("2024-01-01 00:01"), invoice.Date)
			assert.Equal(t, entities.InvoiceStatusPending, invoice.Status)
			assert.Nil(t, invoice.Template)
		})

		contacts.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(1).Return(&entities.Contact{
			ID: 1,
		}, nil)

		calendar.EXPECT().Sync(gomock.Any()).Times(1).Return(errors.New("sync error"))

		logger.EXPECT().Warn(gomock.Any(), gomock.Any()).Times(1).Do(func(msg string, fields ...log.Field) {
			assert.Equal(t, "error", fields[0].Key)
			assert.Equal(t, "cannot sync calendar", msg)
		})

		useCase := update_invoice.NewUseCase(
			periods,
			invoices,
			contacts,
			service,
			calendar,
			logger,
			uow,
		)

		err := useCase.Execute(context.Background(), port.UpdateInvoiceRequest{
			InvoiceID: 1,
			ContactID: 1,
			Date:      datetime.MustParseDateTime("2024-01-01 00:01"),
		})

		assert.NoError(t, err)
	})
}
