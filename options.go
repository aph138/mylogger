package mylogger

import "io"

type Option func(*Logger)

func WithCustomErrorWriter(w io.Writer) Option {
	return func(l *Logger) {
		l.ErrWriter = w
	}
}
func WithCustomInfoWriter(w io.Writer) Option {
	return func(l *Logger) {
		l.InfoWriter = w
	}
}

// WithVerbosity enables Debug's level loging
func WithVerbosity() Option {
	return func(l *Logger) {
		l.Verbose = true
	}
}
