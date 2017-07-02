package config

import (
	"flag"

	"encoding/json"
	"io/ioutil"
)

type Config struct {
	IpfsPath string
	IpfsCmd string

	EthereumApiUrl string

	RepoPath string
}

var (
	Conf Config
	ConfigPath string
)

func init() {
	const (
		configPathDefault = "config.json"
		configPathUsage   = "path to config file"
	)

	flag.StringVar(&configPath, "config", configPathDefault, configPathUsage)
	flag.StringVar(&configPath, "c", configPathDefault, configPathUsage)
}

func LoadFrom(fpath string) error {
        var content []byte

	content, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, &Conf)
}

