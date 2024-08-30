package mylogger

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

// pre-definend levels
const (
	Info = iota
	Debug
	Error
	Critical
)

// Logger defines Logger strcutrue
type Logger struct {
	ErrWriter  io.Writer
	InfoWriter io.Writer
	Handler    Handler
	Verbose    bool
}

// New returns a new Logger
//
// Without any options, its default values are ErrWriter=os.Stderr,
// InfoWriter=os.Stdout and Verbose=false
func New(handler Handler, opts ...Option) *Logger {
	l := &Logger{
		ErrWriter:  os.Stderr,
		InfoWriter: os.Stdout,
		Handler:    handler,
		Verbose:    false,
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
		l.InfoWriter.Write(l.Handler.Handle(m))
	case Debug:
		m.Level = "Debug"
		//check for verbosity
		if l.Verbose {
			l.InfoWriter.Write(l.Handler.Handle(m))
		}
	case Error:
		m.Level = "Error"
		l.ErrWriter.Write(l.Handler.Handle(m))
	case Critical:
		m.Level = "Critical"
		l.InfoWriter.Write(l.Handler.Handle(m))
		//exit app with error status code
		os.Exit(1)
	default:
		m.Level = fmt.Sprintf("Unkown %d", level)
		l.InfoWriter.Write(l.Handler.Handle(m))
	}

}

// Info uses InfoWriter
func (l *Logger) Info(c string) {
	l.Log(Info, c, "")
}

// Debug use InfoWriter
//
// Debug only logs if you pass WithVerbosity when you're making new Logger
// or when Verbose is true.
func (l *Logger) Debug(c string) {
	l.Log(Debug, c, "")
}

// Error uses ErrWriter
func (l *Logger) Error(c string) {
	l.Log(Error, c, "")
}

// Critical uses ErrWriter
//
// Critical will log the given content and then close the app
func (l *Logger) Critical(c string) {
	l.Log(Critical, c, "")
}

// ErrorWithPrefix generates a random prefix and returns it.
// It can be useful for returning HTTP error to user, when you
// want to track the error but don't want to show user critical information
// with error.
func (l *Logger) ErrorWithPrefix(c string) string {
	p := generateRandomString(5) + fmt.Sprint(time.Now().Unix())
	l.Log(Error, c, p)
	return p
}

const RandomStringSeed = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890!@#-_."

func generateRandomString(n int) string {
	result := make([]byte, n)
	for i := range n {
		result[i] = RandomStringSeed[rand.Intn(len(RandomStringSeed))]
	}
	return string(result)
}
