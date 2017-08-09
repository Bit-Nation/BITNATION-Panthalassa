/*
Copyright 2017 Eliott Teissonniere

Permission is hereby granted, free of charge, to any person
obtaining a copy of this software and associated documentation
files (the "Software"), to deal in the Software without restriction,
including without limitation the rights to use, copy, modify, merge,
publish, distribute, sublicense, and/or sell copies of the Software,
and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included
in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package repo

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"encoding/base64"
	"encoding/json"
	"io/ioutil"

	"github.com/DeveloppSoft/go-ipfs-api"

	"golang.org/x/crypto/blake2b"
)

type LedgerInterface interface {
	Sync() error
	GetMessage(peer_name string, sequence string) (string, error)
	GetLastSeq(peer_name string) (string, error)
	GetFeed(peer_name string) ([]string, error)
	Whoami() string
	About(peer_name string) (string, error)
	SetAbout(about About) error
	Publish(data string) error
	AddRessource(b64 string) (string, error)
	GetRessource(id string) (string, error)
	Resolve(name string) (string, error)
}

type Ledger struct {
	Repo string

	sh *shell.Shell // IPFS api
}

func NewLedger(repo_path string, ipfs_api string) *Ledger {
	// Create some files if needed
	checkAndMake(repo_path)
	checkAndMake(repo_path + "/feed")
	checkAndMake(repo_path + "/ressources")
	checkAndMakeFile(repo_path+"/lastseq", []byte("0"))
	checkAndMakeFile(repo_path+"/about.json", []byte("{}"))

	return &Ledger{Repo: repo_path, sh: shell.NewShell(ipfs_api)}
}

// Recursively add stuff in the repo and do an `ipfs name publish`
func (l *Ledger) Sync() error {
	// First, add the repo to ipfs
	id, err := l.sh.AddDir(l.Repo)
	if err != nil {
		return err
	}

	// Do the ipfs name publish <id>
	// Publish for 365 days
	return l.sh.Publish(id, "8760h")
}

// Get a message, returned as a reader
func (l *Ledger) GetMessage(peer_name string, sequence string) (string, error) {
	peer, err := l.Resolve(peer_name)
	if err != nil {
		return "", err
	}

	reader, err := l.sh.Cat(peer + "/feed/" + sequence + ".json")
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(reader)
	return string(bytes), err
}

// Get the last seq number, as a string (no need to convert)
func (l *Ledger) GetLastSeq(peer_name string) (string, error) {
	peer, err := l.Resolve(peer_name)
	if err != nil {
		return "", err
	}

	reader, err := l.sh.Cat(peer + "/lastseq")
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(reader)
	return string(bytes), err
}

// Get all messages from a peer, return a slice of them, ordered from the more recent to the oldest
func (l *Ledger) GetFeed(peer_name string) ([]string, error) {
	result := make([]string, 0)

	seq_str, err := l.GetLastSeq(peer_name)
	if err != nil {
		return result, err
	}

	seq, err := strconv.Atoi(seq_str)
	if err != nil {
		return result, err
	}

	for i := seq; i > 0; i-- {
		msg, err := l.GetMessage(peer_name, strconv.Itoa(i))
		if err != nil {
			return result, err
		}

		result = append(result, msg)
	}

	return result, nil
}

// Return our id or ""
func (l *Ledger) Whoami() string {
	id, err := l.sh.ID()

	if err != nil {
		return ""
	} else {
		return id.ID
	}
}

// Just retrieve about.json
func (l *Ledger) About(peer_name string) (string, error) {
	peer, err := l.Resolve(peer_name)
	if err != nil {
		return "", err
	}

	reader, err := l.sh.Cat(peer + "/about.json")
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(reader)
	return string(bytes), err
}

// Fill the profile of our user
func (l *Ledger) SetAbout(about About) error {
	bytes, err := about.Export()
	if err != nil {
		return err
	}

	// Write that to about.json
	return ioutil.WriteFile(l.Repo+"/about.json", bytes, os.ModePerm)
}

type Message struct {
	Seq       int
	Timestamp time.Time

	Data string
}

// Add a message and increase the lastseq
func (l *Ledger) Publish(data string) error {
	seq_str, err := l.GetLastSeq(l.Whoami())
	if err != nil {
		return err
	}

	seq, err := strconv.Atoi(seq_str)
	if err != nil {
		return err
	}

	seq++
	seq_str = strconv.Itoa(seq)

	// Build the message
	msg := Message{Seq: seq, Timestamp: time.Now(), Data: data}
	msg_byte, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// Just write it to the repo
	err = ioutil.WriteFile(l.Repo+"/feed/"+seq_str+".json", msg_byte, os.ModePerm) // TODO: better perm
	if err != nil {
		return err
	}

	// Increment lastseq
	return ioutil.WriteFile(l.Repo+"/lastseq", []byte(seq_str), os.ModePerm) // TODO: better perm
}

func (l *Ledger) AddRessource(b64 string) (string, error) {
	// Unpack data
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}

	// Calculate checksum (no need for a mega high algo here, let's use md5)
	hash_bytes := blake2b.Sum512(data)
	hash := fmt.Sprintf("%x", hash_bytes)

	err = ioutil.WriteFile(l.Repo+"/ressources/"+hash, data, os.ModePerm) // Need better perms
	return hash, err
}

func (l *Ledger) GetRessource(id string) (string, error) {
	reader, err := l.sh.Cat(id)
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

func (l *Ledger) Resolve(name string) (string, error) {
	return l.sh.Resolve(name)
}
