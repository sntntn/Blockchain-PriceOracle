package main

import (
	"fmt"
	"log"
	"math/big"
	"time"

	"Blockchain-PriceOracle/app/api"
	"Blockchain-PriceOracle/app/utils"
	"Blockchain-PriceOracle/internal/coingecko"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func printPrices(prices *map[string]*big.Int) {
	for symbol, price := range *prices {
		fmt.Printf("price -> %s: $%s\n", symbol, price.String())
	}
}

func CoinGeckoLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("CoinGecko loop: every 1min (background)")
	fmt.Println("----------------------------")
	fmt.Println("1min intervals started - CoinGecko")

	for {
		fmt.Println("\n=== NEW MINUTE ===")

		cgPrices, err := coingecko.FetchPrices()
		if err != nil {
			log.Printf("Fetch error: %v\n", err)
			// TO DO fallback prices
			<-ticker.C
			continue
		}

		printPrices(&cgPrices)

		for symbol, price := range cgPrices {
			if utils.CheckPriceCriteria(symbol, price) {
				log.Printf("%s - SEND TX NOW!\n", symbol)
			}
		}

		<-ticker.C
	}
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	utils.GetOracleClient()
	// utils.TestOracle()

	// -------------------------------------------

	go CoinGeckoLoop()

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/prices/:symbol", api.GetPricesHandler)
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "OK", "price_checker": "running"})
		})
	}

	log.Println("API Server: http://localhost:8080")

	log.Fatal(r.Run(":8080"))
}
