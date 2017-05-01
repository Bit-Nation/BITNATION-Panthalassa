// Implements the Chacha RNG
package pqcrypto

import (
	"errors"

	rand "github.com/DeveloppSoft/go-rand"
	log "github.com/inconshreveable/log15"
)

var logger_rng = log.New("pqcrypto", "rng")

// Return 'nb' random bytes, thanks to chacha20
func GetRandom(nb int) ([]byte, error) {
	var seed [32]byte

	r, err := rand.New(seed[:])
	if err != nil {
		logger_rng.Error(err.Error())
		return nil, err
	}

	var buffer = make([]byte, nb)
	n, err := r.Read(buffer[:])
	if err != nil {
		logger_rng.Error(err.Error())
		return nil, err
	}

	if n != len(buffer) {
		err = errors.New("didn't got enough bytes")
		logger_rng.Error(err.Error())
		return nil, err
	}

	return buffer, nil
}
