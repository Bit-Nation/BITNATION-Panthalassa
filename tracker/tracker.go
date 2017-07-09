// Implements a small module to track updates of users followed (and download/broadcast them)

package tracker

import (
	"context"

	"github.com/Bit-Nation/BITNATION-Panthalassa/repo"

	"github.com/DeveloppSoft/go-ipfs-api"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("tracker")
)

type Tracker struct {
	c  context.Context
	sh *shell.Shell
}

func NewTracker(ctx context.Context, ipfs_api string, rep repo.Ledger) Tracker {
	return Tracker{c: ctx, sh: shell.NewShell(ipfs_api)}
}

func (t *Tracker) Start() {
	id := 0

	for {
		select {
		case <-t.c.Done():
			return
		default: // Do the work
			following, err := t.GetUsersFollowed()
			if err != nil {
				log.Error(err)
				return
			}

			t.update(following)
		}
	}
}

func (t *Tracker) update(users []string) {
	for _, id := range users {
		log.Debug("checking user " + id)
		resolve, err := t.sh.Resolve(id)
		if err != nil {
			log.Error(err)
			return
		}

		log.Debug("user " + id + "resolve to " + resolve)

		actual, err := t.r.GetRepoID(id)
		if err != nil {
			log.Error(err)
			return
		}

		if actual != resolve {
			log.Info("updating " + id)
			// Unfollow the previous repo
			t.UnFollow(actual)

			// Just pin it
			err = t.sh.Pin(id)
			if err != nil {
				log.Error(err)
				return
			}
		}
	}
}

// Just pin the id
func (t *Tracker) Follow(id string) error {
	log.Debug("attempting to follow " + id)
	return t.sh.Pin(id)
}

// Not implemented yet
// TODO: should unpin the id
func (t *Tracker) UnFollow(id string) error {
	log.Debug("unfollowing " + id)

	return t.sh.Unpin(id)
}

// Just check the list of pinned items, return a slice of id
func (t *Tracker) GetUsersFollowed() ([]string, error) {
	result := make([]string, 0)

	pinned_items, err := t.sh.Pins()
	if err != nil {
		return result, err
	}

	for id, _ := range pinned_items {
		result = append(result, id)
	}

	return result, nil
}
