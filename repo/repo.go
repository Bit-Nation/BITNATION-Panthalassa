package repo

import (
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
func (l *Ledger) SyncRepo() error {
	// First, add the repo to ipfs
	id, err := l.sh.AddDir(l.Repo)
	if err != nil {
		return err
	}

	// Do the ipfs name publish <id>
	// Publish for 365 days
	err = l.sh.Publish(id, "8760h")
	if err != nil {
		return err
	}
}
