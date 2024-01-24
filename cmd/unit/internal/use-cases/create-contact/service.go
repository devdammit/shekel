package create_contact

import (
	"context"

	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	serviceport "github.com/devdammit/shekel/cmd/unit/internal/ports/services"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/use-cases"
)

type ContactsRepository interface {
	Create(ctx context.Context, contact entities.Contact) (*entities.Contact, error)
}

type QRCodesService interface {
	Parse(ctx context.Context, image serviceport.Image) (string, error)
}

type UseCase struct {
	contacts ContactsRepository
	qrcodes  QRCodesService
}

func NewUseCase(contactsRepo ContactsRepository, service QRCodesService) *UseCase {
	return &UseCase{
		contacts: contactsRepo,
		qrcodes:  service,
	}
}

func (uc *UseCase) Execute(ctx context.Context, request port.CreateContactRequest) error {
	qrCodes := make([]entities.QRCode, 0, len(request.QRCodes))

	for _, qrCode := range request.QRCodes {
		str, err := uc.qrcodes.Parse(ctx, serviceport.Image{
			Content:     qrCode.Image.Content,
			Name:        qrCode.Image.Name,
			Size:        qrCode.Image.Size,
			ContentType: qrCode.Image.ContentType,
		})
		if err != nil {
			return err
		}
		qrCodes = append(qrCodes, entities.QRCode{
			BankName: qrCode.BankName,
			Content:  str,
		})
	}

	_, err := uc.contacts.Create(ctx, entities.Contact{
		Name:    request.Name,
		Text:    request.Text,
		QRCodes: qrCodes,
	})

	if err != nil {
		return err
	}

	return nil
}
