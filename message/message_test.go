package message

import (
	"time"
	"errors"
	"testing"
)

func TestEncodeAndDecode(t *testing.T) {
	msg := Message{
		From: "eliott",
		Previous: "00000000",
		Seq: 42,
		Timestamp: time.Now(),
		Content: "BITNATION rocks",
		Hash: "0123456789abcdef",
		Signature: "0123456789abcdef",
	}

	encoded, err := msg.ToBytes()
	if err != nil {
		t.Error(err)
	}

	t.Logf("Encode: %x", encoded)

	decoded, err := FromBytes(encoded)
	if err != nil {
		t.Error(err)
	}

	if msg != decoded {
		t.Error(errors.New("decoded message differs from encoded one"))
	}
}
