package message

import (
	"time"
)

type Message struct {
	From      string
	Previous  string
	Seq       int
	Timestamp time.Time

	Content string

	Hash      string
	Signature string
}

func (m *Message) IsValid() bool {
	return true
}
