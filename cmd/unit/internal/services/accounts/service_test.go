package accounts_test

import (
	"context"
	"testing"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	mocks "github.com/devdammit/shekel/cmd/unit/internal/mocks/services"
	service "github.com/devdammit/shekel/cmd/unit/internal/services/accounts"
	"github.com/devdammit/shekel/pkg/currency"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_CalculateBalances(t *testing.T) {
	t.Run("should calculate balances", func(t *testing.T) {
		testCases := []struct {
			name         string
			accounts     []entities.Account
			chunk        []entities.Transaction
			convertCount int
			expected     map[uint64]entities.Account
		}{
			{
				name: "should calculate balances",
				accounts: []entities.Account{
					{
						ID: 1,
						Balance: currency.Amount{
							Value:        100,
							CurrencyCode: currency.USD,
						},
					},
				},
				chunk: []entities.Transaction{
					{
						ID: 1,
						Amount: currency.Amount{
							Value:        100,
							CurrencyCode: currency.USD,
						},
						From: &entities.Account{
							ID: 1,
						},
					},
				},
				convertCount: 0,
				expected: map[uint64]entities.Account{
					1: {
						ID: 1,
						Balance: currency.Amount{
							Value:        0,
							CurrencyCode: currency.USD,
						},
					},
				},
			},
			{
				name: "should calculate balances with currency conversion",
				accounts: []entities.Account{
					{
						ID: 1,
						Balance: currency.Amount{
							Value:        100,
							CurrencyCode: currency.USD,
						},
					},
				},
				chunk: []entities.Transaction{
					{
						ID: 1,
						Amount: currency.Amount{
							Value:        100,
							CurrencyCode: currency.EUR,
						},
						From: &entities.Account{
							ID: 1,
						},
					},
				},
				convertCount: 1,
				expected: map[uint64]entities.Account{
					1: {
						ID: 1,
						Balance: currency.Amount{
							Value:        0,
							CurrencyCode: currency.USD,
						},
					},
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var (
					mockController = gomock.NewController(t)
					repository     = mocks.NewMockRepository(mockController)
					currencyRates  = mocks.NewMockCurrencyRates(mockController)
				)

				repository.EXPECT().GetAll().Return(tc.accounts, nil)

				currencyRates.EXPECT().Convert(gomock.Any(), gomock.Any(), gomock.Any()).Times(tc.convertCount).Return(&currency.Amount{
					Value:        100,
					CurrencyCode: currency.USD,
				}, nil)

				s := service.NewService(repository, currencyRates)

				accounts, err := s.CalculateBalances(context.Background(), tc.chunk)

				assert.NoError(t, err)
				assert.Equal(t, tc.expected, accounts)
			})
		}
	})
}
