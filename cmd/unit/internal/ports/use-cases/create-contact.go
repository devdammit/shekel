package use_cases

import (
	"io"
)

type ContactImage struct {
	Content     io.ReadSeeker
	Name        string
	Size        int64
	ContentType string
}

type ContactQRCode struct {
	Image    ContactImage
	BankName string
}

type CreateContactRequest struct {
	Name    string          `json:"name"`
	Text    string          `json:"text"`
	QRCodes []ContactQRCode `json:"qr_codes"`
}
