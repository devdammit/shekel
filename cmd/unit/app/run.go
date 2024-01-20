package app

import (
	"flag"
	"github.com/devdammit/shekel/pkg/service"
	"os"
	"sync"
	"time"
)

const appName = "shekel_unit"

var (
	fs = flag.NewFlagSet(appName, flag.ExitOnError)

	dbPath          = fs.String("db-path", "data/unit.db", "Database path")
	env             = fs.String("env", "dev", "Environment")
	addr            = fs.String("addr", ":8080", "Kitchen addr")
	diagnosticsAddr = fs.String("diagnostics-addr", ":7071", "Kitchen diagnostics addr")
	shutdownTimeout = fs.Duration("shutdown-timeout", time.Second*30, "Graceful shutdown timeout")
)

func Run() *sync.WaitGroup {
	_ = fs.Parse(os.Args[1:]) // exit on error

	service.Init(appName, *env)
	service.StartDiagnosticsServer(*diagnosticsAddr)

	return service.RunWait()
}
