package bbolt

import "github.com/devdammit/shekel/pkg/log"

type Repository interface {
	Start() error
	GetName() string
}

type Service struct {
	repositories []Repository
}

func NewService(repositories ...Repository) *Service {
	return &Service{
		repositories: repositories,
	}
}

func (s *Service) Start() {
	log.Info("bootstrapping repositories")

	for _, r := range s.repositories {
		err := r.Start()
		if err != nil {
			log.
				With(log.Err(err), log.String("name", r.GetName())).
				Fatal("could not bootstrap repository service")
		}

		log.With(log.String("name", r.GetName())).Info("repository started")
	}

	log.Info("repositories bootstrapped")
}

func (s *Service) Stop() {}
