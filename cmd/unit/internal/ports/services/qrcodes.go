package services

import (
	"errors"
	"io"
)

type Image struct {
	Content     io.ReadSeeker
	Name        string
	Size        int64
	ContentType string
}

var (
	ErrImageNotFound = errors.New("image not found")
)
