package entities

import "github.com/devdammit/shekel/pkg/types/datetime"

type Contact struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`

	Text   string `json:"test"`
	QRCode []byte `json:"qr_code"`

	DeleteAt  *datetime.Time `json:"delete_at"`
	CreatedAt datetime.Time  `json:"created_at"`
	UpdateAt  datetime.Time  `json:"updated_at"`
}
