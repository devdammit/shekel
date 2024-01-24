package qrcodes

import (
	"bytes"
	"context"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/services"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"image"
	"image/jpeg"
)

const (
	QrWidth  = 256
	QrHeight = 256
)

type Repository interface {
	GetByID(ctx context.Context, id uint64) (*entities.Contact, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Parse(_ context.Context, i port.Image) (string, error) {
	img, _, err := image.Decode(i.Content)
	if err != nil {
		return "", err
	}

	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}

	qrReader := qrcode.NewQRCodeReader()
	qrCode, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return "", err
	}

	return qrCode.GetText(), nil
}

func (s *Service) GetImage(_ context.Context, contactID uint64, bankName string) ([]byte, error) {
	contact, err := s.repository.GetByID(context.Background(), contactID)
	if err != nil {
		return nil, err
	}

	var text string

	for _, qrCode := range contact.QRCodes {
		if qrCode.BankName == bankName {
			text = qrCode.Content
			break
		}
	}

	if text == "" {
		return nil, port.ErrImageNotFound
	}

	buffer := new(bytes.Buffer)

	encoder := qrcode.NewQRCodeWriter()
	img, err := encoder.Encode(text, gozxing.BarcodeFormat_QR_CODE, QrWidth, QrHeight, nil)
	if err != nil {
		return nil, err
	}

	err = jpeg.Encode(buffer, img, &jpeg.Options{
		Quality: 100,
	})

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
