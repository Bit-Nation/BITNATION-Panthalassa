package database

import (
	"errors"
	"os"
	"testing"
	"time"

	"context"

	"github.com/Bit-Nation/BITNATION-Panthalassa/message"
)

// TODO: avoid opening and closing the DB each times
// TODO: avoid filling the DB each times

var (
	ErrMatch = errors.New("messages doesn't match")
	ErrFeed  = errors.New("unexpected feed results")
	ErrSeq   = errors.New("unexpected message sequence")

	messages = []message.Message{
		message.Message{From: "<sample pubkey>", Previous: "", Seq: 1, Timestamp: time.Now(), Content: message.MessageContent{Type: "test", Data: "<content to retrieve>"}, Hash: "<sample hash 1>", Signature: "<sample sign>"},
		message.Message{From: "<sample pubkey>", Previous: "<sample hash 1>", Seq: 2, Timestamp: time.Now(), Content: message.MessageContent{Type: "test", Data: "<content to retrieve>"}, Hash: "<sample hash 2>", Signature: "<sample sign>"},
		message.Message{From: "<sample pubkey>", Previous: "<sample hash 2>", Seq: 3, Timestamp: time.Now(), Content: message.MessageContent{Type: "test", Data: "<sample content>"}, Hash: "<sample hash 3>", Signature: "<sample sign>"},
		message.Message{From: "<sample pubkey>", Previous: "<sample hash 3>", Seq: 4, Timestamp: time.Now(), Content: message.MessageContent{Type: "test", Data: "<sample content>"}, Hash: "<sample hash 4>", Signature: "<sample sign>"},

		message.Message{From: "<sample pubkey remove>", Previous: "", Seq: 1, Timestamp: time.Now(), Content: message.MessageContent{Type: "test", Data: "<sample content>"}, Hash: "<sample hash>", Signature: "<sample sign>"},
	}

	invalidMessages = []message.Message{
		message.Message{From: "<sample new pubkey remove>", Previous: "", Seq: 1, Timestamp: time.Now(), Content: message.MessageContent{Type: "test", Data: "<new content>"}, Hash: "<sample new hash 1>", Signature: "<sample new sign>"},
		// Invalid sequence
		message.Message{From: "<sample new pubkey remove>", Previous: "<sample new hash 1", Seq: 1, Timestamp: time.Now(), Content: message.MessageContent{Type: "test", Data: "<new content>"}, Hash: "<sample new hash 2>", Signature: "<sample new sign>"},
	}
)

func init() {
	os.RemoveAll("/tmp/test.db") // Delete previous attempts
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

	for idx, msg := range invalidMessages {
		err := db.AddMessage(msg)
		if idx == 0 && err != nil {
			t.Error(err)
		} else if idx == 1 && err == nil {
			t.Error(ErrSeq)
		}
	}
}

// This test is stupid but who knows...
func TestGetLastMessage(t *testing.T) {
	db := DB{File: "/tmp/test.db"}
	db.Open()
	defer db.Close()

	msg, err := db.GetLastMessage(messages[3].From)
	if err != nil {
		t.Error(err)
	}

	if msg != messages[3] {
		t.Error(ErrMatch)
	}
}

func TestRemove_ExistingMessage(t *testing.T) {
	db := DB{File: "/tmp/test.db"}
	db.Open()
	defer db.Close()

	err := db.Remove("<sample hash>")
	if err != nil {
		t.Error(err)
	}
}

// TODO: better tests
func TestGetFeed(t *testing.T) {
	db := DB{File: "/tmp/test.db"}
	db.Open()
	defer db.Close()

	doFeedTest(db, t, 4, "<sample pubkey>", "", -1, "", "")
	doFeedTest(db, t, 1, "", "", 2, "", "")
	doFeedTest(db, t, 2, "", "", -1, "", ".*retrieve.*")
	doFeedTest(db, t, 4, "<sample pubkey>", "", -1, "test", "")
}

func doFeedTest(db DB, t *testing.T, goal int, from string, previous string, seq int, msg_type string, data string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dst := make(chan message.Message)

	go db.GetFeed(ctx, dst, from, previous, seq, msg_type, data)

	counter := 0

	msg := messages[0] // Doesn't matter
	for (msg != message.Message{}) {
		msg = <-dst

		// An empty message means the GetFeed func finished its work
		if (msg != message.Message{}) {
			counter += 1
			t.Log(msg)
		}
	}

	if counter != goal {
		t.Error(ErrFeed)
	}
}

func TestGetMessage(t *testing.T) {
	db := DB{File: "/tmp/test.db"}
	db.Open()
	defer db.Close()

	msg, err := db.GetMessage("<sample hash 1>")
	if err != nil {
		t.Error(err)
	}

	if msg != messages[0] {
		t.Error(ErrMatch)
	}
}
