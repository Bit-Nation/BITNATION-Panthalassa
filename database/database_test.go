package database

import (
	"errors"
	"testing"
	"time"

	"github.com/Bit-Nation/BITNATION-Panthalassa/message"
)

// TODO: avoid opening and closing the DB each times
// TODO: avoid filling the DB each times

var messages = []message.Message{
	message.Message{From: "<sample pubkey>", Previous: "<sample previous msg 1>", Seq: 2, Timestamp: time.Now(), Content: "<sample content>", Hash: "<sample hash 2>", Signature: "<sample sign>"},
	message.Message{From: "<sample pubkey>", Previous: "<sample previous msg 2>", Seq: 3, Timestamp: time.Now(), Content: "<sample content>", Hash: "<sample hash 3>", Signature: "<sample sign>"},
	message.Message{From: "<sample pubkey>", Previous: "<sample previous msg 3>", Seq: 4, Timestamp: time.Now(), Content: "<sample content>", Hash: "<sample hash 4>", Signature: "<sample sign>"},
}

func TestOpenAndClose(t *testing.T) {
	db := DB{File: "/tmp/test.db"}

	err := db.Open()
	if err != nil {
		t.Error(err)
	}

	err = db.Close()
	if err != nil {
		t.Error(err)
	}
}

// TODO: check adding valid and invalid msg
func TestAdd(t *testing.T) {
	db := DB{File: "/tmp/test.db"}
	db.Open()
	defer db.Close()

	for _, msg := range messages {
		err := db.AddMessage(msg)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestRemove_ExistingMessage(t *testing.T) {
	db := DB{File: "/tmp/test.db"}
	db.Open()
	defer db.Close()

	for _, msg := range messages {
		err := db.AddMessage(msg)
		if err != nil {
			t.Error(err)
		}
	}

	err := db.Remove("<sample hash>")
	if err != nil {
		t.Error(err)
	}
}

func TestSearch(t *testing.T) {
	db := DB{File: "/tmp/test.db"}
	db.Open()
	defer db.Close()

	for _, msg := range messages {
		err := db.AddMessage(msg)
		if err != nil {
			t.Error(err)
		}
	}

	// Everything which contains the word "content"
	result, err := db.Search(".*content.*")

	if err != nil {
		t.Error(err)
	}

	if len(result) != 3 {
		t.Error(errors.New("unable to find messages"))
	}
}

// Ok, I admit that one is nearly the same as TestSearch
func TestGetFeed(t *testing.T) {
	db := DB{File: "/tmp/test.db"}
	db.Open()
	defer db.Close()

	for _, msg := range messages {
		err := db.AddMessage(msg)
		if err != nil {
			t.Error(err)
		}
	}

	result, err := db.GetFeed("<sample pubkey>")
	if err != nil {
		t.Error(err)
	}

	if len(result) != 3 {
		t.Error(errors.New("unable to retrieve messages"))
	}
}
