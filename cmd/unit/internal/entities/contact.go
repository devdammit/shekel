package entities

import "github.com/devdammit/shekel/pkg/types/datetime"

type Contact struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`

	Text   string `json:"text"`
	QRCode []byte `json:"qr_code"`

	DeleteAt  *datetime.DateTime `json:"delete_at"`
	CreatedAt datetime.DateTime  `json:"created_at"`
	UpdateAt  datetime.DateTime  `json:"updated_at"`
}
