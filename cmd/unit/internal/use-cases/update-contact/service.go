package update_contact

import (
	"context"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
)

type ContractsRepository interface {
	GetByID(ctx context.Context, id uint64) (*entities.Contact, error)
	Update(ctx context.Context, contact entities.Contact) error
}

type UseCase struct {
	contracts ContractsRepository
}

func NewUseCase(contracts ContractsRepository) *UseCase {
	return &UseCase{
		contracts: contracts,
	}
}

func (u *UseCase) Execute(ctx context.Context, req port.UpdateContactRequest) error {
	contact, err := u.contracts.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	contact.Name = req.Name
	contact.Text = req.Text

	err = u.contracts.Update(ctx, *contact)
	if err != nil {
		return err
	}

	return nil
}
