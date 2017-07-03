package repo

import (
	"os"
	"io/ioutil"
	"strconv"

	"github.com/DeveloppSoft/go-ipfs-api"
)

type Ledger struct {
	Repo string

	sh *shell.Shell // IPFS api
}

func NewLedger(repo_path string, ipfs_api string) Ledger {
	// Create some files if needed
	checkAndMake(repo_path)
	checkAndMake(repo_path + "/feed")
	checkAndMake(repo_path + "/ressources")

	return Ledger{Repo: repo_path, sh: shell.NewShell(ipfs_api)}
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
func (l *Ledger) GetMessage(peer string, sequence string) (string, error) {
	reader, err := l.sh.Cat(peer + "/feed/" + sequence + ".json")
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(reader)
	return string(bytes), err
}

// Get the last seq number, as a string (no need to convert)
func (l *Ledger) GetLastSeq(peer string) (string, error) {
	reader, err := l.sh.Cat(peer + "/lastseq")
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(reader)
	return string(bytes), err
}

// Get all messages from a peer, return a slice of them, ordered from the more recent to the oldest 
func (l *Ledger) GetFeed(peer string) ([]string, error) {
	result := make([]string, 0)

	seq_str, err := l.GetLastSeq(peer)
	if err != nil {
		return result, err
	}

	seq, err := strconv.Atoi(seq_str)
	if err != nil {
		return result, err
	}

	for i := seq; i > 0; i-- {
		msg, err := l.GetMessage(peer, strconv.Itoa(i))
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
func (l *Ledger) About(peer string) (string, error) {
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
	return ioutil.WriteFile(l.Repo + "/about.json", bytes, os.ModePerm)
}
