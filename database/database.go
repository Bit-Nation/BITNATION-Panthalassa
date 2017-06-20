// Implement the database of panthalassa
package database

import (
	"errors"
	"regexp"

	"context"

	"github.com/syndtr/goleveldb/leveldb"

	"github.com/Bit-Nation/BITNATION-Panthalassa/message"
)

type DB struct {
	*leveldb.DB

	File string
}

func (d *DB) Open() error {
	db, err := leveldb.OpenFile(d.File, nil)
	if err != nil {
		return err
	}

	d.DB = db

	return nil
}

// Add a message, only if it's a vaild one
// TODO: check seq numbers
// TODO: check previous message
func (d *DB) AddMessage(msg message.Message) error {
	// Messages are identified by their hash
	if !msg.IsValid() {
		return errors.New("invalid message")
	}

	msg_bytes, err := msg.ToBytes()
	if err != nil {
		return err
	}

	return d.Put([]byte(msg.Hash), msg_bytes, nil)
}

// Remove a message by its id (hash)
// Currently hash are represented as strings, but it may change
func (d *DB) Remove(id string) error {
	return d.Delete([]byte(id), nil)
}

func (d *DB) GetMessage(id string) (message.Message, error) {
	msg_bytes, err := d.Get([]byte(id), nil)
	if err != nil {
		return message.Message{}, err
	}

	return message.FromBytes(msg_bytes)
}

// Get feed which can be filtered by some parameters, currently:
//	- from: message author
//	- previous: get next message
//	- seq: messages with a specific seq number, useful to get the very first message, if < 0, it is ignored
//	- msg_type: filter by type
//	- data: filter by data, regexp allowed
// You cannot filter by hash or signature
// Results are sent via dst
// TODO: filter via timestamp
// TODO: seq comparison
func (d *DB) GetFeed(ctx context.Context, dst chan <- message.Message, from string, previous string, seq int, msg_type string, data string) error {
	re, err := regexp.Compile(data)
	if err != nil {
		return err
	}

	iter := d.NewIterator(nil, nil)
	for iter.Next() {
		select {
		case <-ctx.Done(): return nil
		default:
			msg, err := message.FromBytes(iter.Value())
			if err != nil {
				continue // Ok, can't load that one but who knows...
			}

			// Regexp test
			match := (data == "")
			if data != "" {
				match = re.Match([]byte(msg.Content.Data))
			}

			// We have got something, time to compare!
			if (msg.From == from || from == "") && (msg.Previous == previous || previous == "") && (msg.Seq == seq || seq < 0 ) && (msg.Content.Type == msg_type || msg_type == "") && (match) {
				dst <-msg
			}
		}
	}

	iter.Release()

	return iter.Error()
}
