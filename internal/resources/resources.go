package resources

import "github.com/devdammit/shekel/pkg/log"

type ResourceService interface {
	GetName() string
	Start() error
	Stop() error
}

type Service struct {
	resources []ResourceService
}

func NewService(resources ...ResourceService) *Service {
	return &Service{
		resources: resources,
	}
}

func (s *Service) Start() {
	log.Info("starting resources service")
	for _, resource := range s.resources {
		err := resource.Start()
		if err != nil {
			log.With(
				log.Err(err),
				log.String("resource", resource.GetName()),
			).Fatal("could not start resource")
		}

		log.With(
			log.String("resource", resource.GetName()),
		).Info("resource started")
	}
}

func (s *Service) Stop() {
	log.Info("stopping resources service")
	for _, resource := range s.resources {
		err := resource.Stop()
		if err != nil {
			log.With(
				log.Err(err),
				log.String("resource", resource.GetName()),
			).Fatal("could not stop resource")
		}

		log.With(
			log.String("resource", resource.GetName()),
		).Info("resource stopped")
	}
}
