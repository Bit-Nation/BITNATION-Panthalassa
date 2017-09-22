> **Disclaimer**: Panthalassa is not yet post quantum safe, this part is work in progress.
> Panthalassa is currently being refatored into TypeScript, follow the progress [here](https://github.com/Bit-Nation/BITNATION-Panthalassa-TS). 

# Panthalassa [![Build Status](https://travis-ci.org/Bit-Nation/BITNATION-Panthalassa.svg?branch=ipfs)](https://travis-ci.org/Bit-Nation/BITNATION-Panthalassa)

The Panthalassa mesh is the backend of the Pangea Jurisdiction. It's built using Secure Scuttlebutt (SSB) and Interplanetary File System (IPFS) protocols. This enables Pangea to be highly resilient and secure, even conferring resistance to emergent threats such as attacks based on high-performance quantum cryptography.


### Development

We use docker for development. Follow the steps to setup your working copy:

1. Make sure you got docker (and docker-compose, it should be included by default).
2. Run ```docker-compose up -d``` this will create the containers, setup the network and start the containers.
3. Run ```docker-compose exec app bash``` to enter into the Panthalassa container.
4. Have fun.

#### Tips and tricks
* Stop the containers ```docker-compose stop```
* Destroy the containers ```docker-compose down```
