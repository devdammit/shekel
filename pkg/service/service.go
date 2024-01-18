package service

import (
	"os"
	"sync"
	"syscall"

	"github.com/devdammit/shekel/pkg/log"
)

type Service struct {
	appName         string
	diagnosticsAddr string
}

// service singleton.
var svc = Service{}

type Option func(*Service)

func WithDiagnosticsServer(addr string) Option {
	return func(s *Service) {
		s.diagnosticsAddr = addr
	}
}

func Init(appName, env string, opts ...Option) {
	SetAlive(true)

	svc.appName = appName
	for _, opt := range opts {
		opt(&svc)
	}

	if env == "dev" {
		log.SetDefault(log.NewLogger(log.WithFormatter(log.FormatterText)))
	}

	log.Info("initializing app")

	if svc.diagnosticsAddr != "" {
		StartDiagnosticsServer(svc.diagnosticsAddr)
	}
}

type StartStopper interface {
	Start()
	Stop()
}

func Run(services ...StartStopper) {
	log.Info("starting app")
	for _, s := range services {
		s.Start()
	}
	SetReady(true)
	log.Info("app ready")
	log.With(log.String("signal", Wait([]os.Signal{syscall.SIGTERM, syscall.SIGINT}).String())).Info("received signal")
	log.Info("stopping")
	SetReady(false)
	// stop in reverse order
	for i := range services {
		services[len(services)-i-1].Stop()
	}
	log.Info("bye 👋")
}

func RunWait(services ...StartStopper) *sync.WaitGroup {
	log.Info("starting app")
	for _, s := range services {
		s.Start()
	}
	SetReady(true)
	log.Info("app ready")

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.With(log.String("signal", Wait([]os.Signal{syscall.SIGTERM, syscall.SIGINT}).String())).Info("received signal")
		log.Info("stopping")
		SetReady(false)
		// stop in reverse order
		for i := range services {
			services[len(services)-i-1].Stop()
		}
		log.Info("bye 👋")
	}()
	return wg
}
