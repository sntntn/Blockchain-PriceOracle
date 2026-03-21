package main

import (
	"log"
	"os"

	"Blockchain-PriceOracle/app/automation"
	"Blockchain-PriceOracle/app/server"
	"Blockchain-PriceOracle/app/websocket"
	"Blockchain-PriceOracle/internal/oracle"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	oracle.GetOracleClient()

	go automation.CoinGeckoLoop()

	server := server.SetupServer()
	wsURL := os.Getenv("WSS")
	contractAddr := os.Getenv("CONTRACT_ADDR")
	if contractAddr != "" && wsURL != "" {
		go func() {
			log.Printf("Starting Ethereum listener: %s", wsURL)
			websocket.StartEventListener(contractAddr, wsURL)
		}()
	} else {
		log.Println("Skipping Ethereum listener - missing ORACLE_CONTRACT_ADDR or ETHEREUM_WSS_URL")
	}

	log.Println("API Server: http://localhost:8080")

	log.Fatal(server.Run(":8080"))
}
