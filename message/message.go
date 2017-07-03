package message

import (
	"time"

	"encoding/json"
)

type Message struct {
	Type string
	Data string

	Tag string // used to set channels

	Timestamp time.Time
	Seq       uint
}

func (m *Message) Export() ([]byte, error) {
	return json.Marshal(m)
}

func Import(b []byte) (Message, error) {
	var m Message
	err := json.Unmarshal(b, &m)

	return m, err
}
