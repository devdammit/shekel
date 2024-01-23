package configs

import (
	"github.com/devdammit/shekel/pkg/types/datetime"
)

type AppConfig struct {
	DateStart *datetime.Date `json:"financial_date_start"`
}
