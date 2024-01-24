package remove_qrcode_from_contact

import "context"

type ContactsRepository interface {
	RemoveQRCode(ctx context.Context, contactID uint64, bankName string) error
}

type UseCase struct {
	contacts ContactsRepository
}

func NewUseCase(contactsRepo ContactsRepository) *UseCase {
	return &UseCase{
		contacts: contactsRepo,
	}
}

func (uc *UseCase) Execute(ctx context.Context, contactID uint64, bankName string) error {
	return uc.contacts.RemoveQRCode(ctx, contactID, bankName)
}
