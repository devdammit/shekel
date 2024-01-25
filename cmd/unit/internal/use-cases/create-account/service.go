package create_account

import (
	"context"
	"errors"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
)

type Repository interface {
	Create(ctx context.Context, account *entities.Account) (*entities.Account, error)
}

type CreateAccountUseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		repo: repo,
	}
}

func (u *CreateAccountUseCase) Execute(ctx context.Context, params port.CreateAccountParams) (bool, error) {
	account := &entities.Account{
		Name:        params.Name,
		Description: params.Description,
		Type:        params.Type,
		Balance:     params.Balance,
	}

	_, err := u.repo.Create(ctx, account)
	if err != nil {
		if errors.Is(err, entities.ErrorAccountExists) {
			return false, errors.New("account already exists")
		}

		return false, err
	}

	return true, nil
}
