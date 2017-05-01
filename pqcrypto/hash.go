// Makes use of the Blake2 hash algorithm
package pqcrypto

import (
	"golang.org/x/crypto/blake2b"
)

// OK, nothing special here
func Hash(data []byte) [blake2b.Size]byte {
	return blake2b.Sum512(data)
}

// Compute and check the hash
func CheckHash(data []byte, target [blake2b.Size]byte) bool {
	if Hash(data) == target {
		return true
	} else {
		return false
	}
}
