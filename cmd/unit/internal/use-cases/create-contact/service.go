package create_contact

import (
	"context"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
)

type ContactsRepository interface {
	Create(ctx context.Context, contact entities.Contact) (*entities.Contact, error)
}

type CreateContactUseCase struct {
	contacts ContactsRepository
}

func NewCreateContactUseCase(contactsRepo ContactsRepository) *CreateContactUseCase {
	return &CreateContactUseCase{
		contacts: contactsRepo,
	}
}

func (uc *CreateContactUseCase) Execute(ctx context.Context, request port.CreateContactRequest) error {
	_, err := uc.contacts.Create(ctx, entities.Contact{
		Name:   request.Name,
		Text:   request.Text,
		QRCode: request.QRCode,
	})
	if err != nil {
		return err
	}

	return nil
}
