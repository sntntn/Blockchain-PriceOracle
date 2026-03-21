package main

import (
	"log"

	"Blockchain-PriceOracle/app/automation"
	"Blockchain-PriceOracle/app/server"
	"Blockchain-PriceOracle/app/utils"
	"Blockchain-PriceOracle/internal/oracle"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	oracle.GetOracleClient()
	utils.GetPriceHistory()

	go automation.CoinGeckoLoop()

	server := server.SetupServer()
	automation.StartEthereumListener()

	log.Println("API Server: http://localhost:8080")

	log.Fatal(server.Run(":8080"))
}
