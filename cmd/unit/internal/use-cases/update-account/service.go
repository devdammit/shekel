package update_account

import (
	"context"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
	"github.com/devdammit/shekel/pkg/currency"
)

type Repository interface {
	GetByID(ctx context.Context, accountID uint64) (*entities.Account, error)
	Update(ctx context.Context, account *entities.Account) (*entities.Account, error)
}

type UpdateAccountUseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *UpdateAccountUseCase {
	return &UpdateAccountUseCase{
		repo: repo,
	}
}

// Execute
// Можно обновить только имя, описание и баланс.
func (u *UpdateAccountUseCase) Execute(
	ctx context.Context,
	ID uint64,
	name string,
	description *string,
	balance currency.Amount,
) (bool, error) {
	account, err := u.repo.GetByID(ctx, ID)
	if err != nil {
		return false, err
	}

	if account.Balance.CurrencyCode != balance.CurrencyCode {
		return false, port.ErrorCannotUpdateCurrencyInAccount
	}

	account.Name = name
	account.Description = description
	account.Balance = balance

	_, err = u.repo.Update(ctx, account)
	if err != nil {
		return false, err
	}

	return true, nil
}
