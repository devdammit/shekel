package app

import (
	"flag"
	"github.com/devdammit/shekel/cmd/unit/internal/handlers/fileserve"
	"github.com/devdammit/shekel/cmd/unit/internal/services/app"
	"github.com/devdammit/shekel/cmd/unit/internal/services/calendar"
	"github.com/devdammit/shekel/cmd/unit/internal/services/invoices"
	"github.com/devdammit/shekel/cmd/unit/internal/services/periods"
	"github.com/devdammit/shekel/cmd/unit/internal/services/qrcodes"
	create_invoice2 "github.com/devdammit/shekel/cmd/unit/internal/uow/bbolt/create-invoice"
	initialize2 "github.com/devdammit/shekel/cmd/unit/internal/uow/bbolt/initialize"
	update_invoice2 "github.com/devdammit/shekel/cmd/unit/internal/uow/bbolt/update-invoice"
	create_account "github.com/devdammit/shekel/cmd/unit/internal/use-cases/create-account"
	create_contact "github.com/devdammit/shekel/cmd/unit/internal/use-cases/create-contact"
	create_invoice "github.com/devdammit/shekel/cmd/unit/internal/use-cases/create-invoice"
	delete_contact "github.com/devdammit/shekel/cmd/unit/internal/use-cases/delete-contact"
	remove_qrcode_from_contact "github.com/devdammit/shekel/cmd/unit/internal/use-cases/remove-qrcode-from-contact"
	set_qrcode_to_contact "github.com/devdammit/shekel/cmd/unit/internal/use-cases/set-qrcode-to-contact"
	update_account "github.com/devdammit/shekel/cmd/unit/internal/use-cases/update-account"
	update_contact "github.com/devdammit/shekel/cmd/unit/internal/use-cases/update-contact"
	update_invoice "github.com/devdammit/shekel/cmd/unit/internal/use-cases/update-invoice"
	"github.com/devdammit/shekel/pkg/log"
	"os"
	"sync"
	"time"

	"github.com/devdammit/shekel/cmd/unit/internal/handlers/graphql"
	boltrepo "github.com/devdammit/shekel/cmd/unit/internal/repositories/bbolt"
	"github.com/devdammit/shekel/cmd/unit/internal/use-cases/initialize"
	"github.com/devdammit/shekel/internal/resources"
	"github.com/devdammit/shekel/pkg/service"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

const appName = "shekel_unit"

var (
	fs              = flag.NewFlagSet(appName, flag.ExitOnError)
	addr            = fs.String("graphql-addr", ":8080", "GraphQL addr")
	dbPath          = fs.String("db-path", "data/unit.db", "Database path")
	env             = fs.String("env", "dev", "Environment")
	diagnosticsAddr = fs.String("diagnostics-addr", ":7071", "Kitchen diagnostics addr")
	shutdownTimeout = fs.Duration("shutdown-timeout", time.Second*30, "Graceful shutdown timeout")
)

func Run() *sync.WaitGroup {
	_ = fs.Parse(os.Args[1:]) // exit on error

	service.Init(appName, *env)
	service.StartDiagnosticsServer(*diagnosticsAddr)

	bbolt := resources.NewBolt(*dbPath)

	appConfig := boltrepo.NewAppConfigRepository(bbolt)
	periodsRepo := boltrepo.NewPeriodsRepository(bbolt, datetime.DateTimeProvider{})
	contactsRepo := boltrepo.NewContactsRepository(bbolt)
	accountsRepo := boltrepo.NewAccountsRepository(bbolt)
	invoicesRepo := boltrepo.NewInvoicesRepository(bbolt)
	invoicesTemplatesRepo := boltrepo.NewInvoicesTemplatesRepository(bbolt)

	invoicesSrv := invoices.NewService(appConfig)
	calendarSrv := calendar.NewService()
	appSrv := app.NewService(periodsRepo)
	qrCodesSrv := qrcodes.NewService(contactsRepo)
	periodsSrv := periods.NewService(appConfig)

	return service.RunWait(
		resources.NewService(bbolt),
		boltrepo.NewService(appConfig, periodsRepo, contactsRepo, accountsRepo, invoicesRepo, invoicesTemplatesRepo),
		graphql.NewServer(*addr, *shutdownTimeout, &graphql.Resolver{
			UseCases: graphql.UseCases{
				Initialize: initialize.NewUseCase(
					periodsRepo,
					datetime.DateTimeProvider{},
					initialize2.NewUoW(bbolt, appConfig, periodsRepo),
				),
				CreateContact: create_contact.NewUseCase(contactsRepo, qrCodesSrv),
				SetQRCodeToContact: set_qrcode_to_contact.NewUseCase(
					contactsRepo,
					qrCodesSrv,
				),
				RemoveQRCodeFromContact: remove_qrcode_from_contact.NewUseCase(contactsRepo),
				DeleteContact:           delete_contact.NewUseCase(contactsRepo),
				UpdateContact:           update_contact.NewUseCase(contactsRepo),
				CreateAccount:           create_account.NewUseCase(accountsRepo),
				UpdateAccount:           update_account.NewUseCase(accountsRepo),
				//DeleteAccount:           delete_account.NewUseCase(accountsRepo), @TODO need transactions
				CreateInvoice: create_invoice.NewUseCase(invoicesSrv, periodsRepo, calendarSrv, create_invoice2.NewUoW(bbolt, invoicesRepo, invoicesTemplatesRepo)),
				UpdateInvoice: update_invoice.NewUseCase(periodsRepo, invoicesRepo, contactsRepo, invoicesSrv, calendarSrv, log.Default(), update_invoice2.NewUoW(bbolt, invoicesTemplatesRepo, invoicesRepo)),
			},
			Reader: graphql.Reader{
				Contacts: contactsRepo,
				App:      appSrv,
				Periods:  periodsRepo,
				Accounts: accountsRepo,
				Invoices: invoicesRepo,
			},
			Services: graphql.Services{
				Periods: periodsSrv,
			},
		}),
		fileserve.NewServer(":8081", *shutdownTimeout, qrCodesSrv),
	)
}
