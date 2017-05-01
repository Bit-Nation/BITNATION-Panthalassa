package pqcrypto

import (
	"testing"
)

// Test GetRandom with different buffer sizes
func TestGetRandom(t *testing.T) {
	get_random_size_test(t, 0)
	get_random_size_test(t, 1)
	get_random_size_test(t, 10)
	get_random_size_test(t, 100)
	get_random_size_test(t, 1000)
}

func get_random_size_test(t *testing.T, nb int) {
	random, err := GetRandom(nb)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("got: %x", random)
	}
}
