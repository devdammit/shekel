package app

import (
	"flag"
	"github.com/devdammit/shekel/cmd/unit/internal/handlers/graphql"
	boltrepo "github.com/devdammit/shekel/cmd/unit/internal/repositories/bbolt"
	"github.com/devdammit/shekel/cmd/unit/internal/use-cases/initialize"
	"github.com/devdammit/shekel/internal/resources"
	"github.com/devdammit/shekel/pkg/service"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"os"
	"sync"
	"time"
)

const appName = "shekel_unit"

var (
	fs = flag.NewFlagSet(appName, flag.ExitOnError)

	graphQLPort     = fs.String("graphql-addr", "8080", "GraphQL addr")
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

	return service.RunWait(
		resources.NewService(bbolt),
		boltrepo.NewService(appConfig, periodsRepo),
		graphql.NewServer(*graphQLPort, &graphql.Resolver{
			UseCases: graphql.UseCases{
				Initialize: initialize.NewUseCase(periodsRepo, appConfig, datetime.DateTimeProvider{}),
				//CreateAccount: create_account.NewUseCase(), @TODO waiting repository implementation
			},
		}),
	)
}
