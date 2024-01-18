package use_cases

type CreateContactRequest struct {
	Name   string `json:"name"`
	Text   string `json:"text"`
	QRCode []byte `json:"qr_code"`
}
