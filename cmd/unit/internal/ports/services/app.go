package services

import "github.com/devdammit/shekel/cmd/unit/internal/entities"

type AppInfo struct {
	Initialized  bool
	ActivePeriod *entities.Period

	Version string
}
