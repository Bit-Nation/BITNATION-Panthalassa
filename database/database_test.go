package database

import (
	"errors"
	"testing"
	"time"

	"context"

	"github.com/Bit-Nation/BITNATION-Panthalassa/message"
)

// TODO: avoid opening and closing the DB each times
// TODO: avoid filling the DB each times

var messages = []message.Message{
	message.Message{From: "<sample pubkey>", Previous: "<sample previous msg 1>", Seq: 2, Timestamp: time.Now(), Content: message.MessageContent{Type: "test", Data:"<content to retrieve>"}, Hash: "<sample hash 2>", Signature: "<sample sign>"},
	message.Message{From: "<sample pubkey>", Previous: "<sample previous msg 2>", Seq: 3, Timestamp: time.Now(), Content: message.MessageContent{Type: "test", Data:"<sample content>"}, Hash: "<sample hash 3>", Signature: "<sample sign>"},
	message.Message{From: "<sample pubkey>", Previous: "<sample previous msg 3>", Seq: 4, Timestamp: time.Now(), Content: message.MessageContent{Type: "test", Data:"<sample content>"}, Hash: "<sample hash 4>", Signature: "<sample sign>"},
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

	doFeedTest(db, t, 3, "<sample pubkey>", "", -1, "", "")
	doFeedTest(db, t, 1, "", "", 2, "", "")
	doFeedTest(db, t, 1, "", "", -1, "", ".*retrieve.*")
	doFeedTest(db, t, 3, "<sample pubkey>", "", -1, "test", "")
}

func doFeedTest(db DB, t *testing.T, goal int, from string, previous string, seq int, msg_type string, data string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dst := make(chan message.Message)

	go db.GetFeed(ctx, dst, from, previous, seq, msg_type, data)

	counter := 0

	for counter < goal {
		select {
		case <-dst: {
			counter += 1
		}
		case <-time.After(3 * time.Second): t.Error(errors.New("feed retrieval is taking to much time"))
		}
	}
}

func TestGetMessage(t *testing.T) {
	db := DB{File: "/tmp/test.db"}
	db.Open()
	defer db.Close()

	msg, err := db.GetMessage("<sample hash 2>")
	if err != nil {
		t.Error(err)
	}

	if msg != messages[0] {
		t.Error(errors.New("messages doesn't match"))
	}
}
