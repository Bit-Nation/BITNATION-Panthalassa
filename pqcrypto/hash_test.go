package pqcrypto

import (
	"testing"
	"errors"
)

var (
	ErrHash = errors.New("hash verification failed")
)

const test_data = "bitnation rocks"

func TestVerify(t *testing.T) {
	hash := Hash([]byte(test_data))
	t.Logf("hash: %x", hash)

	if CheckHash([]byte(test_data), hash) != true {
		t.Error(ErrHash)
	}

	// Alter the Hash
	hash[0] = 1
	if CheckHash([]byte(test_data), hash) == true {
		t.Error(ErrHash)
	}
}
