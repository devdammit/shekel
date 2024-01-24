package set_qrcode_to_contact

import (
	"context"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	serviceport "github.com/devdammit/shekel/cmd/unit/internal/ports/services"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
)

type contactsRepository interface {
	SetQRCode(ctx context.Context, contactID uint64, file entities.QRCode) error
}

type qrCodesService interface {
	Parse(ctx context.Context, image serviceport.Image) (string, error)
}

type UseCase struct {
	contacts contactsRepository
	qrcodes  qrCodesService
}

func NewUseCase(contactsRepo contactsRepository, service qrCodesService) *UseCase {
	return &UseCase{
		contacts: contactsRepo,
		qrcodes:  service,
	}
}

func (uc *UseCase) Execute(ctx context.Context, contactID uint64, file port.ContactQRCode) error {
	qrContent, err := uc.qrcodes.Parse(ctx, serviceport.Image{
		Content:     file.Image.Content,
		Name:        file.Image.Name,
		Size:        file.Image.Size,
		ContentType: file.Image.ContentType,
	})
	if err != nil {
		return err
	}

	err = uc.contacts.SetQRCode(ctx, contactID, entities.QRCode{
		BankName: file.BankName,
		Content:  qrContent,
	})
	if err != nil {
		return err
	}

	return nil
}
