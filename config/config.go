/*
Copyright 2017 Eliott Teissonniere

Permission is hereby granted, free of charge, to any person
obtaining a copy of this software and associated documentation
files (the "Software"), to deal in the Software without restriction,
including without limitation the rights to use, copy, modify, merge,
publish, distribute, sublicense, and/or sell copies of the Software,
and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included
in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package config

import (
	"flag"

	"encoding/json"
	"io/ioutil"
)

type Config struct {
	IpfsApi  string
	RepoPath string
	MetaPath string

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
