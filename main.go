package main

import (
	"g00chy.com/onamaedns/v2/network"
	"log"
)

func main() {
	var network = network.NewNetwork()
	network.CheckIp()
	log.Print("%s", *(network.OnameId))
}
