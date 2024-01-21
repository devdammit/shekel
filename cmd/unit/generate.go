// mock services
//go:generate mockgen -source=internal/services/accounts/service.go -destination=internal/mocks/services/accounts.go -package=mocks

// mock repositories
//go:generate mockgen -source=internal/repositories/bbolt/currency_rates.go -destination=internal/mocks/repositories/bbolt/currency_rates.go -package=mocks

// mock use cases
//go:generate mockgen -source=internal/use-cases/close-invoice/use_case.go -destination=internal/mocks/use-cases/close-invoice/use_case.go -package=mocks
//go:generate mockgen -source=internal/use-cases/create-invoice/use_case.go -destination=internal/mocks/use-cases/create-invoice/use_case.go -package=mocks
//go:generate mockgen -source=internal/use-cases/update-invoice/use_case.go -destination=internal/mocks/use-cases/update-invoice/use_case.go -package=mocks
//go:generate mockgen -source=internal/use-cases/delete-invoice/use_case.go -destination=internal/mocks/use-cases/delete-invoice/use_case.go -package=mocks

// graphql
//go:generate go run github.com/99designs/gqlgen generate

package main
