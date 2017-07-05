package main

import (
	"log"
	"flag"
	"context"

	"github.com/Bit-Nation/BITNATION-Panthalassa/api"
	"github.com/Bit-Nation/BITNATION-Panthalassa/repo"
	"github.com/Bit-Nation/BITNATION-Panthalassa/config"
	"github.com/Bit-Nation/BITNATION-Panthalassa/tracker"
)

func main() {
	// TODO: clean exit
	ctx, cancel := context.WithCancel(context.WithBackground())
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
	t := tracker.NewTracker(ctx, config.Conf.IpfsApi)
	go t.Start()

	// Start the remote api
	my_api := api.NewAPI(config.Conf.APIListenAddr, r, t)
	my_api.Run()
}
