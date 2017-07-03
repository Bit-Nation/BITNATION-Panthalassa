package repo

import (
	"bufio"

	"os"
)

func checkAndMake(path string) {
	if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
		os.Mkdir(path, nil)
	}
}
