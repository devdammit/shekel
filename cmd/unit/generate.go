//go:generate mockgen -source=internal/services/accounts/service.go -destination=internal/mocks/services/accounts.go -package=mocks
//go:generate mockgen -source=internal/repositories/bbolt/currency_rates.go -destination=internal/mocks/repositories/bbolt/currency_rates.go -package=mocks
package main
