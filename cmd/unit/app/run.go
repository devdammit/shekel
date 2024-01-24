package app

import (
	"flag"
	"github.com/devdammit/shekel/cmd/unit/internal/handlers/fileserve"
	"github.com/devdammit/shekel/cmd/unit/internal/services/app"
	"github.com/devdammit/shekel/cmd/unit/internal/services/qrcodes"
	create_contact "github.com/devdammit/shekel/cmd/unit/internal/use-cases/create-contact"
	delete_contact "github.com/devdammit/shekel/cmd/unit/internal/use-cases/delete-contact"
	remove_qrcode_from_contact "github.com/devdammit/shekel/cmd/unit/internal/use-cases/remove-qrcode-from-contact"
	set_qrcode_to_contact "github.com/devdammit/shekel/cmd/unit/internal/use-cases/set-qrcode-to-contact"
	update_contact "github.com/devdammit/shekel/cmd/unit/internal/use-cases/update-contact"
	"os"
	"sync"
	"time"

	"github.com/devdammit/shekel/cmd/unit/internal/handlers/graphql"
	boltrepo "github.com/devdammit/shekel/cmd/unit/internal/repositories/bbolt"
	uowbolt "github.com/devdammit/shekel/cmd/unit/internal/uow/bbolt"
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

	initializeUow := uowbolt.NewInitializeUow(bbolt, appConfig, periodsRepo)

	appSrv := app.NewService(periodsRepo)
	qrCodesSrv := qrcodes.NewService(contactsRepo)

	return service.RunWait(
		resources.NewService(bbolt),
		boltrepo.NewService(appConfig, periodsRepo, contactsRepo),
		graphql.NewServer(*addr, *shutdownTimeout, &graphql.Resolver{
			UseCases: graphql.UseCases{
				Initialize:    initialize.NewUseCase(periodsRepo, datetime.DateTimeProvider{}, initializeUow),
				CreateContact: create_contact.NewUseCase(contactsRepo, qrCodesSrv),
				SetQRCodeToContact: set_qrcode_to_contact.NewUseCase(
					contactsRepo,
					qrCodesSrv,
				),
				RemoveQRCodeFromContact: remove_qrcode_from_contact.NewUseCase(contactsRepo),
				DeleteContact:           delete_contact.NewUseCase(contactsRepo),
				UpdateContact:           update_contact.NewUseCase(contactsRepo),
				//CreateAccount: create_account.NewUseCase(), @TODO waiting repository implementation
			},
			Reader: graphql.Reader{
				Contacts: contactsRepo,
				App:      appSrv,
				Periods:  periodsRepo,
			},
		}),
		fileserve.NewServer(":8081", *shutdownTimeout, qrCodesSrv),
	)
}
