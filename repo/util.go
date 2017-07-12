package repo

import (
	"os"
	"io/ioutil"
)

func checkAndMake(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm) // TODO: better mode
	}
}

func checkAndMakeFile(path string, content []byte) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		ioutil.WriteFile(path, content, os.ModePerm)
	}
}
