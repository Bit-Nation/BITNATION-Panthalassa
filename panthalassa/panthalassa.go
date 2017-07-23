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

package main

import (
	"context"
	"flag"
	"log"

	"github.com/Bit-Nation/BITNATION-Panthalassa/api"
	"github.com/Bit-Nation/BITNATION-Panthalassa/config"
	"github.com/Bit-Nation/BITNATION-Panthalassa/repo"
	"github.com/Bit-Nation/BITNATION-Panthalassa/tracker"
)

func main() {
	// TODO: clean exit
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flag.Parse()

	// Be clear and honest!
	log.Println("WARNING: the post quantum safe functions are not yet implemented, use at your own risks\n")

	// Banner ;)
	// TODO: something more "kickass"
	log.Println("#######################################################################################")
	log.Println("#                                                                                     #")
	log.Println("#                               Welcome to Panthalassa                                #")
	log.Println("#                                                                                     #")
	log.Println("#                   Brought to you by BITNATION (https://bitnation.co)                #")
	log.Println("#                           A creation of Eliott Teissonniere                         #")
	log.Println("#                                                                                     #")
	log.Println("#######################################################################################")

	// Load the config
	config.Load()

	// Make the repo
	r := repo.NewLedger(config.Conf.RepoPath, config.Conf.IpfsApi)

	// Load tracker
	t, err := tracker.NewTracker(ctx, config.Conf.MetaPath, config.Conf.IpfsApi)
	if err != nil {
		log.Println(err)
		return
	}
	go t.Start()

	// Start the remote api
	my_api := api.NewAPI(config.Conf.APIListenAddr, r, t)
	my_api.Run()
}
