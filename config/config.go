package config

import (
	"flag"

	"encoding/json"
	"io/ioutil"
)

type Config struct {
	IpfsApi string
	RepoPath string

	APIListenAddr string
}

var (
	Conf       Config
	ConfigPath string
)

func init() {
	const (
		configPathDefault = "config.json"
		configPathUsage   = "path to config file"
	)

	flag.StringVar(&ConfigPath, "config", configPathDefault, configPathUsage)
	flag.StringVar(&ConfigPath, "c", configPathDefault, configPathUsage)
}

func LoadFrom(fpath string) error {
	var content []byte

	content, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, &Conf)
}

func Load() error {
	return LoadFrom(ConfigPath)
}
