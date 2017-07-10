// Implements a small module to track updates of users followed (and download/broadcast them)

package tracker

import (
	"context"

	"github.com/DeveloppSoft/go-ipfs-api"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("tracker")
)

type Tracker struct {
	c  context.Context
	sh *shell.Shell
	meta MetaTracker
}

func NewTracker(ctx context.Context, meta_path string, ipfs_api string) (Tracker, error) {
	checkAndMake(meta_path)

	my_meta, err := MetaOpen(meta_path + "/following.json")
	if err != nil {
		return Tracker{}, err
	}

	err = my_meta.Parse()
	if err != nil {
		return Tracker{}, err
	}

	return Tracker{c: ctx, sh: shell.NewShell(ipfs_api), meta: my_meta}, nil
}

func (t *Tracker) Start() {
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

		// Get the id of the number actually pinned
		actual := t.meta.Following[id]

		if actual != resolve {
			log.Info("updating " + id)
			// Quite simple: unfollow the old one, follow the new one
			t.UnFollow(actual)
			err = t.Follow(id)
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

	// We save the current hash
	// We need data's hash
	hash, err := t.sh.Resolve(id)
	if err != nil {
		return err
	}

	// Now we write it
	err = t.meta.Append(id, hash)
	if err != nil {
		return err
	}

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
