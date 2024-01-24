package entities

import "github.com/devdammit/shekel/pkg/types/datetime"

type QRCode struct {
	BankName string `json:"bank_name"`
	Content  string `json:"content"`
}

type Contact struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`

	Text    string   `json:"text"`
	QRCodes []QRCode `json:"qr_code"`

	DeletedAt *datetime.DateTime `json:"delete_at"`
	CreatedAt datetime.DateTime  `json:"created_at"`
	UpdatedAt datetime.DateTime  `json:"updated_at"`
}
