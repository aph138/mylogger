package mylogger

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

// pre-defined levels
const (
	Info = iota
	Debug
	Error
	Critical
)

// Logger defines Logger structure
type Logger struct {
	errWriter  io.Writer
	infoWriter io.Writer
	handler    Handler
	verbose    bool
}

// New returns a new Logger
//
// Without any options, its default values are ErrWriter=os.Stderr,
// InfoWriter=os.Stdout and Verbose=false
func New(handler Handler, opts ...Option) *Logger {
	l := &Logger{
		errWriter:  os.Stderr,
		infoWriter: os.Stdout,
		handler:    handler,
		verbose:    false,
	}
	for _, o := range opts {
		o(l)
	}
	return l
}

type Message struct {
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Content string    `json:"message"`
	Prefix  string    `json:"prefix,omitempty"`
}

func (l *Logger) Log(level int, content string, prefix string) {
	m := &Message{
		Time:    time.Now(),
		Content: content,
		Prefix:  prefix,
	}
	switch level {
	case Info:
		m.Level = "Info"
		l.infoWriter.Write(l.handler.Handle(m))
	case Debug:
		m.Level = "Debug"
		if l.verbose { //check for verbosity
			l.infoWriter.Write(l.handler.Handle(m))
		}
	case Error:
		m.Level = "Error"
		l.errWriter.Write(l.handler.Handle(m))
	case Critical:
		m.Level = "Critical"
		l.errWriter.Write(l.handler.Handle(m))
		os.Exit(1) //exit app with error status code
	default:
		m.Level = fmt.Sprintf("Unkown %d", level)
		l.errWriter.Write(l.handler.Handle(m))
	}

}

// Info uses InfoWriter
// Info uses InfoWriter
func (l *Logger) Info(p string, args ...any) {
	l.Log(Info, fmt.Sprintf(p, args...), "")
}

// Debug use InfoWriter
//
// Debug only logs if you pass WithVerbosity when you're making new Logger
// or when Verbose is true.
func (l *Logger) Debug(p string, args ...any) {
	l.Log(Debug, fmt.Sprintf(p, args...), "")
}

// Error uses ErrWriter
func (l *Logger) Error(p string, args ...any) {
	l.Log(Error, fmt.Sprintf(p, args...), "")
}

// Critical uses ErrWriter
//
// Critical will log the given content and then close the app
func (l *Logger) Critical(p string, args ...any) {
	l.Log(Critical, fmt.Sprintf(p, args...), "")
}

// ErrorWithPrefix generates a random prefix and returns it.
// It can be useful for returning HTTP error to user, when you
// want to track the error but don't want to show user critical information
// with error.
func (l *Logger) ErrorWithPrefix(p string, args ...any) string {
	c := generateRandomString(5) + fmt.Sprint(time.Now().Unix())
	l.Log(Error, fmt.Sprintf(p, args...), c)
	return c
}

const RandomStringSeed = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890!@#-_."

func generateRandomString(n int) string {
	result := make([]byte, n)
	for i := range n {
		result[i] = RandomStringSeed[rand.Intn(len(RandomStringSeed))]
	}
	return string(result)
}
