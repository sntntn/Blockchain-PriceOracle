package main

import (
	"log"

	"Blockchain-PriceOracle/app/automation"
	"Blockchain-PriceOracle/app/server"
)

func main() {

	automation.Init()

	go automation.CoinGeckoLoop()

	server := server.SetupServer()
	automation.StartEthereumListener()

	log.Println("API Server: http://localhost:8080")

	log.Fatal(server.Run(":8080"))
}
