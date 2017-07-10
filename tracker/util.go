package tracker

import (
	"os"
)

func checkAndMake(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm) // TODO: better mode
	}
}
