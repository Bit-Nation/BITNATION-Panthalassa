package message

import (
	"time"
	"bytes"

	"encoding/gob"
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

func (m *Message) ToBytes() ([]byte, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	err := enc.Encode(m)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func FromBytes(data []byte) (Message, error) {
	var buf bytes.Buffer
	var m Message

	buf.Write(data)

	dec := gob.NewDecoder(&buf)

	err := dec.Decode(&m)
	if err != nil {
		// Return an empty message
		return Message{}, err
	}

	return m, nil
}
