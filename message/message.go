package message

import (
	"fmt"
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

func (m *Message) Bytes() []byte {
	// Small hack
	return []byte(fmt.Sprintf("%v", m))
}
