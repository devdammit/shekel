package configs

import "github.com/devdammit/shekel/pkg/types/datetime"

type AppConfig struct {
	CheckpointPeriodAt datetime.Date
	FinancialYearStart datetime.Date
}
