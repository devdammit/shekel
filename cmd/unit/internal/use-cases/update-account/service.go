package update_account

import (
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
	"github.com/devdammit/shekel/pkg/currency"
)

type Repository interface {
	GetByID(accountID uint64) (*entities.Account, error)
	Update(account *entities.Account) error
}

type UpdateAccountUseCase struct {
	repo Repository
}

// Execute
// Можно обновить только имя, описание и баланс.
func (u *UpdateAccountUseCase) Execute(ID uint64, name string, description *string, balance currency.Amount) (bool, error) {
	account, err := u.repo.GetByID(ID)
	if err != nil {
		return false, err
	}

	if account.Balance.CurrencyCode != balance.CurrencyCode {
		return false, port.ErrorCannotUpdateCurrencyInAccount
	}

	account.Name = name
	account.Description = description
	account.Balance = balance

	err = u.repo.Update(account)
	if err != nil {
		return false, err
	}

	return true, nil
}
