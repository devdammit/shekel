package entities

import (
	"github.com/devdammit/shekel/pkg/pointer"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"time"
)

type Period struct {
	ID uint64 `json:"id"`

	CreatedAt datetime.Time  `json:"created_at"`
	ClosedAt  *datetime.Time `json:"closed_at"`
}

func (p *Period) Close() {
	p.ClosedAt = pointer.ToDateTime(time.Now())
}
