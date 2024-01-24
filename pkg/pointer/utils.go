package pointer

import (
	"time"

	"github.com/devdammit/shekel/pkg/types/datetime"
)

func Ptr[T any](x T) *T {
	return &x
}

func PtrOrNil[T comparable](x T) *T {
	var zero T
	if x == zero {
		return nil
	}
	return &x
}

func ToDateTime(t time.Time) *datetime.Time {
	return &datetime.Time{Time: t}
}
