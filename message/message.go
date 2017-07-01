package message

import (
	"time"

	"encoding/json"
)

type Message struct {
	Type string
	Data string

	Timestamp time.Time
}

func (m *Message) Export() ([]byte, error) {
	return json.Marshal(m)
}

func Import(b []byte) (Message, error) {
	var m Message
	err := json.Unmarshal(b, &m)

	return m, err
}
