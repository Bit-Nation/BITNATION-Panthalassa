// Implements the Chacha RNG
package pqcrypto

import (
	"errors"

	rand "github.com/DeveloppSoft/go-rand"
)

var (
	ErrBytes = errors.New("didn't got enough bytes")
)

// Return 'nb' random bytes, thanks to chacha20
func GetRandom(nb int) ([]byte, error) {
	var seed [32]byte

	r, err := rand.New(seed[:])
	if err != nil {
		return nil, err
	}

	var buffer = make([]byte, nb)
	n, err := r.Read(buffer[:])
	if err != nil {
		return nil, err
	}

	if n != len(buffer) {
		return nil, ErrBytes
	}

	return buffer, nil
}
