package log

import (
	"io"
)

type options struct {
	Writer    io.Writer
	Level     *Level
	Formatter *Formatter
}

type Option func(o *options)

func WithWriter(w io.Writer) Option {
	return func(o *options) {
		o.Writer = w
	}
}

func WithLevel(l Level) Option {
	return func(o *options) {
		o.Level = &l
	}
}

func WithFormatter(f Formatter) Option {
	return func(o *options) {
		o.Formatter = &f
	}
}
