package main

import (
	"log"
	"flag"

	"github.com/Bit-Nation/BITNATION-Panthalassa/api"
	"github.com/Bit-Nation/BITNATION-Panthalassa/repo"
	"github.com/Bit-Nation/BITNATION-Panthalassa/config"

)

func main() {
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

	// Start the remote api
	my_api := api.NewAPI(config.Conf.APIListenAddr, r)
	my_api.Run()
}
