package repo

import (
	"encoding/json"
)

type About struct {
	Pseudo      string
	Image       string
	ETHAddress  string
	Description string
}

func (a *About) Export() ([]byte, error) {
	return json.Marshal(a)
}

func Import(b []byte) (About, error) {
	var a About
	err := json.Unmarshal(b, &a)

	return a, err
}
