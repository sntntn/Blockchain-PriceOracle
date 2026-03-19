package main

import (
	"fmt"
	"log"
	"time"

	"Blockchain-PriceOracle/internal/coingecko"
	"Blockchain-PriceOracle/internal/oracle"

	"github.com/joho/godotenv"
)

func printPrices(prices *map[string]float64) {
	for symbol, price := range *prices {
		fmt.Printf("price -> %s: $%.2f\n", symbol, price)
	}
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Oracle Test - Smart Contract PRICE")

	client, err := oracle.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	// Test BTC
	btcPrice, err := client.GetPrice("BTC")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("BTC Onchain: %s\n", btcPrice.String())

	btcCL, err := client.GetChainlinkPrice("BTC")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("BTC Chainlink: %s\n", btcCL.String())

	// -------------------------------------------
	fmt.Println("----------------------------")
	fmt.Println("CoinGecko Test - PRICE FETCH")

	ticker := time.NewTicker(1 * time.Minute)
	for {
		fmt.Println("\n=== NEW MINUTE ===")
		if cgPrices, err := coingecko.FetchPrices(); err != nil {
			fmt.Printf("Fetch error: %v\n", err)
			// TO DO fallback prices
		} else {
			printPrices(&cgPrices)
		}
		<-ticker.C
	}

}
