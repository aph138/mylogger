package mylogger

import (
	"encoding/json"
	"fmt"
	"log"
)

type Handler interface {
	Handle(*Message) []byte
}

type JsonHandler struct{}

func (j *JsonHandler) Handle(m *Message) []byte {
	result, err := json.Marshal(m)
	if err != nil {
		msg := fmt.Sprintf("err when encoding log message to json: %s \n", err)
		log.Print(msg)
		return []byte(msg)
	}
	result = append(result, '\n')
	return result
}

type TextHandler struct{}

func (t *TextHandler) Handle(m *Message) []byte {
	/// Time: 0000, Level: LLL, Prefix: AAA, Message: BBB
	var result string
	if len(m.Prefix) > 0 {
		result = fmt.Sprintf("Time=%s, Level=%s, Prefix=%s, Message=%s\n",
			m.Time, m.Level, m.Prefix, m.Content)
	} else {
		result = fmt.Sprintf("Time=%s, Level=%s, Message=%s\n",
			m.Time, m.Level, m.Content)
	}
	return []byte(result)
}
