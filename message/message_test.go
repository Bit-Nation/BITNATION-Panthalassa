package message

import (
	"errors"
	"testing"
	"time"
)

var (
	ErrDecoded = errors.New("decoded message differs from encoded one")
)

func TestImportExport(t *testing.T) {
	msg := Message{
		Type:      "test",
		Data:      "guess what? This is a test!",
		Timestamp: time.Now(),
		Tag:       "test",
		Seq:       1,
	}

	encoded, err := msg.Export()
	if err != nil {
		t.Error(err)
	}

	decoded, err := Import(encoded)
	if err != nil {
		t.Error(err)
	}

	if msg != decoded {
		t.Error(ErrDecoded)
	}
}
