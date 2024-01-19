//go:generate mockgen -source=internal/services/accounts/service.go -destination=internal/mocks/services/accounts.go -package=mocks
//go:generate mockgen -source=internal/repositories/bbolt/currency_rates.go -destination=internal/mocks/repositories/bbolt/currency_rates.go -package=mocks

// use cases
//go:generate mockgen -source=internal/use-cases/create-invoice/use_case.go -destination=internal/mocks/use-cases/create-invoice/use_case.go -package=mocks
//go:generate mockgen -source=internal/use-cases/update-invoice/use_case.go -destination=internal/mocks/use-cases/update-invoice/use_case.go -package=mocks
//go:generate mockgen -source=internal/use-cases/delete-invoice/use_case.go -destination=internal/mocks/use-cases/delete-invoice/use_case.go -package=mocks

package main
