package main

import (
	"log"

	"Blockchain-PriceOracle/app/automation"
	"Blockchain-PriceOracle/app/history"
	"Blockchain-PriceOracle/app/server"
	"Blockchain-PriceOracle/app/websocket"
	"Blockchain-PriceOracle/internal/coingecko"
	"Blockchain-PriceOracle/internal/oracle"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	priceHistory := history.GetPriceHistory()
	revertHistory := oracle.GetRevertHistory()
	oracleLimiter, syncLimiter, cgLimiter := automation.InitLimiters()

	oracleClient, err := oracle.GetOracleClient(revertHistory, oracleLimiter)
	if err != nil {
		log.Fatalf("failed to init oracle client: %v", err)
	}
	cgClient, err := coingecko.GetCoinGeckoClient(cgLimiter)
	if err != nil {
		log.Fatalf("failed to init CoinGecko client: %v", err)
	}

	clientWebsocketsManager := websocket.NewClientManager()

	automation.Sync(oracleClient, priceHistory, syncLimiter)

	go automation.CoinGeckoLoop(oracleClient, cgClient)

	server := server.SetupServer(oracleClient, cgClient, priceHistory, revertHistory, clientWebsocketsManager)
	automation.StartEthereumListener(priceHistory, clientWebsocketsManager)

	log.Println("API Server: http://localhost:8080")

	log.Fatal(server.Run(":8080"))
}
