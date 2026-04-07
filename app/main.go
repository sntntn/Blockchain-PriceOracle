package main

import (
	"log"
	"time"

	"Blockchain-PriceOracle/app/automation"
	"Blockchain-PriceOracle/app/history"
	"Blockchain-PriceOracle/app/server"
	"Blockchain-PriceOracle/app/websocket"
	"Blockchain-PriceOracle/internal/coingecko"
	"Blockchain-PriceOracle/internal/oracle"
	"Blockchain-PriceOracle/internal/ratelimit"

	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	priceHistory := history.GetPriceHistory()
	revertHistory := oracle.GetRevertHistory()
	oracleClient := oracle.GetOracleClient(revertHistory)
	clientWebsocketsManager := websocket.GetClientManager()
	cgLimiter := ratelimit.NewLocalLimiter(
		rate.Every(time.Minute/coingecko.CoinGeckoRateLimitPerMinute),
		coingecko.CoinGeckoRateLimitBurst,
	)
	cgClient := coingecko.NewClient(cgLimiter)

	automation.Sync(oracleClient, priceHistory)

	go automation.CoinGeckoLoop(oracleClient, cgClient)

	server := server.SetupServer(oracleClient, cgClient, priceHistory, revertHistory, clientWebsocketsManager)
	automation.StartEthereumListener(priceHistory, clientWebsocketsManager)

	log.Println("API Server: http://localhost:8080")

	log.Fatal(server.Run(":8080"))
}
