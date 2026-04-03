package main

import (
	"log"

	"Blockchain-PriceOracle/app/automation"
	"Blockchain-PriceOracle/app/history"
	"Blockchain-PriceOracle/app/server"
	"Blockchain-PriceOracle/internal/oracle"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	oracle.GetOracleClient()
	history.GetPriceHistory()

	go automation.CoinGeckoLoop()

	server := server.SetupServer()
	automation.StartEthereumListener()

	log.Println("API Server: http://localhost:8080")

	log.Fatal(server.Run(":8080"))
}
