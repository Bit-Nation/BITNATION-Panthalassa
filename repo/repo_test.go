package repo

import (
	"testing"
)

func TestWhoami(t *testing.T) {
	repo := NewLedger("/tmp/test_panthalssa", "localhost:5001")
	t.Logf(repo.Whoami())
}
