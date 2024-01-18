package configs

type AppConfig struct {
	CheckpointPeriodAt uint8 `env:"CHECKPOINT_PERIOD_AT" envDefault:"15"`
	FinancialYearStart uint8 `env:"FINANCIAL_YEAR_START" envDefault:"1"`
}
