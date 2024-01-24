package entities

import (
	"time"

	"github.com/devdammit/shekel/pkg/pointer"
	"github.com/devdammit/shekel/pkg/types/datetime"
)

type Period struct {
	ID uint64 `json:"id"`

	SequenceOfYear uint8 `json:"sequence_of_year"`

	CreatedAt datetime.DateTime  `json:"created_at"`
	ClosedAt  *datetime.DateTime `json:"closed_at"`
}

func (p *Period) Close() {
	p.ClosedAt = pointer.Ptr(datetime.NewDateTime(time.Now()))
}
