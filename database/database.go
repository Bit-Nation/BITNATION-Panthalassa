// Implement the database of panthalassa
package database

import (
	"errors"
	"regexp"

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

// Return all messages linked a specific pubkey (currently represented as a string)
func (d *DB) GetFeed(pubkey string) ([]message.Message, error) {
	// Make an empty slice
	result := make([]message.Message, 0)

	iter := d.NewIterator(nil, nil)
	for iter.Next() {
		msg, err := message.FromBytes(iter.Value())
		if err != nil {
			continue // Ok, can't load that one but who knows...
		}

		if msg.From == pubkey {
			result = append(result, msg)
		}
	}

	iter.Release()

	err := iter.Error()
	if err != nil {
		// We return result, just in case...
		return result, err
	}

	return result, nil
}

// Search some messages by using a regexp
func (d *DB) Search(query string) ([]message.Message, error) {
	re, err := regexp.Compile(query)
	if err != nil {
		return nil, err
	}

	// Make an empty slice
	result := make([]message.Message, 0)

	iter := d.NewIterator(nil, nil)
	for iter.Next() {
		msg_bytes := iter.Value()

		if re.Match(msg_bytes) {
			msg, err := message.FromBytes(msg_bytes)
			if err != nil {
				continue // That one failed, but maybe some won't
			}

			result = append(result, msg)
		}
	}

	iter.Release()

	err = iter.Error()
	if err != nil {
		// We return result, just in case...
		return result, err
	}

	return result, nil
}
