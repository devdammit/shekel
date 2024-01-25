package graphql

import (
	"context"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	repoport "github.com/devdammit/shekel/cmd/unit/internal/ports/repositories"
	serviceport "github.com/devdammit/shekel/cmd/unit/internal/ports/services"
	"github.com/devdammit/shekel/pkg/currency"

	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

type UseCases struct {
	Initialize interface {
		Execute(ctx context.Context, startDate datetime.Date) error
	}
	CreateAccount interface {
		Execute(ctx context.Context, params port.CreateAccountParams) (bool, error)
	}
	CreateContact interface {
		Execute(ctx context.Context, request port.CreateContactRequest) error
	}
	SetQRCodeToContact interface {
		Execute(ctx context.Context, contactID uint64, file port.ContactQRCode) error
	}
	RemoveQRCodeFromContact interface {
		Execute(ctx context.Context, contactID uint64, bankName string) error
	}
	DeleteContact interface {
		Execute(ctx context.Context, contactID uint64) error
	}
	UpdateContact interface {
		Execute(ctx context.Context, req port.UpdateContactRequest) error
	}
	UpdateAccount interface {
		Execute(ctx context.Context, ID uint64, name string, description *string, balance currency.Amount) (bool, error)
	}
	DeleteAccount interface {
		Execute(ctx context.Context, ID uint64) (bool, error)
	}
	CreateInvoice interface {
		Execute(ctx context.Context, request port.CreateInvoiceRequest) error
	}
	UpdateInvoice interface {
		Execute(ctx context.Context, request port.UpdateInvoiceRequest) error
	}
}

type Reader struct {
	Contacts interface {
		GetAll(ctx context.Context, withDeleted *bool) ([]entities.Contact, error)
	}
	App interface {
		GetInfo(ctx context.Context) (*serviceport.AppInfo, error)
	}
	Periods interface {
		GetAll(ctx context.Context, limit *uint64, offset *uint64) ([]entities.Period, error)
		GetByID(ctx context.Context, ID uint64) (*entities.Period, error)
	}
	Accounts interface {
		GetAll(ctx context.Context, withDeleted *bool) ([]entities.Account, error)
	}
	Invoices interface {
		FindByDates(ctx context.Context, req repoport.FindByDatesRequest) ([]entities.Invoice, error)
	}
}

type Services struct {
	Periods interface {
		GetEstimatedEndDate(ctx context.Context) (datetime.DateTime, error)
	}
}

type Resolver struct {
	UseCases UseCases
	Reader   Reader
	Services Services
}
