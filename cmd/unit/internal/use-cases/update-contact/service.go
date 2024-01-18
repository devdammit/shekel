package update_contact

import (
	"context"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
)

type ContractsRepository interface {
	Find(ctx context.Context, id uint64) (*entities.Contact, error)
	Update(ctx context.Context, contact *entities.Contact) error
}

type UpdateContractUseCase struct {
	contracts ContractsRepository
}

func NewUpdateContractUseCase(contracts ContractsRepository) *UpdateContractUseCase {
	return &UpdateContractUseCase{
		contracts: contracts,
	}
}

func (u *UpdateContractUseCase) Execute(ctx context.Context, req port.UpdateContactRequest) error {
	contact, err := u.contracts.Find(ctx, req.ID)
	if err != nil {
		return err
	}

	contact.Name = req.Name
	contact.Text = req.Text
	contact.QRCode = req.QRCode

	err = u.contracts.Update(ctx, contact)
	if err != nil {
		return err
	}

	return nil
}
